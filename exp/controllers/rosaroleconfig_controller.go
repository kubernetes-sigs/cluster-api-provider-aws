/*
Copyright 2025 The Kubernetes Authors.

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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/go-logr/logr"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	accountroles "github.com/openshift/rosa/cmd/create/accountroles"
	oidcconfig "github.com/openshift/rosa/cmd/create/oidcconfig"
	oidcprovider "github.com/openshift/rosa/cmd/create/oidcprovider"
	operatorroles "github.com/openshift/rosa/cmd/create/operatorroles"
	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/helper"
	interactive "github.com/openshift/rosa/pkg/interactive"
	rosalogging "github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/ocm"
	rosacli "github.com/openshift/rosa/pkg/rosa"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsiface "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSARoleConfigReconciler reconciles a ROSARoleConfig object.
type ROSARoleConfigReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSClient
	NewOCMClient     func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error)
	Runtime          *rosacli.Runtime
}

func (r *ROSARoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClientWithoutControlPlane
	r.NewStsClient = scope.NewSTSClient

	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSARoleConfig{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs/finalizers,verbs=update

func (r *ROSARoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	roleConfig := &expinfrav1.ROSARoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, roleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get ROSARoleConfig")
		return ctrl.Result{Requeue: true}, nil
	}

	log = log.WithValues("cluster", klog.KObj(roleConfig))
	scope, err := scope.NewRosaRoleConfigScope(scope.RosaRoleConfigScopeParams{
		Client:         r.Client,
		RosaRoleConfig: roleConfig,
		ControllerName: "rosaroleconfig",
		Logger:         log,
	})

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosaroleconfig scope: %w", err)
	}

	// Always close the scope and set summary condition
	defer func() {
		conditions.SetSummary(scope.RosaRoleConfig, conditions.WithConditions(expinfrav1.RosaRoleConfigReadyCondition), conditions.WithStepCounter())

		// Delete is already patched
		if roleConfig.ObjectMeta.DeletionTimestamp.IsZero() {
			if err := scope.PatchObject(); err != nil {
				reterr = errors.Join(reterr, err)
			}
		}
	}()

	err = r.setUpRuntime(ctx, scope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up runtime: %w", err)
	}

	if !roleConfig.ObjectMeta.DeletionTimestamp.IsZero() {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionStarted, clusterv1.ConditionSeverityInfo, "Deletion of RosaRolesConfig started")
		return ctrl.Result{}, r.reconcileDelete(scope)
	}

	if controllerutil.AddFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer) {
		return ctrl.Result{}, err
	}

	err = r.createAccountRoles(roleConfig, scope)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Failed to create Account Roles: %v", err)
		return ctrl.Result{}, fmt.Errorf("failed to Create AccountRoles: %w", err)
	}

	err = r.reconcileOIDCConfig(roleConfig, scope)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Failed to create OIDC Config: %v", err)
		return ctrl.Result{}, fmt.Errorf("failed to OICD Config: %w", err)
	}

	err = r.createOIDCProvider(scope)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Failed to create OIDC provider: %v", err)
		return ctrl.Result{}, fmt.Errorf("failed to Create OIDC provider: %w", err)
	}

	err = r.createOperatorRoles(roleConfig, scope)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigReconciliationFailedReason, clusterv1.ConditionSeverityError, "Failed to create Operator Roles: %v", err)
		return ctrl.Result{}, fmt.Errorf("failed to Create OperatorRoles: %w", err)
	}

	if r.rosaRolesConfigReady(scope) {
		conditions.MarkTrue(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition)
		conditions.Set(scope.RosaRoleConfig,
			&clusterv1.Condition{
				Type:     expinfrav1.RosaRoleConfigReadyCondition,
				Status:   corev1.ConditionTrue,
				Reason:   expinfrav1.RosaRoleConfigCreatedReason,
				Severity: clusterv1.ConditionSeverityInfo,
				Message:  "RosaRoleConfig is ready to be used.",
			})
	}
	return ctrl.Result{}, nil
}

func (r *ROSARoleConfigReconciler) reconcileDelete(scope *scope.RosaRoleConfigScope) error {
	err := r.deleteOperatorRoles(scope.RosaRoleConfig.Spec.AccountRoleConfig.Prefix)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete operator roles: %v", err)
		return err
	}

	oidcID := scope.RosaRoleConfig.Status.OIDCID
	if scope.RosaRoleConfig.Spec.OperatorRoleConfig.OIDCID == "" {
		err = r.deleteOIDCProvider(oidcID)
		if err != nil {
			conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete OIDC provider: %v", err)
			return err
		}
	}

	err = r.deleteAccountRoles(scope)
	if err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete account roles: %v", err)
		return err
	}

	if scope.RosaRoleConfig.Spec.OperatorRoleConfig.OIDCID == "" {
		err = r.deleteOIDCConfig(oidcID)
		if err != nil {
			conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to delete OIDC config: %v", err)
			return err
		}
	}

	controllerutil.RemoveFinalizer(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigFinalizer)
	// Explicitly patch the object to persist the finalizer removal
	if err := scope.PatchObject(); err != nil {
		conditions.MarkFalse(scope.RosaRoleConfig, expinfrav1.RosaRoleConfigReadyCondition, expinfrav1.RosaRoleConfigDeletionFailedReason, clusterv1.ConditionSeverityError, "Failed to remove finalizer: %v", err)
		return fmt.Errorf("failed to remove finalizer from ROSARoleConfig: %w", err)
	}

	return nil
}

func (r *ROSARoleConfigReconciler) createOperatorRoles(roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope) error {
	installerRoleArn := scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN
	if installerRoleArn == "" {
		return fmt.Errorf("installer role is empty")
	}
	oidcConfigID := scope.RosaRoleConfig.Status.OIDCID
	if oidcConfigID == "" {
		return fmt.Errorf("OIDCID is required to create operator roles")
	}

	policies, err := r.Runtime.OCMClient.GetPolicies("OperatorRole")
	if err != nil {
		return err
	}

	config := roleConfig.Spec.OperatorRoleConfig
	version := roleConfig.Spec.AccountRoleConfig.Version
	hostedCp := true
	forcePolicyCreation := true

	operatorRoles, err := r.Runtime.AWSClient.ListOperatorRoles(version, "", config.Prefix)
	if err != nil {
		return err
	}

	for _, roles := range operatorRoles {
		for _, role := range roles {
			roleSuffix := strings.TrimPrefix(role.RoleName, config.Prefix)
			if roleSuffix == role.RoleName {
				continue
			}
			switch roleSuffix {
			case expinfrav1.IngressOperatorARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.IngressARN = role.RoleARN
			case expinfrav1.ImageRegistryARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.ImageRegistryARN = role.RoleARN
			case expinfrav1.StorageARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.StorageARN = role.RoleARN
			case expinfrav1.NetworkARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.NetworkARN = role.RoleARN
			case expinfrav1.KubeCloudControllerARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.KubeCloudControllerARN = role.RoleARN
			case expinfrav1.NodePoolManagementARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.NodePoolManagementARN = role.RoleARN
			case expinfrav1.ControlPlaneOperatorARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.ControlPlaneOperatorARN = role.RoleARN
			case expinfrav1.KMSProviderARNSuffix:
				scope.RosaRoleConfig.Status.OperatorRolesRef.KMSProviderARN = role.RoleARN
			}
		}
	}

	if !r.operatorRolesReady(&scope.RosaRoleConfig.Status.OperatorRolesRef) {
		// not all operator roles are set, operator roles are not ready yet.
		r.clearOperatorRolesRef(&scope.RosaRoleConfig.Status.OperatorRolesRef)
		err = operatorroles.CreateOperatorRoles(r.Runtime, ocm.Production, config.PermissionsBoundaryARN, interactive.ModeAuto, policies, version, config.SharedVPCConfig.IsSharedVPC(), config.Prefix, hostedCp, installerRoleArn, forcePolicyCreation,
			oidcConfigID, config.SharedVPCConfig.RouteRoleARN, ocm.DefaultChannelGroup, config.SharedVPCConfig.VPCEndpointRoleARN)
		return err
	}

	return nil
}

func (r *ROSARoleConfigReconciler) reconcileOIDCConfig(roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope) error {
	oidcID := ""
	if scope.RosaRoleConfig.Status.OIDCID != "" {
		oidcID = scope.RosaRoleConfig.Status.OIDCID
	} else if roleConfig.Spec.OperatorRoleConfig.OIDCID != "" {
		oidcID = roleConfig.Spec.OperatorRoleConfig.OIDCID
	}

	if oidcID != "" {
		oidcConfig, err := r.Runtime.OCMClient.GetOidcConfig(oidcID)
		if err != nil || oidcConfig == nil {
			return fmt.Errorf("failed to get OIDC config: %w", err)
		}
		scope.RosaRoleConfig.Status.OIDCID = oidcID
		return nil
	}

	return r.createOIDCConfig(scope)
}

func (r *ROSARoleConfigReconciler) createOIDCProvider(scope *scope.RosaRoleConfigScope) error {
	oidcID := scope.RosaRoleConfig.Status.OIDCID
	if oidcID == "" {
		return nil
	}

	oidcConfig, err := r.Runtime.OCMClient.GetOidcConfig(oidcID)
	if err != nil {
		return err
	}

	providers, err := r.Runtime.AWSClient.ListOidcProviders("", oidcConfig)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		if strings.Contains(provider.Arn, oidcID) {
			scope.RosaRoleConfig.Status.OIDCProviderARN = provider.Arn
			return nil
		}
	}

	return oidcprovider.CreateOIDCProvider(r.Runtime, oidcID, "", true)
}

func (r *ROSARoleConfigReconciler) createAccountRoles(roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope) error {
	config := roleConfig.Spec.AccountRoleConfig
	policies, err := r.Runtime.OCMClient.GetPolicies("AccountRole")
	if err != nil {
		return err
	}

	createRoles := true
	accountRoles, err := r.Runtime.AWSClient.ListAccountRoles(config.Version)
	if err != nil {
		// Let create account roles continue if no account roles are found
		if !strings.Contains(err.Error(), "no account roles found") {
			return err
		}
	}

	for _, role := range accountRoles {
		if role.RoleName == fmt.Sprintf("%s-HCP-ROSA-Installer-Role", config.Prefix) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN = role.RoleARN
		}
		if role.RoleName == fmt.Sprintf("%s-HCP-ROSA-Support-Role", config.Prefix) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.SupportRoleARN = role.RoleARN
		}
		if role.RoleName == fmt.Sprintf("%s-HCP-ROSA-Worker-Role", config.Prefix) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.WorkerRoleARN = role.RoleARN
		}
	}
	if createRoles {
		managedPolicies := true
		err := accountroles.CreateHCPRoles(r.Runtime, config.Prefix, managedPolicies, config.PermissionsBoundaryARN, ocm.Production, policies, config.Version, config.Path, config.SharedVPCConfig.IsSharedVPC(), config.SharedVPCConfig.RouteRoleARN, config.SharedVPCConfig.VPCEndpointRoleARN)
		return err
	}

	return nil
}

func (r *ROSARoleConfigReconciler) createOIDCConfig(scope *scope.RosaRoleConfigScope) error {
	// userPrefix, region are used only for unmanaged OIDC config
	oidcID, createErr := oidcconfig.CreateOIDCConfig(r.Runtime, true, "", "")
	if createErr != nil {
		return fmt.Errorf("failed to Create OIDC config: %w", createErr)
	}

	scope.RosaRoleConfig.Status.OIDCID = oidcID
	return createErr
}

func (r *ROSARoleConfigReconciler) deleteAccountRoles(scope *scope.RosaRoleConfigScope) error {
	roles := scope.RosaRoleConfig.Status.AccountRolesRef
	config := scope.RosaRoleConfig.Spec.AccountRoleConfig
	deleteHcpSharedVpcPolicies := config.SharedVPCConfig.VPCEndpointRoleARN != "" && config.SharedVPCConfig.RouteRoleARN != ""

	clusters, err := r.Runtime.OCMClient.GetAllClusters(r.Runtime.Creator)
	if err != nil {
		return err
	}

	if canDeleteRole(clusters, roles.InstallerRoleARN) {
		err = errors.Join(err, r.Runtime.AWSClient.DeleteAccountRole(strings.Split(roles.InstallerRoleARN, "/")[1], config.Prefix, true, deleteHcpSharedVpcPolicies))
	}
	if canDeleteRole(clusters, roles.WorkerRoleARN) {
		err = errors.Join(err, r.Runtime.AWSClient.DeleteAccountRole(strings.Split(roles.WorkerRoleARN, "/")[1], config.Prefix, true, deleteHcpSharedVpcPolicies))
	}
	if canDeleteRole(clusters, roles.SupportRoleARN) {
		err = errors.Join(err, r.Runtime.AWSClient.DeleteAccountRole(strings.Split(roles.SupportRoleARN, "/")[1], config.Prefix, true, deleteHcpSharedVpcPolicies))
	}
	return err
}

func (r *ROSARoleConfigReconciler) deleteOIDCProvider(oidcConfigID string) error {
	if oidcConfigID == "" {
		return nil
	}

	oidcConfig, err := r.Runtime.OCMClient.GetOidcConfig(oidcConfigID)
	if err != nil {
		return err
	}

	oidcEndpointURL := oidcConfig.IssuerUrl()
	parsedURI, _ := url.ParseRequestURI(oidcEndpointURL)
	if parsedURI.Scheme != helper.ProtocolHttps {
		return fmt.Errorf("expected OIDC endpoint URL '%s' to use an https:// scheme", oidcEndpointURL)
	}
	providerArn, err := r.Runtime.AWSClient.GetOpenIDConnectProviderByOidcEndpointUrl(oidcEndpointURL)
	if err != nil {
		return err
	}

	if providerArn == "" {
		return nil
	}
	hasClusterUsingOidcProvider, err := r.Runtime.OCMClient.HasAClusterUsingOidcProvider(oidcEndpointURL, r.Runtime.Creator.AccountID)
	if err != nil {
		return err
	}

	if hasClusterUsingOidcProvider {
		return fmt.Errorf("there are clusters using OIDC config '%s', can't delete the provider", oidcEndpointURL)
	}

	return r.Runtime.AWSClient.DeleteOpenIDConnectProvider(providerArn)
}

func (r *ROSARoleConfigReconciler) deleteOperatorRoles(prefix string) error {
	hasClusterUsingOperatorRolesPrefix, err := r.Runtime.OCMClient.HasAClusterUsingOperatorRolesPrefix(prefix)
	if err != nil {
		return err
	}
	if hasClusterUsingOperatorRolesPrefix {
		return fmt.Errorf("there are clusters using Operator Roles Prefix '%s', can't delete the IAM roles", prefix)
	}

	credRequests, err := r.Runtime.OCMClient.GetAllCredRequests()
	if err != nil {
		return err
	}

	foundOperatorRoles, err := r.Runtime.AWSClient.GetOperatorRolesFromAccountByPrefix(prefix, credRequests)
	if err != nil {
		return err
	}

	if len(foundOperatorRoles) == 0 {
		return nil
	}

	_, roleARN, err := r.Runtime.AWSClient.CheckRoleExists(foundOperatorRoles[0])
	if err != nil {
		return err
	}

	managedPolicies, err := r.Runtime.AWSClient.HasManagedPolicies(roleARN)
	if err != nil {
		return err
	}

	allSharedVpcPoliciesNotDeleted := make(map[string]bool)
	for _, role := range foundOperatorRoles {
		sharedVpcPoliciesNotDeleted, _ := r.Runtime.AWSClient.DeleteOperatorRole(role, managedPolicies, true)
		for key, value := range sharedVpcPoliciesNotDeleted {
			allSharedVpcPoliciesNotDeleted[key] = value
		}
	}

	for policyOutput, notDeleted := range allSharedVpcPoliciesNotDeleted {
		if notDeleted {
			return fmt.Errorf("unable to delete policy %s: Policy still attached to other resources", policyOutput)
		}
	}
	return nil
}

func (r *ROSARoleConfigReconciler) deleteOIDCConfig(oidcConfigID string) error {
	if oidcConfigID == "" {
		return nil
	}
	return r.Runtime.OCMClient.DeleteOidcConfig(oidcConfigID)
}

func canDeleteRole(clusters []*cmv1.Cluster, roleARN string) bool {
	if roleARN == "" {
		return false
	}
	for _, cluster := range clusters {
		if cluster.AWS().STS().RoleARN() == roleARN ||
			cluster.AWS().STS().SupportRoleARN() == roleARN ||
			cluster.AWS().STS().InstanceIAMRoles().MasterRoleARN() == roleARN ||
			cluster.AWS().STS().InstanceIAMRoles().WorkerRoleARN() == roleARN {
			return false
		}
	}
	return true
}

func (r ROSARoleConfigReconciler) rosaRolesConfigReady(scope *scope.RosaRoleConfigScope) bool {
	if scope.RosaRoleConfig.Status.OIDCID == "" ||
		scope.RosaRoleConfig.Status.OIDCProviderARN == "" ||
		scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN == "" ||
		scope.RosaRoleConfig.Status.AccountRolesRef.SupportRoleARN == "" ||
		scope.RosaRoleConfig.Status.AccountRolesRef.WorkerRoleARN == "" ||
		!r.operatorRolesReady(&scope.RosaRoleConfig.Status.OperatorRolesRef) {
		return false
	}
	return true
}

func (r ROSARoleConfigReconciler) operatorRolesReady(operatorRolesRef *v1beta2.AWSRolesRef) bool {
	if operatorRolesRef.ControlPlaneOperatorARN == "" ||
		operatorRolesRef.ImageRegistryARN == "" ||
		operatorRolesRef.IngressARN == "" ||
		operatorRolesRef.KMSProviderARN == "" ||
		operatorRolesRef.KubeCloudControllerARN == "" ||
		operatorRolesRef.NetworkARN == "" ||
		operatorRolesRef.NodePoolManagementARN == "" ||
		operatorRolesRef.StorageARN == "" {
		return false
	}
	return true
}

// GetOIDCIDFromOperatorRole extracts the OIDC UUID from the operator role policy document.
func (r *ROSARoleConfigReconciler) GetOIDCIDFromOperatorRole(scope *scope.RosaRoleConfigScope, roleDetails *iamv2.GetRoleOutput) (string, error) {
	decodedString, err := url.QueryUnescape(*roleDetails.Role.AssumeRolePolicyDocument)
	if err != nil {
		return "", err
	}

	var policyDoc struct {
		Statement []struct {
			Principal struct {
				Federated string `json:"Federated"`
			} `json:"Principal"`
			Condition map[string]map[string]string `json:"Condition"`
		} `json:"Statement"`
	}

	err = json.Unmarshal([]byte(decodedString), &policyDoc)
	if err != nil {
		return "", err
	}

	// Extract from the 'Federated' ARN
	if len(policyDoc.Statement) > 0 {
		federatedARN := policyDoc.Statement[0].Principal.Federated
		// The format is arn:aws:iam::ACCOUNT_ID:oidc-provider/OIDC_PROVIDER_URL
		// OIDC_PROVIDER_URL ends with /OIDCID
		parts := strings.Split(federatedARN, "/")
		if len(parts) > 1 {
			oidcUUID := parts[len(parts)-1]
			return oidcUUID, nil
		}
	}

	return "", fmt.Errorf("cant extract oidc uuid from the %s policy document", *roleDetails.Role.RoleName)
}

// setUpRuntime sets up the ROSA runtime if it doesn't exist.
func (r *ROSARoleConfigReconciler) setUpRuntime(ctx context.Context, scope *scope.RosaRoleConfigScope) error {
	if r.Runtime != nil {
		return nil
	}

	// Create OCM client
	ocm, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := rosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	r.Runtime = rosacli.NewRuntime()
	r.Runtime.OCMClient = ocmClient
	r.Runtime.Reporter = &rosa.Reporter{}
	r.Runtime.Logger = rosalogging.NewLogger()

	r.Runtime.AWSClient, err = aws.NewClient().Logger(r.Runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	r.Runtime.Creator, err = r.Runtime.AWSClient.GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get creator: %w", err)
	}

	return nil
}

// clearOperatorRolesRef clears all field values in the OperatorRolesRef by setting them to empty strings.
func (r ROSARoleConfigReconciler) clearOperatorRolesRef(operatorRolesRef *v1beta2.AWSRolesRef) {
	if operatorRolesRef == nil {
		return
	}

	operatorRolesRef.IngressARN = ""
	operatorRolesRef.ImageRegistryARN = ""
	operatorRolesRef.StorageARN = ""
	operatorRolesRef.NetworkARN = ""
	operatorRolesRef.KubeCloudControllerARN = ""
	operatorRolesRef.NodePoolManagementARN = ""
	operatorRolesRef.ControlPlaneOperatorARN = ""
	operatorRolesRef.KMSProviderARN = ""
}
