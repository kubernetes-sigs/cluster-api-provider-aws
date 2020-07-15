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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks"
)

// AWSManagedControlPlaneReconciler reconciles a AWSManagedControlPlane object
type AWSManagedControlPlaneReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder

	scheme     *runtime.Scheme
	controller controller.Controller
}

func (r *AWSManagedControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1exp.AWSManagedControlPlane{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPaused(r.Log)).
		Watches(
			&source.Kind{Type: &infrav1exp.AWSManagedCluster{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.managedClusterToManagedControlPlane),
			},
		).
		Build(r)

	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	r.scheme = mgr.GetScheme()
	r.controller = c

	return nil
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=exp.infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=exp.infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes/status,verbs=get;update;patch

func (r *AWSManagedControlPlaneReconciler) Reconcile(req ctrl.Request) (res ctrl.Result, reterr error) {
	logger := r.Log.WithValues("namespace", req.Namespace, "eksControlPlane", req.Name)
	ctx := context.Background()

	// Get the control plane instance
	awsControlPlane := &infrav1exp.AWSManagedControlPlane{}
	if err := r.Client.Get(ctx, req.NamespacedName, awsControlPlane); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsControlPlane.ObjectMeta)
	if err != nil {
		logger.Error(err, "Failed to retrieve owner Cluster from the API Server")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		logger.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if util.IsPaused(cluster, awsControlPlane) {
		logger.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	awsManagedCluster := &infrav1exp.AWSManagedCluster{}
	awsManagedClusterRef := types.NamespacedName{
		Name:      cluster.Spec.InfrastructureRef.Name,
		Namespace: cluster.Spec.InfrastructureRef.Namespace,
	}
	if err := r.Get(ctx, awsManagedClusterRef, awsManagedCluster); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to get AWSManagedCluster")
	}

	// Wait for the aws managed cluster infrastructure to be initialized before creating the EKS control plane
	if !awsManagedCluster.Status.Initialized {
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:            r.Client,
		Logger:            logger,
		Cluster:           cluster,
		AWSManagedCluster: awsManagedCluster,
		ControlPlane:      awsControlPlane,
		ControllerName:    "awsmanagedcontrolplane",
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope
	defer func() {
		if err := managedScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !awsControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, managedScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, managedScope)
}

func (r *AWSManagedControlPlaneReconciler) reconcileNormal(ctx context.Context, managedScope *scope.ManagedControlPlaneScope) (res ctrl.Result, reterr error) {
	managedScope.Info("Reconciling AWSManagedControlPlane")

	// TODO: check failed???

	controllerutil.AddFinalizer(managedScope.ControlPlane, infrav1exp.EKSControlPlaneFinalizer)
	if err := managedScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	ekssvc := eks.NewService(managedScope)

	if err := ekssvc.ReconcileControlPlane(ctx); err != nil {
		return reconcile.Result{}, err
	}

	managedScope.ControlPlane.Status.Ready = true

	return reconcile.Result{}, nil
}

func (r *AWSManagedControlPlaneReconciler) reconcileDelete(_ context.Context, managedScope *scope.ManagedControlPlaneScope) (_ ctrl.Result, reterr error) {
	managedScope.Info("Reconciling AWSManagedClusterPlane delete")

	ekssvc := eks.NewService(managedScope)
	controlPlane := managedScope.ControlPlane

	if err := ekssvc.DeleteControlPlane(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting EKS cluster for EKS control plane %s/%s", managedScope.Namespace(), managedScope.Name())
	}

	controllerutil.RemoveFinalizer(controlPlane, infrav1exp.EKSControlPlaneFinalizer)

	return reconcile.Result{}, nil
}

// ClusterToAWSManagedControlPlane is a handler.ToRequestsFunc to be used to enqueue requests for reconciliation
// for AWSManagedControlPlane based on updates to a Cluster.
func (r *AWSManagedControlPlaneReconciler) ClusterToAWSManagedControlPlane(o handler.MapObject) []ctrl.Request {
	c, ok := o.Object.(*clusterv1.Cluster)
	if !ok {
		r.Log.Error(nil, fmt.Sprintf("Expected a Cluster but got a %T", o.Object))
		return nil
	}

	controlPlaneRef := c.Spec.ControlPlaneRef
	if controlPlaneRef != nil && controlPlaneRef.Kind == "AWSManagedControlPlane" {
		return []ctrl.Request{{NamespacedName: client.ObjectKey{Namespace: controlPlaneRef.Namespace, Name: controlPlaneRef.Name}}}
	}

	return nil
}

func (r *AWSManagedControlPlaneReconciler) managedClusterToManagedControlPlane(o handler.MapObject) []ctrl.Request {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	awsManagedCluster, ok := o.Object.(*infrav1exp.AWSManagedCluster)
	if !ok {
		r.Log.Error(nil, fmt.Sprintf("Expected a AWSManagedCluster but got a %T", o.Object))
		return nil
	}

	if !awsManagedCluster.ObjectMeta.DeletionTimestamp.IsZero() {
		r.Log.V(4).Info("AWSManagedCluster has a deletion timestamp, skipping mapping")
		return nil
	}

	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsManagedCluster.ObjectMeta)
	if err != nil {
		r.Log.Error(err, "failed to get owning cluster")
		return nil
	}

	controlPlaneRef := cluster.Spec.ControlPlaneRef
	if controlPlaneRef == nil || controlPlaneRef.Kind != "AWSManagedControlPlane" {
		r.Log.V(4).Info("ControlPlaneRef is nil or not AWSManagedControlPlane, skipping mapping")
		return nil
	}

	return []ctrl.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      controlPlaneRef.Name,
				Namespace: controlPlaneRef.Namespace,
			},
		},
	}
}
