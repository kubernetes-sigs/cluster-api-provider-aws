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
	"k8s.io/utils/ptr"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
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

func TestControlPlaneScalingConfigToSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *ekscontrolplanev1.ControlPlaneScalingConfig
		expected *ekstypes.ControlPlaneScalingConfig
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "nil tier returns nil",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: nil,
			},
			expected: nil,
		},
		{
			name: "standard tier",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTierStandard),
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTierStandard,
			},
		},
		{
			name: "tier-xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTierXL),
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-xl"),
			},
		},
		{
			name: "tier-2xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTier2XL),
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
		},
		{
			name: "tier-4xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTier4XL),
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-4xl"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ControlPlaneScalingConfigToSDK(tt.input)
			if !cmp.Equal(result, tt.expected, cmpopts.IgnoreUnexported(ekstypes.ControlPlaneScalingConfig{})) {
				t.Errorf("ControlPlaneScalingConfigToSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result, cmpopts.IgnoreUnexported(ekstypes.ControlPlaneScalingConfig{})))
			}
		})
	}
}

func TestControlPlaneScalingConfigFromSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *ekstypes.ControlPlaneScalingConfig
		expected *ekscontrolplanev1.ControlPlaneScalingConfig
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "standard tier",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTierStandard,
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTierStandard),
			},
		},
		{
			name: "tier-xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTierXL),
			},
		},
		{
			name: "tier-2xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTier2XL),
			},
		},
		{
			name: "tier-4xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-4xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ptr.To(ekscontrolplanev1.ControlPlaneScalingTier4XL),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ControlPlaneScalingConfigFromSDK(tt.input)
			if !cmp.Equal(result, tt.expected) {
				t.Errorf("ControlPlaneScalingConfigFromSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result))
			}
		})
	}
}
