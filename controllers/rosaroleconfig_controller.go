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
	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	"sigs.k8s.io/cluster-api/util"
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

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs/status,verbs=get;update;patch

func (r *ROSARoleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	fmt.Println()
	fmt.Println("====================================================================================ROSARoleConfigReconciler======================================================")
	log := logger.FromContext(ctx)
	// runtime := rosacli.NewRuntime()

	// TODO: Get the ROLE CONFIG SCOPE
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}

	// Set CP for easier testing
	cpName := types.NamespacedName{
		Name:      "rosa-rknaur2241-control-plane",
		Namespace: "default",
	}
	if err := r.Client.Get(ctx, cpName, rosaControlPlane); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaControlPlane.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner Cluster from the API Server")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, rosaControlPlane); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	rosaControlPlaneKind := "ROSAControlPlane"
	rosaScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   rosaControlPlane,
		ControllerName: strings.ToLower(rosaControlPlaneKind),
		Endpoints:      []scope.ServiceEndpoint{},
		Logger:         log,
		NewStsClient:   r.NewStsClient,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	roleConfig := &infrav1.ROSARoleConfig{}
	if err := r.Get(ctx, req.NamespacedName, roleConfig); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	roleConfigScope, err := scope.NewRosaRoleConfigScope(scope.RosaRoleConfigScopeParams{
		Client:         r.Client,
		RosaRoleConfig: roleConfig,
		ControllerName: "rosanetwork",
		Endpoints:      r.Endpoints,
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosanetwork scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := rosaScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
		if err := roleConfigScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	ocm, err := r.NewOCMClient(ctx, roleConfigScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	ocmClient, err := rosa.ConvertToRosaOcmClient(ocm)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	fmt.Println("ROLE CONFIG:", roleConfig.Spec)
	// roleConfig.Spec.AccountRoleConfig.Prefix = "rk2"

	roleConfigScope.RosaRoleConfig.Status.OIDCID, err = r.createOIDCConfig(ctx, roleConfig, roleConfigScope, ocmClient)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to Create OIDC provider: %w", err)
	}

	// errAccountRoles := r.createAccountRoles2(ctx, roleConfig, rosaScope, ocmClient)
	// if errAccountRoles != nil {
	// 	return ctrl.Result{}, fmt.Errorf("failed to Create AccountRoles: %w", errAccountRoles)
	// }

	// errOperatorRoles := r.createOperatorRoles(ctx, roleConfig, rosaScope, ocmClient)
	// if errOperatorRoles != nil {
	// 	return ctrl.Result{}, fmt.Errorf("failed to Create OperatorRoles: %w", errOperatorRoles)
	// }

	// errOIDCProvider := r.createOIDCProvider(ctx, roleConfig, rosaScope, ocmClient)
	// if errOIDCProvider != nil {
	// 	return ctrl.Result{}, fmt.Errorf("failed to Create OIDC provider: %w", errOIDCProvider)
	// }

	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA, END OF RECONCILE", err, roleConfigScope.RosaRoleConfig.Status.OIDCID)

	return ctrl.Result{}, nil
}

func (r *ROSARoleConfigReconciler) createOperatorRoles(ctx context.Context, roleConfig *infrav1.ROSARoleConfig, scope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient) error {
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

	permissionBoundary := ""
	// path := ""
	// managedPolicies := true
	version := "4.18"
	isSharedVpc := false
	hostedCp := true
	installerRoleArn := "arn:aws:iam::471112697682:role/rk-HCP-ROSA-Installer-Role"
	forcePolicyCreation := true
	oidcConfigId := "2gf8rmqmkhic4egagjg1hj13relln2ds"
	sharedVpcRoleArn := ""
	vpcEndpointRoleArn := ""

	// OPERATOR ROLES
	//  Knedlik(r *rosa.Runtime, env string, permissionsBoundary string, mode string, policies map[string]*cmv1.AWSSTSPolicy, defaultPolicyVersion string, isSharedVpc bool,
	// prefix string, hostedCp bool, installerRoleArn string, forcePolicyCreation bool, oidcConfigId string, sharedVpcRoleArn string, channelGroup string, vpcEndpointRoleArn string)
	err = operatorroles.Knedlik(runtime, ocm.Production, permissionBoundary, interactive.ModeAuto, policies, version, isSharedVpc, roleConfig.Spec.OperatorRoleConfig.Prefix, hostedCp, installerRoleArn, forcePolicyCreation,
		oidcConfigId, sharedVpcRoleArn, ocm.DefaultChannelGroup, vpcEndpointRoleArn)
	fmt.Println("OPERATOR ROLES err: ", err)

	return err
}

func (r *ROSARoleConfigReconciler) createOIDCConfig(ctx context.Context, roleConfig *infrav1.ROSARoleConfig, scope *scope.RosaRoleConfigScope, ocmClient *ocm.Client) (string, error) {
	fmt.Println("OIDC CONFIG")
	return "2itmoq75h48lb6ssar3ijscm87uu557s", nil

	oicdConfig := roleConfig.Spec.OIDCConfig
	runtime := rosacli.NewRuntime()
	var err error
	// ocm, err := r.NewOCMClient(ctx, scope)
	// if err != nil {
	// 	return err
	// }

	// runtime.OCMClient, err = rosa.ConvertToRosaOcmClient(ocm)
	// if err != nil {
	// 	return err
	// }

	runtime.Reporter = (&NilReporter{})
	runtime.OCMClient = ocmClient
	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return "", fmt.Errorf("failed to create aws client: %w", err)
	}

	runtime.Creator, err = runtime.AWSClient.GetCreator()
	if err != nil {
		return "", err
	}

	// userPrefix, region are used only for unmanaged OIDC config
	configId, createErr := oidcconfig.CreateOIDCConfig(runtime, oicdConfig.ManagedOIDC, oicdConfig.Prefix, oicdConfig.Region)
	// TODO save configId to status
	fmt.Println("CONFIG ID:", configId)
	return configId, createErr
}

