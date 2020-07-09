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
	"net"
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

	allErrs = append(allErrs, r.validateAllowedCIDRBlocks()...)

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

	allErrs = append(allErrs, r.validateAllowedCIDRBlocks()...)

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSCluster) validateAllowedCIDRBlocks() []*field.Error {
	var errs field.ErrorList

	if r.Spec.Bastion.DisableIngressRules && len(r.Spec.Bastion.AllowedCIDRBlocks) > 0 {
		errs = append(errs,
			field.Forbidden(field.NewPath("spec", "bastion", "allowedCIDRBlocks"), "cannot be set if spec.bastion.disableIngressRules is true"),
		)
		return errs
	}

	for i, cidr := range r.Spec.Bastion.AllowedCIDRBlocks {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "bastion", fmt.Sprintf("allowedCIDRBlocks[%d]", i)), cidr, "must be a valid CIDR block"),
			)
		}
	}
	return errs
}

func (r *AWSCluster) Default() {
	// Default to allow open access to the bastion host if no CIDR Blocks have been set
	if len(r.Spec.Bastion.AllowedCIDRBlocks) == 0 && !r.Spec.Bastion.DisableIngressRules {
		r.Spec.Bastion.AllowedCIDRBlocks = []string{"0.0.0.0/0"}
	}

	// Default to Calico ingress rules if no rules have been set
	if r.Spec.NetworkSpec.CNI == nil {
		r.Spec.NetworkSpec.CNI = &CNISpec{
			CNIIngressRules: CNIIngressRules{
				{
					Description: "bgp (calico)",
					Protocol:    SecurityGroupProtocolTCP,
					FromPort:    179,
					ToPort:      179,
				},
				{
					Description: "IP-in-IP (calico)",
					Protocol:    SecurityGroupProtocolIPinIP,
					FromPort:    -1,
					ToPort:      65535,
				},
			},
		}
	}
}
