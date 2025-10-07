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

package v1beta2

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func TestValidateROSANetwork(t *testing.T) {
	g := NewGomegaWithT(t)

	rosaCP := &ROSAControlPlane{
		Spec:   RosaControlPlaneSpec{},
		Status: RosaControlPlaneStatus{},
	}

	t.Run("Validation error when no ROSANetworkRef, no subnets, no AZs", func(t *testing.T) {
		err := rosaCP.validateROSANetwork()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets cannot be empty"))
	})

	t.Run("Validation error when no ROSANetworkRef, subnets present, no AZs", func(t *testing.T) {
		rosaCP.Spec.Subnets = []string{"subnet01", "subnet02"}
		err := rosaCP.validateROSANetwork()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.availabilityZones cannot be empty"))
	})

	t.Run("Validation succeeds when no ROSANetworkRef, subnets and AZs are present", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = []string{"AZ01", "AZ02"}
		err := rosaCP.validateROSANetwork()
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Validation error when ROSANetworkRef, subnets and AZs are present", func(t *testing.T) {
		rosaCP.Spec.ROSANetworkRef = &corev1.LocalObjectReference{}
		err := rosaCP.validateROSANetwork()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation error when ROSANetworkRef and subnets are present, no AZs", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = nil
		err := rosaCP.validateROSANetwork()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation error when ROSANetworkRef and AZs are present, no subnets", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = []string{"AZ01", "AZ02"}
		rosaCP.Spec.Subnets = nil
		err := rosaCP.validateROSANetwork()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.availabilityZones and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation succeeds when ROSANetworkRef is present, no subnets and no AZs", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = nil
		rosaCP.Spec.Subnets = nil
		err := rosaCP.validateROSANetwork()
		g.Expect(err).NotTo(HaveOccurred())
	})
}
