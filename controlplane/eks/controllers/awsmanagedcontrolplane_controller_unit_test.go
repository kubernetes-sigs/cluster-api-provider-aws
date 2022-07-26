/*
Copyright 2022 The Kubernetes Authors.

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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

func TestSecurityGroupRolesForCluster(t *testing.T) {
	tests := []struct {
		name           string
		bastionEnabled bool
	}{
		{
			name:           "Should use bastion security group when bastion is enabled",
			bastionEnabled: true,
		},
		{
			name:           "Should not use bastion security group when bastion is disabled",
			bastionEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			c := getAWSManagedControlPlane("test", "test")
			c.Spec.Bastion.Enabled = tt.bastionEnabled
			s, err := getManagedControlPlaneScope(c)
			g.Expect(err).To(BeNil(), "failed to create cluster scope for test")

			got := securityGroupRolesForControlPlane(s)
			if tt.bastionEnabled {
				g.Expect(got).To(ContainElement(infrav1.SecurityGroupBastion))
			} else {
				g.Expect(got).ToNot(ContainElement(infrav1.SecurityGroupBastion))
			}

			// Verify that function does not modify the package-level variable.
			gotAgain := securityGroupRolesForControlPlane(s)
			g.Expect(gotAgain).To(BeEquivalentTo(got), "two identical calls return different values")
		})
	}
}
