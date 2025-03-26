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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	ec2Service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	elbService "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func TestAWSMachineReconcilerIntegrationTests(t *testing.T) {
	var (
		reconciler AWSMachineReconciler
		mockCtrl   *gomock.Controller
		recorder   *record.FakeRecorder
	)

	setup := func(t *testing.T, g *WithT) {
		t.Helper()
		mockCtrl = gomock.NewController(t)
		recorder = record.NewFakeRecorder(10)
		reconciler = AWSMachineReconciler{
			Client:   testEnv.Client,
			Recorder: recorder,
		}
	}

	teardown := func(g *WithT) {
		mockCtrl.Finish()
	}

	t.Run("Should successfully reconcile control plane machine creation", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		secretMock := mock_services.NewMockSecretInterface(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)

		expect := func(m *mocks.MockEC2APIMockRecorder, s *mock_services.MockSecretInterfaceMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedCreateInstanceCalls(m)
			mockedCreateSecretCall(s)
			mockedCreateLBCalls(t, e, false)
		}
		expect(ec2Mock.EXPECT(), secretMock.EXPECT(), elbMock.EXPECT())

		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bootstrap-data",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		setup(t, g)
		awsMachine := getAWSMachine()
		awsMachine.Namespace = ns.Name
		createAWSMachine(g, awsMachine)

		defer teardown(g)
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, awsMachine, ns, secret)).To(Succeed())
		})

		cs, err := getClusterScope(infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}, Spec: infrav1.AWSClusterSpec{NetworkSpec: infrav1.NetworkSpec{
			Subnets: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
				},
			},
		}}})
		g.Expect(err).To(BeNil())
		cs.Cluster = &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"}}
		cs.AWSCluster.Spec.NetworkSpec.VPC = infrav1.VPCSpec{
			ID:        "vpc-exists",
			CidrBlock: "10.0.0.0/16",
		}
		cs.AWSCluster.Status.Network.APIServerELB.DNSName = DNSName
		cs.AWSCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeClassic,
		}
		cs.AWSCluster.Status.Network.SecurityGroups = map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
			infrav1.SecurityGroupNode: {
				ID: "1",
			},
			infrav1.SecurityGroupLB: {
				ID: "2",
			},
			infrav1.SecurityGroupControlPlane: {
				ID: "3",
			},
		}
		ms, err := getMachineScope(cs, awsMachine)
		g.Expect(err).To(BeNil())

		ms.Machine.Spec.Bootstrap.DataSecretName = aws.String("bootstrap-data")
		ms.Machine.Spec.Version = aws.String("test")
		ms.AWSMachine.Spec.Subnet = &infrav1.AWSResourceReference{ID: aws.String("subnet-1")}
		ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateRunning
		ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		reconciler.elbServiceFactory = func(scope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		ec2Mock.EXPECT().AssociateAddressWithContext(context.TODO(), gomock.Any()).MaxTimes(1)

		reconciler.secretsManagerServiceFactory = func(clusterScope cloud.ClusterScoper) services.SecretInterface {
			return secretMock
		}

		_, err = reconciler.reconcileNormal(ctx, ms, cs, cs, cs, cs)
		g.Expect(err).To(BeNil())
		expectConditions(g, ms.AWSMachine, []conditionAssertion{
			{infrav1.SecurityGroupsReadyCondition, corev1.ConditionTrue, "", ""},
			{infrav1.InstanceReadyCondition, corev1.ConditionTrue, "", ""},
			{infrav1.ELBAttachedCondition, corev1.ConditionTrue, "", ""},
		})
		g.Expect(ms.AWSMachine.Finalizers).Should(ContainElement(infrav1.MachineFinalizer))
	})
	t.Run("Should successfully reconcile control plane machine deletion", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		elbv2Mock := mocks.NewMockELBV2API(mockCtrl)

		expect := func(m *mocks.MockEC2APIMockRecorder, ev2 *mocks.MockELBV2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedDescribeInstanceCalls(m)
			mockedDeleteLBCalls(false, ev2, e)
			mockedDeleteInstanceCalls(m)
		}
		expect(ec2Mock.EXPECT(), elbv2Mock.EXPECT(), elbMock.EXPECT())

		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bootstrap-data",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		setup(t, g)
		awsMachine := getAWSMachine()
		awsMachine.Namespace = ns.Name
		createAWSMachine(g, awsMachine)

		defer teardown(g)
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, awsMachine, ns)).To(Succeed())
		})

		cs, err := getClusterScope(infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}})
		g.Expect(err).To(BeNil())
		cs.Cluster = &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"}}
		cs.AWSCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeClassic,
		}
		ms, err := getMachineScope(cs, awsMachine)
		g.Expect(err).To(BeNil())

		ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateRunning
		ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
		ms.AWSMachine.Spec.ProviderID = aws.String("aws:////myMachine")

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		elbSvc.ELBV2Client = elbv2Mock
		reconciler.elbServiceFactory = func(scope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		_, err = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
		g.Expect(err).To(BeNil())
		expectConditions(g, ms.AWSMachine, []conditionAssertion{
			{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
			{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
		})
		g.Expect(ms.AWSMachine.Finalizers).ShouldNot(ContainElement(infrav1.MachineFinalizer))
	})
	t.Run("Should fail reconciling control-plane machine creation while attaching load balancer", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		secretMock := mock_services.NewMockSecretInterface(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)

		expect := func(m *mocks.MockEC2APIMockRecorder, s *mock_services.MockSecretInterfaceMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedCreateInstanceCalls(m)
			mockedCreateSecretCall(s)
			e.DescribeLoadBalancers(gomock.Eq(&elb.DescribeLoadBalancersInput{
				LoadBalancerNames: aws.StringSlice([]string{"test-cluster-apiserver"}),
			})).
				Return(&elb.DescribeLoadBalancersOutput{}, nil)
		}
		expect(ec2Mock.EXPECT(), secretMock.EXPECT(), elbMock.EXPECT())

		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bootstrap-data",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		setup(t, g)
		awsMachine := getAWSMachine()
		awsMachine.Namespace = ns.Name
		createAWSMachine(g, awsMachine)

		defer teardown(g)
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, awsMachine, ns, secret)).To(Succeed())
		})

		cs, err := getClusterScope(infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}, Spec: infrav1.AWSClusterSpec{NetworkSpec: infrav1.NetworkSpec{
			Subnets: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
				},
			},
		}}})
		g.Expect(err).To(BeNil())
		cs.Cluster = &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"}}
		cs.AWSCluster.Status.Network.APIServerELB.DNSName = DNSName
		cs.AWSCluster.Spec.NetworkSpec.VPC = infrav1.VPCSpec{
			ID:        "vpc-exists",
			CidrBlock: "10.0.0.0/16",
		}
		cs.AWSCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeClassic,
		}
		cs.AWSCluster.Status.Network.SecurityGroups = map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
			infrav1.SecurityGroupNode: {
				ID: "1",
			},
			infrav1.SecurityGroupLB: {
				ID: "2",
			},
			infrav1.SecurityGroupControlPlane: {
				ID: "3",
			},
		}
		ms, err := getMachineScope(cs, awsMachine)
		g.Expect(err).To(BeNil())

		ms.Machine.Spec.Bootstrap.DataSecretName = aws.String("bootstrap-data")
		ms.Machine.Spec.Version = aws.String("test")
		ms.AWSMachine.Spec.Subnet = &infrav1.AWSResourceReference{ID: aws.String("subnet-1")}
		ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateRunning
		ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		reconciler.elbServiceFactory = func(scope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		reconciler.secretsManagerServiceFactory = func(clusterScope cloud.ClusterScoper) services.SecretInterface {
			return secretMock
		}

		ec2Mock.EXPECT().AssociateAddressWithContext(context.TODO(), gomock.Any()).MaxTimes(1)

		_, err = reconciler.reconcileNormal(ctx, ms, cs, cs, cs, cs)
		g.Expect(err).Should(HaveOccurred())
		expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionTrue, "", ""}})
		g.Expect(ms.AWSMachine.Finalizers).Should(ContainElement(infrav1.MachineFinalizer))
	})
	t.Run("Should fail in reconciling control-plane machine deletion while terminating instance ", func(t *testing.T) {
		g := NewWithT(t)
		mockCtrl = gomock.NewController(t)
		ec2Mock := mocks.NewMockEC2API(mockCtrl)
		elbMock := mocks.NewMockELBAPI(mockCtrl)
		elbv2Mock := mocks.NewMockELBV2API(mockCtrl)

		expect := func(m *mocks.MockEC2APIMockRecorder, ev2 *mocks.MockELBV2APIMockRecorder, e *mocks.MockELBAPIMockRecorder) {
			mockedDescribeInstanceCalls(m)
			mockedDeleteLBCalls(false, ev2, e)
			m.TerminateInstancesWithContext(context.TODO(),
				gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: aws.StringSlice([]string{"id-1"}),
				}),
			).
				Return(nil, errors.New("Failed to delete instance"))
		}
		expect(ec2Mock.EXPECT(), elbv2Mock.EXPECT(), elbMock.EXPECT())

		ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("integ-test-%s", util.RandomString(5)))
		g.Expect(err).To(BeNil())

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bootstrap-data",
				Namespace: ns.Name,
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		setup(t, g)
		awsMachine := getAWSMachine()
		awsMachine.Namespace = ns.Name
		createAWSMachine(g, awsMachine)

		defer teardown(g)
		defer t.Cleanup(func() {
			g.Expect(testEnv.Cleanup(ctx, awsMachine, ns)).To(Succeed())
		})

		cs, err := getClusterScope(infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}})
		g.Expect(err).To(BeNil())
		cs.Cluster = &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"}}
		cs.AWSCluster.Spec.ControlPlaneLoadBalancer = &infrav1.AWSLoadBalancerSpec{
			LoadBalancerType: infrav1.LoadBalancerTypeClassic,
		}
		ms, err := getMachineScope(cs, awsMachine)
		g.Expect(err).To(BeNil())

		ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateRunning
		ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
		ms.AWSMachine.Spec.ProviderID = aws.String("aws:////myMachine")

		ec2Svc := ec2Service.NewService(cs)
		ec2Svc.EC2Client = ec2Mock
		reconciler.ec2ServiceFactory = func(scope scope.EC2Scope) services.EC2Interface {
			return ec2Svc
		}

		elbSvc := elbService.NewService(cs)
		elbSvc.EC2Client = ec2Mock
		elbSvc.ELBClient = elbMock
		elbSvc.ELBV2Client = elbv2Mock
		reconciler.elbServiceFactory = func(scope scope.ELBScope) services.ELBInterface {
			return elbSvc
		}

		_, err = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
		g.Expect(err).Should(HaveOccurred())
		expectConditions(g, ms.AWSMachine, []conditionAssertion{
			{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, "DeletingFailed"},
			{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason},
		})
		g.Expect(ms.AWSMachine.Finalizers).ShouldNot(ContainElement(infrav1.MachineFinalizer))
	})
}

