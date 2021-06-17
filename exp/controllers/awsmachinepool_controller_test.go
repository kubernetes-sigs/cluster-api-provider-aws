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

	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/mock_services"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	// "sigs.k8s.io/cluster-api/controllers/noderefutil" //nolint:godot.
	capierrors "sigs.k8s.io/cluster-api/errors"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

func TestAWSMachinePoolReconciler(t *testing.T) {
	var (
		reconciler     AWSMachinePoolReconciler
		cs             *scope.ClusterScope
		ms             *scope.MachinePoolScope
		mockCtrl       *gomock.Controller
		ec2Svc         *mock_services.MockEC2MachineInterface
		asgSvc         *mock_services.MockASGInterface
		recorder       *record.FakeRecorder
		awsMachinePool *expinfrav1.AWSMachinePool
		secret         *corev1.Secret
	)
	setup := func(t *testing.T, g *WithT) {
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
				MinSize: int32(1),
				MaxSize: int32(1),
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

		g.Expect(testEnv.Create(ctx, awsMachinePool)).To(Succeed())
		g.Expect(testEnv.Create(ctx, secret)).To(Succeed())

		ms, err = scope.NewMachinePoolScope(
			scope.MachinePoolScopeParams{
				Client: testEnv.Client,
				Cluster: &clusterv1.Cluster{
					Status: clusterv1.ClusterStatus{
						InfrastructureReady: true,
					},
				},
				MachinePool: &expclusterv1.MachinePool{
					Spec: expclusterv1.MachinePoolSpec{
						Template: clusterv1.MachineTemplateSpec{
							Spec: clusterv1.MachineSpec{
								Bootstrap: clusterv1.Bootstrap{
									DataSecretName: pointer.StringPtr("bootstrap-data"),
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

		cs, err = setupCluster("test-cluster")
		g.Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(t)
		ec2Svc = mock_services.NewMockEC2MachineInterface(mockCtrl)
		asgSvc = mock_services.NewMockASGInterface(mockCtrl)

		// If the test hangs for 9 minutes, increase the value here to the number of events during a reconciliation loop
		recorder = record.NewFakeRecorder(2)

		reconciler = AWSMachinePoolReconciler{
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2MachineInterface {
				return ec2Svc
			},
			asgServiceFactory: func(cloud.ClusterScoper) services.ASGInterface {
				return asgSvc
			},
			Recorder: recorder,
		}
	}

	teardown := func(t *testing.T, g *WithT) {
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
				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", expectedErr).AnyTimes()
				asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, expectedErr).AnyTimes()
			}
			t.Run("should exit immediately on an error state", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				er := capierrors.CreateMachineError
				ms.AWSMachinePool.Status.FailureReason = &er
				ms.AWSMachinePool.Status.FailureMessage = pointer.StringPtr("Couldn't create machine pool")

				buf := new(bytes.Buffer)
				klog.SetOutput(buf)

				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
			})
			t.Run("should add our finalizer to the machinepool", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				getASG(t, g)

				_, _ = reconciler.reconcileNormal(context.Background(), ms, cs, cs)

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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
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

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)

				g.Expect(err).To(BeNil())
				g.Expect(buf.String()).To(ContainSubstring("Bootstrap data secret reference is not yet available"))
				expectConditions(g, ms.AWSMachinePool, []conditionAssertion{{expinfrav1.ASGReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitingForBootstrapDataReason}})
			})
		})
		t.Run("there's a provider ID", func(t *testing.T) {
			id := "<cloudProvider>://<optional>/<segments>/<providerid>"
			setProviderID := func(t *testing.T, g *WithT) {
				_, err := noderefutil.NewProviderID(id)
				g.Expect(err).To(BeNil())

				ms.AWSMachinePool.Spec.ProviderID = id
			}
			t.Run("should look up by provider ID when one exists", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setProviderID(t, g)

				expectedErr := errors.New("no connection available ")
				var launchtemplate *expinfrav1.AWSLaunchTemplate
				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(launchtemplate, "", expectedErr)
				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
			t.Run("should try to create a new machinepool if none exists", func(t *testing.T) {
				g := NewWithT(t)
				setup(t, g)
				defer teardown(t, g)
				setProviderID(t, g)

				expectedErr := errors.New("Invalid instance")
				asgSvc.EXPECT().ASGIfExists(gomock.Any()).Return(nil, nil).AnyTimes()
				ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil)
				ec2Svc.EXPECT().DiscoverLaunchTemplateAMI(gomock.Any()).Return(nil, nil)
				ec2Svc.EXPECT().CreateLaunchTemplate(gomock.Any(), gomock.Any(), gomock.Any()).Return("", expectedErr).AnyTimes()

				_, err := reconciler.reconcileNormal(context.Background(), ms, cs, cs)
				g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
			})
		})
	})

	t.Run("Deleting an AWSMachinePool", func(t *testing.T) {
		finalizer := func(t *testing.T, g *WithT) {
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

			_, err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(errors.Cause(err)).To(MatchError(expectedErr))
		})
		t.Run("should log and remove finalizer when no machinepool exists", func(t *testing.T) {
			g := NewWithT(t)
			setup(t, g)
			defer teardown(t, g)
			finalizer(t, g)

			asgSvc.EXPECT().GetASGByName(gomock.Any()).Return(nil, nil)
			ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)

			_, err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(buf.String()).To(ContainSubstring("Unable to locate ASG"))
			g.Expect(ms.AWSMachinePool.Finalizers).To(ConsistOf(metav1.FinalizerDeleteDependents))
			g.Eventually(recorder.Events).Should(Receive(ContainSubstring("NoASGFound")))
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
			ec2Svc.EXPECT().GetLaunchTemplate(gomock.Any()).Return(nil, "", nil).AnyTimes()

			buf := new(bytes.Buffer)
			klog.SetOutput(buf)
			_, err := reconciler.reconcileDelete(ms, cs, cs)
			g.Expect(err).To(BeNil())
			g.Expect(ms.AWSMachinePool.Status.Ready).To(Equal(false))
			g.Eventually(recorder.Events).Should(Receive(ContainSubstring("DeletionInProgress")))
		})
	})
}

//TODO: This was taken from awsmachine_controller_test, i think it should be moved to elsewhere in both locations like test/helpers

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
