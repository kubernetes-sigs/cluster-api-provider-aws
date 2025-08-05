/*
Copyright The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	. "github.com/onsi/gomega"
	rosaAWSClient "github.com/openshift/rosa/pkg/aws"
	rosaMocks "github.com/openshift/rosa/pkg/aws/mocks"
	"github.com/sirupsen/logrus"
	gomock "go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func TestROSANetworkReconciler_Reconcile(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, "test-namespace")
	g.Expect(err).ToNot(HaveOccurred())

	mockCtrl := gomock.NewController(t)
	ctx := context.TODO()

	identity := &infrav1.AWSClusterControllerIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}

	name := "test-rosa-network"
	rosaNetwork := &expinfrav1.ROSANetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns.Name},
		Spec: expinfrav1.ROSANetworkSpec{
			StackName:             name,
			CIDRBlock:             "10.0.0.0/8",
			AvailabilityZoneCount: 1,
			Region:                "test-region",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: identity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}

	createObject(g, identity, ns.Name)
	createObject(g, rosaNetwork, ns.Name)

	nameDeleted := "test-rosa-network-deleted"
	rosaNetworkDeleted := &expinfrav1.ROSANetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nameDeleted,
			Namespace: ns.Name},
		Spec: expinfrav1.ROSANetworkSpec{
			StackName:             nameDeleted,
			CIDRBlock:             "10.0.0.0/8",
			AvailabilityZoneCount: 1,
			Region:                "test-region",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: identity.Name,
				Kind: infrav1.ControllerIdentityKind,
			},
		},
	}
	controllerutil.AddFinalizer(rosaNetworkDeleted, expinfrav1.ROSANetworkFinalizer)
	createObject(g, rosaNetworkDeleted, ns.Name)
	err = deleteROSANetwork(ctx, rosaNetworkDeleted)
	g.Expect(err).NotTo(HaveOccurred())

	t.Run("Empty result when ROSANetwork object not found", func(t *testing.T) {
		_, _, _, reconciler := createMockClients(mockCtrl)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: "non-existent-object", Namespace: "non-existent-namespace"}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).ToNot(HaveOccurred())
	})

	t.Run("Error result when CF stack GET returns error", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)
		mockDescribeStacksCall(mockCFClient, &cloudformation.DescribeStacksOutput{}, fmt.Errorf("test-error"), 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).To(MatchError(ContainSubstring("error fetching CF stack details:")))
	})

	t.Run("Initial CF stack creation fails", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{}
		validationErr := &smithy.GenericAPIError{
			Code:    "ValidationError",
			Message: "ValidationError",
			Fault:   smithy.FaultServer,
		}

		mockDescribeStacksCall(mockCFClient, describeStacksOutput, validationErr, 1)
		mockCreateStackCall(mockCFClient, &cloudformation.CreateStackOutput{}, fmt.Errorf("test-error"), 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).To(MatchError(ContainSubstring("failed to start CF stack creation:")))

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetwork)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkFailedReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityError))
		g.Expect(cnd.Message).To(Equal("test-error"))
	})

	t.Run("Initial CF stack creation succeeds", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{}
		validationErr := &smithy.GenericAPIError{
			Code:    "ValidationError",
			Message: "ValidationError",
			Fault:   smithy.FaultServer,
		}

		mockDescribeStacksCall(mockCFClient, describeStacksOutput, validationErr, 1)
		mockCreateStackCall(mockCFClient, &cloudformation.CreateStackOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).ToNot(HaveOccurred())

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetwork)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkCreatingReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityInfo))
	})

	t.Run("CF stack creation is in progress", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &name,
					StackStatus: cloudformationtypes.StackStatusCreateInProgress,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Second * 60))
		g.Expect(errReconcile).ToNot(HaveOccurred())

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetwork)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkCreatingReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityInfo))
	})

	t.Run("CF stack creation completed", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &name,
					StackStatus: cloudformationtypes.StackStatusCreateComplete,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).ToNot(HaveOccurred())

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetwork)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkCreatedReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityInfo))
	})

	t.Run("CF stack creation failed", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &name,
					StackStatus: cloudformationtypes.StackStatusCreateFailed,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: rosaNetwork.Name, Namespace: rosaNetwork.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).To(MatchError(ContainSubstring("creation failed")))

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetwork)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkFailedReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityError))
	})

	t.Run("CF stack deletion start failed", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &nameDeleted,
					StackStatus: cloudformationtypes.StackStatusCreateComplete,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		mockDeleteStackCall(mockCFClient, &cloudformation.DeleteStackOutput{}, fmt.Errorf("test-error"), 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: nameDeleted, Namespace: rosaNetworkDeleted.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).To(MatchError(ContainSubstring("failed to start CF stack deletion:")))

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetworkDeleted)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkDeletionFailedReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityError))
	})

	t.Run("CF stack deletion start succeeded", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &nameDeleted,
					StackStatus: cloudformationtypes.StackStatusCreateComplete,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		mockDeleteStackCall(mockCFClient, &cloudformation.DeleteStackOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: nameDeleted, Namespace: rosaNetworkDeleted.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(60 * time.Second))
		g.Expect(errReconcile).NotTo(HaveOccurred())

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetworkDeleted)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkDeletingReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityInfo))
	})

	t.Run("CF stack deletion in progress", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &nameDeleted,
					StackStatus: cloudformationtypes.StackStatusDeleteInProgress,
				},
			},
		}
		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: nameDeleted, Namespace: rosaNetworkDeleted.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(60 * time.Second))
		g.Expect(errReconcile).NotTo(HaveOccurred())
	})

	t.Run("CF stack deletion failed", func(t *testing.T) {
		_, mockCFClient, mockSTSClient, reconciler := createMockClients(mockCtrl)

		mockSTSIdentity(mockSTSClient)

		describeStacksOutput := &cloudformation.DescribeStacksOutput{
			Stacks: []cloudformationtypes.Stack{
				{
					StackName:   &nameDeleted,
					StackStatus: cloudformationtypes.StackStatusDeleteFailed,
				},
			},
		}

		mockDescribeStacksCall(mockCFClient, describeStacksOutput, nil, 1)

		describeStackResourcesOutput := &cloudformation.DescribeStackResourcesOutput{
			StackResources: []cloudformationtypes.StackResource{},
		}

		mockDescribeStackResourcesCall(mockCFClient, describeStackResourcesOutput, nil, 1)

		req := ctrl.Request{}
		req.NamespacedName = types.NamespacedName{Name: nameDeleted, Namespace: rosaNetworkDeleted.Namespace}
		reqReconcile, errReconcile := reconciler.Reconcile(ctx, req)

		g.Expect(reqReconcile.Requeue).To(BeFalse())
		g.Expect(reqReconcile.RequeueAfter).To(Equal(time.Duration(0)))
		g.Expect(errReconcile).To(MatchError(ContainSubstring("CF stack deletion failed")))

		cnd, err := getROSANetworkReadyCondition(reconciler, rosaNetworkDeleted)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(cnd).ToNot(BeNil())
		g.Expect(cnd.Reason).To(Equal(expinfrav1.ROSANetworkDeletionFailedReason))
		g.Expect(cnd.Severity).To(Equal(clusterv1.ConditionSeverityError))
	})

	cleanupObject(g, rosaNetwork)
	cleanupObject(g, rosaNetworkDeleted)
	cleanupObject(g, identity)
}

func TestROSANetworkReconciler_updateROSANetworkResources(t *testing.T) {
	g := NewWithT(t)
	mockCtrl := gomock.NewController(t)
	ctx := context.TODO()

	rosaNetwork := &expinfrav1.ROSANetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-rosa-network",
			Namespace: "test-namespace",
		},
		Spec:   expinfrav1.ROSANetworkSpec{},
		Status: expinfrav1.ROSANetworkStatus{},
	}

	t.Run("Handle cloudformation client error", func(t *testing.T) {
		_, mockCFClient, _, reconciler := createMockClients(mockCtrl)

		mockDescribeStackResourcesCall(mockCFClient, &cloudformation.DescribeStackResourcesOutput{}, fmt.Errorf("test-error"), 1)

		err := reconciler.updateROSANetworkResources(ctx, rosaNetwork)
		g.Expect(err).To(HaveOccurred())
		g.Expect(len(rosaNetwork.Status.Resources)).To(Equal(0))
	})

	t.Run("Update ROSANetwork.Status.Resources", func(t *testing.T) {
		_, mockCFClient, _, reconciler := createMockClients(mockCtrl)

		logicalResourceID := "logical-resource-id"
		resourceStatus := cloudformationtypes.ResourceStatusCreateComplete
		resourceType := "resource-type"
		resourceStatusReason := "resource-status-reason"
		physicalResourceID := "physical-resource-id"

		describeStackResourcesOutput := &cloudformation.DescribeStackResourcesOutput{
			StackResources: []cloudformationtypes.StackResource{
				{
					LogicalResourceId:    &logicalResourceID,
					ResourceStatus:       resourceStatus,
					ResourceType:         &resourceType,
					ResourceStatusReason: &resourceStatusReason,
					PhysicalResourceId:   &physicalResourceID,
				},
			},
		}

		mockDescribeStackResourcesCall(mockCFClient, describeStackResourcesOutput, nil, 1)

		err := reconciler.updateROSANetworkResources(ctx, rosaNetwork)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rosaNetwork.Status.Resources[0].LogicalID).To(Equal(logicalResourceID))
		g.Expect(rosaNetwork.Status.Resources[0].Status).To(Equal(string(resourceStatus)))
		g.Expect(rosaNetwork.Status.Resources[0].ResourceType).To(Equal(resourceType))
		g.Expect(rosaNetwork.Status.Resources[0].Reason).To(Equal(resourceStatusReason))
		g.Expect(rosaNetwork.Status.Resources[0].PhysicalID).To(Equal(physicalResourceID))
	})
}

func TestROSANetworkReconciler_parseSubnets(t *testing.T) {
	g := NewWithT(t)
	mockCtrl := gomock.NewController(t)

	subnet1Id := "subnet1-physical-id"
	subnet2Id := "subnet2-physical-id"

	rosaNetwork := &expinfrav1.ROSANetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-rosa-network",
			Namespace: "test-namespace",
		},
		Spec: expinfrav1.ROSANetworkSpec{},
		Status: expinfrav1.ROSANetworkStatus{
			Resources: []expinfrav1.CFResource{
				{
					ResourceType: "AWS::EC2::Subnet",
					LogicalID:    "SubnetPrivate",
					PhysicalID:   subnet1Id,
					Status:       "subnet1-status",
					Reason:       "subnet1-reason",
				},
				{
					ResourceType: "AWS::EC2::Subnet",
					LogicalID:    "SubnetPublic",
					PhysicalID:   subnet2Id,
					Status:       "subnet2-status",
					Reason:       "subnet2-reason",
				},
				{
					ResourceType: "bogus-type",
					LogicalID:    "bogus-logical-id",
					PhysicalID:   "bugus-physical-id",
					Status:       "bogus-status",
					Reason:       "bogus-reason",
				},
			},
		},
	}

	t.Run("Handle EC2 client error", func(t *testing.T) {
		mockEC2Client, _, _, reconciler := createMockClients(mockCtrl)

		mockDescribeSubnetsCall(mockEC2Client, &ec2.DescribeSubnetsOutput{}, nil, 1)

		err := reconciler.parseSubnets(rosaNetwork)
		g.Expect(err).To(HaveOccurred())
		g.Expect(len(rosaNetwork.Status.Subnets)).To(Equal(0))
	})

	t.Run("Update ROSANetwork.Status.Subnets", func(t *testing.T) {
		mockEC2Client, _, _, reconciler := createMockClients(mockCtrl)

		az := "az01"

		describeSubnetsOutput := &ec2.DescribeSubnetsOutput{
			Subnets: []ec2Types.Subnet{
				{
					AvailabilityZone: &az,
				},
			},
		}

		mockDescribeSubnetsCall(mockEC2Client, describeSubnetsOutput, nil, 2)

		err := reconciler.parseSubnets(rosaNetwork)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rosaNetwork.Status.Subnets[0].AvailabilityZone).To(Equal(az))
		g.Expect(rosaNetwork.Status.Subnets[0].PrivateSubnet).To(Equal(subnet1Id))
		g.Expect(rosaNetwork.Status.Subnets[0].PublicSubnet).To(Equal(subnet2Id))
	})
}

func createMockClients(mockCtrl *gomock.Controller) (*rosaMocks.MockEc2ApiClient, *rosaMocks.MockCloudFormationApiClient, *rosaMocks.MockStsApiClient, *ROSANetworkReconciler) {
	mockEC2Client := rosaMocks.NewMockEc2ApiClient(mockCtrl)
	mockCFClient := rosaMocks.NewMockCloudFormationApiClient(mockCtrl)
	mockSTSClient := rosaMocks.NewMockStsApiClient(mockCtrl)
	awsClient := rosaAWSClient.New(
		awsSdk.Config{},
		rosaAWSClient.NewLoggerWrapper(logrus.New(), nil),
		rosaMocks.NewMockIamApiClient(mockCtrl),
		mockEC2Client,
		rosaMocks.NewMockOrganizationsApiClient(mockCtrl),
		rosaMocks.NewMockS3ApiClient(mockCtrl),
		rosaMocks.NewMockSecretsManagerApiClient(mockCtrl),
		mockSTSClient,
		mockCFClient,
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		rosaMocks.NewMockServiceQuotasApiClient(mockCtrl),
		&rosaAWSClient.AccessKey{},
		false,
	)

	reconciler := &ROSANetworkReconciler{
		Client:    testEnv.Client,
		awsClient: awsClient,
	}

	return mockEC2Client, mockCFClient, mockSTSClient, reconciler
}

func mockSTSIdentity(mockSTSClient *rosaMocks.MockStsApiClient) {
	getCallerIdentityResult := &stsv2.GetCallerIdentityOutput{
		Account: awsSdk.String("foo"),
		Arn:     awsSdk.String("arn:aws:iam::123456789012:rosa/foo"),
	}
	mockSTSClient.
		EXPECT().
		GetCallerIdentity(gomock.Any(), gomock.Any()).
		Return(getCallerIdentityResult, nil).
		AnyTimes()
}

func mockDescribeStacksCall(mockCFClient *rosaMocks.MockCloudFormationApiClient, output *cloudformation.DescribeStacksOutput, err error, times int) {
	mockCFClient.
		EXPECT().
		DescribeStacks(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context,
			_ *cloudformation.DescribeStacksInput,
			_ ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
			return output, err
		}).
		Times(times)
}

func mockCreateStackCall(mockCFClient *rosaMocks.MockCloudFormationApiClient, output *cloudformation.CreateStackOutput, err error, times int) {
	mockCFClient.
		EXPECT().
		CreateStack(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context,
			_ *cloudformation.CreateStackInput,
			_ ...func(*cloudformation.Options)) (*cloudformation.CreateStackOutput, error) {
			return output, err
		}).
		Times(times)
}

func mockDescribeStackResourcesCall(mockCFClient *rosaMocks.MockCloudFormationApiClient, output *cloudformation.DescribeStackResourcesOutput, err error, times int) {
	mockCFClient.
		EXPECT().
		DescribeStackResources(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context,
			_ *cloudformation.DescribeStackResourcesInput,
			_ ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourcesOutput, error) {
			return output, err
		}).
		Times(times)
}

func mockDeleteStackCall(mockCFClient *rosaMocks.MockCloudFormationApiClient, output *cloudformation.DeleteStackOutput, err error, times int) {
	mockCFClient.
		EXPECT().
		DeleteStack(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context,
			_ *cloudformation.DeleteStackInput,
			_ ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error) {
			return output, err
		}).
		Times(times)
}

func mockDescribeSubnetsCall(mockEc2Client *rosaMocks.MockEc2ApiClient, output *ec2.DescribeSubnetsOutput, err error, times int) {
	mockEc2Client.
		EXPECT().
		DescribeSubnets(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context,
			_ *ec2.DescribeSubnetsInput,
			_ ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
			return output, err
		}).
		Times(times)
}

func deleteROSANetwork(ctx context.Context, rosaNetwork *expinfrav1.ROSANetwork) error {
	if err := testEnv.Client.Get(ctx, client.ObjectKeyFromObject(rosaNetwork), rosaNetwork); err != nil {
		return err
	}

	if !rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() {
		return nil
	}

	if err := testEnv.Client.Delete(ctx, rosaNetwork); err != nil {
		return err
	}

	for {
		if err := testEnv.Client.Get(ctx, client.ObjectKeyFromObject(rosaNetwork), rosaNetwork); err != nil {
			return err
		}

		if !rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

func getROSANetworkReadyCondition(reconciler *ROSANetworkReconciler, rosaNet *expinfrav1.ROSANetwork) (*clusterv1.Condition, error) {
	updatedROSANetwork := &expinfrav1.ROSANetwork{}

	if err := reconciler.Client.Get(ctx, client.ObjectKeyFromObject(rosaNet), updatedROSANetwork); err != nil {
		return nil, err
	}

	return conditions.Get(updatedROSANetwork, expinfrav1.ROSANetworkReadyCondition), nil
}
