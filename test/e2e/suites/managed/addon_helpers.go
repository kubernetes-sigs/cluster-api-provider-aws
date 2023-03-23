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
	"fmt"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

type waitForEKSAddonToHaveStatusInput struct {
	ControlPlane *ekscontrolplanev1.AWSManagedControlPlane
	AWSSession   client.ConfigProvider
	AddonName    string
	AddonVersion string
	AddonStatus  []string
}

func waitForEKSAddonToHaveStatus(input waitForEKSAddonToHaveStatusInput, intervals ...interface{}) {
	Expect(input.ControlPlane).ToNot(BeNil(), "Invalid argument. input.ControlPlane can't be nil")
	Expect(input.AWSSession).ToNot(BeNil(), "Invalid argument. input.AWSSession can't be nil")
	Expect(input.AddonName).ShouldNot(BeEmpty(), "Invalid argument. input.AddonName can't be empty")
	Expect(input.AddonVersion).ShouldNot(BeEmpty(), "Invalid argument. input.AddonVersion can't be empty")
	Expect(input.AddonStatus).ShouldNot(BeEmpty(), "Invalid argument. input.AddonStatus can't be empty")

	ginkgo.By(fmt.Sprintf("Ensuring EKS addon %s has status in %q for EKS cluster %s", input.AddonName, input.AddonStatus, input.ControlPlane.Spec.EKSClusterName))

	Eventually(func() (bool, error) {
		installedAddon, err := getEKSClusterAddon(input.ControlPlane.Spec.EKSClusterName, input.AddonName, input.AWSSession)
		if err != nil {
			return false, err
		}

		if installedAddon == nil {
			return false, err
		}

		for i := range input.AddonStatus {
			wantedStatus := input.AddonStatus[i]

			if wantedStatus == *installedAddon.Status {
				return true, nil
			}
		}

		return false, nil
	}, intervals...).Should(BeTrue())
}
