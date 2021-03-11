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

package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/feature"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
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
	_ = controlplanev1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = infrav1exp.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	metricsAddr                string
	enableLeaderElection       bool
	watchNamespace             string
	profilerAddress            string
	eksControlPlaneConcurrency int
	syncPeriod                 time.Duration
	webhookPort                int
	healthAddr                 string
	serviceEndpoints           string

	maxEKSSyncPeriod         = time.Minute * 10
	errMaxSyncPeriodExceeded = errors.New("sync period greater than maximum allowed")
	errEKSInvalidFlags       = errors.New("invalid EKS flag combination")
)

func InitFlags(fs *pflag.FlagSet) {
	fs.StringVar(&metricsAddr, "metrics-addr", ":8080",
		"The address the metric endpoint binds to.")

	fs.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	fs.StringVar(&watchNamespace, "namespace", "",
		"Namespace that the controller watches to reconcile objects. If unspecified, the controller watches for objects across all namespaces.")

	fs.StringVar(&profilerAddress, "profiler-address", "",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")

	fs.IntVar(&eksControlPlaneConcurrency, "ekscontrolplane-concurrency", 10,
		"Number of EKS control planes to process simultaneously")

	fs.DurationVar(&syncPeriod, "sync-period", 10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)")

	fs.IntVar(&webhookPort, "webhook-port", 0,
		"Webhook Server port, disabled by default. When enabled, the manager will only work as webhook server, no reconcilers are installed.")

	fs.StringVar(&serviceEndpoints, "service-endpoints", "",
		"Set custom AWS service endpoins in semi-colon separated format: ${SigningRegion1}:${ServiceID1}=${URL},${ServiceID2}=${URL};${SigningRegion2}...")

	feature.MutableGates.AddFlag(fs)
}

func main() {
	klog.InitFlags(nil)

	rand.Seed(time.Now().UnixNano())
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())

	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	if profilerAddress != "" {
		klog.Infof("Profiler listening for requests at %s", profilerAddress)
		go func() {
			klog.Info(http.ListenAndServe(profilerAddress, nil))
		}()
	}

	if syncPeriod > maxEKSSyncPeriod {
		setupLog.Error(errMaxSyncPeriodExceeded, "sync period exceeded maximum allowed when using EKS", "max-sync-period", maxEKSSyncPeriod)
		os.Exit(1)
	}

	// Parse service endpoints.
	AWSServiceEndpoints, err := endpoints.ParseFlag(serviceEndpoints)
	if err != nil {
		setupLog.Error(err, "unable to parse service endpoints", "controller", "AWSCluster")
		os.Exit(1)
	}

	enableIAM := feature.Gates.Enabled(feature.EKSEnableIAM)
	allowAddRoles := feature.Gates.Enabled(feature.EKSAllowAddRoles)
	setupLog.Info("EKS IAM role creation", "enabled", enableIAM)
	setupLog.Info("EKS IAM additional roles", "enabled", allowAddRoles)
	if allowAddRoles && !enableIAM {
		setupLog.Error(errEKSInvalidFlags, "cannot use EKSAllowAddRoles flag without EKSEnableIAM")
		os.Exit(1)
	}

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-aws-controller"
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "eks-controlplane-manager-leader-election-capa",
		SyncPeriod:             &syncPeriod,
		Namespace:              watchNamespace,
		EventBroadcaster:       broadcaster,
		Port:                   webhookPort,
		HealthProbeBindAddress: healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("aws-controller"))

	setupLog.V(1).Info(fmt.Sprintf("%+v\n", feature.Gates))

	setupReconcilers(mgr, enableIAM, allowAddRoles, AWSServiceEndpoints)
	setupWebhooks(mgr)

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

func setupReconcilers(mgr ctrl.Manager, enableIAM bool, allowAddRoles bool, serviceEndpoints []scope.ServiceEndpoint) {
	if webhookPort != 0 {
		return
	}

	if err := (&controllers.AWSManagedControlPlaneReconciler{
		Client:               mgr.GetClient(),
		Log:                  ctrl.Log.WithName("controllers").WithName("AWSManagedControlPlane"),
		EnableIAM:            enableIAM,
		AllowAdditionalRoles: allowAddRoles,
		Endpoints:            serviceEndpoints,
	}).SetupWithManager(mgr, concurrency(eksControlPlaneConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AWSManagedControlPlane")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if webhookPort == 0 {
		return
	}

	if err := (&controlplanev1.AWSManagedControlPlane{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSManagedControlPlane")
		os.Exit(1)
	}
}

func concurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}
