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

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	clusterapis "sigs.k8s.io/cluster-api/pkg/apis"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	capicluster "sigs.k8s.io/cluster-api/pkg/controller/cluster"
	capimachine "sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/providerconfig"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	elbsvc "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
)

func main() {
	cfg := config.GetConfigOrDie()

	flag.Parse()
	log := logf.Log.WithName("example-controller")
	logf.SetLogger(logf.ZapLogger(false))
	entryLog := log.WithName("entrypoint")

	// Setup a Manager
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		panic(err)
	}
	cs, err := clientset.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	clusterActuator, _ := cluster.NewActuator(cluster.ActuatorParams{
		Codec:          codec,
		ClustersGetter: cs.ClusterV1alpha1(),
	})

	sess := session.Must(session.NewSession())
	ec2client := ec2.New(sess)
	elbclient := elb.New(sess)

	machineActuator, _ := machine.NewActuator(machine.ActuatorParams{
		Codec:          codec,
		MachinesGetter: cs.ClusterV1alpha1(),
		EC2Service:     ec2svc.NewService(ec2client),
		ELBService:     elbsvc.NewService(elbclient),
	})

	// Register our cluster deployer (the interface is in clusterctl and we define the Deployer interface on the actuator)
	common.RegisterClusterProvisioner("aws", clusterActuator)
	if err := providerconfig.AddToScheme(mgr.GetScheme()); err != nil {
		panic(err)
	}
	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		panic(err)
	}
	capimachine.AddWithActuator(mgr, machineActuator)
	capicluster.AddWithActuator(mgr, clusterActuator)

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
