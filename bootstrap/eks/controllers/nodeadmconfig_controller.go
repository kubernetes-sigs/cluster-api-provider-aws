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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/internal/userdata"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	kubeconfigutil "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

const (
	// ssmActivationFinalizer is the finalizer for SSM activations created by the controller.
	ssmActivationFinalizer = "nodeadmconfig.bootstrap.cluster.x-k8s.io/ssm-activation"

	// ssmActivationSecretSuffix is the suffix for the secret storing SSM activation credentials.
	ssmActivationSecretSuffix = "-ssm-activation"

	// SSM activation secret keys.
	ssmActivationIDKey   = "activationId"
	ssmActivationCodeKey = "activationCode"
)

// SSMServiceFactory is a function that creates an SSM service from a cluster scope.
type SSMServiceFactory func(cloud.ClusterScoper) *ssm.Service

// NodeadmConfigReconciler reconciles a NodeadmConfig object.
type NodeadmConfigReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	WatchFilterValue string

	// SSMServiceFactory creates SSM services. Used for test injection.
	SSMServiceFactory SSMServiceFactory
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=nodeadmconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=nodeadmconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machinepools;clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;

func (r *NodeadmConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := logger.FromContext(ctx)

	// get NodeadmConfig
	config := &eksbootstrapv1.NodeadmConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}
	log = log.WithValues("NodeadmConfig", config.GetName())

	// check owner references and look up owning Machine object
	configOwner, err := bsutil.GetTypedConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// no error here, requeue until we find an owner
		log.Debug("NodeadmConfig failed to look up owner reference, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}
	if err != nil {
		log.Error(err, "NodeadmConfig failed to get owner")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		// no error, requeue until we find an owner
		log.Debug("NodeadmConfig has no owner reference set, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	log = log.WithValues(configOwner.GetKind(), configOwner.GetName())

	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Is(err, util.ErrNoCluster) {
			log.Info("NodeadmConfig does not belong to a cluster yet, re-queuing until it's part of a cluster")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		if apierrors.IsNotFound(err) {
			log.Info("Cluster does not exist yet, re-queueing until it is created")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		log.Error(err, "Could not get cluster with metadata")
		return ctrl.Result{}, err
	}
	log = log.WithValues("cluster", klog.KObj(cluster))

	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, config); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	patchHelper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// set up defer block for updating config
	defer func() {
		conditions := []clusterv1beta1.ConditionType{
			eksbootstrapv1.DataSecretAvailableCondition,
		}
		// Include SSM condition in summary for hybrid nodes
		if config.Spec.Hybrid != nil && config.Spec.Hybrid.SSM != nil && config.Spec.Hybrid.SSM.ActivationConfig != nil {
			conditions = append(conditions, eksbootstrapv1.SSMActivationReadyCondition)
		}
		v1beta1conditions.SetSummary(config,
			v1beta1conditions.WithConditions(conditions...),
			v1beta1conditions.WithStepCounter(),
		)

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

	// Handle deletion for hybrid nodes with auto-created SSM activations
	if !config.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, cluster, config)
	}

	// Add finalizer for hybrid nodes with auto-created activations
	if config.Spec.Hybrid != nil && config.Spec.Hybrid.SSM != nil && config.Spec.Hybrid.SSM.ActivationConfig != nil {
		if !controllerutil.ContainsFinalizer(config, ssmActivationFinalizer) {
			controllerutil.AddFinalizer(config, ssmActivationFinalizer)
		}
	}

	// Route to appropriate reconciliation based on hybrid mode
	if config.Spec.Hybrid != nil {
		return r.reconcileHybridNode(ctx, cluster, config, configOwner)
	}

	return r.joinWorker(ctx, cluster, config, configOwner)
}

func (r *NodeadmConfigReconciler) joinWorker(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, configOwner *bsutil.ConfigOwner) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// only need to reconcile the secret for Machine kinds once, but MachinePools need updates for new launch templates
	if config.Status.DataSecretName != nil && configOwner.GetKind() == "Machine" {
		secretKey := client.ObjectKey{Namespace: config.Namespace, Name: *config.Status.DataSecretName}
		log = log.WithValues("data-secret-name", secretKey.Name)
		existingSecret := &corev1.Secret{}

		err := r.Client.Get(ctx, secretKey, existingSecret)
		if err != nil && !apierrors.IsNotFound(err) {
			log.Error(err, "unable to check for existing bootstrap secret")
			return ctrl.Result{}, err
		}
		if err == nil {
			// We already have a secret that we don't need to regenerate
			return ctrl.Result{}, nil
		}
	}

	if cluster.Spec.ControlPlaneRef.Kind != "AWSManagedControlPlane" {
		return ctrl.Result{}, errors.New("Cluster's controlPlaneRef needs to be an AWSManagedControlPlane in order to use the EKS bootstrap provider")
	}

	if !ptr.Deref(cluster.Status.Initialization.InfrastructureProvisioned, false) {
		log.Info("Cluster infrastructure is not ready")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.WaitingForClusterInfrastructureReason,
			clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if !ptr.Deref(cluster.Status.Initialization.ControlPlaneInitialized, false) {
		log.Info("Control Plane has not yet been initialized")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.WaitingForControlPlaneInitializationReason, clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Namespace}, controlPlane); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to get control plane")
	}
	// Check if control plane is ready
	if !v1beta1conditions.IsTrue(controlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition) {
		log.Info("Waiting for control plane to be ready")
		v1beta1conditions.MarkFalse(
			config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityInfo,
			"Control plane is not initialized yet",
		)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}
	log.Info("Control plane is ready, proceeding with userdata generation")

	log.Info("Generating userdata")
	fileResolver := FileResolver{Client: r.Client}
	files, err := fileResolver.ResolveFiles(ctx, config.Namespace, config.Spec.Files)
	if err != nil {
		log.Info("Failed to resolve files for user data")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return ctrl.Result{}, err
	}

	serviceCIDR := ""
	if len(cluster.Spec.ClusterNetwork.Services.CIDRBlocks) > 0 {
		serviceCIDR = cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]
	}
	nodeInput := &userdata.NodeadmInput{
		// AWSManagedControlPlane webhooks default and validate EKSClusterName
		ClusterName:        controlPlane.Spec.EKSClusterName,
		PreNodeadmCommands: config.Spec.PreNodeadmCommands,
		Users:              config.Spec.Users,
		NTP:                config.Spec.NTP,
		DiskSetup:          config.Spec.DiskSetup,
		Mounts:             config.Spec.Mounts,
		Files:              files,
		ServiceCIDR:        serviceCIDR,
		APIServerEndpoint:  cluster.Spec.ControlPlaneEndpoint.Host,
	}
	if config.Spec.Kubelet != nil {
		nodeInput.KubeletFlags = config.Spec.Kubelet.Flags
		if config.Spec.Kubelet.Config != nil {
			nodeInput.KubeletConfig = config.Spec.Kubelet.Config
		}
	}
	if config.Spec.Containerd != nil {
		nodeInput.ContainerdConfig = config.Spec.Containerd.Config
		if config.Spec.Containerd.BaseRuntimeSpec != nil {
			nodeInput.ContainerdBaseRuntimeSpec = config.Spec.Containerd.BaseRuntimeSpec
		}
	}
	if config.Spec.FeatureGates != nil {
		nodeInput.FeatureGates = config.Spec.FeatureGates
	}

	// Fetch CA cert from KubeConfig secret
	obj := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Name,
	}
	ca, err := extractCAFromSecret(ctx, r.Client, obj)
	if err != nil {
		log.Error(err, "Failed to extract CA from kubeconfig secret")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityWarning,
			"Failed to extract CA from kubeconfig secret: %v", err)
		return ctrl.Result{}, err
	}
	nodeInput.CACert = ca

	log.Info("Generating nodeadm userdata",
		"cluster", controlPlane.Spec.EKSClusterName,
		"endpoint", nodeInput.APIServerEndpoint)
	// generate userdata
	userDataScript, err := userdata.NewNodeadmUserdata(nodeInput)
	if err != nil {
		log.Error(err, "Failed to create a worker join configuration")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	// store userdata as secret
	if err := r.storeBootstrapData(ctx, cluster, config, userDataScript); err != nil {
		log.Error(err, "Failed to store bootstrap data")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	v1beta1conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return ctrl.Result{}, nil
}

// reconcileHybridNode handles bootstrap data generation for hybrid nodes.
func (r *NodeadmConfigReconciler) reconcileHybridNode(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, configOwner *bsutil.ConfigOwner) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// Validate control plane reference
	if cluster.Spec.ControlPlaneRef.Kind != "AWSManagedControlPlane" {
		return ctrl.Result{}, errors.New("Cluster's controlPlaneRef needs to be an AWSManagedControlPlane in order to use the EKS bootstrap provider")
	}

	// Skip if already ready
	if config.Status.DataSecretName != nil && configOwner.GetKind() == "Machine" {
		secretKey := client.ObjectKey{Namespace: config.Namespace, Name: *config.Status.DataSecretName}
		existingSecret := &corev1.Secret{}
		err := r.Client.Get(ctx, secretKey, existingSecret)
		switch {
		case err == nil:
			return ctrl.Result{}, nil
		case !apierrors.IsNotFound(err):
			log.Error(err, "unable to check for existing bootstrap secret")
			return ctrl.Result{}, err
		}
	}

	// Check cluster infrastructure readiness
	if !ptr.Deref(cluster.Status.Initialization.InfrastructureProvisioned, false) {
		log.Info("Cluster infrastructure is not ready")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.WaitingForClusterInfrastructureReason,
			clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	// Check control plane initialization
	if !ptr.Deref(cluster.Status.Initialization.ControlPlaneInitialized, false) {
		log.Info("Control Plane has not yet been initialized")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.WaitingForControlPlaneInitializationReason,
			clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Get control plane
	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Namespace}, controlPlane); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to get control plane")
	}

	// Check if control plane is ready
	if !v1beta1conditions.IsTrue(controlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition) {
		log.Info("Waiting for control plane to be ready")
		v1beta1conditions.MarkFalse(
			config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityInfo,
			"Control plane is not initialized yet",
		)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Control plane is ready, proceeding with hybrid node userdata generation")

	// Get or create SSM activation
	activationID, activationCode, err := r.getOrCreateSSMActivation(ctx, cluster, controlPlane, config, configOwner)
	if err != nil {
		log.Error(err, "Failed to get or create SSM activation")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityWarning,
			"Failed to get SSM activation: %v", err)
		return ctrl.Result{}, err
	}

	// Determine which userdata generation path to use
	var userDataScript []byte

	if config.Spec.Hybrid.CustomUserData != nil {
		// Custom template mode - user provides their own template
		log.Info("Generating custom hybrid userdata from user-provided template",
			"cluster", controlPlane.Spec.EKSClusterName,
			"region", controlPlane.Spec.Region)
		userDataScript, err = r.generateCustomHybridUserdata(config, controlPlane, activationID, activationCode)
	} else {
		// Default nodeadm mode - generate MIME multipart userdata
		log.Info("Generating hybrid nodeadm userdata",
			"cluster", controlPlane.Spec.EKSClusterName,
			"region", controlPlane.Spec.Region)
		userDataScript, err = r.generateNodeadmHybridUserdata(ctx, config, controlPlane, activationID, activationCode)
	}

	if err != nil {
		log.Error(err, "Failed to create hybrid node join configuration")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	// Store userdata as secret
	if err := r.storeBootstrapData(ctx, cluster, config, userDataScript); err != nil {
		log.Error(err, "Failed to store bootstrap data")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	v1beta1conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return ctrl.Result{}, nil
}

// generateCustomHybridUserdata generates userdata from a user-provided template.
// This completely replaces the default nodeadm MIME multipart userdata generation.
func (r *NodeadmConfigReconciler) generateCustomHybridUserdata(
	config *eksbootstrapv1.NodeadmConfig,
	controlPlane *ekscontrolplanev1.AWSManagedControlPlane,
	activationID, activationCode string,
) ([]byte, error) {
	customInput := &userdata.CustomHybridInput{
		ClusterName:    controlPlane.Spec.EKSClusterName,
		Region:         controlPlane.Spec.Region,
		ActivationID:   activationID,
		ActivationCode: activationCode,
	}

	// Add Kubernetes version if specified
	if controlPlane.Spec.Version != nil {
		customInput.KubernetesVersion = *controlPlane.Spec.Version
	}

	// Add optional kubelet configuration
	if config.Spec.Kubelet != nil {
		customInput.KubeletFlags = config.Spec.Kubelet.Flags
		customInput.KubeletConfig = config.Spec.Kubelet.Config
	}

	// Add optional containerd configuration
	if config.Spec.Containerd != nil {
		customInput.ContainerdConfig = config.Spec.Containerd.Config
	}

	return userdata.NewCustomHybridUserdata(
		config.Spec.Hybrid.CustomUserData.Template,
		customInput,
	)
}

// generateNodeadmHybridUserdata generates the default nodeadm MIME multipart userdata.
func (r *NodeadmConfigReconciler) generateNodeadmHybridUserdata(
	ctx context.Context,
	config *eksbootstrapv1.NodeadmConfig,
	controlPlane *ekscontrolplanev1.AWSManagedControlPlane,
	activationID, activationCode string,
) ([]byte, error) {
	// Resolve files from secrets if needed
	fileResolver := FileResolver{Client: r.Client}
	files, err := fileResolver.ResolveFiles(ctx, config.Namespace, config.Spec.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve files for user data: %w", err)
	}

	// Build hybrid node input
	nodeInput := &userdata.NodeadmInput{
		ClusterName:        controlPlane.Spec.EKSClusterName,
		Region:             controlPlane.Spec.Region,
		ActivationID:       activationID,
		ActivationCode:     activationCode,
		PreNodeadmCommands: config.Spec.PreNodeadmCommands,
		Files:              files,
		DiskSetup:          config.Spec.DiskSetup,
		Mounts:             config.Spec.Mounts,
		Users:              config.Spec.Users,
		NTP:                config.Spec.NTP,
	}

	// Add kubelet configuration
	if config.Spec.Kubelet != nil {
		nodeInput.KubeletFlags = config.Spec.Kubelet.Flags
		if config.Spec.Kubelet.Config != nil {
			nodeInput.KubeletConfig = config.Spec.Kubelet.Config
		}
	}

	// Add containerd configuration
	if config.Spec.Containerd != nil {
		nodeInput.ContainerdConfig = config.Spec.Containerd.Config
	}

	return userdata.NewNodeadmUserdata(nodeInput)
}

// getOrCreateSSMActivation retrieves or creates SSM activation credentials.
func (r *NodeadmConfigReconciler) getOrCreateSSMActivation(
	ctx context.Context,
	cluster *clusterv1.Cluster,
	controlPlane *ekscontrolplanev1.AWSManagedControlPlane,
	config *eksbootstrapv1.NodeadmConfig,
	configOwner *bsutil.ConfigOwner,
) (string, string, error) {
	log := logger.FromContext(ctx)

	if config.Spec.Hybrid == nil || config.Spec.Hybrid.SSM == nil {
		return "", "", errors.New("hybrid SSM configuration is required")
	}

	ssmOpts := config.Spec.Hybrid.SSM

	// Option 1: Use pre-created activation secret
	if ssmOpts.ActivationRef != nil {
		log.Info("Using referenced SSM activation secret", "secretName", ssmOpts.ActivationRef.Name)
		return r.getActivationFromSecret(ctx, config.Namespace, ssmOpts.ActivationRef.Name, config)
	}

	// Option 2: Auto-create activation
	if ssmOpts.ActivationConfig == nil {
		return "", "", errors.New("either activationRef or activationConfig must be specified")
	}

	// Check if we already have an activation stored
	if config.Status.SSMActivation != nil && config.Status.SSMActivation.SecretName != nil {
		secretName := *config.Status.SSMActivation.SecretName
		activationID, activationCode, err := r.getActivationFromSecret(ctx, config.Namespace, secretName, config)
		if err == nil {
			log.Info("Using existing auto-created SSM activation", "secretName", secretName)
			return activationID, activationCode, nil
		}
		// If secret not found, we need to create a new activation
		if !apierrors.IsNotFound(errors.Cause(err)) {
			return "", "", err
		}
		log.Info("SSM activation secret not found, creating new activation")
	}

	// Create SSM service
	ssmService, err := r.getSSMService(ctx, cluster, controlPlane)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create SSM service")
	}

	// Prepare activation parameters
	registrationLimit := int32(1)
	if ssmOpts.ActivationConfig.RegistrationLimit != nil {
		registrationLimit = *ssmOpts.ActivationConfig.RegistrationLimit
	}

	expirationHours := int32(24)
	if ssmOpts.ActivationConfig.ExpirationHours != nil {
		expirationHours = *ssmOpts.ActivationConfig.ExpirationHours
	}

	// Convert tags to infrav1.Tags
	tags := make(infrav1.Tags)
	for k, v := range ssmOpts.ActivationConfig.Tags {
		tags[k] = v
	}

	params := &ssm.HybridActivationParams{
		IAMRoleName:       ssmOpts.ActivationConfig.IAMRoleName,
		RegistrationLimit: registrationLimit,
		ExpirationHours:   expirationHours,
		Tags:              tags,
		ClusterName:       cluster.Name,
		Namespace:         config.Namespace,
		ConfigName:        config.Name,
		MachineName:       configOwner.GetName(),
		Description:       fmt.Sprintf("CAPA hybrid node activation for %s/%s", config.Namespace, config.Name),
	}

	log.Info("Creating SSM activation",
		"iamRoleName", params.IAMRoleName,
		"registrationLimit", params.RegistrationLimit,
		"expirationHours", params.ExpirationHours)

	// Create activation
	result, err := ssmService.CreateHybridActivation(ctx, params)
	if err != nil {
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.SSMActivationReadyCondition,
			eksbootstrapv1.SSMActivationCreationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"Failed to create SSM activation: %v", err)
		return "", "", errors.Wrap(err, "failed to create SSM activation")
	}

	log.Info("Created SSM activation", "activationID", result.ActivationID)

	// Store activation in secret
	secretName := config.Name + ssmActivationSecretSuffix
	if err := r.storeActivationSecret(ctx, cluster, config, secretName, result.ActivationID, result.ActivationCode); err != nil {
		// Try to delete the activation since we couldn't store the secret
		if delErr := ssmService.DeleteHybridActivation(ctx, result.ActivationID); delErr != nil {
			log.Error(delErr, "Failed to cleanup SSM activation after secret creation failure")
		}
		return "", "", errors.Wrap(err, "failed to store SSM activation secret")
	}

	// Update status
	config.Status.SSMActivation = &eksbootstrapv1.SSMActivationStatus{
		ActivationID:   aws.String(result.ActivationID),
		SecretName:     aws.String(secretName),
		ExpirationTime: &metav1.Time{Time: result.ExpirationTime},
	}

	v1beta1conditions.MarkTrue(config, eksbootstrapv1.SSMActivationReadyCondition)
	return result.ActivationID, result.ActivationCode, nil
}

