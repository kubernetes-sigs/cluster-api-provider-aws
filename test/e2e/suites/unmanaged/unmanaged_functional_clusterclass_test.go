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
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
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
			requiredResources = &shared.TestResource{EC2Normal: 1 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
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
					ControlPlaneMachineCount: ptr.To[int64](1),
					WorkerMachineCount:       ptr.To[int64](0),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			}, result)

			ginkgo.By("Checking if bastion host is running")
			awsCluster, err := GetAWSClusterByName(ctx, e2eCtx.Environment.BootstrapClusterProxy, namespace.Name, clusterName)
			Expect(err).To(BeNil())
			Expect(awsCluster.Status.Bastion.State).To(Equal(infrav1.InstanceStateRunning))
			expectAWSClusterConditions(awsCluster, []conditionAssertion{{infrav1.BastionHostReadyCondition, corev1.ConditionTrue, "", ""}})

			ginkgo.By("PASSED!")
		})
	})

	ginkgo.Describe("Workload cluster with AWS SSM Parameter as the Secret Backend [ClusterClass]", func() {
		ginkgo.It("should be creatable and deletable", func() {
			specName := "functional-test-ssm-parameter-store-clusterclass"
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = ptr.To[int64](1)
			configCluster.WorkerMachineCount = ptr.To[int64](1)
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

	// This test creates a workload cluster using an externally managed VPC and subnets. CAPA is still handling security group
	// creation for the cluster. All applicable resources are restricted to us-west-2a for simplicity.
	ginkgo.Describe("Workload cluster with external infrastructure [ClusterClass]", func() {
		var namespace *corev1.Namespace
		var requiredResources *shared.TestResource
		specName := "functional-test-extinfra-cc"
		mgmtClusterName := fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		mgmtClusterInfra := new(shared.AWSInfrastructure)

		// Some infrastructure creation was moved to a setup node to better organize the test.
		ginkgo.JustBeforeEach(func() {
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 5, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			ginkgo.By("Creating the management cluster infrastructure")
			mgmtClusterInfra.New(shared.AWSInfrastructureSpec{
				ClusterName:       mgmtClusterName,
				VpcCidr:           "10.0.0.0/23",
				PublicSubnetCidr:  "10.0.0.0/24",
				PrivateSubnetCidr: "10.0.1.0/24",
				AvailabilityZone:  "us-west-2a",
			}, e2eCtx)
			mgmtClusterInfra.CreateInfrastructure()
		})

		// Infrastructure cleanup is done in setup node so it is not bypassed if there is a test failure in the subject node.
		ginkgo.JustAfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
			shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			if !e2eCtx.Settings.SkipCleanup {
				ginkgo.By("Deleting the management cluster infrastructure")
				mgmtClusterInfra.DeleteInfrastructure()
			}
		})

		ginkgo.It("should create workload cluster in external VPC", func() {
			ginkgo.By("Validating management infrastructure")
			Expect(mgmtClusterInfra.VPC).NotTo(BeNil())
			Expect(*mgmtClusterInfra.State.VpcState).To(Equal("available"))
			Expect(len(mgmtClusterInfra.Subnets)).To(Equal(2))
			Expect(mgmtClusterInfra.InternetGateway).NotTo(BeNil())
			Expect(mgmtClusterInfra.ElasticIP).NotTo(BeNil())
			Expect(mgmtClusterInfra.NatGateway).NotTo(BeNil())
			Expect(len(mgmtClusterInfra.RouteTables)).To(Equal(2))

			shared.SetEnvVar("BYO_VPC_ID", *mgmtClusterInfra.VPC.VpcId, false)
			shared.SetEnvVar("BYO_PUBLIC_SUBNET_ID", *mgmtClusterInfra.State.PublicSubnetID, false)
			shared.SetEnvVar("BYO_PRIVATE_SUBNET_ID", *mgmtClusterInfra.State.PrivateSubnetID, false)

			ginkgo.By("Creating a management cluster in a peered VPC")
			mgmtConfigCluster := defaultConfigCluster(mgmtClusterName, namespace.Name)
			mgmtConfigCluster.WorkerMachineCount = ptr.To[int64](1)
			mgmtConfigCluster.Flavor = "external-vpc-clusterclass"
			mgmtCluster, mgmtMD, _ := createCluster(ctx, mgmtConfigCluster, result)

			mgmtWM := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       mgmtClusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *mgmtMD[0],
			})
			mgmtCPM := framework.GetControlPlaneMachinesByCluster(ctx, framework.GetControlPlaneMachinesByClusterInput{
				Lister:      e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName: mgmtClusterName,
				Namespace:   namespace.Name,
			})
			Expect(len(mgmtWM)).To(Equal(1))
			Expect(len(mgmtCPM)).To(Equal(1))
			ginkgo.By("Deleting the management cluster")
			deleteCluster(ctx, mgmtCluster)
		})
	})
})
