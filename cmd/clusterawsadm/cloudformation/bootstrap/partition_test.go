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

package bootstrap

import (
	"testing"

	. "github.com/onsi/gomega"

	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/bootstrap/v1beta1"
)

func TestEKSRolePoliciesUsePartition(t *testing.T) {
	tests := []struct {
		partition  string
		wantPrefix string
	}{
		{partition: "aws", wantPrefix: "arn:aws:iam::aws:policy/"},
		{partition: bootstrapv1.PartitionNameUSGov, wantPrefix: "arn:aws-us-gov:iam::aws:policy/"},
		{partition: "aws-cn", wantPrefix: "arn:aws-cn:iam::aws:policy/"},
	}
	for _, tc := range tests {
		t.Run(tc.partition, func(t *testing.T) {
			g := NewWithT(t)
			tmpl := NewTemplate()
			tmpl.Spec.Partition = tc.partition

			nodegroup := tmpl.eksMachinePoolPolicies()
			g.Expect(nodegroup).To(HaveLen(3))
			for _, p := range nodegroup {
				g.Expect(p).To(HavePrefix(tc.wantPrefix))
			}

			fargate := tmpl.fargateProfilePolicies(&bootstrapv1.AWSIAMRoleSpec{})
			g.Expect(fargate).To(ConsistOf(tc.wantPrefix + "AmazonEKSFargatePodExecutionRolePolicy"))
		})
	}
}
