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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	utildefaulting "sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
)

var (
	vV1_17_1 = "v1.17.1"
	vV1_17   = "v1.17"
	vV1_16   = "v1.16"
)

func TestDefaultingWebhook(t *testing.T) {
	defaultTestBastion := infrav1.Bastion{
		AllowedCIDRBlocks: []string{"0.0.0.0/0", "::/0"},
	}
	AZUsageLimit := 3
	defaultVPCSpec := infrav1.VPCSpec{
		AvailabilityZoneUsageLimit: &AZUsageLimit,
		AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
		SubnetSchema:               &infrav1.SubnetSchemaPreferPrivate,
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
		spec         ekscontrolplanev1.AWSManagedControlPlaneSpec
		expectSpec   ekscontrolplanev1.AWSManagedControlPlaneSpec
	}{
		{
			name:         "less than 100 chars",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
		{
			name:         "less than 100 chars, dot in name",
			resourceName: "team1.cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "default_team1_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
		{
			name:         "more than 100 chars",
			resourceName: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			resourceNS:   "default",
			expectHash:   true,
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "capi_", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
		{
			name:         "with patch",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         ekscontrolplanev1.AWSManagedControlPlaneSpec{Version: &vV1_17_1},
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", Version: &vV1_17_1, IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
		{
			name:         "with allowed ip on bastion",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				Bastion: infrav1.Bastion{
					AllowedCIDRBlocks: []string{"100.100.100.100/0", "2001:1234:5678:9a40::/56"},
				},
			},
			expectSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				IdentityRef:    defaultIdentityRef,
				Bastion: infrav1.Bastion{
					AllowedCIDRBlocks: []string{"100.100.100.100/0", "2001:1234:5678:9a40::/56"},
				},
				NetworkSpec:                defaultNetworkSpec,
				TokenMethod:                &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator,
				BootstrapSelfManagedAddons: true,
			},
		},
		{
			name:         "with CNI on network",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			spec:         ekscontrolplanev1.AWSManagedControlPlaneSpec{NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}}},
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: infrav1.NetworkSpec{CNI: &infrav1.CNISpec{}, VPC: defaultVPCSpec}, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
		{
			name:         "secondary CIDR",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectHash:   false,
			expectSpec:   ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "default_cluster1", IdentityRef: defaultIdentityRef, Bastion: defaultTestBastion, NetworkSpec: defaultNetworkSpec, SecondaryCidrBlock: nil, TokenMethod: &ekscontrolplanev1.EKSTokenMethodIAMAuthenticator, BootstrapSelfManagedAddons: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			g := NewWithT(t)

			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tc.resourceName,
					Namespace: tc.resourceNS,
				},
			}
			t.Run("for AWSManagedMachinePool", utildefaulting.DefaultValidateTest(context.Background(), mcp, &AWSManagedControlPlane{}))
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
		name                 string
		eksClusterName       string
		expectError          bool
		expectErrorToContain string // if non-empty, the error message must contain this substring
		eksVersion           string
		hasAddons            bool
		vpcCNI               ekscontrolplanev1.VpcCni
		additionalTags       infrav1.Tags
		secondaryCidr        *string
		secondaryCidrBlocks  []infrav1.VpcCidrBlock
		kubeProxy            ekscontrolplanev1.KubeProxy
		accessConfig         *ekscontrolplanev1.AccessConfig
	}{
		{
			name:           "ekscluster specified",
			eksClusterName: "default_cluster1",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
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
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
		},
		{
			name:           "invalid version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.x17",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
		},
		{
			name:           "addons with allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.18",
			expectError:    false,
			hasAddons:      true,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
		},
		{
			name:           "addons with not allowed k8s version",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.17",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
		},
		{
			name:           "disable vpc cni allowed with no addons or secondary cidr",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
		},
		{
			name:           "disable vpc cni not allowed with vpc cni addon",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: true},
		},
		{
			name:           "disable vpc cni allowed with valid secondary",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: true},
			secondaryCidr:  aws.String("100.64.0.0/16"),
		},
		{
			name:           "disable vpc cni allowed with invalid secondary",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: true},
			secondaryCidr:  aws.String("100.64.0.0/10"),
		},
		{
			name:                 "secondary CIDR block not listed in NetworkSpec.VPC.SecondaryCidrBlocks",
			eksClusterName:       "default_cluster1",
			eksVersion:           "v1.19",
			expectError:          true,
			expectErrorToContain: "100.64.0.0/16 must be listed in AWSManagedControlPlane.spec.network.vpc.secondaryCidrBlocks",
			secondaryCidr:        aws.String("100.64.0.0/16"),
			secondaryCidrBlocks:  []infrav1.VpcCidrBlock{{IPv4CidrBlock: "123.456.0.0/16"}},
		},
		{
			name:           "invalid tags not allowed",
			eksClusterName: "default_cluster1",
			expectError:    true,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
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
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
			kubeProxy: ekscontrolplanev1.KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable kube-proxy not allowed with kube-proxy addon",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: false},
			kubeProxy: ekscontrolplanev1.KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable vpc cni and disable kube-proxy allowed with no addons",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			hasAddons:      false,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: true},
			kubeProxy: ekscontrolplanev1.KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "disable vpc cni and disable kube-proxy not allowed with vpc cni and kube-proxy addons",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    true,
			hasAddons:      true,
			vpcCNI:         ekscontrolplanev1.VpcCni{Disable: true},
			kubeProxy: ekscontrolplanev1.KubeProxy{
				Disable: true,
			},
		},
		{
			name:           "BootstrapClusterCreatorAdminPermissions true with EKSAuthenticationModeConfigMap",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode:                      ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				BootstrapClusterCreatorAdminPermissions: ptr.To(true),
			},
		},
		{
			name:                 "BootstrapClusterCreatorAdminPermissions false with EKSAuthenticationModeConfigMap",
			eksClusterName:       "default_cluster1",
			eksVersion:           "v1.19",
			expectError:          true,
			expectErrorToContain: "bootstrapClusterCreatorAdminPermissions must be true if cluster authentication mode is set to config_map",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode:                      ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				BootstrapClusterCreatorAdminPermissions: ptr.To(false),
			},
		},
		{
			name:           "BootstrapClusterCreatorAdminPermissions false with EKSAuthenticationModeAPIAndConfigMap",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode:                      ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
				BootstrapClusterCreatorAdminPermissions: ptr.To(false),
			},
		},
		{
			name:           "BootstrapClusterCreatorAdminPermissions false with EKSAuthenticationModeAPI",
			eksClusterName: "default_cluster1",
			eksVersion:     "v1.19",
			expectError:    false,
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode:                      ekscontrolplanev1.EKSAuthenticationModeAPI,
				BootstrapClusterCreatorAdminPermissions: ptr.To(false),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			g := NewWithT(t)

			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "mcp-",
					Namespace:    "default",
				},
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName: tc.eksClusterName,
					KubeProxy:      tc.kubeProxy,
					AdditionalTags: tc.additionalTags,
					VpcCni:         tc.vpcCNI,
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							SecondaryCidrBlocks: tc.secondaryCidrBlocks,
						},
					},
				},
			}
			if tc.eksVersion != "" {
				mcp.Spec.Version = aws.String(tc.eksVersion)
			}
			if tc.hasAddons {
				testAddons := []ekscontrolplanev1.Addon{
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
			if tc.accessConfig != nil {
				mcp.Spec.AccessConfig = tc.accessConfig
			}

			err := testEnv.Create(ctx, mcp)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())

				if tc.expectErrorToContain != "" && err != nil {
					g.Expect(err.Error()).To(ContainSubstring(tc.expectErrorToContain))
				}
			} else {
				if tc.expectErrorToContain != "" {
					t.Error("Logic error: expectError=false means that expectErrorToContain must be empty")
					t.FailNow()
				}

				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestWebhookCreateIPv6Details(t *testing.T) {
	tests := []struct {
		name        string
		addons      *[]ekscontrolplanev1.Addon
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
			addons: &[]ekscontrolplanev1.Addon{
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
			addons: &[]ekscontrolplanev1.Addon{
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
		{
			name:        "both ipv6 poolId and ipamPool are set",
			kubeVersion: "v1.22",
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{
						PoolID:   "not-empty",
						IPAMPool: &infrav1.IPAMPool{},
					},
				},
			},
			err: "poolId and ipamPool cannot be used together",
		},
		{
			name:        "both ipv6 cidrBlock and ipamPool are set",
			kubeVersion: "v1.22",
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{
						CidrBlock: "not-empty",
						IPAMPool:  &infrav1.IPAMPool{},
					},
				},
			},
			err: "cidrBlock and ipamPool cannot be used together",
		},
		{
			name:        "Id or name are not set for IPAMPool",
			kubeVersion: "v1.22",
			networkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					IPv6: &infrav1.IPv6{
						IPAMPool: &infrav1.IPAMPool{},
					},
				},
			},
			err: "ipamPool must have either id or name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			g := NewWithT(t)

			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "mcp-",
					Namespace:    "default",
				},
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName: "test-cluster",
					Addons:         tc.addons,
					NetworkSpec:    tc.networkSpec,
					Version:        aws.String(tc.kubeVersion),
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
		oldClusterSpec ekscontrolplanev1.AWSManagedControlPlaneSpec
		newClusterSpec ekscontrolplanev1.AWSManagedControlPlaneSpec
		oldClusterName string
		newClusterName string
		oldEksVersion  string
		newEksVersion  string
		expectError    bool
	}{
		{
			name: "ekscluster specified, same cluster names",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			expectError: false,
		},
		{
			name: "ekscluster specified, different cluster names",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster2",
			},
			expectError: true,
		},
		{
			name: "old ekscluster specified, no new cluster name",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "",
			},
			expectError: true,
		},
		{
			name: "older version",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_16,
			},
			expectError: true,
		},
		{
			name: "same version",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			expectError: false,
		},
		{
			name: "newer version",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_16,
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				Version:        &vV1_17,
			},
			expectError: false,
		},
		{
			name: "no change in access config",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			expectError: false,
		},
		{
			name: "change in access config to nil",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			expectError: true,
		},
		{
			name: "change in access config from nil to valid",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			expectError: false,
		},
		{
			name: "change in access config auth mode from ApiAndConfigMap to API is allowed",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
				},
			},
			expectError: false,
		},
		{
			name: "change in access config auth mode from API to Config Map is denied",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			expectError: true,
		},
		{
			name: "change in access config auth mode from APIAndConfigMap to Config Map is denied",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
				},
			},
			expectError: true,
		},
		{
			name: "change in access config bootstrap admin permissions is ignored",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					BootstrapClusterCreatorAdminPermissions: ptr.To(true),
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				AccessConfig: &ekscontrolplanev1.AccessConfig{
					BootstrapClusterCreatorAdminPermissions: ptr.To(false),
				},
			},
			expectError: false,
		},
		{
			name: "change in encryption config to nil",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider:  ptr.To[string]("provider"),
					Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			expectError: true,
		},
		{
			name: "change in encryption config from nil to valid encryption-config",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider:  ptr.To[string]("provider"),
					Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
				},
			},
			expectError: false,
		},
		{
			name: "change in provider of encryption config",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider:  ptr.To[string]("provider"),
					Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider:  ptr.To[string]("new-provider"),
					Resources: []*string{ptr.To[string]("foo"), ptr.To[string]("bar")},
				},
			},
			expectError: true,
		},
		{
			name: "no change in provider of encryption config",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider: ptr.To[string]("provider"),
				},
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				EncryptionConfig: &ekscontrolplanev1.EncryptionConfig{
					Provider: ptr.To[string]("provider"),
				},
			},
			expectError: false,
		},
		{
			name: "ekscluster specified, same name, invalid tags",
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
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
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{},
				},
				Version: ptr.To[string]("1.22"),
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
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
			oldClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						IPv6: &infrav1.IPv6{},
					},
				},
				Addons: &[]ekscontrolplanev1.Addon{
					{
						Name:    vpcCniAddon,
						Version: "1.11.0",
					},
				},
				Version: ptr.To[string]("v1.22.0"),
			},
			newClusterSpec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				EKSClusterName: "default_cluster1",
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{},
				},
				Version: ptr.To[string]("v1.22.0"),
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ctx := context.TODO()
			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
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

			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName: "default_cluster1",
				},
			}
			if tc.cidrRange != "" {
				mcp.Spec.SecondaryCidrBlock = aws.String(tc.cidrRange)
			}

			warn, err := (&AWSManagedControlPlane{}).ValidateCreate(context.Background(), mcp)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
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

			newMCP := &ekscontrolplanev1.AWSManagedControlPlane{
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName:     "default_cluster1",
					SecondaryCidrBlock: aws.String(tc.cidrRange),
				},
			}
			oldMCP := &ekscontrolplanev1.AWSManagedControlPlane{
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName:     "default_cluster1",
					SecondaryCidrBlock: nil,
				},
			}

			warn, err := (&AWSManagedControlPlane{}).ValidateUpdate(context.Background(), oldMCP, newMCP)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}

