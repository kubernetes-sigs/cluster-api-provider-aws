//go:build e2e
// +build e2e

/*
Copyright 2024 The Kubernetes Authors.

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

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

// Externally managed AWSManagedControlPlane e2e test.
var _ = ginkgo.Describe("[managed] [general] EKS externally managed control plane", func() {
	var (
		namespace   *corev1.Namespace
		ctx         context.Context
		specName    = "eks-externally-managed"
		clusterName string
	)

	shared.ConditionalIt(runGeneralTests, "should skip reconciliation when AWSManagedControlPlane is externally managed", func() {
		ginkgo.By("should have a valid test configuration")
		Expect(e2eCtx.Environment.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. BootstrapClusterProxy can't be nil")
		Expect(e2eCtx.E2EConfig).ToNot(BeNil(), "Invalid argument. e2eConfig can't be nil when calling %s spec", specName)
		Expect(e2eCtx.E2EConfig.Variables).To(HaveKey(shared.KubernetesVersion))

		ctx = context.TODO()
		namespace = shared.SetupSpecNamespace(ctx, specName, e2eCtx)
		clusterName = fmt.Sprintf("%s-%s", specName, util.RandomString(6))

		ginkgo.By("applying the externally managed cluster template")
		configCluster := defaultConfigCluster(clusterName, namespace.Name)
		configCluster.Flavor = EKSExternallyManagedFlavor
		configCluster.ControlPlaneMachineCount = ptr.To[int64](1)
		configCluster.WorkerMachineCount = ptr.To[int64](0)
		err := shared.ApplyTemplate(ctx, configCluster, e2eCtx.Environment.BootstrapClusterProxy)
		Expect(err).ShouldNot(HaveOccurred())

		mgmtClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
		controlPlaneName := getControlPlaneName(clusterName)

		ginkgo.By("waiting for AWSManagedControlPlane to exist")
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		Eventually(func() error {
			return mgmtClient.Get(ctx, crclient.ObjectKey{
				Namespace: namespace.Name,
				Name:      controlPlaneName,
			}, controlPlane)
		}, 2*time.Minute, 5*time.Second).Should(Succeed(), "failed to get AWSManagedControlPlane")

		ginkgo.By("verifying the externally managed annotation is present")
		Expect(controlPlane.Annotations).To(HaveKeyWithValue("cluster.x-k8s.io/managed-by", "e2e-test"))

		ginkgo.By("verifying no EKS cluster was created in AWS")
		eksClusterName := getEKSClusterName(namespace.Name, clusterName)
		eksClient := eks.NewFromConfig(*e2eCtx.AWSSession)
		Consistently(func() error {
			_, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{
				Name: &eksClusterName,
			})
			if err != nil {
				return nil
			}
			return fmt.Errorf("EKS cluster %q should not exist since control plane is externally managed", eksClusterName)
		}, 30*time.Second, 5*time.Second).Should(Succeed())

		ginkgo.By("patching AWSManagedCluster status to simulate external management")
		awsManagedCluster := &infrav1.AWSManagedCluster{}
		Expect(mgmtClient.Get(ctx, crclient.ObjectKey{
			Namespace: namespace.Name,
			Name:      clusterName,
		}, awsManagedCluster)).To(Succeed())
		awsManagedCluster.Status.Ready = true
		Expect(mgmtClient.Status().Update(ctx, awsManagedCluster)).To(Succeed())

		ginkgo.By("patching AWSManagedControlPlane status to simulate external management")
		// Re-fetch the control plane to get the latest version
		Expect(mgmtClient.Get(ctx, crclient.ObjectKey{
			Namespace: namespace.Name,
			Name:      controlPlaneName,
		}, controlPlane)).To(Succeed())

		// Patch status fields
		controlPlane.Status.Ready = true
		controlPlane.Status.Initialized = true
		controlPlane.Status.ExternalManagedControlPlane = ptr.To(true)
		controlPlane.Status.FailureDomains = clusterv1beta1.FailureDomains{
			"us-east-1a": clusterv1beta1.FailureDomainSpec{
				ControlPlane: true,
			},
		}
		// Set the EKSControlPlaneReady condition
		v1beta1conditions.MarkTrue(controlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition)
		Expect(mgmtClient.Status().Update(ctx, controlPlane)).To(Succeed())

		ginkgo.By("waiting for CAPI Cluster to report infrastructure provisioned")
		Eventually(func() bool {
			cluster := &clusterv1.Cluster{}
			if err := mgmtClient.Get(ctx, crclient.ObjectKey{
				Namespace: namespace.Name,
				Name:      clusterName,
			}, cluster); err != nil {
				return false
			}
			return ptr.Deref(cluster.Status.Initialization.InfrastructureProvisioned, false)
		}, 5*time.Minute, 10*time.Second).Should(BeTrue(), "CAPI Cluster should report InfrastructureProvisioned=true")

		ginkgo.By("waiting for CAPI Cluster to report control plane initialized")
		Eventually(func() bool {
			cluster := &clusterv1.Cluster{}
			if err := mgmtClient.Get(ctx, crclient.ObjectKey{
				Namespace: namespace.Name,
				Name:      clusterName,
			}, cluster); err != nil {
				return false
			}
			return ptr.Deref(cluster.Status.Initialization.ControlPlaneInitialized, false)
		}, 5*time.Minute, 10*time.Second).Should(BeTrue(), "CAPI Cluster should report ControlPlaneInitialized=true")

		ginkgo.By("verifying AWSManagedControlPlane has no finalizer (reconciliation was skipped)")
		Expect(mgmtClient.Get(ctx, crclient.ObjectKey{
			Namespace: namespace.Name,
			Name:      controlPlaneName,
		}, controlPlane)).To(Succeed())
		Expect(controlPlane.Finalizers).To(BeEmpty(), "externally managed control plane should have no finalizer")

		ginkgo.By("deleting the cluster")
		cluster := framework.GetClusterByName(ctx, framework.GetClusterByNameInput{
			Getter:    mgmtClient,
			Namespace: namespace.Name,
			Name:      clusterName,
		})
		Expect(cluster).NotTo(BeNil(), "couldn't find CAPI cluster")

		framework.DeleteCluster(ctx, framework.DeleteClusterInput{
			Deleter: mgmtClient,
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
