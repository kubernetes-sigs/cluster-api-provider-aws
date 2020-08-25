/*
Copyright 2020 The Kubernetes Authors.

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
	"fmt"
	"strings"

	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/version"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/hash"
)

const (
	// maxCharsName maximum number of characters for the name
	maxCharsName = 100

	clusterPrefix = "capa_"
)

// log is for logging in this package.
var mcpLog = logf.Log.WithName("awsmanagedcontrolplane-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSManagedControlPlane
func (r *AWSManagedControlPlane) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmanagedcontrolplane,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1alpha3,name=validation.awsmanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha3-awsmanagedcontrolplane,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,versions=v1alpha3,name=default.awsmanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None

var _ webhook.Defaulter = &AWSManagedControlPlane{}
var _ webhook.Validator = &AWSManagedControlPlane{}

func parseEKSVersion(raw string) (*version.Version, error) {
	v, err := version.ParseGeneric(raw)
	if err != nil {
		return nil, err
	}
	return version.MustParseGeneric(fmt.Sprintf("%d.%d", v.Major(), v.Minor())), nil
}

func normalizeVersion(raw string) (string, error) {
	// Normalize version (i.e. remove patch, add "v" prefix) if necessary
	eksV, err := parseEKSVersion(raw)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("v%d.%d", eksV.Major(), eksV.Minor()), nil
}

// ValidateCreate will do any extra validation when creating a AWSManagedControlPlane
func (r *AWSManagedControlPlane) ValidateCreate() error {
	mcpLog.Info("AWSManagedControlPlane validate create", "name", r.Name)

	var allErrs field.ErrorList

	if r.Spec.EKSClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksClusterName"), "eksClusterName is required"))
	}

	allErrs = append(allErrs, r.validateEKSVersion(nil)...)

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedControlPlane
func (r *AWSManagedControlPlane) ValidateUpdate(old runtime.Object) error {
	mcpLog.Info("AWSManagedControlPlane validate update", "name", r.Name)
	oldAWSManagedControlplane, ok := old.(*AWSManagedControlPlane)
	if !ok {
		return apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedControlPlane").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedControlPlane to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateEKSClusterName()...)
	allErrs = append(allErrs, r.validateEKSClusterNameSame(oldAWSManagedControlplane)...)
	allErrs = append(allErrs, r.validateEKSVersion(oldAWSManagedControlplane)...)

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting
func (r *AWSManagedControlPlane) ValidateDelete() error {
	mcpLog.Info("AWSManagedControlPlane validate delete", "name", r.Name)

	return nil
}

func (r *AWSManagedControlPlane) validateEKSClusterName() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.EKSClusterName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksClusterName"), "eksClusterName is required"))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSClusterNameSame(old *AWSManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.EKSClusterName != old.Spec.EKSClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec.eksClusterName"), r.Spec.EKSClusterName, "eksClusterName is different to current cluster name"))
	}

	return allErrs
}

func (r *AWSManagedControlPlane) validateEKSVersion(old *AWSManagedControlPlane) field.ErrorList {
	path := field.NewPath("spec.version")
	var allErrs field.ErrorList

	if r.Spec.Version == nil {
		return allErrs
	}

	v, err := parseEKSVersion(*r.Spec.Version)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(path, *r.Spec.Version, err.Error()))
	}

	if old != nil {
		oldV, err := parseEKSVersion(*old.Spec.Version)
		if err == nil && (v.Major() < oldV.Major() || v.Minor() < oldV.Minor()) {
			allErrs = append(allErrs, field.Invalid(path, *r.Spec.Version, "new version less than old version"))
		}
	}

	return allErrs
}

// Default will set default values for the AWSManagedControlPlane
func (r *AWSManagedControlPlane) Default() {
	mcpLog.Info("AWSManagedControlPlane setting defaults", "name", r.Name)

	if r.Spec.EKSClusterName == "" {
		mcpLog.Info("EKSClusterName is empty, generating name")
		name, err := generateEKSName(r.Name, r.Namespace)
		if err != nil {
			mcpLog.Error(err, "failed to create EKS cluster name")
			return
		}

		mcpLog.Info("defaulting EKS cluster name", "cluster-name", name)
		r.Spec.EKSClusterName = name
	}

	// Normalize version (i.e. remove patch, add "v" prefix) if necessary
	if r.Spec.Version != nil {
		normalizedV, err := normalizeVersion(*r.Spec.Version)
		if err != nil {
			mcpLog.Error(err, "couldn't parse version")
			return
		}
		r.Spec.Version = &normalizedV
	}
}

// generateEKSName generates a name of the EKS cluster
func generateEKSName(clusterName, namespace string) (string, error) {
	escapedName := strings.Replace(clusterName, ".", "_", -1)
	eksName := fmt.Sprintf("%s_%s", namespace, escapedName)

	if len(eksName) < maxCharsName {
		return eksName, nil
	}

	hashLength := 32 - len(clusterPrefix)
	hashedName, err := hash.Base36TruncatedHash(eksName, hashLength)
	if err != nil {
		return "", fmt.Errorf("creating hash from cluster name: %w", err)
	}

	return fmt.Sprintf("%s%s", clusterPrefix, hashedName), nil
}