// getActivationFromSecret retrieves SSM activation credentials from a secret.
func (r *NodeadmConfigReconciler) getActivationFromSecret(ctx context.Context, namespace, secretName string, config *eksbootstrapv1.NodeadmConfig) (string, string, error) {
	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: secretName}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			v1beta1conditions.MarkFalse(config,
				eksbootstrapv1.SSMActivationReadyCondition,
				eksbootstrapv1.SSMActivationSecretNotFoundReason,
				clusterv1beta1.ConditionSeverityError,
				"SSM activation secret %s not found", secretName)
		}
		return "", "", errors.Wrapf(err, "failed to get SSM activation secret %s", secretName)
	}

	activationID, ok := secret.Data[ssmActivationIDKey]
	if !ok || len(activationID) == 0 {
		return "", "", fmt.Errorf("SSM activation secret %s missing %s key", secretName, ssmActivationIDKey)
	}

	activationCode, ok := secret.Data[ssmActivationCodeKey]
	if !ok || len(activationCode) == 0 {
		return "", "", fmt.Errorf("SSM activation secret %s missing %s key", secretName, ssmActivationCodeKey)
	}

	v1beta1conditions.MarkTrue(config, eksbootstrapv1.SSMActivationReadyCondition)
	return string(activationID), string(activationCode), nil
}

