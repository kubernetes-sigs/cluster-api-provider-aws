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
	"time"

	corev1 "k8s.io/api/core/v1"
	"github.com/aws/smithy-go"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	rosaAWSClient "github.com/openshift/rosa/pkg/aws"
	rosaCFNetwork "github.com/openshift/rosa/cmd/create/network"
	"github.com/aws/aws-sdk-go-v2/aws"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api/util/conditions"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

// RosaNetworkReconciler reconciles a RosaNetwork object.
type RosaNetworkReconciler struct {
	client.Client
	Endpoints     []scope.ServiceEndpoint
	Log           logr.Logger
	Scheme        *runtime.Scheme
	awsClient     rosaAWSClient.Client
}

const RosaNetworkFinalizer = "rosanetwork.infrastructure.cluster.x-k8s.io"

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks/status,verbs=get;update;patch

func (r *RosaNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
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

	// Create a new AWS/CloudFormation Client using the session cache
	session := rosaNetworkScope.SessionV2()
	logger := rosaNetworkScope.Logger.GetLogger()
	awsClient, err := rosaAWSClient.NewClient().
		LogrLogger(&logger).
		ExternalConfig(&session).
		Build()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("Failed to create AWS Client: %w", err)
	} else {
		r.awsClient = awsClient
	}

	// Always close the scope
	defer func() {
		if err := rosaNetworkScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, rosaNetworkScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaNetworkScope)
}

func (r *RosaNetworkReconciler) reconcileNormal(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	rosaNetScope.Info("Reconciling RosaNetwork")

	if controllerutil.AddFinalizer(rosaNetScope.RosaNetwork, RosaNetworkFinalizer) {
		if err := rosaNetScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Try to fetch CF stack with a given name
	stackName := rosaNetScope.RosaNetwork.Spec.Name
	cfStack, err := r.awsClient.GetCFStack(ctx, stackName)
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
		_, err = r.awsClient.CreateStackWithParamsTags(ctx, templateBody, stackName, cfParams, cfTags)
		if err != nil {
			conditions.MarkFalse(rosaNetScope.RosaNetwork,
				expinfrav1.RosaNetworkReadyCondition,
				expinfrav1.RosaNetworkFailedReason,
				clusterv1.ConditionSeverityError,
				"%s",
				err.Error())
			return ctrl.Result{}, fmt.Errorf("Failed to start CF stack creation: %w", err)
		} else {
			conditions.MarkFalse(rosaNetScope.RosaNetwork,
				expinfrav1.RosaNetworkReadyCondition,
				expinfrav1.RosaNetworkCreatingReason,
				clusterv1.ConditionSeverityInfo,
				"")
			return ctrl.Result{}, nil
		}
	} else { // The CF stack already exists
		if err := r.updateRosaNetworkResources(ctx, rosaNetScope.RosaNetwork); err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 60}, fmt.Errorf("Error fetching CF stack resources: %w", err)
		}

		switch cfStack.StackStatus {
		case cloudformationtypes.StackStatusCreateInProgress: // Create in progress
			// Set the reason of false RosaNetworkReadyCondition to Creating
			conditions.MarkFalse(rosaNetScope.RosaNetwork,
				expinfrav1.RosaNetworkReadyCondition,
				expinfrav1.RosaNetworkCreatingReason,
				clusterv1.ConditionSeverityInfo,
				"")
			return ctrl.Result{RequeueAfter: time.Second * 60}, nil
		case cloudformationtypes.StackStatusCreateComplete: // Create complete
			if err := r.parseSubnets(ctx, rosaNetScope.RosaNetwork); err != nil {
				return ctrl.Result{}, fmt.Errorf("Parsing stack subnets failed: %w", err)
			}

			// Set the reason of true RosaNetworkReadyCondition to Created
			// We have to use conditions.Set(), since conditions.MarkTrue() does not support setting reason
			conditions.Set(rosaNetScope.RosaNetwork,
				&clusterv1.Condition{
					Type:     expinfrav1.RosaNetworkReadyCondition,
					Status:   corev1.ConditionTrue,
					Reason:   expinfrav1.RosaNetworkCreatedReason,
					Severity: clusterv1.ConditionSeverityInfo,
					// FIXME: Message:  "",
				})
			return ctrl.Result{}, nil
		case cloudformationtypes.StackStatusCreateFailed: // Create failed
			// Set the reason of false RosaNetworkReadyCondition to Failed
			conditions.MarkFalse(rosaNetScope.RosaNetwork,
				expinfrav1.RosaNetworkReadyCondition,
				expinfrav1.RosaNetworkFailedReason,
				clusterv1.ConditionSeverityError,
				"")
			// FIXME: fix the error msg in fmt.Errorf()
			return ctrl.Result{}, fmt.Errorf("CloudFormation stack %s creation failed: ...", cfStack.StackName)
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) reconcileDelete(ctx context.Context, rosaNetScope *scope.RosaNetworkScope) (res ctrl.Result, reterr error) {
	rosaNetScope.Info("Reconciling RosaNetwork delete")

	// Try to fetch CF stack with a given name
	stackName := rosaNetScope.RosaNetwork.Spec.Name
	cfStack, err := r.awsClient.GetCFStack(ctx, stackName)
	if err != nil {
		var apiErr smithy.APIError // in case the stack does not exist, AWS returns ValidationError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ValidationError" {
			cfStack = nil
		} else {
			return ctrl.Result{}, fmt.Errorf("Error fetching CF stack details: %w", err)
		}
	}

	if cfStack != nil { // The CF stack still exists
		if err := r.updateRosaNetworkResources(ctx, rosaNetScope.RosaNetwork); err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 60}, fmt.Errorf("Error fetching CF stack resources: %w", err)
		}

		switch cfStack.StackStatus {
		case cloudformationtypes.StackStatusDeleteInProgress: // Deletion in progress
			return ctrl.Result{RequeueAfter: time.Second * 60}, nil
		case cloudformationtypes.StackStatusDeleteFailed:     // Deletion failed
			conditions.MarkFalse(rosaNetScope.RosaNetwork,
				expinfrav1.RosaNetworkReadyCondition,
				expinfrav1.RosaNetworkDeletionFailedReason,
				clusterv1.ConditionSeverityError,
				"")
			return ctrl.Result{}, fmt.Errorf("CF stack deletion failed")
		default:                                              // All the other states
			err = r.awsClient.DeleteCFStack(ctx, stackName)
			if err != nil {
				conditions.MarkFalse(rosaNetScope.RosaNetwork,
					expinfrav1.RosaNetworkReadyCondition,
					expinfrav1.RosaNetworkDeletionFailedReason,
					clusterv1.ConditionSeverityError,
					"%s",
					err.Error())
				return ctrl.Result{}, fmt.Errorf("Failed to start CF stack deletion: %w", err)
			} else {
				conditions.MarkFalse(rosaNetScope.RosaNetwork,
					expinfrav1.RosaNetworkReadyCondition,
					expinfrav1.RosaNetworkDeletingReason,
					clusterv1.ConditionSeverityInfo,
					"")
				return ctrl.Result{RequeueAfter: time.Second * 60}, nil
			}
		}
	} else {
		controllerutil.RemoveFinalizer(rosaNetScope.RosaNetwork, RosaNetworkFinalizer)
	}

	return ctrl.Result{}, nil
}

