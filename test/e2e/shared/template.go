//go:build e2e
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
	"fmt"
	"os"
	"path"

	"github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
	cfn_bootstrap "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cloudformation/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/credentials"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

const (
	// MultiTenancyJumpPolicy is the policy name for jump host to be used in multi-tenancy test.
	MultiTenancyJumpPolicy = "CAPAMultiTenancyJumpPolicy"
)

var (
	accountRef = cloudformation.Sub("arn:${AWS::Partition}:iam::${AWS::AccountId}:root")
)

// newBootstrapTemplate generates a clusterawsadm configuration, and prints it
// and the resultant cloudformation template to the artifacts directory.
func newBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	By("Creating a bootstrap AWSIAMConfiguration")
	t := cfn_bootstrap.NewTemplate()
	t.Spec.BootstrapUser.Enable = true
	t.Spec.BootstrapUser.ExtraStatements = []iamv1.StatementEntry{
		{
			Effect: iamv1.EffectAllow,
			Resource: iamv1.Resources{
				iamv1.Any,
			},
			Action: iamv1.Actions{
				"ecr-public:*",
				"sts:*",
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
				"ec2:DescribeSubnets",
				"ec2:DescribeNetworkInterfaces",
				"ec2:CreateNetworkInterface",
				"ec2:DescribeAvailabilityZones",
				"ec2:DeleteNetworkInterface",
				"elasticfilesystem:DescribeMountTargets",
				"elasticfilesystem:CreateMountTarget",
				"elasticfilesystem:DeleteMountTarget",
				"elasticfilesystem:DescribeFileSystems",
				"elasticfilesystem:CreateFileSystem",
				"elasticfilesystem:DeleteFileSystem",
				"elasticfilesystem:DescribeAccessPoints",
				"elasticfilesystem:CreateAccessPoint",
				"elasticfilesystem:DeleteAccessPoint",
			},
		},
		{
			Effect: iamv1.EffectAllow,
			Resource: iamv1.Resources{
				"arn:*:iam::*:role/aws-service-role/servicequotas.amazonaws.com/AWSServiceRoleForServiceQuotas",
				"arn:*:iam::*:role/aws-service-role/support.amazonaws.com/AWSServiceRoleForSupport",
				"arn:*:iam::*:role/aws-service-role/trustedadvisor.amazonaws.com/AWSServiceRoleForTrustedAdvisor",
			},
			Action: iamv1.Actions{
				"iam:CreateServiceLinkedRole",
			},
		},
	}
	t.Spec.SecureSecretsBackends = []infrav1.SecretBackend{
		infrav1.SecretBackendSecretsManager,
		infrav1.SecretBackendSSMParameterStore,
	}
	t.Spec.EventBridge = &bootstrapv1.EventBridgeConfig{
		Enable: true,
	}

	region, err := credentials.ResolveRegion("")
	Expect(err).NotTo(HaveOccurred())
	t.Spec.Region = region
	t.Spec.EKS.Disable = false
	t.Spec.EKS.AllowIAMRoleCreation = false
	t.Spec.EKS.DefaultControlPlaneRole.Disable = false
	t.Spec.EKS.ManagedMachinePool.Disable = false
	t.Spec.S3Buckets.Enable = true
	t.Spec.Nodes.ExtraStatements = []iamv1.StatementEntry{
		{
			Effect: iamv1.EffectAllow,
			Resource: iamv1.Resources{
				iamv1.Any,
			},
			Action: iamv1.Actions{
				"elasticfilesystem:DescribeMountTargets",
				"elasticfilesystem:DeleteAccessPoint",
				"elasticfilesystem:DescribeAccessPoints",
				"elasticfilesystem:DescribeFileSystems",
				"elasticfilesystem:CreateAccessPoint",
				"elasticfilesystem:TagResource",
				"ec2:DescribeAvailabilityZones",
			},
		},
	}
	str, err := yaml.Marshal(t.Spec)
	Expect(err).NotTo(HaveOccurred())
	Expect(os.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "awsiamconfiguration.yaml"), str, 0644)).To(Succeed()) //nolint:gosec
	cloudformationTemplate := renderCustomCloudFormation(&t)
	cfnData, err := cloudformationTemplate.YAML()
	Expect(err).NotTo(HaveOccurred())
	Expect(os.WriteFile(path.Join(e2eCtx.Settings.ArtifactFolder, "cloudformation.yaml"), cfnData, 0644)).To(Succeed()) //nolint:gosec
	return &t
}

