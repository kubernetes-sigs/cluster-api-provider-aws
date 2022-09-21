/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSMachineTemplateSpec defines the desired state of AWSMachineTemplate
type AWSMachineTemplateSpec struct {
	Template AWSMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmachinetemplates,scope=Namespaced,categories=cluster-api,shortName=awsmt

// AWSMachineTemplate is the Schema for the awsmachinetemplates API
type AWSMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSMachineTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSMachineTemplateList contains a list of AWSMachineTemplate.
type AWSMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSMachineTemplate{}, &AWSMachineTemplateList{})
}
