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
	"os"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"

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

	ginkgo.Describe("Workload cluster in multiple AZs [ClusterClass]", func() {
		ginkgo.It("It should be creatable and deletable", func() {
			specName := "functional-test-multi-az-clusterclass"
			requiredResources = &shared.TestResource{EC2Normal: 3 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(3)
			configCluster.Flavor = shared.MultiAzClusterClassFlavor
			cluster, _, _ := createCluster(ctx, configCluster, result)

			ginkgo.By("Adding worker nodes to additional subnets")
			mdName1 := clusterName + "-md-1"
			mdName2 := clusterName + "-md-2"
			md1 := makeMachineDeployment(namespace.Name, mdName1, clusterName, 1)
			md2 := makeMachineDeployment(namespace.Name, mdName2, clusterName, 1)
			az1 := os.Getenv(shared.AwsAvailabilityZone1)
			az2 := os.Getenv(shared.AwsAvailabilityZone2)

			// private CIDRs set in cluster-template-multi-az.yaml.
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       md1,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName1),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName1, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), pointer.StringPtr(az1), getSubnetID("cidr-block", "10.0.0.0/24", clusterName)),
			})
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       md2,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName2),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName2, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), pointer.StringPtr(az2), getSubnetID("cidr-block", "10.0.2.0/24", clusterName)),
			})

			ginkgo.By("Waiting for new worker nodes to become ready")
			k8sClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md1}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md2}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		})
	})
})
