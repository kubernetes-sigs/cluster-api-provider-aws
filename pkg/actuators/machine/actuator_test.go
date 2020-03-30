package machine

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
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
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	awsCredentialsSecret := stubAwsCredentialsSecret()
	userDataSecret := stubUserDataSecret()

	machineInvalidProviderConfig := machine.DeepCopy()
	machineInvalidProviderConfig.Spec.ProviderSpec = machinev1.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw: []byte("test"),
		},
	}

	workerMachine := machine.DeepCopy()
	workerMachine.Spec.Labels["node-role.kubernetes.io/worker"] = ""

	cases := []struct {
		name                  string
		machine               *machinev1.Machine
		error                 string
		operation             func(actuator *Actuator, machine *machinev1.Machine)
		event                 string
		awsError              bool
		runInstancesErr       error
		terminateInstancesErr error
		lbErr                 error
		regInstancesWithLbErr error
	}{
		{
			name:    "Create machine event failed on invalid machine scope",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event: "Warning FailedCreate failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.AWSMachineProviderConfig",
		},
		{
			name:    "Create machine event failed, reconciler's create failed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event:    "Warning FailedCreate unable to remove stopped machines: error getting stopped instances: error",
			awsError: true,
		},
		{
			name:    "Create machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event: "Normal Create Created Machine aws-actuator-testing-machine",
		},
		{
			name:    "Update machine event failed on invalid machine scope",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
			event: "Warning FailedUpdate aws-actuator-testing-machine: failed to create scope for machine: failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.AWSMachineProviderConfig",
		},
		{
			name:    "Update machine event failed, reconciler's update failed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
			event:    "Warning FailedCreate error",
			awsError: true,
		},
		{
			name:    "Update machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
			event: "Normal Update Updated Machine aws-actuator-testing-machine",
		},
		{
			name:    "Delete machine event failed on invalid machine scope",
			machine: machineInvalidProviderConfig,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event: "Warning FailedDelete aws-actuator-testing-machine: failed to create scope for machine: failed to get machine config: error unmarshalling providerSpec: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.AWSMachineProviderConfig",
		},
		{
			name:    "Delete machine event failed, reconciler's delete failed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event:    "Warning FailedDelete error",
			awsError: true,
		},
		{
			name:    "Delete machine event succeed",
			machine: machine,
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event: "Normal Delete Deleted machine aws-actuator-testing-machine",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)

			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				Client: fake.NewFakeClient(tc.machine, awsCredentialsSecret, userDataSecret),
				// use fake recorder and store an event into one item long buffer for subsequent check
				EventRecorder: &record.FakeRecorder{
					Events: eventsChannel,
				},
				AwsClientBuilder: func(client runtimeclient.Client, secretName, namespace, region string) (awsclient.Client, error) {
					return mockAWSClient, nil
				},
			}

			if tc.awsError {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
			} else {

				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), nil).AnyTimes()
			}

			mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), tc.runInstancesErr).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil)
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, tc.terminateInstancesErr).AnyTimes()
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, tc.lbErr).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), tc.lbErr)
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()

			actuator := NewActuator(params)

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

func TestHandleMachineErrors(t *testing.T) {
	machine, err := stubMachine()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name        string
		eventAction string
		event       string
	}{
		{
			name:        "Create event when event action is present",
			eventAction: "testAction",
			event:       "Warning FailedtestAction testError",
		},
		{
			name:        "Don't event when there is no event action",
			eventAction: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eventsChannel := make(chan string, 1)

			params := ActuatorParams{
				// use fake recorder and store an event into one item long buffer for subsequent check
				EventRecorder: &record.FakeRecorder{
					Events: eventsChannel,
				},
			}

			actuator := NewActuator(params)

			actuator.handleMachineError(machine, errors.New("testError"), tc.eventAction)

			select {
			case event := <-eventsChannel:
				if event != tc.event {
					t.Errorf("Expected %q event, got %q", tc.event, event)
				}
			default:
				if tc.event != "" {
					t.Errorf("Expected %q event, got none", tc.event)
				}
			}
		})
	}
}
