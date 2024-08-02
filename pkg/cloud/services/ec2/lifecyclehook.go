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

package ec2

import (
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ReconcileLifecycleHooks periodically reconciles a lifecycle hook for the ASG.
func (s *Service) ReconcileLifecycleHooks(scope scope.LifecycleHookScope, asgsvc services.ASGInterface) error {
	lifecyleHooks := scope.GetLifecycleHooks()
	for i := range lifecyleHooks {
		if err := s.reconcileLifecycleHook(scope, asgsvc, &lifecyleHooks[i]); err != nil {
			return err
		}
	}

	// Get a list of lifecycle hooks that are registered with the ASG but not defined in the MachinePool and delete them.
	hooks, err := asgsvc.GetLifecycleHooks(scope)
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		found := false
		for _, definedHook := range scope.GetLifecycleHooks() {
			if hook.Name == definedHook.Name {
				found = true
				break
			}
		}
		if !found {
			scope.Info("Deleting lifecycle hook", "hook", hook.Name)
			if err := asgsvc.DeleteLifecycleHook(scope, hook); err != nil {
				conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookDeletionFailedReason, clusterv1.ConditionSeverityError, err.Error())
				return err
			}
		}
	}

	return nil
}

func (s *Service) reconcileLifecycleHook(scope scope.LifecycleHookScope, asgsvc services.ASGInterface, hook *expinfrav1.AWSLifecycleHook) error {
	scope.Info("Checking for existing lifecycle hook")
	existingHook, err := asgsvc.GetLifecycleHook(scope, hook)
	if err != nil {
		conditions.MarkUnknown(scope.GetMachinePool(), expinfrav1.LifecycleHookReadyCondition, expinfrav1.LifecycleHookNotFoundReason, err.Error())
		return err
	}

	if existingHook == nil {
		scope.Info("Creating lifecycle hook")
		if err := asgsvc.CreateLifecycleHook(scope, hook); err != nil {
			conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookCreationFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
		return nil
	}

	// If the lifecycle hook exists, we need to check if it's up to date
	needsUpdate := asgsvc.LifecycleHookNeedsUpdate(scope, existingHook, hook)

	if needsUpdate {
		scope.Info("Updating lifecycle hook")
		if err := asgsvc.UpdateLifecycleHook(scope, hook); err != nil {
			conditions.MarkFalse(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition, expinfrav1.LifecycleHookUpdateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
	}

	conditions.MarkTrue(scope.GetMachinePool(), expinfrav1.LifecycleHookExistsCondition)
	return nil
}
