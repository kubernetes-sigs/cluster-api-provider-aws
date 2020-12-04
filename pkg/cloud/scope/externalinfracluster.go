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
	"fmt"

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/klogr"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ExternalInfraClusterScopeParams defines the input parameters used to create a new Scope.
type ExternalInfraClusterScopeParams struct {
	Client               client.Client
	Logger               logr.Logger
	Cluster              *clusterv1.Cluster
	ExternalInfraCluster *unstructured.Unstructured
	ControllerName       string
	Endpoints            []ServiceEndpoint
	Session              awsclient.ConfigProvider
}

// NewExternalInfraClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewExternalInfraClusterScope(params ExternalInfraClusterScopeParams) (*ExternalInfraClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.ExternalInfraCluster == nil {
		return nil, errors.New("failed to generate new scope from nil ExternalInfraCluster")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	region, found, err := unstructured.NestedString(params.ExternalInfraCluster.Object, "spec", "region")
	if err != nil || !found {
		return nil, fmt.Errorf("error getting region: %w", err)
	}
	session, err := sessionForRegion(region, params.Endpoints)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.ExternalInfraCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &ExternalInfraClusterScope{
		Logger:               params.Logger,
		client:               params.Client,
		Cluster:              params.Cluster,
		ExternalInfraCluster: &ExternalInfraClusterObject{params.ExternalInfraCluster},
		patchHelper:          helper,
		session:              session,
		controllerName:       params.ControllerName,
	}, nil
}

// ExternalInfraClusterScope defines the basic context for an actuator to operate upon.
type ExternalInfraClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster              *clusterv1.Cluster
	ExternalInfraCluster *ExternalInfraClusterObject

	session        awsclient.ConfigProvider
	controllerName string
}

// Network returns the cluster network object.
func (s *ExternalInfraClusterScope) Network() *infrav1.Network {
	return nil
}

// VPC returns the cluster VPC.
func (s *ExternalInfraClusterScope) VPC() *infrav1.VPCSpec {
	return &infrav1.VPCSpec{}
}

// Subnets returns the cluster subnets.
func (s *ExternalInfraClusterScope) Subnets() infrav1.Subnets {
	return nil
}

// SetSubnets updates the clusters subnets.
func (s *ExternalInfraClusterScope) SetSubnets(subnets infrav1.Subnets) {
}

// CNIIngressRules returns the CNI spec ingress rules.
func (s *ExternalInfraClusterScope) CNIIngressRules() infrav1.CNIIngressRules {
	return infrav1.CNIIngressRules{}
}

// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
func (s *ExternalInfraClusterScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return nil
}

// Name returns the CAPI cluster name.
func (s *ExternalInfraClusterScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ExternalInfraClusterScope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *ExternalInfraClusterScope) Region() string {
	region, found, err := unstructured.NestedString(s.ExternalInfraCluster.Object, "spec", "region")
	if err != nil || !found {
		s.Error(err, "error getting region")
		return ""
	}
	return region
}

// KubernetesClusterName is the name of the Kubernetes cluster. For the cluster
// scope this is the same as the CAPI cluster name
func (s *ExternalInfraClusterScope) KubernetesClusterName() string {
	return s.Cluster.Name
}

// ControlPlaneLoadBalancer returns the AWSLoadBalancerSpec
func (s *ExternalInfraClusterScope) ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec {
	return nil
}

// ControlPlaneLoadBalancerScheme returns the Classic ELB scheme (public or internal facing)
func (s *ExternalInfraClusterScope) ControlPlaneLoadBalancerScheme() infrav1.ClassicELBScheme {
	if s.ControlPlaneLoadBalancer() != nil && s.ControlPlaneLoadBalancer().Scheme != nil {
		return *s.ControlPlaneLoadBalancer().Scheme
	}
	return infrav1.ClassicELBSchemeInternetFacing
}

// ControlPlaneConfigMapName returns the name of the ConfigMap used to
// coordinate the bootstrapping of control plane nodes.
func (s *ExternalInfraClusterScope) ControlPlaneConfigMapName() string {
	return fmt.Sprintf("%s-controlplane", s.Cluster.UID)
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ExternalInfraClusterScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1.ClusterLabelName: s.Cluster.Name,
	})
}

// PatchObject persists the cluster configuration and status.
func (s *ExternalInfraClusterScope) PatchObject() error {
	return nil
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ExternalInfraClusterScope) Close() error {
	return s.PatchObject()
}

// AdditionalTags returns AdditionalTags from the scope's ExternalInfraCluster. The returned value will never be nil.
func (s *ExternalInfraClusterScope) AdditionalTags() infrav1.Tags {
	return nil
}

// APIServerPort returns the APIServerPort to use when creating the load balancer.
func (s *ExternalInfraClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return 6443
}

// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
func (s *ExternalInfraClusterScope) SetFailureDomain(id string, spec clusterv1.FailureDomainSpec) {
}

type ExternalInfraClusterObject struct {
	*unstructured.Unstructured
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *ExternalInfraClusterScope) InfraCluster() cloud.ClusterObject {
	return s.ExternalInfraCluster
}

func (r *ExternalInfraClusterObject) GetConditions() clusterv1.Conditions {
	return nil
}

func (r *ExternalInfraClusterObject) SetConditions(conditions clusterv1.Conditions) {
}

// Session returns the AWS SDK session. Used for creating clients
func (s *ExternalInfraClusterScope) Session() awsclient.ConfigProvider {
	return s.session
}

// Bastion returns the bastion details.
func (s *ExternalInfraClusterScope) Bastion() *infrav1.Bastion {
	return nil
}

// SetBastionInstance sets the bastion instance in the status of the cluster.
func (s *ExternalInfraClusterScope) SetBastionInstance(instance *infrav1.Instance) {
}

// SSHKeyName returns the SSH key name to use for instances.
func (s *ExternalInfraClusterScope) SSHKeyName() *string {
	return nil
}

// ControllerName returns the name of the controller that
// created the ExternalInfraClusterScope.
func (s *ExternalInfraClusterScope) ControllerName() string {
	return s.controllerName
}

// ImageLookupFormat returns the format string to use when looking up AMIs
func (s *ExternalInfraClusterScope) ImageLookupFormat() string {
	return ""
}

// ImageLookupOrg returns the organization name to use when looking up AMIs
func (s *ExternalInfraClusterScope) ImageLookupOrg() string {
	return ""
}

// ImageLookupBaseOS returns the base operating system name to use when looking up AMIs
func (s *ExternalInfraClusterScope) ImageLookupBaseOS() string {
	return ""
}
