/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s *Service) reconcileAuthenticator(ctx context.Context) error {
	s.scope.V(2).Info("Reconciling aws-iam-authenticator configuration", "cluster-name", s.scope.KubernetesClusterName())

	clusterKey := client.ObjectKey{
		Name:      s.scope.Cluster.Name,
		Namespace: s.scope.Cluster.Namespace,
	}

	restConfig, err := remote.RESTConfig(ctx, s.scope.Client, clusterKey)
	if err != nil {
		return fmt.Errorf("getting remote client for %s/%s: %w", s.scope.Cluster.Namespace, s.scope.Cluster.Name, err)
	}

	remoteClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

}
