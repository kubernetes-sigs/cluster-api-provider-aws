//go:build e2e
// +build e2e

/*
Copyright 2026 The Kubernetes Authors.

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
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
)

// EKS control plane component config test.
var _ = ginkgo.Describe("EKS control plane component config test", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "cluster"
		clusterName string
	)

	ginkgo.It("[managed] [control-plane-config] Able to create cluster with scheduler and API server config and update the event TTL", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)

		eventTTL := "30m"
		shared.SetEnvVar(shared.EventTTL, eventTTL, false)

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
				Flavour:                  EKSControlPlaneConfigFlavor,
				ControlPlaneMachineCount: 1, // NOTE: this cannot be zero as clusterctl returns an error
				WorkerMachineCount:       0,
			}
		}

		ginkgo.By("should create an EKS control plane with scheduler and API server config")
		ManagedClusterSpec(ctx, getManagedClusterSpec)

		ginkgo.By(fmt.Sprintf("getting cluster with name %s", clusterName))
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    e2eCtx.Environment.BootstrapClusterProxy.GetClient(),
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find cluster")

		WaitForEKSClusterControlPlaneComponentConfig(ctx, e2eCtx.BootstrapUserAWSSession, eksClusterName, ekstypes.ScoringStrategyTypeMostAllocated, eventTTL)

		changedEventTTL := "15m"
		ginkgo.By(fmt.Sprintf("Changing the event TTL from %s to %s", eventTTL, changedEventTTL))

		mgmtClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}

		Eventually(func() error {
			return mgmtClient.Get(ctx, crclient.ObjectKey{Namespace: namespace.Name, Name: getControlPlaneName(clusterName)}, controlPlane)
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-client-request")...).Should(Succeed(), "eventually failed trying to get the AWSManagedControlPlane")

		patchHelper, err := patch.NewHelper(controlPlane, mgmtClient)
		Expect(err).ToNot(HaveOccurred())
		controlPlane.Spec.KubeAPIServerConfig.EventTTL = changedEventTTL

		Eventually(func() error {
			return patchHelper.Patch(ctx, controlPlane)
		}, e2eCtx.E2EConfig.GetIntervals("", "wait-client-request")...).Should(Succeed(), "eventually failed patching the AWSManagedControlPlane")

		WaitForEKSClusterControlPlaneComponentConfig(ctx, e2eCtx.BootstrapUserAWSSession, eksClusterName, ekstypes.ScoringStrategyTypeMostAllocated, changedEventTTL)

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

// WaitForEKSClusterControlPlaneComponentConfig polls the AWS EKS API until the cluster's scheduler scoring
// strategy and event TTL match the expected values, failing early if the cluster is not found.
func WaitForEKSClusterControlPlaneComponentConfig(ctx context.Context, sess *aws.Config, eksClusterName string, scoringStrategy ekstypes.ScoringStrategyType, eventTTL string) {
	ginkgo.By(fmt.Sprintf("Checking EKS control plane scoring strategy is %s and event TTL is %s", scoringStrategy, eventTTL))
	expectedTTL, err := time.ParseDuration(eventTTL)
	Expect(err).ToNot(HaveOccurred())

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

		if cluster.KubeSchedulerConfig == nil || cluster.KubeSchedulerConfig.NodeResourcesFit == nil || cluster.KubeSchedulerConfig.NodeResourcesFit.ScoringStrategy == nil {
			return fmt.Errorf("scheduler config is not reported for cluster %q", eksClusterName)
		}
		if actual := cluster.KubeSchedulerConfig.NodeResourcesFit.ScoringStrategy.Type; actual != scoringStrategy {
			return fmt.Errorf("scoring strategy mismatch: expected %s, but found %s", scoringStrategy, actual)
		}

		if cluster.KubeApiServerConfig == nil || cluster.KubeApiServerConfig.EventTtl == nil {
			return fmt.Errorf("API server config is not reported for cluster %q", eksClusterName)
		}
		actualTTL, err := time.ParseDuration(aws.ToString(cluster.KubeApiServerConfig.EventTtl))
		if err != nil {
			return fmt.Errorf("failed to parse event TTL %q: %w", aws.ToString(cluster.KubeApiServerConfig.EventTtl), err)
		}
		if actualTTL != expectedTTL {
			return fmt.Errorf("event TTL mismatch: expected %s, but found %s", expectedTTL, actualTTL)
		}

		return nil
	}, 10*time.Minute, 10*time.Second).Should(Succeed(), fmt.Sprintf("eventually failed checking EKS Cluster %q control plane component config", eksClusterName))
}
