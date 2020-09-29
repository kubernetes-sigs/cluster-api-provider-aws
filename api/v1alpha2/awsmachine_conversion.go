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

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this AWSMachine to the Hub version (v1alpha3).
func (src *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.AWSMachine)

	if err := Convert_v1alpha2_AWSMachine_To_v1alpha3_AWSMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data from annotations
	restored := &infrav1alpha3.AWSMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	restoreAWSMachineSpec(&restored.Spec, &dst.Spec, &src.Spec)
	restoreAWSMachineStatus(&restored.Status, &dst.Status)

	// Manual conversion for conditions
	dst.SetConditions(restored.GetConditions())
	return nil
}

func restoreAWSMachineSpec(restored, dst *infrav1alpha3.AWSMachineSpec, src *AWSMachineSpec) {
	dst.ImageLookupFormat = restored.ImageLookupFormat
	dst.ImageLookupBaseOS = restored.ImageLookupBaseOS
	dst.InstanceID = restored.InstanceID

	// Note this may override the manual conversion in Convert_v1alpha2_AWSMachineSpec_To_v1alpha3_AWSMachineSpec.
	if restored.RootVolume != nil {
		restored.RootVolume.DeepCopyInto(dst.RootVolume)
	}

	// override the SSHKeyName conversion if we are roundtripping from v1alpha3 and the v1alpha3 value is nil
	if dst.SSHKeyName != nil && *dst.SSHKeyName == "" && restored.SSHKeyName == nil {
		dst.SSHKeyName = nil
	}

	// manual conversion for UncompressedUserData
	dst.UncompressedUserData = restored.UncompressedUserData

	if restored.SpotMarketOptions != nil {
		dst.SpotMarketOptions = restored.SpotMarketOptions.DeepCopy()
	}

	if restored.NonRootVolumes != nil {
		dst.NonRootVolumes = []*infrav1alpha3.Volume{}
		for _, volume := range restored.NonRootVolumes {
			dst.NonRootVolumes = append(dst.NonRootVolumes, volume.DeepCopy())
		}
	}
	if src.RootDeviceSize != 0 {
		if dst.RootVolume == nil {
			dst.RootVolume = &infrav1alpha3.Volume{
				Size: src.RootDeviceSize,
			}
		} else {
			dst.RootVolume.Size = src.RootDeviceSize
		}
	}

	dst.Tenancy = restored.Tenancy

	if restored.CloudInit.SecureSecretsBackend != "" {
		if src.CloudInit != nil {
			dst.CloudInit.SecureSecretsBackend = restored.CloudInit.SecureSecretsBackend
		}
	}
}

func restoreAWSMachineStatus(restored, dst *infrav1alpha3.AWSMachineStatus) {
	dst.Interruptible = restored.Interruptible
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1alpha3.AWSMachine)
	if err := Convert_v1alpha3_AWSMachine_To_v1alpha2_AWSMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this AWSMachineList to the Hub version (v1alpha3).
func (src *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1alpha3.AWSMachineList)
	return Convert_v1alpha2_AWSMachineList_To_v1alpha3_AWSMachineList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error { // nolint:golint,stylecheck
	src := srcRaw.(*infrav1alpha3.AWSMachineList)
	return Convert_v1alpha3_AWSMachineList_To_v1alpha2_AWSMachineList(src, dst, nil)
}

// Convert_v1alpha2_AWSMachineSpec_To_v1alpha3_AWSMachineSpec converts this AWSMachineSpec to the Hub version (v1alpha3).
// Requires manual conversion as infrav1alpha2.AWSMachineSpec.RootDeviceSize does not exist in AWSMachineSpec.
func Convert_v1alpha2_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in *AWSMachineSpec, out *infrav1alpha3.AWSMachineSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert dst.Spec.FailureDomain.
	out.FailureDomain = in.AvailabilityZone

	// Manually convert SSHKeyName
	out.SSHKeyName = pointer.StringPtr(in.SSHKeyName)

	if in.CloudInit != nil {
		if err := Convert_v1alpha2_CloudInit_To_v1alpha3_CloudInit(in.CloudInit, &out.CloudInit, s); err != nil {
			return err
		}
	}

	// Manually convert RootDeviceSize. This may be overridden by restoring / upconverting from annotation.
	if in.RootDeviceSize != 0 {
		out.RootVolume = &infrav1alpha3.Volume{
			Size: in.RootDeviceSize,
		}
	}

	return nil
}

