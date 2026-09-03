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
	"k8s.io/apimachinery/pkg/util/validation/field"

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

func TestValidateComponentRoutes(t *testing.T) {
	w := &ROSAControlPlane{}

	t.Run("valid console and downloads", func(t *testing.T) {
		g := NewGomegaWithT(t)
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				//nolint:gosec // G101: test fixture, not real credentials
				ComponentRoutes: []rosacontrolplanev1.ComponentRouteSpec{
					{Name: rosacontrolplanev1.ComponentRouteConsole, Hostname: "console.example.com", TLSSecretRef: "console-tls"},
					{Name: rosacontrolplanev1.ComponentRouteDownloads, Hostname: "downloads.example.com", TLSSecretRef: "downloads-tls"},
				},
			},
		}
		errs := w.validateComponentRoutes(rosaCP)
		g.Expect(errs).To(BeEmpty())
	})

	t.Run("valid console only", func(t *testing.T) {
		g := NewGomegaWithT(t)
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ComponentRoutes: []rosacontrolplanev1.ComponentRouteSpec{
					{Name: rosacontrolplanev1.ComponentRouteConsole, Hostname: "console.example.com", TLSSecretRef: "console-tls"},
				},
			},
		}
		errs := w.validateComponentRoutes(rosaCP)
		g.Expect(errs).To(BeEmpty())
	})

	t.Run("no componentRoutes", func(t *testing.T) {
		g := NewGomegaWithT(t)
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{},
		}
		errs := w.validateComponentRoutes(rosaCP)
		g.Expect(errs).To(BeEmpty())
	})

	t.Run("duplicate console key", func(t *testing.T) {
		g := NewGomegaWithT(t)
		rosaCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				ComponentRoutes: []rosacontrolplanev1.ComponentRouteSpec{
					{Name: rosacontrolplanev1.ComponentRouteConsole, Hostname: "console1.example.com", TLSSecretRef: "tls1"},
					{Name: rosacontrolplanev1.ComponentRouteConsole, Hostname: "console2.example.com", TLSSecretRef: "tls2"},
				},
			},
		}
		errs := w.validateComponentRoutes(rosaCP)
		g.Expect(errs).ToNot(BeEmpty())
		g.Expect(errs[0].Type).To(Equal(field.ErrorTypeDuplicate))
	})
}

func TestValidateEc2MetadataHttpTokensImmutability(t *testing.T) {
	w := &ROSAControlPlane{}

	t.Run("Update blocked when changing required to optional", func(t *testing.T) {
		g := NewGomegaWithT(t)
		oldCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensOptional,
			},
		}
		err := w.validateEc2MetadataHTTPTokensImmutability(oldCP, newCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens is immutable"))
	})

	t.Run("Update blocked when adding field to object that was created without it", func(t *testing.T) {
		g := NewGomegaWithT(t)
		oldCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{},
		}
		newCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		err := w.validateEc2MetadataHTTPTokensImmutability(oldCP, newCP)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens is immutable"))
	})

	t.Run("Update allowed when value unchanged", func(t *testing.T) {
		g := NewGomegaWithT(t)
		oldCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		err := w.validateEc2MetadataHTTPTokensImmutability(oldCP, newCP)
		g.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("Update allowed when both old and new have no value", func(t *testing.T) {
		g := NewGomegaWithT(t)
		oldCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{},
		}
		newCP := &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{},
		}
		err := w.validateEc2MetadataHTTPTokensImmutability(oldCP, newCP)
		g.Expect(err).NotTo(HaveOccurred())
	})
}
