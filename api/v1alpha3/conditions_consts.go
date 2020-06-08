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

import clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"

const (
	WorkerConditionCount             = 2
	ControlPlaneConditionCount       = 3
	ClusterConditionCount            = 2
	ClusterConditionWithBastionCount = 3
)

// AWSCluster Conditions and Reasons
const (
	NetworkInfrastructureReadyCondition = "NetworkInfrastructureReady"
	NetworkInfrastructureFailedReason   = "NetworkInfrastructureFailed"
)

const (
	BastionHostReadyCondition = "BastionHostReady"
	BastionHostFailedReason   = "BastionHostFailed"
)

const (
	LoadBalancerReadyCondition = "LoadBalancerReady"
	LoadBalancerFailedReason   = "LoadBalancerFailed"
)

// AWSMachine Conditions and Reasons
const (
	// InstanceReadyCondition reports on current status of the EC2 instance. Ready indicates the instance is in a Running state.
	InstanceReadyCondition clusterv1.ConditionType = "InstanceReady"

	InstanceNotFoundReason   = "InstanceNotFound"
	InstanceTerminatedReason = "InstanceTerminated"
	InstanceStoppedReason    = "InstanceStopped"
	InstanceNotReadyReason   = "InstanceNotReady"
)

const (
	// SecurityGroupsReadyCondition indicates the security groups are up to date on the AWSMachine.
	SecurityGroupsReadyCondition clusterv1.ConditionType = "SecurityGroupsReady"

	SecurityGroupsFailedReason = "SecurityGroupsSyncFailed"
)

const (
	// Only applicable to control plane machines. ELBAttachedCondition will report true when a control plane is successfully registered with an ELB
	// When set to false, severity can be an Error if the subnet is not found or unavailable in the instance's AZ
	ELBAttachedCondition clusterv1.ConditionType = "ELBAttached"

	ELBAttachFailedReason = "ELBAttachFailed"
	ELBDetachFailedReason = "ELBDetachFailed"
)
