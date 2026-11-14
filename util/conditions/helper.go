/*
Copyright 2021 The Kubernetes Authors.

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

// Package conditions provides helper functions for working with conditions.
package conditions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	conditions "sigs.k8s.io/cluster-api/util/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
)

// SetV1Beta2 sets a standard condition when the object supports standard conditions.
// Shared services also accept resources which only expose legacy conditions.
func SetV1Beta2(obj any, condition metav1.Condition) {
	switch target := obj.(type) {
	case v1beta2conditions.Setter:
		v1beta2conditions.Set(target, condition)
	case conditions.Setter:
		conditions.Set(target, condition)
	}
}

// SetV1Beta2Summary summarizes the applicable conditions into the standard Ready condition.
func SetV1Beta2Summary(obj v1beta2conditions.Setter, applicableConditions ...clusterv1beta1.ConditionType) error {
	conditionTypes := make([]string, 0, len(applicableConditions))
	for _, conditionType := range applicableConditions {
		standardType := string(conditionType)
		if conditionType == infrav1.VpcEndpointsReadyCondition {
			standardType = infrav1.AWSClusterVpcEndpointsReadyV1Beta2Condition
		}
		conditionTypes = append(conditionTypes, standardType)
	}
	return v1beta2conditions.SetSummaryCondition(obj, obj, clusterv1.ReadyCondition,
		v1beta2conditions.ForConditionTypes(conditionTypes),
		v1beta2conditions.CustomMergeStrategy{MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
			v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
				clusterv1.NotReadyReason, clusterv1.ReadyUnknownReason, clusterv1.ReadyReason,
			)),
		)},
	)
}

// ErrorConditionAfterInit returns severity error, if the control plane is initialized; otherwise, returns severity warning.
// Failures after control plane is initialized is likely to be non-transient,
// hence conditions severities should be set to Error.
func ErrorConditionAfterInit(cluster *clusterv1.Cluster) clusterv1beta1.ConditionSeverity {
	if ptr.Deref(cluster.Status.Initialization.ControlPlaneInitialized, false) {
		return clusterv1beta1.ConditionSeverityError
	}
	return clusterv1beta1.ConditionSeverityWarning
}
