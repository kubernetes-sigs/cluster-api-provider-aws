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

func TestAWSMachinePoolConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	machinepool := &AWSMachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-machinepool-%s", util.RandomString(5)),
			Namespace: ns.Name,
		},
		Spec: AWSMachinePoolSpec{
			MinSize: 1,
			MaxSize: 3,
		},
	}

	g.Expect(testEnv.Create(ctx, machinepool)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, machinepool)
}

func TestAWSManagedMachinePoolConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	managedMachinepool := &AWSManagedMachinePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-managedmachinepool-%s", util.RandomString(5)),
			Namespace: ns.Name,
		},
	}

	g.Expect(testEnv.Create(ctx, managedMachinepool)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, managedMachinepool)
}

func TestAWSFargateProfileConversion(t *testing.T) {
	g := NewWithT(t)
	ns, err := testEnv.CreateNamespace(ctx, fmt.Sprintf("conversion-webhook-%s", util.RandomString(5)))
	g.Expect(err).ToNot(HaveOccurred())
	fargateProfile := &AWSFargateProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("test-fargate-%s", util.RandomString(5)),
			Namespace: ns.Name,
		},
		Spec: FargateProfileSpec{
			ClusterName: "cluster-name",
			ProfileName: "name",
		},
	}

	g.Expect(testEnv.Create(ctx, fargateProfile)).To(Succeed())
	defer func(do ...client.Object) {
		g.Expect(testEnv.Cleanup(ctx, do...)).To(Succeed())
	}(ns, fargateProfile)
}
