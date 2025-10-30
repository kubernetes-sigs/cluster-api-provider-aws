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
	"context"
	"fmt"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the webhooks for the AWSManagedMachinePool.
func (rt *AWSManagedMachinePoolTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsManagedMachinePoolTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(rt).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepooltemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepooltemplates,versions=v1beta2,name=validation.awsmanagedmachinepooltemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepooltemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepooltemplates,versions=v1beta2,name=default.awsmanagedmachinepooltemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
type awsManagedMachinePoolTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &awsManagedMachinePoolTemplateWebhook{}
var _ webhook.CustomValidator = &awsManagedMachinePoolTemplateWebhook{}

// ValidateCreate will do any extra validation when creating a AWSManagedMachinePoolTemplate.
func (*awsManagedMachinePoolTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", r)
	}
	mmpLog.Info("AWSManagedMachinePoolTemplate validate create", "managed-machine-pool", klog.KObj(r))

	var allErrs field.ErrorList

	if errs := validateScaling(r.Spec.Template); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateRemoteAccess(r.Spec.Template); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateNodegroupUpdateConfig(r.Spec.Template); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateLaunchTemplate(r.Spec.Template); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	allErrs = append(allErrs, r.Spec.Template.Spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when creating a AWSManagedMachinePoolTemplate.
func (*awsManagedMachinePoolTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePoolTemplate validate update", "managed-machine-pool", klog.KObj(r))

	mmpLog.Info("AWSManagedMachinePool validate update", "managed-machine-pool", klog.KObj(r))
	oldPool, ok := oldObj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedMachinePool").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedMachinePool to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateAMPImmutable(oldPool.Spec.Template, r.Spec.Template)...)
	allErrs = append(allErrs, r.Spec.Template.Spec.AdditionalTags.Validate()...)

	if errs := validateScaling(r.Spec.Template); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateNodegroupUpdateConfig(r.Spec.Template); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateLaunchTemplate(r.Spec.Template); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete will do any extra validation when creating a AWSManagedMachinePoolTemplate.
func (*awsManagedMachinePoolTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Default will set default values for the AWSManagedMachinePool.
func (*awsManagedMachinePoolTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSManagedMachinePoolTemplate)
	if !ok {
		return fmt.Errorf("expected an AWSManagedMachinePoolTemplate object but got %T", r)
	}
	if r.Spec.Template.Spec.UpdateConfig == nil {
		r.Spec.Template.Spec.UpdateConfig = defaultManagedMachinePoolUpdateConfig()
	}
	return nil
}
