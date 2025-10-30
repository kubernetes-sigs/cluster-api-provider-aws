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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=awsmanagedmachinepooltemplates,scope=Namespaced,categories=cluster-api,shortName=awsmmpt

// AWSManagedMachinePoolTemplate is the Schema for the awsmanagedmachinepooltemplates API.
type AWSManagedMachinePoolTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSManagedMachinePoolTemplateSpec `json:"spec,omitempty"`
}

//type AWSManagedMachinePoolTemplateResource struct {
//	Spec AWSManagedMachinePoolSpec `json:"spec"`
//}

// AWSManagedMachinePoolTemplateSpec defines the desired state of AWSManagedMachinePoolTemplate.
type AWSManagedMachinePoolTemplateSpec struct {
	Template *AWSManagedMachinePool `json:"template"`
}

// +kubebuilder:object:root=true

// AWSManagedMachinePoolTemplateList contains a list of AWSManagedMachinePoolTemplates.
type AWSManagedMachinePoolTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedMachinePoolTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSManagedMachinePoolTemplate{}, &AWSManagedMachinePoolTemplateList{})
}
