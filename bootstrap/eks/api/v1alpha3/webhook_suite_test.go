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

package v1alpha3

import (
	"fmt"
	"os"
	"path"
	"testing"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	bootstrapv1alpha4 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/test/helpers"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	testEnv *helpers.TestEnvironment
	ctx     = ctrl.SetupSignalHandler()
)

func TestMain(m *testing.M) {
	code := 0
	defer func() { os.Exit(code) }()
	setup()
	defer teardown()
	code = m.Run()
}

func setup() {
	utilruntime.Must(AddToScheme(scheme.Scheme))
	utilruntime.Must(bootstrapv1alpha4.AddToScheme(scheme.Scheme))

	testEnvConfig := helpers.NewTestEnvironmentConfiguration([]string{
		path.Join("bootstrap", "eks", "config", "crd", "bases"),
	},
	).WithWebhookConfiguration("unmanaged", path.Join("bootstrap", "eks", "config", "webhook", "manifests.yaml"))
	var err error
	testEnv, err = testEnvConfig.Build()
	if err != nil {
		panic(err)
	}
	if err := (&bootstrapv1alpha4.EKSConfig{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSCluster webhook: %v", err))
	}
	if err := (&bootstrapv1alpha4.EKSConfigTemplate{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup AWSCluster webhook: %v", err))
	}
	go func() {
		fmt.Println("Starting the manager")
		if err := testEnv.StartManager(ctx); err != nil {
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
