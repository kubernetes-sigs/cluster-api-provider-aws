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

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Context("[managed] [EKS] [smoke] [PR-Blocking]", func() {
	var (
		namespace         *corev1.Namespace
		ctx               context.Context
		requiredResources *shared.TestResource
		specName          = "eks-quick-start"
		clusterName       string
	)

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
	})

	ginkgo.Describe("Running the EKS quick-start spec", func() {
		ginkgo.BeforeEach(func() {
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "eks-quick-start-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		ginkgo.It("should create an EKS cluster with control plane and managed node group", func() {
			eksClusterName := getEKSClusterName(namespace.Name, clusterName)

			ginkgo.By("default iam role should exist")
			VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, eksClusterName, false, e2eCtx.AWSSession)

			ginkgo.By("should create an EKS control plane")
			ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
				return ManagedClusterSpecInput{
					E2EConfig:                e2eCtx.E2EConfig,
					ConfigClusterFn:          defaultConfigCluster,
					BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
					AWSSession:               e2eCtx.BootstrapUserAWSSession,
					Namespace:                namespace,
					ClusterName:              clusterName,
					Flavour:                  EKSControlPlaneOnlyFlavor,
					ControlPlaneMachineCount: 1,
					WorkerMachineCount:       0,
				}
			})

			ginkgo.By("should create a managed node pool")
			MachinePoolSpec(ctx, func() MachinePoolSpecInput {
				return MachinePoolSpecInput{
					E2EConfig:             e2eCtx.E2EConfig,
					ConfigClusterFn:       defaultConfigCluster,
					BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
					AWSSession:            e2eCtx.BootstrapUserAWSSession,
					Namespace:             namespace,
					ClusterName:           clusterName,
					IncludeScaling:        false,
					Cleanup:               false,
					ManagedMachinePool:    true,
					Flavor:                EKSManagedMachinePoolOnlyFlavor,
				}
			})

			ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
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

		ginkgo.AfterEach(func() {
			_ = shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.AfterEach(func() {
		shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
	})
})
