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

// ROSAControlPlane v1beta2 condition types.
const (
	// ROSAControlPlaneReadyCondition defines the Ready condition type that summarizes the operational state of a ROSAControlPlane.
	ROSAControlPlaneReadyCondition = clusterv1.ReadyCondition

	// ROSAControlPlaneControlPlaneReadyCondition reports on the successful reconciliation of ROSAControlPlane.
	ROSAControlPlaneControlPlaneReadyCondition = "ROSAControlPlaneReady"

	// ROSAControlPlaneValidCondition reports whether ROSAControlPlane configuration is valid.
	ROSAControlPlaneValidCondition = "ROSAControlPlaneValid"

	// ROSAControlPlaneUpgradingCondition reports whether ROSAControlPlane is upgrading or not.
	ROSAControlPlaneUpgradingCondition = "ROSAControlPlaneUpgrading"

	// ROSAControlPlaneExternalAuthConfiguredCondition reports whether external auth has been correctly configured.
	ROSAControlPlaneExternalAuthConfiguredCondition = "ExternalAuthConfigured"

	// ROSAControlPlaneRoleConfigReadyCondition reports whether the referenced RosaRoleConfig is ready.
	ROSAControlPlaneRoleConfigReadyCondition = "ROSARoleConfigReady"
)

// ROSAControlPlane v1beta2 reason constants.
const (
	// ROSAControlPlaneReadyReason indicates the ROSAControlPlane is ready.
	ROSAControlPlaneReadyReason = clusterv1.ReadyReason

	// ROSAControlPlaneNotReadyReason indicates the ROSAControlPlane is not ready.
	ROSAControlPlaneNotReadyReason = clusterv1.NotReadyReason

	// ROSAControlPlaneDeletingReason indicates the ROSAControlPlane is being deleted.
	ROSAControlPlaneDeletingReason = clusterv1.DeletingReason

	// ROSAControlPlaneReconciliationFailedReason used to report reconciliation failures.
	ROSAControlPlaneReconciliationFailedReason = "ReconciliationFailed"

	// ROSAControlPlaneDeletionFailedReason used to report failures while deleting ROSAControlPlane.
	ROSAControlPlaneDeletionFailedReason = "DeletionFailed"

	// ROSAControlPlaneInvalidConfigurationReason used to report invalid user input.
	ROSAControlPlaneInvalidConfigurationReason = "InvalidConfiguration"

	// ROSAControlPlaneRoleConfigNotReadyReason used to report when referenced RosaRoleConfig is not ready.
	ROSAControlPlaneRoleConfigNotReadyReason = "ROSARoleConfigNotReady"

	// ROSAControlPlaneRoleConfigNotFoundReason used to report when referenced RosaRoleConfig is not found.
	ROSAControlPlaneRoleConfigNotFoundReason = "ROSARoleConfigNotFound"
)
