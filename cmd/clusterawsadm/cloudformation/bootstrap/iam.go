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
	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

// PolicyName defines the name of a managed IAM policy.
type PolicyName string

// ManagedIAMPolicyNames slice of managed IAM policies.
var ManagedIAMPolicyNames = []PolicyName{ControllersPolicy, ControllersPolicyEKS, ControlPlanePolicy, NodePolicy, CSIPolicy}

// IsValid will check if a given policy name is valid. That is, it will check if the given policy name is
// one of the ManagedIAMPolicyNames.
func (p PolicyName) IsValid() bool {
	for i := range ManagedIAMPolicyNames {
		if ManagedIAMPolicyNames[i] == p {
			return true
		}
	}
	return false
}

// RenderManagedIAMPolicies returns all the managed IAM Policies that would be rendered by the template.
func (t Template) RenderManagedIAMPolicies() map[string]*iam.ManagedPolicy {
	cft := t.RenderCloudFormation()

	return cft.GetAllIAMManagedPolicyResources()
}

// RenderManagedIAMPolicy returns a specific managed IAM Policy by name, or nil if the policy is not found.
func (t Template) RenderManagedIAMPolicy(name PolicyName) *iam.ManagedPolicy {
	cft := t.RenderCloudFormation()

	p, err := cft.GetIAMManagedPolicyWithName(string(name))
	if err != nil {
		// Return error only if the policy is not found.
		return nil
	}
	return p
}
