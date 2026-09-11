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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

var mpTemplateLog = ctrl.Log.WithName("awsmachinepooltemplate-resource")

const infrastructureGroup = "infrastructure.cluster.x-k8s.io"

// AWSMachinePoolTemplate implements a custom validation webhook for AWSMachinePoolTemplate.
type AWSMachinePoolTemplate struct{}

// SetupWebhookWithManager sets up the webhook with the Manager.
func (w *AWSMachinePoolTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &expinfrav1.AWSMachinePoolTemplate{}).
		WithCustomValidator(w).
		WithCustomDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepooltemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepooltemplates,versions=v1beta2,name=validation.awsmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinepooltemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepooltemplates,versions=v1beta2,name=default.awsmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &AWSMachinePoolTemplate{}
var _ webhook.CustomValidator = &AWSMachinePoolTemplate{}

func validateMachinePoolTemplateProviderIDs(spec *expinfrav1.AWSMachinePoolSpec, path *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if spec.ProviderID != "" {
		allErrs = append(allErrs, field.Forbidden(path.Child("providerID"), "cannot be set in templates"))
	}
	if spec.ProviderIDList != nil {
		allErrs = append(allErrs, field.Forbidden(path.Child("providerIDList"), "cannot be set in templates"))
	}

	return allErrs
}

// ValidateCreate will do any extra validation when creating a AWSMachinePoolTemplate.
func (w *AWSMachinePoolTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*expinfrav1.AWSMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachinePoolTemplate object but got %T", obj)
	}

	mpTemplateLog.Info("Validating AWSMachinePoolTemplate create", "name", r.Name)

	var allErrs field.ErrorList

	spec := &r.Spec.Template.Spec
	specPath := field.NewPath("spec", "template", "spec")

	allErrs = append(allErrs, validateMachinePoolTemplateProviderIDs(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolDefaultCoolDown(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolRootVolume(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolNonRootVolumes(spec, specPath)...)
	allErrs = append(allErrs, spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, validateMachinePoolSubnets(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolAdditionalSecurityGroups(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolSpotInstances(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolRefreshPreferences(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolInstanceMarketType(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolCapacityReservation(spec, specPath)...)
	allErrs = append(allErrs, validateLifecycleHooks(spec.AWSLifecycleHooks)...)
	allErrs = append(allErrs, validateMachinePoolIgnition(spec, specPath)...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: infrastructureGroup, Kind: "AWSMachinePoolTemplate"},
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSMachinePoolTemplate.
// AWSMachinePool enforces no immutability, so neither does its template: fields
// stay mutable and changes propagate in place to stamped machine pools. The
// resulting spec must still be valid, so update runs the same spec validation as
// create.
func (w *AWSMachinePoolTemplate) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*expinfrav1.AWSMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSMachinePoolTemplate object but got %T", newObj)
	}

	mpTemplateLog.Info("Validating AWSMachinePoolTemplate update", "name", r.Name)

	var allErrs field.ErrorList

	spec := &r.Spec.Template.Spec
	specPath := field.NewPath("spec", "template", "spec")

	allErrs = append(allErrs, validateMachinePoolTemplateProviderIDs(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolDefaultCoolDown(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolRootVolume(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolNonRootVolumes(spec, specPath)...)
	allErrs = append(allErrs, spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, validateMachinePoolSubnets(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolAdditionalSecurityGroups(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolSpotInstances(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolRefreshPreferences(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolInstanceMarketType(spec, specPath)...)
	allErrs = append(allErrs, validateMachinePoolCapacityReservation(spec, specPath)...)
	allErrs = append(allErrs, validateLifecycleHooks(spec.AWSLifecycleHooks)...)
	allErrs = append(allErrs, validateMachinePoolIgnition(spec, specPath)...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: infrastructureGroup, Kind: "AWSMachinePoolTemplate"},
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting.
func (w *AWSMachinePoolTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default will set default values for the AWSMachinePoolTemplate.
func (w *AWSMachinePoolTemplate) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*expinfrav1.AWSMachinePoolTemplate)
	if !ok {
		return fmt.Errorf("expected an AWSMachinePoolTemplate object but got %T", obj)
	}

	mpTemplateLog.Info("AWSMachinePoolTemplate setting defaults", "name", klog.KObj(r))

	spec := &r.Spec.Template.Spec

	if int(spec.DefaultCoolDown.Duration.Seconds()) == 0 {
		spec.DefaultCoolDown.Duration = 300 * time.Second
	}

	if int(spec.DefaultInstanceWarmup.Duration.Seconds()) == 0 {
		spec.DefaultInstanceWarmup.Duration = 300 * time.Second
	}

	if spec.Ignition != nil && spec.Ignition.StorageType == "" {
		spec.Ignition.StorageType = infrav1.DefaultMachinePoolIgnitionStorageType
	}
	// Defaults the version field if StorageType is not set to `UnencryptedUserData`.
	// When using `UnencryptedUserData` the version field is ignored because the userdata defines its version itself.
	if spec.Ignition != nil && spec.Ignition.Version == "" && spec.Ignition.StorageType != infrav1.IgnitionStorageTypeOptionUnencryptedUserData {
		spec.Ignition.Version = infrav1.DefaultIgnitionVersion
	}

	return nil
}
