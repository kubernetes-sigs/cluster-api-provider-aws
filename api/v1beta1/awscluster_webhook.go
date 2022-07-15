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
	"context"
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/annotations"
)

// log is for logging in this package.
var _ = logf.Log.WithName("awscluster-resource")

func (r *AWSCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/mutate-infrastructure-cluster-x-k8s-io-v1beta1-dockercluster", &webhook.Admission{
		Handler: &AWSClusterMutator{},
	})
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-awscluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1beta1,name=validation.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-awscluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,versions=v1beta1,name=default.awscluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Validator = &AWSCluster{}

// AWSClusterMutator used for defaulting AWSCluster.
// +kubebuilder:object:generate=false
type AWSClusterMutator struct {
	decoder *admission.Decoder
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSCluster) ValidateCreate() error {
	var allErrs field.ErrorList

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateSSHKeyName()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.Spec.S3Bucket.Validate()...)

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSCluster) ValidateDelete() error {
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
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

	newLoadBalancer := &AWSLoadBalancerSpec{}

	if r.Spec.ControlPlaneLoadBalancer != nil {
		newLoadBalancer = r.Spec.ControlPlaneLoadBalancer.DeepCopy()
	}

	if oldC.Spec.ControlPlaneLoadBalancer == nil {
		// If old scheme was nil, the only value accepted here is the default value: internet-facing
		if newLoadBalancer.Scheme != nil && newLoadBalancer.Scheme.String() != ClassicELBSchemeInternetFacing.String() {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "scheme"),
					r.Spec.ControlPlaneLoadBalancer.Scheme, "field is immutable, default value was set to internet-facing"),
			)
		}
	} else {
		// If old scheme was not nil, the new scheme should be the same.
		existingLoadBalancer := oldC.Spec.ControlPlaneLoadBalancer.DeepCopy()
		if !cmp.Equal(existingLoadBalancer.Scheme, newLoadBalancer.Scheme) {
			// Only allow changes from Internet-facing scheme to internet-facing.
			if !(existingLoadBalancer.Scheme.String() == ClassicELBSchemeIncorrectInternetFacing.String() &&
				newLoadBalancer.Scheme.String() == ClassicELBSchemeInternetFacing.String()) {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "scheme"),
						r.Spec.ControlPlaneLoadBalancer.Scheme, "field is immutable"),
				)
			}
		}
		// The name must be defined when the AWSCluster is created. If it is not defined,
		// then the controller generates a default name at runtime, but does not store it,
		// so the name remains nil. In either case, the name cannot be changed.
		if !cmp.Equal(existingLoadBalancer.Name, newLoadBalancer.Name) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "name"),
					r.Spec.ControlPlaneLoadBalancer.Name, "field is immutable"),
			)
		}

		// Block the update for HealthCheckProtocol :
		// - if it was not set in old spec but added in new spec
		// - if it was set in old spec but changed in new spec
		if !cmp.Equal(newLoadBalancer.HealthCheckProtocol, existingLoadBalancer.HealthCheckProtocol) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "healthCheckProtocol"),
					newLoadBalancer.HealthCheckProtocol, "field is immutable once set"),
			)
		}
	}

	if !cmp.Equal(oldC.Spec.ControlPlaneEndpoint, clusterv1.APIEndpoint{}) &&
		!cmp.Equal(r.Spec.ControlPlaneEndpoint, oldC.Spec.ControlPlaneEndpoint) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneEndpoint"), r.Spec.ControlPlaneEndpoint, "field is immutable"),
		)
	}

	// Modifying VPC id is not allowed because it will cause a new VPC creation if set to nil.
	if !cmp.Equal(oldC.Spec.NetworkSpec, NetworkSpec{}) &&
		!cmp.Equal(oldC.Spec.NetworkSpec.VPC, VPCSpec{}) &&
		oldC.Spec.NetworkSpec.VPC.ID != "" {
		if cmp.Equal(r.Spec.NetworkSpec, NetworkSpec{}) ||
			cmp.Equal(r.Spec.NetworkSpec.VPC, VPCSpec{}) ||
			oldC.Spec.NetworkSpec.VPC.ID != r.Spec.NetworkSpec.VPC.ID {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "network", "vpc", "id"),
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

	if annotations.IsExternallyManaged(oldC) && !annotations.IsExternallyManaged(r) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("metadata", "annotations"),
				r.Annotations, "removal of externally managed annotation is not allowed"),
		)
	}

	allErrs = append(allErrs, r.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.Spec.S3Bucket.Validate()...)

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// Default satisfies the defaulting webhook interface.
func (r *AWSCluster) Default() {
	SetObjectDefaults_AWSCluster(r)
}

