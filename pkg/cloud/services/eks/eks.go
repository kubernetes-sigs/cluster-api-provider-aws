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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
)

// ReconcileControlPlane reconciles a EKS control plane.
func (s *Service) ReconcileControlPlane(ctx context.Context) error {
	s.scope.Debug("Reconciling EKS control plane", "cluster", klog.KRef(s.scope.Cluster.Namespace, s.scope.Cluster.Name))

	// Control Plane IAM Role
	if err := s.reconcileControlPlaneIAMRole(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.IAMControlPlaneRolesReadyCondition, ekscontrolplanev1.IAMControlPlaneRolesReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		conditions.Set(s.scope.ControlPlane, metav1.Condition{
			Type:    ekscontrolplanev1.AWSManagedControlPlaneIAMControlPlaneRolesReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  ekscontrolplanev1.AWSManagedControlPlaneIAMControlPlaneRolesReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.IAMControlPlaneRolesReadyCondition)
	conditions.Set(s.scope.ControlPlane, metav1.Condition{
		Type:   ekscontrolplanev1.AWSManagedControlPlaneIAMControlPlaneRolesReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// EKS Cluster
	if err := s.reconcileCluster(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition, ekscontrolplanev1.EKSControlPlaneReconciliationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		conditions.Set(s.scope.ControlPlane, metav1.Condition{
			Type:    ekscontrolplanev1.AWSManagedControlPlaneEKSControlPlaneReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  ekscontrolplanev1.AWSManagedControlPlaneEKSControlPlaneReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition)
	conditions.Set(s.scope.ControlPlane, metav1.Condition{
		Type:   ekscontrolplanev1.AWSManagedControlPlaneEKSControlPlaneReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// EKS Addons
	if err := s.reconcileAddons(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSAddonsConfiguredCondition, ekscontrolplanev1.EKSAddonsConfiguredFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		conditions.Set(s.scope.ControlPlane, metav1.Condition{
			Type:    ekscontrolplanev1.AWSManagedControlPlaneEKSAddonsConfiguredV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  ekscontrolplanev1.AWSManagedControlPlaneEKSAddonsConfiguredFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrap(err, "failed reconciling eks addons")
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSAddonsConfiguredCondition)
	conditions.Set(s.scope.ControlPlane, metav1.Condition{
		Type:   ekscontrolplanev1.AWSManagedControlPlaneEKSAddonsConfiguredV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// EKS Identity Provider
	if err := s.reconcileIdentityProvider(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSIdentityProviderConfiguredCondition, ekscontrolplanev1.EKSIdentityProviderConfiguredFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		conditions.Set(s.scope.ControlPlane, metav1.Condition{
			Type:    ekscontrolplanev1.AWSManagedControlPlaneEKSIdentityProviderConfiguredV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  ekscontrolplanev1.AWSManagedControlPlaneEKSIdentityProviderConfiguredFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrap(err, "failed reconciling eks identity provider")
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSIdentityProviderConfiguredCondition)
	conditions.Set(s.scope.ControlPlane, metav1.Condition{
		Type:   ekscontrolplanev1.AWSManagedControlPlaneEKSIdentityProviderConfiguredV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	if err := s.reconcilePodIdentityAssociations(ctx); err != nil {
		v1beta1conditions.MarkFalse(s.scope.ControlPlane, ekscontrolplanev1.EKSPodIdentityAssociationConfiguredCondition, ekscontrolplanev1.EKSPodIdentityAssociationFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		conditions.Set(s.scope.ControlPlane, metav1.Condition{
			Type:    ekscontrolplanev1.AWSManagedControlPlaneEKSPodIdentityAssociationConfiguredV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  ekscontrolplanev1.AWSManagedControlPlaneEKSPodIdentityAssociationConfigurationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrap(err, "failed reconciling eks pod identity associations")
	}
	v1beta1conditions.MarkTrue(s.scope.ControlPlane, ekscontrolplanev1.EKSPodIdentityAssociationConfiguredCondition)
	conditions.Set(s.scope.ControlPlane, metav1.Condition{
		Type:   ekscontrolplanev1.AWSManagedControlPlaneEKSPodIdentityAssociationConfiguredV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

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
			expinfrav1.IAMNodegroupRolesReadyCondition,
			expinfrav1.IAMNodegroupRolesReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		conditions.Set(s.scope.ManagedMachinePool, metav1.Condition{
			Type:    expinfrav1.AWSManagedMachinePoolIAMNodegroupRolesReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  expinfrav1.AWSManagedMachinePoolIAMNodegroupRolesReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1.IAMNodegroupRolesReadyCondition)
	conditions.Set(s.scope.ManagedMachinePool, metav1.Condition{
		Type:   expinfrav1.AWSManagedMachinePoolIAMNodegroupRolesReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	if err := s.reconcileNodegroup(ctx); err != nil {
		v1beta1conditions.MarkFalse(
			s.scope.ManagedMachinePool,
			expinfrav1.EKSNodegroupReadyCondition,
			expinfrav1.EKSNodegroupReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"%s",
			err.Error(),
		)
		conditions.Set(s.scope.ManagedMachinePool, metav1.Condition{
			Type:    expinfrav1.AWSManagedMachinePoolEKSNodegroupReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  expinfrav1.AWSManagedMachinePoolEKSNodegroupReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.ManagedMachinePool, expinfrav1.EKSNodegroupReadyCondition)
	conditions.Set(s.scope.ManagedMachinePool, metav1.Condition{
		Type:   expinfrav1.AWSManagedMachinePoolEKSNodegroupReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

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
