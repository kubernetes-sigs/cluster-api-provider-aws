//go:build e2e
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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kindcluster "sigs.k8s.io/kind/pkg/cluster"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

const selfHostedManagementClusterNamespace = "self-hosted-management-cluster"
const selfHostedManagementClusterStateFile = "self-hosted-management-cluster.json"
const selfHostedManagementClusterKubeconfigFile = "self-hosted-management-cluster.kubeconfig"

type selfHostedManagementClusterState struct {
	KindClusterName          string `json:"kindClusterName"`
	KindKubeconfigPath       string `json:"kindKubeconfigPath"`
	ManagementKubeconfigPath string `json:"managementKubeconfigPath"`
	ClusterName              string `json:"clusterName"`
	Namespace                string `json:"namespace"`
	ClusterctlConfigPath     string `json:"clusterctlConfigPath"`
}

func validateManagementClusterSettings(settings Settings) error {
	if settings.UseExistingCluster && settings.ProvisionSelfHostedManagementCluster {
		return errors.New("use-existing-cluster and provision-self-hosted-management-cluster cannot be enabled together")
	}
	if settings.ProvisionSelfHostedManagementCluster && settings.TeardownSelfHostedManagementCluster {
		return errors.New("provision-self-hosted-management-cluster and teardown-self-hosted-management-cluster cannot be enabled together")
	}
	return nil
}

type selfHostedManagementClusterProvider struct {
	kindProvider        bootstrap.ClusterProvider
	kindProxy           framework.ClusterProxy
	managementProxy     framework.ClusterProxy
	cluster             *clusterv1.Cluster
	namespace           *corev1.Namespace
	cancelKindWatches   context.CancelFunc
	cancelTargetWatches context.CancelFunc
	clusterctlConfig    string
	artifactFolder      string
	statePath           string
	kindClusterName     string
	intervals           func(string, string) []interface{}
}

func (p *selfHostedManagementClusterProvider) Create(context.Context) {}

func (p *selfHostedManagementClusterProvider) GetKubeconfigPath() string {
	return p.managementProxy.GetKubeconfigPath()
}

func (p *selfHostedManagementClusterProvider) Dispose(ctx context.Context) {
	clusterctl.Move(ctx, clusterctl.MoveInput{
		LogFolder:            filepath.Join(p.artifactFolder, "clusters", p.cluster.Name),
		ClusterctlConfigPath: p.clusterctlConfig,
		FromKubeconfigPath:   p.managementProxy.GetKubeconfigPath(),
		ToKubeconfigPath:     p.kindProxy.GetKubeconfigPath(),
		Namespace:            p.namespace.Name,
	})

	p.cluster = framework.DiscoveryAndWaitForCluster(ctx, framework.DiscoveryAndWaitForClusterInput{
		Getter:    p.kindProxy.GetClient(),
		Namespace: p.cluster.Namespace,
		Name:      p.cluster.Name,
	}, p.intervals("self-hosted", "wait-cluster")...)

	framework.DeleteClusterAndWait(ctx, framework.DeleteClusterAndWaitInput{
		ClusterProxy:         p.kindProxy,
		ClusterctlConfigPath: p.clusterctlConfig,
		Cluster:              p.cluster,
		ArtifactFolder:       p.artifactFolder,
	}, p.intervals("self-hosted", "wait-delete-cluster")...)

	if p.cancelTargetWatches != nil {
		p.cancelTargetWatches()
	}
	if p.cancelKindWatches != nil {
		p.cancelKindWatches()
	}
	p.managementProxy.Dispose(ctx)
	p.kindProxy.Dispose(ctx)
	if p.kindProvider != nil {
		p.kindProvider.Dispose(ctx)
	} else {
		Expect(kindcluster.NewProvider().Delete(p.kindClusterName, p.kindProxy.GetKubeconfigPath())).To(Succeed())
	}
	Expect(os.Remove(filepath.Join(p.artifactFolder, selfHostedManagementClusterKubeconfigFile))).To(Succeed())
	Expect(os.Remove(p.statePath)).To(Succeed())
}

