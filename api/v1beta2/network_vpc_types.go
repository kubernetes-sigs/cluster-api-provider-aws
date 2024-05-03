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

import (
	"fmt"
)

// VPCSpec configures an AWS VPC.
type VPCSpec struct {
	// ID is the vpc-id of the VPC this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	// Defaults to 10.0.0.0/16.
	// Mutually exclusive with IPAMPool.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// IPAMPool defines the IPAMv4 pool to be used for VPC.
	// Mutually exclusive with CidrBlock.
	IPAMPool *IPAMPool `json:"ipamPool,omitempty"`

	// IPv6 contains ipv6 specific settings for the network. Supported only in managed clusters.
	// This field cannot be set on AWSCluster object.
	// +optional
	IPv6 *IPv6 `json:"ipv6,omitempty"`

	// InternetGatewayID is the id of the internet gateway associated with the VPC.
	// +optional
	InternetGatewayID *string `json:"internetGatewayId,omitempty"`

	// CarrierGatewayID is the id of the internet gateway associated with the VPC,
	// for carrier network (Wavelength Zones).
	// +optional
	// +kubebuilder:validation:XValidation:rule="self.startsWith('cagw-')",message="Carrier Gateway ID must start with 'cagw-'"
	CarrierGatewayID *string `json:"carrierGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`

	// AvailabilityZoneUsageLimit specifies the maximum number of availability zones (AZ) that
	// should be used in a region when automatically creating subnets. If a region has more
	// than this number of AZs then this number of AZs will be picked randomly when creating
	// default subnets. Defaults to 3
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=1
	AvailabilityZoneUsageLimit *int `json:"availabilityZoneUsageLimit,omitempty"`

	// AvailabilityZoneSelection specifies how AZs should be selected if there are more AZs
	// in a region than specified by AvailabilityZoneUsageLimit. There are 2 selection schemes:
	// Ordered - selects based on alphabetical order
	// Random - selects AZs randomly in a region
	// Defaults to Ordered
	// +kubebuilder:default=Ordered
	// +kubebuilder:validation:Enum=Ordered;Random
	AvailabilityZoneSelection *AZSelectionScheme `json:"availabilityZoneSelection,omitempty"`

	// EmptyRoutesDefaultVPCSecurityGroup specifies whether the default VPC security group ingress
	// and egress rules should be removed.
	//
	// By default, when creating a VPC, AWS creates a security group called `default` with ingress and egress
	// rules that allow traffic from anywhere. The group could be used as a potential surface attack and
	// it's generally suggested that the group rules are removed or modified appropriately.
	//
	// NOTE: This only applies when the VPC is managed by the Cluster API AWS controller.
	//
	// +optional
	EmptyRoutesDefaultVPCSecurityGroup bool `json:"emptyRoutesDefaultVPCSecurityGroup,omitempty"`

	// PrivateDNSHostnameTypeOnLaunch is the type of hostname to assign to instances in the subnet at launch.
	// For IPv4-only and dual-stack (IPv4 and IPv6) subnets, an instance DNS name can be based on the instance IPv4 address (ip-name)
	// or the instance ID (resource-name). For IPv6 only subnets, an instance DNS name must be based on the instance ID (resource-name).
	// +optional
	// +kubebuilder:validation:Enum:=ip-name;resource-name
	PrivateDNSHostnameTypeOnLaunch *string `json:"privateDnsHostnameTypeOnLaunch,omitempty"`
}

// String returns a string representation of the VPC.
func (v *VPCSpec) String() string {
	return fmt.Sprintf("id=%s", v.ID)
}

// IsUnmanaged returns true if the VPC is unmanaged.
func (v *VPCSpec) IsUnmanaged(clusterName string) bool {
	return v.ID != "" && !v.Tags.HasOwned(clusterName)
}

// IsManaged returns true if VPC is managed.
func (v *VPCSpec) IsManaged(clusterName string) bool {
	return !v.IsUnmanaged(clusterName)
}

// IsIPv6Enabled returns true if the IPv6 block is defined on the network spec.
func (v *VPCSpec) IsIPv6Enabled() bool {
	return v.IPv6 != nil
}
