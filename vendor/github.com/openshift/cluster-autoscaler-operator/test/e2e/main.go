package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-autoscaler-operator/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	caName = "default"
)

var focus string
var namespace string

func init() {
	flag.StringVar(&focus, "focus", "[openshift]", "If set, run only tests containing focus string. E.g. [k8s]")
	flag.StringVar(&namespace, "namespace", "openshift-machine-api", "cluster-autoscaler-operator namespace")
}

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

	if err := osconfigv1.AddToScheme(scheme.Scheme); err != nil {
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

	expectations := []struct {
		expect func() error
		name   string
	}{
		{
			expect: ExpectOperatorAvailable,
			name:   "[k8s][openshift] Expect operator to be available",
		},
		{
			expect: ExpectClusterOperatorStatusAvailable,
			name:   "[openshift] Expect Cluster Operator status to be available",
		},
		{
			expect: CreateClusterAutoscaler,
			name:   "[openshift] Create Cluster Autoscaler resource",
		},
		{
			expect: ExpectClusterAutoscalerAvailable,
			name:   "[openshift] Expect Cluster Autoscaler available",
		},
		{
			expect: ExpectToScaleUpAndDown,
			name:   "[k8s] Expect to scale up and down",
		},
	}

	for _, tc := range expectations {
		if strings.HasPrefix(tc.name, focus) {
			if err := tc.expect(); err != nil {
				glog.Errorf("FAIL: %v: %v", tc.name, err)
				return err
			}
			glog.Infof("PASS: %v", tc.name)
		} else {
			glog.Infof("SKIPPING: %v", tc.name)
		}
	}

	return nil
}
