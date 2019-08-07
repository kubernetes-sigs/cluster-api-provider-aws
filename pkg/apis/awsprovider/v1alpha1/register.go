/*
Copyright 2018 The Kubernetes Authors.

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

// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the awsprovider v1alpha1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider
// +k8s:defaulter-gen=TypeMeta
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
	"sigs.k8s.io/yaml"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "awsprovider.k8s.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// ClusterConfigFromProviderSpec unmarshals a provider config into an AWS Cluster type
func ClusterConfigFromProviderSpec(providerConfig clusterv1.ProviderSpec) (*AWSClusterProviderSpec, error) {
	var config AWSClusterProviderSpec
	if providerConfig.Value == nil {
		return &config, nil
	}

	if err := yaml.Unmarshal(providerConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ClusterStatusFromProviderStatus unmarshals a raw extension into an AWS Cluster type
func ClusterStatusFromProviderStatus(extension *runtime.RawExtension) (*AWSClusterProviderStatus, error) {
	if extension == nil {
		return &AWSClusterProviderStatus{}, nil
	}

	status := new(AWSClusterProviderStatus)
	if err := yaml.Unmarshal(extension.Raw, status); err != nil {
		return nil, err
	}

	return status, nil
}

// MachineStatusFromProviderStatus unmarshals a raw extension into an AWS machine type
func MachineStatusFromProviderStatus(extension *runtime.RawExtension) (*AWSMachineProviderStatus, error) {
	if extension == nil {
		return &AWSMachineProviderStatus{}, nil
	}

	status := new(AWSMachineProviderStatus)
	if err := yaml.Unmarshal(extension.Raw, status); err != nil {
		return nil, err
	}

	return status, nil
}

// EncodeMachineStatus marshals the machine status
func EncodeMachineStatus(status *AWSMachineProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeMachineSpec marshals the machine provider spec.
func EncodeMachineSpec(spec *AWSMachineProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeClusterStatus marshals the cluster status.
func EncodeClusterStatus(status *AWSClusterProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

// EncodeClusterSpec marshals the cluster provider spec.
func EncodeClusterSpec(spec *AWSClusterProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}
