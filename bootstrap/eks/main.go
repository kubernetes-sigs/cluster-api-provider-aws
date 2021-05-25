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
	"context"
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
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	bootstrapv1alpha3 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
	bootstrapv1alpha4 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	bootstrapv1controllers "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/controllers"
	controlplanev1alpha3 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	controlplanev1alpha4 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	expinfrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	expinfrav1alpha4 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/version"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	klog.InitFlags(nil)

	_ = clientgoscheme.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = bootstrapv1alpha3.AddToScheme(scheme)
	_ = bootstrapv1alpha4.AddToScheme(scheme)
	_ = expinfrav1alpha3.AddToScheme(scheme)
	_ = expinfrav1alpha4.AddToScheme(scheme)
	_ = expclusterv1.AddToScheme(scheme)
	_ = controlplanev1alpha3.AddToScheme(scheme)
	_ = controlplanev1alpha4.AddToScheme(scheme)

	// +kubebuilder:scaffold:scheme
}

var (
	metricsAddr                 string
	enableLeaderElection        bool
	leaderElectionLeaseDuration time.Duration
	leaderElectionRenewDeadline time.Duration
	leaderElectionRetryPeriod   time.Duration
	watchNamespace              string
	watchFilterValue            string
	profilerAddress             string
	eksConfigConcurrency        int
	syncPeriod                  time.Duration
	webhookPort                 int
	webhookCertDir              string
)

// InitFlags initializes this manager's flags.
func InitFlags(fs *pflag.FlagSet) {
	fs.StringVar(&metricsAddr, "metrics-bind-addr", ":8080",
		"The address the metric endpoint binds to.")

	fs.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	fs.DurationVar(&leaderElectionLeaseDuration, "leader-elect-lease-duration", 15*time.Second,
		"Interval at which non-leader candidates will wait to force acquire leadership (duration string)")

	fs.DurationVar(&leaderElectionRenewDeadline, "leader-elect-renew-deadline", 10*time.Second,
		"Duration that the leading controller manager will retry refreshing leadership before giving up (duration string)")

	fs.DurationVar(&leaderElectionRetryPeriod, "leader-elect-retry-period", 2*time.Second,
		"Duration the LeaderElector clients should wait between tries of actions (duration string)")

	fs.StringVar(&watchNamespace, "namespace", "",
		"Namespace that the controller watches to reconcile objects. If unspecified, the controller watches for objects across all namespaces.")

	fs.StringVar(&profilerAddress, "profiler-address", "",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")

	fs.IntVar(&eksConfigConcurrency, "eksconfig-concurrency", 10,
		"Number of EKS configs to process simultaneously")

	fs.DurationVar(&syncPeriod, "sync-period", 10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)")

	fs.IntVar(&webhookPort, "webhook-port", 9443,
		"Webhook Server port, disabled by default. When enabled, the manager will only work as webhook server, no reconcilers are installed.")

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")

	fs.StringVar(
		&watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel),
	)

	feature.MutableGates.AddFlag(fs)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	InitFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())

	if profilerAddress != "" {
		klog.Infof("Profiler listening for requests at %s", profilerAddress)
		go func() {
			klog.Info(http.ListenAndServe(profilerAddress, nil))
		}()
	}

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-aws-controller"
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "eks-bootstrap-manager-leader-elect-capa",
		LeaseDuration:      &leaderElectionLeaseDuration,
		RenewDeadline:      &leaderElectionRenewDeadline,
		RetryPeriod:        &leaderElectionRetryPeriod,
		Namespace:          watchNamespace,
		SyncPeriod:         &syncPeriod,
		Port:               webhookPort,
		CertDir:            webhookCertDir,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	ctx := ctrl.SetupSignalHandler()
	setupReconcilers(ctx, mgr)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupReconcilers(ctx context.Context, mgr ctrl.Manager) {
	if err := (&bootstrapv1controllers.EKSConfigReconciler{
		Client:           mgr.GetClient(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, concurrency(eksConfigConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "EKSConfig")
		os.Exit(1)
	}
}

func concurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}
