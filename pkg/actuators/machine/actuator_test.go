package machine

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"k8s.io/utils/pointer"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	machineapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
}

const (
	noError             = ""
	awsServiceError     = "error creating aws service"
	launchInstanceError = "error launching instance"
)

func TestMachineEvents(t *testing.T) {
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec.Value = nil

	workerMachine := machine.DeepCopy()
	workerMachine.Spec.Labels["node-role.kubernetes.io/worker"] = ""

	cases := []struct {
		name                    string
		machine                 *machinev1.Machine
		error                   string
		operation               func(actuator *Actuator, machine *machinev1.Machine)
		event                   string
		describeInstancesOutput *ec2.DescribeInstancesOutput
		describeInstancesErr    error
		runInstancesErr         error
		terminateInstancesErr   error
		lbErr                   error
		regInstancesWithLbErr   error
	}{
		{
			name:    "Create machine event failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Warning FailedCreate error decoding MachineProviderConfig: unable to find machine provider config: Spec.ProviderSpec.Value is not set",
		},
		{
			name:    "Create machine event failed (error creating aws service)",
			machine: machine,
			error:   awsServiceError,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Warning FailedCreate error creating aws service",
		},
		{
			name:            "Create machine event failed (error launching instance)",
			machine:         machine,
			runInstancesErr: fmt.Errorf("error"),
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Warning FailedCreate error creating EC2 instance: error",
		},
		{
			name:    "Create machine event failed (error updating load balancers)",
			machine: machine,
			lbErr:   fmt.Errorf("lb error"),
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Warning FailedCreate lb error",
		},
		{
			name:    "Create machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Normal Created Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Create worker machine event succeed",
			machine: workerMachine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.CreateMachine(machine)
			},
			event: "Normal Created Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Delete machine event failed",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.DeleteMachine(machine)
			},
			event: "Warning FailedDelete error decoding MachineProviderConfig: unable to find machine provider config: Spec.ProviderSpec.Value is not set",
		},
		{
			name:    "Delete machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.DeleteMachine(machine)
			},
			event: "Normal Deleted Deleted machine aws-actuator-testing-machine",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				Client: fake.NewFakeClient(tc.machine, awsCredentialsSecret, userDataSecret),
				AwsClientBuilder: func(client client.Client, secretName, namespace, region string) (awsclient.Client, error) {
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

			mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), tc.runInstancesErr).AnyTimes()
			if tc.describeInstancesOutput == nil {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), tc.describeInstancesErr).AnyTimes()
			} else {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(tc.describeInstancesOutput, tc.describeInstancesErr).AnyTimes()
			}

			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil)
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, tc.terminateInstancesErr).AnyTimes()
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr)
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()

			actuator, err := NewActuator(params)
			if err != nil {
				t.Fatalf("Could not create AWS machine actuator: %v", err)
			}

			tc.operation(actuator, tc.machine)
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