func (r *ROSARoleConfigReconciler) createOIDCProvider(ctx context.Context, roleConfig *infrav1.ROSARoleConfig, scope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient) error {
	fmt.Println("OIDC PROVIDER")

	runtime := rosacli.NewRuntime()

	var err error
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

	return oidcprovider.CreateOIDCProvider(runtime, "2hov94ls1872f8cfv31ihmstsl2d5n8j", string(scope.Cluster.UID), true)
}

func (r *ROSARoleConfigReconciler) createAccountRoles2(ctx context.Context, roleConfig *infrav1.ROSARoleConfig, scope *scope.ROSAControlPlaneScope, ocmClient rosa.OCMClient) error {
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
	fmt.Println("POLICIES", policies)
	permissionBoundary := ""
	path := ""
	managedPolicies := true
	version := "4.18"
	isSharedVpc := false
	// TODO: set runtime.Creator.Partition here
	// assumeRolePolicy := getAssumeRolePolicy(runtime.Creator.Partition, file, a)
	// assumeRolePolicy := getAssumeRolePolicy("aws", file, a)

	runtime.Reporter = (&NilReporter{})

	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	accountroles.CreateHCPRoles(runtime, roleConfig.Spec.AccountRoleConfig.Prefix, managedPolicies, permissionBoundary, ocm.Production, policies, version, path, isSharedVpc)
	return nil
}

