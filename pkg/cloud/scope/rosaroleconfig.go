/*
Copyright 2023 The Kubernetes Authors.

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

	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ROSARoleConfigScopeParams defines the input parameters used to create a new ROSARoleConfigScope.
type ROSARoleConfigScopeParams struct {
	Client         client.Client
	Logger         *logger.Logger
	Cluster        *clusterv1.Cluster
	ControlPlane   *rosacontrolplanev1.ROSAControlPlane
	ControllerName string
	Endpoints      []ServiceEndpoint
	NewStsClient   func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI
}

// NewROSARoleConfigScope creates a new ROSARoleConfigScope from the supplied parameters.
func NewROSARoleConfigScope(params ROSARoleConfigScopeParams) (*ROSARoleConfigScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil AWSManagedControlPlane")
	}
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	managedScope := &ROSARoleConfigScope{
		Logger:         *params.Logger,
		Client:         params.Client,
		Cluster:        params.Cluster,
		ControlPlane:   params.ControlPlane,
		patchHelper:    nil,
		controllerName: params.ControllerName,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, managedScope, params.ControlPlane.Spec.Region, params.Endpoints, params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	helper, err := patch.NewHelper(params.ControlPlane, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	managedScope.patchHelper = helper
	managedScope.session = session
	managedScope.serviceLimiters = serviceLimiters

	stsClient := params.NewStsClient(managedScope, managedScope, managedScope, managedScope.ControlPlane)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to identify the AWS caller: %w", err)
	}
	managedScope.Identity = identity

	return managedScope, nil
}

// ROSARoleConfigScope defines the basic context for an actuator to operate upon.
type ROSARoleConfigScope struct {
	logger.Logger
	Client      client.Client
	patchHelper *patch.Helper

	Cluster      *clusterv1.Cluster
	ControlPlane *rosacontrolplanev1.ROSAControlPlane

	session         awsclient.ConfigProvider
	serviceLimiters throttle.ServiceLimiters
	controllerName  string
	Identity        *sts.GetCallerIdentityOutput
}

// InfraCluster returns the AWSManagedControlPlane object.
func (s *ROSARoleConfigScope) InfraCluster() cloud.ClusterObject {
	return s.ControlPlane
}

// IdentityRef returns the AWSIdentityReference object.
func (s *ROSARoleConfigScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.ControlPlane.Spec.IdentityRef
}

// Session returns the AWS SDK session. Used for creating clients.
func (s *ROSARoleConfigScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ServiceLimiter returns the AWS SDK session. Used for creating clients.
func (s *ROSARoleConfigScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller.
func (s *ROSARoleConfigScope) ControllerName() string {
	return s.controllerName
}

var _ cloud.ScopeUsage = (*ROSARoleConfigScope)(nil)
var _ cloud.Session = (*ROSARoleConfigScope)(nil)
var _ cloud.SessionMetadata = (*ROSARoleConfigScope)(nil)

// Name returns the CAPI cluster name.
func (s *ROSARoleConfigScope) Name() string {
	return s.Cluster.Name
}

// InfraClusterName returns the AWS cluster name.
func (s *ROSARoleConfigScope) InfraClusterName() string {
	return s.ControlPlane.Name
}

// RosaClusterName returns the ROSA cluster name.
func (s *ROSARoleConfigScope) RosaClusterName() string {
	return s.ControlPlane.Spec.RosaClusterName
}

// Namespace returns the cluster namespace.
func (s *ROSARoleConfigScope) Namespace() string {
	return s.Cluster.Namespace
}

// CredentialsSecret returns the CredentialsSecret object.
func (s *ROSARoleConfigScope) CredentialsSecret() *corev1.Secret {
	secretRef := s.ControlPlane.Spec.CredentialsSecretRef
	if secretRef == nil {
		return nil
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.ControlPlane.Spec.CredentialsSecretRef.Name,
			Namespace: s.ControlPlane.Namespace,
		},
	}
}

// ClusterAdminPasswordSecret returns the corev1.Secret object for the cluster admin password.
func (s *ROSARoleConfigScope) ClusterAdminPasswordSecret() *corev1.Secret {
	return s.secretWithOwnerReference(fmt.Sprintf("%s-admin-password", s.Cluster.Name))
}

// ExternalAuthBootstrapKubeconfigSecret returns the corev1.Secret object for the external auth bootstrap kubeconfig.
// This is a temporarily admin kubeconfig generated using break-glass credentials for the user to bootstreap their environment like setting up RBAC for oidc users/groups.
// This Kubeonconfig will be created only once initially and be valid for only 24h.
// The kubeconfig secret will not be autoamticallty rotated and will be invalid after the 24h. However, users can opt to manually delete the secret to trigger the generation of a new one which will be valid for another 24h.
func (s *ROSARoleConfigScope) ExternalAuthBootstrapKubeconfigSecret() *corev1.Secret {
	return s.secretWithOwnerReference(fmt.Sprintf("%s-bootstrap-kubeconfig", s.Cluster.Name))
}

func (s *ROSARoleConfigScope) secretWithOwnerReference(name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.ControlPlane.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(s.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane")),
			},
		},
	}
}

// PatchObject persists the control plane configuration and status.
func (s *ROSARoleConfigScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ControlPlane,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			rosacontrolplanev1.ROSAControlPlaneReadyCondition,
			rosacontrolplanev1.ROSAControlPlaneValidCondition,
			rosacontrolplanev1.ROSAControlPlaneUpgradingCondition,
		}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *ROSARoleConfigScope) Close() error {
	return s.PatchObject()
}
