/*
Copyright 2018 The Kubernetes Authors.
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

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/controller/machine"
	"github.com/openshift/machine-api-operator/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"
	machinesetcontroller "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machineset"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/version"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// The default durations for the leader electrion operations.
var (
	leaseDuration = 120 * time.Second
	renewDealine  = 110 * time.Second
	retryPeriod   = 20 * time.Second
)

func main() {
	printVersion := flag.Bool(
		"version",
		false,
		"print version and exit",
	)

	metricsAddress := flag.String(
		"metrics-bind-address",
		metrics.DefaultMachineMetricsAddress,
		"Address for hosting metrics",
	)

	watchNamespace := flag.String(
		"namespace",
		"",
		"Namespace that the controller watches to reconcile machine-api objects. If unspecified, the controller watches for machine-api objects across all namespaces.",
	)

	leaderElectResourceNamespace := flag.String(
		"leader-elect-resource-namespace",
		"",
		"The namespace of resource object that is used for locking during leader election. If unspecified and running in cluster, defaults to the service account namespace for the controller. Required for leader-election outside of a cluster.",
	)

	leaderElect := flag.Bool(
		"leader-elect",
		false,
		"Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.",
	)

	leaderElectLeaseDuration := flag.Duration(
		"leader-elect-lease-duration",
		leaseDuration,
		"The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled.",
	)

	healthAddr := flag.String(
		"health-addr",
		":9440",
		"The address for health checking.",
	)

	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()

	if *printVersion {
		fmt.Println(version.String)
		os.Exit(0)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Fatalf("Error getting configuration: %v", err)
	}

	// Setup a Manager
	syncPeriod := 10 * time.Minute
	opts := manager.Options{
		LeaderElection:          *leaderElect,
		LeaderElectionNamespace: *leaderElectResourceNamespace,
		LeaderElectionID:        "cluster-api-provider-aws-leader",
		LeaseDuration:           leaderElectLeaseDuration,
		HealthProbeBindAddress:  *healthAddr,
		SyncPeriod:              &syncPeriod,
		MetricsBindAddress:      *metricsAddress,
		// Slow the default retry and renew election rate to reduce etcd writes at idle: BZ 1858400
		RetryPeriod:   &retryPeriod,
		RenewDeadline: &renewDealine,
	}

	if *watchNamespace != "" {
		opts.Namespace = *watchNamespace
		klog.Infof("Watching machine-api objects only in namespace %q for reconciliation.", opts.Namespace)
	}

	mgr, err := manager.New(cfg, opts)
	if err != nil {
		klog.Fatalf("Error creating manager: %v", err)
	}

	// Setup Scheme for all resources
	if err := machinev1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatalf("Error setting up scheme: %v", err)
	}

	if err := configv1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	if err := corev1.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	configManagedClient, startCache, err := newConfigManagedClient(mgr)
	if err != nil {
		klog.Fatal(err)
	}
	mgr.Add(startCache)

	// Initialize machine actuator.
	machineActuator := machineactuator.NewActuator(machineactuator.ActuatorParams{
		Client:              mgr.GetClient(),
		EventRecorder:       mgr.GetEventRecorderFor("awscontroller"),
		AwsClientBuilder:    awsclient.NewValidatedClient,
		ConfigManagedClient: configManagedClient,
	})

	if err := machine.AddWithActuator(mgr, machineActuator); err != nil {
		klog.Fatalf("Error adding actuator: %v", err)
	}

	ctrl.SetLogger(klogr.New())
	setupLog := ctrl.Log.WithName("setup")
	if err = (&machinesetcontroller.Reconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("MachineSet"),
	}).SetupWithManager(mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MachineSet")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("ping", healthz.Ping); err != nil {
		klog.Fatal(err)
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		klog.Fatal(err)
	}

	// Start the Cmd
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		klog.Fatalf("Error starting manager: %v", err)
	}
}

// newConfigManagedClient returns a controller-runtime client that can be used to access the openshift-config-managed
// namespace.
func newConfigManagedClient(mgr manager.Manager) (runtimeclient.Client, manager.Runnable, error) {
	cacheOpts := cache.Options{
		Scheme:    mgr.GetScheme(),
		Mapper:    mgr.GetRESTMapper(),
		Namespace: awsclient.KubeCloudConfigNamespace,
	}
	cache, err := cache.New(mgr.GetConfig(), cacheOpts)
	if err != nil {
		return nil, nil, err
	}

	clientOpts := runtimeclient.Options{
		Scheme: mgr.GetScheme(),
		Mapper: mgr.GetRESTMapper(),
	}

	cachedClient, err := cluster.DefaultNewClient(cache, config.GetConfigOrDie(), clientOpts)
	if err != nil {
		return nil, nil, err
	}

	return cachedClient, cache, nil
}
