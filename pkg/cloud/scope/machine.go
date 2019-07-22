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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/klog"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/cluster-api/pkg/controller/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/pkg/errors"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	AWSClients
	Client          client.Client
	Logger          logr.Logger
	ProviderMachine *infrav1.AWSMachine
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	ctx := context.Background()

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, params.Client, params.ProviderMachine.ObjectMeta)
	if err != nil {
		return nil, err
	} else if machine == nil {
		klog.Infof("Waiting for Machine Controller to set OwnerRef on AWSMachine %q/%q",
			params.ProviderMachine.Namespace, params.ProviderMachine.Name)
		return nil, &capierrors.RequeueAfterError{RequeueAfter: 10 * time.Second}
	}

	// Fetch the cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, params.Client, machine.ObjectMeta)
	if err != nil {
		return nil, err
	}

	clusterScope, err := NewClusterScope(ClusterScopeParams{
		AWSClients: params.AWSClients,
		Client:     params.Client,
		Cluster:    cluster,
		Logger:     params.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &MachineScope{
		client:          params.Client,
		patch:           client.MergeFrom(params.ProviderMachine.DeepCopy()),
		Parent:          clusterScope,
		Machine:         machine,
		ProviderMachine: params.ProviderMachine,
		Logger: clusterScope.Logger.
			WithName(fmt.Sprintf("machine=%s", machine.Name)).
			WithName(fmt.Sprintf("providerMachine=%s", params.ProviderMachine.Name)),
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logr.Logger
	patch  client.Patch
	client client.Client

	Parent          *ClusterScope
	Machine         *clusterv1.Machine
	ProviderMachine *infrav1.AWSMachine
}

// Name returns the AWSMachine name.
func (m *MachineScope) Name() string {
	return m.ProviderMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.ProviderMachine.Namespace
}

// ClusterName returns the parent Cluster name.
func (m *MachineScope) ClusterName() string {
	return m.Parent.Name()
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
	if m.ProviderMachine.Spec.ProviderID != nil {
		return *m.ProviderMachine.Spec.ProviderID
	}
	return ""
}

// SetProviderID sets the AWSMachine providerID in spec.
func (m *MachineScope) SetProviderID(v string) {
	m.ProviderMachine.Spec.ProviderID = pointer.StringPtr(v)
}

// GetInstanceID returns the AWSMachine instance state from the status.
func (m *MachineScope) GetInstanceState() *infrav1.InstanceState {
	return m.ProviderMachine.Status.InstanceState
}

// SetInstanceID sets the AWSMachine instance id.
func (m *MachineScope) SetInstanceState(v infrav1.InstanceState) {
	m.ProviderMachine.Status.InstanceState = &v
}

// SetErrorMessage sets the AWSMachine status error message.
func (m *MachineScope) SetErrorMessage(v error) {
	m.ProviderMachine.Status.ErrorMessage = pointer.StringPtr(v.Error())
}

// SetErrorReason sets the AWSMachine status error reason.
func (m *MachineScope) SetErrorReason(v common.MachineStatusError) {
	m.ProviderMachine.Status.ErrorReason = &v
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.ProviderMachine.Annotations == nil {
		m.ProviderMachine.Annotations = map[string]string{}
	}
	m.ProviderMachine.Annotations[key] = value
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() {
	ctx := context.Background()

	if err := m.client.Patch(ctx, m.ProviderMachine, m.patch); err != nil {
		m.Logger.Error(err, "error patching object")
	}

	if err := m.client.Status().Patch(ctx, m.ProviderMachine, m.patch); err != nil {
		m.Logger.Error(err, "error patching object status")
	}
}
