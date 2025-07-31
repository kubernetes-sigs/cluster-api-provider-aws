/*
Copyright 2025 The Kubernetes Authors.

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

package converters

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestNodeRepairConfigToSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *expinfrav1.NodeRepairConfig
		expected *ekstypes.NodeRepairConfig
	}{
		{
			name:     "nil input returns default disabled",
			input:    nil,
			expected: &ekstypes.NodeRepairConfig{Enabled: aws.Bool(false)},
		},
		{
			name: "enabled repair config",
			input: &expinfrav1.NodeRepairConfig{
				Enabled: aws.Bool(true),
			},
			expected: &ekstypes.NodeRepairConfig{Enabled: aws.Bool(true)},
		},
		{
			name: "disabled repair config",
			input: &expinfrav1.NodeRepairConfig{
				Enabled: aws.Bool(false),
			},
			expected: &ekstypes.NodeRepairConfig{Enabled: aws.Bool(false)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NodeRepairConfigToSDK(tt.input)
			if !cmp.Equal(result, tt.expected, cmpopts.IgnoreUnexported(ekstypes.NodeRepairConfig{})) {
				t.Errorf("NodeRepairConfigToSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result, cmpopts.IgnoreUnexported(ekstypes.NodeRepairConfig{})))
			}
		})
	}
}
