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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	infrautilconditions "sigs.k8s.io/cluster-api-provider-aws/v2/util/conditions"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

const deletingFailedReason = "DeletingFailed"

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.Debug("Reconciling network for cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterVpcReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterVpcReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.VpcReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Secondary CIDRs
	if err := s.associateSecondaryCidrs(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, infrav1.SecondaryCidrReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterSecondaryCidrsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterSecondaryCidrReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterSecondaryCidrsReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, infrav1.SubnetsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterSubnetsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterSubnetsReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterSubnetsReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, infrav1.InternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterInternetGatewayFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Carrier Gateway.
	if err := s.reconcileCarrierGateway(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, infrav1.CarrierGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterCarrierGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterCarrierGatewayFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterCarrierGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Egress Only Internet Gateways.
	if err := s.reconcileEgressOnlyInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, infrav1.EgressOnlyInternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterEgressOnlyInternetGatewayFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, infrav1.NatGatewaysReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterNatGatewaysReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, infrav1.RouteTableReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterRouteTablesReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterRouteTableReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterRouteTablesReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	// VPC Endpoints.
	if err := s.reconcileVPCEndpoints(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, infrav1.VpcEndpointsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.AWSClusterVpcEndpointsReconciliationFailedV1Beta2Reason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkTrue(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition)
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ReadyReason,
	})

	s.scope.Debug("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork(ctx context.Context) (err error) {
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
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPCEndpoints(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// Routing tables.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterRouteTablesReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteRouteTables(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterRouteTablesReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterRouteTablesReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// NAT Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteNatGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterNatGatewaysReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// Carrier Gateway.
	if s.scope.VPC().CarrierGatewayID != nil {
		if err := s.deleteCarrierGateway(); err != nil {
			v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
			infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
				Type:    infrav1.AWSClusterCarrierGatewayReadyV1Beta2Condition,
				Status:  metav1.ConditionFalse,
				Reason:  deletingFailedReason,
				Message: err.Error(),
			})
			return err
		}
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:   infrav1.AWSClusterCarrierGatewayReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: clusterv1beta1.DeletedReason,
		})
	}

	// Egress Only Internet Gateways.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteEgressOnlyInternetGateways(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// Orphaned ENIs — clean up detached network interfaces that reference
	// cluster security groups before attempting subnet deletion.
	if err := s.deleteOrphanedENIs(ctx); err != nil {
		s.scope.Debug("Non-fatal: failed to delete orphaned ENIs", "error", err)
	}

	// Subnets.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterSubnetsReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteSubnets(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterSubnetsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterSubnetsReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	// Secondary CIDR.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterSecondaryCidrsReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.disassociateSecondaryCidrs(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, "DisassociateFailed", clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterSecondaryCidrsReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  "DisassociateFailed",
			Message: err.Error(),
		})
		return err
	}

	// VPC.
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1beta1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1.DeletingReason,
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPC(); err != nil {
		v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, deletingFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.AWSClusterVpcReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  deletingFailedReason,
			Message: err.Error(),
		})
		return err
	}
	v1beta1conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1beta1.DeletedReason, clusterv1beta1.ConditionSeverityInfo, "")
	infrautilconditions.SetV1Beta2(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.AWSClusterVpcReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: clusterv1beta1.DeletedReason,
	})

	s.scope.Debug("Delete network completed successfully")
	return nil
}
