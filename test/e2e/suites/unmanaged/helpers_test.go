//go:build e2e
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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/blang/semver"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/utils/pointer"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util/conditions"
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
	isInTreeCSI               bool
}

// GetClusterByName returns a Cluster object given his name.
func GetAWSClusterByName(ctx context.Context, namespace, name string) (*infrav1.AWSCluster, error) {
	cluster := &clusterv1.Cluster{}
	key := crclient.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	if err := e2eCtx.Environment.BootstrapClusterProxy.GetClient().Get(ctx, key, cluster); err != nil {
		return nil, err
	}

	awsCluster := &infrav1.AWSCluster{}
	awsClusterKey := crclient.ObjectKey{
		Namespace: namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	err := e2eCtx.Environment.BootstrapClusterProxy.GetClient().Get(ctx, awsClusterKey, awsCluster)
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
		ControlPlaneMachineCount: pointer.Int64(1),
		WorkerMachineCount:       pointer.Int64(0),
	}
}

func createLBService(svcNamespace string, svcName string, k8sclient crclient.Client) string {
	ginkgo.By(fmt.Sprintf("Creating service of type Load Balancer with name: %s under namespace: %s", svcName, svcNamespace))
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
	ginkgo.By(fmt.Sprintf("Created Load Balancer service and ELB name is: %s", elbName))

	return elbName
}

func deleteLBService(svcNamespace string, svcName string, k8sclient crclient.Client) {
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
	deleteService(svcName, svcNamespace, nil, svcSpec, k8sclient)
}

func createPodTemplateSpec(statefulsetinfo statefulSetInfo) corev1.PodTemplateSpec {
	ginkgo.By("Creating PodTemplateSpec config object")
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:   statefulsetinfo.name,
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
			Volumes: []corev1.Volume{
				{
					Name: statefulsetinfo.volumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: statefulsetinfo.volumeName},
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
					corev1.ResourceStorage: resource.MustParse("4Gi"),
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

func deleteService(svcName string, svcNamespace string, labels map[string]string, serviceSpec corev1.ServiceSpec, k8sClient crclient.Client) {
	svcToDelete := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: svcNamespace,
			Name:      svcName,
		},
		Spec: serviceSpec,
	}
	if len(labels) > 0 {
		svcToDelete.ObjectMeta.Labels = labels
	}
	Expect(k8sClient.Delete(context.TODO(), &svcToDelete)).NotTo(HaveOccurred())
}

func createStatefulSet(statefulsetinfo statefulSetInfo, k8sclient crclient.Client) {
	ginkgo.By("Creating statefulset")
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
	createStorageClass(statefulsetinfo.isInTreeCSI, statefulsetinfo.storageClassName, k8sclient)
	podTemplateSpec := createPodTemplateSpec(statefulsetinfo)
	volClaimTemplate := createPVC(statefulsetinfo)
	deployStatefulSet(statefulsetinfo, volClaimTemplate, podTemplateSpec, k8sclient)
	waitForStatefulSetRunning(statefulsetinfo, k8sclient)
}

func createStorageClass(isIntree bool, storageClassName string, k8sclient crclient.Client) {
	ginkgo.By(fmt.Sprintf("Creating StorageClass object with name: %s", storageClassName))
	volExpansion := true
	bindingMode := storagev1.VolumeBindingWaitForFirstConsumer
	azs := shared.GetAvailabilityZones(e2eCtx.AWSSession)

	provisioner := "ebs.csi.aws.com"
	params := map[string]string{
		"csi.storage.k8s.io/fstype": "xfs",
		"type":                      "io1",
		"iopsPerGB":                 "100",
	}
	allowedTopo := []corev1.TopologySelectorTerm{{
		MatchLabelExpressions: []corev1.TopologySelectorLabelRequirement{{
			Key:    shared.StorageClassOutTreeZoneLabel,
			Values: []string{*azs[0].ZoneName},
		}},
	}}
	if isIntree {
		provisioner = "kubernetes.io/aws-ebs"
		params = map[string]string{
			"type": "gp2",
		}

		allowedTopo = nil
	}
	storageClass := &storagev1.StorageClass{}
	if err := k8sclient.Get(context.TODO(), crclient.ObjectKey{
		Name:      storageClassName,
		Namespace: metav1.NamespaceDefault,
	}, storageClass); err != nil {
		if apierrors.IsNotFound(err) {
			storageClass = &storagev1.StorageClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "storage.k8s.io/v1",
					Kind:       "StorageClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: storageClassName,
				},
				Parameters:           params,
				Provisioner:          provisioner,
				AllowVolumeExpansion: &volExpansion,
				VolumeBindingMode:    &bindingMode,
				AllowedTopologies:    allowedTopo,
			}
			Expect(k8sclient.Create(context.TODO(), storageClass)).NotTo(HaveOccurred())
		}
	}
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

