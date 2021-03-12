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

package v1beta1

// Conditions and condition Reasons for the MachineHealthCheck object

const (
	// RemediationAllowedCondition is set on MachineHealthChecks to show the status of whether the MachineHealthCheck is
	// allowed to remediate any Machines or whether it is blocked from remediating any further.
	RemediationAllowedCondition ConditionType = "RemediationAllowed"

	// TooManyUnhealthy is the reason used when too many Machines are unhealthy and the MachineHealthCheck is blocked
	// from making any further remediations.
	TooManyUnhealthyReason = "TooManyUnhealthy"
)

const (
	// InstanceExistsCondition is set on the Machine to show whether a virtual mahcine has been created by the cloud provider.
	InstanceExistsCondition ConditionType = "InstanceExists"

	// ErrorCheckingProviderReason is the reason used when the exist operation fails.
	// This would normally be because we cannot contact the provider.
	ErrorCheckingProviderReason = "ErrorCheckingProvider"

	// InstanceMissingReason is the reason used when the machine was provisioned, but the instance has gone missing.
	InstanceMissingReason = "InstanceMissing"

	// InstanceNotCreatedReason is the reason used when the machine has not yet been provisioned.
	InstanceNotCreatedReason = "InstanceNotCreated"
)
