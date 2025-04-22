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

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ROSARoleConfigSpec defines the desired state of ROSARoleConfig
type ROSARoleConfigSpec struct {
	AccountRoleConfig  AccountRoleConfig             `json:"accountRoleConfig"`
	OperatorRoleConfig OperatorRoleConfig            `json:"operatorRoleConfig"`
	OIDCConfig         OIDCConfig                    `json:"oidcConfig"`
	IdentityRef        *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`
	Region             string                        `json:"region,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosaroleconfigs,scope=Namespaced,categories=cluster-api,shortName=rosarole
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// ROSARoleConfig is the Schema for the rosaroleconfigs API
type ROSARoleConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSARoleConfigSpec   `json:"spec,omitempty"`
	Status ROSARoleConfigStatus `json:"status,omitempty"`
}

// AccountRoleConfig defines account-wide IAM roles before creating your ROSA cluster.
type AccountRoleConfig struct {
	// User-defined prefix for all generated AWS resources
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Required
	Prefix string `json:"prefix"`
	// The ARN of the policy that is used to set the permissions boundary for the account roles.
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`
	// The arn path for the account/operator roles as well as their policies.
	// +optional
	Path string `json:"path,omitempty"`
	//  Version of OpenShift that will be used to setup policy tag, for example "4.11"
	// +kubebuilder:validation:Required
	Version string `json:"version"`
	// SharedVPCConfig is used to set up shared VPC.
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

// OperatorRoleConfig defines cluster-specific operator IAM roles based on your cluster configuration.
type OperatorRoleConfig struct {
	//  User-defined prefix for generated AWS operator policies.
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Required
	Prefix string `json:"prefix"`
	// The ARN of the policy that is used to set the permissions boundary for the operator roles.
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`
	// Registered OIDC configuration ID to add its issuer URL as the trusted relationship to the operator roles.''
	// +kubebuilder:validation:Required
	OIDCConfigID string `json:"oidcConfigId,omitempty"`
	// SharedVPCConfig is used to set up shared VPC.
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

// SharedVPCConfig is used to set up shared VPC.
type SharedVPCConfig struct {
	//  Role ARN associated with the private hosted zone used for Hosted Control Plane cluster shared VPC, this role contains policies to be used with Route 53
	RouteRoleARN string `json:"routeRoleARN,omitempty"`
	//  Role ARN associated with the shared VPC used for Hosted Control Plane clusters, this role contains policies to be used with the VPC endpoint
	VPCEndpointRoleARN string `json:"vpcEndpointRoleArn,omitempty"`
}

// OIDCConfig creates OIDC config in a S3 bucket for the client AWS account and populates it to be compliant with OIDC protocol.
// It also creates a Secret in Secrets Manager containing the private key.
type OIDCConfig struct {
	// ManagedOIDC indicates whether it is a Red Hat managed or unmanaged (Customer hosted) OIDC Configuration. Default is true.
	// +kubebuilder:default=true
	ManagedOIDC bool `json:"managedOIDC"`
	// ExternalAuthProviders are external OIDC identity providers that can issue tokens for this cluster.
	// Can only be set if "enableExternalAuthProviders" is set to "True".
	//
	// At most one provider can be configured.
	//
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=1
	ExternalAuthProviders []rosacontrolplanev1.ExternalAuthProvider `json:"externalAuthProviders,omitempty"`
}

// ROSARoleConfigStatus defines the observed state of ROSARoleConfig
type ROSARoleConfigStatus struct {
	// ID of created OIDC config
	OIDCID string `json:"oidcID,omitempty"`
	// Create OIDC provider for operators to authenticate against in an STS cluster.
	OIDCProviderARN string `json:"oidcProviderARN,omitempty"`
	// Created Account roles that can be used to
	AccountRolesRef AccountRolesRef `json:"accountRolesRef,omitempty"`
	// AWS IAM roles used to perform credential requests by the openshift operators.
	OperatorRolesRef rosacontrolplanev1.AWSRolesRef `json:"operatorRolesRef"`
	Conditions       clusterv1.Conditions           `json:"conditions,omitempty"`
}

// AccountRolesRef defscribes ARNs used as Account roles.
type AccountRolesRef struct {
	// InstallerRoleARN is an AWS IAM role that OpenShift Cluster Manager will assume to create the cluster..
	InstallerRoleARN string `json:"installerRoleARN"`
	// SupportRoleARN is an AWS IAM role used by Red Hat SREs to enable
	// access to the cluster account in order to provide support.
	SupportRoleARN string `json:"supportRoleARN"`
	// WorkerRoleARN is an AWS IAM role that will be attached to worker instances.
	WorkerRoleARN string `json:"workerRoleARN,omitempty"`
}

// ROSARoleConfigList contains a list of ROSARoleConfig
// +kubebuilder:object:root=true
type ROSARoleConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSARoleConfig `json:"items"`
}

const (
	// RosaRoleConfigReadyCondition condition reports on the successful reconciliation of RosaRoleConfig.
	RosaRoleConfigReadyCondition clusterv1.ConditionType = "RosaRoleConfigReady"

	// RosaRoleConfigDeletionFailedReason used to report failures while deleting RosaRoleConfig.
	RosaRoleConfigDeletionFailedReason = "DeletionFailed"
)

// GetConditions returns the observations of the operational state of the RosaNetwork resource.
func (r *ROSARoleConfig) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the RosaRoleConfig to the predescribed clusterv1.Conditions.
func (r *ROSARoleConfig) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSARoleConfig{}, &ROSARoleConfigList{})
}
