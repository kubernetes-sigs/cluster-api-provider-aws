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

package kubetest

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api/test/framework"
)

const (
	standardImageRepo = "gcr.io/google-containers/conformance"
	ciImageRepo       = "gcr.io/kubernetes-ci-images/conformance"
)

// Run executes kube-test given an artifact directory, and sets settings
// required for kubetest to work with Cluster API. JUnit files are
// also gathered for inclusion in Prow.
func Run(c *Configuration, proxy framework.ClusterProxy, numNodes int64) {
	outputDir := path.Join(c.reportDir(), "e2e-output")
	Expect(os.MkdirAll(outputDir, 0750)).To(Succeed(), "Invalid test suite argument. Can't create output directory %q", outputDir)
	ginkgoVars := map[string]string{
		"nodes":             strconv.Itoa(c.GinkgoNodes),
		"slowSpecThreshold": strconv.Itoa(c.GinkgoSlowSpecThreshold),
	}
	e2eVars := map[string]string{
		"kubeconfig":           "/tmp/kubeconfig",
		"provider":             "skeleton",
		"report-dir":           "/output",
		"e2e-output-dir":       "/output/e2e-output",
		"dump-logs-on-failure": "false",
		"report-prefix":        "kubetest.",
		"num-nodes":            strconv.FormatInt(numNodes, 10),
		"viper-config":         "/tmp/viper-config.yaml",
	}
	ginkgoArgs := buildArgs(ginkgoVars, "-")
	e2eArgs := buildArgs(e2eVars, "--")
	conformanceImage := c.conformanceImage()
	kubeConfigVolumeMount := volumeArg(proxy.GetKubeconfigPath(), "/tmp/kubeconfig")
	outputVolumeMount := volumeArg(c.reportDir(), "/output")
	viperVolumeMount := volumeArg(c.ConfigFilePath, "/tmp/viper-config.yaml")
	user, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	userArg := user.Uid + ":" + user.Gid
	e2eCmd := exec.Command("docker", "run", "--user", userArg, kubeConfigVolumeMount, outputVolumeMount, viperVolumeMount, "-t", conformanceImage)
	e2eCmd.Args = append(e2eCmd.Args, "/usr/local/bin/ginkgo")
	e2eCmd.Args = append(e2eCmd.Args, ginkgoArgs...)
	e2eCmd.Args = append(e2eCmd.Args, "/usr/local/bin/e2e.test")
	e2eCmd.Args = append(e2eCmd.Args, "--")
	e2eCmd.Args = append(e2eCmd.Args, e2eArgs...)
	e2eCmd = completeCommand(e2eCmd, "Running e2e test", false)
	e2eErr := e2eCmd.Run()
	// Rename files for consumption by prow
	Expect(GatherReports(c)).To(Succeed())
	Expect(e2eErr).NotTo(HaveOccurred())
}

func isSELinuxEnforcing() bool {
	dat, err := ioutil.ReadFile("/sys/fs/selinux/enforce")
	if err != nil {
		return false
	}
	return string(dat) == "1"
}

func volumeArg(src, dest string) string {
	volumeArg := "-v" + src + ":" + dest
	if isSELinuxEnforcing() {
		return volumeArg + ":z"
	}
	return volumeArg
}

func (c *Configuration) reportDir() string {
	return path.Join(c.ArtifactsDirectory, "kubetest")
}

func GatherReports(c *Configuration) error {
	if err := os.MkdirAll(c.reportDir(), 0750); err != nil {
		return err
	}
	return filepath.Walk(c.reportDir(), func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && p != c.reportDir() {
			return filepath.SkipDir
		}
		if filepath.Ext(p) != ".xml" {
			return nil
		}
		base := filepath.Base(p)
		if strings.HasPrefix(base, "junit") {
			newName := strings.ReplaceAll(base, "_", ".")
			if err := os.Rename(p, path.Join(filepath.Dir(p), "..", newName)); err != nil {
				return err
			}
		}
		return nil
	})

}

func (c *Configuration) conformanceImage() string {
	k8sVersion := strings.ReplaceAll(c.E2EConfig.GetVariable("KUBERNETES_VERSION"), "+", "_")
	if c.UseCIArtifacts {
		return ciImageRepo + ":" + k8sVersion
	}
	return standardImageRepo + ":" + k8sVersion
}

// buildArgs converts a string map to the format --key=value
func buildArgs(kv map[string]string, flagMarker string) []string {
	args := make([]string, len(kv))
	i := 0
	for k, v := range kv {
		args[i] = flagMarker + k + "=" + v
		i++
	}
	return args
}

// completeCommand prints a command before running it. Acts as a helper function.
// privateArgs when true will not print arguments.
func completeCommand(cmd *exec.Cmd, desc string, privateArgs bool) *exec.Cmd {
	cmd.Stderr = GinkgoWriter
	cmd.Stdout = GinkgoWriter
	if privateArgs {
		Byf("%s: dir=%s, command=%s", desc, cmd.Dir, cmd)
	} else {
		Byf("%s: dir=%s, command=%s", desc, cmd.Dir, cmd.String())
	}
	return cmd
}
