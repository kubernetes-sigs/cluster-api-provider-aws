/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

// ControlPlaneLoggingSpec defines what EKS control plane logs that should be enabled
type ControlPlaneLoggingSpec struct {
	// APIServer indicates if the Kubernetes API Server log (kube-apiserver) shoulkd be enabled
	// +kubebuilder:default=false
	APIServer bool `json:"apiServer"`
	// Audit indicates if the Kubernetes API audit log should be enabled
	// +kubebuilder:default=false
	Audit bool `json:"audit"`
	// Authenticator indicates if the iam authenticator log should be enabled
	// +kubebuilder:default=false
	Authenticator bool `json:"authenticator"`
	//ControllerManager indicates if the controller manager (kube-controller-manager) log should be enabled
	// +kubebuilder:default=false
	ControllerManager bool `json:"controllerManager"`
	// Scheduler indicates if the Kubernetes scheduler (kube-scheduler) log should be enabled
	// +kubebuilder:default=false
	Scheduler bool `json:"scheduler"`
}

// IsLogEnabled returns true if the log is enabled
func (s *ControlPlaneLoggingSpec) IsLogEnabled(logName string) bool {
	if s == nil {
		return false
	}

	switch logName {
	case eks.LogTypeApi:
		return s.APIServer
	case eks.LogTypeAudit:
		return s.Audit
	case eks.LogTypeAuthenticator:
		return s.Authenticator
	case eks.LogTypeControllerManager:
		return s.ControllerManager
	case eks.LogTypeScheduler:
		return s.Scheduler
	default:
		return false
	}
}

// EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.
type EKSTokenMethod string

var (
	// EKSTokenMethodIAMAuthenticator indicates that IAM autenticator will be used to get a token
	EKSTokenMethodIAMAuthenticator = EKSTokenMethod("iam-authenticator")

	// EKSTokenMethodAWSCli indicates that the AWS CLI will be used to get a token
	// Version 1.16.156 or greater is required of the AWS CLI
	EKSTokenMethodAWSCli = EKSTokenMethod("aws-cli")
)

var (
	// DefaultEKSControlPlaneRole is the name of the default IAM role to use for the EKS control plane
	// if no other role is supplied in the spec and if iam role creation is not enabled. The default
	// can be created using clusterawsadm or created manually
	DefaultEKSControlPlaneRole = fmt.Sprintf("eks-controlplane%s", infrav1.DefaultNameSuffix)
)

// IAMAuthenticatorConfig represents an aws-iam-authenticator configuration
type IAMAuthenticatorConfig struct {
	// RoleMappings is a list of role mappings
	// +optional
	RoleMappings []RoleMapping `json:"mapRoles,omitempty"`
	// UserMappings is a list of user mappings
	// +optional
	UserMappings []UserMapping `json:"mapUsers,omitempty"`
}

// KubernetesMapping represents the kubernetes RBAC mapping
type KubernetesMapping struct {
	// UserName is a kubernetes RBAC user subject
	UserName string `json:"username"`
	// Groups is a list of kubernetes RBAC groups
	Groups []string `json:"groups"`
}

// RoleMapping represents a mapping from a IAM role to Kubernetes users and groups
type RoleMapping struct {
	// RoleARN is the AWS ARN for the role to map
	// +kubebuilder:validation:MinLength:=31
	RoleARN string `json:"rolearn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}

// UserMapping represents a mapping from an IAM user to Kubernetes users and groups
type UserMapping struct {
	// UserARN is the AWS ARN for the user to map
	// +kubebuilder:validation:MinLength:=31
	UserARN string `json:"userarn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}
