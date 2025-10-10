//go:build e2e
// +build e2e

/*
Copyright 2025 The Kubernetes Authors.

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

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// EKS authentication mode e2e tests.
var _ = ginkgo.Describe("[managed] [auth] EKS authentication mode tests", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "auth"
		clusterName string
	)

	shared.ConditionalIt(runGeneralTests, "should create a cluster with access entries", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		eksClusterName := getEKSClusterName(namespace.Name, clusterName)

		ginkgo.By("should create an EKS control plane with access entries enabled")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSAuthAPIAndConfigMapFlavor,
				ControlPlaneMachineCount: 1,
				WorkerMachineCount:       0,
			}
		})

		ginkgo.By("EKS cluster should be active")
		verifyClusterActiveAndOwned(ctx, eksClusterName, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("should create a cluster with access entries")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
						return ManagedClusterSpecInput{
										E2EConfig:                e2eCtx.E2EConfig,
										ConfigClusterFn:          defaultConfigCluster,
										BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
										AWSSession:               e2eCtx.BootstrapUserAWSSession,
										Namespace:                namespace,
										ClusterName:              clusterName,
										Flavour:                  EKSControlPlaneOnlyWithAccessEntriesFlavor,
										ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
										WorkerMachineCount:       0,
						}
		})

		ginkgo.By("should have created the expected access entries")
		expectedEntries := []ekscontrolplanev1.AccessEntry{
			{
				PrincipalARN:     "arn:aws:iam::123456789012:role/KubernetesAdmin",
				Type:             "STANDARD",
				Username:         "kubernetes-admin",
				KubernetesGroups: []string{"system:masters"},
				AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
					{
						PolicyARN: "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy",
						AccessScope: ekscontrolplanev1.AccessScope{
							Type: "cluster",
						},
					},
				},
			},
			{
				PrincipalARN:     "arn:aws:iam::123456789012:role/DeveloperRole",
				Type:             "STANDARD",
				Username:         "developer",
				KubernetesGroups: []string{"developers"},
				AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
					{
						PolicyARN: "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy",
						AccessScope: ekscontrolplanev1.AccessScope{
							Type:       "namespace",
							Namespaces: []string{"default"},
						},
					},
				},
			},
		}
		verifyAccessEntries(ctx, eksClusterName, expectedEntries, e2eCtx.BootstrapUserAWSSession)


		ginkgo.By("EKS cluster should be active")
		verifyClusterActiveAndOwned(ctx, eksClusterName, e2eCtx.BootstrapUserAWSSession)

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
