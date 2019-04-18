package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&MachineAutoscaler{}, &MachineAutoscalerList{})
}

// MachineAutoscalerSpec defines the desired state of MachineAutoscaler
type MachineAutoscalerSpec struct {
	// +kubebuilder:validation:Minimum=0
	MinReplicas int32 `json:"minReplicas"`
	// +kubebuilder:validation:Minimum=1
	MaxReplicas    int32                       `json:"maxReplicas"`
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef"`
}

// MachineAutoscalerStatus defines the observed state of MachineAutoscaler
type MachineAutoscalerStatus struct {
	LastTargetRef *CrossVersionObjectReference `json:"lastTargetRef,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoscaler is the Schema for the machineautoscalers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type MachineAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineAutoscalerSpec   `json:"spec,omitempty"`
	Status MachineAutoscalerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoscalerList contains a list of MachineAutoscaler
type MachineAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineAutoscaler `json:"items"`
}

// CrossVersionObjectReference identifies another object by name, API version,
// and kind.
type CrossVersionObjectReference struct {
	APIVersion string `json:"apiVersion,omitempty"`
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind"`
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}
