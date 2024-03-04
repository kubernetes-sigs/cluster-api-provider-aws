/*
Copyright 2018 The Kubernetes Authors.

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

package scope

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGenerateProviderID(t *testing.T) {
	testCases := []struct {
		ids []string

		expectedProviderID string
	}{
		{
			ids: []string{
				"eu-west-1a",
				"instance-id",
			},
			expectedProviderID: "aws:///eu-west-1a/instance-id",
		},
		{
			ids: []string{
				"eu-west-1a",
				"test-id1",
				"test-id2",
				"instance-id",
			},
			expectedProviderID: "aws:///eu-west-1a/test-id1/test-id2/instance-id",
		},
	}

	for _, tc := range testCases {
		g := NewGomegaWithT(t)
		providerID := GenerateProviderID(tc.ids...)

		g.Expect(providerID).To(Equal(tc.expectedProviderID))
	}
}
