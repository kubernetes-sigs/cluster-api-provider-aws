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
	"encoding/base64"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/klogr"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Client       client.Client
	Logger       logr.Logger
	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	InfraCluster EC2Scope
	AWSMachine   *infrav1.AWSMachine
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
	if params.InfraCluster == nil {
		return nil, errors.New("aws cluster is required when creating a MachineScope")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	helper, err := patch.NewHelper(params.AWSMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &MachineScope{
		Logger:      params.Logger,
		client:      params.Client,
		patchHelper: helper,

		Cluster:      params.Cluster,
		Machine:      params.Machine,
		InfraCluster: params.InfraCluster,
		AWSMachine:   params.AWSMachine,
	}, nil
}

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	Machine      *clusterv1.Machine
	InfraCluster EC2Scope
	AWSMachine   *infrav1.AWSMachine
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
func (m *MachineScope) SetProviderID(instanceID, availabilityZone string) {
	providerID := fmt.Sprintf("aws:///%s/%s", availabilityZone, instanceID)
	m.AWSMachine.Spec.ProviderID = pointer.StringPtr(providerID)
}

// GetInstanceState returns the AWSMachine instance state from the status.
func (m *MachineScope) GetInstanceState() *infrav1.InstanceState {
	return m.AWSMachine.Status.InstanceState
}

// SetInstanceState sets the AWSMachine status instance state.
func (m *MachineScope) SetInstanceState(v infrav1.InstanceState) {
	m.AWSMachine.Status.InstanceState = &v
}

// SetReady sets the AWSMachine Ready Status
func (m *MachineScope) SetReady() {
	m.AWSMachine.Status.Ready = true
}

// SetNotReady sets the AWSMachine Ready Status to false
func (m *MachineScope) SetNotReady() {
	m.AWSMachine.Status.Ready = false
}

// SetFailureMessage sets the AWSMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.AWSMachine.Status.FailureMessage = pointer.StringPtr(v.Error())
}

// SetFailureReason sets the AWSMachine status failure reason.
func (m *MachineScope) SetFailureReason(v capierrors.MachineStatusError) {
	m.AWSMachine.Status.FailureReason = &v
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachineScope) SetAnnotation(key, value string) {
	if m.AWSMachine.Annotations == nil {
		m.AWSMachine.Annotations = map[string]string{}
	}
	m.AWSMachine.Annotations[key] = value
}

// UseSecretsManager returns the computed value of whether or not
// userdata should be stored using AWS Secrets Manager.
func (m *MachineScope) UseSecretsManager() bool {
	return !m.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager
}

// SecureSecretsBackend returns the chosen secret backend.
func (m *MachineScope) SecureSecretsBackend() infrav1.SecretBackend {
	return m.AWSMachine.Spec.CloudInit.SecureSecretsBackend
}

// UserDataIsUncompressed returns the computed value of whether or not
// userdata should be compressed using gzip.
func (m *MachineScope) UserDataIsUncompressed() bool {
	return m.AWSMachine.Spec.UncompressedUserData != nil && *m.AWSMachine.Spec.UncompressedUserData
}

// GetSecretPrefix returns the prefix for the secrets belonging
// to the AWSMachine in AWS Secrets Manager
func (m *MachineScope) GetSecretPrefix() string {
	return m.AWSMachine.Spec.CloudInit.SecretPrefix
}

// SetSecretPrefix sets the prefix for the secrets belonging
// to the AWSMachine in AWS Secrets Manager
func (m *MachineScope) SetSecretPrefix(value string) {
	m.AWSMachine.Spec.CloudInit.SecretPrefix = value
}

// DeleteSecretPrefix deletes the prefix for the secret belonging
// to the AWSMachine in AWS Secrets Manager
func (m *MachineScope) DeleteSecretPrefix() {
	m.AWSMachine.Spec.CloudInit.SecretPrefix = ""
}

// GetSecretCount returns the number of AWS Secret Manager entries making up
// the complete userdata
func (m *MachineScope) GetSecretCount() int32 {
	return m.AWSMachine.Spec.CloudInit.SecretCount
}

// SetSecretCount sets the number of AWS Secret Manager entries making up
// the complete userdata
func (m *MachineScope) SetSecretCount(i int32) {
	m.AWSMachine.Spec.CloudInit.SecretCount = i
}

// SetAddresses sets the AWSMachine address status.
func (m *MachineScope) SetAddresses(addrs []clusterv1.MachineAddress) {
	m.AWSMachine.Status.Addresses = addrs
}

// GetBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName as base64.
func (m *MachineScope) GetBootstrapData() (string, error) {
	value, err := m.GetRawBootstrapData()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(value), nil
}

// GetRawBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetRawBootstrapData() ([]byte, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.client.Get(context.TODO(), key, secret); err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, nil
}

// PatchObject persists the machine spec and status.
func (m *MachineScope) PatchObject() error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding during the deletion process).
	applicableConditions := []clusterv1.ConditionType{
		infrav1.InstanceReadyCondition,
		infrav1.SecurityGroupsReadyCondition,
	}

	if m.IsControlPlane() {
		applicableConditions = append(applicableConditions, infrav1.ELBAttachedCondition)
	}

	conditions.SetSummary(m.AWSMachine,
		conditions.WithConditions(applicableConditions...),
		conditions.WithStepCounterIf(m.AWSMachine.ObjectMeta.DeletionTimestamp.IsZero()),
		conditions.WithStepCounter(),
	)

	return m.patchHelper.Patch(
		context.TODO(),
		m.AWSMachine,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.InstanceReadyCondition,
			infrav1.SecurityGroupsReadyCondition,
			infrav1.ELBAttachedCondition,
		}})
}

// Close the MachineScope by updating the machine spec, machine status.
func (m *MachineScope) Close() error {
	return m.PatchObject()
}

// AdditionalTags merges AdditionalTags from the scope's AWSCluster and AWSMachine. If the same key is present in both,
// the value from AWSMachine takes precedence. The returned Tags will never be nil.
func (m *MachineScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)

	// Start with the cluster-wide tags...
	tags.Merge(m.InfraCluster.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(m.AWSMachine.Spec.AdditionalTags)

	return tags
}

func (m *MachineScope) HasFailed() bool {
	return m.AWSMachine.Status.FailureReason != nil || m.AWSMachine.Status.FailureMessage != nil
}

func (m *MachineScope) InstanceIsRunning() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceRunningStates.Has(string(*state))
}

func (m *MachineScope) InstanceIsOperational() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceOperationalStates.Has(string(*state))
}

func (m *MachineScope) InstanceIsInKnownState() bool {
	state := m.GetInstanceState()
	return state != nil && infrav1.InstanceKnownStates.Has(string(*state))
}

func (m *MachineScope) AWSMachineIsDeleted() bool {
	return !m.AWSMachine.ObjectMeta.DeletionTimestamp.IsZero()
}

func (m *MachineScope) IsEKSManaged() bool {
	return m.InfraCluster.InfraCluster().GetObjectKind().GroupVersionKind().Kind == "AWSManagedControlPlane"
}
