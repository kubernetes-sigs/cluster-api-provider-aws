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

package conformance

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/test/framework/kubernetesversions"
	"sigs.k8s.io/cluster-api/test/framework/kubetest"
	"sigs.k8s.io/cluster-api/util"
)

// TODO @randomvariable: Replace with CAPI e2e framework ClusterUpgradeConformanceSpec
var _ = ginkgo.Describe("[unmanaged] [conformance] tests", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
		specName  = "conformance"
		result    *clusterctl.ApplyClusterTemplateAndWaitResult
	)

	ginkgo.BeforeEach(func() {
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		result = new(clusterctl.ApplyClusterTemplateAndWaitResult)
	})
	ginkgo.Measure(specName, func(b ginkgo.Benchmarker) {

		name := fmt.Sprintf("cluster-%s", util.RandomString(6))
		shared.SetEnvVar("USE_CI_ARTIFACTS", "true", false)
		kubernetesVersion := e2eCtx.E2EConfig.GetVariable(shared.KubernetesVersion)
		flavor := clusterctl.DefaultFlavor
		if e2eCtx.Settings.UseCIArtifacts {
			flavor = "conformance-ci-artifacts"
			var err error
			kubernetesVersion, err = kubernetesversions.LatestCIRelease()
			Expect(err).NotTo(HaveOccurred())
		}
		workerMachineCount, err := strconv.ParseInt(e2eCtx.E2EConfig.GetVariable("CONFORMANCE_WORKER_MACHINE_COUNT"), 10, 64)
		Expect(err).NotTo(HaveOccurred())
		controlPlaneMachineCount, err := strconv.ParseInt(e2eCtx.E2EConfig.GetVariable("CONFORMANCE_CONTROL_PLANE_MACHINE_COUNT"), 10, 64)
		Expect(err).NotTo(HaveOccurred())

		runtime := b.Time("cluster creation", func() {
			clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ConfigCluster: clusterctl.ConfigClusterInput{
					LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
					ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
					KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
					InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
					Flavor:                   flavor,
					Namespace:                namespace.Name,
					ClusterName:              name,
					KubernetesVersion:        kubernetesVersion,
					ControlPlaneMachineCount: pointer.Int64Ptr(controlPlaneMachineCount),
					WorkerMachineCount:       pointer.Int64Ptr(workerMachineCount),
				},
				WaitForClusterIntervals:      e2eCtx.E2EConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eCtx.E2EConfig.GetIntervals(specName, "wait-control-plane"),
				WaitForMachineDeployments:    e2eCtx.E2EConfig.GetIntervals(specName, "wait-worker-nodes"),
			}, result)
		})
		b.RecordValue("cluster creation", runtime.Seconds())
		workloadProxy := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, name)
		runtime = b.Time("conformance suite", func() {
			kubetest.Run(ctx,
				kubetest.RunInput{
					ClusterProxy:   workloadProxy,
					NumberOfNodes:  int(workerMachineCount),
					ConfigFilePath: e2eCtx.Settings.KubetestConfigFilePath,
				},
			)
		})
		b.RecordValue("conformance suite run time", runtime.Seconds())
	}, 1)

	ginkgo.AfterEach(func() {
		shared.SetEnvVar("USE_CI_ARTIFACTS", "false", false)
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		shared.DumpSpecResourcesAndCleanup(ctx, "", namespace, e2eCtx)
	})

})
