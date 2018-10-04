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

package v1alpha1

import (
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSMachineProviderConfig is the type that will be embedded in a Machine.Spec.ProviderConfig field
// for an AWS instance. It is used by the AWS machine actuator to create a single machine instance,
// using the RunInstances call (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RunInstances.html)
// Required parameters such as region that are not specified by this configuration, will be defaulted
// by the actuator.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// AMI is the reference to the AMI from which to create the machine instance.
	AMI AWSResourceReference `json:"ami"`

	// InstanceType is the type of instance to create. Example: m4.xlarge
	InstanceType string `json:"instanceType"`

	// AdditionalTags is the set of tags to add to an instance, in addition to the ones
	// added by default by the actuator. These tags are additive. The actuator will ensure
	// these tags are present, but will not remove any other tags that may exist on the
	// instance.
	// +optional
	AdditionalTags map[string]string `json:"additionalTags,omitempty"`

	// IAMInstanceProfile is a reference to an IAM role to assign to the instance
	// +optional
	IAMInstanceProfile *AWSResourceReference `json:"iamInstanceProfile,omitempty"`

	// PublicIP specifies whether the instance should get a public IP.
	// Precedence for this setting is as follows:
	// 1. This field if set
	// 2. Cluster/flavor setting
	// 3. Subnet default
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`

	// AdditionalSecurityGroups is an array of references to security groups that should be applied to the
	// instance. These security groups would be set in addition to any security groups defined
	// at the cluster level or in the actuator.
	// +optional
	AdditionalSecurityGroups []AWSResourceReference `json:"additionalSecurityGroups,omitempty"`

	// Subnet is a reference to the subnet to use for this instance. If not specified,
	// the cluster subnet will be used.
	// +optional
	Subnet *AWSResourceReference `json:"subnet,omitempty"`

	// KeyName is the name of the SSH key to install on the instance.
	// +optional
	KeyName string `json:"keyName"`
}

// AWSResourceReference is a reference to a specific AWS resource by ID, ARN, or filters.
// Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
// a validation error.
type AWSResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// ARN of resource
	// +optional
	ARN *string `json:"arn,omitempty"`

	// Filters is a set of key/value pairs used to identify a resource
	// They are applied according to the rules defined by the AWS API:
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html
	// +optional
	Filters []Filter `json:"filters"`
}

// Filter is a filter used to identify an AWS resource
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	// The AWS Region the cluster lives in.
	Region string `json:"region"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host.
	SSHKeyName string `json:"sshKeyName,omitempty"`
}

// AWSMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It containsk AWS-specific status information.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceID is the instance ID of the machine created in AWS
	// +optional
	InstanceID *string `json:"instanceID,omitempty"`

	// InstanceState is the state of the AWS instance for this machine
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	// +optional
	Conditions []AWSMachineProviderCondition `json:"conditions,omitempty"`
}

// AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type
type AWSMachineProviderConditionType string

// Valid conditions for an AWS machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated AWSMachineProviderConditionType = "MachineCreated"
)

// AWSMachineProviderCondition is a condition in a AWSMachineProviderStatus
type AWSMachineProviderCondition struct {
	// Type is the type of the condition.
	Type AWSMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	Region  string   `json:"region"`
	Network Network  `json:"network"`
	Bastion Instance `json:"bastion"`
}

// Network encapsulates AWS networking resources.
type Network struct {

	// VPC defines the cluster vpc.
	VPC VPC `json:"vpc"`

	// InternetGatewayID is the id of the internet gateway associated with the VPC.
	InternetGatewayID *string `json:"internetGatewayId"`

	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]*SecurityGroup `json:"securityGroups"`

	// Subnets includes all the subnets defined inside the VPC.
	Subnets Subnets `json:"subnets"`
}

// VPC defines an AWS vpc.
type VPC struct {
	ID string `json:"id"`

	CidrBlock string `json:"cidrBlock"`
}

// String returns a string representation of the VPC.
func (v *VPC) String() string {
	return fmt.Sprintf("id=%s", v.ID)
}

// Subnet defines an AWS subnet attached to a VPC.
type Subnet struct {
	ID string `json:"id"`

	VpcID            string  `json:"vpcId"`
	AvailabilityZone string  `json:"availabilityZone"`
	CidrBlock        string  `json:"cidrBlock"`
	IsPublic         bool    `json:"public"`
	RouteTableID     *string `json:"routeTableId"`
	NatGatewayID     *string `json:"natGatewayId"`
}

// String returns a string representation of the subnet.
func (s *Subnet) String() string {
	return fmt.Sprintf("id=%s/az=%s/public=%v", s.ID, s.AvailabilityZone, s.IsPublic)
}

