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

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

// General EKS e2e test
var _ = ginkgo.Describe("[managed] [general] EKS cluster tests", func() {
	var (
		namespace       *corev1.Namespace
		ctx             context.Context
		specName        = "eks-nodes"
		clusterName     string
		cniAddonName    = "vpc-cni"
		cniAddonVersion = "v1.6.3-eksbuild.1"
	)

	shared.ConditionalIt(runGeneralTests, "should create a cluster and add nodes", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("cluster-%s", util.RandomString(6))

		ginkgo.By("setting up AWS static credentials")
		shared.SetupStaticCredentials(ctx, namespace, e2eCtx)

		ginkgo.By("default iam role should exist")
		verifyRoleExistsAndOwned(controlplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSControlPlaneOnlyWithAddonFlavor,
				ControlPlaneMachineCount: 1, //NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
				CNIManifestPath:          e2eCtx.E2EConfig.GetVariable(shared.CNIPath),
			}
		})

		ginkgo.By("should have the VPC CNI installed")
		CheckAddonExistsSpec(ctx, func() CheckAddonExistsSpecInput {
			return CheckAddonExistsSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				AddonName:             cniAddonName,
				AddonVersion:          cniAddonVersion,
			}
		})

		ginkgo.By("should create a MachineDeployment")
		MachineDeploymentSpec(ctx, func() MachineDeploymentSpecInput {
			return MachineDeploymentSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				Replicas:              1,
				Cleanup:               true,
			}
		})

		ginkgo.By("should create a managed node pool and scale")
		ManagedMachinePoolSpec(ctx, func() ManagedMachinePoolSpecInput {
			return ManagedMachinePoolSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				IncludeScaling:        true,
				Cleanup:               true,
			}
		})

		shared.Byf("getting cluster with name %s", clusterName)
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
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)

		ginkgo.By("Deleting AWS static credentials")
		shared.CleanupStaticCredentials(ctx, namespace, e2eCtx)
	})
})
