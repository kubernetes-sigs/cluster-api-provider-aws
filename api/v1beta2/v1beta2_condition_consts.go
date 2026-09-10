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

// AWSCluster v1beta2 condition types.
const (
	// AWSClusterReadyCondition defines the Ready condition type that summarizes the operational state of an AWSCluster.
	AWSClusterReadyCondition = clusterv1.ReadyCondition

	// AWSClusterVpcReadyCondition reports on the successful reconciliation of a VPC.
	AWSClusterVpcReadyCondition = "VpcReady"

	// AWSClusterSubnetsReadyCondition reports on the successful reconciliation of subnets.
	AWSClusterSubnetsReadyCondition = "SubnetsReady"

	// AWSClusterInternetGatewayReadyCondition reports on the successful reconciliation of internet gateways.
	// Only applicable to managed clusters.
	AWSClusterInternetGatewayReadyCondition = "InternetGatewayReady"

	// AWSClusterEgressOnlyInternetGatewayReadyCondition reports on the successful reconciliation of egress only internet gateways.
	// Only applicable to managed clusters.
	AWSClusterEgressOnlyInternetGatewayReadyCondition = "EgressOnlyInternetGatewayReady"

	// AWSClusterCarrierGatewayReadyCondition reports on the successful reconciliation of carrier gateways.
	// Only applicable to managed clusters.
	AWSClusterCarrierGatewayReadyCondition = "CarrierGatewayReady"

	// AWSClusterNatGatewaysReadyCondition reports on the successful reconciliation of NAT gateways.
	// Only applicable to managed clusters.
	AWSClusterNatGatewaysReadyCondition = "NatGatewaysReady"

	// AWSClusterRouteTablesReadyCondition reports on the successful reconciliation of route tables.
	// Only applicable to managed clusters.
	AWSClusterRouteTablesReadyCondition = "RouteTablesReady"

	// AWSClusterVpcEndpointsReadyCondition reports on the successful reconciliation of VPC endpoints.
	// Only applicable to managed clusters.
	AWSClusterVpcEndpointsReadyCondition = "VpcEndpointsReady"

	// AWSClusterSecondaryCidrsReadyCondition reports on the successful reconciliation of secondary CIDR blocks.
	// Only applicable to managed clusters.
	AWSClusterSecondaryCidrsReadyCondition = "SecondaryCidrsReady"

	// AWSClusterSecurityGroupsReadyCondition reports on the successful reconciliation of cluster security groups.
	AWSClusterSecurityGroupsReadyCondition = "ClusterSecurityGroupsReady"

	// AWSClusterBastionHostReadyCondition reports whether a bastion host is ready. Depending on the configuration,
	// a cluster may not require a bastion host and this condition will be skipped.
	AWSClusterBastionHostReadyCondition = "BastionHostReady"

	// AWSClusterLoadBalancerReadyCondition reports on whether a control plane load balancer was successfully reconciled.
	AWSClusterLoadBalancerReadyCondition = "LoadBalancerReady"

	// AWSClusterS3BucketReadyCondition reports on the successful reconciliation of an S3 bucket.
	AWSClusterS3BucketReadyCondition = "S3BucketReady"

	// AWSClusterPrincipalCredentialRetrievedCondition reports on whether principal credentials could be retrieved successfully.
	// A possible scenario, where retrieval is unsuccessful, is when SourcePrincipal is not authorized for assume role.
	AWSClusterPrincipalCredentialRetrievedCondition = "PrincipalCredentialRetrieved"

	// AWSClusterPrincipalUsageAllowedCondition reports on whether the principal and all nested source identities
	// are allowed to be used in the AWSCluster namespace.
	AWSClusterPrincipalUsageAllowedCondition = "PrincipalUsageAllowed"
)

