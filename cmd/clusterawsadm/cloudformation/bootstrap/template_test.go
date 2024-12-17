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
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/sergi/go-diff/diffmatchpatch"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

func TestRenderCloudformation(t *testing.T) {
	cases := []struct {
		fixture  string
		template func() Template
	}{
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
			fixture: "with_s3_bucket",
			template: func() Template {
				t := NewTemplate()
				t.Spec.S3Buckets.Enable = true
				return t
			},
		},
		{
			fixture: "customsuffix",
			template: func() Template {
				t := NewTemplate()
				t.Spec.NameSuffix = ptr.To[string](".custom-suffix.com")
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
			fixture: "with_eks_default_roles",
			template: func() Template {
				t := NewTemplate()
				t.Spec.Nodes.EC2ContainerRegistryReadOnly = true
				t.Spec.EKS.DefaultControlPlaneRole.Disable = false
				t.Spec.EKS.ManagedMachinePool.Disable = false
				t.Spec.EKS.Fargate.Disable = false
				return t
			},
		},
		{
			fixture: "with_eks_kms_prefix",
			template: func() Template {
				t := NewTemplate()
				t.Spec.Nodes.EC2ContainerRegistryReadOnly = true
				t.Spec.EKS.KMSAliasPrefix = "custom-prefix-*"
				return t
			},
		},
		{
			fixture: "with_extra_statements",
			template: func() Template {
				t := NewTemplate()
				t.Spec.BootstrapUser.Enable = true
				t.Spec.ControlPlane.ExtraStatements = iamv1.Statements{
					{
						Effect:   iamv1.EffectAllow,
						Resource: iamv1.Resources{iamv1.Any},
						Action:   iamv1.Actions{"test:action"},
					},
				}
				t.Spec.Nodes.ExtraStatements = iamv1.Statements{
					{
						Effect:   iamv1.EffectAllow,
						Resource: iamv1.Resources{iamv1.Any},
						Action:   iamv1.Actions{"test:node-action"},
					},
				}
				t.Spec.BootstrapUser.ExtraStatements = iamv1.Statements{
					{
						Effect:   iamv1.EffectAllow,
						Resource: iamv1.Resources{iamv1.Any},
						Action:   iamv1.Actions{"test:user-action"},
					},
				}
				t.Spec.ClusterAPIControllers.ExtraStatements = iamv1.Statements{
					{
						Effect:   iamv1.EffectAllow,
						Resource: iamv1.Resources{iamv1.Any},
						Action:   iamv1.Actions{"test:controller-action"},
					},
				}
				return t
			},
		},
		{
			fixture: "with_eks_disable",
			template: func() Template {
				t := NewTemplate()
				t.Spec.EKS.Disable = true
				return t
			},
		},
		{
			fixture: "with_eks_console",
			template: func() Template {
				t := NewTemplate()
				t.Spec.EKS.EnableUserEKSConsolePolicy = true
				return t
			},
		},
		{
			fixture: "with_allow_assume_role",
			template: func() Template {
				t := NewTemplate()
				t.Spec.AllowAssumeRole = true
				return t
			},
		},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			cfn := cloudformation.Template{}
			data, err := os.ReadFile(path.Join("fixtures", c.fixture+".yaml"))
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

			if !bytes.Equal(tData, data) {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(string(tData), string(data), false)
				out := dmp.DiffPrettyText(diffs)
				t.Fatalf("Differing output (%s):\n%s", c.fixture, out)
			}
		})
	}
}
