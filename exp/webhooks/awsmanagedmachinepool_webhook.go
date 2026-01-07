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

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
)

const (
	maxNodegroupNameLength = 64
)

// log is for logging in this package.
var mmpLog = ctrl.Log.WithName("awsmanagedmachinepool-resource")

// AWSManagedMachinePool implements a custom validation webhook for AWSManagedMachinePool.
type AWSManagedMachinePool struct{}

// SetupWebhookWithManager will setup the webhooks for the AWSManagedMachinePool.
func (w *AWSManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&expinfrav1.AWSManagedMachinePool{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=validation.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=default.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &AWSManagedMachinePool{}
var _ webhook.CustomValidator = &AWSManagedMachinePool{}

func (w *AWSManagedMachinePool) validateScaling(r *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolScaling(r.Spec.Scaling, field.NewPath("spec", "scaling"))
}

func (w *AWSManagedMachinePool) validateNodegroupUpdateConfig(r *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolUpdateConfig(r.Spec.UpdateConfig, field.NewPath("spec", "updateConfig"))
}

func (w *AWSManagedMachinePool) validateRemoteAccess(r *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolRemoteAccess(r.Spec.RemoteAccess, field.NewPath("spec", "remoteAccess"))
}

func (w *AWSManagedMachinePool) validateLaunchTemplate(r *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolLaunchTemplate(r.Spec.AWSLaunchTemplate, r.Spec.InstanceType, r.Spec.DiskSize, field.NewPath("spec"))
}

func (w *AWSManagedMachinePool) validateLifecycleHooks(r *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateLifecycleHooks(r.Spec.AWSLifecycleHooks)
}

// ValidateCreate will do any extra validation when creating a AWSManagedMachinePool.
func (w *AWSManagedMachinePool) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*expinfrav1.AWSManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePool object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePool validate create", "managed-machine-pool", klog.KObj(r))

	var allErrs field.ErrorList

	if r.Spec.EKSNodegroupName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksNodegroupName"), "eksNodegroupName is required"))
	}
	if errs := w.validateScaling(r); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateRemoteAccess(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateNodegroupUpdateConfig(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateLaunchTemplate(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateLifecycleHooks(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedMachinePool.
func (w *AWSManagedMachinePool) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*expinfrav1.AWSManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePool object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePool validate update", "managed-machine-pool", klog.KObj(r))
	oldPool, ok := oldObj.(*expinfrav1.AWSManagedMachinePool)
	if !ok {
		return nil, apierrors.NewInvalid(expinfrav1.GroupVersion.WithKind("AWSManagedMachinePool").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedMachinePool to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, w.validateImmutable(r, oldPool)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if errs := w.validateScaling(r); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateNodegroupUpdateConfig(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateLaunchTemplate(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := w.validateLifecycleHooks(r); len(errs) > 0 {
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

// ValidateDelete allows you to add any extra validation when deleting.
func (w *AWSManagedMachinePool) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (w *AWSManagedMachinePool) validateImmutable(r *expinfrav1.AWSManagedMachinePool, old *expinfrav1.AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolSpecImmutable(&old.Spec, &r.Spec, field.NewPath("spec"))
}

// Default will set default values for the AWSManagedMachinePool.
func (w *AWSManagedMachinePool) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*expinfrav1.AWSManagedMachinePool)
	if !ok {
		return fmt.Errorf("expected an AWSManagedMachinePool object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePool setting defaults", "managed-machine-pool", klog.KObj(r))

	if r.Spec.EKSNodegroupName == "" {
		mmpLog.Info("EKSNodegroupName is empty, generating name")
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxNodegroupNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
			return nil
		}

		mmpLog.Info("Generated EKSNodegroupName", "nodegroup", klog.KRef(r.Namespace, name))
		r.Spec.EKSNodegroupName = name
	}

	if r.Spec.UpdateConfig == nil {
		r.Spec.UpdateConfig = &expinfrav1.UpdateConfig{
			MaxUnavailable: ptr.To[int](1),
		}
	}
	return nil
}
