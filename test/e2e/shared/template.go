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

package shared

import (
	"io/ioutil"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v2"

	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
)

// newBootstrapTemplate generates a clusterawsadm configuration, and prints it
// and the resultant cloudformation template to the artifacts directory
func newBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	By("Creating a bootstrap AWSIAMConfiguration")
	t := cfn_bootstrap.NewTemplate()
	t.Spec.BootstrapUser.Enable = true
	t.Spec.SecureSecretsBackends = []v1alpha3.SecretBackend{
		v1alpha3.SecretBackendSecretsManager,
		v1alpha3.SecretBackendSSMParameterStore,
	}
	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	t.Spec.Region = region
	str, err := yaml.Marshal(t.Spec)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "awsiamconfiguration.yaml"), str, 0644)).To(Succeed())
	cfnData, err := t.RenderCloudFormation().YAML()
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "cloudformation.yaml"), cfnData, 0644)).To(Succeed())
	return &t
}

// getBootstrapTemplate gets or generates a new bootstrap template
func getBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	if e2eCtx.Environment.BootstrapTemplate == nil {
		e2eCtx.Environment.BootstrapTemplate = newBootstrapTemplate(e2eCtx)
	}
	return e2eCtx.Environment.BootstrapTemplate
}
