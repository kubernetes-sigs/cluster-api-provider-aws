package operator

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/clusterautoscaler"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/machineautoscaler"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// OperatorName is the name of this operator.
const OperatorName = "cluster-autoscaler"

// Operator represents an instance of the cluster-autoscaler-operator.
type Operator struct {
	config  *Config
	manager manager.Manager
}

// New returns a new Operator instance with the given config and a
// manager configured with the various controllers.
func New(cfg *Config) (*Operator, error) {
	operator := &Operator{config: cfg}

	// Get a config to talk to the apiserver.
	clientConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// Create the controller-manager.
	managerOptions := manager.Options{
		Namespace:               cfg.WatchNamespace,
		LeaderElection:          cfg.LeaderElection,
		LeaderElectionNamespace: cfg.LeaderElectionNamespace,
		LeaderElectionID:        cfg.LeaderElectionID,
	}

	operator.manager, err = manager.New(clientConfig, managerOptions)
	if err != nil {
		return nil, err
	}

	// Setup Scheme for all resources.
	if err := apis.AddToScheme(operator.manager.GetScheme()); err != nil {
		return nil, err
	}

	if err := operator.AddControllers(); err != nil {
		return nil, err
	}

	statusConfig := &StatusReporterConfig{
		ClusterAutoscalerName:      cfg.ClusterAutoscalerName,
		ClusterAutoscalerNamespace: cfg.ClusterAutoscalerNamespace,
		ReleaseVersion:             cfg.ReleaseVersion,
		RelatedObjects:             operator.RelatedObjects(),
	}

	statusReporter, err := NewStatusReporter(operator.manager, statusConfig)
	if err != nil {
		return nil, err
	}

	operator.manager.Add(statusReporter)

	return operator, nil
}

// RelatedObjects returns a list of objects related to the operator and its
// operands.  These are used in the ClusterOperator status.
func (o *Operator) RelatedObjects() []configv1.ObjectReference {
	relatedNamespaces := map[string]string{}

	relatedNamespaces[o.config.WatchNamespace] = ""
	relatedNamespaces[o.config.LeaderElectionNamespace] = ""
	relatedNamespaces[o.config.ClusterAutoscalerNamespace] = ""

	relatedObjects := []configv1.ObjectReference{}

	for namespace := range relatedNamespaces {
		relatedObjects = append(relatedObjects, configv1.ObjectReference{
			Resource: "namespaces",
			Name:     namespace,
		})
	}

	return relatedObjects
}

// AddControllers configures the various controllers and adds them to
// the operator's manager instance.
func (o *Operator) AddControllers() error {
	// Setup ClusterAutoscaler controller.
	ca := clusterautoscaler.NewReconciler(o.manager, &clusterautoscaler.Config{
		ReleaseVersion: o.config.ReleaseVersion,
		Name:           o.config.ClusterAutoscalerName,
		Image:          o.config.ClusterAutoscalerImage,
		Replicas:       o.config.ClusterAutoscalerReplicas,
		Namespace:      o.config.ClusterAutoscalerNamespace,
		CloudProvider:  o.config.ClusterAutoscalerCloudProvider,
		Verbosity:      o.config.ClusterAutoscalerVerbosity,
		ExtraArgs:      o.config.ClusterAutoscalerExtraArgs,
	})

	if err := ca.AddToManager(o.manager); err != nil {
		return err
	}

	// Setup MachineAutoscaler controller.
	ma := machineautoscaler.NewReconciler(o.manager, &machineautoscaler.Config{
		Namespace:           o.config.ClusterAutoscalerNamespace,
		SupportedTargetGVKs: machineautoscaler.DefaultSupportedTargetGVKs(),
	})

	if err := ma.AddToManager(o.manager); err != nil {
		return err
	}

	return nil
}

// Start starts the operator's controller-manager.
func (o *Operator) Start() error {
	stopCh := signals.SetupSignalHandler()

	return o.manager.Start(stopCh)
}
