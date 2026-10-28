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

package controllers

import (
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

func TestEKSConfigServiceIPv6CIDR(t *testing.T) {
	const (
		vpcIPv6Block   = "2600:1f14:3356:2100::/56"
		eksServiceCIDR = "fdbb:74f2:123e::/108"
		configuredCIDR = "fd00:1234::/108"
	)
	controlPlane := func(ipv6 bool, statusCIDR *string) *ekscontrolplanev1.AWSManagedControlPlane {
		cp := &ekscontrolplanev1.AWSManagedControlPlane{}
		if ipv6 {
			cp.Spec.NetworkSpec.VPC.IPv6 = &infrav1.IPv6{CidrBlock: vpcIPv6Block}
		}
		cp.Status.ServiceCIDR = statusCIDR
		return cp
	}
	config := func(cidr *string) *eksbootstrapv1.EKSConfig {
		c := &eksbootstrapv1.EKSConfig{}
		c.Spec.ServiceIPV6Cidr = cidr
		return c
	}

	cases := []struct {
		name         string
		config       *eksbootstrapv1.EKSConfig
		controlPlane *ekscontrolplanev1.AWSManagedControlPlane
		expect       *string
	}{
		{
			name:         "IPv4 cluster renders no IPv6 settings",
			config:       config(nil),
			controlPlane: controlPlane(false, ptr.To("10.100.0.0/16")),
			expect:       nil,
		},
		{
			name:         "IPv6 cluster uses the service CIDR EKS reported",
			config:       config(nil),
			controlPlane: controlPlane(true, ptr.To(eksServiceCIDR)),
			expect:       ptr.To(eksServiceCIDR),
		},
		{
			name:         "IPv6 cluster falls back to the VPC block until EKS reports the CIDR",
			config:       config(nil),
			controlPlane: controlPlane(true, nil),
			expect:       ptr.To(vpcIPv6Block),
		},
		{
			name:         "CIDR set on the EKSConfig wins",
			config:       config(ptr.To(configuredCIDR)),
			controlPlane: controlPlane(true, ptr.To(eksServiceCIDR)),
			expect:       ptr.To(configuredCIDR),
		},
		{
			name:         "empty CIDR set on the EKSConfig renders no IPv6 settings",
			config:       config(ptr.To("")),
			controlPlane: controlPlane(true, ptr.To(eksServiceCIDR)),
			expect:       nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(eksConfigServiceIPv6CIDR(tc.config, tc.controlPlane)).To(Equal(tc.expect))
		})
	}
}
