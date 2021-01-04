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
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/termination"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/version"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	var printVersion bool
	flag.BoolVar(&printVersion, "version", false, "print version and exit")

	klog.InitFlags(nil)
	logger := klogr.New()

	pollIntervalSeconds := flag.Int64("poll-interval-seconds", 5, "interval in seconds at which termination notice endpoint should be checked (Default: 5)")
	nodeName := flag.String("node-name", "", "name of the node that the termination handler is running on")
	namespace := flag.String("namespace", "", "namespace that the machine for the node should live in. If unspecified, the look for machines across all namespaces.")
	flag.Set("logtostderr", "true")
	flag.Parse()

	if printVersion {
		fmt.Println(version.String)
		os.Exit(0)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(err, "Error getting configuration")
		return
	}

	// Get the poll interval as a duration from the `poll-interval-seconds` flag
	pollInterval := time.Duration(*pollIntervalSeconds) * time.Second

	// Construct a termination handler
	handler, err := termination.NewHandler(logger, cfg, pollInterval, *namespace, *nodeName)
	if err != nil {
		logger.Error(err, "Error constructing termination handler")
		return
	}

	// Start the termination handler
	if err := handler.Run(ctrl.SetupSignalHandler().Done()); err != nil {
		logger.Error(err, "Error starting termination handler")
		return
	}
}
