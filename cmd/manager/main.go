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

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterapis "sigs.k8s.io/cluster-api/pkg/apis"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	capicluster "sigs.k8s.io/cluster-api/pkg/controller/cluster"
	capimachine "sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	watchNamespace := flag.String("namespace", "",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")
	flag.Parse()

	cfg := config.GetConfigOrDie()

	// Setup a Manager
	syncPeriod := 10 * time.Minute
	opts := manager.Options{
		SyncPeriod: &syncPeriod,
	}

	if *watchNamespace != "" {
		opts.Namespace = *watchNamespace
		klog.Infof("Watching cluster-api objects only in namespace %q for reconciliation.", opts.Namespace)
	}

	mgr, err := manager.New(cfg, opts)
	if err != nil {
		klog.Fatalf("Failed to set up overall controller manager: %v", err)
	}

	cs, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Failed to create client from configuration: %v", err)
	}

	coreClient, err := corev1.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Failed to create corev1 client from configuration: %v", err)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetRecorder("aws-controller"))

	// Initialize cluster actuator.
	clusterActuator := cluster.NewActuator(cluster.ActuatorParams{
		Client:         cs.ClusterV1alpha1(),
		LoggingContext: "[cluster-actuator]",
	})

	// Initialize machine actuator.
	machineActuator := machine.NewActuator(machine.ActuatorParams{
		CoreClient:     coreClient,
		ClusterClient:  cs.ClusterV1alpha1(),
		LoggingContext: "[machine-actuator]",
	})

	// Register our cluster deployer (the interface is in clusterctl and we define the Deployer interface on the actuator)
	common.RegisterClusterProvisioner("aws", clusterActuator)

	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Fatal(err)
	}

	capimachine.AddWithActuator(mgr, machineActuator)
	capicluster.AddWithActuator(mgr, clusterActuator)

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Fatalf("Failed to run manager: %v", err)
	}
}