func getMachineScope(cs *scope.ClusterScope, awsMachine *infrav1.AWSMachine) (*scope.MachineScope, error) {
	return scope.NewMachineScope(
		scope.MachineScopeParams{
			Client: testEnv,
			Cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Status: clusterv1.ClusterStatus{
					InfrastructureReady: true,
				},
			},
			Machine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: clusterv1.MachineSpec{
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: ptr.To[string]("bootstrap-data"),
					},
				},
			},
			InfraCluster: cs,
			AWSMachine:   awsMachine,
		},
	)
}

func createAWSMachine(g *WithT, awsMachine *infrav1.AWSMachine) {
	g.Expect(testEnv.Create(ctx, awsMachine)).To(Succeed())
	g.Eventually(func() bool {
		machine := &infrav1.AWSMachine{}
		key := client.ObjectKey{
			Name:      awsMachine.Name,
			Namespace: awsMachine.Namespace,
		}
		return testEnv.Get(ctx, key, machine) == nil
	}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed get the newly created machine %q", awsMachine.Name))
}

func getAWSMachine() *infrav1.AWSMachine {
	return &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: infrav1.AWSMachineSpec{
			CloudInit: infrav1.CloudInit{
				SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				SecretPrefix:         "prefix",
				SecretCount:          1000,
			},
			InstanceType: "test",
			Subnet:       &infrav1.AWSResourceReference{ID: aws.String("subnet-1")},
		},
	}
}

