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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta2"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1beta1 AWSMachinePool receiver to a v1beta2 AWSMachinePool.
func (src *AWSMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePool)
	if err := Convert_v1beta1_AWSMachinePool_To_v1beta2_AWSMachinePool(src, dst, nil); err != nil {
		return err
	}
	
	return nil
}

// ConvertFrom converts the v1beta2 AWSMachinePool receiver to v1beta1 AWSMachinePool.
func (r *AWSMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePool)

	if err := Convert_v1beta2_AWSMachinePool_To_v1beta1_AWSMachinePool(src, r, nil); err != nil {
		return err
	}
	
	return nil
}

// ConvertTo converts the v1beta1 AWSMachinePoolList receiver to a v1beta2 AWSMachinePoolList.
func (src *AWSMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePoolList)
	return Convert_v1beta1_AWSMachinePoolList_To_v1beta2_AWSMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSMachinePoolList receiver to v1beta1 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePoolList)

	return Convert_v1beta2_AWSMachinePoolList_To_v1beta1_AWSMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSManagedMachinePool receiver to a v1beta2 AWSManagedMachinePool.
func (src *AWSManagedMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePool)
	if err := Convert_v1beta1_AWSManagedMachinePool_To_v1beta2_AWSManagedMachinePool(src, dst, nil); err != nil {
		return err
	}

	return nil
}

// ConvertFrom converts the v1beta2 AWSManagedMachinePool receiver to v1beta1 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePool)

	if err := Convert_v1beta2_AWSManagedMachinePool_To_v1beta1_AWSManagedMachinePool(src, r, nil); err != nil {
		return err
	}
	
	return nil
}

// Convert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec is a conversion function.
func Convert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec(in *infrav1exp.AWSManagedMachinePoolSpec, out *AWSManagedMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec(in, out, s)
}

// ConvertTo converts the v1beta1 AWSManagedMachinePoolList receiver to a v1beta2 AWSManagedMachinePoolList.
func (src *AWSManagedMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePoolList)
	return Convert_v1beta1_AWSManagedMachinePoolList_To_v1beta2_AWSManagedMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSManagedMachinePoolList receiver to v1beta1 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePoolList)

	return Convert_v1beta2_AWSManagedMachinePoolList_To_v1beta1_AWSManagedMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSFargateProfile receiver to a v1beta2 AWSFargateProfile.
func (src *AWSFargateProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfile)
	return Convert_v1beta1_AWSFargateProfile_To_v1beta2_AWSFargateProfile(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSFargateProfile receiver to v1beta1 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfile)

	return Convert_v1beta2_AWSFargateProfile_To_v1beta1_AWSFargateProfile(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSFargateProfileList receiver to a v1beta2 AWSFargateProfileList.
func (src *AWSFargateProfileList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfileList)
	return Convert_v1beta1_AWSFargateProfileList_To_v1beta2_AWSFargateProfileList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSFargateProfileList receiver to v1beta1 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfileList)

	return Convert_v1beta2_AWSFargateProfileList_To_v1beta1_AWSFargateProfileList(src, r, nil)
}

// Convert_v1beta1_AMIReference_To_v1beta2_AMIReference converts the v1beta1 AMIReference receiver to a v1beta2 AMIReference.
func Convert_v1beta1_AMIReference_To_v1beta2_AMIReference(in *infrav1beta1.AMIReference, out *infrav1.AMIReference, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta1_AMIReference_To_v1beta2_AMIReference(in, out, s)
}

// Convert_v1beta2_AMIReference_To_v1beta1_AMIReference converts the v1beta2 AMIReference receiver to a v1beta1 AMIReference.
func Convert_v1beta2_AMIReference_To_v1beta1_AMIReference(in *infrav1.AMIReference, out *infrav1beta1.AMIReference, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta2_AMIReference_To_v1beta1_AMIReference(in, out, s)
}

// Convert_v1beta2_Instance_To_v1beta1_Instance is a conversion function.
func Convert_v1beta2_Instance_To_v1beta1_Instance(in *infrav1.Instance, out *infrav1beta1.Instance, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta2_Instance_To_v1beta1_Instance(in, out, s)
}

// Convert_v1beta1_Instance_To_v1beta2_Instance is a conversion function.
func Convert_v1beta1_Instance_To_v1beta2_Instance(in *infrav1beta1.Instance, out *infrav1.Instance, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta1_Instance_To_v1beta2_Instance(in, out, s)
}

// Convert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate converts the v1beta2 AWSLaunchTemplate receiver to a v1beta1 AWSLaunchTemplate.
func Convert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate(in *infrav1exp.AWSLaunchTemplate, out *AWSLaunchTemplate, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate(in, out, s)
}
