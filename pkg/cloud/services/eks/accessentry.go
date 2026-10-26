/*
Copyright 2025 The Kubernetes Authors.

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

package eks

import (
	"context"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) reconcileAccessEntries(ctx context.Context) error {
	if len(s.scope.ControlPlane.Spec.AccessEntries) == 0 {
		s.scope.Info("no access entries defined, skipping reconcile")
		return nil
	}

	if s.scope.ControlPlane.Spec.AccessConfig == nil ||
		s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode == "" ||
		(s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode != ekscontrolplanev1.EKSAuthenticationModeAPI &&
			s.scope.ControlPlane.Spec.AccessConfig.AuthenticationMode != ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap) {
		s.scope.Info("access mode is not api or api_and_config_map, skipping reconcile")
		return nil
	}

	existingAccessEntries, err := s.getExistingAccessEntries(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list existing access entries")
	}

	for _, accessEntry := range s.scope.ControlPlane.Spec.AccessEntries {
		managed, exists := existingAccessEntries[accessEntry.PrincipalARN]
		switch {
		case !exists:
			if err := s.createAccessEntry(ctx, accessEntry); err != nil {
				return errors.Wrapf(err, "failed to create access entry for principal %s", accessEntry.PrincipalARN)
			}
		case managed:
			if err := s.updateAccessEntry(ctx, accessEntry); err != nil {
				return errors.Wrapf(err, "failed to update access entry for principal %s", accessEntry.PrincipalARN)
			}
		default:
			// The entry exists but was not created by CAPA, e.g. the cluster creator entry that
			// EKS adds for bootstrapClusterCreatorAdminPermissions. Creating it would fail with
			// ResourceInUseException on every reconcile, and updating or deleting it could remove
			// the controller's own access, so only associate the policies it is missing.
			s.scope.Debug("Access entry exists but is not managed by CAPA, only associating missing access policies",
				"principalARN", accessEntry.PrincipalARN)
			if err := s.associateMissingAccessPolicies(ctx, accessEntry); err != nil {
				return errors.Wrapf(err, "failed to associate access policies for unmanaged principal %s", accessEntry.PrincipalARN)
			}
		}
		delete(existingAccessEntries, accessEntry.PrincipalARN)
	}

	for principalArn, managed := range existingAccessEntries {
		if !managed {
			continue
		}
		if err := s.deleteAccessEntry(ctx, principalArn); err != nil {
			return errors.Wrapf(err, "failed to delete access entry for principal %s", principalArn)
		}
	}

	record.Event(s.scope.ControlPlane, "SuccessfulReconcileAccessEntries", "Reconciled access entries")
	return nil
}

// getExistingAccessEntries returns every access entry on the cluster, keyed by principal ARN,
// with true for entries CAPA manages (tagged as owned) and false for all others.
func (s *Service) getExistingAccessEntries(ctx context.Context) (map[string]bool, error) {
	existingAccessEntries := make(map[string]bool)
	var nextToken *string

	clusterName := s.scope.KubernetesClusterName()
	managedTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())

	for {
		input := &eks.ListAccessEntriesInput{
			ClusterName: &clusterName,
			NextToken:   nextToken,
		}

		output, err := s.EKSClient.ListAccessEntries(ctx, input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list access entries")
		}

		for _, principalArn := range output.AccessEntries {
			describeOutput, err := s.EKSClient.DescribeAccessEntry(ctx, &eks.DescribeAccessEntryInput{
				ClusterName:  &clusterName,
				PrincipalArn: &principalArn,
			})
			if err != nil {
				// Skipping the entry would make a managed entry look unmanaged or missing.
				return nil, errors.Wrapf(err, "failed to describe access entry for principal %s", principalArn)
			}

			_, managed := describeOutput.AccessEntry.Tags[managedTag]
			existingAccessEntries[principalArn] = managed
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return existingAccessEntries, nil
}

func (s *Service) createAccessEntry(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	clusterName := s.scope.KubernetesClusterName()

	additionalTags := s.scope.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := make(map[string]string)
	for k, v := range additionalTags {
		tags[k] = v
	}

	createInput := &eks.CreateAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
		Tags:         tags,
	}

	if len(accessEntry.KubernetesGroups) > 0 {
		createInput.KubernetesGroups = accessEntry.KubernetesGroups
	}

	if accessEntry.Type != "" {
		createInput.Type = accessEntry.Type.APIValue()
	}

	if accessEntry.Username != "" {
		createInput.Username = &accessEntry.Username
	}

	if _, err := s.EKSClient.CreateAccessEntry(ctx, createInput); err != nil {
		return errors.Wrapf(err, "failed to create access entry for principal %s", accessEntry.PrincipalARN)
	}

	if err := s.reconcileAccessPolicies(ctx, accessEntry); err != nil {
		return errors.Wrapf(err, "failed to reconcile access policies for principal %s", accessEntry.PrincipalARN)
	}

	return nil
}

func (s *Service) updateAccessEntry(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	clusterName := s.scope.KubernetesClusterName()
	describeInput := &eks.DescribeAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
	}

	describeOutput, err := s.EKSClient.DescribeAccessEntry(ctx, describeInput)
	if err != nil {
		return errors.Wrapf(err, "failed to describe access entry for principal %s", accessEntry.PrincipalARN)
	}

	// Normalize empty spec Type to STANDARD to match what EKS reports for a
	// default STANDARD entry, so an unset Type is not seen as drift.
	desiredType := accessEntry.Type
	if desiredType == "" {
		desiredType = ekscontrolplanev1.AccessEntryTypeStandard
	}

	// Type is the only immutable field, so it is the only one requiring a recreate.
	// UpdateAccessEntry accepts username and Kubernetes groups.
	if *desiredType.APIValue() != *describeOutput.AccessEntry.Type {
		if err = s.deleteAccessEntry(ctx, accessEntry.PrincipalARN); err != nil {
			return errors.Wrapf(err, "failed to delete access entry for principal %s during recreation", accessEntry.PrincipalARN)
		}

		if err = s.createAccessEntry(ctx, accessEntry); err != nil {
			return errors.Wrapf(err, "failed to recreate access entry for principal %s", accessEntry.PrincipalARN)
		}
		return nil
	}

	slices.Sort(accessEntry.KubernetesGroups)
	slices.Sort(describeOutput.AccessEntry.KubernetesGroups)

	updateInput := &eks.UpdateAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &accessEntry.PrincipalARN,
	}
	needsUpdate := false

	if !slices.Equal(accessEntry.KubernetesGroups, describeOutput.AccessEntry.KubernetesGroups) {
		updateInput.KubernetesGroups = accessEntry.KubernetesGroups
		needsUpdate = true
	}

	// An unset username means EKS generates one, so the generated value is not drift.
	existingUsername := ""
	if describeOutput.AccessEntry.Username != nil {
		existingUsername = *describeOutput.AccessEntry.Username
	}

	if accessEntry.Username != "" && accessEntry.Username != existingUsername {
		updateInput.Username = &accessEntry.Username
		needsUpdate = true
	}

	if needsUpdate {
		if _, err := s.EKSClient.UpdateAccessEntry(ctx, updateInput); err != nil {
			return errors.Wrapf(err, "failed to update access entry for principal %s", accessEntry.PrincipalARN)
		}
	}

	if err := s.reconcileAccessPolicies(ctx, accessEntry); err != nil {
		return errors.Wrapf(err, "failed to reconcile access policies for principal %s", accessEntry.PrincipalARN)
	}

	return nil
}

func (s *Service) deleteAccessEntry(ctx context.Context, principalArn string) error {
	clusterName := s.scope.KubernetesClusterName()

	if _, err := s.EKSClient.DeleteAccessEntry(ctx, &eks.DeleteAccessEntryInput{
		ClusterName:  &clusterName,
		PrincipalArn: &principalArn,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete access entry for principal %s", principalArn)
	}

	return nil
}

func (s *Service) reconcileAccessPolicies(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	// Normalize type: empty means STANDARD, which needs policy reconciliation
	entryType := accessEntry.Type
	if entryType == "" {
		entryType = ekscontrolplanev1.AccessEntryTypeStandard
	}
	if entryType == ekscontrolplanev1.AccessEntryTypeEC2Linux || entryType == ekscontrolplanev1.AccessEntryTypeEC2Windows {
		s.scope.Info("Skipping access policy reconciliation for EC2 access type", "principalARN", accessEntry.PrincipalARN)
		return nil
	}

	existingPolicies, err := s.getExistingAccessPolicies(ctx, accessEntry.PrincipalARN)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing access policies for principal %s", accessEntry.PrincipalARN)
	}

	clusterName := s.scope.KubernetesClusterName()

	for _, policy := range accessEntry.AccessPolicies {
		existing, alreadyAssociated := existingPolicies[policy.PolicyARN]
		if alreadyAssociated && s.policyScopeMatches(policy.AccessScope, existing.AccessScope) {
			delete(existingPolicies, policy.PolicyARN)
			continue
		}

		input := &eks.AssociateAccessPolicyInput{
			ClusterName:  &clusterName,
			PrincipalArn: &accessEntry.PrincipalARN,
			PolicyArn:    &policy.PolicyARN,
			AccessScope: &ekstypes.AccessScope{
				Type: ekstypes.AccessScopeType(policy.AccessScope.Type),
			},
		}

		if policy.AccessScope.Type == "namespace" && len(policy.AccessScope.Namespaces) > 0 {
			input.AccessScope.Namespaces = policy.AccessScope.Namespaces
		}

		if _, err := s.EKSClient.AssociateAccessPolicy(ctx, input); err != nil {
			return errors.Wrapf(err, "failed to associate access policy %s", policy.PolicyARN)
		}

		delete(existingPolicies, policy.PolicyARN)
	}

	for policyARN := range existingPolicies {
		if _, err := s.EKSClient.DisassociateAccessPolicy(ctx, &eks.DisassociateAccessPolicyInput{
			ClusterName:  &clusterName,
			PrincipalArn: &accessEntry.PrincipalARN,
			PolicyArn:    &policyARN,
		}); err != nil {
			return errors.Wrapf(err, "failed to disassociate access policy %s", policyARN)
		}
	}

	return nil
}

// associateMissingAccessPolicies associates the policies of an access entry that CAPA does not
// manage when they are missing or have a different scope. Unlike reconcileAccessPolicies it never
// disassociates, so access granted outside CAPA (for example the cluster creator's admin policy)
// is kept.
func (s *Service) associateMissingAccessPolicies(ctx context.Context, accessEntry ekscontrolplanev1.AccessEntry) error {
	if accessEntry.Type == ekscontrolplanev1.AccessEntryTypeEC2Linux || accessEntry.Type == ekscontrolplanev1.AccessEntryTypeEC2Windows {
		return nil
	}
	if len(accessEntry.AccessPolicies) == 0 {
		return nil
	}

	existingPolicies, err := s.getExistingAccessPolicies(ctx, accessEntry.PrincipalARN)
	if err != nil {
		return errors.Wrapf(err, "failed to get existing access policies for principal %s", accessEntry.PrincipalARN)
	}

	clusterName := s.scope.KubernetesClusterName()

	for _, policy := range accessEntry.AccessPolicies {
		if existing, ok := existingPolicies[policy.PolicyARN]; ok && s.policyScopeMatches(policy.AccessScope, existing.AccessScope) {
			continue
		}

		input := &eks.AssociateAccessPolicyInput{
			ClusterName:  &clusterName,
			PrincipalArn: &accessEntry.PrincipalARN,
			PolicyArn:    &policy.PolicyARN,
			AccessScope: &ekstypes.AccessScope{
				Type: ekstypes.AccessScopeType(policy.AccessScope.Type),
			},
		}

		if policy.AccessScope.Type == ekscontrolplanev1.AccessScopeTypeNamespace && len(policy.AccessScope.Namespaces) > 0 {
			input.AccessScope.Namespaces = policy.AccessScope.Namespaces
		}

		if _, err := s.EKSClient.AssociateAccessPolicy(ctx, input); err != nil {
			return errors.Wrapf(err, "failed to associate access policy %s", policy.PolicyARN)
		}
	}

	return nil
}

func (s *Service) policyScopeMatches(desired ekscontrolplanev1.AccessScope, existing *ekstypes.AccessScope) bool {
	if existing == nil {
		return false
	}
	if string(desired.Type) != string(existing.Type) {
		return false
	}
	if desired.Type == ekscontrolplanev1.AccessScopeTypeNamespace {
		return slices.Equal(desired.Namespaces, existing.Namespaces)
	}
	return true
}

func (s *Service) getExistingAccessPolicies(ctx context.Context, principalARN string) (map[string]ekstypes.AssociatedAccessPolicy, error) {
	existingPolicies := map[string]ekstypes.AssociatedAccessPolicy{}
	var nextToken *string
	clusterName := s.scope.KubernetesClusterName()

	for {
		input := &eks.ListAssociatedAccessPoliciesInput{
			ClusterName:  &clusterName,
			PrincipalArn: &principalARN,
			NextToken:    nextToken,
		}

		output, err := s.EKSClient.ListAssociatedAccessPolicies(ctx, input)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list associated access policies for principal %s", principalARN)
		}

		for _, policy := range output.AssociatedAccessPolicies {
			existingPolicies[*policy.PolicyArn] = policy
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return existingPolicies, nil
}