// Handle will be used for defaulting.
func (r *AWSClusterMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	dct := &AWSCluster{}

	err := r.decoder.DecodeRaw(req.Object, dct)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, errors.Wrapf(err, "failed to decode DockerCluster resource"))
	}

	oldDct := &AWSCluster{}
	if req.Operation == admissionv1.Update {
		if err := r.decoder.DecodeRaw(req.OldObject, oldDct); err != nil {
			return admission.Errored(http.StatusBadRequest, errors.Wrapf(err, "failed to decode DockerCluster resource"))
		}
	}

	SetObjectDefaults_AWSCluster(dct)

	if req.Operation == admissionv1.Update {
		updateOpts := &metav1.UpdateOptions{}
		if err := r.decoder.DecodeRaw(req.Options, updateOpts); err != nil {
			return admission.Errored(http.StatusBadRequest, errors.Wrapf(err, "failed to decode UpdateOptions resource"))
		}

		// Variant 1: custom merge logic.
		// Assumptions about CAPA behavior:
		// * If subnets are empty => CAPA creates default subnets and adds them to subnets spec
		//   => Fine as in that case topology controller doesn't set subnets at all.
		// * If SecondaryCIDR block is set => CAPA adds a corresponding subnet, if there is no subnet with that CIDR yet
		//   => 1. We have to preserve additional subnets with the SecondaryCIDR block.
		// * CAPA queries AWS for the subnet and syncs fields from AWS to subnet spec
		//   => 2. We have to carry over fields from the current AWSCluster
		// * If subnet.ID == "" => create subnet in AWS and sync fields from new subnet from AWS to subnet spec
		//   => 2. We have to carry over fields from the current AWSCluster
		if updateOpts.FieldManager == "capi-topology" {
			for i := range dct.Spec.NetworkSpec.Subnets {
				oldSubnet := oldDct.Spec.NetworkSpec.Subnets.FindEqual(&dct.Spec.NetworkSpec.Subnets[i])
				if oldSubnet == nil {
					// Subnet has been newly added => nothing to carry-over.
					continue
				}

				// Subnet has been updated => 2. carry-over fields from old subnet, if they are not overwritten.
				if oldSubnet != nil {
					// This case should never happen when VPC is managed, because not intended to allow setting fields non BYOI case.
					if dct.Spec.NetworkSpec.Subnets[i].ID == "" {
						dct.Spec.NetworkSpec.Subnets[i].ID = oldSubnet.ID
					}
					if dct.Spec.NetworkSpec.Subnets[i].CidrBlock == "" {
						dct.Spec.NetworkSpec.Subnets[i].CidrBlock = oldSubnet.CidrBlock
					}
					if dct.Spec.NetworkSpec.Subnets[i].IsPublic == false {
						if oldSubnet.IsPublic {
							dct.Spec.NetworkSpec.Subnets[i].IsPublic = oldSubnet.IsPublic
						}
					}
					if dct.Spec.NetworkSpec.Subnets[i].AvailabilityZone == "" {
						dct.Spec.NetworkSpec.Subnets[i].AvailabilityZone = oldSubnet.AvailabilityZone
					}
					if dct.Spec.NetworkSpec.Subnets[i].RouteTableID == nil {
						dct.Spec.NetworkSpec.Subnets[i].RouteTableID = oldSubnet.RouteTableID
					}
					if dct.Spec.NetworkSpec.Subnets[i].NatGatewayID == nil {
						dct.Spec.NetworkSpec.Subnets[i].NatGatewayID = oldSubnet.NatGatewayID
					}
				}
			}
			// TODO: check if SecondaryCIDR subnets need to be checked.
		}
	}

	// Create the patch
	marshalled, err := json.Marshal(dct)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshalled)
}

// InjectDecoder injects the decoder.
// AWSClusterMutator implements admission.DecoderInjector.
// A decoder will be automatically injected.
func (r *AWSClusterMutator) InjectDecoder(d *admission.Decoder) error {
	r.decoder = d
	return nil
}

func (r *AWSCluster) validateSSHKeyName() field.ErrorList {
	return validateSSHKeyName(r.Spec.SSHKeyName)
}
