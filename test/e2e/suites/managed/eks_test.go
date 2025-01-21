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

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// General EKS e2e test.
var _ = ginkgo.Describe("[managed] [general] EKS cluster tests", func() {
	var (
		namespace    *corev1.Namespace
		ctx          context.Context
		specName     = "cluster"
		clusterName  string
		cniAddonName = "vpc-cni"
	)

	shared.ConditionalIt(runGeneralTests, "should create a cluster and add nodes", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.CNIAddonVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		eksClusterName := getEKSClusterName(namespace.Name, clusterName)

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ekscontrolplanev1.DefaultEKSControlPlaneRole, eksClusterName, false, e2eCtx.BootstrapUserAWSSession)

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
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
			}
		})

		ginkgo.By("should set environment variables on the aws-node daemonset")
		CheckAwsNodeEnvVarsSet(ctx, func() UpdateAwsNodeVersionSpecInput {
			return UpdateAwsNodeVersionSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
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
				AddonVersion:          e2eCtx.E2EConfig.GetVariable(shared.CNIAddonVersion),
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
		MachinePoolSpec(ctx, func() MachinePoolSpecInput {
			return MachinePoolSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				IncludeScaling:        true,
				Cleanup:               true,
				ManagedMachinePool:    true,
				Flavor:                EKSManagedMachinePoolOnlyFlavor,
			}
		})

		ginkgo.By("should create a machine pool and scale")
		MachinePoolSpec(ctx, func() MachinePoolSpecInput {
			return MachinePoolSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				IncludeScaling:        true,
				Cleanup:               true,
				ManagedMachinePool:    false,
				Flavor:                EKSMachinePoolOnlyFlavor,
				UsesLaunchTemplate:    false,
			}
		})

		ginkgo.By("should create a managed node pool with launch template and scale")
		MachinePoolSpec(ctx, func() MachinePoolSpecInput {
			return MachinePoolSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				IncludeScaling:        true,
				Cleanup:               true,
				ManagedMachinePool:    true,
				Flavor:                EKSManagedMachinePoolWithLaunchTemplateOnlyFlavor,
				UsesLaunchTemplate:    true,
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
			Client:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)
	})
})
