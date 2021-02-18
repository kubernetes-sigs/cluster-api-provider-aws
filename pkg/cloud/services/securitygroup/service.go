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

package securitygroup

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// Scope is a scope for use with the security group reconciling service
type Scope interface {
	cloud.ClusterScoper

	// Network returns the cluster network object.
	Network() *infrav1.Network

	// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup

	// SecurityGroupOverrides returns the security groups that are overridden in the cluster spec
	SecurityGroupOverrides() map[infrav1.SecurityGroupRole]string

	// VPC returns the cluster VPC.
	VPC() *infrav1.VPCSpec

	// CNIIngressRules returns the CNI spec ingress rules.
	CNIIngressRules() infrav1.CNIIngressRules

	// Bastion returns the bastion details for the cluster.
	Bastion() *infrav1.Bastion
}

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     Scope
	roles     []infrav1.SecurityGroupRole
	EC2Client ec2iface.EC2API
}

// NewService returns a new service given the api clients.
func NewService(sgScope Scope) *Service {
	return &Service{
		scope:     sgScope,
		roles:     defaultRoles,
		EC2Client: scope.NewEC2Client(sgScope, sgScope, sgScope, sgScope.InfraCluster()),
	}
}

// NewServiceWithRoles returns a new service given the api clients with a defined
// set of roles
func NewServiceWithRoles(sgScope Scope, roles []infrav1.SecurityGroupRole) *Service {
	return &Service{
		scope:     sgScope,
		roles:     roles,
		EC2Client: scope.NewEC2Client(sgScope, sgScope, sgScope, sgScope.InfraCluster()),
	}
}
