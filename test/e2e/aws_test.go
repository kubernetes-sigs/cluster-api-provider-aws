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

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type statefulSetInfo struct {
	name                      string
	namespace                 string
	replicas                  int32
	selector                  map[string]string
	storageClassName          string
	volumeName                string
	svcName                   string
	svcPort                   int32
	svcPortName               string
	containerName             string
	containerImage            string
	containerPort             int32
	podTerminationGracePeriod int64
	volMountPath              string
}

var _ = Describe("functional tests", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
		specName  = "functional-tests"
	)

	BeforeEach(func() {
		Expect(bootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = setupSpecNamespace(ctx, specName, bootstrapClusterProxy, artifactFolder)
		Expect(e2eConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eConfig.Variables).To(HaveKey(KubernetesVersion))
	})

	Describe("Workload cluster with AWS SSM Parameter as the Secret Backend", func() {
		It("It should be creatable and deletable", func() {
			By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(1)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			configCluster.Flavor = SSMFlavor
			_, md := createCluster(ctx, configCluster)

			workerMachines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            bootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *md[0],
			})
			controlPlaneMachines := framework.GetControlPlaneMachinesByCluster(ctx, framework.GetControlPlaneMachinesByClusterInput{
				Lister:      bootstrapClusterProxy.GetClient(),
				ClusterName: clusterName,
				Namespace:   namespace.Name,
			})
			Expect(len(workerMachines)).To(Equal(1))
			Expect(len(controlPlaneMachines)).To(Equal(1))
		})
	})

	Describe("Cluster name validations and provisioning extra AWS resources", func() {
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
		It("Should create a cluster, AWS load balancer, and volume", func() {
			By("Creating cluster with name starting with 'sg-', having '.' and more than 22 characters")
			// Tests a cluster name that satisfies the following cases:
			// - name with more than 22 characters
			// - name with '.'
			// - name that starts with 'sg-'
			clusterName := fmt.Sprintf("sg-test.%s", util.RandomString(20))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			cluster, md := createCluster(ctx, configCluster)
			clusterClient := bootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, clusterName).GetClient()

			By("Waiting for worker nodes to be in Running phase")
			machines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            bootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *md[0],
			})

			Expect(len(machines)).Should(BeNumerically(">", 0))
			statusChecks := []framework.MachineStatusCheck{framework.MachinePhaseCheck(string(clusterv1.MachinePhaseRunning))}
			machineStatusInput := framework.WaitForMachineStatusCheckInput{
				Getter:       bootstrapClusterProxy.GetClient(),
				Machine:      &machines[0],
				StatusChecks: statusChecks,
			}
			framework.WaitForMachineStatusCheck(ctx, machineStatusInput, e2eConfig.GetIntervals("", "wait-machine-status")...)

			By("Creating the LB service")
			lbServiceName := "test-svc-" + util.RandomString(6)
			elbName := createLBService(metav1.NamespaceDefault, lbServiceName, clusterClient)
			verifyElbExists(elbName, true)

			By("Deploying StatefulSet on infra")
			createStatefulSet(nginxStatefulsetInfo, clusterClient)
			awsVolIds := getVolumeIds(nginxStatefulsetInfo, clusterClient)
			verifyVolumesExists(awsVolIds)

			By("Deleting the Cluster")
			deleteCluster(ctx, cluster)

			By("Verifying whether provisioned LB deleted")
			verifyElbExists(elbName, false)

			By("Verifying dynamically provisioned volumes retention")
			verifyVolumesExists(awsVolIds)

			By("Deleting retained dynamically provisioned volumes")
			deleteRetainedVolumes(awsVolIds)
		})
	})

	Describe("Creating cluster after reaching vpc maximum limit", func() {
		It("Cluster created after reaching vpc limit should be in provisioning", func() {
			By("Create VPCs until limit is reached")
			sess := getSession()
			limit := getElasticIPsLimit(sess)
			var vpcsCreated []string
			for getCurrentVPCsCount(sess) < limit {
				vpcsCreated = append(vpcsCreated, createVPC(sess, "10.0.0.0/16"))
			}

			By("Creating cluster beyond vpc limit")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			workloadClusterTemplate := clusterctl.ConfigCluster(ctx, clusterctl.ConfigClusterInput{
				LogFolder:                filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
				ClusterctlConfigPath:     clusterctlConfigPath,
				KubeconfigPath:           bootstrapClusterProxy.GetKubeconfigPath(),
				InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
				Flavor:                   clusterctl.DefaultFlavor,
				Namespace:                namespace.Name,
				ClusterName:              clusterName,
				KubernetesVersion:        e2eConfig.GetVariable(KubernetesVersion),
				ControlPlaneMachineCount: pointer.Int64Ptr(1),
				WorkerMachineCount:       pointer.Int64Ptr(0),
			})
			Expect(bootstrapClusterProxy.Apply(ctx, workloadClusterTemplate)).ShouldNot(HaveOccurred())

			By("Checking cluster gets provisioned when resources available")
			if len(vpcsCreated) > 0 {
				By("Deleting VPCs")
				deleteVPCs(sess, vpcsCreated)
			}

			By("Waiting for cluster to reach infrastructure ready")
			Eventually(func() bool {
				cluster := &clusterv1.Cluster{}
				if err := bootstrapClusterProxy.GetClient().Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace.Name, Name: clusterName}, cluster); nil == err {
					if cluster.Status.InfrastructureReady {
						return true
					}
				}
				return false
			}, e2eConfig.GetIntervals("", "wait-cluster")...).Should(Equal(true))
		})
	})

	Describe("MachineDeployment misconfigurations", func() {
		It("Should fail to create MachineDeployment with invalid subnet or non-configured Availability Zone", func() {
			By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			_, _ = createCluster(ctx, configCluster)

			By("Creating Machine Deployment with invalid subnet ID")
			md1Name := clusterName + "-md-1"
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 bootstrapClusterProxy.GetClient(),
				MachineDeployment:       makeMachineDeployment(namespace.Name, md1Name, clusterName, 1),
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, md1Name),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, md1Name, e2eConfig.GetVariable(AwsNodeMachineType), nil, pointer.StringPtr("invalid-subnet")),
			})

			By("Looking for failure event to be reported")
			Eventually(func() bool {
				eventList := getEvents(namespace.Name)
				subnetError := "Failed to create instance: failed to run instance: InvalidSubnetID.NotFound: " +
					"The subnet ID '%s' does not exist"
				return isErrorEventExists(namespace.Name, md1Name, "FailedCreate", fmt.Sprintf(subnetError, "invalid-subnet"), eventList)
			}, e2eConfig.GetIntervals("", "wait-worker-nodes")...).Should(BeTrue())

			By("Creating Machine Deployment in non-configured Availability Zone")
			md2Name := clusterName + "-md-2"
			//By default, first availability zone will be used for cluster resources. This step attempts to create a machine deployment in the second availability zone
			invalidAz := getAvailabilityZones()[1].ZoneName
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 bootstrapClusterProxy.GetClient(),
				MachineDeployment:       makeMachineDeployment(namespace.Name, md2Name, clusterName, 1),
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, md2Name),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, md2Name, e2eConfig.GetVariable(AwsNodeMachineType), invalidAz, nil),
			})

			By("Looking for failure event to be reported")
			Eventually(func() bool {
				eventList := getEvents(namespace.Name)
				azError := "Failed to create instance: no subnets available in availability zone \"%s\""
				return isErrorEventExists(namespace.Name, md2Name, "FailedCreate", fmt.Sprintf(azError, *invalidAz), eventList)
			}, e2eConfig.GetIntervals("", "wait-worker-nodes")...).Should(BeTrue())
		})
	})

	Describe("Workload cluster in multiple AZs", func() {
		It("It should be creatable and deletable", func() {
			By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.ControlPlaneMachineCount = pointer.Int64Ptr(3)
			configCluster.Flavor = MultiAzFlavor
			cluster, _ := createCluster(ctx, configCluster)

			By("Adding worker nodes to additional subnets")
			mdName1 := clusterName + "-md-1"
			mdName2 := clusterName + "-md-2"
			md1 := makeMachineDeployment(namespace.Name, mdName1, clusterName, 1)
			md2 := makeMachineDeployment(namespace.Name, mdName2, clusterName, 1)
			az1 := os.Getenv(AwsAvailabilityZone1)
			az2 := os.Getenv(AwsAvailabilityZone2)

			//private CIDRs set in cluster-template-multi-az.yaml
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 bootstrapClusterProxy.GetClient(),
				MachineDeployment:       md1,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName1),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName1, e2eConfig.GetVariable(AwsNodeMachineType), pointer.StringPtr(az1), getSubnetId("cidr-block", "10.0.0.0/24")),
			})
			framework.CreateMachineDeployment(ctx, framework.CreateMachineDeploymentInput{
				Creator:                 bootstrapClusterProxy.GetClient(),
				MachineDeployment:       md2,
				BootstrapConfigTemplate: makeJoinBootstrapConfigTemplate(namespace.Name, mdName2),
				InfraMachineTemplate:    makeAWSMachineTemplate(namespace.Name, mdName2, e2eConfig.GetVariable(AwsNodeMachineType), pointer.StringPtr(az2), getSubnetId("cidr-block", "10.0.2.0/24")),
			})

			By("Waiting for new worker nodes to become ready")
			k8sClient := bootstrapClusterProxy.GetClient()
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md1}, e2eConfig.GetIntervals("", "wait-worker-nodes")...)
			framework.WaitForMachineDeploymentNodesToExist(ctx, framework.WaitForMachineDeploymentNodesToExistInput{Lister: k8sClient, Cluster: cluster, MachineDeployment: md2}, e2eConfig.GetIntervals("", "wait-worker-nodes")...)
		})
	})

	Describe("multiple workload clusters", func() {
		Context("in different namespaces with machine failures", func() {
			It("should setup namespaces correctly for the two clusters", func() {
				By("Creating first cluster with single control plane")
				ns1, cf1 := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
					Creator:   bootstrapClusterProxy.GetClient(),
					ClientSet: bootstrapClusterProxy.GetClientSet(),
					Name:      fmt.Sprintf("multi-workload-%s", util.RandomString(6)),
					LogFolder: filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
				})
				ns2, cf2 := framework.CreateNamespaceAndWatchEvents(ctx, framework.CreateNamespaceAndWatchEventsInput{
					Creator:   bootstrapClusterProxy.GetClient(),
					ClientSet: bootstrapClusterProxy.GetClientSet(),
					Name:      fmt.Sprintf("multi-workload-%s", util.RandomString(6)),
					LogFolder: filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
				})

				By("Creating first cluster")
				cluster1Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster := defaultConfigCluster(cluster1Name, ns1.Name)
				configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
				configCluster.Flavor = LimitAzFlavor
				cluster1, md1 := createCluster(ctx, configCluster)
				Expect(len(md1)).To(Equal(1), "Expecting one MachineDeployment")

				By("Deleting a worker node machine")
				deleteMachine(ns1, md1[0])
				time.Sleep(10 * time.Second)

				By("Verifying MachineDeployment is running.")
				framework.DiscoveryAndWaitForMachineDeployments(ctx, framework.DiscoveryAndWaitForMachineDeploymentsInput{Cluster: cluster1, Lister: bootstrapClusterProxy.GetClient()}, e2eConfig.GetIntervals("", "wait-worker-nodes")...)

				By("Creating second cluster")
				cluster2Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster = defaultConfigCluster(cluster2Name, ns2.Name)
				configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
				configCluster.Flavor = LimitAzFlavor
				cluster2, md2 := createCluster(ctx, configCluster)
				Expect(len(md2)).To(Equal(1), "Expecting one MachineDeployment")

				By("Deleting node directly from infra cloud")
				machines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
					Lister:            bootstrapClusterProxy.GetClient(),
					ClusterName:       cluster1Name,
					Namespace:         ns2.Name,
					MachineDeployment: *md2[0],
				})
				Expect(len(machines)).Should(BeNumerically(">", 0))
				terminateInstance(*machines[0].Spec.ProviderID)

				By("Waiting for machine to reach Failed state")
				statusChecks := []framework.MachineStatusCheck{framework.MachinePhaseCheck(string(clusterv1.MachinePhaseFailed))}
				machineStatusInput := framework.WaitForMachineStatusCheckInput{
					Getter:       bootstrapClusterProxy.GetClient(),
					Machine:      &machines[0],
					StatusChecks: statusChecks,
				}
				framework.WaitForMachineStatusCheck(ctx, machineStatusInput, e2eConfig.GetIntervals("", "wait-machine-status")...)

				By("Deleting the clusters and namespaces")
				deleteCluster(ctx, cluster1)
				deleteCluster(ctx, cluster2)
				framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{Deleter: bootstrapClusterProxy.GetClient(), Name: ns1.Name})
				framework.DeleteNamespace(ctx, framework.DeleteNamespaceInput{Deleter: bootstrapClusterProxy.GetClient(), Name: ns2.Name})
				cf1()
				cf2()
			})
		})

		Context("in same namespace", func() {
			It("should create the clusters", func() {
				By("Creating first cluster with single control plane")
				cluster1Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster := defaultConfigCluster(cluster1Name, namespace.Name)
				configCluster.Flavor = LimitAzFlavor
				cluster1, _ := createCluster(ctx, configCluster)

				By("Creating second cluster with single control plane")
				cluster2Name := fmt.Sprintf("cluster-%s", util.RandomString(6))
				configCluster = defaultConfigCluster(cluster2Name, namespace.Name)
				configCluster.Flavor = LimitAzFlavor
				cluster2, _ := createCluster(ctx, configCluster)

				By("Deleting the Clusters")
				deleteCluster(ctx, cluster1)
				deleteCluster(ctx, cluster2)
			})
		})
	})

	Describe("Workload cluster with spot instances", func() {
		It("It should be creatable and deletable", func() {
			By("Creating a cluster")
			clusterName := fmt.Sprintf("cluster-%s", util.RandomString(6))
			configCluster := defaultConfigCluster(clusterName, namespace.Name)
			configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
			configCluster.Flavor = SpotInstancesFlavor
			_, md := createCluster(ctx, configCluster)

			workerMachines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
				Lister:            bootstrapClusterProxy.GetClient(),
				ClusterName:       clusterName,
				Namespace:         namespace.Name,
				MachineDeployment: *md[0],
			})
			controlPlaneMachines := framework.GetControlPlaneMachinesByCluster(ctx, framework.GetControlPlaneMachinesByClusterInput{
				Lister:      bootstrapClusterProxy.GetClient(),
				ClusterName: clusterName,
				Namespace:   namespace.Name,
			})
			Expect(len(workerMachines)).To(Equal(1))
			assertSpotInstanceType(*workerMachines[0].Spec.ProviderID)
			Expect(len(controlPlaneMachines)).To(Equal(1))
			assertSpotInstanceType(*controlPlaneMachines[0].Spec.ProviderID)
		})
	})

	AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		dumpSpecResourcesAndCleanup(ctx, "", bootstrapClusterProxy, artifactFolder, namespace, e2eConfig.GetIntervals, skipCleanup)
	})
})

