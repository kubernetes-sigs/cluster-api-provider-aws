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

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/eks"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

type CheckAddonExistsSpecInput struct {
	E2EConfig             *clusterctl.E2EConfig
	BootstrapClusterProxy framework.ClusterProxy
	AWSSession            client.ConfigProvider
	Namespace             *corev1.Namespace
	ClusterName           string
	AddonName             string
	AddonVersion          string
}

// CheckAddonExistsSpec implements a test for a cluster having an addon
func CheckAddonExistsSpec(ctx context.Context, inputGetter func() CheckAddonExistsSpecInput) {
	var (
		input CheckAddonExistsSpecInput
	)

	input = inputGetter()
	Expect(input.E2EConfig).ToNot(BeNil(), "Invalid argument. input.E2EConfig can't be nil")
	Expect(input.BootstrapClusterProxy).ToNot(BeNil(), "Invalid argument. input.BootstrapClusterProxy can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.Namespace).NotTo(BeNil(), "Invalid argument. input.Namespace can't be nil")
	Expect(input.ClusterName).ShouldNot(HaveLen(0), "Invalid argument. input.ClusterName can't be empty")
	Expect(input.AddonName).ShouldNot(HaveLen(0), "Invalid argument. input.AddonName can't be empty")
	Expect(input.AddonVersion).ShouldNot(HaveLen(0), "Invalid argument. input.AddonVersion can't be empty")

	mgmtClient := input.BootstrapClusterProxy.GetClient()
	controlPlaneName := getControlPlaneName(input.ClusterName)

	shared.Byf("Getting control plane: %s", controlPlaneName)
	controlPlane := &controlplanev1.AWSManagedControlPlane{}
	err := mgmtClient.Get(ctx, crclient.ObjectKey{Namespace: input.Namespace.Name, Name: controlPlaneName}, controlPlane)
	Expect(err).ToNot(HaveOccurred())

	shared.Byf("Checking EKS addon %s is installed on cluster %s and is active", input.AddonName, input.ClusterName)
	waitForEKSAddonToHaveStatus(ctx, waitForEKSAddonToHaveStatusInput{
		ControlPlane: controlPlane,
		AWSSession:   input.AWSSession,
		AddonName:    input.AddonName,
		AddonVersion: input.AddonVersion,
		AddonStatus:  eks.AddonStatusActive,
	}, input.E2EConfig.GetIntervals("", "wait-addon-status")...)
}
