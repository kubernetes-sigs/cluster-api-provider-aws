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

	VPC VPC `json:"vpc"`
}

// VPC is the aws vpc
type VPC struct {
	ID        string `json:"id"`
	CidrBlock string `json:"cidrBlock"`
}
