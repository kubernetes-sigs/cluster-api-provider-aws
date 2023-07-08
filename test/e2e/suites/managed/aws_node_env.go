//go:build e2e
// +build e2e

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

package managed

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

type UpdateAwsNodeVersionSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	BootstrapClusterProxy framework.ClusterProxy
	AWSSession            client.ConfigProvider
	Namespace             *corev1.Namespace
	ClusterName           string
}

// CheckAwsNodeEnvVarsSet implements a test for setting environment variables on the `aws-node` DaemonSet.
func CheckAwsNodeEnvVarsSet(ctx context.Context, inputGetter func() UpdateAwsNodeVersionSpecInput) {
	input := inputGetter()
	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.ClusterName).ShouldNot(BeEmpty(), "Invalid argument. input.ClusterName can't be empty")

	mgmtClient := input.BootstrapClusterProxy.GetClient()
	controlPlaneName := getControlPlaneName(input.ClusterName)

	By(fmt.Sprintf("Getting control plane: %s", controlPlaneName))
	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	err := mgmtClient.Get(ctx, crclient.ObjectKey{Namespace: input.Namespace.Name, Name: controlPlaneName}, controlPlane)
	Expect(err).ToNot(HaveOccurred())

	By(fmt.Sprintf("Checking environment variables are set on AWSManagedControlPlane: %s", controlPlaneName))
	Expect(controlPlane.Spec.VpcCni.Env).NotTo(BeNil())
	Expect(len(controlPlane.Spec.VpcCni.Env)).Should(BeNumerically(">", 1))

	By("Checking if aws-node has been updated with the defined environment variables on the workload cluster")
	daemonSet := &appsv1.DaemonSet{}

	clusterClient := input.BootstrapClusterProxy.GetWorkloadCluster(ctx, input.Namespace.Name, input.ClusterName).GetClient()
	err = clusterClient.Get(ctx, crclient.ObjectKey{Namespace: "kube-system", Name: "aws-node"}, daemonSet)
	Expect(err).ToNot(HaveOccurred())

	for _, container := range daemonSet.Spec.Template.Spec.Containers {
		if container.Name == "aws-node" {
			Expect(matchEnvVar(container.Env, corev1.EnvVar{Name: "FOO", Value: "BAR"})).Should(BeTrue())
			Expect(matchEnvVar(container.Env, corev1.EnvVar{Name: "ENABLE_PREFIX_DELEGATION", Value: "true"})).Should(BeTrue())
			break
		}
	}
}

func matchEnvVar(s []corev1.EnvVar, ev corev1.EnvVar) bool {
	for _, e := range s {
		if e == ev {
			return true
		}
	}
	return false
}
