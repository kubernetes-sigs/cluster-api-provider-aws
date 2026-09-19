/*
Copyright 2022 The Kubernetes Authors.

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

import clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"

const (
	// ROSAControlPlaneReadyV1Beta1Condition condition reports on the successful reconciliation of ROSAControlPlane.
	// Use ROSAControlPlaneReadyCondition from v1beta2_condition_consts.go instead.
	ROSAControlPlaneReadyV1Beta1Condition clusterv1beta1.ConditionType = "ROSAControlPlaneReady"

	// ROSAControlPlaneValidV1Beta1Condition condition reports whether ROSAControlPlane configuration is valid.
	// Use ROSAControlPlaneValidCondition from v1beta2_condition_consts.go instead.
	ROSAControlPlaneValidV1Beta1Condition clusterv1beta1.ConditionType = "ROSAControlPlaneValid"

	// ROSAControlPlaneUpgradingV1Beta1Condition condition reports whether ROSAControlPlane is upgrading or not.
	// Use ROSAControlPlaneUpgradingCondition from v1beta2_condition_consts.go instead.
	ROSAControlPlaneUpgradingV1Beta1Condition clusterv1beta1.ConditionType = "ROSAControlPlaneUpgrading"

	// ExternalAuthConfiguredCondition condition reports whether external auth has beed correctly configured.
	ExternalAuthConfiguredCondition clusterv1beta1.ConditionType = "ExternalAuthConfigured"

	// ROSARoleConfigReadyCondition condition reports whether the referenced RosaRoleConfig is ready.
	ROSARoleConfigReadyCondition clusterv1beta1.ConditionType = "ROSARoleConfigReady"

	// ReconciliationFailedReason used to report reconciliation failures.
	ReconciliationFailedReason = "ReconciliationFailed"

	// ROSAControlPlaneDeletionFailedV1Beta1Reason used to report failures while deleting ROSAControlPlane.
	// Use ROSAControlPlaneDeletionFailedReason from v1beta2_condition_consts.go instead.
	ROSAControlPlaneDeletionFailedV1Beta1Reason = "DeletionFailed"

	// ROSAControlPlaneInvalidConfigurationV1Beta1Reason used to report invalid user input.
	// Use ROSAControlPlaneInvalidConfigurationReason from v1beta2_condition_consts.go instead.
	ROSAControlPlaneInvalidConfigurationV1Beta1Reason = "InvalidConfiguration"

	// ROSARoleConfigNotReadyReason used to report when referenced RosaRoleConfig is not ready.
	ROSARoleConfigNotReadyReason = "ROSARoleConfigNotReady"

	// ROSARoleConfigNotFoundReason used to report when referenced RosaRoleConfig is not found.
	ROSARoleConfigNotFoundReason = "ROSARoleConfigNotFound"
)
