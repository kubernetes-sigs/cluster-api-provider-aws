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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func Convert_v1beta2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in *v1beta2.AWSClusterSpec, out *AWSClusterSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in, out, s)
}

func Convert_v1beta1_AWSResourceReference_To_v1beta2_AWSResourceReference(in *AWSResourceReference, out *v1beta2.AWSResourceReference, s conversion.Scope) error {
	return autoConvert_v1beta1_AWSResourceReference_To_v1beta2_AWSResourceReference(in, out, s)
}

func Convert_v1beta1_AWSMachineSpec_To_v1beta2_AWSMachineSpec(in *AWSMachineSpec, out *v1beta2.AWSMachineSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_AWSMachineSpec_To_v1beta2_AWSMachineSpec(in, out, s)
}
