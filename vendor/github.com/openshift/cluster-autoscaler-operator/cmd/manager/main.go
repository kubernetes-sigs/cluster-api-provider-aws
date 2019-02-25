package main

import (
	"flag"
	"runtime"

	"github.com/golang/glog"
	"github.com/openshift/cluster-autoscaler-operator/pkg/operator"
	"github.com/openshift/cluster-autoscaler-operator/pkg/version"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func printVersion() {
	glog.Infof("Go Version: %s", runtime.Version())
	glog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	glog.Infof("Version: %s", version.String)
}

func main() {
	flag.Parse()
	printVersion()

	config := operator.ConfigFromEnvironment()

	operator, err := operator.New(config)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Starting cluster-autoscaler-operator")
	if err := operator.Start(); err != nil {
		glog.Fatal(err)
	}
}
