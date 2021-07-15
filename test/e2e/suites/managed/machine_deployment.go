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
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

// MachineDeploymentSpecInput is the input for MachineDeploymentSpec
type MachineDeploymentSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	ConfigClusterFn       DefaultConfigClusterFn
	BootstrapClusterProxy framework.ClusterProxy
	AWSSession            client.ConfigProvider
	Namespace             *corev1.Namespace
	Replicas              int64
	ClusterName           string
	Cleanup               bool
}

// MachineDeploymentSpec implements a test for creating a machine deployment for use with CAPA
func MachineDeploymentSpec(ctx context.Context, inputGetter func() MachineDeploymentSpecInput) {
	var (
		input MachineDeploymentSpecInput
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

	shared.Byf("creating an applying the %s template", EKSMachineDeployOnlyFlavor)
	configCluster := input.ConfigClusterFn(input.ClusterName, input.Namespace.Name)
	configCluster.Flavor = EKSMachineDeployOnlyFlavor
	configCluster.WorkerMachineCount = pointer.Int64Ptr(input.Replicas)
	err := shared.ApplyTemplate(ctx, configCluster, input.BootstrapClusterProxy)
	Expect(err).ShouldNot(HaveOccurred())

	shared.Byf("Waiting for the worker node to be running")
	md := framework.DiscoveryAndWaitForMachineDeployments(ctx, framework.DiscoveryAndWaitForMachineDeploymentsInput{
		Lister:  input.BootstrapClusterProxy.GetClient(),
		Cluster: cluster,
	}, input.E2EConfig.GetIntervals("", "wait-worker-nodes")...)

	workerMachines := framework.GetMachinesByMachineDeployments(ctx, framework.GetMachinesByMachineDeploymentsInput{
		Lister:            input.BootstrapClusterProxy.GetClient(),
		ClusterName:       input.ClusterName,
		Namespace:         input.Namespace.Name,
		MachineDeployment: *md[0],
	})

	Expect(len(workerMachines)).To(Equal(1))

	statusChecks := []framework.MachineStatusCheck{framework.MachinePhaseCheck(string(clusterv1.MachinePhaseRunning))}
	machineStatusInput := framework.WaitForMachineStatusCheckInput{
		Getter:       input.BootstrapClusterProxy.GetClient(),
		Machine:      &workerMachines[0],
		StatusChecks: statusChecks,
	}
	framework.WaitForMachineStatusCheck(ctx, machineStatusInput, input.E2EConfig.GetIntervals("", "wait-machine-status")...)

	if input.Cleanup {
		deleteMachineDeployment(ctx, deleteMachineDeploymentInput{
			Deleter:           input.BootstrapClusterProxy.GetClient(),
			MachineDeployment: md[0],
		})
		// deleteMachine(ctx, deleteMachineInput{
		// 	Deleter: input.BootstrapClusterProxy.GetClient(),
		// 	Machine: &workerMachines[0],
		// })

		waitForMachineDeploymentDeleted(ctx, waitForMachineDeploymentDeletedInput{
			Getter:            input.BootstrapClusterProxy.GetClient(),
			MachineDeployment: md[0],
		}, input.E2EConfig.GetIntervals("", "wait-delete-machine-deployment")...)

		waitForMachineDeleted(ctx, waitForMachineDeletedInput{
			Getter:  input.BootstrapClusterProxy.GetClient(),
			Machine: &workerMachines[0],
		}, input.E2EConfig.GetIntervals("", "wait-delete-machine")...)
	}
}
