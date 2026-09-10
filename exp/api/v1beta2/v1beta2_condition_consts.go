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

// AWSMachinePool v1beta2 condition types.
const (
	// AWSMachinePoolReadyCondition defines the Ready condition type that summarizes the operational state of an AWSMachinePool.
	AWSMachinePoolReadyCondition = clusterv1.ReadyCondition

	// AWSMachinePoolASGReadyCondition reports on current status of the autoscaling group. Ready indicates the group is provisioned.
	AWSMachinePoolASGReadyCondition = "ASGReady"

	// AWSMachinePoolLaunchTemplateReadyCondition reports on the status of an AWSMachinePool's associated launch template.
	AWSMachinePoolLaunchTemplateReadyCondition = "LaunchTemplateReady"

	// AWSMachinePoolPreLaunchTemplateUpdateCheckCondition reports if all prerequisites are met for launch template update.
	AWSMachinePoolPreLaunchTemplateUpdateCheckCondition = "PreLaunchTemplateUpdateCheckSuccess"

	// AWSMachinePoolPostLaunchTemplateUpdateOperationCondition reports on successfully completed post launch template update operation.
	AWSMachinePoolPostLaunchTemplateUpdateOperationCondition = "PostLaunchTemplateUpdateOperationSuccess"

	// AWSMachinePoolInstanceRefreshStartedCondition reports on successfully starting instance refresh.
	AWSMachinePoolInstanceRefreshStartedCondition = "InstanceRefreshStarted"

	// AWSMachinePoolLifecycleHookReadyCondition reports on the status of the lifecycle hook.
	AWSMachinePoolLifecycleHookReadyCondition = "LifecycleHookReady"
)

// AWSMachinePool v1beta2 reason constants.
const (
	// AWSMachinePoolReadyReason indicates the AWSMachinePool is ready.
	AWSMachinePoolReadyReason = clusterv1.ReadyReason

	// AWSMachinePoolNotReadyReason indicates the AWSMachinePool is not ready.
	AWSMachinePoolNotReadyReason = clusterv1.NotReadyReason

	// AWSMachinePoolDeletingReason indicates the AWSMachinePool is being deleted.
	AWSMachinePoolDeletingReason = clusterv1.DeletingReason

	// AWSMachinePoolASGNotFoundReason used when the autoscaling group couldn't be retrieved.
	AWSMachinePoolASGNotFoundReason = "ASGNotFound"

	// AWSMachinePoolASGProvisionFailedReason used for failures during autoscaling group provisioning.
	AWSMachinePoolASGProvisionFailedReason = "ASGProvisionFailed"

	// AWSMachinePoolASGDeletionInProgressReason used when the autoscaling group is in a deletion in progress state.
	AWSMachinePoolASGDeletionInProgressReason = "ASGDeletionInProgress"

	// AWSMachinePoolLaunchTemplateNotFoundReason used when an associated launch template can't be found.
	AWSMachinePoolLaunchTemplateNotFoundReason = "LaunchTemplateNotFound"

	// AWSMachinePoolLaunchTemplateCreateFailedReason used for failures during launch template creation.
	AWSMachinePoolLaunchTemplateCreateFailedReason = "LaunchTemplateCreateFailed"

	// AWSMachinePoolLaunchTemplateReconcileFailedReason used for failures during launch template reconciliation.
	AWSMachinePoolLaunchTemplateReconcileFailedReason = "LaunchTemplateReconcileFailed"

	// AWSMachinePoolLaunchTemplateNitroEnclaveEdgeZoneReason used when enclaveOptions is enabled but the pool
	// targets a Local Zone or Wavelength Zone, which does not support Nitro Enclaves.
	AWSMachinePoolLaunchTemplateNitroEnclaveEdgeZoneReason = "NitroEnclaveEdgeZoneUnsupported"

	// AWSMachinePoolPreLaunchTemplateUpdateCheckFailedReason used to report when not all prerequisites are met for launch template update.
	AWSMachinePoolPreLaunchTemplateUpdateCheckFailedReason = "PreLaunchTemplateUpdateCheckFailed"

	// AWSMachinePoolPostLaunchTemplateUpdateOperationFailedReason used to report when post launch template update operation failed.
	AWSMachinePoolPostLaunchTemplateUpdateOperationFailedReason = "PostLaunchTemplateUpdateOperationFailed"

	// AWSMachinePoolInstanceRefreshNotReadyReason used to report instance refresh is not initiated.
	// If there are instance refreshes that are in progress, then a new instance refresh request will fail.
	AWSMachinePoolInstanceRefreshNotReadyReason = "InstanceRefreshNotReady"

	// AWSMachinePoolInstanceRefreshFailedReason used to report when instance refresh is not initiated.
	AWSMachinePoolInstanceRefreshFailedReason = "InstanceRefreshFailed"

	// AWSMachinePoolMachineCreationFailedReason used to report when creating AWSMachines to represent ASG machines failed.
	AWSMachinePoolMachineCreationFailedReason = "AWSMachineCreationFailed"

	// AWSMachinePoolMachineDeletionFailedReason used to report when deleting AWSMachines failed.
	AWSMachinePoolMachineDeletionFailedReason = "AWSMachineDeletionFailed"

	// AWSMachinePoolLifecycleHookCreationFailedReason used for failures during lifecycle hook creation.
	AWSMachinePoolLifecycleHookCreationFailedReason = "LifecycleHookCreationFailed"

	// AWSMachinePoolLifecycleHookUpdateFailedReason used for failures during lifecycle hook update.
	AWSMachinePoolLifecycleHookUpdateFailedReason = "LifecycleHookUpdateFailed"

	// AWSMachinePoolLifecycleHookDeletionFailedReason used for failures during lifecycle hook deletion.
	AWSMachinePoolLifecycleHookDeletionFailedReason = "LifecycleHookDeletionFailed"
)

