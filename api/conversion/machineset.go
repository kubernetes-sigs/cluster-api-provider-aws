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
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
)

type MachineSetConverter struct {
	machineSpecConverter

	oldMachineSet *capiv1a1.MachineSet
}

func NewMachineSetConverter(cluster *capiv1a1.Cluster, machineSet *capiv1a1.MachineSet) *MachineSetConverter {
	return &MachineSetConverter{
		machineSpecConverter: machineSpecConverter{
			ClusterConverter: ClusterConverter{
				oldCluster: cluster,
			},
			oldMachine: &machineSet.Spec.Template.Spec,
		},
		oldMachineSet: machineSet,
	}
}

func (m *MachineSetConverter) GetMachineSet(machineSet *capiv1a2.MachineSet) error {
	if err := capiv1a2.Convert_v1alpha1_MachineSet_To_v1alpha2_MachineSet(m.oldMachineSet, machineSet, nil); err != nil {
		return errors.Wrap(err, "Failed to convert MachineSet")
	}

	machineSet.Spec.Template.Spec.InfrastructureRef = corev1.ObjectReference{
		Name:       m.oldMachine.Name,
		Namespace:  m.oldMachine.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSMachineTemplate",
	}

	return nil
}

func (m *MachineSetConverter) GetAWSMachineTemplate(machineTemplate *capav1a2.AWSMachineTemplate) error {
	if err := m.getAWSMachineSpec(&machineTemplate.Spec.Template.Spec); err != nil {
		return errors.WithStack(err)
	}

	machineTemplate.Name = m.oldMachineSet.Name
	machineTemplate.Namespace = m.oldMachineSet.Namespace

	return nil
}
