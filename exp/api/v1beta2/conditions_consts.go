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

const (
	// ASGReadyCondition reports on current status of the autoscaling group. Ready indicates the group is provisioned.
	ASGReadyCondition = "ASGReady"
	// ASGNotFoundReason used when the autoscaling group couldn't be retrieved.
	ASGNotFoundReason = "ASGNotFound"
	// ASGProvisionFailedReason used for failures during autoscaling group provisioning.
	ASGProvisionFailedReason = "ASGProvisionFailed"
	// ASGDeletionInProgress ASG is in a deletion in progress state.
	ASGDeletionInProgress = "ASGDeletionInProgress"

	// LaunchTemplateReadyCondition represents the status of an AWSMachinePool's associated Launch Template.
	LaunchTemplateReadyCondition = "LaunchTemplateReady"
	// LaunchTemplateNotFoundReason is used when an associated Launch Template can't be found.
	LaunchTemplateNotFoundReason = "LaunchTemplateNotFound"
	// LaunchTemplateCreateFailedReason used for failures during Launch Template creation.
	LaunchTemplateCreateFailedReason = "LaunchTemplateCreateFailed"
	// LaunchTemplateReconcileFailedReason used for failures during Launch Template reconciliation.
	LaunchTemplateReconcileFailedReason = "LaunchTemplateReconcileFailed"

	// PreLaunchTemplateUpdateCheckCondition reports if all prerequisite are met for launch template update.
	PreLaunchTemplateUpdateCheckCondition = "PreLaunchTemplateUpdateCheckSuccess"
	// PostLaunchTemplateUpdateOperationCondition reports on successfully completes post launch template update operation.
	PostLaunchTemplateUpdateOperationCondition = "PostLaunchTemplateUpdateOperationSuccess"

	// PreLaunchTemplateUpdateCheckFailedReason used to report when not all prerequisite are met for launch template update.
	PreLaunchTemplateUpdateCheckFailedReason = "PreLaunchTemplateUpdateCheckFailed"
	// PostLaunchTemplateUpdateOperationFailedReason used to report when post launch template update operation failed.
	PostLaunchTemplateUpdateOperationFailedReason = "PostLaunchTemplateUpdateOperationFailed"

	// InstanceRefreshStartedCondition reports on successfully starting instance refresh.
	InstanceRefreshStartedCondition = "InstanceRefreshStarted"
	// InstanceRefreshNotReadyReason used to report instance refresh is not initiated.
	// If there are instance refreshes that are in progress, then a new instance refresh request will fail.
	InstanceRefreshNotReadyReason = "InstanceRefreshNotReady"
	// InstanceRefreshFailedReason used to report when there instance refresh is not initiated.
	InstanceRefreshFailedReason = "InstanceRefreshFailed"

	// AWSMachineCreationFailed reports if creating AWSMachines to represent ASG (machine pool) machines failed.
	AWSMachineCreationFailed = "AWSMachineCreationFailed"
	// AWSMachineDeletionFailed reports if deleting AWSMachines failed.
	AWSMachineDeletionFailed = "AWSMachineDeletionFailed"
	// LifecycleHookReadyCondition reports on the status of the lifecycle hook.
	LifecycleHookReadyCondition = "LifecycleHookReady"
	// LifecycleHookCreationFailedReason used for failures during lifecycle hook creation.
	LifecycleHookCreationFailedReason = "LifecycleHookCreationFailed"
	// LifecycleHookUpdateFailedReason used for failures during lifecycle hook update.
	LifecycleHookUpdateFailedReason = "LifecycleHookUpdateFailed"
	// LifecycleHookDeletionFailedReason used for failures during lifecycle hook deletion.
	LifecycleHookDeletionFailedReason = "LifecycleHookDeletionFailed"
)

const (
	// EKSNodegroupReadyCondition condition reports on the successful reconciliation of eks control plane.
	EKSNodegroupReadyCondition = "EKSNodegroupReady"
	// EKSNodegroupReconciliationFailedReason used to report failures while reconciling EKS control plane.
	EKSNodegroupReconciliationFailedReason = "EKSNodegroupReconciliationFailed"
	// WaitingForEKSControlPlaneReason used when the machine pool is waiting for
	// EKS control plane infrastructure to be ready before proceeding.
	WaitingForEKSControlPlaneReason = "WaitingForEKSControlPlane"
)

const (
	// EKSFargateProfileReadyCondition condition reports on the successful reconciliation of eks control plane.
	EKSFargateProfileReadyCondition = "EKSFargateProfileReady"
	// EKSFargateCreatingCondition condition reports on whether the fargate
	// profile is creating.
	EKSFargateCreatingCondition = "EKSFargateCreating"
	// EKSFargateDeletingCondition used to report that the profile is deleting.
	EKSFargateDeletingCondition = "EKSFargateDeleting"
	// EKSFargateReconciliationFailedReason used to report failures while reconciling EKS control plane.
	EKSFargateReconciliationFailedReason = "EKSFargateReconciliationFailed"
	// EKSFargateDeletingReason used when the profile is deleting.
	EKSFargateDeletingReason = "Deleting"
	// EKSFargateCreatingReason used when the profile is creating.
	EKSFargateCreatingReason = "Creating"
	// EKSFargateCreatedReason used when the profile is created.
	EKSFargateCreatedReason = "Created"
	// EKSFargateDeletedReason used when the profile is deleted.
	EKSFargateDeletedReason = "Deleted"
	// EKSFargateFailedReason used when the profile failed.
	EKSFargateFailedReason = "Failed"
)

const (
	// IAMNodegroupRolesReadyCondition condition reports on the successful
	// reconciliation of EKS nodegroup iam roles.
	IAMNodegroupRolesReadyCondition = "IAMNodegroupRolesReady"
	// IAMNodegroupRolesReconciliationFailedReason used to report failures while
	// reconciling EKS nodegroup iam roles.
	IAMNodegroupRolesReconciliationFailedReason = "IAMNodegroupRolesReconciliationFailed"
	// IAMFargateRolesReadyCondition condition reports on the successful
	// reconciliation of EKS nodegroup iam roles.
	IAMFargateRolesReadyCondition = "IAMFargateRolesReady"
	// IAMFargateRolesReconciliationFailedReason used to report failures while
	// reconciling EKS nodegroup iam roles.
	IAMFargateRolesReconciliationFailedReason = "IAMFargateRolesReconciliationFailed"
)

const (
	// RosaMachinePoolReadyCondition condition reports on the successful reconciliation of rosa machinepool.
	RosaMachinePoolReadyCondition = "RosaMachinePoolReady"
	// RosaMachinePoolUpgradingCondition condition reports whether ROSAMachinePool is upgrading or not.
	RosaMachinePoolUpgradingCondition = "RosaMachinePoolUpgrading"

	// WaitingForRosaControlPlaneReason used when the machine pool is waiting for
	// ROSA control plane infrastructure to be ready before proceeding.
	WaitingForRosaControlPlaneReason = "WaitingForRosaControlPlane"

	// RosaMachinePoolReconciliationFailedReason used to report failures while reconciling ROSAMachinePool.
	RosaMachinePoolReconciliationFailedReason = "ReconciliationFailed"
)
