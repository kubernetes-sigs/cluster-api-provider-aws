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

func (t Template) csiControlPlaneAwsRoles() []string {
	roles := []string{}
	if !t.Spec.ControlPlane.EnableCSIPolicy {
		roles = append(roles, cloudformation.Ref(AWSIAMRoleControlPlane))
	}
	return roles
}

// From https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/master/docs/example-iam-policy.json
func (t Template) csiControllerPolicy() *v1alpha4.PolicyDocument {
	return &v1alpha4.PolicyDocument{
		Version: v1alpha4.CurrentVersion,
		Statement: []v1alpha4.StatementEntry{
			{
				Effect:   v1alpha4.EffectAllow,
				Resource: v1alpha4.Resources{v1alpha4.Any},
				Action: v1alpha4.Actions{
					"ec2:AttachVolume",
					"ec2:CreateSnapshot",
					"ec2:CreateTags",
					"ec2:CreateVolume",
					"ec2:DeleteSnapshot",
					"ec2:DeleteTags",
					"ec2:DeleteVolume",
					"ec2:DescribeAvailabilityZones",
					"ec2:DescribeInstances",
					"ec2:DescribeSnapshots",
					"ec2:DescribeTags",
					"ec2:DescribeVolumes",
					"ec2:DescribeVolumesModifications",
					"ec2:DetachVolume",
					"ec2:ModifyVolume",
				},
			},
		},
	}
}
