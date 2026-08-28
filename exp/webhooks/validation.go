/*
Copyright 2026 The Kubernetes Authors.

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

package webhooks

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/validation/field"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
)

func validateManagedMachinePoolScaling(scaling *expinfrav1.ManagedMachinePoolScaling, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if scaling != nil { //nolint:nestif
		minField := path.Child("minSize")
		maxField := path.Child("maxSize")
		minSize := scaling.MinSize
		maxSize := scaling.MaxSize
		if minSize != nil {
			if *minSize < 0 {
				allErrs = append(allErrs, field.Invalid(minField, *minSize, "must be greater or equal zero"))
			}
			if maxSize != nil && *maxSize < *minSize {
				allErrs = append(allErrs, field.Invalid(maxField, *maxSize, fmt.Sprintf("must be greater than field %s", minField.String())))
			}
		}
		if maxSize != nil && *maxSize < 0 {
			allErrs = append(allErrs, field.Invalid(maxField, *maxSize, "must be greater than zero"))
		}
	}
	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func validateManagedMachinePoolUpdateConfig(config *expinfrav1.UpdateConfig, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if config != nil {
		if config.MaxUnavailable == nil && config.MaxUnavailablePercentage == nil {
			allErrs = append(allErrs, field.Invalid(path, "", "must specify one of maxUnavailable or maxUnavailablePercentage when using nodegroup updateconfig"))
		}

		if config.MaxUnavailable != nil && config.MaxUnavailablePercentage != nil {
			allErrs = append(allErrs, field.Invalid(path, fmt.Sprintf("maxUnavailable=%d, maxUnavailablePercentage=%d", *config.MaxUnavailable, *config.MaxUnavailablePercentage), "cannot specify both maxUnavailable and maxUnavailablePercentage"))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func validateManagedMachinePoolRemoteAccess(access *expinfrav1.ManagedRemoteAccess, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if access == nil {
		return allErrs
	}
	sourceSecurityGroups := access.SourceSecurityGroups

	if public := access.Public; public && len(sourceSecurityGroups) > 0 {
		allErrs = append(
			allErrs,
			field.Invalid(path.Child("sourceSecurityGroups"), sourceSecurityGroups, "must be empty if public is set"),
		)
	}

	return allErrs
}

func validateManagedMachinePoolLaunchTemplate(lt *expinfrav1.AWSLaunchTemplate, instanceType *string, diskSize *int32, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if lt == nil {
		return allErrs
	}

	if instanceType != nil {
		allErrs = append(allErrs, field.Invalid(path.Child("InstanceType"), instanceType, "InstanceType cannot be specified when LaunchTemplate is specified"))
	}
	if diskSize != nil {
		allErrs = append(allErrs, field.Invalid(path.Child("DiskSize"), diskSize, "DiskSize cannot be specified when LaunchTemplate is specified"))
	}

	if lt.IamInstanceProfile != "" {
		allErrs = append(allErrs, field.Invalid(path.Child("AWSLaunchTemplate", "IamInstanceProfile"), lt.IamInstanceProfile, "IAM instance profile in launch template is prohibited in EKS managed node group"))
	}

	return allErrs
}

func validateLifecycleHooks(hooks []expinfrav1.AWSLifecycleHook) field.ErrorList {
	var allErrs field.ErrorList

	for _, hook := range hooks {
		if hook.Name == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.name"), "Name is required"))
		}
		if hook.NotificationTargetARN != nil && hook.RoleARN == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.roleARN"), "RoleARN is required if NotificationTargetARN is provided"))
		}
		if hook.RoleARN != nil && hook.NotificationTargetARN == nil {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.lifecycleHooks.notificationTargetARN"), "NotificationTargetARN is required if RoleARN is provided"))
		}
		if hook.LifecycleTransition != expinfrav1.LifecycleHookTransitionInstanceLaunching && hook.LifecycleTransition != expinfrav1.LifecycleHookTransitionInstanceTerminating {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.lifecycleTransition"), hook.LifecycleTransition, fmt.Sprintf("LifecycleTransition must be either %q or %q", expinfrav1.LifecycleHookTransitionInstanceLaunching, expinfrav1.LifecycleHookTransitionInstanceTerminating)))
		}
		if hook.DefaultResult != nil && (*hook.DefaultResult != expinfrav1.LifecycleHookDefaultResultContinue && *hook.DefaultResult != expinfrav1.LifecycleHookDefaultResultAbandon) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.defaultResult"), *hook.DefaultResult, fmt.Sprintf("DefaultResult must be either %s or %s", expinfrav1.LifecycleHookDefaultResultContinue, expinfrav1.LifecycleHookDefaultResultAbandon)))
		}
		if hook.HeartbeatTimeout != nil && (hook.HeartbeatTimeout.Seconds() < float64(30) || hook.HeartbeatTimeout.Seconds() > float64(172800)) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.heartbeatTimeout"), *hook.HeartbeatTimeout, "HeartbeatTimeout must be between 30 and 172800 seconds"))
		}
	}

	return allErrs
}

func validateManagedMachinePoolSpecImmutable(oldSpec, newSpec *expinfrav1.AWSManagedMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	appendErrorIfMutated := func(oldVal, newVal interface{}, name string) {
		if !cmp.Equal(oldVal, newVal) {
			allErrs = append(
				allErrs,
				field.Invalid(path.Child(name), newVal, "field is immutable"),
			)
		}
	}

	appendErrorIfSetAndMutated := func(oldVal, newVal interface{}, name string) {
		if !reflect.ValueOf(oldVal).IsZero() && !cmp.Equal(oldVal, newVal) {
			allErrs = append(
				allErrs,
				field.Invalid(path.Child(name), newVal, "field is immutable"),
			)
		}
	}

	if oldSpec.EKSNodegroupName != "" {
		appendErrorIfMutated(oldSpec.EKSNodegroupName, newSpec.EKSNodegroupName, "eksNodegroupName")
	}
	appendErrorIfMutated(oldSpec.SubnetIDs, newSpec.SubnetIDs, "subnetIDs")
	appendErrorIfSetAndMutated(oldSpec.RoleName, newSpec.RoleName, "roleName")
	appendErrorIfMutated(oldSpec.DiskSize, newSpec.DiskSize, "diskSize")
	appendErrorIfMutated(oldSpec.AMIType, newSpec.AMIType, "amiType")
	appendErrorIfMutated(oldSpec.RemoteAccess, newSpec.RemoteAccess, "remoteAccess")
	appendErrorIfSetAndMutated(oldSpec.CapacityType, newSpec.CapacityType, "capacityType")
	appendErrorIfMutated(oldSpec.AvailabilityZones, newSpec.AvailabilityZones, "availabilityZones")
	appendErrorIfMutated(oldSpec.AvailabilityZoneSubnetType, newSpec.AvailabilityZoneSubnetType, "availabilityZoneSubnetType")

	if (oldSpec.AWSLaunchTemplate != nil && newSpec.AWSLaunchTemplate == nil) ||
		(oldSpec.AWSLaunchTemplate == nil && newSpec.AWSLaunchTemplate != nil) {
		allErrs = append(
			allErrs,
			field.Invalid(path.Child("awsLaunchTemplate"), newSpec.AWSLaunchTemplate, "field is immutable"),
		)
	}
	if oldSpec.AWSLaunchTemplate != nil && newSpec.AWSLaunchTemplate != nil {
		appendErrorIfMutated(oldSpec.AWSLaunchTemplate.Name, newSpec.AWSLaunchTemplate.Name, "awsLaunchTemplate.name")
	}

	return allErrs
}

func validateMachinePoolDefaultCoolDown(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if int(spec.DefaultCoolDown.Duration.Seconds()) < 0 {
		allErrs = append(allErrs, field.Required(path.Child("defaultCoolDown"), "DefaultCoolDown must be greater than zero"))
	}

	return allErrs
}

func validateMachinePoolRootVolume(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if spec.AWSLaunchTemplate.RootVolume == nil {
		return allErrs
	}

	rootVolumePath := path.Child("awsLaunchTemplate", "rootVolume")

	if infrav1.VolumeTypesProvisioned.Has(string(spec.AWSLaunchTemplate.RootVolume.Type)) && spec.AWSLaunchTemplate.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(rootVolumePath.Child("iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if spec.AWSLaunchTemplate.RootVolume.Throughput != nil {
		if spec.AWSLaunchTemplate.RootVolume.Type != infrav1.VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(rootVolumePath.Child("throughput"), "throughput is valid only for type 'gp3'"))
		}
		// See https://aws.amazon.com/ebs/general-purpose/ for gp3 limits
		if *spec.AWSLaunchTemplate.RootVolume.Throughput < 125 || *spec.AWSLaunchTemplate.RootVolume.Throughput > 2000 {
			allErrs = append(allErrs, field.Required(rootVolumePath.Child("throughput"), "throughput must be between 125 Mib/s and 2000 MiB/s"))
		}
	}

	if spec.AWSLaunchTemplate.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func validateMachinePoolNonRootVolumes(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	nonRootVolumesPath := path.Child("awsLaunchTemplate", "nonRootVolumes")

	for _, volume := range spec.AWSLaunchTemplate.NonRootVolumes {
		if infrav1.VolumeTypesProvisioned.Has(string(volume.Type)) && volume.IOPS == 0 {
			allErrs = append(allErrs, field.Required(nonRootVolumesPath.Child("iops"), "iops required if type is 'io1' or 'io2'"))
		}

		if volume.Throughput != nil {
			if volume.Type != infrav1.VolumeTypeGP3 {
				allErrs = append(allErrs, field.Required(nonRootVolumesPath.Child("throughput"), "throughput is valid only for type 'gp3'"))
			}
			if *volume.Throughput < 0 {
				allErrs = append(allErrs, field.Required(nonRootVolumesPath.Child("throughput"), "throughput must be nonnegative"))
			}
		}

		if volume.DeviceName == "" {
			allErrs = append(allErrs, field.Required(nonRootVolumesPath.Child("deviceName"), "non root volume should have device name"))
		}
	}

	return allErrs
}

func validateMachinePoolSubnets(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if spec.Subnets == nil {
		return allErrs
	}

	for _, subnet := range spec.Subnets {
		if subnet.ID != nil && subnet.Filters != nil {
			allErrs = append(allErrs, field.Forbidden(path.Child("subnets", "filters"), "providing either subnet ID or filter is supported, should not provide both"))
			break
		}
	}

	return allErrs
}

func validateMachinePoolAdditionalSecurityGroups(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	for _, sg := range spec.AWSLaunchTemplate.AdditionalSecurityGroups {
		if sg.ID != nil && sg.Filters != nil {
			allErrs = append(allErrs, field.Forbidden(path.Child("awsLaunchTemplate", "additionalSecurityGroups"), "either ID or filters should be used"))
		}
	}
	return allErrs
}

func validateMachinePoolSpotInstances(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if spec.AWSLaunchTemplate.SpotMarketOptions != nil && spec.MixedInstancesPolicy != nil {
		allErrs = append(allErrs, field.Forbidden(path.Child("awsLaunchTemplate", "spotMarketOptions"), "either spec.awsLaunchTemplate.spotMarketOptions or spec.mixedInstancesPolicy should be used"))
	}
	return allErrs
}

func validateMachinePoolRefreshPreferences(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if spec.RefreshPreferences == nil {
		return allErrs
	}

	refreshPreferencesPath := path.Child("refreshPreferences")

	if spec.RefreshPreferences.MaxHealthyPercentage != nil && spec.RefreshPreferences.MinHealthyPercentage == nil {
		allErrs = append(allErrs, field.Forbidden(refreshPreferencesPath.Child("maxHealthyPercentage"), "If you specify spec.refreshPreferences.maxHealthyPercentage, you must also specify spec.refreshPreferences.minHealthyPercentage"))
	}

	if spec.RefreshPreferences.MaxHealthyPercentage != nil && spec.RefreshPreferences.MinHealthyPercentage != nil {
		if *spec.RefreshPreferences.MaxHealthyPercentage-*spec.RefreshPreferences.MinHealthyPercentage > 100 {
			allErrs = append(allErrs, field.Forbidden(refreshPreferencesPath.Child("maxHealthyPercentage"), "the difference between spec.refreshPreferences.maxHealthyPercentage and spec.refreshPreferences.minHealthyPercentage cannot be greater than 100"))
		}
	}

	return allErrs
}

func validateMachinePoolInstanceMarketType(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	marketTypePath := path.Child("awsLaunchTemplate", "marketType")

	if spec.AWSLaunchTemplate.MarketType == infrav1.MarketTypeCapacityBlock && spec.AWSLaunchTemplate.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(marketTypePath, "setting marketType to CapacityBlock and spotMarketOptions cannot be used together"))
	}
	if spec.AWSLaunchTemplate.MarketType == infrav1.MarketTypeOnDemand && spec.AWSLaunchTemplate.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(marketTypePath, "setting marketType to OnDemand and spotMarketOptions cannot be used together"))
	}

	if spec.AWSLaunchTemplate.MarketType == infrav1.MarketTypeCapacityBlock && spec.AWSLaunchTemplate.CapacityReservationID == nil {
		allErrs = append(allErrs, field.Forbidden(path.Child("awsLaunchTemplate", "capacityReservationID"), "is required when CapacityBlock is provided"))
	}
	switch spec.AWSLaunchTemplate.MarketType {
	case "", infrav1.MarketTypeOnDemand, infrav1.MarketTypeSpot, infrav1.MarketTypeCapacityBlock:
	default:
		allErrs = append(allErrs, field.Invalid(marketTypePath, spec.AWSLaunchTemplate.MarketType, fmt.Sprintf("Valid values are: %s, %s, %s and omitted", infrav1.MarketTypeOnDemand, infrav1.MarketTypeSpot, infrav1.MarketTypeCapacityBlock)))
	}
	if spec.AWSLaunchTemplate.MarketType == infrav1.MarketTypeSpot && spec.AWSLaunchTemplate.CapacityReservationID != nil {
		allErrs = append(allErrs, field.Forbidden(marketTypePath, "cannot be set to 'Spot' when CapacityReservationID is specified"))
	}

	if spec.AWSLaunchTemplate.CapacityReservationID != nil && spec.AWSLaunchTemplate.SpotMarketOptions != nil {
		allErrs = append(allErrs, field.Forbidden(path.Child("awsLaunchTemplate", "spotMarketOptions"), "cannot be set to when CapacityReservationID is specified"))
	}

	return allErrs
}

func validateMachinePoolCapacityReservation(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if spec.AWSLaunchTemplate.CapacityReservationID != nil &&
		spec.AWSLaunchTemplate.CapacityReservationPreference != infrav1.CapacityReservationPreferenceOnly &&
		spec.AWSLaunchTemplate.CapacityReservationPreference != "" {
		allErrs = append(allErrs, field.Forbidden(path.Child("awsLaunchTemplate", "capacityReservationPreference"), "when capacityReservationId is specified, capacityReservationPreference may only be `CapacityReservationsOnly` or empty"))
	}
	return allErrs
}

func validateMachinePoolIgnition(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) && spec.Ignition != nil {
		allErrs = append(allErrs, field.Forbidden(path.Child("ignition"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	return allErrs
}
