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

package scope

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/klog/klogr"

	"github.com/go-logr/logr"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/controller/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/pkg/errors"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Client  client.Client
	Logger  logr.Logger
	Cluster *clusterv1.Cluster
	//AWSCluster *infrav1.AWSCluster
	Machine    *clusterv1.Machine
	AWSMachine *infrav1.AWSMachine
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, errors.New("machine is required when creating a MachineScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("cluster is required when creating a MachineScope")
	}
	if params.AWSMachine == nil {
		return nil, errors.New("aws machine is required when creating a MachineScope")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	return &MachineScope{
		client:     params.Client,
		patch:      client.MergeFrom(params.AWSMachine.DeepCopy()),
		Cluster:    params.Cluster,
		Machine:    params.Machine,
		AWSMachine: params.AWSMachine,
		Logger:     params.Logger,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logr.Logger
	patch  client.Patch
	client client.Client

	Cluster    *clusterv1.Cluster
	Machine    *clusterv1.Machine
	AWSMachine *infrav1.AWSMachine
}

// Name returns the AWSMachine name.
func (m *MachineScope) Name() string {
	return m.AWSMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.AWSMachine.Namespace
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

// GetInstanceID returns the AWSMachine instance id by parsing Spec.ProviderID.
func (m *MachineScope) GetInstanceID() *string {
	parsed, err := noderefutil.NewProviderID(m.GetProviderID())
	if err != nil {
		return nil
	}
	return pointer.StringPtr(parsed.ID())
}

// GetProviderID returns the AWSMachine providerID from the spec.
func (m *MachineScope) GetProviderID() string {
	if m.AWSMachine.Spec.ProviderID != nil {
		return *m.AWSMachine.Spec.ProviderID
	}
	return ""
}

// SetProviderID sets the AWSMachine providerID in spec.
func (m *MachineScope) SetProviderID(v string) {
	m.AWSMachine.Spec.ProviderID = pointer.StringPtr(v)
}

// GetInstanceID returns the AWSMachine instance state from the status.
func (m *MachineScope) GetInstanceState() *infrav1.InstanceState {
	return m.AWSMachine.Status.InstanceState
}

// SetInstanceID sets the AWSMachine instance id.
func (m *MachineScope) SetInstanceState(v infrav1.InstanceState) {
	m.AWSMachine.Status.InstanceState = &v
}

// SetReady sets the AWSMachine Ready Status
func (m *MachineScope) SetReady() {
	m.AWSMachine.Status.Ready = true
}

// SetErrorMessage sets the AWSMachine status error message.
func (m *MachineScope) SetErrorMessage(v error) {
	m.AWSMachine.Status.ErrorMessage = pointer.StringPtr(v.Error())
}

// SetErrorReason sets the AWSMachine status error reason.
func (m *MachineScope) SetErrorReason(v capierrors.MachineStatusError) {
	m.AWSMachine.Status.ErrorReason = &v
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.AWSMachine.Annotations == nil {
		m.AWSMachine.Annotations = map[string]string{}
	}
	m.AWSMachine.Annotations[key] = value
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() error {
	ctx := context.Background()

	if err := m.client.Patch(ctx, m.AWSMachine, m.patch); err != nil {
		return errors.WithStack(err)
	}

	if err := m.client.Status().Patch(ctx, m.AWSMachine, m.patch); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
