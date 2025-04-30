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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ROSARoleConfigSpec defines the desired state of ROSARoleConfig

type ROSARoleConfigSpec struct {
	AccountRoleConfig  AccountRoleConfig  `json:"accountRoleConfig"`
	OperatorRoleConfig OperatorRoleConfig `json:"operatorRoleConfig"`
	OIDCConfig         OIDCConfig         `json:"oidcConfig"`

	// IdentityRef is a reference to an identity to be used when reconciling rosa roles config.
	// If no identity is specified, the default identity for this controller will be used.
	//
	// +optional
	IdentityRef *AWSIdentityReference `json:"identityRef,omitempty"`
	Region      string                `json:"region,omitempty"`
}

// +kubebuilder:object:root=true

// ROSARoleConfig is the Schema for the rosaroleconfigs API
type ROSARoleConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSARoleConfigSpec   `json:"spec,omitempty"`
	Status ROSARoleConfigStatus `json:"status,omitempty"`
}

type AccountRoleConfig struct {
	// TODO: verify max len
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Required
	Prefix string `json:"prefix"`
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`
	// +optional
	Path string `json:"path,omitempty"`
	// +kubebuilder:validation:Required
	Version string `json:"version"`
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

type OperatorRoleConfig struct {
	// TODO: verify max len
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Required
	Prefix string `json:"prefix"`
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`
	OIDCConfigID           string `json:"oidcConfigId,omitempty"`
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

type SharedVPCConfig struct {
	RouteRoleARN       string `json:"routeRoleARN,omitempty"`
	VPCEndpointRoleARN string `json:"vpcEndpointRoleArn,omitempty"`
}

type OIDCConfig struct {
	ManagedOIDC bool `json:"managedOIDC"`
	// Prefix is required for Unmanaged OIDC
	// +optional
	Prefix string `json:"prefix"`
	// Region is required for Unmanaged OIDC
	// +optional
	Region                string                 `json:"region"`
	ExternalAuthProviders []ExternalAuthProvider `json:"externalAuthProviders,omitempty"`
}

type ExternalAuthProvider struct {
	Name          string        `json:"name"`
	Issuer        Issuer        `json:"issuer"`
	ClaimMappings ClaimMappings `json:"claimMappings"`
	OIDCClients   []OIDCClient  `json:"oidcClients,omitempty"`
}

type Issuer struct {
	IssuerURL string   `json:"issuerURL"`
	Audiences []string `json:"audiences,omitempty"`
}

type ClaimMappings struct {
	Username Mapping `json:"username"`
	Groups   Mapping `json:"groups,omitempty"`
}

type Mapping struct {
	Claim        string `json:"claim"`
	PrefixPolicy string `json:"prefixPolicy,omitempty"`
}

type OIDCClient struct {
	ComponentName      string          `json:"componentName"`
	ComponentNamespace string          `json:"componentNamespace"`
	ClientID           string          `json:"clientID"`
	ClientSecret       SecretReference `json:"clientSecret"`
}

type SecretReference struct {
	Name string `json:"name"`
}

// ROSARoleConfigStatus defines the observed state of ROSARoleConfig
type ROSARoleConfigStatus struct {
	// Conditions       []metav1.Condition `json:"conditions,omitempty"`
	Conditions       clusterv1.Conditions `json:"conditions,omitempty"`
	OIDCID           string               `json:"oidcID,omitempty"`
	OIDCProviderARN  string               `json:"oidcProviderARN,omitempty"`
	AccountRolesRef  AccountRolesRef      `json:"accountRolesRef,omitempty"`
	OperatorRolesRef OperatorRolesRef     `json:"operatorRolesRef,omitempty"`
}

type AccountRolesRef struct {
	InstallerRoleARN string `json:"installerRoleARN,omitempty"`
	SupportRoleARN   string `json:"supportRoleARN,omitempty"`
	WorkerRoleARN    string `json:"workerRoleARN,omitempty"`
}

type OperatorRolesRef struct {
	IngressARN              string `json:"ingressARN,omitempty"`
	ImageRegistryARN        string `json:"imageRegistryARN,omitempty"`
	StorageARN              string `json:"storageARN,omitempty"`
	NetworkARN              string `json:"networkARN,omitempty"`
	KubeCloudControllerARN  string `json:"kubeCloudControllerARN,omitempty"`
	NodePoolManagementARN   string `json:"nodePoolManagementARN,omitempty"`
	ControlPlaneOperatorARN string `json:"controlPlaneOperatorARN,omitempty"`
	KMSProviderARN          string `json:"kmsProviderARN,omitempty"`
}

// +kubebuilder:object:root=true

// ROSARoleConfigList contains a list of ROSARoleConfig
type ROSARoleConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSARoleConfig `json:"items"`
}

const (
	// RosaRoleConfigReadyCondition condition reports on the successful reconciliation of RosaNetwork.
	RosaRoleConfigReadyCondition clusterv1.ConditionType = "RosaRoleConfigReady"

	// RosaRoleConfigDeletionFailedReason used to report failures while deleting RosaNetwork.
	RosaRoleConfigDeletionFailedReason = "DeletionFailed"
)

// GetConditions returns the observations of the operational state of the RosaNetwork resource.
func (r *ROSARoleConfig) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the RosaNetwork to the predescribed clusterv1.Conditions.
func (r *ROSARoleConfig) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSARoleConfig{}, &ROSARoleConfigList{})
}
