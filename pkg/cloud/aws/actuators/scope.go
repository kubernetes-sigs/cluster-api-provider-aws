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
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	sessionCache sync.Map
)

const apiEndpointPort = 6443

func sessionForRegion(region string) (*session.Session, error) {
	s, ok := sessionCache.Load(region)
	if ok {
		return s.(*session.Session), nil
	}

	ns, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}

	sessionCache.Store(region, ns)
	return ns, nil
}

// ScopeParams defines the input parameters used to create a new Scope.
type ScopeParams struct {
	AWSClients
	Client  client.Client
	Logger  logr.Logger
	Cluster *clusterv1.Cluster
}

// NewScope creates a new Scope from the supplied parameters.
// This is meant to be called for each different actuator iteration.
func NewScope(params ScopeParams) (*Scope, error) {
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
		params.Logger = klogr.New().WithName("default-logger")
	}

	return &Scope{
		client:        params.Client,
		AWSClients:    params.AWSClients,
		Cluster:       params.Cluster,
		ClusterConfig: clusterConfig,
		ClusterStatus: clusterStatus,
		Logger: params.Logger.
			WithName(params.Cluster.APIVersion).
			WithName(params.Cluster.Namespace).
			WithName(fmt.Sprintf("cluster=%s", params.Cluster.Name)),
	}, nil
}

// Scope defines the basic context for an actuator to operate upon.
type Scope struct {
	logr.Logger
	client       client.Client
	clusterPatch client.Patch

	AWSClients
	Cluster       *clusterv1.Cluster
	ClusterConfig *v1alpha2.AWSClusterProviderSpec
	ClusterStatus *v1alpha2.AWSClusterProviderStatus
}

// Network returns the cluster network object.
func (s *Scope) Network() *v1alpha2.Network {
	return &s.ClusterStatus.Network
}

// VPC returns the cluster VPC.
func (s *Scope) VPC() *v1alpha2.VPCSpec {
	return &s.ClusterConfig.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *Scope) Subnets() v1alpha2.Subnets {
	return s.ClusterConfig.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *Scope) SecurityGroups() map[v1alpha2.SecurityGroupRole]v1alpha2.SecurityGroup {
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

	// ext, err := v1alpha2.EncodeClusterSpec(s.ClusterConfig)
	// if err != nil {
	// 	s.Error(err, "failed encoding cluster spec")
	// 	return
	// }
	// newStatus, err := v1alpha2.EncodeClusterStatus(s.ClusterStatus)
	// if err != nil {
	// 	s.Error(err, "failed encoding cluster status")
	// 	return
	// }

	// s.Cluster.Spec.ProviderSpec.Value = ext

	// // Do not update Machine if nothing has changed
	// var p []byte // TEMPORARY
	// if len(p) != 0 {
	// 	pb, err := json.MarshalIndent(p, "", "  ")
	// 	if err != nil {
	// 		s.Error(err, "failed to json marshal patch")
	// 		return
	// 	}
	// 	s.Logger.V(1).Info("Patching cluster")
	// 	result, err := s.ClusterClient.Patch(s.Cluster.Name, types.JSONPatchType, pb)
	// 	if err != nil {
	// 		s.Error(err, "failed to patch cluster")
	// 		return
	// 	}
	// 	// Keep the resource version updated so the status update can succeed
	// 	s.Cluster.ResourceVersion = result.ResourceVersion
	// }

	// // Set the APIEndpoint.
	// if s.ClusterStatus.Network.APIServerELB.DNSName != "" {
	// 	s.Cluster.Status.APIEndpoints = []clusterv1.APIEndpoint{
	// 		{
	// 			Host: s.ClusterStatus.Network.APIServerELB.DNSName,
	// 			Port: apiEndpointPort,
	// 		},
	// 	}
	// }
	// s.Cluster.Status.ProviderStatus = newStatus

	// if !reflect.DeepEqual(s.Cluster.Status, s.ClusterCopy.Status) {
	// 	s.Logger.V(1).Info("updating cluster status")
	// 	if _, err := s.ClusterClient.UpdateStatus(s.Cluster); err != nil {
	// 		s.Error(err, "failed to update cluster status")
	// 		return
	// 	}
	// }
}