// storeActivationSecret creates a secret containing SSM activation credentials.
func (r *NodeadmConfigReconciler) storeActivationSecret(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, secretName, activationID, activationCode string) error {
	log := logger.FromContext(ctx)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: eksbootstrapv1.GroupVersion.String(),
					Kind:       "NodeadmConfig",
					Name:       config.Name,
					UID:        config.UID,
					Controller: ptr.To(true),
				},
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			ssmActivationIDKey:   []byte(activationID),
			ssmActivationCodeKey: []byte(activationCode),
		},
	}

	if err := r.Client.Create(ctx, secret); err != nil {
		if apierrors.IsAlreadyExists(err) {
			// Update existing secret
			existing := &corev1.Secret{}
			if err := r.Client.Get(ctx, client.ObjectKey{Namespace: config.Namespace, Name: secretName}, existing); err != nil {
				return errors.Wrap(err, "failed to get existing activation secret")
			}
			existing.Data = secret.Data
			existing.Labels = secret.Labels
			existing.OwnerReferences = secret.OwnerReferences
			if err := r.Client.Update(ctx, existing); err != nil {
				return errors.Wrap(err, "failed to update activation secret")
			}
			log.Info("Updated SSM activation secret", "secret", secretName)
			return nil
		}
		return errors.Wrap(err, "failed to create activation secret")
	}

	log.Info("Created SSM activation secret", "secret", secretName)
	return nil
}

