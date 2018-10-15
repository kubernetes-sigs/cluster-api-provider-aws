/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSClusterProviderConfigSpec defines the desired state of AWSClusterProviderConfig
type AWSClusterProviderConfigSpec struct {
	// The AWS Region the cluster lives in.
	Region string `json:"region"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host.
	SSHKeyName string `json:"sshKeyName,omitempty"`
}

// AWSClusterProviderConfigStatus defines the observed state of AWSClusterProviderConfig
type AWSClusterProviderConfigStatus struct {
	// Region is the AWS region the cluster is deployed into.
	Region string `json:"region"`

	// Network is a reference to the network infrastructure that is created for this cluster.
	Network Network `json:"network,omitempty"`

	// Bastion is a reference to the bastion host for this cluster.
	Bastion Instance `json:"bastion,omitempty"`

	// CACertificate is a PEM encoded CA Certificate for the control plane nodes.
	CACertificate []byte

	// CAPrivateKey is a PEM encoded PKCS1 CA PrivateKey for the control plane nodes.
	CAPrivateKey []byte
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSClusterProviderConfig is the Schema for the awsclusterproviderconfigs API
// +k8s:openapi-gen=true
type AWSClusterProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSClusterProviderConfigSpec   `json:"spec,omitempty"`
	Status AWSClusterProviderConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSClusterProviderConfigList contains a list of AWSClusterProviderConfig
type AWSClusterProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSClusterProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSClusterProviderConfig{}, &AWSClusterProviderConfigList{})
}
