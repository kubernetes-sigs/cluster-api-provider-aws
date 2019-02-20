package main

import (
	"flag"

	"github.com/golang/glog"
	mapiv1beta1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	namespace = "openshift-machine-api"
)

func init() {
	if err := mapiv1beta1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}
}

type testConfig struct {
	client client.Client
}

func newClient() (client.Client, error) {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{})

}

func main() {
	flag.Parse()
	if err := runSuite(); err != nil {
		glog.Fatal(err)
	}
}

func runSuite() error {

	client, err := newClient()
	if err != nil {
		return err
	}
	testConfig := &testConfig{
		client: client,
	}
	glog.Info("RUN: ExpectProviderAvailable")
	if err := testConfig.ExpectProviderAvailable(); err != nil {
		glog.Errorf("FAIL: ExpectProviderAvailable: %v", err)
		return err
	}
	glog.Info("PASS: ExpectProviderAvailable")

	glog.Info("RUN: ExpectAllMachinesLinkedToANode")
	if err := testConfig.ExpectAllMachinesLinkedToANode(); err != nil {
		glog.Errorf("FAIL: ExpectAllMachinesLinkedToANode: %v", err)
		return err
	}
	glog.Info("PASS: ExpectAllMachinesLinkedToANode")

	glog.Info("RUN: ExpectNewNodeWhenDeletingMachine")
	if err := testConfig.ExpectNewNodeWhenDeletingMachine(); err != nil {
		glog.Errorf("FAIL: ExpectNewNodeWhenDeletingMachine: %v", err)
		return err
	}
	glog.Info("PASS: ExpectNewNodeWhenDeletingMachine")

	glog.Info("RUN: ExpectNodeToBeDrainedBeforeDeletingMachine")
	if err := testConfig.ExpectNodeToBeDrainedBeforeDeletingMachine(); err != nil {
		glog.Errorf("FAIL: ExpectNodeToBeDrainedBeforeDeletingMachine: %v", err)
		return err
	}
	glog.Info("PASS: ExpectNodeToBeDrainedBeforeDeletingMachine")

	glog.Info("RUN: ExpectNodeToBeDrainedBeforeMachineIsDeleted")
	if err := testConfig.ExpectNodeToBeDrainedBeforeMachineIsDeleted(); err != nil {
		glog.Errorf("FAIL: ExpectNodeToBeDrainedBeforeMachineIsDeleted: %v", err)
		return err
	}
	glog.Info("PASS: ExpectNodeToBeDrainedBeforeMachineIsDeleted")

	return nil
}
