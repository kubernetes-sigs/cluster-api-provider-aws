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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
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
	t.Run("Should wait for external Control Plane endpoint when LoadBalancer is disabled, and eventually succeed when patched", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder) {
			// First iteration, when the AWS Cluster is missing a valid Control Plane Endpoint
			mockedVPCCallsForExistingVPCAndSubnets(m)
			mockedCreateSGCalls(false, "vpc-exists", m)
			mockedDescribeInstanceCall(m)
			mockedDescribeAvailabilityZones(m, []string{"us-east-1c", "us-east-1a"})

			// Second iteration: the AWS Cluster object has been patched,
			// thus a valid Control Plane Endpoint has been provided
			mockedVPCCallsForExistingVPCAndSubnets(m)
			mockedCreateSGCalls(false, "vpc-exists", m)
			mockedDescribeInstanceCall(m)
		}
		expect(ec2Mock.EXPECT())

		setup(t)
		controllerIdentity := createControllerIdentity(g)
		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())
		// Creating the AWS cluster with a disabled Load Balancer:
		// no ALB, ELB, or NLB specified, the AWS cluster must consistently be reported
		// waiting for the control Plane endpoint.
		awsCluster := getAWSCluster("test", ns.Name)
		awsCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
		}

		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())

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

		_, err = reconciler.reconcileNormal(context.TODO(), cs)
		g.Expect(err).To(BeNil())

		cluster := &infrav1.AWSCluster{}
		g.Expect(testEnv.Get(ctx, client.ObjectKey{Name: cs.AWSCluster.Name, Namespace: cs.AWSCluster.Namespace}, cluster)).ToNot(HaveOccurred())
		g.Expect(cluster.Spec.ControlPlaneEndpoint.Host).To(BeEmpty())
		g.Expect(cluster.Spec.ControlPlaneEndpoint.Port).To(BeZero())
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.LoadBalancerReadyCondition, status: corev1.ConditionFalse, severity: clusterv1.ConditionSeverityInfo, reason: infrav1.WaitForExternalControlPlaneEndpointReason},
		})
		// Mimicking an external operator patching the cluster with an already provisioned Load Balancer:
		// this could be done by a human who provisioned a LB, or by a Control Plane provider.
		g.Expect(retry.RetryOnConflict(retry.DefaultRetry, func() error {
			if err = testEnv.Get(ctx, client.ObjectKey{Name: cs.AWSCluster.Name, Namespace: cs.AWSCluster.Namespace}, cs.AWSCluster); err != nil {
				return err
			}

			cs.AWSCluster.Spec.ControlPlaneEndpoint.Host = "10.0.10.1"
			cs.AWSCluster.Spec.ControlPlaneEndpoint.Port = 6443

			return testEnv.Update(ctx, cs.AWSCluster)
		})).To(Succeed())
		// Executing back a second reconciliation:
		// the AWS Cluster should be ready with no LoadBalancer false condition.
		_, err = reconciler.reconcileNormal(context.TODO(), cs)
		g.Expect(err).To(BeNil())
		g.Expect(cs.VPC().ID).To(Equal("vpc-exists"))
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.ClusterSecurityGroupsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.BastionHostReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.VpcReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.SubnetsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.LoadBalancerReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
		})
	})
	t.Run("Should successfully reconcile AWSCluster creation with unmanaged VPC", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedVPCCallsForExistingVPCAndSubnets(m)
			mockedCreateSGCalls(false, "vpc-exists", m)
			mockedCreateLBCalls(t, e, true)
			mockedDescribeInstanceCall(m)
			mockedDescribeAvailabilityZones(m, []string{"us-east-1c", "us-east-1a"})
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
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", awsCluster.Name))

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
		_, err = reconciler.reconcileNormal(context.TODO(), cs)
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
			mockedVPCCallsForExistingVPCAndSubnets(m)
			mockedCreateSGCalls(true, "vpc-exists", m)
			mockedCreateLBV2Calls(t, e)
			mockedDescribeInstanceCall(m)
			mockedDescribeAvailabilityZones(m, []string{"us-east-1c", "us-east-1a"})
			mockedDescribeTargetGroupsCall(t, e)
			mockedCreateTargetGroupCall(t, e)
			mockedModifyTargetGroupAttributes(t, e)
			mockedDescribeListenersCall(t, e)
			mockedCreateListenerCall(t, e)
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
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", awsCluster.Name))

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
		_, err = reconciler.reconcileNormal(context.TODO(), cs)
		g.Expect(err).To(BeNil())
		g.Expect(cs.VPC().ID).To(Equal("vpc-exists"))
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.ClusterSecurityGroupsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.BastionHostReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.VpcReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.SubnetsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
		})
	})
	t.Run("Should successfully reconcile AWSCluster creation with managed VPC", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		expect := func(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedCallsForMissingEverything(m, e, "my-managed-subnet-priv", "my-managed-subnet-pub")
			mockedCreateSGCalls(false, "vpc-new", m)
			mockedDescribeInstanceCall(m)
			mockedDescribeAvailabilityZones(m, []string{"us-east-1a"})
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

		// Make controller manage resources
		awsCluster.Spec.NetworkSpec.VPC.ID = ""
		awsCluster.Spec.NetworkSpec.Subnets[0].ID = "my-managed-subnet-priv"
		awsCluster.Spec.NetworkSpec.Subnets[1].ID = "my-managed-subnet-pub"

		// NAT gateway of the public subnet will be accessed by the private subnet in the same zone,
		// so use same zone for the 2 test subnets
		awsCluster.Spec.NetworkSpec.Subnets[0].AvailabilityZone = "us-east-1a"
		awsCluster.Spec.NetworkSpec.Subnets[1].AvailabilityZone = "us-east-1a"

		g.Expect(testEnv.Create(ctx, &awsCluster)).To(Succeed())
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: ns.Name,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", awsCluster.Name))

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
		_, err = reconciler.reconcileNormal(context.TODO(), cs)
		g.Expect(err).To(BeNil())
		g.Expect(cs.VPC().ID).To(Equal("vpc-new"))
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{conditionType: infrav1.ClusterSecurityGroupsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.BastionHostReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.VpcReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
			{conditionType: infrav1.SubnetsReadyCondition, status: corev1.ConditionTrue, severity: "", reason: ""},
		})

		// Information should get written back into the `ClusterScope` object. Keeping it up to date means that
		// reconciliation functionality will always work on the latest-known status of AWS cloud resources.

		// Private subnet
		g.Expect(cs.Subnets()[0].ID).To(Equal("my-managed-subnet-priv"))
		g.Expect(cs.Subnets()[0].ResourceID).To(Equal("subnet-1"))
		g.Expect(cs.Subnets()[0].IsPublic).To(BeFalse())
		g.Expect(cs.Subnets()[0].NatGatewayID).To(BeNil())
		g.Expect(cs.Subnets()[0].RouteTableID).To(Equal(aws.String("rtb-1")))

		// Public subnet
		g.Expect(cs.Subnets()[1].ID).To(Equal("my-managed-subnet-pub"))
		g.Expect(cs.Subnets()[1].ResourceID).To(Equal("subnet-2"))
		g.Expect(cs.Subnets()[1].IsPublic).To(BeTrue())
		g.Expect(cs.Subnets()[1].NatGatewayID).To(Equal(aws.String("nat-01")))
		g.Expect(cs.Subnets()[1].RouteTableID).To(Equal(aws.String("rtb-2")))
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
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", awsCluster.Name))

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

		_, err = reconciler.reconcileNormal(context.TODO(), cs)
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
		}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", awsCluster.Name))

		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, &awsCluster, controllerIdentity, ns)).To(Succeed())
		})

		awsCluster.Finalizers = []string{infrav1.ClusterFinalizer}
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
		expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{
			{infrav1.LoadBalancerReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.BastionHostReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.SecondaryCidrsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletingReason},
			{infrav1.RouteTablesReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.VpcEndpointsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.NatGatewaysReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.InternetGatewayReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.SubnetsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.VpcReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
		})
	})
}

func mockedDeleteSGCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSecurityGroups(context.TODO(), gomock.Any(), gomock.Any()).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)
}

func mockedDescribeAvailabilityZones(m *mocks.MockEC2APIMockRecorder, zones []string) {
	output := &ec2.DescribeAvailabilityZonesOutput{}
	matcher := gomock.Any()

	if len(zones) > 0 {
		input := &ec2.DescribeAvailabilityZonesInput{}
		for _, zone := range zones {
			input.ZoneNames = append(input.ZoneNames, zone)
			output.AvailabilityZones = append(output.AvailabilityZones, ec2types.AvailabilityZone{
				ZoneName: aws.String(zone),
				ZoneType: aws.String("availability-zone"),
			})
		}

		matcher = gomock.Eq(input)
	}
	m.DescribeAvailabilityZones(context.TODO(), matcher).AnyTimes().
		Return(output, nil)
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
	m.DescribeInstances(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: []string{"bastion"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"pending", "running", "stopping", "stopped"},
			},
		},
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{
			{
				Instances: []ec2types.Instance{
					{
						InstanceId:   aws.String("id-1"),
						InstanceType: ec2types.InstanceTypeM5Large,
						SubnetId:     aws.String("subnet-1"),
						ImageId:      aws.String("ami-1"),
						IamInstanceProfile: &ec2types.IamInstanceProfile{
							Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
						},
						State: &ec2types.InstanceState{
							Code: aws.Int32(16),
							Name: ec2types.InstanceStateNameRunning,
						},
						RootDeviceName: aws.String("device-1"),
						BlockDeviceMappings: []ec2types.InstanceBlockDeviceMapping{
							{
								DeviceName: aws.String("device-1"),
								Ebs: &ec2types.EbsInstanceBlockDevice{
									VolumeId: aws.String("volume-1"),
								},
							},
						},
						Placement: &ec2types.Placement{
							AvailabilityZone: aws.String("us-east-1a"),
						},
					},
				},
			},
		},
	}, nil)
}

