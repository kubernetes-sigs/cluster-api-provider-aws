/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (r *AWSMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates,versions=v1alpha3,name=validation.awsmachinetemplate.infrastructure.x-k8s.io,sideEffects=None

var (
	_ webhook.Validator = &AWSMachineTemplate{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachineTemplate) ValidateCreate() error {
	var allErrs field.ErrorList
	spec := r.Spec.Template.Spec

	if spec.CloudInit.SecretPrefix != "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretPrefix"), "cannot be set in templates"))
	}

	if spec.CloudInit.SecretCount != 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "cannot be set in templates"))
	}

	if spec.ProviderID != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "providerID"), "cannot be set in templates"))
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachineTemplate) ValidateUpdate(old runtime.Object) error {
	oldAWSMachineTemplate := old.(*AWSMachineTemplate)

	// Allow setting of cloudInit.secureSecretsBackend to "secrets-manager" only to handle v1alpha3 upgrade
	if oldAWSMachineTemplate.Spec.Template.Spec.CloudInit.SecureSecretsBackend == "" && r.Spec.Template.Spec.CloudInit.SecureSecretsBackend == SecretBackendSecretsManager {
		r.Spec.Template.Spec.CloudInit.SecureSecretsBackend = ""
	}

	if !reflect.DeepEqual(r.Spec, oldAWSMachineTemplate.Spec) {
		return apierrors.NewBadRequest("AWSMachineTemplate.Spec is immutable")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachineTemplate) ValidateDelete() error {
	return nil
}
