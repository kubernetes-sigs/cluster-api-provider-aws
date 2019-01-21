/*
Copyright 2019 The Kubernetes Authors.

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
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

const (
	kubeconfigEnvVar = "KUBECONFIG"
	kindClusterName  = "clusterapi"
)

type Kind struct {
	options []string
	// execFunc implemented as function variable for testing hooks
	execFunc func(args ...string) (string, error)
}

func New() *Kind {
	return WithOptions([]string{})
}

func WithOptions(options []string) *Kind {
	return &Kind{
		execFunc: execFunc,
		options:  append(options, fmt.Sprintf("name=%s", kindClusterName)),
	}
}

var execFunc = func(args ...string) (string, error) {
	const executable = "kind"
	klog.V(3).Infof("Running: %v %v", executable, args)
	cmd := exec.Command(executable, args...)
	cmdOut, err := cmd.CombinedOutput()
	klog.V(2).Infof("Ran: %v %v Output: %v", executable, args, string(cmdOut))
	if err != nil {
		err = errors.Wrapf(err, "error running command '%v %v'", executable, strings.Join(args, " "))
	}
	return string(cmdOut), err
}

func (k *Kind) Create() error {
	args := []string{"create", "cluster"}
	for _, opt := range k.options {
		args = append(args, fmt.Sprintf("--%v", opt))
	}

	_, err := k.exec(args...)
	return err
}

func (k *Kind) Delete() error {
	args := []string{"delete", "cluster"}
	for _, opt := range k.options {
		args = append(args, fmt.Sprintf("--%v", opt))
	}

	_, err := k.exec(args...)
	return err
}

func (k *Kind) GetKubeconfig() (string, error) {
	path, err := k.getKubeConfigPath()
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (k *Kind) getKubeConfigPath() (string, error) {
	args := []string{"get", "kubeconfig-path"}
	for _, opt := range k.options {
		args = append(args, fmt.Sprintf("--%v", opt))
	}

	out, err := k.exec(args...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func (k *Kind) exec(args ...string) (string, error) {
	return k.execFunc(args...)
}
