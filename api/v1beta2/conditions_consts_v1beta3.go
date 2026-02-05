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

// AWSCluster's v1beta3 conditions and corresponding reasons.
// These will be used with the V1Beta3 API version.

// AWSCluster's Ready condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterReadyCondition is true if the AWSCluster's deletionTimestamp is not set, and all
	// infrastructure conditions (VpcReady, SubnetsReady, etc.) are true.
	AWSClusterReadyCondition = clusterv1beta1.ReadyV1Beta2Condition

	// AWSClusterReadyReason surfaces when the AWSCluster readiness criteria is met.
	AWSClusterReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterNotReadyReason surfaces when the AWSCluster readiness criteria is not met.
	AWSClusterNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterReadyUnknownReason surfaces when at least one AWSCluster readiness criteria is unknown
	// and no AWSCluster readiness criteria is not met.
	AWSClusterReadyUnknownReason = clusterv1beta1.ReadyUnknownV1Beta2Reason
)

// AWSCluster's VpcReady condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterVpcReadyCondition documents the status of the VPC for an AWSCluster.
	AWSClusterVpcReadyCondition = "VpcReady"

	// AWSClusterVpcReadyReason surfaces when the VPC for an AWSCluster is ready.
	AWSClusterVpcReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterVpcNotReadyReason surfaces when the VPC for an AWSCluster is not ready.
	AWSClusterVpcNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterVpcDeletingReason surfaces when the VPC for an AWSCluster is being deleted.
	AWSClusterVpcDeletingReason = clusterv1beta1.DeletingV1Beta2Reason
)

// AWSCluster's SubnetsReady condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterSubnetsReadyCondition documents the status of subnets for an AWSCluster.
	AWSClusterSubnetsReadyCondition = "SubnetsReady"

	// AWSClusterSubnetsReadyReason surfaces when the subnets for an AWSCluster are ready.
	AWSClusterSubnetsReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterSubnetsNotReadyReason surfaces when the subnets for an AWSCluster are not ready.
	AWSClusterSubnetsNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterSubnetsDeletingReason surfaces when the subnets for an AWSCluster are being deleted.
	AWSClusterSubnetsDeletingReason = clusterv1beta1.DeletingV1Beta2Reason
)

// AWSCluster's LoadBalancerReady condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterLoadBalancerReadyCondition documents the status of the load balancer for an AWSCluster.
	AWSClusterLoadBalancerReadyCondition = "LoadBalancerReady"

	// AWSClusterLoadBalancerReadyReason surfaces when the load balancer for an AWSCluster is ready.
	AWSClusterLoadBalancerReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterLoadBalancerNotReadyReason surfaces when the load balancer for an AWSCluster is not ready.
	AWSClusterLoadBalancerNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterLoadBalancerDeletingReason surfaces when the load balancer for an AWSCluster is being deleted.
	AWSClusterLoadBalancerDeletingReason = clusterv1beta1.DeletingV1Beta2Reason
)

// AWSCluster's ClusterSecurityGroupsReady condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterSecurityGroupsReadyCondition documents the status of security groups for an AWSCluster.
	AWSClusterSecurityGroupsReadyCondition = "ClusterSecurityGroupsReady"

	// AWSClusterSecurityGroupsReadyReason surfaces when the security groups for an AWSCluster are ready.
	AWSClusterSecurityGroupsReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterSecurityGroupsNotReadyReason surfaces when the security groups for an AWSCluster are not ready.
	AWSClusterSecurityGroupsNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterSecurityGroupsDeletingReason surfaces when the security groups for an AWSCluster are being deleted.
	AWSClusterSecurityGroupsDeletingReason = clusterv1beta1.DeletingV1Beta2Reason
)

// AWSCluster's BastionHostReady condition and corresponding reasons that will be used in v1Beta3 API version.
const (
	// AWSClusterBastionHostReadyCondition documents the status of the bastion host for an AWSCluster.
	AWSClusterBastionHostReadyCondition = "BastionHostReady"

	// AWSClusterBastionHostReadyReason surfaces when the bastion host for an AWSCluster is ready.
	AWSClusterBastionHostReadyReason = clusterv1beta1.ReadyV1Beta2Reason

	// AWSClusterBastionHostNotReadyReason surfaces when the bastion host for an AWSCluster is not ready.
	AWSClusterBastionHostNotReadyReason = clusterv1beta1.NotReadyV1Beta2Reason

	// AWSClusterBastionHostDeletingReason surfaces when the bastion host for an AWSCluster is being deleted.
	AWSClusterBastionHostDeletingReason = clusterv1beta1.DeletingV1Beta2Reason
)
