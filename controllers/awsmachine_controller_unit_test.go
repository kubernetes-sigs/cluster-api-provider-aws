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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	ec2Service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	elbService "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
)

const providerID = "aws:////myMachine"

func TestAWSMachineReconciler(t *testing.T) {
	var (
		reconciler     AWSMachineReconciler
		cs             *scope.ClusterScope
		ms             *scope.MachineScope
		mockCtrl       *gomock.Controller
		ec2Svc         *mock_services.MockEC2Interface
		elbSvc         *mock_services.MockELBInterface
		secretSvc      *mock_services.MockSecretInterface
		objectStoreSvc *mock_services.MockObjectStoreInterface
		recorder       *record.FakeRecorder
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

		secretIgnition := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap-data-ignition",
			},
			Data: map[string][]byte{
				"value":  []byte("ignitionJSON"),
				"format": []byte("ignition"),
			},
		}

		client := fake.NewClientBuilder().WithObjects(awsMachine, secret, secretIgnition).WithStatusSubresource(awsMachine).Build()
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
						ClusterName: "capi-test",
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: ptr.To[string]("bootstrap-data"),
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
				Client:     fake.NewClientBuilder().WithObjects(awsMachine, secret).WithStatusSubresource(awsMachine).Build(),
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "test"}},
			},
		)
		g.Expect(err).To(BeNil())
		cs.AWSCluster = &infrav1.AWSCluster{
			Spec: infrav1.AWSClusterSpec{
				ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
					LoadBalancerType: infrav1.LoadBalancerTypeClassic,
				},
			},
		}
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
						ClusterName: "capi-test",
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: ptr.To[string]("bootstrap-data"),
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
		objectStoreSvc = mock_services.NewMockObjectStoreInterface(mockCtrl)

		// If your test hangs for 9 minutes, increase the value here to the number of events during a reconciliation loop
		recorder = record.NewFakeRecorder(2)

		reconciler = AWSMachineReconciler{
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2Interface {
				return ec2Svc
			},
			secretsManagerServiceFactory: func(cloud.ClusterScoper) services.SecretInterface {
				return secretSvc
			},
			objectStoreServiceFactory: func(cloud.ClusterScoper) services.ObjectStoreInterface {
				return objectStoreSvc
			},
			Recorder: recorder,
			Log:      klog.Background(),
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
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()
			}

			t.Run("should exit immediately on an error state", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				er := "CreateError"
				ms.AWSMachine.Status.FailureReason = &er
				ms.AWSMachine.Status.FailureMessage = ptr.To[string]("Couldn't create machine")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)

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
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			t.Run("shouldn't add our finalizer to the machine", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				runningInstance(t, g)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)

				g.Expect(len(ms.AWSMachine.Finalizers)).To(Equal(0))
			})
		})

		t.Run("when provider ID is populated correctly", func(t *testing.T) {
			id := providerID
			providerID := func(t *testing.T, g *WithT) {
				t.Helper()
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			t.Run("should fail to create instance and keep the finalizers as is", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)

				providerID(t, g)
				expectedErr := errors.New("Invalid instance")
				ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

			_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil)
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)

					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStatePending)))
					g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
					g.Expect(buf.String()).To(ContainSubstring("EC2 instance state changed"))

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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)

					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
					g.Expect(ms.AWSMachine.Status.Ready).To(BeTrue())
					g.Expect(buf.String()).To(ContainSubstring("EC2 instance state changed"))
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
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
				g.Expect(buf.String()).To(ContainSubstring("EC2 instance state is undefined"))
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
							ID: ptr.To[string]("sg-2345"),
						},
					}
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(instance.ID, []string{"sg-2345"})
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return([]string{"sg-2345"}, nil)

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{conditionType: infrav1.SecurityGroupsReadyCondition, status: corev1.ConditionTrue}})
				})

				t.Run("should not tag instances if there's no tags", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Times(0)
					if _, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs); err != nil {
						_ = fmt.Errorf("reconcileNormal reutrned an error during test")
					}
				})

				t.Run("should tag instances and volumes with machine and cluster tags", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					ms.AWSMachine.Spec.AdditionalTags = infrav1.Tags{"kind": "alicorn", "colour": "pink"} // takes precedence
					cs.AWSCluster.Spec.AdditionalTags = infrav1.Tags{"colour": "lavender", "shape": "round"}

					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

					// expect one call first to tag the instance and two calls for tagging each of two volumes
					// the volumes get the tags from the AWSCluster _and_ the AWSMachine

					ec2Svc.EXPECT().UpdateResourceTags(
						PointsTo("myMachine"),
						map[string]string{
							"colour": "pink",
							"shape":  "round",
							"kind":   "alicorn",
						},
						map[string]string{},
					).Return(nil)

					ec2Svc.EXPECT().UpdateResourceTags(
						gomock.Any(),
						map[string]string{
							"colour": "pink",
							"shape":  "round",
							"kind":   "alicorn",
						},
						map[string]string{},
					).Return(nil).Times(2)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().UpdateResourceTags(
						gomock.Any(),
						map[string]string{
							"rootDeviceID":   "id1",
							"rootDeviceSize": "30",
						},
						map[string]string{},
					).Return(nil).Times(3)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})
			})

			t.Run("temporarily stopping then starting the AWSMachine(stateless)", func(t *testing.T) {
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
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
				}

				t.Run("should set instance to stopping and unready", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					instanceCreate(t, g)
					getCoreSecurityGroups(t, g)

					instance.State = infrav1.InstanceStateStopping
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopping)))
					g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
					g.Expect(buf.String()).To(ContainSubstring("EC2 instance state changed"))
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopped)))
					g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
					g.Expect(buf.String()).To(ContainSubstring("EC2 instance state changed"))
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
					g.Expect(ms.AWSMachine.Status.Ready).To(BeTrue())
					g.Expect(buf.String()).To(ContainSubstring("EC2 instance state changed"))
				})
			})
			t.Run("deleting the AWSMachine manually", func(t *testing.T) {
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
					g.Expect(buf.String()).To(ContainSubstring("Unexpected EC2 instance termination"))
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
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(ms.AWSMachine.Status.Ready).To(BeFalse())
					g.Expect(buf.String()).To(ContainSubstring("Unexpected EC2 instance termination"))
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

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)
				elbSvc.EXPECT().RegisterInstanceWithAPIServerELB(gomock.Any()).Return(nil)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.ELBAttachedCondition, corev1.ConditionTrue, "", ""}})
				expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
			})
			t.Run("should store userdata for CloudInit using AWS Secrets Manager only when not skipped", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				instanceCreate(t, g)

				// Explicitly skip AWS Secrets Manager.
				ms.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager = true
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil)
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(false, errors.New("error describing ELB"))
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil)
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)
				elbSvc.EXPECT().RegisterInstanceWithAPIServerELB(gomock.Any()).Return(errors.New("failed to attach ELB"))
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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
				ms.AWSMachine.Status.FailureReason = aws.String("error in AWSMachine")
				ms.SetSecretCount(0)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("secretPrefix present, but secretCount is not set")))
			})
			t.Run("Should fail in ensureTag", func(t *testing.T) {
				id := providerID
				ensureTag := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil)
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

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil)
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

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
				})
				t.Run("Should fail if ensureSecurityGroups fails to fetch additional security groups", func(t *testing.T) {
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
							Filters: []infrav1.Filter{
								{
									Name:   "example-name",
									Values: []string{"example-value"},
								},
							},
						},
					}

					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return([]string{"sg-1"}, errors.New("failed to get filtered SGs"))

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
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
							Filters: []infrav1.Filter{
								{
									Name:   "id",
									Values: []string{"sg-1"},
								},
							},
						},
					}

					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Return(errors.New("failed to update security groups"))
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return([]string{"sg-1"}, nil)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
					expectConditions(g, ms.AWSMachine, []conditionAssertion{{infrav1.SecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.SecurityGroupsFailedReason}})
				})
			})
		})
		t.Run("when BYOIP is set", func(t *testing.T) {
			var instance *infrav1.Instance
			secretPrefix := "test/secret"

			t.Run("should succeed", func(t *testing.T) {
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
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil).Times(1)

				ms.AWSMachine.Spec.PublicIP = aws.Bool(false)
				ms.AWSMachine.Spec.ElasticIPPool = &infrav1.ElasticIPPool{
					PublicIpv4Pool:              aws.String("ipv4pool-ec2-0123456789abcdef0"),
					PublicIpv4PoolFallBackOrder: ptr.To(infrav1.PublicIpv4PoolFallbackOrderAmazonPool),
				}
				ec2Svc.EXPECT().ReconcileElasticIPFromPublicPool(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(err).To(BeNil())
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
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

				ms.AWSMachine.ObjectMeta.Labels = map[string]string{
					clusterv1.MachineControlPlaneLabel: "",
				}
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
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
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the instance is terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is deleted", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is in a failure condition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				ms.AWSMachine.Status.FailureReason = ptr.To("UpdateError")
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			})
			t.Run("should not attempt to delete the secret if InsecureSkipSecretsManager is set on CloudInit", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				ms.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager = true

				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(0)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()

				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			})
			t.Run("should delete the secret from the S3 bucket if StorageType ClusterObjectStore is set for Ignition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				ms.AWSMachine.Spec.CloudInit = infrav1.CloudInit{}
				ms.AWSMachine.Spec.Ignition = &infrav1.Ignition{
					Version:     "2.3",
					StorageType: infrav1.IgnitionStorageTypeOptionClusterObjectStore,
				}

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				g.Expect(err).To(BeNil())
			})
			t.Run("should not delete the secret from the S3 bucket if StorageType UnencryptedUserData is set for Ignition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setNodeRef(t, g)

				ms.AWSMachine.Spec.CloudInit = infrav1.CloudInit{}
				ms.AWSMachine.Spec.Ignition = &infrav1.Ignition{
					Version:     "2.3",
					StorageType: infrav1.IgnitionStorageTypeOptionUnencryptedUserData,
				}

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(0)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				g.Expect(err).To(BeNil())
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
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).MaxTimes(0)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the instance is terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is deleted", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			})

			t.Run("should delete the secret if the AWSMachine is in a failure condition", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				setSSM(t, g)

				ms.AWSMachine.Status.FailureReason = ptr.To("UpdateError")
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring("connection error"))
				g.Expect(ms.GetSecretPrefix()).To(Equal("prefix"))
				g.Expect(ms.GetSecretCount()).To(Equal(int32(1000)))
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
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)

				g.Expect(err).To(BeNil())
				g.Expect(ms.GetSecretPrefix()).To(Equal(secretPrefix))
				g.Expect(ms.GetSecretCount()).To(Equal(int32(1)))
			})
		})
	})

	t.Run("Object storage lifecycle for Ignition's userdata", func(t *testing.T) {
		t.Run("when Ignition's StorageType is ClusterObjectStore", func(t *testing.T) {
			useIgnitionWithClusterObjectStore := func(t *testing.T, g *WithT) {
				t.Helper()

				ms.Machine.Spec.Bootstrap.DataSecretName = ptr.To[string]("bootstrap-data-ignition")
				ms.AWSMachine.Spec.CloudInit.SecretCount = 0
				ms.AWSMachine.Spec.CloudInit.SecretPrefix = ""
				ms.AWSMachine.Spec.Ignition = &infrav1.Ignition{
					Version:     "2.3",
					StorageType: infrav1.IgnitionStorageTypeOptionClusterObjectStore,
				}
			}

			t.Run("creating EC2 instances", func(t *testing.T) {
				var instance *infrav1.Instance

				getInstances := func(t *testing.T, g *WithT) {
					t.Helper()

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
				}

				t.Run("should leverage a Cluster Object Store", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance = &infrav1.Instance{
						ID:    "myMachine",
						State: infrav1.InstanceStatePending,
					}
					fakeS3URL := "s3://foo"

					objectStoreSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fakeS3URL, nil).Times(1)
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

					ms.AWSMachine.ObjectMeta.Labels = map[string]string{
						clusterv1.MachineControlPlaneLabel: "",
					}

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})

				t.Run("should leverage a Cluster Object Store with presigned urls", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					if cs.AWSCluster.Spec.S3Bucket == nil {
						cs.AWSCluster.Spec.S3Bucket = &infrav1.S3Bucket{}
					}
					cs.AWSCluster.Spec.S3Bucket.PresignedURLDuration = &metav1.Duration{Duration: 1 * time.Hour}

					instance = &infrav1.Instance{
						ID:    "myMachine",
						State: infrav1.InstanceStatePending,
					}

					//nolint:gosec
					presigned := "https://cluster-api-aws.s3.us-west-2.amazonaws.com/bootstrap-data.yaml?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA3SGQVQG7FGA6KKA6%2F20221104%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20221104T140227Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=b228dbec8c1008c80c162e1210e4503dceead1e4d4751b4d9787314fd6da4d55"

					objectStoreSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(presigned, nil).Times(1)
					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

					ms.AWSMachine.ObjectMeta.Labels = map[string]string{
						clusterv1.MachineControlPlaneLabel: "",
					}

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})
			})

			t.Run("there's a node ref and a secret ARN", func(t *testing.T) {
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

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(instance, nil).AnyTimes()
				}

				t.Run("should delete the object if the instance is running", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					setNodeRef(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance.State = infrav1.InstanceStateRunning
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the instance is terminated", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					setNodeRef(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance.State = infrav1.InstanceStateTerminated
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the instance is deleted", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					setNodeRef(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance.State = infrav1.InstanceStateRunning
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()

					_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the AWSMachine is in a failure condition", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					setNodeRef(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					// TODO: This seems to have no effect on the test result.
					ms.AWSMachine.Status.FailureReason = ptr.To("UpdateError")

					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()

					_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				})
			})

			t.Run("there's only a secret ARN and no node ref", func(t *testing.T) {
				var instance *infrav1.Instance

				getInstances := func(t *testing.T, g *WithT) {
					t.Helper()

					instance = &infrav1.Instance{
						ID: "myMachine",
					}
					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(instance, nil).AnyTimes()
				}

				t.Run("should not delete the object if the instance is running", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)

					instance.State = infrav1.InstanceStateRunning
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).MaxTimes(0)
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the instance is terminated", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance.State = infrav1.InstanceStateTerminated
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the AWSMachine is deleted", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					instance.State = infrav1.InstanceStateRunning
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
					_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				})

				t.Run("should delete the object if the AWSMachine is in a failure condition", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					// TODO: This seems to have no effect on the test result.
					ms.AWSMachine.Status.FailureReason = ptr.To("UpdateError")
					objectStoreSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil).AnyTimes()
					_, _ = reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				})
			})

			t.Run("there is an intermittent connection issue and no object could be created", func(t *testing.T) {
				t.Run("should error if object could not be created", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					useIgnitionWithClusterObjectStore(t, g)

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
					objectStoreSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errors.New("connection error")).Times(1)
					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).ToNot(BeNil())
					g.Expect(err.Error()).To(ContainSubstring("connection error"))
				})
			})
		})

		t.Run("when Ignition's StorageType is UnencryptedUserData", func(t *testing.T) {
			useIgnitionAndUnencryptedUserData := func(t *testing.T, g *WithT) {
				t.Helper()

				ms.Machine.Spec.Bootstrap.DataSecretName = ptr.To[string]("bootstrap-data-ignition")
				ms.AWSMachine.Spec.CloudInit.SecretCount = 0
				ms.AWSMachine.Spec.CloudInit.SecretPrefix = ""
				ms.AWSMachine.Spec.Ignition = &infrav1.Ignition{
					Version:     "2.3",
					StorageType: infrav1.IgnitionStorageTypeOptionUnencryptedUserData,
				}
			}
			t.Run("creating EC2 instances", func(t *testing.T) {
				var instance *infrav1.Instance

				getInstances := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
				}
				t.Run("should NOT leverage a Cluster Object Store", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					getInstances(t, g)
					useIgnitionAndUnencryptedUserData(t, g)

					instance = &infrav1.Instance{
						ID:    "myMachine",
						State: infrav1.InstanceStatePending,
					}
					fakeS3URL := "s3://foo"

					// Expect no Cluster Object Store to be created.
					objectStoreSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fakeS3URL, nil).Times(0)

					ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any(), gomock.Any()).Return(instance, nil).AnyTimes()
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(map[string][]string{"eid": {}}, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
					ec2Svc.EXPECT().GetAdditionalSecurityGroupsIDs(gomock.Any()).Return(nil, nil)

					ms.AWSMachine.ObjectMeta.Labels = map[string]string{
						clusterv1.MachineControlPlaneLabel: "",
					}
					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})
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
			secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
		})
		t.Run("should log and remove finalizer when no machine exists", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
			secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
			secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("EC2 instance is shutting down or already terminated"))
		})
		t.Run("should ignore instances in terminated state", func(t *testing.T) {
			g := NewWithT(t)
			awsMachine := getAWSMachine()
			setup(t, g, awsMachine)
			defer teardown(t, g)
			finalizer(t, g)

			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
				State: infrav1.InstanceStateTerminated,
			}, nil)
			secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("EC2 instance terminated successfully"))
			g.Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
		})
		t.Run("instance not shutting down yet", func(t *testing.T) {
			id := "aws:////myid"
			getRunningInstance := func(t *testing.T, g *WithT) {
				t.Helper()
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{ID: id}, nil)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()
			}
			t.Run("should return an error when the instance can't be terminated", func(t *testing.T) {
				g := NewWithT(t)
				awsMachine := getAWSMachine()
				setup(t, g, awsMachine)
				defer teardown(t, g)
				finalizer(t, g)
				getRunningInstance(t, g)

				expected := errors.New("can't reach AWS to terminate machine")
				ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(expected)

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expected))
				g.Expect(buf.String()).To(ContainSubstring("Terminating EC2 instance"))
				g.Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedTerminate")))
			})
			t.Run("when instance can be shut down", func(t *testing.T) {
				terminateInstance := func(t *testing.T, g *WithT) {
					t.Helper()
					ec2Svc.EXPECT().TerminateInstance(gomock.Any()).Return(nil)
					secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()
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

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
					g.Expect(err).To(BeNil())
				})

				t.Run("should fail to detach control plane ELB from instance", func(t *testing.T) {
					g := NewWithT(t)
					awsMachine := getAWSMachine()
					setup(t, g, awsMachine)
					defer teardown(t, g)
					finalizer(t, g)
					ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
					ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
					reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
						return elbSvc
					}

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
						State: infrav1.InstanceStateTerminated,
					}, nil)
					elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(false, errors.New("error describing ELB"))

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
					ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
					ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
					reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
						return elbSvc
					}

					ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
						State: infrav1.InstanceStateTerminated,
					}, nil)
					elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(false, nil)

					_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
					State: infrav1.InstanceStateTerminated,
				}, nil)
				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				elbSvc.EXPECT().DeregisterInstanceFromAPIServerELB(gomock.Any()).Return(nil)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
				ms.Machine.Labels = map[string]string{clusterv1.MachineControlPlaneLabel: ""}
				ms.AWSMachine.Status.InstanceState = &infrav1.InstanceStateStopping
				reconciler.elbServiceFactory = func(elbScope scope.ELBScope) services.ELBInterface {
					return elbSvc
				}

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
					State: infrav1.InstanceStateTerminated,
				}, nil)
				elbSvc.EXPECT().IsInstanceRegisteredWithAPIServerELB(gomock.Any()).Return(true, nil)
				elbSvc.EXPECT().DeregisterInstanceFromAPIServerELB(gomock.Any()).Return(errors.New("Duplicate access point name for load balancer"))

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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
				ms.SetSecretCount(0)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
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

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs, cs)
				g.Expect(err).To(MatchError(ContainSubstring("Hierarchy Type Mismatch Exception")))
				g.Expect(ms.AWSMachine.Finalizers).To(ContainElement(metav1.FinalizerDeleteDependents))
			})
		})
	})
}

