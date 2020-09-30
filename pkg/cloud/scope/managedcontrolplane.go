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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
)

// ManagedControlPlaneScopeParams defines the input parameters used to create a new Scope.
type ManagedControlPlaneScopeParams struct {
	Client         client.Client
	Logger         logr.Logger
	Cluster        *clusterv1.Cluster
	ControlPlane   *controlplanev1.AWSManagedControlPlane
	ControllerName string
	Endpoints      []ServiceEndpoint
	Session        awsclient.ConfigProvider

	EnableIAM            bool
	AllowAdditionalRoles bool
}

// NewManagedControlPlaneScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewManagedControlPlaneScope(params ManagedControlPlaneScopeParams) (*ManagedControlPlaneScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSManagedControlPlane")
	}
	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	session, err := sessionForRegion(params.ControlPlane.Spec.Region, params.Endpoints)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.ControlPlane, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ManagedControlPlaneScope{
		Logger:               params.Logger,
		Client:               params.Client,
		Cluster:              params.Cluster,
		ControlPlane:         params.ControlPlane,
		patchHelper:          helper,
		session:              session,
		controllerName:       params.ControllerName,
		allowAdditionalRoles: params.AllowAdditionalRoles,
		enableIAM:            params.EnableIAM,
	}, nil
}

// ManagedControlPlaneScope defines the basic context for an actuator to operate upon.
type ManagedControlPlaneScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	ControlPlane *controlplanev1.AWSManagedControlPlane

	session        awsclient.ConfigProvider
	controllerName string

	enableIAM            bool
	allowAdditionalRoles bool
}

// Network returns the control plane network object.
func (s *ManagedControlPlaneScope) Network() *infrav1.Network {
	return &s.ControlPlane.Status.Network
}

// VPC returns the control plane VPC.
func (s *ManagedControlPlaneScope) VPC() *infrav1.VPCSpec {
	return &s.ControlPlane.Spec.NetworkSpec.VPC
}

// Subnets returns the control plane subnets.
func (s *ManagedControlPlaneScope) Subnets() infrav1.Subnets {
	return s.ControlPlane.Spec.NetworkSpec.Subnets
}

// SetSubnets updates the control planes subnets.
func (s *ManagedControlPlaneScope) SetSubnets(subnets infrav1.Subnets) {
	s.ControlPlane.Spec.NetworkSpec.Subnets = subnets
}

// CNIIngressRules returns the CNI spec ingress rules.
func (s *ManagedControlPlaneScope) CNIIngressRules() infrav1.CNIIngressRules {
	if s.ControlPlane.Spec.NetworkSpec.CNI != nil {
		return s.ControlPlane.Spec.NetworkSpec.CNI.CNIIngressRules
	}
	return infrav1.CNIIngressRules{}
}

// SecurityGroups returns the control plane security groups as a map, it creates the map if empty.
func (s *ManagedControlPlaneScope) SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup {
	return s.ControlPlane.Status.Network.SecurityGroups
}

// Name returns the CAPI cluster name.
func (s *ManagedControlPlaneScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ManagedControlPlaneScope) Namespace() string {
	return s.Cluster.Namespace
}

// Region returns the cluster region.
func (s *ManagedControlPlaneScope) Region() string {
	return s.ControlPlane.Spec.Region
}

// ListOptionsLabelSelector returns a ListOptions with a label selector for clusterName.
func (s *ManagedControlPlaneScope) ListOptionsLabelSelector() client.ListOption {
	return client.MatchingLabels(map[string]string{
		clusterv1.ClusterLabelName: s.Cluster.Name,
	})
}

// PatchObject persists the control plane configuration and status.
func (s *ManagedControlPlaneScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ControlPlane,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			infrav1.VpcReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.InternetGatewayReadyCondition,
			infrav1.NatGatewaysReadyCondition,
			infrav1.RouteTablesReadyCondition,
			infrav1.BastionHostReadyCondition,
			controlplanev1.EKSControlPlaneReadyCondition,
			controlplanev1.IAMControlPlaneRolesReadyCondition,
		}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *ManagedControlPlaneScope) Close() error {
	return s.PatchObject()
}

