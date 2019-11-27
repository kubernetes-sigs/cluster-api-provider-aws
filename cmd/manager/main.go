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

	"github.com/golang/glog"
	mapiv1beta1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/controller/machine"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/client"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/version"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	watchNamespace := flag.String("namespace", "", "Namespace that the controller watches to reconcile machine-api objects. If unspecified, the controller watches for machine-api objects across all namespaces.")

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	flag.Parse()

	if printVersion {
		fmt.Println(version.String)
		os.Exit(0)
	}

	flag.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		glog.Fatalf("Error getting configuration: %v", err)
	}

	// Setup a Manager
	syncPeriod := 10 * time.Minute
	opts := manager.Options{
		SyncPeriod: &syncPeriod,
		// Disable metrics serving
		MetricsBindAddress: "0",
	}
	if *watchNamespace != "" {
		opts.Namespace = *watchNamespace
		klog.Infof("Watching machine-api objects only in namespace %q for reconciliation.", opts.Namespace)
	}

	mgr, err := manager.New(cfg, opts)
	if err != nil {
		glog.Fatalf("Error creating manager: %v", err)
	}

	// Setup Scheme for all resources
	if err := mapiv1beta1.AddToScheme(mgr.GetScheme()); err != nil {
		glog.Fatalf("Error setting up scheme: %v", err)
	}

	machineActuator, err := initActuator(mgr)
	if err != nil {
		glog.Fatalf("Error initializing actuator: %v", err)
	}

	if err := machine.AddWithActuator(mgr, machineActuator); err != nil {
		glog.Fatalf("Error adding actuator: %v", err)
	}

	// Start the Cmd
	err = mgr.Start(signals.SetupSignalHandler())
	if err != nil {
		glog.Fatalf("Error starting manager: %v", err)
	}
}

func initActuator(mgr manager.Manager) (*machineactuator.Actuator, error) {
	codec, err := v1beta1.NewCodec()
	if err != nil {
		return nil, fmt.Errorf("unable to create codec: %v", err)
	}

	params := machineactuator.ActuatorParams{
		Client:           mgr.GetClient(),
		Config:           mgr.GetConfig(),
		AwsClientBuilder: awsclient.NewClient,
		Codec:            codec,
		EventRecorder:    mgr.GetEventRecorderFor("aws-controller"),
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		return nil, fmt.Errorf("could not create AWS machine actuator: %v", err)
	}

	return actuator, nil
}