func (r *ROSARoleConfigReconciler) createAccountRoles(ctx context.Context, roleConfig *infrav1.ROSARoleConfig, ocmClient rosa.OCMClient) error {
	runtime := rosacli.NewRuntime()
	policies, err := ocmClient.GetPolicies("AccountRole")
	if err != nil {
		return fmt.Errorf("failed to get policies for AccountRole: %w", err)
	}
	fmt.Println("POLICIES", policies)
	tagList := map[string]string{
		"rosa_role_prefix":       roleConfig.Spec.AccountRoleConfig.Prefix,
		"red-hat-managed":        "true",
		"rosa_hcp_policies":      "true",
		"rosa_role_type":         "installer",
		"rosa_managed_policies":  "true",
		"rosa_openshift_version": "4.18",
	}
	permissionBoundary := ""
	path := ""
	managedPolicies := true
	version := "4.18"
	// accRoleName := common.GetRoleName(roleConfig.Spec.AccountRoleConfig.Prefix, "HCP-ROSA-Installer")
	// a := &accountRolesCreationInput{
	// 	policies: policies,
	// 	env:      ocm.Production,
	// }

	// TODO: set runtime.Creator.Partition here
	// assumeRolePolicy := getAssumeRolePolicy(runtime.Creator.Partition, file, a)
	// assumeRolePolicy := getAssumeRolePolicy("aws", file, a)

	runtime.Reporter = (&NilReporter{})

	runtime.Logger = rosalogging.NewLogger()
	runtime.AWSClient, err = aws.NewClient().Logger(runtime.Logger).Build()
	if err != nil {
		return fmt.Errorf("failed to create aws client: %w", err)
	}

	rolesARNs := []string{}
	// TODO: make it nicer
	files := []string{"installer", "instance_worker", "support"}
	roleNames := []string{"HCP-ROSA-Installer", "HCP-ROSA-Worker", "HCP-ROSA-Support"}
	for i, roleName := range roleNames {
		fmt.Println("GETTING ROLE NAME", roleConfig.Spec.AccountRoleConfig.Prefix, roleName)
		accRoleName := common.GetRoleName(roleConfig.Spec.AccountRoleConfig.Prefix, roleName)
		fmt.Println("ACC NAME: ", accRoleName)

		a := &accountRolesCreationInput{
			policies: policies,
			env:      ocm.Production,
		}
		fmt.Println("GETTING POLICY", files[i])
		assumeRolePolicy := getAssumeRolePolicy("aws", files[i], a)
		fmt.Println("POLICY: ", assumeRolePolicy)
		tagList["rosa_role_type"] = files[i]
		roleARN, err := runtime.AWSClient.EnsureRole(runtime.Reporter, accRoleName, assumeRolePolicy, permissionBoundary, version, tagList, path, managedPolicies)
		if err != nil {
			return err
		}
		rolesARNs = append(rolesARNs, roleARN)

	}

	// TODO save roles to status
	fmt.Println(rolesARNs)
	return nil
}

