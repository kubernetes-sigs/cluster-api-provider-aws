/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileAWSMachine to clean up AWS resources associated with AWSMachine before
	// removing it from the apiserver.
	MachineFinalizer = "awsmachine.infrastructure.cluster.x-k8s.io"
)

// AWSMachineSpec defines the desired state of AWSMachine
type AWSMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	AMI AWSResourceReference `json:"ami,omitempty"`

	// ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// InstanceType is the type of instance to create. Example: m4.xlarge
	InstanceType string `json:"instanceType,omitempty"`

	// AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
	// AWSMachine's value takes precedence.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// IAMInstanceProfile is a name of an IAM instance profile to assign to the instance
	// +optional
	IAMInstanceProfile string `json:"iamInstanceProfile,omitempty"`

	// PublicIP specifies whether the instance should get a public IP.
	// Precedence for this setting is as follows:
	// 1. This field if set
	// 2. Cluster/flavor setting
	// 3. Subnet default
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`

	// AdditionalSecurityGroups is an array of references to security groups that should be applied to the
	// instance. These security groups would be set in addition to any security groups defined
	// at the cluster level or in the actuator.
	// +optional
	AdditionalSecurityGroups []AWSResourceReference `json:"additionalSecurityGroups,omitempty"`

	// AvailabilityZone is references the AWS availability zone to use for this instance.
	// If multiple subnets are matched for the availability zone, the first one return is picked.
	// +optional
	AvailabilityZone *string `json:"availabilityZone,omitempty"`

	// Subnet is a reference to the subnet to use for this instance. If not specified,
	// the cluster subnet will be used.
	// +optional
	Subnet *AWSResourceReference `json:"subnet,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the instance.
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// RootDeviceSize is the size of the root volume in gigabytes(GB).
	// +optional
	RootDeviceSize int64 `json:"rootDeviceSize,omitempty"`

	// NetworkInterfaces is a list of ENIs to associate with the instance.
	// A maximum of 2 may be specified.
	// +optional
	// +kubebuilder:validation:MaxItems=2
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`

	// CloudInit defines options related to the bootstrapping systems where
	// CloudInit is used.
	// +optional
	CloudInit *CloudInit `json:"cloudInit,omitempty"`
}

// CloudInit defines options related to the bootstrapping systems where
// CloudInit is used.
type CloudInit struct {
	// enableSecureSecretsManager, when set to true will use AWS Secrets Manager to ensure
	// userdata privacy. A cloud-init boothook shell script is prepended to download
	// the userdata from Secrets Manager and additionally delete the secret.
	// +optional
	EnableSecureSecretsManager bool `json:"enableSecureSecretsManager,omitempty"`

	// SecretCount is the number of secrets used to form the complete secret
	// +optional
	SecretCount int32 `json:"secretCount,omitempty"`

	// SecretPrefix is the prefix for the secret name. This is stored
	// temporarily, and deleted when the machine registers as a node against
	// the workload cluster.
	// +optional
	SecretPrefix string `json:"secretPrefix,omitempty"`
}

// AWSMachineStatus defines the observed state of AWSMachine
type AWSMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the AWS instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// InstanceState is the state of the AWS instance for this machine.
	// +optional
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	// ErrorReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorReason *errors.MachineStatusError `json:"errorReason,omitempty"`

	// ErrorMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// AWSMachine is the Schema for the awsmachines API
type AWSMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSMachineSpec   `json:"spec,omitempty"`
	Status AWSMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSMachineList contains a list of AWSMachine
type AWSMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSMachine{}, &AWSMachineList{})
}
