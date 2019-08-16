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
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
	"sigs.k8s.io/yaml"
)

// ConvertMachine turns a CAPI v1alpha1 cluster with embedded CAPA ProviderSpec into a AWSMachine and a v1alpha2 Machine object that references it
func ConvertMachine(in *capiv1a1.Machine) (*capiv1a2.Machine, *capav1a2.AWSMachine, error) {
	var (
		out    capiv1a2.Machine
		awsOut capav1a2.AWSMachine
	)

	if err := capiv1a2.Convert_v1alpha1_Machine_To_v1alpha2_Machine(in, &out, nil); err != nil {
		return nil, nil, errors.Wrap(err, "failed to CAPI machine")
	}

	err := getAWSMachine(&in.Spec, &awsOut.Spec)
	if err != nil {
		return nil, nil, err
	}

	awsOut.Spec.ProviderID = out.Spec.ProviderID

	awsOut.Name = in.Name
	awsOut.Namespace = in.Namespace

	ref := corev1.ObjectReference{
		Name:       awsOut.Name,
		Namespace:  awsOut.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSMachine",
	}

	out.Spec.InfrastructureRef = ref

	return &out, &awsOut, nil
}

func getAWSMachine(in *capiv1a1.MachineSpec, out *capav1a2.AWSMachineSpec) error {
	var awsIn capav1a1.AWSMachineProviderSpec

	if in.ProviderSpec.Value == nil {
		return nil
	}

	if err := yaml.Unmarshal(in.ProviderSpec.Value.Raw, &awsIn); err != nil {
		return errors.Wrap(err, "couldn't decode providerSpec")
	}

	if err := capav1a2.Convert_v1alpha1_AWSMachineProviderSpec_To_v1alpha2_AWSMachineSpec(&awsIn, out, nil); err != nil {
		return errors.Wrap(err, "couldn't convert ProviderSpec")
	}

	return nil
}
