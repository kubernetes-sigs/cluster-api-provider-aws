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
	"context"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
)

func TestEKSConfigReconcilerReturnEarlyIfClusterInfraNotReady(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newEKSConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		Initialization: clusterv1.ClusterInitializationStatus{
			InfrastructureProvisioned: ptr.To(false),
		},
	}

	reconciler := EKSConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
	g.Expect(v1beta1conditions.IsFalse(config, eksbootstrapv1.DataSecretAvailableCondition)).To(BeTrue())
	condition := conditions.Get(config, eksbootstrapv1.EKSConfigDataSecretAvailableV1Beta2Condition)
	g.Expect(condition).ToNot(BeNil())
	g.Expect(condition.Status).To(Equal(metav1.ConditionFalse))
	g.Expect(condition.Reason).To(Equal(eksbootstrapv1.EKSConfigWaitingForClusterInfrastructureV1Beta2Reason))
}

func TestEKSConfigReconcilerReturnEarlyIfClusterControlPlaneNotInitialized(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newEKSConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		Initialization: clusterv1.ClusterInitializationStatus{
			InfrastructureProvisioned: ptr.To(true),
		},
	}

	reconciler := EKSConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
	g.Expect(v1beta1conditions.IsFalse(config, eksbootstrapv1.DataSecretAvailableCondition)).To(BeTrue())
	condition := conditions.Get(config, eksbootstrapv1.EKSConfigDataSecretAvailableV1Beta2Condition)
	g.Expect(condition).ToNot(BeNil())
	g.Expect(condition.Reason).To(Equal(eksbootstrapv1.EKSConfigWaitingForControlPlaneInitializationV1Beta2Reason))
}

func configOwner(kind string) *bsutil.ConfigOwner {
	unstructuredOwner := unstructured.Unstructured{
		Object: map[string]interface{}{"kind": kind},
	}
	configOwner := bsutil.ConfigOwner{Unstructured: &unstructuredOwner}
	return &configOwner
}

func TestBootstrapConditionsForExistingSecret(t *testing.T) {
	for _, kind := range []string{eksConfigKind, nodeadmConfigKind} {
		for _, stale := range []bool{false, true} {
			name := kind + "/missing condition"
			if stale {
				name = kind + "/recover failed condition"
			}
			t.Run(name, func(t *testing.T) {
				g := NewWithT(t)
				ctx := context.Background()
				scheme := runtime.NewScheme()
				g.Expect(eksbootstrapv1.AddToScheme(scheme)).To(Succeed())
				g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())
				g.Expect(corev1.AddToScheme(scheme)).To(Succeed())
				cluster := &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster", Namespace: "default"}}
				machine := &clusterv1.Machine{ObjectMeta: metav1.ObjectMeta{Name: "machine", Namespace: "default"}, Spec: clusterv1.MachineSpec{ClusterName: cluster.Name}}
				metadata := metav1.ObjectMeta{Name: "bootstrap", Namespace: "default", Generation: 4, OwnerReferences: []metav1.OwnerReference{{APIVersion: clusterv1.GroupVersion.String(), Kind: "Machine", Name: machine.Name}}}
				secretName := "existing-secret"
				var obj paused.ConditionSetter
				if kind == eksConfigKind {
					config := &eksbootstrapv1.EKSConfig{ObjectMeta: metadata}
					config.Status.DataSecretName = &secretName
					obj = config
				} else {
					config := &eksbootstrapv1.NodeadmConfig{ObjectMeta: metadata}
					config.Status.DataSecretName = &secretName
					obj = config
				}
				v1beta1conditions.MarkTrue(obj, eksbootstrapv1.DataSecretAvailableCondition)
				v1beta1conditions.MarkFalse(obj, clusterv1beta1.PausedV1Beta2Condition, clusterv1.NotPausedReason, clusterv1beta1.ConditionSeverityInfo, "")
				standard := obj.(conditions.Setter)
				conditions.Set(standard, metav1.Condition{Type: "OtherController", Status: metav1.ConditionTrue, Reason: "Observed"})
				if stale {
					conditions.Set(standard, metav1.Condition{Type: eksbootstrapv1.EKSConfigDataSecretAvailableV1Beta2Condition, Status: metav1.ConditionFalse, Reason: "GenerationFailed"})
				}
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: "default"}, Data: map[string][]byte{bootstrapDataKey: []byte("existing bootstrap data")}}
				c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj, cluster, machine, secret).Build()
				req := ctrl.Request{NamespacedName: client.ObjectKeyFromObject(obj)}
				var reconcile func(context.Context, ctrl.Request) (ctrl.Result, error)
				if kind == eksConfigKind {
					reconcile = (&EKSConfigReconciler{Client: c}).Reconcile
				} else {
					reconcile = (&NodeadmConfigReconciler{Client: c}).Reconcile
				}
				// First backfill Paused, then observe the existing secret and persist Ready.
				for range 2 {
					_, err := reconcile(ctx, req)
					g.Expect(err).ToNot(HaveOccurred())
				}
				g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), obj)).To(Succeed())
				ready := conditions.Get(standard, clusterv1.ReadyCondition)
				g.Expect(ready).ToNot(BeNil())
				g.Expect(ready.Status).To(Equal(metav1.ConditionTrue))
				g.Expect(ready.ObservedGeneration).To(Equal(int64(4)))
				g.Expect(conditions.IsTrue(standard, eksbootstrapv1.EKSConfigDataSecretAvailableV1Beta2Condition)).To(BeTrue())
				g.Expect(v1beta1conditions.IsTrue(obj, eksbootstrapv1.DataSecretAvailableCondition)).To(BeTrue())
				g.Expect(v1beta1conditions.IsTrue(obj, clusterv1beta1.ReadyCondition)).To(BeTrue())
				g.Expect(conditions.IsTrue(standard, "OtherController")).To(BeTrue())
				transition := ready.LastTransitionTime
				_, err := reconcile(ctx, req)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), obj)).To(Succeed())
				g.Expect(conditions.Get(standard, clusterv1.ReadyCondition).LastTransitionTime).To(Equal(transition))
				storedSecret := &corev1.Secret{}
				g.Expect(c.Get(ctx, client.ObjectKeyFromObject(secret), storedSecret)).To(Succeed())
				g.Expect(storedSecret.Data).To(Equal(secret.Data))
				g.Expect(storedSecret.ResourceVersion).To(Equal(secret.ResourceVersion))
			})
		}
	}
}
