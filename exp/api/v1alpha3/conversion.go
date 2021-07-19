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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	clusterapiapiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clusterapiapiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1alpha4 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
	infrav1alpha4exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
)

// ConvertTo converts the v1alpha3 AWSMachinePool receiver to a v1alpha4 AWSMachinePool.
func (r *AWSMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachinePool)

	if err := Convert_v1alpha3_AWSMachinePool_To_v1alpha4_AWSMachinePool(r, dst, nil); err != nil {
		return err
	}
	restored := &v1alpha4.AWSMachinePool{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	infrav1alpha3.RestoreAMIReference(&restored.Spec.AWSLaunchTemplate.AMI, &dst.Spec.AWSLaunchTemplate.AMI)
	return nil
}

// ConvertFrom converts the v1alpha4 AWSMachinePool receiver to a v1alpha3 AWSMachinePool.
func (r *AWSMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachinePool)

	if err := Convert_v1alpha4_AWSMachinePool_To_v1alpha3_AWSMachinePool(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha3 AWSMachinePoolList receiver to a v1alpha4 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachinePoolList)

	return Convert_v1alpha3_AWSMachinePoolList_To_v1alpha4_AWSMachinePoolList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSMachinePoolList receiver to a v1alpha3 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachinePoolList)

	return Convert_v1alpha4_AWSMachinePoolList_To_v1alpha3_AWSMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSManagedMachinePool receiver to a v1alpha4 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSManagedMachinePool)
	if err := Convert_v1alpha3_AWSManagedMachinePool_To_v1alpha4_AWSManagedMachinePool(r, dst, nil); err != nil {
		return err
	}

	restored := &infrav1alpha4exp.AWSManagedMachinePool{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Taints = restored.Spec.Taints

	return nil
}

// ConvertFrom converts the v1alpha4 AWSManagedMachinePool receiver to a v1alpha3 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSManagedMachinePool)

	if err := Convert_v1alpha4_AWSManagedMachinePool_To_v1alpha3_AWSManagedMachinePool(src, r, nil); err != nil {
		return err
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha3 AWSManagedMachinePoolList receiver to a v1alpha4 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSManagedMachinePoolList)

	return Convert_v1alpha3_AWSManagedMachinePoolList_To_v1alpha4_AWSManagedMachinePoolList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSManagedMachinePoolList receiver to a v1alpha3 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSManagedMachinePoolList)

	return Convert_v1alpha4_AWSManagedMachinePoolList_To_v1alpha3_AWSManagedMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSManagedCluster receiver to a v1alpha4 AWSManagedCluster.
func (r *AWSManagedCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSManagedCluster)

	return Convert_v1alpha3_AWSManagedCluster_To_v1alpha4_AWSManagedCluster(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSManagedCluster receiver to a v1alpha3 AWSManagedCluster.
func (r *AWSManagedCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSManagedCluster)

	return Convert_v1alpha4_AWSManagedCluster_To_v1alpha3_AWSManagedCluster(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSManagedClusterList receiver to a v1alpha4 AWSManagedClusterList.
func (r *AWSManagedClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSManagedClusterList)

	return Convert_v1alpha3_AWSManagedClusterList_To_v1alpha4_AWSManagedClusterList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSManagedClusterList receiver to a v1alpha3 AWSManagedClusterList.
func (r *AWSManagedClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSManagedClusterList)

	return Convert_v1alpha4_AWSManagedClusterList_To_v1alpha3_AWSManagedClusterList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSFargateProfile receiver to a v1alpha4 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSFargateProfile)

	return Convert_v1alpha3_AWSFargateProfile_To_v1alpha4_AWSFargateProfile(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSFargateProfile receiver to a v1alpha3 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSFargateProfile)

	return Convert_v1alpha4_AWSFargateProfile_To_v1alpha3_AWSFargateProfile(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSFargateProfileList receiver to a v1alpha4 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSFargateProfileList)

	return Convert_v1alpha3_AWSFargateProfileList_To_v1alpha4_AWSFargateProfileList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSFargateProfileList receiver to a v1alpha3 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSFargateProfileList)

	return Convert_v1alpha4_AWSFargateProfileList_To_v1alpha3_AWSFargateProfileList(src, r, nil)
}

// Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint is an autogenerated conversion function.
func Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint(in *clusterapiapiv1alpha3.APIEndpoint, out *clusterapiapiv1alpha4.APIEndpoint, s apiconversion.Scope) error {
	return clusterapiapiv1alpha3.Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint(in, out, s)
}

// Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint is an autogenerated conversion function.
func Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint(in *clusterapiapiv1alpha4.APIEndpoint, out *clusterapiapiv1alpha3.APIEndpoint, s apiconversion.Scope) error {
	return clusterapiapiv1alpha3.Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint(in, out, s)
}

// Convert_v1alpha3_AWSResourceReference_To_v1alpha4_AWSResourceReference is an autogenerated conversion function.
func Convert_v1alpha3_AWSResourceReference_To_v1alpha4_AWSResourceReference(in *infrav1alpha3.AWSResourceReference, out *infrav1alpha4.AWSResourceReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_AWSResourceReference_To_v1alpha4_AWSResourceReference(in, out, s)
}

// Convert_v1alpha4_AWSResourceReference_To_v1alpha3_AWSResourceReference conversion function.
func Convert_v1alpha4_AWSResourceReference_To_v1alpha3_AWSResourceReference(in *infrav1alpha4.AWSResourceReference, out *infrav1alpha3.AWSResourceReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha4_AWSResourceReference_To_v1alpha3_AWSResourceReference(in, out, s)
}

// Convert_v1alpha4_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec is an autogenerated conversion function.
func Convert_v1alpha4_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec(in *infrav1alpha4exp.AWSManagedMachinePoolSpec, out *AWSManagedMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha4_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec(in, out, s)
}

// Convert_v1alpha4_Instance_To_v1alpha3_Instance is an autogenerated conversion function.
func Convert_v1alpha4_Instance_To_v1alpha3_Instance(in *infrav1alpha4.Instance, out *infrav1alpha3.Instance, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha4_Instance_To_v1alpha3_Instance(in, out, s)
}

// Convert_v1alpha3_Instance_To_v1alpha4_Instance is an autogenerated conversion function.
func Convert_v1alpha3_Instance_To_v1alpha4_Instance(in *infrav1alpha3.Instance, out *infrav1alpha4.Instance, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Instance_To_v1alpha4_Instance(in, out, s)
}

func Convert_v1alpha3_AWSResourceReference_To_v1alpha4_AMIReference(in *infrav1alpha3.AWSResourceReference, out *infrav1alpha4.AMIReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_AWSResourceReference_To_v1alpha4_AMIReference(in, out, s)
}

func Convert_v1alpha4_AMIReference_To_v1alpha3_AWSResourceReference(in *infrav1alpha4.AMIReference, out *infrav1alpha3.AWSResourceReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha4_AMIReference_To_v1alpha3_AWSResourceReference(in, out, s)
}
