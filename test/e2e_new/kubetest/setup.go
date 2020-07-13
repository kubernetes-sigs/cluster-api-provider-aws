// +build e2e

/*
Copyright 2020 The Kubernetes Authors.

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

package kubetest

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/yaml"
)

var (
	configFilePath          string
	useCIArtifacts          bool
	ginkgoNodes             int
	ginkgoSlowSpecThreshold int
)

func init() {
	flag.BoolVar(&useCIArtifacts, "kubetest.use-ci-artifacts", false, "use the latest build from the main branch of the Kubernetes repository")
	flag.StringVar(&configFilePath, "kubetest.config-file", "", "path to the kubetest configuration file")
	flag.IntVar(&ginkgoNodes, "kubetest.ginkgo-nodes", 1, "number of ginkgo nodes to use")
	flag.IntVar(&ginkgoSlowSpecThreshold, "kubetest.ginkgo-slowSpecThreshold", 120, "time in s before spec is marked as slow")
}

const (
	ciVersionURL = "https://dl.k8s.io/ci/k8s-master.txt"
)

// Configuration creates a new kubetest configuration
type Configuration struct {
	// UseCIArtifacts will fetch the latest build from the main branch of kubernetes/kubernetes
	UseCIArtifacts bool `json:"useCIArtifacts,omitempty"`
	// ArtifactsDirectory is where conformance suite output will go
	ArtifactsDirectory string `json:"artifactsDirectory,omitempty"`
	// E2EConfig is the clusterctl E2E config, which will be modified for conformance
	E2EConfig *clusterctl.E2EConfig `json:"e2eConfig,omitempty"`
	// Path to e2e config file
	ConfigFilePath string `json:"configFilePath,omitempty"`
	// GinkgoNodes is the number of Ginkgo nodes to use
	GinkgoNodes int `json:"ginkgoNodes,omitempty"`
	// GinkgoSlowSpecThreshold is time in s before spec is marked as slow
	GinkgoSlowSpecThreshold int `json:"ginkgoSlowSpecThreshold,omitempty"`
}

// NewConfiguration creates a new kubetest configuration. This may
// also mutate the clusterctl e2econfig if UseCIArtifacts is enabled.
func NewConfiguration(e2eConfig *clusterctl.E2EConfig, artifactsDirectory string) *Configuration {
	flag.Parse()
	c := Configuration{
		E2EConfig:               e2eConfig,
		UseCIArtifacts:          useCIArtifacts,
		ArtifactsDirectory:      artifactsDirectory,
		ConfigFilePath:          configFilePath,
		GinkgoNodes:             ginkgoNodes,
		GinkgoSlowSpecThreshold: ginkgoSlowSpecThreshold,
	}

	if c.UseCIArtifacts {
		newVersion := fetchKubernetesCIVersion()
		c.E2EConfig.Variables["KUBERNETES_VERSION"] = newVersion
		c.E2EConfig.Variables["CI_VERSION"] = newVersion
	}
	c.DumpE2EConfig()
	return &c
}

// DumpE2EConfig will dump the running e2e config for debugging
func (c Configuration) DumpE2EConfig() {
	yaml, err := yaml.Marshal(c.E2EConfig)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(c.ArtifactsDirectory, "e2econfig.yaml"), yaml, os.ModePerm)).To(Succeed())
}

// fetchKubernetesCIVersion fetches the latest main branch Kubernetes version
func fetchKubernetesCIVersion() string {
	resp, err := http.Get(ciVersionURL)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	return strings.TrimSpace(string(b))
}
