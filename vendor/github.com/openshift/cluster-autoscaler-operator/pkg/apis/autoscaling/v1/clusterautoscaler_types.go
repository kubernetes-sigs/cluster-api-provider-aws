package v1

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
// +kubebuilder:subresource:status
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
	// +kubebuilder:validation:Minimum=0
	MaxNodesTotal *int32         `json:"maxNodesTotal,omitempty"`
	Cores         *ResourceRange `json:"cores,omitempty"`
	Memory        *ResourceRange `json:"memory,omitempty"`
	GPUS          []GPULimit     `json:"gpus,omitempty"`
}

type GPULimit struct {
	// +kubebuilder:validation:MinLength=1
	Type string `json:"type"`

	// +kubebuilder:validation:Minimum=0
	Min int32 `json:"min"`
	// +kubebuilder:validation:Minimum=1
	Max int32 `json:"max"`
}

type ResourceRange struct {
	// +kubebuilder:validation:Minimum=0
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}

type ScaleDownConfig struct {
	Enabled bool `json:"enabled"`
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterAdd *string `json:"delayAfterAdd,omitempty"`
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterDelete *string `json:"delayAfterDelete,omitempty"`
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	DelayAfterFailure *string `json:"delayAfterFailure,omitempty"`
	// +kubebuilder:validation:Pattern=([0-9]*(\.[0-9]*)?[a-z]+)+
	UnneededTime *string `json:"unneededTime,omitempty"`
}