func createCluster(ctx context.Context, configCluster clusterctl.ConfigClusterInput) (*clusterv1.Cluster, []*clusterv1.MachineDeployment) {

	res := clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		ClusterProxy:                 bootstrapClusterProxy,
		ConfigCluster:                configCluster,
		WaitForClusterIntervals:      e2eConfig.GetIntervals("", "wait-cluster"),
		WaitForControlPlaneIntervals: e2eConfig.GetIntervals("", "wait-control-plane"),
		WaitForMachineDeployments:    e2eConfig.GetIntervals("", "wait-worker-nodes"),
	})

	return res.Cluster, res.MachineDeployments
}

func defaultConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	return clusterctl.ConfigClusterInput{
		LogFolder:                filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
		ClusterctlConfigPath:     clusterctlConfigPath,
		KubeconfigPath:           bootstrapClusterProxy.GetKubeconfigPath(),
		InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
		Flavor:                   clusterctl.DefaultFlavor,
		Namespace:                namespace,
		ClusterName:              clusterName,
		KubernetesVersion:        e2eConfig.GetVariable(KubernetesVersion),
		ControlPlaneMachineCount: pointer.Int64Ptr(1),
		WorkerMachineCount:       pointer.Int64Ptr(0),
	}
}

func createLBService(svcNamespace string, svcName string, k8sclient crclient.Client) string {
	Byf("Creating service of type Load Balancer with name: %s under namespace: %s", svcName, svcNamespace)
	svcSpec := corev1.ServiceSpec{
		Type: corev1.ServiceTypeLoadBalancer,
		Ports: []corev1.ServicePort{
			{
				Port:     80,
				Protocol: corev1.ProtocolTCP,
			},
		},
		Selector: map[string]string{
			"app": "nginx",
		},
	}
	createService(svcName, svcNamespace, nil, svcSpec, k8sclient)
	// this sleep is required for the service to get updated with ingress details
	time.Sleep(15 * time.Second)
	svcCreated := &corev1.Service{}
	err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: svcNamespace, Name: svcName}, svcCreated)
	Expect(err).NotTo(HaveOccurred())
	elbName := ""
	if lbs := len(svcCreated.Status.LoadBalancer.Ingress); lbs > 0 {
		ingressHostname := svcCreated.Status.LoadBalancer.Ingress[0].Hostname
		elbName = strings.Split(ingressHostname, "-")[0]
	}
	Byf("Created Load Balancer service and ELB name is: %s", elbName)

	return elbName
}

