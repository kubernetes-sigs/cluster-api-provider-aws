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

// MachinePhase is a string representation of a Machine Phase.
//
// This type is a high-level indicator of the status of the Machine as it is provisioned,
// from the API user’s perspective.
//
// The value should not be interpreted by any software components as a reliable indication
// of the actual state of the Machine, and controllers should not use the Machine Phase field
// value when making decisions about what action to take.
//
// Controllers should always look at the actual state of the Machine’s fields to make those decisions.
type MachinePhase string

var (
	// MachinePhasePending is the first state a Machine is assigned by
	// Cluster API Machine controller after being created.
	MachinePhasePending = MachinePhase("pending")

	// MachinePhaseProvisioning is the state when the
	// Machine infrastructure is being created.
	MachinePhaseProvisioning = MachinePhase("provisioning")

	// MachinePhaseProvisioned is the state when its
	// infrastructure has been created and configured.
	MachinePhaseProvisioned = MachinePhase("provisioned")

	// MachinePhaseRunning is the Machine state when it has
	// become a Kubernetes Node in a Ready state.
	MachinePhaseRunning = MachinePhase("running")

	// MachinePhaseDeleting is the Machine state when a delete
	// request has been sent to the API Server,
	// but its infrastructure has not yet been fully deleted.
	MachinePhaseDeleting = MachinePhase("deleting")

	// MachinePhaseDeleted is the Machine state when the object
	// and the related infrastructure is deleted and
	// ready to be garbage collected by the API Server.
	MachinePhaseDeleted = MachinePhase("deleted")

	// MachinePhaseFailed is the Machine state when the system
	// might require user intervention.
	MachinePhaseFailed = MachinePhase("failed")

	// MachinePhaseUnknown is returned if the Machine state cannot be determined.
	MachinePhaseUnknown = MachinePhase("")
)

// MachinePhaseStringPtr is a helper method to convert MachinePhase to a string pointer.
func MachinePhaseStringPtr(p MachinePhase) *string {
	s := string(p)
	return &s
}