// AWSManagedMachinePool v1beta2 condition types.
const (
	// AWSManagedMachinePoolReadyCondition defines the Ready condition type that summarizes the operational state of an AWSManagedMachinePool.
	AWSManagedMachinePoolReadyCondition = clusterv1.ReadyCondition

	// AWSManagedMachinePoolEKSNodegroupReadyCondition reports on the successful reconciliation of the EKS node group.
	AWSManagedMachinePoolEKSNodegroupReadyCondition = "EKSNodegroupReady"

	// AWSManagedMachinePoolIAMNodegroupRolesReadyCondition reports on the successful reconciliation of EKS node group IAM roles.
	AWSManagedMachinePoolIAMNodegroupRolesReadyCondition = "IAMNodegroupRolesReady"

	// AWSManagedMachinePoolLaunchTemplateReadyCondition reports on the status of the associated launch template.
	AWSManagedMachinePoolLaunchTemplateReadyCondition = "LaunchTemplateReady"
)

// AWSManagedMachinePool v1beta2 reason constants.
const (
	// AWSManagedMachinePoolReadyReason indicates the AWSManagedMachinePool is ready.
	AWSManagedMachinePoolReadyReason = clusterv1.ReadyReason

	// AWSManagedMachinePoolNotReadyReason indicates the AWSManagedMachinePool is not ready.
	AWSManagedMachinePoolNotReadyReason = clusterv1.NotReadyReason

	// AWSManagedMachinePoolDeletingReason indicates the AWSManagedMachinePool is being deleted.
	AWSManagedMachinePoolDeletingReason = clusterv1.DeletingReason

	// AWSManagedMachinePoolEKSNodegroupReconciliationFailedReason used to report failures while reconciling EKS control plane.
	AWSManagedMachinePoolEKSNodegroupReconciliationFailedReason = "EKSNodegroupReconciliationFailed"

	// AWSManagedMachinePoolWaitingForEKSControlPlaneReason used when the machine pool is waiting for
	// EKS control plane infrastructure to be ready before proceeding.
	AWSManagedMachinePoolWaitingForEKSControlPlaneReason = "WaitingForEKSControlPlane"

	// AWSManagedMachinePoolIAMNodegroupRolesReconciliationFailedReason used to report failures while reconciling EKS nodegroup IAM roles.
	AWSManagedMachinePoolIAMNodegroupRolesReconciliationFailedReason = "IAMNodegroupRolesReconciliationFailed"

	// AWSManagedMachinePoolLaunchTemplateNotFoundReason used when an associated launch template can't be found.
	AWSManagedMachinePoolLaunchTemplateNotFoundReason = "LaunchTemplateNotFound"

	// AWSManagedMachinePoolLaunchTemplateCreateFailedReason used for failures during launch template creation.
	AWSManagedMachinePoolLaunchTemplateCreateFailedReason = "LaunchTemplateCreateFailed"

	// AWSManagedMachinePoolLaunchTemplateReconcileFailedReason used for failures during launch template reconciliation.
	AWSManagedMachinePoolLaunchTemplateReconcileFailedReason = "LaunchTemplateReconcileFailed"
)