func createPodTemplateSpec(statefulsetinfo statefulSetInfo) corev1.PodTemplateSpec {
	By("Creating PodTemplateSpec config object")
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: statefulsetinfo.selector,
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &statefulsetinfo.podTerminationGracePeriod,
			Containers: []corev1.Container{
				{
					Name:  statefulsetinfo.containerName,
					Image: statefulsetinfo.containerImage,
					Ports: []corev1.ContainerPort{{Name: statefulsetinfo.svcPortName, ContainerPort: statefulsetinfo.containerPort}},
					VolumeMounts: []corev1.VolumeMount{
						{Name: statefulsetinfo.volumeName, MountPath: statefulsetinfo.volMountPath},
					},
				},
			},
		},
	}
	return podTemplateSpec
}

func createPVC(statefulsetinfo statefulSetInfo) corev1.PersistentVolumeClaim {
	By("Creating PersistentVolumeClaim config object")
	volClaimTemplate := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: statefulsetinfo.volumeName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: &statefulsetinfo.storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("1Gi"),
				},
			},
		},
	}
	return volClaimTemplate
}

func createService(svcName string, svcNamespace string, labels map[string]string, serviceSpec corev1.ServiceSpec, k8sClient crclient.Client) {
	svcToCreate := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: svcNamespace,
			Name:      svcName,
		},
		Spec: serviceSpec,
	}
	if len(labels) > 0 {
		svcToCreate.ObjectMeta.Labels = labels
	}
	Expect(k8sClient.Create(context.TODO(), &svcToCreate)).NotTo(HaveOccurred())
}

