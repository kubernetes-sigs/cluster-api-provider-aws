package machine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client/mock"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	clusterv1.AddToScheme(scheme.Scheme)
}

func TestMachineEvents(t *testing.T) {
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	machine, cluster, awsCredentialsSecret, userDataSecret, err := stubMachineAPIResources()
	if err != nil {
		t.Fatal(err)
	}

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderConfig.Value = nil
	machineInvalidProviderConfig.Spec.ProviderConfig.ValueFrom = nil

	const (
		noError            = ""
		awsServiceError    = "error creating aws service"
		launcInstanceError = "error launching instance"
		lbError            = "error updating load balancers"
	)

	cases := []struct {
		name      string
		machine   *clusterv1.Machine
		error     string
		operation func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine)
		event     string
	}{
		{
			name:    "Create machine event failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate InvalidConfiguration",
		},
		{
			name:    "Create machine event failed (error creating aws service)",
			machine: machine,
			error:   awsServiceError,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event failed (error launching instance)",
			machine: machine,
			error:   launcInstanceError,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event failed (error updating load balancers)",
			machine: machine,
			error:   lbError,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Normal Created Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Delete machine event failed",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.DeleteMachine(cluster, machine)
			},
			event: "Warning FailedDelete InvalidConfiguration",
		},
		{
			name:    "Delete machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *clusterv1.Cluster, machine *clusterv1.Machine) {
				actuator.DeleteMachine(cluster, machine)
			},
			event: "Normal Deleted Deleted Machine aws-actuator-testing-machine",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				Client:     fake.NewFakeClient(tc.machine),
				KubeClient: kubernetesfake.NewSimpleClientset(awsCredentialsSecret, userDataSecret),
				AwsClientBuilder: func(kubeClient kubernetes.Interface, secretName, namespace, region string) (awsclient.Client, error) {
					if tc.error == awsServiceError {
						return nil, fmt.Errorf(awsServiceError)
					}
					return mockAWSClient, nil
				},
				Codec: codec,
				// use fake recorder and store an event into one item long buffer for subsequent check
				EventRecorder: &record.FakeRecorder{
					Events: eventsChannel,
				},
			}

			mockRunInstances(mockAWSClient, tc.error == launcInstanceError)
			mockDescribeInstances(mockAWSClient, false)
			mockTerminateInstances(mockAWSClient)
			mockRegisterInstancesWithLoadBalancer(mockAWSClient, tc.error == lbError)

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create AWS machine actuator: %v", err)
			}

			tc.operation(actuator, cluster, tc.machine)
			select {
			case event := <-eventsChannel:
				if event != tc.event {
					t.Errorf("Expected %q event, got %q", tc.event, event)
				}
			default:
				t.Errorf("Expected %q event, got none", tc.event)
			}
		})
	}
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
			machine, cluster, awsCredentialsSecret, userDataSecret, err := stubMachineAPIResources()
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
				// use empty recorder dropping any event recorded
				EventRecorder: &record.FakeRecorder{},
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
			createErr := actuator.Create(context.TODO(), cluster, machine)

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

			machineStatus := &providerconfigv1.AWSMachineProviderStatus{}
			if err := codec.DecodeProviderStatus(updatedMachine.Status.ProviderStatus, machineStatus); err != nil {
				t.Fatalf("error decoding machine provider status: %v", err)
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
				if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// TODO(jchaloup): Wait until the machine is ready

				// Update a machine
				if err := actuator.Update(context.TODO(), cluster, machine); err != nil {
					t.Errorf("Unable to create instance for machine: %v", err)
				}

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Delete a machine
				if err := actuator.Delete(context.TODO(), cluster, machine); err != nil {
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
		mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
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
			machine, cluster, awsCredentialsSecret, userDataSecret, err := stubMachineAPIResources()
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
				// use empty recorder dropping any event recorded
				EventRecorder: &record.FakeRecorder{},
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

			actuator.Create(context.TODO(), cluster, machine)
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
