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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

func TestManagedMachinePoolWaitingConditionPersistence(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()
	scheme := runtime.NewScheme()
	g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())
	g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
	g.Expect(ekscontrolplanev1.AddToScheme(scheme)).To(Succeed())
	g.Expect(expinfrav1.AddToScheme(scheme)).To(Succeed())
	cluster := &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster", Namespace: "default"}, Spec: clusterv1.ClusterSpec{ControlPlaneRef: clusterv1.ContractVersionedObjectReference{Kind: "AWSManagedControlPlane", Name: "control-plane"}}}
	pool := &clusterv1.MachinePool{ObjectMeta: metav1.ObjectMeta{Name: "pool", Namespace: "default", Labels: map[string]string{clusterv1.ClusterNameLabel: cluster.Name}}}
	infra := &expinfrav1.AWSManagedMachinePool{ObjectMeta: metav1.ObjectMeta{Name: "infra", Namespace: "default", OwnerReferences: []metav1.OwnerReference{{APIVersion: clusterv1.GroupVersion.String(), Kind: "MachinePool", Name: pool.Name}}}}
	cp := &ekscontrolplanev1.AWSManagedControlPlane{ObjectMeta: metav1.ObjectMeta{Name: "control-plane", Namespace: "default"}, Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{Region: "eu-west-1"}}
	v1beta1conditions.MarkTrue(infra, "Unrelated")
	c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(infra, cp).WithObjects(cluster, pool, infra, cp).Build()
	reconciler := &AWSManagedMachinePoolReconciler{Client: c}
	req := ctrl.Request{NamespacedName: client.ObjectKeyFromObject(infra)}
	// The first reconciliation backfills Paused, the second reaches the waiting path.
	_, err := reconciler.Reconcile(ctx, req)
	g.Expect(err).ToNot(HaveOccurred())
	_, err = reconciler.Reconcile(ctx, req)
	g.Expect(err).ToNot(HaveOccurred())
	stored := &expinfrav1.AWSManagedMachinePool{}
	g.Expect(c.Get(ctx, req.NamespacedName, stored)).To(Succeed())
	reported := v1beta1conditions.Get(stored, expinfrav1.EKSNodegroupReadyCondition)
	g.Expect(reported).ToNot(BeNil())
	g.Expect(reported.Status).To(Equal(corev1.ConditionFalse))
	g.Expect(reported.Reason).To(Equal(expinfrav1.WaitingForEKSControlPlaneReason))
	g.Expect(v1beta1conditions.IsTrue(stored, "Unrelated")).To(BeTrue())
}
