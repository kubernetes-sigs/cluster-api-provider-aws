/*
Copyright 2021 The Kubernetes Authors.

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

package controlleridentitycreator

import (
	"fmt"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/test/helpers"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"sigs.k8s.io/controller-runtime/pkg/log"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	testEnv *helpers.TestEnvironment
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

func TestMain(m *testing.M) {
	setup()
	defer func() {
		teardown()
	}()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	utilruntime.Must(infrav1.AddToScheme(scheme.Scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme.Scheme))
	utilruntime.Must(controlplanev1.AddToScheme(scheme.Scheme))

	testEnvConfig := helpers.NewTestEnvironmentConfiguration([]string{
		path.Join("config", "crd", "bases"),
		path.Join("controlplane", "eks", "config", "crd", "bases"),
	},
	).WithWebhookConfiguration("unmanaged", path.Join("config", "webhook", "manifests.yaml"))
	var err error
	testEnv, err = testEnvConfig.Build()
	if err != nil {
		panic(err)
	}
	if err := (&infrav1.AWSCluster{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSCluster webhook: %v", err))
	}
	if err := (&infrav1.AWSClusterControllerIdentity{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSClusterControllerIdentity webhook: %v", err))
	}
	if err = (&infrav1.AWSClusterRoleIdentity{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSClusterRoleIdentity webhook: %v", err))
	}
	if err := (&controlplanev1.AWSManagedControlPlane{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSManagedControlPlane webhook: %v", err))
	}

	err = (&AWSControllerIdentityReconciler{
		Client: testEnv,
		Log:    log.Log,
	}).SetupWithManager(testEnv.Manager, controller.Options{})
	if err != nil {
		panic(fmt.Sprintf("Failed to add AWSControllerIdentityReconciler to the envtest manager: %v", err))
	}

	go func() {
		fmt.Println("Starting the manager")
		if err := testEnv.StartManager(); err != nil {
			panic(fmt.Sprintf("Failed to start the envtest manager: %v", err))
		}
	}()
	testEnv.WaitForWebhooks()
}

func teardown() {
	if err := testEnv.Stop(); err != nil {
		panic(fmt.Sprintf("Failed to stop envtest: %v", err))
	}
}
