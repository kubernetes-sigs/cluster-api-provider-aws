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
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/network"
)

// AWSManagedClusterReconciler reconciles AWSManagedCluster
type AWSManagedClusterReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=exp.infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=exp.infrastructure.cluster.x-k8s.io,resources=awsmanagedclusters/status,verbs=get;update;patch
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

	controlPlane := &infrav1exp.AWSManagedControlPlane{}
	controlPlaneRef := types.NamespacedName{
		Name:      cluster.Spec.ControlPlaneRef.Name,
		Namespace: cluster.Spec.ControlPlaneRef.Namespace,
	}

	if err := r.Get(ctx, controlPlaneRef, controlPlane); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get control plane ref: %w", err)
	}

	managedClusterScope, err := scope.NewManagedClusterScope(scope.ManagedClusterScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		AWSManagedCluster: awsManagedCluster,
		Controlplane:      controlPlane,
		ControllerName:    "awsmanagedcluster",
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create managed scope: %w", err)
	}

	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			infrav1.VpcReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.ClusterSecurityGroupsReadyCondition,
		}

		if managedClusterScope.VPC().IsManaged(managedClusterScope.Name()) {
			applicableConditions = append(applicableConditions,
				infrav1.InternetGatewayReadyCondition,
				infrav1.NatGatewaysReadyCondition,
				infrav1.RouteTablesReadyCondition,
			)
			if managedClusterScope.Bastion().Enabled {
				applicableConditions = append(applicableConditions, infrav1.BastionHostReadyCondition)
			}
		}

		conditions.SetSummary(managedClusterScope.AWSManagedCluster, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := managedClusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Hnadle deleted clusters
	if !awsManagedCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(managedClusterScope)
	}

	return r.reconcileNormal(managedClusterScope)
}

func (r *AWSManagedClusterReconciler) reconcileDelete(managedClusterScope *scope.ManagedClusterScope) (reconcile.Result, error) {
	managedClusterScope.Info("Reconciling AWSManagedCluster delete")

	ec2svc := ec2.NewService(managedClusterScope)
	networkSvc := network.NewService(managedClusterScope)

	awsManagedCluster := managedClusterScope.InfraCluster().(*infrav1exp.AWSManagedCluster)

	//TODO: wait for the control plane to delete

	if err := ec2svc.DeleteBastion(); err != nil {
		return reconcile.Result{}, fmt.Errorf("error deleting bastion for AWSManagedCluster %s/%s: %w", awsManagedCluster.Namespace, awsManagedCluster.Name, err)
	}

	if err := networkSvc.DeleteNetwork(); err != nil {
		return reconcile.Result{}, fmt.Errorf("error deleting network for AWSManagedCluster %s/%s: %w", awsManagedCluster.Namespace, awsManagedCluster.Name, err)
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(awsManagedCluster, infrav1exp.ManagedClusterFinalizer)

	return reconcile.Result{}, nil
}

func (r *AWSManagedClusterReconciler) reconcileNormal(managedClusterScope *scope.ManagedClusterScope) (reconcile.Result, error) {
	managedClusterScope.Info("Reconciling AWSManagedCluster")

	awsManagedCluster := managedClusterScope.AWSManagedCluster

	// If the AWSManagedCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(awsManagedCluster, infrav1exp.ManagedClusterFinalizer)
	// Register the finalizer immediately to avoid orphaning AWS resources on delete
	if err := managedClusterScope.PatchObject(); err != nil {
		return reconcile.Result{}, err
	}

	ec2Service := ec2.NewService(managedClusterScope)
	networkSvc := network.NewService(managedClusterScope)

	if err := networkSvc.ReconcileNetwork(); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile network for AWSManagedCluster %s/%s: %w", awsManagedCluster.Namespace, awsManagedCluster.Name, err)
	}

	if err := ec2Service.ReconcileBastion(); err != nil {
		conditions.MarkFalse(awsManagedCluster, infrav1.BastionHostReadyCondition, infrav1.BastionHostFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return reconcile.Result{}, fmt.Errorf("failed to reconcile bastion host for AWSManagedCluster %s/%s: %w", awsManagedCluster.Namespace, awsManagedCluster.Name, err)
	}

	for _, subnet := range managedClusterScope.Subnets().FilterPrivate() {
		managedClusterScope.SetFailureDomain(subnet.AvailabilityZone, clusterv1.FailureDomainSpec{
			ControlPlane: true,
		})
	}

	// We have initialized the infra - its ok
	// for the control plane to reconcile
	awsManagedCluster.Status.Initialized = true

	// Check the control plane and see if we are ready
	controlPlane := managedClusterScope.Controlplane
	if controlPlane.Status.Ready {
		awsManagedCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
			Host: controlPlane.Spec.ControlPlaneEndpoint.Host,
			Port: controlPlane.Spec.ControlPlaneEndpoint.Port,
		}
		if !awsManagedCluster.Spec.ControlPlaneEndpoint.IsZero() {
			awsManagedCluster.Status.Ready = true
		}
	}

	return reconcile.Result{}, nil
}

func (r *AWSManagedClusterReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	awsManagedCluster := &infrav1exp.AWSManagedCluster{}

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(awsManagedCluster).
		WithEventFilter(predicates.ResourceNotPaused(r.Log)).
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
		predicates.ClusterUnpaused(r.Log),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	// Add a watch for AWSManagedControlPlane
	if err = controller.Watch(
		&source.Kind{Type: &infrav1exp.AWSManagedControlPlane{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(r.managedControlPlaneToManagedCluster),
		},
		predicates.ClusterUnpaused(r.Log),
	); err != nil {
		return fmt.Errorf("failed adding watch on AWSManagedControlPlane: %w", err)
	}

	return nil
}

func (r *AWSManagedClusterReconciler) managedControlPlaneToManagedCluster(o handler.MapObject) []ctrl.Request {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	awsManagedControlPlane, ok := o.Object.(*infrav1exp.AWSManagedControlPlane)
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
