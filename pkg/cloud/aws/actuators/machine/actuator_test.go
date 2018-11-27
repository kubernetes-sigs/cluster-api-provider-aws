package machine

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"

	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client/mock"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/cluster-api-provider-aws/test/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	clusterv1.AddToScheme(scheme.Scheme)
}

const (
	controllerLogName = "awsMachine"

	defaultNamespace         = "default"
	defaultAvailabilityZone  = "us-east-1a"
	region                   = "us-east-1"
	awsCredentialsSecretName = "aws-credentials-secret"
	userDataSecretName       = "aws-actuator-user-data-secret"

	keyName   = "aws-actuator-key-name"
	clusterID = "aws-actuator-cluster"
)

const userDataBlob = `#cloud-config
write_files:
- path: /root/node_bootstrap/node_settings.yaml
  owner: 'root:root'
  permissions: '0640'
  content: |
    node_config_name: node-config-master
runcmd:
- [ cat, /root/node_bootstrap/node_settings.yaml]
`

func testMachineAPIResources(clusterID string) (*clusterv1.Machine, *clusterv1.Cluster, *apiv1.Secret, *apiv1.Secret, error) {
	awsCredentialsSecret := utils.GenerateAwsCredentialsSecretFromEnv(awsCredentialsSecretName, defaultNamespace)

	userDataSecret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userDataSecretName,
			Namespace: defaultNamespace,
		},
		Data: map[string][]byte{
			userDataSecretKey: []byte(userDataBlob),
		},
	}

	machinePc := &providerconfigv1.AWSMachineProviderConfig{
		AMI: providerconfigv1.AWSResourceReference{
			ID: aws.String("ami-a9acbbd6"),
		},
		CredentialsSecret: &corev1.LocalObjectReference{
			Name: awsCredentialsSecretName,
		},
		InstanceType: "m4.xlarge",
		Placement: providerconfigv1.Placement{
			Region:           region,
			AvailabilityZone: defaultAvailabilityZone,
		},
		Subnet: providerconfigv1.AWSResourceReference{
			ID: aws.String("subnet-0e56b13a64ff8a941"),
		},
		IAMInstanceProfile: &providerconfigv1.AWSResourceReference{
			ID: aws.String("openshift_master_launch_instances"),
		},
		KeyName: aws.String(keyName),
		UserDataSecret: &corev1.LocalObjectReference{
			Name: userDataSecretName,
		},
		Tags: []providerconfigv1.TagSpecification{
			{Name: "openshift-node-group-config", Value: "node-config-master"},
			{Name: "host-type", Value: "master"},
			{Name: "sub-host-type", Value: "default"},
		},
		SecurityGroups: []providerconfigv1.AWSResourceReference{
			{ID: aws.String("sg-00868b02fbe29de17")}, // aws-actuator
			{ID: aws.String("sg-0a4658991dc5eb40a")}, // aws-actuator_master
			{ID: aws.String("sg-009a70e28fa4ba84e")}, // aws-actuator_master_k8s
			{ID: aws.String("sg-07323d56fb932c84c")}, // aws-actuator_infra
			{ID: aws.String("sg-08b1ffd32874d59a2")}, // aws-actuator_infra_k8s
		},
		PublicIP: aws.Bool(true),
		LoadBalancers: []providerconfigv1.LoadBalancerReference{
			{
				Name: "cluster-con",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
			{
				Name: "cluster-ext",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
			{
				Name: "cluster-int",
				Type: providerconfigv1.ClassicLoadBalancerType,
			},
		},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed creating codec: %v", err)
	}
	config, err := codec.EncodeProviderConfig(machinePc)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("encodeToProviderConfig failed: %v", err)
	}

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-actuator-testing-machine",
			Namespace: defaultNamespace,
			Labels: map[string]string{
				providerconfigv1.ClusterIDLabel:   clusterID,
				providerconfigv1.MachineRoleLabel: "infra",
				providerconfigv1.MachineTypeLabel: "master",
			},
		},

		Spec: clusterv1.MachineSpec{
			ProviderConfig: *config,
		},
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID,
			Namespace: defaultNamespace,
		},
	}

	return machine, cluster, awsCredentialsSecret, userDataSecret, nil
}