func mockedDeleteInstanceAndAwaitTerminationCalls(m *mocks.MockEC2APIMockRecorder) {
	m.TerminateInstances(context.TODO(),
		gomock.Eq(&ec2.TerminateInstancesInput{
			InstanceIds: []string{"id-1"},
		}),
	).Return(&ec2.TerminateInstancesOutput{
		TerminatingInstances: []ec2types.InstanceStateChange{
			{
				InstanceId: aws.String("id-1"),
				CurrentState: &ec2types.InstanceState{
					Name: ec2types.InstanceStateNameShuttingDown,
				},
				PreviousState: &ec2types.InstanceState{
					Name: ec2types.InstanceStateNameRunning,
				},
			},
		},
	}, nil)

	m.DescribeInstances(gomock.Any(), gomock.Eq(&ec2.DescribeInstancesInput{
		InstanceIds: []string{"id-1"},
	}), gomock.Any()).Return(&ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{
			{
				Instances: []ec2types.Instance{
					{
						InstanceId: aws.String("id-1"),
						State: &ec2types.InstanceState{
							Name: ec2types.InstanceStateNameTerminated,
						},
					},
				},
			},
		},
	}, nil)
}

func mockedDeleteInstanceCalls(m *mocks.MockEC2APIMockRecorder) {
	m.TerminateInstances(context.TODO(),
		gomock.Eq(&ec2.TerminateInstancesInput{
			InstanceIds: []string{"id-1"},
		}),
	).Return(nil, nil)
}

