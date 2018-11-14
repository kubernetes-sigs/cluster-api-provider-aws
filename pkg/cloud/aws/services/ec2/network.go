// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ec2

import (
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork(clusterName string, network *v1alpha1.Network) (err error) {
	klog.V(2).Info("Reconciling network")

	// VPC.
	if err := s.reconcileVPC(clusterName, &network.VPC); err != nil {
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(clusterName, network); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(clusterName, network); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(clusterName, network.Subnets, &network.VPC); err != nil {
		return err
	}

	// Routing tables.
	if err := s.reconcileRouteTables(clusterName, network); err != nil {
		return err
	}

	// Security groups.
	if err := s.reconcileSecurityGroups(clusterName, network); err != nil {
		return err
	}

	klog.V(2).Info("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork(clusterName string, network *v1alpha1.Network) (err error) {
	klog.V(2).Info("Deleting network")

	// Security groups.
	if err := s.deleteSecurityGroups(clusterName, network); err != nil {
		return err
	}

	// Routing tables.
	if err := s.deleteRouteTables(clusterName, network); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.deleteNatGateways(clusterName, network.Subnets, &network.VPC); err != nil {
		return err
	}

	// EIPs.
	if err := s.releaseAddresses(clusterName); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.deleteInternetGateways(clusterName, network); err != nil {
		return err
	}

	// Subnets.
	if err := s.deleteSubnets(clusterName, network); err != nil {
		return err
	}

	// VPC.
	if err := s.deleteVPC(&network.VPC); err != nil {
		return err
	}

	klog.V(2).Info("Delete network completed successfully")
	return nil
}
