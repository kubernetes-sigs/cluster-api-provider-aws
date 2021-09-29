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

func (t Template) cloudProviderNodeAwsRoles() []string {
	roles := []string{}
	if !t.Spec.ControlPlane.DisableCloudProviderPolicy {
		roles = append(roles, cloudformation.Ref(AWSIAMRoleControlPlane))
	}
	if !t.Spec.Nodes.DisableCloudProviderPolicy {
		roles = append(roles, cloudformation.Ref(AWSIAMRoleNodes))
	}

	return roles
}

// From https://github.com/kubernetes/cloud-provider-aws
func (t Template) cloudProviderNodeAwsPolicy() *v1alpha4.PolicyDocument {
	return &v1alpha4.PolicyDocument{
		Version: v1alpha4.CurrentVersion,
		Statement: []v1alpha4.StatementEntry{
			{
				Effect:   v1alpha4.EffectAllow,
				Resource: v1alpha4.Resources{v1alpha4.Any},
				Action: v1alpha4.Actions{
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
