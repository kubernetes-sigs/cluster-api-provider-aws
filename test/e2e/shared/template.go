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
	"context"
	"io/ioutil"
	"path"

	"github.com/awslabs/goformation/v4/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v2"

	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1alpha1"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/credentials"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

const (
	MultiTenancyJumpPolicy = "CAPAMultiTenancyJumpPolicy"
)

var (
	accountRef = cloudformation.Sub("arn:${AWS::Partition}:iam::${AWS::AccountId}:root")
)

// newBootstrapTemplate generates a clusterawsadm configuration, and prints it
// and the resultant cloudformation template to the artifacts directory
func newBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	By("Creating a bootstrap AWSIAMConfiguration")
	t := cfn_bootstrap.NewTemplate()
	t.Spec.BootstrapUser.Enable = true
	t.Spec.BootstrapUser.ExtraStatements = []v1alpha4.StatementEntry{
		{
			Effect: "Allow",
			Action: []string{"sts:AssumeRole"},
			Resource: []string{
				cloudformation.GetAtt(MultiTenancySimpleRole.RoleName(), "Arn"),
				cloudformation.GetAtt(MultiTenancyJumpRole.RoleName(), "Arn"),
			},
		},
	}
	t.Spec.SecureSecretsBackends = []v1alpha4.SecretBackend{
		v1alpha4.SecretBackendSecretsManager,
		v1alpha4.SecretBackendSSMParameterStore,
	}
	t.Spec.EventBridge = &v1alpha1.EventBridgeConfig{
		Enable: true,
	}

	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	t.Spec.Region = region
	t.Spec.EKS.Disable = false
	t.Spec.EKS.AllowIAMRoleCreation = false
	t.Spec.EKS.DefaultControlPlaneRole.Disable = false
	t.Spec.EKS.ManagedMachinePool.Disable = false
	str, err := yaml.Marshal(t.Spec)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "awsiamconfiguration.yaml"), str, 0644)).To(Succeed())
	cloudformationTemplate := renderCustomCloudFormation(&t)
	cfnData, err := cloudformationTemplate.YAML()
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "cloudformation.yaml"), cfnData, 0644)).To(Succeed())
	return &t
}

func renderCustomCloudFormation(t *cfn_bootstrap.Template) *cloudformation.Template {
	cloudformationTemplate := t.RenderCloudFormation()
	appendMultiTenancyRoles(t, cloudformationTemplate)
	appendExtraPoliciesToBootstrapUser(t)
	return cloudformationTemplate
}

func appendExtraPoliciesToBootstrapUser(t *cfn_bootstrap.Template) {
	t.Spec.BootstrapUser.ExtraStatements = append(t.Spec.BootstrapUser.ExtraStatements, v1alpha4.StatementEntry{
		Effect: v1alpha4.EffectAllow,
		Resource: v1alpha4.Resources{
			"*",
		},
		Action: v1alpha4.Actions{
			"servicequotas:GetServiceQuota",
			"servicequotas:RequestServiceQuotaIncrease",
			"servicequotas:ListRequestedServiceQuotaChangeHistory",
			"elasticloadbalancing:DescribeAccountLimits",
			"ec2:DescribeAccountLimits",
			"cloudtrail:LookupEvents",
			"ssm:StartSession",
			"ssm:DescribeSessions",
			"ssm:GetConnectionStatus",
			"ssm:DescribeInstanceProperties",
			"ssm:GetDocument",
			"ssm:TerminateSession",
			"ssm:ResumeSession",
		},
	})
	t.Spec.BootstrapUser.ExtraStatements = append(t.Spec.BootstrapUser.ExtraStatements, v1alpha4.StatementEntry{
		Effect: v1alpha4.EffectAllow,
		Resource: v1alpha4.Resources{
			"arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
			"arn:*:iam::*:role/aws-service-role/servicequotas.amazonaws.com/AWSServiceRoleForServiceQuotas",
			"arn:*:iam::*:role/aws-service-role/support.amazonaws.com/AWSServiceRoleForSupport",
			"arn:*:iam::*:role/aws-service-role/trustedadvisor.amazonaws.com/AWSServiceRoleForTrustedAdvisor",
		},
		Action: v1alpha4.Actions{
			"iam:CreateServiceLinkedRole",
		},
	})
}

func appendMultiTenancyRoles(t *cfn_bootstrap.Template, cfnt *cloudformation.Template) {

	controllersPolicy := cfnt.Resources[string(cfn_bootstrap.ControllersPolicy)].(*cfn_iam.ManagedPolicy)
	controllersPolicy.Roles = append(
		controllersPolicy.Roles,
		cloudformation.Ref(MultiTenancySimpleRole.RoleName()),
		cloudformation.Ref(MultiTenancyNestedRole.RoleName()),
	)
	cfnt.Resources[MultiTenancyJumpPolicy] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: MultiTenancyJumpPolicy,
		PolicyDocument: &v1alpha4.PolicyDocument{
			Version: v1alpha4.CurrentVersion,
			Statement: []v1alpha4.StatementEntry{
				{
					Effect:   v1alpha4.EffectAllow,
					Resource: v1alpha4.Resources{cloudformation.GetAtt(MultiTenancyNestedRole.RoleName(), "Arn")},
					Action:   v1alpha4.Actions{"sts:AssumeRole"},
				},
			},
		},
		Roles: []string{cloudformation.Ref(MultiTenancyJumpRole.RoleName())},
	}

	cfnt.Resources[MultiTenancySimpleRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancySimpleRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(v1alpha4.PrincipalAWS, []string{accountRef}),
	}
	cfnt.Resources[MultiTenancyJumpRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancyJumpRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(v1alpha4.PrincipalAWS, []string{accountRef}),
	}
	cfnt.Resources[MultiTenancyNestedRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancyNestedRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(v1alpha4.PrincipalAWS, []string{accountRef}),
	}
}

// getBootstrapTemplate gets or generates a new bootstrap template
func getBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	if e2eCtx.Environment.BootstrapTemplate == nil {
		e2eCtx.Environment.BootstrapTemplate = newBootstrapTemplate(e2eCtx)
	}
	return e2eCtx.Environment.BootstrapTemplate
}

// ApplyTemplate will render a cluster template and apply it to the management cluster
func ApplyTemplate(ctx context.Context, configCluster clusterctl.ConfigClusterInput, clusterProxy framework.ClusterProxy) error {
	Expect(ctx).NotTo(BeNil(), "ctx is required for ApplyClusterTemplateAndWait")

	Byf("Getting the cluster template yaml")
	workloadClusterTemplate := clusterctl.ConfigCluster(ctx, clusterctl.ConfigClusterInput{
		KubeconfigPath:           configCluster.KubeconfigPath,
		ClusterctlConfigPath:     configCluster.ClusterctlConfigPath,
		Flavor:                   configCluster.Flavor,
		Namespace:                configCluster.Namespace,
		ClusterName:              configCluster.ClusterName,
		KubernetesVersion:        configCluster.KubernetesVersion,
		ControlPlaneMachineCount: configCluster.ControlPlaneMachineCount,
		WorkerMachineCount:       configCluster.WorkerMachineCount,
		InfrastructureProvider:   configCluster.InfrastructureProvider,
		LogFolder:                configCluster.LogFolder,
	})
	Expect(workloadClusterTemplate).ToNot(BeNil(), "Failed to get the cluster template")

	Byf("Applying the %s cluster template yaml to the cluster", configCluster.Flavor)
	return clusterProxy.Apply(ctx, workloadClusterTemplate)
}