func deleteRetainedVolumes(awsVolIDs []*string) {
	ginkgo.By("Deleting dynamically provisioned volumes")
	ec2Client := ec2.New(e2eCtx.AWSSession)
	for _, volumeID := range awsVolIDs {
		input := &ec2.DeleteVolumeInput{
			VolumeId: aws.String(*volumeID),
		}
		_, err := ec2Client.DeleteVolume(input)
		Expect(err).NotTo(HaveOccurred())
		ginkgo.By(fmt.Sprintf("Deleted dynamically provisioned volume with ID: %s", *volumeID))
	}
}

func deployStatefulSet(statefulsetinfo statefulSetInfo, volClaimTemp corev1.PersistentVolumeClaim, podTemplate corev1.PodTemplateSpec, k8sclient crclient.Client) {
	ginkgo.By(fmt.Sprintf("Deploying Statefulset with name: %s under namespace: %s", statefulsetinfo.name, statefulsetinfo.namespace))
	statefulset := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: statefulsetinfo.name, Namespace: statefulsetinfo.namespace},
		Spec: appsv1.StatefulSetSpec{
			ServiceName:          statefulsetinfo.svcName,
			Replicas:             &statefulsetinfo.replicas,
			Selector:             &metav1.LabelSelector{MatchLabels: statefulsetinfo.selector},
			Template:             podTemplate,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{volClaimTemp},
		},
	}
	err := k8sclient.Create(context.TODO(), &statefulset)
	Expect(err).NotTo(HaveOccurred())
}

func getEvents(namespace string) *corev1.EventList {
	eventsList := &corev1.EventList{}
	if err := e2eCtx.Environment.BootstrapClusterProxy.GetClient().List(context.TODO(), eventsList, crclient.InNamespace(namespace), crclient.MatchingLabels{}); err != nil {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Got error while fetching events of namespace: %s, %s \n", namespace, err.Error())
	}

	return eventsList
}

func getSubnetID(filterKey, filterValue, clusterName string) *string {
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
			{
				Name:   aws.String("tag-key"),
				Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/" + clusterName}),
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
	volIDs := make([]*string, len(pvcList.Items))
	for i, pvc := range pvcList.Items {
		volName := pvc.Spec.VolumeName
		volDescription := &corev1.PersistentVolume{}
		err = k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: info.namespace, Name: volName}, volDescription)
		Expect(err).NotTo(HaveOccurred())

		url := ""
		// Out-of-tree ebs CSI use .Spec.PersistentVolumeSource.CSI path
		// In-tree ebs CSI use .Spec.PersistentVolumeSource.AWSElasticBlockStore path
		if volDescription.Spec.PersistentVolumeSource.CSI != nil {
			url = volDescription.Spec.PersistentVolumeSource.CSI.VolumeHandle
		} else if volDescription.Spec.PersistentVolumeSource.AWSElasticBlockStore != nil {
			str := strings.Split(volDescription.Spec.PersistentVolumeSource.AWSElasticBlockStore.VolumeID, "vol-")
			url = "vol-" + str[1]
		}
		volIDs[i] = &url
	}
	return volIDs
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
	return len(awsMachineList.Items) == eventMachinesCnt
}

func getAWSMachinesForDeployment(namespace string, machineDeployment clusterv1.MachineDeployment) *infrav1.AWSMachineList {
	k8sClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	selector, err := metav1.LabelSelectorAsMap(&machineDeployment.Spec.Selector)
	Expect(err).NotTo(HaveOccurred())
	awsMachineList := &infrav1.AWSMachineList{}
	Expect(k8sClient.List(context.TODO(), awsMachineList, crclient.InNamespace(namespace), crclient.MatchingLabels(selector))).NotTo(HaveOccurred())
	return awsMachineList
}

