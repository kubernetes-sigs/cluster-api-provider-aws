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
			name: "standard tier",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTierStandard,
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTierStandard,
			},
		},
		{
			name: "tier-xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTierXL,
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-xl"),
			},
		},
		{
			name: "tier-2xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier2XL,
			},
			expected: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
		},
		{
			name: "tier-4xl",
			input: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier4XL,
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
				Tier: ekscontrolplanev1.ControlPlaneScalingTierStandard,
			},
		},
		{
			name: "tier-xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTierXL,
			},
		},
		{
			name: "tier-2xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-2xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier2XL,
			},
		},
		{
			name: "tier-4xl",
			input: &ekstypes.ControlPlaneScalingConfig{
				Tier: ekstypes.ProvisionedControlPlaneTier("tier-4xl"),
			},
			expected: &ekscontrolplanev1.ControlPlaneScalingConfig{
				Tier: ekscontrolplanev1.ControlPlaneScalingTier4XL,
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

func TestTaintsFromSDK(t *testing.T) {
	tests := []struct {
		name      string
		input     []ekstypes.Taint
		expected  expinfrav1.Taints
		expectErr bool
	}{
		{
			name: "single taint",
			input: []ekstypes.Taint{
				{Key: aws.String("key1"), Value: aws.String("value1"), Effect: ekstypes.TaintEffectNoSchedule},
			},
			expected: expinfrav1.Taints{
				{Key: "key1", Value: "value1", Effect: expinfrav1.TaintEffectNoSchedule},
			},
		},
		{
			name: "multiple taints",
			input: []ekstypes.Taint{
				{Key: aws.String("key1"), Value: aws.String("value1"), Effect: ekstypes.TaintEffectNoExecute},
				{Key: aws.String("key2"), Value: aws.String("value2"), Effect: ekstypes.TaintEffectPreferNoSchedule},
			},
			expected: expinfrav1.Taints{
				{Key: "key1", Value: "value1", Effect: expinfrav1.TaintEffectNoExecute},
				{Key: "key2", Value: "value2", Effect: expinfrav1.TaintEffectPreferNoSchedule},
			},
		},
		{
			name: "nil value",
			input: []ekstypes.Taint{
				{Key: aws.String("key1"), Value: nil, Effect: ekstypes.TaintEffectNoSchedule},
			},
			expected: expinfrav1.Taints{
				{Key: "key1", Value: "", Effect: expinfrav1.TaintEffectNoSchedule},
			},
		},
		{
			name: "nil key",
			input: []ekstypes.Taint{
				{Key: nil, Value: aws.String("value1"), Effect: ekstypes.TaintEffectNoSchedule},
			},
			expected: expinfrav1.Taints{
				{Key: "", Value: "value1", Effect: expinfrav1.TaintEffectNoSchedule},
			},
		},
		{
			name:     "empty input",
			input:    []ekstypes.Taint{},
			expected: expinfrav1.Taints{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: expinfrav1.Taints{},
		},
		{
			name: "unknown effect",
			input: []ekstypes.Taint{
				{Key: aws.String("key1"), Value: aws.String("value1"), Effect: ekstypes.TaintEffect("UnknownEffect")},
			},
			expectErr: true,
		},
		{
			name: "mixed valid and invalid taints",
			input: []ekstypes.Taint{
				{Key: aws.String("key1"), Value: aws.String("value1"), Effect: ekstypes.TaintEffectNoExecute},
				{Key: aws.String("key2"), Value: aws.String("value2"), Effect: ekstypes.TaintEffect("InvalidEffect")},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TaintsFromSDK(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("TaintsFromSDK() expected an error, got nil")
				}
				if result != nil {
					t.Errorf("TaintsFromSDK() expected nil result on error, got %v", result)
				}
				return
			}
			if err != nil {
				t.Fatalf("TaintsFromSDK() unexpected error: %v", err)
			}
			if !cmp.Equal(result, tt.expected) {
				t.Errorf("TaintsFromSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result))
			}
		})
	}
}
