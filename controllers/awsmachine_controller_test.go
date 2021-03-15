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

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/mock_services"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("AWSMachineReconciler", func() {
	var (
		reconciler AWSMachineReconciler
		cs         *scope.ClusterScope
		ms         *scope.MachineScope
		mockCtrl   *gomock.Controller
		ec2Svc     *mock_services.MockEC2MachineInterface
		secretSvc  *mock_services.MockSecretInterface
		recorder   *record.FakeRecorder
	)

	BeforeEach(func() {
		// https://github.com/kubernetes/klog/issues/87#issuecomment-540153080
		// TODO: Replace with LogToOutput when https://github.com/kubernetes/klog/pull/99 merges
		var err error

		if err := flag.Set("logtostderr", "false"); err != nil {
			_ = fmt.Errorf("Error setting logtostderr flag")
		}
		if err := flag.Set("v", "2"); err != nil {
			_ = fmt.Errorf("Error setting v flag")
		}
		klog.SetOutput(GinkgoWriter)

		awsMachine := &infrav1.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: infrav1.AWSMachineSpec{
				CloudInit: infrav1.CloudInit{
					SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				},
			},
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bootstrap-data",
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}

		ms, err = scope.NewMachineScope(
			scope.MachineScopeParams{
				Client: fake.NewFakeClient([]runtime.Object{awsMachine, secret}...),
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
		Expect(err).To(BeNil())

		cs, err = scope.NewClusterScope(
			scope.ClusterScopeParams{
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			},
		)
		Expect(err).To(BeNil())

		client := fake.NewFakeClient()
		ctx := context.TODO()
		client.Create(ctx, awsMachine)
		client.Create(ctx, secret)
		client.Create(ctx, secret)

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
		Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(GinkgoT())
		ec2Svc = mock_services.NewMockEC2MachineInterface(mockCtrl)
		secretSvc = mock_services.NewMockSecretInterface(mockCtrl)

		// If your test hangs for 9 minutes, increase the value here to the number of events during a reconciliation loop
		recorder = record.NewFakeRecorder(2)

		reconciler = AWSMachineReconciler{
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2MachineInterface {
				return ec2Svc
			},
			secretsManagerServiceFactory: func(cloud.ClusterScoper) services.SecretInterface {
				return secretSvc
			},
			Recorder: recorder,
			Log:      klogr.New(),
		}
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Reconciling an AWSMachine", func() {
		When("we can't reach amazon", func() {
			expectedErr := errors.New("no connection available ")

			BeforeEach(func() {
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, expectedErr).AnyTimes()
			})

			It("should exit immediately on an error state", func() {
				er := capierrors.CreateMachineError
				ms.AWSMachine.Status.FailureReason = &er
				ms.AWSMachine.Status.FailureMessage = pointer.StringPtr("Couldn't create machine")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
			})

			It("should add our finalizer to the machine", func() {
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)

				Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
			})

			It("should exit immediately if cluster infra isn't ready", func() {
				ms.Cluster.Status.InfrastructureReady = false

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(err).To(BeNil())
				Expect(buf.String()).To(ContainSubstring("Cluster infrastructure is not ready yet"))
				expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForClusterInfrastructureReason}})
			})

			It("should exit immediately if bootstrap data secret reference isn't available", func() {
				ms.Machine.Spec.Bootstrap.DataSecretName = nil

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)

				Expect(err).To(BeNil())
				Expect(buf.String()).To(ContainSubstring("Bootstrap data secret reference is not yet available"))
				expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForBootstrapDataReason}})
			})

			It("should return an error when we can't list instances by tags", func() {
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})

		When("there's a provider ID", func() {
			id := "aws:////myMachine"
			BeforeEach(func() {
				_, err := noderefutil.NewProviderID(id)
				Expect(err).To(BeNil())

				ms.AWSMachine.Spec.ProviderID = &id
			})

			It("it should look up by provider ID when one exists", func() {
				expectedErr := errors.New("no connection available ")
				ec2Svc.EXPECT().InstanceIfExists(PointsTo("myMachine")).Return(nil, expectedErr)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})

			It("should try to create a new machine if none exists", func() {
				expectedErr := errors.New("Invalid instance")
				ec2Svc.EXPECT().InstanceIfExists(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(nil, expectedErr)
				secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})

		When("instance creation succeeds", func() {
			var instance *infrav1.Instance
			BeforeEach(func() {
				instance = &infrav1.Instance{
					ID: "myMachine",
				}
				instance.State = infrav1.InstanceStatePending

				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test", int32(1), nil).Times(1)
				ec2Svc.EXPECT().CreateInstance(gomock.Any(), gomock.Any()).Return(instance, nil)
			})

			Context("instance security group errors", func() {
				BeforeEach(func() {
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).Return(nil, errors.New("stop here"))
				})

				It("should set attributes after creating an instance", func() {
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Spec.ProviderID).To(PointTo(Equal("aws:////myMachine")))
				})

				Context("with captured logging", func() {
					var buf *bytes.Buffer

					BeforeEach(func() {
						buf = new(bytes.Buffer)
						klog.SetOutput(buf)
					})

					It("should set instance to pending", func() {
						secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
						instance.State = infrav1.InstanceStatePending
						_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
						Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStatePending)))
						Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
						Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
						expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.InstanceNotReadyReason}})
					})

					It("should set instance to running", func() {
						secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
						instance.State = infrav1.InstanceStateRunning
						_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
						Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
						Expect(ms.AWSMachine.Status.Ready).To(Equal(true))
						Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
						expectConditions(ms.AWSMachine, []conditionAssertion{
							{conditionType: infrav1.InstanceReadyCondition, status: corev1.ConditionTrue},
						})
					})
				})
			})

			Context("New EC2 instance state", func() {
				It("should error when the instance state is a new unseen one", func() {
					buf := new(bytes.Buffer)
					klog.SetOutput(buf)
					instance.State = "NewAWSMachineState"
					secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					Expect(buf.String()).To(ContainSubstring(("EC2 instance state is undefined")))
					Eventually(recorder.Events).Should(Receive(ContainSubstring("InstanceUnhandledState")))
					Expect(ms.AWSMachine.Status.FailureMessage).To(PointTo(Equal("EC2 instance state \"NewAWSMachineState\" is undefined")))
					expectConditions(ms.AWSMachine, []conditionAssertion{{conditionType: infrav1.InstanceReadyCondition, status: corev1.ConditionUnknown}})
				})
			})

			Context("Security Groups succeed", func() {
				BeforeEach(func() {
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
						Return(map[string][]string{"eid": {}}, nil)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil)
				})

				It("should reconcile security groups", func() {
					ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1.AWSResourceReference{
						{
							ID: pointer.StringPtr("sg-2345"),
						},
					}
					// ms.AWSMachine.Spec.AdditionalSecurityGroups = []infrav1
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(instance.ID, []string{"sg-2345"})

					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					expectConditions(ms.AWSMachine, []conditionAssertion{{conditionType: infrav1.SecurityGroupsReadyCondition, status: corev1.ConditionTrue}})
				})

				It("should not tag anything if there's not tags", func() {
					ec2Svc.EXPECT().UpdateInstanceSecurityGroups(gomock.Any(), gomock.Any()).Times(0)
					if _, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs); err != nil {
						_ = fmt.Errorf("reconcileNormal reutrned an error during test")
					}
				})

				It("should tag instances from machine and cluster tags", func() {
					ms.AWSMachine.Spec.AdditionalTags = infrav1.Tags{"kind": "alicorn"}
					cs.AWSCluster.Spec.AdditionalTags = infrav1.Tags{"colour": "lavender"}

					ec2Svc.EXPECT().UpdateResourceTags(
						PointsTo("myMachine"),
						map[string]string{
							"kind":   "alicorn",
							"colour": "lavender",
						},
						map[string]string{},
					).Return(nil)

					_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(err).To(BeNil())
				})
			})

			When("temporarily stopping then starting the AWSMachine", func() {
				var buf *bytes.Buffer
				BeforeEach(func() {
					buf = new(bytes.Buffer)
					klog.SetOutput(buf)
					ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
						Return(map[string][]string{"eid": {}}, nil).Times(1)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
					ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				})

				It("should set instance to stopping and unready", func() {
					instance.State = infrav1.InstanceStateStopping
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopping)))
					Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
					expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceStoppedReason}})
				})

				It("should then set instance to stopped and unready", func() {
					instance.State = infrav1.InstanceStateStopped
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateStopped)))
					Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
					expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceStoppedReason}})
				})

				It("should then set instance to running and ready once it is restarted", func() {
					instance.State = infrav1.InstanceStateRunning
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.InstanceState).To(PointTo(Equal(infrav1.InstanceStateRunning)))
					Expect(ms.AWSMachine.Status.Ready).To(Equal(true))
					Expect(buf.String()).To(ContainSubstring(("EC2 instance state changed")))
				})
			})

			When("deleting the AWSMachine outside of Kubernetes", func() {
				var buf *bytes.Buffer
				BeforeEach(func() {
					buf = new(bytes.Buffer)
					klog.SetOutput(buf)
					secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
					secretSvc.EXPECT().UserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				})

				It("should warn if an instance is shutting-down", func() {
					instance.State = infrav1.InstanceStateShuttingDown
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					Expect(buf.String()).To(ContainSubstring(("Unexpected EC2 instance termination")))
					Eventually(recorder.Events).Should(Receive(ContainSubstring("UnexpectedTermination")))
				})

				It("should error when the instance is seen as terminated", func() {
					instance.State = infrav1.InstanceStateTerminated
					_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
					Expect(ms.AWSMachine.Status.Ready).To(Equal(false))
					Expect(buf.String()).To(ContainSubstring(("Unexpected EC2 instance termination")))
					Eventually(recorder.Events).Should(Receive(ContainSubstring("UnexpectedTermination")))
					Expect(ms.AWSMachine.Status.FailureMessage).To(PointTo(Equal("EC2 instance state \"terminated\" is unexpected")))
					expectConditions(ms.AWSMachine, []conditionAssertion{{infrav1.InstanceReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityError, infrav1.InstanceTerminatedReason}})
				})
			})
		})
	})

	Context("secrets management lifecycle", func() {
		var instance *infrav1.Instance
		secretPrefix := "test/secret"
		When("creating EC2 instances", func() {
			It("should leverage AWS Secrets Manager", func() {
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

		When("there's a node ref and a secret ARN", func() {
			BeforeEach(func() {
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
			})

			It("should delete the secret if the instance is running", func() {
				instance.State = infrav1.InstanceStateRunning
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
					Return(map[string][]string{"eid": {}}, nil).Times(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			It("should delete the secret if the instance is terminated", func() {
				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			It("should delete the secret if the AWSMachine is deleted", func() {
				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})

			It("should delete the secret if the AWSMachine is in a failure condition", func() {
				ms.AWSMachine.Status.FailureReason = capierrors.MachineStatusErrorPtr(capierrors.UpdateMachineError)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})
		})

		When("there's only a secret ARN and no node ref", func() {
			BeforeEach(func() {
				instance = &infrav1.Instance{
					ID: "myMachine",
				}
				ms.AWSMachine.Spec.CloudInit = infrav1.CloudInit{
					SecretPrefix:         "secret",
					SecretCount:          5,
					SecureSecretsBackend: infrav1.SecretBackendSecretsManager,
				}
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(instance, nil).AnyTimes()
			})

			It("should not delete the secret if the instance is running", func() {
				instance.State = infrav1.InstanceStateRunning
				ec2Svc.EXPECT().GetInstanceSecurityGroups(gomock.Any()).
					Return(map[string][]string{"eid": {}}, nil).Times(1)
				ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{}, nil).Times(1)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).MaxTimes(0)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			It("should delete the secret if the instance is terminated", func() {
				instance.State = infrav1.InstanceStateTerminated
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
			})

			It("should delete the secret if the AWSMachine is deleted", func() {
				instance.State = infrav1.InstanceStateRunning
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})

			It("should delete the secret if the AWSMachine is in a failure condition", func() {
				ms.AWSMachine.Status.FailureReason = capierrors.MachineStatusErrorPtr(capierrors.UpdateMachineError)
				secretSvc.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil).AnyTimes()
				_, _ = reconciler.reconcileDelete(ms, cs, cs, cs)
			})
		})

		When("there is an intermittent connection issue and no secret could be stored", func() {
			BeforeEach(func() {
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil).AnyTimes()
			})
			It("should error if secret could not be created", func() {
				secretSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(secretPrefix, int32(0), errors.New("connection error")).Times(1)
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs, cs)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("connection error"))
				Expect(ms.GetSecretPrefix()).To(Equal(""))
				Expect(ms.GetSecretCount()).To(Equal(int32(0)))
			})

			It("should update prefix and count on successful creation", func() {
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

				Expect(err).To(BeNil())
				Expect(ms.GetSecretPrefix()).To(Equal(secretPrefix))
				Expect(ms.GetSecretCount()).To(Equal(int32(1)))
			})
		})
	})

	Context("deleting an AWSMachine", func() {
		BeforeEach(func() {
			ms.AWSMachine.Finalizers = []string{
				infrav1.MachineFinalizer,
				metav1.FinalizerDeleteDependents,
			}
		})

		It("should exit immediately on an error state", func() {
			expectedErr := errors.New("no connection available ")
			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, expectedErr).AnyTimes()

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			Expect(errors.Cause(err)).To(MatchError(expectedErr))
		})

		It("should log and remove finalizer when no machine exists", func() {
			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(nil, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("Unable to locate EC2 instance by ID or tags"))
			Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
			Eventually(recorder.Events).Should(Receive(ContainSubstring("NoInstanceFound")))
		})

		It("should ignore instances in shutting down state", func() {
			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
				State: infrav1.InstanceStateShuttingDown,
			}, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("EC2 instance is shutting down or already terminated"))
			Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
		})

		It("should ignore instances in terminated down state", func() {
			ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{
				State: infrav1.InstanceStateTerminated,
			}, nil)

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("EC2 instance is shutting down or already terminated"))
			Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
		})

		Context("Instance not shutting down yet", func() {
			id := "aws:////myid"

			BeforeEach(func() {
				ec2Svc.EXPECT().GetRunningInstanceByTags(gomock.Any()).Return(&infrav1.Instance{ID: id}, nil)
			})

			It("should return an error when the instance can't be terminated", func() {
				expected := errors.New("can't reach AWS to terminate machine")
				ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(expected)

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
				Expect(errors.Cause(err)).To(MatchError(expected))
				Expect(buf.String()).To(ContainSubstring("Terminating EC2 instance"))
				Eventually(recorder.Events).Should(Receive(ContainSubstring("FailedTerminate")))
			})

			When("instance can be shut down", func() {
				BeforeEach(func() {
					ec2Svc.EXPECT().TerminateInstanceAndWait(gomock.Any()).Return(nil)
				})

				When("there are network interfaces", func() {
					BeforeEach(func() {
						ms.AWSMachine.Spec.NetworkInterfaces = []string{
							"eth0",
							"eth1",
						}
					})

					It("should error when it can't retrieve security groups", func() {
						expected := errors.New("can't reach AWS to list security groups")
						ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return(nil, expected)

						_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
						Expect(errors.Cause(err)).To(MatchError(expected))
					})

					It("should error when it can't detach a security group from an interface", func() {
						expected := errors.New("can't reach AWS to detach security group")
						ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{"sg0", "sg1"}, nil)
						ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(gomock.Any(), gomock.Any()).Return(expected)

						_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
						Expect(errors.Cause(err)).To(MatchError(expected))
					})

					It("should detach all combinations of network interfaces", func() {
						groups := []string{"sg0", "sg1"}
						ec2Svc.EXPECT().GetCoreSecurityGroups(gomock.Any()).Return([]string{"sg0", "sg1"}, nil)
						ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(groups, "eth0").Return(nil)
						ec2Svc.EXPECT().DetachSecurityGroupsFromNetworkInterface(groups, "eth1").Return(nil)

						_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
						Expect(err).To(BeNil())
					})
				})

				It("should remove security groups", func() {
					_, err := reconciler.reconcileDelete(ms, cs, cs, cs)
					Expect(err).To(BeNil())
					Expect(ms.AWSMachine.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
				})
			})
		})
	})
})

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

func expectConditions(m *infrav1.AWSMachine, expected []conditionAssertion) {
	Expect(len(m.Status.Conditions)).To(BeNumerically(">=", len(expected)), "number of conditions")
	for _, c := range expected {
		actual := conditions.Get(m, c.conditionType)
		Expect(actual).To(Not(BeNil()))
		Expect(actual.Type).To(Equal(c.conditionType))
		Expect(actual.Status).To(Equal(c.status))
		Expect(actual.Severity).To(Equal(c.severity))
		Expect(actual.Reason).To(Equal(c.reason))
	}
}
