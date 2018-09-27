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

import "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"

func (s *Service) ReconcileNetwork(network *v1alpha1.Network) (err error) {

	// VPC.
	if err := s.reconcileVPC(&network.VPC); err != nil {
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(network.Subnets, &network.VPC); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(network); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(network.Subnets, &network.VPC); err != nil {
		return err
	}

	// API server loadbalancer.
	if err := s.reconcileAPIServerLoadbalancer(network.Loadbalancer); err != nil {
		return err
	}

	return nil
}
