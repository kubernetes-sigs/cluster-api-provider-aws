/*
Copyright 2025 The Kubernetes Authors.

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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// AWSManagedControlPlaneClassSpec defines the AWSManagedControlPlane properties that may be shared across several AWS managed control planes.
type AWSManagedControlPlaneClassSpec struct {
	// MachineTemplate contains information about how machines
	// should be shaped when creating or updating a control plane.
	// For the AWSManagedControlPlaneTemplate, this field is used
	// only to fulfill the CAPI contract.
	// +optional
	MachineTemplate *AWSManagedControlPlaneTemplateMachineTemplate `json:"machineTemplate,omitempty"`

	// IdentityRef is a reference to an identity to be used when reconciling the managed control plane.
	// If no identity is specified, the default identity for this controller will be used.
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// NetworkSpec encapsulates all things related to AWS networking.
	NetworkSpec infrav1.NetworkSpec `json:"network,omitempty"`

	// SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
	// Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.
	// +optional
	SecondaryCidrBlock *string `json:"secondaryCidrBlock,omitempty"`

	// Region is the AWS Region in which the cluster is located.
	Region string `json:"region,omitempty"`

	// Partition is the AWS security partition being used. Defaults to "aws".
	// +optional
	Partition string `json:"partition,omitempty"`

	// SSHKeyName is the name of the SSH key to attach to the bastion host.
	// Valid values are empty string (do not use SSH keys), a valid SSH key name,
	// or omitted (use the default SSH key name).
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version.
	// If no version is supplied, then the latest version of Kubernetes that EKS supports will be used.
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Pattern:=^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.?(\.0|[1-9][0-9]*)?$
	// +optional
	Version *string `json:"version,omitempty"`

	// RoleName specifies the name of the IAM role that gives EKS permission to make API calls.
	// If the role is pre-existing, it will be treated as unmanaged and not deleted on cluster deletion.
	// +kubebuilder:validation:MinLength:=2
	// +optional
	RoleName *string `json:"roleName,omitempty"`

	// RoleAdditionalPolicies allows you to attach additional policies to the control plane role.
	// You must enable the EKSAllowAddRoles feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies *[]string `json:"roleAdditionalPolicies,omitempty"`

	// RolePath sets the path to the role. For more information about paths, see IAM Identifiers
	// (https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html)
	// in the IAM User Guide.
	//
	// This parameter is optional. If it is not included, it defaults to a slash
	// (/).
	// +optional
	RolePath string `json:"rolePath,omitempty"`

	// RolePermissionsBoundary sets the ARN of the managed policy that is used
	// to set the permissions boundary for the role.
	//
	// A permissions boundary policy defines the maximum permissions that identity-based
	// policies can grant to an entity, but does not grant permissions. Permissions
	// boundaries do not define the maximum permissions that a resource-based policy
	// can grant to an entity. To learn more, see Permissions boundaries for IAM
	// entities (https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html)
	// in the IAM User Guide.
	//
	// For more information about policy types, see Policy types (https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#access_policy-types)
	// in the IAM User Guide.
	// +optional
	RolePermissionsBoundary string `json:"rolePermissionsBoundary,omitempty"`

	// Logging specifies which EKS Cluster logs should be enabled.
	// Entries for each enabled log will be sent to CloudWatch.
	// +optional
	Logging *ControlPlaneLoggingSpec `json:"logging,omitempty"`

	// EncryptionConfig specifies the encryption configuration for the cluster.
	// +optional
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider,
	// in addition to the ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// IAMAuthenticatorConfig allows the specification of any additional user or role mappings
	// for generating the aws-iam-authenticator configuration.
	// If this is nil, the default configuration is still generated for the cluster.
	// +optional
	IAMAuthenticatorConfig *IAMAuthenticatorConfig `json:"iamAuthenticatorConfig,omitempty"`

	// EndpointAccess specifies how the control plane endpoints are accessible.
	// +optional
	EndpointAccess EndpointAccess `json:"endpointAccess,omitempty"`

	// ImageLookupFormat is the AMI naming format used to look up machine images when a machine does not specify an AMI.
	// Supports substitutions for {{.BaseOS}} and {{.K8sVersion}}.
	// +optional
	ImageLookupFormat string `json:"imageLookupFormat,omitempty"`

	// ImageLookupOrg is the AWS Organization ID used to look up machine images when a machine does not specify an AMI.
	// +optional
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// ImageLookupBaseOS is the name of the base operating system used to look up machine images when a machine does not specify an AMI.
	// +optional
	ImageLookupBaseOS string `json:"imageLookupBaseOS,omitempty"`

	// Bastion contains options to configure the bastion host.
	// +optional
	Bastion infrav1.Bastion `json:"bastion"`

	// TokenMethod is used to specify the method for obtaining a client token for communicating with EKS.
	// Defaults to iam-authenticator.
	// +kubebuilder:default=iam-authenticator
	// +kubebuilder:validation:Enum=iam-authenticator;aws-cli
	// +optional
	TokenMethod *EKSTokenMethod `json:"tokenMethod,omitempty"`

	// AssociateOIDCProvider can be enabled to automatically create an identity provider
	// for the controller for use with IAM roles for service accounts.
	// +kubebuilder:default=false
	// +optional
	AssociateOIDCProvider bool `json:"associateOIDCProvider,omitempty"`

	// Addons defines the EKS addons to enable with the cluster.
	// +optional
	Addons *[]Addon `json:"addons,omitempty"`

	// OIDCIdentityProviderConfig is used to specify the OIDC provider configuration to attach to this EKS cluster.
	// +optional
	OIDCIdentityProviderConfig *OIDCIdentityProviderConfig `json:"oidcIdentityProviderConfig,omitempty"`

	// VpcCni specifies configuration options for the VPC CNI plugin.
	// +optional
	VpcCni VpcCni `json:"vpcCni,omitempty"`

	// BootstrapSelfManagedAddons is used to set configuration options for
	// bare EKS cluster without EKS default networking addons
	// If you set this value to false when creating a cluster, the default networking add-ons will not be installed
	// +kubebuilder:default=true
	BootstrapSelfManagedAddons bool `json:"bootstrapSelfManagedAddons,omitempty"`

	// RestrictPrivateSubnets indicates that the EKS control plane should only use private subnets.
	// +kubebuilder:default=false
	// +optional
	RestrictPrivateSubnets bool `json:"restrictPrivateSubnets,omitempty"`

	// KubeProxy defines managed attributes of the kube-proxy daemonset.
	// +optional
	KubeProxy KubeProxy `json:"kubeProxy,omitempty"`
}
