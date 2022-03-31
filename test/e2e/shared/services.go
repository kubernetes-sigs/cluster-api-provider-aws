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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
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
