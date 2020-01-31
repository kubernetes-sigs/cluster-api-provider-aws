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

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var _ = logf.Log.WithName("awsmachine-resource")

func (r *AWSMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmachine,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,versions=v1alpha3,name=validation.awsmachine.infrastructure.cluster.x-k8s.io

var _ webhook.Validator = &AWSMachine{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachine) ValidateCreate() error {
	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, r.validateCloudInitSecret())
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachine) ValidateUpdate(old runtime.Object) error {
	newAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		return apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new AWSMachine to unstructured object")),
		})
	}
	oldAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(old)
	if err != nil {
		return apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old AWSMachine to unstructured object")),
		})
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateCloudInitSecret()...)

	newAWSMachineSpec := newAWSMachine["spec"].(map[string]interface{})
	oldAWSMachineSpec := oldAWSMachine["spec"].(map[string]interface{})

	// allow changes to providerID
	delete(oldAWSMachineSpec, "providerID")
	delete(newAWSMachineSpec, "providerID")

	// allow changes to additionalTags
	delete(oldAWSMachineSpec, "additionalTags")
	delete(newAWSMachineSpec, "additionalTags")

	// allow changes to additionalSecurityGroups
	delete(oldAWSMachineSpec, "additionalSecurityGroups")
	delete(newAWSMachineSpec, "additionalSecurityGroups")

	// allow changes to secretARN
	if cloudInit, ok := oldAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretARN")
	}

	if cloudInit, ok := newAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretARN")
	}

	if !reflect.DeepEqual(oldAWSMachineSpec, newAWSMachineSpec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSMachine) validateCloudInitSecret() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.CloudInit.SecretARN != "" && r.Spec.CloudInit.InsecureSkipSecretsManager {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretARN"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
	}

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AWSMachine) ValidateDelete() error {
	return nil
}