// AWSFargateProfile v1beta2 condition types.
const (
	// AWSFargateProfileReadyCondition defines the Ready condition type that summarizes the operational state of an AWSFargateProfile.
	AWSFargateProfileReadyCondition = clusterv1.ReadyCondition

	// AWSFargateProfileEKSFargateProfileReadyCondition reports on the successful reconciliation of the EKS Fargate profile.
	AWSFargateProfileEKSFargateProfileReadyCondition = "EKSFargateProfileReady"

	// AWSFargateProfileEKSFargateCreatingCondition reports whether the Fargate profile is being created.
	AWSFargateProfileEKSFargateCreatingCondition = "EKSFargateCreating"

	// AWSFargateProfileEKSFargateDeletingCondition used to report that the profile is deleting.
	AWSFargateProfileEKSFargateDeletingCondition = "EKSFargateDeleting"

	// AWSFargateProfileIAMFargateRolesReadyCondition reports on the successful reconciliation of Fargate IAM roles.
	AWSFargateProfileIAMFargateRolesReadyCondition = "IAMFargateRolesReady"
)

// AWSFargateProfile v1beta2 reason constants.
const (
	// AWSFargateProfileReadyReason indicates the AWSFargateProfile is ready.
	AWSFargateProfileReadyReason = clusterv1.ReadyReason

	// AWSFargateProfileNotReadyReason indicates the AWSFargateProfile is not ready.
	AWSFargateProfileNotReadyReason = clusterv1.NotReadyReason

	// AWSFargateProfileDeletingReason indicates the AWSFargateProfile is being deleted.
	AWSFargateProfileDeletingReason = clusterv1.DeletingReason

	// AWSFargateProfileReconciliationFailedReason used to report failures while reconciling EKS Fargate profile.
	AWSFargateProfileReconciliationFailedReason = "EKSFargateReconciliationFailed"

	// AWSFargateProfileCreatingReason used when the profile is creating.
	AWSFargateProfileCreatingReason = "Creating"

	// AWSFargateProfileCreatedReason used when the profile is created.
	AWSFargateProfileCreatedReason = "Created"

	// AWSFargateProfileFailedReason used when the profile failed.
	AWSFargateProfileFailedReason = "Failed"

	// AWSFargateProfileDeletedReason used when the profile is deleted.
	AWSFargateProfileDeletedReason = "Deleted"

	// AWSFargateProfileIAMFargateRolesReconciliationFailedReason used to report failures while reconciling Fargate IAM roles.
	AWSFargateProfileIAMFargateRolesReconciliationFailedReason = "IAMFargateRolesReconciliationFailed"
)

