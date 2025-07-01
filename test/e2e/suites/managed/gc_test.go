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

package managed

import (
	"context"
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// EKS cluster external resource gc tests.
var _ = ginkgo.Describe("[managed] [gc] EKS Cluster external resource GC tests", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "eks-extresgc"
		clusterName string
	)

	ginkgo.It("[managed] [gc] should cleanup a cluster that has ELB/NLB load balancers", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.AWSSessionV2)

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				AWSSessionV2:             e2eCtx.BootstrapUserAWSSessionV2,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSManagedPoolFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       1,
			}
		})

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		ginkgo.By("getting AWSManagedControlPlane")
		cp := GetControlPlaneByName(ctx, GetControlPlaneByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: cluster.Spec.InfrastructureRef.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		})

		ginkgo.By("Waiting for the machine pool to be running")
		mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
			Lister:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		Expect(len(mp)).To(Equal(1))

		workloadClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name)
		workloadYamlPath := e2eCtx.E2EConfig.MustGetVariable(shared.GcWorkloadPath)
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
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)

		ginkgo.By(fmt.Sprintf("Deleting workload/tenant cluster %s", clusterName))
		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			ClusterProxy:         e2eCtx.Environment.BootstrapClusterProxy,
			Cluster:              cluster,
			ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
			ArtifactFolder:       e2eCtx.Settings.ArtifactFolder,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		ginkgo.By("Getting counts of service load balancers")
		arns, err := shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (nlb) still", len(arns))
		arns, err = shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (elb) still", len(arns))
	})

	ginkgo.It("[managed] [gc] should cleanup a cluster that has ELB/NLB load balancers using AlternativeGCStrategy", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		shared.ReconfigureDeployment(ctx, shared.ReconfigureDeploymentInput{
			Getter:       e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			ClientSet:    e2eCtx.Environment.BootstrapClusterProxy.GetClientSet(),
			Name:         "capa-controller-manager",
			Namespace:    "capa-system",
			WaitInterval: e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready"),
		}, shared.EnableAlternativeGCStrategy, shared.ValidateAlternativeGCStrategyEnabled)

		specName += "-alterstrategy"
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.AWSSessionV2)

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				AWSSessionV2:             e2eCtx.BootstrapUserAWSSessionV2,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSManagedPoolFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       1,
			}
		})

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		ginkgo.By("getting AWSManagedControlPlane")
		cp := GetControlPlaneByName(ctx, GetControlPlaneByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: cluster.Spec.InfrastructureRef.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		})

		ginkgo.By("Waiting for the machine pool to be running")
		mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
			Lister:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		Expect(len(mp)).To(Equal(1))

		workloadClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name)
		workloadYamlPath := e2eCtx.E2EConfig.MustGetVariable(shared.GcWorkloadPath)
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
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)
		shared.WaitForLoadBalancerToExistForService(shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeELB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)

		ginkgo.By(fmt.Sprintf("Deleting workload/tenant cluster %s", clusterName))
		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			ClusterProxy:         e2eCtx.Environment.BootstrapClusterProxy,
			Cluster:              cluster,
			ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
			ArtifactFolder:       e2eCtx.Settings.ArtifactFolder,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		ginkgo.By("Getting counts of service load balancers")
		arns, err := shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             infrav1.LoadBalancerTypeNLB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(BeEmpty(), "there are %d service load balancers (nlb) still", len(arns))
		arns, err = shared.GetLoadBalancerARNs(shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
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