// reconcileDelete handles cleanup when a NodeadmConfig is being deleted.
func (r *NodeadmConfigReconciler) reconcileDelete(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// Only process deletion for hybrid nodes with auto-created activations
	if !controllerutil.ContainsFinalizer(config, ssmActivationFinalizer) {
		return ctrl.Result{}, nil
	}

	// Check if we have an activation to clean up
	if config.Status.SSMActivation != nil && config.Status.SSMActivation.ActivationID != nil {
		activationID := *config.Status.SSMActivation.ActivationID
		log.Info("Deleting SSM activation", "activationID", activationID)

		// Get control plane for AWS credentials
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Namespace}, controlPlane); err != nil {
			if !apierrors.IsNotFound(err) {
				return ctrl.Result{}, errors.Wrap(err, "failed to get control plane")
			}
			// Control plane already deleted, can't delete SSM activation
			log.Info("Control plane not found, skipping SSM activation deletion", "activationID", activationID)
		} else {
			// Create SSM service and delete activation
			ssmService, err := r.getSSMService(ctx, cluster, controlPlane)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, "failed to create SSM service for cleanup")
			}

			if err := ssmService.DeleteHybridActivation(ctx, activationID); err != nil {
				log.Error(err, "Failed to delete SSM activation", "activationID", activationID)
				return ctrl.Result{}, errors.Wrap(err, "failed to delete SSM activation")
			}
			log.Info("Deleted SSM activation", "activationID", activationID)
		}
	}

	// Remove finalizer
	controllerutil.RemoveFinalizer(config, ssmActivationFinalizer)
	return ctrl.Result{}, nil
}

