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
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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

func validateScaling(r *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.Scaling != nil { //nolint:nestif
		minField := field.NewPath("spec", "scaling", "minSize")
		maxField := field.NewPath("spec", "scaling", "maxSize")
		minSize := r.Spec.Scaling.MinSize
		maxSize := r.Spec.Scaling.MaxSize
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

func validateNodegroupUpdateConfig(r *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.UpdateConfig != nil {
		nodegroupUpdateConfigField := field.NewPath("spec", "updateConfig")

		if r.Spec.UpdateConfig.MaxUnavailable == nil && r.Spec.UpdateConfig.MaxUnavailablePercentage == nil {
			allErrs = append(allErrs, field.Invalid(nodegroupUpdateConfigField, r.Spec.UpdateConfig, "must specify one of maxUnavailable or maxUnavailablePercentage when using nodegroup updateconfig"))
		}

		if r.Spec.UpdateConfig.MaxUnavailable != nil && r.Spec.UpdateConfig.MaxUnavailablePercentage != nil {
			allErrs = append(allErrs, field.Invalid(nodegroupUpdateConfigField, r.Spec.UpdateConfig, "cannot specify both maxUnavailable and maxUnavailablePercentage"))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func validateRemoteAccess(r *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.RemoteAccess == nil {
		return allErrs
	}
	remoteAccessPath := field.NewPath("spec", "remoteAccess")
	sourceSecurityGroups := r.Spec.RemoteAccess.SourceSecurityGroups

	if public := r.Spec.RemoteAccess.Public; public && len(sourceSecurityGroups) > 0 {
		allErrs = append(
			allErrs,
			field.Invalid(remoteAccessPath.Child("sourceSecurityGroups"), sourceSecurityGroups, "must be empty if public is set"),
		)
	}

	return allErrs
}

func validateLaunchTemplate(r *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.AWSLaunchTemplate == nil {
		return allErrs
	}

	if r.Spec.InstanceType != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "InstanceType"), r.Spec.InstanceType, "InstanceType cannot be specified when LaunchTemplate is specified"))
	}
	if r.Spec.DiskSize != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "DiskSize"), r.Spec.DiskSize, "DiskSize cannot be specified when LaunchTemplate is specified"))
	}

	if r.Spec.AWSLaunchTemplate.IamInstanceProfile != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "AWSLaunchTemplate", "IamInstanceProfile"), r.Spec.AWSLaunchTemplate.IamInstanceProfile, "IAM instance profile in launch template is prohibited in EKS managed node group"))
	}

	return allErrs
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
	if errs := validateScaling(r); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateRemoteAccess(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateNodegroupUpdateConfig(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateLaunchTemplate(r); len(errs) > 0 {
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
	allErrs = append(allErrs, validateAMPImmutable(oldPool, r)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if errs := validateScaling(r); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateNodegroupUpdateConfig(r); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateLaunchTemplate(r); len(errs) > 0 {
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

func validateAMPImmutable(old *AWSManagedMachinePool, current *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList

	appendErrorIfMutated := func(old, update interface{}, name string) {
		if !cmp.Equal(old, update) {
			allErrs = append(
				allErrs,
				field.Invalid(field.NewPath("spec", name), update, "field is immutable"),
			)
		}
	}
	appendErrorIfSetAndMutated := func(old, update interface{}, name string) {
		if !reflect.ValueOf(old).IsZero() && !cmp.Equal(old, update) {
			allErrs = append(
				allErrs,
				field.Invalid(field.NewPath("spec", name), update, "field is immutable"),
			)
		}
	}

	if old.Spec.EKSNodegroupName != "" {
		appendErrorIfMutated(old.Spec.EKSNodegroupName, current.Spec.EKSNodegroupName, "eksNodegroupName")
	}
	appendErrorIfMutated(old.Spec.SubnetIDs, current.Spec.SubnetIDs, "subnetIDs")
	appendErrorIfSetAndMutated(old.Spec.RoleName, current.Spec.RoleName, "roleName")
	appendErrorIfMutated(old.Spec.DiskSize, current.Spec.DiskSize, "diskSize")
	appendErrorIfMutated(old.Spec.AMIType, current.Spec.AMIType, "amiType")
	appendErrorIfMutated(old.Spec.RemoteAccess, current.Spec.RemoteAccess, "remoteAccess")
	appendErrorIfSetAndMutated(old.Spec.CapacityType, current.Spec.CapacityType, "capacityType")
	appendErrorIfMutated(old.Spec.AvailabilityZones, current.Spec.AvailabilityZones, "availabilityZones")
	appendErrorIfMutated(old.Spec.AvailabilityZoneSubnetType, current.Spec.AvailabilityZoneSubnetType, "availabilityZoneSubnetType")
	if (old.Spec.AWSLaunchTemplate != nil && current.Spec.AWSLaunchTemplate == nil) ||
		(old.Spec.AWSLaunchTemplate == nil && current.Spec.AWSLaunchTemplate != nil) {
		allErrs = append(
			allErrs,
			field.Invalid(field.NewPath("spec", "AWSLaunchTemplate"), old.Spec.AWSLaunchTemplate, "field is immutable"),
		)
	}
	if old.Spec.AWSLaunchTemplate != nil && current.Spec.AWSLaunchTemplate != nil {
		appendErrorIfMutated(old.Spec.AWSLaunchTemplate.Name, current.Spec.AWSLaunchTemplate.Name, "awsLaunchTemplate.name")
	}

	return allErrs
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
		r.Spec.UpdateConfig = defaultManagedMachinePoolUpdateConfig()
	}
	return nil
}

func defaultManagedMachinePoolUpdateConfig() *UpdateConfig {
	return &UpdateConfig{
		MaxUnavailable: ptr.To[int](1),
	}
}
