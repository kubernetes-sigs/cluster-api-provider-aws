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

package managed

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/net"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// General EKS e2e test.
var _ = ginkgo.Describe("[managed] [general] [ipv6] EKS cluster tests", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "eks-nodes"
		clusterName string
	)

	shared.ConditionalIt(runGeneralTests, "should create a cluster and add nodes", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.CNIAddonVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.CorednsAddonVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubeproxyAddonVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ekscontrolplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.BootstrapUserAWSSession)

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSIPv6ClusterFlavor,
				ControlPlaneMachineCount: 1, //NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       1,
			}
		})

		ginkgo.By("should create a managed node pool and scale")
		MachinePoolSpec(ctx, func() MachinePoolSpecInput {
			return MachinePoolSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				ManagedMachinePool:    true,
				Flavor:                EKSIPv6ClusterFlavor,
			}
		})

		shared.Byf("getting cluster with name %s", clusterName)
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find CAPI cluster")
		// Don't replace this with var pods *corev1.PodsList as that results in a nil pointer error from List.
		pods := &corev1.PodList{}
		listOptions := []client.ListOption{
			client.InNamespace(namespace.Namespace),
			client.MatchingLabels(map[string]string{"app": "aws-node"}),
		}
		clusterClient := e2eCtx.Environment.BootstrapClusterProxy.GetWorkloadCluster(ctx, namespace.Name, clusterName).GetClient()
		err := clusterClient.List(ctx, pods, listOptions...)
		Expect(err).ToNot(HaveOccurred())

		for _, pod := range pods.Items {
			Expect(net.IsIPv6String(pod.Status.PodIP)).To(BeTrue())
		}

		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			Getter:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)
	})
})
