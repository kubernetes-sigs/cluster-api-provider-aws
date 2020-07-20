/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// ManagedControlPlaneFinalizer allows the controller to clean up resources on delete
	ManagedControlPlaneFinalizer = "awsmanagedcontrolplane.infrastructure.cluster.x-k8s.io"
)

// AWSManagedControlPlaneSpec defines the desired state of AWSManagedControlPlane
type AWSManagedControlPlaneSpec struct {
	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec infrav1.NetworkSpec `json:"networkSpec,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetesthat EKS supports
	// will be used.
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Pattern:=^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.?$
	// +optional
	Version *string `json:"version"`

	// RoleName specifies the name of IAM role that gives EKS
	// permission to make API calls. If no name is supplied
	// then a role is created. If the role is pre-existing
	// we will treat it as unmanaged and not delete it on
	// deletion
	// +kubebuilder:validation:MinLength:=2
	// +optional
	RoleName *string `json:"roleName"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the control plane role.
	// +optional
	RoleAdditionalPolicies *[]string `json:"roleAdditionalPolicies"`

	// Logging specifies which EKS Cluster logs should be enabled. Entries for
	// each of the enabled logs will be sent to CloudWatch
	// +optional
	Logging *ControlPlaneLoggingSpec `json:"logging,omitempty"`

	// EncryptionConfig specifies the encryption configuration for the cluster
	// +optional
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// Endpoints specifies access to this cluster's control plane endpoints
	// +optional
	EndpointAccess EndpointAccess `json:"endpointAccess,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// Bastion contains options to configure the bastion host.
	// +optional
	Bastion infrav1.Bastion `json:"bastion"`

	// TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
	// iam-authenticator - obtains a client token using iam-authentictor
	// aws-cli - obtains a client token using the AWS CLI
	// Defaults to iam-authenticator
	// +kubebuilder:default=iam-authenticator
	// +kubebuilder:validation:Enum=iam-authenticator;aws-cli
	TokenMethod *EKSTokenMethod `json:"tokenMethod,omitempty"`
}

// EndpointAccess specifies how control plane endpoints are accessible
type EndpointAccess struct {
	// Public controls whether control plane endpoints are publicly accessible
	// +optional
	Public *bool `json:"public,omitempty"`
	// PublicCIDRs specifies which blocks can access the public endpoint
	// +optional
	PublicCIDRs []*string `json:"publicCIDRs,omitempty"`
	// Private points VPC-internal control plane access to the private endpoint
	// +optional
	Private *bool `json:"private,omitempty"`
}

// EncryptionConfig specifies the encryption configuration for the EKS clsuter
type EncryptionConfig struct {
	// Provider specifies the ARN or alias of the CMK (in AWS KMS)
	Provider *string `json:"provider,omitempty"`
	//Resources specifies the resources to be encrypted
	Resources []*string `json:"resources,omitempty"`
}

// AWSManagedControlPlaneStatus defines the observed state of AWSManagedControlPlane
type AWSManagedControlPlaneStatus struct {
	// Networks holds details about the AWS networking resources used by the control plane
	// +optional
	Network infrav1.Network `json:"network,omitempty"`
	// FailureDomains specifies a list fo available availability zones that can be used
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	// Bastion holds details of the instance that is used as a bastion jump box
	// +optional
	Bastion *infrav1.Instance `json:"bastion,omitempty"`
	// Ready denotes that the AWSManagedControlPlane API Server is ready to
	// receive requests and that the VPC infra is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// ErrorMessage indicates that there is a terminal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
	// Conditions specifies the cpnditions for the managed control plane
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedcontrolplanes,shortName=awsmcp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSManagedControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.networkSpec.vpc.id",description="AWS VPC the control plane is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSManagedControlPlane is the Schema for the awsmanagedcontrolplanes API
type AWSManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status AWSManagedControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedControlPlaneList contains a list of AWSManagedControlPlane
type AWSManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedControlPlane `json:"items"`
}

// GetConditions returns the control planes conditions
func (r *AWSManagedControlPlane) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the AWSManagedControlPlane
func (r *AWSManagedControlPlane) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSManagedControlPlane{}, &AWSManagedControlPlaneList{})
}
