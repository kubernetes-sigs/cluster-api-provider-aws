/*
Copyright 2026 The Kubernetes Authors.

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

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// AWSManagedMachinePoolTemplateSpec defines the desired state of AWSManagedMachinePoolTemplate.
type AWSManagedMachinePoolTemplateSpec struct {
	Template AWSManagedMachinePoolTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedmachinepooltemplates,scope=Namespaced,categories=cluster-api,shortName=awsmmpt
// +kubebuilder:storageversion

// AWSManagedMachinePoolTemplate is the schema for the Amazon EKS Managed Machine Pool Templates API.
type AWSManagedMachinePoolTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSManagedMachinePoolTemplateSpec `json:"spec,omitempty"`
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

// AWSManagedMachinePoolTemplateResource describes the data needed to create an AWSManagedMachinePool from a template.
type AWSManagedMachinePoolTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1beta1.ObjectMeta `json:"metadata,omitempty"`
	Spec       AWSManagedMachinePoolSpec `json:"spec"`
}
