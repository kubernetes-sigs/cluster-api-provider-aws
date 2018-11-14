// +build integration

// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"bytes"
	"flag"
	"io"
	"os/exec"
	"strings"
	"testing"
)

var (
	kindBinary      = flag.String("kindBinary", "kind", "path to the kind binary")
	kubectlBinary   = flag.String("kubectlBinary", "kubectl", "path to the kubectl binary")
	awsProviderYAML = flag.String("awsProviderYAML", "", "path to the Kubernetes YAML for the aws provider")
	clusterAPIYAML  = flag.String("clusterAPIYAML", "", "path to the Kubernetes YAML for the cluster API")
)

func TestYAMLAppliesToRealCluster(t *testing.T) {
	kubeConfig := setupKIND(t)
	defer teardownKIND(t)
	applyYAML(t, kubeConfig)
}

func setupKIND(t *testing.T) string {
	run(t, exec.Command(*kindBinary, "create", "cluster"))
	path, err := exec.Command(*kindBinary, "get", "kubeconfig-path").Output()
	if err != nil {
		t.Fatalf("couldn't get kubeconfig path: %v", err)
	}
	return strings.TrimSpace(string(path))
}

func teardownKIND(t *testing.T) {
	run(t, exec.Command(*kindBinary, "delete", "cluster"))
}

func applyYAML(t *testing.T, kubeConfig string) {
	run(t, exec.Command(
		*kubectlBinary,
		"create",
		"--kubeconfig="+kubeConfig,
		"-f", *awsProviderYAML,
		"-f", *clusterAPIYAML,
	))
}

func run(t *testing.T, cmd *exec.Cmd) {
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("couldn't get pipe: %v", err)
	}
	go captureOutput(t, errPipe, "stderr")
	if cmd.Stdout == nil {
		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatalf("couldn't get pipe: %v", err)
		}
		go captureOutput(t, outPipe, "stdout")
	}

	err = cmd.Run()

	if err != nil {
		t.Fatalf("Couldn't run command %s %s: %v", cmd.Path, cmd.Args, err)
	}
}

func captureOutput(t *testing.T, pipe io.ReadCloser, label string) {
	buffer := &bytes.Buffer{}
	defer pipe.Close()
	for {
		n, _ := buffer.ReadFrom(pipe)
		if n == 0 {
			return
		}
		t.Logf("[%s] %s\n", label, buffer.String())
	}
}