func TestCreateAndDeleteMachine(t *testing.T) {
	cases := []struct {
		name                string
		createErrorExpected bool
	}{
		{
			name:                "machine creation succeeds",
			createErrorExpected: false,
		},
		{
			name:                "machine creation fails",
			createErrorExpected: true,
		},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// kube client is needed to fetch aws credentials:
			// - kubeClient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
			// cluster client for updating machine statues
			// - clusterClient.ClusterV1alpha1().Machines(machineCopy.Namespace).UpdateStatus(machineCopy)
			machine, cluster, awsCredentialsSecret, userDataSecret, err := testMachineAPIResources(clusterID)
			if err != nil {
				t.Fatal(err)
			}

			fakeKubeClient := kubernetesfake.NewSimpleClientset(awsCredentialsSecret, userDataSecret)
			//fakeClient := fake.NewSimpleClientset(machine)

			fakeClient := fake.NewFakeClient(machine)

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			params := ActuatorParams{
				Client:     fakeClient,
				KubeClient: fakeKubeClient,
				AwsClientBuilder: func(kubeClient kubernetes.Interface, secretName, namespace, region string) (awsclient.Client, error) {
					return mockAWSClient, nil
				},
				Codec: codec,
			}

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create AWS machine actuator: %v", err)
			}

			mockRunInstances(mockAWSClient, tc.createErrorExpected)
			mockDescribeInstances(mockAWSClient, tc.createErrorExpected)
			mockTerminateInstances(mockAWSClient)
			mockRegisterInstancesWithLoadBalancer(mockAWSClient, tc.createErrorExpected)

			// Create the machine
			createErr := actuator.Create(cluster, machine)

			// Get updated machine object from the cluster client
			key := types.NamespacedName{
				Namespace: machine.Namespace,
				Name:      machine.Name,
			}
			updatedMachine := clusterv1.Machine{}
			err = fakeClient.Get(context.Background(), client.ObjectKey(key), &updatedMachine)
			if err != nil {
				t.Fatalf("Unable to retrieve machine: %v", err)
			}

			codec, err := providerconfigv1.NewCodec()
			if err != nil {
				t.Fatalf("error creating codec: %v", err)
			}
			machineStatus, err := ProviderStatusFromMachine(codec, &updatedMachine)
			if err != nil {
				t.Fatalf("error getting machineStatus: %v", err)
			}

			if tc.createErrorExpected {
				assert.Error(t, createErr)
				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationFailed)
			} else {
				assert.NoError(t, createErr)
				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationSucceeded)
			}
			if !tc.createErrorExpected {
				// Get the machine
				if exists, err := actuator.Exists(cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// TODO(jchaloup): Wait until the machine is ready

				// Update a machine
				if err := actuator.Update(cluster, machine); err != nil {
					t.Errorf("Unable to create instance for machine: %v", err)
				}

				// Get the machine
				if exists, err := actuator.Exists(cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Delete a machine
				if err := actuator.Delete(cluster, machine); err != nil {
					t.Errorf("Unable to delete instance for machine: %v", err)
				}
			}
		})
	}
}

func mockRunInstances(mockAWSClient *mockaws.MockClient, genError bool) {
	var err error

	if genError {
		err = errors.New("requested RunInstances error")
	}

	mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(
		&ec2.Reservation{
			Instances: []*ec2.Instance{
				{
					ImageId:    aws.String("ami-a9acbbd6"),
					InstanceId: aws.String("i-02fcb933c5da7085c"),
					State: &ec2.InstanceState{
						Name: aws.String("Running"),
						Code: aws.Int64(16),
					},
					LaunchTime: aws.Time(time.Now()),
				},
			},
		}, err)
}

func mockDescribeInstances(mockAWSClient *mockaws.MockClient, genError bool) {
	var err error

	if genError {
		err = errors.New("requested RunInstances error")
	}

	mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(
		&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							ImageId:    aws.String("ami-a9acbbd6"),
							InstanceId: aws.String("i-02fcb933c5da7085c"),
							State: &ec2.InstanceState{
								Name: aws.String("Running"),
								Code: aws.Int64(16),
							},
							LaunchTime: aws.Time(time.Now()),
						},
					},
				},
			},
		}, err).AnyTimes()
}

