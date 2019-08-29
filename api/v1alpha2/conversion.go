/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha2

import (
	"errors"
	unsafe "unsafe"

	conversion "k8s.io/apimachinery/pkg/conversion"
	v1alpha1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

// nolint
func Convert_v1alpha1_AWSClusterProviderSpec_To_v1alpha2_AWSClusterSpec(in *v1alpha1.AWSClusterProviderSpec, out *AWSClusterSpec, s conversion.Scope) error {
	if err := Convert_v1alpha1_NetworkSpec_To_v1alpha2_NetworkSpec(&in.NetworkSpec, &out.NetworkSpec, s); err != nil {
		return err
	}

	out.Region = in.Region
	out.SSHKeyName = in.SSHKeyName

	// DISCARDS:
	// CAKeyPair
	// EtcdCAKeyPair
	// FrontProxyCAKeyPair
	// SAKeyPair
	// ClusterConfiguration
	// AdditionalUserDataFiles

	return nil
}

// nolint
func Convert_v1alpha1_AWSClusterProviderStatus_To_v1alpha2_AWSClusterStatus(in *v1alpha1.AWSClusterProviderStatus, out *AWSClusterStatus, s conversion.Scope) error {
	if err := Convert_v1alpha1_Network_To_v1alpha2_Network(&in.Network, &out.Network, s); err != nil {
		return err
	}

	if err := Convert_v1alpha1_Instance_To_v1alpha2_Instance(&in.Bastion, &out.Bastion, s); err != nil {
		return err
	}

	return nil
}

// nolint
func Convert_v1alpha1_AWSMachineProviderSpec_To_v1alpha2_AWSMachineSpec(in *v1alpha1.AWSMachineProviderSpec, out *AWSMachineSpec, s conversion.Scope) error {

	if err := Convert_v1alpha1_AWSResourceReference_To_v1alpha2_AWSResourceReference(&in.AMI, &out.AMI, s); err != nil {
		return err
	}

	out.ImageLookupOrg = in.ImageLookupOrg
	out.InstanceType = in.InstanceType
	out.AdditionalTags = *(*map[string]string)(unsafe.Pointer(&in.AdditionalTags))
	out.IAMInstanceProfile = in.IAMInstanceProfile
	out.PublicIP = (*bool)(unsafe.Pointer(in.PublicIP))

	out.AdditionalSecurityGroups = make([]AWSResourceReference, len(in.AdditionalSecurityGroups))

	for i, sg := range in.AdditionalSecurityGroups {
		if err := Convert_v1alpha1_AWSResourceReference_To_v1alpha2_AWSResourceReference(&sg, &out.AdditionalSecurityGroups[i], s); err != nil {
			return err
		}
	}

	out.AvailabilityZone = (*string)(unsafe.Pointer(in.AvailabilityZone))

	var subnet AWSResourceReference

	if err := Convert_v1alpha1_AWSResourceReference_To_v1alpha2_AWSResourceReference(in.Subnet, &subnet, s); err != nil {
		return err
	}

	out.Subnet = &subnet

	out.SSHKeyName = in.KeyName
	out.RootDeviceSize = in.RootDeviceSize

	return nil
}

// nolint
func Convert_v1alpha2_Instance_To_v1alpha1_Instance(in *Instance, out *v1alpha1.Instance, s conversion.Scope) error {
	return errors.New("not implemented")
}

// nolint
func Convert_v1alpha1_Instance_To_v1alpha2_Instance(in *v1alpha1.Instance, out *Instance, s conversion.Scope) error {
	if err := autoConvert_v1alpha1_Instance_To_v1alpha2_Instance(in, out, s); err != nil {
		return err
	}

	out.SSHKeyName = in.KeyName
	return nil
}
