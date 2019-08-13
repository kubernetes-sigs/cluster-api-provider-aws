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
)

// AWSMachineTemplateSpec defines the desired state of AWSMachineTemplate
type AWSMachineTemplateSpec struct {
	Template AWSMachineTemplateResource `json:"template"`
}

// AWSMachineTemplateResource describes the data needed to create am AWSMachine from a template
type AWSMachineTemplateResource struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	clusterv1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec AWSMachineSpec `json:"spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineTemplate is the Schema for the awsmachinetemplates API
// +k8s:openapi-gen=true
type AWSMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSMachineTemplateSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineTemplateList contains a list of AWSMachineTemplate
type AWSMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSMachineTemplate{}, &AWSMachineTemplateList{})
}