func mockedVPCCallsForExistingVPCAndSubnets(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
		NatGateways: []ec2types.NatGateway{
			{
				NatGatewayId: aws.String("nat-01"),
				SubnetId:     aws.String("subnet-2"),
				State:        ec2types.NatGatewayStateAvailable,
			},
		},
	}, nil)
	m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
		Resources: []string{"subnet-1"},
		Tags: []ec2types.Tag{
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
	m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
		Resources: []string{"subnet-2"},
		Tags: []ec2types.Tag{
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
	m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
		},
	})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{
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
				Tags: []ec2types.Tag{
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
	m.DescribeRouteTables(context.TODO(), gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
		},
	})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{
			{
				Routes: []ec2types.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
			},
		},
	}, nil)
	m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
		NatGateways: []ec2types.NatGateway{},
	}, nil)
	m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
		VpcIds: []string{
			"vpc-exists",
		},
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	})).
		Return(&ec2.DescribeVpcsOutput{
			Vpcs: []ec2types.Vpc{
				{
					State:     ec2types.VpcStateAvailable,
					VpcId:     aws.String("vpc-exists"),
					CidrBlock: aws.String("10.0.0.0/8"),
					Tags: []ec2types.Tag{
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

// mockedCallsForMissingEverything mocks most of the AWSCluster reconciliation calls to the AWS API,
// except for what other functions provide (see `mockedCreateSGCalls` and `mockedDescribeInstanceCall`).
func mockedCallsForMissingEverything(m *mocks.MockEC2APIMockRecorder, e *mocks.MockELBAPIMockRecorder, privateSubnetName string, publicSubnetName string) {
	describeVPCByNameCall := m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []string{"test-cluster-vpc"},
			},
		},
	})).Return(&ec2.DescribeVpcsOutput{Vpcs: []ec2types.Vpc{}}, nil)
	m.CreateVpc(context.TODO(), gomock.Eq(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/8"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVpc,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-vpc"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
				},
			},
		},
	})).After(describeVPCByNameCall).Return(&ec2.CreateVpcOutput{
		Vpc: &ec2types.Vpc{
			State:     ec2types.VpcStateAvailable,
			VpcId:     aws.String("vpc-new"),
			CidrBlock: aws.String("10.0.0.0/8"),
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-cluster-vpc"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
					Value: aws.String("common"),
				},
			},
		},
	}, nil)

	m.DescribeVpcAttribute(context.TODO(), gomock.Eq(&ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String("vpc-new"),
		Attribute: ec2types.VpcAttributeNameEnableDnsHostnames,
	})).Return(&ec2.DescribeVpcAttributeOutput{
		EnableDnsHostnames: &ec2types.AttributeBooleanValue{Value: aws.Bool(true)},
	}, nil)

	m.DescribeVpcAttribute(context.TODO(), gomock.Eq(&ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String("vpc-new"),
		Attribute: ec2types.VpcAttributeNameEnableDnsSupport,
	})).Return(&ec2.DescribeVpcAttributeOutput{
		EnableDnsSupport: &ec2types.AttributeBooleanValue{Value: aws.Bool(true)},
	}, nil)

	m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{},
	}, nil)

	subnet1 := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
		VpcId:            aws.String("vpc-new"),
		CidrBlock:        aws.String("10.0.10.0/24"),
		AvailabilityZone: aws.String("us-east-1a"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSubnet,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(privateSubnetName),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("kubernetes.io/role/internal-elb"),
						Value: aws.String("1"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("private"),
					},
				},
			},
		},
	})).Return(&ec2.CreateSubnetOutput{
		Subnet: &ec2types.Subnet{
			VpcId:               aws.String("vpc-new"),
			SubnetId:            aws.String("subnet-1"),
			CidrBlock:           aws.String("10.0.10.0/24"),
			AvailabilityZone:    aws.String("us-east-1a"),
			MapPublicIpOnLaunch: aws.Bool(false),
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(privateSubnetName),
				},
				{
					Key:   aws.String("kubernetes.io/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("kubernetes.io/role/internal-elb"),
					Value: aws.String("1"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
					Value: aws.String("private"),
				},
			},
		},
	}, nil)

	m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
		SubnetIds: []string{"subnet-1"},
	}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{
			{
				VpcId:            aws.String("vpc-new"),
				SubnetId:         aws.String("subnet-1"),
				CidrBlock:        aws.String("10.0.10.0/24"),
				AvailabilityZone: aws.String("us-east-1a"),
				State:            ec2types.SubnetStateAvailable,
			},
		},
	}, nil).After(subnet1)

	subnet2 := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
		VpcId:            aws.String("vpc-new"),
		CidrBlock:        aws.String("10.0.11.0/24"),
		AvailabilityZone: aws.String("us-east-1a"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSubnet,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(publicSubnetName),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("kubernetes.io/role/elb"),
						Value: aws.String("1"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("public"),
					},
				},
			},
		},
	})).Return(&ec2.CreateSubnetOutput{
		Subnet: &ec2types.Subnet{
			VpcId:               aws.String("vpc-new"),
			SubnetId:            aws.String("subnet-2"),
			CidrBlock:           aws.String("10.0.11.0/24"),
			AvailabilityZone:    aws.String("us-east-1a"),
			MapPublicIpOnLaunch: aws.Bool(false),
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(publicSubnetName),
				},
				{
					Key:   aws.String("kubernetes.io/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("kubernetes.io/role/elb"),
					Value: aws.String("1"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
					Value: aws.String("owned"),
				},
				{
					Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
					Value: aws.String("public"),
				},
			},
		},
	}, nil)

	m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
		SubnetIds: []string{"subnet-2"},
	}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{
			{
				VpcId:               aws.String("vpc-new"),
				SubnetId:            aws.String("subnet-2"),
				CidrBlock:           aws.String("10.0.11.0/24"),
				AvailabilityZone:    aws.String("us-east-1a"),
				MapPublicIpOnLaunch: aws.Bool(false),
				State:               ec2types.SubnetStateAvailable,
			},
		},
	}, nil).After(subnet2)

	m.ModifySubnetAttribute(context.TODO(), gomock.Eq(&ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String("subnet-2"),
		MapPublicIpOnLaunch: &ec2types.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	})).Return(&ec2.ModifySubnetAttributeOutput{}, nil)

	m.DescribeRouteTables(context.TODO(), gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{
			{
				Routes: []ec2types.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
			},
		},
	}, nil).MinTimes(1).MaxTimes(2)

	m.DescribeInternetGateways(context.TODO(), gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{"vpc-new"},
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []ec2types.InternetGateway{},
	}, nil)

	m.CreateInternetGateway(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateInternetGatewayInput{})).
		Return(&ec2.CreateInternetGatewayOutput{
			InternetGateway: &ec2types.InternetGateway{
				InternetGatewayId: aws.String("igw-1"),
				Tags: []ec2types.Tag{
					{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster")),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-igw"),
					},
				},
			},
		}, nil)

	m.AttachInternetGateway(context.TODO(), gomock.Eq(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String("igw-1"),
		VpcId:             aws.String("vpc-new"),
	})).
		Return(&ec2.AttachInternetGatewayOutput{}, nil)

	m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-new"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil).MinTimes(1).MaxTimes(2)

	m.DescribeAddresses(context.TODO(), gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: []string{"common"},
			},
		},
	})).Return(&ec2.DescribeAddressesOutput{
		Addresses: []ec2types.Address{},
	}, nil)

	m.AllocateAddress(context.TODO(), gomock.Eq(&ec2.AllocateAddressInput{
		Domain: ec2types.DomainTypeVpc,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeElasticIp,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-eip-common"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
				},
			},
		},
	})).Return(&ec2.AllocateAddressOutput{
		AllocationId: aws.String("1234"),
	}, nil)

	m.CreateNatGateway(context.TODO(), gomock.Eq(&ec2.CreateNatGatewayInput{
		AllocationId: aws.String("1234"),
		SubnetId:     aws.String("subnet-2"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeNatgateway,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-nat"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("common"),
					},
				},
			},
		},
	})).Return(&ec2.CreateNatGatewayOutput{
		NatGateway: &ec2types.NatGateway{
			NatGatewayId: aws.String("nat-01"),
			SubnetId:     aws.String("subnet-2"),
		},
	}, nil)

	m.DescribeNatGateways(gomock.Any(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []string{"nat-01"},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
		NatGateways: []ec2types.NatGateway{
			{
				NatGatewayId: aws.String("nat-01"),
				SubnetId:     aws.String("subnet-2"),
				State:        ec2types.NatGatewayStateAvailable,
			},
		},
	}, nil)

	m.CreateRouteTable(context.TODO(), gomock.Eq(&ec2.CreateRouteTableInput{
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeRouteTable,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-rt-private-us-east-1a"),
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
						Value: aws.String("common"),
					},
				},
			},
		},
		VpcId: aws.String("vpc-new"),
	})).Return(&ec2.CreateRouteTableOutput{
		RouteTable: &ec2types.RouteTable{
			RouteTableId: aws.String("rtb-1"),
		},
	}, nil)

	m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		NatGatewayId:         aws.String("nat-01"),
		RouteTableId:         aws.String("rtb-1"),
	})).Return(&ec2.CreateRouteOutput{}, nil)

	m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String("rtb-1"),
		SubnetId:     aws.String("subnet-1"),
	})).Return(&ec2.AssociateRouteTableOutput{}, nil)

	m.CreateRouteTable(context.TODO(), gomock.Eq(&ec2.CreateRouteTableInput{
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeRouteTable,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-rt-public-us-east-1a"),
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
						Value: aws.String("common"),
					},
				},
			},
		},
		VpcId: aws.String("vpc-new"),
	})).Return(&ec2.CreateRouteTableOutput{
		RouteTable: &ec2types.RouteTable{
			RouteTableId: aws.String("rtb-2"),
		},
	}, nil)

	m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String("igw-1"),
		RouteTableId:         aws.String("rtb-2"),
	})).Return(&ec2.CreateRouteOutput{}, nil)

	m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String("rtb-2"),
		SubnetId:     aws.String("subnet-2"),
	})).Return(&ec2.AssociateRouteTableOutput{}, nil)

	e.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{"test-cluster-apiserver"},
	})).Return(&elb.DescribeLoadBalancersOutput{
		LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{},
	}, nil)

	e.CreateLoadBalancer(gomock.Any(), gomock.Eq(&elb.CreateLoadBalancerInput{
		Listeners: []elbtypes.Listener{
			{
				InstancePort:     aws.Int32(6443),
				InstanceProtocol: aws.String("TCP"),
				LoadBalancerPort: 6443,
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: aws.String("test-cluster-apiserver"),
		Scheme:           aws.String("internet-facing"),
		SecurityGroups:   []string{"sg-apiserver-lb"},
		Subnets:          []string{"subnet-2"},
		Tags: []elbtypes.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-cluster-apiserver"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Value: aws.String("owned"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
				Value: aws.String("apiserver"),
			},
		},
	})).Return(&elb.CreateLoadBalancerOutput{
		DNSName: aws.String("unittest24.de"),
	}, nil)

	e.ConfigureHealthCheck(gomock.Any(), gomock.Eq(&elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String("test-cluster-apiserver"),
		HealthCheck: &elbtypes.HealthCheck{
			Target:             aws.String("TCP:6443"),
			Interval:           aws.Int32(10),
			Timeout:            aws.Int32(5),
			HealthyThreshold:   aws.Int32(5),
			UnhealthyThreshold: aws.Int32(3),
		},
	})).Return(&elb.ConfigureHealthCheckOutput{}, nil)
}

