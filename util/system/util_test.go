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

package system

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetManagerNamespace(t *testing.T) {
	cases := []struct {
		Name      string
		runBefore func()
		Expected  string
	}{
		{
			Name: "env var set to custom namespace",
			runBefore: func() {
				os.Setenv(namespaceEnvVarName, "capa")
			},
			Expected: "capa",
		},
		{
			Name: "env var empty",
			runBefore: func() {
				os.Unsetenv(namespaceEnvVarName)
			},
			Expected: "capa-system",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.runBefore()
			g := NewWithT(t)
			g.Expect(GetManagerNamespace()).To(Equal(c.Expected))
		})
	}
}

func TestGetNamespaceFromFile(t *testing.T) {
	g := NewWithT(t)
	path, err := os.Getwd()
	g.Expect(err).NotTo(HaveOccurred())
	nsPath := path + "namespace"
	_, err = os.OpenFile(filepath.Clean(nsPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	g.Expect(err).NotTo(HaveOccurred())
	ns := []byte("different-ns")
	g.Expect(os.WriteFile(nsPath, ns, 0644)).NotTo(HaveOccurred()) //nolint:gosec
	g.Expect(GetNamespaceFromFile(nsPath)).To(Equal("different-ns"))
	g.Expect(os.Remove(nsPath)).NotTo(HaveOccurred())
}
