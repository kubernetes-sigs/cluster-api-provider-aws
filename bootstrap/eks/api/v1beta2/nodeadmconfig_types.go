package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeadmConfigSpec defines the desired state of NodeadmConfig.
type NodeadmConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NodeadmConfig. Edit nodeadmconfig_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// NodeadmConfigStatus defines the observed state of NodeadmConfig.
type NodeadmConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NodeadmConfig is the Schema for the nodeadmconfigs API.
type NodeadmConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeadmConfigSpec   `json:"spec,omitempty"`
	Status NodeadmConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeadmConfigList contains a list of NodeadmConfig.
type NodeadmConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeadmConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeadmConfig{}, &NodeadmConfigList{})
}
