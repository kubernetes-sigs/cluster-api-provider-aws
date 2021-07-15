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

	//. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/test/framework/clusterctl"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

const (
	EKSManagedPoolFlavor = "eks-managedmachinepool"

	EKSControlPlaneOnlyFlavor          = "eks-control-plane-only"
	EKSControlPlaneOnlyWithAddonFlavor = "eks-control-plane-only-withaddon"
	EKSMachineDeployOnlyFlavor         = "eks-machine-deployment-only"
	EKSManagedPoolOnlyFlavor           = "eks-managed-machinepool-only"
)

type DefaultConfigClusterFn func(clusterName, namespace string) clusterctl.ConfigClusterInput

func getEKSClusterName(namespace, clusterName string) string {
	return fmt.Sprintf("%s_%s-control-plane", namespace, clusterName)
}

func getEKSNodegroupName(namespace, clusterName string) string {
	return fmt.Sprintf("%s_%s-pool-0", namespace, clusterName)
}

func getControlPlaneName(clusterName string) string {
	return fmt.Sprintf("%s-control-plane", clusterName)
}

func verifyClusterActiveAndOwned(eksClusterName, clusterName string, sess client.ConfigProvider) {
	cluster, err := getEKSCluster(eksClusterName, sess)
	Expect(err).NotTo(HaveOccurred())

	tagName := infrav1.ClusterTagKey(clusterName)
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
	err := k8sclient.Get(ctx, apimachinerytypes.NamespacedName{Name: "aws-auth", Namespace: metav1.NamespaceSystem}, cm)

	Expect(err).ShouldNot(HaveOccurred())
}

func verifyRoleExistsAndOwned(roleName string, clusterName string, checkOwned bool, sess client.ConfigProvider) {
	iamClient := iam.New(sess)
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	output, err := iamClient.GetRole(input)
	Expect(err).ShouldNot(HaveOccurred())

	if checkOwned {
		found := false
		expectedTagName := infrav1.ClusterAWSCloudProviderTagKey(clusterName)
		for _, tag := range output.Role.Tags {
			if *tag.Key == expectedTagName && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
				found = true
				break
			}
		}
		Expect(found).To(BeTrue(), "expecting the cluster owned tag to exist")
	}
}

func verifyManagedNodeGroup(clusterName, eksClusterName, nodeGroupName string, checkOwned bool, sess client.ConfigProvider) {
	eksClient := eks.New(sess)
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodeGroupName),
	}
	result, err := eksClient.DescribeNodegroup(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(*result.Nodegroup.Status).To(BeEquivalentTo(eks.NodegroupStatusActive))

	if checkOwned {
		tagName := infrav1.ClusterAWSCloudProviderTagKey(clusterName)
		tagValue, ok := result.Nodegroup.Tags[tagName]
		Expect(ok).To(BeTrue(), "expecting the cluster owned tag to exist")
		Expect(*tagValue).To(BeEquivalentTo(string(infrav1.ResourceLifecycleOwned)))
	}
}
