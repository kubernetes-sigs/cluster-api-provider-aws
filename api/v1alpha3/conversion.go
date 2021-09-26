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

package v1alpha3

import (
	"unsafe"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	apiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	apiUpstreamv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 AWSCluster receiver to a v1beta1 AWSCluster.
func (r *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSCluster)

	if err := Convert_v1alpha3_AWSCluster_To_v1beta1_AWSCluster(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &v1beta1.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	if restored.Status.Bastion != nil {
		if dst.Status.Bastion == nil {
			dst.Status.Bastion = &v1beta1.Instance{}
		}
		restoreInstance(restored.Status.Bastion, dst.Status.Bastion)
	}

	return nil
}

// ConvertFrom converts the v1beta1 AWSCluster receiver to a v1alpha3 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSCluster)

	if err := Convert_v1beta1_AWSCluster_To_v1alpha3_AWSCluster(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha3 AWSClusterList receiver to a v1beta1 AWSClusterList.
func (r *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterList)

	return Convert_v1alpha3_AWSClusterList_To_v1beta1_AWSClusterList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterList receiver to a v1alpha3 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterList)

	return Convert_v1beta1_AWSClusterList_To_v1alpha3_AWSClusterList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachine receiver to a v1beta1 AWSMachine.
func (r *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachine)
	if err := Convert_v1alpha3_AWSMachine_To_v1beta1_AWSMachine(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &v1beta1.AWSMachine{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	RestoreAMIReference(&restored.Spec.AMI, &dst.Spec.AMI)
	if restored.Spec.RootVolume != nil {
		if dst.Spec.RootVolume == nil {
			dst.Spec.RootVolume = &v1beta1.Volume{}
		}
		RestoreRootVolume(restored.Spec.RootVolume, dst.Spec.RootVolume)
	}
	if restored.Spec.NonRootVolumes != nil {
		if dst.Spec.NonRootVolumes == nil {
			dst.Spec.NonRootVolumes = []v1beta1.Volume{}
		}
		restoreNonRootVolumes(restored.Spec.NonRootVolumes, dst.Spec.NonRootVolumes)
	}
	return nil
}

// ConvertFrom converts the v1beta1 AWSMachine receiver to a v1alpha3 AWSMachine.
func (r *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSMachine)

	if err := Convert_v1beta1_AWSMachine_To_v1alpha3_AWSMachine(src, r, nil); err != nil {
		return err
	}
	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts the v1alpha3 AWSMachineList receiver to a v1beta1 AWSMachineList.
func (r *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineList)

	return Convert_v1alpha3_AWSMachineList_To_v1beta1_AWSMachineList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineList receiver to a v1alpha3 AWSMachineList.
func (r *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSMachineList)

	return Convert_v1beta1_AWSMachineList_To_v1alpha3_AWSMachineList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachineTemplate receiver to a v1beta1 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineTemplate)
	if err := Convert_v1alpha3_AWSMachineTemplate_To_v1beta1_AWSMachineTemplate(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &v1beta1.AWSMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	RestoreAMIReference(&restored.Spec.Template.Spec.AMI, &dst.Spec.Template.Spec.AMI)
	if restored.Spec.Template.Spec.RootVolume != nil {
		if dst.Spec.Template.Spec.RootVolume == nil {
			dst.Spec.Template.Spec.RootVolume = &v1beta1.Volume{}
		}
		RestoreRootVolume(restored.Spec.Template.Spec.RootVolume, dst.Spec.Template.Spec.RootVolume)
	}
	if restored.Spec.Template.Spec.NonRootVolumes != nil {
		if dst.Spec.Template.Spec.NonRootVolumes == nil {
			dst.Spec.Template.Spec.NonRootVolumes = []v1beta1.Volume{}
		}
		restoreNonRootVolumes(restored.Spec.Template.Spec.NonRootVolumes, dst.Spec.Template.Spec.NonRootVolumes)
	}

	return nil
}

// ConvertFrom converts the v1beta1 AWSMachineTemplate receiver to a v1alpha3 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSMachineTemplate)

	if err := Convert_v1beta1_AWSMachineTemplate_To_v1alpha3_AWSMachineTemplate(src, r, nil); err != nil {
		return err
	}
	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts the v1alpha3 AWSMachineTemplateList receiver to a v1beta1 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineTemplateList)

	return Convert_v1alpha3_AWSMachineTemplateList_To_v1beta1_AWSMachineTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineTemplateList receiver to a v1alpha3 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSMachineTemplateList)

	return Convert_v1beta1_AWSMachineTemplateList_To_v1alpha3_AWSMachineTemplateList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentity receiver to a v1beta1 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterStaticIdentity)
	if err := Convert_v1alpha3_AWSClusterStaticIdentity_To_v1beta1_AWSClusterStaticIdentity(r, dst, nil); err != nil {
		return err
	}

	dst.Spec.SecretRef = r.Spec.SecretRef.Name
	return nil
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentity receiver to a v1alpha3 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterStaticIdentity)

	if err := Convert_v1beta1_AWSClusterStaticIdentity_To_v1alpha3_AWSClusterStaticIdentity(src, r, nil); err != nil {
		return err
	}

	r.Spec.SecretRef.Name = src.Spec.SecretRef
	return nil
}

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentityList receiver to a v1beta1 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterStaticIdentityList)

	return Convert_v1alpha3_AWSClusterStaticIdentityList_To_v1beta1_AWSClusterStaticIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentityList receiver to a v1alpha3 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterStaticIdentityList)

	return Convert_v1beta1_AWSClusterStaticIdentityList_To_v1alpha3_AWSClusterStaticIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentity receiver to a v1beta1 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterRoleIdentity)

	return Convert_v1alpha3_AWSClusterRoleIdentity_To_v1beta1_AWSClusterRoleIdentity(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentity receiver to a v1alpha3 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterRoleIdentity)

	return Convert_v1beta1_AWSClusterRoleIdentity_To_v1alpha3_AWSClusterRoleIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentityList receiver to a v1beta1 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterRoleIdentityList)

	return Convert_v1alpha3_AWSClusterRoleIdentityList_To_v1beta1_AWSClusterRoleIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentityList receiver to a v1alpha3 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterRoleIdentityList)

	return Convert_v1beta1_AWSClusterRoleIdentityList_To_v1alpha3_AWSClusterRoleIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentity receiver to a v1beta1 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterControllerIdentity)

	return Convert_v1alpha3_AWSClusterControllerIdentity_To_v1beta1_AWSClusterControllerIdentity(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentity receiver to a v1alpha3 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterControllerIdentity)

	return Convert_v1beta1_AWSClusterControllerIdentity_To_v1alpha3_AWSClusterControllerIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentityList receiver to a v1beta1 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterControllerIdentityList)

	return Convert_v1alpha3_AWSClusterControllerIdentityList_To_v1beta1_AWSClusterControllerIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentityList receiver to a v1alpha3 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterControllerIdentityList)

	return Convert_v1beta1_AWSClusterControllerIdentityList_To_v1alpha3_AWSClusterControllerIdentityList(src, r, nil)
}

