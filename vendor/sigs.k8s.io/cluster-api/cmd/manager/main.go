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
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/cluster-api/pkg/apis"
	"sigs.k8s.io/cluster-api/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	flag.Set("logtostderr", "true")
	klog.InitFlags(nil)
	watchNamespace := flag.String("namespace", "",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")

	flag.Parse()
	if *watchNamespace != "" {
		klog.Infof("Watching cluster-api objects only in namespace %q for reconciliation", *watchNamespace)
	}

	// Setup controller-runtime logger.
	log.SetLogger(klogr.New())

	// Get a config to talk to the api-server.
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components.
	syncPeriod := 10 * time.Minute
	mgr, err := manager.New(cfg, manager.Options{
		SyncPeriod: &syncPeriod,
		Namespace:  *watchNamespace,
	})

	if err != nil {
		klog.Fatalf("Failed to create new Manager: %v", err)
	}

	klog.Info("Registering Components")

	// Setup Scheme for all resources.
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	// Setup all Controllers.
	if err := controller.AddToManager(mgr); err != nil {
		klog.Fatal(err)
	}

	klog.Info("Starting the Cmd")

	// Start the Cmd
	klog.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
