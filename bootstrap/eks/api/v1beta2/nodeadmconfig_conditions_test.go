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

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

func TestNodeadmConfigV1Beta2Conditions(t *testing.T) {
	obj := &NodeadmConfig{}
	if got := obj.GetV1Beta2Conditions(); got != nil {
		t.Fatalf("expected nil conditions before setting them, got %v", got)
	}

	legacy := clusterv1beta1.Conditions{{Type: "Ready"}}
	obj.SetConditions(legacy)
	conditions := []metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Ready"}}
	obj.SetV1Beta2Conditions(conditions)
	if got := obj.GetV1Beta2Conditions(); !reflect.DeepEqual(got, conditions) {
		t.Fatalf("v1beta2 conditions = %v, want %v", got, conditions)
	}
	if got := obj.GetConditions(); !reflect.DeepEqual(got, legacy) {
		t.Fatalf("legacy conditions changed: %v", got)
	}

	cloned := obj.DeepCopy()
	cloned.Status.V1Beta2.Conditions[0].Reason = "Changed"
	if got := obj.GetV1Beta2Conditions()[0].Reason; got != "Ready" {
		t.Fatalf("deep copy changed original condition reason to %q", got)
	}
}
