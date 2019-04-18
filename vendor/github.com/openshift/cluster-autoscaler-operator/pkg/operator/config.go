package operator

import (
	"os"
	"strconv"

	"k8s.io/klog"
)

const (
	// DefaultWatchNamespace is the default namespace the operator
	// will watch for instances of its custom resources.
	DefaultWatchNamespace = "openshift-machine-api"

	// DefaultLeaderElection controls whether the default
	// configuration performs leader-election on startup.
	DefaultLeaderElection = true

	// DefaultLeaderElectionNamespace is the default namespace in
	// which the leader-election ConfigMap will be created.
	DefaultLeaderElectionNamespace = "openshift-machine-api"

	// DefaultLeaderElectionID is the default name for the ConfigMap
	// used for leader-election.
	DefaultLeaderElectionID = "cluster-autoscaler-operator-leader"

	// DefaultClusterAutoscalerNamespace is the default namespace for
	// cluster-autoscaler deployments.
	DefaultClusterAutoscalerNamespace = "openshift-machine-api"

	// DefaultClusterAutoscalerName is the default ClusterAutoscaler
	// object watched by the operator.
	DefaultClusterAutoscalerName = "default"

	// DefaultClusterAutoscalerImage is the default image used in
	// ClusterAutoscaler deployments.
	DefaultClusterAutoscalerImage = "quay.io/openshift/origin-cluster-autoscaler:v4.0"

	// DefaultClusterAutoscalerReplicas is the default number of
	// replicas in ClusterAutoscaler deployments.
	DefaultClusterAutoscalerReplicas = 1

	// DefaultClusterAutoscalerCloudProvider is the default name for
	// the CloudProvider beeing used.
	DefaultClusterAutoscalerCloudProvider = "openshift-machine-api"

	// DefaultClusterAutoscalerVerbosity is the default logging
	// verbosity level for ClusterAutoscaler deployments.
	DefaultClusterAutoscalerVerbosity = 1
)

// Config represents the runtime configuration for the operator.
type Config struct {
	// ReleaseVersion is the version the operator is expected
	// to report once it has reached level.
	ReleaseVersion string

	// WatchNamespace is the namespace the operator will watch for
	// ClusterAutoscaler and MachineAutoscaler instances.
	WatchNamespace string

	// LeaderElection indicates whether to perform leader-election.
	LeaderElection bool

	// LeaderElectionNamespace is the namespace in which the
	// leader-election ConfigMap will be created.
	LeaderElectionNamespace string

	// LeaderElectionID is the name of the leader-election ConfigMap.
	LeaderElectionID string

	// ClusterAutoscalerNamespace is the namespace in which
	// cluster-autoscaler deployments will be created.
	ClusterAutoscalerNamespace string

	// ClusterAutoscalerName is the name of the ClusterAutoscaler
	// resource that will be watched by the operator.
	ClusterAutoscalerName string

	// ClusterAutoscalerImage is the image to be used in
	// ClusterAutoscaler deployments.
	ClusterAutoscalerImage string

	// ClusterAutoscalerReplicas is the number of replicas to be
	// configured in ClusterAutoscaler deployments.
	ClusterAutoscalerReplicas int32

	// ClusterAutoscalerCloudProvider is the name for the
	// CloudProvider beeing used.
	ClusterAutoscalerCloudProvider string

	// ClusterAutoscalerVerbosity is the logging verbosity level for
	// ClusterAutoscaler deployments.
	ClusterAutoscalerVerbosity int

	// ClusterAutoscalerExtraArgs is a string of additional arguments
	// passed to all ClusterAutoscaler deployments.
	//
	// This is not exposed in the CRD.  It is only configurable via
	// environment variable, and in a normal OpenShift install the CVO
	// will remove it if set manually.  It is only for development and
	// debugging purposes.
	ClusterAutoscalerExtraArgs string
}

// NewConfig returns a new Config object with defaults set.
func NewConfig() *Config {
	return &Config{
		WatchNamespace:                 DefaultWatchNamespace,
		LeaderElection:                 DefaultLeaderElection,
		LeaderElectionNamespace:        DefaultLeaderElectionNamespace,
		LeaderElectionID:               DefaultLeaderElectionID,
		ClusterAutoscalerNamespace:     DefaultClusterAutoscalerNamespace,
		ClusterAutoscalerName:          DefaultClusterAutoscalerName,
		ClusterAutoscalerImage:         DefaultClusterAutoscalerImage,
		ClusterAutoscalerReplicas:      DefaultClusterAutoscalerReplicas,
		ClusterAutoscalerCloudProvider: DefaultClusterAutoscalerCloudProvider,
		ClusterAutoscalerVerbosity:     DefaultClusterAutoscalerVerbosity,
	}
}

// ConfigFromEnvironment returns a new Config object with defaults
// overridden by environment variables when set.
func ConfigFromEnvironment() *Config {
	config := NewConfig()

	if releaseVersion, ok := os.LookupEnv("RELEASE_VERSION"); ok {
		config.ReleaseVersion = releaseVersion
	}

	if watchNamespace, ok := os.LookupEnv("WATCH_NAMESPACE"); ok {
		config.WatchNamespace = watchNamespace
	}

	if leaderElection, ok := os.LookupEnv("LEADER_ELECTION"); ok {
		le, err := strconv.ParseBool(leaderElection)
		if err != nil {
			le = DefaultLeaderElection
			klog.Errorf("Error parsing LEADER_ELECTION environment variable: %v", err)
		}

		config.LeaderElection = le
	}

	if leNamespace, ok := os.LookupEnv("LEADER_ELECTION_NAMESPACE"); ok {
		config.LeaderElectionNamespace = leNamespace
	}

	if leID, ok := os.LookupEnv("LEADER_ELECTION_ID"); ok {
		config.LeaderElectionID = leID
	}

	if caName, ok := os.LookupEnv("CLUSTER_AUTOSCALER_NAME"); ok {
		config.ClusterAutoscalerName = caName
	}

	if caImage, ok := os.LookupEnv("CLUSTER_AUTOSCALER_IMAGE"); ok {
		config.ClusterAutoscalerImage = caImage
	}

	if cloudProvider, ok := os.LookupEnv("CLUSTER_AUTOSCALER_CLOUD_PROVIDER"); ok {
		config.ClusterAutoscalerCloudProvider = cloudProvider
	}

	if caNamespace, ok := os.LookupEnv("CLUSTER_AUTOSCALER_NAMESPACE"); ok {
		config.ClusterAutoscalerNamespace = caNamespace
	}

	if caVerbosity, ok := os.LookupEnv("CLUSTER_AUTOSCALER_VERBOSITY"); ok {
		v, err := strconv.Atoi(caVerbosity)
		if err != nil {
			v = DefaultClusterAutoscalerVerbosity
			klog.Errorf("Error parsing CLUSTER_AUTOSCALER_VERBOSITY environment variable: %v", err)
		}

		config.ClusterAutoscalerVerbosity = v
	}

	if caExtraArgs, ok := os.LookupEnv("CLUSTER_AUTOSCALER_EXTRA_ARGS"); ok {
		config.ClusterAutoscalerExtraArgs = caExtraArgs
	}

	return config
}
