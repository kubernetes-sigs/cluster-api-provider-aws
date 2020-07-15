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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

const (
	// ManagedClusterFinalizer allows ReconcileAWSManagedCluster to clean up AWS resources associated with AWSManagedCluster before
	// removing it from the apiserver.
	ManagedClusterFinalizer = "awsmanagedcluster.exp.infrastructure.cluster.x-k8s.io"
)

// AWSManagedClusterSpec defines the desired state of AWSManagedCluster
type AWSManagedClusterSpec struct {
	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec infrav1.NetworkSpec `json:"networkSpec,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// Bastion contains options to configure the bastion host.
	// +optional
	Bastion infrav1.Bastion `json:"bastion"`
}

// AWSManagedClusterStatus defines the observed state of AWSManagedCluster
type AWSManagedClusterStatus struct {
	// Ready is true when the provider resources are ready and also
	// the AWSManagedControlPlane has a API server URL .
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Initialized is true when the cluster infrastructure has been created. This will be true
	// before the EKS control plane is ready and is used to signal
	// that its safe for the AWSManagedControlPlane to be reconciled.
	Initialized bool `json:"initialized,omitempty"`

	Network        infrav1.Network          `json:"network,omitempty"`
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	Bastion        *infrav1.Instance        `json:"bastion,omitempty"`
	Conditions     clusterv1.Conditions     `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedclusters,scope=Namespaced,categories=cluster-api,shortName=awsmc
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for worker nodes"
// +kubebuilder:printcolumn:name="Initialized",type="string",JSONPath=".status.initialized",description="Cluster infrastructure is ready for EKS control plane"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.networkSpec.vpc.id",description="AWS VPC the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.apiEndpoints[0]",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSManagedCluster is the Schema for the awsmanagedclusters API
type AWSManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedClusterSpec   `json:"spec,omitempty"`
	Status AWSManagedClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedClusterList contains a list of AWSManagedCluster
type AWSManagedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedCluster `json:"items"`
}

// GetConditions gets the status conditions for the AWSManagedCluster
func (r *AWSManagedCluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the AWSManagedCluster
func (r *AWSManagedCluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSManagedCluster{}, &AWSManagedClusterList{})
}
