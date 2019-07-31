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

package awscluster

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb"

	"sigs.k8s.io/cluster-api/pkg/util"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	infrastructurev1alpha2 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName  = "awscluster-controller"
	apiEndpointPort = 6443
)

// Add creates a new AWSCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAWSCluster{
		Client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		recorder: mgr.GetEventRecorderFor(controllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to AWSCluster
	err = c.Watch(
		&source.Kind{Type: &infrastructurev1alpha2.AWSCluster{}},
		&handler.EnqueueRequestForObject{},
	)
	if err != nil {
		return err
	}

	// Watch for changes to Cluster and enqueue the associated AWSCluster
	return c.Watch(
		&source.Kind{Type: &v1alpha2.Cluster{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: util.ClusterToInfrastructureMapFunc(schema.GroupVersionKind{
				Group:   infrastructurev1alpha2.SchemeGroupVersion.Group,
				Version: infrastructurev1alpha2.SchemeGroupVersion.Version,
				Kind:    "AWSCluster",
			}),
		},
	)
}

var _ reconcile.Reconciler = &ReconcileAWSCluster{}

// ReconcileAWSCluster reconciles a AWSCluster object
type ReconcileAWSCluster struct {
	client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a AWSCluster object and makes changes based on the state read
// and what is in the AWSCluster.Spec
func (r *ReconcileAWSCluster) Reconcile(request reconcile.Request) (_ reconcile.Result, reterr error) {
	ctx := context.TODO()
	logger := log.Log.WithName(controllerName).
		WithName(fmt.Sprintf("namespace=%s", request.Namespace)).
		WithName(fmt.Sprintf("awsCluster=%s", request.Name))

	// Fetch the AWSCluster instance
	awsCluster := &infrastructurev1alpha2.AWSCluster{}
	err := r.Get(ctx, request.NamespacedName, awsCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger = logger.WithName(awsCluster.APIVersion)

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		logger.Info("Waiting for Cluster Controller to set OwnerRef on AWSCluuster")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	logger = logger.WithName(fmt.Sprintf("cluster=%s", cluster.Name))

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:     r.Client,
		Logger:     logger,
		Cluster:    cluster,
		AWSCluster: awsCluster,
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
	clusterScope.Info("Reconciling Cluster delete")

	ec2svc := ec2.NewService(clusterScope)
	elbsvc := elb.NewService(clusterScope)

	if err := elbsvc.DeleteLoadbalancers(); err != nil {
		clusterScope.Error(err, "Error deleting cluster load balancer")
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if err := ec2svc.DeleteBastion(); err != nil {
		clusterScope.Error(err, "Error deleting cluster bastion")
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if err := ec2svc.DeleteNetwork(); err != nil {
		clusterScope.Error(err, "Error deleting cluster network")
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	}

	// Cluster is deleted so remove the finalizer.
	clusterScope.AWSCluster.Finalizers = util.Filter(clusterScope.AWSCluster.Finalizers, infrastructurev1alpha2.ClusterFinalizer)

	return reconcile.Result{}, nil
}

// TODO(ncdc): should this be a function on ClusterScope?
func reconcileNormal(scope *scope.ClusterScope) (reconcile.Result, error) {
	scope.Info("Reconciling Cluster")

	awsCluster := scope.AWSCluster

	// If the AWSCluster doesn't have our finalizer, add it.
	if !util.Contains(awsCluster.Finalizers, infrastructurev1alpha2.ClusterFinalizer) {
		awsCluster.Finalizers = append(awsCluster.Finalizers, infrastructurev1alpha2.ClusterFinalizer)
	}

	ec2Service := ec2.NewService(scope)
	elbService := elb.NewService(scope)

	if err := ec2Service.ReconcileNetwork(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile network for cluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := ec2Service.ReconcileBastion(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile bastion host for cluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if err := elbService.ReconcileLoadbalancers(); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile load balancers for cluster %s/%s", awsCluster.Namespace, awsCluster.Name)
	}

	if awsCluster.Status.Network.APIServerELB.DNSName == "" {
		scope.Info("Waiting on API server ELB DNS name")
		return reconcile.Result{RequeueAfter: 15 * time.Second}, nil
	}

	// Set APIEndpoints so the Cluster API Cluster Controller can pull them
	awsCluster.Status.APIEndpoints = []infrastructurev1alpha2.APIEndpoint{
		{
			Host: awsCluster.Status.Network.APIServerELB.DNSName,
			// TODO(ncdc): should this come from awsCluster.Status.Network.APIServerELB.Listeners[0].Port?
			Port: apiEndpointPort,
		},
	}

	// No errors, so mark us ready so the Cluster API Cluster Controller can pull it
	awsCluster.Status.Ready = true

	return reconcile.Result{}, nil
}
