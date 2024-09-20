//go:build e2e
// +build e2e

/*
Copyright 2022 The Kubernetes Authors.

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

package gc_unmanaged //nolint:stylecheck

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [gc]", func() {
	var (
		ctx               context.Context
		result            *clusterctl.ApplyClusterTemplateAndWaitResult
		requiredResources *shared.TestResource
	)

	ginkgo.It("[unmanaged] [gc] should cleanup a cluster that has ELB/NLB load balancers", func() {
		ginkgo.By("should have a valid test configuration")
		specName := "unmanaged-gc-cluster"

		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}

		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
		requiredResources.WriteRequestedResources(e2eCtx, specName)
		Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		namespace := shared.SetupNamespace(ctx, specName, e2eCtx)
		defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
		ginkgo.By("Creating cluster with single control plane")
		clusterName := fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		configCluster := defaultConfigCluster(clusterName, namespace.Name)
		configCluster.WorkerMachineCount = ptr.To[int64](1)
		createCluster(ctx, configCluster, result)

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		workloadClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name)
		workloadYamlPath := e2eCtx.E2EConfig.GetVariable(shared.GcWorkloadPath)
		ginkgo.By(fmt.Sprintf("Installing sample workload with load balancer services: %s", workloadYamlPath))
		workloadYaml, err := os.ReadFile(workloadYamlPath) //nolint:gosec
		Expect(err).ShouldNot(HaveOccurred())
		Expect(workloadClusterProxy.CreateOrUpdate(ctx, workloadYaml)).ShouldNot(HaveOccurred())

		ginkgo.By("Waiting for the Deployment to be available")
		shared.WaitForDeploymentsAvailable(ctx, shared.WaitForDeploymentsAvailableInput{
			Getter:    workloadClusterProxy.GetClient(),
			Name:      "podinfo",
			Namespace: "default",
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready")...)

		ginkgo.By("Checking we have the load balancers in AWS")
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)

		ginkgo.By(fmt.Sprintf("Deleting workload/tenant cluster %s", clusterName))
		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			Client:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		ginkgo.By("Getting counts of service load balancers")
		arns, err := shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (nlb) still", len(arns))
		arns, err = shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (elb) still", len(arns))
	})

	ginkgo.It("[unmanaged] [gc] should cleanup a cluster that has ELB/NLB load balancers using AlternativeGCStrategy", func() {
		ginkgo.By("should have a valid test configuration")
		specName := "unmanaged-gc-alterstrategy-cluster"

		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}

		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		shared.ReconfigureDeployment(ctx, shared.ReconfigureDeploymentInput{
			Getter:       e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			ClientSet:    e2eCtx.Environment.BootstrapClusterProxy.GetClientSet(),
			Name:         "capa-controller-manager",
			Namespace:    "capa-system",
			WaitInterval: e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready"),
		}, shared.EnableAlternativeGCStrategy, shared.ValidateAlternativeGCStrategyEnabled)

		requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
		requiredResources.WriteRequestedResources(e2eCtx, specName)
		Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		namespace := shared.SetupNamespace(ctx, specName, e2eCtx)
		defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
		ginkgo.By("Creating cluster with single control plane")
		clusterName := fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		configCluster := defaultConfigCluster(clusterName, namespace.Name)
		configCluster.WorkerMachineCount = ptr.To[int64](1)
		c, md, cp := createCluster(ctx, configCluster, result)
		Expect(c).NotTo(BeNil(), "Expecting cluster created")
		Expect(len(md)).To(Equal(1), "Expecting one MachineDeployment")
		Expect(cp).NotTo(BeNil(), "Expecting control plane created")

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		workloadClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name)
		workloadYamlPath := e2eCtx.E2EConfig.GetVariable(shared.GcWorkloadPath)
		ginkgo.By(fmt.Sprintf("Installing sample workload with load balancer services: %s", workloadYamlPath))
		workloadYaml, err := os.ReadFile(workloadYamlPath) //nolint:gosec
		Expect(err).ShouldNot(HaveOccurred())
		Expect(workloadClusterProxy.CreateOrUpdate(ctx, workloadYaml)).ShouldNot(HaveOccurred())

		ginkgo.By("Waiting for the Deployment to be available")
		shared.WaitForDeploymentsAvailable(ctx, shared.WaitForDeploymentsAvailableInput{
			Getter:    workloadClusterProxy.GetClient(),
			Name:      "podinfo",
			Namespace: "default",
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready")...)

		ginkgo.By("Checking we have the load balancers in AWS")
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)

		ginkgo.By(fmt.Sprintf("Deleting workload/tenant cluster %s", clusterName))
		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			Client:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		ginkgo.By("Getting counts of service load balancers")
		arns, err := shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (nlb) still", len(arns))
		arns, err = shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      clusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (elb) still", len(arns))

		shared.ReconfigureDeployment(ctx, shared.ReconfigureDeploymentInput{
			Getter:       e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			ClientSet:    e2eCtx.Environment.BootstrapClusterProxy.GetClientSet(),
			Name:         "capa-controller-manager",
			Namespace:    "capa-system",
			WaitInterval: e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready"),
		}, shared.DisableAlternativeGCStrategy, shared.ValidateAlternativeGCStrategyDisabled)
	})
})

// TODO (richardcase): remove this when we merge these tests with the main eks e2e tests.
func defaultConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	return clusterctl.ConfigClusterInput{
		LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
		ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
		KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
		InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
		Flavor:                   clusterctl.DefaultFlavor,
		Namespace:                namespace,
		ClusterName:              clusterName,
		KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
		ControlPlaneMachineCount: ptr.To[int64](1),
		WorkerMachineCount:       ptr.To[int64](0),
	}
}

// TODO (richardcase): remove this when we merge these tests with the main eks e2e tests.
func createCluster(ctx context.Context, configCluster clusterctl.ConfigClusterInput, result *clusterctl.ApplyClusterTemplateAndWaitResult) (*clusterv1.Cluster, []*clusterv1.MachineDeployment, *controlplanev1.KubeadmControlPlane) {
	clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		ClusterProxy:                 e2eCtx.Environment.BootstrapClusterProxy,
		ConfigCluster:                configCluster,
		WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals("", "wait-cluster"),
		WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals("", "wait-control-plane"),
		WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes"),
	}, result)

	return result.Cluster, result.MachineDeployments, result.ControlPlane
}
