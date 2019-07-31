/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha2

// ClusterPhase is a string representation of a Cluster Phase.
//
// This type is a high-level indicator of the status of the Cluster as it is provisioned,
// from the API user’s perspective.
//
// The value should not be interpreted by any software components as a reliable indication
// of the actual state of the Cluster, and controllers should not use the Cluster Phase field
// value when making decisions about what action to take.
//
// Controllers should always look at the actual state of the Cluster’s fields to make those decisions.
type ClusterPhase string

var (
	// ClusterPhasePending is the first state a Cluster is assigned by
	// Cluster API Cluster controller after being created.
	ClusterPhasePending = ClusterPhase("pending")

	// ClusterPhaseProvisioning is the state when the Cluster has a provider infrastructure
	// object associated and can start provisioning.
	ClusterPhaseProvisioning = ClusterPhase("provisioning")

	// ClusterPhaseProvisioned is the state when its
	// infrastructure has been created and configured.
	ClusterPhaseProvisioned = ClusterPhase("provisioned")

	// ClusterPhaseDeleting is the Cluster state when a delete
	// request has been sent to the API Server,
	// but its infrastructure has not yet been fully deleted.
	ClusterPhaseDeleting = ClusterPhase("deleting")

	// ClusterPhaseFailed is the Cluster state when the system
	// might require user intervention.
	ClusterPhaseFailed = ClusterPhase("failed")

	// ClusterPhaseUnknown is returned if the Cluster state cannot be determined.
	ClusterPhaseUnknown = ClusterPhase("")
)

// ClusterPhaseStringPtr is a helper method to convert ClusterPhase to a string pointer.
func ClusterPhaseStringPtr(p ClusterPhase) *string {
	s := string(p)
	return &s
}
