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
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

func defaultConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	return clusterctl.ConfigClusterInput{
		LogFolder:                filepath.Join(e2eCtx.Settings.ArtifactFolder, "clusters", e2eCtx.Environment.BootstrapClusterProxy.GetName()),
		ClusterctlConfigPath:     e2eCtx.Environment.ClusterctlConfigPath,
		KubeconfigPath:           e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath(),
		InfrastructureProvider:   "aws",
		Flavor:                   EKSManagedPoolFlavor,
		Namespace:                namespace,
		ClusterName:              clusterName,
		KubernetesVersion:        e2eCtx.E2EConfig.MustGetVariable(shared.KubernetesVersion),
		ControlPlaneMachineCount: ptr.To[int64](1),
		WorkerMachineCount:       ptr.To[int64](0),
	}
}

func upgradeFromConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	cfg := defaultConfigCluster(clusterName, namespace)
	cfg.KubernetesVersion = e2eCtx.E2EConfig.MustGetVariable(shared.EksUpgradeFromVersion)
	return cfg
}

func nitroEnclaveConfigCluster(clusterName, namespace string) clusterctl.ConfigClusterInput {
	cfg := defaultConfigCluster(clusterName, namespace)
	cfg.KubernetesVersion = "v1.35.0"
	return cfg
}

func assertEKSNitroEnclaveEnabled(ctx context.Context, namespace, clusterName string) {
	ginkgo.By(fmt.Sprintf("Verifying Nitro Enclave enabled on instances in cluster: %s", clusterName))
	eksClusterName := getEKSClusterName(namespace, clusterName)
	ec2Client := ec2.NewFromConfig(*e2eCtx.AWSSession)
	describeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/cluster/" + eksClusterName),
				Values: []string{"owned"},
			},
			{
				Name:   aws.String("tag:sigs.k8s.io/cluster-api-provider-aws/role"),
				Values: []string{"node"},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	}
	// EKS managed node group instances may not be running immediately after the cluster is
	// active, so poll until they appear and all have EnclaveOptions.Enabled=true.
	Eventually(func(g Gomega) {
		result, err := ec2Client.DescribeInstances(ctx, describeInput)
		g.Expect(err).To(BeNil())
		g.Expect(result.Reservations).ToNot(BeEmpty())
		for _, r := range result.Reservations {
			for _, inst := range r.Instances {
				g.Expect(inst.EnclaveOptions).ToNot(BeNil())
				g.Expect(aws.ToBool(inst.EnclaveOptions.Enabled)).To(BeTrue(),
					"expected instance %s to have EnclaveOptions.Enabled=true", aws.ToString(inst.InstanceId))
			}
		}
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...).Should(Succeed())
}