// Convert_v1beta1_Volume_To_v1alpha3_Volume .
func Convert_v1beta1_Volume_To_v1alpha3_Volume(in *v1beta1.Volume, out *Volume, s apiconversion.Scope) error {
	return autoConvert_v1beta1_Volume_To_v1alpha3_Volume(in, out, s)
}

// Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint .
func Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint(in *apiv1alpha3.APIEndpoint, out *apiUpstreamv1beta1.APIEndpoint, s apiconversion.Scope) error {
	return apiv1alpha3.Convert_v1alpha3_APIEndpoint_To_v1beta1_APIEndpoint(in, out, s)
}

// Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint .
func Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint(in *apiUpstreamv1beta1.APIEndpoint, out *apiv1alpha3.APIEndpoint, s apiconversion.Scope) error {
	return apiv1alpha3.Convert_v1beta1_APIEndpoint_To_v1alpha3_APIEndpoint(in, out, s)
}

// Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec .
func Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec(in *AWSClusterStaticIdentitySpec, out *v1beta1.AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1beta1_AWSClusterStaticIdentitySpec(in, out, s)
}

// Convert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec .
func Convert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in *v1beta1.AWSClusterStaticIdentitySpec, out *AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in, out, s)
}

// Convert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec .
func Convert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in *v1beta1.AWSMachineSpec, out *AWSMachineSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in, out, s)
}

// Convert_v1beta1_Instance_To_v1alpha3_Instance .
func Convert_v1beta1_Instance_To_v1alpha3_Instance(in *v1beta1.Instance, out *Instance, s apiconversion.Scope) error {
	return autoConvert_v1beta1_Instance_To_v1alpha3_Instance(in, out, s)
}