// ROSAMachinePool v1beta2 condition types.
const (
	// ROSAMachinePoolReadyCondition defines the Ready condition type that summarizes the operational state of a ROSAMachinePool.
	ROSAMachinePoolReadyCondition = clusterv1.ReadyCondition

	// ROSAMachinePoolNodePoolReadyCondition reports on the successful reconciliation of rosa machinepool.
	ROSAMachinePoolNodePoolReadyCondition = "RosaMachinePoolReady"

	// ROSAMachinePoolUpgradingCondition reports whether ROSAMachinePool is upgrading or not.
	ROSAMachinePoolUpgradingCondition = "RosaMachinePoolUpgrading"
)

// ROSAMachinePool v1beta2 reason constants.
const (
	// ROSAMachinePoolReadyReason indicates the ROSAMachinePool is ready.
	ROSAMachinePoolReadyReason = clusterv1.ReadyReason

	// ROSAMachinePoolNotReadyReason indicates the ROSAMachinePool is not ready.
	ROSAMachinePoolNotReadyReason = clusterv1.NotReadyReason

	// ROSAMachinePoolDeletingReason indicates the ROSAMachinePool is being deleted.
	ROSAMachinePoolDeletingReason = clusterv1.DeletingReason

	// ROSAMachinePoolWaitingForRosaControlPlaneReason used when the machine pool is waiting for
	// ROSA control plane infrastructure to be ready before proceeding.
	ROSAMachinePoolWaitingForRosaControlPlaneReason = "WaitingForRosaControlPlane"

	// ROSAMachinePoolReconciliationFailedReason used to report failures while reconciling ROSAMachinePool.
	ROSAMachinePoolReconciliationFailedReason = "ReconciliationFailed"
)

// ROSACluster v1beta2 condition types.
const (
	// ROSAClusterReadyCondition defines the Ready condition type that summarizes the operational state of a ROSACluster.
	ROSAClusterReadyCondition = clusterv1.ReadyCondition
)

// ROSACluster v1beta2 reason constants.
const (
	// ROSAClusterReadyReason indicates the ROSACluster is ready.
	ROSAClusterReadyReason = clusterv1.ReadyReason

	// ROSAClusterNotReadyReason indicates the ROSACluster is not ready.
	ROSAClusterNotReadyReason = clusterv1.NotReadyReason
)

// ROSANetwork v1beta2 condition types.
const (
	// ROSANetworkReadyCondition defines the Ready condition type that summarizes the operational state of a ROSANetwork.
	ROSANetworkReadyCondition = clusterv1.ReadyCondition

	// ROSANetworkNetworkReadyCondition reports on the successful reconciliation of ROSANetwork.
	ROSANetworkNetworkReadyCondition = "ROSANetworkReady"
)

// ROSANetwork v1beta2 reason constants.
const (
	// ROSANetworkReadyReason indicates the ROSANetwork is ready.
	ROSANetworkReadyReason = clusterv1.ReadyReason

	// ROSANetworkNotReadyReason indicates the ROSANetwork is not ready.
	ROSANetworkNotReadyReason = clusterv1.NotReadyReason

	// ROSANetworkCreatingReason used when ROSANetwork is being created.
	ROSANetworkCreatingReason = "Creating"

	// ROSANetworkCreatedReason used when ROSANetwork is created.
	ROSANetworkCreatedReason = "Created"

	// ROSANetworkFailedReason used when ROSANetwork creation failed.
	ROSANetworkFailedReason = "Failed"

	// ROSANetworkDeletingReason indicates the ROSANetwork is being deleted.
	ROSANetworkDeletingReason = clusterv1.DeletingReason

	// ROSANetworkDeletionFailedReason used to report failures while deleting ROSANetwork.
	ROSANetworkDeletionFailedReason = "DeletionFailed"
)

