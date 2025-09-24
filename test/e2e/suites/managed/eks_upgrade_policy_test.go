//go:build e2e
// +build e2e

/*
Copyright 2025 The Kubernetes Authors.

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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
)

// EKS upgrade policy test.
var _ = ginkgo.Describe("EKS upgrade policy test", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "cluster"
		clusterName string
	)

	ginkgo.It("[managed] [upgrade-policy] Able to update cluster upgrade policy from STANDARD to EXTENDED", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)

		upgradePolicy := ekscontrolplanev1.UpgradePolicyStandard
		shared.SetEnvVar(shared.UpgradePolicy, upgradePolicy.String(), false)

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))
		eksClusterName := getEKSClusterName(namespace.Name, clusterName)

		ginkgo.By("default iam role should exist")
		VerifyRoleExistsAndOwned(ctx, ekscontrolplanev1.DefaultEKSControlPlaneRole, eksClusterName, false, e2eCtx.AWSSession)

		getManagedClusterSpec := func() ManagedClusterSpecInput {
			return ManagedClusterSpecInput{
				E2EConfig:                e2eCtx.E2EConfig,
				ConfigClusterFn:          defaultConfigCluster,
				BootstrapClusterProxy:    e2eCtx.Environment.BootstrapClusterProxy,
				AWSSession:               e2eCtx.BootstrapUserAWSSession,
				Namespace:                namespace,
				ClusterName:              clusterName,
				Flavour:                  EKSUpgradePolicyFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
			}
		}

		ginkgo.By("should create an EKS control plane")
		ManagedClusterSpec(ctx, getManagedClusterSpec)

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		WaitForEKSClusterUpgradePolicy(ctx, e2eCtx.BootstrapUserAWSSession, eksClusterName, upgradePolicy)

		changedUpgradePolicy := ekscontrolplanev1.UpgradePolicyExtended
		ginkgo.By(fmt.Sprintf("Changing the UpgradePolicy from %s to %s", upgradePolicy, changedUpgradePolicy))
		shared.SetEnvVar(shared.UpgradePolicy, changedUpgradePolicy.String(), false)
		ManagedClusterSpec(ctx, getManagedClusterSpec)
		WaitForEKSClusterUpgradePolicy(ctx, e2eCtx.BootstrapUserAWSSession, eksClusterName, changedUpgradePolicy)

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

func WaitForEKSClusterUpgradePolicy(ctx context.Context, sess *aws.Config, eksClusterName string, upgradePolicy ekscontrolplanev1.UpgradePolicy) {
	ginkgo.By(fmt.Sprintf("Checking EKS control plane upgrade policy matches %s", upgradePolicy))
	Eventually(func() error {
		cluster, err := getEKSCluster(ctx, eksClusterName, sess)
		if err != nil {
			smithyErr := awserrors.ParseSmithyError(err)
			notFoundErr := &ekstypes.ResourceNotFoundException{}
			if smithyErr.ErrorCode() == notFoundErr.ErrorCode() {
				// Unrecoverable error stop trying and fail early.
				return StopTrying(fmt.Sprintf("unrecoverable error: cluster %q not found: %s", eksClusterName, smithyErr.ErrorMessage()))
			}
			return err // For transient errors, retry
		}

		expectedPolicy := converters.SupportTypeToSDK(upgradePolicy)
		actualPolicy := cluster.UpgradePolicy.SupportType

		if actualPolicy != expectedPolicy {
			// The upgrade policy change hasn't been reflected in EKS yet, error and try again.
			return fmt.Errorf("upgrade policy mismatch: expected %s, but found %s", expectedPolicy, actualPolicy)
		}

		// Success in finding the change has been reflected in EKS.
		return nil
	}, 5*time.Minute, 10*time.Second).Should(Succeed(), fmt.Sprintf("eventually failed checking EKS Cluster %q upgrade policy is %s", eksClusterName, upgradePolicy))
}
