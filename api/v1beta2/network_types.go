/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

// NetworkSpec encapsulates all things related to AWS network.
type NetworkSpec struct {
	// VPC configuration.
	// +optional
	VPC VPCSpec `json:"vpc,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`

	// CNI configuration
	// +optional
	CNI *CNISpec `json:"cni,omitempty"`

	// SecurityGroupOverrides is an optional set of security groups to use for cluster instances
	// This is optional - if not provided new security groups will be created for the cluster
	// +optional
	SecurityGroupOverrides map[SecurityGroupRole]string `json:"securityGroupOverrides,omitempty"`

	// AdditionalControlPlaneIngressRules is an optional set of ingress rules to add to the control plane
	// +optional
	AdditionalControlPlaneIngressRules []IngressRule `json:"additionalControlPlaneIngressRules,omitempty"`
}

// NetworkStatus encapsulates AWS networking resources.
type NetworkStatus struct {
	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]SecurityGroup `json:"securityGroups,omitempty"`

	// APIServerELB is the Kubernetes api server load balancer.
	APIServerELB LoadBalancer `json:"apiServerElb,omitempty"`

	// SecondaryAPIServerELB is the secondary Kubernetes api server load balancer.
	SecondaryAPIServerELB LoadBalancer `json:"secondaryAPIServerELB,omitempty"`

	// NatGatewaysIPs contains the public IPs of the NAT Gateways
	NatGatewaysIPs []string `json:"natGatewaysIPs,omitempty"`
}

// IPv6 contains ipv6 specific settings for the network.
type IPv6 struct {
	// CidrBlock is the CIDR block provided by Amazon when VPC has enabled IPv6.
	// Mutually exclusive with IPAMPool.
	// +optional
	CidrBlock string `json:"cidrBlock,omitempty"`

	// PoolID is the IP pool which must be defined in case of BYO IP is defined.
	// Must be specified if CidrBlock is set.
	// Mutually exclusive with IPAMPool.
	// +optional
	PoolID string `json:"poolId,omitempty"`

	// EgressOnlyInternetGatewayID is the id of the egress only internet gateway associated with an IPv6 enabled VPC.
	// +optional
	EgressOnlyInternetGatewayID *string `json:"egressOnlyInternetGatewayId,omitempty"`

	// IPAMPool defines the IPAMv6 pool to be used for VPC.
	// Mutually exclusive with CidrBlock.
	// +optional
	IPAMPool *IPAMPool `json:"ipamPool,omitempty"`
}

// IPAMPool defines the IPAM pool to be used for VPC.
type IPAMPool struct {
	// ID is the ID of the IPAM pool this provider should use to create VPC.
	ID string `json:"id,omitempty"`
	// Name is the name of the IPAM pool this provider should use to create VPC.
	Name string `json:"name,omitempty"`
	// The netmask length of the IPv4 CIDR you want to allocate to VPC from
	// an Amazon VPC IP Address Manager (IPAM) pool.
	// Defaults to /16 for IPv4 if not specified.
	NetmaskLength int64 `json:"netmaskLength,omitempty"`
}

// RouteTable defines an AWS routing table.
type RouteTable struct {
	ID string `json:"id"`
}
