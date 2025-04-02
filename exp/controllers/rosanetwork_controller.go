/*
Copyright The Kubernetes Authors.

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
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/smithy-go"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	rosaAWSClient "github.com/openshift/rosa/pkg/aws"
	rosaCFNetwork "github.com/openshift/rosa/cmd/create/network"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

// RosaNetworkReconciler reconciles a RosaNetwork object.
type RosaNetworkReconciler struct {
	client.Client
	Endpoints []scope.ServiceEndpoint
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks/status,verbs=get;update;patch
func (r *RosaNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	// _ = context.Background()
	//_ = r.Log.WithValues("rosanetwork", req.NamespacedName)
	log := logger.FromContext(ctx)

	// Get the rosanetwork instance
	rosaNetwork := &expinfrav1.RosaNetwork{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaNetwork); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	rosaNetworkScope, err := scope.NewRosaNetworkScope(scope.RosaNetworkScopeParams{
		Client:         r.Client,
		RosaNetwork:    rosaNetwork,
		ControllerName: "rosanetwork",
		Endpoints:      r.Endpoints,
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosanetwork scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := rosaNetworkScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, rosaNetworkScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaNetworkScope)
}

func (r *RosaNetworkReconciler) reconcileNormal(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	rosaNetScope.Info("Reconciling RosaNetwork")

	// FIXME: add finalizer
	// if controllerutil.AddFinalizer() {
	// 	if err := rosaNetScope.PatchObject(); err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	// }

	// Create a new AWS/CloudFormation Client using the session cache
	session := rosaNetScope.SessionV2()
	logger := rosaNetScope.Logger.GetLogger()
	awsClient, err := rosaAWSClient.NewClient().
		LogrLogger(&logger).
		ExternalConfig(&session).
		Build()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("Failed to create AWS Client: %w", err)
	}

	cfStack, err := awsClient.GetCFStack(ctx, rosaNetScope.RosaNetwork.Spec.Name)
	if err != nil {
		var apiErr smithy.APIError // in case the stack does not exist, AWS returns ValidationError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ValidationError" {
			cfStack = nil
		} else {
			return ctrl.Result{}, fmt.Errorf("Error fetching CF stack details: %w", err)
		}
	}

	if cfStack == nil { // The CF stack does not exist yet
		templateBody := string(rosaCFNetwork.CloudFormationTemplateFile)
		stackName := rosaNetScope.RosaNetwork.Spec.Name
		cfParams := map[string]string{
			"AvailabilityZoneCount": strconv.Itoa(rosaNetScope.RosaNetwork.Spec.AvailabilityZoneCount),
			"Region":                rosaNetScope.RosaNetwork.Spec.Region,
			"Name":                  rosaNetScope.RosaNetwork.Spec.Name,
			"VpcCidr":               rosaNetScope.RosaNetwork.Spec.CIDRBlock,
		}
		// Explicitly specified AZs
		for i := 1; i <= len(rosaNetScope.RosaNetwork.Spec.AvailabilityZones); i++ {
			cfParams[fmt.Sprintf("AZ%s", i)] = rosaNetScope.RosaNetwork.Spec.AvailabilityZones[i-1]
		}
		cfTags := map[string]string{}

		// Call the AWS CF stack create API
		_, err = awsClient.CreateStackWithParamsTags(ctx, templateBody, stackName, cfParams, cfTags)
		if err != nil {
			// FIXME: update conditions
			return ctrl.Result{}, fmt.Errorf("failed to create CF stack: %w", err)
		} else { // The CF stack already exists
		}
	}

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) reconcileDelete(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	// TODO
	return ctrl.Result{}, nil
}

//func (r *RosaNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
func (r *RosaNetworkReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&expinfrav1.RosaNetwork{}).
		Complete(r)
}
