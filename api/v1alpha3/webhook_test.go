/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha3

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAWSClusterConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	clusterName := fmt.Sprintf("test-cluster-%s", util.RandomString(5))
	cluster := &AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: ns.Name,
		},
	}

	g.Expect(testEnv.Create(ctx, cluster)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, cluster)
}

func TestAWSMachineConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	machineName := fmt.Sprintf("test-machine-%s", util.RandomString(5))
	machine := &AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineName,
			Namespace: ns.Name,
		},
	}

	g.Expect(testEnv.Create(ctx, machine)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, machine)
}

func TestAWSMachineTemplateConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	machineTemplateName := fmt.Sprintf("test-machine-%s", util.RandomString(5))
	machine := &AWSMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineTemplateName,
			Namespace: ns.Name,
		},
	}

	g.Expect(testEnv.Create(ctx, machine)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, machine)
}
