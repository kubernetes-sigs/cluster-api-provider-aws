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
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

var _ = ginkgo.Context("[unmanaged] [Cluster API Framework]", func() {

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
	})

	// DEPRECATED. Should be replaced with the conformance upgrade spec
	ginkgo.Describe("KCP Upgrade Spec - Single Control Plane Cluster", func() {
		// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
		requiredResources := &shared.TestResource{EC2: 2, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
		ginkgo.BeforeEach(func() {
			requiredResources.WriteRequestedResources(e2eCtx, "capi-kcp-single-cp-upgrade-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})
		capi_e2e.KCPUpgradeSpec(context.TODO(), func() capi_e2e.KCPUpgradeSpecInput {
			return capi_e2e.KCPUpgradeSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				Flavor:                   clusterctl.DefaultFlavor,
				ControlPlaneMachineCount: 1,
				ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:              e2eCtx.Settings.SkipCleanup,
			}
		})
	})

	// DEPRECATED. Should be replaced with the conformance upgrade spec
	ginkgo.Describe("KCP Upgrade Spec - HA Control Plane Cluster using Scale-In", func() {
		// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
		requiredResources := &shared.TestResource{EC2: 4, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
		ginkgo.BeforeEach(func() {
			requiredResources.WriteRequestedResources(e2eCtx, "capi-kcp-scale-in-upgrade-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})
		capi_e2e.KCPUpgradeSpec(context.TODO(), func() capi_e2e.KCPUpgradeSpecInput {
			return capi_e2e.KCPUpgradeSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				Flavor:                   shared.KCPScaleInFlavor,
				ControlPlaneMachineCount: 3,
				ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:              e2eCtx.Settings.SkipCleanup,
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Machine Remediation Spec", func() {
		// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
		requiredResources := &shared.TestResource{EC2: 4, IGW: 1, NGW: 3, VPC: 1, ClassicLB: 1, EIP: 3}
		ginkgo.BeforeEach(func() {
			requiredResources.WriteRequestedResources(e2eCtx, "capi-remediation-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.MachineRemediationSpec(context.TODO(), func() capi_e2e.MachineRemediationSpecInput {
			return capi_e2e.MachineRemediationSpecInput{
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

	ginkgo.Describe("Machine Pool Spec", func() {
		// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
		requiredResources := &shared.TestResource{EC2: 4, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3}
		ginkgo.BeforeEach(func() {
			requiredResources.WriteRequestedResources(e2eCtx, "capi-machinepool-test")
			Expect(shared.AcquireResources(requiredResources, config.GinkgoConfig.ParallelNode, flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.MachinePoolSpec(context.TODO(), func() capi_e2e.MachinePoolInput {
			return capi_e2e.MachinePoolInput{
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
})
