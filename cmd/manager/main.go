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
	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/controller"
	clusterapis "sigs.k8s.io/cluster-api/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var (
	logLevel string
)

const (
	defaultLogLevel = "info"
)

func init() {
	pflag.CommandLine.StringVar(&logLevel, "log-level", defaultLogLevel, "Log level (debug,info,warn,error,fatal)")
}

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	//flag.CommandLine.Parse([]string{})
	pflag.Parse()

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	initActuator(mgr)
	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting the Cmd.")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}

func initActuator(m manager.Manager) {
	config := m.GetConfig()

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
	}

	log.SetOutput(os.Stdout)
	if lvl, err := log.ParseLevel(logLevel); err != nil {
		log.Panic(err)
	} else {
		log.SetLevel(lvl)
	}

	logger := log.WithField("controller", "awsMachine")

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		glog.Fatal(err)
	}

	params := machineactuator.ActuatorParams{
		Client:           m.GetClient(),
		KubeClient:       kubeClient,
		AwsClientBuilder: awsclient.NewClient,
		Logger:           logger,
		Codec:            codec,
	}

	machineactuator.MachineActuator, err = machineactuator.NewActuator(params)
	if err != nil {
		glog.Fatalf("Could not create AWS machine actuator: %v", err)
	}
}
