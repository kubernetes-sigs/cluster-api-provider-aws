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
	"unsafe"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	v1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this AWSCluster to the Hub version (v1beta1).
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.AWSCluster)

	if err := Convert_v1alpha2_AWSCluster_To_v1beta1_AWSCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually convert Status.APIEndpoints to Spec.ControlPlaneEndpoint.
	if len(src.Status.APIEndpoints) > 0 {
		endpoint := src.Status.APIEndpoints[0]
		dst.Spec.ControlPlaneEndpoint.Host = endpoint.Host
		dst.Spec.ControlPlaneEndpoint.Port = int32(endpoint.Port)
	}

	// Manually restore data.
	restored := &infrav1beta1.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	// override the SSHKeyName conversion if we are roundtripping from v1beta1 and the v1beta1 value is nil
	if src.Spec.SSHKeyName == "" && restored.Spec.SSHKeyName == nil {
		dst.Spec.SSHKeyName = nil
	}

	dst.Spec.Bastion.AllowedCIDRBlocks = restored.Spec.Bastion.AllowedCIDRBlocks
	dst.Spec.Bastion.AMI = restored.Spec.Bastion.AMI
	dst.Spec.Bastion.DisableIngressRules = restored.Spec.Bastion.DisableIngressRules
	dst.Spec.Bastion.InstanceType = restored.Spec.Bastion.InstanceType
	dst.Spec.ImageLookupFormat = restored.Spec.ImageLookupFormat
	dst.Spec.ImageLookupOrg = restored.Spec.ImageLookupOrg
	dst.Spec.ImageLookupBaseOS = restored.Spec.ImageLookupBaseOS

	// If src ControlPlaneLoadBalancer is nil, do not copy restored ControlPlaneLoadBalancer into it.
	if src.Spec.ControlPlaneLoadBalancer != nil {
		if restored.Spec.ControlPlaneLoadBalancer != nil {
			// If both restored and src ControlPlaneLoadBalancer is non-nil, only copy the missing part from restored.
			// Scheme is already copied from src in Convert_v1alpha2_AWSLoadBalancerSpec_To_v1beta1_AWSLoadBalancerSpec.
			dst.Spec.ControlPlaneLoadBalancer.CrossZoneLoadBalancing = restored.Spec.ControlPlaneLoadBalancer.CrossZoneLoadBalancing
			dst.Spec.ControlPlaneLoadBalancer.Subnets = restored.Spec.ControlPlaneLoadBalancer.Subnets
			dst.Spec.ControlPlaneLoadBalancer.AdditionalSecurityGroups = restored.Spec.ControlPlaneLoadBalancer.AdditionalSecurityGroups
			dst.Spec.ControlPlaneLoadBalancer.Name = restored.Spec.ControlPlaneLoadBalancer.Name
		}
	}

	dst.Spec.NetworkSpec.CNI = restored.Spec.NetworkSpec.CNI
	dst.Status.FailureDomains = restored.Status.FailureDomains
	dst.Status.Network.APIServerELB.AvailabilityZones = restored.Status.Network.APIServerELB.AvailabilityZones
	dst.Status.Network.APIServerELB.Attributes.CrossZoneLoadBalancing = restored.Status.Network.APIServerELB.Attributes.CrossZoneLoadBalancing
	dst.Spec.NetworkSpec.SecurityGroupOverrides = restored.Spec.NetworkSpec.SecurityGroupOverrides

	restoreInstance(restored.Status.Bastion, dst.Status.Bastion)

	// Manually set RootDeviceSize after restoring Bastion instance.
	if src.Status.Bastion.RootDeviceSize != 0 {
		if dst.Status.Bastion.RootVolume == nil {
			dst.Status.Bastion.RootVolume = &infrav1beta1.Volume{
				Size: src.Status.Bastion.RootDeviceSize,
			}
		} else {
			dst.Status.Bastion.RootVolume.Size = src.Status.Bastion.RootDeviceSize
		}
	}

	if restored.Spec.NetworkSpec.VPC.AvailabilityZoneUsageLimit != nil {
		dst.Spec.NetworkSpec.VPC.AvailabilityZoneUsageLimit = restored.Spec.NetworkSpec.VPC.AvailabilityZoneUsageLimit
	}
	if restored.Spec.NetworkSpec.VPC.AvailabilityZoneSelection != nil {
		dst.Spec.NetworkSpec.VPC.AvailabilityZoneSelection = restored.Spec.NetworkSpec.VPC.AvailabilityZoneSelection
	}
	// Manually convert conditions
	dst.SetConditions(restored.GetConditions())

	dst.Spec.IdentityRef = restored.Spec.IdentityRef
	return nil
}