func persistSelfHostedManagementClusterKubeconfig(sourcePath, artifactFolder string) (string, error) {
	managementKubeconfig, err := os.ReadFile(sourcePath) //nolint:gosec // The source path is created by the e2e framework.
	if err != nil {
		return "", err
	}
	managementKubeconfigPath := filepath.Join(artifactFolder, selfHostedManagementClusterKubeconfigFile)
	if err := os.WriteFile(managementKubeconfigPath, managementKubeconfig, 0o600); err != nil { //nolint:gosec // The artifact directory is supplied by the e2e runner.
		return "", err
	}
	return managementKubeconfigPath, nil
}

func writeSelfHostedManagementClusterState(provider *selfHostedManagementClusterProvider) {
	managementKubeconfigPath, err := persistSelfHostedManagementClusterKubeconfig(provider.managementProxy.GetKubeconfigPath(), provider.artifactFolder)
	Expect(err).NotTo(HaveOccurred())
	state := selfHostedManagementClusterState{
		KindClusterName:          provider.kindClusterName,
		KindKubeconfigPath:       provider.kindProxy.GetKubeconfigPath(),
		ManagementKubeconfigPath: managementKubeconfigPath,
		ClusterName:              provider.cluster.Name,
		Namespace:                provider.namespace.Name,
		ClusterctlConfigPath:     provider.clusterctlConfig,
	}
	data, err := json.MarshalIndent(state, "", "  ")
	Expect(err).NotTo(HaveOccurred())
	Expect(os.WriteFile(provider.statePath, append(data, '\n'), 0o600)).To(Succeed())
}

func loadSelfHostedManagementCluster(e2eCtx *E2EContext) (bootstrap.ClusterProvider, framework.ClusterProxy) {
	statePath := filepath.Join(e2eCtx.Settings.ArtifactFolder, selfHostedManagementClusterStateFile)
	data, err := os.ReadFile(statePath) //nolint:gosec // The path is fixed under the configured e2e artifact directory.
	Expect(err).NotTo(HaveOccurred(), "Failed to read self-hosted management-cluster state")
	state := selfHostedManagementClusterState{}
	Expect(json.Unmarshal(data, &state)).To(Succeed())

	kindProxy := framework.NewClusterProxy("bootstrap", state.KindKubeconfigPath, e2eCtx.Environment.Scheme)
	managementProxy := framework.NewClusterProxy("self-hosted-management", state.ManagementKubeconfigPath, e2eCtx.Environment.Scheme)
	provider := &selfHostedManagementClusterProvider{
		kindProxy:        kindProxy,
		managementProxy:  managementProxy,
		cluster:          &clusterv1.Cluster{ObjectMeta: metav1.ObjectMeta{Name: state.ClusterName, Namespace: state.Namespace}},
		namespace:        &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: state.Namespace}},
		clusterctlConfig: state.ClusterctlConfigPath,
		artifactFolder:   e2eCtx.Settings.ArtifactFolder,
		statePath:        statePath,
		kindClusterName:  state.KindClusterName,
		intervals:        e2eCtx.E2EConfig.GetIntervals,
	}
	return provider, managementProxy
}

// createClusterctlLocalRepository generates a clusterctl repository.
// Must always be run after kubetest.NewConfiguration.
func createClusterctlLocalRepository(e2eCtx *E2EContext, repositoryFolder string) string {
	createRepositoryInput := clusterctl.CreateRepositoryInput{
		E2EConfig:        e2eCtx.E2EConfig,
		RepositoryFolder: repositoryFolder,
	}

	if !e2eCtx.IsManaged {
		// Ensuring a CNI file is defined in the config and register a FileTransformation to inject the referenced file as in place of the CNI_RESOURCES envSubst variable.
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(capi_e2e.CNIPath), "Missing %s variable in the config", capi_e2e.CNIPath)
		cniPath := e2eCtx.E2EConfig.MustGetVariable(capi_e2e.CNIPath)
		Expect(cniPath).To(BeAnExistingFile(), "The %s variable should resolve to an existing file", capi_e2e.CNIPath)
		createRepositoryInput.RegisterClusterResourceSetConfigMapTransformation(cniPath, capi_e2e.CNIResources)
	}

	clusterctlConfig := clusterctl.CreateRepository(context.TODO(), createRepositoryInput)
	Expect(clusterctlConfig).To(BeAnExistingFile(), "The clusterctl generate file does not exists in the local repository %s", repositoryFolder)
	return clusterctlConfig
}