func TestAWSMachineReconcilerAWSClusterToAWSMachines(t *testing.T) {
	testCases := []struct {
		name         string
		ownerCluster *clusterv1.Cluster
		awsCluster   *infrav1.AWSCluster
		awsMachine   *clusterv1.Machine
		requests     []reconcile.Request
	}{
		{
			name:         "Should create reconcile request successfully",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-6"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-6",
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-6",
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-6",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-6",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-6",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "aws-machine-6",
					},
				},
			},
		},
		{
			name:         "Should not create reconcile request for deleted clusters",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", DeletionTimestamp: &metav1.Time{Time: time.Now()}}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "aws-test-1",
					},
					Name: "aws-test-1",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-1",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-1",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-1",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
					DeletionTimestamp: &metav1.Time{Time: time.Now()},
				},
			},
		},
		{
			name: "Should not create reconcile request if ownerCluster not found",
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "aws-test-2",
					},
					Name: "aws-test-2",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-2",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-2",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-2",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
		},
		{
			name:         "Should not create reconcile request if owned Machines not found",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-3"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-3",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						Name:       "aws-machine-3",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-3",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-3",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
		{
			name:         "Should not create reconcile request if owned Machine type is not AWSMachine",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-4"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-4",
					},
					Name:      "aws-test-4",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					Kind: "Machine",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "Machine",
						Name:       "aws-machine-4",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-4",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-4",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
		{
			name:         "Should not create reconcile request if name for machine in infrastructure ref not found",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-5"}},
			awsMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5",
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-5",
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachine",
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5",
					OwnerReferences: []metav1.OwnerReference{
						{
							Name:       "capi-test-5",
							Kind:       "Cluster",
							APIVersion: clusterv1.GroupVersion.String(),
						},
					},
				},
			},
			requests: []reconcile.Request{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
				Log:    klog.Background(),
			}
			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())

			createObject(g, tc.ownerCluster, ns.Name)
			defer cleanupObject(g, tc.ownerCluster)

			createObject(g, tc.awsMachine, ns.Name)
			defer cleanupObject(g, tc.awsMachine)

			tc.awsCluster.Namespace = ns.Name
			defer t.Cleanup(func() {
				g.Expect(testEnv.Cleanup(ctx, tc.awsCluster, ns)).To(Succeed())
			})

			requests := reconciler.AWSClusterToAWSMachines(logger.NewLogger(klog.Background()))(ctx, tc.awsCluster)
			if tc.requests != nil {
				if len(tc.requests) > 0 {
					tc.requests[0].Namespace = ns.Name
				}
				g.Expect(requests).To(ConsistOf(tc.requests))
			} else {
				g.Expect(requests).To(BeNil())
			}
		})
	}
}