func (r *ROSARoleConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClientWithoutControlPlane
	r.NewStsClient = scope.NewSTSClient

	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ROSARoleConfig{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Complete(r)
}

// func (mp *managedPoliciesCreator) createRoles(r *rosa.Runtime, input *accountRolesCreationInput) error {
// 	r.Reporter.Infof("Creating classic account roles using '%s'", r.Creator.ARN)

// 	for file, role := range aws.AccountRoles {
// 		accRoleName := common.GetRoleName(input.prefix, role.Name)
// 		assumeRolePolicy := getAssumeRolePolicy(r.Creator.Partition, file, input)

// 		r.Reporter.Debugf("Creating role '%s'", accRoleName)
// 		tagsList := mp.getRoleTags(file, input)
// 		r.Reporter.Debugf("start to EnsureRole")
// 		roleARN, err := r.AWSClient.EnsureRole(r.Reporter, accRoleName, assumeRolePolicy, input.permissionsBoundary,
// 			input.defaultPolicyVersion, tagsList, input.path, true)
// 		if err != nil {
// 			return err
// 		}
// 		r.Reporter.Infof("Created role '%s' with ARN '%s'", accRoleName, roleARN)

// 		err = attachManagedPolicies(r, input, file, accRoleName)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }.

type accountRolesCreationInput struct {
	prefix               string
	permissionsBoundary  string
	accountID            string
	env                  string
	policies             map[string]*cmv1.AWSSTSPolicy
	defaultPolicyVersion string
	path                 string
	isSharedVpc          bool
}

func getAssumeRolePolicy(partition string, file string, input *accountRolesCreationInput) string {
	filename := fmt.Sprintf("sts_%s_trust_policy", file)
	policyDetail := aws.GetPolicyDetails(input.policies, filename)
	return aws.InterpolatePolicyDocument(partition, policyDetail, map[string]string{
		"partition":      partition,
		"aws_account_id": aws.GetJumpAccount(input.env),
	})
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

// aws iam create-role \
// 	--assume-role-policy-document file://sts_installer_trust_policy.json \
// 	--role-name rk-HCP-ROSA-Installer-Role \
// 	--tags Key=rosa_role_prefix,Value=rk Key=red-hat-managed,Value=true Key=rosa_hcp_policies,Value=true Key=rosa_role_type,Value=installer Key=rosa_managed_policies,Value=true Key=rosa_openshift_version,Value=4.18

// aws iam attach-role-policy \
// 	--policy-arn arn:aws:iam::aws:policy/service-role/ROSAInstallerPolicy \
// 	--role-name rk-HCP-ROSA-Installer-Role

// aws iam create-role \
// 	--assume-role-policy-document file://sts_support_trust_policy.json \
// 	--role-name rk-HCP-ROSA-Support-Role \
// 	--tags Key=rosa_role_prefix,Value=rk Key=red-hat-managed,Value=true Key=rosa_hcp_policies,Value=true Key=rosa_role_type,Value=support Key=rosa_managed_policies,Value=true Key=rosa_openshift_version,Value=4.18

// aws iam attach-role-policy \
// 	--policy-arn arn:aws:iam::aws:policy/service-role/ROSASRESupportPolicy \
// 	--role-name rk-HCP-ROSA-Support-Role

// aws iam create-role \
// 	--assume-role-policy-document file://sts_instance_worker_trust_policy.json \
// 	--role-name rk-HCP-ROSA-Worker-Role \
// 	--tags Key=rosa_role_prefix,Value=rk Key=red-hat-managed,Value=true Key=rosa_hcp_policies,Value=true Key=rosa_managed_policies,Value=true Key=rosa_openshift_version,Value=4.18 Key=rosa_role_type,Value=instance_worker

// aws iam attach-role-policy \
// 	--policy-arn arn:aws:iam::aws:policy/service-role/ROSAWorkerInstancePolicy \
// 	--role-name rk-HCP-ROSA-Worker-Role

// INPUT: {rk02  471112697682 staging map[sts_hcp_installer_permission_policy:0x14001586280 sts_hcp_instance_worker_permission_policy:0x140015862d0 sts_hcp_support_permission_policy:0x14001586320 sts_installer_core_permission_policy:0x14001586370
// sts_installer_permission_policy:0x140015863c0 sts_installer_privatelink_permission_policy:0x14001586410 sts_installer_trust_policy:0x14001586460 sts_installer_vpc_permission_policy:0x140015864b0 sts_instance_controlplane_permission_policy:0x14001586500
// sts_instance_controlplane_trust_policy:0x14001586550 sts_instance_worker_permission_policy:0x140015865a0 sts_instance_worker_trust_policy:0x140015865f0 sts_support_permission_policy:0x14001586640 sts_support_trust_policy:0x14001586690] 4.18  false}
// I: Creating hosted CP account roles using 'arn:aws:iam::471112697682:role/471112697682-poweruser'
// CALLING GET ROLE: rk02 HCP-ROSA-Support
// CALLING GET ASSUME ROLE POLICY: rk02 HCP-ROSA-Support map[sts_hcp_installer_permission_policy:0x14001586280 sts_hcp_instance_worker_permission_policy:0x140015862d0 sts_hcp_support_permission_policy:0x14001586320 sts_installer_core_permission_policy:0x14001586370 sts_installer_permission_policy:0x140015863c0 sts_installer_privatelink_permission_policy:0x14001586410 sts_installer_trust_policy:0x14001586460 sts_installer_vpc_permission_policy:0x140015864b0 sts_instance_controlplane_permission_policy:0x14001586500 sts_instance_controlplane_trust_policy:0x14001586550 sts_instance_worker_permission_policy:0x140015865a0 sts_instance_worker_trust_policy:0x140015865f0 sts_support_permission_policy:0x14001586640 sts_support_trust_policy:0x14001586690]
// POLICY DETAIL: {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Technical-Support-13829321"]}}]}
// CALLING ENSURE ROLE: rk02-HCP-ROSA-Support-Role rk02-HCP-ROSA-Support-Role {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Technical-Support-13829321"]}}]}  4.18 map[red-hat-managed:true rosa_hcp_policies:true rosa_managed_policies:true rosa_openshift_version:4.18 rosa_role_prefix:rk02 rosa_role_type:support]
// I: Attached trust policy to role 'rk02-HCP-ROSA-Support-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Support-Role)': {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Technical-Support-13829321"]}}]}
// I: Created role 'rk02-HCP-ROSA-Support-Role' with ARN 'arn:aws:iam::471112697682:role/rk02-HCP-ROSA-Support-Role'
// I: Attached policy 'ROSASRESupportPolicy(https://docs.aws.amazon.com/aws-managed-policy/latest/reference/ROSASRESupportPolicy)' to role 'rk02-HCP-ROSA-Support-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Support-Role)'

// CALLING GET ROLE: rk02 HCP-ROSA-Worker
// CALLING GET ASSUME ROLE POLICY: rk02 HCP-ROSA-Worker map[sts_hcp_installer_permission_policy:0x14001586280 sts_hcp_instance_worker_permission_policy:0x140015862d0 sts_hcp_support_permission_policy:0x14001586320 sts_installer_core_permission_policy:0x14001586370 sts_installer_permission_policy:0x140015863c0 sts_installer_privatelink_permission_policy:0x14001586410 sts_installer_trust_policy:0x14001586460 sts_installer_vpc_permission_policy:0x140015864b0 sts_instance_controlplane_permission_policy:0x14001586500 sts_instance_controlplane_trust_policy:0x14001586550 sts_instance_worker_permission_policy:0x140015865a0 sts_instance_worker_trust_policy:0x140015865f0 sts_support_permission_policy:0x14001586640 sts_support_trust_policy:0x14001586690]
// POLICY DETAIL: {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"Service": ["ec2.amazonaws.com"]}}]}
// CALLING ENSURE ROLE: rk02-HCP-ROSA-Worker-Role rk02-HCP-ROSA-Worker-Role {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"Service": ["ec2.amazonaws.com"]}}]}  4.18 map[red-hat-managed:true rosa_hcp_policies:true rosa_managed_policies:true rosa_openshift_version:4.18 rosa_role_prefix:rk02 rosa_role_type:instance_worker]
// I: Attached trust policy to role 'rk02-HCP-ROSA-Worker-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Worker-Role)': {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"Service": ["ec2.amazonaws.com"]}}]}
// I: Created role 'rk02-HCP-ROSA-Worker-Role' with ARN 'arn:aws:iam::471112697682:role/rk02-HCP-ROSA-Worker-Role'
// I: Attached policy 'ROSAWorkerInstancePolicy(https://docs.aws.amazon.com/aws-managed-policy/latest/reference/ROSAWorkerInstancePolicy)' to role 'rk02-HCP-ROSA-Worker-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Worker-Role)'

// I: Skip attaching policy 'sts_hcp_ec2_registry_permission_policy' to role 'rk02-HCP-ROSA-Worker-Role' (zero egress feature toggle is not enabled)
// CALLING GET ROLE: rk02 HCP-ROSA-Installer
// CALLING GET ASSUME ROLE POLICY: rk02 HCP-ROSA-Installer map[sts_hcp_installer_permission_policy:0x14001586280 sts_hcp_instance_worker_permission_policy:0x140015862d0 sts_hcp_support_permission_policy:0x14001586320 sts_installer_core_permission_policy:0x14001586370 sts_installer_permission_policy:0x140015863c0 sts_installer_privatelink_permission_policy:0x14001586410 sts_installer_trust_policy:0x14001586460 sts_installer_vpc_permission_policy:0x140015864b0 sts_instance_controlplane_permission_policy:0x14001586500 sts_instance_controlplane_trust_policy:0x14001586550 sts_instance_worker_permission_policy:0x140015865a0 sts_instance_worker_trust_policy:0x140015865f0 sts_support_permission_policy:0x14001586640 sts_support_trust_policy:0x14001586690]
// POLICY DETAIL: {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Managed-OpenShift-Installer"]}}]}
// CALLING ENSURE ROLE: rk02-HCP-ROSA-Installer-Role rk02-HCP-ROSA-Installer-Role {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Managed-OpenShift-Installer"]}}]}  4.18 map[red-hat-managed:true rosa_hcp_policies:true rosa_managed_policies:true rosa_openshift_version:4.18 rosa_role_prefix:rk02 rosa_role_type:installer]
// I: Attached trust policy to role 'rk02-HCP-ROSA-Installer-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Installer-Role)': {"Version": "2012-10-17", "Statement": [{"Action": ["sts:AssumeRole"], "Effect": "Allow", "Principal": {"AWS": ["arn:aws:iam::644306948063:role/RH-Managed-OpenShift-Installer"]}}]}
// I: Created role 'rk02-HCP-ROSA-Installer-Role' with ARN 'arn:aws:iam::471112697682:role/rk02-HCP-ROSA-Installer-Role'
// I: Attached policy 'ROSAInstallerPolicy(https://docs.aws.amazon.com/aws-managed-policy/latest/reference/ROSAInstallerPolicy)' to role 'rk02-HCP-ROSA-Installer-Role(https://console.aws.amazon.com/iam/home?#/roles/rk02-HCP-ROSA-Installer-Role)'
