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

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go/aws"
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

		ec2Mock.EXPECT().AssociateAddress(context.TODO(), gomock.Any()).MaxTimes(1)

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

		_, err = reconciler.reconcileDelete(context.TODO(), ms, cs, cs, cs, cs)
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
			e.DescribeLoadBalancers(ctx, gomock.Eq(&elb.DescribeLoadBalancersInput{
				LoadBalancerNames: []string{"test-cluster-apiserver"},
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

		ec2Mock.EXPECT().AssociateAddress(context.TODO(), gomock.Any()).MaxTimes(1)

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
			m.TerminateInstances(context.TODO(),
				gomock.Eq(&ec2.TerminateInstancesInput{
					InstanceIds: []string{"id-1"},
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

		_, err = reconciler.reconcileDelete(context.TODO(), ms, cs, cs, cs, cs)
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
	m.DescribeInstances(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Values: []string{"owned"},
			},
			{
				Name:   aws.String("tag:Name"),
				Values: []string{"test"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"pending", "running"},
			},
		},
	})).Return(&ec2.DescribeInstancesOutput{}, nil)
	m.DescribeInstanceTypes(context.TODO(), gomock.Any()).
		Return(&ec2.DescribeInstanceTypesOutput{
			InstanceTypes: []ec2types.InstanceTypeInfo{
				{
					ProcessorInfo: &ec2types.ProcessorInfo{
						SupportedArchitectures: []ec2types.ArchitectureType{
							ec2types.ArchitectureTypeX8664,
						},
					},
				},
			},
		}, nil)
	m.DescribeImages(context.TODO(), gomock.Eq(&ec2.DescribeImagesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("owner-id"),
				Values: []string{"819546954734"},
			},
			{
				Name:   aws.String("name"),
				Values: []string{"capa-ami-ubuntu-24.04-?test-*"},
			},
			{
				Name:   aws.String("architecture"),
				Values: []string{"x86_64"},
			},
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
	})).Return(&ec2.DescribeImagesOutput{Images: []ec2types.Image{
		{
			ImageId:      aws.String("latest"),
			CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
		},
	}}, nil)
	m.RunInstances(context.TODO(), gomock.Any()).Return(&ec2.RunInstancesOutput{
		Instances: []ec2types.Instance{
			{
				State: &ec2types.InstanceState{
					Name: ec2types.InstanceStateNameRunning,
				},
				IamInstanceProfile: &ec2types.IamInstanceProfile{
					Arn: aws.String("arn:aws:iam::123456789012:instance-profile/foo"),
				},
				InstanceId:     aws.String("two"),
				InstanceType:   ec2types.InstanceTypeM5Large,
				SubnetId:       aws.String("subnet-1"),
				ImageId:        aws.String("ami-1"),
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
	}, nil)
	m.DescribeNetworkInterfaces(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfacesInput{Filters: []ec2types.Filter{
		{
			Name:   aws.String("attachment.instance-id"),
			Values: []string{"two"},
		},
	}})).Return(&ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []ec2types.NetworkInterface{
			{
				NetworkInterfaceId: aws.String("eni-1"),
				Groups: []ec2types.GroupIdentifier{
					{
						GroupId: aws.String("3"),
					},
				},
			},
		},
	}, nil).MaxTimes(3)
	m.DescribeNetworkInterfaceAttribute(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("eni-1"),
		Attribute:          ec2types.NetworkInterfaceAttributeGroupSet,
	})).Return(&ec2.DescribeNetworkInterfaceAttributeOutput{Groups: []ec2types.GroupIdentifier{{GroupId: aws.String("3")}}}, nil).MaxTimes(1)
	m.ModifyNetworkInterfaceAttribute(context.TODO(), gomock.Any()).AnyTimes()
	m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{Filters: []ec2types.Filter{
		{
			Name:   aws.String("state"),
			Values: []string{"pending", "available"},
		},
		{
			Name:   aws.String("subnet-id"),
			Values: []string{"subnet-1"},
		},
	}})).Return(&ec2.DescribeSubnetsOutput{Subnets: []ec2types.Subnet{
		{
			SubnetId: aws.String("subnet-1"),
		},
	}}, nil)
}

func mockedDescribeInstanceCalls(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeInstances(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		InstanceIds: []string{"myMachine"},
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{{Instances: []ec2types.Instance{{Placement: &ec2types.Placement{AvailabilityZone: aws.String("us-east-1a")}, InstanceId: aws.String("id-1"), State: &ec2types.InstanceState{Name: "id-1", Code: aws.Int32(16)}}}}},
	}, nil)
}