func getAWSMachineWithAdditionalTags() *infrav1.AWSMachine {
	return &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: infrav1.AWSMachineSpec{
			CloudInit: infrav1.CloudInit{
				SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
			},
			AdditionalTags: map[string]string{"foo": "bar"},
		},
	}
}

func PointsTo(s string) gomock.Matcher {
	return &pointsTo{
		val: s,
	}
}

type pointsTo struct {
	val string
}

func (p *pointsTo) Matches(x interface{}) bool {
	ptr, ok := x.(*string)
	if !ok {
		return false
	}

	if ptr == nil {
		return false
	}

	return *ptr == p.val
}

func (p *pointsTo) String() string {
	return fmt.Sprintf("Pointer to string %q", p.val)
}

type conditionAssertion struct {
	conditionType clusterv1.ConditionType
	status        corev1.ConditionStatus
	severity      clusterv1.ConditionSeverity
	reason        string
}

func expectConditions(g *WithT, m *infrav1.AWSMachine, expected []conditionAssertion) {
	g.Expect(len(m.Status.Conditions)).To(BeNumerically(">=", len(expected)), "number of conditions")
	for _, c := range expected {
		actual := conditions.Get(m, c.conditionType)
		g.Expect(actual).To(Not(BeNil()))
		g.Expect(actual.Type).To(Equal(c.conditionType))
		g.Expect(actual.Status).To(Equal(c.status))
		g.Expect(actual.Severity).To(Equal(c.severity))
		g.Expect(actual.Reason).To(Equal(c.reason))
	}
}

