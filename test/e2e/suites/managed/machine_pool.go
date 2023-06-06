//go:build e2e
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
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// MachinePoolSpecInput is the input for MachinePoolSpec.
type MachinePoolSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	ConfigClusterFn       DefaultConfigClusterFn
	BootstrapClusterProxy framework.ClusterProxy
	AWSSession            client.ConfigProvider
	Namespace             *corev1.Namespace
	ClusterName           string
	IncludeScaling        bool
	Cleanup               bool
	ManagedMachinePool    bool
	Flavor                string
	UsesLaunchTemplate    bool
}

// MachinePoolSpec implements a test for creating a machine pool.
func MachinePoolSpec(ctx context.Context, inputGetter func() MachinePoolSpecInput) {
	input := inputGetter()
	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.ConfigClusterFn).ToNot(BeNil(), "Invalid argument. input.ConfigClusterFn can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.ClusterName).ShouldNot(BeEmpty(), "Invalid argument. input.ClusterName can't be empty")
	Expect(input.Flavor).ShouldNot(BeEmpty(), "Invalid argument. input.Flavor can't be empty")

	ginkgo.By(fmt.Sprintf("getting cluster with name %s", input.ClusterName))
	cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
		Getter:    input.BootstrapClusterProxy.GetClient(),
		Namespace: input.Namespace.Name,
		Name:      input.ClusterName,
	})
	Expect(cluster).NotTo(BeNil(), "couldn't find CAPI cluster")

	ginkgo.By(fmt.Sprintf("creating an applying the %s template", input.Flavor))
	configCluster := input.ConfigClusterFn(input.ClusterName, input.Namespace.Name)
	configCluster.Flavor = input.Flavor
	configCluster.WorkerMachineCount = pointer.Int64(1)
	workloadClusterTemplate := shared.GetTemplate(ctx, configCluster)
	if input.UsesLaunchTemplate {
		userDataTemplate := `#!/bin/bash
/etc/eks/bootstrap.sh %s \
  --container-runtime containerd
`
		eksClusterName := getEKSClusterName(input.Namespace.Name, input.ClusterName)
		userData := fmt.Sprintf(userDataTemplate, eksClusterName)
		userDataEncoded := base64.StdEncoding.EncodeToString([]byte(userData))
		workloadClusterTemplate = []byte(strings.ReplaceAll(string(workloadClusterTemplate), "USER_DATA", userDataEncoded))
	}
	ginkgo.By(string(workloadClusterTemplate))
	ginkgo.By(fmt.Sprintf("Applying the %s cluster template yaml to the cluster", configCluster.Flavor))
	err := input.BootstrapClusterProxy.Apply(ctx, workloadClusterTemplate)
	Expect(err).ShouldNot(HaveOccurred())

	ginkgo.By("Waiting for the machine pool to be running")
	mp := framework.DiscoveryAndWaitForMachinePools(ctx, framework.DiscoveryAndWaitForMachinePoolsInput{
		Lister:  input.BootstrapClusterProxy.GetClient(),
		Getter:  input.BootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	}, input.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
	Expect(len(mp)).To(Equal(1))

	ginkgo.By("Check the status of the node group")
	eksClusterName := getEKSClusterName(input.Namespace.Name, input.ClusterName)
	if input.ManagedMachinePool {
		var nodeGroupName string
		if input.UsesLaunchTemplate {
			nodeGroupName = getEKSNodegroupWithLaunchTemplateName(input.Namespace.Name, input.ClusterName)
		} else {
			nodeGroupName = getEKSNodegroupName(input.Namespace.Name, input.ClusterName)
		}
		verifyManagedNodeGroup(eksClusterName, nodeGroupName, true, input.AWSSession)
	} else {
		asgName := getASGName(input.ClusterName)
		verifyASG(eksClusterName, asgName, true, input.AWSSession)
	}

	if input.IncludeScaling { // TODO (richardcase): should this be a separate spec?
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