// Convert_v1alpha3_AWSMachineSpec_To_v1alpha2_AWSMachineSpec converts from the Hub version (v1alpha3) of the AWSMachineSpec to this version.
// Requires manual conversion as infrav1alpha3.AWSMachineSpec.ImageLookupBaseOS does not exist in AWSMachineSpec.
func Convert_v1alpha3_AWSMachineSpec_To_v1alpha2_AWSMachineSpec(in *infrav1alpha3.AWSMachineSpec, out *AWSMachineSpec, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_AWSMachineSpec_To_v1alpha2_AWSMachineSpec(in, out, s); err != nil {
		return err
	}

	// Manually convert FailureDomain to AvailabilityZone.
	out.AvailabilityZone = in.FailureDomain

	// Manually convert SSHKeyName
	if in.SSHKeyName != nil {
		out.SSHKeyName = *in.SSHKeyName
	}

	if !reflect.DeepEqual(in.CloudInit, infrav1alpha3.CloudInit{}) {
		out.CloudInit = &CloudInit{}
		if err := Convert_v1alpha3_CloudInit_To_v1alpha2_CloudInit(&in.CloudInit, out.CloudInit, s); err != nil {
			return err
		}
	}

	if in.RootVolume != nil {
		out.RootDeviceSize = in.RootVolume.Size
	}

	// Discards ImageLookupBaseOS & ImageLookupFormat

	return nil
}

// Convert_v1alpha2_AWSMachineStatus_To_v1alpha3_AWSMachineStatus converts this AWSMachineStatus to the Hub version (v1alpha3).
func Convert_v1alpha2_AWSMachineStatus_To_v1alpha3_AWSMachineStatus(in *AWSMachineStatus, out *infrav1alpha3.AWSMachineStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_AWSMachineStatus_To_v1alpha3_AWSMachineStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Error fields to the Failure fields
	out.FailureMessage = in.ErrorMessage
	out.FailureReason = in.ErrorReason

	return nil
}

// Convert_v1alpha3_AWSMachineStatus_To_v1alpha2_AWSMachineStatus converts from the Hub version (v1alpha3) of the AWSMachineStatus to this version.
func Convert_v1alpha3_AWSMachineStatus_To_v1alpha2_AWSMachineStatus(in *infrav1alpha3.AWSMachineStatus, out *AWSMachineStatus, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_AWSMachineStatus_To_v1alpha2_AWSMachineStatus(in, out, s); err != nil {
		return err
	}

	// Manually convert the Failure fields to the Error fields
	out.ErrorMessage = in.FailureMessage
	out.ErrorReason = in.FailureReason

	return nil
}

// Convert_v1alpha2_Instance_To_v1alpha3_Instance converts this Instance to the Hub version (v1alpha3).
func Convert_v1alpha2_Instance_To_v1alpha3_Instance(in *Instance, out *infrav1alpha3.Instance, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_Instance_To_v1alpha3_Instance(in, out, s); err != nil {
		return err
	}

	// Manually convert RootDeviceSize. This may be overridden by restoring / upconverting from annotation.
	if in.RootDeviceSize != 0 {
		out.RootVolume = &infrav1alpha3.Volume{
			Size: in.RootDeviceSize,
		}
	}

	return nil
}

// Convert_v1alpha3_Instance_To_v1alpha2_Instance converts from the Hub version (v1alpha3) of the Instance to this version.
func Convert_v1alpha3_Instance_To_v1alpha2_Instance(in *infrav1alpha3.Instance, out *Instance, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_Instance_To_v1alpha2_Instance(in, out, s); err != nil {
		return err
	}

	if in.RootVolume != nil {
		out.RootDeviceSize = in.RootVolume.Size
	}

	return nil
}

func Convert_v1alpha2_CloudInit_To_v1alpha3_CloudInit(in *CloudInit, out *infrav1alpha3.CloudInit, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha2_CloudInit_To_v1alpha3_CloudInit(in, out, s); err != nil {
		return err
	}

	out.InsecureSkipSecretsManager = !in.EnableSecureSecretsManager

	return nil
}

func Convert_v1alpha3_CloudInit_To_v1alpha2_CloudInit(in *infrav1alpha3.CloudInit, out *CloudInit, s apiconversion.Scope) error {
	if err := autoConvert_v1alpha3_CloudInit_To_v1alpha2_CloudInit(in, out, s); err != nil {
		return err
	}

	out.EnableSecureSecretsManager = !in.InsecureSkipSecretsManager

	return nil
}