// Convert_v1alpha3_Network_To_v1beta1_NetworkStatus is based on the autogenerated function and handles the renaming of the Network struct to NetworkStatus
func Convert_v1alpha3_Network_To_v1beta1_NetworkStatus(in *Network, out *v1beta1.NetworkStatus, s apiconversion.Scope) error {
	out.SecurityGroups = *(*map[v1beta1.SecurityGroupRole]v1beta1.SecurityGroup)(unsafe.Pointer(&in.SecurityGroups))
	if err := Convert_v1alpha3_ClassicELB_To_v1beta1_ClassicELB(&in.APIServerELB, &out.APIServerELB, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_NetworkStatus_To_v1alpha3_Network is based on the autogenerated function and handles the renaming of the NetworkStatus struct to Network
func Convert_v1beta1_NetworkStatus_To_v1alpha3_Network(in *v1beta1.NetworkStatus, out *Network, s apiconversion.Scope) error {
	out.SecurityGroups = *(*map[SecurityGroupRole]SecurityGroup)(unsafe.Pointer(&in.SecurityGroups))
	if err := Convert_v1beta1_ClassicELB_To_v1alpha3_ClassicELB(&in.APIServerELB, &out.APIServerELB, s); err != nil {
		return err
	}
	return nil
}

// Manually restore the instance root device data.
// Assumes restored and dst are non-nil.
func restoreInstance(restored, dst *v1beta1.Instance) {
	dst.VolumeIDs = restored.VolumeIDs

	if restored.RootVolume != nil {
		if dst.RootVolume == nil {
			dst.RootVolume = &v1beta1.Volume{}
		}
		RestoreRootVolume(restored.RootVolume, dst.RootVolume)
	}

	if restored.NonRootVolumes != nil {
		if dst.NonRootVolumes == nil {
			dst.NonRootVolumes = []v1beta1.Volume{}
		}
		restoreNonRootVolumes(restored.NonRootVolumes, dst.NonRootVolumes)
	}
}

// Convert_v1alpha3_AWSResourceReference_To_v1beta1_AMIReference is a conversion function.
func Convert_v1alpha3_AWSResourceReference_To_v1beta1_AMIReference(in *AWSResourceReference, out *v1beta1.AMIReference, s apiconversion.Scope) error {
	out.ID = (*string)(unsafe.Pointer(in.ID))
	return nil
}

// Convert_v1beta1_AMIReference_To_v1alpha3_AWSResourceReference is a conversion function.
func Convert_v1beta1_AMIReference_To_v1alpha3_AWSResourceReference(in *v1beta1.AMIReference, out *AWSResourceReference, s apiconversion.Scope) error {
	out.ID = (*string)(unsafe.Pointer(in.ID))
	return nil
}

// RestoreAMIReference manually restore the EKSOptimizedLookupType for AWSMachine and AWSMachineTemplate
// Assumes both restored and dst are non-nil.
func RestoreAMIReference(restored, dst *v1beta1.AMIReference) {
	if restored.EKSOptimizedLookupType == nil {
		return
	}
	dst.EKSOptimizedLookupType = restored.EKSOptimizedLookupType
}

// restoreNonRootVolumes manually restores the non-root volumes
// Assumes both restoredVolumes and dstVolumes are non-nil.
func restoreNonRootVolumes(restoredVolumes, dstVolumes []v1beta1.Volume) {
	// restoring the nonrootvolumes which are missing in dstVolumes
	// restoring dstVolumes[i].Encrypted to nil in order to avoid v1beta1 --> v1alpha3 --> v1beta1 round trip errors
	for i := range restoredVolumes {
		if restoredVolumes[i].Encrypted == nil {
			if len(dstVolumes) <= i {
				dstVolumes = append(dstVolumes, restoredVolumes[i])
			} else {
				dstVolumes[i].Encrypted = nil
			}
		}
		dstVolumes[i].Throughput = restoredVolumes[i].Throughput
	}
}

// RestoreRootVolume manually restores the root volumes.
// Assumes both restored and dst are non-nil.
// Volume.Encrypted type changed from bool in v1alpha3 to *bool in v1beta1
// Volume.Encrypted value as nil/&false in v1beta1 will convert to false in v1alpha3 by auto-conversion, so restoring it to nil in order to avoid v1beta1 --> v1alpha3 --> v1beta1 round trip errors
func RestoreRootVolume(restored, dst *v1beta1.Volume) {
	if dst.Encrypted == pointer.BoolPtr(true) {
		return
	}
	if restored.Encrypted == nil {
		dst.Encrypted = nil
	}
	dst.Throughput = restored.Throughput
}
