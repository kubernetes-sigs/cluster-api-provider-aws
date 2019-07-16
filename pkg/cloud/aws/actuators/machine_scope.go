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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Client          client.Client
	Logger          logr.Logger
	Cluster         *clusterv1.Cluster
	Machine         *clusterv1.Machine
	ProviderMachine *v1alpha2.AWSMachine
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	// TODO get cluster
	clusterScope, err := NewClusterScope(ClusterScopeParams{
		AWSClients: params.AWSClients,
		Client:     params.Client,
		Cluster:    params.Cluster,
		Logger:     params.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &MachineScope{
		client:          params.Client,
		matchinePatch:   client.MergeFrom(params.Machine),
		Parent:          clusterScope,
		Machine:         params.Machine,
		ProviderMachine: params.ProviderMachine,
		Logger: clusterScope.Logger.
			WithName(fmt.Sprintf("machine=%s", params.Machine.Name)).
			WithName(fmt.Sprintf("providerMachine=%s", params.ProviderMachine.Name)),
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logr.Logger
	matchinePatch client.Patch
	client        client.Client

	Parent          *ClusterScope
	Machine         *clusterv1.Machine
	ProviderMachine *v1alpha2.AWSMachine
}

// Name returns the machine name.
func (m *MachineScope) Name() string {
	return m.Machine.Name
}

// Namespace returns the machine namespace.
func (m *MachineScope) Namespace() string {
	return m.Machine.Namespace
}

// IsControlPlane returns true if the machine is a control plane.
func (m *MachineScope) IsControlPlane() bool {
	return util.IsControlPlaneMachine(m.Machine)
}

// Role returns the machine role from the labels.
func (m *MachineScope) Role() string {
	if util.IsControlPlaneMachine(m.Machine) {
		return "control-plane"
	}
	return "node"
}

// Region returns the machine region.
func (m *MachineScope) Region() string {
	return m.Parent.Region()
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() {
	ctx := context.Background()

	if err := m.client.Patch(ctx, m.ProviderMachine, m.matchinePatch); err != nil {
		m.Logger.Error(err, "error patching object")
		return
	}

	if err := m.client.Status().Patch(ctx, m.ProviderMachine, m.matchinePatch); err != nil {
		m.Logger.Error(err, "error patching object status")
	}
}
