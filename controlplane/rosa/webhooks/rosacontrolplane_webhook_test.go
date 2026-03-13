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

package webhooks

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
)

func TestValidateROSANetwork(t *testing.T) {
	g := NewGomegaWithT(t)

	rosaCP := &rosacontrolplanev1.ROSAControlPlane{
		Spec:   rosacontrolplanev1.RosaControlPlaneSpec{},
		Status: rosacontrolplanev1.RosaControlPlaneStatus{},
	}

	w := &ROSAControlPlane{}

	t.Run("Validation error when no ROSANetworkRef, no subnets, no AZs", func(t *testing.T) {
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets cannot be empty"))
	})

	t.Run("Validation error when no ROSANetworkRef, subnets present, no AZs", func(t *testing.T) {
		rosaCP.Spec.Subnets = []string{"subnet01", "subnet02"}
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.availabilityZones cannot be empty"))
	})

	t.Run("Validation succeeds when no ROSANetworkRef, subnets and AZs are present", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = []string{"AZ01", "AZ02"}
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Validation error when ROSANetworkRef, subnets and AZs are present", func(t *testing.T) {
		rosaCP.Spec.ROSANetworkRef = &corev1.LocalObjectReference{}
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation error when ROSANetworkRef and subnets are present, no AZs", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = nil
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.subnets and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation error when ROSANetworkRef and AZs are present, no subnets", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = []string{"AZ01", "AZ02"}
		rosaCP.Spec.Subnets = nil
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spec.availabilityZones and spec.rosaNetworkRef are mutually exclusive"))
	})

	t.Run("Validation succeeds when ROSANetworkRef is present, no subnets and no AZs", func(t *testing.T) {
		rosaCP.Spec.AvailabilityZones = nil
		rosaCP.Spec.Subnets = nil
		err := w.validateROSANetworkRef(rosaCP)
		g.Expect(err).NotTo(HaveOccurred())
	})
}

func TestValidateChannel(t *testing.T) {
	g := NewGomegaWithT(t)

	w := &ROSAControlPlane{}

	t.Run("Validation succeeds when channel is not specified", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Validation succeeds for valid stable channel", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "stable-4.16",
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Validation succeeds for valid eus channel", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "eus-4.16",
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Validation error for invalid channel format", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "invalid",
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("must be in format '<channelGroup>-<major>.<minor>'"))
	})

	t.Run("Validation error for invalid channel group", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "invalid-4.16",
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("channel group must be one of"))
	})

	t.Run("Validation error for invalid Y-stream format", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "stable-4",
				Version: "4.16.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("Y-stream must be in format '<major>.<minor>'"))
	})

	t.Run("Validation error when channel Y-stream doesn't match version Y-stream", func(t *testing.T) {
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Channel: "stable-4.16",
				Version: "4.15.5",
			},
		}
		err := w.validateChannel(rosaCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("channel Y-stream '4.16' must match version Y-stream '4.15'"))
	})
}
