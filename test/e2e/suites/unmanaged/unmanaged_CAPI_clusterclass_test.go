//go:build e2e
// +build e2e

/*
Copyright 2022 The Kubernetes Authors.

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
)

var _ = ginkgo.Context("[unmanaged] [Cluster API Framework] [ClusterClass]", func() {
	var (
		ctx               = context.TODO()
		requiredResources *shared.TestResource
	)

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
	})

	ginkgo.Describe("Self Hosted Spec [ClusterClass]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 2 * e2eCtx.Settings.InstanceVCPU, IGW: 1, NGW: 1, VPC: 1, ClassicLB: 1, EIP: 1, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-clusterctl-self-hosted-test-clusterclass")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.SelfHostedSpec(ctx, func() capi_e2e.SelfHostedSpecInput {
			return capi_e2e.SelfHostedSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
				Flavor:                shared.SelfHostedClusterClassFlavor,
			}
		})

		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("Cluster Upgrade Spec - HA control plane with workers [K8s-Upgrade] [ClusterClass]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 5 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-cluster-upgrade-clusterclass-test")
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
				Flavor:                   pointer.String(shared.TopologyFlavor),
			}
		})

		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})

	ginkgo.Describe("ClusterClass Changes Spec - SSA immutability checks [ClusterClass]", func() {
		ginkgo.BeforeEach(func() {
			// As the resources cannot be defined by the It() clause in CAPI tests, using the largest values required for all It() tests in this CAPI test.
			requiredResources = &shared.TestResource{EC2Normal: 5 * e2eCtx.Settings.InstanceVCPU, IGW: 2, NGW: 2, VPC: 2, ClassicLB: 2, EIP: 2, EventBridgeRules: 50}
			requiredResources.WriteRequestedResources(e2eCtx, "capi-cluster-ssa-clusterclass-test")
			Expect(shared.AcquireResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))).To(Succeed())
		})

		capi_e2e.ClusterClassChangesSpec(ctx, func() capi_e2e.ClusterClassChangesSpecInput {
			return capi_e2e.ClusterClassChangesSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ClusterctlConfigPath:  e2eCtx.Environment.ClusterctlConfigPath,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ArtifactFolder:        e2eCtx.Settings.ArtifactFolder,
				SkipCleanup:           e2eCtx.Settings.SkipCleanup,
				Flavor:                shared.TopologyFlavor,
				// ModifyControlPlaneFields are the ControlPlane fields which will be set on the
				// ControlPlaneTemplate of the ClusterClass after the initial Cluster creation.
				// The test verifies that these fields are rolled out to the ControlPlane.
				ModifyControlPlaneFields: map[string]interface{}{
					"spec.machineTemplate.nodeDrainTimeout": "10s",
				},
				// ModifyMachineDeploymentBootstrapConfigTemplateFields are the fields which will be set on the
				// BootstrapConfigTemplate of all MachineDeploymentClasses of the ClusterClass after the initial Cluster creation.
				// The test verifies that these fields are rolled out to the MachineDeployments.
				ModifyMachineDeploymentBootstrapConfigTemplateFields: map[string]interface{}{
					"spec.template.spec.verbosity": int64(4),
				},

				// ModifyMachineDeploymentInfrastructureMachineTemplateFields are the fields which will be set on the
				// InfrastructureMachineTemplate of all MachineDeploymentClasses of the ClusterClass after the initial Cluster creation.
				// The test verifies that these fields are rolled out to the MachineDeployments.
				ModifyMachineDeploymentInfrastructureMachineTemplateFields: map[string]interface{}{
					"spec.template.spec.additionalTags": map[string]interface{}{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			}
		})

		ginkgo.AfterEach(func() {
			shared.ReleaseResources(requiredResources, ginkgo.GinkgoParallelProcess(), flock.New(shared.ResourceQuotaFilePath))
		})
	})
})
