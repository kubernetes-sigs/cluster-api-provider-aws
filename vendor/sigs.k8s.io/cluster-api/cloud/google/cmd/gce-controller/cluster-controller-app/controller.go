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

package cluster_controller_app

import (
	"os"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/apiserver-builder/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/cloud/google"
	"sigs.k8s.io/cluster-api/cloud/google/cmd/gce-controller/cluster-controller-app/options"
	"sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	clusterapiclientsetscheme "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/scheme"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	"sigs.k8s.io/cluster-api/pkg/controller/config"
	"sigs.k8s.io/cluster-api/pkg/controller/sharedinformers"
)

const (
	gceClusterControllerName = "gce-cluster-controller"
)

func StartClusterController(server *options.ClusterControllerServer, shutdown <-chan struct{}) {
	config, err := controller.GetConfig(server.CommonConfig.Kubeconfig)
	if err != nil {
		glog.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create client for talking to the apiserver: %v", err)
	}

	params := google.ClusterActuatorParams{
		ClusterClient: client.ClusterV1alpha1().Clusters(corev1.NamespaceDefault),
	}
	actuator, err := google.NewClusterActuator(params)
	if err != nil {
		glog.Fatalf("Could not create Google cluster actuator: %v", err)
	}

	si := sharedinformers.NewSharedInformers(config, shutdown)
	// If this doesn't compile, the code generator probably
	// overwrote the customized NewClusterController function.
	c := cluster.NewClusterController(config, si, actuator)
	c.RunAsync(shutdown)

	select {}
}

func RunClusterController(server *options.ClusterControllerServer) error {
	// TODO refactor to remove duplication of code between this and the machine
	// controller server. Issue #262
	kubeConfig, err := controller.GetConfig(server.CommonConfig.Kubeconfig)
	if err != nil {
		glog.Errorf("Could not create Config for talking to the apiserver: %v", err)
		return err
	}

	kubeClientControl, err := kubernetes.NewForConfig(
		rest.AddUserAgent(kubeConfig, "cluster-controller-manager"),
	)
	if err != nil {
		glog.Errorf("Invalid API configuration for kubeconfig-control: %v", err)
		return err
	}

	recorder, err := createRecorder(kubeClientControl)
	if err != nil {
		glog.Errorf("Could not create event recorder : %v", err)
		return err
	}

	// run function will block and never return.
	run := func(stop <-chan struct{}) {
		StartClusterController(server, stop)
	}

	leaderElectConfig := config.GetLeaderElectionConfig()
	if !leaderElectConfig.LeaderElect {
		run(make(<-chan (struct{})))
	}

	// Identity used to distinguish between multiple cluster controller instances.
	id, err := os.Hostname()
	if err != nil {
		return err
	}

	leaderElectionClient := kubernetes.NewForConfigOrDie(rest.AddUserAgent(kubeConfig, "cluster-leader-election"))

	id = id + "-" + string(uuid.NewUUID())
	// Lock required for leader election
	rl, err := resourcelock.New(
		leaderElectConfig.ResourceLock,
		metav1.NamespaceSystem,
		gceClusterControllerName,
		leaderElectionClient.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id + "-" + gceClusterControllerName,
			EventRecorder: recorder,
		})
	if err != nil {
		return err
	}

	// Try and become the leader and start cluster controller loops
	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: leaderElectConfig.LeaseDuration.Duration,
		RenewDeadline: leaderElectConfig.RenewDeadline.Duration,
		RetryPeriod:   leaderElectConfig.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				glog.Fatalf("leaderelection lost")
			},
		},
	})
	panic("unreachable")
}

func createRecorder(kubeClient *kubernetes.Clientset) (record.EventRecorder, error) {

	eventsScheme := runtime.NewScheme()
	if err := corev1.AddToScheme(eventsScheme); err != nil {
		return nil, err
	}
	// We also emit events for our own types
	clusterapiclientsetscheme.AddToScheme(eventsScheme)

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.CoreV1().RESTClient()).Events("")})
	return eventBroadcaster.NewRecorder(eventsScheme, corev1.EventSource{Component: gceClusterControllerName}), nil
}
