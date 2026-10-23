/*
Copyright 2026 The Kubernetes Authors.

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

package conditions

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	currentconditions "sigs.k8s.io/cluster-api/util/conditions"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
)

func TestSetV1Beta2(t *testing.T) {
	g := NewWithT(t)
	obj := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Generation: 3}}
	v1beta1conditions.MarkFalse(obj, infrav1.VpcReadyCondition, infrav1.VpcCreationStartedReason, clusterv1beta1.ConditionSeverityInfo, "legacy")
	legacy := obj.Status.Conditions.DeepCopy()
	transition := metav1.NewTime(time.Unix(100, 0))
	condition := metav1.Condition{Type: infrav1.AWSClusterVpcReadyV1Beta2Condition, Status: metav1.ConditionFalse, Reason: infrav1.AWSClusterVpcCreationStartedV1Beta2Reason, Message: "creating", LastTransitionTime: transition}
	SetV1Beta2(obj, condition)
	g.Expect(obj.Status.Conditions).To(Equal(legacy))
	actual := v1beta2conditions.Get(obj, condition.Type)
	g.Expect(actual.ObservedGeneration).To(Equal(int64(3)))
	g.Expect(actual.LastTransitionTime).To(Equal(transition))
	condition.Reason = infrav1.AWSClusterVpcReconciliationFailedV1Beta2Reason
	condition.Message = "failed"
	condition.LastTransitionTime = metav1.Time{}
	SetV1Beta2(obj, condition)
	g.Expect(v1beta2conditions.Get(obj, condition.Type).LastTransitionTime).To(Equal(transition))
	condition.Status = metav1.ConditionTrue
	condition.Reason = clusterv1.ReadyReason
	condition.Message = ""
	SetV1Beta2(obj, condition)
	g.Expect(v1beta2conditions.Get(obj, condition.Type).LastTransitionTime.After(transition.Time)).To(BeTrue())
	g.Expect(obj.Status.Conditions).To(Equal(legacy))

	pool := &clusterv1.MachinePool{ObjectMeta: metav1.ObjectMeta{Generation: 2}}
	SetV1Beta2(pool, condition)
	g.Expect(currentconditions.Get(pool, condition.Type).ObservedGeneration).To(Equal(int64(2)))
	// Shared services can still accept an object without condition storage.
	SetV1Beta2(&corev1.ConfigMap{}, condition)
}

func TestSetV1Beta2Summary(t *testing.T) {
	for _, tt := range []struct {
		name    string
		status  metav1.ConditionStatus
		missing bool
		want    metav1.ConditionStatus
		reason  string
	}{
		{name: "ready", status: metav1.ConditionTrue, want: metav1.ConditionTrue, reason: clusterv1.ReadyReason},
		{name: "failed", status: metav1.ConditionFalse, want: metav1.ConditionFalse, reason: clusterv1.NotReadyReason},
		{name: "unknown", status: metav1.ConditionUnknown, want: metav1.ConditionUnknown, reason: clusterv1.ReadyUnknownReason},
		{name: "missing dependency", missing: true, want: metav1.ConditionUnknown, reason: clusterv1.ReadyUnknownReason},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			obj := &infrav1.AWSMachine{}
			if !tt.missing {
				v1beta2conditions.Set(obj, metav1.Condition{Type: infrav1.AWSMachineInstanceReadyV1Beta2Condition, Status: tt.status, Reason: "Observed", Message: "instance state"})
			}
			v1beta2conditions.Set(obj, metav1.Condition{Type: infrav1.AWSMachineSecurityGroupsReadyV1Beta2Condition, Status: metav1.ConditionTrue, Reason: clusterv1.ReadyReason})
			g.Expect(SetV1Beta2Summary(obj, infrav1.InstanceReadyCondition, infrav1.SecurityGroupsReadyCondition)).To(Succeed())
			ready := v1beta2conditions.Get(obj, clusterv1.ReadyCondition)
			g.Expect(ready.Status).To(Equal(tt.want))
			g.Expect(ready.Reason).To(Equal(tt.reason))
			if tt.want == metav1.ConditionFalse {
				g.Expect(ready.Message).To(ContainSubstring("instance state"))
			}
		})
	}
}