func makeAWSMachineTemplate(namespace, name, instanceType string, subnetID *string) *infrav1.AWSMachineTemplate {
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
					SSHKeyName:         pointer.String(os.Getenv("AWS_SSH_KEY_NAME")),
				},
			},
		},
	}

	if subnetID != nil {
		resRef := &infrav1.AWSResourceReference{
			ID: subnetID,
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

func makeMachineDeployment(namespace, mdName, clusterName string, az *string, replicas int32) *clusterv1.MachineDeployment {
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
					Version: pointer.String(e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion)),
				},
			},
		},
	}
	if az != nil {
		machineDeployment.Spec.Template.Spec.FailureDomain = az
	}
	return machineDeployment
}

func assertSpotInstanceType(instanceID string) {
	ginkgo.By(fmt.Sprintf("Finding EC2 spot instance with ID: %s", instanceID))
	ec2Client := ec2.New(e2eCtx.AWSSession)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID[strings.LastIndex(instanceID, "/")+1:]),
		},
		Filters: []*ec2.Filter{{Name: aws.String("instance-lifecycle"), Values: aws.StringSlice([]string{"spot"})}},
	}

	result, err := ec2Client.DescribeInstances(input)
	Expect(err).To(BeNil())
	Expect(len(result.Reservations)).To(Equal(1))
	Expect(len(result.Reservations[0].Instances)).To(Equal(1))
}

func assertInstanceMetadataOptions(instanceID string, expected infrav1.InstanceMetadataOptions) {
	ginkgo.By(fmt.Sprintf("Finding EC2 instance with ID: %s", instanceID))
	ec2Client := ec2.New(e2eCtx.AWSSession)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID[strings.LastIndex(instanceID, "/")+1:]),
		},
	}

	result, err := ec2Client.DescribeInstances(input)
	Expect(err).To(BeNil())
	Expect(len(result.Reservations)).To(Equal(1))
	Expect(len(result.Reservations[0].Instances)).To(Equal(1))

	metadataOptions := result.Reservations[0].Instances[0].MetadataOptions
	Expect(metadataOptions).ToNot(BeNil())

	Expect(metadataOptions.HttpTokens).To(HaveValue(Equal(string(expected.HTTPTokens))))
	Expect(metadataOptions.HttpEndpoint).To(HaveValue(Equal(string(expected.HTTPEndpoint))))
	Expect(metadataOptions.InstanceMetadataTags).To(HaveValue(Equal(string(expected.InstanceMetadataTags))))
	Expect(metadataOptions.HttpPutResponseHopLimit).To(HaveValue(Equal(expected.HTTPPutResponseHopLimit)))
}

func terminateInstance(instanceID string) {
	ginkgo.By(fmt.Sprintf("Terminating EC2 instance with ID: %s", instanceID))
	ec2Client := ec2.New(e2eCtx.AWSSession)
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID[strings.LastIndex(instanceID, "/")+1:]),
		},
	}

	result, err := ec2Client.TerminateInstances(input)
	Expect(err).To(BeNil())
	Expect(len(result.TerminatingInstances)).To(Equal(1))
	termCode := int64(32)
	Expect(*result.TerminatingInstances[0].CurrentState.Code).To(Equal(termCode))
}

func verifyElbExists(elbName string, exists bool) {
	ginkgo.By(fmt.Sprintf("Verifying ELB with name %s present", elbName))
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
		ginkgo.By(fmt.Sprintf("ELB with name %s exists", elbName))
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException))
		ginkgo.By(fmt.Sprintf("ELB with name %s doesn't exists", elbName))
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
	ginkgo.By(fmt.Sprintf("Ensuring Statefulset(%s) is running", info.name))
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
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

type conditionAssertion struct {
	conditionType clusterv1.ConditionType
	status        corev1.ConditionStatus
	severity      clusterv1.ConditionSeverity
	reason        string
}

func expectAWSClusterConditions(m *infrav1.AWSCluster, expected []conditionAssertion) {
	Expect(len(m.Status.Conditions)).To(BeNumerically(">=", len(expected)), "number of conditions")
	for _, c := range expected {
		actual := conditions.Get(m, c.conditionType)
		Expect(actual).To(Not(BeNil()))
		Expect(actual.Type).To(Equal(c.conditionType))
		Expect(actual.Status).To(Equal(c.status))
		Expect(actual.Severity).To(Equal(c.severity))
		Expect(actual.Reason).To(Equal(c.reason))
	}
}

