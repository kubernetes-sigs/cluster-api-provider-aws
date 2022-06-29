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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/workload"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	cannonations "sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ExternalResourceGCReconciler is a reconciler for garbage collecting external resources. For example.
// if you create a Service of type load balancer that results in an ELB/NLB
// being created in AWS.
type ExternalResourceGCReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
}

// SetupWithManager is used to setup the controller.
func (r *ExternalResourceGCReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.Cluster{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		Complete(r)
}

// Reconcile will handle the garbage collection of external resources.
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,verbs=get;list;watch;update;patch
func (r *ExternalResourceGCReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling external resource garbage collection")

	// get the cluster
	cluster := &clusterv1.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	if cannonations.IsPaused(cluster, cluster) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	ref := cluster.Spec.InfrastructureRef
	infraCluster, err := external.Get(ctx, r.Client, ref, cluster.Namespace)
	if err != nil {
		if apierrors.IsNotFound(errors.Cause(err)) {
			log.Info("Could not find infra cluster, requeuing", "refGroupVersionKind", ref.GroupVersionKind(), "refName", ref.Name)
			return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
		}
		return reconcile.Result{}, err
	}

	extScope, err := scope.NewExternalResourceGCScope(scope.ExternalResourceGCScopeParams{
		Client:       r.Client,
		Cluster:      cluster,
		Logger:       &log,
		InfraCluster: infraCluster,
	})
	if err != nil {
		log.Error(err, "error creating external resource gc scope")
		return reconcile.Result{}, fmt.Errorf("creating external resource gc scope: %w", err)
	}

	if !cluster.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, extScope)
	}

	return r.reconcileNormal(ctx, extScope)
}

func (r *ExternalResourceGCReconciler) reconcileNormal(_ context.Context, extScope *scope.ExternalResourceGCScope) (_ ctrl.Result, reterr error) {
	extScope.Info("Reconciling external resources")

	if !extScope.InfraCluster().GetDeletionTimestamp().IsZero() {
		extScope.V(2).Info("infra cluster has been marked for deletion, not taking action to add gc finalizer")

		return reconcile.Result{}, nil
	}

	shouldGC, err := extScope.ShouldGarbageCollect()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("determining if cluster needs garbage collecting: %w", err)
	}
	if !shouldGC {
		extScope.V(2).Info("infra cluster has been marked (via annotation) as not requiring garbage collection")

		return reconcile.Result{}, nil
	}

	extScope.V(2).Info("Adding garbage collection finalizer")
	controllerutil.AddFinalizer(extScope.InfraCluster(), expinfrav1.ExternalResourceGCFinalizer)
	if err := extScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	r.Recorder.Event(
		extScope.InfraCluster(),
		"Normal",
		"MarkedForExtResourceGC",
		"the infrastructure cluster has been marked for external resource garbage collection")

	return reconcile.Result{}, nil
}

func (r *ExternalResourceGCReconciler) reconcileDelete(ctx context.Context, extScope *scope.ExternalResourceGCScope) (_ ctrl.Result, reterr error) {
	extScope.Info("Reconciling delete for external resources")

	if !controllerutil.ContainsFinalizer(extScope.InfraCluster(), expinfrav1.ExternalResourceGCFinalizer) {
		extScope.Info("infra cluster has no garbage collection finalizer, no action required")

		return reconcile.Result{}, nil
	}

	shouldGC, err := extScope.ShouldGarbageCollect()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("determining if cluster needs garbage collecting: %w", err)
	}

	if shouldGC {
		extScope.Info("Deleting from workload/tenant cluster")

		wkSvc := workload.NewService(extScope)
		res, err := wkSvc.ReconcileDelete(ctx)
		if err != nil {
			extScope.Error(err, "error deleting remote resources")
			return reconcile.Result{}, fmt.Errorf("deleting remote resources: %w", err)
		}
		if !res.IsZero() {
			return res, nil
		}
	} else {
		extScope.Info("Infra cluster has GC finalizer but not deleting from workload/tenant cluster as annotation says don't garbage collect")
	}

	extScope.V(2).Info("Removing garbage collection finalizer")
	controllerutil.RemoveFinalizer(extScope.InfraCluster(), expinfrav1.ExternalResourceGCFinalizer)
	if err := extScope.PatchObject(); err != nil {
		return reconcile.Result{}, fmt.Errorf("patching infra cluster: %w", err)
	}

	r.Recorder.Event(
		extScope.InfraCluster(),
		"Normal",
		"CompletedExtResourceGC",
		"external resource garbage collection for the infrastructure cluster has completed")

	return reconcile.Result{}, nil
}
