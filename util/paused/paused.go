/*
Copyright 2025 The Kubernetes Authors.

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

// This is a copy of the paused helpers from cluster-api which supports both
// legacy-only resources and resources transitioning to standard conditions.

// Package paused implements paused helper functions.
package paused

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	infrautilconditions "sigs.k8s.io/cluster-api-provider-aws/v2/util/conditions"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/annotations"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
)

// ConditionSetter combines the client.Object and Setter interface.
type ConditionSetter interface {
	v1beta1conditions.Setter
	client.Object
}

// EnsurePausedCondition sets the paused condition on the object and returns if it should be considered as paused.
func EnsurePausedCondition(ctx context.Context, c client.Client, cluster *clusterv1.Cluster, obj ConditionSetter) (isPaused bool, conditionChanged bool, err error) {
	oldCondition := v1beta1conditions.Get(obj, clusterv1beta1.PausedV1Beta2Condition)
	newCondition := pausedCondition(c.Scheme(), cluster, obj, string(clusterv1beta1.PausedV1Beta2Condition))
	standardCondition := metav1.Condition{
		Type:    clusterv1.PausedCondition,
		Status:  metav1.ConditionStatus(newCondition.Status),
		Reason:  newCondition.Reason,
		Message: newCondition.Message,
	}
	standardConditionChanged := false
	standardConditionGenerationChanged := false
	if target, ok := obj.(conditions.Setter); ok {
		old := conditions.Get(target, standardCondition.Type)
		standardConditionChanged = old == nil || old.Status != standardCondition.Status ||
			old.Reason != standardCondition.Reason || old.Message != standardCondition.Message
		standardConditionGenerationChanged = old != nil && old.ObservedGeneration != obj.GetGeneration()
	}

	isPaused = newCondition.Status == corev1.ConditionTrue

	log := ctrl.LoggerFrom(ctx)

	conditionChanged = oldCondition == nil || !v1beta1conditions.HasSameState(oldCondition, &newCondition) || standardConditionChanged
	// Return early if neither the paused condition nor its observed generation changed.
	if !conditionChanged && !standardConditionGenerationChanged {
		if isPaused {
			log.V(6).Info("Reconciliation is paused for this object", "reason", newCondition.Message)
		}
		return isPaused, false, nil
	}

	patchHelper, err := v1beta1patch.NewHelper(obj, c)
	if err != nil {
		return isPaused, false, err
	}

	if conditionChanged && isPaused {
		log.V(4).Info("Pausing reconciliation for this object", "reason", newCondition.Message)
	} else if conditionChanged {
		log.V(4).Info("Unpausing reconciliation for this object")
	}

	v1beta1conditions.Set(obj, &newCondition)
	infrautilconditions.SetV1Beta2(obj, standardCondition)

	if err := patchHelper.Patch(ctx, obj, v1beta1patch.WithOwnedConditions{Conditions: []clusterv1beta1.ConditionType{
		clusterv1beta1.PausedV1Beta2Condition,
	}}, v1beta1patch.WithOwnedV1Beta2Conditions{Conditions: []string{clusterv1.PausedCondition}}); err != nil {
		return isPaused, false, err
	}

	// Persist generation-only updates here because callers create their patch helpers afterward.
	// Report no state change so unpaused reconciliation can continue without another pass.
	return isPaused, conditionChanged, nil
}

// pausedCondition sets the paused condition on the object and returns if it should be considered as paused.
func pausedCondition(scheme *runtime.Scheme, cluster *clusterv1.Cluster, obj ConditionSetter, targetConditionType string) clusterv1beta1.Condition {
	if (cluster != nil && ptr.Deref(cluster.Spec.Paused, false)) || annotations.HasPaused(obj) {
		var messages []string
		if cluster != nil && ptr.Deref(cluster.Spec.Paused, false) {
			messages = append(messages, "Cluster spec.paused is set to true")
		}
		if annotations.HasPaused(obj) {
			kind := "Object"
			if gvk, err := apiutil.GVKForObject(obj, scheme); err == nil {
				kind = gvk.Kind
			}
			messages = append(messages, fmt.Sprintf("%s has the cluster.x-k8s.io/paused annotation", kind))
		}

		return clusterv1beta1.Condition{
			Type:    clusterv1beta1.ConditionType(targetConditionType),
			Status:  corev1.ConditionTrue,
			Reason:  clusterv1beta1.PausedV1Beta2Reason,
			Message: strings.Join(messages, ", "),
		}
	}

	return clusterv1beta1.Condition{
		Type:   clusterv1beta1.ConditionType(targetConditionType),
		Status: corev1.ConditionFalse,
		Reason: clusterv1beta1.NotPausedV1Beta2Reason,
	}
}