// AWSCluster v1beta2 reason constants.
const (
	// AWSClusterReadyReason indicates the AWSCluster infrastructure is ready.
	AWSClusterReadyReason = clusterv1.ReadyReason

	// AWSClusterNotReadyReason indicates the AWSCluster infrastructure is not ready.
	AWSClusterNotReadyReason = clusterv1.NotReadyReason

	// AWSClusterDeletingReason indicates the AWSCluster is being deleted.
	AWSClusterDeletingReason = clusterv1.DeletingReason

	// AWSClusterVpcReconciliationFailedReason used when errors occur during VPC reconciliation.
	AWSClusterVpcReconciliationFailedReason = "VpcReconciliationFailed"

	// AWSClusterVpcCreationStartedReason used when attempting to create a VPC for a managed cluster.
	// Will not be applied to unmanaged clusters.
	AWSClusterVpcCreationStartedReason = "VpcCreationStarted"

	// AWSClusterSubnetsReconciliationFailedReason used to report failures while reconciling subnets.
	AWSClusterSubnetsReconciliationFailedReason = "SubnetsReconciliationFailed"

	// AWSClusterInternetGatewayFailedReason used when errors occur during internet gateway reconciliation.
	AWSClusterInternetGatewayFailedReason = "InternetGatewayFailed"

	// AWSClusterEgressOnlyInternetGatewayFailedReason used when errors occur during egress only internet gateway reconciliation.
	AWSClusterEgressOnlyInternetGatewayFailedReason = "EgressOnlyInternetGatewayFailed"

	// AWSClusterCarrierGatewayFailedReason used when errors occur during carrier gateway reconciliation.
	AWSClusterCarrierGatewayFailedReason = "CarrierGatewayFailed"

	// AWSClusterNatGatewaysReconciliationFailedReason used when any errors occur during reconciliation of NAT gateways.
	AWSClusterNatGatewaysReconciliationFailedReason = "NatGatewaysReconciliationFailed"

	// AWSClusterNatGatewaysCreationStartedReason set once when creating new NAT gateways.
	AWSClusterNatGatewaysCreationStartedReason = "NatGatewaysCreationStarted"

	// AWSClusterRouteTableReconciliationFailedReason used when any errors occur during reconciliation of route tables.
	AWSClusterRouteTableReconciliationFailedReason = "RouteTableReconciliationFailed"

	// AWSClusterVpcEndpointsReconciliationFailedReason used when any errors occur during reconciliation of VPC endpoints.
	AWSClusterVpcEndpointsReconciliationFailedReason = "VpcEndpointsReconciliationFailed"

	// AWSClusterSecondaryCidrReconciliationFailedReason used when any errors occur during reconciliation of secondary CIDR blocks.
	AWSClusterSecondaryCidrReconciliationFailedReason = "SecondaryCidrReconciliationFailed"

	// AWSClusterSecurityGroupReconciliationFailedReason used when any errors occur during reconciliation of security groups.
	AWSClusterSecurityGroupReconciliationFailedReason = "SecurityGroupReconciliationFailed"

	// AWSClusterBastionHostFailedReason used when an error occurs during the creation of a bastion host.
	AWSClusterBastionHostFailedReason = "BastionHostFailed"

	// AWSClusterBastionCreationStartedReason used when creating a new bastion host.
	AWSClusterBastionCreationStartedReason = "BastionCreationStarted"

	// AWSClusterLoadBalancerFailedReason used when an error occurs during load balancer reconciliation.
	AWSClusterLoadBalancerFailedReason = "LoadBalancerFailed"

	// AWSClusterWaitForDNSNameReason used while waiting for a DNS name for the API server to be populated.
	AWSClusterWaitForDNSNameReason = "WaitForDNSName"

	// AWSClusterWaitForExternalControlPlaneEndpointReason is set when the AWSCluster is waiting for an externally managed
	// load balancer, such as an external control plane provider.
	AWSClusterWaitForExternalControlPlaneEndpointReason = "WaitForExternalControlPlaneEndpoint"

	// AWSClusterWaitForDNSNameResolveReason used while waiting for DNS name to resolve.
	AWSClusterWaitForDNSNameResolveReason = "WaitForDNSNameResolve"

	// AWSClusterS3BucketFailedReason used when any errors occur during reconciliation of an S3 bucket.
	AWSClusterS3BucketFailedReason = "S3BucketCreationFailed"

	// AWSClusterPrincipalCredentialRetrievalFailedReason used when errors occur during identity credential retrieval.
	AWSClusterPrincipalCredentialRetrievalFailedReason = "PrincipalCredentialRetrievalFailed"

	// AWSClusterCredentialProviderBuildFailedReason used when errors occur during building providers before trying credential retrieval.
	//nolint:gosec
	AWSClusterCredentialProviderBuildFailedReason = "CredentialProviderBuildFailed"

	// AWSClusterPrincipalUsageUnauthorizedReason used when AWSCluster namespace is not in the identity's allowed namespaces list.
	AWSClusterPrincipalUsageUnauthorizedReason = "PrincipalUsageUnauthorized"

	// AWSClusterSourcePrincipalUsageUnauthorizedReason used when AWSCluster namespace is not in the intersection
	// of source identity allowed namespaces and allowed namespaces of the identities that source identity depends on.
	AWSClusterSourcePrincipalUsageUnauthorizedReason = "SourcePrincipalUsageUnauthorized"
)

