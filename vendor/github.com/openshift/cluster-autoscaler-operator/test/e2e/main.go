package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	namespace = "openshift-machine-api"
	caName    = "default"
)

var F *Framework

type Framework struct {
	Client client.Client
}

func newClient() error {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client, err := client.New(cfg, client.Options{})
	if err != nil {
		return err
	}

	F = &Framework{Client: client}

	return nil
}

func main() {
	flag.Parse()

	if err := apis.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}

	if err := newClient(); err != nil {
		glog.Fatal(err)
	}

	if err := runSuite(); err != nil {
		glog.Fatal(err)
	}
}

func runSuite() error {
	if err := ExpectOperatorAvailable(); err != nil {
		glog.Errorf("FAIL: ExpectOperatorAvailable: %v", err)
		return err
	}
	glog.Info("PASS: ExpectOperatorAvailable")

	if err := CreateClusterAutoscaler(); err != nil {
		glog.Errorf("FAIL: CreateClusterAutoscaler: %v", err)
		return err
	}
	glog.Info("PASS: CreateClusterAutoscaler")

	if err := ExpectClusterAutoscalerAvailable(); err != nil {
		glog.Errorf("FAIL: ExpectClusterAutoscalerAvailable: %v", err)
		return err
	}
	glog.Info("PASS: ExpectClusterAutoscalerAvailable")

	return nil
}
