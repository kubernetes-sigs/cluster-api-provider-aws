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

package conversiontest

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// fuzzCoverageSpy records Errorf calls for verification.
type fuzzCoverageSpy struct {
	testing.TB
	errors []string
}

func (s *fuzzCoverageSpy) Errorf(format string, args ...any) {
	s.errors = append(s.errors, fmt.Sprintf(format, args...))
}

func (s *fuzzCoverageSpy) Helper() {}

// TestSpoke implements conversion.Convertible for testing.
type TestSpoke struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

func (s *TestSpoke) GetObjectKind() schema.ObjectKind { return s }
func (s *TestSpoke) DeepCopyObject() runtime.Object {
	c := *s
	return &c
}
func (s *TestSpoke) ConvertFrom(hub conversion.Hub) error { return nil }
func (s *TestSpoke) ConvertTo(hub conversion.Hub) error   { return nil }

func TestRequireFuzzCoverage(t *testing.T) {
	t.Run("missing coverage fails", func(t *testing.T) {
		scheme := runtime.NewScheme()
		gv := schema.GroupVersion{Group: "test.example.com", Version: "v1"}

		// Register a Convertible type without coverage
		scheme.AddKnownTypeWithName(gv.WithKind("TestSpoke"), &TestSpoke{})

		spy := &fuzzCoverageSpy{TB: t}
		cases := []utilconversion.FuzzTestFuncInput{}
		knownGaps := map[string]string{}

		// Should error because TestSpoke is not covered
		RequireFuzzCoverage(spy, scheme, gv, cases, knownGaps)
		if len(spy.errors) == 0 {
			t.Error("expected errors for missing coverage, got none")
		}
	})

	t.Run("stale gap entries fail", func(t *testing.T) {
		scheme := runtime.NewScheme()
		gv := schema.GroupVersion{Group: "test.example.com", Version: "v1"}

		// Register a Convertible type with coverage
		scheme.AddKnownTypeWithName(gv.WithKind("TestSpoke"), &TestSpoke{})

		spy := &fuzzCoverageSpy{TB: t}
		cases := []utilconversion.FuzzTestFuncInput{
			{
				Spoke: &TestSpoke{},
			},
		}
		knownGaps := map[string]string{
			"TestSpoke": "test reason", // Stale: it's now covered
		}

		// Should error because TestSpoke is covered but still in knownGaps
		RequireFuzzCoverage(spy, scheme, gv, cases, knownGaps)
		if len(spy.errors) == 0 {
			t.Error("expected errors for stale gap entries, got none")
		}
	})

	t.Run("full coverage passes", func(t *testing.T) {
		scheme := runtime.NewScheme()
		gv := schema.GroupVersion{Group: "test.example.com", Version: "v1"}

		// Register a Convertible type
		scheme.AddKnownTypeWithName(gv.WithKind("TestSpoke"), &TestSpoke{})

		// Test that full coverage with correct gaps passes
		cases := []utilconversion.FuzzTestFuncInput{
			{
				Spoke: &TestSpoke{},
			},
		}
		knownGaps := map[string]string{}

		// This should not fail
		RequireFuzzCoverage(t, scheme, gv, cases, knownGaps)
	})
}