func restoreInstance(restored, dst *infrav1beta1.Instance) {
	if restored != nil {
		dst.AvailabilityZone = restored.AvailabilityZone
		dst.NonRootVolumes = restored.NonRootVolumes
		dst.SpotMarketOptions = restored.SpotMarketOptions

		// Note this may override the manual conversion in Convert_v1alpha2_Instance_To_v1beta1_Instance.
		if restored.RootVolume != nil {
			restored.RootVolume.DeepCopyInto(dst.RootVolume)
		}

		dst.Tenancy = restored.Tenancy
		dst.VolumeIDs = restored.VolumeIDs
	}
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1beta1.AWSCluster)

	if err := Convert_v1beta1_AWSCluster_To_v1alpha2_AWSCluster(src, dst, nil); err != nil {
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

// ConvertTo converts this AWSClusterList to the Hub version (v1beta1).
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.AWSClusterList)
	return Convert_v1alpha2_AWSClusterList_To_v1beta1_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1beta1.AWSClusterList)
	return Convert_v1beta1_AWSClusterList_To_v1alpha2_AWSClusterList(src, dst, nil)
}

// Convert_v1alpha2_AWSClusterStatus_To_v1beta1_AWSClusterStatus converts AWSCluster.Status from v1alpha2 to v1beta1.
func Convert_v1alpha2_AWSClusterStatus_To_v1beta1_AWSClusterStatus(in *AWSClusterStatus, out *v1beta1.AWSClusterStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_AWSClusterStatus_To_v1beta1_AWSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert Status.Bastion.
	if !reflect.DeepEqual(in.Bastion, Instance{}) {
		out.Bastion = &v1beta1.Instance{}
		if err := Convert_v1alpha2_Instance_To_v1beta1_Instance(&in.Bastion, out.Bastion, s); err != nil {
			return err
		}
	}

	return nil
}

// Convert_v1alpha2_AWSClusterSpec_To_v1beta1_AWSClusterSpec.
func Convert_v1alpha2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in *AWSClusterSpec, out *infrav1beta1.AWSClusterSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_AWSClusterSpec_To_v1beta1_AWSClusterSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert Bastion.
	out.Bastion.Enabled = !in.DisableBastionHost

	// Manually convert SSHKeyName
	out.SSHKeyName = pointer.StringPtr(in.SSHKeyName)

	return nil
}