func TestWebhookValidateAccessEntries(t *testing.T) {
	tests := []struct {
		name          string
		accessConfig  *ekscontrolplanev1.AccessConfig
		accessEntries []ekscontrolplanev1.AccessEntry
		expectError   bool
		errorSubstr   string
	}{
		{
			name: "valid access entries with api auth mode",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:             ekscontrolplanev1.AccessEntryTypeStandard,
					KubernetesGroups: []string{"system:masters"},
				},
			},
			expectError: false,
		},
		{
			name: "valid access entries with api_and_config_map auth mode",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:             ekscontrolplanev1.AccessEntryTypeStandard,
					KubernetesGroups: []string{"system:masters"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid access entries with config_map auth mode",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeConfigMap,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:             ekscontrolplanev1.AccessEntryTypeStandard,
					KubernetesGroups: []string{"system:masters"},
				},
			},
			expectError: true,
			errorSubstr: "accessEntries can only be used when authenticationMode is set to api or api_and_config_map",
		},
		{
			name: "invalid ec2_linux access entry with kubernetes groups",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:             ekscontrolplanev1.AccessEntryTypeEC2Linux,
					KubernetesGroups: []string{"system:masters"},
				},
			},
			expectError: true,
			errorSubstr: "kubernetesGroups cannot be specified when type is ec2_linux or ec2_windows",
		},
		{
			name: "invalid ec2_windows access entry with access policies",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN: "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:         ekscontrolplanev1.AccessEntryTypeEC2Windows,
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy",
							AccessScope: ekscontrolplanev1.AccessScope{
								Type: ekscontrolplanev1.AccessScopeTypeCluster,
							},
						},
					},
				},
			},
			expectError: true,
			errorSubstr: "accessPolicies cannot be specified when type is ec2_linux or ec2_windows",
		},
		{
			name: "invalid access policy with namespace type and no namespaces",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN: "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:         ekscontrolplanev1.AccessEntryTypeStandard,
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy",
							AccessScope: ekscontrolplanev1.AccessScope{
								Type: ekscontrolplanev1.AccessScopeTypeNamespace,
							},
						},
					},
				},
			},
			expectError: true,
			errorSubstr: "at least one value must be provided when accessScope type is namespace",
		},
		{
			name: "valid access policy with namespace type and namespaces",
			accessConfig: &ekscontrolplanev1.AccessConfig{
				AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPI,
			},
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN: "arn:aws:iam::123456789012:role/EKSAdmin",
					Type:         ekscontrolplanev1.AccessEntryTypeStandard,
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: "arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy",
							AccessScope: ekscontrolplanev1.AccessScope{
								Type:       ekscontrolplanev1.AccessScopeTypeNamespace,
								Namespaces: []string{"default", "kube-system"},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &ekscontrolplanev1.AWSManagedControlPlane{
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName: "default_cluster1",
					AccessConfig:   tc.accessConfig,
					AccessEntries:  tc.accessEntries,
				},
			}

			warn, err := (&AWSManagedControlPlane{}).ValidateCreate(context.Background(), mcp)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				if tc.errorSubstr != "" {
					g.Expect(err.Error()).To(ContainSubstring(tc.errorSubstr))
				}
			} else {
				g.Expect(err).To(BeNil())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}
