// +build e2e

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

package shared

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/yaml"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// NOTE: these are temporary local functions as the CAPI e2e framework assumes kubeadm
// in a number of places. These can be removed when the following is fixed:
// https://github.com/kubernetes-sigs/cluster-api/issues/3983

type InitManagementClusterAndWatchControllerLogsInput struct {
	ClusterProxy             framework.ClusterProxy
	ClusterctlConfigPath     string
	InfrastructureProviders  []string
	BootstrapProviders       []string
	ControlPlaneProviders    []string
	LogFolder                string
	DisableMetricsCollection bool
}

func InitManagementClusterAndWatchControllerLogs(ctx context.Context, input InitManagementClusterAndWatchControllerLogsInput, intervals ...interface{}) {
	Expect(ctx).NotTo(BeNil(), "ctx is required for InitManagementClusterAndWatchControllerLogs")
	Expect(input.ClusterProxy).ToNot(BeNil(), "Invalid argument. input.ClusterProxy can't be nil when calling InitManagementClusterAndWatchControllerLogs")
	Expect(input.ClusterctlConfigPath).To(BeAnExistingFile(), "Invalid argument. input.ClusterctlConfigPath must be an existing file when calling InitManagementClusterAndWatchControllerLogs")
	Expect(input.InfrastructureProviders).ToNot(BeEmpty(), "Invalid argument. input.InfrastructureProviders can't be empty when calling InitManagementClusterAndWatchControllerLogs")
	Expect(os.MkdirAll(input.LogFolder, 0755)).To(Succeed(), "Invalid argument. input.LogFolder can't be created for InitManagementClusterAndWatchControllerLogs")

	client := input.ClusterProxy.GetClient()
	controllersDeployments := framework.GetControllerDeployments(context.TODO(), framework.GetControllerDeploymentsInput{
		Lister: client,
	})
	if len(controllersDeployments) == 0 {
		clusterctl.Init(context.TODO(), clusterctl.InitInput{
			// pass reference to the management cluster hosting this test
			KubeconfigPath: input.ClusterProxy.GetKubeconfigPath(),
			// pass the clusterctl generate file that points to the local provider repository created for this test
			ClusterctlConfigPath: input.ClusterctlConfigPath,
			// setup the desired list of providers for a single-tenant management cluster
			CoreProvider:            config.ClusterAPIProviderName,
			BootstrapProviders:      input.BootstrapProviders,
			ControlPlaneProviders:   input.ControlPlaneProviders,
			InfrastructureProviders: input.InfrastructureProviders,
			// setup clusterctl logs folder
			LogFolder: input.LogFolder,
		})
	}

	By("Waiting for provider controllers to be running")
	controllersDeployments = framework.GetControllerDeployments(context.TODO(), framework.GetControllerDeploymentsInput{
		Lister: client,
	})
	Expect(controllersDeployments).ToNot(BeEmpty(), "The list of controller deployments should not be empty")
	for _, deployment := range controllersDeployments {
		framework.WaitForDeploymentsAvailable(context.TODO(), framework.WaitForDeploymentsAvailableInput{
			Getter:     client,
			Deployment: deployment,
		}, intervals...)

		// Start streaming logs from all controller providers
		framework.WatchDeploymentLogs(context.TODO(), framework.WatchDeploymentLogsInput{
			GetLister:  client,
			ClientSet:  input.ClusterProxy.GetClientSet(),
			Deployment: deployment,
			LogPath:    filepath.Join(input.LogFolder, "controllers"),
		})

		if input.DisableMetricsCollection {
			return
		}
		framework.WatchPodMetrics(ctx, framework.WatchPodMetricsInput{
			GetLister:   client,
			ClientSet:   input.ClusterProxy.GetClientSet(),
			Deployment:  deployment,
			MetricsPath: filepath.Join(input.LogFolder, "controllers"),
		})
	}
}

func localLoadE2EConfig(configPath string) *clusterctl.E2EConfig {
	configData, err := ioutil.ReadFile(configPath)
	Expect(err).ToNot(HaveOccurred(), "Failed to read the e2e test config file")
	Expect(configData).ToNot(BeEmpty(), "The e2e test config file should not be empty")

	config := &clusterctl.E2EConfig{}
	Expect(yaml.Unmarshal(configData, config)).To(Succeed(), "Failed to convert the e2e test config file to yaml")

	config.Defaults()
	config.AbsPaths(filepath.Dir(configPath))

	//TODO: this is the reason why we can't use this at present for the EKS tests
	//Expect(config.Validate()).To(Succeed(), "The e2e test config file is not valid")

	return config
}