// AWSMachine v1beta2 condition types.
const (
	// AWSMachineReadyCondition defines the Ready condition type that summarizes the operational state of an AWSMachine.
	AWSMachineReadyCondition = clusterv1.ReadyCondition

	// AWSMachineInstanceReadyCondition reports on current status of the EC2 instance. Ready indicates the instance is in a Running state.
	AWSMachineInstanceReadyCondition = "InstanceReady"

	// AWSMachineSecurityGroupsReadyCondition indicates the security groups are up to date on the AWSMachine.
	AWSMachineSecurityGroupsReadyCondition = "SecurityGroupsReady"

	// AWSMachineELBAttachedCondition will report true when a control plane is successfully registered with an ELB.
	// When set to false, the subnet may not be found or unavailable in the instance's AZ.
	// Only applicable to control plane machines.
	AWSMachineELBAttachedCondition = "ELBAttached"

	// AWSMachineDedicatedHostReleaseCondition reports on the status of dedicated host release operations.
	// This condition tracks whether the dedicated host has been successfully released or if there are failures.
	AWSMachineDedicatedHostReleaseCondition = "DedicatedHostRelease"
)

// AWSMachine v1beta2 reason constants.
const (
	// AWSMachineReadyReason indicates the AWSMachine is ready.
	AWSMachineReadyReason = clusterv1.ReadyReason

	// AWSMachineNotReadyReason indicates the AWSMachine is not ready.
	AWSMachineNotReadyReason = clusterv1.NotReadyReason

	// AWSMachineDeletingReason indicates the AWSMachine is being deleted.
	AWSMachineDeletingReason = clusterv1.DeletingReason

	// AWSMachineInstanceNotFoundReason used when the instance couldn't be retrieved.
	AWSMachineInstanceNotFoundReason = "InstanceNotFound"

	// AWSMachineInstanceTerminatedReason used when the instance is in a terminated state.
	AWSMachineInstanceTerminatedReason = "InstanceTerminated"

	// AWSMachineInstanceStoppedReason used when the instance is in a stopped state.
	AWSMachineInstanceStoppedReason = "InstanceStopped"

	// AWSMachineInstanceNotReadyReason used when the instance is in a pending state.
	AWSMachineInstanceNotReadyReason = "InstanceNotReady"

	// AWSMachineInstanceProvisionStartedReason set when the provisioning of an instance started.
	AWSMachineInstanceProvisionStartedReason = "InstanceProvisionStarted"

	// AWSMachineInstanceProvisionFailedReason used for failures during instance provisioning.
	AWSMachineInstanceProvisionFailedReason = "InstanceProvisionFailed"

	// AWSMachineWaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	AWSMachineWaitingForClusterInfrastructureReason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// AWSMachineWaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	AWSMachineWaitingForBootstrapDataReason = clusterv1.WaitingForBootstrapDataReason

	// AWSMachineSecurityGroupsFailedReason used when the security groups could not be synced.
	AWSMachineSecurityGroupsFailedReason = "SecurityGroupsSyncFailed"

	// AWSMachineELBAttachFailedReason used when a control plane node fails to attach to the ELB.
	AWSMachineELBAttachFailedReason = "ELBAttachFailed"

	// AWSMachineELBDetachFailedReason used when a control plane node fails to detach from an ELB.
	AWSMachineELBDetachFailedReason = "ELBDetachFailed"

	// AWSMachineDedicatedHostReleaseFailedReason used when the dedicated host release fails.
	AWSMachineDedicatedHostReleaseFailedReason = "DedicatedHostReleaseFailed"
)
