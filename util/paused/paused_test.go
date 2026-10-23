/*
Copyright 2025 The Kubernetes Authors.

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

// Package paused implements paused helper functions.
package paused

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/test/builder"
)

func TestEnsurePausedCondition(t *testing.T) {
	g := NewWithT(t)

	scheme := runtime.NewScheme()
	g.Expect(builder.AddTransitionV1Beta2ToScheme(scheme)).To(Succeed())
	g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())

	// Cluster Case 1: unpaused
	normalCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-cluster",
			Namespace: "default",
		},
	}

	// Cluster Case 2: paused
	pausedCluster := normalCluster.DeepCopy()
	pausedCluster.Spec.Paused = ptr.To(true)

	// Object case 1: unpaused
	obj := &builder.Phase2Obj{ObjectMeta: metav1.ObjectMeta{
		Name:      "some-object",
		Namespace: "default",
	}}

	// Object case 2: paused
	pausedObj := obj.DeepCopy()
	pausedObj.SetAnnotations(map[string]string{clusterv1.PausedAnnotation: ""})

	tests := []struct {
		name         string
		cluster      *clusterv1.Cluster
		object       ConditionSetter
		wantIsPaused bool
	}{
		{
			name:         "unpaused cluster and unpaused object",
			cluster:      normalCluster.DeepCopy(),
			object:       obj.DeepCopy(),
			wantIsPaused: false,
		},
		{
			name:         "paused cluster and unpaused object",
			cluster:      pausedCluster.DeepCopy(),
			object:       obj.DeepCopy(),
			wantIsPaused: true,
		},
		{
			name:         "unpaused cluster and paused object",
			cluster:      normalCluster.DeepCopy(),
			object:       pausedObj.DeepCopy(),
			wantIsPaused: true,
		},
		{
			name:         "paused cluster and paused object",
			cluster:      pausedCluster.DeepCopy(),
			object:       pausedObj.DeepCopy(),
			wantIsPaused: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			ctx := context.Background()

			c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(&clusterv1.Cluster{}, &builder.Phase2Obj{}).
				WithObjects(tt.object, tt.cluster).Build()

			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(tt.object), tt.object)).To(Succeed())

			// The first run should set the condition.
			gotIsPaused, gotConditionChanged, err := EnsurePausedCondition(ctx, c, tt.cluster, tt.object)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(gotConditionChanged).To(BeTrue(), "The first reconcile should set the Paused condition")
			g.Expect(gotIsPaused).To(Equal(tt.wantIsPaused))

			// The second reconcile should be a no-op.
			gotIsPaused, gotConditionChanged, err = EnsurePausedCondition(ctx, c, tt.cluster, tt.object)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(gotConditionChanged).To(BeFalse(), "The second reconcile should not change the Paused condition")
			g.Expect(gotIsPaused).To(Equal(tt.wantIsPaused))
		})
	}
}

func TestEnsurePausedConditionDualWrite(t *testing.T) {
	for _, initiallyPaused := range []bool{false, true} {
		name := "unpaused"
		if initiallyPaused {
			name = "paused"
		}
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			ctx := context.Background()
			scheme := runtime.NewScheme()
			g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
			cluster := &clusterv1.Cluster{Spec: clusterv1.ClusterSpec{Paused: ptr.To(initiallyPaused)}}
			obj := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "cluster", Namespace: "default", Generation: 4}}
			// A matching legacy condition must not prevent backfilling standard conditions.
			legacy := pausedCondition(scheme, cluster, obj, clusterv1.PausedCondition)
			v1beta1conditions.Set(obj, &legacy)
			unrelated := metav1.Condition{Type: "OtherController", Status: metav1.ConditionTrue, Reason: "Observed"}
			conditions.Set(obj, unrelated)
			c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(obj).WithObjects(obj).Build()

			paused, changed, err := EnsurePausedCondition(ctx, c, cluster, obj)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(paused).To(Equal(initiallyPaused))
			g.Expect(changed).To(BeTrue())
			stored := &infrav1.AWSCluster{}
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
			actual := conditions.Get(stored, clusterv1.PausedCondition)
			g.Expect(actual).ToNot(BeNil())
			g.Expect(actual.Status).To(Equal(metav1.ConditionStatus(legacy.Status)))
			g.Expect(actual.Reason).To(Equal(legacy.Reason))
			g.Expect(actual.ObservedGeneration).To(Equal(int64(4)))
			g.Expect(conditions.Get(stored, unrelated.Type)).ToNot(BeNil())
			transition := actual.LastTransitionTime

			_, changed, err = EnsurePausedCondition(ctx, c, cluster, obj)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(changed).To(BeFalse())
			g.Expect(conditions.Get(obj, clusterv1.PausedCondition).LastTransitionTime).To(Equal(transition))

			// A spec update refreshes observedGeneration without stopping unpaused reconciliation.
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), obj)).To(Succeed())
			obj.Generation++
			g.Expect(c.Update(ctx, obj)).To(Succeed())
			resourceVersion := obj.ResourceVersion
			paused, changed, err = EnsurePausedCondition(ctx, c, cluster, obj)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(paused).To(Equal(initiallyPaused))
			g.Expect(changed).To(BeFalse())
			g.Expect(conditions.Get(obj, clusterv1.PausedCondition).ObservedGeneration).To(Equal(int64(5)))
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
			g.Expect(stored.ResourceVersion).ToNot(Equal(resourceVersion))
			g.Expect(conditions.Get(stored, clusterv1.PausedCondition).ObservedGeneration).To(Equal(int64(5)))
			g.Expect(conditions.Get(stored, clusterv1.PausedCondition).LastTransitionTime).To(Equal(transition))
			g.Expect(conditions.Get(stored, unrelated.Type)).To(Equal(conditions.Get(obj, unrelated.Type)))
			g.Expect(v1beta1conditions.HasSameState(v1beta1conditions.Get(stored, clusterv1beta1.PausedV1Beta2Condition), &legacy)).To(BeTrue())

			// Once the new generation is saved, reconciliation must not patch again.
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), obj)).To(Succeed())
			resourceVersion = obj.ResourceVersion
			_, changed, err = EnsurePausedCondition(ctx, c, cluster, obj)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(changed).To(BeFalse())
			g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
			g.Expect(stored.ResourceVersion).To(Equal(resourceVersion))

			cluster.Spec.Paused = ptr.To(!initiallyPaused)
			paused, changed, err = EnsurePausedCondition(ctx, c, cluster, obj)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(changed).To(BeTrue())
			g.Expect(paused).To(Equal(!initiallyPaused))
			g.Expect(conditions.Get(obj, clusterv1.PausedCondition).Status).To(Equal(metav1.ConditionStatus(v1beta1conditions.Get(obj, clusterv1beta1.PausedV1Beta2Condition).Status)))
		})
	}
}
