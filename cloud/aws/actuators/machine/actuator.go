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
	"encoding/json"
	"fmt"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

const (
	// SecurityGroupsLastAppliedAnnotation is the key for the machine object
	// annotation which tracks the SecurityGroups that the machine actuator is
	// responsible for. These are the SecurityGroups that have been handled by
	// the AdditionalSecurityGroups in the Machine Provider Config.
	SecurityGroupsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws/last-applied/security-groups"

	// TagsLastAppliedAnnotation is the key for the machine object annotation
	// which tracks the SecurityGroups that the machine actuator is responsible
	// for. These are the SecurityGroups that have been handled by the
	// AdditionalTags in the Machine Provider Config.
	TagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws/last-applied/tags"
)

// ec2Svc are the functions from the ec2 service, not the client, this actuator needs.
// This should never need to import the ec2 sdk.
type ec2Svc interface {
	CreateInstance(*clusterv1.Machine, *v1alpha1.AWSMachineProviderConfig, *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error)
	InstanceIfExists(*string) (*v1alpha1.Instance, error)
	TerminateInstance(*string) error
	CreateOrGetMachine(*clusterv1.Machine, *v1alpha1.AWSMachineProviderStatus, *v1alpha1.AWSMachineProviderConfig, *v1alpha1.AWSClusterProviderStatus) (*v1alpha1.Instance, error)
	UpdateInstanceSecurityGroups(*string, []*string) error
	UpdateResourceTags(*string, map[string]string, map[string]string) error
}

// codec are the functions off the generated codec that this actuator uses.
type codec interface {
	DecodeFromProviderConfig(clusterv1.ProviderConfig, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	codec codec

	// Services
	ec2            ec2Svc
	machinesGetter client.MachinesGetter
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	// Codec is needed to work with the provider configs and statuses.
	Codec codec

	// Services

	// ClusterService is the interface to cluster-api.
	MachinesGetter client.MachinesGetter
	// EC2Service is the interface to ec2.
	EC2Service ec2Svc
}

// NewActuator returns an actuator.
func NewActuator(params ActuatorParams) (*Actuator, error) {
	return &Actuator{
		codec:          params.Codec,
		ec2:            params.EC2Service,
		machinesGetter: params.MachinesGetter,
	}, nil
}

// Create creates a machine and is invoked by the machine controller.
func (a *Actuator) Create(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine provider status")
	}
	clusterStatus, err := a.clusterProviderStatus(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster provider status")
	}
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine config")
	}

	i, err := a.ec2.CreateOrGetMachine(machine, status, config, clusterStatus)
	if err != nil {
		return errors.Wrap(err, "failed to create or get machine")
	}

	status.InstanceID = &i.ID
	status.InstanceState = aws.String(string(i.State))

	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}
	machine.Annotations["cluster-api-provider-aws"] = "true"

	return a.updateStatus(machine, status)
}

// Delete deletes a machine and is invoked by the Machine Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Deleting machine %v for cluster %v.", machine.Name, cluster.Name)

	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine provider status")
	}

	instance, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	// The machine hasn't been created yet
	if instance == nil {
		return nil
	}

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case v1alpha1.InstanceStateShuttingDown, v1alpha1.InstanceStateTerminated:
		return nil
	default:
		err = a.ec2.TerminateInstance(status.InstanceID)
		if err != nil {
			return errors.Wrap(err, "failed to terminate instance")
		}
	}

	return nil
}

// Update updates a machine and is invoked by the Machine Controller.
// If the Update attempts to mutate any immutable state, the method will error
// and no updates will be performed.
func (a *Actuator) Update(cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Updating machine %v for cluster %v.", machine.Name, cluster.Name)

	// Get the updated config. We're going to compare parts of this to the
	// current Tags and Security Groups that AWS is aware of.
	config, err := a.machineProviderConfig(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine config")
	}

	// Get the new status from the provided machine object.
	// We need this in case any of it has changed, and we also require the
	// instanceID for various AWS queries.
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return errors.Wrap(err, "failed to get machine status")
	}

	// Get the current instance description from AWS.
	instanceDescription, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	// We can now compare the various AWS state to the state we were passed.
	// We will check immutable state first, in order to fail quickly before
	// moving on to state that we can mutate.
	// TODO: Implement immutable state check.

	// Ensure that the security groups are correct.
	securityGroupsChanged, err := a.ensureSecurityGroups(
		machine,
		status.InstanceID,
		config.AdditionalSecurityGroups,
		instanceDescription.SecurityGroups,
	)
	if err != nil {
		return errors.Wrap(err, "failed to ensure security groups")
	}

	// Ensure that the tags are correct.
	tagsChanged, err := a.ensureTags(machine, status.InstanceID, config.AdditionalTags)
	if err != nil {
		return errors.Wrap(err, "failed to ensure tags")
	}

	// We need to update the machine since annotations may have changed.
	if securityGroupsChanged || tagsChanged {
		err = a.updateMachine(machine)
		if err != nil {
			return errors.Wrap(err, "failed to update machine")
		}
	}

	// Finally update the machine status.
	err = a.updateStatus(machine, status)
	if err != nil {
		return errors.Wrap(err, "failed to update machine status")
	}

	return nil
}

