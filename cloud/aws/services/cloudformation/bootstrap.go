// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudformation

import (
	"github.com/awslabs/goformation/cloudformation"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/iam"
)

// BootstrapTemplate is an AWS CloudFormation template to bootstrap
// IAM policies, users and roles for use by Cluster API Provider AWS
func BootstrapTemplate() *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources["AWSIAMManagedPolicyClusterController"] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("cluster-controller"),
		Description:       `For the Kubernetes Cluster API Provider AWS Cluster Controller`,
		PolicyDocument:    clusterControllerPolicy(),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleClusterController"),
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources["AWSIAMManagedPolicyMachineController"] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("machine-controller"),
		Description:       `For the Kubernetes Cluster API Provider AWS Machine Controller`,
		PolicyDocument:    machineControllerPolicy(),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleMachineController"),
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources["AWSIAMManagedPolicyCloudProviderControlPlane"] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: "control-plane-cloud-provider-aws.k8s.io",
		Description:       `For the Kubernetes Cloud Provider AWS Control Plane`,
		PolicyDocument:    cloudProviderControlPlaneAwsPolicy(),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources["AWSIAMManagedPolicyCloudProviderNodes"] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: "nodes.cloud-provider-aws.k8s.io",
		Description:       `For the Kubernetes Cloud Provider AWS nodes`,
		PolicyDocument:    cloudProviderNodeAwsPolicy(),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
			cloudformation.Ref("AWSIAMRoleNodes"),
		},
	}

	template.Resources["AWSIAMUserBootstrapper"] = cloudformation.AWSIAMUser{
		UserName: iam.NewManagedName("bootstrapper"),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
	}

	template.Resources["AWSIAMGroupBootstrapper"] = cloudformation.AWSIAMGroup{
		GroupName: iam.NewManagedName("bootstrapper"),
	}

	template.Resources["AWSIAMRoleControlPlane"] = cloudformation.AWSIAMRole{
		RoleName:                 iam.NewManagedName("control-plane"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
	}

	template.Resources["AWSIAMRoleClusterController"] = cloudformation.AWSIAMRole{
		RoleName:                 iam.NewManagedName("cluster-controller"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
	}

	template.Resources["AWSIAMRoleMachineController"] = cloudformation.AWSIAMRole{
		RoleName:                 iam.NewManagedName("machine-controller"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
	}

	template.Resources["AWSIAMRoleNodes"] = cloudformation.AWSIAMRole{
		RoleName:                 iam.NewManagedName("nodes"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
	}

	template.Resources["AWSIAMInstanceProfileControlPlane"] = cloudformation.AWSIAMInstanceProfile{
		InstanceProfileName: iam.NewManagedName("control-plane"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources["AWSIAMInstanceProfileClusterController"] = cloudformation.AWSIAMInstanceProfile{
		InstanceProfileName: iam.NewManagedName("cluster-controller"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleClusterController"),
		},
	}

	template.Resources["AWSIAMInstanceProfileMachineController"] = cloudformation.AWSIAMInstanceProfile{
		InstanceProfileName: iam.NewManagedName("machine-controller"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleMachineController"),
		},
	}

	template.Resources["AWSIAMInstanceProfileNodes"] = cloudformation.AWSIAMInstanceProfile{
		InstanceProfileName: iam.NewManagedName("nodes"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleNodes"),
		},
	}

	return template
}

func ec2AssumeRolePolicy() *iam.PolicyDocument {
	return &iam.PolicyDocument{
		Version: iam.CurrentVersion,
		Statement: []iam.StatementEntry{
			{
				Effect:    "Allow",
				Principal: iam.Principals{"Service": iam.PrincipalID{"ec2.amazonaws.com"}},
				Action:    iam.Actions{"sts:AssumeRole"},
			},
		},
	}
}

func clusterControllerPolicy() *iam.PolicyDocument {
	return &iam.PolicyDocument{
		Version: iam.CurrentVersion,
		Statement: []iam.StatementEntry{
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"ec2:AllocateAddress",
					"ec2:AssociateRouteTable",
					"ec2:AttachInternetGateway",
					"ec2:AuthorizeSecurityGroupIngress",
					"ec2:CreateInternetGateway",
					"ec2:CreateNatGateway",
					"ec2:CreateRoute",
					"ec2:CreateRouteTable",
					"ec2:CreateSecurityGroup",
					"ec2:CreateSubnet",
					"ec2:CreateTags",
					"ec2:CreateVpc",
					"ec2:DeleteInternetGateway",
					"ec2:DeleteNatGateway",
					"ec2:DeleteRouteTable",
					"ec2:DeleteSecurityGroup",
					"ec2:DeleteSubnet",
					"ec2:DeleteVpc",
					"ec2:DescribeAddresses",
					"ec2:DescribeAvailabilityZones",
					"ec2:DescribeInternetGateways",
					"ec2:DescribeNatGateways",
					"ec2:DescribeRouteTables",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeSubnets",
					"ec2:DescribeVpcs",
					"ec2:DetachInternetGateway",
					"ec2:DisassociateRouteTable",
					"ec2:ModifySubnetAttribute",
					"ec2:ReleaseAddress",
					"ec2:RevokeSecurityGroupIngress",
					"elasticloadbalancing:CreateLoadBalancer",
					"elasticloadbalancing:ConfigureHealthCheck",
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:DescribeLoadBalancers",
				},
			},
		},
	}
}

func machineControllerPolicy() *iam.PolicyDocument {
	return &iam.PolicyDocument{
		Version: iam.CurrentVersion,
		Statement: []iam.StatementEntry{
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"ec2:CreateTags",
					"ec2:DescribeInstances",
					"ec2:RunInstances",
					"ec2:TerminateInstances",
					"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
				},
			},
		},
	}
}

// From https://github.com/kubernetes/cloud-provider-aws
func cloudProviderControlPlaneAwsPolicy() *iam.PolicyDocument {
	return &iam.PolicyDocument{
		Version: iam.CurrentVersion,
		Statement: []iam.StatementEntry{
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"autoscaling:DescribeAutoScalingGroups",
					"autoscaling:DescribeLaunchConfigurations",
					"autoscaling:DescribeTags",
					"ec2:DescribeInstances",
					"ec2:DescribeRegions",
					"ec2:DescribeRouteTables",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeSubnets",
					"ec2:DescribeVolumes",
					"ec2:CreateSecurityGroup",
					"ec2:CreateTags",
					"ec2:CreateVolume",
					"ec2:ModifyInstanceAttribute",
					"ec2:ModifyVolume",
					"ec2:AttachVolume",
					"ec2:AuthorizeSecurityGroupIngress",
					"ec2:CreateRoute",
					"ec2:DeleteRoute",
					"ec2:DeleteSecurityGroup",
					"ec2:DeleteVolume",
					"ec2:DetachVolume",
					"ec2:RevokeSecurityGroupIngress",
					"ec2:DescribeVpcs",
					"elasticloadbalancing:AddTags",
					"elasticloadbalancing:AttachLoadBalancerToSubnets",
					"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
					"elasticloadbalancing:CreateLoadBalancer",
					"elasticloadbalancing:CreateLoadBalancerPolicy",
					"elasticloadbalancing:CreateLoadBalancerListeners",
					"elasticloadbalancing:ConfigureHealthCheck",
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:DeleteLoadBalancerListeners",
					"elasticloadbalancing:DescribeLoadBalancers",
					"elasticloadbalancing:DescribeLoadBalancerAttributes",
					"elasticloadbalancing:DetachLoadBalancerFromSubnets",
					"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
					"elasticloadbalancing:ModifyLoadBalancerAttributes",
					"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
					"elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
					"elasticloadbalancing:AddTags",
					"elasticloadbalancing:CreateListener",
					"elasticloadbalancing:CreateTargetGroup",
					"elasticloadbalancing:DeleteListener",
					"elasticloadbalancing:DeleteTargetGroup",
					"elasticloadbalancing:DescribeListeners",
					"elasticloadbalancing:DescribeLoadBalancerPolicies",
					"elasticloadbalancing:DescribeTargetGroups",
					"elasticloadbalancing:DescribeTargetHealth",
					"elasticloadbalancing:ModifyListener",
					"elasticloadbalancing:ModifyTargetGroup",
					"elasticloadbalancing:RegisterTargets",
					"elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
					"iam:CreateServiceLinkedRole",
					"kms:DescribeKey",
				},
			},
		},
	}
}

