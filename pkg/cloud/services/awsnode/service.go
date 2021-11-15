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

package awsnode

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
)

// Scope is a scope for use with the awsnode reconciling service.
type Scope interface {
	cloud.ClusterScoper

	// RemoteClient returns the Kubernetes client for connecting to the workload cluster.
	RemoteClient() (client.Client, error)
	// Subnets returns the cluster subnets.
	Subnets() infrav1.Subnets
	// SecondaryCidrBlock returns the optional secondary CIDR block to use for pod IPs
	SecondaryCidrBlock() *string
	// SecurityGroups returns the control plane security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	// DisableVPCCNI returns whether the AWS VPC CNI should be disabled
	DisableVPCCNI() bool
}

// Service defines the spec for a service.
type Service struct {
	scope Scope
}

// NewService will create a new service.
func NewService(awsnodeScope Scope) *Service {
	return &Service{
		scope: awsnodeScope,
	}
}