func mockedCreateMaximumVPCCalls(m *mocks.MockEC2APIMockRecorder) {
	describeVPCByNameCall := m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{},
	}, nil)
	m.CreateVpc(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).After(describeVPCByNameCall).Return(nil, errors.New("The maximum number of VPCs has been reached"))
}

func mockedDeleteVPCCallsForNonExistentVPC(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{},
	}, nil).AnyTimes()
	m.DescribeRouteTables(context.TODO(), gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{""},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{},
	}, nil).AnyTimes()
	m.DescribeInternetGateways(context.TODO(), gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{""},
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []ec2types.InternetGateway{},
	}, nil)
	m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{""},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
		NatGateways: []ec2types.NatGateway{},
	}, nil).AnyTimes()
	m.DescribeAddresses(context.TODO(), gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Values: []string{"owned"},
			},
		},
	})).Return(nil, nil)
	m.DeleteVpc(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteVpcInput{
		VpcId: aws.String("vpc-exists"),
	})).Return(nil, nil)
}

func mockedDeleteVPCCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeVpcEndpoints(context.TODO(), gomock.Eq(&ec2.DescribeVpcEndpointsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Values: []string{"owned"},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeVpcEndpointsOutput{}, nil).AnyTimes()
	m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
		},
	})).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{
			{
				VpcId:               aws.String("vpc-exists"),
				SubnetId:            aws.String("subnet-1"),
				AvailabilityZone:    aws.String("us-east-1a"),
				CidrBlock:           aws.String("10.0.10.0/24"),
				MapPublicIpOnLaunch: aws.Bool(false),
			},
		},
	}, nil).AnyTimes()
	m.DescribeRouteTables(context.TODO(), gomock.Eq(&ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(&ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{
			{
				Routes: []ec2types.Route{
					{
						GatewayId: aws.String("igw-12345"),
					},
				},
				RouteTableId: aws.String("rt-12345"),
			},
		},
	}, nil).AnyTimes()
	m.DeleteRouteTable(context.TODO(), gomock.Eq(&ec2.DeleteRouteTableInput{
		RouteTableId: aws.String("rt-12345"),
	}))
	m.DescribeInternetGateways(context.TODO(), gomock.Eq(&ec2.DescribeInternetGatewaysInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{"vpc-exists"},
			},
		},
	})).Return(&ec2.DescribeInternetGatewaysOutput{
		InternetGateways: []ec2types.InternetGateway{
			{
				Attachments:       nil,
				InternetGatewayId: aws.String("ig-12345"),
			},
		},
	}, nil)
	m.DetachInternetGateway(context.TODO(), gomock.Eq(&ec2.DetachInternetGatewayInput{
		VpcId:             aws.String("vpc-exists"),
		InternetGatewayId: aws.String("ig-12345"),
	}))
	m.DeleteInternetGateway(context.TODO(), gomock.Eq(&ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String("ig-12345"),
	}))
	m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
		Filter: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-exists"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil).AnyTimes()
	m.DescribeAddresses(context.TODO(), gomock.Eq(&ec2.DescribeAddressesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Values: []string{"owned"},
			},
		},
	})).Return(&ec2.DescribeAddressesOutput{
		Addresses: []ec2types.Address{
			{
				AssociationId: aws.String("1234"),
				AllocationId:  aws.String("1234"),
				PublicIp:      aws.String("1.2.3.4"),
			},
		},
	}, nil)
	m.DisassociateAddress(context.TODO(), &ec2.DisassociateAddressInput{
		AssociationId: aws.String("1234"),
	})
	m.ReleaseAddress(context.TODO(), &ec2.ReleaseAddressInput{
		AllocationId: aws.String("1234"),
	})
	m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
		VpcIds: []string{
			"vpc-exists",
		},
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{string(ec2types.VpcStatePending), string(ec2types.VpcStateAvailable)},
			},
		},
	})).
		Return(&ec2.DescribeVpcsOutput{
			Vpcs: []ec2types.Vpc{
				{
					State:     ec2types.VpcStateAvailable,
					VpcId:     aws.String("vpc-exists"),
					CidrBlock: aws.String("10.0.0.0/8"),
					Tags: []ec2types.Tag{
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
	m.DeleteSubnet(context.TODO(), gomock.Eq(&ec2.DeleteSubnetInput{
		SubnetId: aws.String("subnet-1"),
	}))
	m.DeleteVpc(context.TODO(), gomock.Eq(&ec2.DeleteVpcInput{
		VpcId: aws.String("vpc-exists"),
	}))
}

func mockedCreateSGCalls(recordLBV2 bool, vpcID string, m *mocks.MockEC2APIMockRecorder) {
	m.DescribeSecurityGroups(context.TODO(), gomock.Eq(&ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
			{
				Name:   aws.String("tag-key"),
				Values: []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"},
			},
		},
	})).Return(
		&ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []ec2types.SecurityGroup{
				{
					GroupId:   aws.String("1"),
					GroupName: aws.String("test-sg"),
				},
			},
		}, nil)
	m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String("test-cluster-bastion"),
		Description: aws.String("Kubernetes cluster test-cluster: bastion"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
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
	m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String("test-cluster-apiserver-lb"),
		Description: aws.String("Kubernetes cluster test-cluster: apiserver-lb"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
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
	m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String("test-cluster-lb"),
		Description: aws.String("Kubernetes cluster test-cluster: lb"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
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
	securityGroupControl := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String("test-cluster-controlplane"),
		Description: aws.String("Kubernetes cluster test-cluster: controlplane"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
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
	securityGroupNode := m.CreateSecurityGroup(context.TODO(), gomock.Eq(&ec2.CreateSecurityGroupInput{
		VpcId:       aws.String(vpcID),
		GroupName:   aws.String("test-cluster-node"),
		Description: aws.String("Kubernetes cluster test-cluster: node"),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSecurityGroup,
				Tags: []ec2types.Tag{
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
	m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-controlplane"),
	})).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
		After(securityGroupControl).Times(2)
	m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-node"),
	})).
		Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
		After(securityGroupNode).Times(2)
	if recordLBV2 {
		m.AuthorizeSecurityGroupIngress(context.TODO(), gomock.AssignableToTypeOf(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String("sg-lb"),
		})).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).
			After(securityGroupNode).Times(1)
	}
}
