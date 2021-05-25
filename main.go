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
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"k8s.io/klog/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clusterv1exp "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	infrav1alpha2 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/controllers"
	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1alpha3exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	controllersexp "sigs.k8s.io/cluster-api-provider-aws/exp/controllers"
	"sigs.k8s.io/cluster-api-provider-aws/exp/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/feature"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/version"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = infrav1alpha2.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	_ = infrav1alpha3exp.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = controlplanev1.AddToScheme(scheme)
	_ = clusterv1exp.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	metricsAddr              string
	enableLeaderElection     bool
	leaderElectionNamespace  string
	watchNamespace           string
	profilerAddress          string
	awsClusterConcurrency    int
	instanceStateConcurrency int
	awsMachineConcurrency    int
	syncPeriod               time.Duration
	webhookPort              int
	healthAddr               string
	serviceEndpoints         string
)

func main() {
	klog.InitFlags(nil)

	rand.Seed(time.Now().UnixNano())
	initFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())

	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	if profilerAddress != "" {
		setupLog.Info("Profiler listening for requests", "profiler-address", profilerAddress)
		go func() {
			setupLog.Error(http.ListenAndServe(profilerAddress, nil), "listen and serve error")
		}()
	}

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-aws-controller"
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      metricsAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "controller-leader-election-capa",
		LeaderElectionNamespace: leaderElectionNamespace,
		SyncPeriod:              &syncPeriod,
		Namespace:               watchNamespace,
		EventBroadcaster:        broadcaster,
		Port:                    webhookPort,
		HealthProbeBindAddress:  healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("aws-controller"))

	setupLog.V(1).Info(fmt.Sprintf("feature gates: %+v\n", feature.Gates))

	// Parse service endpoints.
	AWSServiceEndpoints, err := endpoints.ParseFlag(serviceEndpoints)
	if err != nil {
		setupLog.Error(err, "unable to parse service endpoints", "controller", "AWSCluster")
		os.Exit(1)
	}

	if webhookPort == 0 {
		if err = (&controllers.AWSMachineReconciler{
			Client:    mgr.GetClient(),
			Log:       ctrl.Log.WithName("controllers").WithName("AWSMachine"),
			Recorder:  mgr.GetEventRecorderFor("awsmachine-controller"),
			Endpoints: AWSServiceEndpoints,
		}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: awsMachineConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSMachine")
			os.Exit(1)
		}
		if err = (&controllers.AWSClusterReconciler{
			Client:    mgr.GetClient(),
			Log:       ctrl.Log.WithName("controllers").WithName("AWSCluster"),
			Recorder:  mgr.GetEventRecorderFor("awscluster-controller"),
			Endpoints: AWSServiceEndpoints,
		}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSCluster")
			os.Exit(1)
		}

		if feature.Gates.Enabled(feature.EKS) {
			setupLog.Info("enabling EKS controllers")

			enableIAM := feature.Gates.Enabled(feature.EKSEnableIAM)

			if err = (&controllersexp.AWSManagedMachinePoolReconciler{
				Client:    mgr.GetClient(),
				Log:       ctrl.Log.WithName("controllers").WithName("AWSManagedMachinePool"),
				Recorder:  mgr.GetEventRecorderFor("awsmanagedmachinepool-reconciler"),
				EnableIAM: enableIAM,
				Endpoints: AWSServiceEndpoints,
			}).SetupWithManager(mgr, controller.Options{}); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSManagedMachinePool")
				os.Exit(1)
			}
			if err = (&controllersexp.AWSManagedClusterReconciler{
				Client:   mgr.GetClient(),
				Log:      ctrl.Log.WithName("controllers").WithName("AWSManagedCluster"),
				Recorder: mgr.GetEventRecorderFor("awsmanagedcluster-reconciler"),
			}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency}); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSManagedCluster")
			}
		}
		if feature.Gates.Enabled(feature.MachinePool) {
			if err = (&controllersexp.AWSMachinePoolReconciler{
				Client:   mgr.GetClient(),
				Log:      ctrl.Log.WithName("controllers").WithName("AWSMachinePool"),
				Recorder: mgr.GetEventRecorderFor("awsmachinepool-controller"),
			}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSMachinePool")
				os.Exit(1)
			}
		}

		if feature.Gates.Enabled(feature.EventBridgeInstanceState) {
			setupLog.Info("EventBridge notifications enabled. enabling AWSInstanceStateController")
			if err = (&instancestate.AwsInstanceStateReconciler{
				Client:    mgr.GetClient(),
				Log:       ctrl.Log.WithName("controllers").WithName("AWSInstanceStateController"),
				Endpoints: AWSServiceEndpoints,
			}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: instanceStateConcurrency}); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSInstanceStateController")
				os.Exit(1)
			}
		}
	} else {
		if err = (&infrav1alpha3.AWSMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachineTemplate")
			os.Exit(1)
		}
		if err = (&infrav1alpha3.AWSMachineTemplateList{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachineTemplateList")
			os.Exit(1)
		}
		if err = (&infrav1alpha3.AWSCluster{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSCluster")
			os.Exit(1)
		}
		if err = (&infrav1alpha3.AWSMachine{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachine")
			os.Exit(1)
		}
		if err = (&infrav1alpha3.AWSMachineList{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachineList")
			os.Exit(1)
		}
		if err = (&infrav1alpha3.AWSClusterList{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSClusterList")
			os.Exit(1)
		}
		if feature.Gates.Enabled(feature.EKS) {
			setupLog.Info("enabling EKS webhooks")
			if err = (&infrav1alpha3exp.AWSManagedMachinePool{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "AWSManagedMachinePool")
				os.Exit(1)
			}
		}
		if feature.Gates.Enabled(feature.MachinePool) {
			setupLog.Info("enabling webhook for AWSMachinePool")
			if err = (&infrav1alpha3exp.AWSMachinePool{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachinePool")
				os.Exit(1)
			}
		}
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddReadyzCheck("ping", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&metricsAddr,
		"metrics-addr",
		":8080",
		"The address the metric endpoint binds to.",
	)

	fs.BoolVar(
		&enableLeaderElection,
		"enable-leader-election",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
	)

	fs.StringVar(
		&watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	fs.StringVar(
		&leaderElectionNamespace,
		"leader-election-namespace",
		"",
		"Namespace that the controller performs leader election in. If unspecified, the controller will discover which namespace it is running in.",
	)

	fs.StringVar(
		&profilerAddress,
		"profiler-address",
		"",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)",
	)

	fs.IntVar(&awsClusterConcurrency,
		"awscluster-concurrency",
		5,
		"Number of AWSClusters to process simultaneously",
	)

	fs.IntVar(&instanceStateConcurrency,
		"instance-state-concurrency",
		5,
		"Number of concurrent watches for instance state changes",
	)

	fs.IntVar(&awsMachineConcurrency,
		"awsmachine-concurrency",
		10,
		"Number of AWSMachines to process simultaneously",
	)

	fs.DurationVar(&syncPeriod,
		"sync-period",
		10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)",
	)

	fs.IntVar(&webhookPort,
		"webhook-port",
		0,
		"Webhook Server port, disabled by default. When enabled, the manager will only work as webhook server, no reconcilers are installed.",
	)

	fs.StringVar(&healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	fs.StringVar(&serviceEndpoints,
		"service-endpoints",
		"",
		"Set custom AWS service endpoins in semi-colon separated format: ${SigningRegion1}:${ServiceID1}=${URL},${ServiceID2}=${URL};${SigningRegion2}...",
	)

	feature.MutableGates.AddFlag(fs)
}
