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
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/annotations"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ExternalResourceGCScopeOption is an option when creating a ExternalResourceGCScope.
type ExternalResourceGCScopeOption func(*ExternalResourceGCScope)

// WithExternalResourceGCScopeLogger is an option to set a specific logger.
func WithExternalResourceGCScopeLogger(logger *logr.Logger) ExternalResourceGCScopeOption {
	return func(rcs *ExternalResourceGCScope) {
		rcs.Logger = *logger
	}
}

// WithExternalResourceGCScopeTenantClient is an option to use a pre-existing client
// for the remote/tenant/workload cluster.
func WithExternalResourceGCScopeTenantClient(client client.Client) ExternalResourceGCScopeOption {
	return func(rcs *ExternalResourceGCScope) {
		rcs.remoteClient = client
	}
}

// ExternalResourceGCScopeParams are inputs for creating a new ExternalResourceGCScope.
type ExternalResourceGCScopeParams struct {
	Client       client.Client
	RemoteClient client.Client
	Cluster      *clusterv1.Cluster
	InfraCluster client.Object
	Logger       *logr.Logger
}

// NewExternalResourceGCScope creates a new ExternalResourceGCScope.
func NewExternalResourceGCScope(params ExternalResourceGCScopeParams, opts ...ExternalResourceGCScopeOption) (*ExternalResourceGCScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.InfraCluster == nil {
		return nil, errors.New("failed to generate new scope from nil infra cluster")
	}
	if params.Client == nil {
		return nil, errors.New("failed to generate new scope from nil client")
	}

	if params.Logger == nil {
		log := klogr.New()
		params.Logger = &log
	}

	scopeLogger := params.Logger.WithValues(
		"name", params.InfraCluster.GetName(),
		"namespace", params.InfraCluster.GetNamespace(),
		"kind", params.InfraCluster.GetObjectKind(),
	)

	extScope := &ExternalResourceGCScope{
		cluster:      params.Cluster,
		client:       params.Client,
		remoteClient: params.RemoteClient,
		infraCluster: params.InfraCluster,
		Logger:       scopeLogger,
	}

	for _, opt := range opts {
		opt(extScope)
	}

	helper, err := patch.NewHelper(extScope.infraCluster, extScope.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	extScope.patchHelper = helper

	return extScope, nil
}

// ExternalResourceGCScope is a scope that represents performing operations
// in relation to garbage collecting resources in a tenant/workload/child cluster.
type ExternalResourceGCScope struct {
	logr.Logger

	client       client.Client
	remoteClient client.Client
	cluster      *clusterv1.Cluster
	infraCluster client.Object

	patchHelper *patch.Helper
}

// RemoteClient will get the client for the workload/tenant/remote cluster.
func (s *ExternalResourceGCScope) RemoteClient() (client.Client, error) {
	if s.remoteClient == nil {
		clusterKey := client.ObjectKey{
			Name:      s.cluster.Name,
			Namespace: s.cluster.Namespace,
		}

		restConfig, err := remote.RESTConfig(context.Background(), s.cluster.Name, s.client, clusterKey)
		if err != nil {
			return nil, fmt.Errorf("getting remote rest config for %s/%s: %w", s.cluster.Namespace, s.cluster.Name, err)
		}
		restConfig.Timeout = 1 * time.Minute
		remoteClient, err := client.New(restConfig, client.Options{Scheme: scheme})
		if err != nil {
			return nil, fmt.Errorf("creating remote client: %w", err)
		}
		s.remoteClient = remoteClient
	}

	return s.remoteClient, nil
}

// InfraCluster will get the infra cluster.
func (s *ExternalResourceGCScope) InfraCluster() client.Object {
	return s.infraCluster
}

// PatchObject persists the infra cluster configuration and status.
func (s *ExternalResourceGCScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.infraCluster)
}

// Close closes the current scope persisting the infra cluster configuration and status.
func (s *ExternalResourceGCScope) Close() error {
	return s.PatchObject()
}

// ShouldGarbageCollect returns whether a cluster should be garbage collected. Garbage collection is using an opt-out model
// initially if the feature is enabled.
func (s *ExternalResourceGCScope) ShouldGarbageCollect() (bool, error) {
	val, found := annotations.Get(s.infraCluster, annotations.ExternalResourceGCAnnotation)
	if !found {
		return true, nil
	}

	shouldGC, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("failed to convert annotation %s to a bool: %w", annotations.ExternalResourceGCAnnotation, err)
	}

	return shouldGC, nil
}
