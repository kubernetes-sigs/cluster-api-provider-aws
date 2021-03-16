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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"

	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/internal/userdata"
)

// EKSConfigReconciler reconciles a EKSConfig object
type EKSConfigReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machinepools;clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=exp.cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;

func (r *EKSConfigReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, rerr error) {
	ctx := context.Background()
	log := r.Log.WithValues("namespace", req.NamespacedName.Namespace, "name", req.NamespacedName.Name)

	// get EKSConfig
	config := &bootstrapv1.EKSConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}

	// check owner references and look up owning Machine object
	configOwner, err := bsutil.GetConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// no error here, requeue until we find an owner
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "Failed to get owner")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		// no error, requeue until we find an owner
		return ctrl.Result{}, nil
	}

	log = log.WithValues(configOwner.GetKind(), configOwner.GetName())

	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Is(err, util.ErrNoCluster) {
			log.Info("EKSConfig does not belong to a cluster yet, re-queuing until it's partof a cluster")
			return ctrl.Result{}, nil
		}
		if apierrors.IsNotFound(err) {
			log.Info("Cluster does not exist yet, re-queueing until it is created")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Could not get cluster with metadata")
		return ctrl.Result{}, err
	}
	log = log.WithValues("cluster", cluster.Name)

	if annotations.IsPaused(cluster, config) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	patchHelper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// set up defer block for updating config
	defer func() {
		conditions.SetSummary(config,
			conditions.WithConditions(
				bootstrapv1.DataSecretAvailableCondition,
			),
			conditions.WithStepCounter(),
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

	return r.joinWorker(ctx, log, cluster, config)
}

func (r *EKSConfigReconciler) joinWorker(ctx context.Context, log logr.Logger, cluster *clusterv1.Cluster, config *bootstrapv1.EKSConfig) (ctrl.Result, error) {
	if config.Status.DataSecretName != nil {
		secretKey := client.ObjectKey{Namespace: config.Namespace, Name: *config.Status.DataSecretName}
		log = log.WithValues("data-secret-name", secretKey.Name)
		existingSecret := &corev1.Secret{}

		// No error here means the Secret exists and we have no
		// reason to proceed.
		err := r.Client.Get(ctx, secretKey, existingSecret)
		switch {
		case err == nil:
			return ctrl.Result{}, nil
		case !apierrors.IsNotFound(err):
			log.Error(err, "unable to check for existing bootstrap secret")
			return ctrl.Result{}, err
		}
	}

	if cluster.Spec.ControlPlaneRef == nil || cluster.Spec.ControlPlaneRef.Kind != "AWSManagedControlPlane" {
		return ctrl.Result{}, errors.New("Cluster's controlPlaneRef needs to be an AWSManagedControlPlane in order to use the EKS bootstrap provider")
	}

	if !cluster.Status.InfrastructureReady {
		log.Info("Cluster infrastructure is not ready")
		conditions.MarkFalse(config,
			bootstrapv1.DataSecretAvailableCondition,
			bootstrapv1.WaitingForClusterInfrastructureReason,
			clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if !cluster.Status.ControlPlaneInitialized {
		log.Info("Control Plane has not yet been initialized")
		conditions.MarkFalse(config, bootstrapv1.DataSecretAvailableCondition, bootstrapv1.WaitingForControlPlaneInitializationReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Spec.ControlPlaneRef.Namespace}, controlPlane); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Generating userdata")

	// generate userdata
	userDataScript, err := userdata.NewNode(&userdata.NodeInput{
		// AWSManagedControlPlane webhooks default and validate EKSClusterName
		ClusterName: controlPlane.Spec.EKSClusterName,

		KubeletExtraArgs: config.Spec.KubeletExtraArgs,
	})
	if err != nil {
		log.Error(err, "Failed to create a worker join configuration")
		conditions.MarkFalse(config, bootstrapv1.DataSecretAvailableCondition, bootstrapv1.DataSecretGenerationFailedReason, clusterv1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	// store userdata as secret
	if err := r.storeBootstrapData(ctx, log, cluster, config, userDataScript); err != nil {
		log.Error(err, "Failed to store bootstrap data")
		conditions.MarkFalse(config, bootstrapv1.DataSecretAvailableCondition, bootstrapv1.DataSecretGenerationFailedReason, clusterv1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EKSConfigReconciler) SetupWithManager(mgr ctrl.Manager, option controller.Options) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&bootstrapv1.EKSConfig{}).
		WithOptions(option).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(r.Log, r.WatchFilterValue)).
		Watches(
			&source.Kind{Type: &clusterv1.Machine{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.MachineToBootstrapMapFunc),
			},
		)

	if feature.Gates.Enabled(feature.MachinePool) {
		b = b.Watches(
			&source.Kind{Type: &expv1.MachinePool{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.MachinePoolToBootstrapMapFunc),
			},
		)
	}

	c, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	err = c.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(r.ClusterToEKSConfigs),
		},
		predicates.ClusterUnpausedAndInfrastructureReady(r.Log),
	)
	if err != nil {
		return errors.Wrap(err, "failed adding watch for Clusters to controller manager")
	}

	return nil
}

// storeBootstrapData creates a new secret with the data passed in as input,
// sets the reference in the configuration status and ready to true.
func (r *EKSConfigReconciler) storeBootstrapData(ctx context.Context, log logr.Logger, cluster *clusterv1.Cluster, config *bootstrapv1.EKSConfig, data []byte) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterLabelName: cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: bootstrapv1.GroupVersion.String(),
					Kind:       "EKSConfig",
					Name:       config.Name,
					UID:        config.UID,
					Controller: pointer.BoolPtr(true),
				},
			},
		},
		Data: map[string][]byte{
			"value": data,
		},
		Type: clusterv1.ClusterSecretType,
	}

	// as secret creation and scope.Config status patch are not atomic operations
	// it is possible that secret creation happens but the config.Status patches are not applied
	if err := r.Client.Create(ctx, secret); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrap(err, "failed to create bootstrap data secret for EKSConfig")
		}
		log.Info("bootstrap data secret for EKSConfig already exists", "secret", secret.Name)
	}
	config.Status.DataSecretName = pointer.StringPtr(secret.Name)
	config.Status.Ready = true
	conditions.MarkTrue(config, bootstrapv1.DataSecretAvailableCondition)
	return nil
}

