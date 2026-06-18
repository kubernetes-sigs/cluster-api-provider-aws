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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

const (
	// NodeadmConfigKind is the Kind for the NodeadmConfig resource.
	NodeadmConfigKind = "NodeadmConfig"
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

	// Hybrid contains configuration for EKS Hybrid Nodes.
	// When specified, the NodeadmConfig generates userdata for hybrid node
	// bootstrapping instead of standard EC2 node bootstrapping.
	// Hybrid mode is mutually exclusive with EC2-specific options.
	// +optional
	Hybrid *HybridOptions `json:"hybrid,omitempty"`
}

// HybridOptions defines configuration for EKS Hybrid Nodes.
// When specified, the NodeadmConfig generates userdata for hybrid node
// bootstrapping instead of standard EC2 node bootstrapping.
type HybridOptions struct {
	// SSM configures SSM-based authentication for hybrid nodes.
	// This is required for hybrid node support.
	// +kubebuilder:validation:Required
	SSM HybridSSMOptions `json:"ssm"`

	// CustomUserData allows providing a custom userdata template with variable interpolation.
	// When specified, the default nodeadm MIME multipart userdata generation is completely
	// replaced with the user-provided template. The template uses Go text/template syntax.
	// This is useful for hybrid nodes in environments where cloud-init or the standard
	// bootstrap approach is not suitable.
	// +optional
	CustomUserData *CustomUserDataOptions `json:"customUserData,omitempty"`
}

// CustomUserDataOptions defines a custom userdata template for hybrid nodes.
// When specified, the user-provided template completely replaces the default
// nodeadm MIME multipart userdata generation.
type CustomUserDataOptions struct {
	// Template is a Go text/template format string that will be rendered with
	// available runtime variables. The output is stored as-is in the bootstrap secret.
	//
	// Available template variables:
	//   - {{.ClusterName}}        - EKS cluster name
	//   - {{.Region}}             - AWS region
	//   - {{.KubernetesVersion}}  - Kubernetes version (e.g., "1.34")
	//   - {{.ActivationID}}       - SSM activation ID
	//   - {{.ActivationCode}}     - SSM activation code
	//   - {{.KubeletFlags}}       - []string of kubelet flags (if configured)
	//   - {{.KubeletConfig}}      - string of kubelet config YAML (if configured)
	//   - {{.ContainerdConfig}}   - string of containerd config (if configured)
	//
	// Template functions available:
	//   - {{Indent N .Content}}      - Indent content by N spaces
	//   - {{join .Slice ","}}        - Join slice with delimiter
	//   - {{base64Encode .Data}}     - Base64 encode data
	//   - {{default "val" .Var}}     - Return default if Var is empty
	//   - {{trimSpace .Content}}     - Trim leading/trailing whitespace
	//
	// Example template:
	//   #!/bin/bash
	//   cat <<EOF > /etc/nodeadm/config.yaml
	//   apiVersion: node.eks.aws/v1alpha1
	//   kind: NodeConfig
	//   spec:
	//     cluster:
	//       name: {{.ClusterName}}
	//       region: {{.Region}}
	//     hybrid:
	//       ssm:
	//         activationId: {{.ActivationID}}
	//         activationCode: {{.ActivationCode}}
	//   EOF
	//   /usr/local/bin/nodeadm init -c file:///etc/nodeadm/config.yaml
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Template string `json:"template"`
}

// HybridSSMOptions configures SSM activation-based authentication for hybrid nodes.
// Either ActivationRef or ActivationConfig must be specified, but not both.
// +kubebuilder:validation:XValidation:rule="has(self.activationRef) || has(self.activationConfig)",message="one of activationRef or activationConfig must be specified"
// +kubebuilder:validation:XValidation:rule="!has(self.activationRef) || !has(self.activationConfig)",message="activationRef and activationConfig are mutually exclusive"
type HybridSSMOptions struct {
	// ActivationRef references an existing Secret containing SSM activation credentials.
	// The Secret must contain 'activationId' and 'activationCode' keys.
	// When specified, the controller will not create or manage SSM activations.
	// +optional
	ActivationRef *corev1.LocalObjectReference `json:"activationRef,omitempty"`

	// ActivationConfig specifies parameters for automatically creating an SSM activation.
	// The controller will create the activation and store credentials in a Secret.
	// The activation will be deleted when the NodeadmConfig is deleted.
	// +optional
	ActivationConfig *SSMActivationConfig `json:"activationConfig,omitempty"`
}

// SSMActivationConfig specifies parameters for creating an SSM hybrid activation.
type SSMActivationConfig struct {
	// IAMRoleName is the name of the IAM role that hybrid nodes will assume.
	// This role must have the necessary permissions for EKS hybrid nodes
	// and trust policy allowing ssm.amazonaws.com to assume it.
	// See: https://docs.aws.amazon.com/eks/latest/userguide/hybrid-nodes-creds.html
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	IAMRoleName string `json:"iamRoleName"`

	// RegistrationLimit is the maximum number of hybrid nodes that can register
	// using this activation. Each NodeadmConfig creates its own activation,
	// so this typically corresponds to the number of nodes using this config.
	// Minimum: 1, Maximum: 1000, Default: 1
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +optional
	RegistrationLimit *int32 `json:"registrationLimit,omitempty"`

	// ExpirationHours is the number of hours until the activation expires.
	// After expiration, no new nodes can register using this activation,
	// but already-registered nodes are unaffected.
	// Minimum: 1, Maximum: 720 (30 days), Default: 24
	// +kubebuilder:default=24
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=720
	// +optional
	ExpirationHours *int32 `json:"expirationHours,omitempty"`

	// Tags are key-value pairs to apply to the SSM activation.
	// These tags help identify and manage activations in the AWS console.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
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
func (r *NodeadmConfig) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the NodeadmConfig to the predescribed clusterv1.Conditions.
func (r *NodeadmConfig) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// NodeadmConfigStatus defines the observed state of NodeadmConfig.
type NodeadmConfigStatus struct {
	// Deprecated: This field will be removed with the CAPI v1beta2 transition
	// Ready indicates the BootstrapData secret is ready to be consumed.
	// +optional
	Ready bool `json:"ready,omitempty"`
	// Initialization provides observations of the NodeadmConfig initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Machine provisioning.
	// +optional
	Initialization NodeadmConfigInitializationStatus `json:"initialization,omitempty"`

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
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// SSMActivation contains status information about the SSM activation
	// when hybrid mode is enabled with auto-created activations.
	// +optional
	SSMActivation *SSMActivationStatus `json:"ssmActivation,omitempty"`
}

// SSMActivationStatus contains status information for an auto-created SSM activation.
type SSMActivationStatus struct {
	// ActivationID is the ID of the SSM activation.
	// +optional
	ActivationID *string `json:"activationID,omitempty"`

	// SecretName is the name of the Secret containing activation credentials.
	// +optional
	SecretName *string `json:"secretName,omitempty"`

	// ExpirationTime is when the activation expires.
	// +optional
	ExpirationTime *metav1.Time `json:"expirationTime,omitempty"`
}

// NodeadmConfigInitializationStatus provides observations of the NodeadmConfig initialization process.
type NodeadmConfigInitializationStatus struct {
	// DataSecretCreated is true when the Machine's bootstrap secret is created.
	// NOTE: This field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.
	// +optional
	DataSecretCreated *bool `json:"dataSecretCreated,omitempty"`
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