// From https://github.com/kubernetes/cloud-provider-aws
func cloudProviderNodeAwsPolicy() *iam.PolicyDocument {
	return &iam.PolicyDocument{
		Version: iam.CurrentVersion,
		Statement: []iam.StatementEntry{
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"ec2:DescribeInstances",
					"ec2:DescribeRegions",
					"ecr:GetAuthorizationToken",
					"ecr:BatchCheckLayerAvailability",
					"ecr:GetDownloadUrlForLayer",
					"ecr:GetRepositoryPolicy",
					"ecr:DescribeRepositories",
					"ecr:ListImages",
					"ecr:BatchGetImage",
				},
			},
		},
	}
}

// ReconcileBootstrapStack creates or updates bootstrap CloudFormation
func (s *Service) ReconcileBootstrapStack(stackName string) error {

	template := BootstrapTemplate()
	yaml, err := template.YAML()
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	if err := s.createStack(stackName, string(yaml)); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == "AlreadyExistsException" {
			glog.Infof("AWS Cloudformation stack %q already exists, updating", stackName)
			updateErr := s.updateStack(stackName, string(yaml))
			if updateErr != nil {
				code, ok := awserrors.Code(errors.Cause(updateErr))
				message := awserrors.Message(errors.Cause(updateErr))
				if !ok || code != "ValidationError" || message != "No updates are to be performed." {
					return updateErr
				}
			}
			return nil
		}
		return err
	}
	return nil
}
