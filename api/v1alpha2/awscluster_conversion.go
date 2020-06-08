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
	"reflect"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	v1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this AWSCluster to the Hub version (v1alpha3).
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1alpha3.AWSCluster)

	if err := Convert_v1alpha2_AWSCluster_To_v1alpha3_AWSCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually convert Status.APIEndpoints to Spec.ControlPlaneEndpoint.
	if len(src.Status.APIEndpoints) > 0 {
		endpoint := src.Status.APIEndpoints[0]
		dst.Spec.ControlPlaneEndpoint.Host = endpoint.Host
		dst.Spec.ControlPlaneEndpoint.Port = int32(endpoint.Port)
	}

	// Manually restore data.
	restored := &infrav1alpha3.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	// override the SSHKeyName conversion if we are roundtripping from v1alpha3 and the v1alpha3 value is nil
	if src.Spec.SSHKeyName == "" && restored.Spec.SSHKeyName == nil {
		dst.Spec.SSHKeyName = nil
	}

	dst.Spec.ImageLookupFormat = restored.Spec.ImageLookupFormat
	dst.Spec.ImageLookupOrg = restored.Spec.ImageLookupOrg
	dst.Spec.ImageLookupBaseOS = restored.Spec.ImageLookupBaseOS
	if restored.Spec.ControlPlaneLoadBalancer != nil {
		dst.Spec.ControlPlaneLoadBalancer = restored.Spec.ControlPlaneLoadBalancer
	}
	dst.Status.FailureDomains = restored.Status.FailureDomains
	dst.Status.Network.APIServerELB.AvailabilityZones = restored.Status.Network.APIServerELB.AvailabilityZones
	dst.Status.Network.APIServerELB.Attributes.CrossZoneLoadBalancing = restored.Status.Network.APIServerELB.Attributes.CrossZoneLoadBalancing

	if restored.Status.Bastion != nil {
		restored.Status.Bastion.DeepCopyInto(dst.Status.Bastion)
	}

	// Manually convert conditions
	dst.SetConditions(restored.GetConditions())

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1alpha3.AWSCluster)

	if err := Convert_v1alpha3_AWSCluster_To_v1alpha2_AWSCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually convert Spec.ControlPlaneEndpoint to Status.APIEndpoints.
	if !src.Spec.ControlPlaneEndpoint.IsZero() {
		dst.Status.APIEndpoints = []APIEndpoint{
			{
				Host: src.Spec.ControlPlaneEndpoint.Host,
				Port: int(src.Spec.ControlPlaneEndpoint.Port),
			},
		}
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this AWSClusterList to the Hub version (v1alpha3).
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1alpha3.AWSClusterList)
	return Convert_v1alpha2_AWSClusterList_To_v1alpha3_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1alpha3.AWSClusterList)
	return Convert_v1alpha3_AWSClusterList_To_v1alpha2_AWSClusterList(src, dst, nil)
}

// Convert_v1alpha2_AWSClusterStatus_To_v1alpha3_AWSClusterStatus converts AWSCluster.Status from v1alpha2 to v1alpha3.
func Convert_v1alpha2_AWSClusterStatus_To_v1alpha3_AWSClusterStatus(in *AWSClusterStatus, out *v1alpha3.AWSClusterStatus, s apiconversion.Scope) error { // nolint
	if err := autoConvert_v1alpha2_AWSClusterStatus_To_v1alpha3_AWSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert Status.Bastion.
	if !reflect.DeepEqual(in.Bastion, Instance{}) {
		out.Bastion = &v1alpha3.Instance{}
		if err := Convert_v1alpha2_Instance_To_v1alpha3_Instance(&in.Bastion, out.Bastion, s); err != nil {
			return err
		}
	}

	return nil
}

// Convert_v1alpha2_AWSClusterSpec_To_v1alpha3_AWSClusterSpec.
func Convert_v1alpha2_AWSClusterSpec_To_v1alpha3_AWSClusterSpec(in *AWSClusterSpec, out *infrav1alpha3.AWSClusterSpec, s apiconversion.Scope) error { //nolint
	if err := autoConvert_v1alpha2_AWSClusterSpec_To_v1alpha3_AWSClusterSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert Bastion.
	out.Bastion.Enabled = !in.DisableBastionHost

	// Manually convert SSHKeyName
	out.SSHKeyName = pointer.StringPtr(in.SSHKeyName)

	return nil
}

// Convert_v1alpha3_AWSClusterSpec_To_v1alpha2_AWSClusterSpec converts from the Hub version (v1alpha3) of the AWSClusterSpec to this version.
// Requires manual conversion as infrav1alpha3.AWSClusterSpec.ImageLookupOrg does not exist in AWSClusterSpec.
func Convert_v1alpha3_AWSClusterSpec_To_v1alpha2_AWSClusterSpec(in *infrav1alpha3.AWSClusterSpec, out *AWSClusterSpec, s apiconversion.Scope) error { // nolint
	if err := autoConvert_v1alpha3_AWSClusterSpec_To_v1alpha2_AWSClusterSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert DisableBastionHost.
	out.DisableBastionHost = !in.Bastion.Enabled

	// Manually convert SSHKeyName
	if in.SSHKeyName != nil {
		out.SSHKeyName = *in.SSHKeyName
	}

	return nil
}

// Convert_v1alpha3_AWSClusterStatus_To_v1alpha2_AWSClusterStatus.
func Convert_v1alpha3_AWSClusterStatus_To_v1alpha2_AWSClusterStatus(in *infrav1alpha3.AWSClusterStatus, out *AWSClusterStatus, s apiconversion.Scope) error { //nolint
	if err := autoConvert_v1alpha3_AWSClusterStatus_To_v1alpha2_AWSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert Status.Bastion.
	if in.Bastion != nil {
		if err := Convert_v1alpha3_Instance_To_v1alpha2_Instance(in.Bastion, &out.Bastion, s); err != nil {
			return err
		}
	}

	return nil
}

// Convert_v1alpha3_ClassicELB_To_v1alpha2_ClassicELB.
func Convert_v1alpha3_ClassicELB_To_v1alpha2_ClassicELB(in *infrav1alpha3.ClassicELB, out *ClassicELB, s apiconversion.Scope) error { //nolint
	return autoConvert_v1alpha3_ClassicELB_To_v1alpha2_ClassicELB(in, out, s)
}

// Convert_v1alpha3_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec.
func Convert_v1alpha3_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec(in *infrav1alpha3.AWSLoadBalancerSpec, out *AWSLoadBalancerSpec, s apiconversion.Scope) error { //nolint
	return autoConvert_v1alpha3_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec(in, out, s)
}

func Convert_v1alpha3_ClassicELBAttributes_To_v1alpha2_ClassicELBAttributes(in *infrav1alpha3.ClassicELBAttributes, out *ClassicELBAttributes, s apiconversion.Scope) error { //nolint
	return autoConvert_v1alpha3_ClassicELBAttributes_To_v1alpha2_ClassicELBAttributes(in, out, s)
}
