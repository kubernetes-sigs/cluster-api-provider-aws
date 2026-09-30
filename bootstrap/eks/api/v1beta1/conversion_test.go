/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(eksbootstrapv1.AddToScheme(scheme)).To(Succeed())

	t.Run("for EKSConfig", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &eksbootstrapv1.EKSConfig{},
		Spoke:  &EKSConfig{},
	}))

	t.Run("for EKSConfigTemplate", utilconversion.FuzzTestFunc(utilconversion.FuzzTestFuncInput{
		Scheme: scheme,
		Hub:    &eksbootstrapv1.EKSConfigTemplate{},
		Spoke:  &EKSConfigTemplate{},
	}))
}

func TestEKSConfigV1Beta2ConditionsRoundTrip(t *testing.T) {
	g := NewWithT(t)
	hub := &eksbootstrapv1.EKSConfig{}
	hub.SetV1Beta2Conditions([]metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Available"}})

	spoke := &EKSConfig{}
	g.Expect(spoke.ConvertFrom(hub)).To(Succeed())
	restored := &eksbootstrapv1.EKSConfig{}
	g.Expect(spoke.ConvertTo(restored)).To(Succeed())
	g.Expect(restored.GetV1Beta2Conditions()).To(Equal(hub.GetV1Beta2Conditions()))
}
