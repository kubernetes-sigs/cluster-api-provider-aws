// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machine

// should not need to import the ec2 sdk here
import (
	"sort"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/providerconfig/v1alpha1"
	service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	// SecurityGroupsLastAppliedAnnotation is the key for the machine object
	// annotation which tracks the SecurityGroups that the machine actuator is
	// responsible for. These are the SecurityGroups that have been handled by
	// the AdditionalSecurityGroups in the Machine Provider Config.
	SecurityGroupsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws/last-applied/security-groups"
)

// Ensures that the security groups of the machine are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (a *Actuator) ensureSecurityGroups(ec2svc service.EC2MachineInterface, machine *clusterv1.Machine, instanceID string, additional []v1alpha1.AWSResourceReference, existing []string) (bool, error) {
	annotation, err := a.machineAnnotationJSON(machine, SecurityGroupsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	changed, ids := a.securityGroupsChanged(annotation, additional, existing)
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

	if err := a.updateMachineAnnotationJSON(machine, SecurityGroupsLastAppliedAnnotation, newAnnotation); err != nil {
		return false, err
	}

	return true, nil
}

// securityGroupsChanged determines which security groups to delete and which to add.
func (a *Actuator) securityGroupsChanged(annotation map[string]interface{}, additional []v1alpha1.AWSResourceReference, existing []string) (bool, []string) {
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

	// Build the security group list.
	res := []string{}
	for id, keep := range state {
		if keep {
			res = append(res, id)
		}
	}

	// Add groups managed externally (i.e. not in the state).
	for _, id := range existing {
		if _, ok := state[id]; !ok {
			res = append(res, id)
		}
	}

	changed := len(existing) != len(res)

	if !changed {
		// Length is the same, check if the ids are the same too.
		sort.Strings(existing)
		sort.Strings(res)
		for i, id := range res {
			if existing[i] != id {
				changed = true
				break
			}
		}
	}

	return changed, res
}

// tagsChanged determines which tags to delete and which to add.
func (a *Actuator) tagsChanged(annotation map[string]interface{}, src map[string]string) (bool, map[string]string, map[string]string, map[string]interface{}) {
	// Bool tracking if we found any changed state.
	changed := false

	// Tracking for created/updated
	created := map[string]string{}

	// Tracking for tags that were deleted.
	deleted := map[string]string{}

	// The new annotation that we need to set if anything is created/updated.
	newAnnotation := map[string]interface{}{}

	// Loop over annotation, checking if entries are in src.
	// If an entry is present in annotation but not src, it has been deleted
	// since last time. We flag this in the deleted map.
	for t, v := range annotation {
		_, ok := src[t]

		// Entry isn't in src, it has been deleted.
		if !ok {
			// Cast v to a string here. This should be fine, tags are always
			// strings.
			deleted[t] = v.(string)
			changed = true
		}
	}

	// Loop over src, checking for entries in annotation.
	//
	// If an entry is in src, but not annotation, it has been created since
	// last time.
	//
	// If an entry is in both src and annotation, we compare their values, if
	// the value in src differs from that in annotation, the tag has been
	// updated since last time.
	for t, v := range src {
		av, ok := annotation[t]

		// Entries in the src always need to be noted in the newAnnotation. We
		// know they're going to be created or updated.
		newAnnotation[t] = v

		// Entry isn't in annotation, it's new.
		if !ok {
			created[t] = v
			newAnnotation[t] = v
			changed = true
			continue
		}

		// Entry is in annotation, has the value changed?
		if v != av {
			created[t] = v
			changed = true
		}

		// Entry existed in both src and annotation, and their values were
		// equal. Nothing to do.
	}

	// We made it through the loop, and everything that was in src, was also
	// in dst. Nothing changed.
	return changed, created, deleted, newAnnotation
}
