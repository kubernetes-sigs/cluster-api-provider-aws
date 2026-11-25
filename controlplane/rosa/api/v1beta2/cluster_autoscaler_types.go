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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterAutoscaler configures the cluster-wide behaviour of the cluster autoscaler of a ROSA HCP cluster.
//
// The cluster autoscaler always exists on a ROSA HCP cluster. When this configuration is set,
// the controller owns the whole autoscaler configuration: fields that are omitted are reconciled
// to the ROSA HCP default values, and changes made outside of CAPA (for example with the ROSA CLI)
// are overwritten.
//
// Machine pool level limits (minimum and maximum replicas) are configured on each machine pool.
type ClusterAutoscaler struct {
	// resourceLimits configures cluster-wide resource limits for the cluster autoscaler.
	// +optional
	ResourceLimits *ClusterAutoscalerResourceLimits `json:"resourceLimits,omitempty"`

	// maxPodGracePeriod is the maximum time, in seconds, the cluster autoscaler waits for pods to
	// terminate gracefully before removing a node during scale down.
	// Must be 0 or at least 600. Defaults to 600 when omitted.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:XValidation:rule="self == 0 || self >= 600",message="maxPodGracePeriod must be 0 or at least 600"
	// +optional
	MaxPodGracePeriod *int32 `json:"maxPodGracePeriod,omitempty"`

	// maxNodeProvisionTime is the maximum time the cluster autoscaler waits for a node to be provisioned
	// before considering the scale up failed, for example "15m".
	// Must be between 15m and 60m. Defaults to 15m when omitted.
	// +kubebuilder:validation:XValidation:rule="duration(self) >= duration('15m') && duration(self) <= duration('60m')",message="maxNodeProvisionTime must be between 15m and 60m"
	// +optional
	MaxNodeProvisionTime *metav1.Duration `json:"maxNodeProvisionTime,omitempty"`

	// podPriorityThreshold is the priority below which pods are considered expendable: they don't trigger
	// a scale up and don't prevent a node from being removed during scale down.
	// Must be at most 1000000000, the highest priority allowed for user-defined PriorityClasses.
	// Defaults to -10 when omitted.
	// +kubebuilder:validation:Maximum=1000000000
	// +optional
	PodPriorityThreshold *int32 `json:"podPriorityThreshold,omitempty"`
}

// ClusterAutoscalerResourceLimits configures cluster-wide resource limits for the cluster autoscaler.
type ClusterAutoscalerResourceLimits struct {
	// maxNodesTotal is the maximum number of worker nodes across all machine pools.
	// The cluster autoscaler does not scale up beyond this value.
	// Must be 0 (no limit) or at least 2. Defaults to 0 when omitted.
	// The total cluster size is also limited by ROSA and is validated by OCM.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:XValidation:rule="self == 0 || self >= 2",message="maxNodesTotal must be 0 or at least 2"
	// +optional
	MaxNodesTotal *int32 `json:"maxNodesTotal,omitempty"`
}

// ClusterAutoscalerStatus reports the cluster autoscaler configuration currently stored in OCM,
// including default values and values set outside of CAPA.
// It intentionally has no validation so that unexpected values reported by OCM never block status updates.
type ClusterAutoscalerStatus struct {
	// resourceLimits reports the cluster-wide resource limits.
	// +optional
	ResourceLimits *ClusterAutoscalerResourceLimitsStatus `json:"resourceLimits,omitempty"`

	// maxPodGracePeriod is the maximum time, in seconds, given to pods to terminate during scale down.
	// +optional
	MaxPodGracePeriod *int32 `json:"maxPodGracePeriod,omitempty"`

	// maxNodeProvisionTime is the maximum time to wait for a node to be provisioned.
	// +optional
	MaxNodeProvisionTime *metav1.Duration `json:"maxNodeProvisionTime,omitempty"`

	// podPriorityThreshold is the priority below which pods are considered expendable.
	// +optional
	PodPriorityThreshold *int32 `json:"podPriorityThreshold,omitempty"`
}

// ClusterAutoscalerResourceLimitsStatus reports the cluster-wide resource limits currently stored in OCM.
type ClusterAutoscalerResourceLimitsStatus struct {
	// maxNodesTotal is the maximum number of worker nodes across all machine pools. 0 means no limit.
	// +optional
	MaxNodesTotal *int32 `json:"maxNodesTotal,omitempty"`
}
