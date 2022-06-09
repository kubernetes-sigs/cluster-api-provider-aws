/*
Copyright 2022 The Kubernetes Authors.

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
	"github.com/pkg/errors"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
)

var (
	errClientRequired = errors.New("you must supply a client")
)

// RemoteScopeOption is an option when creating a RemoteClusterScope.
type RemoteScopeOption func(*RemoteClusterScope)

// WithRemoteScopeLogger is an option to set a specific logger.
func WithRemoteScopeLogger(logger *logr.Logger) RemoteScopeOption {
	return func(rcs *RemoteClusterScope) {
		rcs.Logger = *logger
	}
}

// WithRemoteScopeTenantClient is an option to use a pre-existing client
// for the remote/tenant/workload cluster.
func WithRemoteScopeTenantClient(client client.Client) RemoteScopeOption {
	return func(rcs *RemoteClusterScope) {
		rcs.Client = client
	}
}

// RemoteClusterScopeParams are inputs for creating a new RemoteClusterScope.
type RemoteClusterScopeParams struct {
	Client  client.Client
	Cluster *clusterv1.Cluster
	Logger  *logr.Logger
}

// NewRemoteClusterScope creates a new RemoteClusterScope.
func NewRemoteClusterScope(params RemoteClusterScopeParams, opts ...RemoteScopeOption) (*RemoteClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}

	if params.Logger == nil {
		log := klogr.New()
		params.Logger = &log
	}

	remoteScope := &RemoteClusterScope{
		Cluster: params.Cluster,
		Logger:  *params.Logger,
	}

	for _, opt := range opts {
		opt(remoteScope)
	}

	if remoteScope.Client == nil {
		if params.Client == nil {
			return nil, errClientRequired
		}

		clusterKey := client.ObjectKey{
			Name:      params.Cluster.Name,
			Namespace: params.Cluster.Namespace,
		}

		restConfig, err := remote.RESTConfig(context.Background(), params.Cluster.Name, params.Client, clusterKey)
		if err != nil {
			return nil, fmt.Errorf("getting remote rest config for %s/%s: %w", params.Cluster.Namespace, params.Cluster.Name, err)
		}
		restConfig.Timeout = 1 * time.Minute
		remoteClient, err := client.New(restConfig, client.Options{Scheme: scheme})
		if err != nil {
			return nil, fmt.Errorf("creating remote client: %w", err)
		}
		remoteScope.Client = remoteClient
	}

	return remoteScope, nil
}

// RemoteClusterScope is a scope that represents performing operations
// against a tenant/workload/child cluster.
type RemoteClusterScope struct {
	logr.Logger
	Client  client.Client
	Cluster *clusterv1.Cluster
}
