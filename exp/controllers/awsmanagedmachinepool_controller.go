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

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capiv1exp "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks"
)

// AWSManagedMachinePoolReconciler reconciles a AWSManagedMachinePool object
type AWSManagedMachinePoolReconciler struct {
	client.Client
	Log       logr.Logger
	Recorder  record.EventRecorder
	Endpoints []scope.ServiceEndpoint

	EnableIAM        bool
	WatchFilterValue string
}

// SetupWithManager is used to setup the controller
func (r *AWSManagedMachinePoolReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	gvk, err := apiutil.GVKForObject(new(infrav1exp.AWSManagedMachinePool), mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "failed to find GVK for AWSManagedMachinePool")
	}
	managedControlPlaneToManagedMachinePoolMap := managedControlPlaneToManagedMachinePoolMapFunc(r.Client, gvk, r.Log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1exp.AWSManagedMachinePool{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(r.Log, r.WatchFilterValue)).
		Watches(
			&source.Kind{Type: &capiv1exp.MachinePool{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: machinePoolToInfrastructureMapFunc(gvk),
			},
		).
		Watches(
			&source.Kind{Type: &controlplanev1.AWSManagedControlPlane{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: managedControlPlaneToManagedMachinePoolMap,
			},
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=exp.cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes;awsmanagedcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools/status,verbs=get;update;patch

// Reconcile reconciles AWSManagedMachinePools
func (r *AWSManagedMachinePoolReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := r.Log.WithValues("namespace", req.Namespace, "AWSManagedMachinePool", req.Name)
	ctx := context.Background()

	awsPool := &infrav1exp.AWSManagedMachinePool{}
	if err := r.Get(ctx, req.NamespacedName, awsPool); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	machinePool, err := getOwnerMachinePool(ctx, r.Client, awsPool.ObjectMeta)
	if err != nil {
		logger.Error(err, "Failed to retrieve owner MachinePool from the API Server")
		return ctrl.Result{}, err
	}
	if machinePool == nil {
		logger.Info("MachinePool Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("MachinePool", machinePool.Name)

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		logger.Info("Failed to retrieve Cluster from MachinePool")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("Cluster", cluster.Name)

	controlPlaneKey := client.ObjectKey{
		Namespace: awsPool.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &controlplanev1.AWSManagedControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		logger.Info("Failed to retrieve ControlPlane from MachinePool")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("AWSManagedControlPlane", controlPlane.Name)

	if !controlPlane.Status.Ready {
		logger.Info("Control plane is not ready yet")
		conditions.MarkFalse(awsPool, infrav1exp.EKSNodegroupReadyCondition, infrav1exp.WaitingForEKSControlPlaneReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	machinePoolScope, err := scope.NewManagedMachinePoolScope(scope.ManagedMachinePoolScopeParams{
		Logger:             logger,
		Client:             r.Client,
		ControllerName:     "awsmanagedmachinepool",
		Cluster:            cluster,
		ControlPlane:       controlPlane,
		MachinePool:        machinePool,
		ManagedMachinePool: awsPool,
		EnableIAM:          r.EnableIAM,
		Endpoints:          r.Endpoints,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			infrav1exp.EKSNodegroupReadyCondition,
			infrav1exp.IAMNodegroupRolesReadyCondition,
		}

		conditions.SetSummary(machinePoolScope.ManagedMachinePool, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := machinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !awsPool.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machinePoolScope)
	}

	return r.reconcileNormal(ctx, machinePoolScope)
}

func (r *AWSManagedMachinePoolReconciler) reconcileNormal(
	_ context.Context,
	machinePoolScope *scope.ManagedMachinePoolScope,
) (ctrl.Result, error) {
	machinePoolScope.Info("Reconciling AWSManagedMachinePool")

	controllerutil.AddFinalizer(machinePoolScope.ManagedMachinePool, infrav1exp.ManagedMachinePoolFinalizer)
	if err := machinePoolScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	ekssvc := eks.NewNodegroupService(machinePoolScope)

	if err := ekssvc.ReconcilePool(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile machine pool for AWSManagedMachinePool %s/%s", machinePoolScope.ManagedMachinePool.Namespace, machinePoolScope.ManagedMachinePool.Name)
	}

	return ctrl.Result{}, nil
}

func (r *AWSManagedMachinePoolReconciler) reconcileDelete(
	_ context.Context,
	machinePoolScope *scope.ManagedMachinePoolScope,
) (ctrl.Result, error) {
	machinePoolScope.Info("Reconciling deletion of AWSManagedMachinePool")

	ekssvc := eks.NewNodegroupService(machinePoolScope)

	if err := ekssvc.ReconcilePoolDelete(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile machine pool deletion for AWSManagedMachinePool %s/%s", machinePoolScope.ManagedMachinePool.Namespace, machinePoolScope.ManagedMachinePool.Name)
	}

	controllerutil.RemoveFinalizer(machinePoolScope.ManagedMachinePool, infrav1exp.ManagedMachinePoolFinalizer)

	return reconcile.Result{}, nil
}

// GetOwnerClusterKey returns only the Cluster name and namespace
func GetOwnerClusterKey(obj metav1.ObjectMeta) (*client.ObjectKey, error) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "Cluster" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == clusterv1.GroupVersion.Group {
			return &client.ObjectKey{
				Namespace: obj.Namespace,
				Name:      ref.Name,
			}, nil
		}
	}
	return nil, nil
}

func managedControlPlaneToManagedMachinePoolMapFunc(c client.Client, gvk schema.GroupVersionKind, log logr.Logger) handler.ToRequestsFunc {
	return func(o handler.MapObject) []reconcile.Request {
		ctx := context.Background()
		awsControlPlane, ok := o.Object.(*controlplanev1.AWSManagedControlPlane)
		if !ok {
			return nil
		}
		if !awsControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			return nil
		}

		clusterKey, err := GetOwnerClusterKey(awsControlPlane.ObjectMeta)
		if err != nil {
			log.Error(err, "couldn't get AWS control plane owner ObjectKey")
			return nil
		}
		if clusterKey == nil {
			return nil
		}

		managedPoolForClusterList := capiv1exp.MachinePoolList{}
		if err := c.List(
			ctx, &managedPoolForClusterList, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterLabelName: clusterKey.Name},
		); err != nil {
			log.Error(err, "couldn't list pools for cluster")
			return nil
		}

		mapFunc := machinePoolToInfrastructureMapFunc(gvk)

		var results []ctrl.Request
		for i := range managedPoolForClusterList.Items {
			managedPool := mapFunc.Map(handler.MapObject{
				Object: &managedPoolForClusterList.Items[i],
			})
			results = append(results, managedPool...)
		}

		return results
	}
}
