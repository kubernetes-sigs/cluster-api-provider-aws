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
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks"
)

func fargateProfilePolicies(roleSpec *bootstrapv1.AWSIAMRoleSpec) []string {
	policies := eks.FargateRolePolicies()
	if roleSpec.ExtraPolicyAttachments != nil {
		policies = append(policies, roleSpec.ExtraPolicyAttachments...)
	}

	return policies
}
