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
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	service "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services"

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
func (a *Actuator) ensureSecurityGroups(ec2svc service.EC2MachineInterface, machine *clusterv1.Machine, instanceID *string, additionalSecurityGroups []v1alpha1.AWSResourceReference, instanceSecurityGroups map[string]string) (bool, error) {
	// Get the SecurityGroup annotations
	annotation, err := a.machineAnnotationJSON(machine, SecurityGroupsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Now we check to see if we need to mutate any mutable state.
	// Check if any additions have been made to the instance Security Groups.
	changed, groupIDs, newAnnotation := a.securityGroupsChanged(
		annotation,
		additionalSecurityGroups,
		instanceSecurityGroups,
	)
	if changed {
		// Finally, update the instance with our new security groups.
		err := ec2svc.UpdateInstanceSecurityGroups(instanceID, groupIDs)
		if err != nil {
			return false, err
		}

		// And then update the annotation to reflect our new managed state.
		err = a.updateMachineAnnotationJSON(machine, SecurityGroupsLastAppliedAnnotation, newAnnotation)
		if err != nil {
			return false, err
		}
	}

	return changed, nil
}

// securityGroupsChanged determines which security groups to delete and which to
// add.
func (a *Actuator) securityGroupsChanged(annotation map[string]interface{}, src []v1alpha1.AWSResourceReference, dst map[string]string) (bool, []*string, map[string]interface{}) {
	// State tracking for decisions later in the method.
	// The bool here indicates if a listed groupID was deleted or not, which is
	// used while generating our final array later.
	state := map[string]bool{}

	// Get the `src` into an easily usable map. We call this `newAnnotation` as
	// this map will also be returned at the end to be used as the new
	// machine annotation.
	//
	// We can also set each entry in the `src` to true in the `state`, because
	// we know we want to keep these. These additions to the state could
	// represent new or unchanged securityGroups. It doesn't matter which.
	newAnnotation := map[string]interface{}{}
	for _, s := range src {
		newAnnotation[*s.ID] = struct{}{}
		state[*s.ID] = true
	}

	// Loop over `annotation`, checking `newAnnotation` for things that were
	// deleted since last time.
	// If we find something in the `annotation`, but not in the `newAnnotation`,
	// we flag it as `false` (not found, deleted).
	for groupID := range annotation {
		_, ok := newAnnotation[groupID]

		if !ok {
			state[groupID] = false
		}
	}

	// At this point our `state` variable now represents which groupIDs we wish
	// to keep, and which we wish to drop, based on the incoming Kubernetes
	// state in `newAnnotation` (`src`) and the state from last time in the
	// `annotation`.
	// This `state` variable represents the state of what we have been and are
	// responsible for managing.
	// We are now able to build up a new map, which will represent the combined
	// state of what the machine actuator is managing, as well as groups being
	// managed externally.
	groupsMap := map[string]struct{}{}

	// Add groups that we want to keep from the new state to groupsMap.
	for groupID, keep := range state {
		if keep {
			groupsMap[groupID] = struct{}{}
		}
	}

	// Loop over `dst` (AWS) comparing entries to the `state` map.
	//
	// If we find a `dst` entry in the `state` map and it's set to `true`, add
	// it to the `groupsMap`. (We're managing it and wish to keep it)
	//
	// If we find an entry in the `state` map and it's set to `false`, don't add
	// it to the `groupsMap`. (We're managing it and wish to delete it)
	//
	// If we don't find an entry in the `state` map, we aren't managing this
	// groupID, add it to `groupsMap`. (It's externally managed, keep it)
	for groupID := range dst {
		keep, ok := state[groupID]

		// If we don't find this in the state map, we aren't managing it. Add
		// it to the groupsMap and continue.
		if !ok {
			groupsMap[groupID] = struct{}{}
			continue
		}

		// Ok, we found it in the state map, what do we want to do?
		// If we're keeping it, add it to the groupsMap, otherwise do nothing
		// and let it drop.
		if keep {
			groupsMap[groupID] = struct{}{}
		}
	}

	// Variable indicating if state was changed compared to AWS.
	changed := false

	// At this point, we know exactly what we're going to send to AWS and what
	// our new annotation should be.
	// We perform a final comparison of the dst (AWS) state to the groupsMap,
	// to see if we really need to send a request to AWS.
	if len(dst) == len(groupsMap) {
		// We need to compare each item, in case any changed.
		for groupID := range dst {
			_, ok := groupsMap[groupID]
			if !ok {
				// A change was detected.
				changed = true
				break
			}
		}
	} else {
		// If the lengths are different, the state definitely changed. No
		// comparison needed.
		changed = true
	}

	// Generate a groups array to send back to the calling function.
	groupsArray := []*string{}
	for groupID := range groupsMap {
		groupsArray = append(groupsArray, &groupID)
	}

	// Finally, we're done.
	return changed, groupsArray, newAnnotation
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
