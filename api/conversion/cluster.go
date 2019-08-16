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
	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	capav1a2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	capav1a1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	capiv1a2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1a1 "sigs.k8s.io/cluster-api/pkg/apis/deprecated/v1alpha1"
	"sigs.k8s.io/yaml"
)

// ConvertCluster turns a CAPI v1alpha1 cluster with embedded CAPA ProviderSpec into a AWSCluster and a v1alpha2 Cluster object that references it
func ConvertCluster(in *capiv1a1.Cluster) (*capiv1a2.Cluster, *capav1a2.AWSCluster, error) {
	var out capiv1a2.Cluster

	if err := capiv1a2.Convert_v1alpha1_Cluster_To_v1alpha2_Cluster(in, &out, nil); err != nil {
		return nil, nil, err
	}

	awsOut, err := getAWSCluster(&in.Spec)
	if err != nil {
		return nil, nil, err
	}

	if awsOut == nil {
		return &out, nil, nil
	}

	awsOut.Name = in.Name
	awsOut.Namespace = in.Namespace

	ref := corev1.ObjectReference{
		Name:       awsOut.Name,
		Namespace:  awsOut.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSCluster",
	}

	out.Spec.InfrastructureRef = &ref

	return &out, awsOut, nil
}

func getAWSCluster(in *capiv1a1.ClusterSpec) (*capav1a2.AWSCluster, error) {
	var (
		awsIn  capav1a1.AWSClusterProviderSpec
		awsOut capav1a2.AWSCluster
	)

	if in.ProviderSpec.Value == nil {
		return nil, nil
	}

	if err := yaml.Unmarshal(in.ProviderSpec.Value.Raw, &awsIn); err != nil {
		return nil, errors.Wrap(err, "couldn't decode ProviderSpec")
	}

	if err := capav1a2.Convert_v1alpha1_AWSClusterProviderSpec_To_v1alpha2_AWSClusterSpec(&awsIn, &awsOut.Spec, nil); err != nil {
		return nil, errors.Wrap(err, "couldn't convert ProviderSpec")
	}

	return &awsOut, nil
}
