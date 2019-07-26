/*
Copyright 2018 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const apiEndpointPort = 6443

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	AWSClients
	Client  client.Client
	Logger  logr.Logger
	Cluster *clusterv1.Cluster
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each different actuator iteration.
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil cluster")
	}

	clusterConfig, err := v1alpha2.ClusterConfigFromProviderSpec(params.Cluster.Spec.ProviderSpec)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider config: %v", err)
	}

	clusterStatus, err := v1alpha2.ClusterStatusFromProviderStatus(params.Cluster.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider status: %v", err)
	}

	session, err := sessionForRegion(clusterConfig.Region)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	if params.AWSClients.EC2 == nil {
		params.AWSClients.EC2 = ec2.New(session)
	}

	if params.AWSClients.ELB == nil {
		params.AWSClients.ELB = elb.New(session)
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	return &ClusterScope{
		client:       params.Client,
		clusterPatch: client.MergeFrom(params.Cluster.DeepCopy()),

		AWSClients:    params.AWSClients,
		Cluster:       params.Cluster,
		ClusterConfig: clusterConfig,
		ClusterStatus: clusterStatus,
		Logger: params.Logger.
			WithName(params.Cluster.APIVersion).
			WithName(fmt.Sprintf("namespace=%s", params.Cluster.Namespace)).
			WithName(fmt.Sprintf("cluster=%s", params.Cluster.Name)),
	}, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client       client.Client
	clusterPatch client.Patch

	AWSClients
	Cluster       *clusterv1.Cluster
	ClusterConfig *v1alpha2.AWSClusterProviderSpec
	ClusterStatus *v1alpha2.AWSClusterProviderStatus
}

// Network returns the cluster network object.
func (s *ClusterScope) Network() *v1alpha2.Network {
	return &s.ClusterStatus.Network
}

// VPC returns the cluster VPC.
func (s *ClusterScope) VPC() *v1alpha2.VPCSpec {
	return &s.ClusterConfig.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ClusterScope) Subnets() v1alpha2.Subnets {
	return s.ClusterConfig.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ClusterScope) SecurityGroups() map[v1alpha2.SecurityGroupRole]v1alpha2.SecurityGroup {
	return s.ClusterStatus.Network.SecurityGroups
}

// Name returns the cluster name.
func (s *ClusterScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ClusterScope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *ClusterScope) Region() string {
	return s.ClusterConfig.Region
}

// ControlPlaneConfigMapName returns the name of the ConfigMap used to
// coordinate the bootstrapping of control plane nodes.
func (s *ClusterScope) ControlPlaneConfigMapName() string {
	return fmt.Sprintf("%s-controlplane", s.Cluster.UID)
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ClusterScope) ListOptionsLabelSelector() client.ListOptionFunc {
	return client.MatchingLabels(map[string]string{
		clusterv1.MachineClusterLabelName: s.Cluster.Name,
	})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() {
	ctx := context.Background()

	// Patch Cluster object.
	ext, err := v1alpha2.EncodeClusterSpec(s.ClusterConfig)
	if err != nil {
		s.Error(err, "failed encoding cluster spec")
		return
	}
	s.Cluster.Spec.ProviderSpec.Value = ext
	if err := s.client.Patch(ctx, s.Cluster, s.clusterPatch); err != nil {
		s.Error(err, "failed to patch object")
		return
	}

	// Patch Cluster status.
	newStatus, err := v1alpha2.EncodeClusterStatus(s.ClusterStatus)
	if err != nil {
		s.Error(err, "failed encoding cluster status")
		return
	}
	if s.ClusterStatus.Network.APIServerELB.DNSName != "" {
		s.Cluster.Status.APIEndpoints = []clusterv1.APIEndpoint{
			{
				Host: s.ClusterStatus.Network.APIServerELB.DNSName,
				Port: apiEndpointPort,
			},
		}
	}
	s.Cluster.Status.ProviderStatus = newStatus
	if err := s.client.Status().Patch(ctx, s.Cluster, s.clusterPatch); err != nil {
		s.Error(err, "failed to patch object status")
		return
	}
}
