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

package managed

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

// ManagedMachinePoolSpecInput is the input for ManagedMachinePoolSpec
type ManagedMachinePoolSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	ConfigClusterFn       DefaultConfigClusterFn
	BootstrapClusterProxy framework.ClusterProxy
	AWSSession            client.ConfigProvider
	Namespace             *corev1.Namespace
	ClusterName           string
	IncludeScaling        bool
	Cleanup               bool
}

// ManagedMachinePoolSpec implements a test for creating a managed machine pool
func ManagedMachinePoolSpec(ctx context.Context, inputGetter func() ManagedMachinePoolSpecInput) {
	var (
		input ManagedMachinePoolSpecInput
	)

	input = inputGetter()
	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.ConfigClusterFn).ToNot(BeNil(), "Invalid argument. input.ConfigClusterFn can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.ClusterName).ShouldNot(HaveLen(0), "Invalid argument. input.ClusterName can't be empty")

	shared.Byf("getting cluster with name %s", input.ClusterName)
	cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
		Getter:    input.BootstrapClusterProxy.GetClient(),
		Namespace: input.Namespace.Name,
		Name:      input.ClusterName,
	})
	Expect(cluster).NotTo(BeNil(), "couldn't find CAPI cluster")

	shared.Byf("creating an applying the %s template", EKSManagedPoolOnlyFlavor)
	configCluster := input.ConfigClusterFn(input.ClusterName, input.Namespace.Name)
	configCluster.Flavor = EKSManagedPoolOnlyFlavor
	configCluster.WorkerMachineCount = pointer.Int64Ptr(1)
	err := shared.ApplyTemplate(ctx, configCluster, input.BootstrapClusterProxy)
	Expect(err).ShouldNot(HaveOccurred())

	shared.Byf("Waiting for the machine pool to be running")
	mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
		Lister:  input.BootstrapClusterProxy.GetClient(),
		Getter:  input.BootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	}, input.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
	Expect(len(mp)).To(Equal(1))

	shared.Byf("Check the status of the node group")
	nodeGroupName := getEKSNodegroupName(input.Namespace.Name, input.ClusterName)
	eksClusterName := getEKSClusterName(input.Namespace.Name, input.ClusterName)
	verifyManagedNodeGroup(input.ClusterName, eksClusterName, nodeGroupName, true, input.AWSSession)

	if input.IncludeScaling { //TODO (richardcase): should this be a separate spec?
		ginkgo.By("Scaling the machine pool up")
		framework.ScaleMachinePoolAndWait(ctx, framework.ScaleMachinePoolAndWaitInput{
			ClusterProxy:              input.BootstrapClusterProxy,
			Cluster:                   cluster,
			Replicas:                  2,
			MachinePools:              mp,
			WaitForMachinePoolToScale: input.E2EConfig.GetIntervals("", "wait-worker-nodes"),
		})

		ginkgo.By("Scaling the machine pool down")
		framework.ScaleMachinePoolAndWait(ctx, framework.ScaleMachinePoolAndWaitInput{
			ClusterProxy:              input.BootstrapClusterProxy,
			Cluster:                   cluster,
			Replicas:                  1,
			MachinePools:              mp,
			WaitForMachinePoolToScale: input.E2EConfig.GetIntervals("", "wait-worker-nodes"),
		})
	}

	if input.Cleanup {
		deleteMachinePool(ctx, deleteMachinePoolInput{
			Deleter:     input.BootstrapClusterProxy.GetClient(),
			MachinePool: mp[0],
		})

		waitForMachinePoolDeleted(ctx, waitForMachinePoolDeletedInput{
			Getter:      input.BootstrapClusterProxy.GetClient(),
			MachinePool: mp[0],
		}, input.E2EConfig.GetIntervals("", "wait-delete-machine-pool")...)
	}
}
