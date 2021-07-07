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

package unmanaged

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/exp/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [functional]", func() {
	var (
		ctx    context.Context
		result *clusterctl.ApplyClusterTemplateAndWaitResult
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}
	})

	ginkgo.Describe("Multitenancy test", func() {

		ginkgo.It("should create cluster with nested assumed role", func() {
			// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
			specName := "functional-multitenancy-nested"
			requiredResources := &shared.TestResource{EC2: 1, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1}
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
					Flavor:                   shared.NestedMultitenancyFlavor,
					Namespace:                namespace.Name,
					ClusterName:              clusterName,
					KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
					ControlPlaneMachineCount: pointer.Int64Ptr(1),
					WorkerMachineCount:       pointer.Int64Ptr(0),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
				WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			}, result)

			ginkgo.By("PASSED!")
		})
	})

	// // TODO: @sedefsavas: Requires env var logic to be removed
	ginkgo.PDescribe("[Serial] Upgrade to main branch Kubernetes", func() {
		ginkgo.Context("in same namespace", func() {
			ginkgo.It("should create the clusters", func() {
				specName := "upgrade-to-main-branch-k8s"
				requiredResources := &shared.TestResource{EC2: 2, IGW: 1, NGW: 3, VPC: 1, ClassicLB: 1, EIP: 3}
				requiredResources.WriteRequestedResources(e2eCtx, "upgrade-to-master-test")
				Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
				defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
				namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
				defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
				ginkgo.By("Creating first cluster with single control plane")
				cluster1Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				shared.SetEnvVar("USE_CI_ARTIFACTS", "true", false)
				tagPrefix := "v"
				searchSemVer, err := semver.Make(strings.TrimPrefix(e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion), tagPrefix))
				Expect(err).NotTo(HaveOccurred())

				shared.SetEnvVar(shared.KubernetesVersion, "v"+searchSemVer.String(), false)
				configCluster := defaultConfigCluster(cluster1Name, namespace.Name)

				configCluster.Flavor = shared.UpgradeToMain
				configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
				createCluster(ctx, configCluster, result)

				kubernetesUgradeVersion, err := LatestCIReleaseForVersion("v" + searchSemVer.String())
				Expect(err).NotTo(HaveOccurred())
				configCluster.KubernetesVersion = kubernetesUgradeVersion
				configCluster.Flavor = "upgrade-ci-artifacts"
				cluster2, md, kcp := createCluster(ctx, configCluster, result)

				ginkgo.By(fmt.Sprintf("Waiting for Kubernetes versions of machines in MachineDeployment %s/%s to be upgraded from %s to %s",
					md[0].Namespace, md[0].Name, e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion), kubernetesUgradeVersion))

				framework.WaitForMachineDeploymentMachinesToBeUpgraded(ctx, framework.WaitForMachineDeploymentMachinesToBeUpgradedInput{
					Lister:                   e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					Cluster:                  cluster2,
					MachineCount:             int(*md[0].Spec.Replicas),
					KubernetesUpgradeVersion: kubernetesUgradeVersion,
					MachineDeployment:        *md[0],
				}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade")...)

				ginkgo.By("Waiting for control-plane machines to have the upgraded kubernetes version")
				framework.WaitForControlPlaneMachinesToBeUpgraded(ctx, framework.WaitForControlPlaneMachinesToBeUpgradedInput{
					Lister:                   e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					Cluster:                  cluster2,
					MachineCount:             int(*kcp.Spec.Replicas),
					KubernetesUpgradeVersion: kubernetesUgradeVersion,
				}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-machine-upgrade")...)

				ginkgo.By("Deleting the Clusters")
				shared.SetEnvVar("USE_CI_ARTIFACTS", "false", false)
				deleteCluster(ctx, cluster2)
			})
		})
	})

	ginkgo.Describe("Workload cluster with AWS SSM Parameter as the Secret Backend", func() {
		ginkgo.It("should be creatable and deletable", func() {
			specName := "functional-test-ssm-parameter-store"
			requiredResources := &shared.TestResource{EC2: 2, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
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
			configCluster.Flavor = shared.SSMFlavor
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

	ginkgo.Describe("MachineDeployment misconfigurations", func() {
		ginkgo.It("MachineDeployment misconfigurations", func() {
			specName := "functional-test-md-misconfigurations"
			requiredResources := &shared.TestResource{EC2: 1, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			_, _, _ = createCluster(ctx, configCluster, result)

			ginkgo.By("Creating Machine Deployment with invalid subnet ID")
			md1Name := clusterName + "-md-1"
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       makeMachineDeployment(namespace.Name, md1Name, clusterName, 1),
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, md1Name),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, md1Name, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), nil, pointer.StringPtr("invalid-subnet")),
			})

			ginkgo.By("Looking for failure event to be reported")
			Eventually(func() bool {
				eventList := getEvents(namespace.Name)
				subnetError := "Failed to create instance: failed to run instance: InvalidSubnetID.NotFound: " +
					"The subnet ID '%s' does not exist"
				return isErrorEventExists(namespace.Name, md1Name, "FailedCreate", fmt.Sprintf(subnetError, "invalid-subnet"), eventList)
			}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...).Should(BeTrue())

			ginkgo.By("Creating Machine Deployment in non-configured Availability Zone")
			md2Name := clusterName + "-md-2"
			//By default, first availability zone will be used for cluster resources. This step attempts to create a machine deployment in the second availability zone
			invalidAz := shared.GetAvailabilityZones(e2eCtx.AWSSession)[1].ZoneName
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       makeMachineDeployment(namespace.Name, md2Name, clusterName, 1),
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, md2Name),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, md2Name, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), invalidAz, nil),
			})

			ginkgo.By("Looking for failure event to be reported")
			Eventually(func() bool {
				eventList := getEvents(namespace.Name)
				azError := "Failed to create instance: no subnets available in availability zone \"%s\""
				return isErrorEventExists(namespace.Name, md2Name, "FailedCreate", fmt.Sprintf(azError, *invalidAz), eventList)
			}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...).Should(BeTrue())
		})
	})

	ginkgo.Describe("Workload cluster in multiple AZs", func() {
		ginkgo.It("It should be creatable and deletable", func() {
			specName := "functional-test-multi-az"
			requiredResources := &shared.TestResource{EC2: 3, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(3)
			configCluster.Flavor = shared.MultiAzFlavor
			cluster, _, _ := createCluster(ctx, configCluster, result)

			ginkgo.By("Adding worker nodes to additional subnets")
			mdName1 := clusterName + "-md-1"
			mdName2 := clusterName + "-md-2"
			md1 := makeMachineDeployment(namespace.Name, mdName1, clusterName, 1)
			md2 := makeMachineDeployment(namespace.Name, mdName2, clusterName, 1)
			az1 := os.Getenv(shared.AwsAvailabilityZone1)
			az2 := os.Getenv(shared.AwsAvailabilityZone2)

			//private CIDRs set in cluster-template-multi-az.yaml
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       md1,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName1),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName1, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), pointer.StringPtr(az1), getSubnetId("cidr-block", "10.0.0.0/24")),
			})
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				MachineDeployment:       md2,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName2),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName2, e2eCtx.E2EConfig.GetVariable(shared.AwsNodeMachineType), pointer.StringPtr(az2), getSubnetId("cidr-block", "10.0.2.0/24")),
			})

			ginkgo.By("Waiting for new worker nodes to become ready")
			k8sClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md1}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md2}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		})
	})

	// TODO @randomvariable: Await more resources
	ginkgo.PDescribe("Multiple workload clusters", func() {
		ginkgo.Context("in different namespaces with machine failures", func() {
			ginkgo.It("should setup namespaces correctly for the two clusters", func() {
				specName := "functional-test-multi-namespace"
				requiredResources := &shared.TestResource{EC2: 4, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 6}
				requiredResources.WriteRequestedResources(e2eCtx, specName)
				Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
				defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))

				ginkgo.By("Creating first cluster with single control plane")
				ns1, cf1 := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
					Creator:   e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					ClientSet: e2eCtx.Environment.BootstrapClusterProxy.GetClientSet(),
					Name:      fmt.Sprintf("functional-multi-namespace-1-%s", util.RandomString(6)),
					LogFolder: filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
				})
				e2eCtx.Environment.Namespaces[ns1] = cf1
				ns2, cf2 := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
					Creator:   e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					ClientSet: e2eCtx.Environment.BootstrapClusterProxy.GetClientSet(),
					Name:      fmt.Sprintf("functional-multi-namespace-2-%s", util.RandomString(6)),
					LogFolder: filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
				})
				e2eCtx.Environment.Namespaces[ns2] = cf2

				ginkgo.By("Creating first cluster")
				cluster1Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster := defaultConfigCluster(cluster1Name, ns1.Name)
				configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
				configCluster.Flavor = shared.LimitAzFlavor
				cluster1, md1, _ := createCluster(ctx, configCluster, result)
				Expect(len(md1)).To(Equal(1), "Expecting one MachineDeployment")

				ginkgo.By("Deleting a worker node machine")
				deleteMachine(ns1, md1[0])
				time.Sleep(10 * time.Second)

				ginkgo.By("Verifying MachineDeployment is running.")
				framework.DiscoveryAndWaitForMachineDeployments(ctx, framework.DiscoveryAndWaitForMachineDeploymentsInput{Cluster: cluster1, Lister: e2eCtx.Environment.BootstrapClusterProxy.GetClient()}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)

				ginkgo.By("Creating second cluster")
				cluster2Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster = defaultConfigCluster(cluster2Name, ns2.Name)
				configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
				configCluster.Flavor = shared.LimitAzFlavor
				cluster2, md2, _ := createCluster(ctx, configCluster, result)
				Expect(len(md2)).To(Equal(1), "Expecting one MachineDeployment")

				ginkgo.By("Deleting node directly from infra cloud")
				machines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
					Lister:            e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					ClusterName:       cluster1Name,
					Namespace:         ns2.Name,
					MachineDeployment: *md2[0],
				})
				Expect(len(machines)).Should(BeNumerically(">", 0))
				terminateInstance(*machines[0].Spec.ProviderID)

				ginkgo.By("Waiting for AWSMachine to be labelled as terminated")
				Eventually(func() bool {
					machineList := getAWSMachinesForDeployment(ns2.Name, *md2[0])
					labels := machineList.Items[0].GetLabels()
					return labels[instancestate.Ec2InstanceStateLabelKey] == string(infrav1.InstanceStateTerminated)
				}, e2eCtx.E2EConfig.GetIntervals("", "wait-machine-status")...).Should(Equal(true))

				ginkgo.By("Waiting for machine to reach Failed state")
				statusChecks := []framework.MachineStatusCheck{framework.MachinePhaseCheck(string(clusterv1.MachinePhaseFailed))}
				machineStatusInput := framework.WaitForMachineStatusCheckInput{
					Getter:       e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
					Machine:      &machines[0],
					StatusChecks: statusChecks,
				}
				framework.WaitForMachineStatusCheck(ctx, machineStatusInput, e2eCtx.E2EConfig.GetIntervals("", "wait-machine-status")...)

				ginkgo.By("Deleting the clusters and namespaces")
				deleteCluster(ctx, cluster1)
				deleteCluster(ctx, cluster2)
				framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(), Name: ns1.Name})
				framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(), Name: ns2.Name})
				cf1()
				cf2()
			})
		})

		ginkgo.Context("Defining clusters in the same namespace", func() {
			specName := "functional-test-multi-cluster-single-namespace"
			ginkgo.It("should create the clusters", func() {
				requiredResources := &shared.TestResource{EC2: 2, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 6}
				requiredResources.WriteRequestedResources(e2eCtx, specName)
				Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
				defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
				namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
				defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
				ginkgo.By("Creating first cluster with single control plane")
				cluster1Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster := defaultConfigCluster(cluster1Name, namespace.Name)
				configCluster.Flavor = shared.LimitAzFlavor
				cluster1, _, _ := createCluster(ctx, configCluster, result)

				ginkgo.By("Creating second cluster with single control plane")
				cluster2Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster = defaultConfigCluster(cluster2Name, namespace.Name)
				configCluster.Flavor = shared.LimitAzFlavor
				cluster2, _, _ := createCluster(ctx, configCluster, result)

				ginkgo.By("Deleting the Clusters")
				deleteCluster(ctx, cluster1)
				deleteCluster(ctx, cluster2)
			})
		})
	})

	ginkgo.Describe("Workload cluster with spot instances", func() {
		ginkgo.It("should be creatable and deletable", func() {
			specName := "functional-test-spot-instances"
			requiredResources := &shared.TestResource{EC2: 2, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
			ginkgo.By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			configCluster.Flavor = shared.SpotInstancesFlavor
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
			assertSpotInstanceType(*workerMachines[0].Spec.ProviderID)
			Expect(len(controlPlaneMachines)).To(Equal(1))
		})
	})

})
