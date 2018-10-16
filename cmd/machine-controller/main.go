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
	"os"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/client-go/kubernetes"
	machineactuator "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/cluster-api/pkg/controller/config"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/cluster-api/pkg/controller/sharedinformers"
)

var (
	logLevel string
)

func init() {
	config.ControllerConfig.AddFlags(pflag.CommandLine)
	pflag.CommandLine.StringVar(&logLevel, "log-level", defaultLogLevel, "Log level (debug,info,warn,error,fatal)")
}

const (
	controllerLogName = "awsMachine"
	defaultLogLevel   = "info"
)

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
	pflag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	config, err := controller.GetConfig(config.ControllerConfig.Kubeconfig)
	if err != nil {
		glog.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
	}

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

	logger := log.WithField("controller", controllerLogName)

	params := machineactuator.ActuatorParams{
		ClusterClient:    client,
		KubeClient:       kubeClient,
		AwsClientBuilder: awsclient.NewClient,
		Logger:           logger,
	}

	actuator, err := machineactuator.NewActuator(params)
	if err != nil {
		glog.Fatalf("Could not create AWS machine actuator: %v", err)
	}

	shutdown := make(chan struct{})
	si := sharedinformers.NewSharedInformers(config, shutdown)
	// If this doesn't compile, the code generator probably
	// overwrote the customized NewMachineController function.
	c := machine.NewMachineController(config, si, actuator)
	c.Run(shutdown)
	select {}
}