// Exists test for the existence of a machine and is invoked by the Machine Controller
func (a *Actuator) Exists(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (bool, error) {
	glog.Infof("Checking if machine %v for cluster %v exists.", machine.Name, cluster.Name)
	status, err := a.machineProviderStatus(machine)
	if err != nil {
		return false, err
	}

	instance, err := a.ec2.InstanceIfExists(status.InstanceID)
	if err != nil {
		return false, err
	}
	if instance == nil {
		return false, nil
	}
	// TODO update status here
	switch instance.State {
	case v1alpha1.InstanceStateRunning, v1alpha1.InstanceStatePending:
		return true, nil
	default:
		return false, nil
	}
}

// Ensures that the security groups of the machine are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (a *Actuator) ensureSecurityGroups(machine *clusterv1.Machine, instanceID *string, additionalSecurityGroups []v1alpha1.AWSResourceReference, instanceSecurityGroups map[string]string) (bool, error) {
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
		err := a.ec2.UpdateInstanceSecurityGroups(instanceID, groupIDs)
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

// Ensure that the tags of the machine are correct
// Returns bool, error
// Bool indicates if changes were made or not, allowing the caller to decide
// if the machine should be updated.
func (a *Actuator) ensureTags(machine *clusterv1.Machine, instanceID *string, additionalTags map[string]string) (bool, error) {
	annotation, err := a.machineAnnotationJSON(machine, TagsLastAppliedAnnotation)
	if err != nil {
		return false, err
	}

	// Check if the instance tags were changed. If they were, update them.
	// It would be possible here to only send new/updated tags, but for the
	// moment we send everything, even if only a single tag was created or
	// upated.
	changed, created, deleted, newAnnotation := a.tagsChanged(annotation, additionalTags)
	if changed {
		err = a.ec2.UpdateResourceTags(instanceID, created, deleted)
		if err != nil {
			return false, err
		}

		// We also need to update the annotation if anything changed.
		err = a.updateMachineAnnotationJSON(machine, TagsLastAppliedAnnotation, newAnnotation)
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

// Returns a map[string]interface from a JSON annotation.
// This method gets the given `annotation` from the `machine` and unmarshalls it
// from a JSON string into a `map[string]interface{}`.
func (a *Actuator) machineAnnotationJSON(machine *clusterv1.Machine, annotation string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	jsonAnnotation := a.machineAnnotation(machine, annotation)
	if len(jsonAnnotation) == 0 {
		return out, nil
	}

	err := json.Unmarshal([]byte(jsonAnnotation), &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

// Fetches the specific machine annotation.
func (a *Actuator) machineAnnotation(machine *clusterv1.Machine, annotation string) string {
	return machine.GetAnnotations()[annotation]
}

// updateMachineAnnotationJSON updates the `annotation` on `machine` with
// `content`. `content` in this case should be a `map[string]interface{}`
// suitable for turning into JSON. This `content` map will be marshalled into a
// JSON string before being set as the given `annotation`.
func (a *Actuator) updateMachineAnnotationJSON(machine *clusterv1.Machine, annotation string, content map[string]interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	a.updateMachineAnnotation(machine, annotation, string(b))

	return nil
}

// updateMachineAnnotation updates the `annotation` on the given `machine` with
// `content`.
func (a *Actuator) updateMachineAnnotation(machine *clusterv1.Machine, annotation string, content string) {
	// Get the annotations
	annotations := machine.GetAnnotations()

	// Set our annotation to the given content.
	annotations[annotation] = content

	// Update the machine object with these annotations
	machine.SetAnnotations(annotations)
}

func (a *Actuator) machineProviderConfig(machine *clusterv1.Machine) (*v1alpha1.AWSMachineProviderConfig, error) {
	machineProviderCfg := &v1alpha1.AWSMachineProviderConfig{}
	err := a.codec.DecodeFromProviderConfig(machine.Spec.ProviderConfig, machineProviderCfg)
	return machineProviderCfg, err
}

func (a *Actuator) machineProviderStatus(machine *clusterv1.Machine) (*v1alpha1.AWSMachineProviderStatus, error) {
	status := &v1alpha1.AWSMachineProviderStatus{}
	err := a.codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)
	return status, err
}

func (a *Actuator) clusterProviderStatus(cluster *clusterv1.Cluster) (*v1alpha1.AWSClusterProviderStatus, error) {
	providerStatus := &v1alpha1.AWSClusterProviderStatus{}
	err := a.codec.DecodeProviderStatus(cluster.Status.ProviderStatus, providerStatus)
	return providerStatus, err
}

func (a *Actuator) updateStatus(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus) error {
	machinesClient := a.machinesGetter.Machines(machine.Namespace)
	encodedProviderStatus, err := a.codec.EncodeProviderStatus(status)
	if err != nil {
		return fmt.Errorf("failed to encode machine status: %v", err)
	}
	if encodedProviderStatus != nil {
		machine.Status.ProviderStatus = encodedProviderStatus
		if _, err := machinesClient.UpdateStatus(machine); err != nil {
			return fmt.Errorf("failed to update machine status: %v", err)
		}
	}
	return nil
}

func (a *Actuator) updateMachine(machine *clusterv1.Machine) error {
	machinesClient := a.machinesGetter.Machines(machine.Namespace)

	_, err := machinesClient.Update(machine)
	if err != nil {
		return fmt.Errorf("failed to update machine: %v", err)
	}

	return nil
}
