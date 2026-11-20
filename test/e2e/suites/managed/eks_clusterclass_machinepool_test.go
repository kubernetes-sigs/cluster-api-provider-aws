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

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Describe("[managed] [general] EKS clusterclass self-managed machine pool tests", func() {
	const specName = "cluster"
	var (
		ctx         context.Context
		clusterName string
		namespace   *corev1.Namespace
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()

		if !runGeneralTests() {
			ginkgo.Skip("skipping due to unmet condition")
		}

		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "E2EConfig can't be nil")
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.CNIAddonVersion))

		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, "", false, e2eCtx.AWSSession)
	})

	ginkgo.AfterEach(func() {
		shared.DumpSpecResourcesAndCleanup(ctx, specName, namespace, e2eCtx)
	})

	// QuickStartSpec does not wait for machine pools, so drive the apply directly and
	// wait for the pool with the EKS worker-nodes interval.
	ginkgo.It("should create a workload cluster with a self-managed machine pool", func() {
		configCluster := defaultConfigCluster(clusterName, namespace.Name)
		configCluster.Flavor = EKSClusterClassMachinePoolFlavor
		configCluster.WorkerMachineCount = ptr.To[int64](1)

		result := &clusterctl.ApplyClusterTemplateAndWaitResult{}
		clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
			ClusterProxy:  e2eCtx.Environment.BootstrapClusterProxy,
			ConfigCluster: configCluster,
			ControlPlaneWaiters: clusterctl.ControlPlaneWaiters{
				WaitForControlPlaneInitialized:   WaitForEKSControlPlaneInitialized,
				WaitForControlPlaneMachinesReady: WaitForEKSControlPlaneMachinesReady,
			},
			WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
			WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
			WaitForMachinePools:          e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
		}, result)
	})
})
