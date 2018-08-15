/*
Copyright YEAR The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"
	"log"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/olympus"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Poseidon
// +k8s:openapi-gen=true
// +resource:path=poseidons,strategy=PoseidonStrategy
type Poseidon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PoseidonSpec   `json:"spec,omitempty"`
	Status PoseidonStatus `json:"status,omitempty"`
}

// PoseidonSpec defines the desired state of Poseidon
type PoseidonSpec struct {
	PodSpec    v1.PodTemplate
	Deployment v1beta1.Deployment
}

// PoseidonStatus defines the observed state of Poseidon
type PoseidonStatus struct {
}

// Validate checks that an instance of Poseidon is well formed
func (PoseidonStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*olympus.Poseidon)
	log.Printf("Validating fields for Poseidon %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default Poseidon field values
func (PoseidonSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*Poseidon)
	// set default field values here
	log.Printf("Defaulting fields for Poseidon %s\n", obj.Name)
}

func (b PoseidonStrategy) TriggerFunc(obj runtime.Object) []storage.MatchValue {
	// Change this function to override the trigger fn that is used
	value := b.DefaultStorageStrategy.TriggerFunc(obj)
	return value
}

// The following functions allow spec.deployment.name to be selected when listing
// or watching resources
func (b PoseidonStrategy) GetAttrs(o runtime.Object) (labels.Set, fields.Set, bool, error) {
	// Change this function to override the attributes that are matched
	l, _, uninit, e := b.DefaultStorageStrategy.GetAttrs(o)
	obj := o.(*olympus.Poseidon)

	fs := fields.Set{"spec.deployment.name": obj.Spec.Deployment.Name}
	fs = generic.AddObjectMetaFieldsSet(fs, &obj.ObjectMeta, true)
	return l, fs, uninit, e
}

func (b PoseidonStrategy) BasicMatch(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:       label,
		Field:       field,
		GetAttrs:    b.GetAttrs,
		IndexFields: []string{"spec.deployment.name"},
	}
}

// All field selector fields must appear in this function
func (b PoseidonSchemeFns) FieldSelectorConversion(label, value string) (string, string, error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	case "spec.deployment.name":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector: only %q, %q, %q", label, "metadata.name", "metadata.namespace", "spec.deployment.name")
	}
}
