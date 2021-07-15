// +build e2e

/*
Copyright 2021 The Kubernetes Authors.

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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/blang/semver"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	client_runtime "sigs.k8s.io/controller-runtime/pkg/client"
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

// GetClusterByName returns a Cluster object given his name
func GetAWSClusterByName(ctx context.Context, namespace, name string) (*infrav1.AWSCluster, error) {
	awsCluster := &infrav1.AWSCluster{}
	key := client_runtime.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	err := e2eCtx.Environment.BootstrapClusterProxy.GetClient().Get(ctx, key, awsCluster)
	return awsCluster, err
}

func createCluster(ctx context.Context, configCluster clusterctl.ConfigClusterInput, result *clusterctl.ApplyClusterTemplateAndWaitResult) (*clusterv1.Cluster, []*clusterv1.MachineDeployment, *controlplanev1.KubeadmControlPlane) {
	clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		ClusterProxy:                 e2eCtx.Environment.BootstrapClusterProxy,
		ConfigCluster:                configCluster,
		WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals("", "wait-cluster"),
		WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals("", "wait-control-plane"),
		WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes"),
	}, result)

	return result.Cluster, result.MachineDeployments, result.ControlPlane
}

func defaultConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	return clusterctl.ConfigClusterInput{
		LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
		ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
		KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
		InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
		Flavor:                   clusterctl.DefaultFlavor,
		Namespace:                namespace,
		ClusterName:              clusterName,
		KubernetesVersion:        e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion),
		ControlPlaneMachineCount: pointer.Int64Ptr(1),
		WorkerMachineCount:       pointer.Int64Ptr(0),
	}
}

func createLBService(svcNamespace string, svcName string, k8sclient crclient.Client) string {
	shared.Byf("Creating service of type Load Balancer with name: %s under namespace: %s", svcName, svcNamespace)
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
	shared.Byf("Created Load Balancer service and ELB name is: %s", elbName)

	return elbName
}

func createPodTemplateSpec(statefulsetinfo statefulSetInfo) corev1.PodTemplateSpec {
	ginkgo.By("Creating PodTemplateSpec config object")
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
	ginkgo.By("Creating PersistentVolumeClaim config object")
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
	ginkgo.By("Creating statefulset")
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
	shared.Byf("Creating StorageClass object with name: %s", storageClassName)
	volExpansion := true
	bindingMode := storagev1.VolumeBindingImmediate
	azs := shared.GetAvailabilityZones(e2eCtx.AWSSession)
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
				Key:    shared.StorageClassFailureZoneLabel,
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
		Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	})

	framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
		Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)
}

func deleteMachine(namespace *corev1.Namespace, md *clusterv1.MachineDeployment) {
	machineList := &clusterv1.MachineList{}
	selector, err := metav1.LabelSelectorAsMap(&md.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())

	bootstrapClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
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
	ginkgo.By("Deleting dynamically provisioned volumes")
	ec2Client := ec2.New(e2eCtx.AWSSession)
	for _, volumeId := range awsVolIds {
		input := &ec2.DeleteVolumeInput{
			VolumeId: aws.String(*volumeId),
		}
		_, err := ec2Client.DeleteVolume(input)
		Expect(err).NotTo(HaveOccurred())
		shared.Byf("Deleted dynamically provisioned volume with ID: %s", *volumeId)
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
	shared.Byf("Deploying Statefulset with name: %s under namespace: %s", statefulsetinfo.name, statefulsetinfo.namespace)
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
	if err := e2eCtx.Environment.BootstrapClusterProxy.GetClient().List(context.TODO(), eventsList, crclient.InNamespace(namespace), crclient.MatchingLabels{}); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Got error while fetching events of namespace: %s, %s \n", namespace, err.Error())
	}

	return eventsList
}

func getSubnetId(filterKey, filterValue string) *string {
	var subnetOutput *ec2.DescribeSubnetsOutput
	var err error

	ec2Client := ec2.New(e2eCtx.AWSSession)
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
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-infra-subnets")...).Should(Equal(1))

	return subnetOutput.Subnets[0].SubnetId
}

func getVolumeIds(info statefulSetInfo, k8sclient crclient.Client) []*string {
	ginkgo.By("Retrieving IDs of dynamically provisioned volumes.")
	statefulset := &appsv1.StatefulSet{}
	err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: info.namespace, Name: info.name}, statefulset)
	Expect(err).NotTo(HaveOccurred())
	podSelector, err := metav1.LabelSelectorAsMap(statefulset.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())
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
	k8sClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	machineDeployment := &clusterv1.MachineDeployment{}
	if err := k8sClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: namespace, Name: machineDeploymentName}, machineDeployment); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Got error while getting machinedeployment %s \n", machineDeploymentName)
		return false
	}

	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	if err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Got error while reading lables of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
		return false
	}

	awsMachineList := &infrav1.AWSMachineList{}
	if err := k8sClient.List(context.TODO(), awsMachineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector)); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Got error while getting awsmachines of machinedeployment: %s, %s \n", machineDeploymentName, err.Error())
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

func getAWSMachinesForDeployment(namespace string, machineDeployment clusterv1.MachineDeployment) *infrav1.AWSMachineList {
	k8sClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())
	awsMachineList := &infrav1.AWSMachineList{}
	Expect(k8sClient.List(context.TODO(), awsMachineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector))).NotTo(HaveOccurred())
	return awsMachineList
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
					JoinConfiguration: &bootstrapv1.JoinConfiguration{
						NodeRegistration: bootstrapv1.NodeRegistrationOptions{
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
					Version: pointer.StringPtr(e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion)),
				},
			},
		},
	}
}

func assertSpotInstanceType(instanceId string) {
	shared.Byf("Finding EC2 spot instance with ID: %s", instanceId)
	ec2Client := ec2.New(e2eCtx.AWSSession)
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
	shared.Byf("Terminating EC2 instance with ID: %s", instanceId)
	ec2Client := ec2.New(e2eCtx.AWSSession)
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
	shared.Byf("Verifying ELB with name %s present", elbName)
	elbClient := elb.New(e2eCtx.AWSSession)
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(elbName),
		},
	}
	elbsOutput, err := elbClient.DescribeLoadBalancers(input)
	if exists {
		Expect(err).NotTo(HaveOccurred())
		Expect(len(elbsOutput.LoadBalancerDescriptions)).To(Equal(1))
		shared.Byf("ELB with name %s exists", elbName)
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException))
		shared.Byf("ELB with name %s doesn't exists", elbName)
	}
}

func verifyVolumesExists(awsVolumeIds []*string) {
	ginkgo.By("Ensuring dynamically provisioned volumes exists")
	ec2Client := ec2.New(e2eCtx.AWSSession)
	input := &ec2.DescribeVolumesInput{
		VolumeIds: awsVolumeIds,
	}
	_, err := ec2Client.DescribeVolumes(input)
	Expect(err).NotTo(HaveOccurred())
}

func waitForStatefulSetRunning(info statefulSetInfo, k8sclient crclient.Client) {
	shared.Byf("Ensuring Statefulset(%s) is running", info.name)
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

// LatestCIReleaseForVersion returns the latest ci release of a specific version.
func LatestCIReleaseForVersion(searchVersion string) (string, error) {
	ciVersionURL := "https://dl.k8s.io/ci/latest-%d.%d.txt"
	tagPrefix := "v"
	searchSemVer, err := semver.Make(strings.TrimPrefix(searchVersion, tagPrefix))
	if err != nil {
		return "", err
	}
	searchSemVer.Minor++
	resp, err := http.Get(fmt.Sprintf(ciVersionURL, searchSemVer.Major, searchSemVer.Minor))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

func getStatefulSetInfo() statefulSetInfo {
	return statefulSetInfo{
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
}
