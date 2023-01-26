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

package managed

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// EKS related constants.
const (
	EKSManagedPoolFlavor                              = "eks-managedmachinepool"
	EKSControlPlaneOnlyFlavor                         = "eks-control-plane-only"
	EKSControlPlaneOnlyWithAddonFlavor                = "eks-control-plane-only-withaddon"
	EKSMachineDeployOnlyFlavor                        = "eks-machine-deployment-only"
	EKSManagedMachinePoolOnlyFlavor                   = "eks-managed-machinepool-only"
	EKSManagedMachinePoolWithLaunchTemplateOnlyFlavor = "eks-managed-machinepool-with-launch-template-only"
	EKSMachinePoolOnlyFlavor                          = "eks-machinepool-only"
	EKSIPv6ClusterFlavor                              = "eks-ipv6-cluster"
	EKSControlPlaneOnlyLegacyFlavor                   = "eks-control-plane-only-legacy"
)

type DefaultConfigClusterFn func(clusterName, namespace string) clusterctl.ConfigClusterInput

func getEKSClusterName(namespace, clusterName string) string {
	return fmt.Sprintf("%s_%s-control-plane", namespace, clusterName)
}

func getEKSNodegroupName(namespace, clusterName string) string {
	return fmt.Sprintf("%s_%s-pool-0", namespace, clusterName)
}

func getEKSNodegroupWithLaunchTemplateName(namespace, clusterName string) string {
	return fmt.Sprintf("%s_%s-pool-lt-0", namespace, clusterName)
}

func getControlPlaneName(clusterName string) string {
	return fmt.Sprintf("%s-control-plane", clusterName)
}

func getASGName(clusterName string) string {
	return fmt.Sprintf("%s-mp-0", clusterName)
}

func verifyClusterActiveAndOwned(eksClusterName string, sess client.ConfigProvider) {
	cluster, err := getEKSCluster(eksClusterName, sess)
	Expect(err).NotTo(HaveOccurred())

	tagName := infrav1.ClusterTagKey(eksClusterName)
	tagValue, ok := cluster.Tags[tagName]
	Expect(ok).To(BeTrue(), "expecting the cluster owned tag to exist")
	Expect(*tagValue).To(BeEquivalentTo(string(infrav1.ResourceLifecycleOwned)))

	Expect(*cluster.Status).To(BeEquivalentTo(eks.ClusterStatusActive))
}

func getEKSCluster(eksClusterName string, sess client.ConfigProvider) (*eks.Cluster, error) {
	eksClient := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	}
	result, err := eksClient.DescribeCluster(input)

	return result.Cluster, err
}

func getEKSClusterAddon(eksClusterName, addonName string, sess client.ConfigProvider) (*eks.Addon, error) {
	eksClient := eks.New(sess)

	describeInput := &eks.DescribeAddonInput{
		AddonName:   &addonName,
		ClusterName: &eksClusterName,
	}
	describeOutput, err := eksClient.DescribeAddon(describeInput)
	if err != nil {
		return nil, fmt.Errorf("describing eks addon %s: %w", addonName, err)
	}

	return describeOutput.Addon, nil
}

func verifySecretExists(ctx context.Context, secretName, namespace string, k8sclient crclient.Client) {
	secret := &corev1.Secret{}
	err := k8sclient.Get(ctx, apimachinerytypes.NamespacedName{Name: secretName, Namespace: namespace}, secret)

	Expect(err).ShouldNot(HaveOccurred())
}

func verifyConfigMapExists(ctx context.Context, name, namespace string, k8sclient crclient.Client) {
	cm := &corev1.ConfigMap{}
	Eventually(func() error {
		return k8sclient.Get(ctx, apimachinerytypes.NamespacedName{Name: name, Namespace: namespace}, cm)
	}, 2*time.Minute, 5*time.Second).Should(Succeed())
}

func VerifyRoleExistsAndOwned(roleName string, eksClusterName string, checkOwned bool, sess client.ConfigProvider) {
	iamClient := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	output, err := iamClient.GetRole(input)
	Expect(err).ShouldNot(HaveOccurred())

	if checkOwned {
		found := false
		expectedTagName := infrav1.ClusterAWSCloudProviderTagKey(eksClusterName)
		for _, tag := range output.Role.Tags {
			if *tag.Key == expectedTagName && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
				found = true
				break
			}
		}
		Expect(found).To(BeTrue(), "expecting the cluster owned tag to exist")
	}
}

func verifyManagedNodeGroup(eksClusterName, nodeGroupName string, checkOwned bool, sess client.ConfigProvider) {
	eksClient := eks.New(sess)
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodeGroupName),
	}
	result, err := eksClient.DescribeNodegroup(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(*result.Nodegroup.Status).To(BeEquivalentTo(eks.NodegroupStatusActive))

	if checkOwned {
		tagName := infrav1.ClusterAWSCloudProviderTagKey(eksClusterName)
		tagValue, ok := result.Nodegroup.Tags[tagName]
		Expect(ok).To(BeTrue(), "expecting the cluster owned tag to exist")
		Expect(*tagValue).To(BeEquivalentTo(string(infrav1.ResourceLifecycleOwned)))
	}
}

func verifyASG(eksClusterName, asgName string, checkOwned bool, sess client.ConfigProvider) {
	asgClient := autoscaling.New(sess)
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	result, err := asgClient.DescribeAutoScalingGroups(input)
	Expect(err).NotTo(HaveOccurred())
	for _, instance := range result.AutoScalingGroups[0].Instances {
		Expect(*instance.LifecycleState).To(Equal("InService"), "expecting the instance in service")
	}

	if checkOwned {
		found := false
		for _, tag := range result.AutoScalingGroups[0].Tags {
			if *tag.Key == infrav1.ClusterAWSCloudProviderTagKey(eksClusterName) {
				Expect(*tag.Value).To(Equal(string(infrav1.ResourceLifecycleOwned)))
				found = true
				break
			}
		}
		Expect(found).To(BeTrue(), "expecting the cluster owned tag to exist")
	}
}