// getSSMService creates an SSM service using control plane credentials.
func (r *NodeadmConfigReconciler) getSSMService(ctx context.Context, cluster *clusterv1.Cluster, controlPlane *ekscontrolplanev1.AWSManagedControlPlane) (*ssm.Service, error) {
	// Create managed control plane scope for AWS credentials
	managedScope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		ControllerName: "nodeadmconfig",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create managed control plane scope")
	}

	if r.SSMServiceFactory != nil {
		return r.SSMServiceFactory(managedScope), nil
	}
	return ssm.NewService(managedScope), nil
}

// storeBootstrapData creates a new secret with the data passed in as input,
// sets the reference in the configuration status and ready to true.
func (r *NodeadmConfigReconciler) storeBootstrapData(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, data []byte) error {
	log := logger.FromContext(ctx)

	// as secret creation and scope.Config status patch are not atomic operations
	// it is possible that secret creation happens but the config.Status patches are not applied
	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      config.Name,
		Namespace: config.Namespace,
	}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			if secret, err = r.createBootstrapSecret(ctx, cluster, config, data); err != nil {
				return errors.Wrap(err, "failed to create bootstrap data secret for NodeadmConfig")
			}
			log.Info("created bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		} else {
			return errors.Wrap(err, "failed to get data secret for NodeadmConfig")
		}
	} else {
		updated, err := r.updateBootstrapSecret(ctx, secret, data)
		if err != nil {
			return errors.Wrap(err, "failed to update data secret for NodeadmConfig")
		}
		if updated {
			log.Info("updated bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		} else {
			log.Trace("no change in bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		}
	}

	config.Status.DataSecretName = ptr.To(secret.Name)
	config.Status.Initialization.DataSecretCreated = ptr.To(true)
	//nolint:staticcheck // we will support this implementation until CAPA is v1beta2 compliant
	config.Status.Ready = true
	v1beta1conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return nil
}

func (r *NodeadmConfigReconciler) createBootstrapSecret(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, data []byte) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: eksbootstrapv1.GroupVersion.String(),
					Kind:       "NodeadmConfig",
					Name:       config.Name,
					UID:        config.UID,
					Controller: ptr.To[bool](true),
				},
			},
		},
		Data: map[string][]byte{
			"value": data,
		},
		Type: clusterv1.ClusterSecretType,
	}
	return secret, r.Client.Create(ctx, secret)
}

