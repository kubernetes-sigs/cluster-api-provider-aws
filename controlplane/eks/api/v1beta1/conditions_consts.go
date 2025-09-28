/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

const (
	// EKSControlPlaneReadyCondition condition reports on the successful reconciliation of eks control plane.
	EKSControlPlaneReadyCondition = "EKSControlPlaneReady"
	// EKSControlPlaneCreatingCondition condition reports on whether the eks
	// control plane is creating.
	EKSControlPlaneCreatingCondition = "EKSControlPlaneCreating"
	// EKSControlPlaneUpdatingCondition condition reports on whether the eks
	// control plane is updating.
	EKSControlPlaneUpdatingCondition = "EKSControlPlaneUpdating"
	// EKSControlPlaneReconciliationFailedReason used to report failures while reconciling EKS control plane.
	EKSControlPlaneReconciliationFailedReason = "EKSControlPlaneReconciliationFailed"
)

const (
	// IAMControlPlaneRolesReadyCondition condition reports on the successful reconciliation of eks control plane iam roles.
	IAMControlPlaneRolesReadyCondition = "IAMControlPlaneRolesReady"
	// IAMControlPlaneRolesReconciliationFailedReason used to report failures while reconciling EKS control plane iam roles.
	IAMControlPlaneRolesReconciliationFailedReason = "IAMControlPlaneRolesReconciliationFailed"
)

const (
	// IAMAuthenticatorConfiguredCondition condition reports on the successful reconciliation of aws-iam-authenticator config.
	IAMAuthenticatorConfiguredCondition = "IAMAuthenticatorConfigured"
	// IAMAuthenticatorConfigurationFailedReason used to report failures while reconciling the aws-iam-authenticator config.
	IAMAuthenticatorConfigurationFailedReason = "IAMAuthenticatorConfigurationFailed"
)

const (
	// EKSAddonsConfiguredCondition condition reports on the successful reconciliation of EKS addons.
	EKSAddonsConfiguredCondition = "EKSAddonsConfigured"
	// EKSAddonsConfiguredFailedReason used to report failures while reconciling the EKS addons.
	EKSAddonsConfiguredFailedReason = "EKSAddonsConfiguredFailed"
)

const (
	// EKSIdentityProviderConfiguredCondition condition reports on the successful association of identity provider config.
	EKSIdentityProviderConfiguredCondition = "EKSIdentityProviderConfigured"
	// EKSIdentityProviderConfiguredFailedReason used to report failures while reconciling the identity provider config association.
	EKSIdentityProviderConfiguredFailedReason = "EKSIdentityProviderConfiguredFailed"
)
