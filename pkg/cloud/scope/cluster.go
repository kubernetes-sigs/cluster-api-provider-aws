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

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	AWSClients
	Client     client.Client
	Logger     logr.Logger
	Cluster    *clusterv1.Cluster
	AWSCluster *infrav1.AWSCluster
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
		ec2Client := ec2.New(session)
		ec2Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.AWSCluster))
		params.AWSClients.EC2 = ec2Client
	}

	if params.AWSClients.ELB == nil {
		elbClient := elb.New(session)
		elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.AWSCluster))
		params.AWSClients.ELB = elbClient
	}

	if params.AWSClients.ResourceTagging == nil {
		resourceTagging := resourcegroupstaggingapi.New(session)
		resourceTagging.Handlers.Complete.PushBack(recordAWSPermissionsIssue(params.AWSCluster))
		params.AWSClients.ResourceTagging = resourceTagging
	}

	helper, err := patch.NewHelper(params.AWSCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &ClusterScope{
		Logger:      params.Logger,
		client:      params.Client,
		AWSClients:  params.AWSClients,
		Cluster:     params.Cluster,
		AWSCluster:  params.AWSCluster,
		patchHelper: helper,
	}, nil
}

func recordAWSPermissionsIssue(target runtime.Object) func(r *request.Request) {
	return func(r *request.Request) {
		if awsErr, ok := r.Error.(awserr.Error); ok {
			switch awsErr.Code() {
			case "AuthFailure", "UnauthorizedOperation":
				record.Warnf(target, awsErr.Code(), "Operation %s failed with a credentials or permission issue", r.Operation.Name)
			}
		}
	}
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	AWSClients
	Cluster    *clusterv1.Cluster
	AWSCluster *infrav1.AWSCluster
}

// Network returns the cluster network object.
func (s *ClusterScope) Network() *infrav1.Network {
	return &s.AWSCluster.Status.Network
}

// VPC returns the cluster VPC.
func (s *ClusterScope) VPC() *infrav1.VPCSpec {
	return &s.AWSCluster.Spec.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ClusterScope) Subnets() infrav1.Subnets {
	return s.AWSCluster.Spec.NetworkSpec.Subnets
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ClusterScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
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

// ControlPlaneLoadBalancer returns the AWSLoadBalancerSpec
func (s *ClusterScope) ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec {
	return s.AWSCluster.Spec.ControlPlaneLoadBalancer
}

// ControlPlaneLoadBalancerScheme returns the Classic ELB scheme (public or internal facing)
func (s *ClusterScope) ControlPlaneLoadBalancerScheme() infrav1.ClassicELBScheme {
	if s.ControlPlaneLoadBalancer() != nil && s.ControlPlaneLoadBalancer().Scheme != nil {
		return *s.ControlPlaneLoadBalancer().Scheme
	}
	return infrav1.ClassicELBSchemeInternetFacing
}

// ControlPlaneConfigMapName returns the name of the ConfigMap used to
// coordinate the bootstrapping of control plane nodes.
func (s *ClusterScope) ControlPlaneConfigMapName() string {
	return fmt.Sprintf("%s-controlplane", s.Cluster.UID)
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ClusterScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1.MachineClusterLabelName: s.Cluster.Name,
	})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.patchHelper.Patch(context.TODO(), s.AWSCluster)
}

// AdditionalTags returns AdditionalTags from the scope's AWSCluster. The returned value will never be nil.
func (s *ClusterScope) AdditionalTags() infrav1.Tags {
	if s.AWSCluster.Spec.AdditionalTags == nil {
		s.AWSCluster.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.AWSCluster.Spec.AdditionalTags.DeepCopy()
}

// APIServerPort returns the APIServerPort to use when creating the load balancer.
func (s *ClusterScope) APIServerPort() int64 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return int64(*s.Cluster.Spec.ClusterNetwork.APIServerPort)
	}
	return 6443
}
