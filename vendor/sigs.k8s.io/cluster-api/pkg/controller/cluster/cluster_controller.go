/*
Copyright 2018 The Kubernetes Authors.

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

package cluster

import (
	"context"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
	capierrors "sigs.k8s.io/cluster-api/pkg/errors"
	"sigs.k8s.io/cluster-api/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var DefaultActuator Actuator

func AddWithActuator(mgr manager.Manager, actuator Actuator) error {
	return add(mgr, newReconciler(mgr, actuator))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, actuator Actuator) reconcile.Reconciler {
	return &ReconcileCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme(), actuator: actuator}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cluster_controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Cluster
	err = c.Watch(&source.Kind{Type: &clusterv1.Cluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCluster{}

// ReconcileCluster reconciles a Cluster object
type ReconcileCluster struct {
	client.Client
	scheme   *runtime.Scheme
	actuator Actuator
}

func (r *ReconcileCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	cluster := &clusterv1.Cluster{}
	err := r.Get(context.Background(), request.NamespacedName, cluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	name := cluster.Name
	klog.Infof("Running reconcile Cluster for %q", name)

	// If object hasn't been deleted and doesn't have a finalizer, add one
	// Add a finalizer to newly created objects.
	if cluster.ObjectMeta.DeletionTimestamp.IsZero() {
		finalizerCount := len(cluster.Finalizers)

		if !util.Contains(cluster.Finalizers, metav1.FinalizerDeleteDependents) {
			cluster.Finalizers = append(cluster.ObjectMeta.Finalizers, metav1.FinalizerDeleteDependents)
		}

		if !util.Contains(cluster.Finalizers, clusterv1.ClusterFinalizer) {
			cluster.Finalizers = append(cluster.ObjectMeta.Finalizers, clusterv1.ClusterFinalizer)
		}

		if len(cluster.Finalizers) > finalizerCount {
			if err := r.Update(context.Background(), cluster); err != nil {
				klog.Infof("Failed to add finalizer to cluster %q: %v", name, err)
				return reconcile.Result{}, err
			}

			// Since adding the finalizer updates the object return to avoid later update issues.
			return reconcile.Result{Requeue: true}, nil
		}

	}

	if !cluster.ObjectMeta.DeletionTimestamp.IsZero() {
		// no-op if finalizer has been removed.
		if !util.Contains(cluster.ObjectMeta.Finalizers, clusterv1.ClusterFinalizer) {
			klog.Infof("reconciling cluster object %v causes a no-op as there is no finalizer.", name)
			return reconcile.Result{}, nil
		}

		klog.Infof("reconciling cluster object %v triggers delete.", name)
		if err := r.actuator.Delete(cluster); err != nil {
			klog.Errorf("Error deleting cluster object %v; %v", name, err)
			return reconcile.Result{}, err
		}
		// Remove finalizer on successful deletion.
		klog.Infof("cluster object %v deletion successful, removing finalizer.", name)
		cluster.ObjectMeta.Finalizers = util.Filter(cluster.ObjectMeta.Finalizers, clusterv1.ClusterFinalizer)
		if err := r.Client.Update(context.Background(), cluster); err != nil {
			klog.Errorf("Error removing finalizer from cluster object %v; %v", name, err)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	klog.Infof("reconciling cluster object %v triggers idempotent reconcile.", name)
	if err := r.actuator.Reconcile(cluster); err != nil {
		if requeueErr, ok := errors.Cause(err).(capierrors.HasRequeueAfterError); ok {
			klog.Infof("Actuator returned requeue-after error: %v", requeueErr)
			return reconcile.Result{Requeue: true, RequeueAfter: requeueErr.GetRequeueAfter()}, nil
		}
		klog.Errorf("Error reconciling cluster object %v; %v", name, err)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
