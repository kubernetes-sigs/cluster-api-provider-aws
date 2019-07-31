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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capierrors "sigs.k8s.io/cluster-api/pkg/errors"
)

const (
	ClusterFinalizer                   = "cluster.cluster.x-k8s.io"
	ClusterAnnotationControlPlaneReady = "cluster.x-k8s.io/control-plane-ready"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

/// [Cluster]
// Cluster is the Schema for the clusters API
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=clusters,shortName=cl
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

/// [Cluster]

/// [ClusterSpec]
// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// Cluster network configuration
	// +optional
	ClusterNetwork *ClusterNetworkingConfig `json:"clusterNetwork,omitempty"`

	// InfrastructureRef is a reference to a provider-specific resource that holds the details
	// for provisioning infrastructure for a cluster in said provider.
	// +optional
	InfrastructureRef *corev1.ObjectReference `json:"infrastructureRef,omitempty"`
}

/// [ClusterSpec]

/// [ClusterNetworkingConfig]
// ClusterNetworkingConfig specifies the different networking
// parameters for a cluster.
type ClusterNetworkingConfig struct {
	// The network ranges from which service VIPs are allocated.
	Services NetworkRanges `json:"services"`

	// The network ranges from which Pod networks are allocated.
	Pods NetworkRanges `json:"pods"`

	// Domain name for services.
	ServiceDomain string `json:"serviceDomain"`
}

/// [ClusterNetworkingConfig]

/// [NetworkRanges]
// NetworkRanges represents ranges of network addresses.
type NetworkRanges struct {
	CIDRBlocks []string `json:"cidrBlocks"`
}

/// [NetworkRanges]

/// [ClusterStatus]
// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// APIEndpoints represents the endpoints to communicate with the control plane.
	// +optional
	APIEndpoints []APIEndpoint `json:"apiEndpoints,omitempty"`

	// ErrorReason indicates that there is a problem reconciling the
	// state, and will be set to a token value suitable for
	// programmatic interpretation.
	// +optional
	ErrorReason *capierrors.ClusterStatusError `json:"errorReason,omitempty"`

	// ErrorMessage indicates that there is a problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// Phase represents the current phase of cluster actuation.
	// E.g. Pending, Running, Terminating, Failed etc.
	// +optional
	Phase *string `json:"phase,omitempty"`

	// InfrastructureReady is the state of the infrastructure provider.
	// +optional
	InfrastructureReady bool `json:"infrastructureReady"`
}

// SetTypedPhase sets the Phase field to the string representation of ClusterPhase.
func (c *ClusterStatus) SetTypedPhase(p ClusterPhase) {
	c.Phase = ClusterPhaseStringPtr(p)
}

// GetTypedPhase attempts to parse the Phase field and return
// the typed ClusterPhase representation as described in `machine_phase_types.go`.
func (c *ClusterStatus) GetTypedPhase() ClusterPhase {
	if c.Phase == nil {
		return ClusterPhaseUnknown
	}

	switch phase := ClusterPhase(*c.Phase); phase {
	case
		ClusterPhasePending,
		ClusterPhaseProvisioning,
		ClusterPhaseProvisioned,
		ClusterPhaseDeleting,
		ClusterPhaseFailed:
		return phase
	default:
		return ClusterPhaseUnknown
	}
}

/// [ClusterStatus]

/// [APIEndpoint]
// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int `json:"port"`
}

/// [APIEndpoint]

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
