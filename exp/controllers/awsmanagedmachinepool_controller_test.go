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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
)

func TestManagedMachinePoolReadyConditions(t *testing.T) {
	for _, tt := range []struct {
		name           string
		configured     bool
		templateStatus metav1.ConditionStatus
		want           metav1.ConditionStatus
	}{
		{name: "no custom launch template", want: metav1.ConditionTrue},
		{name: "removed launch template with stale failure", templateStatus: metav1.ConditionFalse, want: metav1.ConditionTrue},
		{name: "configured launch template not yet reported", configured: true, want: metav1.ConditionUnknown},
		{name: "configured launch template unknown", configured: true, templateStatus: metav1.ConditionUnknown, want: metav1.ConditionUnknown},
		{name: "configured launch template failed", configured: true, templateStatus: metav1.ConditionFalse, want: metav1.ConditionFalse},
		{name: "configured launch template ready", configured: true, templateStatus: metav1.ConditionTrue, want: metav1.ConditionTrue},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			obj := &expinfrav1.AWSManagedMachinePool{}
			if tt.configured {
				obj.Spec.AWSLaunchTemplate = &expinfrav1.AWSLaunchTemplate{}
			}
			for _, typ := range []clusterv1beta1.ConditionType{expinfrav1.EKSNodegroupReadyCondition, expinfrav1.IAMNodegroupRolesReadyCondition} {
				v1beta1conditions.MarkTrue(obj, typ)
				conditions.Set(obj, metav1.Condition{Type: string(typ), Status: metav1.ConditionTrue, Reason: clusterv1.ReadyReason})
			}
			if tt.templateStatus != "" {
				switch tt.templateStatus {
				case metav1.ConditionTrue:
					v1beta1conditions.MarkTrue(obj, expinfrav1.LaunchTemplateReadyCondition)
				case metav1.ConditionFalse:
					v1beta1conditions.MarkFalse(obj, expinfrav1.LaunchTemplateReadyCondition, "Observed", clusterv1beta1.ConditionSeverityInfo, "")
				case metav1.ConditionUnknown:
					v1beta1conditions.MarkUnknown(obj, expinfrav1.LaunchTemplateReadyCondition, "Observed", "")
				}
				conditions.Set(obj, metav1.Condition{Type: expinfrav1.AWSManagedMachinePoolLaunchTemplateReadyV1Beta2Condition, Status: tt.templateStatus, Reason: "Observed"})
			}
			g.Expect(setManagedMachinePoolReadyCondition(obj)).To(Succeed())
			g.Expect(conditions.Get(obj, clusterv1.ReadyCondition).Status).To(Equal(tt.want))
			if tt.want == metav1.ConditionTrue {
				g.Expect(v1beta1conditions.IsTrue(obj, clusterv1beta1.ReadyCondition)).To(BeTrue())
			}
			if tt.want == metav1.ConditionFalse {
				g.Expect(v1beta1conditions.IsFalse(obj, clusterv1beta1.ReadyCondition)).To(BeTrue())
			}
		})
	}
}
