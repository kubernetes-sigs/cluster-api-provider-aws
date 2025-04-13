/*
Copyright The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// RosaNetworkSpec defines the desired state of RosaNetwork
type RosaNetworkSpec struct {
	// The AWS region in which the components of ROSA network infrastruture are to be crated
	// +immutable
	Region string `json:"region"`

	// The number of availability zones to be used for creation of the network infrastructure.
	// You can specify anything between one and four, depending on the chosen AWS region.
	// +kubebuilder:default=1
	// +optional
	// +immutable
	AvailabilityZoneCount int `json:"availabilityZoneCount"`

	// The list of availability zones to be used for creation of the network infrastructure.
	// You can specify anything between one and four valid availability zones from a given region.
	// Should you specify both the availabilityZoneCount and availabilityZones, the list of availability zones takes preference.
	// +optional
	// +immutable
	AvailabilityZones []string `json:"availabilityZones"`

	// CIDR block to be used for the VPC
	// +kubebuilder:validation:Format=cidr
	// +immutable
	CIDRBlock string `json:"cidrBlock"`
}

// RosaNetworkSubnet groups public and private subnet and the availability zone in which the two subnets got created
type RosaNetworkSubnet struct {
	// Availability zone of the subnet pair
	AvailabilityZone string `json:"availabilityZone"`

	// ID of the public subnet
	PublicSubnet     string `json:"publicSubnet"`

	// ID of the private subnet
	PrivateSubnet    string `json:"privateSubnet"`
}

// CFResource groups information pertaining to a resource created as a part of a cloudformation stack
type CFResource struct {
	// Name of the created resource: NATGateway1, VPC, SecurityGroup, ...
	Resource string `json:"resource"`

	// Identified of the created resource. Will be filled in once the resource is created & ready
	ID string `json:"ID"`

	// Status of the resource: CREATE_IN_PROGRESS, CREATE_COMPLETE, ...
	Status string `json:"status"`

	// Message pertaining to the status of the resource
	Reason string `json:"reason"`
}

// RosaNetworkStatus defines the observed state of RosaNetwork
type RosaNetworkStatus struct {
	// Array of created private, public subnets and availability zones, grouped by availability zones
	Subnets []RosaNetworkSubnet `json:"subnets"`

	// Resources created in the cloudformation stack
	Resources []CFResource `json:"resources"`

	// Conditions specifies the conditions for RosaNetwork
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosanetworks,shortName=rosanet,scope=Namespaced,categories=cluster-api

// RosaNetwork is the Schema for the rosanetworks API
type RosaNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RosaNetworkSpec   `json:"spec,omitempty"`
	Status RosaNetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RosaNetworkList contains a list of RosaNetwork
type RosaNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RosaNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RosaNetwork{}, &RosaNetworkList{})
}
