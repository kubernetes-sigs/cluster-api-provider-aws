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

package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	bootstrapv1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/conversiontest"
)

func TestGeneratedUnsafeConversions(t *testing.T) {
	scheme := runtime.NewScheme()
	g := NewWithT(t)
	g.Expect(AddToScheme(scheme)).To(Succeed())
	g.Expect(bootstrapv1beta1.AddToScheme(scheme)).To(Succeed())

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
		PackagePath:   "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1alpha1",
		GeneratedFile: filepath.Join(root, "cmd", "clusterawsadm", "api", "bootstrap", "v1alpha1", "zz_generated.conversion.go"),
		ManualConversionDirs: []string{
			filepath.Join(root, "cmd", "clusterawsadm", "api", "bootstrap", "v1alpha1"),
		},
		ExtraTypes: []any{
			BootstrapUser{},
			bootstrapv1beta1.BootstrapUser{},
			EventBridgeConfig{},
			bootstrapv1beta1.EventBridgeConfig{},
		},
	})
}
