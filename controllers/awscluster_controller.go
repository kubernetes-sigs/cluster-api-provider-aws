/*
Copyright 2019 The Kubernetes Authors.

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
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AWSClusterReconciler reconciles a AwsCluster object
type AWSClusterReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *AWSClusterReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.TODO()
	log := r.Log.WithValues("namespace", req.Namespace, "awsCluster", req.Name)

	// Fetch the AWSCluster instance
	awsCluster := &infrav1.AWSCluster{}
	err := r.Get(ctx, req.NamespacedName, awsCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:     r.Client,
		Logger:     log,
		Cluster:    cluster,
		AWSCluster: awsCluster,
		Recorder:   r.Recorder,
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any AWSMachine changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted clusters
	if !awsCluster.DeletionTimestamp.IsZero() {
		return reconcileDelete(clusterScope)
	}

	// Handle non-deleted clusters
	return reconcileNormal(clusterScope)
}

// TODO(ncdc): should this be a function on ClusterScope?
func reconcileDelete(clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	clusterScope.Info("Reconciling AWSCluster delete")

	ec2svc := ec2.NewService(clusterScope)
	elbsvc := elb.NewService(clusterScope)
	awsCluster := clusterScope.AWSCluster

	if err := elbsvc.DeleteLoadbalancers(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting load balancer for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := ec2svc.DeleteBastion(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting bastion for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := ec2svc.DeleteNetwork(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting network for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	// Cluster is deleted so remove the finalizer.
	clusterScope.AWSCluster.Finalizers = util.Filter(clusterScope.AWSCluster.Finalizers, infrav1.ClusterFinalizer)

	return reconcile.Result{}, nil
}

// TODO(ncdc): should this be a function on ClusterScope?
func reconcileNormal(clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	clusterScope.Info("Reconciling AWSCluster")

	awsCluster := clusterScope.AWSCluster

	// If the AWSCluster doesn't have our finalizer, add it.
	if !util.Contains(awsCluster.Finalizers, infrav1.ClusterFinalizer) {
		awsCluster.Finalizers = append(awsCluster.Finalizers, infrav1.ClusterFinalizer)
	}

	ec2Service := ec2.NewService(clusterScope)
	elbService := elb.NewService(clusterScope)

	if err := ec2Service.ReconcileNetwork(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile network for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := ec2Service.ReconcileBastion(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile bastion host for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := elbService.ReconcileLoadbalancers(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile load balancers for AWSCluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if awsCluster.Status.Network.APIServerELB.DNSName == "" {
		clusterScope.Info("Waiting on API server ELB DNS name")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	// Set APIEndpoints so the Cluster API Cluster Controller can pull them
	// TODO: should we get the Port from the first listener on the ELB?
	awsCluster.Status.APIEndpoints = []infrav1.APIEndpoint{
		{
			Host: awsCluster.Status.Network.APIServerELB.DNSName,
			Port: int(clusterScope.APIServerPort()),
		},
	}

	// No errors, so mark us ready so the Cluster API Cluster Controller can pull it
	awsCluster.Status.Ready = true

	return reconcile.Result{}, nil
}

func (r *AWSClusterReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AWSCluster{}).
		Complete(r)
}
