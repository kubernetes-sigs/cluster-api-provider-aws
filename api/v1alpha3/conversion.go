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
	"sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	apiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	apiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 AWSCluster receiver to a v1alpha4 AWSCluster.
func (r *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSCluster)

	return Convert_v1alpha3_AWSCluster_To_v1alpha4_AWSCluster(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSCluster receiver to a v1alpha3 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSCluster)

	return Convert_v1alpha4_AWSCluster_To_v1alpha3_AWSCluster(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterList receiver to a v1alpha4 AWSClusterList.
func (r *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterList)

	return Convert_v1alpha3_AWSClusterList_To_v1alpha4_AWSClusterList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterList receiver to a v1alpha3 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterList)

	return Convert_v1alpha4_AWSClusterList_To_v1alpha3_AWSClusterList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachine receiver to a v1alpha4 AWSMachine.
func (r *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachine)

	return Convert_v1alpha3_AWSMachine_To_v1alpha4_AWSMachine(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSMachine receiver to a v1alpha3 AWSMachine.
func (r *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachine)

	return Convert_v1alpha4_AWSMachine_To_v1alpha3_AWSMachine(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachineList receiver to a v1alpha4 AWSMachineList.
func (r *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachineList)

	return Convert_v1alpha3_AWSMachineList_To_v1alpha4_AWSMachineList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSMachineList receiver to a v1alpha3 AWSMachineList.
func (r *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachineList)

	return Convert_v1alpha4_AWSMachineList_To_v1alpha3_AWSMachineList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachineTemplate receiver to a v1alpha4 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachineTemplate)

	return Convert_v1alpha3_AWSMachineTemplate_To_v1alpha4_AWSMachineTemplate(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSMachineTemplate receiver to a v1alpha3 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachineTemplate)

	return Convert_v1alpha4_AWSMachineTemplate_To_v1alpha3_AWSMachineTemplate(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachineTemplateList receiver to a v1alpha4 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSMachineTemplateList)

	return Convert_v1alpha3_AWSMachineTemplateList_To_v1alpha4_AWSMachineTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSMachineTemplateList receiver to a v1alpha3 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSMachineTemplateList)

	return Convert_v1alpha4_AWSMachineTemplateList_To_v1alpha3_AWSMachineTemplateList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentity receiver to a v1alpha4 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterStaticIdentity)
	if err := Convert_v1alpha3_AWSClusterStaticIdentity_To_v1alpha4_AWSClusterStaticIdentity(r, dst, nil); err != nil {
		return err
	}

	dst.Spec.SecretRef = r.Spec.SecretRef.Name
	return nil
}

// ConvertFrom converts the v1alpha4 AWSClusterStaticIdentity receiver to a v1alpha3 AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterStaticIdentity)

	if err := Convert_v1alpha4_AWSClusterStaticIdentity_To_v1alpha3_AWSClusterStaticIdentity(src, r, nil); err != nil {
		return err
	}

	r.Spec.SecretRef.Name = src.Spec.SecretRef
	return nil
}

// ConvertTo converts the v1alpha3 AWSClusterStaticIdentityList receiver to a v1alpha4 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterStaticIdentityList)

	return Convert_v1alpha3_AWSClusterStaticIdentityList_To_v1alpha4_AWSClusterStaticIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterStaticIdentityList receiver to a v1alpha3 AWSClusterStaticIdentityList.
func (r *AWSClusterStaticIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterStaticIdentityList)

	return Convert_v1alpha4_AWSClusterStaticIdentityList_To_v1alpha3_AWSClusterStaticIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentity receiver to a v1alpha4 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterRoleIdentity)

	return Convert_v1alpha3_AWSClusterRoleIdentity_To_v1alpha4_AWSClusterRoleIdentity(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterRoleIdentity receiver to a v1alpha3 AWSClusterRoleIdentity.
func (r *AWSClusterRoleIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterRoleIdentity)

	return Convert_v1alpha4_AWSClusterRoleIdentity_To_v1alpha3_AWSClusterRoleIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterRoleIdentityList receiver to a v1alpha4 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterRoleIdentityList)

	return Convert_v1alpha3_AWSClusterRoleIdentityList_To_v1alpha4_AWSClusterRoleIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterRoleIdentityList receiver to a v1alpha3 AWSClusterRoleIdentityList.
func (r *AWSClusterRoleIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterRoleIdentityList)

	return Convert_v1alpha4_AWSClusterRoleIdentityList_To_v1alpha3_AWSClusterRoleIdentityList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentity receiver to a v1alpha4 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterControllerIdentity)

	return Convert_v1alpha3_AWSClusterControllerIdentity_To_v1alpha4_AWSClusterControllerIdentity(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterControllerIdentity receiver to a v1alpha3 AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterControllerIdentity)

	return Convert_v1alpha4_AWSClusterControllerIdentity_To_v1alpha3_AWSClusterControllerIdentity(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSClusterControllerIdentityList receiver to a v1alpha4 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.AWSClusterControllerIdentityList)

	return Convert_v1alpha3_AWSClusterControllerIdentityList_To_v1alpha4_AWSClusterControllerIdentityList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterControllerIdentityList receiver to a v1alpha3 AWSClusterControllerIdentityList.
func (r *AWSClusterControllerIdentityList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.AWSClusterControllerIdentityList)

	return Convert_v1alpha4_AWSClusterControllerIdentityList_To_v1alpha3_AWSClusterControllerIdentityList(src, r, nil)
}

// Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint is an autogenerated conversion function.
func Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint(in *apiv1alpha3.APIEndpoint, out *apiv1alpha4.APIEndpoint, s apiconversion.Scope) error {
	return apiv1alpha3.Convert_v1alpha3_APIEndpoint_To_v1alpha4_APIEndpoint(in, out, s)
}

// Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint is an autogenerated conversion function.
func Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint(in *apiv1alpha4.APIEndpoint, out *apiv1alpha3.APIEndpoint, s apiconversion.Scope) error {
	return apiv1alpha3.Convert_v1alpha4_APIEndpoint_To_v1alpha3_APIEndpoint(in, out, s)
}

// Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1alpha4_AWSClusterStaticIdentitySpec is an autogenerated conversion function.
func Convert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1alpha4_AWSClusterStaticIdentitySpec(in *AWSClusterStaticIdentitySpec, out *v1alpha4.AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_AWSClusterStaticIdentitySpec_To_v1alpha4_AWSClusterStaticIdentitySpec(in, out, s)
}

// Convert_v1alpha4_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec is an autogenerated conversion function.
func Convert_v1alpha4_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in *v1alpha4.AWSClusterStaticIdentitySpec, out *AWSClusterStaticIdentitySpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha4_AWSClusterStaticIdentitySpec_To_v1alpha3_AWSClusterStaticIdentitySpec(in, out, s)
}
