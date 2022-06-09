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

package workload

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileDelete performs any operations that relate to the reconciliation of a cluster delete as it relates to
// the remote/child/tenant/workload cluster. For example, it will delete Services of type load balancer to ensure
// that an ELB/NLB in AWS have been deleted.
func (s *Service) ReconcileDelete(ctx context.Context) (reconcile.Result, error) {
	s.scope.Info("reconciling deletion of external resources")

	requeue, err := s.deleteServices(ctx, corev1.ServiceTypeLoadBalancer)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("deleting workload services of type load balancer: %w", err)
	}

	if requeue {
		s.scope.V(2).Info("requeue after deleting workload services")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// TODO: add additional deletions in the future, like PVC

	return reconcile.Result{}, err
}

func (s *Service) deleteServices(ctx context.Context, serviceType corev1.ServiceType) (requeue bool, err error) {
	s.scope.Info("deleting workload services", "type", serviceType)

	services := &corev1.ServiceList{}
	if listErr := s.scope.Client.List(ctx, services); listErr != nil {
		return false, fmt.Errorf("listing services in remote cluster: %w", err)
	}

	items := []*corev1.Service{}
	for i := range services.Items {
		svc := services.Items[i]
		if svc.Spec.Type == serviceType {
			items = append(items, &svc)
		}
	}

	if len(items) == 0 {
		s.scope.V(2).Info("no services found to delete", "type", serviceType)
		return false, nil
	}

	for i := range items {
		svc := services.Items[i]
		if !svc.DeletionTimestamp.IsZero() {
			s.scope.V(2).Info("load balancer service is deleting already", "name", svc.Name, "namespace", svc.Namespace)
			continue
		}

		s.scope.V(2).Info("deleting service", "name", svc.Name, "namespace", svc.Namespace, "type", serviceType)
		if deleteErr := s.scope.Client.Delete(ctx, &svc); deleteErr != nil {
			return false, fmt.Errorf("deleting load balancer service %s/%s: %w", svc.Namespace, svc.Name, deleteErr)
		}
	}
	s.scope.V(2).Info("requeueing after deleting services", "type", serviceType)

	return true, nil
}
