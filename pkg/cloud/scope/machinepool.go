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

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/klogr"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capierrors "sigs.k8s.io/cluster-api/errors"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachinePoolScope defines a scope defined around a machine and its cluster.
type MachinePoolScope struct {
	logr.Logger
	client      client.Client
	PatchHelper *patch.Helper

	Cluster        *clusterv1.Cluster
	MachinePool    *expclusterv1.MachinePool
	InfraCluster   EC2Scope
	AWSMachinePool *expinfrav1.AWSMachinePool
}

// MachinePoolScopeParams defines a scope defined around a machine and its cluster.
type MachinePoolScopeParams struct {
	Client client.Client
	Logger logr.Logger

	Cluster        *clusterv1.Cluster
	MachinePool    *expclusterv1.MachinePool
	InfraCluster   EC2Scope
	AWSMachinePool *expinfrav1.AWSMachinePool
}

// GetProviderID returns the AWSMachine providerID from the spec.
func (m *MachinePoolScope) GetProviderID() string {
	if m.AWSMachinePool.Spec.ProviderID != "" {
		return m.AWSMachinePool.Spec.ProviderID
	}
	return ""
}

// NewMachinePoolScope creates a new MachinePoolScope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachinePoolScope(params MachinePoolScopeParams) (*MachinePoolScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachinePoolScope")
	}
	if params.MachinePool == nil {
		return nil, errors.New("machinepool is required when creating a MachinePoolScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("cluster is required when creating a MachinePoolScope")
	}
	if params.AWSMachinePool == nil {
		return nil, errors.New("aws machine pool is required when creating a MachinePoolScope")
	}
	if params.InfraCluster == nil {
		return nil, errors.New("aws cluster is required when creating a MachinePoolScope")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	helper, err := patch.NewHelper(params.AWSMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachinePoolScope{
		Logger:      params.Logger,
		client:      params.Client,
		PatchHelper: helper,

		Cluster:        params.Cluster,
		MachinePool:    params.MachinePool,
		InfraCluster:   params.InfraCluster,
		AWSMachinePool: params.AWSMachinePool,
	}, nil
}

// Name returns the AWSMachinePool name.
func (m *MachinePoolScope) Name() string {
	return m.AWSMachinePool.Name
}

// Namespace returns the namespace name.
func (m *MachinePoolScope) Namespace() string {
	return m.AWSMachinePool.Namespace
}

// GetRawBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
// todo(rudoi): stolen from MachinePool - any way to reuse?
func (m *MachinePoolScope) GetRawBootstrapData() ([]byte, error) {
	if m.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		return nil, errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName}

	if err := m.client.Get(context.TODO(), key, secret); err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve bootstrap data secret for AWSMachine %s/%s", m.Namespace(), m.Name())
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, nil
}

// AdditionalTags merges AdditionalTags from the scope's AWSCluster and AWSMachinePool. If the same key is present in both,
// the value from AWSMachinePool takes precedence. The returned Tags will never be nil.
func (m *MachinePoolScope) AdditionalTags() infrav1.Tags {
	tags := make(infrav1.Tags)

	// Start with the cluster-wide tags...
	tags.Merge(m.InfraCluster.AdditionalTags())
	// ... and merge in the Machine's
	tags.Merge(m.AWSMachinePool.Spec.AdditionalTags)

	return tags
}

// PatchObject persists the machinepool spec and status.
func (m *MachinePoolScope) PatchObject() error {
	return m.PatchHelper.Patch(
		context.TODO(),
		m.AWSMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.ASGReadyCondition,
			expinfrav1.LaunchTemplateReadyCondition,
		}})
}

// Close the MachinePoolScope by updating the machinepool spec, machine status.
func (m *MachinePoolScope) Close() error {
	return m.PatchObject()
}

// SetAnnotation sets a key value annotation on the AWSMachine.
func (m *MachinePoolScope) SetAnnotation(key, value string) {
	if m.AWSMachinePool.Annotations == nil {
		m.AWSMachinePool.Annotations = map[string]string{}
	}
	m.AWSMachinePool.Annotations[key] = value
}

// SetFailureMessage sets the AWSMachine status failure message.
func (m *MachinePoolScope) SetFailureMessage(v error) {
	m.AWSMachinePool.Status.FailureMessage = pointer.StringPtr(v.Error())
}

// SetFailureReason sets the AWSMachine status failure reason.
func (m *MachinePoolScope) SetFailureReason(v capierrors.MachineStatusError) {
	m.AWSMachinePool.Status.FailureReason = &v
}

// HasFailed returns true when the AWSMachinePool's Failure reason or Failure message is populated
func (m *MachinePoolScope) HasFailed() bool {
	return m.AWSMachinePool.Status.FailureReason != nil || m.AWSMachinePool.Status.FailureMessage != nil
}

// SetNotReady sets the AWSMachinePool Ready Status to false
func (m *MachinePoolScope) SetNotReady() {
	m.AWSMachinePool.Status.Ready = false
}

// GetASGStatus returns the AWSMachinePool instance state from the status.
func (m *MachinePoolScope) GetASGStatus() *expinfrav1.ASGStatus {
	return m.AWSMachinePool.Status.ASGStatus
}

// SetASGStatus sets the AWSMachine status instance state.
func (m *MachinePoolScope) SetASGStatus(v expinfrav1.ASGStatus) {
	m.AWSMachinePool.Status.ASGStatus = &v
}

func (m *MachinePoolScope) IsEKSManaged() bool {
	return m.InfraCluster.InfraCluster().GetObjectKind().GroupVersionKind().Kind == "AWSManagedControlPlane"
}
