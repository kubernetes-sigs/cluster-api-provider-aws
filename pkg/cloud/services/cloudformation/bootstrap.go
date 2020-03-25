/*
Copyright 2018 The Kubernetes Authors.

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

package cloudformation

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/pkg/errors"
	"k8s.io/klog"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/iam"
)

const (
	ControllersPolicy  = "AWSIAMManagedPolicyControllers"
	ControlPlanePolicy = "AWSIAMManagedPolicyCloudProviderControlPlane"
	NodePolicy         = "AWSIAMManagedPolicyCloudProviderNodes"
)

// ManagedIAMPolicyNames slice of managed IAM policies
var ManagedIAMPolicyNames = [...]string{ControllersPolicy, ControlPlanePolicy, NodePolicy}

// BootstrapTemplate is an AWS CloudFormation template to bootstrap
// IAM policies, users and roles for use by Cluster API Provider AWS
func BootstrapTemplate(accountID, partition string, extraControlPlanePolicies, extraNodePolicies []string) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[ControllersPolicy] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("controllers"),
		Description:       `For the Kubernetes Cluster API Provider AWS Controllers`,
		PolicyDocument:    controllersPolicy(accountID, partition),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControllers"),
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources[ControlPlanePolicy] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("control-plane"),
		Description:       `For the Kubernetes Cloud Provider AWS Control Plane`,
		PolicyDocument:    cloudProviderControlPlaneAwsPolicy(),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources[NodePolicy] = &cfn_iam.ManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("nodes"),
		Description:       `For the Kubernetes Cloud Provider AWS nodes`,
		PolicyDocument:    nodePolicy(partition),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
			cloudformation.Ref("AWSIAMRoleNodes"),
		},
	}

	template.Resources["AWSIAMUserBootstrapper"] = &cfn_iam.User{
		UserName: iam.NewManagedName("bootstrapper"),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
	}

	template.Resources["AWSIAMGroupBootstrapper"] = &cfn_iam.Group{
		GroupName: iam.NewManagedName("bootstrapper"),
	}

	template.Resources["AWSIAMRoleControlPlane"] = &cfn_iam.Role{
		RoleName:                 iam.NewManagedName("control-plane"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
		ManagedPolicyArns:        extraControlPlanePolicies,
	}

	template.Resources["AWSIAMRoleControllers"] = &cfn_iam.Role{
		RoleName:                 iam.NewManagedName("controllers"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
	}

	template.Resources["AWSIAMRoleNodes"] = &cfn_iam.Role{
		RoleName:                 iam.NewManagedName("nodes"),
		AssumeRolePolicyDocument: ec2AssumeRolePolicy(),
		ManagedPolicyArns:        extraNodePolicies,
	}

	template.Resources["AWSIAMInstanceProfileControlPlane"] = &cfn_iam.InstanceProfile{
		InstanceProfileName: iam.NewManagedName("control-plane"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources["AWSIAMInstanceProfileControllers"] = &cfn_iam.InstanceProfile{
		InstanceProfileName: iam.NewManagedName("controllers"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControllers"),
		},
	}

	template.Resources["AWSIAMInstanceProfileNodes"] = &cfn_iam.InstanceProfile{
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

func controllersPolicy(accountID, partition string) *iam.PolicyDocument {
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
					"ec2:ModifyVpcAttribute",
					"ec2:DeleteInternetGateway",
					"ec2:DeleteNatGateway",
					"ec2:DeleteRouteTable",
					"ec2:DeleteSecurityGroup",
					"ec2:DeleteSubnet",
					"ec2:DeleteTags",
					"ec2:DeleteVpc",
					"ec2:DescribeAccountAttributes",
					"ec2:DescribeAddresses",
					"ec2:DescribeAvailabilityZones",
					"ec2:DescribeInstances",
					"ec2:DescribeInternetGateways",
					"ec2:DescribeImages",
					"ec2:DescribeNatGateways",
					"ec2:DescribeNetworkInterfaces",
					"ec2:DescribeNetworkInterfaceAttribute",
					"ec2:DescribeRouteTables",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeSubnets",
					"ec2:DescribeVpcs",
					"ec2:DescribeVpcAttribute",
					"ec2:DescribeVolumes",
					"ec2:DetachInternetGateway",
					"ec2:DisassociateRouteTable",
					"ec2:DisassociateAddress",
					"ec2:ModifyInstanceAttribute",
					"ec2:ModifyNetworkInterfaceAttribute",
					"ec2:ModifySubnetAttribute",
					"ec2:ReleaseAddress",
					"ec2:RevokeSecurityGroupIngress",
					"ec2:RunInstances",
					"ec2:TerminateInstances",
					"tag:GetResources",
					"elasticloadbalancing:AddTags",
					"elasticloadbalancing:CreateLoadBalancer",
					"elasticloadbalancing:ConfigureHealthCheck",
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:DescribeLoadBalancers",
					"elasticloadbalancing:DescribeLoadBalancerAttributes",
					"elasticloadbalancing:DescribeTags",
					"elasticloadbalancing:ModifyLoadBalancerAttributes",
					"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
					"elasticloadbalancing:RemoveTags",
				},
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{fmt.Sprintf(
					"arn:%s:iam::%s:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing",
					partition,
					accountID,
				)},
				Action: iam.Actions{
					"iam:CreateServiceLinkedRole",
				},
				Condition: iam.Conditions{
					"StringLike": map[string]string{"iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"},
				},
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{fmt.Sprintf(
					"arn:%s:iam::%s:role/%s",
					partition,
					accountID,
					iam.NewManagedName("*"),
				)},
				Action: iam.Actions{
					"iam:PassRole",
				},
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{fmt.Sprintf(
					"arn:%s:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*",
					partition,
				)},
				Action: iam.Actions{
					"secretsmanager:CreateSecret",
					"secretsmanager:DeleteSecret",
					"secretsmanager:TagResource",
				},
			},
		},
	}
}

func bootstrapSecretPolicy(partition string) iam.StatementEntry {
	return iam.StatementEntry{
		Effect: iam.EffectAllow,
		Resource: iam.Resources{fmt.Sprintf(
			"arn:%s:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*",
			partition,
		)},
		Action: iam.Actions{
			"secretsmanager:DeleteSecret",
			"secretsmanager:GetSecretValue",
		},
	}
}

func sessionManagerPolicy() iam.StatementEntry {
	return iam.StatementEntry{
		Effect:   iam.EffectAllow,
		Resource: iam.Resources{"*"},
		Action: iam.Actions{
			"ssm:UpdateInstanceInformation",
			"ssmmessages:CreateControlChannel",
			"ssmmessages:CreateDataChannel",
			"ssmmessages:OpenControlChannel",
			"ssmmessages:OpenDataChannel",
			"s3:GetEncryptionConfiguration",
		},
	}
}

func nodePolicy(partition string) *iam.PolicyDocument {
	policyDocument := cloudProviderNodeAwsPolicy()
	policyDocument.Statement = append(
		policyDocument.Statement,
		bootstrapSecretPolicy(partition),
		sessionManagerPolicy(),
	)
	return policyDocument
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
					"ec2:DescribeImages",
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

func getPolicyDocFromPolicyName(policyName, accountID, partition string) (*iam.PolicyDocument, error) {
	switch policyName {
	case ControllersPolicy:
		return controllersPolicy(accountID, partition), nil
	case ControlPlanePolicy:
		return cloudProviderControlPlaneAwsPolicy(), nil
	case NodePolicy:
		return cloudProviderNodeAwsPolicy(), nil
	}
	return nil, fmt.Errorf("PolicyName %q did not match with any ManagedIAMPolicy", policyName)
}

// GenerateManagedIAMPolicyDocuments generates JSON representation of policy documents for all ManagedIAMPolicy
func (s *Service) GenerateManagedIAMPolicyDocuments(policyDocDir, accountID, partition string) error {
	for _, pn := range ManagedIAMPolicyNames {
		pd, err := getPolicyDocFromPolicyName(pn, accountID, partition)
		if err != nil {
			return fmt.Errorf("failed to get PolicyDocument for ManagedIAMPolicy %q, %v", pn, err)
		}

		pds, err := pd.JSON()
		if err != nil {
			return fmt.Errorf("failed to marshal policy document for ManagedIAMPolicy %q: %v", pn, err)
		}

		fn := path.Join(policyDocDir, fmt.Sprintf("%s.json", pn))
		err = ioutil.WriteFile(fn, []byte(pds), 0755)
		if err != nil {
			return fmt.Errorf("failed to generate policy document for ManagedIAMPolicy %q: %v", pn, err)
		}
	}
	return nil
}

// ReconcileBootstrapStack creates or updates bootstrap CloudFormation
func (s *Service) ReconcileBootstrapStack(stackName, accountID, partition string, extraControlPlanePolicies, extraNodePolicies []string) error {

	template := BootstrapTemplate(accountID, partition, extraControlPlanePolicies, extraNodePolicies)
	yaml, err := template.YAML()
	processedYaml := string(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	if err := s.createStack(stackName, processedYaml); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == "AlreadyExistsException" {
			klog.Infof("AWS Cloudformation stack %q already exists, updating", stackName)
			updateErr := s.updateStack(stackName, processedYaml)
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
