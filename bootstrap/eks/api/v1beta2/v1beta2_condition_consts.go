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

// EKSConfig v1beta2 condition types.
const (
	// EKSConfigReadyCondition defines the Ready condition type that summarizes the operational state of an EKSConfig.
	EKSConfigReadyCondition = clusterv1.ReadyCondition

	// EKSConfigDataSecretAvailableCondition reports on the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	EKSConfigDataSecretAvailableCondition = "DataSecretAvailable"
)

// EKSConfig v1beta2 reason constants.
const (
	// EKSConfigReadyReason indicates the EKSConfig is ready.
	EKSConfigReadyReason = clusterv1.ReadyReason

	// EKSConfigNotReadyReason indicates the EKSConfig is not ready.
	EKSConfigNotReadyReason = clusterv1.NotReadyReason

	// EKSConfigDeletingReason indicates the EKSConfig is being deleted.
	EKSConfigDeletingReason = clusterv1.DeletingReason

	// EKSConfigDataSecretGenerationFailedReason indicates an error while generating a data secret;
	// those kind of errors are usually due to misconfigurations and user intervention is required to get them fixed.
	EKSConfigDataSecretGenerationFailedReason = "DataSecretGenerationFailed"

	// EKSConfigWaitingForClusterInfrastructureReason indicates waiting for cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the EKSConfig controller ensures this pre-condition is satisfied.
	EKSConfigWaitingForClusterInfrastructureReason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// EKSConfigWaitingForControlPlaneInitializationReason indicates waiting for the control plane to be initialized.
	//
	// NOTE: This is a pre-condition for starting to create machines;
	// the EKSConfig controller ensures this pre-condition is satisfied.
	EKSConfigWaitingForControlPlaneInitializationReason = clusterv1.WaitingForControlPlaneInitializedReason
)

// NodeadmConfig v1beta2 condition types.
const (
	// NodeadmConfigReadyCondition defines the Ready condition type that summarizes the operational state of a NodeadmConfig.
	NodeadmConfigReadyCondition = clusterv1.ReadyCondition

	// NodeadmConfigDataSecretAvailableCondition reports on the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	NodeadmConfigDataSecretAvailableCondition = "DataSecretAvailable"
)

// NodeadmConfig v1beta2 reason constants.
const (
	// NodeadmConfigReadyReason indicates the NodeadmConfig is ready.
	NodeadmConfigReadyReason = clusterv1.ReadyReason

	// NodeadmConfigNotReadyReason indicates the NodeadmConfig is not ready.
	NodeadmConfigNotReadyReason = clusterv1.NotReadyReason

	// NodeadmConfigDeletingReason indicates the NodeadmConfig is being deleted.
	NodeadmConfigDeletingReason = clusterv1.DeletingReason

	// NodeadmConfigDataSecretGenerationFailedReason indicates an error while generating a data secret;
	// those kind of errors are usually due to misconfigurations and user intervention is required to get them fixed.
	NodeadmConfigDataSecretGenerationFailedReason = "DataSecretGenerationFailed"

	// NodeadmConfigWaitingForClusterInfrastructureReason indicates waiting for cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the NodeadmConfig controller ensures this pre-condition is satisfied.
	NodeadmConfigWaitingForClusterInfrastructureReason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// NodeadmConfigWaitingForControlPlaneInitializationReason indicates waiting for the control plane to be initialized.
	//
	// NOTE: This is a pre-condition for starting to create machines;
	// the NodeadmConfig controller ensures this pre-condition is satisfied.
	NodeadmConfigWaitingForControlPlaneInitializationReason = clusterv1.WaitingForControlPlaneInitializedReason
)
