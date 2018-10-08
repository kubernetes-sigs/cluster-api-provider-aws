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

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	// TagsLastAppliedAnnotation is the key for the machine object annotation
	// which tracks the SecurityGroups that the machine actuator is responsible
	// for. These are the SecurityGroups that have been handled by the
	// AdditionalTags in the Machine Provider Config.
	TagsLastAppliedAnnotation = "sigs.k8s.io/cluster-api-provider-aws/last-applied/tags"
)

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
