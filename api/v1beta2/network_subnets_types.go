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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/utils/ptr"
)

const (
	// ZoneTypeAvailabilityZone defines the regular AWS zones in the Region.
	ZoneTypeAvailabilityZone ZoneType = "availability-zone"
	// ZoneTypeLocalZone defines the AWS zone type in Local Zone infrastructure.
	ZoneTypeLocalZone ZoneType = "local-zone"
	// ZoneTypeWavelengthZone defines the AWS zone type in Wavelength infrastructure.
	ZoneTypeWavelengthZone ZoneType = "wavelength-zone"
)

// SubnetSpec configures an AWS Subnet.
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	// If you're bringing your subnet, set the AWS subnet-id here, it must start with `subnet-`.
	//
	// When the VPC is managed by CAPA, and you'd like the provider to create a subnet for you,
	// the id can be set to any placeholder value that does not start with `subnet-`;
	// upon creation, the subnet AWS identifier will be populated in the `ResourceID` field and
	// the `id` field is going to be used as the subnet name. If you specify a tag
	// called `Name`, it takes precedence.
	ID string `json:"id"`

	// ResourceID is the subnet identifier from AWS, READ ONLY.
	// This field is populated when the provider manages the subnet.
	// +optional
	ResourceID string `json:"resourceID,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// IPv6CidrBlock is the IPv6 CIDR block to be used when the provider creates a managed VPC.
	// A subnet can have an IPv4 and an IPv6 address.
	// IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.
	// +optional
	IPv6CidrBlock string `json:"ipv6CidrBlock,omitempty"`

	// AvailabilityZone defines the availability zone to use for this subnet in the cluster's region.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// IsPublic defines the subnet as a public subnet. A subnet is public when it is associated with a route table that has a route to an internet gateway.
	// +optional
	IsPublic bool `json:"isPublic"`

	// IsIPv6 defines the subnet as an IPv6 subnet. A subnet is IPv6 when it is associated with a VPC that has IPv6 enabled.
	// IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.
	// +optional
	IsIPv6 bool `json:"isIpv6,omitempty"`

	// RouteTableID is the routing table id associated with the subnet.
	// +optional
	RouteTableID *string `json:"routeTableId,omitempty"`

	// NatGatewayID is the NAT gateway id associated with the subnet.
	// Ignored unless the subnet is managed by the provider, in which case this is set on the public subnet where the NAT gateway resides. It is then used to determine routes for private subnets in the same AZ as the public subnet.
	// +optional
	NatGatewayID *string `json:"natGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`

	// ZoneType defines the type of the zone where the subnet is created.
	//
	// The valid values are availability-zone, local-zone, and wavelength-zone.
	//
	// Subnet with zone type availability-zone (regular) is always selected to create cluster
	// resources, like Load Balancers, NAT Gateways, Contol Plane nodes, etc.
	//
	// Subnet with zone type local-zone or wavelength-zone is not eligible to automatically create
	// regular cluster resources.
	//
	// The public subnet in availability-zone or local-zone is associated with regular public
	// route table with default route entry to a Internet Gateway.
	//
	// The public subnet in wavelength-zone is associated with a carrier public
	// route table with default route entry to a Carrier Gateway.
	//
	// The private subnet in the availability-zone is associated with a private route table with
	// the default route entry to a NAT Gateway created in that zone.
	//
	// The private subnet in the local-zone or wavelength-zone is associated with a private route table with
	// the default route entry re-using the NAT Gateway in the Region (preferred from the
	// parent zone, the zone type availability-zone in the region, or first table available).
	//
	// +kubebuilder:validation:Enum=availability-zone;local-zone;wavelength-zone
	// +optional
	ZoneType *ZoneType `json:"zoneType,omitempty"`

	// ParentZoneName is the zone name where the current subnet's zone is tied when
	// the zone is a Local Zone.
	//
	// The subnets in Local Zone or Wavelength Zone locations consume the ParentZoneName
	// to select the correct private route table to egress traffic to the internet.
	//
	// +optional
	ParentZoneName *string `json:"parentZoneName,omitempty"`
}

// GetResourceID returns the identifier for this subnet,
// if the subnet was not created or reconciled, it returns the subnet ID.
func (s *SubnetSpec) GetResourceID() string {
	if s.ResourceID != "" {
		return s.ResourceID
	}
	return s.ID
}

// String returns a string representation of the subnet.
func (s *SubnetSpec) String() string {
	return fmt.Sprintf("id=%s/az=%s/public=%v", s.GetResourceID(), s.AvailabilityZone, s.IsPublic)
}

// IsEdge returns the true when the subnet is created in the edge zone,
// Local Zones.
func (s *SubnetSpec) IsEdge() bool {
	if s.ZoneType == nil {
		return false
	}
	if s.ZoneType.Equal(ZoneTypeLocalZone) {
		return true
	}
	if s.ZoneType.Equal(ZoneTypeWavelengthZone) {
		return true
	}
	return false
}

// IsEdgeWavelength returns true only when the subnet is created in Wavelength Zone.
func (s *SubnetSpec) IsEdgeWavelength() bool {
	if s.ZoneType == nil {
		return false
	}
	if *s.ZoneType == ZoneTypeWavelengthZone {
		return true
	}
	return false
}

