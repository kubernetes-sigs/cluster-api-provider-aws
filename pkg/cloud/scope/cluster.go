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
	clusterv1alpha2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	AWSClients
	Client     client.Client
	Logger     logr.Logger
	Cluster    *clusterv1alpha2.Cluster
	AWSCluster *v1alpha2.AWSCluster
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.AWSCluster == nil {
		return nil, errors.New("failed to generate new scope from nil AWSCluster")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	session, err := sessionForRegion(params.AWSCluster.Spec.Region)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	if params.AWSClients.EC2 == nil {
		params.AWSClients.EC2 = ec2.New(session)
	}

	if params.AWSClients.ELB == nil {
		params.AWSClients.ELB = elb.New(session)
	}

	return &ClusterScope{
		Logger:          params.Logger,
		client:          params.Client,
		AWSClients:      params.AWSClients,
		Cluster:         params.Cluster,
		AWSCluster:      params.AWSCluster,
		awsClusterPatch: client.MergeFrom(params.AWSCluster.DeepCopy()),
	}, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client client.Client
	AWSClients
	Cluster         *clusterv1alpha2.Cluster
	AWSCluster      *v1alpha2.AWSCluster
	awsClusterPatch client.Patch
}

// Network returns the cluster network object.
func (s *ClusterScope) Network() *v1alpha2.Network {
	return &s.AWSCluster.Status.Network
}

// VPC returns the cluster VPC.
func (s *ClusterScope) VPC() *v1alpha2.VPCSpec {
	return &s.AWSCluster.Spec.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ClusterScope) Subnets() v1alpha2.Subnets {
	return s.AWSCluster.Spec.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ClusterScope) SecurityGroups() map[v1alpha2.SecurityGroupRole]v1alpha2.SecurityGroup {
	return s.AWSCluster.Status.Network.SecurityGroups
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
	return s.AWSCluster.Spec.Region
}

// ControlPlaneConfigMapName returns the name of the ConfigMap used to
// coordinate the bootstrapping of control plane nodes.
func (s *ClusterScope) ControlPlaneConfigMapName() string {
	return fmt.Sprintf("%s-controlplane", s.Cluster.UID)
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ClusterScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1alpha2.MachineClusterLabelName: s.Cluster.Name,
	})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	ctx := context.TODO()

	// Patch Cluster object.
	if err := s.client.Patch(ctx, s.AWSCluster, s.awsClusterPatch); err != nil {
		return errors.Wrapf(err, "error patching AWSCluster %s/%s", s.Cluster.Namespace, s.Cluster.Name)
	}

	// Patch Cluster status.
	if err := s.client.Status().Patch(ctx, s.AWSCluster, s.awsClusterPatch); err != nil {
		return errors.Wrapf(err, "error patching AWSCluster %s/%s status", s.Cluster.Namespace, s.Cluster.Name)
	}

	return nil
}
