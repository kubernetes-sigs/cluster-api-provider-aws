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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

// EKS cluster upgrade tests
var _ = Describe("EKS Cluster upgrade test", func() {
	const (
		initialVersion   = "v1.17.0"
		upgradeToversion = "v1.18.0"
	)
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "eks-upgrade"
		clusterName string
	)

	shared.ConditionalIt(runUpgradeTests, "[managed] [upgrade] should create a cluster and upgrade the kubernetes version", func() {
		By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("cluster-%s", util.RandomString(6))

		By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSControlPlaneOnlyFlavor, //TODO (richardcase) - change in the future when upgrades to machinepools work
				ControlPlaneMachineCount: 1,                         //NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       1,
				CNIManifestPath:          e2eCtx.E2EConfig.GetVariable(shared.CNIPath),
				KubernetesVersion:        initialVersion,
			}
		})

		// TODO: should cluster be returned from the ManagedClusterSpec as a convenience
		shared.Byf("getting cluster with name %s", clusterName)
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		//TODO (richardcase) - uncomment when we use machine pools again
		// shared.Byf("Waiting for the machine pool to be running")
		// mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
		// 	Lister:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
		// 	Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
		// 	Cluster: cluster,
		// }, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		// Expect(len(mp)).To(Equal(1))

		shared.Byf("should upgrade control plane to version %s", upgradeToversion)
		UpgradeControlPlaneVersionSpec(ctx, func() UpgradeControlPlaneVersionSpecInput {
			return UpgradeControlPlaneVersionSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ClusterName:           clusterName,
				Namespace:             namespace,
				UpgradeVersion:        upgradeToversion,
			}
		})

		//TODO (richardcase): add test for the node group upgrade

		shared.Byf("should delete cluster %s", clusterName)
		DeleteClusterSpec(ctx, func() DeleteClusterSpecInput {
			return DeleteClusterSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ClusterName:           clusterName,
				Namespace:             namespace,
			}
		})

	})

})
