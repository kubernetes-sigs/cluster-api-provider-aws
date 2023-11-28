/*
Copyright 2023 The Kubernetes Authors.
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSManagedControlPlaneTemplateSpec defines the desired state of AWSManagedControlPlaneTemplate.
type AWSManagedControlPlaneTemplateSpec struct {
	Template AWSManagedControlPlaneTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedcontrolplanetemplatess,shortName=awsmcpt,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

// AWSManagedControlPlaneTemplate is the Schema for the AWSManagedControlPlaneTemplates API.
type AWSManagedControlPlaneTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSManagedControlPlaneTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedControlPlaneTemplateList contains a list of AWSManagedControlPlaneTemplates.
type AWSManagedControlPlaneTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedControlPlaneTemplate `json:"items"`
}

// AWSManagedControlPlaneTemplateResource describes the data needed to create an AWSManagedControlPlane from a template.
type AWSManagedControlPlaneTemplateResource struct {
	Spec AWSManagedControlPlaneTemplateResourceSpec `json:"spec"`
}

func init() {
	SchemeBuilder.Register(&AWSManagedControlPlaneTemplate{}, &AWSManagedControlPlaneTemplateList{})
}
