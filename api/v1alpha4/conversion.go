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
	"sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha4 AWSCluster receiver to a v1beta1 AWSCluster.
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSCluster)

	return Convert_v1alpha4_AWSCluster_To_v1beta1_AWSCluster(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSCluster receiver to a v1alpha4 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSCluster)

	return Convert_v1beta1_AWSCluster_To_v1alpha4_AWSCluster(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterList receiver to a v1beta1 AWSClusterList.
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterList)

	return Convert_v1alpha4_AWSClusterList_To_v1beta1_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterList receiver to a v1alpha4 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.AWSClusterList)

	return Convert_v1beta1_AWSClusterList_To_v1alpha4_AWSClusterList(src, r, nil)
}

// ConvertTo converts the v1alpha4 AWSMachine receiver to a v1beta1 AWSMachine.
func (src *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachine)
	return Convert_v1alpha4_AWSMachine_To_v1beta1_AWSMachine(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachine to a v1alpha4 AWSMachine.
func (dst *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachine)

	return Convert_v1beta1_AWSMachine_To_v1alpha4_AWSMachine(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSMachineList receiver to a v1beta1 AWSMachineList.
func (src *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineList)
	return Convert_v1alpha4_AWSMachineList_To_v1beta1_AWSMachineList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineList to a v1alpha4 AWSMachineList.
func (dst *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachineList)

	return Convert_v1beta1_AWSMachineList_To_v1alpha4_AWSMachineList(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSMachineTemplate receiver to a v1beta1 AWSMachineTemplate.
func (src *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineTemplate)
	return Convert_v1alpha4_AWSMachineTemplate_To_v1beta1_AWSMachineTemplate(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineTemplate to a v1alpha4 AWSMachineTemplate.
func (dst *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachineTemplate)

	return Convert_v1beta1_AWSMachineTemplate_To_v1alpha4_AWSMachineTemplate(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSMachineTemplateList receiver to a v1beta1 AWSMachineTemplateList.
func (src *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSMachineTemplateList)
	return Convert_v1alpha4_AWSMachineTemplateList_To_v1beta1_AWSMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineTemplateList to a v1alpha4 AWSMachineTemplateList.
func (dst *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSMachineTemplateList)

	return Convert_v1beta1_AWSMachineTemplateList_To_v1alpha4_AWSMachineTemplateList(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterStaticIdentity receiver to a v1beta1 AWSClusterStaticIdentity.
func (src *AWSClusterStaticIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterStaticIdentity)
	return Convert_v1alpha4_AWSClusterStaticIdentity_To_v1beta1_AWSClusterStaticIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentity to a v1alpha4 AWSClusterStaticIdentity.
func (dst *AWSClusterStaticIdentity) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterStaticIdentity)

	return Convert_v1beta1_AWSClusterStaticIdentity_To_v1alpha4_AWSClusterStaticIdentity(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterStaticIdentityList receiver to a v1beta1 AWSClusterStaticIdentityList.
func (src *AWSClusterStaticIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterStaticIdentityList)
	return Convert_v1alpha4_AWSClusterStaticIdentityList_To_v1beta1_AWSClusterStaticIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterStaticIdentityList to a v1alpha4 AWSClusterStaticIdentityList.
func (dst *AWSClusterStaticIdentityList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterStaticIdentityList)

	return Convert_v1beta1_AWSClusterStaticIdentityList_To_v1alpha4_AWSClusterStaticIdentityList(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterRoleIdentity receiver to a v1beta1 AWSClusterRoleIdentity.
func (src *AWSClusterRoleIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterRoleIdentity)
	return Convert_v1alpha4_AWSClusterRoleIdentity_To_v1beta1_AWSClusterRoleIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentity to a v1alpha4 AWSClusterRoleIdentity.
func (dst *AWSClusterRoleIdentity) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterRoleIdentity)

	return Convert_v1beta1_AWSClusterRoleIdentity_To_v1alpha4_AWSClusterRoleIdentity(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterRoleIdentityList receiver to a v1beta1 AWSClusterRoleIdentityList.
func (src *AWSClusterRoleIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterRoleIdentityList)
	return Convert_v1alpha4_AWSClusterRoleIdentityList_To_v1beta1_AWSClusterRoleIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterRoleIdentityList to a v1alpha4 AWSClusterRoleIdentityList.
func (dst *AWSClusterRoleIdentityList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterRoleIdentityList)

	return Convert_v1beta1_AWSClusterRoleIdentityList_To_v1alpha4_AWSClusterRoleIdentityList(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterControllerIdentity receiver to a v1beta1 AWSClusterControllerIdentity.
func (src *AWSClusterControllerIdentity) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterControllerIdentity)
	return Convert_v1alpha4_AWSClusterControllerIdentity_To_v1beta1_AWSClusterControllerIdentity(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentity to a v1alpha4 AWSClusterControllerIdentity.
func (dst *AWSClusterControllerIdentity) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterControllerIdentity)

	return Convert_v1beta1_AWSClusterControllerIdentity_To_v1alpha4_AWSClusterControllerIdentity(src, dst, nil)
}

// ConvertTo converts the v1alpha4 AWSClusterControllerIdentityList receiver to a v1beta1 AWSClusterControllerIdentityList.
func (src *AWSClusterControllerIdentityList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.AWSClusterControllerIdentityList)
	return Convert_v1alpha4_AWSClusterControllerIdentityList_To_v1beta1_AWSClusterControllerIdentityList(src, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSClusterControllerIdentityList to a v1alpha4 AWSClusterControllerIdentityList.
func (dst *AWSClusterControllerIdentityList) ConvertFrom(srcRaw conversion.Hub) error{
	src := srcRaw.(*v1beta1.AWSClusterControllerIdentityList)

	return Convert_v1beta1_AWSClusterControllerIdentityList_To_v1alpha4_AWSClusterControllerIdentityList(src, dst, nil)
}