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
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

func (t Template) bootstrapUserGroups() []string {
	groups := []string{
		cloudformation.Ref(AWSIAMGroupBootstrapper),
	}
	groups = append(groups, t.Spec.BootstrapUser.ExtraGroups...)
	return groups
}

func (t Template) bootstrapUserPolicy() []cfn_iam.User_Policy {
	userPolicies := []cfn_iam.User_Policy{}
	if t.Spec.BootstrapUser.ExtraStatements != nil {
		userPolicies = append(userPolicies,
			cfn_iam.User_Policy{
				PolicyName: t.Spec.StackName,
				PolicyDocument: v1alpha4.PolicyDocument{
					Statement: t.Spec.BootstrapUser.ExtraStatements,
					Version:   v1alpha4.CurrentVersion,
				},
			},
		)
	}
	return userPolicies
}
