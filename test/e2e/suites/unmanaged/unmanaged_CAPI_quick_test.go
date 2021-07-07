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

package unmanaged

import (
	"context"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/config"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
)

var _ = ginkgo.Context("[unmanaged] [Cluster API Framework] [smoke] [PR-Blocking]", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
	)

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = shared.SetupSpecNamespace(ctx, "capi-quick-start", e2eCtx)
	})

	ginkgo.Describe("Running the quick-start spec", func() {
		// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
		requiredResources := &shared.TestResource{EC2: 2, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
		ginkgo.BeforeEach(func() {
			requiredResources.WriteRequestedResources(e2eCtx, "capi-quick-start-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})
		capi_e2e.QuickStartSpec(context.TODO(), func() capi_e2e.QuickStartSpecInput {
			return capi_e2e.QuickStartSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
		})
	})
	ginkgo.AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
	})
})