func TestActuator(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	getMachineStatus := func(objectClient client.Client, machine *machinev1.Machine) (*providerconfigv1.AWSMachineProviderStatus, error) {
		// Get updated machine object from the cluster client
		key := types.NamespacedName{
			Namespace: machine.Namespace,
			Name:      machine.Name,
		}
		updatedMachine := machinev1.Machine{}
		err := objectClient.Get(context.Background(), client.ObjectKey(key), &updatedMachine)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve machine: %v", err)
		}

		machineStatus := &providerconfigv1.AWSMachineProviderStatus{}
		if err := codec.DecodeProviderStatus(updatedMachine.Status.ProviderStatus, machineStatus); err != nil {
			return nil, fmt.Errorf("error decoding machine provider status: %v", err)
		}
		return machineStatus, nil
	}

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec.Value = nil

	machineNoClusterID := machine.DeepCopy()
	delete(machineNoClusterID.Labels, providerconfigv1.ClusterIDLabel)

	pendingInstance := stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c")
	pendingInstance.State = &ec2.InstanceState{
		Name: aws.String(ec2.InstanceStateNamePending),
	}

	cases := []struct {
		name                    string
		machine                 *machinev1.Machine
		error                   string
		operation               func(client client.Client, actuator *Actuator, machine *machinev1.Machine)
		describeInstancesOutput *ec2.DescribeInstancesOutput
		runInstancesErr         error
		describeInstancesErr    error
		terminateInstancesErr   error
		lbErr                   error
	}{
		{
			name:    "Create machine with success",
			machine: machine,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), machine)
				assert.NoError(t, createErr)

				machineStatus, err := getMachineStatus(objectClient, machine)
				if err != nil {
					t.Fatalf("Unable to get machine status: %v", err)
				}

				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationSucceeded)

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Update a machine
				if err := actuator.Update(context.TODO(), machine); err != nil {
					t.Errorf("Unable to create instance for machine: %v", err)
				}

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

				// Delete a machine
				if err := actuator.Delete(context.TODO(), machine); err != nil {
					t.Errorf("Unable to delete instance for machine: %v", err)
				}
			},
		},
		{
			name:            "Create machine with failure",
			machine:         machine,
			runInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), machine)
				assert.Error(t, createErr)

				machineStatus, err := getMachineStatus(objectClient, machine)
				if err != nil {
					t.Fatalf("Unable to get machine status: %v", err)
				}

				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationFailed)
			},
		},
		{
			name:    "Update machine with success",
			machine: machine,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:    "Update machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:  "Update machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:                 "Update machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name: "Update machine failed (no running instances)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name: "Update machine succeeds (two running instances)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c"),
							stubInstance("ami-a9acbbd7", "i-02fcb933c5da7085d"),
						},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name: "Update machine status fails (instance pending)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							pendingInstance,
						},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name: "Update machine failed (two running instances, error terminating one)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							stubInstance("ami-a9acbbd6", "i-02fcb933c5da7085c"),
							stubInstance("ami-a9acbbd7", "i-02fcb933c5da7085d"),
						},
					},
				},
			},
			terminateInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:    "Update machine with failure (cluster ID missing)",
			machine: machineNoClusterID,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:  "Update machine failed (error updating load balancers)",
			lbErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
		},
		{
			name:                 "Describe machine fails (error getting running instance)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Describe(machine)
			},
		},
		{
			name:    "Describe machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name:  "Describe machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name: "Describe machine fails (no running instance)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name: "Describe machine succeeds",
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name:    "Exists machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name: "Exists machine fails (no running instance)",
			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), machine)
			},
		},
		{
			name:    "Delete machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
		},
		{
			name:  "Delete machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
		},
		{
			name:                 "Delete machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
		},
		{
			name: "Delete machine failed (no running instances)",

			describeInstancesOutput: &ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{},
					},
				},
			},
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
		},
		{
			name: "Delete machine failed (error terminating instances)",

			terminateInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := fake.NewFakeClient(machine, awsCredentialsSecret, userDataSecret)
			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			params := ActuatorParams{
				Client: fakeClient,
				AwsClientBuilder: func(client client.Client, secretName, namespace, region string) (awsclient.Client, error) {
					if tc.error == awsServiceError {
						return nil, fmt.Errorf(awsServiceError)
					}
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

			mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), tc.runInstancesErr).AnyTimes()

			if tc.describeInstancesOutput == nil {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), tc.describeInstancesErr).AnyTimes()
			} else {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(tc.describeInstancesOutput, tc.describeInstancesErr).AnyTimes()
			}

			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, tc.terminateInstancesErr).AnyTimes()
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()

			if tc.machine == nil {
				tc.operation(fakeClient, actuator, machine)
			} else {
				tc.operation(fakeClient, actuator, tc.machine)
			}
		})
	}
}

