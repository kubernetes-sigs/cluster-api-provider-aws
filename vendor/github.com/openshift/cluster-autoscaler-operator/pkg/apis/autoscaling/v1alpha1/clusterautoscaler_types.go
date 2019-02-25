package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&ClusterAutoscaler{}, &ClusterAutoscalerList{})
}

// ClusterAutoscalerSpec defines the desired state of ClusterAutoscaler
type ClusterAutoscalerSpec struct {
	ResourceLimits       *ResourceLimits  `json:"resourceLimits,omitempty"`
	ScaleDown            *ScaleDownConfig `json:"scaleDown,omitempty"`
	MaxPodGracePeriod    *int32           `json:"maxPodGracePeriod,omitempty"`
	PodPriorityThreshold *int32           `json:"podPriorityThreshold,omitempty"`
}

// ClusterAutoscalerStatus defines the observed state of ClusterAutoscaler
type ClusterAutoscalerStatus struct {
	// TODO: Add status fields.
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAutoscaler is the Schema for the clusterautoscalers API
// +k8s:openapi-gen=true
type ClusterAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterAutoscalerSpec   `json:"spec,omitempty"`
	Status ClusterAutoscalerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAutoscalerList contains a list of ClusterAutoscaler
type ClusterAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterAutoscaler `json:"items"`
}

type ResourceLimits struct {
	MaxNodesTotal *int32         `json:"maxNodesTotal,omitempty"`
	Cores         *ResourceRange `json:"cores,omitempty"`
	Memory        *ResourceRange `json:"memory,omitempty"`
	GPUS          []GPULimit     `json:"gpus,omitempty"`
}

type GPULimit struct {
	Type string `json:"type"`
	ResourceRange
}

type ResourceRange struct {
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}

type ScaleDownConfig struct {
	Enabled           bool   `json:"enabled"`
	DelayAfterAdd     string `json:"delayAfterAdd"`
	DelayAfterDelete  string `json:"delayAfterDelete"`
	DelayAfterFailure string `json:"delayAfterFailure"`
	UnneededTime      string `json:"unneededTime,omitempty"`
}