func TestAWSMachineReconcilerRequeueAWSMachinesForUnpausedCluster(t *testing.T) {
	testCases := []struct {
		name         string
		ownerCluster *clusterv1.Cluster
		requests     []reconcile.Request
	}{
		{
			name:         "Should not create reconcile request for deleted clusters",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default", DeletionTimestamp: &metav1.Time{Time: time.Now()}}},
		},
		{
			name:         "Should create reconcile request successfully",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}},
			requests:     []reconcile.Request{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
				Log:    klog.Background(),
			}
			requests := reconciler.requeueAWSMachinesForUnpausedCluster(logger.NewLogger(klog.Background()))(ctx, tc.ownerCluster)
			if tc.requests != nil {
				g.Expect(requests).To(ConsistOf(tc.requests))
			} else {
				g.Expect(requests).To(BeNil())
			}
		})
	}
}

func TestAWSMachineReconcilerIndexAWSMachineByInstanceID(t *testing.T) {
	t.Run("Should not return instance id if cluster type is not AWSCluster", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klog.Background(),
		}
		machine := &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}}
		requests := reconciler.indexAWSMachineByInstanceID(machine)
		g.Expect(requests).To(BeNil())
	})
	t.Run("Should return instance id successfully", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klog.Background(),
		}
		awsMachine := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}, Spec: infrav1.AWSMachineSpec{InstanceID: aws.String("12345")}}
		requests := reconciler.indexAWSMachineByInstanceID(awsMachine)
		g.Expect(requests).To(ConsistOf([]string{"12345"}))
	})
	t.Run("Should not return instance id if instance id is not present", func(t *testing.T) {
		g := NewWithT(t)
		reconciler := &AWSMachineReconciler{
			Client: testEnv.Client,
			Log:    klog.Background(),
		}
		awsMachine := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: "default"}}
		requests := reconciler.indexAWSMachineByInstanceID(awsMachine)
		g.Expect(requests).To(BeNil())
	})
}

