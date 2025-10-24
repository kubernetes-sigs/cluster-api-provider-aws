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

// Package eks provides a service to reconcile EKS control plane and nodegroups.
package eks

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	ekscontrolplanev1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta1"
	expinfrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// ReconcileControlPlane reconciles a EKS control plane.
func (s *Service) ReconcileControlPlane(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS control plane", "cluster", klog.KRef(s.scope.Cluster.Namespace, s.scope.Cluster.Name))

	// Control Plane IAM Role
	if err := s.reconcileControlPlaneIAMRole(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1beta1.IAMControlPlaneRolesReadyCondition, ekscontrolplanev1beta1.IAMControlPlaneRolesReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1beta1.IAMControlPlaneRolesReadyCondition)

	// EKS Cluster
	if err := s.reconcileCluster(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSControlPlaneReadyCondition, ekscontrolplanev1beta1.EKSControlPlaneReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSControlPlaneReadyCondition)

	// EKS Addons
	if err := s.reconcileAddons(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSAddonsConfiguredCondition, ekscontrolplanev1beta1.EKSAddonsConfiguredFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		return errors.Wrap(err, "failed reconciling eks addons")
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSAddonsConfiguredCondition)

	// EKS Identity Provider
	if err := s.reconcileIdentityProvider(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSIdentityProviderConfiguredCondition, ekscontrolplanev1beta1.EKSIdentityProviderConfiguredFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return errors.Wrap(err, "failed reconciling eks identity provider")
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1beta1.EKSIdentityProviderConfiguredCondition)

	s.scope.Debug("Reconcile EKS control plane completed successfully")
	return nil
}

// DeleteControlPlane deletes the EKS control plane.
func (s *Service) DeleteControlPlane(ctx context.Context) (err error) {
	s.scope.Debug("Deleting EKS control plane")

	// EKS Cluster
	if err := s.deleteCluster(ctx); err != nil {
		return err
	}

	// Control Plane IAM role
	if err := s.deleteControlPlaneIAMRole(ctx); err != nil {
		return err
	}

	// OIDC Provider
	if err := s.deleteOIDCProvider(ctx); err != nil {
		return err
	}

	s.scope.Debug("Delete EKS control plane completed successfully")
	return nil
}

// ReconcilePool is the entrypoint for ManagedMachinePool reconciliation.
func (s *NodegroupService) ReconcilePool(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS nodegroup")

	if err := s.reconcileNodegroupIAMRole(ctx); err != nil {
		v1beta1conditions.MarkFalse(
			s.scope.ManagedMachinePool,
			expinfrav1beta1.IAMNodegroupRolesReadyCondition,
			expinfrav1beta1.IAMNodegroupRolesReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1beta1.IAMNodegroupRolesReadyCondition)

	if err := s.reconcileNodegroup(ctx); err != nil {
		v1beta1conditions.MarkFalse(
			s.scope.ManagedMachinePool,
			expinfrav1beta1.EKSNodegroupReadyCondition,
			expinfrav1beta1.EKSNodegroupReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1beta1.EKSNodegroupReadyCondition)

	return nil
}

// ReconcilePoolDelete is the entrypoint for ManagedMachinePool deletion
// reconciliation.
func (s *NodegroupService) ReconcilePoolDelete(ctx context.Context) error {
	s.scope.Debug("Reconciling deletion of EKS nodegroup")

	eksNodegroupName := s.scope.NodegroupName()

	ng, err := s.describeNodegroup(ctx)
	if err != nil {
		if awserrors.IsNotFound(err) {
			s.scope.Trace("EKS nodegroup does not exist")
			return nil
		}
		return errors.Wrap(err, "failed to describe EKS nodegroup")
	}
	if ng == nil {
		return nil
	}

	if err := s.deleteNodegroupAndWait(ctx); err != nil {
		return errors.Wrap(err, "failed to delete nodegroup")
	}

	if err := s.deleteNodegroupIAMRole(ctx); err != nil {
		return errors.Wrap(err, "failed to delete nodegroup IAM role")
	}

	record.Eventf(s.scope.ManagedMachinePool, "SuccessfulDeleteEKSNodegroup", "Deleted EKS nodegroup %s", eksNodegroupName)

	return nil
}
