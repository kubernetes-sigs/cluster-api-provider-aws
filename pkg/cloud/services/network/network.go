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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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

	// Search for a previously created and tagged VPC
	vpc, err := s.describeVPC()
	if err != nil {
		if awserrors.IsNotFound(err) {
			// If the VPC does not exist, nothing to do
			return nil
		}
		return err
	}
	vpc.DeepCopyInto(s.scope.VPC())

	// Routing tables.
	if err := s.deleteRouteTables(); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.deleteNatGateways(); err != nil {
		return err
	}

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.deleteInternetGateways(); err != nil {
		return err
	}

	// Subnets.
	if err := s.deleteSubnets(); err != nil {
		return err
	}

	// VPC.
	if err := s.deleteVPC(); err != nil {
		return err
	}

	s.scope.V(2).Info("Delete network completed successfully")
	return nil
}
