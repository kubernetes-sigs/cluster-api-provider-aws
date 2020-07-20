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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/network"
)

// AWSManagedControlPlaneReconciler reconciles a AWSManagedControlPlane object
type AWSManagedControlPlaneReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
}

// SetupWithManager is used to setup the controller
func (r *AWSManagedControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	awsManagedControlPlane := &infrav1exp.AWSManagedControlPlane{}
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(awsManagedControlPlane).
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
		return fmt.Errorf("failed setting up the AWSManagedControlPlane controller manager: %w", err)
	}

	if err = c.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: util.ClusterToInfrastructureMapFunc(awsManagedControlPlane.GroupVersionKind()),
		},
		predicates.ClusterUnpausedAndInfrastructureReady(r.Log),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	return nil
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedcontrolplanes/status,verbs=get;update;patch

// Reconcile will reconcile AWSManagedControlPlane Resources
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

	logger = logger.WithValues("cluster", cluster.Name)

	managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:         r.Client,
		Logger:         logger,
		Cluster:        cluster,
		ControlPlane:   awsControlPlane,
		ControllerName: "awsmanagedcontrolplane",
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope
	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			infrav1exp.EKSControlPlaneReadyCondition,
			infrav1exp.IAMControlPlaneRolesReadyCondition,
			infrav1.VpcReadyCondition,
			infrav1.SubnetsReadyCondition,
			infrav1.ClusterSecurityGroupsReadyCondition,
		}

		if managedScope.VPC().IsManaged(managedScope.Name()) {
			applicableConditions = append(applicableConditions,
				infrav1.InternetGatewayReadyCondition,
				infrav1.NatGatewaysReadyCondition,
				infrav1.RouteTablesReadyCondition,
			)
			if managedScope.Bastion().Enabled {
				applicableConditions = append(applicableConditions, infrav1.BastionHostReadyCondition)
			}
		}

		conditions.SetSummary(managedScope.ControlPlane, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

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

	awsManagedControlPlane := managedScope.ControlPlane

	controllerutil.AddFinalizer(managedScope.ControlPlane, infrav1exp.ManagedControlPlaneFinalizer)
	if err := managedScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	ec2Service := ec2.NewService(managedScope)
	networkSvc := network.NewService(managedScope)
	ekssvc := eks.NewService(managedScope)

	if err := networkSvc.ReconcileNetwork(); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile network for AWSManagedControlPlane %s/%s: %w", awsManagedControlPlane.Namespace, awsManagedControlPlane.Name, err)
	}

	if err := ec2Service.ReconcileBastion(); err != nil {
		conditions.MarkFalse(awsManagedControlPlane, infrav1.BastionHostReadyCondition, infrav1.BastionHostFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return reconcile.Result{}, fmt.Errorf("failed to reconcile bastion host for AWSManagedControlPlane %s/%s: %w", awsManagedControlPlane.Namespace, awsManagedControlPlane.Name, err)
	}

	if err := ekssvc.ReconcileControlPlane(ctx); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile control plane for AWSManagedControlPlane %s/%s: %s", awsManagedControlPlane.Namespace, awsManagedControlPlane.Name, err)
	}

	for _, subnet := range managedScope.Subnets().FilterPrivate() {
		managedScope.SetFailureDomain(subnet.AvailabilityZone, clusterv1.FailureDomainSpec{
			ControlPlane: true,
		})
	}

	return reconcile.Result{}, nil
}

func (r *AWSManagedControlPlaneReconciler) reconcileDelete(_ context.Context, managedScope *scope.ManagedControlPlaneScope) (_ ctrl.Result, reterr error) {
	managedScope.Info("Reconciling AWSManagedClusterPlane delete")

	controlPlane := managedScope.ControlPlane

	ekssvc := eks.NewService(managedScope)
	ec2svc := ec2.NewService(managedScope)
	networkSvc := network.NewService(managedScope)

	if err := ekssvc.DeleteControlPlane(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting EKS cluster for EKS control plane %s/%s", managedScope.Namespace(), managedScope.Name())
	}

	if err := ec2svc.DeleteBastion(); err != nil {
		return reconcile.Result{}, fmt.Errorf("error deleting bastion for AWSManagedControlPlane %s/%s: %w", controlPlane.Namespace, controlPlane.Name, err)
	}

	if err := networkSvc.DeleteNetwork(); err != nil {
		return reconcile.Result{}, fmt.Errorf("error deleting network for AWSManagedControlPlane %s/%s: %w", controlPlane.Namespace, controlPlane.Name, err)
	}

	controllerutil.RemoveFinalizer(controlPlane, infrav1exp.ManagedControlPlaneFinalizer)

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
	if cluster == nil {
		r.Log.V(4).Info("Owning cluster not set on AWSManagedCluster, skipping mapping")
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
