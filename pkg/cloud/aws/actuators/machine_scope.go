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
	"github.com/go-logr/logr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha2"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Cluster    *clusterv1.Cluster
	Machine    *clusterv1.Machine
	AWSMachine *v1alpha2.AWSMachine
	Client     client.ClusterV1alpha2Interface
	Logger     logr.Logger
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	scope, err := NewScope(ScopeParams{
		AWSClients: params.AWSClients,
		Client:     params.Client,
		Cluster:    params.Cluster,
		Logger:     params.Logger,
	})
	if err != nil {
		return nil, err
	}

	var machineClient client.MachineInterface
	if params.Client != nil {
		machineClient = params.Client.Machines(params.Machine.Namespace)
	}
	scope.Logger = scope.Logger.WithName(params.Machine.Name)
	return &MachineScope{
		Scope:         scope,
		Machine:       params.Machine,
		MachineCopy:   params.Machine.DeepCopy(),
		MachineClient: machineClient,
		MachineConfig: &params.AWSMachine.Spec,
		MachineStatus: &params.AWSMachine.Status,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	*Scope

	Machine *clusterv1.Machine
	// MachineCopy is used to generate a patch diff at the end of the scope's lifecycle.
	MachineCopy   *clusterv1.Machine
	MachineClient client.MachineInterface
	MachineConfig *v1alpha2.AWSMachineSpec
	MachineStatus *v1alpha2.AWSMachineStatus
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
	// TODO
}
