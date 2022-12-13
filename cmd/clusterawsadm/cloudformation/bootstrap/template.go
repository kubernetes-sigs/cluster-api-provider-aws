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

	"github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"

	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/converters"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	eksiam "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/iam"
)

// Constants that define resources for a Template.
const (
	AWSIAMGroupBootstrapper                      = "AWSIAMGroupBootstrapper"
	AWSIAMInstanceProfileControllers             = "AWSIAMInstanceProfileControllers"
	AWSIAMInstanceProfileControlPlane            = "AWSIAMInstanceProfileControlPlane"
	AWSIAMInstanceProfileNodes                   = "AWSIAMInstanceProfileNodes"
	AWSIAMRoleControllers                        = "AWSIAMRoleControllers"
	AWSIAMRoleControlPlane                       = "AWSIAMRoleControlPlane"
	AWSIAMRoleNodes                              = "AWSIAMRoleNodes"
	AWSIAMRoleEKSControlPlane                    = "AWSIAMRoleEKSControlPlane"
	AWSIAMRoleEKSNodegroup                       = "AWSIAMRoleEKSNodegroup"
	AWSIAMRoleEKSFargate                         = "AWSIAMRoleEKSFargate"
	AWSIAMUserBootstrapper                       = "AWSIAMUserBootstrapper"
	ControllersPolicy                 PolicyName = "AWSIAMManagedPolicyControllers"
	ControllersPolicyEKS              PolicyName = "AWSIAMManagedPolicyControllersEKS"
	ControlPlanePolicy                PolicyName = "AWSIAMManagedPolicyCloudProviderControlPlane"
	NodePolicy                        PolicyName = "AWSIAMManagedPolicyCloudProviderNodes"
	CSIPolicy                         PolicyName = "AWSEBSCSIPolicyController"
	EKSConsolePolicy                  PolicyName = "AWSIAMManagedPolicyEKSConsole"
)

// Template is an AWS CloudFormation template to bootstrap
// IAM policies, users and roles for use by Cluster API Provider AWS.
type Template struct {
	Spec *bootstrapv1.AWSIAMConfigurationSpec
}

// NewTemplate will generate a new Template.
func NewTemplate() Template {
	conf := bootstrapv1.NewAWSIAMConfiguration()
	return Template{
		Spec: &conf.Spec,
	}
}

// NewManagedName creates an IAM acceptable name prefixed with this Cluster API
// implementation's prefix.
func (t Template) NewManagedName(name string) string {
	return fmt.Sprintf("%s%s%s", t.Spec.NamePrefix, name, *t.Spec.NameSuffix)
}

