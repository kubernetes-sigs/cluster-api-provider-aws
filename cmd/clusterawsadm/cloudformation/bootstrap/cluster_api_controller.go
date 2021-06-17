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

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

const (
	eksClusterPolicyName = "AmazonEKSClusterPolicy"
)

func (t Template) controllersPolicyGroups() []string {
	groups := []string{}
	if t.Spec.BootstrapUser.Enable {
		groups = append(groups, cloudformation.Ref(AWSIAMGroupBootstrapper))
	}
	return groups
}

func (t Template) controllersPolicyRoleAttachments() []string {
	attachments := []string{
		cloudformation.Ref(AWSIAMRoleControllers),
	}
	if !t.Spec.ControlPlane.DisableClusterAPIControllerPolicyAttachment {
		attachments = append(
			attachments,
			cloudformation.Ref(AWSIAMRoleControlPlane),
		)
	}
	return attachments
}

func (t Template) controllersTrustPolicy() *infrav1.PolicyDocument {
	policyDocument := ec2AssumeRolePolicy()
	policyDocument.Statement = append(policyDocument.Statement, t.Spec.ClusterAPIControllers.TrustStatements...)
	return policyDocument
}

func (t Template) controllersRolePolicy() []cfn_iam.Role_Policy {
	policies := []cfn_iam.Role_Policy{}

	if t.Spec.ClusterAPIControllers.ExtraStatements != nil {
		policies = append(policies,
			cfn_iam.Role_Policy{
				PolicyName: t.Spec.StackName,
				PolicyDocument: infrav1.PolicyDocument{
					Statement: t.Spec.ClusterAPIControllers.ExtraStatements,
					Version:   infrav1.CurrentVersion,
				},
			},
		)
	}
	return policies
}

