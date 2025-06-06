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
	"strings"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/go-logr/logr"
	"github.com/openshift/rosa/pkg/aws"
	interactive "github.com/openshift/rosa/pkg/interactive"
	rosalogging "github.com/openshift/rosa/pkg/logging"
	"github.com/openshift/rosa/pkg/ocm"
	rosacli "github.com/openshift/rosa/pkg/rosa"

	accountroles "github.com/openshift/rosa/cmd/create/accountroles"
	oidcconfig "github.com/openshift/rosa/cmd/create/oidcconfig"
	oidcprovider "github.com/openshift/rosa/cmd/create/oidcprovider"

	operatorroles "github.com/openshift/rosa/cmd/create/operatorroles"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSARoleConfigReconciler reconciles a ROSARoleConfig object.
type ROSARoleConfigReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	Endpoints        []scope.ServiceEndpoint
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsiface.STSAPI
	NewOCMClient     func(ctx context.Context, scope rosa.OCMSecretsRetriever) (rosa.OCMClient, error)
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
func (r *ROSARoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	fmt.Println("====================================================================================ROSARoleConfigReconciler======================================================")
	log := logger.FromContext(ctx)
	// runtime := rosacli.NewRuntime()

	roleConfig := &expinfrav1.ROSARoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, roleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	log = log.WithValues("cluster", klog.KObj(roleConfig))
	scope, err := scope.NewRosaRoleConfigScope(scope.RosaRoleConfigScopeParams{
		Client:         r.Client,
		RosaRoleConfig: roleConfig,
		ControllerName: "rosaroleconfig",
		Endpoints:      r.Endpoints,
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosaroleconfig scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := scope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	ocm, err := r.NewOCMClient(ctx, scope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := rosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	if scope.RosaRoleConfig.Status.OIDCID == "" {
		err = r.createOIDCConfig(ctx, roleConfig, scope, ocmClient)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to OICD Config: %w", err)
		}
	}

	errAccountRoles := r.createAccountRoles(ctx, roleConfig, scope, ocmClient)
	if errAccountRoles != nil {
		return ctrl.Result{}, fmt.Errorf("failed to Create AccountRoles: %w", errAccountRoles)
	}

	errOperatorRoles := r.createOperatorRoles(ctx, roleConfig, scope, ocmClient)
	if errOperatorRoles != nil {
		return ctrl.Result{}, fmt.Errorf("failed to Create OperatorRoles: %w", errOperatorRoles)
	}

	errOIDCProvider := r.createOIDCProvider(ctx, scope, ocmClient)
	if errOIDCProvider != nil {
		return ctrl.Result{}, fmt.Errorf("failed to Create OIDC provider: %w", errOIDCProvider)
	}

	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA, END OF RECONCILE", err, scope.RosaRoleConfig.Status.OIDCID)

	return ctrl.Result{}, nil
}

func (r *ROSARoleConfigReconciler) createOperatorRoles(ctx context.Context, roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope, ocmClient rosa.OCMClient) error {
	fmt.Println("createOperatorRoles", scope.RosaRoleConfig.Status)

	installerRoleArn := scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN
	if installerRoleArn == "" {
		return fmt.Errorf("installer role is empty")
	}
	oidcConfigId := scope.RosaRoleConfig.Status.OIDCID
	if oidcConfigId == "" {
		return fmt.Errorf("OIDCID is empty")
	}

	runtime := rosacli.NewRuntime()
	policies, err := ocmClient.GetPolicies("OperatorRole")
	if err != nil {
		return err
	}
	runtime.OCMClient, err = rosa.NewOCMClient(ctx, scope)
	if err != nil {
		return err
	}

	runtime.Reporter = (&NilReporter{})
	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	runtime.Creator, err = runtime.AWSClient.GetCreator()
	if err != nil {
		return err
	}

	config := roleConfig.Spec.OperatorRoleConfig
	version := roleConfig.Spec.AccountRoleConfig.Version
	hostedCp := true
	forcePolicyCreation := true
	isSharedVpc := config.SharedVPCConfig.VPCEndpointRoleARN != "" && config.SharedVPCConfig.RouteRoleARN != ""

	operatorRoles, err := runtime.AWSClient.ListOperatorRoles(version, "", config.Prefix)
	if err != nil {
		return err
	}

	if len(operatorRoles) > 0 {
		for key, roles := range operatorRoles {
			fmt.Println(key, len(roles))
			for _, role := range roles {
				fmt.Println(role.RoleName)
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-openshift-ingress-operator-cloud-credentials", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.IngressARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-openshift-image-registry-installer-cloud-credentials", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.ImageRegistryARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-openshift-cluster-csi-drivers-ebs-cloud-credentials", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.StorageARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-openshift-cloud-network-config-controller-cloud-credentials", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.NetworkARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-kube-system-kube-controller-manager", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.KubeCloudControllerARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-kube-system-capa-controller-manager", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.NodePoolManagementARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-kube-system-control-plane-operator", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.ControlPlaneOperatorARN = role.RoleARN
				}
				if strings.Contains(role.RoleName, fmt.Sprintf("%s-kube-system-kms-provider", config.Prefix)) {
					scope.RosaRoleConfig.Status.OperatorRolesRef.KMSProviderARN = role.RoleARN
				}
			}

		}
	} else {
		err = operatorroles.CreateOperatorRoles(runtime, ocm.Production, config.PermissionsBoundaryARN, interactive.ModeAuto, policies, version, isSharedVpc, config.Prefix, hostedCp, installerRoleArn, forcePolicyCreation,
			oidcConfigId, config.SharedVPCConfig.RouteRoleARN, ocm.DefaultChannelGroup, config.SharedVPCConfig.VPCEndpointRoleARN)
		return err
	}

	return nil
}

func (r *ROSARoleConfigReconciler) createOIDCConfig(ctx context.Context, roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope, ocmClient *ocm.Client) error {
	fmt.Println("OIDC CONFIG")

	oicdConfig := roleConfig.Spec.OIDCConfig
	runtime := rosacli.NewRuntime()
	var err error
	runtime.Reporter = (&NilReporter{})
	runtime.OCMClient = ocmClient
	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	runtime.Creator, err = runtime.AWSClient.GetCreator()
	if err != nil {
		return err
	}

	// userPrefix, region are used only for unmanaged OIDC config
	configId, createErr := oidcconfig.CreateOIDCConfig(runtime, oicdConfig.ManagedOIDC, oicdConfig.Prefix, oicdConfig.Region)
	if createErr != nil {
		return fmt.Errorf("failed to Create OIDC config: %w", err)
	}

	scope.RosaRoleConfig.Status.OIDCID = configId
	return createErr
}

func (r *ROSARoleConfigReconciler) createOIDCProvider(ctx context.Context, scope *scope.RosaRoleConfigScope, ocmClient *ocm.Client) error {
	if scope.RosaRoleConfig.Status.OIDCProviderARN != "" {
		return nil
	}

	var err error
	oidcId := scope.RosaRoleConfig.Status.OIDCID
	runtime := rosacli.NewRuntime()
	runtime.OCMClient = ocmClient
	runtime.Reporter = (&NilReporter{})

	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	oidcConfig, err := runtime.OCMClient.GetOidcConfig(oidcId)
	if err != nil {
		return err
	}

	providers, err := runtime.AWSClient.ListOidcProviders("", oidcConfig)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		fmt.Println("PROVIDER", provider.Arn, provider.ClusterId)
		if strings.Contains(provider.Arn, oidcId) {
			scope.RosaRoleConfig.Status.OIDCProviderARN = provider.Arn
			return nil
		}
	}

	runtime.Creator, err = runtime.AWSClient.GetCreator()
	if err != nil {
		return err
	}

	fmt.Println("CREATING PROVIDER WITH OIDC ID", oidcId)
	// scope.RosaRoleConfig.Status.OIDCProviderARN
	return oidcprovider.CreateOIDCProvider(runtime, oidcId, "", true)
}