// MachineToBootstrapMapFunc is a handler.ToRequestsFunc to be used to enqueue requests
// for EKSConfig reconciliation
func (r *EKSConfigReconciler) MachineToBootstrapMapFunc(o handler.MapObject) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.Object.(*clusterv1.Machine)
	if !ok {
		return nil
	}
	if m.Spec.Bootstrap.ConfigRef != nil && m.Spec.Bootstrap.ConfigRef.GroupVersionKind() == bootstrapv1.GroupVersion.WithKind("EKSConfig") {
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}
	return result
}

// MachinePoolToBootstrapMapFunc is a handler.ToRequestsFunc to be uses to enqueue requests
// for EKSConfig reconciliation
func (r *EKSConfigReconciler) MachinePoolToBootstrapMapFunc(o handler.MapObject) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.Object.(*expv1.MachinePool)
	if !ok {
		return nil
	}
	configRef := m.Spec.Template.Spec.Bootstrap.ConfigRef
	if configRef != nil && configRef.GroupVersionKind().GroupKind() == bootstrapv1.GroupVersion.WithKind("EKSConfig").GroupKind() {
		name := client.ObjectKey{Namespace: m.Namespace, Name: configRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

// ClusterToEKSConfigs is a handler.ToRequestsFunc to be used to enqueue requests for
// EKSConfig reconciliation
func (r *EKSConfigReconciler) ClusterToEKSConfigs(o handler.MapObject) []ctrl.Request {
	result := []ctrl.Request{}

	c, ok := o.Object.(*clusterv1.Cluster)
	if !ok {
		return nil
	}

	selectors := []client.ListOption{
		client.InNamespace(c.Namespace),
		client.MatchingLabels{
			clusterv1.ClusterLabelName: c.Name,
		},
	}

	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.Background(), machineList, selectors...); err != nil {
		r.Log.Error(err, "failed to list Machines for Cluster", "name", c.Name, "namespace", c.Namespace)
		return nil
	}

	for _, m := range machineList.Items {
		if m.Spec.Bootstrap.ConfigRef != nil &&
			m.Spec.Bootstrap.ConfigRef.GroupVersionKind().GroupKind() == bootstrapv1.GroupVersion.WithKind("EKSConfig").GroupKind() {
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}
	}

	return result
}
