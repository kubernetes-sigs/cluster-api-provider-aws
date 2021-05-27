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
	"github.com/awslabs/goformation/v4/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

func (t Template) cloudProviderControlPlaneAwsRoles() []string {
	roles := []string{}
	if !t.Spec.ControlPlane.DisableCloudProviderPolicy {
		roles = append(roles, cloudformation.Ref(AWSIAMRoleControlPlane))
	}
	return roles
}

// From https://github.com/kubernetes/cloud-provider-aws
func (t Template) cloudProviderControlPlaneAwsPolicy() *v1alpha4.PolicyDocument {
	return &v1alpha4.PolicyDocument{
		Version: v1alpha4.CurrentVersion,
		Statement: []v1alpha4.StatementEntry{
			{
				Effect:   v1alpha4.EffectAllow,
				Resource: v1alpha4.Resources{v1alpha4.Any},
				Action: v1alpha4.Actions{
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
