/*
Copyright 2018 The Kubernetes Authors.

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

package machine

// should not need to import the ec2 sdk here
import (
	"sort"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
)

const (
	// SecurityGroupsLastAppliedAnnotation is the key for the machine object
	// annotation which tracks the SecurityGroups that the machine actuator is
	// responsible for. These are the SecurityGroups that have been handled by
	// the AdditionalSecurityGroups in the Machine Provider Config.
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// for annotation formatting rules.
	SecurityGroupsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws-last-applied-security-groups"
)

// Ensures that the security groups of the machine are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (a *Actuator) ensureSecurityGroups(ec2svc service.EC2MachineInterface, scope *actuators.MachineScope, instanceID string, additional []v1alpha1.AWSResourceReference, existing map[string][]string) (bool, error) {
	annotation, err := a.machineAnnotationJSON(scope.Machine, SecurityGroupsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	core, err := ec2svc.GetCoreSecurityGroups(scope)
	if err != nil {
		return false, err
	}
	changed, ids := a.securityGroupsChanged(annotation, core, additional, existing)
	if !changed {
		return false, nil
	}

	if err := ec2svc.UpdateInstanceSecurityGroups(instanceID, ids); err != nil {
		return false, err
	}

	// Build and store annotation.
	newAnnotation := make(map[string]interface{}, len(additional))
	for _, id := range additional {
		newAnnotation[*id.ID] = struct{}{}
	}

	if err := a.updateMachineAnnotationJSON(scope.Machine, SecurityGroupsLastAppliedAnnotation, newAnnotation); err != nil {
		return false, err
	}

	return true, nil
}

// securityGroupsChanged determines which security groups to delete and which to add.
func (a *Actuator) securityGroupsChanged(annotation map[string]interface{}, core []string, additional []v1alpha1.AWSResourceReference, existing map[string][]string) (bool, []string) {
	state := map[string]bool{}
	for _, s := range additional {
		state[*s.ID] = true
	}

	// Loop over `annotation`, checking the state for things that were deleted since last time.
	// If we find something in the `annotation`, but not in the state, we flag it as `false` (not found, deleted).
	for groupID := range annotation {
		if _, ok := state[groupID]; !ok {
			state[groupID] = false
		}
	}

	// add (or add back) the core security groups
	for _, s := range core {
		state[s] = true
	}

	// Build the security group list.
	res := []string{}
	for id, keep := range state {
		if keep {
			res = append(res, id)
		}
	}

	for _, actual := range existing {
		if len(actual) != len(res) {
			return true, res
		}

		// Length is the same, check if the ids are the same too.
		sort.Strings(actual)
		sort.Strings(res)
		for i, id := range res {
			if actual[i] != id {
				return true, res
			}
		}
	}

	return false, res
}
