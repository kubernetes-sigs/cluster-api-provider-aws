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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
	runtime "k8s.io/apimachinery/pkg/runtime"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func fuzzCases(scheme *runtime.Scheme) []utilconversion.FuzzTestFuncInput {
	return []utilconversion.FuzzTestFuncInput{
		{
			Scheme: scheme,
			Hub:    &eksbootstrapv1.EKSConfig{},
			Spoke:  &EKSConfig{},
		},
		{
			Scheme: scheme,
			Hub:    &eksbootstrapv1.EKSConfigTemplate{},
			Spoke:  &EKSConfigTemplate{},
		},
	}
}

func TestFuzzyConversion(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(eksbootstrapv1.AddToScheme(scheme)).To(Succeed())

	for _, c := range fuzzCases(scheme) {
		t.Run("for "+reflect.TypeOf(c.Spoke).Elem().Name(), utilconversion.FuzzTestFunc(c))
	}
}

func TestFuzzyConversionCoverage(t *testing.T) {
	g := NewWithT(t)
	scheme := runtime.NewScheme()
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(eksbootstrapv1.AddToScheme(scheme)).To(Succeed())

	conversiontest.RequireFuzzCoverage(t, scheme, GroupVersion, fuzzCases(scheme), map[string]string{})
}

func TestGeneratedUnsafeConversions(t *testing.T) {
	scheme := runtime.NewScheme()
	g := NewWithT(t)
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(eksbootstrapv1.AddToScheme(scheme)).To(Succeed())

	root, err := filepath.Abs(".")
	g.Expect(err).NotTo(HaveOccurred())
	// Walk up until we find the repo root (contains go.mod)
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("could not find repo root")
		}
		root = parent
	}

	conversiontest.CheckUnsafeStructCasts(t, conversiontest.UnsafeCastInput{
		Scheme:        scheme,
		PackagePath:   "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta1",
		GeneratedFile: filepath.Join(root, "bootstrap", "eks", "api", "v1beta1", "zz_generated.conversion.go"),
		ManualConversionDirs: []string{
			filepath.Join(root, "bootstrap", "eks", "api", "v1beta1"),
		},
		ExtraTypes: []any{},
	})
}