func TestAWSMachineReconcilerReconcile(t *testing.T) {
	testCases := []struct {
		name         string
		awsMachine   *infrav1.AWSMachine
		ownerMachine *clusterv1.Machine
		ownerCluster *clusterv1.Cluster
		awsCluster   *infrav1.AWSCluster
		expectError  bool
		requeue      bool
	}{
		{
			name:        "Should Reconcile successfully if no AWSMachine found",
			expectError: false,
		},
		{
			name:        "Should Reconcile AWSMachine with requeue",
			awsMachine:  &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "aws-test-1"}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"}},
			requeue:     true,
			expectError: false,
		},
		{
			name: "Should fail Reconcile with GetOwnerMachine failure",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-2",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				},
				Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			expectError: true,
		},
		{
			name: "Should not Reconcile if machine does not contain cluster label",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-3", Annotations: map[string]string{clusterv1.PausedAnnotation: ""}, OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "capi-test-machine",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
				},
			},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"}},
			expectError:  false,
		},
		{
			name: "Should not Reconcile if cluster is paused",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-4", Annotations: map[string]string{clusterv1.PausedAnnotation: ""}, OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-1",
					},
					Name: "capi-test-machine", Namespace: "default",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
				},
			},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"}, Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{Name: "foo"},
			}},
			expectError: false,
		},
		{
			name: "Should not Reconcile if AWSManagedControlPlane is not ready",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-1",
					},
					Name: "capi-test-machine", Namespace: "default",
				}, Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
				},
			},
			ownerCluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					ControlPlaneRef: &corev1.ObjectReference{Kind: AWSManagedControlPlaneRefKind},
				},
			},
			expectError: false,
		},
		{
			name: "Should not Reconcile if AWSCluster is not ready",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-1",
					},
					Name: "capi-test-machine", Namespace: "default",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
				},
			},
			ownerCluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: &corev1.ObjectReference{Name: "aws-test-5"},
				},
			},
			expectError: false,
		},
		{
			name: "Should fail to reconcile while fetching infra cluster",
			awsMachine: &infrav1.AWSMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-test-5", OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: clusterv1.GroupVersion.String(),
							Kind:       "Machine",
							Name:       "capi-test-machine",
							UID:        "1",
						},
					},
				}, Spec: infrav1.AWSMachineSpec{InstanceType: "test"},
			},
			ownerMachine: &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						clusterv1.ClusterNameLabel: "capi-test-1",
					},
					Name: "capi-test-machine", Namespace: "default",
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "capi-test",
				},
			},
			ownerCluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1"},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: &corev1.ObjectReference{Name: "aws-test-5"},
				},
			},
			awsCluster:  &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "aws-test-5"}},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSMachineReconciler{
				Client: testEnv.Client,
			}
			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())
			defer func() {
				g.Expect(testEnv.Cleanup(ctx, ns)).To(Succeed())
			}()

			createObject(g, tc.ownerCluster, ns.Name)
			defer cleanupObject(g, tc.ownerCluster)

			createObject(g, tc.awsCluster, ns.Name)
			defer cleanupObject(g, tc.awsCluster)

			createObject(g, tc.ownerMachine, ns.Name)
			defer cleanupObject(g, tc.ownerMachine)

			createObject(g, tc.awsMachine, ns.Name)
			defer cleanupObject(g, tc.awsMachine)
			if tc.awsMachine != nil {
				g.Eventually(func() bool {
					machine := &infrav1.AWSMachine{}
					key := client.ObjectKey{
						Name:      tc.awsMachine.Name,
						Namespace: ns.Name,
					}
					err = testEnv.Get(ctx, key, machine)
					return err == nil
				}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed get the newly created machine %q", tc.awsMachine.Name))

				result, err := reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: tc.awsMachine.Namespace,
						Name:      tc.awsMachine.Name,
					},
				})
				if tc.expectError {
					g.Expect(err).ToNot(BeNil())
				} else {
					g.Expect(err).To(BeNil())
				}
				if tc.requeue {
					g.Expect(result.RequeueAfter).To(BeZero())
				}
			} else {
				_, err = reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: "default",
						Name:      "test",
					},
				})
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestAWSMachineReconcilerReconcileDefaultsToLoadBalancerTypeClassic(t *testing.T) {
	// When working with an outdated v1beta2 CRD by mistake, it could happen that
	// `AWSCluster.Spec.ControlPlaneLoadBalancer.LoadBalancerType` was not set, but the object still written to etcd.
	// This test simulates this case using a fake client. The controller should still handle that value by assuming
	// classic LB as the type, since that is the default. It should not mistakenly try to reconcile against a v2 LB.

	g := NewWithT(t)

	ns := "testns"

	cp := &kubeadmv1beta1.KubeadmControlPlane{}
	cp.SetName("capi-cp-test-1")
	cp.SetNamespace(ns)

	ownerCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "capi-test-1", Namespace: ns},
		Spec: clusterv1.ClusterSpec{
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				Name:       "capi-test-1", // assuming same name
				Namespace:  ns,
				APIVersion: infrav1.GroupVersion.String(),
			},
			ControlPlaneRef: &corev1.ObjectReference{
				Kind:       "KubeadmControlPlane",
				Namespace:  cp.Namespace,
				Name:       cp.Name,
				APIVersion: kubeadmv1beta1.GroupVersion.String(),
			},
		},
		Status: clusterv1.ClusterStatus{
			InfrastructureReady: true,
		},
	}

	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "capi-test-1",
			Namespace: ns,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       ownerCluster.Name,
					UID:        "1",
				},
			},
		},
		Spec: infrav1.AWSClusterSpec{
			ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
				Scheme: &infrav1.ELBSchemeInternetFacing,
				// `LoadBalancerType` not set (i.e. empty string; must default to attaching instance to classic LB)
			},
			NetworkSpec: infrav1.NetworkSpec{
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:       "subnet-1",
						IsPublic: false,
					},
					infrav1.SubnetSpec{
						IsPublic: false,
					},
				},
			},
		},
		Status: infrav1.AWSClusterStatus{
			Ready: true,
			Network: infrav1.NetworkStatus{
				SecurityGroups: map[infrav1.SecurityGroupRole]infrav1.SecurityGroup{
					infrav1.SecurityGroupControlPlane: {
						ID: "1",
					},
					infrav1.SecurityGroupNode: {
						ID: "2",
					},
					infrav1.SecurityGroupLB: {
						ID: "3",
					},
				},
			},
		},
	}

	ownerMachine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterv1.ClusterNameLabel:         "capi-test-1",
				clusterv1.MachineControlPlaneLabel: "", // control plane node so that controller tries to register it with LB
			},
			Name:      "capi-test-machine",
			Namespace: ns,
		},
		Spec: clusterv1.MachineSpec{
			ClusterName: "capi-test",
			Bootstrap: clusterv1.Bootstrap{
				DataSecretName: aws.String("bootstrap-data"),
			},
		},
	}

	awsMachine := &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-test-7",
			Namespace: ns,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Machine",
					Name:       "capi-test-machine",
					UID:        "1",
				},
			},
		},
		Spec: infrav1.AWSMachineSpec{
			InstanceType: "test",
			ProviderID:   aws.String("aws://the-zone/two"),
			CloudInit: infrav1.CloudInit{
				SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				SecretPrefix:         "prefix",
				SecretCount:          1000,
			},
		},
		Status: infrav1.AWSMachineStatus{
			Conditions: clusterv1.Conditions{
				{
					Type:   "Paused",
					Status: corev1.ConditionFalse,
					Reason: "NotPaused",
				},
			},
		},
	}

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

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap-data",
			Namespace: ns,
		},
		Data: map[string][]byte{
			"value": []byte("shell-script"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithObjects(ownerCluster, awsCluster, ownerMachine, awsMachine, controllerIdentity, secret, cp).WithStatusSubresource(awsCluster, awsMachine).Build()

	recorder := record.NewFakeRecorder(10)
	reconciler := &AWSMachineReconciler{
		Client:   fakeClient,
		Recorder: recorder,
	}

	mockCtrl := gomock.NewController(t)
	ec2Mock := mocks.NewMockEC2API(mockCtrl)
	elbMock := mocks.NewMockELBAPI(mockCtrl)
	secretMock := mock_services.NewMockSecretInterface(mockCtrl)

	cs, err := getClusterScope(*awsCluster)
	g.Expect(err).To(BeNil())

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

	ec2Mock.EXPECT().DescribeInstancesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{"two"}),
	})).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("two"),
						InstanceType: aws.String("m5.large"),
						SubnetId:     aws.String("subnet-1"),
						ImageId:      aws.String("ami-1"),
						State: &ec2.InstanceState{
							Name: aws.String(ec2.InstanceStateNameRunning),
						},
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("thezone"),
						},
						MetadataOptions: &ec2.InstanceMetadataOptionsResponse{
							HttpEndpoint:            aws.String(string(infrav1.InstanceMetadataEndpointStateEnabled)),
							HttpPutResponseHopLimit: aws.Int64(1),
							HttpTokens:              aws.String(string(infrav1.HTTPTokensStateOptional)),
							InstanceMetadataTags:    aws.String(string(infrav1.InstanceMetadataEndpointStateDisabled)),
						},
					},
				},
			},
		},
	}, nil)

	// Must attach to a classic LB, not another type. Only these mock calls are therefore expected.
	mockedCreateLBCalls(t, elbMock.EXPECT(), false)

	ec2Mock.EXPECT().DescribeNetworkInterfacesWithContext(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfacesInput{Filters: []*ec2.Filter{
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
	ec2Mock.EXPECT().DescribeNetworkInterfaceAttributeWithContext(context.TODO(), gomock.Eq(&ec2.DescribeNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String("eni-1"),
		Attribute:          aws.String("groupSet"),
	})).Return(&ec2.DescribeNetworkInterfaceAttributeOutput{Groups: []*ec2.GroupIdentifier{{GroupId: aws.String("3")}}}, nil).MaxTimes(1)
	ec2Mock.EXPECT().ModifyNetworkInterfaceAttributeWithContext(context.TODO(), gomock.Any()).AnyTimes()
	ec2Mock.EXPECT().AssociateAddressWithContext(context.TODO(), gomock.Any()).MaxTimes(1)

	_, err = reconciler.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: awsMachine.Namespace,
			Name:      awsMachine.Name,
		},
	})

	g.Expect(err).To(BeNil())
}

func createObject(g *WithT, obj client.Object, namespace string) {
	if obj.DeepCopyObject() != nil {
		obj.SetNamespace(namespace)
		g.Expect(testEnv.Create(ctx, obj)).To(Succeed())
	}
}

func cleanupObject(g *WithT, obj client.Object) {
	if obj.DeepCopyObject() != nil {
		g.Expect(testEnv.Cleanup(ctx, obj)).To(Succeed())
	}
}
