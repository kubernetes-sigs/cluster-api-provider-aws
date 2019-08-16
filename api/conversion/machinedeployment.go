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

// ConvertMachineDeployment turns a CAPI v1alpha1 MachineDeployment with embedded AWS MachineProviderSpec into a
// CAPA v1alpha2 AWSMachine and a CAPI v1alpha2 MachineDeployment that references it
func ConvertMachineDeployment(in *capiv1a1.MachineDeployment) (*capiv1a2.MachineDeployment, *capav1a2.AWSMachineTemplate, error) {
	var (
		out    capiv1a2.MachineDeployment
		awsOut capav1a2.AWSMachineTemplate
	)

	if err := capiv1a2.Convert_v1alpha1_MachineDeployment_To_v1alpha2_MachineDeployment(in, &out, nil); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to convert MachineSet")
	}

	err := getAWSMachine(&in.Spec.Template.Spec, &awsOut.Spec.Template.Spec)
	if err != nil {
		return nil, nil, err
	}

	awsOut.Name = in.Name
	awsOut.Namespace = in.Namespace

	ref := corev1.ObjectReference{
		Name:       awsOut.Name,
		Namespace:  awsOut.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSMachineTemplate",
	}

	out.Spec.Template.Spec.InfrastructureRef = ref

	return &out, &awsOut, nil
}
