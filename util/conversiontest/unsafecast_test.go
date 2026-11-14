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
	"path/filepath"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata"
)

// testSpy records Errorf calls for verification.
type testSpy struct {
	testing.TB
	errors []string
}

func (s *testSpy) Errorf(format string, args ...any) {
	s.errors = append(s.errors, fmt.Sprintf(format, args...))
}

func (s *testSpy) Helper() {}

func TestCheckUnsafeStructCasts(t *testing.T) {
	scheme := runtime.NewScheme()

	// Get absolute path to testdata
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	testDataDir := filepath.Join(root, "testdata")
	generatedIdenticalFile := filepath.Join(testDataDir, "fake_generated_identical.go")
	generatedDifferentFile := filepath.Join(testDataDir, "fake_generated_different.go")
	generatedBypassFile := filepath.Join(testDataDir, "fake_generated_bypass.go")

	t.Run("layout identical types pass", func(t *testing.T) {
		t.Run("identical layout — no errors", func(t *testing.T) {
			// This test verifies that types with identical layout pass the check
			input := UnsafeCastInput{
				Scheme:        scheme,
				PackagePath:   "sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata",
				GeneratedFile: generatedIdenticalFile,
				ExtraTypes: []any{
					testdata.LayoutIdentical{},
					testdata.Hub{},
				},
			}
			spy := &testSpy{TB: t}
			CheckUnsafeStructCasts(spy, input)
			if len(spy.errors) != 0 {
				t.Errorf("expected no errors for identical layout types, got %d errors: %v", len(spy.errors), spy.errors)
			}
		})

		t.Run("different layout — errors reported", func(t *testing.T) {
			// This test verifies that types with different layout are detected
			input := UnsafeCastInput{
				Scheme:        scheme,
				PackagePath:   "sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata",
				GeneratedFile: generatedDifferentFile,
				ExtraTypes: []any{
					testdata.LayoutDifferent{},
					testdata.Hub{},
				},
			}
			spy := &testSpy{TB: t}
			CheckUnsafeStructCasts(spy, input)
			if len(spy.errors) == 0 {
				t.Error("expected error for layout-different type")
			}
		})
	})

	t.Run("bypass detection", func(t *testing.T) {
		// This test verifies that unsafe casts that bypass hand-written field converters
		// are detected. The testdata has hand-written Convert_InnerType_To_HubInner,
		// and autoConvert_SpokeWithBypass_To_HubWithBypass does an unsafe cast that
		// would skip this conversion.
		spy := &testSpy{TB: t}
		input := UnsafeCastInput{
			Scheme:        scheme,
			PackagePath:   "sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata",
			GeneratedFile: generatedBypassFile,
			ManualConversionDirs: []string{
				testDataDir,
			},
			ExtraTypes: []any{
				testdata.SpokeWithBypass{},
				testdata.HubWithBypass{},
				testdata.InnerType{},
				testdata.HubInner{},
			},
		}
		// Should error because unsafe cast bypasses the manual Convert_InnerType_To_HubInner
		CheckUnsafeStructCasts(spy, input)
		if len(spy.errors) == 0 {
			t.Error("expected error for unsafe cast that bypasses manual conversion, got none")
		}
	})
}

func TestWalkTypes(t *testing.T) {
	scheme := runtime.NewScheme()
	reg := buildTypeRegistry(scheme, []any{
		testdata.LayoutIdentical{},
		testdata.Hub{},
	})

	// Verify that types are registered
	if reg["sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata.LayoutIdentical"] == nil {
		t.Error("expected LayoutIdentical to be registered")
	}
	if reg["sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest/testdata.Hub"] == nil {
		t.Error("expected Hub to be registered")
	}
}

func TestSameLayout(t *testing.T) {
	type A struct {
		Field int
	}
	type B struct {
		Field int
	}

	aType := reflect.TypeOf(A{})
	bType := reflect.TypeOf(B{})

	if !sameLayout(aType, bType, make(map[typePair]bool)) {
		t.Error("expected same layout for identical struct types")
	}
}

func TestSameLayoutDifferent(t *testing.T) {
	type A struct {
		Field int
	}
	type B struct {
		Field int
		Extra string
	}

	aType := reflect.TypeOf(A{})
	bType := reflect.TypeOf(B{})

	if sameLayout(aType, bType, make(map[typePair]bool)) {
		t.Error("expected different layout for structs with different fields")
	}
}
