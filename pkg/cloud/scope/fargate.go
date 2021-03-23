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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/throttle"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	controlplanev1exp "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
)

// FargateProfileScopeParams defines the input parameters used to create a new Scope.
type FargateProfileScopeParams struct {
	Client         client.Client
	Logger         logr.Logger
	Cluster        *clusterv1.Cluster
	ControlPlane   *controlplanev1exp.AWSManagedControlPlane
	FargateProfile *infrav1exp.AWSFargateProfile
	ControllerName string
	Endpoints      []ServiceEndpoint
	Session        awsclient.ConfigProvider

	EnableIAM bool
}

// NewFargateProfileScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewFargateProfileScope(params FargateProfileScopeParams) (*FargateProfileScope, error) {
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSFargateProfile")
	}
	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	session, serviceLimiters, err := sessionForRegion(params.ControlPlane.Spec.Region, params.Endpoints)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.FargateProfile, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &FargateProfileScope{
		Logger:          params.Logger,
		Client:          params.Client,
		Cluster:         params.Cluster,
		ControlPlane:    params.ControlPlane,
		FargateProfile:  params.FargateProfile,
		patchHelper:     helper,
		session:         session,
		serviceLimiters: serviceLimiters,
		controllerName:  params.ControllerName,
		enableIAM:       params.EnableIAM,
	}, nil
}

// FargateProfileScope defines the basic context for an actuator to operate upon.
type FargateProfileScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	Cluster        *clusterv1.Cluster
	ControlPlane   *controlplanev1exp.AWSManagedControlPlane
	FargateProfile *infrav1exp.AWSFargateProfile

	session         awsclient.ConfigProvider
	serviceLimiters throttle.ServiceLimiters
	controllerName  string

	enableIAM bool
}

// ManagedPoolName returns the managed machine pool name.
func (s *FargateProfileScope) ManagedPoolName() string {
	return s.FargateProfile.Name
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients
func (s *FargateProfileScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ClusterName returns the cluster name.
func (s *FargateProfileScope) ClusterName() string {
	return s.Cluster.Name
}

// EnableIAM indicates that reconciliation should create IAM roles
func (s *FargateProfileScope) EnableIAM() bool {
	return s.enableIAM
}

// AdditionalTags returns AdditionalTags from the scope's FargateProfile
// The returned value will never be nil.
func (s *FargateProfileScope) AdditionalTags() infrav1.Tags {
	if s.FargateProfile.Spec.AdditionalTags == nil {
		s.FargateProfile.Spec.AdditionalTags = infrav1.Tags{}
	}

	return s.FargateProfile.Spec.AdditionalTags.DeepCopy()
}

// RoleName returns the node group role name
func (s *FargateProfileScope) RoleName() string {
	return s.FargateProfile.Spec.RoleName
}

// ControlPlaneSubnets returns the control plane subnets.
func (s *FargateProfileScope) ControlPlaneSubnets() infrav1.Subnets {
	return s.ControlPlane.Spec.NetworkSpec.Subnets
}

// SubnetIDs returns the machine pool subnet IDs.
func (s *FargateProfileScope) SubnetIDs() []string {
	return s.FargateProfile.Spec.SubnetIDs
}

// IAMReadyFalse marks the ready condition false using warning if error isn't
// empty
func (s *FargateProfileScope) IAMReadyFalse(reason string, err string) error {
	severity := clusterv1.ConditionSeverityWarning
	if err == "" {
		severity = clusterv1.ConditionSeverityInfo
	}
	conditions.MarkFalse(
		s.FargateProfile,
		infrav1exp.IAMFargateRolesReadyCondition,
		reason,
		severity,
		err,
	)
	if err := s.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to mark role not ready")
	}
	return nil
}

// PatchObject persists the control plane configuration and status.
func (s *FargateProfileScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.FargateProfile,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			infrav1exp.EKSFargateProfileReadyCondition,
			infrav1exp.EKSFargateCreatingCondition,
			infrav1exp.EKSFargateDeletingCondition,
			infrav1exp.IAMFargateRolesReadyCondition,
		}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *FargateProfileScope) Close() error {
	return s.PatchObject()
}

// InfraCluster returns the AWS infrastructure cluster or control plane object.
func (s *FargateProfileScope) InfraCluster() cloud.ClusterObject {
	return s.ControlPlane
}

// Session returns the AWS SDK session. Used for creating clients
func (s *FargateProfileScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ControllerName returns the name of the controller that
// created the FargateProfile.
func (s *FargateProfileScope) ControllerName() string {
	return s.controllerName
}

// KubernetesClusterName is the name of the EKS cluster name.
func (s *FargateProfileScope) KubernetesClusterName() string {
	return s.ControlPlane.Spec.EKSClusterName
}
