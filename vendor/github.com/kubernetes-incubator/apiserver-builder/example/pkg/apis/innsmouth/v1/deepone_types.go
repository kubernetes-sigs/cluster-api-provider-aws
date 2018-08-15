/*
Copyright 2017 The Kubernetes Authors.

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

package v1

import (
	"github.com/kubernetes-incubator/apiserver-builder/example/pkg/apis/innsmouth/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=deepones
// DeepOne defines a resident of innsmouth
type DeepOne struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeepOneSpec   `json:"spec,omitempty"`
	Status DeepOneStatus `json:"status,omitempty"`
}

// DeepOnesSpec defines the desired state of DeepOne
type DeepOneSpec struct {
	// fish_required defines the number of fish required by the DeepOne.
	FishRequired int `json:"fish_required,omitempty"`

	Sample            SampleElem                       `json:"sample,omitempty"`
	SamplePointer     *SamplePointerElem               `json:"sample_pointer,omitempty"`
	SampleList        []SampleListElem                 `json:"sample_list,omitempty"`
	SamplePointerList []*SampleListPointerElem         `json:"sample_pointer_list,omitempty"`
	SampleMap         map[string]SampleMapElem         `json:"sample_map,omitempty"`
	SamplePointerMap  map[string]*SampleMapPointerElem `json:"sample_pointer_map,omitempty"`

	// Example of using a constant
	Const      common.CustomType            `json:"const,omitempty"`
	ConstPtr   *common.CustomType           `json:"constPtr,omitempty"`
	ConstSlice []common.CustomType          `json:"constSlice,omitempty"`
	ConstMap   map[string]common.CustomType `json:"constMap,omitempty"`

	// TODO: Fix issues with deep copy to make these work
	//ConstSlicePtr []*common.CustomType          `json:"constSlicePtr,omitempty"`
	//ConstMapPtr map[string]*common.CustomType `json:"constMapPtr,omitempty"`
}

type SampleListElem struct {
	Sub []SampleListSubElem `json:"sub,omitempty"`
}

type SampleListSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleListPointerElem struct {
	Sub []*SampleListPointerSubElem `json:"sub,omitempty"`
}

type SampleListPointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleMapElem struct {
	Sub map[string]SampleMapSubElem `json:"sub,omitempty"`
}

type SampleMapSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleMapPointerElem struct {
	Sub map[string]*SampleMapPointerSubElem `json:"sub,omitempty"`
}

type SampleMapPointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SamplePointerElem struct {
	Sub *SamplePointerSubElem `json:"sub,omitempty"`
}

type SamplePointerSubElem struct {
	Foo string `json:"foo,omitempty"`
}

type SampleElem struct {
	Sub SampleSubElem `json:"sub,omitempty"`
}

type SampleSubElem struct {
	Foo string `json:"foo,omitempty"`
}

// DeepOneStatus defines the observed state of DeepOne
type DeepOneStatus struct {
	// actual_fish defines the number of fish caught by the DeepOne.
	ActualFish int `json:"actual_fish,omitempty"`
}
