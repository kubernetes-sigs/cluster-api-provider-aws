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

package e2e_new

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("functional tests", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
	)

	BeforeEach(func() {
		Expect(bootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = setupSpecNamespace(ctx, "functional-tests", bootstrapClusterProxy, artifactFolder)
	})

	Describe("Create cluster with name having more than 22 characters", func() {
		It("Cluster should be provisioned and deleted", func() {
			By("Creating a workload cluster with single control plane")
			clusterName := fmt.Sprintf("long-cluster-name-cluster-%s", util.RandomString(20))
			cluster := createCluster(ctx, clusterName, namespace.Name)

			By("Deleting the Cluster")
			framework.DeleteClusterAndWait(ctx, framework.DeleteClusterAndWaitInput{Client: bootstrapClusterProxy.GetClient(), Cluster: cluster}, e2eConfig.GetIntervals("", "wait-delete-cluster")...)
		})
	})

	Describe("Provisioning LoadBalancer dynamically and deleting on cluster deletion", func() {
		It("It should create and delete Load Balancer", func() {
			By("Creating a cluster")
			clusterName := fmt.Sprintf("service-lb-cluster-%s", util.RandomString(8))
			cluster := createCluster(ctx, clusterName, namespace.Name)
			clusterClient := bootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, clusterName).GetClient()

			By("Creating the LB service")
			lbServiceName := "test-svc-" + util.RandomString(6)
			elbName := createLBService(metav1.NamespaceDefault, lbServiceName, clusterClient)
			verifyElbExists(elbName, true)

			By("Deleting the Cluster")
			framework.DeleteClusterAndWait(ctx, framework.DeleteClusterAndWaitInput{Client: bootstrapClusterProxy.GetClient(), Cluster: cluster}, e2eConfig.GetIntervals("", "wait-delete-cluster")...)

			By("Verifying whether provisioned LB deleted")
			verifyElbExists(elbName, false)
		})
	})

	AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		dumpSpecResourcesAndCleanup(ctx, "", bootstrapClusterProxy, artifactFolder, namespace, e2eConfig.GetIntervals, skipCleanup)
	})
})

func createCluster(ctx context.Context, clusterName string, namespace string) *clusterv1.Cluster {
	cluster, _, _ := clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
		ClusterProxy: bootstrapClusterProxy,
		ConfigCluster: clusterctl.ConfigClusterInput{
			LogFolder:                filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
			ClusterctlConfigPath:     clusterctlConfigPath,
			KubeconfigPath:           bootstrapClusterProxy.GetKubeconfigPath(),
			InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
			Flavor:                   clusterctl.DefaultFlavor,
			Namespace:                namespace,
			ClusterName:              clusterName,
			KubernetesVersion:        e2eConfig.GetVariable(KubernetesVersion),
			ControlPlaneMachineCount: pointer.Int64Ptr(1),
			WorkerMachineCount:       pointer.Int64Ptr(1),
		},
		CNIManifestPath:              e2eConfig.GetVariable(CNIPath),
		WaitForClusterIntervals:      e2eConfig.GetIntervals("", "wait-cluster"),
		WaitForControlPlaneIntervals: e2eConfig.GetIntervals("", "wait-control-plane"),
		WaitForMachineDeployments:    e2eConfig.GetIntervals("", "wait-worker-nodes"),
	})

	return cluster
}

func createLBService(svcNamespace string, svcName string, k8sclient crclient.Client) string {
	Byf("Creating service of type Kubernetes Load Balancer: name=%s, namespace=%s", svcName, svcNamespace)
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
	return elbName
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

func verifyElbExists(elbName string, exists bool) {
	Byf("Verifying existence of AWS ELB Load Balancer: name=%s", elbName)
	elbClient := elb.New(awsSession)
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(elbName),
		},
	}
	elbsOutput, err := elbClient.DescribeLoadBalancers(input)
	if exists {
		Expect(err).NotTo(HaveOccurred())
		count := len(elbsOutput.LoadBalancerDescriptions)
		Expect(count).To(Equal(1), "Number of AWS ELBs not equal to 1: name=%s, count=%s", elbName, count)
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException), "AWS ELB found when it shouldn't: name=%s", elbName)
	}
}
