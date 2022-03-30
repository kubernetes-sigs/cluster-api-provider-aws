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

package shared

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateLBService(e2eCtx *E2EContext, svcNamespace string, svcName string, k8sclient crclient.Client) string {
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
	CreateService(svcName, svcNamespace, nil, svcSpec, k8sclient)
	elbName := ""
	Eventually(func() bool {
		svcCreated := &corev1.Service{}
		err := k8sclient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: svcNamespace, Name: svcName}, svcCreated)
		Expect(err).NotTo(HaveOccurred())
		if lbs := len(svcCreated.Status.LoadBalancer.Ingress); lbs > 0 {
			ingressHostname := svcCreated.Status.LoadBalancer.Ingress[0].Hostname
			elbName = strings.Split(ingressHostname, "-")[0]
			return true
		}
		return false
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-create-service")...).Should(BeTrue())
	Byf("Created Load Balancer service and ELB name is: %s", elbName)
	return elbName
}

func DeleteLBService(svcNamespace string, svcName string, k8sclient crclient.Client) {
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

func CreateService(svcName string, svcNamespace string, labels map[string]string, serviceSpec corev1.ServiceSpec, k8sClient crclient.Client) {
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

func CreateDefaultNginxDeployment(deploymentNamespace, deploymentName string, k8sClient crclient.Client) {
	Byf("Creating Deployment with name: %s under namespace: %s", deploymentName, deploymentNamespace)
	deployment := defaultNginxDeployment(deploymentName, deploymentNamespace)
	Expect(k8sClient.Create(context.TODO(), &deployment)).NotTo(HaveOccurred())
	Eventually(func() bool {
		getDeployment := &v1.Deployment{}
		err := k8sClient.Get(context.TODO(), apimachinerytypes.NamespacedName{Namespace: deploymentNamespace, Name: deploymentName}, getDeployment)
		Expect(err).NotTo(HaveOccurred())
		for _, c := range getDeployment.Status.Conditions {
			if c.Type == v1.DeploymentAvailable && c.Status == corev1.ConditionTrue {
				return getDeployment.Status.AvailableReplicas > 0
			}
		}
		return false
	}, 60*time.Second).Should(BeTrue())
}

func DeleteDefaultNginxDeployment(deploymentNamespace, deploymentName string, k8sClient crclient.Client) {
	Byf("Deleting Deployment with name: %s under namespace: %s", deploymentName, deploymentNamespace)
	deployment := defaultNginxDeployment(deploymentName, deploymentNamespace)
	Expect(k8sClient.Delete(context.TODO(), &deployment)).NotTo(HaveOccurred())
}

func deleteService(svcName, svcNamespace string, labels map[string]string, serviceSpec corev1.ServiceSpec, k8sClient crclient.Client) {
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

func VerifyElbExists(e2eCtx *E2EContext, elbName string, exists bool) {
	Byf("Verifying ELB with name %s present and instances are attached", elbName)
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(elbName),
		},
	}
	elbClient := elb.New(e2eCtx.AWSSession)
	elbsOutput, err := elbClient.DescribeLoadBalancers(input)
	if exists {
		Expect(err).NotTo(HaveOccurred())
		Expect(len(elbsOutput.LoadBalancerDescriptions)).To(Equal(1))
		Expect(len(elbsOutput.LoadBalancerDescriptions[0].Instances)).Should(BeNumerically(">=", 1))
		Byf("ELB with name %s exists", elbName)
	} else {
		aerr, ok := err.(awserr.Error)
		Expect(ok).To(BeTrue())
		Expect(aerr.Code()).To(Equal(elb.ErrCodeAccessPointNotFoundException))
		Byf("ELB with name %s doesn't exists", elbName)
	}
}

func defaultNginxDeployment(deploymentName, deploymentNamespace string) v1.Deployment {
	selector, err := metav1.ParseToLabelSelector("app=nginx")
	Expect(err).To(BeNil())
	deploymentSpec := v1.DeploymentSpec{
		Selector: selector,
		Replicas: pointer.Int32(1),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nginx",
				Labels: map[string]string{
					"app": "nginx",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "k8s.gcr.io/nginx-slim:0.8",
						Ports: []corev1.ContainerPort{{
							Name: "nginx-port", ContainerPort: int32(80),
						}},
					},
				},
			},
		},
	}
	return v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: deploymentNamespace,
			Name:      deploymentName,
			Labels: map[string]string{
				"app": "nginx",
			},
		},
		Spec: deploymentSpec,
	}
}