// AdditionalTags returns AdditionalTags from the scope's EksControlPlane. The returned value will never be nil.
func (s *ManagedControlPlaneScope) AdditionalTags() infrav1.Tags {
	if s.ControlPlane.Spec.AdditionalTags == nil {
		s.ControlPlane.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.ControlPlane.Spec.AdditionalTags.DeepCopy()
}

// APIServerPort returns the port to use when communicating with the API server
func (s *ManagedControlPlaneScope) APIServerPort() int32 {
	return 443
}

// SetFailureDomain sets the infrastructure provider failure domain key to the spec given as input.
func (s *ManagedControlPlaneScope) SetFailureDomain(id string, spec clusterv1.FailureDomainSpec) {
	if s.ControlPlane.Status.FailureDomains == nil {
		s.ControlPlane.Status.FailureDomains = make(clusterv1.FailureDomains)
	}
	s.ControlPlane.Status.FailureDomains[id] = spec
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *ManagedControlPlaneScope) InfraCluster() cloud.ClusterObject {
	return s.ControlPlane
}

// Session returns the AWS SDK session. Used for creating clients
func (s *ManagedControlPlaneScope) Session() awsclient.ConfigProvider {
	return s.session
}

// Bastion returns the bastion details.
func (s *ManagedControlPlaneScope) Bastion() *infrav1.Bastion {
	return &s.ControlPlane.Spec.Bastion
}

// SetBastionInstance sets the bastion instance in the status of the cluster.
func (s *ManagedControlPlaneScope) SetBastionInstance(instance *infrav1.Instance) {
	s.ControlPlane.Status.Bastion = instance
}

// SSHKeyName returns the SSH key name to use for instances.
func (s *ManagedControlPlaneScope) SSHKeyName() *string {
	return s.ControlPlane.Spec.SSHKeyName
}

// ControllerName returns the name of the controller that
// created the ManagedControlPlane.
func (s *ManagedControlPlaneScope) ControllerName() string {
	return s.controllerName
}

// TokenMethod returns the token method to use in the kubeconfig
func (s *ManagedControlPlaneScope) TokenMethod() controlplanev1.EKSTokenMethod {
	if s.ControlPlane.Spec.TokenMethod != nil {
		return *s.ControlPlane.Spec.TokenMethod
	}

	return controlplanev1.EKSTokenMethodIAMAuthenticator
}

// KubernetesClusterName is the name of the Kubernetes cluster. For the managed
// scope this is the different to the CAPI cluster name and is the EKS cluster name
func (s *ManagedControlPlaneScope) KubernetesClusterName() string {
	return s.ControlPlane.Spec.EKSClusterName
}

// EnableIAM indicates that reconciliation should create IAM roles
func (s *ManagedControlPlaneScope) EnableIAM() bool {
	return s.enableIAM
}

// AllowAdditionalRoles indicates if additional roles can be added to the created IAM roles
func (s *ManagedControlPlaneScope) AllowAdditionalRoles() bool {
	return s.allowAdditionalRoles
}

// ImageLookupFormat returns the format string to use when looking up AMIs
func (s *ManagedControlPlaneScope) ImageLookupFormat() string {
	return s.ControlPlane.Spec.ImageLookupFormat
}

// ImageLookupOrg returns the organization name to use when looking up AMIs
func (s *ManagedControlPlaneScope) ImageLookupOrg() string {
	return s.ControlPlane.Spec.ImageLookupOrg
}

// ImageLookupBaseOS returns the base operating system name to use when looking up AMIs
func (s *ManagedControlPlaneScope) ImageLookupBaseOS() string {
	return s.ControlPlane.Spec.ImageLookupBaseOS
}
