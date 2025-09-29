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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.Debug("Reconciling network for cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VpcReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.VpcReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Secondary CIDRs
	if err := s.associateSecondaryCidrs(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.SecondaryCidrReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.SecondaryCidrsReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.SubnetsReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.SubnetsReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.InternetGatewayFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.InternetGatewayReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Carrier Gateway.
	if err := s.reconcileCarrierGateway(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.CarrierGatewayFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.CarrierGatewayReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Egress Only Internet Gateways.
	if err := s.reconcileEgressOnlyInternetGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.EgressOnlyInternetGatewayFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.EgressOnlyInternetGatewayReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.NatGatewaysReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.NatGatewaysReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.RouteTableReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.RouteTablesReadyCondition,
		Status: metav1.ConditionTrue,
	})

	// VPC Endpoints.
	if err := s.reconcileVPCEndpoints(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VpcEndpointsReconciliationFailedReason,
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:   infrav1.VpcEndpointsReadyCondition,
		Status: metav1.ConditionTrue,
	})

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
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPCEndpoints(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// Routing tables.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteRouteTables(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// NAT Gateways.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteNatGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// Carrier Gateway.
	if s.scope.VPC().CarrierGatewayID != nil {
		if err := s.deleteCarrierGateway(); err != nil {
			conditions.Set(s.scope.InfraCluster(), metav1.Condition{
				Type:    infrav1.VpcEndpointsReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  "DeletingFailed",
				Message: fmt.Sprintf("%s", err),
			})
			return err
		}
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  clusterv1.DeletedV1Beta1Reason,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Egress Only Internet Gateways.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteEgressOnlyInternetGateways(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// Subnets.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteSubnets(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	// Secondary CIDR.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.disassociateSecondaryCidrs(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DisassociateFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}

	// VPC.
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletingReason,
		Message: fmt.Sprintf("%s", err),
	})
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPC(); err != nil {
		conditions.Set(s.scope.InfraCluster(), metav1.Condition{
			Type:    infrav1.VpcEndpointsReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  "DeletingFailed",
			Message: fmt.Sprintf("%s", err),
		})
		return err
	}
	conditions.Set(s.scope.InfraCluster(), metav1.Condition{
		Type:    infrav1.VpcEndpointsReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  clusterv1.DeletedV1Beta1Reason,
		Message: fmt.Sprintf("%s", err),
	})

	s.scope.Debug("Delete network completed successfully")
	return nil
}
