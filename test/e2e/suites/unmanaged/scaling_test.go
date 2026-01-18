//go:build e2e
// +build e2e

/*
Copyright 2026 The Kubernetes Authors.

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

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [scaling]", func() {
	var (
		ctx               context.Context
		result            *clusterctl.ApplyClusterTemplateAndWaitResult
		requiredResources *shared.TestResource
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}
	})

	ginkgo.Describe("MachineDeployment scaling operations", func() {
		ginkgo.It("should support scaling up, down, and to zero", func() {
			specName := "functional-scaling"
			if !e2eCtx.Settings.SkipQuotas {
				requiredResources = &shared.TestResource{EC2Normal: 3 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
				requiredResources.WriteRequestedResources(e2eCtx, specName)
				Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
				defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
			}

			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating a cluster with control plane")
			clusterName := fmt.Sprintf("%s-%s", specName, util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = ptr.To[int64](1)
			configCluster.WorkerMachineCount = ptr.To[int64](0)
			cluster, _, _ := createCluster(ctx, configCluster, result)

			ginkgo.By("Verifying AWSMachineTemplate has capacity and nodeInfo set for autoscaling from zero")
			Eventually(func(g Gomega) {
				awsMachineTemplateList := &infrav1.AWSMachineTemplateList{}
				g.Expect(e2eCtx.Environment.BootstrapClusterProxy.GetClient().List(ctx, awsMachineTemplateList, client.InNamespace(namespace.Name))).To(Succeed())
				g.Expect(awsMachineTemplateList.Items).ToNot(BeEmpty())

				for _, template := range awsMachineTemplateList.Items {
					capacity := template.Status.Capacity
					_, hasCPU := capacity[corev1.ResourceCPU]
					_, hasMemory := capacity[corev1.ResourceMemory]
					g.Expect(hasCPU).To(BeTrue(), "Expected AWSMachineTemplate %s to have .status.capacity for CPU set", template.Name)
					g.Expect(hasMemory).To(BeTrue(), "Expected AWSMachineTemplate %s to have .status.capacity for memory set", template.Name)
					g.Expect(template.Status.NodeInfo).ToNot(BeNil(), "Expected AWSMachineTemplate %s to have .status.nodeInfo set", template.Name)
					g.Expect(template.Status.NodeInfo.Architecture).ToNot(BeEmpty(), "Expected AWSMachineTemplate %s to have .status.nodeInfo.architecture set", template.Name)
					g.Expect(template.Status.NodeInfo.OperatingSystem).ToNot(BeEmpty(), "Expected AWSMachineTemplate %s to have .status.nodeInfo.operatingSystem set", template.Name)
				}
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-deployment")...).Should(Succeed())

			ginkgo.By("Creating a MachineDeployment with annotations for autoscaler")
			mdName := clusterName + "-md-0"
			machineTemplate := makeAWSMachineTemplate(namespace.Name, mdName, e2eCtx.E2EConfig.MustGetVariable(shared.AwsNodeMachineType), nil)

			machineDeployment := makeMachineDeployment(namespace.Name, mdName, clusterName, nil, 0)
			machineDeployment.Annotations = map[string]string{
				"cluster.x-k8s.io/cluster-api-autoscaler-node-group-max-size": "3",
				"cluster.x-k8s.io/cluster-api-autoscaler-node-group-min-size": "0",
			}

			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       machineDeployment,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName),
				InfraMachineTemplate:    machineTemplate,
			})

			ginkgo.By("Verifying MachineDeployment is created with 0 replicas")
			Eventually(func(g Gomega) {
				md := &clusterv1.MachineDeployment{}
				key := client.ObjectKey{Namespace: namespace.Name, Name: mdName}
				g.Expect(e2eCtx.Environment.BootstrapClusterProxy.GetClient().Get(ctx, key, md)).To(Succeed())
				g.Expect(md.Spec.Replicas).ToNot(BeNil())
				g.Expect(*md.Spec.Replicas).To(Equal(int32(0)))
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-deployment")...).Should(Succeed())

			ginkgo.By("Scaling MachineDeployment from 0 to 1")
			framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
				ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
				Cluster:                   cluster,
				MachineDeployment:         machineDeployment,
				Replicas:                  1,
				WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("Verifying 1 worker node exists")
			workerMachines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *machineDeployment,
			})
			Expect(len(workerMachines)).To(Equal(1))

			ginkgo.By("Scaling MachineDeployment from 1 to 3")
			framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
				ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
				Cluster:                   cluster,
				MachineDeployment:         machineDeployment,
				Replicas:                  3,
				WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("Verifying 3 worker nodes exist")
			workerMachines = framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *machineDeployment,
			})
			Expect(len(workerMachines)).To(Equal(3))

			ginkgo.By("Verifying all machines are running")
			for _, machine := range workerMachines {
				Expect(machine.Status.Phase).To(Equal(clusterv1.MachinePhaseRunning))
			}

			ginkgo.By("Scaling MachineDeployment from 3 to 1")
			framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
				ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
				Cluster:                   cluster,
				MachineDeployment:         machineDeployment,
				Replicas:                  1,
				WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("Verifying 1 worker node exists after scaling down")
			workerMachines = framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *machineDeployment,
			})
			Expect(len(workerMachines)).To(Equal(1))

			ginkgo.By("Scaling MachineDeployment from 1 to 0")
			framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
				ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
				Cluster:                   cluster,
				MachineDeployment:         machineDeployment,
				Replicas:                  0,
				WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("Verifying no worker nodes exist after scaling to 0")
			Eventually(func(g Gomega) {
				machines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
					Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					ClusterName:       clusterName,
					Namespace:         namespace.Name,
					MachineDeployment: *machineDeployment,
				})
				g.Expect(len(machines)).To(Equal(0), "Expected 0 machines after scaling to zero")
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes")...).Should(Succeed())

			ginkgo.By("Scaling back up from 0 to 2 to verify scale-from-zero works")
			framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
				ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
				Cluster:                   cluster,
				MachineDeployment:         machineDeployment,
				Replicas:                  2,
				WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			})

			ginkgo.By("Verifying 2 worker nodes exist after scaling from zero")
			workerMachines = framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *machineDeployment,
			})
			Expect(len(workerMachines)).To(Equal(2))

			ginkgo.By("Verifying all machines are running after scaling from zero")
			for _, machine := range workerMachines {
				Expect(machine.Status.Phase).To(Equal(clusterv1.MachinePhaseRunning))
			}

			ginkgo.By("PASSED!")
		})
	})

	ginkgo.Describe("MachineDeployment scaling stress test", func() {
		ginkgo.It("should handle rapid scaling operations", func() {
			specName := "functional-scaling-stress"
			if !e2eCtx.Settings.SkipQuotas {
				requiredResources = &shared.TestResource{EC2Normal: 5 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
				requiredResources.WriteRequestedResources(e2eCtx, specName)
				Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
				defer shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
			}

			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("%s-%s", specName, util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = ptr.To[int64](1)
			configCluster.WorkerMachineCount = ptr.To[int64](1)
			cluster, md, _ := createCluster(ctx, configCluster, result)
			Expect(len(md)).To(Equal(1))

			scalingSequence := []int32{3, 1, 5, 2, 0, 3, 1, 0}
			for _, targetReplicas := range scalingSequence {
				ginkgo.By(fmt.Sprintf("Scaling MachineDeployment to %d replicas", targetReplicas))
				framework.ScaleAndWaitMachineDeployment(ctx, framework.ScaleAndWaitMachineDeploymentInput{
					ClusterProxy:              e2eCtx.Environment.BootstrapClusterProxy,
					Cluster:                   cluster,
					MachineDeployment:         md[0],
					Replicas:                  targetReplicas,
					WaitForMachineDeployments: e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
				})

				ginkgo.By(fmt.Sprintf("Verifying MachineDeployment has %d replicas", targetReplicas))
				Eventually(func(g Gomega) {
					machines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
						Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
						ClusterName:       clusterName,
						Namespace:         namespace.Name,
						MachineDeployment: *md[0],
					})
					g.Expect(len(machines)).To(Equal(int(targetReplicas)), "Expected %d machines", targetReplicas)
				}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes")...).Should(Succeed())
			}

			ginkgo.By("PASSED!")
		})
	})
})
