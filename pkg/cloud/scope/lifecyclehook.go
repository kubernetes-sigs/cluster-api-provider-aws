/*
Copyright 2022 The Kubernetes Authors.

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
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

// LaunchTemplateScope defines a scope defined around a launch template.
type LifecycleHookScope struct {
	client.Client
	logger.Logger

	MachinePool    *expclusterv1.MachinePool
	LifecycleHooks []expinfrav1.AWSLifecycleHook
	AWSMachinePool *expinfrav1.AWSMachinePool
}

type LifecycleHookScopeParams struct {
	Client client.Client
	Logger *logger.Logger

	MachinePool    *expclusterv1.MachinePool
	LifecycleHooks []expinfrav1.AWSLifecycleHook
	AWSMachinePool *expinfrav1.AWSMachinePool
}

// NewLifecycleHookScope creates a new LifecycleHookScope.
func NewLifecycleHookScope(params LifecycleHookScopeParams) (*LifecycleHookScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a LifecycleHookScope")
	}
	if params.MachinePool == nil {
		return nil, errors.New("machine pool is required when creating a LifecycleHookScope")
	}
	if params.LifecycleHooks == nil {
		params.LifecycleHooks = make([]expinfrav1.AWSLifecycleHook, 0)
	}
	if params.AWSMachinePool == nil {
		return nil, errors.New("aws machine pool is required when creating a LifecycleHookScope")
	}

	if params.Logger == nil {
		log := klog.Background()
		params.Logger = logger.NewLogger(log)
	}

	return &LifecycleHookScope{
		Client:         params.Client,
		Logger:         *params.Logger,
		MachinePool:    params.MachinePool,
		LifecycleHooks: params.LifecycleHooks,
		AWSMachinePool: params.AWSMachinePool,
	}, nil
}

func (s *LifecycleHookScope) GetASGName() string {
	return s.AWSMachinePool.Name
}

func (s *LifecycleHookScope) GetLifecycleHooks() []expinfrav1.AWSLifecycleHook {
	return s.LifecycleHooks
}

func (s *LifecycleHookScope) GetMachinePool() *expclusterv1.MachinePool {
	return s.MachinePool
}
