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

var componentConfigIgnoreUnexported = cmpopts.IgnoreUnexported(
	ekstypes.KubeSchedulerConfigRequest{},
	ekstypes.NodeResourcesFitConfig{},
	ekstypes.ScoringStrategy{},
	ekstypes.ResourceWeight{},
	ekstypes.KubeApiServerConfigRequest{},
	ekstypes.ServiceNodePortRange{},
	ekstypes.KubeControllerManagerConfigRequest{},
	ekstypes.HorizontalPodAutoscalerControllerConfigRequest{},
	ekstypes.PodGcControllerConfigRequest{},
)

func TestKubeSchedulerConfigToSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *ekscontrolplanev1.KubeSchedulerConfig
		expected *ekstypes.KubeSchedulerConfigRequest
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty config returns nil",
			input:    &ekscontrolplanev1.KubeSchedulerConfig{},
			expected: nil,
		},
		{
			name: "empty node resources fit returns nil",
			input: &ekscontrolplanev1.KubeSchedulerConfig{
				NodeResourcesFit: &ekscontrolplanev1.NodeResourcesFitConfig{},
			},
			expected: nil,
		},
		{
			name: "scoring strategy type only",
			input: &ekscontrolplanev1.KubeSchedulerConfig{
				NodeResourcesFit: &ekscontrolplanev1.NodeResourcesFitConfig{
					ScoringStrategy: &ekscontrolplanev1.ScoringStrategy{
						Type: ekscontrolplanev1.ScoringStrategyTypeMostAllocated,
					},
				},
			},
			expected: &ekstypes.KubeSchedulerConfigRequest{
				NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
					ScoringStrategy: &ekstypes.ScoringStrategy{
						Type: ekstypes.ScoringStrategyTypeMostAllocated,
					},
				},
			},
		},
		{
			name: "scoring strategy with resources",
			input: &ekscontrolplanev1.KubeSchedulerConfig{
				NodeResourcesFit: &ekscontrolplanev1.NodeResourcesFitConfig{
					ScoringStrategy: &ekscontrolplanev1.ScoringStrategy{
						Type: ekscontrolplanev1.ScoringStrategyTypeMostAllocated,
						Resources: []ekscontrolplanev1.ResourceWeight{
							{Name: "cpu", Weight: 2},
							{Name: "memory", Weight: 1},
						},
					},
				},
			},
			expected: &ekstypes.KubeSchedulerConfigRequest{
				NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
					ScoringStrategy: &ekstypes.ScoringStrategy{
						Type: ekstypes.ScoringStrategyTypeMostAllocated,
						Resources: []ekstypes.ResourceWeight{
							{Name: aws.String("cpu"), Weight: aws.Int32(2)},
							{Name: aws.String("memory"), Weight: aws.Int32(1)},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KubeSchedulerConfigToSDK(tt.input)
			if !cmp.Equal(result, tt.expected, componentConfigIgnoreUnexported) {
				t.Errorf("KubeSchedulerConfigToSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result, componentConfigIgnoreUnexported))
			}
		})
	}
}

func TestKubeAPIServerConfigToSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *ekscontrolplanev1.KubeAPIServerConfig
		expected *ekstypes.KubeApiServerConfigRequest
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty config returns nil",
			input:    &ekscontrolplanev1.KubeAPIServerConfig{},
			expected: nil,
		},
		{
			name: "event ttl only",
			input: &ekscontrolplanev1.KubeAPIServerConfig{
				EventTTL: "15m",
			},
			expected: &ekstypes.KubeApiServerConfigRequest{
				EventTtl: aws.String("15m"),
			},
		},
		{
			name: "service node port range only",
			input: &ekscontrolplanev1.KubeAPIServerConfig{
				ServiceNodePortRange: &ekscontrolplanev1.ServiceNodePortRange{MinPort: 20000, MaxPort: 32767},
			},
			expected: &ekstypes.KubeApiServerConfigRequest{
				ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 20000, MaxPort: 32767},
			},
		},
		{
			name: "all parameters",
			input: &ekscontrolplanev1.KubeAPIServerConfig{
				EventTTL:             "30m",
				ServiceNodePortRange: &ekscontrolplanev1.ServiceNodePortRange{MinPort: 20000, MaxPort: 32767},
			},
			expected: &ekstypes.KubeApiServerConfigRequest{
				EventTtl:             aws.String("30m"),
				ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 20000, MaxPort: 32767},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KubeAPIServerConfigToSDK(tt.input)
			if !cmp.Equal(result, tt.expected, componentConfigIgnoreUnexported) {
				t.Errorf("KubeAPIServerConfigToSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result, componentConfigIgnoreUnexported))
			}
		})
	}
}

func TestKubeControllerManagerConfigToSDK(t *testing.T) {
	tests := []struct {
		name     string
		input    *ekscontrolplanev1.KubeControllerManagerConfig
		expected *ekstypes.KubeControllerManagerConfigRequest
	}{
		{
			name:     "nil input returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "empty nested configs return nil",
			input: &ekscontrolplanev1.KubeControllerManagerConfig{
				HorizontalPodAutoscalerControllerConfig: &ekscontrolplanev1.HorizontalPodAutoscalerControllerConfig{},
				PodGCControllerConfig:                   &ekscontrolplanev1.PodGCControllerConfig{},
			},
			expected: nil,
		},
		{
			name: "horizontal pod autoscaler sync period only",
			input: &ekscontrolplanev1.KubeControllerManagerConfig{
				HorizontalPodAutoscalerControllerConfig: &ekscontrolplanev1.HorizontalPodAutoscalerControllerConfig{
					HorizontalPodAutoscalerSyncPeriod: "10s",
				},
			},
			expected: &ekstypes.KubeControllerManagerConfigRequest{
				HorizontalPodAutoscalerControllerConfig: &ekstypes.HorizontalPodAutoscalerControllerConfigRequest{
					HorizontalPodAutoscalerSyncPeriod: aws.String("10s"),
				},
			},
		},
		{
			name: "all parameters",
			input: &ekscontrolplanev1.KubeControllerManagerConfig{
				HorizontalPodAutoscalerControllerConfig: &ekscontrolplanev1.HorizontalPodAutoscalerControllerConfig{
					HorizontalPodAutoscalerSyncPeriod: "10s",
				},
				PodGCControllerConfig: &ekscontrolplanev1.PodGCControllerConfig{
					TerminatedPodGCThreshold: aws.Int32(10000),
				},
			},
			expected: &ekstypes.KubeControllerManagerConfigRequest{
				HorizontalPodAutoscalerControllerConfig: &ekstypes.HorizontalPodAutoscalerControllerConfigRequest{
					HorizontalPodAutoscalerSyncPeriod: aws.String("10s"),
				},
				PodGcControllerConfig: &ekstypes.PodGcControllerConfigRequest{
					TerminatedPodGcThreshold: aws.Int32(10000),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KubeControllerManagerConfigToSDK(tt.input)
			if !cmp.Equal(result, tt.expected, componentConfigIgnoreUnexported) {
				t.Errorf("KubeControllerManagerConfigToSDK() diff (-want +got):\n%s", cmp.Diff(tt.expected, result, componentConfigIgnoreUnexported))
			}
		})
	}
}
