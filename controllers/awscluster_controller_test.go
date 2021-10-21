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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/mock_services"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAWSClusterReconciler(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	reconciler := &AWSClusterReconciler{
		Client: testEnv.Client,
	}

	instance := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	instance.Default()

	// Create the AWSCluster object and expect the Reconcile and Deployment to be created
	g.Expect(testEnv.Create(ctx, instance)).To(Succeed())

	// Calling reconcile should not error and not requeue the request with insufficient set up
	result, err := reconciler.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		},
	})
	g.Expect(err).To(BeNil())
	g.Expect(result).To(BeZero())
}

func TestAWSReconcileWithIdentity(t *testing.T) {
	var (
		reconciler AWSClusterReconciler
		mockCtrl   *gomock.Controller
		ec2Svc     *mock_services.MockEC2MachineInterface
		elbSvc     *mock_services.MockELBClusterInterface
		networkSvc *mock_services.MockNetworkInterface
		sgSvc      *mock_services.MockSecurityGroupInterface
		recorder   *record.FakeRecorder
	)
	ctx := context.Background()

	setup := func(awsCluster *infrav1.AWSCluster, t *testing.T, g *WithT) client.WithWatch {
		identity := getIdentityRef()
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
		csClient := fake.NewClientBuilder().WithObjects(awsCluster, identity, secret).Build()

		mockCtrl = gomock.NewController(t)
		ec2Svc = mock_services.NewMockEC2MachineInterface(mockCtrl)
		elbSvc = mock_services.NewMockELBClusterInterface(mockCtrl)
		networkSvc = mock_services.NewMockNetworkInterface(mockCtrl)
		sgSvc = mock_services.NewMockSecurityGroupInterface(mockCtrl)

		recorder = record.NewFakeRecorder(2)

		reconciler = AWSClusterReconciler{
			Client: csClient,
			ec2ServiceFactory: func(scope.EC2Scope) services.EC2MachineInterface {
				return ec2Svc
			},
			elbServiceFactory: func(clusterScope scope.ClusterScope) services.ELBClusterInterface {
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

	teardown := func(t *testing.T, g *WithT) {
		mockCtrl.Finish()
	}

	t.Run("Reconciling an AWSCluster", func(t *testing.T) {
		t.Run("with identity ref", func(t *testing.T) {
			runningCluster := func(t *testing.T, g *WithT) {
				ec2Svc.EXPECT().ReconcileBastion().Return(nil)
				elbSvc.EXPECT().ReconcileLoadbalancers().Return(nil)
				networkSvc.EXPECT().ReconcileNetwork().Return(nil)
				sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil)
			}

			t.Run("should add owner Ref as AWSCluster and finalizer to static identity ref", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test")
				csClient := setup(&awsCluster, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(ctx, cs)
				g.Expect(err).To(BeNil())

				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient, string(cs.IdentityRef().Kind), cs.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				g.Expect(unstructuredIdentityObject.GetOwnerReferences()).To(ConsistOf([]metav1.OwnerReference{{
					APIVersion:         "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:               "AWSCluster",
					Name:               "test",
					UID:                "",
					BlockOwnerDeletion: aws.Bool(true),
				}}))
				g.Expect(unstructuredIdentityObject.GetFinalizers()).To(ConsistOf(infrav1.AWSStaticIdentityFinalizer))
			})
			t.Run("should not allow deletion of identityRef if cluster is running", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test")
				csClient := setup(&awsCluster, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())

				_, err = reconciler.reconcileNormal(ctx, cs)
				g.Expect(err).To(BeNil())

				identityRef := getIdentityRef()
				err = csClient.Delete(ctx, identityRef)
				g.Expect(err).To(BeNil())

				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient, string(cs.IdentityRef().Kind), cs.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				// Because of finalizers added to identityRef, DeletionTimestamp is appended to the identityRef object
				// showing that delete operation was not successful and the object still exists
				g.Expect(unstructuredIdentityObject.GetDeletionTimestamp()).ToNot(BeNil())
				g.Expect(unstructuredIdentityObject.GetFinalizers()).NotTo(BeNil())
				g.Expect(unstructuredIdentityObject.GetName()).To(Equal("test-identity"))
			})
		})
	})
	t.Run("Reconciling delete AWSCluster", func(t *testing.T) {
		t.Run("with identity ref", func(t *testing.T) {
			runningCluster := func(t *testing.T, g *WithT) {
				ec2Svc.EXPECT().DeleteBastion().Return(nil).AnyTimes()
				elbSvc.EXPECT().DeleteLoadbalancers().Return(nil).AnyTimes()
				networkSvc.EXPECT().DeleteNetwork().Return(nil).AnyTimes()
				sgSvc.EXPECT().DeleteSecurityGroups().Return(nil).AnyTimes()

				ec2Svc.EXPECT().ReconcileBastion().Return(nil).AnyTimes()
				elbSvc.EXPECT().ReconcileLoadbalancers().Return(nil).AnyTimes()
				networkSvc.EXPECT().ReconcileNetwork().Return(nil).AnyTimes()
				sgSvc.EXPECT().ReconcileSecurityGroups().Return(nil).AnyTimes()
			}

			t.Run("should remove owner Ref as AWSCluster and finalizer from static identity ref", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test")
				csClient := setup(&awsCluster, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())

				_, err = reconciler.reconcileNormal(ctx, cs)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).To(BeNil())

				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient, string(cs.IdentityRef().Kind), cs.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				g.Expect(unstructuredIdentityObject.GetOwnerReferences()).NotTo(ConsistOf([]metav1.OwnerReference{{
					APIVersion:         "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:               "AWSCluster",
					Name:               "test",
					UID:                "",
					BlockOwnerDeletion: aws.Bool(true),
				}}))
				g.Expect(unstructuredIdentityObject.GetFinalizers()).NotTo(ConsistOf(infrav1.AWSStaticIdentityFinalizer))
			})

			t.Run("should ignore finalizer NotFound error and remove AWSCluster finalizer", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test")
				csClient := setup(&awsCluster, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())

				_, err = reconciler.reconcileNormal(ctx, cs)
				g.Expect(err).To(BeNil())

				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient, string(cs.IdentityRef().Kind), cs.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				unstructuredIdentityObject.SetFinalizers(nil)

				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).To(BeNil())
				g.Expect(awsCluster.Finalizers).NotTo(ContainElement(infrav1.ClusterFinalizer))
			})

			t.Run("should ignore ownerRef NotFound error and remove AWSCluster finalizer", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster := getAWSCluster("test")
				csClient := setup(&awsCluster, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster,
					},
				)
				g.Expect(err).To(BeNil())

				_, err = reconciler.reconcileNormal(ctx, cs)
				g.Expect(err).To(BeNil())

				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient, string(cs.IdentityRef().Kind), cs.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				unstructuredIdentityObject.SetOwnerReferences(nil)

				_, err = reconciler.reconcileDelete(ctx, cs)
				g.Expect(err).To(BeNil())
				g.Expect(awsCluster.Finalizers).NotTo(ContainElement(infrav1.ClusterFinalizer))
			})
			t.Run("should remove finalizer only when all ownerRef is deleted from identity", func(t *testing.T) {
				g := NewWithT(t)
				awsCluster1 := getAWSCluster("test-1")
				csClient1 := setup(&awsCluster1, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs1, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient1,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster1,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(ctx, cs1)
				g.Expect(err).To(BeNil())

				awsCluster2 := getAWSCluster("test-2")
				csClient2 := setup(&awsCluster2, t, g)
				defer teardown(t, g)
				runningCluster(t, g)
				cs2, err := scope.NewClusterScope(
					scope.ClusterScopeParams{
						Client:     csClient2,
						Cluster:    &clusterv1.Cluster{},
						AWSCluster: &awsCluster2,
					},
				)
				g.Expect(err).To(BeNil())
				_, err = reconciler.reconcileNormal(ctx, cs2)
				g.Expect(err).To(BeNil())

				_, err = reconciler.reconcileDelete(ctx, cs1)
				g.Expect(err).To(BeNil())
				unstructuredIdentityObject, err := getUnstructuredFromObjectReference(ctx, csClient1, string(cs1.IdentityRef().Kind), cs1.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				g.Expect(unstructuredIdentityObject.GetFinalizers()).NotTo(BeNil())

				_, err = reconciler.reconcileDelete(ctx, cs2)
				g.Expect(err).To(BeNil())
				unstructuredIdentityObject, err = getUnstructuredFromObjectReference(ctx, csClient2, string(cs2.IdentityRef().Kind), cs2.IdentityRef().Name)
				g.Expect(err).To(BeNil())
				g.Expect(unstructuredIdentityObject.GetFinalizers()).To(BeNil())
			})
		})
	})
}

func getAWSCluster(name string) infrav1.AWSCluster {
	return infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: infrav1.AWSClusterSpec{
			Region: "us-east-1",
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: "test-identity",
				Kind: infrav1.ClusterStaticIdentityKind,
			},
		},
	}
}

func getIdentityRef() *infrav1.AWSClusterStaticIdentity {
	return &infrav1.AWSClusterStaticIdentity{
		TypeMeta: metav1.TypeMeta{APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1"},
		Spec: infrav1.AWSClusterStaticIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{NamespaceList: []string{""}},
			},
			SecretRef: "test-secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-identity",
		},
	}
}
