package machine

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/record"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

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

	cluster := stubCluster()
	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec.Value = nil
	machineInvalidProviderConfig.Spec.ProviderSpec.ValueFrom = nil

	workerMachine := machine.DeepCopy()
	workerMachine.Labels[providerconfigv1.MachineTypeLabel] = "worker"

	cases := []struct {
		name                    string
		machine                 *machinev1.Machine
		error                   string
		operation               func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine)
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
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate InvalidConfiguration",
		},
		{
			name:    "Create machine event failed (error creating aws service)",
			machine: machine,
			error:   awsServiceError,
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:            "Create machine event failed (error launching instance)",
			machine:         machine,
			runInstancesErr: fmt.Errorf("error"),
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event failed (error updating load balancers)",
			machine: machine,
			lbErr:   fmt.Errorf("error"),
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Warning FailedCreate CreateError",
		},
		{
			name:    "Create machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Normal Created Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Create worker machine event succeed",
			machine: workerMachine,
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.CreateMachine(cluster, machine)
			},
			event: "Normal Created Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Delete machine event failed",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.DeleteMachine(cluster, machine)
			},
			event: "Warning FailedDelete InvalidConfiguration",
		},
		{
			name:    "Delete machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.DeleteMachine(cluster, machine)
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
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c"), tc.describeInstancesErr).AnyTimes()
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

func TestActuator(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	cluster := stubCluster()
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
	machineInvalidProviderConfig.Spec.ProviderSpec.ValueFrom = nil

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
		operation               func(client client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine)
		describeInstancesOutput *ec2.DescribeInstancesOutput
		runInstancesErr         error
		describeInstancesErr    error
		terminateInstancesErr   error
		lbErr                   error
	}{
		{
			name:    "Create machine with success",
			machine: machine,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), cluster, machine)
				assert.NoError(t, createErr)

				machineStatus, err := getMachineStatus(objectClient, machine)
				if err != nil {
					t.Fatalf("Unable to get machine status: %v", err)
				}

				assert.Equal(t, machineStatus.Conditions[0].Reason, MachineCreationSucceeded)

				// Get the machine
				if exists, err := actuator.Exists(context.TODO(), cluster, machine); err != nil || !exists {
					t.Errorf("Instance for %v does not exists: %v", strings.Join([]string{machine.Namespace, machine.Name}, "/"), err)
				} else {
					t.Logf("Instance for %v exists", strings.Join([]string{machine.Namespace, machine.Name}, "/"))
				}

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
			},
		},
		{
			name:            "Create machine with failure",
			machine:         machine,
			runInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				createErr := actuator.Create(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Update machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:  "Update machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Update machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Update machine with failure (cluster ID missing)",
			machine: machineNoClusterID,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:  "Update machine failed (error updating load balancers)",
			lbErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Describe machine fails (error getting running instance)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:    "Describe machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:  "Describe machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name: "Describe machine succeeds",
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Describe(cluster, machine)
			},
		},
		{
			name:    "Exists machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Exists(context.TODO(), cluster, machine)
			},
		},
		{
			name:    "Delete machine failed (invalid configuration)",
			machine: machineInvalidProviderConfig,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name:  "Delete machine failed (error creating aws service)",
			error: awsServiceError,
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name:                 "Delete machine failed (error getting running instances)",
			describeInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
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
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
			},
		},
		{
			name: "Delete machine failed (error terminating instances)",
			terminateInstancesErr: fmt.Errorf("error"),
			operation: func(objectClient client.Client, actuator *Actuator, cluster *machinev1.Cluster, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), cluster, machine)
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
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c"), tc.describeInstancesErr).AnyTimes()
			} else {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(tc.describeInstancesOutput, tc.describeInstancesErr).AnyTimes()
			}

			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, tc.terminateInstancesErr).AnyTimes()
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()

			if tc.machine == nil {
				tc.operation(fakeClient, actuator, cluster, machine)
			} else {
				tc.operation(fakeClient, actuator, cluster, tc.machine)
			}
		})
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

	cluster := stubCluster()
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

			mockAWSClient.EXPECT().RunInstances(placementMatcher{placement}).Return(
				&ec2.Reservation{
					Instances: []*ec2.Instance{
						{
							ImageId:    aws.String("ami-a9acbbd6"),
							InstanceId: aws.String("i-02fcb933c5da7085c"),
							State: &ec2.InstanceState{
								Name: aws.String("Running"),
							},
							LaunchTime: aws.Time(time.Now()),
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
										Name: aws.String("Running"),
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

			actuator.Create(context.TODO(), cluster, machine)
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