func createStatefulSet(statefulsetinfo statefulSetInfo, k8sclient crclient.Client) {
	By("Creating statefulset")
	createStorageClass(statefulsetinfo.storageClassName, k8sclient)
	svcSpec := corev1.ServiceSpec{
		ClusterIP: "None",
		Ports: []corev1.ServicePort{
			{
				Port: statefulsetinfo.svcPort,
				Name: statefulsetinfo.svcPortName,
			},
		},
		Selector: statefulsetinfo.selector,
	}
	createService(statefulsetinfo.svcName, statefulsetinfo.namespace, statefulsetinfo.selector, svcSpec, k8sclient)
	podTemplateSpec := createPodTemplateSpec(statefulsetinfo)
	volClaimTemplate := createPVC(statefulsetinfo)
	deployStatefulSet(statefulsetinfo, volClaimTemplate, podTemplateSpec, k8sclient)
	waitForStatefulSetRunning(statefulsetinfo, k8sclient)
}

func createStorageClass(storageClassName string, k8sclient crclient.Client) {
	Byf("Creating StorageClass object with name: %s", storageClassName)
	volExpansion := true
	bindingMode := storagev1.VolumeBindingImmediate
	azs := getAvailabilityZones()
	storageClass := storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "storage.k8s.io/v1",
			Kind:       "StorageClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		Parameters: map[string]string{
			"type": "gp2",
		},
		Provisioner:          "kubernetes.io/aws-ebs",
		AllowVolumeExpansion: &volExpansion,
		MountOptions:         []string{"debug"},
		VolumeBindingMode:    &bindingMode,
		AllowedTopologies: []corev1.TopologySelectorTerm{{
			MatchLabelExpressions: []corev1.TopologySelectorLabelRequirement{{
				Key:    StorageClassFailureZoneLabel,
				Values: []string{*azs[0].ZoneName},
			}},
		}},
	}
	Expect(k8sclient.Create(context.TODO(), &storageClass)).NotTo(HaveOccurred())
}

