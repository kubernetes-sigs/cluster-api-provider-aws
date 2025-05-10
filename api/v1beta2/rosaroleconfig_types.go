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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ROSARoleConfigSpec defines the desired state of ROSARoleConfig

type ROSARoleConfigSpec struct {
	AccountRoleConfig  AccountRoleConfig  `json:"accountRoleConfig"`
	OperatorRoleConfig OperatorRoleConfig `json:"operatorRoleConfig"`
	OIDCConfig         OIDCConfig         `json:"oidcConfig"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosaroleconfig,scope=cluster,categories=cluster-api,shortName=rosarole
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

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
	// +kubebuilder:validation:Required
	OIDCConfigID string `json:"oidcConfigId,omitempty"`
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

type SharedVPCConfig struct {
	RouteRoleARN       string `json:"routeRoleARN,omitempty"`
	VPCEndpointRoleARN string `json:"vpcEndpointRoleArn,omitempty"`
}

type OIDCConfig struct {
	CreateManagedOIDC     string                 `json:"createManagedOIDC"`
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
	Conditions       []metav1.Condition `json:"conditions,omitempty"`
	OIDCID           string             `json:"oidcID,omitempty"`
	OIDCProviderARN  string             `json:"oidcProviderARN,omitempty"`
	AccountRolesRef  AccountRolesRef    `json:"accountRolesRef,omitempty"`
	OperatorRolesRef OperatorRolesRef   `json:"operatorRolesRef,omitempty"`
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

func init() {
	SchemeBuilder.Register(&ROSARoleConfig{}, &ROSARoleConfigList{})
}
