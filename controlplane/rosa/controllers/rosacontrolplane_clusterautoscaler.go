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

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// reconcileClusterAutoscaler reports the OCM cluster autoscaler configuration in status and, when
// spec.clusterAutoscaler is set, converges OCM to it by sending a PATCH with only the fields that differ.
// It must only be called once the OCM cluster is ready.
//
// It returns an error only for failures that should be retried (for example OCM being unavailable).
// Errors caused by the configuration itself are reported in the ClusterAutoscalerConfigured condition.
// When spec.clusterAutoscaler is not set, failures only affect status and are never returned.
func (r *ROSAControlPlaneReconciler) reconcileClusterAutoscaler(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope, cluster *cmv1.Cluster) error {
	spec := rosaScope.ControlPlane.Spec.ClusterAutoscaler
	managed := spec != nil
	if !managed {
		v1beta1conditions.Delete(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)
	}

	if r.NewClusterAutoscalerClient == nil {
		return clusterAutoscalerError(rosaScope, managed, fmt.Errorf("failed to create cluster autoscaler client: NewClusterAutoscalerClient is nil"))
	}
	autoscalerClient, err := r.NewClusterAutoscalerClient(ctx, rosaScope)
	if err != nil {
		return clusterAutoscalerError(rosaScope, managed, fmt.Errorf("failed to create cluster autoscaler client: %w", err))
	}
	defer autoscalerClient.Close()

	current, err := autoscalerClient.GetClusterAutoscaler(ctx, cluster.ID())
	if err != nil {
		return clusterAutoscalerError(rosaScope, managed, fmt.Errorf("failed to get cluster autoscaler: %w", err))
	}

	if managed {
		patch, err := rosa.BuildClusterAutoscalerPatch(spec, current)
		if err != nil {
			return clusterAutoscalerError(rosaScope, managed, err)
		}
		if patch != nil {
			rosaScope.Info("updating cluster autoscaler", "clusterID", cluster.ID(),
				"fields", strings.Join(rosa.ClusterAutoscalerPatchFields(patch), ","))
			updated, err := autoscalerClient.UpdateClusterAutoscaler(ctx, cluster.ID(), patch)
			if err != nil {
				rosaScope.ControlPlane.Status.ClusterAutoscaler = rosa.ClusterAutoscalerStatusFromOCM(current)
				return clusterAutoscalerError(rosaScope, managed, fmt.Errorf("failed to update cluster autoscaler: %w", err))
			}
			if updated != nil {
				current = updated
			}
		}
		v1beta1conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ClusterAutoscalerConfiguredCondition)
	}

	rosaScope.ControlPlane.Status.ClusterAutoscaler = rosa.ClusterAutoscalerStatusFromOCM(current)
	return nil
}

// clusterAutoscalerError reports err and decides whether it should be retried.
// When the autoscaler is not managed, the error is only logged.
// When it is managed, the ClusterAutoscalerConfigured condition is set from the OCM HTTP status:
// 400 means OCM rejected the configuration and 403/404 mean the autoscaler API isn't available for
// the cluster; neither is retried. Any other failure (5xx, network, ...) is returned to be retried.
func clusterAutoscalerError(rosaScope *scope.ROSAControlPlaneScope, managed bool, err error) error {
	if !managed {
		rosaScope.Info("failed to read cluster autoscaler, status not updated", "error", err.Error())
		return nil
	}

	switch rosa.OCMAPIErrorStatus(err) {
	case http.StatusBadRequest:
		v1beta1conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ClusterAutoscalerConfiguredCondition,
			rosacontrolplanev1.ClusterAutoscalerValidationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"%s", err.Error())
		return nil
	case http.StatusForbidden, http.StatusNotFound:
		v1beta1conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ClusterAutoscalerConfiguredCondition,
			rosacontrolplanev1.ClusterAutoscalerNotAvailableReason,
			clusterv1beta1.ConditionSeverityError,
			"%s", err.Error())
		return nil
	default:
		v1beta1conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ClusterAutoscalerConfiguredCondition,
			rosacontrolplanev1.ReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityWarning,
			"%s", err.Error())
		return err
	}
}

// markClusterAutoscalerWaitingForCluster reports that spec.clusterAutoscaler will be sent to OCM once the cluster is ready.
func markClusterAutoscalerWaitingForCluster(rosaScope *scope.ROSAControlPlaneScope, state cmv1.ClusterState) {
	if rosaScope.ControlPlane.Spec.ClusterAutoscaler == nil {
		return
	}
	v1beta1conditions.MarkFalse(rosaScope.ControlPlane,
		rosacontrolplanev1.ClusterAutoscalerConfiguredCondition,
		rosacontrolplanev1.ClusterAutoscalerWaitingForClusterReason,
		clusterv1beta1.ConditionSeverityInfo,
		"cluster autoscaler configuration will be sent to OCM once the cluster is ready, current state is %q", state)
}
