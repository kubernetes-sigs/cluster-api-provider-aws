/*
 Copyright The Kubernetes Authors.

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

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

// RosaNetworkScopeParams defines the input parameters used to create a new RosaNetworkScope.
type RosaNetworkScopeParams struct {
	Client         client.Client
	ControllerName string
	Endpoints      []ServiceEndpoint
	Logger         *logger.Logger
	RosaNetwork    *expinfrav1.RosaNetwork
}

// RosaNetworkScope defines the basic context for an actuator to operate upon.
type RosaNetworkScope struct {
	logger.Logger
	Client            client.Client
	controllerName    string
	patchHelper       *patch.Helper
	RosaNetwork       *expinfrav1.RosaNetwork
	serviceLimiters   throttle.ServiceLimiters
	serviceLimitersV2 throttle.ServiceLimiters
	session           awsclient.ConfigProvider
	sessionV2         awsv2.Config
}

// NewRosaNetworkScope creates a new NewRosaNetworkScope from the supplied parameters.
func NewRosaNetworkScope(params RosaNetworkScopeParams) (*RosaNetworkScope, error) {
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	rosaNetworkScope := &RosaNetworkScope{
		Logger:         *params.Logger,
		Client:         params.Client,
		controllerName: params.ControllerName,
		patchHelper:    nil,
		RosaNetwork:    params.RosaNetwork,
	}

	session, serviceLimiters, err := sessionForClusterWithRegion(params.Client, rosaNetworkScope, params.RosaNetwork.Spec.Region, params.Endpoints, params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws session: %v", err)
	}

	sessionv2, serviceLimitersv2, err := sessionForClusterWithRegionV2(params.Client, rosaNetworkScope, params.RosaNetwork.Spec.Region, params.Endpoints, params.Logger)
	if err != nil {
		return nil, errors.Errorf("failed to create aws V2 session: %v", err)
	}

	patchHelper, err := patch.NewHelper(params.RosaNetwork, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	rosaNetworkScope.patchHelper = patchHelper
	rosaNetworkScope.session = session
	rosaNetworkScope.sessionV2 = *sessionv2
	rosaNetworkScope.serviceLimiters = serviceLimiters
	rosaNetworkScope.serviceLimitersV2 = serviceLimitersv2

	return rosaNetworkScope, nil
}

// SessionV2 returns the AWS SDK V2 Config. Used for creating clients.
func (s *RosaNetworkScope) SessionV2() awsv2.Config {
	return s.sessionV2
}

// IdentityRef returns the AWSIdentityReference object.
func (s *RosaNetworkScope) IdentityRef() *infrav1.AWSIdentityReference {
	return s.RosaNetwork.Spec.IdentityRef
}

// Session returns the AWS SDK session (used for creating clients).
func (s *RosaNetworkScope) Session() awsclient.ConfigProvider {
	return s.session
}

// ServiceLimiter returns the AWS SDK session (used for creating clients).
func (s *RosaNetworkScope) ServiceLimiter(service string) *throttle.ServiceLimiter {
	if sl, ok := s.serviceLimiters[service]; ok {
		return sl
	}
	return nil
}

// ControllerName returns the name of the controller.
func (s *RosaNetworkScope) ControllerName() string {
	return s.controllerName
}

// InfraCluster returns the RosaNetwork object.
// The method is then used in session.go to set proper Conditions for the RosaNetwork object.
func (s *RosaNetworkScope) InfraCluster() cloud.ClusterObject {
	return s.RosaNetwork
}

// InfraClusterName returns the name of the RosaNetwork object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *RosaNetworkScope) InfraClusterName() string {
	return s.RosaNetwork.Name
}

// Namespace returns the namespace of the RosaNetwork object.
// The method is then used in session.go to set the key to the AWS session cache.
func (s *RosaNetworkScope) Namespace() string {
	return s.RosaNetwork.Namespace
}

// PatchObject persists the rosanetwork configuration and status.
func (s *RosaNetworkScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.RosaNetwork,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.RosaNetworkReadyCondition,
		}})
}

// Close closes the current scope persisting the rosanetwork configuration and status.
func (s *RosaNetworkScope) Close() error {
	return s.PatchObject()
}
