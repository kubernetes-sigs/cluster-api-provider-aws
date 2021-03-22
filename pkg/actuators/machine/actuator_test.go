package machine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	mockaws "sigs.k8s.io/cluster-api-provider-aws/pkg/client/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	// Add types to scheme
	machinev1.AddToScheme(scheme.Scheme)
	configv1.AddToScheme(scheme.Scheme)
}

func TestMachineEvents(t *testing.T) {
	g := NewWithT(t)

	awsCredentialsSecret := stubAwsCredentialsSecret()
	g.Expect(k8sClient.Create(context.TODO(), awsCredentialsSecret)).To(Succeed())
	defer func() {
		g.Expect(k8sClient.Delete(context.TODO(), awsCredentialsSecret)).To(Succeed())
	}()

	userDataSecret := stubUserDataSecret()
	g.Expect(k8sClient.Create(context.TODO(), userDataSecret)).To(Succeed())
	defer func() {
		g.Expect(k8sClient.Delete(context.TODO(), userDataSecret)).To(Succeed())
	}()

	cases := []struct {
		name                string
		error               string
		operation           func(actuator *Actuator, machine *machinev1.Machine)
		event               string
		awsError            bool
		invalidMachineScope bool
	}{
		{
			name: "Create machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: failed to create scope for machine: failed to create aws client: AWS client error",
			invalidMachineScope: true,
			awsError:            false,
		},
		{
			name: "Create machine event failed, reconciler's create failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: reconciler failed to Create machine: unable to remove stopped machines: error getting stopped instances: AWS error",
			invalidMachineScope: false,
			awsError:            true,
		},
		{
			name: "Create machine event succeed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Create(context.TODO(), machine)
			},
			event:               "Created Machine aws-actuator-testing-machine",
			invalidMachineScope: false,
			awsError:            false,
		},
		{
			name: "Update machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: failed to create scope for machine: failed to create aws client: AWS client error",
			invalidMachineScope: true,
			awsError:            false,
		},
		{
			name: "Update machine event failed, reconciler's update failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: reconciler failed to Update machine: AWS error",
			invalidMachineScope: false,
			awsError:            true,
		},
		{
			name: "Update machine event succeed and only one event is created",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Update(context.TODO(), machine)
				actuator.Update(context.TODO(), machine)
			},
			event: "Updated Machine aws-actuator-testing-machine",
		},
		{
			name: "Delete machine event failed on invalid machine scope",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: failed to create scope for machine: failed to create aws client: AWS client error",
			invalidMachineScope: true,
			awsError:            false,
		},
		{
			name: "Delete machine event failed, reconciler's delete failed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event:               "aws-actuator-testing-machine: reconciler failed to Delete machine: AWS error",
			invalidMachineScope: false,
			awsError:            true,
		},
		{
			name: "Delete machine event succeed",
			operation: func(actuator *Actuator, machine *machinev1.Machine) {
				actuator.Delete(context.TODO(), machine)
			},
			event:               "Deleted machine aws-actuator-testing-machine",
			invalidMachineScope: false,
			awsError:            false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			gs := NewWithT(t)

			machine, err := stubMachine()
			gs.Expect(err).ToNot(HaveOccurred())
			gs.Expect(stubMachine).ToNot(BeNil())

			// Create the machine
			gs.Expect(k8sClient.Create(ctx, machine)).To(Succeed())
			defer func() {
				gs.Expect(k8sClient.Delete(ctx, machine)).To(Succeed())
			}()

			// Create infrastructure object
			infra := &configv1.Infrastructure{ObjectMeta: metav1.ObjectMeta{Name: awsclient.GlobalInfrastuctureName}}
			gs.Expect(k8sClient.Create(ctx, infra)).To(Succeed())
			defer func() {
				gs.Expect(k8sClient.Delete(ctx, infra)).To(Succeed())
			}()

			// Ensure the machine has synced to the cache
			getMachine := func() error {
				machineKey := types.NamespacedName{Namespace: machine.Namespace, Name: machine.Name}
				return k8sClient.Get(ctx, machineKey, machine)
			}
			gs.Eventually(getMachine, timeout).Should(Succeed())

			mockCtrl := gomock.NewController(t)
			mockAWSClient := mockaws.NewMockClient(mockCtrl)
			awsClientBuilder := func(client runtimeclient.Client, secretName, namespace, region string, configManagedClient runtimeclient.Client) (awsclient.Client, error) {
				return mockAWSClient, nil
			}
			if tc.invalidMachineScope {
				awsClientBuilder = func(client runtimeclient.Client, secretName, namespace, region string, configManagedClient runtimeclient.Client) (awsclient.Client, error) {
					return nil, errors.New("AWS client error")
				}
			}

			if tc.awsError {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(nil, errors.New("AWS error")).AnyTimes()
			} else {
				mockAWSClient.EXPECT().DescribeInstances(gomock.Any()).Return(stubDescribeInstancesOutput("ami-a9acbbd6", "i-02fcb933c5da7085c", ec2.InstanceStateNameRunning), nil).AnyTimes()
			}

			mockAWSClient.EXPECT().RunInstances(gomock.Any()).Return(stubReservation("ami-a9acbbd6", "i-02fcb933c5da7085c"), nil).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil)
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().TerminateInstances(gomock.Any()).Return(&ec2.TerminateInstancesOutput{}, nil).AnyTimes()
			mockAWSClient.EXPECT().RegisterInstancesWithLoadBalancer(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeLoadBalancers(gomock.Any()).Return(stubDescribeLoadBalancersOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2DescribeTargetGroups(gomock.Any()).Return(stubDescribeTargetGroupsOutput(), nil).AnyTimes()
			mockAWSClient.EXPECT().ELBv2RegisterTargets(gomock.Any()).Return(nil, nil).AnyTimes()
			mockAWSClient.EXPECT().DescribeVpcs(gomock.Any()).Return(StubDescribeVPCs()).AnyTimes()
			mockAWSClient.EXPECT().DescribeDHCPOptions(gomock.Any()).Return(StubDescribeDHCPOptions()).AnyTimes()
			mockAWSClient.EXPECT().CreateTags(gomock.Any()).Return(&ec2.CreateTagsOutput{}, nil).AnyTimes()

			params := ActuatorParams{
				Client:           k8sClient,
				EventRecorder:    eventRecorder,
				AwsClientBuilder: awsClientBuilder,
			}
			actuator := NewActuator(params)
			tc.operation(actuator, machine)

			eventList := &v1.EventList{}
			waitForEvent := func() error {
				gs.Expect(k8sClient.List(ctx, eventList, client.InNamespace(machine.Namespace))).To(Succeed())
				if len(eventList.Items) != 1 {
					errorMsg := fmt.Sprintf("Expected len 1, got %d", len(eventList.Items))
					return errors.New(errorMsg)
				}
				return nil
			}

			gs.Eventually(waitForEvent, timeout).Should(Succeed())

			gs.Expect(eventList.Items[0].Message).To(Equal(tc.event))

			for i := range eventList.Items {
				gs.Expect(k8sClient.Delete(ctx, &eventList.Items[i])).To(Succeed())
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
