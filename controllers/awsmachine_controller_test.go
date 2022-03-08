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
	"bytes"
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/mock_services"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const providerID = "aws:////myMachine"

func TestAWSMachineReconciler(t *testing.T) {
	var (
		reconciler AWSMachineReconciler
		cs         *scope.ClusterScope
		ms         *scope.MachineScope
		mockCtrl   *gomock.Controller
		ec2Svc     *mock_services.MockEC2Interface
		elbSvc     *mock_services.MockELBInterface
		secretSvc  *mock_services.MockSecretInterface
		recorder   *record.FakeRecorder
	)

	setup := func(t *testing.T, g *WithT, awsMachine *infrav1.AWSMachine) {
		// https://github.com/kubernetes/klog/issues/87#issuecomment-540153080
		// TODO: Replace with LogToOutput when https://github.com/kubernetes/klog/pull/99 merges
		t.Helper()

		var err error

		if err := flag.Set("logtostderr", "false"); err != nil {
			_ = fmt.Errorf("Error setting logtostderr flag")
		}
		if err := flag.Set("v", "2"); err != nil {
			_ = fmt.Errorf("Error setting v flag")
		}
		klog.SetOutput(GinkgoWriter)

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap-data",
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}

		client := fake.NewClientBuilder().WithObjects(awsMachine, secret).Build()
		ms, err = scope.NewMachineScope(
			scope.MachineScopeParams{
				Client: client,
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
							DataSecretName: pointer.StringPtr("bootstrap-data"),
						},
					},
				},
				InfraCluster: cs,
				AWSMachine:   awsMachine,
			},
		)
		g.Expect(err).To(BeNil())

		cs, err = scope.NewClusterScope(
			scope.ClusterScopeParams{
				Client:     fake.NewClientBuilder().WithObjects(awsMachine, secret).Build(),
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}},
			},
		)
		g.Expect(err).To(BeNil())

		ms, err = scope.NewMachineScope(
			scope.MachineScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					Status: clusterv1.ClusterStatus{
						InfrastructureReady: true,
					},
				},
				Machine: &clusterv1.Machine{
					Spec: clusterv1.MachineSpec{
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: pointer.StringPtr("bootstrap-data"),
						},
					},
				},
				InfraCluster: cs,
				AWSMachine:   awsMachine,
			},
		)
		g.Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(t)
		ec2Svc = mock_services.NewMockEC2Interface(mockCtrl)
		secretSvc = mock_services.NewMockSecretInterface(mockCtrl)
		elbSvc = mock_services.NewMockELBInterface(mockCtrl)

		// If your test hangs for 9 minutes, increase the value here to the number of events during a reconciliation loop
		recorder = record.NewFakeRecorder(2)

		reconciler = AWSMachineReconciler{
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2Interface {
				return ec2Svc
			},
			secretsManagerServiceFactory: func(cloud.ClusterScoper) services.SecretInterface {
				return secretSvc
			},
			Recorder: recorder,
			Log:      klogr.New(),
		}
	}
	teardown := func(t *testing.T, g *WithT) {
		t.Helper()
		mockCtrl.Finish()
	}

	t.Run("Reconciling an AWSMachine", func(t *testing.T) {
		t.Run("when can't reach amazon", func(t *testing.T) {
			expectedErr := errors.New("no connection available ")
			runningInstance := func(t *testing.T, g *WithT) {
				t.Helper()
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, expectedErr).AnyTimes()
			}

			t.Run("should exit immediately on an error state", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				er := capierrors.CreateMachineError
				ms.AWSMachine.Status.FailureReason = &er
				ms.AWSMachine.Status.FailureMessage = pointer.StringPtr("Couldn't create machine")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
			})

			t.Run("should exit immediately if cluster infra isn't ready", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				ms.Cluster.Status.InfrastructureReady = false

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(buf.String()).To(ContainSubstring("Cluster infrastructure is not ready yet"))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForClusterInfrastructureReason}})
			})

			t.Run("should exit immediately if bootstrap data secret reference isn't available", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				ms.Machine.Spec.Bootstrap.DataSecretName = nil

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)

				g.Expect(err).To(BeNil())
				g.Expect(buf.String()).To(ContainSubstring("Bootstrap data secret reference is not yet available"))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForBootstrapDataReason}})
			})

			t.Run("should return an error when we can't list instances by tags", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			t.Run("shouldn't add our finalizer to the machine", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)

				g.Expect(len(ms.AWSMachine.Finalizers)).To(Equal(0))
			})
		})

		t.Run("when there's a provider ID", func(t *testing.T) {
			id := providerID
			providerID := func(t *testing.T, g *WithT) {
				t.Helper()
				_, err := noderefutil.NewProviderID(id)
				g.Expect(err).To(BeNil())

				ms.AWSMachine.Spec.ProviderID = &id
			}

			t.Run("should look up by provider ID when one exists", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)

				providerID(t, g)
				expectedErr := errors.New("no connection available ")
				ec2Svc.EXPECT().InstanceIfExists(PointsTo("myMachine")).Return(nil, expectedErr)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			t.Run("should try to create a new machine if none exists and add finalizers", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)

				providerID(t, g)
				expectedErr := errors.New("Invalid instance")
				ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(nil, expectedErr)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})

		t.Run("should fail to find instance if no provider ID provided", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			id := "aws////myMachine"

			ms.AWSMachine.Spec.ProviderID = &id
			expectedErr := "providerID must be of the form <cloudProvider>://<optional>/<segments>/<provider id>"

			_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			g.Expect(err.Error()).To(ContainSubstring(expectedErr))
		})

		t.Run("when instance creation succeeds", func(t *testing.T) {
			var instance *infrav1.Instance

			instanceCreate := func(t *testing.T, g *WithT) {
				t.Helper()

				instance = &infrav1.Instance{
					ID:        "myMachine",
					VolumeIDs: []string{"volume-1", "volume-2"},
				}
				instance.State = infrav1.InstanceStatePending

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
			}

			t.Run("instance security group errors", func(t *testing.T) {
				var buf *bytes.Buffer
				getInstanceSecurityGroups := func(t *testing.T, g *WithT) {
					t.Helper()

					buf = new(bytes.Buffer)
					klog.SetOutput(buf)
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(nil, errors.New("stop here"))
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				}

				t.Run("should set attributes after creating an instance", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getInstanceSecurityGroups(t, g)

					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Spec.ProviderID).To(PointTo(Equal(providerID)))
				})

				t.Run("should set instance to pending", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getInstanceSecurityGroups(t, g)

					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					instance.State = infrav1.InstanceStatePending
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStatePending)))
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))

					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				})

				t.Run("should set instance to running", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)

					instanceCreate(t, g)
					getInstanceSecurityGroups(t, g)

					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					instance.State = infrav1.InstanceStateRunning
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(true))
					g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{
						{conditionType: infrav1.InstanceReadyCondition, status: corev1.ConditionTrue},
					})
				})
			})
			t.Run("new EC2 instance state: should error when the instance state is a new unseen one", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)
				instance.State = "NewAWSMachineState"
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
				g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state is undefined")))
				g.Eventually(recorder.Events).Should(Receive(ContainSubstring("InstanceUnhandledState")))
				g.Expect(ms.AWSMachine.Status.FailureMessage).To(PointTo(Equal("EC2 instance state \"NewAWSMachineState\" is undefined")))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{conditionType: infrav1.InstanceReadyCondition, status: corev1.ConditionUnknown}})
			})
			t.Run("security Groups succeed", func(t *testing.T) {
				getCoreSecurityGroups := func(t *testing.T, g *WithT) {
					t.Helper()

					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
						Return(map[string][]string{"eid": {}}, nil)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				}
				t.Run("should reconcile security groups", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1.AWSResourceReference{
						{
							ID: pointer.StringPtr("sg-2345"),
						},
					}
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(instance.ID, []string{"sg-2345"})

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{conditionType: infrav1.SecurityGroupsReadyCondition, status: corev1.ConditionTrue}})
				})

				t.Run("should not tag anything if there's not tags", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Times(0)
					if _, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs); err != nil {
						_ = fmt.Errorf("reconcileNormal reutrned an error during test")
					}
				})

				t.Run("should tag instances from machine and cluster tags", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ms.AWSMachine.Spec.AdditionalTags = infrav1.Tags{"kind": "alicorn"}
					cs.AWSCluster.Spec.AdditionalTags = infrav1.Tags{"colour": "lavender"}

					ec2Svc.EXPECT().UpdateResourceTags(
						gomock.Any(),
						map[string]string{
							"kind": "alicorn",
						},
						map[string]string{},
					).Return(nil).Times(2)

					ec2Svc.EXPECT().UpdateResourceTags(
						PointsTo("myMachine"),
						map[string]string{
							"colour": "lavender",
							"kind":   "alicorn",
						},
						map[string]string{},
					).Return(nil)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})
				t.Run("should tag instances volume tags", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachineWithAdditionalTags()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ms.AWSMachine.Spec.AdditionalTags = infrav1.Tags{"rootDeviceID": "id1", "rootDeviceSize": "30"}

					ec2Svc.EXPECT().UpdateResourceTags(
						gomock.Any(),
						map[string]string{
							"rootDeviceID":   "id1",
							"rootDeviceSize": "30",
						},
						map[string]string{},
					).Return(nil).Times(3)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})
			})

			t.Run("temporarily stopping then starting the AWSMachine", func(t *testing.T) {
				var buf *bytes.Buffer
				getCoreSecurityGroups := func(t *testing.T, g *WithT) {
					t.Helper()

					buf = new(bytes.Buffer)
					klog.SetOutput(buf)
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
						Return(map[string][]string{"eid": {}}, nil).Times(1)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				}

				t.Run("should set instance to stopping and unready", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					instance.State = infrav1.InstanceStateStopping
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopping)))
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceStoppedReason}})
				})

				t.Run("should then set instance to stopped and unready", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					instance.State = infrav1.InstanceStateStopped
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopped)))
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceStoppedReason}})
				})

				t.Run("should then set instance to running and ready once it is restarted", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					instance.State = infrav1.InstanceStateRunning
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(true))
					g.Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
				})
			})
			t.Run("deleting the AWSMachine outside of Kubernetes", func(t *testing.T) {
				var buf *bytes.Buffer
				deleteMachine := func(t *testing.T, g *WithT) {
					t.Helper()

					buf = new(bytes.Buffer)
					klog.SetOutput(buf)
					secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				}

				t.Run("should warn if an instance is shutting-down", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					deleteMachine(t, g)
					instance.State = infrav1.InstanceStateShuttingDown
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					g.Expect(buf.String()).To(ContainSubstring(("Unexpected EC2 instance termination")))
					g.Eventually(recorder.Events).Should(Receive(ContainSubstring("UnexpectedTermination")))
				})

				t.Run("should error when the instance is seen as terminated", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					deleteMachine(t, g)

					instance.State = infrav1.InstanceStateTerminated
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					g.Expect(buf.String()).To(ContainSubstring(("Unexpected EC2 instance termination")))
					g.Eventually(recorder.Events).Should(Receive(ContainSubstring("UnexpectedTermination")))
					g.Expect(ms.AWSMachine.Status.FailureMessage).To(PointTo(Equal("EC2 instance state \"terminated\" is unexpected")))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceTerminatedReason}})
				})
			})
			t.Run("should not register if control plane ELB is already registered", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
			})
			t.Run("should attach control plane ELB to instance", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)
				elbSvc.EXPECT().RegisterInstanceWithAPIServerELB(gomock.Any()).Return(nil)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionTrue, "", ""}})
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
			})
			t.Run("Should store userdata using AWS Secrets Manager", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ms.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager = true
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).To(BeNil())
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
			})
			t.Run("should fail to delete bootstrap data secret if AWSMachine state is updated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)
				ms.Machine.Status.NodeRef = &corev1.ObjectReference{
					Namespace: "default",
					Name:      "test",
				}

				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(errors.New("failed to delete entries from AWS Secret")).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				g.Expect(err).To(MatchError(ContainSubstring("failed to delete entries from AWS Secret")))
			})
		})
		t.Run("when instance creation fails", func(t *testing.T) {
			var instance *infrav1.Instance
			instanceCreate := func(t *testing.T, g *WithT) {
				t.Helper()
				instance = &infrav1.Instance{
					ID:               "myMachine",
					VolumeIDs:        []string{"volume-1", "volume-2"},
					AvailabilityZone: "us-east-1",
				}
				instance.State = infrav1.InstanceStatePending
			}
			t.Run("Should fail while getting userdata", func(t *testing.T) {
				expectedError := "failed to generate init script"
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(expectedError)).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err.Error()).To(ContainSubstring(expectedError))

				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceProvisionFailedReason}})
			})
			t.Run("should fail to determine the registration status of control plane ELB", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(false, errors.New("error describing ELB"))
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring("error describing ELB"))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				g.Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedAttachControlPlaneELB")))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
			})
			t.Run("should fail to attach control plane ELB to instance", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)
				elbSvc.EXPECT().RegisterInstanceWithAPIServerELB(gomock.Any()).Return(errors.New("failed to attach ELB"))
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring("failed to attach ELB"))
				g.Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedAttachControlPlaneELB")))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
			})
			t.Run("should fail to delete bootstrap data secret if AWSMachine is in failed state", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				ms.SetSecretPrefix("test")
				ms.AWSMachine.Status.FailureReason = (*capierrors.MachineStatusError)(aws.String("error in AWSMachine"))

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("secretPrefix present, but secretCount is not set")))
			})
			t.Run("Should fail in ensureTag", func(t *testing.T) {
				id := providerID
				ensureTag := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				}

				t.Run("Should fail to return machine annotations after instance is created", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureTag(t, g)
					ms.AWSMachine.Annotations = map[string]string{TagsLastAppliedAnnotation: "12345"}

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err.Error()).To(ContainSubstring("json: cannot unmarshal number into Go value of type map[string]interface {}"))
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				})
				t.Run("Should fail to update resource tags after instance is created", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					id := providerID
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureTag(t, g)
					ms.AWSMachine.Annotations = map[string]string{TagsLastAppliedAnnotation: "{\"tag\":\"tag1\"}"}

					ec2Svc.EXPECT().UpdateResourceTags(gomock.Any(), gomock.Any(), map[string]string{"tag": "tag1"}).Return(errors.New("failed to update resource tag"))

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
				})
			})
			t.Run("While ensuring SecurityGroups", func(t *testing.T) {
				id := providerID
				ensureSecurityGroups := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
					secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil)
				}

				t.Run("Should fail to return machine annotations", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					id := providerID
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureSecurityGroups(t, g)
					ms.AWSMachine.Annotations = map[string]string{SecurityGroupsLastAppliedAnnotation: "12345"}

					ec2Svc.EXPECT().UpdateResourceTags(gomock.Any(), map[string]string{"tag": "\"old_tag\"\"\""}, gomock.Any()).Return(nil).AnyTimes()

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
				})
				t.Run("Should fail to fetch core security groups", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureSecurityGroups(t, g)
					ms.AWSMachine.Annotations = map[string]string{SecurityGroupsLastAppliedAnnotation: "{\"tag\":\"tag1\"}"}

					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, errors.New("failed to get core security groups")).Times(1)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
				})
				t.Run("Should return silently if ensureSecurityGroups fails to fetch additional security groups", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					id := providerID
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureSecurityGroups(t, g)
					ms.AWSMachine.Annotations = map[string]string{SecurityGroupsLastAppliedAnnotation: "{\"tag\":\"tag1\"}"}
					ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1.AWSResourceReference{
						{
							ID:  aws.String("sg-1"),
							ARN: aws.String("arn-1"),
							Filters: []infrav1.Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					}

					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
					ec2Svc.EXPECT().GetFilteredSecurityGroupID(gomock.Any()).Return("sg-1", errors.New("failed to get filtered SGs"))

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionTrue, "", ""}})
				})
				t.Run("Should fail to update security group", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					id := providerID
					ms.AWSMachine.Spec.ProviderID = &id

					instanceCreate(t, g)
					ensureSecurityGroups(t, g)
					ms.AWSMachine.Annotations = map[string]string{SecurityGroupsLastAppliedAnnotation: "{\"tag\":\"tag1\"}"}
					ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1.AWSResourceReference{
						{
							ID:  aws.String("sg-1"),
							ARN: aws.String("arn-1"),
							Filters: []infrav1.Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					}

					ec2Svc.EXPECT().GetFilteredSecurityGroupID(gomock.Any()).Return("sg-1", nil)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Return(errors.New("failed to update security groups"))

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
				})
			})
		})
	})

	t.Run("Secrets management lifecycle", func(t *testing.T) {
		t.Run("Secrets management lifecycle when creating EC2 instances", func(t *testing.T) {
			var instance *infrav1.Instance
			secretPrefix := "test/secret"

			t.Run("should leverage AWS Secrets Manager", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)

				instance = &infrav1.Instance{
					ID:    "myMachine",
					State: infrav1.InstanceStatePending,
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(secretPrefix, int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)

				ms.AWSMachine.ObjectMeta.Labels = map[string]string{
					clusterv1.MachineControlPlaneLabelName: "",
				}
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})
		})

		t.Run("Secrets management lifecycle when there's a node ref and a secret ARN", func(t *testing.T) {
			var instance *infrav1.Instance
			setNodeRef := func(t *testing.T, g *WithT) {
				t.Helper()

				instance = &infrav1.Instance{
					ID: "myMachine",
				}

				ms.Machine.Status.NodeRef = &corev1.ObjectReference{
					Kind:       "Node",
					Name:       "myMachine",
					APIVersion: "v1",
				}

				ms.AWSMachine.Spec.CloudInit = infrav1.CloudInit{
					SecretPrefix:         "secret",
					SecretCount:          5,
					SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				}
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(instance, nil).AnyTimes()
			}

			t.Run("should delete the secret if the instance is running", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				instance.State = infrav1.InstanceStateRunning
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
					Return(map[string][]string{"eid": {}}, nil).Times(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the instance is terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is deleted", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is in a failure condition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				ms.AWSMachine.Status.FailureReason = capierrors.MachineStatusErrorPtr(capierrors.UpdateMachineError)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})
		})

		t.Run("Secrets management lifecycle when there's only a secret ARN and no node ref", func(t *testing.T) {
			var instance *infrav1.Instance
			setSSM := func(t *testing.T, g *WithT) {
				t.Helper()

				instance = &infrav1.Instance{
					ID: "myMachine",
				}

				ms.AWSMachine.Spec.CloudInit = infrav1.CloudInit{
					SecretPrefix:         "secret",
					SecretCount:          5,
					SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				}
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(instance, nil).AnyTimes()
			}

			t.Run("should not delete the secret if the instance is running", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				instance.State = infrav1.InstanceStateRunning
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
					Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).MaxTimes(0)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the instance is terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is deleted", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is in a failure condition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				ms.AWSMachine.Status.FailureReason = capierrors.MachineStatusErrorPtr(capierrors.UpdateMachineError)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})
		})

		t.Run("Secrets management lifecycle when there is an intermittent connection issue and no secret could be stored", func(t *testing.T) {
			var instance *infrav1.Instance
			secretPrefix := "test/secret"

			getInstances := func(t *testing.T, g *WithT) {
				t.Helper()
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
			}

			t.Run("should error if secret could not be created", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				getInstances(t, g)

				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(secretPrefix, int32(0), errors.New("connection error")).Times(1)
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring("connection error"))
				g.Expect(ms.GetSecretPrefix()).To(Equal(""))
				g.Expect(ms.GetSecretCount()).To(Equal(int32(0)))
			})
			t.Run("should update prefix and count on successful creation", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				getInstances(t, g)

				instance = &infrav1.Instance{
					ID: "myMachine",
				}
				instance.State = infrav1.InstanceStatePending
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(secretPrefix, int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)

				g.Expect(err).To(BeNil())
				g.Expect(ms.GetSecretPrefix()).To(Equal(secretPrefix))
				g.Expect(ms.GetSecretCount()).To(Equal(int32(1)))
			})
		})
	})

	t.Run("Deleting an AWSMachine", func(t *testing.T) {
		finalizer := func(t *testing.T, g *WithT) {
			t.Helper()

			ms.AWSMachine.Finalizers = []string{
				infrav1.MachineFinalizer,
				metav1.FinalizerDeleteDependents,
			}
		}
		t.Run("should exit immediately on an error state", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			expectedErr := errors.New("no connection available ")
			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, expectedErr).AnyTimes()

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
		})
		t.Run("should log and remove finalizer when no machine exists", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("Unable to locate EC2 instance by ID or tags"))
			g.Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
			g.Eventually(recorder.Events).Should(Receive(ContainSubstring("NoInstanceFound")))
		})
		t.Run("should ignore instances in shutting down state", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
				State: infrav1.InstanceStateShuttingDown,
			}, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("EC2 instance is shutting down or already terminated"))
			g.Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
		})
		t.Run("should ignore instances in terminated down state", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
				State: infrav1.InstanceStateTerminated,
			}, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("EC2 instance is shutting down or already terminated"))
			g.Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
		})
		t.Run("instance not shutting down yet", func(t *testing.T) {
			id := "aws:////myid"
			getRunningInstance := func(t *testing.T, g *WithT) {
				t.Helper()
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{ID: id}, nil)
			}
			t.Run("should return an error when the instance can't be terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				getRunningInstance(t, g)

				expected := errors.New("can't reach AWS to terminate machine")
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(expected)

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expected))
				g.Expect(buf.String()).To(ContainSubstring("Terminating EC2 instance"))
				g.Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedTerminate")))
			})
			t.Run("when instance can be shut down", func(t *testing.T) {
				terminateInstance := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil)
				}

				t.Run("should error when it can't retrieve security groups if there are network interfaces", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					getRunningInstance(t, g)
					terminateInstance(t, g)

					ms.AWSMachine.Spec.NetworkInterfaces = []string{
						"eth0",
						"eth1",
					}
					expected := errors.New("can't reach AWS to list security groups")
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return(nil, expected)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(errors.Cause(err)).To(MatchError(expected))
				})

				t.Run("should error when it can't detach a security group from an interface", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					getRunningInstance(t, g)
					terminateInstance(t, g)

					ms.AWSMachine.Spec.NetworkInterfaces = []string{
						"eth0",
						"eth1",
					}
					expected := errors.New("can't reach AWS to detach security group")
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{"sg0", "sg1"}, nil)
					ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(gomock.Any(), gomock.Any()).Return(expected)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(errors.Cause(err)).To(MatchError(expected))
				})

				t.Run("should detach all combinations of network interfaces", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					getRunningInstance(t, g)
					terminateInstance(t, g)

					ms.AWSMachine.Spec.NetworkInterfaces = []string{
						"eth0",
						"eth1",
					}
					groups := []string{"sg0", "sg1"}
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{"sg0", "sg1"}, nil)
					ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(groups, "eth0").Return(nil)
					ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(groups, "eth1").Return(nil)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})

				t.Run("should remove security groups", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					getRunningInstance(t, g)
					terminateInstance(t, g)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
				})

				t.Run("should fail to detach control plane ELB from instance", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
					ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
					reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
						return elbSvc
					}

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
						State: infrav1.InstanceStateTerminated,
					}, nil)
					elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(false, errors.New("error describing ELB"))

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(err.Error()).To(ContainSubstring("error describing ELB"))
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
					g.Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedDetachControlPlaneELB")))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, "DeletingFailed"}})
				})

				t.Run("should not do anything if control plane ELB is already detached from instance", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
					ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
					reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
						return elbSvc
					}

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
						State: infrav1.InstanceStateTerminated,
					}, nil)
					elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					g.Expect(err).To(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason}})
				})
			})
		})
		t.Run("Reconcile LB detachment", func(t *testing.T) {
			t.Run("should fail to determine registration status of ELB", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
					State: infrav1.InstanceStateTerminated,
				}, nil)
				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				elbSvc.EXPECT().DeregisterInstanceFromAPIServerELB(gomock.Any()).Return(nil)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, clusterv1.DeletedReason}})
			})
			t.Run("should fail to detach control plane ELB from instance", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabelName: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
					State: infrav1.InstanceStateTerminated,
				}, nil)
				elbSvc.EXPECT().InstanceIsRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				elbSvc.EXPECT().DeregisterInstanceFromAPIServerELB(gomock.Any()).Return(errors.New("Duplicate access point name for load balancer"))

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring("Duplicate access point name for load balancer"))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, "DeletingFailed"}})
			})
			t.Run("should fail if secretPrefix present, but secretCount is not set", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				ms.SetSecretPrefix("test")

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("secretPrefix present, but secretCount is not set")))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
			})
			t.Run("should fail if secrets backend is invalid", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				ms.AWSMachine.Spec.CloudInit.SecureSecretsBackend = "InvalidSecretBackend"

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("invalid secret backend")))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
			})
			t.Run("should fail if deleting entries from AWS Secret fails", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				ms.SetSecretPrefix("test")
				ms.SetSecretCount(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(errors.New("Hierarchy Type Mismatch Exception"))

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("Hierarchy Type Mismatch Exception")))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
			})
		})
	})
}

func getAWSMachine() *infrav1.AWSMachine {
	return &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: infrav1.AWSMachineSpec{
			CloudInit: infrav1.CloudInit{
				SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
			},
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
