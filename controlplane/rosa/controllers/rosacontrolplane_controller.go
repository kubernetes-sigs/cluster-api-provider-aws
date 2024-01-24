/*
Copyright 2023 The Kubernetes Authors.

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
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	clustersmgmtv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	capiannotations "sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/secret"
)

const (
	rosaCreatorArnProperty = "rosa_creator_arn"

	rosaControlPlaneKind = "ROSAControlPlane"
	// ROSAControlPlaneFinalizer allows the controller to clean up resources on delete.
	ROSAControlPlaneFinalizer = "rosacontrolplane.controlplane.cluster.x-k8s.io"
)

type ROSAControlPlaneReconciler struct {
	client.Client
	WatchFilterValue string
	WaitInfraPeriod  time.Duration
}

// SetupWithManager is used to setup the controller.
func (r *ROSAControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(rosaControlPlane).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log.GetLogger(), r.WatchFilterValue)).
		Build(r)

	if err != nil {
		return fmt.Errorf("failed setting up the AWSManagedControlPlane controller manager: %w", err)
	}

	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, rosaControlPlane.GroupVersionKind(), mgr.GetClient(), &expinfrav1.ROSACluster{})),
		predicates.ClusterUnpausedAndInfrastructureReady(log.GetLogger()),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	if err = c.Watch(
		source.Kind(mgr.GetCache(), &expinfrav1.ROSACluster{}),
		handler.EnqueueRequestsFromMapFunc(r.rosaClusterToROSAControlPlane(log)),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ROSACluster")
	}

	return nil
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes/status,verbs=get;update;patch

// Reconcile will reconcile RosaControlPlane Resources.
func (r *ROSAControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	// Get the control plane instance
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaControlPlane); err != nil {
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

	if capiannotations.IsPaused(cluster, rosaControlPlane) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	rosaScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   rosaControlPlane,
		ControllerName: strings.ToLower(rosaControlPlaneKind),
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := rosaScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, rosaScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaScope)
}

func (r *ROSAControlPlaneReconciler) reconcileNormal(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane")

	// if !rosaScope.Cluster.Status.InfrastructureReady {
	//	rosaScope.Info("Cluster infrastructure is not ready yet")
	//	return ctrl.Result{RequeueAfter: r.WaitInfraPeriod}, nil
	//}
	if controllerutil.AddFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer) {
		if err := rosaScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}

	rosaClient, err := rosa.NewRosaClient(ctx, rosaScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create a rosa client: %w", err)
	}
	defer rosaClient.Close()

	cluster, err := rosaClient.GetCluster()
	if err != nil {
		return ctrl.Result{}, err
	}

	if clusterID := cluster.ID(); clusterID != "" {
		rosaScope.ControlPlane.Status.ID = &clusterID
		if cluster.Status().State() == clustersmgmtv1.ClusterStateReady {
			conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneReadyCondition)
			rosaScope.ControlPlane.Status.Ready = true

			apiEndpoint, err := buildAPIEndpoint(cluster)
			if err != nil {
				return ctrl.Result{}, err
			}
			rosaScope.ControlPlane.Spec.ControlPlaneEndpoint = *apiEndpoint

			if err := r.reconcileKubeconfig(ctx, rosaScope, rosaClient, cluster); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to reconcile kubeconfig: %w", err)
			}

			return ctrl.Result{}, nil
		}

		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ROSAControlPlaneReadyCondition,
			string(cluster.Status().State()),
			clusterv1.ConditionSeverityInfo,
			"")

		rosaScope.Info("waiting for cluster to become ready", "state", cluster.Status().State())
		// Requeue so that status.ready is set to true when the cluster is fully created.
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	clusterSpec, err := ocmCluster(rosaScope, nil)
	if err != nil {
		return ctrl.Result{}, err
	}

	newCluster, err := rosaClient.CreateCluster(clusterSpec)
	if err != nil {
		rosaScope.Info("error", "error", err)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	rosaScope.Info("cluster created", "state", newCluster.Status().State())
	clusterID := newCluster.ID()
	rosaScope.ControlPlane.Status.ID = &clusterID

	return ctrl.Result{}, nil
}

func ocmCluster(rosaScope *scope.ROSAControlPlaneScope, now func() time.Time) (*clustersmgmtv1.Cluster, error) {
	if now == nil {
		now = time.Now
	}
	// Create the cluster:
	clusterBuilder := clustersmgmtv1.NewCluster().
		Name(rosaScope.RosaClusterName()).
		MultiAZ(true).
		Product(
			clustersmgmtv1.NewProduct().
				ID("rosa"),
		).
		Region(
			clustersmgmtv1.NewCloudRegion().
				ID(*rosaScope.ControlPlane.Spec.Region),
		).
		FIPS(false).
		EtcdEncryption(false).
		DisableUserWorkloadMonitoring(true).
		Version(
			clustersmgmtv1.NewVersion().
				ID(*rosaScope.ControlPlane.Spec.Version).
				ChannelGroup("stable"),
		).
		ExpirationTimestamp(now().Add(1 * time.Hour)).
		Hypershift(clustersmgmtv1.NewHypershift().Enabled(true))

	networkBuilder := clustersmgmtv1.NewNetwork()
	networkBuilder = networkBuilder.Type("OVNKubernetes")
	networkBuilder = networkBuilder.MachineCIDR(*rosaScope.ControlPlane.Spec.MachineCIDR)
	clusterBuilder = clusterBuilder.Network(networkBuilder)

	stsBuilder := clustersmgmtv1.NewSTS().RoleARN(*rosaScope.ControlPlane.Spec.InstallerRoleARN)
	// stsBuilder = stsBuilder.ExternalID(config.ExternalID)
	stsBuilder = stsBuilder.SupportRoleARN(*rosaScope.ControlPlane.Spec.SupportRoleARN)
	roles := []*clustersmgmtv1.OperatorIAMRoleBuilder{}
	for _, role := range []struct {
		Name      string
		Namespace string
		RoleARN   string
		Path      string
	}{
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-ingress-operator",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.IngressARN,
		},
		{
			Name:      "installer-cloud-credentials",
			Namespace: "openshift-image-registry",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.ImageRegistryARN,
		},
		{
			Name:      "ebs-cloud-credentials",
			Namespace: "openshift-cluster-csi-drivers",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.StorageARN,
		},
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-cloud-network-config-controller",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.NetworkARN,
		},
		{
			Name:      "kube-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.KubeCloudControllerARN,
		},
		{
			Name:      "kms-provider",
			Namespace: "kube-system",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.KMSProviderARN,
		},
		{
			Name:      "control-plane-operator",
			Namespace: "kube-system",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.ControlPlaneOperatorARN,
		},
		{
			Name:      "capa-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rosaScope.ControlPlane.Spec.RolesRef.NodePoolManagementARN,
		},
	} {
		roles = append(roles, clustersmgmtv1.NewOperatorIAMRole().
			Name(role.Name).
			Namespace(role.Namespace).
			RoleARN(role.RoleARN))
	}
	stsBuilder = stsBuilder.OperatorIAMRoles(roles...)

	instanceIAMRolesBuilder := clustersmgmtv1.NewInstanceIAMRoles()
	instanceIAMRolesBuilder.WorkerRoleARN(*rosaScope.ControlPlane.Spec.WorkerRoleARN)
	stsBuilder = stsBuilder.InstanceIAMRoles(instanceIAMRolesBuilder)
	stsBuilder.OidcConfig(clustersmgmtv1.NewOidcConfig().ID(*rosaScope.ControlPlane.Spec.OIDCID))
	stsBuilder.AutoMode(true)

	awsBuilder := clustersmgmtv1.NewAWS().
		AccountID(*rosaScope.ControlPlane.Spec.AccountID).
		BillingAccountID(*rosaScope.ControlPlane.Spec.AccountID).
		SubnetIDs(rosaScope.ControlPlane.Spec.Subnets...).
		STS(stsBuilder)
	clusterBuilder = clusterBuilder.AWS(awsBuilder)

	clusterNodesBuilder := clustersmgmtv1.NewClusterNodes()
	clusterNodesBuilder = clusterNodesBuilder.AvailabilityZones(rosaScope.ControlPlane.Spec.AvailabilityZones...)
	clusterBuilder = clusterBuilder.Nodes(clusterNodesBuilder)

	clusterProperties := map[string]string{}
	clusterProperties[rosaCreatorArnProperty] = *rosaScope.ControlPlane.Spec.CreatorARN

	clusterBuilder = clusterBuilder.Properties(clusterProperties)
	return clusterBuilder.Build()
}

func (r *ROSAControlPlaneReconciler) reconcileDelete(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane delete")

	rosaClient, err := rosa.NewRosaClient(ctx, rosaScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create a rosa client: %w", err)
	}
	defer rosaClient.Close()

	cluster, err := rosaClient.GetCluster()
	if err != nil {
		return ctrl.Result{}, err
	}

	if cluster != nil {
		if err := rosaClient.DeleteCluster(cluster.ID()); err != nil {
			return ctrl.Result{}, err
		}
	}

	controllerutil.RemoveFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer)

	return ctrl.Result{}, nil
}

func (r *ROSAControlPlaneReconciler) reconcileKubeconfig(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope, rosaClient *rosa.RosaClient, cluster *clustersmgmtv1.Cluster) error {
	rosaScope.Debug("Reconciling ROSA kubeconfig for cluster", "cluster-name", rosaScope.RosaClusterName())

	clusterRef := client.ObjectKeyFromObject(rosaScope.Cluster)
	kubeconfigSecret, err := secret.GetFromNamespacedName(ctx, r.Client, clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get kubeconfig secret: %w", err)
		}
	}

	// generate a new password for the cluster admin user, or retrieve an existing one.
	password, err := r.reconcileClusterAdminPassword(ctx, rosaScope)
	if err != nil {
		return fmt.Errorf("failed to reconcile cluster admin password secret: %w", err)
	}

	clusterName := rosaScope.RosaClusterName()
	userName := fmt.Sprintf("%s-capi-admin", clusterName)
	apiServerURL := cluster.API().URL()

	// create new user with admin privileges in the ROSA cluster if 'userName' doesn't already exist.
	err = rosaClient.CreateAdminUserIfNotExist(cluster.ID(), userName, password)
	if err != nil {
		return err
	}

	clientConfig := &restclient.Config{
		Host:     apiServerURL,
		Username: userName,
	}
	// request an acccess token using the credentials of the cluster admin user created earlier.
	// this token is used in the kubeconfig to authenticate with the API server.
	token, err := rosa.RequestToken(ctx, apiServerURL, userName, password, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}

	// create the kubeconfig spec.
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)
	cfg := &api.Config{
		APIVersion: api.SchemeGroupVersion.Version,
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server: apiServerURL,
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		CurrentContext: contextName,
		AuthInfos: map[string]*api.AuthInfo{
			userName: {
				Token: token.AccessToken,
			},
		},
	}
	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config to yaml: %w", err)
	}

	if kubeconfigSecret != nil {
		// update existing kubeconfig secret.
		kubeconfigSecret.Data[secret.KubeconfigDataName] = out
		if err := r.Client.Update(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to update kubeconfig secret: %w", err)
		}
	} else {
		// create new kubeconfig secret.
		controllerOwnerRef := *metav1.NewControllerRef(rosaScope.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"))
		kubeconfigSecret = kubeconfig.GenerateSecretWithOwner(clusterRef, out, controllerOwnerRef)
		if err := r.Client.Create(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to create kubeconfig secret: %w", err)
		}
	}

	rosaScope.ControlPlane.Status.Initialized = true
	return nil
}

// reconcileClusterAdminPassword generates and store the password of the cluster admin user in a secret which is used to request a token for kubeconfig auth.
// Since it is not possible to retrieve a user's password through the ocm API once created,
// we have to store the password in a secret as it is needed later to refresh the token.
func (r *ROSAControlPlaneReconciler) reconcileClusterAdminPassword(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, error) {
	passwordSecret := rosaScope.ClusterAdminPasswordSecret()
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(passwordSecret), passwordSecret)
	if err == nil {
		password := string(passwordSecret.Data["value"])
		return password, nil
	} else if !apierrors.IsNotFound(err) {
		return "", fmt.Errorf("failed to get cluster admin password secret: %w", err)
	}
	// Generate a new password and create the secret
	password, err := rosa.GenerateRandomPassword()
	if err != nil {
		return "", err
	}

	controllerOwnerRef := *metav1.NewControllerRef(rosaScope.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"))
	passwordSecret.Data = map[string][]byte{
		"value": []byte(password),
	}
	passwordSecret.OwnerReferences = []metav1.OwnerReference{
		controllerOwnerRef,
	}
	if err := r.Client.Create(ctx, passwordSecret); err != nil {
		return "", err
	}

	return password, nil
}

func (r *ROSAControlPlaneReconciler) rosaClusterToROSAControlPlane(log *logger.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		rosaCluster, ok := o.(*expinfrav1.ROSACluster)
		if !ok {
			log.Error(fmt.Errorf("expected a ROSACluster but got a %T", o), "Expected ROSACluster")
			return nil
		}

		if !rosaCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Debug("ROSACluster has a deletion timestamp, skipping mapping")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaCluster.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get owning cluster")
			return nil
		}
		if cluster == nil {
			log.Debug("Owning cluster not set on ROSACluster, skipping mapping")
			return nil
		}

		controlPlaneRef := cluster.Spec.ControlPlaneRef
		if controlPlaneRef == nil || controlPlaneRef.Kind != rosaControlPlaneKind {
			log.Debug("ControlPlaneRef is nil or not ROSAControlPlane, skipping mapping")
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      controlPlaneRef.Name,
					Namespace: controlPlaneRef.Namespace,
				},
			},
		}
	}
}

func buildAPIEndpoint(cluster *clustersmgmtv1.Cluster) (*clusterv1.APIEndpoint, error) {
	parsedURL, err := url.ParseRequestURI(cluster.API().URL())
	if err != nil {
		return nil, err
	}
	host, portStr, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &clusterv1.APIEndpoint{
		Host: host,
		Port: int32(port), // #nosec G109
	}, nil
}
