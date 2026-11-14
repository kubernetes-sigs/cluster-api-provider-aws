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
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

var log = ctrl.Log.WithName("awsmachinepool-resource")

// AWSMachinePool implements a custom validation webhook for AWSMachinePool.
type AWSMachinePool struct{}

// SetupWebhookWithManager will setup the webhooks for the AWSMachinePool.
func (w *AWSMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &expinfrav1.AWSMachinePool{}).
		WithCustomValidator(w).
		WithCustomDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1beta2,name=validation.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,versions=v1beta2,name=default.awsmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &AWSMachinePool{}
var _ webhook.CustomValidator = &AWSMachinePool{}

func (w *AWSMachinePool) validateDefaultCoolDown(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolDefaultCoolDown(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateRootVolume(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolRootVolume(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateNonRootVolumes(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolNonRootVolumes(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateSubnets(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolSubnets(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateAdditionalSecurityGroups(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolAdditionalSecurityGroups(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateSpotInstances(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolSpotInstances(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateRefreshPreferences(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolRefreshPreferences(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateLifecycleHooks(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateLifecycleHooks(r.Spec.AWSLifecycleHooks)
}

func (w *AWSMachinePool) ignitionEnabled(r *expinfrav1.AWSMachinePool) bool {
	return r.Spec.Ignition != nil
}

func (w *AWSMachinePool) validateIgnition(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolIgnition(&r.Spec, field.NewPath("spec"))
}

// ValidateCreate will do any extra validation when creating a AWSMachinePool.
func (w *AWSMachinePool) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*expinfrav1.AWSMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachinePool object but got %T", r)
	}

	log.Info("AWSMachinePool validate create", "machine-pool", klog.KObj(r))

	var allErrs field.ErrorList

	allErrs = append(allErrs, w.validateDefaultCoolDown(r)...)
	allErrs = append(allErrs, w.validateRootVolume(r)...)
	allErrs = append(allErrs, w.validateNonRootVolumes(r)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, w.validateSubnets(r)...)
	allErrs = append(allErrs, w.validateAdditionalSecurityGroups(r)...)
	allErrs = append(allErrs, w.validateSpotInstances(r)...)
	allErrs = append(allErrs, w.validateRefreshPreferences(r)...)
	allErrs = append(allErrs, w.validateInstanceMarketType(r)...)
	allErrs = append(allErrs, w.validateCapacityReservation(r)...)
	allErrs = append(allErrs, w.validateLifecycleHooks(r)...)
	allErrs = append(allErrs, w.validateIgnition(r)...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

func (w *AWSMachinePool) validateCapacityReservation(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolCapacityReservation(&r.Spec, field.NewPath("spec"))
}

func (w *AWSMachinePool) validateInstanceMarketType(r *expinfrav1.AWSMachinePool) field.ErrorList {
	return validateMachinePoolInstanceMarketType(&r.Spec, field.NewPath("spec"))
}

// ValidateUpdate will do any extra validation when updating a AWSMachinePool.
func (w *AWSMachinePool) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*expinfrav1.AWSMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachinePool object but got %T", r)
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, w.validateDefaultCoolDown(r)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, w.validateSubnets(r)...)
	allErrs = append(allErrs, w.validateAdditionalSecurityGroups(r)...)
	allErrs = append(allErrs, w.validateSpotInstances(r)...)
	allErrs = append(allErrs, w.validateRefreshPreferences(r)...)
	allErrs = append(allErrs, w.validateLifecycleHooks(r)...)

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
func (w *AWSMachinePool) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default will set default values for the AWSMachinePool.
func (w *AWSMachinePool) Default(ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*expinfrav1.AWSMachinePool)
	if !ok {
		return fmt.Errorf("expected an AWSMachinePool object but got %T", r)
	}

	if int(r.Spec.DefaultCoolDown.Duration.Seconds()) == 0 {
		log.Info("DefaultCoolDown is zero, setting 300 seconds as default")
		r.Spec.DefaultCoolDown.Duration = 300 * time.Second
	}

	if int(r.Spec.DefaultInstanceWarmup.Duration.Seconds()) == 0 {
		log.Info("DefaultInstanceWarmup is zero, setting 300 seconds as default")
		r.Spec.DefaultInstanceWarmup.Duration = 300 * time.Second
	}

	if w.ignitionEnabled(r) && r.Spec.Ignition.StorageType == "" {
		r.Spec.Ignition.StorageType = infrav1.DefaultMachinePoolIgnitionStorageType
	}
	// Defaults the version field if StorageType is not set to `UnencryptedUserData`.
	// When using `UnencryptedUserData` the version field is ignored because the userdata defines its version itself.
	if w.ignitionEnabled(r) && r.Spec.Ignition.Version == "" && r.Spec.Ignition.StorageType != infrav1.IgnitionStorageTypeOptionUnencryptedUserData {
		r.Spec.Ignition.Version = infrav1.DefaultIgnitionVersion
	}

	return nil
}
