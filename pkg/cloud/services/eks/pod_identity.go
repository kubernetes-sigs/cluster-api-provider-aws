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

package eks

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	ekspodidentities "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks/podidentities"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) reconcilePodIdentities(ctx context.Context) error {
	s.scope.Info("Reconciling EKS Pod Identities")

	eksClusterName := s.scope.KubernetesClusterName()

	// Get existing eks pod identities on the cluster
	currentAssociations, err := s.listEksPodIdentities(ctx, eksClusterName)
	if err != nil {
		s.Error(err, "failed listing eks pod identity assocations")
		return fmt.Errorf("listing eks pod identity assocations: %w", err)
	}

	if len(currentAssociations) == 0 && len(s.scope.ControlPlane.Spec.PodIdentityAssociations) == 0 {
		s.scope.Debug("no eks pod identities found, no action needed")
		return nil
	}

	s.scope.Debug("eks pod identities found, creating reconciliation plan")
	desiredAssociations := s.translateAPIToPodAssociation(s.scope.ControlPlane.Spec.PodIdentityAssociations)
	existingAssociations := s.translateAWSToPodAssociation(currentAssociations)

	s.scope.Debug("creating eks pod identity association plan", "cluster", eksClusterName)
	podAssociationsPlan := ekspodidentities.NewPlan(eksClusterName, desiredAssociations, existingAssociations, s.EKSClient)
	procedures, err := podAssociationsPlan.Create(ctx)
	if err != nil {
		s.scope.Error(err, "failed creating eks pod identity association plan")
		return fmt.Errorf("creating eks pod identity association plan: %w", err)
	}
	for _, procedure := range procedures {
		s.scope.Debug("Executing pod association procedure", "name", procedure.Name())
		if err := procedure.Do(ctx); err != nil {
			s.scope.Error(err, "failed executing pod association procedure", "name", procedure.Name())
			return fmt.Errorf("executing pod association procedure %s: %w", procedure.Name(), err)
		}
	}

	record.Eventf(s.scope.ControlPlane, "SuccessfulReconcileEKSClusterPodIdentityAssociations", "Reconciled Pod Identity associations for EKS Cluster %s", s.scope.KubernetesClusterName())
	s.scope.Info("Reconcile EKS pod identity associations completed successfully")

	return nil
}

func (s *Service) listEksPodIdentities(ctx context.Context, eksClusterName string) ([]types.PodIdentityAssociationSummary, error) {
	s.Debug("getting list of associated eks pod identities")

	input := &eks.ListPodIdentityAssociationsInput{
		ClusterName: &eksClusterName,
	}

	output, err := s.EKSClient.ListPodIdentityAssociations(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("listing eks pod identity assocations: %w", err)
	}

	return output.Associations, nil
}

func (s *Service) translateAPIToPodAssociation(assocs []ekscontrolplanev1.PodIdentityAssociation) []ekspodidentities.EKSPodIdentityAssociation {
	converted := []ekspodidentities.EKSPodIdentityAssociation{}

	for _, assoc := range assocs {
		a := assoc
		c := ekspodidentities.EKSPodIdentityAssociation{
			ServiceAccountName:      a.ServiceAccountName,
			ServiceAccountNamespace: a.ServiceAccountNamespace,
			RoleARN:                 a.RoleARN,
		}

		converted = append(converted, c)
	}

	return converted
}

func (s *Service) translateAWSToPodAssociation(assocs []types.PodIdentityAssociationSummary) []ekspodidentities.EKSPodIdentityAssociation {
	converted := []ekspodidentities.EKSPodIdentityAssociation{}

	for _, assoc := range assocs {
		c := ekspodidentities.EKSPodIdentityAssociation{
			ServiceAccountName:      *assoc.ServiceAccount,
			ServiceAccountNamespace: *assoc.Namespace,
			RoleARN:                 *assoc.AssociationArn,
			AssociationID:           *assoc.AssociationId,
		}

		converted = append(converted, c)
	}

	return converted
}
