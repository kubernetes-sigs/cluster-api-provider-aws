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
	"reflect"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks"
)

const (
	maxProfileNameLength = 100
)

// SetupWebhookWithManager will setup the webhooks for the AWSFargateProfile
func (r *AWSFargateProfile) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha3-awsfargateprofile,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,versions=v1alpha3,name=default.awsfargateprofile.infrastructure.cluster.x-k8s.io,sideEffects=None
// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awsfargateprofile,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,versions=v1alpha3,name=validation.awsfargateprofile.infrastructure.cluster.x-k8s.io,sideEffects=None

var _ webhook.Defaulter = &AWSFargateProfile{}
var _ webhook.Validator = &AWSFargateProfile{}

// Default will set default values for the AWSFargateProfile
func (r *AWSFargateProfile) Default() {
	if r.Labels == nil {
		r.Labels = make(map[string]string)
	}
	r.Labels[clusterv1.ClusterLabelName] = r.Spec.ClusterName

	if r.Spec.ProfileName == "" {
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxProfileNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
			return
		}

		r.Spec.ProfileName = name
	}
}

func (r *AWSFargateProfile) ValidateUpdate(oldObj runtime.Object) error {
	gv := r.GroupVersionKind().GroupKind()
	old, ok := oldObj.(*AWSFargateProfile)
	if !ok {
		return apierrors.NewInvalid(gv, r.Name, field.ErrorList{
			field.InternalError(nil, errors.Errorf("failed to convert old %s to object", gv.Kind)),
		})
	}

	var allErrs field.ErrorList
	if !reflect.DeepEqual(old.Spec, r.Spec) {
		allErrs = append(
			allErrs,
			field.Invalid(field.NewPath("spec"), r.Spec, "is immutable"),
		)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		gv,
		r.Name,
		allErrs,
	)
}

func (r *AWSFargateProfile) ValidateCreate() error {
	return nil
}

func (r *AWSFargateProfile) ValidateDelete() error {
	return nil
}
