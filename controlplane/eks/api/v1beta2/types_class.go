/*
Copyright 2023 The Kubernetes Authors.
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

import infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

// AWSManagedControlPlaneClassSpec defines the AWSManagedControlPlane properties that may be shared across
// several AWS managed control planes.
type AWSManagedControlPlaneClassSpec struct {

	// IdentityRef is a reference to a identity to be used when reconciling the managed control plane.
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec infrav1.NetworkSpec `json:"network,omitempty"`

	// SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
	// Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.
	// +optional
	SecondaryCidrBlock *string `json:"secondaryCidrBlock,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// Partition is the AWS security partition being used. Defaults to "aws"
	// +optional
	Partition string `json:"partition,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetes that EKS supports
	// will be used.
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Pattern:=^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.?(\.0|[1-9][0-9]*)?$
	// +optional
	Version *string `json:"version,omitempty"`

	// RoleName specifies the name of IAM role that gives EKS
	// permission to make API calls. If the role is pre-existing
	// we will treat it as unmanaged and not delete it on
	// deletion. If the EKSEnableIAM feature flag is true
	// and no name is supplied then a role is created.
	// +kubebuilder:validation:MinLength:=2
	// +optional
	RoleName *string `json:"roleName,omitempty"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the control plane role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies *[]string `json:"roleAdditionalPolicies,omitempty"`

	// Logging specifies which EKS Cluster logs should be enabled. Entries for
	// each of the enabled logs will be sent to CloudWatch
	// +optional
	Logging *ControlPlaneLoggingSpec `json:"logging,omitempty"`

	// EncryptionConfig specifies the encryption configuration for the cluster
	// +optional
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// IAMAuthenticatorConfig allows the specification of any additional user or role mappings
	// for use when generating the aws-iam-authenticator configuration. If this is nil the
	// default configuration is still generated for the cluster.
	// +optional
	IAMAuthenticatorConfig *IAMAuthenticatorConfig `json:"iamAuthenticatorConfig,omitempty"`

	// ImageLookupFormat is the AMI naming format to look up machine images when
	// a machine does not specify an AMI. When set, this will be used for all
	// cluster machines unless a machine specifies a different ImageLookupOrg.
	// Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
	// OS and kubernetes version, respectively. The BaseOS will be the value in
	// ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
	// defined by the packages produced by kubernetes/release without v as a
	// prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
	// image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
	// searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
	// Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
	// also: https://golang.org/pkg/text/template/
	// +optional
	ImageLookupFormat string `json:"imageLookupFormat,omitempty"`

	// ImageLookupOrg is the AWS Organization ID to look up machine images when a
	// machine does not specify an AMI. When set, this will be used for all
	// cluster machines unless a machine specifies a different ImageLookupOrg.
	// +optional
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// ImageLookupBaseOS is the name of the base operating system used to look
	// up machine images when a machine does not specify an AMI. When set, this
	// will be used for all cluster machines unless a machine specifies a
	// different ImageLookupBaseOS.
	ImageLookupBaseOS string `json:"imageLookupBaseOS,omitempty"`

	// TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
	// iam-authenticator - obtains a client token using iam-authentictor
	// aws-cli - obtains a client token using the AWS CLI
	// Defaults to iam-authenticator
	// +kubebuilder:default=iam-authenticator
	// +kubebuilder:validation:Enum=iam-authenticator;aws-cli
	TokenMethod *EKSTokenMethod `json:"tokenMethod,omitempty"`

	// Bastion contains options to configure the bastion host.
	// +optional
	Bastion infrav1.Bastion `json:"bastion"`

	// AssociateOIDCProvider can be enabled to automatically create an identity
	// provider for the controller for use with IAM roles for service accounts
	// +kubebuilder:default=false
	AssociateOIDCProvider bool `json:"associateOIDCProvider,omitempty"`

	// Addons defines the EKS addons to enable with the EKS cluster.
	// +optional
	Addons *[]Addon `json:"addons,omitempty"`

	// IdentityProviderconfig is used to specify the oidc provider config
	// to be attached with this eks cluster
	// +optional
	OIDCIdentityProviderConfig *OIDCIdentityProviderConfig `json:"oidcIdentityProviderConfig,omitempty"`

	// VpcCni is used to set configuration options for the VPC CNI plugin
	// +optional
	VpcCni VpcCni `json:"vpcCni,omitempty"`

	// KubeProxy defines managed attributes of the kube-proxy daemonset
	KubeProxy KubeProxy `json:"kubeProxy,omitempty"`

	// Endpoints specifies access to this cluster's control plane endpoints
	// +optional
	EndpointAccess EndpointAccess `json:"endpointAccess,omitempty"`
}