// setupBootstrapCluster installs Cluster API components via clusterctl.
func setupBootstrapCluster(config *clusterctl.E2EConfig, scheme *runtime.Scheme, useExistingCluster bool, options ...framework.Option) (bootstrap.ClusterProvider, framework.ClusterProxy) {
	var clusterProvider bootstrap.ClusterProvider
	kubeconfigPath := ""
	if !useExistingCluster {
		clusterProvider = bootstrap.CreateKindBootstrapClusterAndLoadImages(context.TODO(), bootstrap.CreateKindBootstrapClusterAndLoadImagesInput{
			Name:               config.ManagementClusterName,
			KubernetesVersion:  config.MustGetVariable(KubernetesVersionManagement),
			RequiresDockerSock: config.HasDockerProvider(),
			Images:             config.Images,
		})
		Expect(clusterProvider).ToNot(BeNil(), "Failed to create a bootstrap cluster")

		kubeconfigPath = clusterProvider.GetKubeconfigPath()
		Expect(kubeconfigPath).To(BeAnExistingFile(), "Failed to get the kubeconfig file for the bootstrap cluster")
	}

	clusterProxy := framework.NewClusterProxy("bootstrap", kubeconfigPath, scheme, options...)
	Expect(clusterProxy).ToNot(BeNil(), "Failed to get a bootstrap cluster proxy")

	return clusterProvider, clusterProxy
}

// initBootstrapCluster uses kind to create a cluster.
func initBootstrapCluster(e2eCtx *E2EContext) {
	// NOTE: the following originally used clusterctl.InitManagementClusterAndWatchControllerLogs.
	// This can be used again when https://github.com/kubernetes-sigs/cluster-api/issues/3983 is completed
	InitManagementClusterAndWatchControllerLogs(context.TODO(), InitManagementClusterAndWatchControllerLogsInput{
		ClusterProxy:            e2eCtx.Environment.BootstrapClusterProxy,
		ClusterctlConfigPath:    e2eCtx.Environment.ClusterctlConfigPath,
		InfrastructureProviders: e2eCtx.E2EConfig.InfrastructureProviders(),
		BootstrapProviders:      e2eCtx.BootstrapProviders(),
		ControlPlaneProviders:   e2eCtx.ControlPlaneProviders(),
		LogFolder:               filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
	}, e2eCtx.E2EConfig.GetIntervals(e2eCtx.Environment.BootstrapClusterProxy.GetName(), "wait-controllers")...)
}