func createEFS() *efs.FileSystemDescription {
	efs, err := shared.CreateEFS(e2eCtx, string(uuid.NewUUID()))
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		state, err := shared.GetEFSState(e2eCtx, aws.StringValue(efs.FileSystemId))
		return aws.StringValue(state), err
	}, 2*time.Minute, 5*time.Second).Should(Equal("available"))
	return efs
}

func createSecurityGroupForEFS(clusterName string, vpc *ec2.Vpc) *ec2.CreateSecurityGroupOutput {
	securityGroup, err := shared.CreateSecurityGroup(e2eCtx, clusterName+"-efs-sg", "security group for EFS Access", *(vpc.VpcId))
	Expect(err).NotTo(HaveOccurred())
	nameFilter := &ec2.Filter{
		Name:   aws.String("tag:Name"),
		Values: aws.StringSlice([]string{clusterName + "-node"}),
	}
	nodeSecurityGroups, err := shared.GetSecurityGroupByFilters(e2eCtx, []*ec2.Filter{
		nameFilter,
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(len(nodeSecurityGroups)).To(Equal(1))
	_, err = shared.CreateSecurityGroupIngressRuleWithSourceSG(e2eCtx, aws.StringValue(securityGroup.GroupId), "tcp", 2049, aws.StringValue(nodeSecurityGroups[0].GroupId))
	Expect(err).NotTo(HaveOccurred())
	return securityGroup
}

func createMountTarget(efs *efs.FileSystemDescription, securityGroup *ec2.CreateSecurityGroupOutput, vpc *ec2.Vpc) *efs.MountTargetDescription {
	mt, err := shared.CreateMountTargetOnEFS(e2eCtx, aws.StringValue(efs.FileSystemId), aws.StringValue(vpc.VpcId), aws.StringValue(securityGroup.GroupId))
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() (string, error) {
		state, err := shared.GetMountTargetState(e2eCtx, *mt.MountTargetId)
		return aws.StringValue(state), err
	}, 5*time.Minute, 10*time.Second).Should(Equal("available"))
	return mt
}

func deleteMountTarget(mountTarget *efs.MountTargetDescription) {
	_, err := shared.DeleteMountTarget(e2eCtx, *mountTarget.MountTargetId)
	Expect(err).NotTo(HaveOccurred())
	Eventually(func(g Gomega) {
		_, err = shared.GetMountTarget(e2eCtx, *mountTarget.MountTargetId)
		g.Expect(err).ShouldNot(Equal(nil))
		aerr, ok := err.(awserr.Error)
		g.Expect(ok).To(BeTrue())
		g.Expect(aerr.Code()).To(Equal(efs.ErrCodeMountTargetNotFound))
	}, 5*time.Minute, 10*time.Second).Should(Succeed())
}

// example taken from aws-efs-csi-driver (https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/master/examples/kubernetes/dynamic_provisioning/specs/storageclass.yaml)
func createEFSStorageClass(storageClassName string, clusterClient crclient.Client, efs *efs.FileSystemDescription) {
	storageClass := &storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "storage.k8s.io/v1",
			Kind:       "StorageClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		MountOptions: []string{"tls"},
		Parameters: map[string]string{
			"provisioningMode": "efs-ap",
			"fileSystemId":     aws.StringValue(efs.FileSystemId),
			"directoryPerms":   "700",
			"gidRangeStart":    "1000",
			"gidRangeEnd":      "2000",
		},
		Provisioner: "efs.csi.aws.com",
	}
	Expect(clusterClient.Create(context.TODO(), storageClass)).NotTo(HaveOccurred())
}

// example taken from aws-efs-csi-driver (https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/master/examples/kubernetes/dynamic_provisioning/specs/pod.yaml)
func createPVCForEFS(storageClassName string, clusterClient crclient.Client) {
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "efs-claim",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			StorageClassName: &storageClassName,
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI),
				},
			},
		},
	}
	Expect(clusterClient.Create(context.TODO(), pvc)).NotTo(HaveOccurred())
}

// example taken from aws-efs-csi-driver (https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/master/examples/kubernetes/dynamic_provisioning/specs/pod.yaml)
func createPodWithEFSMount(clusterClient crclient.Client) {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "efs-app",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "app",
					Image:   "centos",
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", "while true; do echo $(date -u) >> /data/out; sleep 5; done"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "persistent-storage",
							MountPath: "/data",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "persistent-storage",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "efs-claim",
						},
					},
				},
			},
		},
	}
	Expect(clusterClient.Create(context.TODO(), pod)).NotTo(HaveOccurred())
}
