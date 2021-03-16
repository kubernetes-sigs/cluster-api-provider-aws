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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
)

// AWSManagedClusterReconciler reconciles AWSManagedCluster
type AWSManagedClusterReconciler struct {
	client.Client
	Log              logr.Logger
	Recorder         record.EventRecorder
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes;awsmanagedcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *AWSManagedClusterReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	log := r.Log.WithValues("namespace", req.Namespace, "awsManagedCluster", req.Name)

	// Fetch the AWSManagedCluster instance
	awsManagedCluster := &infrav1exp.AWSManagedCluster{}
	err := r.Get(ctx, req.NamespacedName, awsManagedCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsManagedCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	if util.IsPaused(cluster, awsManagedCluster) {
		log.Info("AWSManagedCluster or linked Cluster is marked as paused. Won't reconcile")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	controlPlaneRef := types.NamespacedName{
		Name:      cluster.Spec.ControlPlaneRef.Name,
		Namespace: cluster.Spec.ControlPlaneRef.Namespace,
	}

	if err := r.Get(ctx, controlPlaneRef, controlPlane); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get control plane ref: %w", err)
	}

	log = log.WithValues("controlPlane", controlPlaneRef.Name)

	patchHelper, err := patch.NewHelper(awsManagedCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Set the values from the managed control plane
	awsManagedCluster.Status.Ready = controlPlane.Status.Ready
	awsManagedCluster.Spec.ControlPlaneEndpoint = controlPlane.Spec.ControlPlaneEndpoint
	awsManagedCluster.Status.FailureDomains = controlPlane.Status.FailureDomains

	if err := patchHelper.Patch(ctx, awsManagedCluster); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to patch AWSManagedCluster: %w", err)
	}

	log.Info("Successfully reconciled AWSManagedCluster")

	return reconcile.Result{}, nil

}

func (r *AWSManagedClusterReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	awsManagedCluster := &infrav1exp.AWSManagedCluster{}

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(awsManagedCluster).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(r.Log, r.WatchFilterValue)).
		Build(r)

	if err != nil {
		return fmt.Errorf("error creating controller: %w", err)
	}

	// Add a watch for clusterv1.Cluster unpaise
	if err = controller.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: util.ClusterToInfrastructureMapFunc(awsManagedCluster.GroupVersionKind()),
		},
		predicates.ResourceNotPausedAndHasFilterLabel(r.Log, r.WatchFilterValue),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	// Add a watch for AWSManagedControlPlane
	if err = controller.Watch(
		&source.Kind{Type: &ekscontrolplanev1.AWSManagedControlPlane{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(r.managedControlPlaneToManagedCluster),
		},
	); err != nil {
		return fmt.Errorf("failed adding watch on AWSManagedControlPlane: %w", err)
	}

	return nil
}

func (r *AWSManagedClusterReconciler) managedControlPlaneToManagedCluster(o handler.MapObject) []ctrl.Request {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	awsManagedControlPlane, ok := o.Object.(*ekscontrolplanev1.AWSManagedControlPlane)
	if !ok {
		r.Log.Error(nil, fmt.Sprintf("Expected a AWSManagedControlPlane but got a %T", o.Object))
		return nil
	}

	if !awsManagedControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		r.Log.V(4).Info("AWSManagedControlPlane has a deletion timestamp, skipping mapping")
		return nil
	}

	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsManagedControlPlane.ObjectMeta)
	if err != nil {
		r.Log.Error(err, "failed to get owning cluster")
		return nil
	}
	if cluster == nil {
		r.Log.Info("no owning cluster, skipping mapping")
		return nil
	}

	managedClusterRef := cluster.Spec.InfrastructureRef
	if managedClusterRef == nil || managedClusterRef.Kind != "AWSManagedCluster" {
		r.Log.V(4).Info("InfrastructureRef is nil or not AWSManagedCluster, skipping mapping")
		return nil
	}

	return []ctrl.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      managedClusterRef.Name,
				Namespace: managedClusterRef.Namespace,
			},
		},
	}
}
