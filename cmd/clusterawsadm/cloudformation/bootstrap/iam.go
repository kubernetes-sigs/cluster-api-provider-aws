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
	"io/ioutil"
	"path"

	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/converters"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/iam/v1alpha1"
)

type PolicyName string

// ManagedIAMPolicyNames slice of managed IAM policies
var ManagedIAMPolicyNames = [4]PolicyName{ControllersPolicy, ControlPlanePolicy, NodePolicy, CSIPolicy}

func (p PolicyName) IsValid() bool {
	for i := range ManagedIAMPolicyNames {
		if ManagedIAMPolicyNames[i] == p {
			return true
		}
	}
	return false
}

// GenerateManagedIAMPolicyDocuments generates JSON representation of policy documents for all ManagedIAMPolicy
func (t Template) GenerateManagedIAMPolicyDocuments(policyDocDir string) error {
	for _, pn := range ManagedIAMPolicyNames {
		pd := t.GetPolicyDocFromPolicyName(pn)

		pds, err := converters.IAMPolicyDocumentToJSON(*pd)
		if err != nil {
			return fmt.Errorf("failed to marshal policy document for ManagedIAMPolicy %q: %w", pn, err)
		}

		fn := path.Join(policyDocDir, fmt.Sprintf("%s.json", pn))
		err = ioutil.WriteFile(fn, []byte(pds), 0o600)
		if err != nil {
			return fmt.Errorf("failed to generate policy document for ManagedIAMPolicy %q: %w", pn, err)
		}
	}
	return nil
}

func (t Template) policyFunctionMap() map[PolicyName]func() *iamv1.PolicyDocument {
	return map[PolicyName]func() *iamv1.PolicyDocument{
		ControlPlanePolicy: t.cloudProviderControlPlaneAwsPolicy,
		ControllersPolicy:  t.ControllersPolicy,
		NodePolicy:         t.cloudProviderNodeAwsPolicy,
		CSIPolicy:          t.csiControllerPolicy,
	}
}

func (t Template) GetPolicyDocFromPolicyName(policyName PolicyName) *iamv1.PolicyDocument {
	return t.policyFunctionMap()[policyName]()
}
