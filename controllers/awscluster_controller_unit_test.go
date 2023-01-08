/*
Copyright 2022 The Kubernetes Authors.

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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/mock_services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
)

func TestAWSClusterReconcilerReconcile(t *testing.T) {
	testCases := []struct {
		name         string
		awsCluster   *infrav1.AWSCluster
		ownerCluster *clusterv1.Cluster
		expectError  bool
	}{
		{
			name: "Should fail Reconcile if owner cluster not found",
			awsCluster: &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{GenerateName: "aws-test-", OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: clusterv1.GroupVersion.String(),
					Kind:       "Cluster",
					Name:       "capi-fail-test",
					UID:        "1",
				}}}},
			expectError: true,
		},
		{
			name:        "Should not reconcile if owner reference is not set",
			awsCluster:  &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{GenerateName: "aws-test-"}},
			expectError: false,
		},
		{
			name:         "Should not Reconcile if cluster is paused",
			awsCluster:   &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{GenerateName: "aws-test-", Annotations: map[string]string{clusterv1.PausedAnnotation: ""}}},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{GenerateName: "capi-test-"}},
			expectError:  false,
		},
		{
			name:        "Should Reconcile successfully if no AWSCluster found",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			reconciler := &AWSClusterReconciler{
				Client: testEnv.Client,
			}

			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())

			if tc.ownerCluster != nil {
				tc.ownerCluster.Namespace = ns.Name
				g.Expect(testEnv.Create(ctx, tc.ownerCluster)).To(Succeed())
				defer func(do ...client.Object) {
					g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
				}(tc.ownerCluster)
				tc.awsCluster.OwnerReferences = []metav1.OwnerReference{
					{
						APIVersion: clusterv1.GroupVersion.String(),
						Kind:       "Cluster",
						Name:       tc.ownerCluster.Name,
						UID:        "1",
					},
				}
			}
			createCluster(g, tc.awsCluster, ns.Name)
			defer cleanupCluster(g, tc.awsCluster, ns)

			if tc.awsCluster != nil {
				_, err := reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: tc.awsCluster.Namespace,
						Name:      tc.awsCluster.Name,
					},
				})
				if tc.expectError {
					g.Expect(err).ToNot(BeNil())
				} else {
					g.Expect(err).To(BeNil())
				}
			} else {
				_, err := reconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: client.ObjectKey{
						Namespace: ns.Name,
						Name:      "test",
					},
				})
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestAWSClusterReconcileOperations(t *testing.T) {
	var (
		reconciler AWSClusterReconciler
		mockCtrl   *gomock.Controller
		ec2Svc     *mock_services.MockEC2Interface
		elbSvc     *mock_services.MockELBInterface
		networkSvc *mock_services.MockNetworkInterface
		sgSvc      *mock_services.MockSecurityGroupInterface
		recorder   *record.FakeRecorder
		ctx        context.Context
	)

	setup := func(t *testing.T, awsCluster *infrav1.AWSCluster) client.WithWatch {
		t.Helper()
		ctx = context.TODO()
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "capa-system",
			},
			Data: map[string][]byte{
				"AccessKeyID":     []byte("access-key-id"),
				"SecretAccessKey": []byte("secret-access-key"),
				"SessionToken":    []byte("session-token"),
			},
		}
		csClient := fake.NewClientBuilder().WithObjects(awsCluster, secret).Build()

		mockCtrl = gomock.NewController(t)
		ec2Svc = mock_services.NewMockEC2Interface(mockCtrl)
		elbSvc = mock_services.NewMockELBInterface(mockCtrl)
		networkSvc = mock_services.NewMockNetworkInterface(mockCtrl)
		sgSvc = mock_services.NewMockSecurityGroupInterface(mockCtrl)

		recorder = record.NewFakeRecorder(2)

		reconciler = AWSClusterReconciler{
			Client: csClient,
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2Interface {
				return ec2Svc
			},
			elbServiceFactory: func(elbScope scope.ELBScope) services.ELBInterface {
				return elbSvc
			},
			networkServiceFactory: func(clusterScope scope.ClusterScope) services.NetworkInterface {
				return networkSvc
			},
			securityGroupFactory: func(clusterScope scope.ClusterScope) services.SecurityGroupInterface {
				return sgSvc
			},
			Recorder: recorder,
		}
		return csClient
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	t.Run("Reconciling an AWSCluster", func(t *testing.T) {
		t.Run("Reconcile success", func(t *testing.T) {
			t.Run("Should successfully create AWSCluster with Cluster Finalizer and LoadBalancerReady status true on AWSCluster", func(t *testing.T) {
				g := NewWithT(t)
				runningCluster := func() {
					ec2Svc.EXPECT().ReconcileBastion().Return(nil)
					elbSvc.EXPECT().ReconcileLoadbalancers().Return(nil)
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
				}

				awsCluster := getAWSCluster("test", "test")
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				awsCluster.Status.Network.APIServerELB.DNSName = DNSName
				awsCluster.Status.Network.APIServerELB.AvailabilityZones = []string{"us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e"}
				cs.SetSubnets(infrav1.Subnets{
					{
						ID:               "private-subnet-1",
						AvailabilityZone: "us-east-1b",
						IsPublic:         false,
					},
					{
						ID:               "private-subnet-2",
						AvailabilityZone: "us-east-1a",
						IsPublic:         false,
					},
					{
						ID:               "private-subnet-3",
						AvailabilityZone: "us-east-1c",
						IsPublic:         false,
					},
					{
						ID:               "private-subnet-4",
						AvailabilityZone: "us-east-1d",
						IsPublic:         false,
					},
					{
						ID:               "private-subnet-5",
						AvailabilityZone: "us-east-1e",
						IsPublic:         false,
					},
				})
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).To(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.LoadBalancerReadyCondition, corev1.ConditionTrue, "", ""}})
				g.Expect(awsCluster.GetFinalizers()).To(ContainElement(infrav1.ClusterFinalizer))
			})
		})
		t.Run("Reconcile failure", func(t *testing.T) {
			expectedErr := errors.New("failed to get resource")
			t.Run("Should fail AWSCluster create with reconcile network failure", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(expectedErr)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).Should(Equal(expectedErr))
			})
			t.Run("Should fail AWSCluster create with ClusterSecurityGroupsReadyCondition status false", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(expectedErr)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).ToNot(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.ClusterSecurityGroupsReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.ClusterSecurityGroupReconciliationFailedReason}})
			})
			t.Run("Should fail AWSCluster create with BastionHostReadyCondition status false", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
					ec2Svc.EXPECT().ReconcileBastion().Return(expectedErr)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).ToNot(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.BastionHostReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.BastionHostFailedReason}})
			})
			t.Run("Should fail AWSCluster create with failure in LoadBalancer reconciliation", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
					ec2Svc.EXPECT().ReconcileBastion().Return(nil)
					elbSvc.EXPECT().ReconcileLoadbalancers().Return(expectedErr)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).ToNot(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.LoadBalancerReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityWarning, infrav1.LoadBalancerFailedReason}})
			})
			t.Run("Should fail AWSCluster create with LoadBalancer reconcile failure with WaitForDNSName condition as false", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
					ec2Svc.EXPECT().ReconcileBastion().Return(nil)
					elbSvc.EXPECT().ReconcileLoadbalancers().Return(nil)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).To(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.LoadBalancerReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitForDNSNameReason}})
			})
			t.Run("Should fail AWSCluster create with LoadBalancer reconcile failure with WaitForDNSNameResolve condition as false", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				runningCluster := func() {
					networkSvc.EXPECT().ReconcileNetwork().Return(nil)
					sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
					ec2Svc.EXPECT().ReconcileBastion().Return(nil)
					elbSvc.EXPECT().ReconcileLoadbalancers().Return(nil)
				}
				csClient := setup(t, &awsCluster)
				defer teardown()
				runningCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				awsCluster.Status.Network.APIServerELB.DNSName = "test-apiserver.us-east-1.aws"
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(cs)
				g.Expect(err).To(BeNil())
				expectAWSClusterConditions(g, cs.AWSCluster, []conditionAssertion{{infrav1.LoadBalancerReadyCondition, corev1.ConditionFalse, clusterv1.ConditionSeverityInfo, infrav1.WaitForDNSNameResolveReason}})
			})
		})
	})
	t.Run("Reconcile delete AWSCluster", func(t *testing.T) {
		t.Run("Reconcile success", func(t *testing.T) {
			deleteCluster := func() {
				ec2Svc.EXPECT().DeleteBastion().Return(nil)
				elbSvc.EXPECT().DeleteLoadbalancers().Return(nil)
				networkSvc.EXPECT().DeleteNetwork().Return(nil)
				sgSvc.EXPECT().DeleteSecurityGroups().Return(nil)
			}
			t.Run("Should successfully delete AWSCluster with Cluster Finalizer removed", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test", "test")
				csClient := setup(t, &awsCluster)
				defer teardown()
				deleteCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).To(BeNil())
				g.Expect(awsCluster.GetFinalizers()).ToNot(ContainElement(infrav1.ClusterFinalizer))
			})
		})
		t.Run("Reconcile failure", func(t *testing.T) {
			expectedErr := errors.New("failed to get resource")
			t.Run("Should fail AWSCluster delete with LoadBalancer deletion failed and Cluster Finalizer not removed", func(t *testing.T) {
				g := NewWithT(t)
				deleteCluster := func() {
					t.Helper()
					elbSvc.EXPECT().DeleteLoadbalancers().Return(expectedErr)
				}
				awsCluster := getAWSCluster("test", "test")
				awsCluster.Finalizers = []string{infrav1.ClusterFinalizer}
				csClient := setup(t, &awsCluster)
				defer teardown()
				deleteCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(awsCluster.GetFinalizers()).To(ContainElement(infrav1.ClusterFinalizer))
			})
			t.Run("Should fail AWSCluster delete with Bastion deletion failed and Cluster Finalizer not removed", func(t *testing.T) {
				g := NewWithT(t)
				deleteCluster := func() {
					ec2Svc.EXPECT().DeleteBastion().Return(expectedErr)
					elbSvc.EXPECT().DeleteLoadbalancers().Return(nil)
				}
				awsCluster := getAWSCluster("test", "test")
				awsCluster.Finalizers = []string{infrav1.ClusterFinalizer}
				csClient := setup(t, &awsCluster)
				defer teardown()
				deleteCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(awsCluster.GetFinalizers()).To(ContainElement(infrav1.ClusterFinalizer))
			})
			t.Run("Should fail AWSCluster delete with security group deletion failed and Cluster Finalizer not removed", func(t *testing.T) {
				g := NewWithT(t)
				deleteCluster := func() {
					ec2Svc.EXPECT().DeleteBastion().Return(nil)
					elbSvc.EXPECT().DeleteLoadbalancers().Return(nil)
					sgSvc.EXPECT().DeleteSecurityGroups().Return(expectedErr)
				}
				awsCluster := getAWSCluster("test", "test")
				awsCluster.Finalizers = []string{infrav1.ClusterFinalizer}
				csClient := setup(t, &awsCluster)
				defer teardown()
				deleteCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(awsCluster.GetFinalizers()).To(ContainElement(infrav1.ClusterFinalizer))
			})
			t.Run("Should fail AWSCluster delete with network deletion failed and Cluster Finalizer not removed", func(t *testing.T) {
				g := NewWithT(t)
				deleteCluster := func() {
					ec2Svc.EXPECT().DeleteBastion().Return(nil)
					elbSvc.EXPECT().DeleteLoadbalancers().Return(nil)
					sgSvc.EXPECT().DeleteSecurityGroups().Return(nil)
					networkSvc.EXPECT().DeleteNetwork().Return(expectedErr)
				}
				awsCluster := getAWSCluster("test", "test")
				awsCluster.Finalizers = []string{infrav1.ClusterFinalizer}
				csClient := setup(t, &awsCluster)
				defer teardown()
				deleteCluster()
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).ToNot(BeNil())
				g.Expect(awsCluster.GetFinalizers()).To(ContainElement(infrav1.ClusterFinalizer))
			})
		})
	})
}

func TestAWSClusterReconcilerRequeueAWSClusterForUnpausedCluster(t *testing.T) {
	testCases := []struct {
		name         string
		awsCluster   *infrav1.AWSCluster
		ownerCluster *clusterv1.Cluster
		requeue      bool
	}{
		{
			name: "Should create reconcile request successfully",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{GenerateName: "aws-test-"}, TypeMeta: metav1.TypeMeta{Kind: "AWSCluster", APIVersion: infrav1.GroupVersion.String()},
			},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test"}},
			requeue:      true,
		},
		{
			name: "Should not create reconcile request if AWSCluster is externally managed",
			awsCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{GenerateName: "aws-test-", Annotations: map[string]string{clusterv1.ManagedByAnnotation: "capi-test"}},
				TypeMeta:   metav1.TypeMeta{Kind: "AWSCluster", APIVersion: infrav1.GroupVersion.String()},
			},
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test"}},
			requeue:      false,
		},
		{
			name:         "Should not create reconcile request for deleted clusters",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test", DeletionTimestamp: &metav1.Time{Time: time.Now()}}},
			requeue:      false,
		},
		{
			name:         "Should not create reconcile request if infrastructure ref for AWSCluster on owner cluster is not set",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test"}},
			requeue:      false,
		},
		{
			name: "Should not create reconcile request if infrastructure ref type on owner cluster is not AWSCluster",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test"}, Spec: clusterv1.ClusterSpec{InfrastructureRef: &corev1.ObjectReference{
				APIVersion: clusterv1.GroupVersion.String(),
				Kind:       "Cluster",
				Name:       "aws-test"}}},
			requeue: false,
		},
		{
			name: "Should not create reconcile request if AWSCluster not found",
			ownerCluster: &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "capi-test"}, Spec: clusterv1.ClusterSpec{InfrastructureRef: &corev1.ObjectReference{
				APIVersion: clusterv1.GroupVersion.String(),
				Kind:       "AWSCluster",
				Name:       "aws-test"}}},
			requeue: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			log := logger.FromContext(ctx)
			reconciler := &AWSClusterReconciler{
				Client: testEnv.Client,
			}

			ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("namespace-%s", util.RandomString(5)))
			g.Expect(err).To(BeNil())
			createCluster(g, tc.awsCluster, ns.Name)
			defer cleanupCluster(g, tc.awsCluster, ns)

			if tc.ownerCluster != nil {
				if tc.awsCluster != nil {
					tc.ownerCluster.Spec = clusterv1.ClusterSpec{InfrastructureRef: &corev1.ObjectReference{
						APIVersion: infrav1.GroupVersion.String(),
						Kind:       "AWSCluster",
						Name:       tc.awsCluster.Name,
						Namespace:  ns.Name,
					}}
				}
				tc.ownerCluster.Namespace = ns.Name
			}
			handlerFunc := reconciler.requeueAWSClusterForUnpausedCluster(ctx, log)
			result := handlerFunc(tc.ownerCluster)
			if tc.requeue {
				g.Expect(result).To(ContainElement(reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: ns.Name,
						Name:      tc.awsCluster.Name,
					},
				}))
			} else {
				g.Expect(result).To(BeNil())
			}
		})
	}
}

func createCluster(g *WithT, awsCluster *infrav1.AWSCluster, namespace string) {
	if awsCluster != nil {
		awsCluster.Namespace = namespace
		g.Expect(testEnv.Create(ctx, awsCluster)).To(Succeed())
		g.Eventually(func() bool {
			cluster := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      awsCluster.Name,
				Namespace: namespace,
			}
			err := testEnv.Get(ctx, key, cluster)
			return err == nil
		}, 10*time.Second).Should(Equal(true))
	}
}

func cleanupCluster(g *WithT, awsCluster *infrav1.AWSCluster, namespace *corev1.Namespace) {
	if awsCluster != nil {
		func(do ...client.Object) {
			g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
		}(awsCluster, namespace)
	}
}

func TestSecurityGroupRolesForCluster(t *testing.T) {
	tests := []struct {
		name           string
		bastionEnabled bool
		want           []infrav1.SecurityGroupRole
	}{
		{
			name:           "Should use bastion security group when bastion is enabled",
			bastionEnabled: true,
			want:           append(defaultAWSSecurityGroupRoles, infrav1.SecurityGroupBastion),
		},
		{
			name:           "Should not use bastion security group when bastion is not enabled",
			bastionEnabled: false,
			want:           defaultAWSSecurityGroupRoles,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			c := getAWSCluster("test", "test")
			c.Spec.Bastion.Enabled = tt.bastionEnabled
			s, err := getClusterScope(c)
			g.Expect(err).To(BeNil(), "failed to create cluster scope for test")

			got := securityGroupRolesForCluster(*s)
			g.Expect(got).To(Equal(tt.want))
		})
	}
}
