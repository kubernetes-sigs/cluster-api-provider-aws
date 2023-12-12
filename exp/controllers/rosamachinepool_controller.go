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
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAMachinePoolReconciler reconciles a ROSAMachinePool object.
type ROSAMachinePoolReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
}

// SetupWithManager is used to setup the controller.
func (r *ROSAMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	gvk, err := apiutil.GVKForObject(new(expinfrav1.ROSAMachinePool), mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "failed to find GVK for ROSAMachinePool")
	}
	// TODO (alberto): implement rosaControlPlaneToROSAMachinePoolMapFunc.
	// rosaControlPlaneToROSAMachinePoolMap := rosaControlPlaneToROSAMachinePoolMapFunc(r.Client, gvk, log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSAMachinePool{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log.GetLogger(), r.WatchFilterValue)).
		Watches(
			&expclusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(machinePoolToInfrastructureMapFunc(gvk)),
		).
		// Watches(
		//	&rosacontrolplanev1.ROSAControlPlane{},
		//	handler.EnqueueRequestsFromMapFunc(rosaControlPlaneToROSAMachinePoolMap),
		// ).
		Complete(r)
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes;rosadcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools/status,verbs=get;update;patch

// Reconcile reconciles ROSAMachinePools.
func (r *ROSAMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	rosaMachinePool := &expinfrav1.ROSAMachinePool{}
	if err := r.Get(ctx, req.NamespacedName, rosaMachinePool); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	machinePool, err := getOwnerMachinePool(ctx, r.Client, rosaMachinePool.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner MachinePool from the API Server")
		return ctrl.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("MachinePool", klog.KObj(machinePool))

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("Failed to retrieve Cluster from MachinePool")
		return reconcile.Result{}, nil
	}

	if annotations.IsPaused(cluster, rosaMachinePool) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	controlPlaneKey := client.ObjectKey{
		Namespace: rosaMachinePool.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		log.Info("Failed to retrieve ControlPlane from MachinePool")
		return reconcile.Result{}, nil
	}

	rosaMachinePoolScope, err := scope.NewROSAMachinePoolScope(scope.ROSAMachinePoolScopeParams{
		Client:          r.Client,
		ControllerName:  "rosamachinepool",
		Cluster:         cluster,
		ControlPlane:    controlPlane,
		MachinePool:     machinePool,
		ROSAMachinePool: rosaMachinePool,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	if !controlPlane.Status.Ready {
		log.Info("Control plane is not ready yet")
		// TODO(alberto): markFalse ROSAMachinePoolReadyCondition here.
		// conditions.MarkFalse(rosaMachinePool, expinfrav1.ROSAMachinePoolReadyCondition, "WaitingForControlPane", clusterv1.ConditionSeverityInfo, "")
		// if err := rosaMachinePoolScope.Close(); err != nil && reterr == nil {
		// 	 reterr = err
		// }
		return ctrl.Result{}, nil
	}

	defer func() {
		//  applicableConditions := []clusterv1.ConditionType{
		//	// TODO(alberto): add conditions that impact readiness.
		//  }
		//  conditions.SetSummary(rosaMachinePoolScope.ROSAMachinePool, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := rosaMachinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !rosaMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, rosaMachinePoolScope)
	}

	return ctrl.Result{}, r.reconcileNormal(ctx, rosaMachinePoolScope)
}

func (r *ROSAMachinePoolReconciler) reconcileNormal(ctx context.Context, rosaMachinePoolScope *scope.ROSAMachinePoolScope,
) error {
	rosaMachinePoolScope.Info("Reconciling ROSAdMachinePool")
	log := logger.FromContext(ctx)

	if controllerutil.AddFinalizer(rosaMachinePoolScope.ROSAMachinePool, expinfrav1.ManagedMachinePoolFinalizer) {
		if err := rosaMachinePoolScope.PatchObject(); err != nil {
			return err
		}
	}

	// create OCM NodePool
	ocmClient, err := newOCMClient()
	if err != nil {
		return err
	}
	defer func() {
		ocmClient.ocm.Close()
	}()

	rosaMachinePool := rosaMachinePoolScope.ROSAMachinePool
	machinePool := rosaMachinePoolScope.MachinePool
	controlPlane := rosaMachinePoolScope.ControlPlane

	nodePool, found, err := ocmClient.GetNodePool(*controlPlane.Status.ID, rosaMachinePool.Name[:15])
	if err != nil {
		return err
	}
	if found {
		conditions.MarkFalse(rosaMachinePool,
			expinfrav1.ROSAMachinePoolReadyCondition,
			nodePool.Status().Message(),
			clusterv1.ConditionSeverityInfo,
			"")

		rosaMachinePool.Status.Replicas = int32(nodePool.Status().CurrentReplicas())
		if nodePool.Replicas() == nodePool.Status().CurrentReplicas() && nodePool.Status().Message() == "" {
			conditions.MarkTrue(rosaMachinePool, expinfrav1.ROSAMachinePoolReadyCondition)
		}
		if err := rosaMachinePoolScope.PatchObject(); err != nil {
			return err
		}

		// TODO (alberto): discover and store providerIDs from aws so the CAPI controller can match then to Nodes and report readiness.
		log.Info("NodePool exists", "state", nodePool.Status().Message())
		return nil
	}

	npBuilder := cmv1.NewNodePool()
	npBuilder.ID(rosaMachinePool.Name[:15])
	// .Labels(labelMap).Taints(taintBuilders...)

	if rosaMachinePool.Spec.Autoscaling != nil {
		npBuilder = npBuilder.Autoscaling(
			cmv1.NewNodePoolAutoscaling().
				MinReplica(rosaMachinePool.Spec.Autoscaling.MinReplica).
				MaxReplica(rosaMachinePool.Spec.Autoscaling.MaxReplica))
	} else {
		npBuilder = npBuilder.Replicas(int(*machinePool.Spec.Replicas))
	}

	if rosaMachinePool.Spec.AWS.Subnet != "" {
		npBuilder.Subnet(rosaMachinePool.Spec.AWS.Subnet)
	}

	npBuilder.AWSNodePool(cmv1.NewAWSNodePool().InstanceType(rosaMachinePool.Spec.AWS.InstanceType))

	newNodePool, err := npBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to build rosa nodepool: %w", err)
	}

	_, err = ocmClient.CreateNodePool(*controlPlane.Status.ID, newNodePool)
	if err != nil {
		return fmt.Errorf("failed to build reconcile ROSAMachinePool: %w", err)
	}

	return nil
}

func (r *ROSAMachinePoolReconciler) reconcileDelete(_ context.Context, rosaMachinePoolScope *scope.ROSAMachinePoolScope,
) error {
	// create OCM NodePool
	ocmClient, err := newOCMClient()
	if err != nil {
		return err
	}
	defer ocmClient.Close()

	if err := ocmClient.DeleteNodePool(*rosaMachinePoolScope.ControlPlane.Status.ID, rosaMachinePoolScope.ROSAMachinePool.Name[:15]); err != nil {
		return err
	}
	controllerutil.RemoveFinalizer(rosaMachinePoolScope.ROSAMachinePool, expinfrav1.ManagedMachinePoolFinalizer)
	if err := rosaMachinePoolScope.PatchObject(); err != nil {
		return err
	}
	return nil
}

// OCMClient is a temporary helper to talk to OCM API.
// TODO(alberto): vendor this from https://github.com/openshift/rosa/tree/master/pkg/ocm or build its own package here.
type OCMClient struct {
	ocm *sdk.Connection
}

func newOCMClient() (*OCMClient, error) {
	// Create the connection, and remember to close it:
	token := os.Getenv("OCM_TOKEN")
	ocmAPIUrl := os.Getenv("OCM_API_URL")
	if ocmAPIUrl == "" {
		ocmAPIUrl = "https://api.openshift.com"
	}

	// Create a logger that has the debug level enabled:
	ocmLogger, err := sdk.NewGoLoggerBuilder().
		Debug(false).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	connection, err := sdk.NewConnectionBuilder().
		Logger(ocmLogger).
		Tokens(token).
		URL(ocmAPIUrl).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to ocm client: %w", err)
	}
	ocmClient := OCMClient{ocm: connection}

	return &ocmClient, nil
}

