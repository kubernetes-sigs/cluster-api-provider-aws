/*
Copyright 2021 The Kubernetes Authors.

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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AWSClusterStaticIdentityReconciler reconciles a AWSClusterIdentity object.
type AWSClusterStaticIdentityReconciler struct {
	client.Client
	Log              logr.Logger
	Endpoints        []scope.ServiceEndpoint
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusterstaticidentities,verbs=get;list;watch;update

func (r *AWSClusterStaticIdentityReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the AWSCluster instance
	awsCluster := &infrav1.AWSCluster{}
	err := r.Get(ctx, req.NamespacedName, awsCluster)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	// Fetch all static identities
	staticIdentities := &infrav1.AWSClusterStaticIdentityList{}
	err = r.List(ctx, staticIdentities)
	if err != nil {
		return reconcile.Result{}, err
	}

	ownerRef := metav1.OwnerReference{
		APIVersion: awsCluster.APIVersion,
		Kind:       awsCluster.Kind,
		Name:       awsCluster.Name,
		UID:        awsCluster.UID,
	}

	for _, identity := range staticIdentities.Items {
		identity := identity
		if util.HasOwnerRef(identity.OwnerReferences, ownerRef) && awsCluster.Spec.IdentityRef.Name != identity.Name {
			identity.OwnerReferences = util.RemoveOwnerRef(identity.OwnerReferences, ownerRef)
			controllerutil.RemoveFinalizer(&identity, infrav1.AWSStaticIdentityFinalizer)
			err = r.Client.Update(ctx, &identity)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.V(2).Info("Updating the ownerRefs for staticIdentity ", identity.Name)
		}
	}

	return ctrl.Result{}, nil
}

func (r *AWSClusterStaticIdentityReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	awsController := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.AWSCluster{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue))

	return awsController.Complete(r)
}
