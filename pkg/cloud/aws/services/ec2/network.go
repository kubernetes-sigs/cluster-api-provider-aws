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

package ec2

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.V(2).Info("Reconciling network for cluster", "cluster-name", s.scope.Cluster.Name, "cluster-namespace", s.scope.Cluster.Namespace)

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		return err
	}

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		return err
	}

	// Security groups.
	if err := s.reconcileSecurityGroups(); err != nil {
		return err
	}

	s.scope.V(2).Info("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork() (err error) {
	s.scope.V(2).Info("Deleting network")

	// Security groups.
	if err := s.deleteSecurityGroups(); err != nil {
		return err
	}

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
