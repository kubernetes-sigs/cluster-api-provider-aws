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

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api/util/topology"
)

func (r *AWSMachineTemplateWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AWSMachineTemplate{}).
		WithValidator(r).
		Complete()
}

// AWSMachineTemplateWebhook implements a custom validation webhook for AWSMachineTemplate.
// +kubebuilder:object:generate=false
type AWSMachineTemplateWebhook struct{}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates,versions=v1beta2,name=validation.awsmachinetemplate.infrastructure.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
var _ webhook.CustomValidator = &AWSMachineTemplateWebhook{}

func (r *AWSMachineTemplate) validateRootVolume() field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec
	if spec.RootVolume == nil {
		return allErrs
	}

	if VolumeTypesProvisioned.Has(string(spec.RootVolume.Type)) && spec.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if spec.RootVolume.Throughput != nil {
		if spec.RootVolume.Type != VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.throughput"), "throughput is valid only for type 'gp3'"))
		}
		if *spec.RootVolume.Throughput < 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.throughput"), "throughput must be nonnegative"))
		}
	}

	if spec.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func (r *AWSMachineTemplate) validateNonRootVolumes() field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec

	for _, volume := range spec.NonRootVolumes {
		if VolumeTypesProvisioned.Has(string(volume.Type)) && volume.IOPS == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.iops"), "iops required if type is 'io1' or 'io2'"))
		}

		if volume.Throughput != nil {
			if volume.Type != VolumeTypeGP3 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.throughput"), "throughput is valid only for type 'gp3'"))
			}
			if *volume.Throughput < 0 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.throughput"), "throughput must be nonnegative"))
			}
		}

		if volume.DeviceName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.deviceName"), "non root volume should have device name"))
		}
	}

	return allErrs
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachineTemplateWebhook) ValidateCreate(_ context.Context, raw runtime.Object) error {
	var allErrs field.ErrorList
	obj, ok := raw.(*AWSMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereMachineTemplate but got a %T", raw))
	}

	spec := obj.Spec.Template.Spec

	if spec.CloudInit.SecretPrefix != "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretPrefix"), "cannot be set in templates"))
	}

	if spec.CloudInit.SecretCount != 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "cannot be set in templates"))
	}

	if spec.ProviderID != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "providerID"), "cannot be set in templates"))
	}

	allErrs = append(allErrs, obj.validateRootVolume()...)
	allErrs = append(allErrs, obj.validateNonRootVolumes()...)

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) && spec.Ignition != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ignition"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	cloudInitConfigured := spec.CloudInit.SecureSecretsBackend != "" || spec.CloudInit.InsecureSkipSecretsManager
	if cloudInitConfigured && spec.Ignition != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit"),
			"cannot be set if spec.template.spec.ignition is set"))
	}

	return aggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachineTemplateWebhook) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) error {
	newAWSMachineTemplate, ok := newRaw.(*AWSMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a AWSMachineTemplate but got a %T", newRaw))
	}
	oldAWSMachineTemplate, ok := oldRaw.(*AWSMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a AWSMachineTemplate but got a %T", oldRaw))
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	var allErrs field.ErrorList

	// Allow setting of cloudInit.secureSecretsBackend to "secrets-manager" only to handle v1beta2 upgrade
	if oldAWSMachineTemplate.Spec.Template.Spec.CloudInit.SecureSecretsBackend == "" && newAWSMachineTemplate.Spec.Template.Spec.CloudInit.SecureSecretsBackend == SecretBackendSSMParameterStore {
		newAWSMachineTemplate.Spec.Template.Spec.CloudInit.SecureSecretsBackend = ""
	}

	if !topology.ShouldSkipImmutabilityChecks(req, newAWSMachineTemplate) && !cmp.Equal(newAWSMachineTemplate.Spec, oldAWSMachineTemplate.Spec) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), newAWSMachineTemplate, "AWSMachineTemplate.Spec is immutable"))
	}

	return aggregateObjErrors(newAWSMachineTemplate.GroupVersionKind().GroupKind(), newAWSMachineTemplate.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachineTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) error {
	return nil
}
