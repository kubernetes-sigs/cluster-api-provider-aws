/*
Copyright 2026 The Kubernetes Authors.

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

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// AWSManagedControlPlane v1beta2 condition types.
const (
	// AWSManagedControlPlaneReadyCondition defines the Ready condition type that summarizes the operational state of an AWSManagedControlPlane.
	AWSManagedControlPlaneReadyCondition = clusterv1.ReadyCondition

	// AWSManagedControlPlaneEKSControlPlaneReadyCondition reports on the successful reconciliation of the EKS control plane.
	AWSManagedControlPlaneEKSControlPlaneReadyCondition = "EKSControlPlaneReady"

	// AWSManagedControlPlaneEKSControlPlaneCreatingCondition reports whether the EKS control plane is being created.
	AWSManagedControlPlaneEKSControlPlaneCreatingCondition = "EKSControlPlaneCreating"

	// AWSManagedControlPlaneEKSControlPlaneUpdatingCondition reports whether the EKS control plane is being updated.
	AWSManagedControlPlaneEKSControlPlaneUpdatingCondition = "EKSControlPlaneUpdating"

	// AWSManagedControlPlaneIAMControlPlaneRolesReadyCondition reports on the successful reconciliation of EKS control plane IAM roles.
	AWSManagedControlPlaneIAMControlPlaneRolesReadyCondition = "IAMControlPlaneRolesReady"

	// AWSManagedControlPlaneIAMAuthenticatorConfiguredCondition reports on the successful reconciliation of the aws-iam-authenticator config.
	AWSManagedControlPlaneIAMAuthenticatorConfiguredCondition = "IAMAuthenticatorConfigured"

	// AWSManagedControlPlaneEKSAddonsConfiguredCondition reports on the successful reconciliation of EKS addons.
	AWSManagedControlPlaneEKSAddonsConfiguredCondition = "EKSAddonsConfigured"

	// AWSManagedControlPlaneEKSIdentityProviderConfiguredCondition reports on the successful association of identity provider config.
	AWSManagedControlPlaneEKSIdentityProviderConfiguredCondition = "EKSIdentityProviderConfigured"
)

// AWSManagedControlPlane v1beta2 reason constants.
const (
	// AWSManagedControlPlaneReadyReason indicates the AWSManagedControlPlane is ready.
	AWSManagedControlPlaneReadyReason = clusterv1.ReadyReason

	// AWSManagedControlPlaneNotReadyReason indicates the AWSManagedControlPlane is not ready.
	AWSManagedControlPlaneNotReadyReason = clusterv1.NotReadyReason

	// AWSManagedControlPlaneDeletingReason indicates the AWSManagedControlPlane is being deleted.
	AWSManagedControlPlaneDeletingReason = clusterv1.DeletingReason

	// AWSManagedControlPlaneEKSControlPlaneReconciliationFailedReason used to report failures while reconciling EKS control plane.
	AWSManagedControlPlaneEKSControlPlaneReconciliationFailedReason = "EKSControlPlaneReconciliationFailed"

	// AWSManagedControlPlaneIAMControlPlaneRolesReconciliationFailedReason used to report failures while reconciling EKS control plane IAM roles.
	AWSManagedControlPlaneIAMControlPlaneRolesReconciliationFailedReason = "IAMControlPlaneRolesReconciliationFailed"

	// AWSManagedControlPlaneIAMAuthenticatorConfigurationFailedReason used to report failures while reconciling the aws-iam-authenticator config.
	AWSManagedControlPlaneIAMAuthenticatorConfigurationFailedReason = "IAMAuthenticatorConfigurationFailed"

	// AWSManagedControlPlaneEKSAddonsConfiguredFailedReason used to report failures while reconciling the EKS addons.
	AWSManagedControlPlaneEKSAddonsConfiguredFailedReason = "EKSAddonsConfiguredFailed"

	// AWSManagedControlPlaneEKSIdentityProviderConfiguredFailedReason used to report failures while reconciling the identity provider config association.
	AWSManagedControlPlaneEKSIdentityProviderConfiguredFailedReason = "EKSIdentityProviderConfiguredFailed"
)