func TestAvailabilityZone(t *testing.T) {
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

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			machine, err := stubMachine()
			if err != nil {
				t.Fatal(err)
			}

			machinePc := &providerconfigv1.AWSMachineProviderConfig{}
			if err = codec.DecodeProviderSpec(&machine.Spec.ProviderSpec, machinePc); err != nil {
				t.Fatal(err)
			}

			// no load balancers tested
			machinePc.LoadBalancers = nil

			machinePc.Placement.AvailabilityZone = tc.availabilityZone
			if tc.subnet == "" {
				machinePc.Subnet.ID = nil
			} else {
				machinePc.Subnet.ID = aws.String(tc.subnet)
			}

			config, err := codec.EncodeProviderSpec(machinePc)
			if err != nil {
				t.Fatal(err)
			}
			machine.Spec.ProviderSpec = *config

			fakeClient := fake.NewFakeClient(machine, awsCredentialsSecret, userDataSecret)

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			params := ActuatorParams{
				Client: fakeClient,
				AwsClientBuilder: func(client client.Client, secretName, namespace, region string) (awsclient.Client, error) {
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

			var placement *ec2.Placement
			if tc.availabilityZone != "" && tc.subnet == "" {
				placement = &ec2.Placement{AvailabilityZone: aws.String(tc.availabilityZone)}
			}
			az := "us-east-1a"
			mockAWSClient.EXPECT().RunInstances(placementMatcher{placement}).Return(
				&ec2.Reservation{
					Instances: []*ec2.Instance{
						{
							ImageId:    aws.String("ami-a9acbbd6"),
							InstanceId: aws.String("i-02fcb933c5da7085c"),
							State: &ec2.InstanceState{
								Name: aws.String(ec2.InstanceStateNameRunning),
							},
							LaunchTime: aws.Time(time.Now()),
							Placement: &ec2.Placement{
								AvailabilityZone: &az,
							},
						},
					},
				}, nil)

			mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(
				&ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								{
									ImageId:    aws.String("ami-a9acbbd6"),
									InstanceId: aws.String("i-02fcb933c5da7085c"),
									State: &ec2.InstanceState{
										Name: aws.String(ec2.InstanceStateNameRunning),
										Code: aws.Int64(16),
									},
									LaunchTime: aws.Time(time.Now()),
								},
							},
						},
					},
				}, nil).AnyTimes()

			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil)
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).AnyTimes()
			mockAWSClient.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().DescribeSubnets(gomock.Any()).Return(&ec2.DescribeSubnetsOutput{}, nil)

			actuator.Create(context.TODO(), machine)
		})
	}
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

func TestGetUserData(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}
	providerConfig := stubProviderConfig()
	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	testCases := []struct {
		secret *apiv1.Secret
		error  error
	}{
		{
			secret: &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespace,
				},
				Data: map[string][]byte{
					userDataSecretKey: []byte(userDataBlob),
				},
			},
			error: nil,
		},
		{
			secret: &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "notFound",
					Namespace: defaultNamespace,
				},
				Data: map[string][]byte{
					userDataSecretKey: []byte(userDataBlob),
				},
			},
			error: &machineapierrors.MachineError{},
		},
		{
			secret: &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespace,
				},
				Data: map[string][]byte{
					"badKey": []byte(userDataBlob),
				},
			},
			error: &machineapierrors.MachineError{},
		},
	}

	for _, tc := range testCases {
		params := ActuatorParams{
			Client: fake.NewFakeClient(tc.secret),
			Codec:  codec,
			EventRecorder: &record.FakeRecorder{
				Events: make(chan string, 1),
			},
		}
		actuator, err := NewActuator(params)
		if err != nil {
			t.Fatalf("Could not create AWS machine actuator: %v", err)
		}
		userData, err := actuator.getUserData(machine, providerConfig)
		if tc.error != nil {
			if err == nil {
				t.Fatal("Expected error")
			}
			_, expectMachineError := tc.error.(*machineapierrors.MachineError)
			_, gotMachineError := err.(*machineapierrors.MachineError)
			if expectMachineError && !gotMachineError || !expectMachineError && gotMachineError {
				t.Errorf("Expected %T, got: %T", tc.error, err)
			}
		} else {
			if compare := bytes.Compare(userData, []byte(userDataBlob)); compare != 0 {
				t.Errorf("Expected: %v, got: %v", []byte(userDataBlob), userData)
			}
		}
	}
}

