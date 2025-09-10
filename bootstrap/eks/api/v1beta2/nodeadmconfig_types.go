package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeadmConfigSpec defines the desired state of NodeadmConfig.
type NodeadmConfigSpec struct {
	// Kubelet contains options for kubelet.
	// +optional
	Kubelet *KubeletOptions `json:"kubelet,omitempty"`

	// Containerd contains options for containerd.
	// +optional
	Containerd *ContainerdOptions `json:"containerd,omitempty"`

	// Instance contains options for the node's operating system and devices.
	// +optional
	Instance *InstanceOptions `json:"instance,omitempty"`

	// FeatureGates holds key-value pairs to enable or disable application features.
	// +optional
	FeatureGates map[Feature]bool `json:"featureGates,omitempty"`

	// PreBootstrapCommands specifies extra commands to run before bootstrapping nodes.
	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

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

// InstanceOptions determines how the node's operating system and devices are configured.
type InstanceOptions struct {
	// LocalStorage contains options for configuring EC2 instance stores.
	// +optional
	LocalStorage *LocalStorageOptions `json:"localStorage,omitempty"`
}

// LocalStorageOptions control how EC2 instance stores are used when available.
type LocalStorageOptions struct {
	// Strategy specifies how to handle an instance's local storage devices.
	Strategy LocalStorageStrategy `json:"strategy"`

	// MountPath is the path where the filesystem will be mounted.
	// Defaults to "/mnt/k8s-disks/".
	// +optional
	MountPath string `json:"mountPath,omitempty"`

	// DisabledMounts is a list of directories that will not be mounted to LocalStorage.
	// By default, all mounts are enabled.
	// +optional
	DisabledMounts []DisabledMount `json:"disabledMounts,omitempty"`
}

// Feature specifies which feature gate should be toggled.
// +kubebuilder:validation:Enum=InstanceIdNodeName;FastImagePull
type Feature string

const (
	FeatureInstanceIdNodeName Feature = "InstanceIdNodeName"
	FeatureFastImagePull      Feature = "FastImagePull"
)

// LocalStorageStrategy specifies how to handle an instance's local storage devices.
// +kubebuilder:validation:Enum=RAID0;RAID10;Mount
type LocalStorageStrategy string

const (
	RAID0Strategy  LocalStorageStrategy = "RAID0"
	RAID10Strategy LocalStorageStrategy = "RAID10"
	MountStrategy  LocalStorageStrategy = "Mount"
)

// DisabledMount specifies a directory that should not be mounted onto local storage.
// +kubebuilder:validation:Enum=Containerd;PodLogs
type DisabledMount string

const (
	DisabledMountContainerd DisabledMount = "Containerd"
	DisabledMountPodLogs    DisabledMount = "PodLogs"
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