func createVPC(sess client.ConfigProvider, cidrblock string) string {
	ec2Client := ec2.New(sess)
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidrblock),
	}
	result, err := ec2Client.CreateVpc(input)
	Expect(err).NotTo(HaveOccurred())
	return *result.Vpc.VpcId
}

func deleteCluster(ctx context.Context, cluster *clusterv1.Cluster) {
	framework.DeleteCluster(ctx, framework.DeleteClusterInput{
		Deleter: bootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	})

	framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
		Getter:  bootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	}, e2eConfig.GetIntervals("", "wait-delete-cluster")...)
}

func deleteMachine(namespace *corev1.Namespace, md *clusterv1.MachineDeployment) {
	machineList := &clusterv1.MachineList{}
	selector, err := metav1.LabelSelectorAsMap(&md.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())

	bootstrapClient := bootstrapClusterProxy.GetClient()
	err = bootstrapClient.List(context.TODO(), machineList, crclient.InNamespace(namespace.Name), crclient.MatchingLabels(selector))
	Expect(err).NotTo(HaveOccurred())

	Expect(len(machineList.Items)).ToNot(Equal(0))
	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace.Name,
			Name:      machineList.Items[0].Name,
		},
	}
	Expect(bootstrapClient.Delete(context.TODO(), machine)).To(Succeed())
}

