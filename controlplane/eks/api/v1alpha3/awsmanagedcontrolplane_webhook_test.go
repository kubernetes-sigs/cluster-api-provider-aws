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
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	AZUsageLimit := 3
	defaultVPCSpec := infrav1.VPCSpec{
		AvailabilityZoneUsageLimit: &AZUsageLimit,
		AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
	}
	defaultIdentityRef := &infrav1.AWSIdentityReference{
		Kind: infrav1.ControllerIdentityKind,
		Name: infrav1.AWSClusterControllerIdentityName,
	}
	defaultNetworkSpec := infrav1.NetworkSpec{
		VPC: defaultVPCSpec,
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
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "less than 100 chars, dot in name",
			resourceName: "team1.cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_team1_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "more than 100 chars",
			resourceName: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			resourceNS:   "default",
			expectHash:   true,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "capi_", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "with patch",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{Version: &vV1_17_1},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Version: &vV1_17, IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "with allowed ip on bastion",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{Bastion: infrav1.Bastion{AllowedCIDRBlocks: []string{"100.100.100.100/0"}}},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: infrav1.Bastion{AllowedCIDRBlocks: []string{"100.100.100.100/0"}}, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "with CNI on network",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         AWSManagedControlPlaneSpec{NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}}},
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}, VPC: defaultVPCSpec}, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
		{
			name:         "secondary CIDR",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, SecondaryCidrBlock: nil, TokenMethod: &EKSTokenMethodIAMAuthenticator},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tc.resourceName,
					Namespace: tc.resourceNS,
				},
			}
			mcp.Spec = tc.spec

			g.Expect(testEnv.Create(ctx, mcp)).To(Succeed())

			defer func() {
				testEnv.Delete(ctx, mcp)
			}()

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

func TestWebhookCreate(t *testing.T) {
	tests := []struct {
		name           string
		eksClusterName string
		expectError    bool
		eksVersion     string
		hasAddon       bool
		disableVPCCNI  bool
		secondaryCidr  *string
	}{
		{
			name:           "ekscluster specified",
			eksClusterName: "default_cluster1",
			expectError:    false,
			hasAddon:       false,
			disableVPCCNI:  false,
		},
		{
			name:           "ekscluster NOT specified",
			eksClusterName: "",
			expectError:    false,
			hasAddon:       false,
			disableVPCCNI:  false,
		},
		{
			name:           "invalid version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.x17",
			expectError:    true,
			hasAddon:       false,
			disableVPCCNI:  false,
		},
		{
			name:           "addon with allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.18",
			expectError:    false,
			hasAddon:       true,
			disableVPCCNI:  false,
		},
		{
			name:           "addon with not allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.17",
			expectError:    true,
			hasAddon:       true,
			disableVPCCNI:  false,
		},
		{
			name:           "disable vpc cni allowed with no addon or secondary cidr",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddon:       false,
			disableVPCCNI:  true,
		},
		{
			name:           "disable vpc cni not allowed with vpc cni addon",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddon:       true,
			disableVPCCNI:  true,
		},
		{
			name:           "disable vpc cni not allowed with secondary",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddon:       false,
			disableVPCCNI:  true,
			secondaryCidr:  aws.String("100.64.0.0/10"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "mcp-",
					Namespace:    "default",
				},
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: tc.eksClusterName,
					DisableVPCCNI:  tc.disableVPCCNI,
				},
			}
			if tc.eksVersion != "" {
				mcp.Spec.Version = &tc.eksVersion
			}
			if tc.hasAddon {
				testAddons := []Addon{
					{
						Name:    vpcCniAddon,
						Version: "v1.0.0",
					},
				}
				mcp.Spec.Addons = &testAddons
			}
			if tc.secondaryCidr != nil {
				mcp.Spec.SecondaryCidrBlock = tc.secondaryCidr
			}

			err := testEnv.Create(ctx, mcp)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestWebhookUpdate(t *testing.T) {
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
			ctx := context.TODO()
			mcp := &AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "mcp-",
					Namespace:    "default",
				},
				Spec: tc.oldClusterSpec,
			}
			g.Expect(testEnv.Create(ctx, mcp)).To(Succeed())
			mcp.Spec = tc.newClusterSpec
			err := testEnv.Update(ctx, mcp)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestValidatingWebhookCreate_SecondaryCidr(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
		cidrRange   string
	}{
		{
			name:        "complete range 1",
			cidrRange:   "100.64.0.0/10",
			expectError: true,
		},
		{
			name:        "complete range 2",
			cidrRange:   "198.19.0.0/16",
			expectError: false,
		},
		{
			name:        "subrange",
			cidrRange:   "100.67.0.0/16",
			expectError: false,
		},
		{
			name:        "invalid value",
			cidrRange:   "not a cidr range",
			expectError: true,
		},
		{
			name:        "unsupported range",
			cidrRange:   "10.0.0.1/20",
			expectError: true,
		},
		{
			name:        "too large",
			cidrRange:   "100.64.0.0/15",
			expectError: true,
		},
		{
			name:        "too small",
			cidrRange:   "100.64.0.0/29",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: "default_cluster1",
				},
			}
			if tc.cidrRange != "" {
				mcp.Spec.SecondaryCidrBlock = &tc.cidrRange
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

func TestValidatingWebhookUpdate_SecondaryCidr(t *testing.T) {
	tests := []struct {
		name        string
		cidrRange   string
		expectError bool
	}{
		{
			name:        "complete range 1",
			cidrRange:   "100.64.0.0/10",
			expectError: true,
		},
		{
			name:        "complete range 2",
			cidrRange:   "198.19.0.0/16",
			expectError: false,
		},
		{
			name:        "subrange",
			cidrRange:   "100.67.0.0/16",
			expectError: false,
		},
		{
			name:        "invalid value",
			cidrRange:   "not a cidr range",
			expectError: true,
		},
		{
			name:        "unsupported range",
			cidrRange:   "10.0.0.1/20",
			expectError: true,
		},
		{
			name:        "too large",
			cidrRange:   "100.64.0.0/15",
			expectError: true,
		},
		{
			name:        "too small",
			cidrRange:   "100.64.0.0/29",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			newMCP := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName:     "default_cluster1",
					SecondaryCidrBlock: &tc.cidrRange,
				},
			}
			oldMCP := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName:     "default_cluster1",
					SecondaryCidrBlock: nil,
				},
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
