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

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws/client"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	"sigs.k8s.io/cluster-api/util/patch"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

// UpgradeControlPlaneVersionAndWaitInput is the input type for upgradeControlPlaneVersionAndWait.
type UpgradeControlPlaneVersionSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	AWSSession            client.ConfigProvider
	BootstrapClusterProxy framework.ClusterProxy
	ClusterName           string
	Namespace             *corev1.Namespace
	UpgradeVersion        string
}

// UpgradeControlPlaneVersionSpec updates the EKS control plane version and waits for the upgrade
func UpgradeControlPlaneVersionSpec(ctx context.Context, inputGetter func() UpgradeControlPlaneVersionSpecInput) {
	var (
		input UpgradeControlPlaneVersionSpecInput
	)

	input = inputGetter()

	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.ClusterName).ToNot(BeNil(), "Invalid argument. input.ClusterName can't be nil")
	Expect(input.Namespace).ToNot(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.UpgradeVersion).ToNot(BeNil(), "Invalid argument. input.UpgradeVersion can't be nil")

	mgmtClient := input.BootstrapClusterProxy.GetClient()
	controlPlaneName := getControlPlaneName(input.ClusterName)

	shared.Byf("Getting control plane: %s", controlPlaneName)
	controlPlane := &controlplanev1.AWSManagedControlPlane{}
	err := mgmtClient.Get(ctx, crclient.ObjectKey{Namespace: input.Namespace.Name, Name: controlPlaneName}, controlPlane)
	Expect(err).ToNot(HaveOccurred())

	shared.Byf("Patching control plane %s from %s to %s", controlPlaneName, *controlPlane.Spec.Version, input.UpgradeVersion)
	patchHelper, err := patch.NewHelper(controlPlane, mgmtClient)
	Expect(err).ToNot(HaveOccurred())
	controlPlane.Spec.Version = &input.UpgradeVersion
	Expect(patchHelper.Patch(ctx, controlPlane)).To(Succeed())

	ginkgo.By("Waiting for EKS control-plane to be upgraded to new version")
	waitForControlPlaneToBeUpgraded(ctx, waitForControlPlaneToBeUpgradedInput{
		ControlPlane:   controlPlane,
		AWSSession:     input.AWSSession,
		UpgradeVersion: input.UpgradeVersion,
	}, input.E2EConfig.GetIntervals("", "wait-control-plane-upgrade")...)
}
