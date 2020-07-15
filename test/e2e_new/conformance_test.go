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

package e2e_new

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e_new/kubetest"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util"
)

const (
	AMIPrefix             = "capa-ami-ubuntu-18.04-"
	DefaultImageLookupOrg = "258751437250"
)

var _ = Describe("conformance tests", func() {
	var (
		namespace *corev1.Namespace
		ctx       context.Context
	)

	BeforeEach(func() {
		Expect(bootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(kubetestConfig).ToNot(BeNil(), "Invalid argument. ConformanceConfiguration can't be nil")
		ctx = context.TODO()
		// Setup a Namespace where to host objects for this spec and create a watcher for the namespace events.
		namespace = setupSpecNamespace(ctx, "conformance-tests", bootstrapClusterProxy, artifactFolder)
	})
	specName := "conformance"
	Measure("conformance", func(b Benchmarker) {

		name := fmt.Sprintf("cluster-%s", util.RandomString(6))
		flavor := clusterctl.DefaultFlavor
		if kubetestConfig.UseCIArtifacts {
			flavor = "conformance-ci-artifacts"
		}
		workerMachineCount, err := strconv.ParseInt(e2eConfig.GetVariable("CONFORMANCE_WORKER_MACHINE_COUNT"), 10, 64)
		Expect(err).NotTo(HaveOccurred())
		controlPlaneMachineCount, err := strconv.ParseInt(e2eConfig.GetVariable("CONFORMANCE_CONTROL_PLANE_MACHINE_COUNT"), 10, 64)
		Expect(err).NotTo(HaveOccurred())

		runtime := b.Time("cluster creation", func() {
			_, _, _ = clusterctl.ApplyClusterTemplateAndWait(ctx, clusterctl.ApplyClusterTemplateAndWaitInput{
				ClusterProxy: bootstrapClusterProxy,
				ConfigCluster: clusterctl.ConfigClusterInput{
					LogFolder:                filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
					ClusterctlConfigPath:     clusterctlConfigPath,
					KubeconfigPath:           bootstrapClusterProxy.GetKubeconfigPath(),
					InfrastructureProvider:   clusterctl.DefaultInfrastructureProvider,
					Flavor:                   flavor,
					Namespace:                namespace.Name,
					ClusterName:              name,
					KubernetesVersion:        e2eConfig.GetVariable(KubernetesVersion),
					ControlPlaneMachineCount: pointer.Int64Ptr(controlPlaneMachineCount),
					WorkerMachineCount:       pointer.Int64Ptr(workerMachineCount),
				},
				CNIManifestPath:              e2eConfig.GetVariable(CNIPath),
				WaitForClusterIntervals:      e2eConfig.GetIntervals(specName, "wait-cluster"),
				WaitForControlPlaneIntervals: e2eConfig.GetIntervals(specName, "wait-control-plane"),
				WaitForMachineDeployments:    e2eConfig.GetIntervals(specName, "wait-worker-nodes"),
			})
		})
		b.RecordValue("cluster creation", runtime.Seconds())
		workloadProxy := bootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, name)
		runtime = b.Time("conformance suite", func() {
			kubetest.Run(kubetestConfig, workloadProxy, workerMachineCount)
		})
		b.RecordValue("conformance suite run time", runtime.Seconds())
	}, 1)

	AfterEach(func() {
		// Dumps all the resources in the spec namespace, then cleanups the cluster object and the spec namespace itself.
		dumpSpecResourcesAndCleanup(ctx, "", bootstrapClusterProxy, artifactFolder, namespace, e2eConfig.GetIntervals, skipCleanup)
	})

})

func prepConformanceConfig(artifactsDir string, e2eConfig *clusterctl.E2EConfig) error {
	for i, p := range e2eConfig.Providers {
		if p.Name != "aws" {
			continue
		}
		e2eConfig.Providers[i].Files = append(p.Files, clusterctl.Files{
			SourcePath: path.Join(artifactFolder, "/templates/cluster-template-conformance-ci-artifacts.yaml"),
			TargetName: "cluster-template-conformance-ci-artifacts.yaml",
		},
		)
	}
	templateDir := path.Join(artifactsDir, "templates")
	overlayDir := path.Join(artifactsDir, "overlay")
	srcDir := path.Join("data", "kubetest", "kustomization")
	originalTemplate := path.Join("..", "..", "templates", "cluster-template.yaml")

	if err := copy(path.Join(srcDir, "kustomization.yaml"), path.Join(overlayDir, "kustomization.yaml")); err != nil {
		return err
	}
	if err := copy(path.Join(srcDir, "kustomizeversions.yaml"), path.Join(overlayDir, "kustomizeversions.yaml")); err != nil {
		return err
	}
	if err := copy(originalTemplate, path.Join(overlayDir, "cluster-template.yaml")); err != nil {
		return err
	}
	cmd := exec.Command("kustomize", "build", overlayDir)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	os.MkdirAll(templateDir, 0o750)
	if err := ioutil.WriteFile(path.Join(templateDir, "cluster-template-conformance-ci-artifacts.yaml"), data, 0o640); err != nil {
		return err
	}
	return nil
}

func copy(src, dest string) error {
	os.MkdirAll(path.Dir(dest), 0o750)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	return nil
}
