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

	"github.com/awslabs/goformation/cloudformation"
	"github.com/pkg/errors"
	"k8s.io/klog"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/iam"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

const (
	// ControllersPolicy is the IAM Managed Policy for Cluster API AWS Controllers
	ControllersPolicy = "AWSIAMManagedPolicyControllers"
	// ControlPlanePolicy is the IAM Managed Policy for master nodes
	ControlPlanePolicy = "AWSIAMManagedPolicyCloudProviderControlPlane"
	// NodePolicy is the IAM Managed Policy for worker nodes
	NodePolicy = "AWSIAMManagedPolicyCloudProviderNodes"
)

// ManagedIAMPolicyNames slice of managed IAM policies
var ManagedIAMPolicyNames = [...]string{ControllersPolicy, ControlPlanePolicy, NodePolicy}

// BootstrapTemplate is an AWS CloudFormation template to bootstrap
// IAM policies, users and roles for use by Cluster API Provider AWS
func BootstrapTemplate(accountID string) *cloudformation.Template {
	template := cloudformation.NewTemplate()

	template.Resources[ControllersPolicy] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("controllers"),
		Description:       `For the Kubernetes Cluster API Provider AWS Controllers`,
		PolicyDocument:    controllersPolicy(accountID),
		Groups: []string{
			cloudformation.Ref("AWSIAMGroupBootstrapper"),
		},
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControllers"),
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources[ControlPlanePolicy] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("control-plane"),
		Description:       `For the Kubernetes Cloud Provider AWS Control Plane`,
		PolicyDocument:    cloudProviderControlPlaneAwsPolicy(),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControlPlane"),
		},
	}

	template.Resources[NodePolicy] = cloudformation.AWSIAMManagedPolicy{
		ManagedPolicyName: iam.NewManagedName("nodes"),
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

	template.Resources["AWSIAMRoleControllers"] = cloudformation.AWSIAMRole{
		RoleName:                 iam.NewManagedName("controllers"),
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

	template.Resources["AWSIAMInstanceProfileControllers"] = cloudformation.AWSIAMInstanceProfile{
		InstanceProfileName: iam.NewManagedName("controllers"),
		Roles: []string{
			cloudformation.Ref("AWSIAMRoleControllers"),
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

func tagRequestStringEqualityConditions() map[string][]string {
	return map[string][]string{
		fmt.Sprintf("aws:RequestTag/%s", tags.NameAWSProviderManaged): {"true"},
		fmt.Sprintf("aws:RequestTag/%s", tags.NameAWSClusterAPIRole): {
			tags.ValueAPIServerRole,
			tags.ValueBastionRole,
			tags.ValueCommonRole,
			tags.ValuePrivateRole,
			tags.ValuePublicRole,
		},
	}
}

func tagMandatoryPresenceConditions() map[string][]string {
	return map[string][]string{
		"aws:TagKeys": {
			fmt.Sprintf("%s*", tags.NameKubernetesClusterPrefix),
			tags.NameAWSProviderManaged,
			tags.NameAWSClusterAPIRole,
			"Name",
		},
	}
}

func tagExistingObjectPresenceConditions(service string) map[string][]string {
	return map[string][]string{
		fmt.Sprintf("%s:ResourceTag/%s", service, tags.NameAWSProviderManaged): {"true"},
	}
}

func tagRequestForNewObjectConditions() *iam.Conditions {
	return &iam.Conditions{
		"StringEquals":           tagRequestStringEqualityConditions(),
		"ForAnyValue:StringLike": tagMandatoryPresenceConditions(),
		"Bool":                   secureTransportCondition(),
	}
}

func taggedObjectRequestCondition(service string) *iam.Conditions {
	return &iam.Conditions{
		"StringEquals": tagExistingObjectPresenceConditions(service),
		"Bool":         secureTransportCondition(),
	}
}

func tagRequestForExistingObjectConditions(service string) *iam.Conditions {
	equalityTests := tagExistingObjectPresenceConditions(service)
	for k, v := range tagRequestStringEqualityConditions() {
		equalityTests[k] = v
	}
	return &iam.Conditions{
		"StringEquals":           equalityTests,
		"ForAnyValue:StringLike": tagMandatoryPresenceConditions(),
		"Bool":                   secureTransportCondition(),
	}
}

func secureTransportCondition() map[string][]string {
	return map[string][]string{"aws:SecureTransport": {"true"}}
}

func controllersPolicy(accountID string) *iam.PolicyDocument {
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
					"ec2:DescribeImages",
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
				},
				Condition: iam.Conditions{"Bool": secureTransportCondition()},
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{fmt.Sprintf(
					"arn:aws:iam::%s:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing",
					accountID),
				},
				Action: iam.Actions{
					"iam:CreateServiceLinkedRole",
				},
				Condition: iam.Conditions{
					"StringLike": map[string]string{"iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"},
					"Bool":       secureTransportCondition(),
				},
			},
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"ec2:CreateTags",
					"ec2:RunInstances",
					"elasticloadbalancing:CreateLoadBalancer",
				},
				Condition: *tagRequestForNewObjectConditions(),
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{
					"arn:aws:ec2:::instance/*",
				},
				Action: iam.Actions{
					"ec2:CreateTags",
				},
				Condition: *tagRequestForExistingObjectConditions("ec2"),
			},
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"ec2:TerminateInstances",
					"ec2:DescribeInstances",
				},
				Condition: *taggedObjectRequestCondition("ec2"),
			},
			{
				Effect:   iam.EffectAllow,
				Resource: iam.Resources{"*"},
				Action: iam.Actions{
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:ConfigureHealthCheck",
					"elasticloadbalancing:DescribeLoadBalancers",
					"elasticloadbalancing:DescribeLoadBalancerAttributes",
					"elasticloadbalancing:ModifyLoadBalancerAttributes",
					"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
				},
				Condition: *taggedObjectRequestCondition("elasticloadbalancing"),
			},
			{
				Effect: iam.EffectAllow,
				Resource: iam.Resources{fmt.Sprintf(
					"arn:aws:iam::%s:role/%s",
					accountID,
					iam.NewManagedName("*"),
				)},
				Action: iam.Actions{
					"iam:PassRole",
				},
				Condition: iam.Conditions{"Bool": secureTransportCondition()},
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

func getPolicyDocFromPolicyName(policyName, accountID string) (*iam.PolicyDocument, error) {
	switch policyName {
	case ControllersPolicy:
		return controllersPolicy(accountID), nil
	case ControlPlanePolicy:
		return cloudProviderControlPlaneAwsPolicy(), nil
	case NodePolicy:
		return cloudProviderNodeAwsPolicy(), nil
	}
	return nil, fmt.Errorf("PolicyName %q did not match with any ManagedIAMPolicy", policyName)
}

// GenerateManagedIAMPolicyDocuments generates JSON representation of policy documents for all ManagedIAMPolicy
func (s *Service) GenerateManagedIAMPolicyDocuments(policyDocDir, accountID string) error {
	for _, pn := range ManagedIAMPolicyNames {
		pd, err := getPolicyDocFromPolicyName(pn, accountID)
		if err != nil {
			return fmt.Errorf("Error: failed to get PolicyDocument for ManagedIAMPolicy %q, %v", pn, err)
		}

		pds, err := pd.JSON()
		if err != nil {
			return fmt.Errorf("ERROR: failed to marshal policy document for ManagedIAMPolicy %q: %v", pn, err)
		}

		fn := path.Join(policyDocDir, fmt.Sprintf("%s.json", pn))
		err = ioutil.WriteFile(fn, []byte(pds), 0755)
		if err != nil {
			return fmt.Errorf("ERROR: failed to generate policy document for ManagedIAMPolicy %q: %v", pn, err)
		}
	}
	return nil
}

// ReconcileBootstrapStack creates or updates bootstrap CloudFormation
func (s *Service) ReconcileBootstrapStack(stackName string, accountID string) error {

	template := BootstrapTemplate(accountID)
	yaml, err := template.YAML()
	processedYaml := iam.ProcessPolicyDocument(string(yaml))
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
