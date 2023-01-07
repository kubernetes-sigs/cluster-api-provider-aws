/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSGDifference(t *testing.T) {
	tests := []struct {
		name     string
		self     IngressRules
		input    IngressRules
		expected IngressRules
	}{
		{
			name:     "self and input are nil",
			self:     nil,
			input:    nil,
			expected: nil,
		},
		{
			name: "input is nil",
			self: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
			input: nil,
			expected: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
		},
		{
			name: "self has more rules",
			self: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
				{
					Description: "MY-SSH",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
			},
			input: IngressRules{
				{
					Description:            "SSH",
					Protocol:               SecurityGroupProtocolTCP,
					FromPort:               22,
					ToPort:                 22,
					SourceSecurityGroupIDs: []string{"sg-source-1"},
				},
			},
			expected: IngressRules{
				{
					Description: "MY-SSH",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    22,
					ToPort:      22,
					CidrBlocks:  []string{"0.0.0.0/0"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			out := tc.self.Difference(tc.input)

			g.Expect(out).To(Equal(tc.expected))
		})
	}
}
