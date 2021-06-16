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

package bootstrap

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/sergi/go-diff/diffmatchpatch"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/yaml"
)

type testCase struct {
	fixture  string
	template func() Template
}

func Test_RenderCloudformation(t *testing.T) {
	cases := []testCase{
		{
			fixture:  "default",
			template: NewTemplate,
		},
		{
			fixture: "with_ssm_secret_backend",
			template: func() Template {
				t := NewTemplate()
				t.Spec.SecureSecretsBackends = []infrav1.SecretBackend{
					infrav1.SecretBackendSSMParameterStore,
				}
				return t
			},
		},
		{
			fixture: "with_all_secret_backends",
			template: func() Template {
				t := NewTemplate()
				t.Spec.SecureSecretsBackends = []infrav1.SecretBackend{
					infrav1.SecretBackendSecretsManager,
					infrav1.SecretBackendSSMParameterStore,
				}
				return t
			},
		},
		{
			fixture: "customsuffix",
			template: func() Template {
				t := NewTemplate()
				t.Spec.NameSuffix = pointer.StringPtr(".custom-suffix.com")
				return t
			},
		},
		{
			fixture: "with_bootstrap_user",
			template: func() Template {
				t := NewTemplate()
				t.Spec.BootstrapUser.Enable = true
				return t
			},
		},
		{
			fixture: "with_custom_bootstrap_user",
			template: func() Template {
				t := NewTemplate()
				t.Spec.BootstrapUser.Enable = true
				t.Spec.BootstrapUser.UserName = "custom-bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
				return t
			},
		},
		{
			fixture: "with_different_instance_profiles",
			template: func() Template {
				t := NewTemplate()
				t.Spec.ClusterAPIControllers.AllowedEC2InstanceProfiles = []string{"customrole"}
				return t
			},
		},
		{
			fixture: "with_eks_enable",
			template: func() Template {
				t := NewTemplate()
				t.Spec.EKS.Enable = true
				t.Spec.Nodes.EC2ContainerRegistryReadOnly = true
				return t
			},
		},
		{
			fixture: "with_eks_default_roles",
			template: func() Template {
				t := NewTemplate()
				t.Spec.EKS.Enable = true
				t.Spec.Nodes.EC2ContainerRegistryReadOnly = true
				t.Spec.EKS.DefaultControlPlaneRole.Disable = false
				t.Spec.EKS.ManagedMachinePool.Disable = false
				t.Spec.EKS.Fargate.Disable = false
				return t
			},
		},
		{
			fixture: "with_extra_statements",
			template: func() Template {
				t := NewTemplate()
				t.Spec.BootstrapUser.Enable = true
				t.Spec.ControlPlane.ExtraStatements = infrav1.Statements{
					{
						Effect:   infrav1.EffectAllow,
						Resource: infrav1.Resources{infrav1.Any},
						Action:   infrav1.Actions{"test:action"},
					},
				}
				t.Spec.Nodes.ExtraStatements = infrav1.Statements{
					{
						Effect:   infrav1.EffectAllow,
						Resource: infrav1.Resources{infrav1.Any},
						Action:   infrav1.Actions{"test:node-action"},
					},
				}
				t.Spec.BootstrapUser.ExtraStatements = infrav1.Statements{
					{
						Effect:   infrav1.EffectAllow,
						Resource: infrav1.Resources{infrav1.Any},
						Action:   infrav1.Actions{"test:user-action"},
					},
				}
				t.Spec.ClusterAPIControllers.ExtraStatements = infrav1.Statements{
					{
						Effect:   infrav1.EffectAllow,
						Resource: infrav1.Resources{infrav1.Any},
						Action:   infrav1.Actions{"test:controller-action"},
					},
				}
				return t
			},
		},
	}

	for _, c := range cases {
		cfn := cloudformation.Template{}
		data, err := ioutil.ReadFile(path.Join("fixtures", c.fixture+".yaml"))
		if err != nil {
			t.Fatal(err)
		}
		err = yaml.Unmarshal(data, cfn)
		if err != nil {
			t.Fatal(err)
		}

		tData, err := c.template().RenderCloudFormation().YAML()
		if err != nil {
			t.Fatal(err)
		}
		ioutil.WriteFile("/tmp/tmp1", tData, 0600) // nolint:gosec

		if string(tData) != string(data) {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(string(tData), string(data), false)
			out := dmp.DiffPrettyText(diffs)
			t.Fatal(fmt.Sprintf("Differing output (%s):\n%s", c.fixture, out))
		}
	}
}
