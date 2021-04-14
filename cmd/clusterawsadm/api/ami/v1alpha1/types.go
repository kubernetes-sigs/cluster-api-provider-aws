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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	AWSAMIKind     = "AWSAMI"
	AWSAMIListKind = "AWSAMIList"
)

// AWSAMI defines an AMI
type AWSAMISpec struct {
	OS                string `json:"os"`
	Region            string `json:"region"`
	ImageID           string `json:"imageID"`
	KubernetesVersion string `json:"kubernetesVersion"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// AMI defines an AMI
type AWSAMI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AWSAMISpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// AWSAMIList defines a list of AMIs
type AWSAMIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAMI `json:"items"`
}

func (a *AWSAMIList) ToTable() *metav1.Table {
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{
				Name: "Kubernetes Version",
				Type: "string",
			},
			{
				Name: "Region",
				Type: "string",
			},
			{
				Name: "OS",
				Type: "string",
			},
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "AMI ID",
				Type: "string",
			},
		},
	}

	for _, ami := range a.Items {

		row := metav1.TableRow{
			Cells: []interface{}{ami.Spec.KubernetesVersion, ami.Spec.Region, ami.Spec.OS, ami.GetName(), ami.Spec.ImageID},
		}
		table.Rows = append(table.Rows, row)

	}
	return table
}

func (a *AWSAMI) GetObjectKind() schema.ObjectKind {
	return &metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       AWSAMIKind,
	}
}

func (a *AWSAMIList) GetObjectKind() schema.ObjectKind {
	return &metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       AWSAMIListKind,
	}
}