// Convert_v1beta1_AWSClusterSpec_To_v1alpha2_AWSClusterSpec converts from the Hub version (v1beta1) of the AWSClusterSpec to this version.
// Requires manual conversion as infrav1beta1.AWSClusterSpec.ImageLookupOrg does not exist in AWSClusterSpec.
func Convert_v1beta1_AWSClusterSpec_To_v1alpha2_AWSClusterSpec(in *infrav1beta1.AWSClusterSpec, out *AWSClusterSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_AWSClusterSpec_To_v1alpha2_AWSClusterSpec(in, out, s); err != nil {
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

// Convert_v1beta1_AWSClusterStatus_To_v1alpha2_AWSClusterStatus.
func Convert_v1beta1_AWSClusterStatus_To_v1alpha2_AWSClusterStatus(in *infrav1beta1.AWSClusterStatus, out *AWSClusterStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_AWSClusterStatus_To_v1alpha2_AWSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert Status.Bastion.
	if in.Bastion != nil {
		if err := Convert_v1beta1_Instance_To_v1alpha2_Instance(in.Bastion, &out.Bastion, s); err != nil {
			return err
		}
	}

	return nil
}

// Convert_v1beta1_ClassicELB_To_v1alpha2_ClassicELB.
func Convert_v1beta1_ClassicELB_To_v1alpha2_ClassicELB(in *infrav1beta1.ClassicELB, out *ClassicELB, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_ClassicELB_To_v1alpha2_ClassicELB(in, out, s); err != nil {
		return err
	}

	if in.Listeners != nil {
		out.Listeners = []*ClassicELBListener{}
		for _, inELB := range in.Listeners {
			out.Listeners = append(out.Listeners, &ClassicELBListener{
				Protocol:         ClassicELBProtocol(inELB.Protocol),
				Port:             inELB.Port,
				InstanceProtocol: ClassicELBProtocol(inELB.InstanceProtocol),
				InstancePort:     inELB.InstancePort,
			})
		}
	} else {
		out.Listeners = nil
	}

	return nil
}

func Convert_v1alpha2_ClassicELB_To_v1beta1_ClassicELB(in *ClassicELB, out *v1beta1.ClassicELB, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_ClassicELB_To_v1beta1_ClassicELB(in, out, s); err != nil {
		return err
	}

	if in.Listeners != nil {
		out.Listeners = []v1beta1.ClassicELBListener{}
		for _, inELB := range in.Listeners {
			if inELB != nil {
				out.Listeners = append(out.Listeners, v1beta1.ClassicELBListener{
					Protocol:         v1beta1.ClassicELBProtocol(inELB.Protocol),
					Port:             inELB.Port,
					InstanceProtocol: v1beta1.ClassicELBProtocol(inELB.InstanceProtocol),
					InstancePort:     inELB.InstancePort,
				})
			}
		}
	} else {
		out.Listeners = nil
	}

	out.AvailabilityZones = []string{}

	return nil
}

// Convert_v1beta1_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec.
func Convert_v1beta1_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec(in *infrav1beta1.AWSLoadBalancerSpec, out *AWSLoadBalancerSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSLoadBalancerSpec_To_v1alpha2_AWSLoadBalancerSpec(in, out, s)
}

func Convert_v1beta1_ClassicELBAttributes_To_v1alpha2_ClassicELBAttributes(in *infrav1beta1.ClassicELBAttributes, out *ClassicELBAttributes, s apiconversion.Scope) error {
	return autoConvert_v1beta1_ClassicELBAttributes_To_v1alpha2_ClassicELBAttributes(in, out, s)
}

// Convert_v1beta1_VPCSpec_To_v1alpha2_VPCSpec is an autogenerated conversion function.
func Convert_v1beta1_VPCSpec_To_v1alpha2_VPCSpec(in *infrav1beta1.VPCSpec, out *VPCSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_VPCSpec_To_v1alpha2_VPCSpec(in, out, s)
}

// Convert_v1beta1_NetworkSpec_To_v1alpha2_NetworkSpec
func Convert_v1beta1_NetworkSpec_To_v1alpha2_NetworkSpec(in *infrav1beta1.NetworkSpec, out *NetworkSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_NetworkSpec_To_v1alpha2_NetworkSpec(in, out, s); err != nil {
		return err
	}

	if in.Subnets != nil {
		out.Subnets = Subnets{}
		for _, ir := range in.Subnets {
			outSubnet := SubnetSpec{}
			if err := Convert_v1beta1_SubnetSpec_To_v1alpha2_SubnetSpec(&ir, &outSubnet, s); err != nil {
				return err
			}
			out.Subnets = append(out.Subnets, &outSubnet)
		}
	} else {
		out.Subnets = nil
	}

	return nil
}

func Convert_v1alpha2_NetworkSpec_To_v1beta1_NetworkSpec(in *NetworkSpec, out *v1beta1.NetworkSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_NetworkSpec_To_v1beta1_NetworkSpec(in, out, s); err != nil {
		return err
	}

	if in.Subnets != nil {
		out.Subnets = v1beta1.Subnets{}
		for _, ir := range in.Subnets {
			if ir != nil {
				outSubnet := v1beta1.SubnetSpec{}
				if err := Convert_v1alpha2_SubnetSpec_To_v1beta1_SubnetSpec(ir, &outSubnet, s); err != nil {
					return err
				}
				out.Subnets = append(out.Subnets, outSubnet)
			}
		}
	} else {
		out.Subnets = nil
	}

	out.SecurityGroupOverrides = map[v1beta1.SecurityGroupRole]string{}

	return nil
}

func Convert_v1alpha2_Network_To_v1beta1_NetworkStatus(in *Network, out *v1beta1.NetworkStatus, s apiconversion.Scope) error {
	if err := Convert_v1alpha2_ClassicELB_To_v1beta1_ClassicELB(&in.APIServerELB, &out.APIServerELB, s); err != nil {
		return nil
	}

	if in.SecurityGroups != nil {
		out.SecurityGroups = make(map[v1beta1.SecurityGroupRole]v1beta1.SecurityGroup, len(in.SecurityGroups))
		for role, group := range in.SecurityGroups {
			outGroup := v1beta1.SecurityGroup{}
			if err := Convert_v1alpha2_SecurityGroup_To_v1beta1_SecurityGroup(&group, &outGroup, s); err != nil {
				return err
			}
			out.SecurityGroups[v1beta1.SecurityGroupRole(role)] = outGroup
		}
	} else {
		out.SecurityGroups = nil
	}

	return nil
}

func Convert_v1beta1_NetworkStatus_To_v1alpha2_Network(in *v1beta1.NetworkStatus, out *Network, s apiconversion.Scope) error {
	if err := Convert_v1beta1_ClassicELB_To_v1alpha2_ClassicELB(&in.APIServerELB, &out.APIServerELB, s); err != nil {
		return nil
	}

	if in.SecurityGroups != nil {
		out.SecurityGroups = make(map[SecurityGroupRole]SecurityGroup, len(in.SecurityGroups))
		for role, group := range in.SecurityGroups {
			outGroup := SecurityGroup{}
			if err := Convert_v1beta1_SecurityGroup_To_v1alpha2_SecurityGroup(&group, &outGroup, s); err != nil {
				return err
			}
			out.SecurityGroups[SecurityGroupRole(role)] = outGroup
		}
	} else {
		out.SecurityGroups = nil
	}

	return nil
}

func Convert_v1beta1_SecurityGroup_To_v1alpha2_SecurityGroup(in *v1beta1.SecurityGroup, out *SecurityGroup, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_SecurityGroup_To_v1alpha2_SecurityGroup(in, out, s); err != nil {
		return err
	}

	if in.IngressRules != nil {
		out.IngressRules = IngressRules{}
		for _, ir := range in.IngressRules {
			out.IngressRules = append(out.IngressRules, (*IngressRule)(unsafe.Pointer(&ir)))
		}
	} else {
		out.IngressRules = nil
	}

	return nil
}

func Convert_v1alpha2_SecurityGroup_To_v1beta1_SecurityGroup(in *SecurityGroup, out *v1beta1.SecurityGroup, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_SecurityGroup_To_v1beta1_SecurityGroup(in, out, s); err != nil {
		return err
	}

	if in.IngressRules != nil {
		out.IngressRules = v1beta1.IngressRules{}
		for _, ir := range in.IngressRules {
			out.IngressRules = append(out.IngressRules, *(*v1beta1.IngressRule)(unsafe.Pointer(&ir)))
		}
	} else {
		out.IngressRules = nil
	}

	return nil
}
