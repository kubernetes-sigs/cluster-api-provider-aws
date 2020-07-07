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
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this AWSMachineTemplate to the Hub version (v1alpha3).
func (src *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.AWSMachineTemplate)
	if err := Convert_v1alpha2_AWSMachineTemplate_To_v1alpha3_AWSMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &infrav1alpha3.AWSMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	restoreAWSMachineSpec(&restored.Spec.Template.Spec, &dst.Spec.Template.Spec)

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1alpha3.AWSMachineTemplate)
	if err := Convert_v1alpha3_AWSMachineTemplate_To_v1alpha2_AWSMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this AWSMachineTemplateList to the Hub version (v1alpha3).
func (src *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.AWSMachineTemplateList)
	return Convert_v1alpha2_AWSMachineTemplateList_To_v1alpha3_AWSMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1alpha3.AWSMachineTemplateList)
	return Convert_v1alpha3_AWSMachineTemplateList_To_v1alpha2_AWSMachineTemplateList(src, dst, nil)
}
