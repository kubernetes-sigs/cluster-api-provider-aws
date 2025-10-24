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
	"k8s.io/klog/v2"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	infrautilconditions "sigs.k8s.io/cluster-api-provider-aws/v2/util/conditions"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.Debug("Reconciling network for cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcReadyCondition, infrav1beta1.VpcReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.VpcReadyCondition)

	// Secondary CIDRs
	if err := s.associateSecondaryCidrs(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SecondaryCidrsReadyCondition, infrav1beta1.SecondaryCidrReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.SecondaryCidrsReadyCondition)

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SubnetsReadyCondition, infrav1beta1.SubnetsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.SubnetsReadyCondition)

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.InternetGatewayReadyCondition, infrav1beta1.InternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.InternetGatewayReadyCondition)

	// Carrier Gateway.
	if err := s.reconcileCarrierGateway(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.CarrierGatewayReadyCondition, infrav1beta1.CarrierGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.CarrierGatewayReadyCondition)

	// Egress Only Internet Gateways.
	if err := s.reconcileEgressOnlyInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.EgressOnlyInternetGatewayReadyCondition, infrav1beta1.EgressOnlyInternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.EgressOnlyInternetGatewayReadyCondition)

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.NatGatewaysReadyCondition, infrav1beta1.NatGatewaysReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.NatGatewaysReadyCondition)

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.RouteTablesReadyCondition, infrav1beta1.RouteTableReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.RouteTablesReadyCondition)

	// VPC Endpoints.
	if err := s.reconcileVPCEndpoints(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcEndpointsReadyCondition, infrav1beta1.VpcEndpointsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1beta1.VpcEndpointsReadyCondition)

	s.scope.Debug("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork() (err error) {
	s.scope.Debug("Deleting network")

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

	// VPC Endpoints.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcEndpointsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPCEndpoints(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcEndpointsReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcEndpointsReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// Routing tables.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.RouteTablesReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteRouteTables(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.RouteTablesReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.RouteTablesReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// NAT Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.NatGatewaysReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteNatGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.NatGatewaysReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.NatGatewaysReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.InternetGatewayReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.InternetGatewayReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.InternetGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// Carrier Gateway.
	if s.scope.VPC().CarrierGatewayID != nil {
		if err := s.deleteCarrierGateway(); err != nil {
			v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.CarrierGatewayReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
			return err
		}
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.CarrierGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	}

	// Egress Only Internet Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.EgressOnlyInternetGatewayReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteEgressOnlyInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.EgressOnlyInternetGatewayReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.EgressOnlyInternetGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// Subnets.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SubnetsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteSubnets(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SubnetsReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SubnetsReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	// Secondary CIDR.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SecondaryCidrsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.disassociateSecondaryCidrs(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.SecondaryCidrsReadyCondition, "DisassociateFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}

	// VPC.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPC(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcReadyCondition, "DeletingFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1beta1.VpcReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")

	s.scope.Debug("Delete network completed successfully")
	return nil
}
