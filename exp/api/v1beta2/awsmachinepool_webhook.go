/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

var log = ctrl.Log.WithName("awsmachinepool-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSMachinePool.
func (r *AWSMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1beta2,name=validation.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1beta2,name=default.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var (
	_ webhook.Defaulter = &AWSMachinePool{}
	_ webhook.Validator = &AWSMachinePool{}
)

func (r *AWSMachinePool) validateDefaultCoolDown() field.ErrorList {
	var allErrs field.ErrorList

	if int(r.Spec.DefaultCoolDown.Duration.Seconds()) < 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.DefaultCoolDown"), "DefaultCoolDown must be greater than zero"))
	}

	return allErrs
}

func (r *AWSMachinePool) validateRootVolume() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.AWSLaunchTemplate.RootVolume == nil {
		return allErrs
	}

	if v1beta2.VolumeTypesProvisioned.Has(string(r.Spec.AWSLaunchTemplate.RootVolume.Type)) && r.Spec.AWSLaunchTemplate.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.awsLaunchTemplate.rootVolume.iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if r.Spec.AWSLaunchTemplate.RootVolume.Throughput != nil {
		if r.Spec.AWSLaunchTemplate.RootVolume.Type != v1beta2.VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.awsLaunchTemplate.rootVolume.throughput"), "throughput is valid only for type 'gp3'"))
		}
		if *r.Spec.AWSLaunchTemplate.RootVolume.Throughput < 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.awsLaunchTemplate.rootVolume.throughput"), "throughput must be nonnegative"))
		}
	}

	if r.Spec.AWSLaunchTemplate.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func (r *AWSMachinePool) validateSubnets() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.Subnets == nil {
		return allErrs
	}

	for _, subnet := range r.Spec.Subnets {
		if subnet.ID != nil && subnet.Filters != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.subnets.filters"), "providing either subnet ID or filter is supported, should not provide both"))
			break
		}
	}

	return allErrs
}

func (r *AWSMachinePool) validateAdditionalSecurityGroups() field.ErrorList {
	var allErrs field.ErrorList
	for _, sg := range r.Spec.AWSLaunchTemplate.AdditionalSecurityGroups {
		if sg.ID != nil && sg.Filters != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.awsLaunchTemplate.AdditionalSecurityGroups"), "either ID or filters should be used"))
		}
	}
	return allErrs
}

func (r *AWSMachinePool) validateSpotInstances() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.AWSLaunchTemplate.SpotMarketOptions != nil && r.Spec.MixedInstancesPolicy != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.awsLaunchTemplate.spotMarketOptions"), "either spec.awsLaunchTemplate.spotMarketOptions or spec.mixedInstancesPolicy should be used"))
	}
	return allErrs
}

func (r *AWSMachinePool) validateLifecycleHooks() field.ErrorList {
	return validateLifecycleHooks(r.Spec.AWSLifecycleHooks)
}

func validateLifecycleHooks(hooks []AWSLifecycleHook) field.ErrorList {
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
		if hook.LifecycleTransition != LifecycleTransitionInstanceLaunch && hook.LifecycleTransition != LifecycleTransitionInstanceTerminate {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.lifecycleTransition"), hook.LifecycleTransition, "LifecycleTransition must be either EC2_INSTANCE_LAUNCHING or EC2_INSTANCE_TERMINATING"))
		}
		if hook.DefaultResult != nil && (*hook.DefaultResult != DefaultResultContinue && *hook.DefaultResult != DefaultResultAbandon) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.defaultResult"), *hook.DefaultResult, "DefaultResult must be either CONTINUE or ABANDON"))
		}
		if hook.HeartbeatTimeout != nil && (hook.HeartbeatTimeout.Seconds() < float64(30) || hook.HeartbeatTimeout.Seconds() > float64(172800)) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec.lifecycleHooks.heartbeatTimeout"), *hook.HeartbeatTimeout, "HeartbeatTimeout must be between 30 and 172800 seconds"))
		}
	}

	return allErrs
}

// ValidateCreate will do any extra validation when creating a AWSMachinePool.
func (r *AWSMachinePool) ValidateCreate() (admission.Warnings, error) {
	log.Info("AWSMachinePool validate create", "machine-pool", klog.KObj(r))

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateDefaultCoolDown()...)
	allErrs = append(allErrs, r.validateRootVolume()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateSubnets()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.validateSpotInstances()...)
	allErrs = append(allErrs, r.validateLifecycleHooks()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSMachinePool.
func (r *AWSMachinePool) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateDefaultCoolDown()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateSubnets()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.validateSpotInstances()...)
	allErrs = append(allErrs, r.validateLifecycleHooks()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting.
func (r *AWSMachinePool) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

// Default will set default values for the AWSMachinePool.
func (r *AWSMachinePool) Default() {
	if int(r.Spec.DefaultCoolDown.Duration.Seconds()) == 0 {
		log.Info("DefaultCoolDown is zero, setting 300 seconds as default")
		r.Spec.DefaultCoolDown.Duration = 300 * time.Second
	}

	if int(r.Spec.DefaultInstanceWarmup.Duration.Seconds()) == 0 {
		log.Info("DefaultInstanceWarmup is zero, setting 300 seconds as default")
		r.Spec.DefaultInstanceWarmup.Duration = 300 * time.Second
	}
}