func TestCreate(t *testing.T) {
	// mock aws API calls
	mockCtrl := gomock.NewController(t)
	mockAWSClient := mockaws.NewMockClient(mockCtrl)
	mockAWSClient.EXPECT().DescribeSecurityGroups(gomock.Any()).Return(nil, fmt.Errorf("describeSecurityGroups error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeAvailabilityZones(gomock.Any()).Return(nil, fmt.Errorf("describeAvailabilityZones error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeImages(gomock.Any()).Return(nil, fmt.Errorf("describeImages error")).AnyTimes()
	mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), nil).AnyTimes()
	mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil).AnyTimes()
	mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), nil).AnyTimes()
	mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
	mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), nil)
	mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
	mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()

	testCases := []struct {
		testcase             string
		providerConfig       *providerconfigv1.AWSMachineProviderConfig
		userDataSecret       *corev1.Secret
		awsCredentialsSecret *corev1.Secret
		error                error
	}{
		{
			testcase: "Create succeed",
			providerConfig: &providerconfigv1.AWSMachineProviderConfig{
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
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
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
					{
						Name: "cluster-net-lb",
						Type: providerconfigv1.NetworkLoadBalancerType,
					},
				},
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			error:                nil,
		},
		{
			testcase: "Bad userData",
			providerConfig: &providerconfigv1.AWSMachineProviderConfig{
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
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
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
					{
						Name: "cluster-net-lb",
						Type: providerconfigv1.NetworkLoadBalancerType,
					},
				},
			},
			userDataSecret: &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      userDataSecretName,
					Namespace: defaultNamespace,
				},
				Data: map[string][]byte{
					"badKey": []byte(userDataBlob),
				},
			},
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			error:                &machineapierrors.MachineError{},
		},
		{
			testcase: "Failed security groups return invalid config machine error",
			providerConfig: &providerconfigv1.AWSMachineProviderConfig{
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
				SecurityGroups: []providerconfigv1.AWSResourceReference{{
					Filters: []providerconfigv1.Filter{{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-%s-sg", clusterID, "role")},
					}},
				}},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			error:                &machineapierrors.MachineError{},
		},
		{
			testcase: "Failed Availability zones return invalid config machine error",
			providerConfig: &providerconfigv1.AWSMachineProviderConfig{
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
					Filters: []providerconfigv1.Filter{{
						Name:   "tag:Name",
						Values: []string{fmt.Sprintf("%s-private-%s", clusterID, "az")},
					}},
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
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			error:                &machineapierrors.MachineError{},
		},
		{
			testcase: "Failed BlockDevices return invalid config machine error",
			providerConfig: &providerconfigv1.AWSMachineProviderConfig{
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
				BlockDevices: []providerconfigv1.BlockDeviceMappingSpec{
					{
						EBS: &providerconfigv1.EBSBlockDeviceSpec{
							VolumeType: pointer.StringPtr("type"),
							VolumeSize: pointer.Int64Ptr(int64(1)),
							Iops:       pointer.Int64Ptr(int64(1)),
						},
					},
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
					{ID: aws.String("sg-00868b02fbe29de17")},
					{ID: aws.String("sg-0a4658991dc5eb40a")},
					{ID: aws.String("sg-009a70e28fa4ba84e")},
					{ID: aws.String("sg-07323d56fb932c84c")},
					{ID: aws.String("sg-08b1ffd32874d59a2")},
				},
				PublicIP: aws.Bool(true),
			},
			userDataSecret:       stubUserDataSecret(),
			awsCredentialsSecret: stubAwsCredentialsSecret(),
			error:                &machineapierrors.MachineError{},
		},
	}

	for _, tc := range testCases {
		// create fake resources
		t.Logf("testCase: %v", tc.testcase)
		codec, err := providerconfigv1.NewCodec()
		if err != nil {
			t.Fatalf("unable to build codec: %v", err)
		}
		encodedProviderConfig, err := codec.EncodeProviderSpec(tc.providerConfig)
		if err != nil {
			t.Fatalf("Unexpected error")
		}
		machine, err := stubMachine()
		if err != nil {
			t.Fatal(err)
		}
		machine.Spec.ProviderSpec = *encodedProviderConfig
		fakeClient := fake.NewFakeClientWithScheme(scheme.Scheme, machine, tc.awsCredentialsSecret, tc.userDataSecret)

		// create actuator
		params := ActuatorParams{
			Client: fakeClient,
			AwsClientBuilder: func(client client.Client, secretName, namespace, region string) (awsclient.Client, error) {
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

		// test create
		err = actuator.Create(context.TODO(), machine)
		if tc.error != nil {
			if err == nil {
				t.Fatalf("Expected error")
			}
			_, expectMachineError := tc.error.(*machineapierrors.MachineError)
			_, gotMachineError := err.(*machineapierrors.MachineError)
			if expectMachineError && !gotMachineError || !expectMachineError && gotMachineError {
				t.Fatalf("Expected %T, got: %T", tc.error, err)
			}
		} else if err != nil {
			t.Fatalf("Unexpected error")
		}
	}
}

func TestGetMachineInstances(t *testing.T) {
	clusterID := "aws-actuator-cluster"
	instanceID := "i-02fa4197109214b46"
	imageID := "ami-a9acbbd6"

	machine, err := stubMachine()
	if err != nil {
		t.Fatalf("unable to build stub machine: %v", err)
	}

	codec, err := providerconfigv1.NewCodec()
	if err != nil {
		t.Fatalf("unable to build codec: %v", err)
	}

	testCases := []struct {
		testcase       string
		providerStatus providerconfigv1.AWSMachineProviderStatus
		awsClientFunc  func(*gomock.Controller) awsclient.Client
		exists         bool
	}{
		{
			testcase:       "empty-status-search-by-tag",
			providerStatus: providerconfigv1.AWSMachineProviderStatus{},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				request := &ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						{
							Name:   awsTagFilter("Name"),
							Values: aws.StringSlice([]string{machine.Name}),
						},

						clusterFilter(clusterID),
					},
				}

				mockAWSClient.EXPECT().DescribeInstances(request).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameRunning),
					nil,
				).Times(1)

				return mockAWSClient
			},
			exists: true,
		},
		{
			testcase: "has-status-search-by-id-running",
			providerStatus: providerconfigv1.AWSMachineProviderStatus{
				InstanceID: aws.String(instanceID),
			},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				request := &ec2.DescribeInstancesInput{
					InstanceIds: aws.StringSlice([]string{instanceID}),
				}

				mockAWSClient.EXPECT().DescribeInstances(request).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameRunning),
					nil,
				).Times(1)

				return mockAWSClient
			},
			exists: true,
		},
		{
			testcase: "has-status-search-by-id-terminated",
			providerStatus: providerconfigv1.AWSMachineProviderStatus{
				InstanceID: aws.String(instanceID),
			},
			awsClientFunc: func(ctrl *gomock.Controller) awsclient.Client {
				mockAWSClient := mockaws.NewMockClient(ctrl)

				first := mockAWSClient.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
					InstanceIds: aws.StringSlice([]string{instanceID}),
				}).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameTerminated),
					nil,
				).Times(1)

				mockAWSClient.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
					Filters: []*ec2.Filter{
						{
							Name:   awsTagFilter("Name"),
							Values: aws.StringSlice([]string{machine.Name}),
						},

						clusterFilter(clusterID),
					},
				}).Return(
					stubDescribeInstancesOutput(imageID, instanceID, ec2.InstanceStateNameTerminated),
					nil,
				).Times(1).After(first)

				return mockAWSClient
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testcase, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			awsStatusRaw, err := codec.EncodeProviderStatus(&tc.providerStatus)
			if err != nil {
				t.Errorf("Error encoding ProviderStatus: %v", err)
			}

			machineCopy := machine.DeepCopy()
			machineCopy.Status.ProviderStatus = awsStatusRaw

			awsClient := tc.awsClientFunc(ctrl)

			params := ActuatorParams{
				Codec:            codec,
				AwsClientBuilder: awsClientBuilderFunc(awsClient),
			}

			actuator, err := NewActuator(params)
			if err != nil {
				t.Errorf("Error creating Actuator: %v", err)
			}

			instances, err := actuator.getMachineInstances(machineCopy)
			if err != nil {
				t.Errorf("Unexpected error from getMachineInstances: %v", err)
			}
			if tc.exists != (len(instances) > 0) {
				t.Errorf("Expected instance exists: %t, got instances: %v", tc.exists, instances)
			}
		})
	}
}

func awsClientBuilderFunc(c awsclient.Client) awsclient.AwsClientBuilderFuncType {
	return func(_ client.Client, _, _, _ string) (awsclient.Client, error) {
		return c, nil
	}
}
