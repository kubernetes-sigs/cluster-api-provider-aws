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

package iamauth

import (
	"context"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	iamauthv1 "sigs.k8s.io/aws-iam-authenticator/pkg/mapper/crd/apis/iamauthenticator/v1alpha1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

// capaGenerateName is the GenerateName prefix used for every CAPA-managed
// IAMIdentityMapping CR. ReconcileMappings uses this prefix to distinguish
// CAPA-managed CRs (which are candidates for deletion when the desired set
// shrinks) from out-of-band CRs created directly against the
// aws-iam-authenticator CRD (which must never be touched).
const capaGenerateName = "capa-iamauth-"

type crdBackend struct {
	client crclient.Client
}

func (b *crdBackend) MapRole(mapping ekscontrolplanev1.RoleMapping) error {
	ctx := context.TODO()

	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}

	if err := b.client.List(ctx, &mappingList); err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if roleMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: capaGenerateName,
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.RoleARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

func (b *crdBackend) MapUser(mapping ekscontrolplanev1.UserMapping) error {
	ctx := context.TODO()

	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}

	if err := b.client.List(ctx, &mappingList); err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if userMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: capaGenerateName,
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.UserARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

// ReconcileMappings applies the given role and user mappings as the desired
// state. It lists all CAPA-managed IAMIdentityMapping CRs (identified by the
// capaGenerateName prefix), creates missing ones, and deletes CAPA-managed CRs
// that no longer match any desired entry. Out-of-band CRs (created directly
// against the aws-iam-authenticator CRD without the CAPA GenerateName prefix)
// are never touched.
func (b *crdBackend) ReconcileMappings(
	roles []ekscontrolplanev1.RoleMapping,
	users []ekscontrolplanev1.UserMapping,
) error {
	// Validate all inputs before touching state.
	for _, m := range roles {
		if errs := m.Validate(); errs != nil {
			return kerrors.NewAggregate(errs)
		}
	}
	for _, m := range users {
		if errs := m.Validate(); errs != nil {
			return kerrors.NewAggregate(errs)
		}
	}

	ctx := context.TODO()

	// Note: IAMIdentityMapping is cluster-scoped (`+genclient:nonNamespaced`
	// in aws-iam-authenticator; `scope: Cluster` in the CRD). Do not add an
	// InNamespace(kube-system) selector here — on a real cluster the API
	// server strips the Namespace field, so a namespaced List returns zero
	// items even though the code below sets Namespace on Create (a harmless
	// no-op on real clusters; the controller-runtime fake client preserves
	// the field, which is why the tests would pass under either variant).
	mappingList := iamauthv1.IAMIdentityMappingList{}
	if err := b.client.List(ctx, &mappingList); err != nil {
		return fmt.Errorf("listing IAMIdentityMappings: %w", err)
	}

	// Track which desired entries are already satisfied by an existing CR so
	// we don't create duplicates.
	desiredRoleMatched := make([]bool, len(roles))
	desiredUserMatched := make([]bool, len(users))
	var toDelete []*iamauthv1.IAMIdentityMapping

	for i := range mappingList.Items {
		existing := &mappingList.Items[i]

		// Preserve out-of-band CRs — anything without the CAPA GenerateName
		// prefix is left untouched.
		if !strings.HasPrefix(existing.GenerateName, capaGenerateName) {
			continue
		}

		matched := false
		for j, r := range roles {
			if !desiredRoleMatched[j] && roleMappingMatchesIAMMap(r, existing) {
				desiredRoleMatched[j] = true
				matched = true
				break
			}
		}
		if !matched {
			for j, u := range users {
				if !desiredUserMatched[j] && userMappingMatchesIAMMap(u, existing) {
					desiredUserMatched[j] = true
					matched = true
					break
				}
			}
		}
		if !matched {
			toDelete = append(toDelete, existing)
		}
	}

	// Delete stale CAPA-managed CRs first — reduces the window in which two
	// CRs could co-exist for the same ARN.
	for _, cr := range toDelete {
		if err := b.client.Delete(ctx, cr); err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("deleting stale IAMIdentityMapping %s: %w", cr.Name, err)
		}
	}

	// Create missing role mappings.
	for j, r := range roles {
		if desiredRoleMatched[j] {
			continue
		}
		cr := &iamauthv1.IAMIdentityMapping{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    metav1.NamespaceSystem,
				GenerateName: capaGenerateName,
			},
			Spec: iamauthv1.IAMIdentityMappingSpec{
				ARN:      r.RoleARN,
				Username: r.UserName,
				Groups:   r.Groups,
			},
		}
		if err := b.client.Create(ctx, cr); err != nil {
			return fmt.Errorf("creating role IAMIdentityMapping: %w", err)
		}
	}

	// Create missing user mappings.
	for j, u := range users {
		if desiredUserMatched[j] {
			continue
		}
		cr := &iamauthv1.IAMIdentityMapping{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    metav1.NamespaceSystem,
				GenerateName: capaGenerateName,
			},
			Spec: iamauthv1.IAMIdentityMappingSpec{
				ARN:      u.UserARN,
				Username: u.UserName,
				Groups:   u.Groups,
			},
		}
		if err := b.client.Create(ctx, cr); err != nil {
			return fmt.Errorf("creating user IAMIdentityMapping: %w", err)
		}
	}

	return nil
}

func roleMappingMatchesIAMMap(mapping ekscontrolplanev1.RoleMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.RoleARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func userMappingMatchesIAMMap(mapping ekscontrolplanev1.UserMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.UserARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}
