/*
Copyright 2018 The Kubernetes Authors.

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

package network

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.V(2).Info("Reconciling network for cluster", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.VpcReadyCondition)

	// Secondary CIDR
	if err := s.associateSecondaryCidr(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, infrav1.SecondaryCidrReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, infrav1.SubnetsReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, infrav1.InternetGatewayFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, infrav1.NatGatewaysReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, infrav1.RouteTableReconciliationFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	s.scope.V(2).Info("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork() (err error) {
	s.scope.V(2).Info("Deleting network")

	vpc := &infrav1.VPCSpec{}
	// Get VPC used for the cluster
	if s.scope.VPC().ID != "" {
		var err error
		vpc, err = s.describeVPCByID()
		if err != nil {
			if awserrors.IsNotFound(err) {
				// If the VPC does not exist, nothing to do
				return nil
			}
			return err
		}
	} else {
		s.scope.Error(err, "non-fatal: VPC ID is missing, ")
	}

	vpc.DeepCopyInto(s.scope.VPC())

	// Secondary CIDR
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.disassociateSecondaryCidr(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, "DisassociateFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}

	// Routing tables.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteRouteTables(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// NAT Gateways.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteNatGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// Subnets.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteSubnets(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// VPC.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPC(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	s.scope.V(2).Info("Delete network completed successfully")
	return nil
}
