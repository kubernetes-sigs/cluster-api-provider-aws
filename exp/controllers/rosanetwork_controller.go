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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// RosaNetworkReconciler reconciles a RosaNetwork object
type RosaNetworkReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks/status,verbs=get;update;patch

func (r *RosaNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("rosanetwork", req.NamespacedName)

	// Get the rosanetwork instance
	rosaNetwork := &expinfrav1.RosaNetwork{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaNetwork); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// TODO
	// 1. instantiate rosaNetScope
	// 2. defer scope close
	// 3. if rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() { return r.reconcileDelete(ctx, rosaNetScope) }
	// 4. otherwise r.reconcileNormal(ctx, rosaScope(ctx, rosaNetScope)

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) reconcileNormal(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	// TODO

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) reconcileDelete(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	// TODO

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.RosaNetwork{}).
		Complete(r)
}
