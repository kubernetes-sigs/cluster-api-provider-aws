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

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks"
)

// AWSFargateProfileReconciler reconciles a AWSFargateProfile object
type AWSFargateProfileReconciler struct {
	client.Client
	Log       logr.Logger
	Recorder  record.EventRecorder
	Endpoints []scope.ServiceEndpoint

	EnableIAM bool
}

// SetupWithManager is used to setup the controller
func (r *AWSFargateProfileReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	managedControlPlaneToFargateProfileMap := managedControlPlaneToFargateProfileMapFunc(r.Client, r.Log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1exp.AWSFargateProfile{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPaused(r.Log)).
		Watches(
			&source.Kind{Type: &controlplanev1.AWSManagedControlPlane{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: managedControlPlaneToFargateProfileMap,
			},
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes;awsmanagedcontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles/status,verbs=get;update;patch

// Reconcile reconciles AWSFargateProfiles
func (r *AWSFargateProfileReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := r.Log.WithValues("namespace", req.Namespace, "AWSFargateProfile", req.Name)
	ctx := context.Background()

	fargateProfile := &infrav1exp.AWSFargateProfile{}
	if err := r.Get(ctx, req.NamespacedName, fargateProfile); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	cluster, err := util.GetClusterByName(ctx, r.Client, fargateProfile.Namespace, fargateProfile.Spec.ClusterName)
	if err != nil {
		logger.Info("Failed to retrieve Cluster from AWSFargateProfile")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("Cluster", cluster.Name)

	controlPlaneKey := client.ObjectKey{
		Namespace: fargateProfile.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &controlplanev1.AWSManagedControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		logger.Info("Failed to retrieve ControlPlane from AWSFargateProfile")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("AWSManagedControlPlane", controlPlane.Name)

	fargateProfileScope, err := scope.NewFargateProfileScope(scope.FargateProfileScopeParams{
		Logger:         logger,
		Client:         r.Client,
		ControllerName: "awsfargateprofile",
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		FargateProfile: fargateProfile,
		EnableIAM:      r.EnableIAM,
		Endpoints:      r.Endpoints,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	defer func() {
		applicableConditions := []clusterv1.ConditionType{
			infrav1exp.IAMFargateRolesReadyCondition,
			infrav1exp.EKSFargateProfileReadyCondition,
		}

		conditions.SetSummary(fargateProfileScope.FargateProfile, conditions.WithConditions(applicableConditions...), conditions.WithStepCounter())

		if err := fargateProfileScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !controlPlane.Status.Ready {
		logger.Info("Control plane is not ready yet")
		conditions.MarkFalse(fargateProfile, clusterv1.ReadyCondition, infrav1exp.WaitingForEKSControlPlaneReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if !fargateProfile.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, fargateProfileScope)
	}

	return r.reconcileNormal(ctx, fargateProfileScope)
}

func (r *AWSFargateProfileReconciler) reconcileNormal(
	_ context.Context,
	fargateProfileScope *scope.FargateProfileScope,
) (ctrl.Result, error) {
	fargateProfileScope.Info("Reconciling AWSFargateProfile")

	controllerutil.AddFinalizer(fargateProfileScope.FargateProfile, infrav1exp.FargateProfileFinalizer)
	if err := fargateProfileScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	ekssvc := eks.NewFargateService(fargateProfileScope)

	res, err := ekssvc.Reconcile()
	if err != nil {
		return res, errors.Wrapf(err, "failed to reconcile fargate profile for AWSFargateProfile %s/%s", fargateProfileScope.FargateProfile.Namespace, fargateProfileScope.FargateProfile.Name)
	}

	return res, nil
}

func (r *AWSFargateProfileReconciler) reconcileDelete(
	_ context.Context,
	fargateProfileScope *scope.FargateProfileScope,
) (ctrl.Result, error) {
	fargateProfileScope.Info("Reconciling deletion of AWSFargateProfile")

	ekssvc := eks.NewFargateService(fargateProfileScope)

	res, err := ekssvc.ReconcileDelete()
	if err != nil {
		return res, errors.Wrapf(err, "failed to reconcile fargate profile deletion for AWSFargateProfile %s/%s", fargateProfileScope.FargateProfile.Namespace, fargateProfileScope.FargateProfile.Name)
	}

	if res.IsZero() {
		controllerutil.RemoveFinalizer(fargateProfileScope.FargateProfile, infrav1exp.FargateProfileFinalizer)
	}

	return res, nil
}

func managedControlPlaneToFargateProfileMapFunc(c client.Client, log logr.Logger) handler.ToRequestsFunc {
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

		fargateProfileForClusterList := infrav1exp.AWSFargateProfileList{}
		if err := c.List(
			ctx, &fargateProfileForClusterList, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterLabelName: clusterKey.Name},
		); err != nil {
			log.Error(err, "couldn't list fargate profiles for cluster")
			return nil
		}

		var results []ctrl.Request
		for i := range fargateProfileForClusterList.Items {
			fp := fargateProfileForClusterList.Items[i]
			results = append(results, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: fp.Namespace,
					Name:      fp.Name,
				},
			})
		}

		return results
	}
}
