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

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws/client"

	controlplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
)

type waitForEKSAddonToHaveStatusInput struct {
	ControlPlane *controlplanev1.AWSManagedControlPlane
	AWSSession   client.ConfigProvider
	AddonName    string
	AddonVersion string
	AddonStatus  string
}

func waitForEKSAddonToHaveStatus(ctx context.Context, input waitForEKSAddonToHaveStatusInput, intervals ...interface{}) {
	Expect(input.ControlPlane).ToNot(BeNil(), "Invalid argument. input.ControlPlane can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.AddonName).ShouldNot(HaveLen(0), "Invalid argument. input.AddonName can't be empty")
	Expect(input.AddonVersion).ShouldNot(HaveLen(0), "Invalid argument. input.AddonVersion can't be empty")
	Expect(input.AddonStatus).ShouldNot(HaveLen(0), "Invalid argument. input.AddonStatus can't be empty")

	ginkgo.By(fmt.Sprintf("Ensuring EKS addon %s has status %s for EKS cluster %s", input.AddonName, input.AddonStatus, input.ControlPlane.Spec.EKSClusterName))

	Eventually(func() (bool, error) {
		installedAddon, err := getEKSClusterAddon(input.ControlPlane.Spec.EKSClusterName, input.AddonName, input.AWSSession)
		if err != nil {
			return false, err
		}

		if installedAddon == nil {
			return false, err
		}

		if *installedAddon.Status != input.AddonStatus {
			return false, nil
		}

		return true, nil

	}, intervals...).Should(BeTrue())
}