func (client *OCMClient) Close() error {
	return client.ocm.Close()
}

func (client *OCMClient) CreateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := client.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		Add().Body(nodePool).
		Send()
	if err != nil {
		return nil, fmt.Errorf("failed to create NodePool %w", response.Error())
	}
	return response.Body(), nil
}

func (client *OCMClient) GetNodePools(clusterID string) ([]*cmv1.NodePool, error) {
	response, err := client.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, fmt.Errorf("failed to get NodePools %w", response.Error())
	}
	return response.Items().Slice(), nil
}

func (client *OCMClient) GetNodePool(clusterID string, nodePoolID string) (*cmv1.NodePool, bool, error) {
	response, err := client.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		NodePool(nodePoolID).
		Get().
		Send()
	if response.Status() == 404 {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get NodePool %w", response.Error())
	}
	return response.Body(), true, nil
}

func (client *OCMClient) UpdateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := client.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePool.ID()).
		Update().Body(nodePool).
		Send()
	if err != nil {
		return nil, fmt.Errorf("failed to update NodePool %w", response.Error())
	}
	return response.Body(), nil
}

func (client *OCMClient) DeleteNodePool(clusterID string, nodePoolID string) error {
	response, err := client.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePoolID).
		Delete().
		Send()
	if err != nil {
		return fmt.Errorf("failed to delete NodePool %w", response.Error())
	}
	return nil
}