func (r *ROSARoleConfigReconciler) createAccountRoles(ctx context.Context, roleConfig *expinfrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope, ocmClient rosa.OCMClient) error {
	fmt.Println("createAccountRoles")
	config := roleConfig.Spec.AccountRoleConfig
	runtime := rosacli.NewRuntime()
	policies, err := ocmClient.GetPolicies("AccountRole")
	if err != nil {
		return err
	}
	runtime.OCMClient, err = rosa.NewOCMClient(ctx, scope)
	if err != nil {
		return err
	}

	runtime.Reporter = (&NilReporter{})
	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	runtime.Creator, err = runtime.AWSClient.GetCreator()
	if err != nil {
		return err
	}

	createRoles := true
	accountRoles, err := runtime.AWSClient.ListAccountRoles(config.Version)
	if err != nil {
		// Let create account roles continue if no account roles are found
		if !strings.Contains(err.Error(), "no account roles found") {
			return err
		}
	}

	for _, role := range accountRoles {
		if strings.Contains(role.RoleName, fmt.Sprintf("%s-HCP-ROSA-Installer", config.Prefix)) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.InstallerRoleARN = role.RoleARN
		}
		if strings.Contains(role.RoleName, fmt.Sprintf("%s-HCP-ROSA-Support", config.Prefix)) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.SupportRoleARN = role.RoleARN
		}
		if strings.Contains(role.RoleName, fmt.Sprintf("%s-HCP-ROSA-Worker", config.Prefix)) {
			createRoles = false
			scope.RosaRoleConfig.Status.AccountRolesRef.WorkerRoleARN = role.RoleARN
		}

	}
	if createRoles {
		runtime.Reporter = (&NilReporter{})
		runtime.Logger = rosalogging.NewLogger()
		runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
		if err != nil {
			return fmt.Errorf("failed to create aws client: %w", err)
		}

		managedPolicies := true
		isSharedVpc := config.SharedVPCConfig.VPCEndpointRoleARN != "" && config.SharedVPCConfig.RouteRoleARN != ""
		accountroles.CreateHCPRoles(runtime, config.Prefix, managedPolicies, config.PermissionsBoundaryARN, ocm.Production, policies, config.Version, config.Path, isSharedVpc, config.SharedVPCConfig.RouteRoleARN, config.SharedVPCConfig.VPCEndpointRoleARN)
	}

	return nil
}

type NilReporter struct {
}

// Debugf prints a debug message with the given format and arguments.
func (r *NilReporter) Debugf(format string, args ...interface{}) {
}

// Infof prints an informative message with the given format and arguments.
func (r *NilReporter) Infof(format string, args ...interface{}) {

}

// Warnf prints an warning message with the given format and arguments.
func (r *NilReporter) Warnf(format string, args ...interface{}) {
}

// Errorf prints an error message with the given format and arguments. It also return an error
// containing the same information, which will be usually discarded, except when the caller needs to
// report the error and also return it.
func (r *NilReporter) Errorf(format string, args ...interface{}) error {
	return nil
}
func (r *NilReporter) IsTerminal() bool {
	return true
}