func mockTerminateInstances(mockAWSClient *mockaws.MockClient) {
	mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(
		&ec2.TerminateInstancesOutput{}, nil)
}

func mockRegisterInstancesWithLoadBalancer(mockAWSClient *mockaws.MockClient, createError bool) {
	if createError {
		return
	}
	// RegisterInstancesWithLoadBalancer should be called for every load balancer name in the machine
	// spec for create and for update (3 * 2 = 6)
	for i := 0; i < 6; i++ {
		mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any())
	}
}

func TestAvailabiltyZone(t *testing.T) {
	cases := []struct {
		name             string
		availabilityZone string
		subnet           string
	}{
		{
			name:             "availability zone only",
			availabilityZone: "us-east-1a",
		},
		{
			name:   "subnet only",
			subnet: "subnet-b46032ec",
		},
		{
			name:             "availability zone and subnet",
			availabilityZone: "us-east-1a",
			subnet:           "subnet-b46032ec",
		},
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			machine, cluster, awsCredentialsSecret, userDataSecret, err := testMachineAPIResources(clusterID)
			if err != nil {
				t.Fatal(err)
			}

			machinePc := &providerconfigv1.AWSMachineProviderConfig{}
			if err = codec.DecodeProviderConfig(&machine.Spec.ProviderConfig, machinePc); err != nil {
				t.Fatal(err)
			}

			machinePc.Placement.AvailabilityZone = tc.availabilityZone
			if tc.subnet == "" {
				machinePc.Subnet.ID = nil
			} else {
				machinePc.Subnet.ID = aws.String(tc.subnet)
			}

			config, err := codec.EncodeProviderConfig(machinePc)
			if err != nil {
				t.Fatal(err)
			}
			machine.Spec.ProviderConfig = *config

			fakeKubeClient := kubernetesfake.NewSimpleClientset(awsCredentialsSecret, userDataSecret)

			fakeClient := fake.NewFakeClient(machine)

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			params := ActuatorParams{
				Client:     fakeClient,
				KubeClient: fakeKubeClient,
				AwsClientBuilder: func(kubeClient kubernetes.Interface, secretName, namespace, region string) (awsclient.Client, error) {
					return mockAWSClient, nil
				},
				Codec: codec,
			}

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create AWS machine actuator: %v", err)
			}

			mockRunInstancesForPlacement(mockAWSClient, tc.availabilityZone, tc.subnet)
			mockDescribeInstances(mockAWSClient, false)
			mockTerminateInstances(mockAWSClient)
			mockRegisterInstancesWithLoadBalancer(mockAWSClient, false)
			mockDescribeSubnets(mockAWSClient)

			actuator.Create(cluster, machine)
		})
	}
}

func mockDescribeSubnets(mockAWSClient *mockaws.MockClient) {
	mockAWSClient.EXPECT().DescribeSubnets(gomock.Any()).Return(&ec2.DescribeSubnetsOutput{}, nil)
}

func mockRunInstancesForPlacement(mockAWSClient *mockaws.MockClient, availabilityZone, subnet string) {
	var placement *ec2.Placement
	if availabilityZone != "" && subnet == "" {
		placement = &ec2.Placement{AvailabilityZone: aws.String(availabilityZone)}
	}

	mockAWSClient.EXPECT().RunInstances(Placement(placement)).Return(
		&ec2.Reservation{
			Instances: []*ec2.Instance{
				{
					ImageId:    aws.String("ami-a9acbbd6"),
					InstanceId: aws.String("i-02fcb933c5da7085c"),
					State: &ec2.InstanceState{
						Name: aws.String("Running"),
						Code: aws.Int64(16),
					},
					LaunchTime: aws.Time(time.Now()),
				},
			},
		}, nil)
}

type placementMatcher struct {
	placement *ec2.Placement
}

func (m placementMatcher) Matches(input interface{}) bool {
	runInstancesInput, ok := input.(*ec2.RunInstancesInput)
	if !ok {
		return false
	}
	if runInstancesInput.Placement == m.placement {
		return true
	}
	return false
}

func (m placementMatcher) String() string {
	return fmt.Sprintf("is placement: %#v", m.placement)
}

func Placement(placement *ec2.Placement) gomock.Matcher { return placementMatcher{placement} }
