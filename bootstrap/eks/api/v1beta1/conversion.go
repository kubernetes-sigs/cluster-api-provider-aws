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
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1beta2"
)

// ConvertTo converts the v1beta1 EKSConfig receiver to a v1beta2 EKSConfig.
func (r *EKSConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfig)

	if err := Convert_v1beta1_EKSConfig_To_v1beta2_EKSConfig(r, dst, nil); err != nil {
		return err
	}

	return nil
}

// ConvertFrom converts the v1beta2 EKSConfig receiver to a v1beta1 EKSConfig.
func (r *EKSConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfig)

	if err := Convert_v1beta2_EKSConfig_To_v1beta1_EKSConfig(src, r, nil); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 EKSConfigList receiver to a v1beta2 EKSConfigList.
func (r *EKSConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigList)

	return Convert_v1beta1_EKSConfigList_To_v1beta2_EKSConfigList(r, dst, nil)
}

// ConvertFrom converts the v1beta2 EKSConfigList receiver to a v1beta1 EKSConfigList.
func (r *EKSConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigList)

	return Convert_v1beta2_EKSConfigList_To_v1beta1_EKSConfigList(src, r, nil)
}

// ConvertTo converts the v1beta1 EKSConfigTemplate receiver to a v1beta2 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigTemplate)

	if err := Convert_v1beta1_EKSConfigTemplate_To_v1beta2_EKSConfigTemplate(r, dst, nil); err != nil {
		return err
	}

	return nil
}

// ConvertFrom converts the v1beta2 EKSConfigTemplate receiver to a v1beta1 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigTemplate)

	if err := Convert_v1beta2_EKSConfigTemplate_To_v1beta1_EKSConfigTemplate(src, r, nil); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 EKSConfigTemplateList receiver to a v1beta2 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigTemplateList)

	return Convert_v1beta1_EKSConfigTemplateList_To_v1beta2_EKSConfigTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1beta2 EKSConfigTemplateList receiver to a v1beta1 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigTemplateList)

	return Convert_v1beta2_EKSConfigTemplateList_To_v1beta1_EKSConfigTemplateList(src, r, nil)
}