func deleteRetainedVolumes(awsVolIds []*string) {
	By("Deleting dynamically provisioned volumes")
	ec2Client := ec2.New(getSession())
	for _, volumeId := range awsVolIds {
		input := &ec2.DeleteVolumeInput{
			VolumeId: aws.String(*volumeId),
		}
		_, err := ec2Client.DeleteVolume(input)
		Expect(err).NotTo(HaveOccurred())
		Byf("Deleted dynamically provisioned volume with ID: %s", *volumeId)
	}
}

func deleteVPCs(sess client.ConfigProvider, vpcIds []string) {
	ec2Client := ec2.New(sess)
	for _, vpcId := range vpcIds {
		input := &ec2.DeleteVpcInput{
			VpcId: aws.String(vpcId),
		}
		_, err := ec2Client.DeleteVpc(input)
		Expect(err).NotTo(HaveOccurred())
	}
}

func deployStatefulSet(statefulsetinfo statefulSetInfo, volClaimTemp corev1.PersistentVolumeClaim, podTemplate corev1.PodTemplateSpec, k8sclient crclient.Client) {
	Byf("Deploying Statefulset with name: %s under namespace: %s", statefulsetinfo.name, statefulsetinfo.namespace)
	statefulset := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: statefulsetinfo.name, Namespace: statefulsetinfo.namespace},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             &statefulsetinfo.replicas,
			Selector:             &metav1.LabelSelector{MatchLabels: statefulsetinfo.selector},
			Template:             podTemplate,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{volClaimTemp},
		},
	}
	Expect(k8sclient.Create(context.TODO(), &statefulset)).NotTo(HaveOccurred())
}

func getAvailabilityZones() []*ec2.AvailabilityZone {
	ec2Client := ec2.New(getSession())
	azs, err := ec2Client.DescribeAvailabilityZones(nil)
	Expect(err).NotTo(HaveOccurred())
	return azs.AvailabilityZones
}

func getCurrentVPCsCount(sess client.ConfigProvider) int {
	ec2Client := ec2.New(sess)
	input := &ec2.DescribeVpcsInput{}
	result, err := ec2Client.DescribeVpcs(input)
	Expect(err).NotTo(HaveOccurred())
	return len(result.Vpcs)
}

func getElasticIPsLimit(sess client.ConfigProvider) int {
	ec2Client := ec2.New(sess)
	input := &ec2.DescribeAccountAttributesInput{
		AttributeNames: []*string{
			aws.String("vpc-max-elastic-ips"),
		},
	}
	result, err := ec2Client.DescribeAccountAttributes(input)
	Expect(err).NotTo(HaveOccurred())
	res, err := strconv.Atoi(*result.AccountAttributes[0].AttributeValues[0].AttributeValue)
	Expect(err).NotTo(HaveOccurred())
	return res
}