// Update the userdata in the bootstrap Secret.
func (r *NodeadmConfigReconciler) updateBootstrapSecret(ctx context.Context, secret *corev1.Secret, data []byte) (bool, error) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	if !bytes.Equal(secret.Data["value"], data) {
		secret.Data["value"] = data
		return true, r.Client.Update(ctx, secret)
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeadmConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, option controller.Options) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&eksbootstrapv1.NodeadmConfig{}).
		WithOptions(option).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), logger.FromContext(ctx).GetLogger(), r.WatchFilterValue)).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(r.MachineToBootstrapMapFunc),
		)

	if feature.Gates.Enabled(feature.MachinePool) {
		b = b.Watches(
			&clusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(r.MachinePoolToBootstrapMapFunc),
		)
	}

	c, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc((r.ClusterToNodeadmConfigs)),
			predicates.ClusterPausedTransitionsOrInfrastructureProvisioned(mgr.GetScheme(), logger.FromContext(ctx).GetLogger())),
	)
	if err != nil {
		return errors.Wrap(err, "failed adding watch for Clusters to controller manager")
	}
	return nil
}

// MachineToBootstrapMapFunc is a handler.ToRequestsFunc to be used to enque requests for
// NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) MachineToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*clusterv1.Machine)
	if !ok {
		klog.Errorf("Expected a Machine but got a %T", o)
		return result
	}
	if m.Spec.Bootstrap.ConfigRef.IsDefined() && m.Spec.Bootstrap.ConfigRef.APIGroup == eksbootstrapv1.GroupVersion.Group && m.Spec.Bootstrap.ConfigRef.Kind == eksbootstrapv1.NodeadmConfigKind {
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}
	return result
}

