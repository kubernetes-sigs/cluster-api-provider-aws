/*
Copyright 2019 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	ec2Service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	elbService "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/network"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/securitygroup"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
)

func TestAWSClusterReconcilerIntegrationTests(t *testing.T) {
	var (
		reconciler AWSClusterReconciler
		mockCtrl   *gomock.Controller
		recorder   *record.FakeRecorder
		ctx        context.Context
	)

	setup := func(t *testing.T) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		recorder = record.NewFakeRecorder(10)
		reconciler = AWSClusterReconciler{
			Client:   testEnv.Client,
			Recorder: recorder,
		}
		ctx = context.TODO()
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	t.Run("Should successfully reconcile AWSCluster creation with unmanaged VPC", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedCreateVPCCalls(m)
			mockedCreateSGCalls(false, m)
			mockedCreateLBCalls(t, e)
			mockedDescribeInstanceCall(m)
		}
		expect(ec2Mock.EXPECT(), elbMock.EXPECT())

		setup(t)
		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		awsCluster := getAWSCluster("test", ns.Name)
		awsCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeClassic,
		}

		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(Equal(true))

		defer teardown()
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsCluster, controllerIdentity, ns)).To(Succeed())
		})

		cs, err := getClusterScope(awsCluster)
		g.Expect(err).To(BeNil())
		networkSvc := network.NewService(cs)
		networkSvc.EC2Client = ec2Mock
		reconciler.networkServiceFactory = func(clusterScope scope.ClusterScope) services.NetworkInterface {
			return networkSvc
		}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}
		testSecurityGroupRoles := []infrav1.SecurityGroupRole{
			infrav1.SecurityGroupBastion,
			infrav1.SecurityGroupAPIServerLB,
			infrav1.SecurityGroupLB,
			infrav1.SecurityGroupControlPlane,
			infrav1.SecurityGroupNode,
		}
		sgSvc := securitygroup.NewService(cs, testSecurityGroupRoles)
		sgSvc.EC2Client = ec2Mock

		reconciler.securityGroupFactory = func(clusterScope scope.ClusterScope) services.SecurityGroupInterface {
			return sgSvc
		}
		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock

		reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}
		cs.SetSubnets([]infrav1.SubnetSpec{
			{
				ID:               "subnet-2",
				AvailabilityZone: "us-east-1c",
				IsPublic:         true,
				CidrBlock:        "10.0.11.0/24",
			},
			{
				ID:               "subnet-1",
				AvailabilityZone: "us-east-1a",
				CidrBlock:        "10.0.10.0/24",
				IsPublic:         false,
			},
		})
		_, err = reconciler.reconcileNormal(cs)
		g.Expect(err).To(BeNil())
		g.Expect(cs.VPC().ID).To(Equal("vpc-exists"))
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.ClusterSecurityGroupsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.BastionHostReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.VpcReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.SubnetsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
		})
	})
	t.Run("Should successfully reconcile AWSCluster creation with unmanaged VPC and a network type load balancer", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbv2Mock := mocks.NewMockELBV2API(mockCtrl)

		setup(t)
		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		awsCluster := getAWSCluster("test", ns.Name)
		awsCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeNLB,
			// Overwrite here because otherwise it's longer than 32, and we'll get a hashed name.
			Name: aws.String("test-cluster-apiserver"),
		}

		expect := func(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBV2APIMockRecorder) {
			mockedCreateVPCCalls(m)
			mockedCreateSGCalls(true, m)
			mockedCreateLBV2Calls(t, e)
			mockedDescribeInstanceCall(m)
		}
		expect(ec2Mock.EXPECT(), elbv2Mock.EXPECT())

		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(Equal(true))

		defer teardown()
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsCluster, controllerIdentity, ns)).To(Succeed())
		})

		cs, err := getClusterScope(awsCluster)
		cs.Cluster.Namespace = ns.Name
		g.Expect(err).To(BeNil())
		networkSvc := network.NewService(cs)
		networkSvc.EC2Client = ec2Mock
		reconciler.networkServiceFactory = func(clusterScope scope.ClusterScope) services.NetworkInterface {
			return networkSvc
		}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}
		testSecurityGroupRoles := []infrav1.SecurityGroupRole{
			infrav1.SecurityGroupBastion,
			infrav1.SecurityGroupAPIServerLB,
			infrav1.SecurityGroupLB,
			infrav1.SecurityGroupControlPlane,
			infrav1.SecurityGroupNode,
		}
		sgSvc := securitygroup.NewService(cs, testSecurityGroupRoles)
		sgSvc.EC2Client = ec2Mock

		reconciler.securityGroupFactory = func(clusterScope scope.ClusterScope) services.SecurityGroupInterface {
			return sgSvc
		}
		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBV2Client = elbv2Mock

		reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}
		cs.SetSubnets([]infrav1.SubnetSpec{
			{
				ID:               "subnet-2",
				AvailabilityZone: "us-east-1c",
				IsPublic:         true,
				CidrBlock:        "10.0.11.0/24",
			},
			{
				ID:               "subnet-1",
				AvailabilityZone: "us-east-1a",
				CidrBlock:        "10.0.10.0/24",
				IsPublic:         false,
			},
		})
		_, err = reconciler.reconcileNormal(cs)
		g.Expect(err).To(BeNil())
		g.Expect(cs.VPC().ID).To(Equal("vpc-exists"))
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.ClusterSecurityGroupsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.BastionHostReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.VpcReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.SubnetsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
		})
	})

	t.Run("Should fail on AWSCluster reconciliation if VPC limit exceeded", func(t *testing.T) {
		// Assuming the max VPC limit is 2 and when two VPCs are created, the creation of 3rd VPC throws mocked error from EC2 API
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbv2Mock := mocks.NewMockELBV2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, ev2 *mocks.MockELBV2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedCreateMaximumVPCCalls(m)
			mockedDeleteVPCCallsForNonExistentVPC(m)
			mockedDeleteLBCalls(true, ev2, e)
			mockedDescribeInstanceCall(m)
			mockedDeleteInstanceAndAwaitTerminationCalls(m)
		}
		expect(ec2Mock.EXPECT(), elbv2Mock.EXPECT(), elbMock.EXPECT())

		setup(t)
		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())
		awsCluster := infrav1.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: ns.Name,
			},
			Spec: infrav1.AWSClusterSpec{
				Region: "us-east-1",
			},
		}
		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())

		defer teardown()
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(Equal(true))
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsCluster, controllerIdentity, ns)).To(Succeed())
		})
		cs, err := getClusterScope(awsCluster)
		g.Expect(err).To(BeNil())

		networkSvc := network.NewService(cs)
		networkSvc.EC2Client = ec2Mock
		reconciler.networkServiceFactory = func(clusterScope scope.ClusterScope) services.NetworkInterface {
			return networkSvc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		elbSvc.ELBV2Client = elbv2Mock
		reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(ec2Scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		_, err = reconciler.reconcileNormal(cs)
		g.Expect(err.Error()).To(ContainSubstring("The maximum number of VPCs has been reached"))

		_, err = reconciler.reconcileDelete(ctx, cs)
		g.Expect(err).To(BeNil())
	})
	t.Run("Should successfully delete AWSCluster with managed VPC", func(t *testing.T) {
		g := NewWithT(t)

		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		elbv2Mock := mocks.NewMockELBV2API(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, ev2 *mocks.MockELBV2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedDeleteVPCCalls(m)
			mockedDescribeInstanceCall(m)
			mockedDeleteLBCalls(true, ev2, e)
			mockedDeleteInstanceAndAwaitTerminationCalls(m)
			mockedDeleteSGCalls(m)
		}
		expect(ec2Mock.EXPECT(), elbv2Mock.EXPECT(), elbMock.EXPECT())

		setup(t)
		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())
		awsCluster := getAWSCluster("test", ns.Name)

		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())
		defer teardown()
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(Equal(true))

		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsCluster, controllerIdentity, ns)).To(Succeed())
		})

		cs, err := getClusterScope(awsCluster)
		g.Expect(err).To(BeNil())

		networkSvc := network.NewService(cs)
		networkSvc.EC2Client = ec2Mock
		reconciler.networkServiceFactory = func(clusterScope scope.ClusterScope) services.NetworkInterface {
			return networkSvc
		}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(ec2Scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		elbSvc.ELBV2Client = elbv2Mock
		reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		testSecurityGroupRoles := []infrav1.SecurityGroupRole{
			infrav1.SecurityGroupBastion,
			infrav1.SecurityGroupAPIServerLB,
			infrav1.SecurityGroupLB,
			infrav1.SecurityGroupControlPlane,
			infrav1.SecurityGroupNode,
		}
		sgSvc := securitygroup.NewService(cs, testSecurityGroupRoles)
		sgSvc.EC2Client = ec2Mock
		reconciler.securityGroupFactory = func(clusterScope scope.ClusterScope) services.SecurityGroupInterface {
			return sgSvc
		}

		_, err = reconciler.reconcileDelete(ctx, cs)
		g.Expect(err).To(BeNil())
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.LoadBalancerReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.BastionHostReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.SecondaryCidrsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletingReason},
			{infrav1.RouteTablesReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.NatGatewaysReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.InternetGatewayReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.SubnetsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.VpcReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
		})
	})
}

func mockedDeleteSGCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSecurityGroupsPages(gomock.Any(), gomock.Any()).Return(nil)
}

func createControllerIdentity(g *WithT) *infrav1.AWSClusterControllerIdentity {
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			Kind: string(infrav1.ControllerIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}
	g.Expect(testEnv.Create(ctx, controllerIdentity)).To(Succeed())
	return controllerIdentity
}

func mockedDescribeInstanceCall(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: aws.StringSlice([]string{"bastion"}),
			},
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: aws.StringSlice([]string{"pending", "running", "stopping", "stopped"}),
			},
		},
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("id-1"),
						InstanceType: aws.String("m5.large"),
						SubnetId:     aws.String("subnet-1"),
						ImageId:      aws.String("ami-1"),
						IamInstanceProfile: &ec2.IamInstanceProfile{
							Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
						},
						State: &ec2.InstanceState{
							Code: aws.Int64(16),
							Name: aws.String(ec2.StateAvailable),
						},
						RootDeviceName: aws.String("device-1"),
						BlockDeviceMappings: []*ec2.InstanceBlockDeviceMapping{
							{
								DeviceName: aws.String("device-1"),
								Ebs: &ec2.EbsInstanceBlockDevice{
									VolumeId: aws.String("volume-1"),
								},
							},
						},
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("us-east-1a"),
						},
					},
				},
			},
		},
	}, nil)
}

func mockedDeleteInstanceAndAwaitTerminationCalls(m *mocks.MockEC2APIMockRecorder) {
	m.TerminateInstances(
		gomock.Eq(&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice([]string{"id-1"}),
		}),
	).Return(nil, nil)
	m.WaitUntilInstanceTerminated(
		gomock.Eq(&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{"id-1"}),
		}),
	).Return(nil)
}

func mockedDeleteInstanceCalls(m *mocks.MockEC2APIMockRecorder) {
	m.TerminateInstances(
		gomock.Eq(&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice([]string{"id-1"}),
		}),
	).Return(nil, nil)
}

func mockedCreateVPCCalls(m *mocks.MockEC2APIMockRecorder) {
	m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
		Resources: aws.StringSlice([]string{"subnet-1"}),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("kubernetes.io/cluster/test-cluster"),
				Value: aws.String("shared"),
			},
			{
				Key:   aws.String("kubernetes.io/role/internal-elb"),
				Value: aws.String("1"),
			},
		},
	})).Return(&ec2.CreateTagsOutput{}, nil)
	m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
		Resources: aws.StringSlice([]string{"subnet-2"}),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("kubernetes.io/cluster/test-cluster"),
				Value: aws.String("shared"),
			},
			{
				Key:   aws.String("kubernetes.io/role/elb"),
				Value: aws.String("1"),
			},
		},
	})).Return(&ec2.CreateTagsOutput{}, nil)
	m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
		}})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{
				VpcId:               aws.String("vpc-exists"),
				SubnetId:            aws.String("subnet-1"),
				AvailabilityZone:    aws.String("us-east-1a"),
				CidrBlock:           aws.String("10.0.10.0/24"),
				MapPublicIpOnLaunch: aws.Bool(false),
			},
			{
				VpcId:               aws.String("vpc-exists"),
				SubnetId:            aws.String("subnet-2"),
				AvailabilityZone:    aws.String("us-east-1c"),
				CidrBlock:           aws.String("10.0.11.0/24"),
				MapPublicIpOnLaunch: aws.Bool(false),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("public"),
					},
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-subnet-public"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("shared"),
					},
				},
			},
		},
	}, nil)
	m.DescribeRouteTables(gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
		}})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []*ec2.RouteTable{
			{
				Routes: []*ec2.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
			},
		},
	}, nil)
	m.DescribeNatGatewaysPages(gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String("vpc-exists")},
			},
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
		}}), gomock.Any()).Return(nil)
	m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String("vpc-exists"),
		},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
		},
	})).
		Return(&ec2.DescribeVpcsOutput{
			Vpcs: []*ec2.Vpc{
				{
					State:     aws.String("available"),
					VpcId:     aws.String("vpc-exists"),
					CidrBlock: aws.String("10.0.0.0/8"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("common"),
						},
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster"),
						},
					},
				},
			},
		}, nil)
}

func mockedCreateMaximumVPCCalls(m *mocks.MockEC2APIMockRecorder) {
	m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).Return(nil, errors.New("The maximum number of VPCs has been reached"))
}

func mockedDeleteVPCCallsForNonExistentVPC(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			},
		}})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{},
	}, nil).AnyTimes()
	m.DescribeRouteTables(gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{""}),
		},
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			},
		}})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []*ec2.RouteTable{}}, nil).AnyTimes()
	m.DescribeInternetGateways(gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: aws.StringSlice([]string{""}),
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []*ec2.InternetGateway{},
	}, nil)
	m.DescribeNatGatewaysPages(gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String("")},
			},
		},
	}), gomock.Any()).Return(nil).AnyTimes()
	m.DescribeAddresses(gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			}},
	})).Return(nil, nil)
	m.DeleteVpc(gomock.AssignableToTypeOf(&ec2.DeleteVpcInput{
		VpcId: aws.String("vpc-exists")})).Return(nil, nil)
}

func mockedDeleteVPCCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
		}})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{
				VpcId:               aws.String("vpc-exists"),
				SubnetId:            aws.String("subnet-1"),
				AvailabilityZone:    aws.String("us-east-1a"),
				CidrBlock:           aws.String("10.0.10.0/24"),
				MapPublicIpOnLaunch: aws.Bool(false),
			},
		},
	}, nil).AnyTimes()
	m.DescribeRouteTables(gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			},
		}})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []*ec2.RouteTable{
			{
				Routes: []*ec2.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
				RouteTableId: aws.String("rt-12345"),
			},
		},
	}, nil).AnyTimes()
	m.DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String("rt-12345"),
	}))
	m.DescribeInternetGateways(gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []*ec2.InternetGateway{
			{
				Attachments:       nil,
				InternetGatewayId: aws.String("ig-12345"),
			},
		},
	}, nil)
	m.DetachInternetGateway(gomock.Eq(&ec2.DetachInternetGatewayInput{
		VpcId:             aws.String("vpc-exists"),
		InternetGatewayId: aws.String("ig-12345"),
	}))
	m.DeleteInternetGateway(gomock.Eq(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String("ig-12345"),
	}))
	m.DescribeNatGatewaysPages(gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String("vpc-exists")},
			},
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
		}}), gomock.Any()).Return(nil).AnyTimes()
	m.DescribeAddresses(gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			}},
	})).Return(&ec2.DescribeAddressesOutput{
		Addresses: []*ec2.Address{
			{
				AssociationId: aws.String("1234"),
				AllocationId:  aws.String("1234"),
				PublicIp:      aws.String("1.2.3.4"),
			},
		},
	}, nil)
	m.DisassociateAddress(&ec2.DisassociateAddressInput{
		AssociationId: aws.String("1234"),
	})
	m.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: aws.String("1234"),
	})
	m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String("vpc-exists"),
		},
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
			},
		},
	})).
		Return(&ec2.DescribeVpcsOutput{
			Vpcs: []*ec2.Vpc{
				{
					State:     aws.String("available"),
					VpcId:     aws.String("vpc-exists"),
					CidrBlock: aws.String("10.0.0.0/8"),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							Value: aws.String("common"),
						},
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster"),
						},
						{
							Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
							Value: aws.String("owned"),
						},
					},
				},
			},
		}, nil)
	m.DeleteSubnet(gomock.Eq(&ec2.DeleteSubnetInput{
		SubnetId: aws.String("subnet-1"),
	}))
	m.DeleteVpc(gomock.Eq(&ec2.DeleteVpcInput{
		VpcId: aws.String("vpc-exists"),
	}))
}

func mockedCreateSGCalls(recordLBV2 bool, m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSecurityGroups(gomock.Eq(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{"vpc-exists"}),
			},
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
			},
		},
	})).Return(
		&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []*ec2.SecurityGroup{
				{
					GroupId:   aws.String("1"),
					GroupName: aws.String("test-sg"),
				},
			},
		}, nil)
	m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-exists"),
		GroupName:   aws.String("test-cluster-bastion"),
		Description: aws.String("Kubernetes cluster test-cluster: bastion"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-bastion"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("bastion"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-bastion")}, nil)
	m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-exists"),
		GroupName:   aws.String("test-cluster-apiserver-lb"),
		Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-apiserver-lb"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("apiserver-lb"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-apiserver-lb")}, nil)
	m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-exists"),
		GroupName:   aws.String("test-cluster-lb"),
		Description: aws.String("Kubernetes cluster test-cluster: lb"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-lb"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("lb"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-lb")}, nil)
	securityGroupControl := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-exists"),
		GroupName:   aws.String("test-cluster-controlplane"),
		Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-controlplane"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("controlplane"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-controlplane")}, nil)
	securityGroupNode := m.CreateSecurityGroup(gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String("vpc-exists"),
		GroupName:   aws.String("test-cluster-node"),
		Description: aws.String("Kubernetes cluster test-cluster: node"),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("security-group"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-node"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("node"),
					},
				},
			},
		},
	})).
		Return(&ec2.CreateSecurityGroupOutput{GroupId: aws.String("sg-node")}, nil)
	m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-controlplane"),
	})).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
		After(securityGroupControl).Times(2)
	m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-node"),
	})).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
		After(securityGroupNode).Times(2)
	if recordLBV2 {
		m.AuthorizeSecurityGroupIngress(gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String("sg-lb"),
		})).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
			After(securityGroupNode).Times(1)
	}
}
