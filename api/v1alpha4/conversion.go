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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Convert_v1alpha4_APIEndpoint_To_v1beta1_APIEndpoint .
func Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in *clusterv1alpha4.ObjectMeta, out *clusterv1.ObjectMeta, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in, out, s)
}

// Convert_v1beta1_APIEndpoint_To_v1alpha4_APIEndpoint .
func Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in *clusterv1.ObjectMeta, out *clusterv1alpha4.ObjectMeta, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in, out, s)
}
