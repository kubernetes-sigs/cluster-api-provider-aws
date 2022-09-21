/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var _ = logf.Log.WithName("awsclustercontrolleridentity-resource")

func (r *AWSClusterControllerIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-awsclustercontrolleridentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustercontrolleridentities,versions=v1beta1,name=validation.awsclustercontrolleridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-awsclustercontrolleridentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustercontrolleridentities,versions=v1beta1,name=default.awsclustercontrolleridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var (
	_ webhook.Validator = &AWSClusterControllerIdentity{}
	_ webhook.Defaulter = &AWSClusterControllerIdentity{}
)

// ValidateCreate will do any extra validation when creating an AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ValidateCreate() error {
	// Ensures AWSClusterControllerIdentity being singleton by only allowing "default" as name
	if r.Name != AWSClusterControllerIdentityName {
		return field.Invalid(field.NewPath("name"),
			r.Name, "AWSClusterControllerIdentity is a singleton and only acceptable name is default")
	}

	// Validate selector parses as Selector if AllowedNameSpaces is populated
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil
}

// ValidateDelete allows you to add any extra validation when deleting an AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ValidateDelete() error {
	return nil
}

// ValidateUpdate will do any extra validation when updating an AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) ValidateUpdate(old runtime.Object) error {
	oldP, ok := old.(*AWSClusterControllerIdentity)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an AWSClusterControllerIdentity but got a %T", old))
	}

	if !cmp.Equal(r.Spec, oldP.Spec) {
		return errors.New("AWSClusterControllerIdentity is immutable")
	}

	if r.Name != oldP.Name {
		return field.Invalid(field.NewPath("name"),
			r.Name, "AWSClusterControllerIdentity is a singleton and only acceptable name is default")
	}

	// Validate selector parses as Selector if AllowedNameSpaces is not nil
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return field.Invalid(field.NewPath("spec", "allowedNamespaces", "selectors"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil
}

// Default will set default values for the AWSClusterControllerIdentity.
func (r *AWSClusterControllerIdentity) Default() {
	SetDefaults_Labels(&r.ObjectMeta)
}
