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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
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
	EKSClusterClassFlavor                             = "eks-clusterclass"
	EKSAuthAPIAndConfigMapFlavor                      = "eks-auth-api-and-config-map"
	EKSAuthBootstrapDisabledFlavor                    = "eks-auth-bootstrap-disabled"
)

const (
	clientRequestTimeout       = 2 * time.Minute
	clientRequestCheckInterval = 5 * time.Second
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

func verifyClusterActiveAndOwned(ctx context.Context, eksClusterName string, sess *aws.Config) {
	var (
		cluster *ekstypes.Cluster
		err     error
	)
	Eventually(func() error {
		cluster, err = getEKSCluster(ctx, eksClusterName, sess)
		return err
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), fmt.Sprintf("eventually failed trying to get EKS Cluster %q", eksClusterName))

	tagName := infrav1.ClusterTagKey(eksClusterName)
	tagValue, ok := cluster.Tags[tagName]
	Expect(ok).To(BeTrue(), "expecting the cluster owned tag to exist")
	Expect(tagValue).To(BeEquivalentTo(string(infrav1.ResourceLifecycleOwned)))
	Expect(cluster.Status).To(BeEquivalentTo(ekstypes.ClusterStatusActive))
}

func getEKSCluster(ctx context.Context, eksClusterName string, sess *aws.Config) (*ekstypes.Cluster, error) {
	eksClient := eks.NewFromConfig(*sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(eksClusterName),
	}
	result, err := eksClient.DescribeCluster(ctx, input)

	return result.Cluster, err
}

func verifyClusterAuthenticationMode(ctx context.Context, eksClusterName string, expectedAuthMode ekstypes.AuthenticationMode, sess *aws.Config) {
	var (
		cluster *ekstypes.Cluster
		err     error
	)
	Eventually(func() error {
		cluster, err = getEKSCluster(ctx, eksClusterName, sess)
		return err
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), fmt.Sprintf("eventually failed trying to get EKS Cluster %q", eksClusterName))

	Expect(cluster.AccessConfig).ToNot(BeNil(), "expecting AccessConfig to be set on the cluster")
	Expect(cluster.AccessConfig.AuthenticationMode).To(BeEquivalentTo(expectedAuthMode),
		fmt.Sprintf("expecting authentication mode to be %s, got %s", expectedAuthMode, cluster.AccessConfig.AuthenticationMode))
}

func getEKSClusterAddon(ctx context.Context, eksClusterName, addonName string, sess *aws.Config) (*ekstypes.Addon, error) {
	eksClient := eks.NewFromConfig(*sess)

	describeInput := &eks.DescribeAddonInput{
		AddonName:   &addonName,
		ClusterName: &eksClusterName,
	}

	describeOutput, err := eksClient.DescribeAddon(ctx, describeInput)
	if err != nil {
		return nil, fmt.Errorf("describing eks addon %s: %w", addonName, err)
	}

	return describeOutput.Addon, nil
}

func verifySecretExists(ctx context.Context, secretName, namespace string, k8sclient crclient.Client) {
	secret := &corev1.Secret{}
	Eventually(func() error {
		return k8sclient.Get(ctx, apimachinerytypes.NamespacedName{Name: secretName, Namespace: namespace}, secret)
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), fmt.Sprintf("eventually failed trying to verify Secret %q exists", secretName))
}

func verifyConfigMapExists(ctx context.Context, name, namespace string, k8sclient crclient.Client) {
	cm := &corev1.ConfigMap{}
	Eventually(func() error {
		return k8sclient.Get(ctx, apimachinerytypes.NamespacedName{Name: name, Namespace: namespace}, cm)
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), fmt.Sprintf("eventually failed trying to verify ConfigMap %q exists", name))
}

func VerifyRoleExistsAndOwned(ctx context.Context, roleName string, eksClusterName string, checkOwned bool, sess *aws.Config) {
	iamClient := iam.NewFromConfig(*sess)
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	var (
		output *iam.GetRoleOutput
		err    error
	)

	Eventually(func() error {
		output, err = iamClient.GetRole(ctx, input)
		return err
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), fmt.Sprintf("eventually failed trying to get IAM Role %q", roleName))

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

func verifyManagedNodeGroup(ctx context.Context, eksClusterName, nodeGroupName string, checkOwned bool, sess *aws.Config) {
	eksClient := eks.NewFromConfig(*sess)
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(eksClusterName),
		NodegroupName: aws.String(nodeGroupName),
	}
	var (
		result *eks.DescribeNodegroupOutput
		err    error
	)

	Eventually(func() error {
		result, err = eksClient.DescribeNodegroup(ctx, input)
		if err != nil {
			return fmt.Errorf("error describing nodegroup: %w", err)
		}

		nodeGroupStatus := result.Nodegroup.Status
		if nodeGroupStatus != ekstypes.NodegroupStatusActive {
			return fmt.Errorf("expected nodegroup.Status to be %q, was %q instead", ekstypes.NodegroupStatusActive, nodeGroupStatus)
		}

		return nil
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed(), "eventually failed trying to describe EKS Node group")

	if checkOwned {
		tagName := infrav1.ClusterAWSCloudProviderTagKey(eksClusterName)
		tagValue, ok := result.Nodegroup.Tags[tagName]
		Expect(ok).To(BeTrue(), "expecting the cluster owned tag to exist")
		Expect(tagValue).To(BeEquivalentTo(string(infrav1.ResourceLifecycleOwned)))
	}
}

func verifyASG(eksClusterName, asgName string, checkOwned bool, cfg *aws.Config) {
	asgClient := autoscaling.NewFromConfig(*cfg)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			asgName,
		},
	}
	var (
		result *autoscaling.DescribeAutoScalingGroupsOutput
		err    error
	)

	Eventually(func() error {
		result, err = asgClient.DescribeAutoScalingGroups(context.TODO(), input)
		return err
	}, clientRequestTimeout, clientRequestCheckInterval).Should(Succeed())

	for _, instance := range result.AutoScalingGroups[0].Instances {
		Expect(string(instance.LifecycleState)).To(Equal("InService"), "expecting the instance in service")
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