// SetZoneInfo updates the subnets with zone information.
func (s *SubnetSpec) SetZoneInfo(zones []*ec2.AvailabilityZone) error {
	zoneInfo := func(zoneName string) *ec2.AvailabilityZone {
		for _, zone := range zones {
			if aws.StringValue(zone.ZoneName) == zoneName {
				return zone
			}
		}
		return nil
	}

	zone := zoneInfo(s.AvailabilityZone)
	if zone == nil {
		if len(s.AvailabilityZone) > 0 {
			return fmt.Errorf("unable to update zone information for subnet '%v' and zone '%v'", s.ID, s.AvailabilityZone)
		}
		return fmt.Errorf("unable to update zone information for subnet '%v'", s.ID)
	}
	if zone.ZoneType != nil {
		s.ZoneType = ptr.To(ZoneType(*zone.ZoneType))
	}
	if zone.ParentZoneName != nil {
		s.ParentZoneName = zone.ParentZoneName
	}
	return nil
}

// Subnets is a slice of Subnet.
// +listType=map
// +listMapKey=id
type Subnets []SubnetSpec

// ToMap returns a map from id to subnet.
func (s Subnets) ToMap() map[string]*SubnetSpec {
	res := make(map[string]*SubnetSpec)
	for i := range s {
		x := s[i]
		res[x.GetResourceID()] = &x
	}
	return res
}

// IDs returns a slice of the subnet ids.
func (s Subnets) IDs() []string {
	res := []string{}
	for _, subnet := range s {
		// Prevent returning edge zones (Local Zone) to regular Subnet IDs.
		// Edge zones should not deploy control plane nodes, and does not support Nat Gateway and
		// Network Load Balancers. Any resource for the core infrastructure should not consume edge
		// zones.
		if subnet.IsEdge() {
			continue
		}
		res = append(res, subnet.GetResourceID())
	}
	return res
}

// IDsWithEdge returns a slice of the subnet ids.
func (s Subnets) IDsWithEdge() []string {
	res := []string{}
	for _, subnet := range s {
		res = append(res, subnet.GetResourceID())
	}
	return res
}

// FindByID returns a single subnet matching the given id or nil.
//
// The returned pointer can be used to write back into the original slice.
func (s Subnets) FindByID(id string) *SubnetSpec {
	for i := range s {
		x := &(s[i]) // pointer to original structure
		if x.GetResourceID() == id {
			return x
		}
	}
	return nil
}

// FindEqual returns a subnet spec that is equal to the one passed in.
// Two subnets are defined equal to each other if their id is equal
// or if they are in the same vpc and the cidr block is the same.
//
// The returned pointer can be used to write back into the original slice.
func (s Subnets) FindEqual(spec *SubnetSpec) *SubnetSpec {
	for i := range s {
		x := &(s[i]) // pointer to original structure
		if (spec.GetResourceID() != "" && x.GetResourceID() == spec.GetResourceID()) ||
			(spec.CidrBlock == x.CidrBlock) ||
			(spec.IPv6CidrBlock != "" && spec.IPv6CidrBlock == x.IPv6CidrBlock) {
			return x
		}
	}
	return nil
}

// FilterPrivate returns a slice containing all subnets marked as private.
func (s Subnets) FilterPrivate() (res Subnets) {
	for _, x := range s {
		// Subnets in AWS Local Zones or Wavelength should not be used by core infrastructure.
		if x.IsEdge() {
			continue
		}
		if !x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterPublic returns a slice containing all subnets marked as public.
func (s Subnets) FilterPublic() (res Subnets) {
	for _, x := range s {
		// Subnets in AWS Local Zones or Wavelength should not be used by core infrastructure.
		if x.IsEdge() {
			continue
		}
		if x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterByZone returns a slice containing all subnets that live in the availability zone specified.
func (s Subnets) FilterByZone(zone string) (res Subnets) {
	for _, x := range s {
		if x.AvailabilityZone == zone {
			res = append(res, x)
		}
	}
	return
}

// GetUniqueZones returns a slice containing the unique zones of the subnets.
func (s Subnets) GetUniqueZones() []string {
	keys := make(map[string]bool)
	zones := []string{}
	for _, x := range s {
		if _, value := keys[x.AvailabilityZone]; len(x.AvailabilityZone) > 0 && !value {
			keys[x.AvailabilityZone] = true
			zones = append(zones, x.AvailabilityZone)
		}
	}
	return zones
}

// SetZoneInfo updates the subnets with zone information.
func (s Subnets) SetZoneInfo(zones []*ec2.AvailabilityZone) error {
	for i := range s {
		if err := s[i].SetZoneInfo(zones); err != nil {
			return err
		}
	}
	return nil
}

// HasPublicSubnetWavelength returns true when there are subnets in Wavelength zone.
func (s Subnets) HasPublicSubnetWavelength() bool {
	for _, sub := range s {
		if sub.ZoneType == nil {
			return false
		}
		if sub.IsPublic && *sub.ZoneType == ZoneTypeWavelengthZone {
			return true
		}
	}
	return false
}

// ZoneType defines listener AWS Availability Zone type.
type ZoneType string

// String returns the string representation for the zone type.
func (z ZoneType) String() string {
	return string(z)
}

// Equal compares two zone types.
func (z ZoneType) Equal(other ZoneType) bool {
	return z == other
}