// RenderCloudFormation will render and return a cloudformation Template.
func (t Template) RenderCloudFormation() *cloudformation.Template {
	template := cloudformation.NewTemplate()

	if t.Spec.BootstrapUser.Enable {
		template.Resources[AWSIAMUserBootstrapper] = &cfn_iam.User{
			UserName:          t.Spec.BootstrapUser.UserName,
			Groups:            t.bootstrapUserGroups(),
			ManagedPolicyArns: t.Spec.ControlPlane.ExtraPolicyAttachments,
			Policies:          t.bootstrapUserPolicy(),
			Tags:              converters.MapToCloudFormationTags(t.Spec.BootstrapUser.Tags),
		}

		template.Resources[AWSIAMGroupBootstrapper] = &cfn_iam.Group{
			GroupName: t.Spec.BootstrapUser.GroupName,
		}
	}

	template.Resources[string(ControllersPolicy)] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: t.NewManagedName("controllers"),
		Description:       `For the Kubernetes Cluster API Provider AWS Controllers`,
		PolicyDocument:    t.ControllersPolicy(),
		Groups:            t.controllersPolicyGroups(),
		Roles:             t.controllersPolicyRoleAttachments(),
	}

	if !t.Spec.EKS.Disable {
		template.Resources[string(ControllersPolicyEKS)] = &cfn_iam.ManagedPolicy{
			ManagedPolicyName: t.NewManagedName("controllers-eks"),
			Description:       `For the Kubernetes Cluster API Provider AWS Controllers`,
			PolicyDocument:    t.ControllersPolicyEKS(),
			Groups:            t.controllersPolicyGroups(),
			Roles:             t.controllersPolicyRoleAttachments(),
		}
	}

	if !t.Spec.ControlPlane.DisableCloudProviderPolicy {
		template.Resources[string(ControlPlanePolicy)] = &cfn_iam.ManagedPolicy{
			ManagedPolicyName: t.NewManagedName("control-plane"),
			Description:       `For the Kubernetes Cloud Provider AWS Control Plane`,
			PolicyDocument:    t.cloudProviderControlPlaneAwsPolicy(),
			Roles:             t.cloudProviderControlPlaneAwsRoles(),
		}
	}

	if !t.Spec.Nodes.DisableCloudProviderPolicy {
		template.Resources[string(NodePolicy)] = &cfn_iam.ManagedPolicy{
			ManagedPolicyName: t.NewManagedName("nodes"),
			Description:       `For the Kubernetes Cloud Provider AWS nodes`,
			PolicyDocument:    t.nodePolicy(),
			Roles:             t.cloudProviderNodeAwsRoles(),
		}
	}

	if t.Spec.ControlPlane.EnableCSIPolicy {
		template.Resources[string(CSIPolicy)] = &cfn_iam.ManagedPolicy{
			ManagedPolicyName: t.NewManagedName("csi"),
			Description:       `For the AWS EBS CSI Driver for Kubernetes`,
			PolicyDocument:    t.csiControllerPolicy(),
			Roles:             t.csiControlPlaneAwsRoles(),
		}
	}

	template.Resources[AWSIAMRoleControlPlane] = &cfn_iam.Role{
		RoleName:                 t.NewManagedName("control-plane"),
		AssumeRolePolicyDocument: t.controlPlaneTrustPolicy(),
		ManagedPolicyArns:        t.Spec.ControlPlane.ExtraPolicyAttachments,
		Policies:                 t.controlPlanePolicies(),
		Tags:                     converters.MapToCloudFormationTags(t.Spec.ControlPlane.Tags),
	}

	template.Resources[AWSIAMRoleControllers] = &cfn_iam.Role{
		RoleName:                 t.NewManagedName("controllers"),
		AssumeRolePolicyDocument: t.controllersTrustPolicy(),
		Policies:                 t.controllersRolePolicy(),
		Tags:                     converters.MapToCloudFormationTags(t.Spec.ClusterAPIControllers.Tags),
	}

	template.Resources[AWSIAMRoleNodes] = &cfn_iam.Role{
		RoleName:                 t.NewManagedName("nodes"),
		AssumeRolePolicyDocument: t.nodeTrustPolicy(),
		ManagedPolicyArns:        t.nodeManagedPolicies(),
		Policies:                 t.nodePolicies(),
		Tags:                     converters.MapToCloudFormationTags(t.Spec.Nodes.Tags),
	}

	template.Resources[AWSIAMInstanceProfileControlPlane] = &cfn_iam.InstanceProfile{
		InstanceProfileName: t.NewManagedName("control-plane"),
		Roles: []string{
			cloudformation.Ref(AWSIAMRoleControlPlane),
		},
	}

	template.Resources[AWSIAMInstanceProfileControllers] = &cfn_iam.InstanceProfile{
		InstanceProfileName: t.NewManagedName("controllers"),
		Roles: []string{
			cloudformation.Ref(AWSIAMRoleControllers),
		},
	}

	template.Resources[AWSIAMInstanceProfileNodes] = &cfn_iam.InstanceProfile{
		InstanceProfileName: t.NewManagedName("nodes"),
		Roles: []string{
			cloudformation.Ref(AWSIAMRoleNodes),
		},
	}

	if !t.Spec.EKS.DefaultControlPlaneRole.Disable && !t.Spec.EKS.Disable {
		template.Resources[AWSIAMRoleEKSControlPlane] = &cfn_iam.Role{
			RoleName:                 ekscontrolplanev1.DefaultEKSControlPlaneRole,
			AssumeRolePolicyDocument: AssumeRolePolicy(iamv1.PrincipalService, []string{"eks.amazonaws.com"}),
			ManagedPolicyArns:        t.eksControlPlanePolicies(),
			Tags:                     converters.MapToCloudFormationTags(t.Spec.EKS.DefaultControlPlaneRole.Tags),
		}
	}

	if !t.Spec.EKS.ManagedMachinePool.Disable && !t.Spec.EKS.Disable {
		template.Resources[AWSIAMRoleEKSNodegroup] = &cfn_iam.Role{
			RoleName:                 expinfrav1.DefaultEKSNodegroupRole,
			AssumeRolePolicyDocument: AssumeRolePolicy(iamv1.PrincipalService, []string{"ec2.amazonaws.com", "eks.amazonaws.com"}),
			ManagedPolicyArns:        t.eksMachinePoolPolicies(),
			Tags:                     converters.MapToCloudFormationTags(t.Spec.EKS.ManagedMachinePool.Tags),
		}
	}

	if !t.Spec.EKS.Fargate.Disable && !t.Spec.EKS.Disable {
		template.Resources[AWSIAMRoleEKSFargate] = &cfn_iam.Role{
			RoleName:                 expinfrav1.DefaultEKSFargateRole,
			AssumeRolePolicyDocument: AssumeRolePolicy(iamv1.PrincipalService, []string{eksiam.EKSFargateService}),
			ManagedPolicyArns:        t.fargateProfilePolicies(t.Spec.EKS.Fargate),
			Tags:                     converters.MapToCloudFormationTags(t.Spec.EKS.Fargate.Tags),
		}
	}

	if t.Spec.EKS.EnableUserEKSConsolePolicy && !t.Spec.EKS.Disable {
		template.Resources[string(EKSConsolePolicy)] = &cfn_iam.ManagedPolicy{
			ManagedPolicyName: t.NewManagedName("eks-console"),
			Description:       `For users/groups to view EKS nodes and workloads`,
			PolicyDocument:    t.eksConsolePolicies(),
		}
	}

	return template
}

func ec2AssumeRolePolicy() *iamv1.PolicyDocument {
	return AssumeRolePolicy(iamv1.PrincipalService, []string{"ec2.amazonaws.com"})
}

// AWSArnAssumeRolePolicy will assume Policies using PolicyArns.
func AWSArnAssumeRolePolicy(identityID string) *iamv1.PolicyDocument {
	return AssumeRolePolicy(iamv1.PrincipalAWS, []string{identityID})
}

// AWSServiceAssumeRolePolicy will assume an AWS Service policy.
func AWSServiceAssumeRolePolicy(identityID string) *iamv1.PolicyDocument {
	return AssumeRolePolicy(iamv1.PrincipalService, []string{identityID})
}

// AssumeRolePolicy will create a role session and pass session policies programmatically.
func AssumeRolePolicy(identityType iamv1.PrincipalType, principalIDs []string) *iamv1.PolicyDocument {
	return &iamv1.PolicyDocument{
		Version: iamv1.CurrentVersion,
		Statement: []iamv1.StatementEntry{
			{
				Effect:    iamv1.EffectAllow,
				Principal: iamv1.Principals{identityType: principalIDs},
				Action:    iamv1.Actions{"sts:AssumeRole"},
			},
		},
	}
}
