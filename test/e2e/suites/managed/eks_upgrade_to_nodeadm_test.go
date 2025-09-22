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
	"fmt"

	"github.com/blang/semver"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ref "k8s.io/client-go/tools/reference"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// EKS cluster upgrade tests.
var _ = ginkgo.Describe("EKS Cluster upgrade test", func() {
	var (
		namespace        *corev1.Namespace
		ctx              context.Context
		specName         = "eks-upgrade"
		clusterName      string
		initialVersion   string
		upgradeToVersion string
	)

	shared.ConditionalIt(runUpgradeTests, "[managed] [upgrade] [nodeadm] should create a cluster and upgrade the kubernetes version", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.EksUpgradeFromVersion))
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.EksUpgradeToVersion))
		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		initialVersion = e2eCtx.E2EConfig.MustGetVariable(shared.EksUpgradeFromVersion)
		upgradeToVersion = e2eCtx.E2EConfig.MustGetVariable(shared.EksUpgradeToVersion)

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, clusterName, false, e2eCtx.AWSSession)

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSControlPlaneOnlyFlavor, // TODO (richardcase) - change in the future when upgrades to machinepools work
				ControlPlaneMachineCount: 1,                         // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
				KubernetesVersion:        initialVersion,
			}
		})

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		ginkgo.By("should create a MachineDeployment")
		MachineDeploymentSpec(ctx, func() MachineDeploymentSpecInput {
			return MachineDeploymentSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				ConfigClusterFn:       defaultConfigCluster,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				Namespace:             namespace,
				ClusterName:           clusterName,
				Replicas:              1,
				Cleanup:               false,
			}
		})

		ginkgo.By(fmt.Sprintf("should upgrade control plane to version %s", upgradeToVersion))
		UpgradeControlPlaneVersionSpec(ctx, func() UpgradeControlPlaneVersionSpecInput {
			return UpgradeControlPlaneVersionSpecInput{
				E2EConfig:             e2eCtx.E2EConfig,
				AWSSession:            e2eCtx.BootstrapUserAWSSession,
				BootstrapClusterProxy: e2eCtx.Environment.BootstrapClusterProxy,
				ClusterName:           clusterName,
				Namespace:             namespace,
				UpgradeVersion:        upgradeToVersion,
			}
		})

		ginkgo.By(fmt.Sprintf("should upgrade mahchine deployments to version %s", upgradeToVersion))
		kube133, err := semver.ParseTolerant("1.33.0")
		Expect(err).To(BeNil(), "semver should pass")
		upgradeToVersionParse, err := semver.ParseTolerant(upgradeToVersion)
		Expect(err).To(BeNil(), "semver should pass")

		md := framework.DiscoveryAndWaitForMachineDeployments(ctx, framework.DiscoveryAndWaitForMachineDeploymentsInput{
			Lister:  e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes")...)
		var nodeadmConfigTemplate *eksbootstrapv1.NodeadmConfigTemplate
		if upgradeToVersionParse.GTE(kube133) {
			ginkgo.By("creating a nodeadmconfigtemplate object")
			nodeadmConfigTemplate = &eksbootstrapv1.NodeadmConfigTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-nodeadm-config", clusterName),
					Namespace: namespace.Name,
				},
				Spec: eksbootstrapv1.NodeadmConfigTemplateSpec{
					Template: eksbootstrapv1.NodeadmConfigTemplateResource{
						Spec: eksbootstrapv1.NodeadmConfigSpec{
							PreBootstrapCommands: []string{
								"echo \"hello world\"",
							},
						},
					},
				},
			}
			ginkgo.By("creating the nodeadm config template in the cluster")
			Expect(e2eCtx.Environment.BootstrapClusterProxy.GetClient().Create(ctx, nodeadmConfigTemplate)).To(Succeed())
		}
		ginkgo.By("upgrading machine deployments")
		input := UpgradeMachineDeploymentsAndWaitInput{
			BootstrapClusterProxy:       e2eCtx.Environment.BootstrapClusterProxy,
			Cluster:                     cluster,
			UpgradeVersion:              upgradeToVersion,
			MachineDeployments:          md,
			WaitForMachinesToBeUpgraded: e2eCtx.E2EConfig.GetIntervals("", "wait-worker-nodes"),
		}
		if nodeadmConfigTemplate != nil {
			nodeadmRef, err := ref.GetReference(initScheme(), nodeadmConfigTemplate)
			Expect(err).To(BeNil(), "object should have ref")
			input.UpgradeBootstrapTemplate = nodeadmRef
		}
		UpgradeMachineDeploymentsAndWait(ctx, input)

		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Cluster: cluster,
		})
		framework.WaitForClusterDeleted(ctx, framework.WaitForClusterDeletedInput{
			ClusterProxy:         e2eCtx.Environment.BootstrapClusterProxy,
			Cluster:              cluster,
			ClusterctlConfigPath: e2eCtx.Environment.ClusterctlConfigPath,
			ArtifactFolder:       e2eCtx.Settings.ArtifactFolder,
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-delete-cluster")...)
	})
})
