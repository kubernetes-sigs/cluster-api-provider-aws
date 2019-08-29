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

type ClusterConverter struct {
	oldCluster    *capiv1a1.Cluster
	oldAWSCluster *capav1a1.AWSClusterProviderSpec
}

func NewClusterConvert(cluster *capiv1a1.Cluster) *ClusterConverter {
	return &ClusterConverter{
		oldCluster:    cluster,
		oldAWSCluster: nil,
	}
}

func (c *ClusterConverter) GetCluster(cluster *capiv1a2.Cluster) error {
	if err := capiv1a2.Convert_v1alpha1_Cluster_To_v1alpha2_Cluster(c.oldCluster, cluster, nil); err != nil {
		return errors.WithStack(err)
	}

	ref := corev1.ObjectReference{
		Name:       c.oldCluster.Name,
		Namespace:  c.oldCluster.Namespace,
		APIVersion: capav1a2.GroupVersion.String(),
		Kind:       "AWSCluster",
	}

	cluster.Spec.InfrastructureRef = &ref

	return nil
}

func (c *ClusterConverter) getOldAWSCluster() (*capav1a1.AWSClusterProviderSpec, error) {
	if c.oldAWSCluster == nil {
		var oldAWSCluster capav1a1.AWSClusterProviderSpec
		if c.oldCluster.Spec.ProviderSpec.Value == nil {
			return nil, nil
		}

		if err := yaml.Unmarshal(c.oldCluster.Spec.ProviderSpec.Value.Raw, &oldAWSCluster); err != nil {
			return nil, errors.Wrap(err, "couldn't decode ProviderSpec")
		}

		c.oldAWSCluster = &oldAWSCluster
	}

	return c.oldAWSCluster, nil
}

func (c *ClusterConverter) GetAWSCluster(cluster *capav1a2.AWSCluster) error {
	oldCluster, err := c.getOldAWSCluster()
	if err != nil {
		return err
	}

	if oldCluster == nil {
		return nil
	}

	if err := capav1a2.Convert_v1alpha1_AWSClusterProviderSpec_To_v1alpha2_AWSClusterSpec(oldCluster, &cluster.Spec, nil); err != nil {
		return err
	}

	cluster.Name = c.oldCluster.Name
	cluster.Namespace = c.oldCluster.Namespace

	return nil
}
