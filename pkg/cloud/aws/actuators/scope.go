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

package actuators

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

// ScopeParams defines the input parameters used to create a new Scope.
type ScopeParams struct {
	AWSClients
	Cluster *clusterv1.Cluster
	Client  client.ClusterV1alpha1Interface
	Logger  logr.Logger
}

// NewScope creates a new Scope from the supplied parameters.
// This is meant to be called for each different actuator iteration.
func NewScope(params ScopeParams) (*Scope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil cluster")
	}

	clusterConfig, err := v1alpha1.ClusterConfigFromProviderSpec(params.Cluster.Spec.ProviderSpec)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider config: %v", err)
	}

	clusterStatus, err := v1alpha1.ClusterStatusFromProviderStatus(params.Cluster.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider status: %v", err)
	}

	session, err := session.NewSession(aws.NewConfig().WithRegion(clusterConfig.Region))
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	if params.AWSClients.EC2 == nil {
		params.AWSClients.EC2 = ec2.New(session)
	}

	if params.AWSClients.ELB == nil {
		params.AWSClients.ELB = elb.New(session)
	}

	var clusterClient client.ClusterInterface
	if params.Client != nil {
		clusterClient = params.Client.Clusters(params.Cluster.Namespace)
	}

	if params.Logger == nil {
		params.Logger = klogr.New().WithName("default-logger")
	}

	return &Scope{
		AWSClients:    params.AWSClients,
		Cluster:       params.Cluster,
		ClusterClient: clusterClient,
		ClusterConfig: clusterConfig,
		ClusterStatus: clusterStatus,
		Logger:        params.Logger.WithName(params.Cluster.APIVersion).WithName(params.Cluster.Namespace).WithName(params.Cluster.Name),
	}, nil
}

// Scope defines the basic context for an actuator to operate upon.
type Scope struct {
	AWSClients
	Cluster       *clusterv1.Cluster
	ClusterClient client.ClusterInterface
	ClusterConfig *v1alpha1.AWSClusterProviderSpec
	ClusterStatus *v1alpha1.AWSClusterProviderStatus
	logr.Logger
}

// Network returns the cluster network object.
func (s *Scope) Network() *v1alpha1.Network {
	return &s.ClusterStatus.Network
}

// VPC returns the cluster VPC.
func (s *Scope) VPC() *v1alpha1.VPCSpec {
	return &s.ClusterConfig.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *Scope) Subnets() v1alpha1.Subnets {
	return s.ClusterConfig.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *Scope) SecurityGroups() map[v1alpha1.SecurityGroupRole]*v1alpha1.SecurityGroup {
	return s.ClusterStatus.Network.SecurityGroups
}

// Name returns the cluster name.
func (s *Scope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *Scope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *Scope) Region() string {
	return s.ClusterConfig.Region
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *Scope) Close() {
	if s.ClusterClient == nil {
		return
	}
	ext, err := v1alpha1.EncodeClusterSpec(s.ClusterConfig)
	if err != nil {
		s.Error(err, "failed encoding cluster spec")
		return
	}
	status, err := v1alpha1.EncodeClusterStatus(s.ClusterStatus)
	if err != nil {
		s.Error(err, "failed encoding cluster status")
		return
	}

	// Update the resource version and try again if there is a conflict saving the object.
	if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		s.V(2).Info("Updating cluster", "cluster-resource-version", s.Cluster.ResourceVersion)
		s.Cluster.Spec.ProviderSpec.Value = ext
		s.V(6).Info("Cluster status before update", "cluster-status", s.Cluster.Status)
		latest, err := s.ClusterClient.Update(s.Cluster)
		if err != nil {
			s.V(3).Info("Cluster resource version is out of date")
			// Fetch and update the latest resource version
			newestCluster, err2 := s.ClusterClient.Get(s.Cluster.Name, metav1.GetOptions{})
			if err2 != nil {
				s.Error(err2, "failed to fetch latest cluster")
				return err2
			}
			s.Cluster.ResourceVersion = newestCluster.ResourceVersion
			return err
		}
		s.V(5).Info("Latest cluster status", "cluster-status", latest.Status)
		s.Cluster.Status.DeepCopyInto(&latest.Status)
		latest.Status.ProviderStatus = status
		_, err = s.ClusterClient.UpdateStatus(latest)
		return err
	}); err != nil {
		s.Error(err, "failed to retry on conflict")
	}
	s.V(2).Info("Successfully updated cluster")
}
