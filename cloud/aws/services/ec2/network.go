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
	"context"
	"github.com/golang/glog"
	"go.opencensus.io/trace"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/instrumentation"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork(ctx context.Context, clusterName string, network *v1alpha1.Network) (err error) {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "ReconcileNetwork"),
	)
	defer span.End()

	glog.V(2).Info("Reconciling network")

	// VPC.
	if err := s.reconcileVPC(ctx, clusterName, &network.VPC); err != nil {
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(ctx, clusterName, network); err != nil {
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(ctx, clusterName, network); err != nil {
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(ctx, clusterName, network.Subnets, &network.VPC); err != nil {
		return err
	}

	// Routing tables.
	if err := s.reconcileRouteTables(ctx, clusterName, network); err != nil {
		return err
	}

	// Security groups.
	if err := s.reconcileSecurityGroups(ctx, clusterName, network); err != nil {
		return err
	}

	glog.V(2).Info("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork(ctx context.Context, clusterName string, network *v1alpha1.Network) (err error) {
	ctx, span := trace.StartSpan(
		ctx, instrumentation.MethodName("services", "ec2", "DeleteNetwork"),
	)
	defer span.End()

	glog.V(2).Info("Deleting network")

	// Security groups.
	if err := s.deleteSecurityGroups(ctx, clusterName, network); err != nil {
		return err
	}

	// Routing tables.
	if err := s.deleteRouteTables(ctx, clusterName, network); err != nil {
		return err
	}

	if err := s.deleteNatGateways(ctx, clusterName, network.Subnets, &network.VPC); err != nil {
		return err
	}

	if err := s.releaseAddresses(ctx, clusterName); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(ctx, clusterName, network); err != nil {
		return err
	}

	// Subnets.
	if err := s.deleteSubnets(ctx, clusterName, network); err != nil {
		return err
	}

	// VPC.
	if err := s.deleteVPC(ctx, &network.VPC); err != nil {
		return err
	}

	glog.V(2).Info("Delete network completed successfully")
	return nil
}
