/*
Copyright 2020 The Kubernetes Authors.

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
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

func TestAWSMachinePoolReconciler(t *testing.T) {
	var (
		reconciler        AWSMachinePoolReconciler
		cs                *scope.ClusterScope
		ms                *scope.MachinePoolScope
		mockCtrl          *gomock.Controller
		ec2Svc            *mock_services.MockEC2Interface
		asgSvc            *mock_services.MockASGInterface
		reconSvc          *mock_services.MockMachinePoolReconcileInterface
		recorder          *record.FakeRecorder
		awsMachinePool    *expinfrav1.AWSMachinePool
		secret            *corev1.Secret
		userDataSecretKey apimachinerytypes.NamespacedName
	)
	setup := func(t *testing.T, g *WithT) {
		t.Helper()

		var err error

		if err := flag.Set("logtostderr", "false"); err != nil {
			_ = fmt.Errorf("Error setting logtostderr flag")
		}
		if err := flag.Set("v", "2"); err != nil {
			_ = fmt.Errorf("Error setting v flag")
		}
		ctx := context.TODO()

		awsMachinePool = &expinfrav1.AWSMachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: expinfrav1.AWSMachinePoolSpec{
				MinSize: int32(0),
				MaxSize: int32(100),
				MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
					InstancesDistribution: &expinfrav1.InstancesDistribution{
						OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
						SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyCapacityOptimized,
						OnDemandBaseCapacity:                aws.Int64(0),
						OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
					},
					Overrides: []expinfrav1.Overrides{
						{
							InstanceType: "m6a.32xlarge",
						},
					},
				},
			},
		}

		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bootstrap-data",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"value": []byte("shell-script"),
			},
		}
		userDataSecretKey = apimachinerytypes.NamespacedName{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}

		g.Expect(testEnv.Create(ctx, awsMachinePool)).To(Succeed())
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		cs, err = setupCluster("test-cluster")
		g.Expect(err).To(BeNil())

		ms, err = scope.NewMachinePoolScope(
			scope.MachinePoolScopeParams{
				Client: testEnv.Client,
				Cluster: &clusterv1.Cluster{
					Status: clusterv1.ClusterStatus{
						InfrastructureReady: true,
					},
				},
				MachinePool: &expclusterv1.MachinePool{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mp",
						Namespace: "default",
					},
					Spec: expclusterv1.MachinePoolSpec{
						ClusterName: "test",
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								ClusterName: "test",
								Bootstrap: clusterv1.Bootstrap{
									DataSecretName: ptr.To[string]("bootstrap-data"),
								},
							},
						},
					},
				},
				InfraCluster:   cs,
				AWSMachinePool: awsMachinePool,
			},
		)
		g.Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(t)
		ec2Svc = mock_services.NewMockEC2Interface(mockCtrl)
		asgSvc = mock_services.NewMockASGInterface(mockCtrl)
		reconSvc = mock_services.NewMockMachinePoolReconcileInterface(mockCtrl)

		// If the test hangs for 9 minutes, increase the value here to the number of events during a reconciliation loop
		recorder = record.NewFakeRecorder(2)

		reconciler = AWSMachinePoolReconciler{
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2Interface {
				return ec2Svc
			},
			asgServiceFactory: func(cloud.ClusterScoper) services.ASGInterface {
				return asgSvc
			},
			reconcileServiceFactory: func(scope.EC2Scope) services.MachinePoolReconcileInterface {
				return reconSvc
			},
			Recorder: recorder,
		}
	}

	teardown := func(t *testing.T, g *WithT) {
		t.Helper()

		ctx := context.TODO()
		mpPh, err := patch.NewHelper(awsMachinePool, testEnv)
		g.Expect(err).ShouldNot(HaveOccurred())
		awsMachinePool.SetFinalizers([]string{})
		g.Expect(mpPh.Patch(ctx, awsMachinePool)).To(Succeed())
		g.Expect(testEnv.Delete(ctx, awsMachinePool)).To(Succeed())
		g.Expect(testEnv.Delete(ctx, secret)).To(Succeed())
		mockCtrl.Finish()
	}

	t.Run("Reconciling an AWSMachinePool", func(t *testing.T) {
		t.Run("when can't reach amazon", func(t *testing.T) {
			expectedErr := errors.New("no connection available ")
			getASG := func(t *testing.T, g *WithT) {
				t.Helper()

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil, expectedErr).AnyTimes()
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, expectedErr).AnyTimes()
			}
			t.Run("should exit immediately on an error state", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				er := capierrors.CreateMachineError
				ms.AWSMachinePool.Status.FailureReason = &er
				ms.AWSMachinePool.Status.FailureMessage = ptr.To[string]("Couldn't create machine pool")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_ = reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
			})
			t.Run("should add our finalizer to the machinepool", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				_ = reconciler.reconcileNormal(context.Background(), ms, cs, cs)

				g.Expect(ms.AWSMachinePool.Finalizers).To(ContainElement(expinfrav1.MachinePoolFinalizer))
			})
			t.Run("should exit immediately if cluster infra isn't ready", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				ms.Cluster.Status.InfrastructureReady = false

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(BeNil())
				g.Expect(buf.String()).To(ContainSubstring("Cluster infrastructure is not ready yet"))
				expectConditions(g, ms.AWSMachinePool, []conditionAssertion{{expinfrav1.ASGReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForClusterInfrastructureReason}})
			})
			t.Run("should exit immediately if bootstrap data secret reference isn't available", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				ms.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName = nil
				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)

				g.Expect(err).To(BeNil())
				g.Expect(buf.String()).To(ContainSubstring("Bootstrap data secret reference is not yet available"))
				expectConditions(g, ms.AWSMachinePool, []conditionAssertion{{expinfrav1.ASGReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForBootstrapDataReason}})
			})
		})
		t.Run("there's a provider ID", func(t *testing.T) {
			id := "<cloudProvider>://<optional>/<segments>/<providerid>"
			setProviderID := func(t *testing.T, g *WithT) {
				t.Helper()
				ms.AWSMachinePool.Spec.ProviderID = id
			}
			getASG := func(t *testing.T, g *WithT) {
				t.Helper()

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil, nil).AnyTimes()
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil).AnyTimes()
			}
			t.Run("should look up by provider ID when one exists", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setProviderID(t, g)
				getASG(t, g)

				expectedErr := errors.New("no connection available ")
				reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)
				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})
		t.Run("there's suspended processes provided during ASG creation", func(t *testing.T) {
			setSuspendedProcesses := func(t *testing.T, g *WithT) {
				t.Helper()
				ms.AWSMachinePool.Spec.SuspendProcesses = &expinfrav1.SuspendProcessesTypes{
					Processes: &expinfrav1.Processes{
						Launch:    ptr.To[bool](true),
						Terminate: ptr.To[bool](true),
					},
				}
			}
			t.Run("it should not call suspend as we don't have an ASG yet", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setSuspendedProcesses(t, g)

				reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil)
				asgSvc.EXPECT().CreateASG(gomock.Any()).Return(&expinfrav1.AutoScalingGroup{
					Name: "name",
				}, nil)
				asgSvc.EXPECT().SuspendProcesses("name", []string{"Launch", "Terminate"}).Return(nil).AnyTimes().Times(0)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})
		})
		t.Run("all processes are suspended", func(t *testing.T) {
			setSuspendedProcesses := func(t *testing.T, g *WithT) {
				t.Helper()
				ms.AWSMachinePool.Spec.SuspendProcesses = &expinfrav1.SuspendProcessesTypes{
					All: true,
				}
			}
			t.Run("processes should be suspended during an update call", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setSuspendedProcesses(t, g)
				ms.AWSMachinePool.Spec.SuspendProcesses.All = true
				reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&expinfrav1.AutoScalingGroup{
					Name: "name",
				}, nil)
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{}, nil).Times(1)
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil).AnyTimes()
				asgSvc.EXPECT().SuspendProcesses("name", gomock.InAnyOrder([]string{
					"ScheduledActions",
					"Launch",
					"Terminate",
					"AddToLoadBalancer",
					"AlarmNotification",
					"AZRebalance",
					"InstanceRefresh",
					"HealthCheck",
					"ReplaceUnhealthy",
				})).Return(nil).AnyTimes().Times(1)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})
		})
		t.Run("there are existing processes already suspended", func(t *testing.T) {
			setSuspendedProcesses := func(t *testing.T, g *WithT) {
				t.Helper()

				ms.AWSMachinePool.Spec.SuspendProcesses = &expinfrav1.SuspendProcessesTypes{
					Processes: &expinfrav1.Processes{
						Launch:    ptr.To[bool](true),
						Terminate: ptr.To[bool](true),
					},
				}
			}
			t.Run("it should suspend and resume processes that are desired to be suspended and desired to be resumed", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setSuspendedProcesses(t, g)

				reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&expinfrav1.AutoScalingGroup{
					Name:                      "name",
					CurrentlySuspendProcesses: []string{"Launch", "process3"},
				}, nil)
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{}, nil).Times(1)
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil).AnyTimes()
				asgSvc.EXPECT().SuspendProcesses("name", []string{"Terminate"}).Return(nil).AnyTimes().Times(1)
				asgSvc.EXPECT().ResumeProcesses("name", []string{"process3"}).Return(nil).AnyTimes().Times(1)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})
		})

		t.Run("externally managed annotation", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)

			asg := expinfrav1.AutoScalingGroup{
				Name:            "an-asg",
				DesiredCapacity: ptr.To[int32](1),
			}
			reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&asg, nil)
			asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{}, nil)
			asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil)
			reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)

			ms.MachinePool.Annotations = map[string]string{
				clusterv1.ReplicasManagedByAnnotation: "somehow-externally-managed",
			}
			ms.MachinePool.Spec.Replicas = ptr.To[int32](0)

			g.Expect(testEnv.Create(ctx, ms.MachinePool)).To(Succeed())

			_ = reconciler.reconcileNormal(context.Background(), ms, cs, cs)
			g.Expect(*ms.MachinePool.Spec.Replicas).To(Equal(int32(1)))
		})
		t.Run("No need to update Asg because asgNeedsUpdates is false and no subnets change", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)

			asg := expinfrav1.AutoScalingGroup{
				MinSize: int32(0),
				MaxSize: int32(100),
				MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
					InstancesDistribution: &expinfrav1.InstancesDistribution{
						OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
						SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyCapacityOptimized,
						OnDemandBaseCapacity:                aws.Int64(0),
						OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
					},
					Overrides: []expinfrav1.Overrides{
						{
							InstanceType: "m6a.32xlarge",
						},
					},
				},
				Subnets: []string{"subnet1", "subnet2"}}
			reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&asg, nil).AnyTimes()
			asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet2", "subnet1"}, nil).Times(1)
			asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil).Times(0)

			err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
			g.Expect(err).To(Succeed())
		})
		t.Run("update Asg due to subnet changes", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)

			asg := expinfrav1.AutoScalingGroup{
				MinSize: int32(0),
				MaxSize: int32(100),
				Subnets: []string{"subnet1", "subnet2"}}
			reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&asg, nil).AnyTimes()
			asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet1"}, nil).Times(1)
			asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil).Times(1)

			err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
			g.Expect(err).To(Succeed())
		})
		t.Run("update Asg due to asgNeedsUpdates returns true", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)

			asg := expinfrav1.AutoScalingGroup{
				MinSize: int32(0),
				MaxSize: int32(2),
				Subnets: []string{}}
			reconSvc.EXPECT().ReconcileLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			reconSvc.EXPECT().ReconcileTags(gomock.Any(), gomock.Any()).Return(nil)
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&asg, nil).AnyTimes()
			asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{}, nil).Times(1)
			asgSvc.EXPECT().UpdateASG(gomock.Any()).Return(nil).Times(1)

			err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
			g.Expect(err).To(Succeed())
		})

		t.Run("ReconcileLaunchTemplate not mocked", func(t *testing.T) {
			launchTemplateIDExisting := "lt-existing"

			t.Run("nothing exists, so launch template and ASG must be created", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				reconciler.reconcileServiceFactory = nil // use real implementation, but keep EC2 calls mocked (`ec2ServiceFactory`)
				reconSvc = nil                           // not used
				defer teardown(t, g)

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(nil, "", nil, nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-abcdef123"), nil)
				ec2Svc.EXPECT().CreateLaunchTemplate(gomock.Any(), gomock.Eq(ptr.To[string]("ami-abcdef123")), gomock.Eq(userDataSecretKey), gomock.Eq([]byte("shell-script"))).Return("lt-ghijkl456", nil)
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil)
				asgSvc.EXPECT().CreateASG(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
					}, nil
				})

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})

			t.Run("launch template and ASG exist and need no update", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				reconciler.reconcileServiceFactory = nil // use real implementation, but keep EC2 calls mocked (`ec2ServiceFactory`)
				reconSvc = nil                           // not used
				defer teardown(t, g)

				// Latest ID and version already stored, no need to retrieve it
				ms.AWSMachinePool.Status.LaunchTemplateID = launchTemplateIDExisting
				ms.AWSMachinePool.Status.LaunchTemplateVersion = ptr.To[string]("1")

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(
					&expinfrav1.AWSLaunchTemplate{
						Name: "test",
						AMI: infrav1.AMIReference{
							ID: ptr.To[string]("ami-existing"),
						},
					},
					// No change to user data
					userdata.ComputeHash([]byte("shell-script")),
					&userDataSecretKey,
					nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-existing"), nil) // no change
				ec2Svc.EXPECT().LaunchTemplateNeedsUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)

				asgSvc.EXPECT().GetASGByName(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))

					// No difference to `AWSMachinePool.spec`
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
						Subnets: []string{
							"subnet-1",
						},
						MinSize:              awsMachinePool.Spec.MinSize,
						MaxSize:              awsMachinePool.Spec.MaxSize,
						MixedInstancesPolicy: awsMachinePool.Spec.MixedInstancesPolicy.DeepCopy(),
					}, nil
				})
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet-1"}, nil) // no change
				// No changes, so there must not be an ASG update!
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Times(0)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})

			t.Run("launch template and ASG exist and only AMI ID changed", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				reconciler.reconcileServiceFactory = nil // use real implementation, but keep EC2 calls mocked (`ec2ServiceFactory`)
				reconSvc = nil                           // not used
				defer teardown(t, g)

				// Latest ID and version already stored, no need to retrieve it
				ms.AWSMachinePool.Status.LaunchTemplateID = launchTemplateIDExisting
				ms.AWSMachinePool.Status.LaunchTemplateVersion = ptr.To[string]("1")

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(
					&expinfrav1.AWSLaunchTemplate{
						Name: "test",
						AMI: infrav1.AMIReference{
							ID: ptr.To[string]("ami-existing"),
						},
					},
					// No change to user data
					userdata.ComputeHash([]byte("shell-script")),
					&userDataSecretKey,
					nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-different"), nil)
				ec2Svc.EXPECT().LaunchTemplateNeedsUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				asgSvc.EXPECT().CanStartASGInstanceRefresh(gomock.Any()).Return(true, nil)
				ec2Svc.EXPECT().PruneLaunchTemplateVersions(gomock.Any()).Return(nil)
				ec2Svc.EXPECT().CreateLaunchTemplateVersion(gomock.Any(), gomock.Any(), gomock.Eq(ptr.To[string]("ami-different")), gomock.Eq(apimachinerytypes.NamespacedName{Namespace: "default", Name: "bootstrap-data"}), gomock.Any()).Return(nil)
				ec2Svc.EXPECT().GetLaunchTemplateLatestVersion(gomock.Any()).Return("2", nil)
				// AMI change should trigger rolling out new nodes
				asgSvc.EXPECT().StartASGInstanceRefresh(gomock.Any())

				asgSvc.EXPECT().GetASGByName(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))

					// No difference to `AWSMachinePool.spec`
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
						Subnets: []string{
							"subnet-1",
						},
						MinSize:              awsMachinePool.Spec.MinSize,
						MaxSize:              awsMachinePool.Spec.MaxSize,
						MixedInstancesPolicy: awsMachinePool.Spec.MixedInstancesPolicy.DeepCopy(),
					}, nil
				})
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet-1"}, nil) // no change
				// No changes, so there must not be an ASG update!
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Times(0)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})

			t.Run("launch template and ASG exist and only bootstrap data secret name changed", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				reconciler.reconcileServiceFactory = nil // use real implementation, but keep EC2 calls mocked (`ec2ServiceFactory`)
				reconSvc = nil                           // not used
				defer teardown(t, g)

				// Latest ID and version already stored, no need to retrieve it
				ms.AWSMachinePool.Status.LaunchTemplateID = launchTemplateIDExisting
				ms.AWSMachinePool.Status.LaunchTemplateVersion = ptr.To[string]("1")

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(
					&expinfrav1.AWSLaunchTemplate{
						Name: "test",
						AMI: infrav1.AMIReference{
							ID: ptr.To[string]("ami-existing"),
						},
					},
					// No change to user data
					userdata.ComputeHash([]byte("shell-script")),
					// But the name of the secret changes from `previous-secret-name` to `bootstrap-data`
					&apimachinerytypes.NamespacedName{Namespace: "default", Name: "previous-secret-name"},
					nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-existing"), nil)
				ec2Svc.EXPECT().LaunchTemplateNeedsUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				asgSvc.EXPECT().CanStartASGInstanceRefresh(gomock.Any()).Return(true, nil)
				ec2Svc.EXPECT().PruneLaunchTemplateVersions(gomock.Any()).Return(nil)
				ec2Svc.EXPECT().CreateLaunchTemplateVersion(gomock.Any(), gomock.Any(), gomock.Eq(ptr.To[string]("ami-existing")), gomock.Eq(apimachinerytypes.NamespacedName{Namespace: "default", Name: "bootstrap-data"}), gomock.Any()).Return(nil)
				ec2Svc.EXPECT().GetLaunchTemplateLatestVersion(gomock.Any()).Return("2", nil)
				// Changing the bootstrap data secret name should trigger rolling out new nodes, no matter what the
				// content (user data) is. This way, users can enforce a rollout by changing the bootstrap config
				// reference (`MachinePool.spec.template.spec.bootstrap`).
				asgSvc.EXPECT().StartASGInstanceRefresh(gomock.Any())

				asgSvc.EXPECT().GetASGByName(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))

					// No difference to `AWSMachinePool.spec`
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
						Subnets: []string{
							"subnet-1",
						},
						MinSize:              awsMachinePool.Spec.MinSize,
						MaxSize:              awsMachinePool.Spec.MaxSize,
						MixedInstancesPolicy: awsMachinePool.Spec.MixedInstancesPolicy.DeepCopy(),
					}, nil
				})
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet-1"}, nil) // no change
				// No changes, so there must not be an ASG update!
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Times(0)

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})

			t.Run("launch template and ASG created from zero, then bootstrap config reference changes", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				reconciler.reconcileServiceFactory = nil // use real implementation, but keep EC2 calls mocked (`ec2ServiceFactory`)
				reconSvc = nil                           // not used
				defer teardown(t, g)

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(nil, "", nil, nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-abcdef123"), nil)
				ec2Svc.EXPECT().CreateLaunchTemplate(gomock.Any(), gomock.Eq(ptr.To[string]("ami-abcdef123")), gomock.Eq(userDataSecretKey), gomock.Eq([]byte("shell-script"))).Return("lt-ghijkl456", nil)
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil)
				asgSvc.EXPECT().CreateASG(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
					}, nil
				})

				err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())

				g.Expect(ms.AWSMachinePool.Status.LaunchTemplateID).ToNot(BeEmpty())

				// Data secret name changes
				newBootstrapSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bootstrap-data-new", // changed
						Namespace: "default",
					},
					Data: map[string][]byte{
						"value": secret.Data["value"], // not changed
					},
				}
				g.Expect(testEnv.Create(ctx, newBootstrapSecret)).To(Succeed())
				ms.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName = ptr.To[string](newBootstrapSecret.Name)

				// Since `AWSMachinePool.status.launchTemplateVersion` isn't set yet,
				// the controller will ask for the current version and then set the status.
				ec2Svc.EXPECT().GetLaunchTemplateLatestVersion(gomock.Any()).Return("1", nil)

				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Eq("test")).Return(
					&expinfrav1.AWSLaunchTemplate{
						Name: "test",
						AMI: infrav1.AMIReference{
							ID: ptr.To[string]("ami-existing"),
						},
					},
					// No change to user data content
					userdata.ComputeHash([]byte("shell-script")),
					&apimachinerytypes.NamespacedName{Namespace: "default", Name: "bootstrap-data"},
					nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(ptr.To[string]("ami-existing"), nil)
				ec2Svc.EXPECT().LaunchTemplateNeedsUpdate(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				asgSvc.EXPECT().CanStartASGInstanceRefresh(gomock.Any()).Return(true, nil)
				ec2Svc.EXPECT().PruneLaunchTemplateVersions(gomock.Any()).Return(nil)
				ec2Svc.EXPECT().CreateLaunchTemplateVersion(gomock.Any(), gomock.Any(), gomock.Eq(ptr.To[string]("ami-existing")), gomock.Eq(apimachinerytypes.NamespacedName{Namespace: "default", Name: "bootstrap-data-new"}), gomock.Any()).Return(nil)
				ec2Svc.EXPECT().GetLaunchTemplateLatestVersion(gomock.Any()).Return("2", nil)
				// Changing the bootstrap data secret name should trigger rolling out new nodes, no matter what the
				// content (user data) is. This way, users can enforce a rollout by changing the bootstrap config
				// reference (`MachinePool.spec.template.spec.bootstrap.configRef`).
				asgSvc.EXPECT().StartASGInstanceRefresh(gomock.Any())

				asgSvc.EXPECT().GetASGByName(gomock.Any()).DoAndReturn(func(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
					g.Expect(scope.Name()).To(Equal("test"))

					// No difference to `AWSMachinePool.spec`
					return &expinfrav1.AutoScalingGroup{
						Name: scope.Name(),
						Subnets: []string{
							"subnet-1",
						},
						MinSize:              awsMachinePool.Spec.MinSize,
						MaxSize:              awsMachinePool.Spec.MaxSize,
						MixedInstancesPolicy: awsMachinePool.Spec.MixedInstancesPolicy.DeepCopy(),
					}, nil
				})
				asgSvc.EXPECT().SubnetIDs(gomock.Any()).Return([]string{"subnet-1"}, nil) // no change
				// No changes, so there must not be an ASG update!
				asgSvc.EXPECT().UpdateASG(gomock.Any()).Times(0)

				err = reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(err).To(Succeed())
			})
		})
	})

	t.Run("Deleting an AWSMachinePool", func(t *testing.T) {
		finalizer := func(t *testing.T, g *WithT) {
			t.Helper()

			ms.AWSMachinePool.Finalizers = []string{
				expinfrav1.MachinePoolFinalizer,
				metav1.FinalizerDeleteDependents,
			}
		}
		t.Run("should exit immediately on an error state", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)
			finalizer(t, g)

			expectedErr := errors.New("no connection available ")
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, expectedErr).AnyTimes()

			err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
		})
		t.Run("should log and remove finalizer when no machinepool exists", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)
			finalizer(t, g)

			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil)
			ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil, nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("Unable to locate ASG"))
			g.Expect(ms.AWSMachinePool.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
			g.Eventually(recorder.Events).Should(Receive(ContainSubstring(expinfrav1.ASGNotFoundReason)))
		})
		t.Run("should cause AWSMachinePool to go into NotReady", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)
			finalizer(t, g)

			inProgressASG := expinfrav1.AutoScalingGroup{
				Name:   "an-asg-that-is-currently-being-deleted",
				Status: expinfrav1.ASGStatusDeleteInProgress,
			}
			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(&inProgressASG, nil)
			ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil, nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)
			err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(ms.AWSMachinePool.Status.Ready).To(BeFalse())
			g.Eventually(recorder.Events).Should(Receive(ContainSubstring("DeletionInProgress")))
		})
	})
}

//TODO: This was taken from awsmachine_controller_test, i think it should be moved to elsewhere in both locations like test/helpers.

type conditionAssertion struct {
	conditionType clusterv1.ConditionType
	status        corev1.ConditionStatus
	severity      clusterv1.ConditionSeverity
	reason        string
}

func expectConditions(g *WithT, m *expinfrav1.AWSMachinePool, expected []conditionAssertion) {
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

func setupCluster(clusterName string) (*scope.ClusterScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec:       infrav1.AWSClusterSpec{},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).Build()
	return scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: clusterName},
		},
		AWSCluster: awsCluster,
		Client:     client,
	})
}

func TestDiffASG(t *testing.T) {
	type args struct {
		machinePoolScope *scope.MachinePoolScope
		existingASG      *expinfrav1.AutoScalingGroup
	}
	tests := []struct {
		name           string
		args           args
		wantDifference bool
	}{
		{
			name: "replicas != asg.desiredCapacity",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](0),
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
				},
			},
			wantDifference: true,
		},
		{
			name: "replicas (nil) != asg.desiredCapacity",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: nil,
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
				},
			},
			wantDifference: true,
		},
		{
			name: "replicas != asg.desiredCapacity (nil)",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](0),
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: nil,
				},
			},
			wantDifference: true,
		},
		{
			name: "maxSize != asg.maxSize",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize: 1,
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
					MaxSize:         2,
				},
			},
			wantDifference: true,
		},
		{
			name: "minSize != asg.minSize",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize: 2,
							MinSize: 0,
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
					MaxSize:         2,
					MinSize:         1,
				},
			},
			wantDifference: true,
		},
		{
			name: "capacityRebalance != asg.capacityRebalance",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:   ptr.To[int32](1),
					MaxSize:           2,
					MinSize:           0,
					CapacityRebalance: false,
				},
			},
			wantDifference: true,
		},
		{
			name: "MixedInstancesPolicy != asg.MixedInstancesPolicy",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
							MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
								InstancesDistribution: &expinfrav1.InstancesDistribution{
									OnDemandAllocationStrategy: expinfrav1.OnDemandAllocationStrategyPrioritized,
								},
								Overrides: nil,
							},
						},
					},
					Logger: *logger.NewLogger(logr.Discard()),
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:      ptr.To[int32](1),
					MaxSize:              2,
					MinSize:              0,
					CapacityRebalance:    true,
					MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{},
				},
			},
			wantDifference: true,
		},
		{
			name: "MixedInstancesPolicy.InstancesDistribution != asg.MixedInstancesPolicy.InstancesDistribution",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
							MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
								InstancesDistribution: &expinfrav1.InstancesDistribution{
									OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
									SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyCapacityOptimized,
									OnDemandBaseCapacity:                aws.Int64(0),
									OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
								},
								Overrides: []expinfrav1.Overrides{
									{
										InstanceType: "m6a.32xlarge",
									},
								},
							},
						},
					},
					Logger: *logger.NewLogger(logr.Discard()),
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:   ptr.To[int32](1),
					MaxSize:           2,
					MinSize:           0,
					CapacityRebalance: true,
					MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
						InstancesDistribution: &expinfrav1.InstancesDistribution{
							OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
							SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyLowestPrice,
							OnDemandBaseCapacity:                aws.Int64(0),
							OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
						},
						Overrides: []expinfrav1.Overrides{
							{
								InstanceType: "m6a.32xlarge",
							},
						},
					},
				},
			},
			wantDifference: true,
		},
		{
			name: "MixedInstancesPolicy.InstancesDistribution unset",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
							MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
								Overrides: []expinfrav1.Overrides{
									{
										InstanceType: "m6a.32xlarge",
									},
								},
							},
						},
					},
					Logger: *logger.NewLogger(logr.Discard()),
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:   ptr.To[int32](1),
					MaxSize:           2,
					MinSize:           0,
					CapacityRebalance: true,
					MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
						InstancesDistribution: &expinfrav1.InstancesDistribution{
							OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
							SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyLowestPrice,
							OnDemandBaseCapacity:                aws.Int64(0),
							OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
						},
						Overrides: []expinfrav1.Overrides{
							{
								InstanceType: "m6a.32xlarge",
							},
						},
					},
				},
			},
			wantDifference: false,
		},
		{
			name: "SuspendProcesses != asg.SuspendProcesses",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
							MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
								InstancesDistribution: &expinfrav1.InstancesDistribution{
									OnDemandAllocationStrategy: expinfrav1.OnDemandAllocationStrategyPrioritized,
								},
								Overrides: nil,
							},
							SuspendProcesses: &expinfrav1.SuspendProcessesTypes{
								Processes: &expinfrav1.Processes{
									Launch:    ptr.To[bool](true),
									Terminate: ptr.To[bool](true),
								},
							},
						},
					},
					Logger: *logger.NewLogger(logr.Discard()),
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:           ptr.To[int32](1),
					MaxSize:                   2,
					MinSize:                   0,
					CapacityRebalance:         true,
					MixedInstancesPolicy:      &expinfrav1.MixedInstancesPolicy{},
					CurrentlySuspendProcesses: []string{"Launch", "Terminate"},
				},
			},
			wantDifference: true,
		},
		{
			name: "all matches",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](1),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{
							MaxSize:           2,
							MinSize:           0,
							CapacityRebalance: true,
							MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
								InstancesDistribution: &expinfrav1.InstancesDistribution{
									OnDemandAllocationStrategy: expinfrav1.OnDemandAllocationStrategyPrioritized,
								},
								Overrides: nil,
							},
						},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity:   ptr.To[int32](1),
					MaxSize:           2,
					MinSize:           0,
					CapacityRebalance: true,
					MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
						InstancesDistribution: &expinfrav1.InstancesDistribution{
							OnDemandAllocationStrategy: expinfrav1.OnDemandAllocationStrategyPrioritized,
						},
						Overrides: nil,
					},
				},
			},
			wantDifference: false,
		},
		{
			name: "externally managed annotation ignores difference between desiredCapacity and replicas",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								clusterv1.ReplicasManagedByAnnotation: "", // empty value counts as true (= externally managed)
							},
						},
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](0),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
				},
			},
			wantDifference: false,
		},
		{
			name: "without externally managed annotation ignores difference between desiredCapacity and replicas",
			args: args{
				machinePoolScope: &scope.MachinePoolScope{
					MachinePool: &expclusterv1.MachinePool{
						Spec: expclusterv1.MachinePoolSpec{
							Replicas: ptr.To[int32](0),
						},
					},
					AWSMachinePool: &expinfrav1.AWSMachinePool{
						Spec: expinfrav1.AWSMachinePoolSpec{},
					},
				},
				existingASG: &expinfrav1.AutoScalingGroup{
					DesiredCapacity: ptr.To[int32](1),
				},
			},
			wantDifference: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			if tt.wantDifference {
				g.Expect(diffASG(tt.args.machinePoolScope, tt.args.existingASG)).ToNot(BeEmpty())
			} else {
				g.Expect(diffASG(tt.args.machinePoolScope, tt.args.existingASG)).To(BeEmpty())
			}
		})
	}
}
