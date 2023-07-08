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

package unmanaged

import (
	"context"

	"github.com/gofrs/flock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	capi_e2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

var _ = ginkgo.Context("[unmanaged] [Cluster API Framework]", func() {
	var (
		ctx               = context.TODO()
		requiredResources *shared.TestResource
	)

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
	})

	ginkgo.Describe("MachineDeployment Remediation Spec", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 4 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 3, VPC: 1, ClassicLB: 1, EIP: 3, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-md-remediation-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.MachineDeploymentRemediationSpec(ctx, func() capi_e2e.MachineDeploymentRemediationSpecInput {
			return capi_e2e.MachineDeploymentRemediationSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Machine Pool Spec", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 4 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 3, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-machinepool-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.MachinePoolSpec(ctx, func() capi_e2e.MachinePoolInput {
			return capi_e2e.MachinePoolInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Self Hosted Spec", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-clusterctl-self-hosted-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.SelfHostedSpec(ctx, func() capi_e2e.SelfHostedSpecInput {
			return capi_e2e.SelfHostedSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
				Flavor:                "remote-management-cluster",
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Clusterctl Upgrade Spec [from latest v1beta1 release to v1beta2]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 5 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-clusterctl-upgrade-test-v1beta1-to-v1beta2")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.ClusterctlUpgradeSpec(ctx, func() capi_e2e.ClusterctlUpgradeSpecInput {
			return capi_e2e.ClusterctlUpgradeSpecInput{
				E2EConfig:                 e2eCtx.E2EConfig,
				ClusterctlConfigPath:      e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:     e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:            e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:               e2eCtx.Settings.SkipCleanup,
				MgmtFlavor:                "remote-management-cluster",
				InitWithBinary:            e2eCtx.E2EConfig.GetVariable("INIT_WITH_BINARY_V1BETA1"),
				InitWithProvidersContract: "v1beta1",
			}
		})
		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Cluster Upgrade Spec - Single control plane with workers [K8s-Upgrade]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 5 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-worker-upgrade-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.ClusterUpgradeConformanceSpec(ctx, func() capi_e2e.ClusterUpgradeConformanceSpecInput {
			return capi_e2e.ClusterUpgradeConformanceSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:              e2eCtx.Settings.SkipCleanup,
				SkipConformanceTests:     true,
				ControlPlaneMachineCount: pointer.Int64(1),
			}
		})

		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Cluster Upgrade Spec - HA control plane with scale in rollout [K8s-Upgrade]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 10 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-ha-cluster-upgrade-scale-in-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.ClusterUpgradeConformanceSpec(ctx, func() capi_e2e.ClusterUpgradeConformanceSpecInput {
			return capi_e2e.ClusterUpgradeConformanceSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:              e2eCtx.Settings.SkipCleanup,
				SkipConformanceTests:     true,
				ControlPlaneMachineCount: pointer.Int64(3),
				WorkerMachineCount:       pointer.Int64(0),
				Flavor:                   pointer.String(shared.KCPScaleInFlavor),
			}
		})
	})

	ginkgo.Describe("Cluster Upgrade Spec - HA Control Plane Cluster [K8s-Upgrade]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 10 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-ha-cluster-upgrade-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.ClusterUpgradeConformanceSpec(ctx, func() capi_e2e.ClusterUpgradeConformanceSpecInput {
			return capi_e2e.ClusterUpgradeConformanceSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:           e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:              e2eCtx.Settings.SkipCleanup,
				SkipConformanceTests:     true,
				ControlPlaneMachineCount: pointer.Int64(3),
				WorkerMachineCount:       pointer.Int64(0),
				Flavor:                   pointer.String(clusterctl.DefaultFlavor),
			}
		})

		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})
})
