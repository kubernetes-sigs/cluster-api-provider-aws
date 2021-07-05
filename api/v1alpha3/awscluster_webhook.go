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
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var _ = logf.Log.WithName("awscluster-resource")

func (r *AWSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha3-awscluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1alpha3,name=validation.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha3-awscluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1alpha3,name=default.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None

var (
	_ webhook.Validator = &AWSCluster{}
	_ webhook.Defaulter = &AWSCluster{}
)

func (r *AWSCluster) ValidateCreate() error {
	var allErrs field.ErrorList

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateSSHKeyName()...)

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSCluster) ValidateDelete() error {
	return nil
}

func (r *AWSCluster) ValidateUpdate(old runtime.Object) error {
	var allErrs field.ErrorList

	oldC, ok := old.(*AWSCluster)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an AWSCluster but got a %T", old))
	}

	if r.Spec.Region != oldC.Spec.Region {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "region"), r.Spec.Region, "field is immutable"),
		)
	}

	existingLoadBalancer := &AWSLoadBalancerSpec{}
	newLoadBalancer := &AWSLoadBalancerSpec{}

	if oldC.Spec.ControlPlaneLoadBalancer != nil {
		existingLoadBalancer = oldC.Spec.ControlPlaneLoadBalancer.DeepCopy()
	}
	if r.Spec.ControlPlaneLoadBalancer != nil {
		newLoadBalancer = r.Spec.ControlPlaneLoadBalancer.DeepCopy()
	}
	if !reflect.DeepEqual(existingLoadBalancer.Scheme, newLoadBalancer.Scheme) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "scheme"),
				r.Spec.ControlPlaneLoadBalancer.Scheme, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(oldC.Spec.ControlPlaneEndpoint, clusterv1.APIEndpoint{}) &&
		!reflect.DeepEqual(r.Spec.ControlPlaneEndpoint, oldC.Spec.ControlPlaneEndpoint) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneEndpoint"), r.Spec.ControlPlaneEndpoint, "field is immutable"),
		)
	}

	// Modifying VPC id is not allowed because it will cause a new VPC creation if set to nil.
	if !reflect.DeepEqual(oldC.Spec.NetworkSpec, NetworkSpec{}) &&
		!reflect.DeepEqual(oldC.Spec.NetworkSpec.VPC, VPCSpec{}) &&
		oldC.Spec.NetworkSpec.VPC.ID != "" {
		if reflect.DeepEqual(r.Spec.NetworkSpec, NetworkSpec{}) ||
			reflect.DeepEqual(r.Spec.NetworkSpec.VPC, VPCSpec{}) ||
			oldC.Spec.NetworkSpec.VPC.ID != r.Spec.NetworkSpec.VPC.ID {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "networkSpec", "vpc", "id"),
					r.Spec.IdentityRef, "field cannot be modified once set"))
		}
	}

	// If a identityRef is already set, do not allow removal of it.
	if oldC.Spec.IdentityRef != nil && r.Spec.IdentityRef == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "identityRef"),
				r.Spec.IdentityRef, "field cannot be set to nil"),
		)
	}

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSCluster) Default() {
	SetDefaults_Bastion(&r.Spec.Bastion)
	SetDefaults_NetworkSpec(&r.Spec.NetworkSpec)

	if r.Spec.IdentityRef == nil {
		r.Spec.IdentityRef = &AWSIdentityReference{
			Kind: ControllerIdentityKind,
			Name: AWSClusterControllerIdentityName,
		}
	}
}

func (r *AWSCluster) validateSSHKeyName() field.ErrorList {
	return validateSSHKeyName(r.Spec.SSHKeyName)
}
