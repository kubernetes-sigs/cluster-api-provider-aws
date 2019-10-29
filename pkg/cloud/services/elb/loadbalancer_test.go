/*
Copyright 2019 The Kubernetes Authors.

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

package elb

import (
	"testing"
)

func TestGenerateELBName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "test",
			expected: "test-apiserver",
		},
		{
			name:     "0123456789012345678901",
			expected: "0123456789012345678901-apiserver",
		},
		{
			name:     "01234567890123456789012",
			expected: "26o3cjil5at5qn27vukn5x09b3ql-k8s",
		},
		{
			name:     "anotherverylongtoolongname",
			expected: "t8gnrbbifaaf5d0k4xmwui3xwvip-k8s",
		},
		{
			name:     "anotherverylongtoolongnameanotherverylongtoolongname",
			expected: "tph1huzox1f10z9ow1inrootjws8-k8s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			elbName, err := GenerateELBName(tt.name)

			if err != nil {
				t.Error(err)
			}

			if elbName != tt.expected {
				t.Errorf("expected ELB name: %v, got name: %v", tt.expected, elbName)
			}

			if len(elbName) > 32 {
				t.Errorf("ELB name too long: %v vs. %s", len(elbName), "32")
			}

		})
	}
}
