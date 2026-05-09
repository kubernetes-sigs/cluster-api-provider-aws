//go:build e2e
// +build e2e

/*
Copyright 2026 The Kubernetes Authors.

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
	"os"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
)

// EKS pod identity association e2e tests.
var _ = ginkgo.Describe("[managed] [pod-identities] EKS pod identity association tests", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "pod-identities"
		clusterName string
	)

	shared.ConditionalIt(runGeneralTests, "should create a cluster with pod identity associations", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		eksClusterName := getEKSClusterName(namespace.Name, clusterName)

		ginkgo.By("creating an EKS cluster with pod identity associations")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSControlPlaneOnlyWithPodIdentitiesFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
			}
		})

		ginkgo.By("verifying EKS cluster is active")
		verifyClusterActiveAndOwned(ctx, eksClusterName, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("verifying the expected pod identity associations were created")
		controllersRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/controllers.cluster-api-provider-aws.sigs.k8s.io", os.Getenv(shared.AwsAccountID))
		expected := []ekscontrolplanev1.PodIdentityAssociation{
			{
				ServiceAccountNamespace: "kube-system",
				ServiceAccountName:      "aws-node",
				RoleARN:                 controllersRoleARN,
			},
			{
				ServiceAccountNamespace: "default",
				ServiceAccountName:      "test-sa",
				RoleARN:                 controllersRoleARN,
			},
		}
		verifyPodIdentityAssociations(ctx, eksClusterName, expected, e2eCtx.BootstrapUserAWSSession)

		mgmtClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
		controlPlaneName := getControlPlaneName(clusterName)
		controlPlaneKey := crclient.ObjectKey{Namespace: namespace.Name, Name: controlPlaneName}
		nodesRoleARN := fmt.Sprintf("arn:aws:iam::%s:role/nodes.cluster-api-provider-aws.sigs.k8s.io", os.Getenv(shared.AwsAccountID))

		patchAssociations := func(mutate func(*ekscontrolplanev1.AWSManagedControlPlane)) {
			cp := &ekscontrolplanev1.AWSManagedControlPlane{}
			Expect(mgmtClient.Get(ctx, controlPlaneKey, cp)).To(Succeed())
			helper, err := patch.NewHelper(cp, mgmtClient)
			Expect(err).ToNot(HaveOccurred())
			mutate(cp)
			Expect(helper.Patch(ctx, cp)).To(Succeed())
		}

		ginkgo.By("adding a new pod identity association")
		patchAssociations(func(cp *ekscontrolplanev1.AWSManagedControlPlane) {
			cp.Spec.PodIdentityAssociations = append(cp.Spec.PodIdentityAssociations, ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountNamespace: "kube-system",
				ServiceAccountName:      "coredns",
				RoleARN:                 controllersRoleARN,
			})
		})
		expected = append(expected, ekscontrolplanev1.PodIdentityAssociation{
			ServiceAccountNamespace: "kube-system",
			ServiceAccountName:      "coredns",
			RoleARN:                 controllersRoleARN,
		})
		verifyPodIdentityAssociations(ctx, eksClusterName, expected, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("setting targetRoleARN on an existing pod identity association")
		patchAssociations(func(cp *ekscontrolplanev1.AWSManagedControlPlane) {
			for i := range cp.Spec.PodIdentityAssociations {
				a := &cp.Spec.PodIdentityAssociations[i]
				if a.ServiceAccountNamespace == "default" && a.ServiceAccountName == "test-sa" {
					a.TargetRoleARN = nodesRoleARN
				}
			}
		})
		for i := range expected {
			if expected[i].ServiceAccountNamespace == "default" && expected[i].ServiceAccountName == "test-sa" {
				expected[i].TargetRoleARN = nodesRoleARN
			}
		}
		verifyPodIdentityAssociations(ctx, eksClusterName, expected, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("removing targetRoleARN from an existing pod identity association")
		patchAssociations(func(cp *ekscontrolplanev1.AWSManagedControlPlane) {
			for i := range cp.Spec.PodIdentityAssociations {
				a := &cp.Spec.PodIdentityAssociations[i]
				if a.ServiceAccountNamespace == "default" && a.ServiceAccountName == "test-sa" {
					a.TargetRoleARN = ""
				}
			}
		})
		for i := range expected {
			if expected[i].ServiceAccountNamespace == "default" && expected[i].ServiceAccountName == "test-sa" {
				expected[i].TargetRoleARN = ""
			}
		}
		verifyPodIdentityAssociations(ctx, eksClusterName, expected, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("removing a pod identity association")
		patchAssociations(func(cp *ekscontrolplanev1.AWSManagedControlPlane) {
			filtered := cp.Spec.PodIdentityAssociations[:0]
			for _, a := range cp.Spec.PodIdentityAssociations {
				if a.ServiceAccountNamespace == "kube-system" && a.ServiceAccountName == "aws-node" {
					continue
				}
				filtered = append(filtered, a)
			}
			cp.Spec.PodIdentityAssociations = filtered
		})
		remaining := expected[:0]
		for _, a := range expected {
			if a.ServiceAccountNamespace == "kube-system" && a.ServiceAccountName == "aws-node" {
				continue
			}
			remaining = append(remaining, a)
		}
		expected = remaining
		verifyPodIdentityAssociations(ctx, eksClusterName, expected, e2eCtx.BootstrapUserAWSSession)

		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find CAPI cluster")

		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			ClusterProxy:         e2eCtx.Environment.BootstrapClusterProxy,
			Cluster:              cluster,
			ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
			ArtifactFolder:       e2eCtx.Settings.ArtifactFolder,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)
	})
})
