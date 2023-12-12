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

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rosav1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ROSAMachinePoolScopeParams defines the input parameters used to create a new Scope.
type ROSAMachinePoolScopeParams struct {
	Client          client.Client
	Logger          *logger.Logger
	Cluster         *clusterv1.Cluster
	ControlPlane    *rosav1.ROSAControlPlane
	ROSAMachinePool *expinfrav1.ROSAMachinePool
	MachinePool     *expclusterv1.MachinePool
	ControllerName  string
}

// NewROSAMachinePoolScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewROSAMachinePoolScope(params ROSAMachinePoolScopeParams) (*ROSAMachinePoolScope, error) {
	if params.ControlPlane == nil {
		return nil, errors.New("failed to generate new scope from nil ControlPlane")
	}
	if params.MachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil MachinePool")
	}
	if params.ROSAMachinePool == nil {
		return nil, errors.New("failed to generate new scope from nil ROSAMachinePool")
	}
	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	patchHelper, err := patch.NewHelper(params.ROSAMachinePool, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init ROSAMachinePool patch helper")
	}

	return &ROSAMachinePoolScope{
		Logger:          *params.Logger,
		Client:          params.Client,
		Cluster:         params.Cluster,
		ControlPlane:    params.ControlPlane,
		ROSAMachinePool: params.ROSAMachinePool,
		MachinePool:     params.MachinePool,
		controllerName:  params.ControllerName,
		patchHelper:     patchHelper,
	}, nil
}

// ROSAMachinePoolScope defines the basic context for an actuator to operate upon.
type ROSAMachinePoolScope struct {
	logger.Logger
	client.Client
	Cluster         *clusterv1.Cluster
	ControlPlane    *rosav1.ROSAControlPlane
	ROSAMachinePool *expinfrav1.ROSAMachinePool
	MachinePool     *expclusterv1.MachinePool
	controllerName  string
	patchHelper     *patch.Helper
}

// PatchObject persists the control plane configuration and status.
func (s *ROSAMachinePoolScope) PatchObject() error {
	return s.patchHelper.Patch(
		context.TODO(),
		s.ROSAMachinePool,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			expinfrav1.EKSNodegroupReadyCondition,
			expinfrav1.IAMNodegroupRolesReadyCondition,
		}})
}

// Close closes the current scope persisting the control plane configuration and status.
func (s *ROSAMachinePoolScope) Close() error {
	return s.PatchObject()
}
