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

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"

	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// ManagedClusterScopeParams defines the input parameters used to create a new Scope.
type ManagedClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *clusterv1.Cluster
	AWSManagedCluster *infrav1exp.AWSManagedCluster
	Controlplane      *infrav1exp.AWSManagedControlPlane
	Session           awsclient.ConfigProvider
	ControllerName    string
}

// NewManagedClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedClusterScope(params ManagedClusterScopeParams) (*ManagedClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.AWSManagedCluster == nil {
		return nil, errors.New("failed to generate new scope from nil AWSManagedCluster")
	}
	if params.Controlplane == nil {
		return nil, errors.New("failed to generate new scope from nil ControlPlane")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	session, err := sessionForRegion(params.AWSManagedCluster.Spec.Region)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.AWSManagedCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &ManagedClusterScope{
		Logger:            params.Logger,
		client:            params.Client,
		Cluster:           params.Cluster,
		AWSManagedCluster: params.AWSManagedCluster,
		Controlplane:      params.Controlplane,
		patchHelper:       helper,
		session:           session,
		controllerName:    params.ControllerName,
	}, nil
}

// ManagedClusterScope defines the basic context for an actuator to operate upon.
type ManagedClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster           *clusterv1.Cluster
	AWSManagedCluster *infrav1exp.AWSManagedCluster
	Controlplane      *infrav1exp.AWSManagedControlPlane

	session        awsclient.ConfigProvider
	controllerName string
}

// Network returns the cluster network object.
func (s *ManagedClusterScope) Network() *infrav1.Network {
	return &s.AWSManagedCluster.Status.Network
}

// VPC returns the cluster VPC.
func (s *ManagedClusterScope) VPC() *infrav1.VPCSpec {
	return &s.AWSManagedCluster.Spec.NetworkSpec.VPC
}

// Subnets returns the cluster subnets.
func (s *ManagedClusterScope) Subnets() infrav1.Subnets {
	return s.AWSManagedCluster.Spec.NetworkSpec.Subnets
}

// SetSubnets updates the clusters subnets.
func (s *ManagedClusterScope) SetSubnets(subnets infrav1.Subnets) {
	s.AWSManagedCluster.Spec.NetworkSpec.Subnets = subnets
}

// CNIIngressRules returns the CNI spec ingress rules.
func (s *ManagedClusterScope) CNIIngressRules() infrav1.CNIIngressRules {
	if s.AWSManagedCluster.Spec.NetworkSpec.CNI != nil {
		return s.AWSManagedCluster.Spec.NetworkSpec.CNI.CNIIngressRules
	}
	return infrav1.CNIIngressRules{}
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ManagedClusterScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return s.AWSManagedCluster.Status.Network.SecurityGroups
}

// Name returns the cluster name.
func (s *ManagedClusterScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ManagedClusterScope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *ManagedClusterScope) Region() string {
	return s.AWSManagedCluster.Spec.Region
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ManagedClusterScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1.ClusterLabelName: s.Cluster.Name,
	})
}

// PatchObject persists the cluster configuration and status.
func (s *ManagedClusterScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.AWSManagedCluster,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			infrav1.VpcReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.InternetGatewayReadyCondition,
			infrav1.NatGatewaysReadyCondition,
			infrav1.RouteTablesReadyCondition,
			infrav1.BastionHostReadyCondition,
		}})
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ManagedClusterScope) Close() error {
	return s.PatchObject()
}

// AdditionalTags returns AdditionalTags from the scope's AWSCluster. The returned value will never be nil.
func (s *ManagedClusterScope) AdditionalTags() infrav1.Tags {
	if s.AWSManagedCluster.Spec.AdditionalTags == nil {
		s.AWSManagedCluster.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.AWSManagedCluster.Spec.AdditionalTags.DeepCopy()
}

func (s *ManagedClusterScope) APIServerPort() int32 {
	return 443
}

// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
func (s *ManagedClusterScope) SetFailureDomain(id string, spec clusterv1.FailureDomainSpec) {
	if s.AWSManagedCluster.Status.FailureDomains == nil {
		s.AWSManagedCluster.Status.FailureDomains = make(clusterv1.FailureDomains)
	}
	s.AWSManagedCluster.Status.FailureDomains[id] = spec
}

func (s *ManagedClusterScope) InfraCluster() cloud.ClusterObject {
	return s.AWSManagedCluster
}

func (s *ManagedClusterScope) Session() awsclient.ConfigProvider {
	return s.session
}

func (s *ManagedClusterScope) Bastion() *infrav1.Bastion {
	return &s.AWSManagedCluster.Spec.Bastion
}

func (s *ManagedClusterScope) SetBastionInstance(instance *infrav1.Instance) {
	s.AWSManagedCluster.Status.Bastion = instance
}

func (s *ManagedClusterScope) SSHKeyName() *string {
	return s.AWSManagedCluster.Spec.SSHKeyName
}

// ControllerName returns the name of the controller that
// created the ManagedClusterScope.
func (s *ManagedClusterScope) ControllerName() string {
	return s.controllerName
}