func mockedCreateSecretCall(s *mock_services.MockSecretInterfaceMockRecorder) {
	s.Create(gomock.AssignableToTypeOf(&scope.MachineScope{}), gomock.AssignableToTypeOf([]byte{}))
	s.UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.AssignableToTypeOf([]scope.ServiceEndpoint{}))
}

func mockedCreateInstanceCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Values: aws.StringSlice([]string{"owned"}),
			},
			{
				Name:   aws.String("tag:Name"),
				Values: aws.StringSlice([]string{"test"}),
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: aws.StringSlice([]string{"pending", "running"}),
			},
		},
	})).Return(&ec2.DescribeInstancesOutput{}, nil)
	m.DescribeInstanceTypesWithContext(context.TODO(), gomock.Any()).
		Return(&ec2.DescribeInstanceTypesOutput{
			InstanceTypes: []*ec2.InstanceTypeInfo{
				{
					ProcessorInfo: &ec2.ProcessorInfo{
						SupportedArchitectures: []*string{
							aws.String("x86_64"),
						},
					},
				},
			},
		}, nil)
	m.DescribeImagesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: aws.StringSlice([]string{"819546954734"}),
			},
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{"capa-ami-ubuntu-24.04-?test-*"}),
			},
			{
				Name:   aws.String("architecture"),
				Values: aws.StringSlice([]string{"x86_64"}),
			},
			{
				Name:   aws.String("state"),
				Values: aws.StringSlice([]string{"available"}),
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: aws.StringSlice([]string{"hvm"}),
			},
		},
	})).Return(&ec2.DescribeImagesOutput{Images: []*ec2.Image{
		{
			ImageId:      aws.String("latest"),
			CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
		},
	}}, nil)
	m.RunInstancesWithContext(context.TODO(), gomock.Any()).Return(&ec2.Reservation{
		Instances: []*ec2.Instance{
			{
				State: &ec2.InstanceState{
					Name: aws.String(ec2.InstanceStateNameRunning),
				},
				IamInstanceProfile: &ec2.IamInstanceProfile{
					Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
				},
				InstanceId:     aws.String("two"),
				InstanceType:   aws.String("m5.large"),
				SubnetId:       aws.String("subnet-1"),
				ImageId:        aws.String("ami-1"),
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
	}, nil)
	m.DescribeNetworkInterfacesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfacesInput{Filters: []*ec2.Filter{
		{
			Name:   aws.String("attachment.instance-id"),
			Values: aws.StringSlice([]string{"two"}),
		},
	}})).Return(&ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []*ec2.NetworkInterface{
			{
				NetworkInterfaceId: aws.String("eni-1"),
				Groups: []*ec2.GroupIdentifier{
					{
						GroupId: aws.String("3"),
					},
				},
			},
		},
	}, nil).MaxTimes(3)
	m.DescribeNetworkInterfaceAttributeWithContext(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("eni-1"),
		Attribute:          aws.String("groupSet"),
	})).Return(&ec2.DescribeNetworkInterfaceAttributeOutput{Groups: []*ec2.GroupIdentifier{{GroupId: aws.String("3")}}}, nil).MaxTimes(1)
	m.ModifyNetworkInterfaceAttributeWithContext(context.TODO(), gomock.Any()).AnyTimes()
	m.DescribeSubnetsWithContext(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{
		{
			Name:   aws.String("state"),
			Values: aws.StringSlice([]string{"pending", "available"}),
		},
		{
			Name:   aws.String("subnet-id"),
			Values: aws.StringSlice([]string{"subnet-1"}),
		},
	}})).Return(&ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{
		{
			SubnetId: aws.String("subnet-1"),
		},
	}}, nil)
}

func mockedDescribeInstanceCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{"myMachine"}),
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{{Instances: []*ec2.Instance{{Placement: &ec2.Placement{AvailabilityZone: aws.String("us-east-1a")}, InstanceId: aws.String("id-1"), State: &ec2.InstanceState{Name: aws.String("id-1"), Code: aws.Int64(16)}}}}},
	}, nil)
}
