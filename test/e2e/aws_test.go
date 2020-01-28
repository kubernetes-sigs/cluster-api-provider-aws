// +build e2e

/*
Copyright 2018 The Kubernetes Authors.

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

package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	apirand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	kubeadmv1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	capiFlag "sigs.k8s.io/cluster-api/test/helpers/flag"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	cniManifests      = capiFlag.DefineOrLookupStringFlag("cniManifests", "https://docs.projectcalico.org/v3.9/manifests/calico.yaml", "URL to CNI manifests to load")
	kubectlBinary     = capiFlag.DefineOrLookupStringFlag("kubectlBinary", "kubectl", "path to the kubectl binary")
	availabilityZones []*string
	privateCIDRs      = [3]string{"10.0.0.0/24", "10.0.2.0/24", "10.0.4.0/24"}
	publicCIDRs       = [3]string{"10.0.1.0/24", "10.0.3.0/24", "10.0.5.0/24"}
)

const (
	scaleUpReplicas   = 5
	scaleDownReplicas = 0
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

type testSetup struct {
	namespace               string
	clusterName             string
	awsClusterName          string
	cpMachinePrefix         string
	cpAWSMachinePrefix      string
	cpBootstrapConfigPrefix string
	mdBootstrapConfig       string
	machineDeploymentName   string
	awsMachineTemplateName  string
	testTmpDir              string
	instanceType            string
	initialReplicas         int32
	multipleAZ              bool
	availabilityZone        *string
	subnetId                *string
}

var _ = Describe("functional tests", func() {
	var (
		setup         testSetup
		cancelWatches context.CancelFunc
	)
	BeforeEach(func() {
		var ctx context.Context
		ctx, cancelWatches = context.WithCancel(context.Background())

		var err error
		setup = testSetup{}
		setup.testTmpDir, err = ioutil.TempDir(suiteTmpDir, "aws-test")
		Expect(err).NotTo(HaveOccurred())

		setup.namespace = "test-" + util.RandomString(6)
		createNamespace(setup.namespace)

		go func() {
			defer GinkgoRecover()
			watchEvents(ctx, setup.namespace)
		}()

		setup.clusterName = "test-" + util.RandomString(6)
		setup.awsClusterName = "test-infra-" + util.RandomString(6)
		setup.cpMachinePrefix = "test-" + util.RandomString(6)
		setup.cpAWSMachinePrefix = "test-infra-" + util.RandomString(6)
		setup.cpBootstrapConfigPrefix = "test-boot-" + util.RandomString(6)
		setup.mdBootstrapConfig = "test-boot-md" + util.RandomString(6)
		setup.machineDeploymentName = "test-capa-md" + util.RandomString(6)
		setup.awsMachineTemplateName = "test-infra-capa-mt" + util.RandomString(6)
		setup.initialReplicas = 2
		setup.instanceType = "t3.large"
		setup.multipleAZ = false
	})

	AfterEach(func() {
		defer cancelWatches()
	})

	Describe("workload cluster lifecycle", func() {
		It("It should be creatable and deletable", func() {
			By("Creating a cluster with single control plane")
			makeSingleControlPlaneCluster(setup)

			By("Creating the second Control Plane Machine")
			awsMachineName := setup.cpAWSMachinePrefix + "-1"
			bootstrapConfigName := setup.cpBootstrapConfigPrefix + "-1"
			machineName := setup.cpMachinePrefix + "-1"
			createAdditionalControlPlaneMachine(setup, machineName, awsMachineName, bootstrapConfigName)

			By("Creating the third Control Plane Machine")
			awsMachineName = setup.cpAWSMachinePrefix + "-2"
			bootstrapConfigName = setup.cpBootstrapConfigPrefix + "-2"
			machineName = setup.cpMachinePrefix + "-2"
			createAdditionalControlPlaneMachine(setup, machineName, awsMachineName, bootstrapConfigName)

			By("Creating the MachineDeployment")
			createMachineDeployment(setup)

			By("Scale the MachineDeployment up")
			scaleMachineDeployment(setup.namespace, setup.machineDeploymentName, setup.initialReplicas, scaleUpReplicas)

			By("Scale the MachineDeployment down")
			scaleMachineDeployment(setup.namespace, setup.machineDeploymentName, scaleUpReplicas, scaleDownReplicas)

			By("Scale the MachineDeployment up again")
			scaleMachineDeployment(setup.namespace, setup.machineDeploymentName, scaleDownReplicas, setup.initialReplicas)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Provisioning LoadBalancer dynamically and deleting on cluster deletion", func() {
		lbServiceName := "test-svc-" + util.RandomString(6)
		It("It should create and delete Load Balancer", func() {
			By("Creating a cluster with single control plane")
			clusterK8sClient := makeSingleControlPlaneCluster(setup)

			By("Creating the MachineDeployment")
			createMachineDeployment(setup)

			By("Creating the LB service")
			elbName := createLBService(metav1.NamespaceDefault, lbServiceName, clusterK8sClient)
			verifyElbExists(elbName, true)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)

			By("Verifying whether provisioned LB deleted")
			verifyElbExists(elbName, false)
		})
	})

	Describe("Provisioning volumes dynamically and retain on cluster deletion", func() {
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

		It("It should create volumes and volumes should not be deleted along with cluster infra", func() {
			By("Creating a cluster with single control plane")
			clusterK8sClient := makeSingleControlPlaneCluster(setup)

			By("Creating the MachineDeployment")
			createMachineDeployment(setup)

			By("Deploying StatefulSet on infra")
			createStatefulSet(nginxStatefulsetInfo, clusterK8sClient)
			awsVolIds := getVolumeIds(nginxStatefulsetInfo, clusterK8sClient)
			verifyVolumesExists(awsVolIds)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)

			By("Verifying dynamically provisioned volumes retention")
			verifyVolumesExists(awsVolIds)

			By("Deleting retained dynamically provisioned volumes")
			deleteRetainedVolumes(awsVolIds)
		})
	})

	Describe("MachineDeployment with invalid subnet ID and AZ", func() {
		It("It should be creatable and deletable", func() {
			deployment1 := setup.machineDeploymentName + "-1"
			deployment2 := setup.machineDeploymentName + "-2"
			template1 := setup.awsMachineTemplateName + "-1"
			template2 := setup.awsMachineTemplateName + "-2"
			template3 := setup.awsMachineTemplateName + "-3"
			bsConfig1 := setup.mdBootstrapConfig + "-1"
			bsConfig2 := setup.mdBootstrapConfig + "-2"

			By("Creating a cluster with single control plane")
			makeSingleControlPlaneCluster(setup)

			By("Creating Machine Deployment with invalid subnet ID")
			setup.subnetId = aws.String("notcreated")
			setup.machineDeploymentName = deployment1
			setup.awsMachineTemplateName = template1
			setup.mdBootstrapConfig = bsConfig1
			createPendingMachineDeployment(setup)

			By("Creating Machine Deployment in non-configured Availability Zone")
			setup.availabilityZone = getAvailabilityZones()[1].ZoneName
			setup.machineDeploymentName = deployment2
			setup.awsMachineTemplateName = template2
			setup.mdBootstrapConfig = bsConfig2
			createPendingMachineDeployment(setup)

			By("Ensuring MachineDeployments are not in running state")
			Expect(waitForMachineDeploymentRunning(setup.namespace, deployment1)).To(BeFalse())
			Expect(waitForMachineDeploymentRunning(setup.namespace, deployment2)).To(BeFalse())
			eventList := getEvents(setup.namespace)
			subnetError := "Failed to create instance: failed to run instance: InvalidSubnetID.NotFound: " +
				"The subnet ID '%s' does not exist"
			Expect(isErrorEventExists(setup.namespace, deployment1, "FailedCreate", fmt.Sprintf(subnetError, *setup.subnetId), eventList)).To(BeTrue())

			By("Create new AwsMachineTemplate with correct Subnet ID and update its name in MachineDeployment")
			sess = getSession()
			clusterVpcId := getVpcId(fmt.Sprintf("%s-vpc", setup.clusterName))
			subnetId2 := getSubnetId("vpc-id", clusterVpcId)
			makeAWSMachineTemplate(setup.namespace, template3, setup.instanceType, nil, subnetId2)
			setup.machineDeploymentName = deployment1
			setup.awsMachineTemplateName = template3
			updateMachineDeploymentInfra(setup)

			By("Deleting the cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("multiple workload clusters", func() {
		Context("in different namespaces", func() {
			var ns1, clName1, ns2, clName2 string
			It("should create first cluster", func() {
				ns1 = setup.namespace
				clName1 = setup.clusterName
				By("Creating first cluster with single control plane")
				makeSingleControlPlaneCluster(setup)
			})
			It("should create second cluster in a different namespace", func() {
				ns2 = setup.namespace
				clName2 = setup.clusterName
				By("Creating second cluster with single control plane")
				makeSingleControlPlaneCluster(setup)
			})

			It("should delete both clusters", func() {
				By("Deleting the Clusters")
				deleteCluster(ns1, clName1)
				deleteCluster(ns2, clName2)
			})
		})

		Context("in same namespace", func() {
			var ns, clName1, clName2 string
			It("should create first cluster", func() {
				ns = setup.namespace
				clName1 = setup.clusterName
				By("Creating first cluster with single control plane")
				makeSingleControlPlaneCluster(setup)
			})
			It("should create second cluster in the same namespace", func() {
				setup.namespace = ns
				clName2 = setup.clusterName
				By("Creating second cluster with single control plane")
				makeSingleControlPlaneCluster(setup)
			})

			It("should delete both clusters", func() {
				By("Deleting the Clusters")
				deleteCluster(ns, clName1)
				deleteCluster(ns, clName2)
			})
		})
	})

	Describe("MachineDeployment will replace a deleted Machine", func() {
		It("It should reconcile the deleted machine", func() {
			By("Creating a workload cluster with single control plane")
			makeSingleControlPlaneCluster(setup)

			By("Creating the MachineDeployment")
			createMachineDeployment(setup)

			By("Deleting a worker node machine")
			deleteMachine(setup.namespace, setup.machineDeploymentName)
			time.Sleep(10 * time.Second)

			waitForMachineDeploymentRunning(setup.namespace, setup.machineDeploymentName)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Workload cluster in multiple AZs", func() {
		It("It should be creatable and deletable", func() {
			By("Creating a workload cluster with single control plane")
			setup.multipleAZ = true
			makeSingleControlPlaneCluster(setup)

			By("Creating the second Control Plane Machine in second AZ")
			awsMachineName := setup.cpAWSMachinePrefix + "-1"
			bootstrapConfigName := setup.cpBootstrapConfigPrefix + "-1"
			machineName := setup.cpMachinePrefix + "-1"
			setup.availabilityZone = availabilityZones[1]
			createAdditionalControlPlaneMachine(setup, machineName, awsMachineName, bootstrapConfigName)

			By("Creating the third Control Plane Machine using SubnetID in third AZ")
			awsMachineName = setup.cpAWSMachinePrefix + "-2"
			bootstrapConfigName = setup.cpBootstrapConfigPrefix + "-2"
			machineName = setup.cpMachinePrefix + "-2"
			subnetId3 := getSubnetId("cidr-block", privateCIDRs[2])
			setup.availabilityZone = nil
			setup.subnetId = subnetId3
			createAdditionalControlPlaneMachine(setup, machineName, awsMachineName, bootstrapConfigName)

			By("Creating the MachineDeployment in second AZ")
			setup.availabilityZone = availabilityZones[1]
			setup.subnetId = nil
			createMachineDeployment(setup)

			By("Creating the MachineDeployment using SubnetID in third AZ")
			setup.awsMachineTemplateName = setup.awsMachineTemplateName + "-2"
			setup.mdBootstrapConfig = setup.mdBootstrapConfig + "-2"
			setup.machineDeploymentName = setup.machineDeploymentName + "-2"
			setup.availabilityZone = nil
			setup.subnetId = subnetId3
			createMachineDeployment(setup)

			subnetId2 := getSubnetId("cidr-block", privateCIDRs[1])
			By("Verifying if the nodes are deployed in second AZ")
			verfiyInstancesInSubnet(setup.initialReplicas+1, subnetId2)

			By("Verifying if the nodes are deployed in third AZ")
			verfiyInstancesInSubnet(setup.initialReplicas+1, subnetId3)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Creating cluster after reaching vpc maximum limit", func() {
		It("Cluster created after reaching vpc limit should be in provisioning", func() {
			By("Create clusters till vpc limit")
			sess = getSession()
			limit := getElasticIPsLimit()
			var vpcsCreated []string
			for getCurrentVPCsCount() < limit {
				vpcsCreated = append(vpcsCreated, createVPC("10.0.0.0/16"))
			}

			By("Creating cluster beyond vpc limit")
			Expect(createCluster(setup.namespace, setup.clusterName, setup.awsClusterName, setup.multipleAZ)).Should(BeFalse())
			Expect(getClusterStatus(setup.namespace, setup.clusterName)).Should(Equal(string(clusterv1.ClusterPhaseProvisioning)))

			By("Checking cluster gets provisioned when resources available")
			if len(vpcsCreated) > 0 {
				deleteVPCs(vpcsCreated)
				Expect(waitForClusterInfrastructureReady(setup.namespace, setup.clusterName)).Should(BeTrue())
			}

			By("Deleting the cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Delete infra node directly from infra provider", func() {
		It("Machine referencing deleted infra node should come to failed state", func() {
			By("Creating a workload cluster with single control plane")
			makeSingleControlPlaneCluster(setup)

			By("Creating the MachineDeployment")
			setup.availabilityZone = nil
			setup.subnetId = nil
			createMachineDeployment(setup)

			By("Deleting node directly from infra cloud")
			machines, err := getMachinesOfDeployment(setup.namespace, setup.machineDeploymentName)
			Expect(err).To(BeNil())
			Expect(len(machines.Items)).Should(BeNumerically(">", 0))
			terminateInstance(*machines.Items[0].Spec.ProviderID)
			verifyMachinePhase(setup.namespace, machines.Items[0].Name, clusterv1.MachinePhaseFailed)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Create cluster with name having more than 22 characters", func() {
		It("Cluster should be provisioned and deleted", func() {
			By("Creating a workload cluster with single control plane")
			setup.clusterName = "test-" + util.RandomString(20)
			makeSingleControlPlaneCluster(setup)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})

	Describe("Create cluster with name having '.'", func() {
		It("Cluster should be provisioned and deleted", func() {
			By("Creating a workload cluster with single control plane")
			setup.clusterName = "test." + util.RandomString(6)
			makeSingleControlPlaneCluster(setup)

			By("Deleting the Cluster")
			deleteCluster(setup.namespace, setup.clusterName)
		})
	})
})

func watchEvents(ctx context.Context, namespace string) {
	logFile := path.Join(artifactPath, "resources", namespace, "events.log")
	fmt.Fprintf(GinkgoWriter, "Creating directory: %s\n", filepath.Dir(logFile))
	Expect(os.MkdirAll(filepath.Dir(logFile), 0755)).To(Succeed())

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		clientSet,
		10*time.Minute,
		informers.WithNamespace(namespace),
	)
	eventInformer := informerFactory.Core().V1().Events().Informer()
	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := obj.(*corev1.Event)
			f.WriteString(fmt.Sprintf("[New Event] %s/%s\n\tresource: %s/%s/%s\n\treason: %s\n\tmessage: %s\n\tfull: %#v\n",
				e.Namespace, e.Name, e.InvolvedObject.APIVersion, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message, e))
		},
		UpdateFunc: func(_, obj interface{}) {
			e := obj.(*corev1.Event)
			f.WriteString(fmt.Sprintf("[Updated Event] %s/%s\n\tresource: %s/%s/%s\n\treason: %s\n\tmessage: %s\n\tfull: %#v\n",
				e.Namespace, e.Name, e.InvolvedObject.APIVersion, e.InvolvedObject.Kind, e.InvolvedObject.Name, e.Reason, e.Message, e))
		},
		DeleteFunc: func(obj interface{}) {},
	})

	stopInformer := make(chan struct{})
	defer close(stopInformer)
	informerFactory.Start(stopInformer)
	<-ctx.Done()
	stopInformer <- struct{}{}
}

func updateMachineDeploymentInfra(setup testSetup) {
	deployment := &clusterv1.MachineDeployment{}
	Expect(kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: setup.namespace, Name: setup.machineDeploymentName}, deployment))
	deployment.Spec.Template.Spec.InfrastructureRef.Name = setup.awsMachineTemplateName
	Expect(kindClient.Update(context.TODO(), deployment)).NotTo(HaveOccurred())
	waitForMachinesCountMatch(setup.namespace, setup.machineDeploymentName, setup.initialReplicas, setup.initialReplicas)
	waitForMachineDeploymentRunning(setup.namespace, setup.machineDeploymentName)
}

func verifyMachinePhase(namespace, machineName string, phase clusterv1.MachinePhase) {
	fmt.Fprintf(GinkgoWriter, "Ensuring machine %s's state is %s ... \n", machineName, string(phase))
	machine := &clusterv1.Machine{}
	Eventually(
		func() bool {
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineName}, machine); err != nil {
				return false
			}
			return machine.Status.Phase == string(phase)
		}, 20*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func terminateInstance(instanceId string) {
	fmt.Fprintf(GinkgoWriter, "Terminating EC2 instance with ID: %s ... \n", instanceId)
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

func getMachinesOfDeployment(namespace, machineDeploymentName string) (*clusterv1.MachineList, error) {
	fmt.Fprintf(GinkgoWriter, "Fetching the Machines of MachineDeployment %s: \n", machineDeploymentName)
	machineList := &clusterv1.MachineList{}
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while getting machinedeployment %s \n", machineDeploymentName)
		return nil, err
	}

	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while reading lables of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return nil, err
	}

	if err := kindClient.List(context.TODO(), machineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while getting machines of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return nil, err
	}
	return machineList, nil
}

func deleteVPCs(vpcIds []string) {
	ec2Client := ec2.New(sess)
	for _, vpcId := range vpcIds {
		input := &ec2.DeleteVpcInput{
			VpcId: aws.String(vpcId),
		}
		_, err := ec2Client.DeleteVpc(input)
		Expect(err).NotTo(HaveOccurred())
	}
}

func createVPC(cidrblock string) string {
	ec2Client := ec2.New(sess)
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(cidrblock),
	}
	result, err := ec2Client.CreateVpc(input)
	Expect(err).NotTo(HaveOccurred())
	return *result.Vpc.VpcId
}

func getClusterStatus(namespace, clusterName string) string {
	cluster := &clusterv1.Cluster{}
	if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: clusterName}, cluster); nil == err {
		return cluster.Status.Phase
	}
	return ""
}

func getElasticIPsLimit() int {
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

func getCurrentVPCsCount() int {
	ec2Client := ec2.New(sess)
	input := &ec2.DescribeVpcsInput{}
	result, err := ec2Client.DescribeVpcs(input)
	Expect(err).NotTo(HaveOccurred())
	return len(result.Vpcs)
}

func getAvailabilityZones() []*ec2.AvailabilityZone {
	ec2Client := ec2.New(getSession())
	azs, err := ec2Client.DescribeAvailabilityZones(nil)
	Expect(err).NotTo(HaveOccurred())
	return azs.AvailabilityZones
}

func createStatefulSet(statefulsetinfo statefulSetInfo, k8sclient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "Creating statefulset...\n")
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

func deleteRetainedVolumes(awsVolIds []*string) {
	fmt.Fprintf(GinkgoWriter, "Deleting dynamically provisioned volumes...\n")
	ec2Client := ec2.New(getSession())
	for _, volumeId := range awsVolIds {
		input := &ec2.DeleteVolumeInput{
			VolumeId: aws.String(*volumeId),
		}
		_, err := ec2Client.DeleteVolume(input)
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "Deleted dynamically provisioned volume with ID: %s \n", *volumeId)
	}
}

func verifyVolumesExists(awsVolumeIds []*string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring dynamically provisioned volumes exists..\n")
	ec2Client := ec2.New(getSession())
	input := &ec2.DescribeVolumesInput{
		VolumeIds: awsVolumeIds,
	}
	_, err := ec2Client.DescribeVolumes(input)
	Expect(err).NotTo(HaveOccurred())
}

func getVolumeIds(info statefulSetInfo, k8sclient crclient.Client) []*string {
	fmt.Fprintf(GinkgoWriter, "Retrieving IDs of dynamically provisioned volumes..\n")
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

func waitForStatefulSetRunning(info statefulSetInfo, k8sclient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Statefulset(%s) is running..\n", info.name)
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

func deployStatefulSet(statefulsetinfo statefulSetInfo, volClaimTemp corev1.PersistentVolumeClaim, podTemplate corev1.PodTemplateSpec, k8sclient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "Deploying Statefulset with name: %s under namespace: %s\n", statefulsetinfo.name, statefulsetinfo.namespace)
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

func createPVC(statefulsetinfo statefulSetInfo) corev1.PersistentVolumeClaim {
	fmt.Fprintf(GinkgoWriter, "Creating PersistentVolumeClaim config object\n")
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

func createPodTemplateSpec(statefulsetinfo statefulSetInfo) corev1.PodTemplateSpec {
	fmt.Fprintf(GinkgoWriter, "Creating PodTemplateSpec config object\n")
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

func createStorageClass(storageClassName string, k8sclient crclient.Client) {
	fmt.Fprintf(GinkgoWriter, "Creating StorageClass object with name: %s\n", storageClassName)
	volExpansion := true
	bindingMode := storagev1.VolumeBindingImmediate
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
	}
	Expect(k8sclient.Create(context.TODO(), &storageClass)).NotTo(HaveOccurred())
}

func createService(svcName string, svcNamespace string, labels map[string]string, serviceSpec corev1.ServiceSpec, k8sClient crclient.Client) {
	svcToCreate := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: svcNamespace,
			Name:      svcName,
		},
		Spec: serviceSpec,
	}
	if labels != nil && len(labels) > 0 {
		svcToCreate.ObjectMeta.Labels = labels
	}
	Expect(k8sClient.Create(context.TODO(), &svcToCreate)).NotTo(HaveOccurred())
}

func createLBService(svcNamespace string, svcName string, k8sclient crclient.Client) string {
	fmt.Fprintf(GinkgoWriter, "Creating service of type Load Balancer with name: %s under namespace: %s\n", svcName, svcNamespace)
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
	//this sleep is required for the service to get updated with ingress details
	time.Sleep(15 * time.Second)
	svcCreated := &corev1.Service{}
	err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: svcNamespace, Name: svcName}, svcCreated)
	Expect(err).NotTo(HaveOccurred())
	elbName := ""
	if lbs := len(svcCreated.Status.LoadBalancer.Ingress); lbs > 0 {
		ingressHostname := svcCreated.Status.LoadBalancer.Ingress[0].Hostname
		elbName = strings.Split(ingressHostname, "-")[0]
	}
	fmt.Fprintf(GinkgoWriter, "Created Load Balancer service and ELB name is: %s\n", elbName)
	return elbName
}

func verifyElbExists(elbName string, exists bool) {
	fmt.Fprintf(GinkgoWriter, "Verifying ELB with name %s present\n", elbName)
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
		fmt.Fprintf(GinkgoWriter, "ELB with name %s exists\n", elbName)
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException))
		fmt.Fprintf(GinkgoWriter, "ELB with name %s doesn't exists\n", elbName)
	}
}

func createClusterKubeConfigs(tmpDir, namespace, clusterName string) (string, crclient.Client) {
	cluster := crclient.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}
	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, cluster)
	Expect(err).NotTo(HaveOccurred())

	kubeConfigPath := path.Join(tmpDir, clusterName+".kubeconfig")
	Expect(ioutil.WriteFile(kubeConfigPath, kubeConfigData, 0640)).To(Succeed())

	kubeConfigData, readErr := ioutil.ReadFile(kubeConfigPath)
	Expect(readErr).NotTo(HaveOccurred())

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	Expect(err).NotTo(HaveOccurred())

	k8sclient, err := crclient.New(restConfig, crclient.Options{})
	Expect(err).NotTo(HaveOccurred())

	return kubeConfigPath, k8sclient
}

func getSubnetId(filterKey, filterValue string) *string {
	ec2c := ec2.New(sess)
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
	result, err := ec2c.DescribeSubnets(subnetInput)
	Expect(err).NotTo(HaveOccurred())
	return result.Subnets[0].SubnetId
}

func getVpcId(vpcName string) string {
	ec2c := ec2.New(sess)
	input := &ec2.DescribeVpcsInput{}
	result, err := ec2c.DescribeVpcs(input)
	Expect(err).NotTo(HaveOccurred())
	for _, tmpVpc := range result.Vpcs {
		for _, tmpTag := range tmpVpc.Tags {
			if *tmpTag.Key == "Name" && *tmpTag.Value == vpcName {
				return *tmpVpc.VpcId
			}
		}
	}
	return ""
}

func verfiyInstancesInSubnet(numOfInstances int32, subnetId *string) {
	ec2c := ec2.New(sess)
	instanceInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("subnet-id"),
				Values: []*string{subnetId},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	result, err := ec2c.DescribeInstances(instanceInput)
	Expect(err).NotTo(HaveOccurred())
	Expect(int(numOfInstances)).To(Equal(len(result.Reservations)))
}

func getAvailabilityZone() []*ec2.AvailabilityZone {
	ec2c := ec2.New(sess)
	azs, err := ec2c.DescribeAvailabilityZones(nil)
	Expect(err).NotTo(HaveOccurred())
	return azs.AvailabilityZones
}

func createCluster(namespace, clusterName, awsClusterName string, multiAZ bool) bool {
	By("Creating an AWSCluster")
	makeAWSCluster(namespace, awsClusterName, multiAZ)

	By("Creating a Cluster")
	makeCluster(namespace, clusterName, awsClusterName)

	By("Ensuring Cluster Infrastructure Reports as Ready")
	return waitForClusterInfrastructureReady(namespace, clusterName)
}

func makeSingleControlPlaneCluster(setup testSetup) crclient.Client {

	Expect(createCluster(setup.namespace, setup.clusterName, setup.awsClusterName, setup.multipleAZ)).Should(BeTrue())

	By("Creating the initial Control Plane Machine")
	awsMachineName := setup.cpAWSMachinePrefix + "-0"
	bootstrapConfigName := setup.cpBootstrapConfigPrefix + "-0"
	machineName := setup.cpMachinePrefix + "-0"
	createInitialControlPlaneMachine(setup, machineName, awsMachineName, bootstrapConfigName)
	clusterKubeConfigPath, clusterK8sClient := createClusterKubeConfigs(setup.testTmpDir, setup.namespace, setup.clusterName)

	By("Deploying CNI to created Cluster")
	deployCNI(clusterKubeConfigPath, *cniManifests)
	waitForMachineNodeReady(setup.namespace, machineName)
	return clusterK8sClient
}

func createPendingMachineDeployment(setup testSetup) {
	fmt.Fprintf(GinkgoWriter, "Creating MachineDeployment in namespace %s with name %s\n", setup.namespace, setup.machineDeploymentName)
	makeAWSMachineTemplate(setup.namespace, setup.awsMachineTemplateName, setup.instanceType, setup.availabilityZone, setup.subnetId)
	makeJoinBootstrapConfigTemplate(setup.namespace, setup.mdBootstrapConfig)
	makeMachineDeployment(setup.namespace, setup.machineDeploymentName, setup.awsMachineTemplateName, setup.mdBootstrapConfig, setup.clusterName, setup.initialReplicas)
}

func scaleMachineDeployment(namespace, machineDeployment string, replicasCurrent int32, replicasProposed int32) {
	fmt.Fprintf(GinkgoWriter, "Scaling MachineDeployment from %d to %d\n", replicasCurrent, replicasProposed)
	deployment := &clusterv1.MachineDeployment{}
	Expect(kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeployment}, deployment))
	deployment.Spec.Replicas = &replicasProposed
	Expect(kindClient.Update(context.TODO(), deployment)).NotTo(HaveOccurred())
	waitForMachinesCountMatch(namespace, machineDeployment, replicasCurrent, replicasProposed)
	waitForMachineDeploymentRunning(namespace, machineDeployment)
}

func createMachineDeployment(setup testSetup) {
	fmt.Fprintf(GinkgoWriter, "Creating MachineDeployment in namespace %s with name %s\n", setup.namespace, setup.machineDeploymentName)
	makeAWSMachineTemplate(setup.namespace, setup.awsMachineTemplateName, setup.instanceType, setup.availabilityZone, setup.subnetId)
	makeJoinBootstrapConfigTemplate(setup.namespace, setup.mdBootstrapConfig)
	makeMachineDeployment(setup.namespace, setup.machineDeploymentName, setup.awsMachineTemplateName, setup.mdBootstrapConfig, setup.clusterName, setup.initialReplicas)
	waitForMachinesCountMatch(setup.namespace, setup.machineDeploymentName, setup.initialReplicas, setup.initialReplicas)
	waitForMachineDeploymentRunning(setup.namespace, setup.machineDeploymentName)
}

func isErrorEventExists(namespace, machineDeploymentName, eventReason, errorMsg string, eList *corev1.EventList) bool {
	fmt.Fprintf(GinkgoWriter, "Checking Error event with message %s is present \n", errorMsg)
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while getting machinedeployment %s \n", machineDeploymentName)
		return false
	}

	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while reading lables of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return false
	}

	awsMachineList := &infrav1.AWSMachineList{}
	if err := kindClient.List(context.TODO(), awsMachineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
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

func getEvents(namespace string) *corev1.EventList {
	eventsList := &corev1.EventList{}
	if err := kindClient.List(context.TODO(), eventsList, crclient.InNamespace(namespace), crclient.MatchingLabels{}); err != nil {
		fmt.Fprintf(GinkgoWriter, "Got error while fetching events of namespace: %s, %s \n", namespace, err.Error())
	}

	return eventsList
}

func deleteMachine(namespace, machineDeploymentName string) {
	machineDeployment := &clusterv1.MachineDeployment{}
	Expect(kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment)).To(Succeed())
	machineList := &clusterv1.MachineList{}

	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())

	err = kindClient.List(context.TODO(), machineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector))
	Expect(err).NotTo(HaveOccurred())

	Expect(len(machineList.Items)).ToNot(Equal(0))
	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      machineList.Items[0].Name,
		},
	}
	Expect(kindClient.Delete(context.TODO(), machine)).To(Succeed())
}

func waitForMachinesCountMatch(namespace, machineDeploymentName string, replicasCurrent int32, replicasProposed int32) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Machine count matched to %d \n", replicasProposed)
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		return
	}
	scalingDown := replicasCurrent > replicasProposed
	Eventually(
		func() (int32, error) {
			machineList := &clusterv1.MachineList{}
			selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
			if err != nil {
				return -1, err
			}
			if err := kindClient.List(context.TODO(), machineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
				return -1, err
			} else {
				var runningMachines int32
				var deletingMachines int32
				for _, item := range machineList.Items {
					if string(clusterv1.MachinePhaseRunning) == item.Status.Phase {
						runningMachines += 1
					}
					if string(clusterv1.MachinePhaseDeleting) == item.Status.Phase {
						deletingMachines += 1
					}
				}
				if scalingDown {
					if deletingMachines > 0 {
						return -1, nil
					}
				}
				return runningMachines, nil
			}
		},
		20*time.Minute, 30*time.Second,
	).Should(Equal(replicasProposed))
}

func waitForMachineDeploymentRunning(namespace, machineDeploymentName string) bool {
	fmt.Fprintf(GinkgoWriter, "Verifying MachineDeployment %s/%s is running..\n", namespace, machineDeploymentName)
	endTime := time.Now().Add(5 * time.Minute)
	for time.Now().Before(endTime) {
		machineDeployment := &clusterv1.MachineDeployment{}
		if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); nil == err {
			if *machineDeployment.Spec.Replicas == machineDeployment.Status.ReadyReplicas {
				return true
			}
		}
		time.Sleep(15 * time.Second)
	}
	return false
}

func makeMachineDeployment(namespace, mdName, awsMachineTemplateName, bootstrapConfigName, clusterName string, replicas int32) {
	fmt.Fprintf(GinkgoWriter, "Creating MachineDeployment %s/%s\n", namespace, mdName)
	machineDeployment := &clusterv1.MachineDeployment{
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
							Name:       bootstrapConfigName,
							Namespace:  namespace,
						},
					},
					InfrastructureRef: corev1.ObjectReference{
						Kind:       "AWSMachineTemplate",
						APIVersion: infrav1.GroupVersion.String(),
						Name:       awsMachineTemplateName,
						Namespace:  namespace,
					},
					Version: k8sVersion,
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), machineDeployment)).NotTo(HaveOccurred())
}

func makeJoinBootstrapConfigTemplate(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Join KubeadmConfigTemplate %s/%s\n", namespace, name)
	configTemplate := &bootstrapv1.KubeadmConfigTemplate{
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
					Files: []bootstrapv1.File{
						{
							Path:    "/tmp/userdata-length-test.txt",
							Content: apirand.String(20000),
						},
					},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), configTemplate)).NotTo(HaveOccurred())
}

func makeAWSMachineTemplate(namespace, name, instanceType string, az, subnetId *string) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSMachineTemplate %s/%s\n", namespace, name)
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
					SSHKeyName:         keyPairName,
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
	Expect(kindClient.Create(context.TODO(), awsMachine)).NotTo(HaveOccurred())
}

func deployCNI(kubeConfigPath, manifestPath string) {
	Expect(exec.Command(
		*kubectlBinary,
		"create",
		"--kubeconfig="+kubeConfigPath,
		"-f", manifestPath,
	).Run()).To(Succeed())
}

func waitForMachineNodeReady(namespace, name string) {
	machine := &clusterv1.Machine{}
	Expect(kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine)).To(Succeed())

	nodeName := machine.Status.NodeRef.Name

	cluster := crclient.ObjectKey{
		Namespace: namespace,
		Name:      machine.Labels[clusterv1.ClusterLabelName],
	}

	kubeConfigData, err := kubeconfig.FromSecret(context.TODO(), kindClient, cluster)
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigData)
	Expect(err).NotTo(HaveOccurred())

	nodeClient, err := crclient.New(restConfig, crclient.Options{})
	Expect(err).NotTo(HaveOccurred())

	Eventually(
		func() bool {
			node := &corev1.Node{}
			if err := nodeClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Name: nodeName}, node); err != nil {
				fmt.Fprintf(GinkgoWriter, "Error retrieving node: %v", err)
				return false
			}

			conditionMap := make(map[corev1.NodeConditionType]*corev1.NodeCondition)
			for i, condition := range node.Status.Conditions {
				if condition.Type == corev1.NodeReady {
					conditionMap[condition.Type] = &node.Status.Conditions[i]
				}
			}

			if condition, ok := conditionMap[corev1.NodeReady]; ok {
				if condition.Status == corev1.ConditionTrue {
					return true
				}
			}
			return false
		},
		2*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func createAdditionalControlPlaneMachine(setup testSetup, machineName, awsMachineName, bootstrapConfigName string) {
	makeAWSMachine(setup.namespace, awsMachineName, setup.instanceType, setup.availabilityZone, setup.subnetId)
	makeJoinBootstrapConfig(setup.namespace, bootstrapConfigName)
	makeMachine(setup.namespace, machineName, awsMachineName, bootstrapConfigName, setup.clusterName)
	waitForMachineBootstrapReady(setup.namespace, machineName)
	waitForAWSMachineRunning(setup.namespace, awsMachineName)
	waitForAWSMachineReady(setup.namespace, awsMachineName)
	waitForMachineNodeRef(setup.namespace, machineName)
	waitForMachineNodeReady(setup.namespace, machineName)
}

func createInitialControlPlaneMachine(setup testSetup, machineName, awsMachineName, bootstrapConfigName string) {
	makeAWSMachine(setup.namespace, awsMachineName, setup.instanceType, setup.availabilityZone, setup.subnetId)
	makeInitBootstrapConfig(setup.namespace, bootstrapConfigName)
	makeMachine(setup.namespace, machineName, awsMachineName, bootstrapConfigName, setup.clusterName)
	waitForMachineBootstrapReady(setup.namespace, machineName)
	waitForAWSMachineRunning(setup.namespace, awsMachineName)
	waitForAWSMachineReady(setup.namespace, awsMachineName)
	waitForMachineNodeRef(setup.namespace, machineName)
	waitForClusterControlPlaneInitialized(setup.namespace, setup.clusterName)
}

func deleteCluster(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Deleting Cluster %s/%s\n", namespace, name)
	Expect(kindClient.Delete(context.TODO(), &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	})).To(Succeed())

	Eventually(
		func() bool {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
				if apierrors.IsNotFound(err) {
					return true
				}
				return false
			}
			return false
		},
		20*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForMachineNodeRef(namespace, name string) {
	Eventually(
		func() *corev1.ObjectReference {
			machine := &clusterv1.Machine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return nil
			}
			return machine.Status.NodeRef

		},
		10*time.Minute, 15*time.Second,
	).ShouldNot(BeNil())
}

func waitForClusterControlPlaneInitialized(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring control plane initialized for cluster %s/%s\n", namespace, name)
	Eventually(
		func() (bool, error) {
			cluster := &clusterv1.Cluster{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); err != nil {
				return false, err
			}
			return cluster.Status.ControlPlaneInitialized, nil
		},
		15*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForAWSMachineRunning(namespace, name string) {
	Eventually(
		func() (bool, error) {
			awsMachine := &infrav1.AWSMachine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, awsMachine); err != nil {
				return false, err
			}
			if awsMachine.Status.InstanceState == nil {
				return false, nil
			}
			return *awsMachine.Status.InstanceState == infrav1.InstanceStateRunning, nil
		},
		5*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForClusterInfrastructureReady(namespace, name string) bool {
	fmt.Fprintf(GinkgoWriter, "Ensuring infrastructure is ready for cluster %s/%s\n", namespace, name)
	endTime := time.Now().Add(20 * time.Minute)
	for time.Now().Before(endTime) {
		cluster := &clusterv1.Cluster{}
		if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, cluster); nil == err {
			if cluster.Status.InfrastructureReady {
				return true
			}
		}
		time.Sleep(15 * time.Second)
	}
	return false
}

func waitForMachineBootstrapReady(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Ensuring Machine %s/%s has bootstrapReady\n", namespace, name)
	Eventually(
		func() (bool, error) {
			machine := &clusterv1.Machine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, machine); err != nil {
				return false, err
			}
			return machine.Status.BootstrapReady, nil
		},
		5*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func waitForAWSMachineReady(namespace, name string) {
	Eventually(
		func() (bool, error) {
			awsMachine := &infrav1.AWSMachine{}
			if err := kindClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: name}, awsMachine); err != nil {
				return false, err
			}
			return awsMachine.Status.Ready, nil
		},
		2*time.Minute, 15*time.Second,
	).Should(BeTrue())
}

func makeMachine(namespace, name, awsMachineName, bootstrapConfigName, clusterName string) {
	fmt.Fprintf(GinkgoWriter, "Creating Machine %s/%s\n", namespace, name)
	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				clusterv1.MachineControlPlaneLabelName: "true",
			},
		},
		Spec: clusterv1.MachineSpec{
			ClusterName: clusterName,
			Bootstrap: clusterv1.Bootstrap{
				ConfigRef: &corev1.ObjectReference{
					Kind:       "KubeadmConfig",
					APIVersion: bootstrapv1.GroupVersion.String(),
					Name:       bootstrapConfigName,
					Namespace:  namespace,
				},
			},
			InfrastructureRef: corev1.ObjectReference{
				Kind:       "AWSMachine",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       awsMachineName,
				Namespace:  namespace,
			},
			Version: k8sVersion,
		},
	}
	Expect(kindClient.Create(context.TODO(), machine)).To(Succeed())
}

func makeJoinBootstrapConfig(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Join KubeadmConfig %s/%s\n", namespace, name)
	config := &bootstrapv1.KubeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigSpec{
			JoinConfiguration: &kubeadmv1beta1.JoinConfiguration{
				ControlPlane: &kubeadmv1beta1.JoinControlPlane{},
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:             "{{ ds.meta_data.local_hostname }}",
					KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), config)).To(Succeed())
}

func makeInitBootstrapConfig(namespace, name string) {
	fmt.Fprintf(GinkgoWriter, "Creating Init KubeadmConfig %s/%s\n", namespace, name)
	config := &bootstrapv1.KubeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: bootstrapv1.KubeadmConfigSpec{
			ClusterConfiguration: &kubeadmv1beta1.ClusterConfiguration{
				APIServer: kubeadmv1beta1.APIServer{
					ControlPlaneComponent: kubeadmv1beta1.ControlPlaneComponent{
						ExtraArgs: map[string]string{"cloud-provider": "aws"},
					},
				},
				ControllerManager: kubeadmv1beta1.ControlPlaneComponent{
					ExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
			InitConfiguration: &kubeadmv1beta1.InitConfiguration{
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:             "{{ ds.meta_data.local_hostname }}",
					KubeletExtraArgs: map[string]string{"cloud-provider": "aws"},
				},
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), config)).To(Succeed())
}

func makeAWSMachine(namespace, name, instanceType string, az, subnetId *string) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSMachine %s/%s\n", namespace, name)
	awsMachine := &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSMachineSpec{
			InstanceType:       instanceType,
			IAMInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io",
			SSHKeyName:         keyPairName,
		},
	}
	if az != nil {
		awsMachine.Spec.FailureDomain = az
	}
	if subnetId != nil {
		awsMachine.Spec.Subnet = &infrav1.AWSResourceReference{ID: subnetId}
	}
	Expect(kindClient.Create(context.TODO(), awsMachine)).To(Succeed())
}

func makeCluster(namespace, name, awsClusterName string) {
	fmt.Fprintf(GinkgoWriter, "Creating Cluster %s/%s\n", namespace, name)
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: &clusterv1.ClusterNetwork{
				Pods: &clusterv1.NetworkRanges{
					CIDRBlocks: []string{"192.168.0.0/16"},
				},
			},
			InfrastructureRef: &corev1.ObjectReference{
				Kind:       "AWSCluster",
				APIVersion: infrav1.GroupVersion.String(),
				Name:       awsClusterName,
				Namespace:  namespace,
			},
		},
	}
	Expect(kindClient.Create(context.TODO(), cluster)).To(Succeed())
}

func makeAWSCluster(namespace, name string, multipleAZ bool) {
	fmt.Fprintf(GinkgoWriter, "Creating AWSCluster %s/%s\n", namespace, name)
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSClusterSpec{
			Region:     region,
			SSHKeyName: keyPairName,
		},
	}
	if multipleAZ {
		azs := getAvailabilityZone()
		availabilityZones = append(availabilityZones, azs[0].ZoneName, azs[1].ZoneName, azs[2].ZoneName)
		subnets := []*infrav1.SubnetSpec{
			&infrav1.SubnetSpec{
				CidrBlock:        privateCIDRs[0],
				AvailabilityZone: *availabilityZones[0],
			},
			&infrav1.SubnetSpec{
				CidrBlock:        publicCIDRs[0],
				AvailabilityZone: *availabilityZones[0],
				IsPublic:         true,
			},
			&infrav1.SubnetSpec{
				CidrBlock:        privateCIDRs[1],
				AvailabilityZone: *availabilityZones[1],
			},
			&infrav1.SubnetSpec{
				CidrBlock:        publicCIDRs[1],
				AvailabilityZone: *availabilityZones[1],
				IsPublic:         true,
			},
			&infrav1.SubnetSpec{
				CidrBlock:        privateCIDRs[2],
				AvailabilityZone: *availabilityZones[2],
			},
			&infrav1.SubnetSpec{
				CidrBlock:        publicCIDRs[2],
				AvailabilityZone: *availabilityZones[2],
				IsPublic:         true,
			},
		}
		awsCluster.Spec.NetworkSpec.Subnets = subnets
	}

	Expect(kindClient.Create(context.TODO(), awsCluster)).To(Succeed())
}

func createNamespace(namespace string) {
	fmt.Fprintf(GinkgoWriter, "creating namespace %q\n", namespace)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	Expect(kindClient.Create(context.TODO(), ns)).To(Succeed())
}
