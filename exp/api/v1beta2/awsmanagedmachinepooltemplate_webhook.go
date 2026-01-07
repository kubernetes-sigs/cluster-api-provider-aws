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

package v1beta2

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var mmpTemplateLog = ctrl.Log.WithName("awsmanagedmachinepooltemplate-resource")

// SetupWebhookWithManager sets up the webhook with the Manager.
func (r *AWSManagedMachinePoolTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsManagedMachinePoolTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepooltemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepooltemplates,versions=v1beta2,name=validation.awsmanagedmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepooltemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepooltemplates,versions=v1beta2,name=default.awsmanagedmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsManagedMachinePoolTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &awsManagedMachinePoolTemplateWebhook{}
var _ webhook.CustomValidator = &awsManagedMachinePoolTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsManagedMachinePoolTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", obj)
	}

	mmpTemplateLog.Info("Validating AWSManagedMachinePoolTemplate create", "name", r.Name)

	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec
	specPath := field.NewPath("spec", "template", "spec")

	allErrs = append(allErrs, validateManagedMachinePoolScaling(spec.Scaling, specPath.Child("scaling"))...)
	allErrs = append(allErrs, validateManagedMachinePoolUpdateConfig(spec.UpdateConfig, specPath.Child("updateConfig"))...)
	allErrs = append(allErrs, validateManagedMachinePoolRemoteAccess(spec.RemoteAccess, specPath.Child("remoteAccess"))...)
	allErrs = append(allErrs, validateManagedMachinePoolLaunchTemplate(spec.AWSLaunchTemplate, spec.InstanceType, spec.DiskSize, specPath)...)
	allErrs = append(allErrs, validateLifecycleHooks(spec.AWSLifecycleHooks)...)
	allErrs = append(allErrs, spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "AWSManagedMachinePoolTemplate"},
		r.Name,
		allErrs,
	)
}

func (*awsManagedMachinePoolTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	old, ok := oldObj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", oldObj)
	}

	newTemplate, ok := newObj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", newObj)
	}

	mmpTemplateLog.Info("Validating AWSManagedMachinePoolTemplate update", "name", newTemplate.Name)

	var allErrs field.ErrorList

	// All-or-nothing immutability for v0
	// TODO: Consider nuanced immutability in future to allow updates to:
	// - Scaling (minSize, maxSize)
	// - UpdateConfig
	// - Labels, Taints
	// - LifecycleHooks
	// See AWSManagedMachinePool.validateImmutable()
	if !cmp.Equal(old.Spec, newTemplate.Spec) {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec"),
			newTemplate.Spec,
			"AWSManagedMachinePoolTemplate.Spec is immutable",
		))
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "AWSManagedMachinePoolTemplate"},
			newTemplate.Name,
			allErrs,
		)
	}

	return nil, nil
}

func (*awsManagedMachinePoolTemplateWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", obj)
	}

	mmpTemplateLog.Info("Validating AWSManagedMachinePoolTemplate delete", "name", r.Name)

	return nil, nil
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*awsManagedMachinePoolTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", obj)
	}

	mmpTemplateLog.Info("AWSManagedMachinePoolTemplate setting defaults", "name", klog.KObj(r))

	if r.Spec.Template.Spec.UpdateConfig == nil {
		r.Spec.Template.Spec.UpdateConfig = &UpdateConfig{
			MaxUnavailable: ptr.To[int](1),
		}
	}

	return nil
}
