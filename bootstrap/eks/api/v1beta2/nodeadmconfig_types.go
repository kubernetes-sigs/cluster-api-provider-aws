package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NodeadmConfigSpec defines the desired state of NodeadmConfig.
type NodeadmConfigSpec struct {
	// Kubelet contains options for kubelet.
	// +optional
	Kubelet *KubeletOptions `json:"kubelet,omitempty"`

	// Containerd contains options for containerd.
	// +optional
	Containerd *ContainerdOptions `json:"containerd,omitempty"`

	// FeatureGates holds key-value pairs to enable or disable application features.
	// +optional
	FeatureGates map[Feature]bool `json:"featureGates,omitempty"`

	// PreNodeadmCommands specifies extra commands to run before bootstrapping nodes.
	// +optional
	PreNodeadmCommands []string `json:"preNodeadmCommands,omitempty"`

	// Files specifies extra files to be passed to user_data upon creation.
	// +optional
	Files []File `json:"files,omitempty"`

	// Users specifies extra users to add.
	// +optional
	Users []User `json:"users,omitempty"`

	// NTP specifies NTP configuration.
	// +optional
	NTP *NTP `json:"ntp,omitempty"`

	// DiskSetup specifies options for the creation of partition tables and file systems on devices.
	// +optional
	DiskSetup *DiskSetup `json:"diskSetup,omitempty"`

	// Mounts specifies a list of mount points to be setup.
	// +optional
	Mounts []MountPoints `json:"mounts,omitempty"`
}

// KubeletOptions are additional parameters passed to kubelet.
type KubeletOptions struct {
	// Config is a KubeletConfiguration that will be merged with the defaults.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Config *runtime.RawExtension `json:"config,omitempty"`

	// Flags are command-line kubelet arguments that will be appended to the defaults.
	// +optional
	Flags []string `json:"flags,omitempty"`
}

// ContainerdOptions are additional parameters passed to containerd.
type ContainerdOptions struct {
	// Config is an inline containerd configuration TOML that will be merged with the defaults.
	// +optional
	Config string `json:"config,omitempty"`

	// BaseRuntimeSpec is the OCI runtime specification upon which all containers will be based.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	BaseRuntimeSpec *runtime.RawExtension `json:"baseRuntimeSpec,omitempty"`
}

// Feature specifies which feature gate should be toggled.
// +kubebuilder:validation:Enum=InstanceIdNodeName;FastImagePull
type Feature string

const (
	// FeatureInstanceIDNodeName  will use EC2 instance ID as node name.
	FeatureInstanceIDNodeName Feature = "InstanceIdNodeName"
	// FeatureFastImagePull enables a parallel image pull for container images.
	FeatureFastImagePull Feature = "FastImagePull"
)

// GetConditions returns the observations of the operational state of the NodeadmConfig resource.
func (r *NodeadmConfig) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the NodeadmConfig to the predescribed clusterv1.Conditions.
func (r *NodeadmConfig) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// NodeadmConfigStatus defines the observed state of NodeadmConfig.
type NodeadmConfigStatus struct {
	// Ready indicates the BootstrapData secret is ready to be consumed.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`

	// FailureReason will be set on non-retryable errors.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set on non-retryable errors.
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions defines current service state of the NodeadmConfig.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
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