func getEvents(namespace string) *corev1.EventList {
	eventsList := &corev1.EventList{}
	if err := bootstrapClusterProxy.GetClient().List(context.TODO(), eventsList, crclient.InNamespace(namespace), crclient.MatchingLabels{}); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while fetching events of namespace: %s, %s \n", namespace, err.Error())
	}

	return eventsList
}

func getSession() client.ConfigProvider {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	Expect(err).NotTo(HaveOccurred())
	return sess
}

func getSubnetId(filterKey, filterValue string) *string {
	var subnetOutput *ec2.DescribeSubnetsOutput
	var err error

	ec2Client := ec2.New(getSession())
	subnetInput := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(filterKey),
				Values: []*string{
					aws.String(filterValue),
				},
			},
		},
	}

	Eventually(func() int {
		subnetOutput, err = ec2Client.DescribeSubnets(subnetInput)
		Expect(err).NotTo(HaveOccurred())
		return len(subnetOutput.Subnets)
	}, e2eConfig.GetIntervals("", "wait-infra-subnets")...).Should(Equal(1))

	return subnetOutput.Subnets[0].SubnetId
}

func getVolumeIds(info statefulSetInfo, k8sclient crclient.Client) []*string {
	By("Retrieving IDs of dynamically provisioned volumes.")
	statefulset := &appsv1.StatefulSet{}
	err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: info.namespace, Name: info.name}, statefulset)
	Expect(err).NotTo(HaveOccurred())
	podSelector, err := metav1.LabelSelectorAsMap(statefulset.Spec.Selector)
	pvcList := &corev1.PersistentVolumeClaimList{}
	err = k8sclient.List(context.TODO(), pvcList, crclient.InNamespace(info.namespace), crclient.MatchingLabels(podSelector))
	Expect(err).NotTo(HaveOccurred())
	var volIds []*string
	for _, pvc := range pvcList.Items {
		volName := pvc.Spec.VolumeName
		volDescription := &corev1.PersistentVolume{}
		err = k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: info.namespace, Name: volName}, volDescription)
		Expect(err).NotTo(HaveOccurred())
		urlSlice := strings.Split(volDescription.Spec.AWSElasticBlockStore.VolumeID, "/")
		volIds = append(volIds, &urlSlice[len(urlSlice)-1])
	}
	return volIds
}

func isErrorEventExists(namespace, machineDeploymentName, eventReason, errorMsg string, eList *corev1.EventList) bool {
	k8sClient := bootstrapClusterProxy.GetClient()
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := k8sClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while getting machinedeployment %s \n", machineDeploymentName)
		return false
	}

	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while reading lables of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return false
	}

	awsMachineList := &infrav1.AWSMachineList{}
	if err := k8sClient.List(context.TODO(), awsMachineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while getting awsmachines of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return false
	}

	eventMachinesCnt := 0
	for _, awsMachine := range awsMachineList.Items {
		for _, event := range eList.Items {
			if strings.Contains(event.Name, awsMachine.Name) && event.Reason == eventReason && strings.Contains(event.Message, errorMsg) {
				eventMachinesCnt++
				break
			}
		}
	}
	if len(awsMachineList.Items) == eventMachinesCnt {
		return true
	}
	return false
}

func makeAWSMachineTemplate(namespace, name, instanceType string, az, subnetId *string) *infrav1.AWSMachineTemplate {
	awsMachine := &infrav1.AWSMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSMachineTemplateSpec{
			Template: infrav1.AWSMachineTemplateResource{
				Spec: infrav1.AWSMachineSpec{
					InstanceType:       instanceType,
					IAMInstanceProfile: "nodes.cluster-api-provider-aws.sigs.k8s.io",
					SSHKeyName:         pointer.StringPtr(os.Getenv("AWS_SSH_KEY_NAME")),
				},
			},
		},
	}
	if az != nil {
		awsMachine.Spec.Template.Spec.FailureDomain = az
	}

	if subnetId != nil {
		resRef := &infrav1.AWSResourceReference{
			ID: subnetId,
		}
		awsMachine.Spec.Template.Spec.Subnet = resRef
	}

	return awsMachine
}