// MachinePoolToBootstrapMapFunc is a handler.ToRequestsFunc to be uses to enqueue requests
// for NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) MachinePoolToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*clusterv1.MachinePool)
	if !ok {
		klog.Errorf("Expected a MachinePool but got a %T", o)
		return result
	}
	configRef := m.Spec.Template.Spec.Bootstrap.ConfigRef
	if configRef.IsDefined() && configRef.APIGroup == eksbootstrapv1.GroupVersion.Group && configRef.Kind == eksbootstrapv1.NodeadmConfigKind {
		name := client.ObjectKey{Namespace: m.Namespace, Name: configRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

// ClusterToNodeadmConfigs is a handler.ToRequestsFunc to be used to enqueue requests for
// NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) ClusterToNodeadmConfigs(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	c, ok := o.(*clusterv1.Cluster)
	if !ok {
		klog.Errorf("Expected a Cluster but got a %T", o)
		return result
	}

	selectors := []client.ListOption{
		client.InNamespace(c.Namespace),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: c.Name,
		},
	}

	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.Background(), machineList, selectors...); err != nil {
		return nil
	}

	for _, m := range machineList.Items {
		if m.Spec.Bootstrap.ConfigRef.IsDefined() &&
			m.Spec.Bootstrap.ConfigRef.APIGroup == eksbootstrapv1.GroupVersion.Group &&
			m.Spec.Bootstrap.ConfigRef.Kind == eksbootstrapv1.NodeadmConfigKind {
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}
	}

	return result
}

func extractCAFromSecret(ctx context.Context, c client.Client, obj client.ObjectKey) (string, error) {
	data, err := kubeconfigutil.FromSecret(ctx, c, obj)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get kubeconfig secret %s", obj.Name)
	}
	config, err := clientcmd.Load(data)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse kubeconfig data from secret %s", obj.Name)
	}

	// Iterate through all clusters in the kubeconfig and use the first one with CA data
	for _, cluster := range config.Clusters {
		if len(cluster.CertificateAuthorityData) > 0 {
			return base64.StdEncoding.EncodeToString(cluster.CertificateAuthorityData), nil
		}
	}

	return "", fmt.Errorf("no cluster with CA data found in kubeconfig")
}
