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

package actuators

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"sigs.k8s.io/yaml"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Cluster *clusterv1.Cluster
	Machine *clusterv1.Machine
	Client  client.ClusterV1alpha1Interface
	Logger  logr.Logger
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	scope, err := NewScope(ScopeParams{
		AWSClients: params.AWSClients,
		Client:     params.Client, Cluster: params.Cluster,
		Logger: params.Logger,
	})
	if err != nil {
		return nil, err
	}

	machineConfig, err := MachineConfigFromProviderSpec(params.Client, params.Machine.Spec.ProviderSpec, scope.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine config")
	}

	machineStatus, err := v1alpha1.MachineStatusFromProviderStatus(params.Machine.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine provider status")
	}

	var machineClient client.MachineInterface
	if params.Client != nil {
		machineClient = params.Client.Machines(params.Machine.Namespace)
	}
	scope.Logger = scope.Logger.WithName(params.Machine.Name)
	return &MachineScope{
		Scope:         scope,
		Machine:       params.Machine,
		MachineClient: machineClient,
		MachineConfig: machineConfig,
		MachineStatus: machineStatus,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	*Scope

	Machine       *clusterv1.Machine
	MachineClient client.MachineInterface
	MachineConfig *v1alpha1.AWSMachineProviderSpec
	MachineStatus *v1alpha1.AWSMachineProviderStatus
}

// Name returns the machine name.
func (m *MachineScope) Name() string {
	return m.Machine.Name
}

// Namespace returns the machine namespace.
func (m *MachineScope) Namespace() string {
	return m.Machine.Namespace
}

// Role returns the machine role from the labels.
func (m *MachineScope) Role() string {
	return m.Machine.Labels["set"]
}

// Region returns the machine region.
func (m *MachineScope) Region() string {
	return m.Scope.Region()
}

// GetMachine returns the machine wrapped in the scope.
func (m *MachineScope) GetMachine() *clusterv1.Machine {
	return m.Machine
}

// GetScope returns the scope that is wrapping the machine.
func (m *MachineScope) GetScope() *Scope {
	return m.Scope
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() {
	if m.MachineClient == nil {
		return
	}
	ext, err := v1alpha1.EncodeMachineSpec(m.MachineConfig)
	if err != nil {
		m.Error(err, "failed to encode machine spec")
		return
	}
	status, err := v1alpha1.EncodeMachineStatus(m.MachineStatus)
	if err != nil {
		m.Error(err, "failed to encode machine status")
		return
	}

	// Sometimes when an object gets updated the local copy is out of date with
	// the copy stored on the server. In the case of cluster-api this will
	// always be because the local copy will have an out-of-date resource
	// version. This is because something else has updated the resource version
	// on the server and thus the local copy is behind.
	// This retry function will update the resource version if the local copy is
	// behind and try again.
	// This retry function will *only* update the resource version. If some
	// other data has changed then there is a problem. Nothing else should be
	// updating the object and this function will (correctly) fail.
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		m.V(2).Info("Updating machine", "machine-resource-version", m.Machine.ResourceVersion, "node-ref", m.Machine.Status.NodeRef)
		m.Machine.Spec.ProviderSpec.Value = ext
		m.V(6).Info("Machine status before update", "machine-status", m.Machine.Status)
		latest, err := m.MachineClient.Update(m.Machine)
		if err != nil {
			m.V(3).Info("Machine resource version is out of date")
			// Fetch and update the latest resource version
			newestMachine, err2 := m.MachineClient.Get(m.Machine.Name, metav1.GetOptions{})
			if err2 != nil {
				m.Error(err2, "failed to fetch latest Machine")
				return err2
			}
			m.Machine.ResourceVersion = newestMachine.ResourceVersion
			return err
		}
		m.V(5).Info("Latest machine", "machine", latest)
		// The machine may have status (nodeRef) that the latest doesn't yet
		// have, however some timestamps may be rolled back a bit with this copy.
		m.Machine.Status.DeepCopyInto(&latest.Status)
		latest.Status.ProviderStatus = status
		_, err = m.MachineClient.UpdateStatus(latest)
		return err
	}); err != nil {
		m.Error(err, "error retrying on conflict")
	}
	m.V(2).Info("Successfully updated machine")
}

// MachineConfigFromProviderSpec tries to decode the JSON-encoded spec, falling back on getting a MachineClass if the value is absent.
func MachineConfigFromProviderSpec(clusterClient client.MachineClassesGetter, providerConfig clusterv1.ProviderSpec, log logr.Logger) (*v1alpha1.AWSMachineProviderSpec, error) {
	var config v1alpha1.AWSMachineProviderSpec
	if providerConfig.Value != nil {
		log.V(4).Info("Decoding ProviderConfig from Value")
		return unmarshalProviderSpec(providerConfig.Value, log)
	}

	if providerConfig.ValueFrom != nil && providerConfig.ValueFrom.MachineClass != nil {
		ref := providerConfig.ValueFrom.MachineClass
		log.V(4).Info("Decoding ProviderConfig from MachineClass")
		log.V(6).Info("Machine class reference", "ref", fmt.Sprintf("%+v", ref))
		if ref.Provider != "" && ref.Provider != "aws" {
			return nil, errors.Errorf("Unsupported provider: %q", ref.Provider)
		}

		if len(ref.Namespace) > 0 && len(ref.Name) > 0 {
			log.V(4).Info("Getting MachineClass", "reference-namespace", ref.Namespace, "reference-name", ref.Name)
			mc, err := clusterClient.MachineClasses(ref.Namespace).Get(ref.Name, metav1.GetOptions{})
			log.V(6).Info("Retrieved MachineClass", "machine-class", fmt.Sprintf("%+v", mc))
			if err != nil {
				return nil, err
			}
			providerConfig.Value = &mc.ProviderSpec
			return unmarshalProviderSpec(&mc.ProviderSpec, log)
		}
	}

	return &config, nil
}

func unmarshalProviderSpec(spec *runtime.RawExtension, log logr.Logger) (*v1alpha1.AWSMachineProviderSpec, error) {
	var config v1alpha1.AWSMachineProviderSpec
	if spec != nil {
		if err := yaml.Unmarshal(spec.Raw, &config); err != nil {
			return nil, err
		}
	}
	log.V(6).Info("Found ProviderSpec", "provider-spec", fmt.Sprintf("%+v", config))
	return &config, nil
}
