package main

import (
	"flag"
	"runtime"

	"github.com/openshift/cluster-autoscaler-operator/pkg/operator"
	"github.com/openshift/cluster-autoscaler-operator/pkg/version"
	"k8s.io/klog"
)

func printVersion() {
	klog.Infof("Go Version: %s", runtime.Version())
	klog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	klog.Infof("Version: %s", version.String)
}

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	printVersion()

	config := operator.ConfigFromEnvironment()

	operator, err := operator.New(config)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("Starting cluster-autoscaler-operator")
	if err := operator.Start(); err != nil {
		klog.Fatal(err)
	}
}
