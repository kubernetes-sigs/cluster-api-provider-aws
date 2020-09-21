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

package v1alpha3

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

var (
	vV1_17_1 = "v1.17.1"
	vV1_17   = "v1.17"
	vV1_16   = "v1.16"
)

func TestDefaultingWebhook(t *testing.T) {
	defaultTestBastion := infrav1.Bastion{
		AllowedCIDRBlocks: []string{"0.0.0.0/0"},
	}
	defaultNetworkSpec := infrav1.NetworkSpec{
		CNI: &infrav1.CNISpec{
			CNIIngressRules: []*infrav1.CNIIngressRule{
				{
					Description: "bgp (calico)",
					Protocol:    "tcp",
					FromPort:    179,
					ToPort:      179,
				},
				{
					Description: "IP-in-IP (calico)",
					Protocol:    "4",
					FromPort:    -1,
					ToPort:      65535,
				},
			},
		},
	}

	tests := []struct {
		name         string
		resourceName string
		resourceNS   string
		expectHash   bool
		expect       string
		spec         AWSManagedControlPlaneSpec
		expectSpec   AWSManagedControlPlaneSpec
	}{
		{
			name:         "less than 100 chars",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec},
		},
		{
			name:         "less than 100 chars, dot in name",
			resourceName: "team1.cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_team1_cluster1", Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec},
		},
		{
			name:         "more than 100 chars",
			resourceName: "ABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDE",
			resourceNS:   "default",
			expectHash:   true,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "capi_", Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec},
		},
		{
			name:         "with patch",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{Version: &vV1_17_1},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Version: &vV1_17, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec},
		},
		{
			name:         "with allowed ip on bastion",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{Bastion: infrav1.Bastion{AllowedCIDRBlocks: []string{"100.100.100.100/0"}}},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Bastion: infrav1.Bastion{AllowedCIDRBlocks: []string{"100.100.100.100/0"}}, NetworkSpec: defaultNetworkSpec},
		},
		{
			name:         "with CNI on network",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}}},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Bastion: defaultTestBastion, NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				ObjectMeta: v1.ObjectMeta{
					Name:      tc.resourceName,
					Namespace: tc.resourceNS,
				},
			}
			mcp.Spec = tc.spec
			mcp.Default()

			g.Expect(mcp.Spec.EKSClusterName).ToNot(BeEmpty())

			if tc.expectHash {
				g.Expect(strings.HasPrefix(mcp.Spec.EKSClusterName, "capa_")).To(BeTrue())
				// We don't care about the exact name
				tc.expectSpec.EKSClusterName = mcp.Spec.EKSClusterName
			}
			g.Expect(mcp.Spec).To(Equal(tc.expectSpec))
		})
	}
}

func TestValidatingWebhookCreate(t *testing.T) {
	tests := []struct {
		name           string
		eksClusterName string
		expectError    bool
		eksVersion     string
	}{
		{
			name:           "ekscluster specified",
			eksClusterName: "default_cluster1",
			expectError:    false,
		},
		{
			name:           "ekscluster NOT specified",
			eksClusterName: "",
			expectError:    true,
		},
		{
			name:           "invalid version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.x17",
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: tc.eksClusterName,
				},
			}
			if tc.eksVersion != "" {
				mcp.Spec.Version = &tc.eksVersion
			}
			err := mcp.ValidateCreate()

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestValidatingWebhookUpdate(t *testing.T) {
	tests := []struct {
		name           string
		oldClusterSpec AWSManagedControlPlaneSpec
		newClusterSpec AWSManagedControlPlaneSpec
		oldClusterName string
		newClusterName string
		oldEksVersion  string
		newEksVersion  string
		expectError    bool
	}{
		{
			name: "ekscluster specified, same cluster names",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			expectError: false,
		},
		{
			name: "ekscluster specified, different cluster names",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster2",
			},
			expectError: true,
		},
		{
			name: "old ekscluster specified, no new cluster name",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "",
			},
			expectError: true,
		},
		{
			name: "older version",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_16,
			},
			expectError: true,
		},
		{
			name: "same version",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			expectError: false,
		},
		{
			name: "newer version",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_16,
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			newMCP := &AWSManagedControlPlane{
				Spec: tc.newClusterSpec,
			}
			oldMCP := &AWSManagedControlPlane{
				Spec: tc.oldClusterSpec,
			}
			err := newMCP.ValidateUpdate(oldMCP)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
