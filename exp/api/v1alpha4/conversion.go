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
	infrav1alpha4 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha4 AWSMachinePool receiver to a v1beta1 AWSMachinePool.
func ( src *AWSMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachinePool)
	return Convert_v1alpha4_AWSMachinePool_To_v1beta1_AWSMachinePool(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachinePool receiver to v1alpha4 AWSMachinePool.
func (r *AWSMachinePool) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachinePool)

	return Convert_v1beta1_AWSMachinePool_To_v1alpha4_AWSMachinePool(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSMachinePoolList receiver to a v1beta1 AWSMachinePoolList.
func ( src *AWSMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachinePoolList)
	return Convert_v1alpha4_AWSMachinePoolList_To_v1beta1_AWSMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachinePoolList receiver to v1alpha4 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachinePoolList)

	return Convert_v1beta1_AWSMachinePoolList_To_v1alpha4_AWSMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSManagedMachinePool receiver to a v1beta1 AWSManagedMachinePool.
func ( src *AWSManagedMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedMachinePool)
	return Convert_v1alpha4_AWSManagedMachinePool_To_v1beta1_AWSManagedMachinePool(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedMachinePool receiver to v1alpha4 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSManagedMachinePool)

	return Convert_v1beta1_AWSManagedMachinePool_To_v1alpha4_AWSManagedMachinePool(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSManagedMachinePoolList receiver to a v1beta1 AWSManagedMachinePoolList.
func ( src *AWSManagedMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSManagedMachinePoolList)
	return Convert_v1alpha4_AWSManagedMachinePoolList_To_v1beta1_AWSManagedMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedMachinePoolList receiver to v1alpha4 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSManagedMachinePoolList)

	return Convert_v1beta1_AWSManagedMachinePoolList_To_v1alpha4_AWSManagedMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSFargateProfile receiver to a v1beta1 AWSFargateProfile.
func ( src *AWSFargateProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSFargateProfile)
	return Convert_v1alpha4_AWSFargateProfile_To_v1beta1_AWSFargateProfile(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSFargateProfile receiver to v1alpha4 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSFargateProfile)

	return Convert_v1beta1_AWSFargateProfile_To_v1alpha4_AWSFargateProfile(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSFargateProfileList receiver to a v1beta1 AWSFargateProfileList.
func ( src *AWSFargateProfileList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSFargateProfileList)
	return Convert_v1alpha4_AWSFargateProfileList_To_v1beta1_AWSFargateProfileList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSFargateProfileList receiver to v1alpha4 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSFargateProfileList)

	return Convert_v1beta1_AWSFargateProfileList_To_v1alpha4_AWSFargateProfileList(src, r, nil)
}

// Convert_v1alpha4_AMIReference_To_v1beta1_AMIReference converts the v1alpha4 AMIReference receiver to a v1beta1 AMIReference.
func Convert_v1alpha4_AMIReference_To_v1beta1_AMIReference(in *infrav1alpha4.AMIReference, out *infrav1beta1.AMIReference, s apiconversion.Scope) error {
	return infrav1alpha4.Convert_v1alpha4_AMIReference_To_v1beta1_AMIReference(in, out, s)
}

// Convert_v1beta1_AMIReference_To_v1alpha4_AMIReference converts the v1beta1 AMIReference receiver to a v1alpha4 AMIReference.
func Convert_v1beta1_AMIReference_To_v1alpha4_AMIReference(in *infrav1beta1.AMIReference, out *infrav1alpha4.AMIReference, s apiconversion.Scope) error {
	return infrav1alpha4.Convert_v1beta1_AMIReference_To_v1alpha4_AMIReference(in, out, s)
}
