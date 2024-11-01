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

package addons

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestAddOnEqual(t *testing.T) {
	ptr := func(s string) *string { return &s }
	tags := func(s, t string) map[string]string {
		return map[string]string{
			s: t,
		}
	}
	g := gomega.NewGomegaWithT(t)
	tests := []struct {
		orig        *EKSAddon
		other       *EKSAddon
		result      gomega.OmegaMatcher
		includeTags bool
	}{
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
				ResolveConflict:       ptr("d"),
				Tags:                  tags("a", "1"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
				ResolveConflict:       ptr("d"),
				Tags:                  tags("a", "1"),
				Status:                ptr("e"),
			},
			result: gomega.BeTrueBecause("addon values are equal (except status)"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
			},
			other: &EKSAddon{
				Version:               ptr("b"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
			},
			result: gomega.BeFalseBecause("addon version differs"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("c"),
				Configuration:         ptr("c"),
			},
			result: gomega.BeFalseBecause("addon serviceAccountRoleARN differs"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("d"),
			},
			result: gomega.BeFalseBecause("addon configuration differs"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("c"),
				ResolveConflict:       ptr("d"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Configuration:         ptr("d"),
				ResolveConflict:       ptr("e"),
			},
			result: gomega.BeFalseBecause("addon conflict resolution differs"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Tags:                  tags("a", "1"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Tags:                  tags("a", "2"),
			},
			result: gomega.BeTrueBecause("addon tags differ but not used for comparison"),
		},
		{
			orig: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Tags:                  tags("a", "1"),
			},
			other: &EKSAddon{
				Version:               ptr("a"),
				ServiceAccountRoleARN: ptr("b"),
				Tags:                  tags("a", "2"),
			},
			result:      gomega.BeFalseBecause("addon tags differ and used for comparison"),
			includeTags: true,
		},
	}

	for _, test := range tests {
		g.Expect(test.orig.IsEqual(test.other, test.includeTags)).To(test.result)
	}
}
