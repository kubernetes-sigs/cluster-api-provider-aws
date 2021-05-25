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
	"sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 EKSConfig receiver to a v1alpha4 EKSConfig.
func (r *EKSConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.EKSConfig)

	return Convert_v1alpha3_EKSConfig_To_v1alpha4_EKSConfig(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 EKSConfig receiver to a v1alpha3 EKSConfig.
func (r *EKSConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.EKSConfig)

	return Convert_v1alpha4_EKSConfig_To_v1alpha3_EKSConfig(src, r, nil)
}

// ConvertTo converts the v1alpha3 EKSConfigList receiver to a v1alpha4 EKSConfigList.
func (r *EKSConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.EKSConfigList)

	return Convert_v1alpha3_EKSConfigList_To_v1alpha4_EKSConfigList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 EKSConfigList receiver to a v1alpha3 EKSConfigList.
func (r *EKSConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.EKSConfigList)

	return Convert_v1alpha4_EKSConfigList_To_v1alpha3_EKSConfigList(src, r, nil)
}

// ConvertTo converts the v1alpha3 EKSConfigTemplate receiver to a v1alpha4 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.EKSConfigTemplate)

	return Convert_v1alpha3_EKSConfigTemplate_To_v1alpha4_EKSConfigTemplate(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 EKSConfigTemplate receiver to a v1alpha3 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.EKSConfigTemplate)

	return Convert_v1alpha4_EKSConfigTemplate_To_v1alpha3_EKSConfigTemplate(src, r, nil)
}

// ConvertTo converts the v1alpha3 EKSConfigTemplateList receiver to a v1alpha4 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.EKSConfigTemplateList)

	return Convert_v1alpha3_EKSConfigTemplateList_To_v1alpha4_EKSConfigTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 EKSConfigTemplateList receiver to a v1alpha3 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.EKSConfigTemplateList)

	return Convert_v1alpha4_EKSConfigTemplateList_To_v1alpha3_EKSConfigTemplateList(src, r, nil)
}
