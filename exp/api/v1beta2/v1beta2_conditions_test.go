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

package v1beta2

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

type conditionsObject interface {
	runtime.Object
	GetConditions() clusterv1beta1.Conditions
	SetConditions(clusterv1beta1.Conditions)
	GetV1Beta2Conditions() []metav1.Condition
	SetV1Beta2Conditions([]metav1.Condition)
}

func TestAdditionalV1Beta2Conditions(t *testing.T) {
	for _, tt := range []struct {
		name string
		obj  conditionsObject
	}{
		{name: "AWSFargateProfile", obj: &AWSFargateProfile{}},
		{name: "ROSANetwork", obj: &ROSANetwork{}},
		{name: "ROSARoleConfig", obj: &ROSARoleConfig{}},
		{name: "ROSAOCMRoleConfig", obj: &ROSAOCMRoleConfig{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.obj.GetV1Beta2Conditions(); got != nil {
				t.Fatalf("expected nil conditions before setting them, got %v", got)
			}

			legacy := clusterv1beta1.Conditions{{Type: "Ready"}}
			tt.obj.SetConditions(legacy)
			conditions := []metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Ready"}}
			tt.obj.SetV1Beta2Conditions(conditions)
			if got := tt.obj.GetV1Beta2Conditions(); !reflect.DeepEqual(got, conditions) {
				t.Fatalf("v1beta2 conditions = %v, want %v", got, conditions)
			}
			if got := tt.obj.GetConditions(); !reflect.DeepEqual(got, legacy) {
				t.Fatalf("legacy conditions changed: %v", got)
			}

			cloned := tt.obj.DeepCopyObject().(conditionsObject)
			cloned.GetV1Beta2Conditions()[0].Reason = "Changed"
			if got := tt.obj.GetV1Beta2Conditions()[0].Reason; got != "Ready" {
				t.Fatalf("deep copy changed original condition reason to %q", got)
			}
		})
	}
}
