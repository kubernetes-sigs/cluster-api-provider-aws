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

package gc_managed //nolint:stylecheck

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	ms "sigs.k8s.io/cluster-api-provider-aws/test/e2e/suites/managed"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
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
		ms.VerifyRoleExistsAndOwned(ekscontrolplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("should create an EKS control plane")
		ms.ManagedClusterSpec(ctx, func() ms.ManagedClusterSpecInput {
			return ms.ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  ms.EKSManagedPoolFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       1,
			}
		})

		shared.Byf("getting cluster with name %s", clusterName)
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		ginkgo.By("getting AWSManagedControlPlane")
		cp := ms.GetControlPlaneByName(ctx, ms.GetControlPlaneByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: cluster.Spec.InfrastructureRef.Namespace,
			Name:      cluster.Spec.InfrastructureRef.Name,
		})

		shared.Byf("Waiting for the machine pool to be running")
		mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
			Lister:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		Expect(len(mp)).To(Equal(1))

		workloadClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, cluster.Namespace, cluster.Name)
		workloadYamlPath := e2eCtx.E2EConfig.GetVariable(shared.GcWorkloadPath)
		shared.Byf("Installing sample workload with load balancer services: %s", workloadYamlPath)
		workloadYaml, err := os.ReadFile(workloadYamlPath) //nolint:gosec
		Expect(err).ShouldNot(HaveOccurred())
		Expect(workloadClusterProxy.Apply(ctx, workloadYaml)).ShouldNot(HaveOccurred())

		shared.Byf("Waiting for the Deployment to be available")
		shared.WaitForDeploymentsAvailable(ctx, shared.WaitForDeploymentsAvailableInput{
			Getter:    workloadClusterProxy.GetClient(),
			Name:      "podinfo",
			Namespace: "default",
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-deployment-ready")...)

		shared.Byf("Checking we have the load balancers in AWS")
		shared.WaitForLoadBalancerToExistForService(ctx, shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             shared.LoadBalancerTypeNLB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)
		shared.WaitForLoadBalancerToExistForService(ctx, shared.WaitForLoadBalancerToExistForServiceInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             shared.LoadBalancerTypeELB,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-loadbalancer-ready")...)

		shared.Byf("Deleting workload/tenant cluster %s", clusterName)
		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		shared.Byf("Getting counts of service load balancers")
		arns, err := shared.GetLoadBalancerARNs(ctx, shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-nlb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             shared.LoadBalancerTypeNLB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(HaveLen(0), "there are %d service load balancers (nlb) still", len(arns))
		arns, err = shared.GetLoadBalancerARNs(ctx, shared.GetLoadBalancerARNsInput{
			AWSSession:       e2eCtx.BootstrapUserAWSSession,
			ServiceName:      "podinfo-elb",
			ServiceNamespace: "default",
			ClusterName:      cp.Spec.EKSClusterName,
			Type:             shared.LoadBalancerTypeELB,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(arns).To(HaveLen(0), "there are %d service load balancers (elb) still", len(arns))
	})
})

// TODO (richardcase): remove this when we merge these tests with the main eks e2e tests.
func defaultConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	return clusterctl.ConfigClusterInput{
		LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
		ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
		KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
		InfrastructureProvider:   "aws",
		Flavor:                   ms.EKSManagedPoolFlavor,
		Namespace:                namespace,
		ClusterName:              clusterName,
		KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
		ControlPlaneMachineCount: pointer.Int64Ptr(1),
		WorkerMachineCount:       pointer.Int64Ptr(0),
	}
}
