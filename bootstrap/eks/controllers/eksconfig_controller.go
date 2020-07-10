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

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"

	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
)

// EKSConfigReconciler reconciles a EKSConfig object
type EKSConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type EKSConfigScope struct {
	logr.Logger
	Config      *bootstrapv1.EKSConfig
	ConfigOwner *bsutil.ConfigOwner
	Cluster     *clusterv1.Cluster
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs/status,verbs=get;update;patch

func (r *EKSConfigReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, rerr error) {
	ctx := context.Background()
	log := r.Log.WithValues("eksconfig", req.NamespacedName)

	// get EKSConfig
	config := &bootstrapv1.EKSConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}

	// check owner references and look up owning Machine object
	configOwner, err := bsutil.GetConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// no error here, requeue until we find an owner
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "Failed to get owner")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		// no error, requeue until we find an owner
		return ctrl.Result{}, nil
	}

	log = log.WithValues("kind", configOwner.GetKind(), "version", configOwner.GetResourceVersion(), "name", configOwner.GetName())

	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Cause(err) == util.ErrNoCluster {
			log.Info(fmt.Sprintf("%s does not belong to a cluster yet, requeueing until it's part of a cluster", configOwner.GetKind()))
			return ctrl.Result{}, nil
		}
		if apierrors.IsNotFound(err) {
			log.Info("Cluster does not exist yet, requeueing until it is created")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Could not get cluster with metadata")
		return ctrl.Result{}, err
	}

	if annotations.IsPaused(cluster, config) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	scope := &EKSConfigScope{
		Logger:      log,
		Config:      config,
		ConfigOwner: configOwner,
		Cluster:     cluster,
	}

	patchHelper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// set up defer block for updating config
	defer func() {
		// TODO: update status conditions
		patchOpts := []patch.Option{}
		if rerr == nil {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}
		if err := patchHelper.Patch(ctx, config, patchOpts...); err != nil {
			log.Error(rerr, "Failed to patch config")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	if !cluster.Status.InfrastructureReady {
		log.Info("Cluster infrastructure is not ready, requeueing")
		// TODO: set condition
		return ctrl.Result{}, nil
	}

	if !cluster.Status.ControlPlaneInitialized {
		log.Info("Cluster has not yet been initialized, requeueing")
		// TODO? set condition
		return ctrl.Result{}, nil
	}

	return r.joinWorker(ctx, scope)
}

func (r *EKSConfigReconciler) joinWorker(ctx context.Context, scope *EKSConfigScope) (ctrl.Result, error) {
	// generate userdata
	// store userdata as secret (this can basically be totally copied from kubeadm provider - any way to reuse?)
	// set status.DataSecretName
	// set status.Ready to true
	// mark DataSecretAvailableCondition as true

	return ctrl.Result{}, nil
}

func (r *EKSConfigReconciler) SetupWithManager(mgr ctrl.Manager, option controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bootstrapv1.EKSConfig{}).
		WithOptions(option).
		Complete(r)
}
