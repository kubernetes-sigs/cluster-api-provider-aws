//go:build e2e
// +build e2e

/*
Copyright 2025 The Kubernetes Authors.

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

package unmanaged

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [upgrade]", func() {
	var (
		specName          = "classic-elb-upgrade"
		ctx               context.Context
		testNamespace     *corev1.Namespace
		testCancelWatches context.CancelFunc

		managementClusterName          string
		managementClusterNamespace     *corev1.Namespace
		managementClusterCancelWatches context.CancelFunc
		managementClusterResources     *clusterctl.ApplyClusterTemplateAndWaitResult
		managementClusterProvider      bootstrap.ClusterProvider
		managementClusterProxy         framework.ClusterProxy

		kubernetesVersionFrom string
		kubernetesVersionTo   string
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		managementClusterResources = new(clusterctl.ApplyClusterTemplateAndWaitResult)

		kubernetesVersionFrom = e2eCtx.E2EConfig.GetVariable(shared.ClassicElbTestKubernetesFrom)
		Expect(kubernetesVersionFrom).ToNot(BeEmpty(), "kubernetesVersionFrom is not set")
		kubernetesVersionTo = e2eCtx.E2EConfig.GetVariable(shared.ClassicElbTestKubernetesTo)
		Expect(kubernetesVersionTo).ToNot(BeEmpty(), "kubernetesVersionTo is not set")
	})

	ginkgo.AfterEach(func() {
		if testNamespace != nil {
			// Dump all the logs from the workload cluster before deleting them.
			framework.DumpAllResourcesAndLogs(ctx, managementClusterProxy, e2eCtx.Settings.ArtifactFolder, testNamespace, managementClusterResources.Cluster)

			if !e2eCtx.Settings.SkipCleanup {
				shared.Byf("Deleting all clusters in namespace %s in management cluster %s", testNamespace.Name, managementClusterName)
				deleteAllClustersAndWait(ctx, deleteAllClustersAndWaitInput{
					Client:    managementClusterProxy.GetClient(),
					Namespace: testNamespace.Name,
				}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-delete-cluster")...)

				shared.Byf("Deleting namespace %s used for hosting the %q test", testNamespace.Name, specName)
				framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{
					Deleter: managementClusterProxy.GetClient(),
					Name:    testNamespace.Name,
				})

				shared.Byf("Deleting providers")
				clusterctl.Delete(ctx, clusterctl.DeleteInput{
					LogFolder:            filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", managementClusterResources.Cluster.Name),
					ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
					KubeconfigPath:       managementClusterProxy.GetKubeconfigPath(),
				})
			}
			testCancelWatches()
		}

		if !e2eCtx.Settings.SkipCleanup {
			managementClusterProxy.Dispose(ctx)
			managementClusterProvider.Dispose(ctx)
		} else {
			framework.DumpSpecResourcesAndCleanup(ctx, specName, e2eCtx.Environment.BootstrapClusterProxy, e2eCtx.Settings.ArtifactFolder, managementClusterNamespace, managementClusterCancelWatches, managementClusterResources.Cluster, e2eCtx.E2EConfig.GetIntervals, e2eCtx.Settings.SkipCleanup)
		}
	})

	//  classic elb upgrade to v1.30+
	ginkgo.Describe("Should create a management cluster and upgrade the workload cluster to v1.30+", func() {
		ginkgo.It("Should create a management cluster and upgrade the workload cluster to v1.30+", func() {
			managementClusterName = fmt.Sprintf("%s-management-%s", specName, util.RandomString(6))
			managementClusterLogFolder := filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", managementClusterName)
			managemntClusterVersion := e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersionManagement)

			ginkgo.By("Creating a kind cluster to be used as a new management cluster")

			managementClusterProvider = bootstrap.CreateKindBootstrapClusterAndLoadImages(ctx, bootstrap.CreateKindBootstrapClusterAndLoadImagesInput{
				Name:               managementClusterName,
				KubernetesVersion:  managemntClusterVersion,
				RequiresDockerSock: false,
				Images:             e2eCtx.E2EConfig.Images,
				LogFolder:          filepath.Join(managementClusterLogFolder, "logs-kind"),
			})

			Expect(managementClusterProvider).ToNot(BeNil(), "Failed to create a kind cluster")

			kubeconfigPath := managementClusterProvider.GetKubeconfigPath()
			Expect(kubeconfigPath).To(BeAnExistingFile(), "Failed to get the kubeconfig file for the kind cluster")

			scheme := runtime.NewScheme()
			clusterv1.AddToScheme(scheme)
			infrav1.AddToScheme(scheme)
			framework.TryAddDefaultSchemes(scheme)
			managementClusterProxy = framework.NewClusterProxy(managementClusterName, kubeconfigPath, scheme)
			Expect(managementClusterProxy).ToNot(BeNil(), "Failed to get a kind cluster proxy")

			ginkgo.By("Turning the new cluster into a management cluster with older versions of CAPA")

			clusterctl.InitManagementClusterAndWatchControllerLogs(ctx, clusterctl.InitManagementClusterAndWatchControllerLogsInput{
				ClusterProxy:            managementClusterProxy,
				ClusterctlConfigPath:    e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapProviders:      e2eCtx.BootstrapProviders(),
				ControlPlaneProviders:   e2eCtx.ControlPlaneProviders(),
				InfrastructureProviders: []string{"aws:v2.7.1"},
				LogFolder:               managementClusterLogFolder,
			}, e2eCtx.E2EConfig.GetIntervals(e2eCtx.Environment.BootstrapClusterProxy.GetName(), "wait-controllers")...)

			shared.CreateAWSClusterControllerIdentity(managementClusterProxy.GetClient())

			ginkgo.By("Management cluster with older version of CAPA is running")

			shared.Byf("Creating a namespace for hosting the %s test workload cluster", specName)
			testNamespace, testCancelWatches = framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
				Creator:   managementClusterProxy.GetClient(),
				ClientSet: managementClusterProxy.GetClientSet(),
				Name:      specName,
				LogFolder: filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", "bootstrap"),
			})

			ginkgo.By("Creating a test workload cluster with k8s version < v1.30")
			workloadClusterName := fmt.Sprintf("%s-workload-%s", specName, util.RandomString(6))
			workloadClusterNamespace := testNamespace.Name

			clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy: managementClusterProxy,
				ConfigCluster: clusterctl.ConfigClusterInput{
					LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
					ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
					KubeconfigPath:           managementClusterProxy.GetKubeconfigPath(),
					InfrastructureProvider:   "aws",
					Flavor:                   "classicelb-upgrade",
					Namespace:                workloadClusterNamespace,
					ClusterName:              workloadClusterName,
					KubernetesVersion:        kubernetesVersionFrom,
					ControlPlaneMachineCount: ptr.To[int64](1),
					WorkerMachineCount:       ptr.To[int64](1),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			}, managementClusterResources)

			ginkgo.By("Checking the load balancer for the workload cluster is using SSL for the health check")

			awsCluster, err := GetAWSClusterByName(ctx, managementClusterProxy, workloadClusterNamespace, workloadClusterName)
			Expect(err).NotTo(HaveOccurred(), "failed to get the AWS cluster")

			shared.CheckClassicElbHealthCheck(shared.CheckClassicElbHealthCheckInput{
				AWSSession:       e2eCtx.BootstrapUserAWSSession,
				LoadBalancerName: awsCluster.Status.Network.APIServerELB.Name,
				ExpectedTarget:   "SSL:6443",
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-classic-elb-health-check-short")...)

			ginkgo.By("Now the workload cluster is ready upgrade CAPA to main version")

			clusterctl.UpgradeManagementClusterAndWait(ctx, clusterctl.UpgradeManagementClusterAndWaitInput{
				ClusterctlConfigPath:    e2eCtx.Environment.ClusterctlConfigPath,
				ClusterProxy:            managementClusterProxy,
				InfrastructureProviders: []string{"aws:v9.9.99"},
				LogFolder:               managementClusterLogFolder,
			}, e2eCtx.E2EConfig.GetIntervals(e2eCtx.Environment.BootstrapClusterProxy.GetName(), "wait-controllers")...)

			time.Sleep(1 * time.Minute)

			ginkgo.By("Management cluster is upgraded to main version")

			shared.CheckClassicElbHealthCheck(shared.CheckClassicElbHealthCheckInput{
				AWSSession:       e2eCtx.BootstrapUserAWSSession,
				LoadBalancerName: awsCluster.Status.Network.APIServerELB.Name,
				ExpectedTarget:   "TCP:6443",
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-classic-elb-health-check-long")...)

			shared.Byf("Upgrading the control plane to %s", kubernetesVersionTo)

			framework.UpgradeControlPlaneAndWaitForUpgrade(ctx, framework.UpgradeControlPlaneAndWaitForUpgradeInput{
				ClusterProxy:                managementClusterProxy,
				Cluster:                     managementClusterResources.Cluster,
				ControlPlane:                managementClusterResources.ControlPlane,
				KubernetesUpgradeVersion:    kubernetesVersionTo,
				WaitForMachinesToBeUpgraded: e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade"),
				WaitForKubeProxyUpgrade:     e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade"),
				WaitForDNSUpgrade:           e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade"),
				WaitForEtcdUpgrade:          e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade"),
			})

			shared.Byf("Upgrading the machine deployment to %s", kubernetesVersionTo)

			framework.UpgradeMachineDeploymentsAndWait(ctx, framework.UpgradeMachineDeploymentsAndWaitInput{
				ClusterProxy:                managementClusterProxy,
				Cluster:                     managementClusterResources.Cluster,
				UpgradeVersion:              kubernetesVersionTo,
				MachineDeployments:          managementClusterResources.MachineDeployments,
				WaitForMachinesToBeUpgraded: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("All machines are upgraded")
		})
	})
})

// deleteAllClustersAndWaitInput is the input type for deleteAllClustersAndWait.
type deleteAllClustersAndWaitInput struct {
	Client    client.Client
	Namespace string
}

// deleteAllClustersAndWait deletes all cluster resources in the given namespace and waits for them to be gone.
func deleteAllClustersAndWait(ctx context.Context, input deleteAllClustersAndWaitInput, intervals ...interface{}) {
	Expect(ctx).NotTo(BeNil(), "ctx is required for deleteAllClustersAndWaitOldAPI")
	Expect(input.Client).ToNot(BeNil(), "Invalid argument. input.Client can't be nil when calling deleteAllClustersAndWaitOldAPI")
	Expect(input.Namespace).ToNot(BeEmpty(), "Invalid argument. input.Namespace can't be empty when calling deleteAllClustersAndWaitOldAPI")

	coreCAPIStorageVersion := getCoreCAPIStorageVersion(ctx, input.Client)

	clusterList := &unstructured.UnstructuredList{}
	clusterList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   clusterv1.GroupVersion.Group,
		Version: coreCAPIStorageVersion,
		Kind:    "ClusterList",
	})
	Expect(input.Client.List(ctx, clusterList, client.InNamespace(input.Namespace))).To(Succeed(), "Failed to list clusters in namespace %s", input.Namespace)

	// Enforce restart of kube-controller-manager by stealing its lease.
	// Note: Due to a known issue in the kube-controller-manager we have to restart it
	// in case the kube-controller-manager internally caches ownerRefs of apiVersions
	// which now don't exist anymore (e.g. v1alpha3/v1alpha4).
	// Alternatives to this would be:
	// * some other way to restart the kube-controller-manager (e.g. control plane node rollout)
	// * removing ownerRefs from (at least) MachineDeployments
	Eventually(func(g Gomega) {
		kubeControllerManagerLease := &coordinationv1.Lease{}
		g.Expect(input.Client.Get(ctx, client.ObjectKey{Namespace: metav1.NamespaceSystem, Name: "kube-controller-manager"}, kubeControllerManagerLease)).To(Succeed())
		// As soon as the kube-controller-manager detects it doesn't own the lease anymore it will restart.
		// Once the current lease times out the kube-controller-manager will become leader again.
		kubeControllerManagerLease.Spec.HolderIdentity = ptr.To("e2e-test-client")
		g.Expect(input.Client.Update(ctx, kubeControllerManagerLease)).To(Succeed())
	}, 3*time.Minute, 3*time.Second).Should(Succeed(), "failed to steal lease from kube-controller-manager to trigger restart")

	for _, c := range clusterList.Items {
		shared.Byf("Deleting cluster %s", c.GetName())
		Expect(input.Client.Delete(ctx, c.DeepCopy())).To(Succeed())
	}

	for _, c := range clusterList.Items {
		shared.Byf("Waiting for cluster %s to be deleted", c.GetName())
		Eventually(func() bool {
			cluster := c.DeepCopy()
			key := client.ObjectKey{
				Namespace: c.GetNamespace(),
				Name:      c.GetName(),
			}
			return apierrors.IsNotFound(input.Client.Get(ctx, key, cluster))
		}, intervals...).Should(BeTrue())
	}
}

func getCoreCAPIStorageVersion(ctx context.Context, c client.Client) string {
	clusterCRD := &apiextensionsv1.CustomResourceDefinition{}
	if err := c.Get(ctx, client.ObjectKey{Name: "clusters.cluster.x-k8s.io"}, clusterCRD); err != nil {
		Expect(err).ToNot(HaveOccurred(), "failed to retrieve a machine CRD")
	}
	// Pick the storage version
	for _, version := range clusterCRD.Spec.Versions {
		if version.Storage {
			return version.Name
		}
	}
	ginkgo.Fail("Cluster CRD has no storage version")
	return ""
}