// ControllersPolicy will create a policy from a Template for AWS Controllers.
func (t Template) ControllersPolicy() *infrav1.PolicyDocument {
	statement := []infrav1.StatementEntry{
		{
			Effect:   infrav1.EffectAllow,
			Resource: infrav1.Resources{infrav1.Any},
			Action: infrav1.Actions{
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
				"elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
				"elasticloadbalancing:DescribeTags",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:RegisterInstancesWithLoadBalancer",
				"elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
				"elasticloadbalancing:RemoveTags",
				"autoscaling:DescribeAutoScalingGroups",
				"autoscaling:DescribeInstanceRefreshes",
				"ec2:CreateLaunchTemplate",
				"ec2:CreateLaunchTemplateVersion",
				"ec2:DescribeLaunchTemplates",
				"ec2:DescribeLaunchTemplateVersions",
				"ec2:DeleteLaunchTemplate",
				"ec2:DeleteLaunchTemplateVersions",
				"ec2:DescribeKeyPairs",
			},
		},
		{
			Effect: infrav1.EffectAllow,
			Resource: infrav1.Resources{
				"arn:*:autoscaling:*:*:autoScalingGroup:*:autoScalingGroupName/*",
			},
			Action: infrav1.Actions{
				"autoscaling:CreateAutoScalingGroup",
				"autoscaling:UpdateAutoScalingGroup",
				"autoscaling:CreateOrUpdateTags",
				"autoscaling:StartInstanceRefresh",
				"autoscaling:DeleteAutoScalingGroup",
				"autoscaling:DeleteTags",
			},
		},
		{
			Effect: infrav1.EffectAllow,
			Resource: infrav1.Resources{
				"arn:*:iam::*:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling",
			},
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "autoscaling.amazonaws.com"},
			},
		},
		{
			Effect: infrav1.EffectAllow,
			Resource: infrav1.Resources{
				"arn:*:iam::*:role/aws-service-role/elasticloadbalancing.amazonaws.com/AWSServiceRoleForElasticLoadBalancing",
			},
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"},
			},
		},
		{
			Effect: infrav1.EffectAllow,
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Resource: infrav1.Resources{
				"arn:*:iam::*:role/aws-service-role/spot.amazonaws.com/AWSServiceRoleForEC2Spot",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "spot.amazonaws.com"},
			},
		},
		{
			Effect:   infrav1.EffectAllow,
			Resource: t.allowedEC2InstanceProfiles(),
			Action: infrav1.Actions{
				"iam:PassRole",
			},
		},
	}
	for _, secureSecretBackend := range t.Spec.SecureSecretsBackends {
		switch secureSecretBackend {
		case infrav1.SecretBackendSecretsManager:
			statement = append(statement, infrav1.StatementEntry{
				Effect: infrav1.EffectAllow,
				Resource: infrav1.Resources{
					"arn:*:secretsmanager:*:*:secret:aws.cluster.x-k8s.io/*",
				},
				Action: infrav1.Actions{
					"secretsmanager:CreateSecret",
					"secretsmanager:DeleteSecret",
					"secretsmanager:TagResource",
				},
			})
		case infrav1.SecretBackendSSMParameterStore:
			statement = append(statement, infrav1.StatementEntry{
				Effect: infrav1.EffectAllow,
				Resource: infrav1.Resources{
					"arn:*:ssm:*:*:parameter/cluster.x-k8s.io/*",
				},
				Action: infrav1.Actions{
					"ssm:PutParameter",
					"ssm:DeleteParameter",
					"ssm:AddTagsToResource",
				},
			})
		}
	}
	if t.Spec.EKS.Enable {
		allowedIAMActions := infrav1.Actions{
			"iam:GetRole",
			"iam:ListAttachedRolePolicies",
		}
		statement = append(statement, infrav1.StatementEntry{
			Effect: infrav1.EffectAllow,
			Resource: infrav1.Resources{
				"arn:*:ssm:*:*:parameter/aws/service/eks/optimized-ami/*",
			},
			Action: infrav1.Actions{
				"ssm:GetParameter",
			},
		})

		statement = append(statement, infrav1.StatementEntry{
			Effect: infrav1.EffectAllow,
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Resource: infrav1.Resources{
				"arn:*:iam::*:role/aws-service-role/eks.amazonaws.com/AWSServiceRoleForAmazonEKS",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "eks.amazonaws.com"},
			},
		})

		statement = append(statement, infrav1.StatementEntry{
			Effect: infrav1.EffectAllow,
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Resource: infrav1.Resources{
				"arn:*:iam::*:role/aws-service-role/eks-nodegroup.amazonaws.com/AWSServiceRoleForAmazonEKSNodegroup",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "eks-nodegroup.amazonaws.com"},
			},
		})

		statement = append(statement, infrav1.StatementEntry{
			Effect: infrav1.EffectAllow,
			Action: infrav1.Actions{
				"iam:CreateServiceLinkedRole",
			},
			Resource: infrav1.Resources{
				"arn:aws:iam::*:role/aws-service-role/eks-fargate-pods.amazonaws.com/AWSServiceRoleForAmazonEKSForFargate",
			},
			Condition: infrav1.Conditions{
				infrav1.StringLike: map[string]string{"iam:AWSServiceName": "eks-fargate.amazonaws.com"},
			},
		})

		if t.Spec.EKS.AllowIAMRoleCreation {
			allowedIAMActions = append(allowedIAMActions, infrav1.Actions{
				"iam:DetachRolePolicy",
				"iam:DeleteRole",
				"iam:CreateRole",
				"iam:TagRole",
				"iam:AttachRolePolicy",
			}...)

			statement = append(statement, infrav1.StatementEntry{
				Action: infrav1.Actions{
					"iam:ListOpenIDConnectProviders",
					"iam:CreateOpenIDConnectProvider",
					"iam:AddClientIDToOpenIDConnectProvider",
					"iam:UpdateOpenIDConnectProviderThumbprint",
					"iam:DeleteOpenIDConnectProvider",
				},
				Resource: infrav1.Resources{
					"*",
				},
				Effect: infrav1.EffectAllow,
			})
		}
		statement = append(statement, []infrav1.StatementEntry{
			{
				Action: allowedIAMActions,
				Resource: infrav1.Resources{
					"arn:*:iam::*:role/*",
				},
				Effect: infrav1.EffectAllow,
			}, {
				Action: infrav1.Actions{
					"iam:GetPolicy",
				},
				Resource: infrav1.Resources{
					t.generateAWSManagedPolicyARN(eksClusterPolicyName),
				},
				Effect: infrav1.EffectAllow,
			}, {
				Action: infrav1.Actions{
					"eks:DescribeCluster",
					"eks:ListClusters",
					"eks:CreateCluster",
					"eks:TagResource",
					"eks:UpdateClusterVersion",
					"eks:DeleteCluster",
					"eks:UpdateClusterConfig",
					"eks:UntagResource",
					"eks:UpdateNodegroupVersion",
					"eks:DescribeNodegroup",
					"eks:DeleteNodegroup",
					"eks:UpdateNodegroupConfig",
					"eks:CreateNodegroup",
				},
				Resource: infrav1.Resources{
					"arn:*:eks:*:*:cluster/*",
					"arn:*:eks:*:*:nodegroup/*/*/*",
				},
				Effect: infrav1.EffectAllow,
			}, {
				Action: infrav1.Actions{
					"eks:ListAddons",
					"eks:CreateAddon",
					"eks:DescribeAddonVersions",
					"eks:DescribeAddon",
					"eks:DeleteAddon",
					"eks:UpdateAddon",
					"eks:TagResource",
					"eks:DescribeFargateProfile",
					"eks:CreateFargateProfile",
					"eks:DeleteFargateProfile",
				},
				Resource: infrav1.Resources{
					"*",
				},
				Effect: infrav1.EffectAllow,
			}, {
				Action: infrav1.Actions{
					"iam:PassRole",
				},
				Resource: infrav1.Resources{
					"*",
				},
				Condition: infrav1.Conditions{
					"StringEquals": map[string]string{
						"iam:PassedToService": "eks.amazonaws.com",
					},
				},
				Effect: infrav1.EffectAllow,
			},
		}...)

	}

	if t.Spec.EventBridge.Enable {
		statement = append(statement, infrav1.StatementEntry{
			Effect:   infrav1.EffectAllow,
			Resource: infrav1.Resources{infrav1.Any},
			Action: infrav1.Actions{
				"events:DeleteRule",
				"events:DescribeRule",
				"events:ListTargetsByRule",
				"events:PutRule",
				"events:PutTargets",
				"events:RemoveTargets",
				"sqs:CreateQueue",
				"sqs:DeleteMessage",
				"sqs:DeleteQueue",
				"sqs:GetQueueAttributes",
				"sqs:GetQueueUrl",
				"sqs:ReceiveMessage",
				"sqs:SetQueueAttributes",
			},
		})
	}

	return &infrav1.PolicyDocument{
		Version:   infrav1.CurrentVersion,
		Statement: statement,
	}
}

func (t Template) allowedEC2InstanceProfiles() infrav1.Resources {
	if t.Spec.ClusterAPIControllers.AllowedEC2InstanceProfiles == nil {
		t.Spec.ClusterAPIControllers.AllowedEC2InstanceProfiles = []string{
			t.NewManagedName(infrav1.Any),
		}
	}
	instanceProfiles := make(infrav1.Resources, len(t.Spec.ClusterAPIControllers.AllowedEC2InstanceProfiles))

	for i, p := range t.Spec.ClusterAPIControllers.AllowedEC2InstanceProfiles {
		instanceProfiles[i] = fmt.Sprintf("arn:*:iam::*:role/%s", p)
	}

	return instanceProfiles
}
