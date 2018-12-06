/*
Copyright 2018 The Kubernetes Authors.

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

package kind

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kindBinary      = flag.String("kindBinary", "kind", "path to the kind binary")
	kubectlBinary   = flag.String("kubectlBinary", "kubectl", "path to the kubectl binary")
	awsProviderYAML = flag.String("awsProviderYAML", "", "path to the Kubernetes YAML for the aws provider")
	clusterAPIYAML  = flag.String("clusterAPIYAML", "", "path to the Kubernetes YAML for the cluster API")
	managerImageTar = flag.String("managerImageTar", "", "a script to load the manager Docker image into Docker")
)

const kindContainerName = "kind-1-control-plane"

// Cluster represents the running state of a KIND cluster.
// An empty struct is enough to call Setup() on.
type Cluster struct {
	tmpDir   string
	kubepath string
}

// Setup creates a kind cluster and returns a path to the kubeconfig
func (c *Cluster) Setup() {
	var err error
	c.tmpDir, err = ioutil.TempDir("", "kind-home")
	gomega.Expect(err).To(gomega.BeNil())

	c.run(exec.Command(*kindBinary, "create", "cluster"))
	path := c.runWithOutput(exec.Command(*kindBinary, "get", "kubeconfig-path"))
	c.kubepath = strings.TrimSpace(string(path))
	fmt.Fprintf(ginkgo.GinkgoWriter, "kubeconfig path: %q\n", c.kubepath)

	if *managerImageTar != "" {
		c.loadImage()
	}

	c.applyYAML()
}

func (c *Cluster) loadImage() {
	// TODO(EKF): once kind supports loading images directly, remove this hack
	file, err := os.Open(*managerImageTar)
	gomega.Expect(err).To(gomega.BeNil())

	// Pipe the tar file into the kind container then docker-load it
	cmd := exec.Command("docker", "exec", "--interactive", kindContainerName, "docker", "load")
	cmd.Stdin = file
	c.run(cmd)
}

// Teardown attempts to delete the KIND cluster
func (c *Cluster) Teardown() {
	c.run(exec.Command(*kindBinary, "delete", "cluster"))
	os.RemoveAll(c.tmpDir)
}

// applyYAML takes the provided awsProviderYAML and clusterAPIYAML and applies them to a cluster given by the kubeconfig path kubeConfig.
func (c *Cluster) applyYAML() {
	c.run(exec.Command(
		*kubectlBinary,
		"create",
		"--kubeconfig="+c.kubepath,
		"-f", *awsProviderYAML,
		"-f", *clusterAPIYAML,
	))
}

// KubeConfig returns an absolute path to the Kube Config
func (c *Cluster) KubeClient() kubernetes.Interface {
	cfg, err := clientcmd.BuildConfigFromFlags("", c.kubepath)
	gomega.ExpectWithOffset(1, err).To(gomega.BeNil())
	client, err := kubernetes.NewForConfig(cfg)
	gomega.ExpectWithOffset(1, err).To(gomega.BeNil())
	return client

}

func (c *Cluster) runWithOutput(cmd *exec.Cmd) []byte {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	c.run(cmd)
	return stdout.Bytes()
}

func (c *Cluster) run(cmd *exec.Cmd) {
	errPipe, err := cmd.StderrPipe()
	gomega.ExpectWithOffset(1, err).To(gomega.BeNil())

	cmd.Env = append(
		cmd.Env,
		// KIND positions the configuration file relative to HOME.
		// To prevent clobbering an existing KIND installation, override this
		// n.b. HOME isn't always set inside BAZEL
		fmt.Sprintf("HOME=%s", c.tmpDir),
		//needed for Docker. TODO(EKF) Should be properly hermetic
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	)

	// Log output
	go captureOutput(errPipe, "stderr")
	if cmd.Stdout == nil {
		outPipe, err := cmd.StdoutPipe()
		gomega.ExpectWithOffset(1, err).To(gomega.BeNil())
		go captureOutput(outPipe, "stdout")
	}

	gomega.ExpectWithOffset(1, cmd.Run()).To(gomega.BeNil())
}

func captureOutput(pipe io.ReadCloser, label string) {
	buffer := &bytes.Buffer{}
	defer pipe.Close()
	for {
		n, _ := buffer.ReadFrom(pipe)
		if n == 0 {
			return
		}
		fmt.Fprintf(ginkgo.GinkgoWriter, "[%s] %s\n", label, buffer.String())
	}
}
