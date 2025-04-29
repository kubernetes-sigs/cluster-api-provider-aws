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

package unmanaged

import (
	"context"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
)

// setupNamespace initializes the namespace for the test.
func setupNamespace(ctx context.Context, e2eCtx *shared.E2EContext) *corev1.Namespace {
	Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
	return shared.SetupSpecNamespace(ctx, "capa-dedicate-host", e2eCtx)
}

// setupRequiredResources allocates the required resources for the test.
func setupRequiredResources(e2eCtx *shared.E2EContext) *shared.TestResource {
	requiredResources := &shared.TestResource{
		EC2Normal:        2 * e2eCtx.Settings.InstanceVCPU,
		IGW:              1,
		NGW:              1,
		VPC:              1,
		ClassicLB:        1,
		EIP:              3,
		EventBridgeRules: 50,
	}
	requiredResources.WriteRequestedResources(e2eCtx, "capa-dedicated-hosts-test")

	Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
	return requiredResources
}

// releaseResources releases the resources allocated for the test.
func releaseResources(requiredResources *shared.TestResource, e2eCtx *shared.E2EContext) {
	shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
}

// runQuickStartSpec executes the QuickStartSpec test.
func runQuickStartSpec(e2eCtx *shared.E2EContext) {
	capi_e2e.QuickStartSpec(context.TODO(), func() capi_e2e.QuickStartSpecInput {
		return capi_e2e.QuickStartSpecInput{
			E2EConfig:             e2eCtx.E2EConfig,
			ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
			BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
			ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
			SkipCleanup:           e2eCtx.Settings.SkipCleanup,
		}
	})
}

// cleanupNamespace cleans up the namespace and dumps resources.
func cleanupNamespace(ctx context.Context, namespace *corev1.Namespace, e2eCtx *shared.E2EContext) {
	shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
}

var _ = ginkgo.Context("[unmanaged] [dedicated-host]", func() {
	var (
		namespace         *corev1.Namespace
		ctx               context.Context
		requiredResources *shared.TestResource
		dedicatedHostID   string
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
		namespace = setupNamespace(ctx, e2eCtx)
		dedicatedHostID, _ = shared.GetDedicatedHost(e2eCtx)
	})

	ginkgo.Describe("Running the dedicated-hosts spec", func() {
		ginkgo.BeforeEach(func() {
			requiredResources = setupRequiredResources(e2eCtx)
			// e2eCtx.Settings.DedicatedHostID = dedicatedHostID
		})

		ginkgo.It("should run the QuickStartSpec", func() {
			runQuickStartSpec(e2eCtx)
		})

		ginkgo.AfterEach(func() {
			shared.DeleteDedicatedHost(e2eCtx, dedicatedHostID)
			releaseResources(requiredResources, e2eCtx)
		})
	})

	ginkgo.AfterEach(func() {
		cleanupNamespace(ctx, namespace, e2eCtx)
	})
})