// Subnets is a slice of Subnet.
type Subnets []*Subnet

// ToMap returns a map from id to subnet.
func (s Subnets) ToMap() map[string]*Subnet {
	res := make(map[string]*Subnet)
	for _, x := range s {
		res[x.ID] = x
	}
	return res
}

// FilterPrivate returns a slice containing all subnets marked as private.
func (s Subnets) FilterPrivate() (res Subnets) {
	for _, x := range s {
		if !x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterPublic returns a slice containing all subnets marked as public.
func (s Subnets) FilterPublic() (res Subnets) {
	for _, x := range s {
		if x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// RouteTable defines an AWS routing table.
type RouteTable struct {
	ID string `json:"id"`
}

// SecurityGroupRole defines the unique role of a security group.
type SecurityGroupRole string

var (
	SecurityGroupBastion      = SecurityGroupRole("bastion")
	SecurityGroupNode         = SecurityGroupRole("node")
	SecurityGroupControlPlane = SecurityGroupRole("controlplane")
)

// SecurityGroup defines an AWS security group.
type SecurityGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	IngressRules IngressRules `json:"ingressRule"`
}

// String returns a string representation of the security group.
func (s *SecurityGroup) String() string {
	return fmt.Sprintf("id=%s/name=%s", s.ID, s.Name)
}

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

var (
	SecurityGroupProtocolAll  = SecurityGroupProtocol("-1")
	SecurityGroupProtocolTCP  = SecurityGroupProtocol("tcp")
	SecurityGroupProtocolUDP  = SecurityGroupProtocol("udp")
	SecurityGroupProtocolICMP = SecurityGroupProtocol("icmp")
)

// IngressRule defines an AWS ingress rule for security groups.
type IngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`
	FromPort    int64                 `json:"fromPort"`
	ToPort      int64                 `json:"toPort"`

	// List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.
	CidrBlocks []string `json:"cidrBlocks"`

	// The security group id to allow access from. Cannot be specified with CidrBlocks.
	SourceSecurityGroupID *string `json:"sourceSecurityGroupId"`
}

// String returns a string representation of the ingress rule.
func (i *IngressRule) String() string {
	return fmt.Sprintf("protocol=%s/range=[%d-%d]/description=%s", i.Protocol, i.FromPort, i.ToPort, i.Description)
}

// IngressRules is a slice of AWS ingress rules for security groups.
type IngressRules []*IngressRule

// Returns the difference between this slice and the other slice.
func (i IngressRules) Difference(o IngressRules) (out IngressRules) {
	for _, x := range i {
		found := false
		for _, y := range o {
			if reflect.DeepEqual(x, y) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, x)
		}
	}

	return
}

// InstanceState describes the state of an AWS instance.
type InstanceState string

var (
	InstanceStatePending      = InstanceState("pending")
	InstanceStateRunning      = InstanceState("running")
	InstanceStateShuttingDown = InstanceState("shutting-down")
	InstanceStateTerminated   = InstanceState("terminated")
	InstanceStateStopping     = InstanceState("stopping")
	InstanceStateStopped      = InstanceState("stopped")
)

// Instance describes an AWS instance.
type Instance struct {
	ID string `json:"id"`

	// The current state of the instance.
	State InstanceState `json:"instanceState"`

	// The instance type.
	Type string `json:"type"`

	// The ID of the subnet of the instance.
	SubnetID string `json:"subnetId"`

	// The ID of the AMI used to launch the instance.
	ImageID string `json:"imageId"`

	// The name of the SSH key pair.
	KeyName *string `json:"keyName"`

	// SecurityGroupIDs are one or more security group IDs this instance belongs to.
	SecurityGroupIDs []string `json:"securityGroupIds"`

	// UserData is the raw data script passed to the instance which is run upon bootstrap.
	// This field must not be base64 encoded and should only be used when running a new instance.
	UserData *string `json:"userData"`

	// The ARN of the IAM instance profile associated with the instance, if applicable.
	IAMProfile *AWSResourceReference `json:"iamProfile"`

	// The private IPv4 address assigned to the instance.
	PrivateIP *string `json:"privateIp"`

	// The public IPv4 address assigned to the instance, if applicable.
	PublicIP *string `json:"publicIp"`

	// Specifies whether enhanced networking with ENA is enabled.
	ENASupport *bool `json:"enaSupport"`

	// Indicates whether the instance is optimized for Amazon EBS I/O.
	EBSOptimized *bool `json:"ebsOptimized"`

	// The tags associated with the instance.
	Tags map[string]string `json:"tag"`
}
