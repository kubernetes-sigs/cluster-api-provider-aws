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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/exp/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[unmanaged] [functional]", func() {
	var (
		ctx               context.Context
		result            *clusterctl.ApplyClusterTemplateAndWaitResult
		requiredResources *shared.TestResource
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		result = &clusterctl.ApplyClusterTemplateAndWaitResult{}
	})

	ginkgo.Describe("Create a cluster that uses the external cloud provider", func() {
		ginkgo.It("should create volumes dynamically with external cloud provider", func() {
			specName := "functional-external-cloud-provider"
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1}
			requiredResources.WriteRequestedResources(e2eCtx, "external-cloud-provider-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))

			Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
			Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))
			shared.CreateAWSClusterControllerIdentity(e2eCtx.Environment.BootstrapClusterProxy.GetClient())

			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.Flavor = shared.ExternalCloudProvider
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(1)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			cluster, _, _ := createCluster(ctx, configCluster, result)

			ginkgo.By("Creating the LB service")
			clusterClient := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, clusterName).GetClient()
			lbServiceName := "test-svc-" + util.RandomString(6)
			elbName := createLBService(metav1.NamespaceDefault, lbServiceName, clusterClient)
			verifyElbExists(elbName, true)

			nginxStatefulsetInfo := statefulSetInfo{
				name:                      "nginx-statefulset",
				namespace:                 metav1.NamespaceDefault,
				replicas:                  int32(2),
				selector:                  map[string]string{"app": "nginx"},
				storageClassName:          "aws-ebs-volumes",
				volumeName:                "nginx-volumes",
				svcName:                   "nginx-svc",
				svcPort:                   int32(80),
				svcPortName:               "nginx-web",
				containerName:             "nginx",
				containerImage:            "k8s.gcr.io/nginx-slim:0.8",
				containerPort:             int32(80),
				podTerminationGracePeriod: int64(30),
				volMountPath:              "/usr/share/nginx/html",
			}

			ginkgo.By("Deploying StatefulSet on infra")
			createStatefulSet(nginxStatefulsetInfo, clusterClient)
			awsVolIds := getVolumeIds(nginxStatefulsetInfo, clusterClient)
			verifyVolumesExists(awsVolIds)

			ginkgo.By("Deleting LB service")
			deleteLBService(metav1.NamespaceDefault, lbServiceName, clusterClient)

			deleteCluster(ctx, cluster)
			ginkgo.By("Deleting retained dynamically provisioned volumes")
			deleteRetainedVolumes(awsVolIds)
			ginkgo.By("PASSED!")
		})
	})

	ginkgo.Describe("GPU-enabled cluster test", func() {
		ginkgo.It("should create cluster with single worker", func() {
			specName := "functional-gpu-cluster"
			// Change the multiplier for EC2GPU if GPU type is changed. g4dn.xlarge uses 2 vCPU
			requiredResources = &shared.TestResource{EC2GPU: 2 * 2, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1}
			requiredResources.WriteRequestedResources(e2eCtx, "gpu-test")
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))

			ginkgo.By("Creating cluster with a single worker")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))

			clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ConfigCluster: clusterctl.ConfigClusterInput{
					LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
					ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
					KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
					InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
					Flavor:                   shared.GPUFlavor,
					Namespace:                namespace.Name,
					ClusterName:              clusterName,
					KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
					ControlPlaneMachineCount: pointer.Int64Ptr(1),
					WorkerMachineCount:       pointer.Int64Ptr(1),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
				WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
				// nvidia-gpu flavor creates a config map as part of a crs, that exceeds the annotations size limit when we do kubectl apply.
				// This is because the entire config map is stored in `last-applied` annotation for tracking.
				// The workaround is to use server side apply by passing `--server-side` flag to kubectl apply.
				// More on server side apply here: https://kubernetes.io/docs/reference/using-api/server-side-apply/
				Args: []string{"--server-side"},
			}, result)

			shared.AWSGPUSpec(ctx, e2eCtx, shared.AWSGPUSpecInput{
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				NamespaceName:         namespace.Name,
				ClusterName:           clusterName,
				SkipCleanup:           false,
			})
			ginkgo.By("PASSED!")
		})
	})

	ginkgo.Describe("Multitenancy test", func() {
		ginkgo.It("should create cluster with nested assumed role", func() {
			// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
			specName := "functional-multitenancy-nested"
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
					Flavor:                   shared.NestedMultitenancyFlavor,
					Namespace:                namespace.Name,
					ClusterName:              clusterName,
					KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
					ControlPlaneMachineCount: pointer.Int64Ptr(1),
					WorkerMachineCount:       pointer.Int64Ptr(0),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			}, result)

			ginkgo.By("PASSED!")
		})
	})

	// // TODO: @sedefsavas: Requires env var logic to be removed
	ginkgo.PDescribe("[Serial] Upgrade to main branch Kubernetes", func() {
		ginkgo.Context("in same namespace", func() {
			ginkgo.It("should create the clusters", func() {
				specName := "upgrade-to-main-branch-k8s"
				requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 3, VPC: 1, ClassicLB: 1, EIP: 3}
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
			requiredResources = &shared.TestResource{EC2Normal: 1 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
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
			// By default, first availability zone will be used for cluster resources. This step attempts to create a machine deployment in the second availability zone
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
			configCluster.Flavor = shared.MultiAzFlavor
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

	// TODO @randomvariable: Await more resources
	ginkgo.PDescribe("Multiple workload clusters", func() {
		ginkgo.Context("in different namespaces with machine failures", func() {
			ginkgo.It("should setup namespaces correctly for the two clusters", func() {
				specName := "functional-test-multi-namespace"
				requiredResources = &shared.TestResource{EC2Normal: 4 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 6}
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
				requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 6}
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
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
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

	ginkgo.Describe("Externally managed security groups", func() {
		ginkgo.It("should create a cluster using external security groups", func() {
			specName := "functional-test-external-securitygroups"
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating the VPC and subnets")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			vpc, vpcErr := shared.CreateVPC(e2eCtx, clusterName+"-vpc", "10.0.0.0/22")
			Expect(vpcErr).To(BeNil())
			vpcID := *vpc.VpcId
			shared.SetEnvVar("VPC_ID", vpcID, false)
			pubSubnet, pubSErr := shared.CreateSubnet(e2eCtx, clusterName, "10.0.0.0/24", "", vpcID, "public")
			Expect(pubSErr).To(BeNil())
			pubSubnetID := *pubSubnet.SubnetId
			shared.SetEnvVar("PUBLIC_SUBNET_ID", pubSubnetID, false)
			priSubnet, priSErr := shared.CreateSubnet(e2eCtx, clusterName, "10.0.1.0/24", "", vpcID, "private")
			Expect(priSErr).To(BeNil())
			priSubnetID := *priSubnet.SubnetId
			shared.SetEnvVar("PRIVATE_SUBNET_ID", priSubnetID, false)

			ginkgo.By("Creating security groups")
			cpSG, _ := shared.CreateSecurityGroup(e2eCtx, clusterName+"-controlplane", clusterName+" controlplane", vpcID)
			Expect(cpSG).NotTo(BeNil())
			cpSGID := *cpSG.GroupId
			shared.SetEnvVar("CP_SG_ID", cpSGID, false)
			shared.CreateSecurityGroupIngressRule(e2eCtx, cpSGID, "controlplane default", "0.0.0.0/0", "-1", -1, -1)
			apiSG, _ := shared.CreateSecurityGroup(e2eCtx, clusterName+"-apiserver-lb", clusterName+" apiserver", vpcID)
			Expect(apiSG).NotTo(BeNil())
			apiSGID := *apiSG.GroupId
			shared.SetEnvVar("API_SG_ID", apiSGID, false)
			shared.CreateSecurityGroupIngressRule(e2eCtx, apiSGID, "apiserver default", "0.0.0.0/0", "-1", -1, -1)
			nodeSG, _ := shared.CreateSecurityGroup(e2eCtx, clusterName+"-node", clusterName+" node", vpcID)
			Expect(nodeSG).NotTo(BeNil())
			nodeSGID := *nodeSG.GroupId
			shared.SetEnvVar("NODE_SG_ID", nodeSGID, false)
			shared.CreateSecurityGroupIngressRule(e2eCtx, nodeSGID, "node default", "0.0.0.0/0", "-1", -1, -1)
			lbSG, _ := shared.CreateSecurityGroup(e2eCtx, clusterName+"-lb", clusterName+" load balancer", vpcID)
			Expect(lbSG).NotTo(BeNil())
			lbSGID := *lbSG.GroupId
			shared.SetEnvVar("LB_SG_ID", lbSGID, false)
			shared.CreateSecurityGroupIngressRule(e2eCtx, lbSGID, "load balancer default", "0.0.0.0/0", "-1", -1, -1)

			ginkgo.By("Creating Internet gateway")
			igw, _ := shared.CreateInternetGateway(e2eCtx, clusterName+"-igw")
			Expect(igw).NotTo(BeNil())
			igwID := *igw.InternetGatewayId
			igwA, _ := shared.AttachInternetGateway(e2eCtx, igwID, vpcID)
			Expect(igwA).To(BeTrue())

			ginkgo.By("Allocating Elastic IP")
			eip, _ := shared.AllocateAddress(e2eCtx, clusterName+"-eip")
			Expect(eip).NotTo(BeNil())
			allocationID := *eip.AllocationId

			ginkgo.By("Creating NAT gateway")
			ngw, _ := shared.CreateNatGateway(e2eCtx, clusterName+"-nat", "public", allocationID, pubSubnetID)
			Expect(ngw).NotTo(BeNil())
			ngwID := *ngw.NatGatewayId
			shared.WaitForNatGatewayState(e2eCtx, ngwID, 180, "available")

			ginkgo.By("Creating route tables")
			pubRoute, _ := shared.CreateRouteTable(e2eCtx, clusterName+"-rt-public", vpcID)
			priRoute, _ := shared.CreateRouteTable(e2eCtx, clusterName+"-rt-private", vpcID)

			ginkgo.By("Creating associating routes")
			pubRouteID := *pubRoute.RouteTableId
			priRouteID := *priRoute.RouteTableId
			pubAssocRT, _ := shared.AssociateRouteTable(e2eCtx, pubRouteID, pubSubnetID)
			priAssocRT, _ := shared.AssociateRouteTable(e2eCtx, priRouteID, priSubnetID)

			ginkgo.By("Creating routes")
			shared.CreateRoute(e2eCtx, pubRouteID, "0.0.0.0/0", nil, igw.InternetGatewayId, nil)
			shared.CreateRoute(e2eCtx, priRouteID, "0.0.0.0/0", ngw.NatGatewayId, nil, nil)

			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			configCluster.Flavor = "external-securitygroups"
			cluster, md, _ := createCluster(ctx, configCluster, result)

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

			deleteCluster(ctx, cluster)

			ginkgo.By("Disassociating route tables")
			pubAssocID := *pubAssocRT.AssociationId
			priAssocID := *priAssocRT.AssociationId
			Expect(shared.DisassociateRouteTable(e2eCtx, pubAssocID)).To(BeTrue())
			Expect(shared.DisassociateRouteTable(e2eCtx, priAssocID)).To(BeTrue())

			ginkgo.By("Deleting route tables")
			Expect(shared.DeleteRouteTable(e2eCtx, priRouteID)).To(BeTrue())
			Expect(shared.DeleteRouteTable(e2eCtx, pubRouteID)).To(BeTrue())

			ginkgo.By("Deleting NAT gateway")
			Expect(shared.DeleteNatGateway(e2eCtx, ngwID)).To(BeTrue())
			shared.WaitForNatGatewayState(e2eCtx, ngwID, 180, "deleted")

			ginkgo.By("Releasing Elastic IP")
			Expect(shared.ReleaseAddress(e2eCtx, allocationID)).To(BeTrue())

			ginkgo.By("Detaching Internet gateway")
			Eventually(shared.DetachInternetGateway(e2eCtx, igwID, vpcID), 30*time.Second).Should(BeTrue())

			ginkgo.By("Deleting Internet gateway")
			Expect(shared.DeleteInternetGateway(e2eCtx, igwID)).To(BeTrue())

			ginkgo.By("Deleting security groups")
			Expect(shared.DeleteSecurityGroup(e2eCtx, cpSGID)).To(BeTrue())
			Expect(shared.DeleteSecurityGroup(e2eCtx, apiSGID)).To(BeTrue())
			Expect(shared.DeleteSecurityGroup(e2eCtx, nodeSGID)).To(BeTrue())
			Expect(shared.DeleteSecurityGroup(e2eCtx, lbSGID)).To(BeTrue())

			ginkgo.By("Deleting subnets")
			Expect(shared.DeleteSubnet(e2eCtx, priSubnetID)).To(BeTrue())
			Expect(shared.DeleteSubnet(e2eCtx, pubSubnetID)).To(BeTrue())

			ginkgo.By("Deleting the VPC")
			Expect(shared.DeleteVPC(e2eCtx, vpcID)).To(BeTrue())
		})
	})

	ginkgo.Describe("Peerings and internal ELB", func() {
		ginkgo.It("should create external clusters in peered VPC and with an internal ELB", func() {
			specName := "functional-test-peered-internal-elb"
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 5}
			requiredResources.WriteRequestedResources(e2eCtx, specName)
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
			defer shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)

			ginkgo.By("Creating the management cluster infrastructure")
			mgmtClusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			mgmtClusterInfra := new(shared.AWSInfrastructure)
			mgmtClusterInfra.New(shared.AWSInfrastructureSpec{
				ClusterName:       mgmtClusterName,
				VpcCidr:           "10.0.0.0/23",
				PublicSubnetCidr:  "10.0.0.0/24",
				PrivateSubnetCidr: "10.0.1.0/24",
			}, e2eCtx)
			mgmtClusterInfra.CreateInfrastructure()
			Expect(mgmtClusterInfra.VPC).NotTo(BeNil())
			Expect(*mgmtClusterInfra.State.VpcState).To(Equal("available"))
			Expect(len(mgmtClusterInfra.Subnets)).To(Equal(2))
			Expect(mgmtClusterInfra.InternetGateway).NotTo(BeNil())
			Expect(mgmtClusterInfra.ElasticIP).NotTo(BeNil())
			Expect(mgmtClusterInfra.NatGateway).NotTo(BeNil())
			Expect(len(mgmtClusterInfra.RouteTables)).To(Equal(2))

			ginkgo.By("Creating the workload cluster infrastructure")
			wlClusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			wlClusterInfra := new(shared.AWSInfrastructure)
			wlClusterInfra.New(shared.AWSInfrastructureSpec{
				ClusterName:       wlClusterName,
				VpcCidr:           "10.0.2.0/23",
				PublicSubnetCidr:  "10.0.2.0/24",
				PrivateSubnetCidr: "10.0.3.0/24",
			}, e2eCtx)
			wlClusterInfra.CreateInfrastructure()
			Expect(wlClusterInfra.VPC).NotTo(BeNil())
			Expect(*wlClusterInfra.State.VpcState).To(Equal("available"))
			Expect(len(wlClusterInfra.Subnets)).To(Equal(2))
			Expect(wlClusterInfra.InternetGateway).NotTo(BeNil())
			Expect(wlClusterInfra.ElasticIP).NotTo(BeNil())
			Expect(wlClusterInfra.NatGateway).NotTo(BeNil())
			Expect(len(wlClusterInfra.RouteTables)).To(Equal(2))

			shared.SetEnvVar("MGMT_VPC_ID", *mgmtClusterInfra.VPC.VpcId, false)
			shared.SetEnvVar("WL_VPC_ID", *wlClusterInfra.VPC.VpcId, false)
			shared.SetEnvVar("MGMT_PUBLIC_SUBNET_ID", *mgmtClusterInfra.State.PublicSubnetID, false)
			shared.SetEnvVar("MGMT_PRIVATE_SUBNET_ID", *mgmtClusterInfra.State.PrivateSubnetID, false)
			shared.SetEnvVar("WL_PUBLIC_SUBNET_ID", *wlClusterInfra.State.PublicSubnetID, false)
			shared.SetEnvVar("WL_PRIVATE_SUBNET_ID", *wlClusterInfra.State.PrivateSubnetID, false)

			ginkgo.By("Creating VPC peerings")
			cPeering, _ := shared.CreatePeering(e2eCtx, mgmtClusterName+"-"+wlClusterName, *mgmtClusterInfra.VPC.VpcId, *wlClusterInfra.VPC.VpcId)
			Expect(cPeering).NotTo(BeNil())
			Eventually(func() bool {
				aPeering, err := shared.AcceptPeering(e2eCtx, *cPeering.VpcPeeringConnectionId)
				if err != nil {
					return false
				}
				return aPeering != nil
			}, 60*time.Second).Should(BeTrue())

			ginkgo.By("Creating routes for peerings")
			shared.CreateRoute(e2eCtx, *mgmtClusterInfra.State.PublicRouteTableID, "10.0.2.0/23", nil, nil, cPeering.VpcPeeringConnectionId)
			shared.CreateRoute(e2eCtx, *mgmtClusterInfra.State.PrivateRouteTableID, "10.0.2.0/23", nil, nil, cPeering.VpcPeeringConnectionId)
			shared.CreateRoute(e2eCtx, *wlClusterInfra.State.PublicRouteTableID, "10.0.0.0/23", nil, nil, cPeering.VpcPeeringConnectionId)
			shared.CreateRoute(e2eCtx, *wlClusterInfra.State.PrivateRouteTableID, "10.0.0.0/23", nil, nil, cPeering.VpcPeeringConnectionId)

			ginkgo.By("Creating a management cluster in a peered VPC")
			mgmtConfigCluster := defaultConfigCluster(mgmtClusterName, namespace.Name)
			mgmtConfigCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			mgmtConfigCluster.Flavor = "peered-remote"
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

			mgmtClusterProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, mgmtCluster.Namespace, mgmtCluster.Name)

			shared.Byf("Creating a namespace for hosting the %s test spec", specName)
			mgmtNamespace := framework.CreateNamespace(ctx, framework.CreateNamespaceInput{
				Creator: mgmtClusterProxy.GetClient(),
				Name:    namespace.Name,
			})

			ginkgo.By("Initializing the management cluster")
			clusterctl.InitManagementClusterAndWatchControllerLogs(ctx, clusterctl.InitManagementClusterAndWatchControllerLogsInput{
				ClusterProxy:            mgmtClusterProxy,
				ClusterctlConfigPath:    e2eCtx.Environment.ClusterctlConfigPath,
				InfrastructureProviders: e2eCtx.E2EConfig.InfrastructureProviders(),
				LogFolder:               filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", mgmtCluster.Name),
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-controllers")...)

			ginkgo.By("Ensure API servers are stable before doing the move")
			Consistently(func() error {
				kubeSystem := &corev1.Namespace{}
				return e2eCtx.Environment.BootstrapClusterProxy.GetClient().Get(ctx, client.ObjectKey{Name: "kube-system"}, kubeSystem)
			}, "5s", "100ms").Should(BeNil(), "Failed to assert bootstrap API server stability")
			Consistently(func() error {
				kubeSystem := &corev1.Namespace{}
				return mgmtClusterProxy.GetClient().Get(ctx, client.ObjectKey{Name: "kube-system"}, kubeSystem)
			}, "5s", "100ms").Should(BeNil(), "Failed to assert management API server stability")

			ginkgo.By("Moving the management cluster to be self hosted")
			clusterctl.Move(ctx, clusterctl.MoveInput{
				LogFolder:            filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", "bootstrap"),
				ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
				FromKubeconfigPath:   e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
				ToKubeconfigPath:     mgmtClusterProxy.GetKubeconfigPath(),
				Namespace:            namespace.Name,
			})

			mgmtCluster = framework.DiscoveryAndWaitForCluster(ctx, framework.DiscoveryAndWaitForClusterInput{
				Getter:    mgmtClusterProxy.GetClient(),
				Namespace: mgmtNamespace.Name,
				Name:      mgmtCluster.Name,
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster")...)

			mgmtControlPlane := framework.GetKubeadmControlPlaneByCluster(ctx, framework.GetKubeadmControlPlaneByClusterInput{
				Lister:      mgmtClusterProxy.GetClient(),
				ClusterName: mgmtCluster.Name,
				Namespace:   mgmtCluster.Namespace,
			})
			Expect(mgmtControlPlane).ToNot(BeNil())

			ginkgo.By("Creating a namespace to host the internal-elb spec")
			wlNamespace := framework.CreateNamespace(ctx, framework.CreateNamespaceInput{
				Creator: mgmtClusterProxy.GetClient(),
				Name:    wlClusterName,
			})

			ginkgo.By("Creating workload cluster with internal ELB")
			wlConfigCluster := defaultConfigCluster(wlClusterName, wlNamespace.Name)
			wlConfigCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			wlConfigCluster.Flavor = "internal-elb"
			wlResult := &clusterctl.ApplyClusterTemplateAndWaitResult{}
			clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy:                 mgmtClusterProxy,
				ConfigCluster:                wlConfigCluster,
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals("", "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals("", "wait-control-plane"),
				WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes"),
			}, wlResult)

			wlWM := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            mgmtClusterProxy.GetClient(),
				ClusterName:       mgmtClusterName,
				Namespace:         wlNamespace.Name,
				MachineDeployment: *wlResult.MachineDeployments[0],
			})
			wlCPM := framework.GetControlPlaneMachinesByCluster(ctx, framework.GetControlPlaneMachinesByClusterInput{
				Lister:      mgmtClusterProxy.GetClient(),
				ClusterName: wlClusterName,
				Namespace:   wlNamespace.Name,
			})
			Expect(len(wlWM)).To(Equal(1))
			Expect(len(wlCPM)).To(Equal(1))

			ginkgo.By("Deleting the workload cluster")
			framework.DeleteCluster(ctx, framework.DeleteClusterInput{
				Deleter: mgmtClusterProxy.GetClient(),
				Cluster: wlResult.Cluster,
			})

			framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
				Getter:  mgmtClusterProxy.GetClient(),
				Cluster: wlResult.Cluster,
			}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

			ginkgo.By("Moving the management cluster back to bootstrap")
			clusterctl.Move(ctx, clusterctl.MoveInput{
				LogFolder:            filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", mgmtCluster.Name),
				ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
				FromKubeconfigPath:   mgmtClusterProxy.GetKubeconfigPath(),
				ToKubeconfigPath:     e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
				Namespace:            namespace.Name,
			})

			mgmtCluster = framework.DiscoveryAndWaitForCluster(ctx, framework.DiscoveryAndWaitForClusterInput{
				Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				Namespace: mgmtNamespace.Name,
				Name:      mgmtCluster.Name,
			}, e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster")...)

			mgmtControlPlane = framework.GetKubeadmControlPlaneByCluster(ctx, framework.GetKubeadmControlPlaneByClusterInput{
				Lister:      e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
				ClusterName: mgmtCluster.Name,
				Namespace:   mgmtCluster.Namespace,
			})
			Expect(mgmtControlPlane).ToNot(BeNil())

			ginkgo.By("Deleting the management cluster")
			deleteCluster(ctx, mgmtCluster)

			ginkgo.By("Deleting peering connection")
			Expect(shared.DeletePeering(e2eCtx, *cPeering.VpcPeeringConnectionId)).To(BeTrue())

			ginkgo.By("Deleting the workload cluster infrastructure")
			wlClusterInfra.DeleteInfrastructure()

			ginkgo.By("Deleting the management cluster infrastructure")
			mgmtClusterInfra.DeleteInfrastructure()
		})
	})
})