func renderCustomCloudFormation(t *cfn_bootstrap.Template) *cloudformation.Template {
	cloudformationTemplate := t.RenderCloudFormation()
	appendMultiTenancyRoles(t, cloudformationTemplate)
	return cloudformationTemplate
}

func appendMultiTenancyRoles(_ *cfn_bootstrap.Template, cfnt *cloudformation.Template) {
	controllersPolicy := cfnt.Resources[string(cfn_bootstrap.ControllersPolicy)].(*cfn_iam.ManagedPolicy)
	controllersPolicy.Roles = append(
		controllersPolicy.Roles,
		cloudformation.Ref(MultiTenancySimpleRole.RoleName()),
		cloudformation.Ref(MultiTenancyNestedRole.RoleName()),
	)
	cfnt.Resources[MultiTenancyJumpPolicy] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: MultiTenancyJumpPolicy,
		PolicyDocument: &iamv1.PolicyDocument{
			Version: iamv1.CurrentVersion,
			Statement: []iamv1.StatementEntry{
				{
					Effect:   iamv1.EffectAllow,
					Resource: iamv1.Resources{cloudformation.GetAtt(MultiTenancyNestedRole.RoleName(), "Arn")},
					Action:   iamv1.Actions{"sts:AssumeRole"},
				},
			},
		},
		Roles: []string{cloudformation.Ref(MultiTenancyJumpRole.RoleName())},
	}

	cfnt.Resources[MultiTenancySimpleRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancySimpleRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(iamv1.PrincipalAWS, []string{accountRef}),
	}
	cfnt.Resources[MultiTenancyJumpRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancyJumpRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(iamv1.PrincipalAWS, []string{accountRef}),
	}
	cfnt.Resources[MultiTenancyNestedRole.RoleName()] = &cfn_iam.Role{
		RoleName:                 MultiTenancyNestedRole.RoleName(),
		AssumeRolePolicyDocument: cfn_bootstrap.AssumeRolePolicy(iamv1.PrincipalAWS, []string{accountRef}),
	}
}

// getBootstrapTemplate gets or generates a new bootstrap template.
func getBootstrapTemplate(e2eCtx *E2EContext) *cfn_bootstrap.Template {
	if e2eCtx.Environment.BootstrapTemplate == nil {
		e2eCtx.Environment.BootstrapTemplate = newBootstrapTemplate(e2eCtx)
	}
	return e2eCtx.Environment.BootstrapTemplate
}

// ApplyTemplate will render a cluster template and apply it to the management cluster.
func ApplyTemplate(ctx context.Context, configCluster clusterctl.ConfigClusterInput, clusterProxy framework.ClusterProxy) error {
	workloadClusterTemplate := GetTemplate(ctx, configCluster)
	By(fmt.Sprintf("Applying the %s cluster template yaml to the cluster", configCluster.Flavor))
	return clusterProxy.CreateOrUpdate(ctx, workloadClusterTemplate)
}

// GetTemplate will render a cluster template.
func GetTemplate(ctx context.Context, configCluster clusterctl.ConfigClusterInput) []byte {
	Expect(ctx).NotTo(BeNil(), "ctx is required for ApplyClusterTemplateAndWait")

	By("Getting the cluster template yaml")
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

	return workloadClusterTemplate
}
