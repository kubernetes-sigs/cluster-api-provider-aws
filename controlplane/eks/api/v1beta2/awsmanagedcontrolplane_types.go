/*
Copyright 2022 The Kubernetes Authors.

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

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ManagedControlPlaneFinalizer allows the controller to clean up resources on delete.
	ManagedControlPlaneFinalizer = "awsmanagedcontrolplane.controlplane.cluster.x-k8s.io"

	// AWSManagedControlPlaneKind is the Kind of AWSManagedControlPlane.
	AWSManagedControlPlaneKind = "AWSManagedControlPlane"
)

// AWSManagedControlPlaneSpec defines the desired state of an Amazon EKS Cluster.
type AWSManagedControlPlaneSpec struct { //nolint: maligned
	AWSManagedControlPlaneClassSpec `json:",inline"`

	// EKSClusterName allows you to specify the name of the EKS cluster in
	// AWS. If you don't specify a name then a default name will be created
	// based on the namespace and name of the managed control plane.
	// +optional
	EKSClusterName string `json:"eksClusterName,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// KubeProxy specifies how the kube-proxy daemonset is managed.
type KubeProxy struct {
	// Disable set to true indicates that kube-proxy should be disabled. With EKS clusters
	// kube-proxy is automatically installed into the cluster. For clusters where you want
	// to use kube-proxy functionality that is provided with an alternate CNI, this option
	// provides a way to specify that the kube-proxy daemonset should be deleted. You cannot
	// set this to true if you are using the Amazon kube-proxy addon.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`
}

// VpcCni specifies configuration related to the VPC CNI.
type VpcCni struct {
	// Disable indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
	// Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
	// to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
	// should be deleted. You cannot set this to true if you are using the
	// Amazon VPC CNI addon.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`
	// Env defines a list of environment variables to apply to the `aws-node` DaemonSet
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// EndpointAccess specifies how control plane endpoints are accessible.
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

// EncryptionConfig specifies the encryption configuration for the EKS clsuter.
type EncryptionConfig struct {
	// Provider specifies the ARN or alias of the CMK (in AWS KMS)
	Provider *string `json:"provider,omitempty"`
	// Resources specifies the resources to be encrypted
	Resources []*string `json:"resources,omitempty"`
}

// OIDCProviderStatus holds the status of the AWS OIDC identity provider.
type OIDCProviderStatus struct {
	// ARN holds the ARN of the provider
	ARN string `json:"arn,omitempty"`
	// TrustPolicy contains the boilerplate IAM trust policy to use for IRSA
	TrustPolicy string `json:"trustPolicy,omitempty"`
}

type IdentityProviderStatus struct {
	// ARN holds the ARN of associated identity provider
	ARN string `json:"arn,omitempty"`

	// Status holds current status of associated identity provider
	Status string `json:"status,omitempty"`
}

// AWSManagedControlPlaneStatus defines the observed state of an Amazon EKS Cluster.
type AWSManagedControlPlaneStatus struct {
	// Networks holds details about the AWS networking resources used by the control plane
	// +optional
	Network infrav1.NetworkStatus `json:"networkStatus,omitempty"`
	// FailureDomains specifies a list fo available availability zones that can be used
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	// Bastion holds details of the instance that is used as a bastion jump box
	// +optional
	Bastion *infrav1.Instance `json:"bastion,omitempty"`
	// OIDCProvider holds the status of the identity provider for this cluster
	// +optional
	OIDCProvider OIDCProviderStatus `json:"oidcProvider,omitempty"`
	// ExternalManagedControlPlane indicates to cluster-api that the control plane
	// is managed by an external service such as AKS, EKS, GKE, etc.
	// +kubebuilder:default=true
	ExternalManagedControlPlane *bool `json:"externalManagedControlPlane,omitempty"`
	// Initialized denotes whether or not the control plane has the
	// uploaded kubernetes config-map.
	// +optional
	Initialized bool `json:"initialized"`
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
	// Addons holds the current status of the EKS addons
	// +optional
	Addons []AddonState `json:"addons,omitempty"`
	// IdentityProviderStatus holds the status for
	// associated identity provider
	// +optional
	IdentityProviderStatus IdentityProviderStatus `json:"identityProviderStatus,omitempty"`
	// Version defines the Kubernetes version for the control plane instance.
	// +optional
	Version string `json:"version"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedcontrolplanes,shortName=awsmcp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSManagedControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.network.vpc.id",description="AWS VPC the control plane is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSManagedControlPlane is the schema for the Amazon EKS Managed Control Plane API.
type AWSManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status AWSManagedControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedControlPlaneList contains a list of Amazon EKS Managed Control Planes.
type AWSManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedControlPlane `json:"items"`
}

// GetConditions returns the control planes conditions.
func (r *AWSManagedControlPlane) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the AWSManagedControlPlane.
func (r *AWSManagedControlPlane) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSManagedControlPlane{}, &AWSManagedControlPlaneList{})
}
