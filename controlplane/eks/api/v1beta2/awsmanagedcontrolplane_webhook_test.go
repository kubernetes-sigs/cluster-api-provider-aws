/*
Copyright 2022 The Kubernetes Authors.

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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api/util/defaulting"
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
			CNIIngressRules: []infrav1.CNIIngressRule{
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
			expectSpec:   AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Version: &vV1_17_1, IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &EKSTokenMethodIAMAuthenticator},
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
			t.Run("for AWSManagedMachinePool", utildefaulting.DefaultValidateTest(mcp))
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
	tests := []struct { //nolint:maligned
		name           string
		eksClusterName string
		expectError    bool
		eksVersion     string
		hasAddons      bool
		vpcCNI         VpcCni
		additionalTags infrav1.Tags
		secondaryCidr  *string
		kubeProxy      KubeProxy
	}{
		{
			name:           "ekscluster specified",
			eksClusterName: "default_cluster1",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
			additionalTags: infrav1.Tags{
				"a":     "b",
				"key-2": "value-2",
			},
		},
		{
			name:           "ekscluster NOT specified",
			eksClusterName: "",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
		},
		{
			name:           "invalid version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.x17",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
		},
		{
			name:           "addons with allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.18",
			expectError:    false,
			hasAddons:      true,
			vpcCNI:         VpcCni{Disable: false},
		},
		{
			name:           "addons with not allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.17",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         VpcCni{Disable: false},
		},
		{
			name:           "disable vpc cni allowed with no addons or secondary cidr",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
		},
		{
			name:           "disable vpc cni not allowed with vpc cni addon",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         VpcCni{Disable: true},
		},
		{
			name:           "disable vpc cni allowed with valid secondary",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: true},
			secondaryCidr:  aws.String("100.64.0.0/16"),
		},
		{
			name:           "disable vpc cni allowed with invalid secondary",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: true},
			secondaryCidr:  aws.String("100.64.0.0/10"),
		},
		{
			name:           "invalid tags not allowed",
			eksClusterName: "default_cluster1",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
			additionalTags: infrav1.Tags{
				"key-1":                    "value-1",
				"":                         "value-2",
				strings.Repeat("CAPI", 33): "value-3",
				"key-4":                    strings.Repeat("CAPI", 65),
			},
		},
		{
			name:           "disable kube-proxy allowed with no addons",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: false},
			kubeProxy: KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable kube-proxy not allowed with kube-proxy addon",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         VpcCni{Disable: false},
			kubeProxy: KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable vpc cni and disable kube-proxy allowed with no addons",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         VpcCni{Disable: true},
			kubeProxy: KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable vpc cni and disable kube-proxy not allowed with vpc cni and kube-proxy addons",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         VpcCni{Disable: true},
			kubeProxy: KubeProxy{
				Disable: true,
			},
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
					KubeProxy:      tc.kubeProxy,
					AdditionalTags: tc.additionalTags,
					VpcCni:         tc.vpcCNI,
				},
			}
			if tc.eksVersion != "" {
				mcp.Spec.Version = &tc.eksVersion
			}
			if tc.hasAddons {
				testAddons := []Addon{
					{
						Name:    vpcCniAddon,
						Version: "v1.0.0",
					},
					{
						Name:    kubeProxyAddon,
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

func TestWebhookCreateIPv6Details(t *testing.T) {
	tests := []struct {
		name        string
		addons      []Addon
		kubeVersion string
		networkSpec infrav1.NetworkSpec
		err         string
	}{
		{
			name:        "ipv6 with lower cluster version",
			kubeVersion: "v1.18",
			err:         fmt.Sprintf("IPv6 requires Kubernetes %s or greater", minKubeVersionForIPv6),
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
				},
			},
		},
		{
			name:        "ipv6 no addons",
			kubeVersion: "v1.22",
			err:         "addons are required to be set explicitly if IPv6 is enabled",
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
				},
			},
		},
		{
			name:        "ipv6 with addons but cni version is lower than supported version",
			kubeVersion: "v1.22",
			addons: []Addon{
				{
					Name:    vpcCniAddon,
					Version: "1.9.3",
				},
			},
			err: fmt.Sprintf("vpc-cni version must be above or equal to %s for IPv6", minVpcCniVersionForIPv6),
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
				},
			},
		},
		{
			name:        "ipv6 with addons and correct cni and cluster version",
			kubeVersion: "v1.22",
			addons: []Addon{
				{
					Name:    vpcCniAddon,
					Version: "1.11.0",
				},
			},
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{},
				},
			},
		},
		{
			name:        "ipv6 cidr block is set but pool is left empty",
			kubeVersion: "v1.18",
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{
						CidrBlock: "not-empty",
						// PoolID is empty
					},
				},
			},
			err: "poolId cannot be empty if cidrBlock is set",
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
					EKSClusterName: "test-cluster",
					Addons:         &tc.addons,
					NetworkSpec:    tc.networkSpec,
					Version:        &tc.kubeVersion,
				},
			}
			err := testEnv.Create(ctx, mcp)

			if tc.err != "" {
				g.Expect(err).To(MatchError(ContainSubstring(tc.err)))
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
		{
			name: "change in encryption config to nil",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider:  pointer.String("provider"),
					Resources: []*string{pointer.String("foo"), pointer.String("bar")},
				},
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			expectError: true,
		},
		{
			name: "change in encryption config from nil to valid encryption-config",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider:  pointer.String("provider"),
					Resources: []*string{pointer.String("foo"), pointer.String("bar")},
				},
			},
			expectError: false,
		},
		{
			name: "change in provider of encryption config",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider:  pointer.String("provider"),
					Resources: []*string{pointer.String("foo"), pointer.String("bar")},
				},
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider:  pointer.String("new-provider"),
					Resources: []*string{pointer.String("foo"), pointer.String("bar")},
				},
			},
			expectError: true,
		},
		{
			name: "no change in provider of encryption config",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider: pointer.String("provider"),
				},
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &EncryptionConfig{
					Provider: pointer.String("provider"),
				},
			},
			expectError: false,
		},
		{
			name: "ekscluster specified, same name, invalid tags",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AdditionalTags: infrav1.Tags{
					"key-1":                    "value-1",
					"":                         "value-2",
					strings.Repeat("CAPI", 33): "value-3",
					"key-4":                    strings.Repeat("CAPI", 65),
				},
			},
			expectError: true,
		},
		{
			name: "changing ipv6 enabled is not allowed after it has been set - false, true",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{},
				},
				Version: pointer.String("1.22"),
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						IPv6: &infrav1.IPv6{},
					},
				},
			},
			expectError: true,
		},
		{
			name: "changing ipv6 enabled is not allowed after it has been set - true, false",
			oldClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						IPv6: &infrav1.IPv6{},
					},
				},
				Addons: &[]Addon{
					{
						Name:    vpcCniAddon,
						Version: "1.11.0",
					},
				},
				Version: pointer.String("v1.22.0"),
			},
			newClusterSpec: AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{},
				},
				Version: pointer.String("v1.22.0"),
			},
			expectError: true,
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

func TestValidatingWebhookCreateSecondaryCidr(t *testing.T) {
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

func TestValidatingWebhookUpdateSecondaryCidr(t *testing.T) {
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