func makeJoinBootstrapConfigTemplate(namespace, name string) *bootstrapv1.KubeadmConfigTemplate {
	return &bootstrapv1.KubeadmConfigTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigTemplateSpec{
			Template: bootstrapv1.KubeadmConfigTemplateResource{
				Spec: bootstrapv1.KubeadmConfigSpec{
					JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
						NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
							Name:             "{{ ds.meta_data.local_hostname }}",
							KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
						},
					},
				},
			},
		},
	}
}

func makeMachineDeployment(namespace, mdName, clusterName string, replicas int32) *clusterv1.MachineDeployment {
	return &clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mdName,
			Namespace: namespace,
			Labels: map[string]string{
				"cluster.x-k8s.io/cluster-name": clusterName,
				"nodepool":                      mdName,
			},
		},
		Spec: clusterv1.MachineDeploymentSpec{
			Replicas: &replicas,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster.x-k8s.io/cluster-name": clusterName,
					"nodepool":                      mdName,
				},
			},
			ClusterName: clusterName,
			Template: clusterv1.MachineTemplateSpec{
				ObjectMeta: clusterv1.ObjectMeta{
					Labels: map[string]string{
						"cluster.x-k8s.io/cluster-name": clusterName,
						"nodepool":                      mdName,
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: clusterName,
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: &corev1.ObjectReference{
							Kind:       "KubeadmConfigTemplate",
							APIVersion: bootstrapv1.GroupVersion.String(),
							Name:       mdName,
							Namespace:  namespace,
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachineTemplate",
						APIVersion: infrav1.GroupVersion.String(),
						Name:       mdName,
						Namespace:  namespace,
					},
					Version: pointer.StringPtr(e2eConfig.GetVariable(KubernetesVersion)),
				},
			},
		},
	}
}

func assertSpotInstanceType(instanceId string) {
	Byf("Finding EC2 spot instance with ID: %s", instanceId)
	ec2Client := ec2.New(getSession())
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId[strings.LastIndex(instanceId, "/")+1:]),
		},
		Filters: []*ec2.Filter{{Name: aws.String("instance-lifecycle"), Values: aws.StringSlice([]string{"spot"})}},
	}

	result, err := ec2Client.DescribeInstances(input)
	Expect(err).To(BeNil())
	Expect(len(result.Reservations)).To(Equal(1))
	Expect(len(result.Reservations[0].Instances)).To(Equal(1))
}

func terminateInstance(instanceId string) {
	Byf("Terminating EC2 instance with ID: %s", instanceId)
	ec2Client := ec2.New(getSession())
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId[strings.LastIndex(instanceId, "/")+1:]),
		},
	}

	result, err := ec2Client.TerminateInstances(input)
	Expect(err).To(BeNil())
	Expect(len(result.TerminatingInstances)).To(Equal(1))
	termCode := int64(32)
	Expect(*result.TerminatingInstances[0].CurrentState.Code).To(Equal(termCode))
}

func verifyElbExists(elbName string, exists bool) {
	Byf("Verifying ELB with name %s present", elbName)
	elbClient := elb.New(getSession())
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(elbName),
		},
	}
	elbsOutput, err := elbClient.DescribeLoadBalancers(input)
	if exists {
		Expect(err).NotTo(HaveOccurred())
		Expect(len(elbsOutput.LoadBalancerDescriptions)).To(Equal(1))
		Byf("ELB with name %s exists", elbName)
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException))
		Byf("ELB with name %s doesn't exists", elbName)
	}
}

func verifyVolumesExists(awsVolumeIds []*string) {
	By("Ensuring dynamically provisioned volumes exists")
	ec2Client := ec2.New(getSession())
	input := &ec2.DescribeVolumesInput{
		VolumeIds: awsVolumeIds,
	}
	_, err := ec2Client.DescribeVolumes(input)
	Expect(err).NotTo(HaveOccurred())
}

func waitForStatefulSetRunning(info statefulSetInfo, k8sclient crclient.Client) {
	Byf("Ensuring Statefulset(%s) is running", info.name)
	statefulset := &appsv1.StatefulSet{}
	Eventually(
		func() (bool, error) {
			if err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: info.namespace, Name: info.name}, statefulset); err != nil {
				return false, err
			}
			return *statefulset.Spec.Replicas == statefulset.Status.ReadyReplicas, nil
		}, 10*time.Minute, 30*time.Second,
	).Should(BeTrue())
}
