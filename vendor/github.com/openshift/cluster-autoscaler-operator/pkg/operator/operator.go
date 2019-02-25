package operator

import (
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/clusterautoscaler"
	"github.com/openshift/cluster-autoscaler-operator/pkg/controller/machineautoscaler"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// OperatorName is the name of this operator.
const OperatorName = "cluster-autoscaler-operator"

// Operator represents an instance of the cluster-autoscaler-operator.
type Operator struct {
	config  *Config
	status  *StatusReporter
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

	operator.status, err = NewStatusReporter(clientConfig)
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

	return operator, nil
}

// AddControllers configures the various controllers and adds them to
// the operator's manager instance.
func (o *Operator) AddControllers() error {
	// Setup ClusterAutoscaler controller.
	ca := clusterautoscaler.NewReconciler(o.manager, &clusterautoscaler.Config{
		Name:          o.config.ClusterAutoscalerName,
		Image:         o.config.ClusterAutoscalerImage,
		Replicas:      o.config.ClusterAutoscalerReplicas,
		Namespace:     o.config.ClusterAutoscalerNamespace,
		CloudProvider: o.config.ClusterAutoscalerCloudProvider,
	})

	if err := ca.AddToManager(o.manager); err != nil {
		return err
	}

	// Setup MachineAutoscaler controller.
	ma := machineautoscaler.NewReconciler(o.manager, &machineautoscaler.Config{
		Namespace: o.config.ClusterAutoscalerNamespace,
	})

	if err := ma.AddToManager(o.manager); err != nil {
		return err
	}

	return nil
}

// Start starts the operator's controller-manager.
func (o *Operator) Start() error {
	stopCh := signals.SetupSignalHandler()

	// Report status to the CVO.
	go o.status.Report(stopCh)

	return o.manager.Start(stopCh)
}