// ROSARoleConfig v1beta2 condition types.
const (
	// ROSARoleConfigReadyCondition defines the Ready condition type that summarizes the operational state of a ROSARoleConfig.
	ROSARoleConfigReadyCondition = clusterv1.ReadyCondition

	// ROSARoleConfigRoleConfigReadyCondition reports on the successful reconciliation of ROSARoleConfig.
	ROSARoleConfigRoleConfigReadyCondition = "RosaRoleConfigReady"
)

// ROSARoleConfig v1beta2 reason constants.
const (
	// ROSARoleConfigReadyReason indicates the ROSARoleConfig is ready.
	ROSARoleConfigReadyReason = clusterv1.ReadyReason

	// ROSARoleConfigNotReadyReason indicates the ROSARoleConfig is not ready.
	ROSARoleConfigNotReadyReason = clusterv1.NotReadyReason

	// ROSARoleConfigDeletingReason indicates the ROSARoleConfig is being deleted.
	ROSARoleConfigDeletingReason = clusterv1.DeletingReason

	// ROSARoleConfigReconciliationFailedReason used to report reconciliation failures.
	ROSARoleConfigReconciliationFailedReason = "ReconciliationFailed"

	// ROSARoleConfigDeletionFailedReason used to report failures while deleting ROSARoleConfig.
	ROSARoleConfigDeletionFailedReason = "DeletionFailed"

	// ROSARoleConfigDeletionStartedReason used to indicate that the deletion of ROSARoleConfig has started.
	ROSARoleConfigDeletionStartedReason = "DeletionStarted"

	// ROSARoleConfigCreatedReason used to indicate that the ROSARoleConfig has been created.
	ROSARoleConfigCreatedReason = "Created"
)

// ROSAOCMRoleConfig v1beta2 condition types.
const (
	// ROSAOCMRoleConfigReadyCondition defines the Ready condition type that summarizes the operational state of a ROSAOCMRoleConfig.
	ROSAOCMRoleConfigReadyCondition = clusterv1.ReadyCondition

	// ROSAOCMRoleConfigOCMRoleConfigReadyCondition reports on the successful reconciliation of ROSAOCMRoleConfig.
	ROSAOCMRoleConfigOCMRoleConfigReadyCondition = "ROSAOCMRoleConfigReady"
)

// ROSAOCMRoleConfig v1beta2 reason constants.
const (
	// ROSAOCMRoleConfigReadyReason indicates the ROSAOCMRoleConfig is ready.
	ROSAOCMRoleConfigReadyReason = clusterv1.ReadyReason

	// ROSAOCMRoleConfigNotReadyReason indicates the ROSAOCMRoleConfig is not ready.
	ROSAOCMRoleConfigNotReadyReason = clusterv1.NotReadyReason

	// ROSAOCMRoleConfigDeletingReason indicates the ROSAOCMRoleConfig is being deleted.
	ROSAOCMRoleConfigDeletingReason = clusterv1.DeletingReason

	// ROSAOCMRoleConfigDeletionStartedReason used to indicate that the deletion of ROSAOCMRoleConfig has started.
	ROSAOCMRoleConfigDeletionStartedReason = "DeletionStarted"

	// ROSAOCMRoleConfigReconciliationFailedReason used to report reconciliation failures.
	ROSAOCMRoleConfigReconciliationFailedReason = "ReconciliationFailed"

	// ROSAOCMRoleConfigDeletionFailedReason used to report failures while deleting ROSAOCMRoleConfig.
	ROSAOCMRoleConfigDeletionFailedReason = "DeletionFailed"

	// ROSAOCMRoleConfigCreatedReason used to indicate that the ROSAOCMRoleConfig has been created.
	ROSAOCMRoleConfigCreatedReason = "Created"

	// ROSAOCMRoleConfigLinkedReason used to indicate that the OCM role has been linked to the organization.
	ROSAOCMRoleConfigLinkedReason = "Linked"
)
