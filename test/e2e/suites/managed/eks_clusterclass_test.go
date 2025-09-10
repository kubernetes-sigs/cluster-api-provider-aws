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

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

var _ = ginkgo.Describe("[managed] [general] EKS clusterclass tests", func() {
	const specName = "cluster"
	var (
		ctx         context.Context
		clusterName string
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

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, "", false, e2eCtx.AWSSession)
	})

	capi_e2e.QuickStartSpec(context.TODO(), func() capi_e2e.QuickStartSpecInput {
		return capi_e2e.QuickStartSpecInput{
			E2EConfig:             e2eCtx.E2EConfig,
			ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
			BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
			ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
			SkipCleanup:           e2eCtx.Settings.SkipCleanup,
			Flavor:                ptr.To(EKSClusterClassFlavor),
			ClusterName:           ptr.To(clusterName),
			WorkerMachineCount:    ptr.To(int64(3)),
			ControlPlaneWaiters: clusterctl.ControlPlaneWaiters{
				WaitForControlPlaneInitialized:   WaitForEKSControlPlaneInitialized,
				WaitForControlPlaneMachinesReady: WaitForEKSControlPlaneMachinesReady,
			},
		}
	})
})
