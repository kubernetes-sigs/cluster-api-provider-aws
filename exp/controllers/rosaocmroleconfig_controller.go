/*
Copyright 2026 The Kubernetes Authors.

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
	"sync"

	"github.com/openshift/rosa/pkg/aws"
	rosalogging "github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/ocmrole"
	"github.com/openshift/rosa/pkg/reporter"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsiface "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	caparosa "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAOCMRoleConfigReconciler reconciles an ROSAOCMRoleConfig object.
type ROSAOCMRoleConfigReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSClient
	NewOCMClient     func(ctx context.Context, scope caparosa.OCMSecretsRetriever) (caparosa.OCMClient, error)
	Runtime          *rosacli.Runtime

	// Per-organization linking mutex
	orgLinkMutexes sync.Map // map[orgID]*sync.Mutex
}

func (r *ROSAOCMRoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = caparosa.NewWrappedOCMClientWithoutControlPlane
	r.NewStsClient = scope.NewSTSClient

	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSAOCMRoleConfig{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaocmroleconfigs/finalizers,verbs=update

// Reconcile reconciles ROSAOCMRoleConfig.
func (r *ROSAOCMRoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	ocmRoleConfig := &expinfrav1.ROSAOCMRoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, ocmRoleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ROSAOCMRoleConfig")
		return ctrl.Result{Requeue: true}, nil
	}

	log = log.WithValues("ocmRoleConfig", klog.KObj(ocmRoleConfig))
	scope, err := scope.NewROSAOCMRoleConfigScope(scope.ROSAOCMRoleConfigScopeParams{
		Client:            r.Client,
		ROSAOCMRoleConfig: ocmRoleConfig,
		ControllerName:    "rosaocmroleconfig",
		Logger:            log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosaocmroleconfig scope: %w", err)
	}

	// Always close scope and set summary condition
	defer func() {
		v1beta1conditions.SetSummary(scope.ROSAOCMRoleConfig,
			v1beta1conditions.WithConditions(expinfrav1.ROSAOCMRoleConfigReadyCondition),
			v1beta1conditions.WithStepCounter())
		if err := scope.PatchObject(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if err := r.setUpRuntime(ctx, scope); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up runtime: %w", err)
	}

	if !scope.ROSAOCMRoleConfig.DeletionTimestamp.IsZero() {
		scope.Info("Deleting ROSAOCMRoleConfig.")
		v1beta1conditions.MarkFalse(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition, expinfrav1.ROSAOCMRoleConfigDeletionStarted, clusterv1beta1.ConditionSeverityInfo, "Deletion of ROSAOCMRoleConfig started")
		err = r.reconcileDelete(scope)
		if err == nil {
			controllerutil.RemoveFinalizer(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigFinalizer)
		}
		return ctrl.Result{}, err
	}

	if controllerutil.AddFinalizer(scope.ROSAOCMRoleConfig, expinfrav1.ROSAOCMRoleConfigFinalizer) {
		return ctrl.Result{}, nil
	}

	return r.reconcileOCMRole(ctx, scope)
}

func (r *ROSAOCMRoleConfigReconciler) reconcileOCMRole(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) (ctrl.Result, error) {
	ocmRoleConfig := scope.ROSAOCMRoleConfig

	// Get current OCM organization
	orgID, externalID, err := r.Runtime.OCMClient.GetCurrentOrganization()
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to get OCM organization: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to get OCM organization: %w", err)
	}

	if ocmRoleConfig.Status.OrganizationID == "" {
		ocmRoleConfig.Status.OrganizationID = orgID
	}

	// Build role name
	roleName := aws.GetOCMRoleName(ocmRoleConfig.Spec.RolePrefix, aws.OCMRole, externalID)

	// Convert profile from CRD enum to ROSA CLI enum
	rosaProfile, err := r.convertProfile(ocmRoleConfig.Spec.Profile)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Invalid profile: %s", err.Error())
		return ctrl.Result{}, err
	}

	roleARN, created, err := ocmrole.GetOrCreateOCMRole(
		r.Runtime,
		ocmRoleConfig.Spec.RolePrefix,
		rosaProfile,
		ocmRoleConfig.Spec.PermissionsBoundaryARN,
		ocmRoleConfig.Spec.Path,
		false, // managed policies - we only support customer managed policies on OCM role as today.
	)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to get or create OCM role: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to get or create OCM role: %w", err)
	}

	if created {
		scope.Info("Successfully created OCM role", "roleARN", roleARN)
		r.Recorder.Eventf(ocmRoleConfig, "Normal", expinfrav1.ROSAOCMRoleConfigCreatedReason,
			"Created OCM role: %s", roleARN)
	} else {
		scope.Info("Using existing OCM role", "roleARN", roleARN)
	}

	// Update status with role ARN
	ocmRoleConfig.Status.RoleARN = roleARN

	// Acquire per-organization mutex to prevent concurrent LinkOrgToRole races
	// Multiple ROSAOCMRoleConfigs with different AWS accounts can link to the same org,
	// causing read-modify-write races on the org label. Serialize linking per org.
	linkMu := r.getLinkMutex(orgID)
	linkMu.Lock()
	defer linkMu.Unlock()

	// Check if role is linked to organization
	existsOnOCM, _, linkedARN, err := r.Runtime.OCMClient.CheckRoleExists(orgID, roleName, r.Runtime.Creator.AccountID)
	if err != nil {
		v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
			expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
			"Failed to check OCM role link status: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to check OCM role link status: %w", err)
	}

	// Link role if not already linked
	if !existsOnOCM || linkedARN != roleARN {
		scope.Info("Linking OCM role to organization", "roleARN", roleARN, "orgID", orgID)

		linked, err := r.Runtime.OCMClient.LinkOrgToRole(orgID, roleARN)
		if err != nil {
			v1beta1conditions.MarkFalse(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition,
				expinfrav1.ROSAOCMRoleConfigReconciliationFailedReason, clusterv1beta1.ConditionSeverityError,
				"Failed to link OCM role: %s", err.Error())
			return ctrl.Result{}, fmt.Errorf("failed to link OCM role to organization: %w", err)
		}

		if linked {
			scope.Info("Successfully linked OCM role to organization")
			r.Recorder.Eventf(ocmRoleConfig, "Normal", expinfrav1.ROSAOCMRoleConfigLinkedReason,
				"Linked OCM role %s to organization %s", roleARN, orgID)
		}
	}

	// Set ready condition
	v1beta1conditions.MarkTrue(ocmRoleConfig, expinfrav1.ROSAOCMRoleConfigReadyCondition)

	scope.Info("Successfully reconciled ROSAOCMRoleConfig", "roleARN", roleARN)
	return ctrl.Result{}, nil
}

func (r *ROSAOCMRoleConfigReconciler) reconcileDelete(scope *scope.ROSAOCMRoleConfigScope) error {
	ocmRoleConfig := scope.ROSAOCMRoleConfig

	orgID := ocmRoleConfig.Status.OrganizationID
	if orgID == "" {
		scope.Info("No organization ID in status, skipping cleanup")
		return nil
	}

	roleARN := ocmRoleConfig.Status.RoleARN

	if roleARN == "" {
		scope.Info("No role ARN in status, skipping cleanup")
		return nil
	}

	roleName, err := aws.GetResourceIdFromARN(roleARN)
	if err != nil {
		return fmt.Errorf("failed to parse role ARN: %w", err)
	}

	// Acquire per-organization mutex to prevent concurrent unlink/link races
	// Multiple ROSAOCMRoleConfigs with different AWS accounts can link to the same org,
	// and we need to prevent a race where one CR is unlinking while another is linking.
	linkMu := r.getLinkMutex(orgID)
	linkMu.Lock()

	// Check if role is linked to organization
	existsOnOCM, _, linkedARN, err := r.Runtime.OCMClient.CheckRoleExists(orgID, roleName, r.Runtime.Creator.AccountID)
	if err != nil {
		linkMu.Unlock()
		return fmt.Errorf("failed to check role link status: %w", err)
	}

	if existsOnOCM && linkedARN == roleARN {
		err := r.Runtime.OCMClient.UnlinkOCMRoleFromOrg(orgID, roleARN)
		if err != nil {
			linkMu.Unlock()
			return fmt.Errorf("failed to unlink the role: %w", err)
		}
		scope.Info("Successfully unlinked role", "roleARN", roleARN, "orgID", orgID)
	}

	// Release mutex - OCM operations are complete, IAM deletion doesn't need it
	linkMu.Unlock()

	// Delete IAM role
	roleExists, _, err := r.Runtime.AWSClient.CheckRoleExists(roleName)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}

	if roleExists {
		err = r.Runtime.AWSClient.DeleteOCMRole(roleName, false)
		if err != nil {
			return fmt.Errorf("failed to delete IAM role: %w", err)
		}
		scope.Info("Successfully deleted OCM role", "roleName", roleName)
	}

	return nil
}

func (r *ROSAOCMRoleConfigReconciler) setUpRuntime(ctx context.Context, scope *scope.ROSAOCMRoleConfigScope) error {
	if r.Runtime != nil {
		return nil
	}

	// Create OCM client
	ocm, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := caparosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	runtime := rosacli.NewRuntime()
	runtime.OCMClient = ocmClient
	runtime.Reporter = reporter.CreateReporter()
	runtime.Logger = rosalogging.NewLogger()

	session := scope.Session()
	awsClient, err := aws.NewClient().
		Logger(runtime.Logger).
		ExternalConfig(&session).
		Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}
	runtime.AWSClient = awsClient

	creator, err := awsClient.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get creator: %w", err)
	}
	runtime.Creator = creator

	// atomic assignment - only when fully initialized
	r.Runtime = runtime
	return nil
}

func (r *ROSAOCMRoleConfigReconciler) convertProfile(profile expinfrav1.ROSAOCMRoleProfile) (ocmrole.RoleProfile, error) {
	switch profile {
	case expinfrav1.ROSAOCMRoleProfileStandard:
		return ocmrole.ProfileStandard, nil
	case expinfrav1.ROSAOCMRoleProfileAdmin:
		return ocmrole.ProfileAdmin, nil
	case expinfrav1.ROSAOCMRoleProfileNoConsole:
		return ocmrole.ProfileNoConsole, nil
	default:
		return "", fmt.Errorf("unknown profile: %s", profile)
	}
}

func (r *ROSAOCMRoleConfigReconciler) getLinkMutex(orgID string) *sync.Mutex {
	mu, _ := r.orgLinkMutexes.LoadOrStore(orgID, &sync.Mutex{})
	return mu.(*sync.Mutex)
}
