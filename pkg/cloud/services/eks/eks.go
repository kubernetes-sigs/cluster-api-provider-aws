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

package eks

import (
	"context"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"

	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
)

// ReconcileControlPlane reconciles a EKS control plane
func (s *Service) ReconcileControlPlane(ctx context.Context) error {
	s.scope.V(2).Info("Reconciling EKS control plane", "cluster-name", s.scope.Cluster.Name, "cluster-namespace", s.scope.Cluster.Namespace)

	// Control Plane IAM Role
	if err := s.reconcileControlPlaneIAMRole(); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, infrav1exp.IAMControlPlaneRolesReadyCondition, infrav1exp.IAMControlPlaneRolesReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.ControlPlane, infrav1exp.IAMControlPlaneRolesReadyCondition)

	// EKS Cluster
	if err := s.reconcileCluster(ctx); err != nil {
		conditions.MarkFalse(s.scope.ControlPlane, infrav1exp.EKSControlPlaneReadyCondition, infrav1exp.EKSControlPlaneReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.ControlPlane, infrav1exp.EKSControlPlaneReadyCondition)

	s.scope.V(2).Info("Reconcile EKS control plane completed successfully")
	return nil
}

// DeleteControlPlane deletes the EKS control plane.
func (s *Service) DeleteControlPlane() (err error) {
	s.scope.V(2).Info("Deleting EKS control plane")

	// EKS Cluster
	if err := s.deleteCluster(); err != nil {
		return err
	}

	// Control Plane IAM role
	if err := s.deleteControlPlaneIAMRole(); err != nil {
		return err
	}

	s.scope.V(2).Info("Delete EKS control plane completed successfully")
	return nil
}