func (r *RosaNetworkReconciler) updateRosaNetworkResources(ctx context.Context, rosaNet *expinfrav1.RosaNetwork) error {
		resources, err := r.awsClient.DescribeCFStackResources(ctx, rosaNet.Spec.Name)
		if err != nil {
			return err
		}

		rosaNet.Status.Resources = make([]expinfrav1.CFResource, len(*resources))
		for i, resource := range *resources {
			rosaNet.Status.Resources[i] = expinfrav1.CFResource{
				LogicalId:    aws.ToString(resource.LogicalResourceId),
				PhysicalId:   aws.ToString(resource.PhysicalResourceId),
				ResourceType: aws.ToString(resource.ResourceType),
				Status:       fmt.Sprintf("%s", resource.ResourceStatus),
				Reason:       aws.ToString(resource.ResourceStatusReason),
			}
		}

		return nil
}

func (r *RosaNetworkReconciler) parseSubnets(ctx context.Context, rosaNet *expinfrav1.RosaNetwork) error {
	subnets := make(map[string]*expinfrav1.RosaNetworkSubnet)

	for _, resource := range rosaNet.Status.Resources {
		if (resource.ResourceType != "AWS::EC2::Subnet") { // Skip all non subnets
			continue
		}

		az, err := r.awsClient.GetSubnetAvailabilityZone(resource.PhysicalId)
		if err != nil {
			return err
		}

		if subnets[az] == nil {
			subnets[az] = &expinfrav1.RosaNetworkSubnet{
				AvailabilityZone: az,
				PublicSubnet:     "",
				PrivateSubnet:    "",
			}
		}

		if resource.LogicalId[0:13] == "SubnetPrivate" {
			subnets[az].PrivateSubnet = resource.PhysicalId
		} else {
			subnets[az].PublicSubnet = resource.PhysicalId
		} 
	}

	rosaNet.Status.Subnets = make([]expinfrav1.RosaNetworkSubnet, len(subnets))
	i := 0
	for _, v := range subnets {
		rosaNet.Status.Subnets[i] = *v
		i++
	}

	return nil
}

//func (r *RosaNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
func (r *RosaNetworkReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&expinfrav1.RosaNetwork{}).
		Complete(r)
}
