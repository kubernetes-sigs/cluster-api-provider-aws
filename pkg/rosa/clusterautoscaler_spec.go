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

package rosa

import (
	"fmt"
	"strings"
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
)

// ROSA HCP cluster autoscaler defaults. They mirror the defaults Clusters Service applies when it
// creates the cluster autoscaler of a new HCP cluster (DefaultClusterAutoScaler in Clusters Service).
// They are applied to fields omitted from a set spec.clusterAutoscaler.
const (
	// DefaultClusterAutoscalerMaxNodesTotal is the default maximum number of nodes. 0 means no limit.
	DefaultClusterAutoscalerMaxNodesTotal int32 = 0
	// DefaultClusterAutoscalerMaxPodGracePeriod is the default pod grace period during scale down, in seconds.
	DefaultClusterAutoscalerMaxPodGracePeriod int32 = 600
	// DefaultClusterAutoscalerMaxNodeProvisionTime is the default maximum node provision time.
	DefaultClusterAutoscalerMaxNodeProvisionTime = 15 * time.Minute
	// DefaultClusterAutoscalerPodPriorityThreshold is the default pod priority threshold.
	DefaultClusterAutoscalerPodPriorityThreshold int32 = -10
)

// clusterAutoscalerConfig is a fully defaulted cluster autoscaler configuration.
type clusterAutoscalerConfig struct {
	maxNodesTotal        int32
	maxPodGracePeriod    int32
	maxNodeProvisionTime time.Duration
	podPriorityThreshold int32
}

// desiredClusterAutoscaler returns the configuration requested by spec, with ROSA HCP defaults for omitted fields.
func desiredClusterAutoscaler(spec *rosacontrolplanev1.ClusterAutoscaler) clusterAutoscalerConfig {
	desired := clusterAutoscalerConfig{
		maxNodesTotal:        DefaultClusterAutoscalerMaxNodesTotal,
		maxPodGracePeriod:    DefaultClusterAutoscalerMaxPodGracePeriod,
		maxNodeProvisionTime: DefaultClusterAutoscalerMaxNodeProvisionTime,
		podPriorityThreshold: DefaultClusterAutoscalerPodPriorityThreshold,
	}
	if spec == nil {
		return desired
	}
	if spec.ResourceLimits != nil && spec.ResourceLimits.MaxNodesTotal != nil {
		desired.maxNodesTotal = *spec.ResourceLimits.MaxNodesTotal
	}
	if spec.MaxPodGracePeriod != nil {
		desired.maxPodGracePeriod = *spec.MaxPodGracePeriod
	}
	if spec.MaxNodeProvisionTime != nil {
		desired.maxNodeProvisionTime = spec.MaxNodeProvisionTime.Duration
	}
	if spec.PodPriorityThreshold != nil {
		desired.podPriorityThreshold = *spec.PodPriorityThreshold
	}
	return desired
}

// BuildClusterAutoscalerPatch returns the body of a PATCH that converges the OCM cluster autoscaler
// (current) to spec, with ROSA HCP defaults for omitted fields. Only fields whose value differs are
// included, so fields that are already correct are never sent. It returns nil when nothing differs.
// A field missing from current is treated as different.
func BuildClusterAutoscalerPatch(spec *rosacontrolplanev1.ClusterAutoscaler, current *cmv1.ClusterAutoscaler) (*cmv1.ClusterAutoscaler, error) {
	desired := desiredClusterAutoscaler(spec)
	builder := cmv1.NewClusterAutoscaler()
	changed := false

	currentMaxNodesTotal, ok := current.ResourceLimits().GetMaxNodesTotal()
	if !ok || int32(currentMaxNodesTotal) != desired.maxNodesTotal { //nolint:gosec // values are validated to fit in int32
		builder.ResourceLimits(cmv1.NewAutoscalerResourceLimits().MaxNodesTotal(int(desired.maxNodesTotal)))
		changed = true
	}

	currentMaxPodGracePeriod, ok := current.GetMaxPodGracePeriod()
	if !ok || int32(currentMaxPodGracePeriod) != desired.maxPodGracePeriod { //nolint:gosec // values are validated to fit in int32
		builder.MaxPodGracePeriod(int(desired.maxPodGracePeriod))
		changed = true
	}

	if !durationEquals(current.MaxNodeProvisionTime(), desired.maxNodeProvisionTime) {
		builder.MaxNodeProvisionTime(formatDuration(desired.maxNodeProvisionTime))
		changed = true
	}

	currentPodPriorityThreshold, ok := current.GetPodPriorityThreshold()
	if !ok || int32(currentPodPriorityThreshold) != desired.podPriorityThreshold { //nolint:gosec // values are validated to fit in int32
		builder.PodPriorityThreshold(int(desired.podPriorityThreshold))
		changed = true
	}

	if !changed {
		return nil, nil
	}
	patch, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build cluster autoscaler patch: %w", err)
	}
	return patch, nil
}

// ClusterAutoscalerStatusFromOCM converts the OCM cluster autoscaler into the status representation.
// Fields not reported by OCM are left unset.
func ClusterAutoscalerStatusFromOCM(current *cmv1.ClusterAutoscaler) *rosacontrolplanev1.ClusterAutoscalerStatus {
	if current == nil {
		return nil
	}
	status := &rosacontrolplanev1.ClusterAutoscalerStatus{}
	if maxNodesTotal, ok := current.ResourceLimits().GetMaxNodesTotal(); ok {
		status.ResourceLimits = &rosacontrolplanev1.ClusterAutoscalerResourceLimitsStatus{
			MaxNodesTotal: ptr.To(int32(maxNodesTotal)), //nolint:gosec // OCM validates the value range
		}
	}
	if maxPodGracePeriod, ok := current.GetMaxPodGracePeriod(); ok {
		status.MaxPodGracePeriod = ptr.To(int32(maxPodGracePeriod)) //nolint:gosec // OCM validates the value range
	}
	if raw, ok := current.GetMaxNodeProvisionTime(); ok {
		if d, err := time.ParseDuration(raw); err == nil {
			status.MaxNodeProvisionTime = &metav1.Duration{Duration: d}
		}
	}
	if podPriorityThreshold, ok := current.GetPodPriorityThreshold(); ok {
		status.PodPriorityThreshold = ptr.To(int32(podPriorityThreshold)) //nolint:gosec // OCM validates the value range
	}
	return status
}

// durationEquals reports whether the OCM duration string raw represents d. OCM may return the value
// in a different format than CAPA sends it (for example "15m" vs "15m0s"), so values are compared parsed.
func durationEquals(raw string, d time.Duration) bool {
	if raw == "" {
		return false
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return false
	}
	return parsed == d
}

// formatDuration formats d in the compact form used by OCM and the ROSA CLI, for example "15m" instead of "15m0s".
func formatDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	if strings.HasSuffix(s, "h0m") {
		s = strings.TrimSuffix(s, "0m")
	}
	return s
}