func setupSelfHostedManagementCluster(e2eCtx *E2EContext) (bootstrap.ClusterProvider, framework.ClusterProxy) {
	kindProvider := e2eCtx.Environment.BootstrapClusterProvider
	kindProxy := e2eCtx.Environment.BootstrapClusterProxy
	Expect(kindProvider).ToNot(BeNil(), "A kind bootstrap cluster is required to provision a self-hosted management cluster")

	namespace, cancelKindWatches := framework.CreateNamespaceAndWatchEvents(context.TODO(), framework.CreateNamespaceAndWatchEventsInput{
		Creator:   kindProxy.GetClient(),
		ClientSet: kindProxy.GetClientSet(),
		Name:      selfHostedManagementClusterNamespace,
		LogFolder: filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", kindProxy.GetName()),
	})

	clusterResources := &clusterctl.ApplyClusterTemplateAndWaitResult{}
	controlPlaneMachineCount := int64(1)
	workerMachineCount := int64(1)
	clusterctl.ApplyClusterTemplateAndWait(context.TODO(), clusterctl.ApplyClusterTemplateAndWaitInput{
		ClusterProxy: kindProxy,
		ConfigCluster: clusterctl.ConfigClusterInput{
			LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", kindProxy.GetName()),
			ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
			KubeconfigPath:           kindProxy.GetKubeconfigPath(),
			InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
			Flavor:                   "remote-management-cluster",
			Namespace:                namespace.Name,
			ClusterName:              e2eCtx.E2EConfig.ManagementClusterName + "-aws",
			KubernetesVersion:        e2eCtx.E2EConfig.MustGetVariable(KubernetesVersionManagement),
			ControlPlaneMachineCount: &controlPlaneMachineCount,
			WorkerMachineCount:       &workerMachineCount,
		},
		WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals("self-hosted", "wait-cluster"),
		WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals("self-hosted", "wait-control-plane"),
		WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals("self-hosted", "wait-worker-nodes"),
	}, clusterResources)

	managementProxy := kindProxy.GetWorkloadCluster(context.TODO(), namespace.Name, clusterResources.Cluster.Name,
		framework.WithMachineLogCollector(kindProxy.GetLogCollector()))
	_, cancelTargetWatches := framework.CreateNamespaceAndWatchEvents(context.TODO(), framework.CreateNamespaceAndWatchEventsInput{
		Creator:   managementProxy.GetClient(),
		ClientSet: managementProxy.GetClientSet(),
		Name:      namespace.Name,
		LogFolder: filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", clusterResources.Cluster.Name),
	})

	InitManagementClusterAndWatchControllerLogs(context.TODO(), InitManagementClusterAndWatchControllerLogsInput{
		ClusterProxy:            managementProxy,
		ClusterctlConfigPath:    e2eCtx.Environment.ClusterctlConfigPath,
		InfrastructureProviders: e2eCtx.E2EConfig.InfrastructureProviders(),
		BootstrapProviders:      e2eCtx.BootstrapProviders(),
		ControlPlaneProviders:   e2eCtx.ControlPlaneProviders(),
		LogFolder:               filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", clusterResources.Cluster.Name),
	}, e2eCtx.E2EConfig.GetIntervals("self-hosted", "wait-controllers")...)

	clusterctl.Move(context.TODO(), clusterctl.MoveInput{
		LogFolder:            filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", kindProxy.GetName()),
		ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
		FromKubeconfigPath:   kindProxy.GetKubeconfigPath(),
		ToKubeconfigPath:     managementProxy.GetKubeconfigPath(),
		Namespace:            namespace.Name,
	})

	provider := &selfHostedManagementClusterProvider{
		kindProvider:        kindProvider,
		kindProxy:           kindProxy,
		managementProxy:     managementProxy,
		cluster:             clusterResources.Cluster,
		namespace:           namespace,
		cancelKindWatches:   cancelKindWatches,
		cancelTargetWatches: cancelTargetWatches,
		clusterctlConfig:    e2eCtx.Environment.ClusterctlConfigPath,
		artifactFolder:      e2eCtx.Settings.ArtifactFolder,
		statePath:           filepath.Join(e2eCtx.Settings.ArtifactFolder, selfHostedManagementClusterStateFile),
		kindClusterName:     e2eCtx.E2EConfig.ManagementClusterName,
		intervals:           e2eCtx.E2EConfig.GetIntervals,
	}
	writeSelfHostedManagementClusterState(provider)
	return provider, managementProxy
}

// tearDown the bootstrap kind cluster.
func tearDown(bootstrapClusterProvider bootstrap.ClusterProvider, bootstrapClusterProxy framework.ClusterProxy) {
	if provider, ok := bootstrapClusterProvider.(*selfHostedManagementClusterProvider); ok {
		provider.Dispose(context.TODO())
		return
	}
	if bootstrapClusterProxy != nil {
		bootstrapClusterProxy.Dispose(context.TODO())
	}
	if bootstrapClusterProvider != nil {
		bootstrapClusterProvider.Dispose(context.TODO())
	}
}
