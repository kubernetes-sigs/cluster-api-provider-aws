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
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
)

const (
	maxNodegroupNameLength = 64
)

// log is for logging in this package.
var mmpLog = ctrl.Log.WithName("awsmanagedmachinepool-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSManagedMachinePool.
func (r *AWSManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsManagedMachinePoolWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=validation.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=default.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsManagedMachinePoolWebhook struct{}

var _ webhook.CustomDefaulter = &awsManagedMachinePoolWebhook{}
var _ webhook.CustomValidator = &awsManagedMachinePoolWebhook{}

func (r *AWSManagedMachinePool) validateScaling() field.ErrorList {
	return validateManagedMachinePoolScaling(r.Spec.Scaling, field.NewPath("spec", "scaling"))
}

func (r *AWSManagedMachinePool) validateNodegroupUpdateConfig() field.ErrorList {
	return validateManagedMachinePoolUpdateConfig(r.Spec.UpdateConfig, field.NewPath("spec", "updateConfig"))
}

func (r *AWSManagedMachinePool) validateRemoteAccess() field.ErrorList {
	return validateManagedMachinePoolRemoteAccess(r.Spec.RemoteAccess, field.NewPath("spec", "remoteAccess"))
}

func (r *AWSManagedMachinePool) validateLaunchTemplate() field.ErrorList {
	return validateManagedMachinePoolLaunchTemplate(r.Spec.AWSLaunchTemplate, r.Spec.InstanceType, r.Spec.DiskSize, field.NewPath("spec"))
}

func (r *AWSManagedMachinePool) validateLifecycleHooks() field.ErrorList {
	return validateLifecycleHooks(r.Spec.AWSLifecycleHooks)
}

// ValidateCreate will do any extra validation when creating a AWSManagedMachinePool.
func (*awsManagedMachinePoolWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePool object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePool validate create", "managed-machine-pool", klog.KObj(r))

	var allErrs field.ErrorList

	if r.Spec.EKSNodegroupName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksNodegroupName"), "eksNodegroupName is required"))
	}
	if errs := r.validateScaling(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateRemoteAccess(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateNodegroupUpdateConfig(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateLaunchTemplate(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateLifecycleHooks(); len(errs) > 0 {
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
func (*awsManagedMachinePoolWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedMachinePool object but got %T", r)
	}

	mmpLog.Info("AWSManagedMachinePool validate update", "managed-machine-pool", klog.KObj(r))
	oldPool, ok := oldObj.(*AWSManagedMachinePool)
	if !ok {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedMachinePool").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedMachinePool to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateImmutable(oldPool)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if errs := r.validateScaling(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateNodegroupUpdateConfig(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateLaunchTemplate(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateLifecycleHooks(); len(errs) > 0 {
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
func (*awsManagedMachinePoolWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (r *AWSManagedMachinePool) validateImmutable(old *AWSManagedMachinePool) field.ErrorList {
	return validateManagedMachinePoolSpecImmutable(&old.Spec, &r.Spec, field.NewPath("spec"))
}

// Default will set default values for the AWSManagedMachinePool.
func (*awsManagedMachinePoolWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSManagedMachinePool)
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
		r.Spec.UpdateConfig = &UpdateConfig{
			MaxUnavailable: ptr.To[int](1),
		}
	}
	return nil
}
