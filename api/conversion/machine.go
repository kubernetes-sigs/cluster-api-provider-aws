/*
Copyright 2019 The Kubernetes Authors.

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

package conversion

import (
	"unsafe"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	cabpkv1a2 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	"sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/kubeadm/v1beta1"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
	"sigs.k8s.io/yaml"
)

type machineSpecConverter struct {
	ClusterConverter

	oldMachine    *capiv1a1.MachineSpec
	oldAWSMachine *capav1a1.AWSMachineProviderSpec
}

func (m *machineSpecConverter) getOldAWSMachine() (*capav1a1.AWSMachineProviderSpec, error) {
	if m.oldAWSMachine == nil {
		var oldAWSMachine capav1a1.AWSMachineProviderSpec
		if m.oldMachine.ProviderSpec.Value == nil {
			return nil, nil
		}

		if err := yaml.Unmarshal(m.oldMachine.ProviderSpec.Value.Raw, &oldAWSMachine); err != nil {
			return nil, errors.Wrap(err, "couldn't decode ProviderSpec")
		}

		m.oldAWSMachine = &oldAWSMachine
	}

	return m.oldAWSMachine, nil
}

func (m *machineSpecConverter) getAWSMachineSpec(spec *capav1a2.AWSMachineSpec) error {
	oldMachine, err := m.getOldAWSMachine()
	if err != nil {
		return errors.WithStack(err)
	}

	if oldMachine == nil {
		return nil
	}

	if err := capav1a2.Convert_v1alpha1_AWSMachineProviderSpec_To_v1alpha2_AWSMachineSpec(oldMachine, spec, nil); err != nil {
		return errors.Wrap(err, "couldn't convert ProviderSpec")
	}

	spec.ProviderID = m.oldMachine.ProviderID
	return nil
}

type MachineConverter struct {
	machineSpecConverter

	oldMachine *capiv1a1.Machine
}

func NewMachineConverter(cluster *capiv1a1.Cluster, machine *capiv1a1.Machine) *MachineConverter {
	return &MachineConverter{
		machineSpecConverter: machineSpecConverter{
			ClusterConverter: ClusterConverter{
				oldCluster: cluster,
			},
			oldMachine: &machine.Spec,
		},
		oldMachine: machine,
	}
}

func (m *MachineConverter) GetMachine(machine *capiv1a2.Machine) error {
	if err := capiv1a2.Convert_v1alpha1_Machine_To_v1alpha2_Machine(m.oldMachine, machine, nil); err != nil {
		return errors.Wrap(err, "failed to CAPI machine")
	}

	machine.Spec.InfrastructureRef = corev1.ObjectReference{
		Name:       m.oldMachine.Name,
		Namespace:  m.oldMachine.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSMachine",
	}

	return nil
}

func (m *MachineConverter) GetAWSMachine(machine *capav1a2.AWSMachine) error {
	if err := m.getAWSMachineSpec(&machine.Spec); err != nil {
		return err
	}

	machine.Name = m.oldMachine.Name
	machine.Namespace = m.oldMachine.Namespace

	return nil
}

func (m *MachineConverter) GetKubeadmConfig(cfg *cabpkv1a2.KubeadmConfig) error {
	spec, err := m.getOldAWSMachine()
	if err != nil {
		return err
	}

	cluster, err := m.getOldAWSCluster()
	if err != nil {
		return err
	}

	// I don't like this but the types are equivalent, this should be safe-ish

	cfg.Spec.InitConfiguration = (*v1beta1.InitConfiguration)(unsafe.Pointer(&spec.KubeadmConfiguration.Init))
	cfg.Spec.JoinConfiguration = (*v1beta1.JoinConfiguration)(unsafe.Pointer(&spec.KubeadmConfiguration.Join))
	cfg.Spec.ClusterConfiguration = (*v1beta1.ClusterConfiguration)(unsafe.Pointer(&cluster.ClusterConfiguration))

	cfg.Spec.Files = make([]cabpkv1a2.File, len(cluster.AdditionalUserDataFiles))
	for i, oldFile := range cluster.AdditionalUserDataFiles {
		newFile := &cfg.Spec.Files[i]
		newFile.Path = oldFile.Path
		newFile.Owner = oldFile.Owner
		newFile.Permissions = oldFile.Permissions
		newFile.Content = oldFile.Content
	}

	return nil
}
