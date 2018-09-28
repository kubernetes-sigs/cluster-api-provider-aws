// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"github.com/pkg/errors"
	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

// Context defines the context services accept and work within.
type Context interface {
	Cluster() *clusterv1.Cluster
	Client() client.ClusterInterface
	Config() *providerconfigv1.AWSClusterProviderConfig
	Status() *providerconfigv1.AWSClusterProviderStatus
}

// contextFromCluster returns a new cluster context.
func (a *Actuator) contextFromCluster(cluster *clusterv1.Cluster) (*contextImpl, error) {
	res := &contextImpl{
		cluster:        cluster,
		client:         a.clustersGetter.Clusters(cluster.Namespace),
		providerConfig: &providerconfigv1.AWSClusterProviderConfig{},
		providerStatus: &providerconfigv1.AWSClusterProviderStatus{},
	}

	if err := a.codec.DecodeFromProviderConfig(cluster.Spec.ProviderConfig, res.providerConfig); err != nil {
		return nil, errors.Errorf("failed to load cluster provider config: %v", err)
	}

	if err := a.codec.DecodeProviderStatus(cluster.Status.ProviderStatus, res.providerStatus); err != nil {
		return nil, errors.Errorf("failed to load cluster provider status: %v", err)
	}

	return res, nil
}

type contextImpl struct {
	cluster        *clusterv1.Cluster
	client         client.ClusterInterface
	providerConfig *providerconfigv1.AWSClusterProviderConfig
	providerStatus *providerconfigv1.AWSClusterProviderStatus
}

func (c *contextImpl) Cluster() *clusterv1.Cluster {
	return c.cluster
}

func (c *contextImpl) Client() client.ClusterInterface {
	return c.client
}

func (c *contextImpl) Config() *providerconfigv1.AWSClusterProviderConfig {
	return c.providerConfig
}

func (c *contextImpl) Status() *providerconfigv1.AWSClusterProviderStatus {
	return c.providerStatus
}
