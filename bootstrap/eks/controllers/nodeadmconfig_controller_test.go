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

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestNodeadmServiceCIDR(t *testing.T) {
	clusterWithServices := func(cidrs ...string) *clusterv1.Cluster {
		c := &clusterv1.Cluster{}
		c.Spec.ClusterNetwork.Services.CIDRBlocks = cidrs
		return c
	}
	controlPlaneWithStatus := func(cidr *string) *ekscontrolplanev1.AWSManagedControlPlane {
		cp := &ekscontrolplanev1.AWSManagedControlPlane{}
		cp.Status.ServiceCIDR = cidr
		return cp
	}

	testCases := []struct {
		name         string
		cluster      *clusterv1.Cluster
		controlPlane *ekscontrolplanev1.AWSManagedControlPlane
		expect       string
	}{
		{
			name:         "nothing set falls through to the template default",
			cluster:      clusterWithServices(),
			controlPlane: controlPlaneWithStatus(nil),
			expect:       "",
		},
		{
			name:         "cluster services CIDR when EKS hasn't reported one",
			cluster:      clusterWithServices("10.96.0.0/12"),
			controlPlane: controlPlaneWithStatus(nil),
			expect:       "10.96.0.0/12",
		},
		{
			name:         "EKS-assigned IPv6 CIDR wins over an IPv4 cluster services CIDR",
			cluster:      clusterWithServices("10.96.0.0/12"),
			controlPlane: controlPlaneWithStatus(ptr.To("fdbb:74f2:123e::/108")),
			expect:       "fdbb:74f2:123e::/108",
		},
		{
			name:         "EKS-assigned IPv6 CIDR when the cluster has no services CIDR",
			cluster:      clusterWithServices(),
			controlPlane: controlPlaneWithStatus(ptr.To("fdbb:74f2:123e::/108")),
			expect:       "fdbb:74f2:123e::/108",
		},
		{
			name:         "empty status value is ignored",
			cluster:      clusterWithServices("10.100.0.0/16"),
			controlPlane: controlPlaneWithStatus(ptr.To("")),
			expect:       "10.100.0.0/16",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(nodeadmServiceCIDR(tc.cluster, tc.controlPlane)).To(Equal(tc.expect))
		})
	}
}
