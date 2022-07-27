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

package unmanaged

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [functional] [ClusterClass]", func() {
	var (
		ctx               context.Context
		result            *clusterctl.ApplyClusterTemplateAndWaitResult
		requiredResources *shared.TestResource
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}
	})

	ginkgo.Describe("Multitenancy test [ClusterClass]", func() {
		ginkgo.It("should create cluster with nested assumed role", func() {
			// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
			specName := "functional-multitenancy-nested-clusterclass"
			requiredResources = &shared.TestResource{EC2Normal: 1 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			Expect(shared.SetMultitenancyEnvVars(e2eCtx.AWSSession)).To(Succeed())

			ginkgo.By("Creating cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ConfigCluster: clusterctl.ConfigClusterInput{
					LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
					ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
					KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
					InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
					Flavor:                   shared.NestedMultitenancyClusterClassFlavor,
					Namespace:                namespace.Name,
					ClusterName:              clusterName,
					KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
					ControlPlaneMachineCount: pointer.Int64Ptr(1),
					WorkerMachineCount:       pointer.Int64Ptr(0),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			}, result)

			ginkgo.By("Checking if bastion host is running")
			awsCluster, err := GetAWSClusterByName(ctx, namespace.Name, clusterName)
			Expect(err).To(BeNil())
			Expect(awsCluster.Status.Bastion.State).To(Equal(infrav1.InstanceStateRunning))
			expectAWSClusterConditions(awsCluster, []conditionAssertion{{infrav1.BastionHostReadyCondition, corev1.ConditionTrue, "", ""}})

			ginkgo.By("PASSED!")
		})
	})

	ginkgo.Describe("Workload cluster with AWS SSM Parameter as the Secret Backend [ClusterClass]", func() {
		ginkgo.It("should be creatable and deletable", func() {
			specName := "functional-test-ssm-parameter-store-clusterclass"
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(1)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			configCluster.Flavor = shared.TopologyFlavor
			_, md, _ := createCluster(ctx, configCluster, result)

			workerMachines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *md[0],
			})
			controlPlaneMachines := framework.GetControlPlaneMachinesByCluster(ctx, framework.GetControlPlaneMachinesByClusterInput{
				Lister:      e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName: clusterName,
				Namespace:   namespace.Name,
			})
			Expect(len(workerMachines)).To(Equal(1))
			Expect(len(controlPlaneMachines)).To(Equal(1))
		})
	})
})
