/*
Copyright 2026 The Kubernetes Authors.

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

package eks

import (
	"testing"

	. "github.com/onsi/gomega"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/endpoints"
)

func TestRolePoliciesForPartition(t *testing.T) {
	tests := []struct {
		region        string
		wantNodegroup []string
		wantFargate   []string
	}{
		{
			region: "eu-west-1",
			wantNodegroup: []string{
				"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
				"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
				"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			},
			wantFargate: []string{"arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"},
		},
		{
			region: "us-gov-west-1",
			wantNodegroup: []string{
				"arn:aws-us-gov:iam::aws:policy/AmazonEKSWorkerNodePolicy",
				"arn:aws-us-gov:iam::aws:policy/AmazonEKS_CNI_Policy",
				"arn:aws-us-gov:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			},
			wantFargate: []string{"arn:aws-us-gov:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"},
		},
		{
			region: "cn-north-1",
			wantNodegroup: []string{
				"arn:aws-cn:iam::aws:policy/AmazonEKSWorkerNodePolicy",
				"arn:aws-cn:iam::aws:policy/AmazonEKS_CNI_Policy",
				"arn:aws-cn:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			},
			wantFargate: []string{"arn:aws-cn:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.region, func(t *testing.T) {
			g := NewWithT(t)
			partition := endpoints.GetPartitionFromRegion(tc.region)
			g.Expect(NodegroupRolePoliciesForPartition(partition)).To(Equal(tc.wantNodegroup))
			g.Expect(FargateRolePoliciesForPartition(partition)).To(Equal(tc.wantFargate))
		})
	}

	g := NewWithT(t)
	g.Expect(NodegroupRolePolicies()).To(Equal(tests[0].wantNodegroup))
	g.Expect(NodegroupRolePoliciesUSGov()).To(Equal(tests[1].wantNodegroup))
	g.Expect(FargateRolePolicies()).To(Equal(tests[0].wantFargate))
	g.Expect(FargateRolePoliciesUSGov()).To(Equal(tests[1].wantFargate))
}
