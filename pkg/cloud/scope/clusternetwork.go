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

package scope

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/version"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterNetworkScopeParams defines the input parameters used to create a new Scope.
type ClusterNetworkScopeParams struct {
	AWSClients
	Client  client.Client
	Logger  logr.Logger
	Cluster *clusterv1.Cluster
	Target  runtime.Object

	Region         string
	NetworkSpec    *infrav1.NetworkSpec
	NetworkStatus  *infrav1.Network
	AdditionalTags infrav1.Tags
	APIServerPort  *int32
}

// NewClusterNetworkScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterNetworkScope(params ClusterNetworkScopeParams) (*ClusterNetworkScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.Target == nil {
		return nil, errors.New("failed to generate new scope from nil Target")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	session, err := sessionForRegion(params.Region)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	userAgentHandler := request.NamedHandler{
		Name: "capa/user-agent",
		Fn:   request.MakeAddToUserAgentHandler("aws.cluster.x-k8s.io", version.Get().String()),
	}

	if params.AWSClients.EC2 == nil {
		ec2Client := ec2.New(session)
		ec2Client.Handlers.Build.PushFrontNamed(userAgentHandler)
		ec2Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.Target))
		params.AWSClients.EC2 = ec2Client
	}

	if params.AWSClients.ELB == nil {
		elbClient := elb.New(session)
		elbClient.Handlers.Build.PushFrontNamed(userAgentHandler)
		elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.Target))
		params.AWSClients.ELB = elbClient
	}

	if params.AWSClients.ResourceTagging == nil {
		resourceTagging := resourcegroupstaggingapi.New(session)
		resourceTagging.Handlers.Build.PushFrontNamed(userAgentHandler)
		resourceTagging.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.Target))
		params.AWSClients.ResourceTagging = resourceTagging
	}

	if params.AWSClients.SecretsManager == nil {
		sClient := secretsmanager.New(session)
		sClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.Target))
		params.AWSClients.SecretsManager = sClient
	}

	helper, err := patch.NewHelper(params.Target, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &ClusterNetworkScope{
		Logger:      params.Logger,
		client:      params.Client,
		AWSClients:  params.AWSClients,
		Cluster:     params.Cluster,
		Target:      params.Target,
		patchHelper: helper,

		region:        params.Region,
		NetworkSpec:   params.NetworkSpec,
		networkStatus: params.NetworkStatus,
		apiServerPort: params.APIServerPort,
		tags:          params.AdditionalTags,
	}, nil
}

// ClusterNetworkScope defines the basic context for an actuator to operate upon.
type ClusterNetworkScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	AWSClients
	Cluster *clusterv1.Cluster
	Target  runtime.Object

	region string

	NetworkSpec   *infrav1.NetworkSpec
	networkStatus *infrav1.Network

	apiServerPort *int32
	tags          infrav1.Tags
}

// Network returns the cluster network object.
func (s *ClusterNetworkScope) Network() *infrav1.Network {
	return s.networkStatus
}

// VPC returns the cluster VPC.
func (s *ClusterNetworkScope) VPC() *infrav1.VPCSpec {
	return &s.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ClusterNetworkScope) Subnets() infrav1.Subnets {
	return s.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ClusterNetworkScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return s.networkStatus.SecurityGroups
}

// Name returns the cluster name.
func (s *ClusterNetworkScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ClusterNetworkScope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *ClusterNetworkScope) Region() string {
	return s.region
}

// PatchObject persists the cluster configuration and status.
func (s *ClusterNetworkScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.Target)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterNetworkScope) Close() error {
	return s.PatchObject()
}

// AdditionalTags returns AdditionalTags from the scope's AWSCluster. The returned value will never be nil.
func (s *ClusterNetworkScope) AdditionalTags() infrav1.Tags {
	if s.tags == nil {
		s.tags = infrav1.Tags{}
	}

	return s.tags.DeepCopy()
}

// APIServerPort returns the APIServerPort to use when creating the load balancer.
func (s *ClusterNetworkScope) APIServerPort() int32 {
	if s.apiServerPort != nil {
		return *s.apiServerPort
	}
	return 6443
}
