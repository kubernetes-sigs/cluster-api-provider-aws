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

// Package managed implements a test for creating a managed cluster using CAPA.
package managed

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// ManagedClusterSpecInput is the input for ManagedClusterSpec.
type ManagedClusterSpecInput struct {
	E2EConfig                *clusterctl.E2EConfig
	ConfigClusterFn          DefaultConfigClusterFn
	BootstrapClusterProxy    framework.ClusterProxy
	AWSSession               aws.Config
	AWSSessionV2             *aws.Config
	Namespace                *corev1.Namespace
	ClusterName              string
	Flavour                  string
	ControlPlaneMachineCount int64
	WorkerMachineCount       int64
	KubernetesVersion        string
	CluserSpecificRoles      bool
}

// ManagedClusterSpec implements a test for creating a managed cluster using CAPA.
func ManagedClusterSpec(ctx context.Context, inputGetter func() ManagedClusterSpecInput) {
	input := inputGetter()
	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.ConfigClusterFn).ToNot(BeNil(), "Invalid argument. input.ConfigClusterFn can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.ClusterName).ShouldNot(BeEmpty(), "Invalid argument. input.ClusterName can't be empty")

	ginkgo.By(fmt.Sprintf("creating an applying the %s template", input.Flavour))
	configCluster := input.ConfigClusterFn(input.ClusterName, input.Namespace.Name)
	configCluster.Flavor = input.Flavour
	configCluster.ControlPlaneMachineCount = ptr.To[int64](input.ControlPlaneMachineCount)
	configCluster.WorkerMachineCount = ptr.To[int64](input.WorkerMachineCount)
	if input.KubernetesVersion != "" {
		configCluster.KubernetesVersion = input.KubernetesVersion
	}
	err := shared.ApplyTemplate(ctx, configCluster, input.BootstrapClusterProxy)
	Expect(err).ShouldNot(HaveOccurred())

	ginkgo.By("Waiting for the cluster to be provisioned")
	start := time.Now()
	cluster := framework.DiscoveryAndWaitForCluster(ctx, framework.DiscoveryAndWaitForClusterInput{
		Getter:    input.BootstrapClusterProxy.GetClient(),
		Namespace: configCluster.Namespace,
		Name:      configCluster.ClusterName,
	}, input.E2EConfig.GetIntervals("", "wait-cluster")...)
	duration := time.Since(start)
	ginkgo.By(fmt.Sprintf("Finished waiting for cluster after: %s", duration))
	Expect(cluster).NotTo(BeNil())

	ginkgo.By("Checking EKS cluster is active")
	eksClusterName := getEKSClusterName(input.Namespace.Name, input.ClusterName)
	verifyClusterActiveAndOwned(ctx, eksClusterName, input.AWSSessionV2)

	if input.CluserSpecificRoles {
		ginkgo.By("Checking that the cluster specific IAM role exists")
		VerifyRoleExistsAndOwned(ctx, fmt.Sprintf("%s-iam-service-role", input.ClusterName), eksClusterName, true, input.AWSSessionV2)
	} else {
		ginkgo.By("Checking that the cluster default IAM role exists")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, eksClusterName, false, input.AWSSessionV2)
	}

	ginkgo.By("Checking kubeconfig secrets exist")
	bootstrapClient := input.BootstrapClusterProxy.GetClient()
	verifySecretExists(ctx, fmt.Sprintf("%s-kubeconfig", input.ClusterName), input.Namespace.Name, bootstrapClient)
	verifySecretExists(ctx, fmt.Sprintf("%s-user-kubeconfig", input.ClusterName), input.Namespace.Name, bootstrapClient)

	// this will not be created unless there are worker machines or set by IAMAuthenticatorConfig on the managed control plane spec
	if input.WorkerMachineCount > 0 {
		ginkgo.By("Checking that aws-iam-authenticator config map exists")
		workloadClusterProxy := input.BootstrapClusterProxy.GetWorkloadCluster(ctx, input.Namespace.Name, input.ClusterName)
		workloadClient := workloadClusterProxy.GetClient()
		verifyConfigMapExists(ctx, "aws-auth", metav1.NamespaceSystem, workloadClient)
	}
}

// DeleteClusterSpecInput is the input to DeleteClusterSpec.
type DeleteClusterSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	BootstrapClusterProxy framework.ClusterProxy
	Namespace             *corev1.Namespace
	ClusterName           string
}
