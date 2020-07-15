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
	// EKSControlPlaneFinalizer allows the controller to clean up resources on delete
	EKSControlPlaneFinalizer = "eks.exp.infrastructure.cluster.x-k8s.io"
)

// AWSManagedControlPlaneSpec defines the desired state of AWSManagedControlPlane
type AWSManagedControlPlaneSpec struct {
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

	// Logging specifies whether a logging type is enabled or disabled
	// +optional
	Logging map[string]bool `json:"logging,omitempty"`

	// EncryptionConfig specifies the encryption configuration for the cluster
	// +optional
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// Private indicates if the control plane should be private
	// +optional
	Private *bool `json:"private,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
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
	// Initialized denotes whether or not the control plane has the
	// uploaded kubeadm-config configmap.
	// +kubebuilder:default=false
	Initialized bool `json:"initialized"`

	// Ready denotes that the AWSManagedControlPlane API Server is ready to
	// receive requests.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// FailureReason indicates that there is a terminal problem reconciling the
	// state, and will be set to a token value suitable for
	// programmatic interpretation.
	// +optional
	//FailureReason errors.KubeadmControlPlaneStatusError `json:"failureReason,omitempty"`

	// ErrorMessage indicates that there is a terminal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedcontrolplanes,shortName=awsmcp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=".status.ready",description="AWSManagedControlPlane API Server is ready to receive requests"
// +kubebuilder:printcolumn:name="Initialized",type=boolean,JSONPath=".status.initialized",description="This denotes whether or not the control plane has the uploaded kubeadm-config configmap"

// AWSManagedControlPlane is the Schema for the awsmanagedcontrolplanes API
type AWSManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status AWSManagedControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedControlPlane contains a list of AWSManagedControlPlane
type AWSManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSManagedControlPlane{}, &AWSManagedControlPlaneList{})
}
