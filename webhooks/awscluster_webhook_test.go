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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilfeature "k8s.io/component-base/featuregate/testing"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/defaulting"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

func TestAWSClusterDefault(t *testing.T) {
	cluster := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	t.Run("for AWSCluster", defaultValidateTest(context.Background(), cluster, &AWSCluster{}, true))
	cluster.Default()
	g := NewWithT(t)
	g.Expect(cluster.Spec.IdentityRef).NotTo(BeNil())
}

func TestAWSClusterValidateCreate(t *testing.T) {
	unsupportedIncorrectScheme := infrav1.ELBScheme("any-other-scheme")

	tests := []struct {
		name    string
		cluster *infrav1.AWSCluster
		wantErr bool
		expect  func(g *WithT, res *infrav1.AWSLoadBalancerSpec)
	}{
		{
			name: "No options are allowed when LoadBalancer is disabled (name)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
						Name:             ptr.To("name"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (crossZoneLoadBalancing)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						CrossZoneLoadBalancing: true,
						LoadBalancerType:       infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (subnets)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Subnets:          []string{"foo", "bar"},
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (healthCheckProtocol)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
						LoadBalancerType:    infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (additionalSecurityGroups)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						AdditionalSecurityGroups: []string{"foo", "bar"},
						LoadBalancerType:         infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (additionalListeners)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						AdditionalListeners: []infrav1.AdditionalListenerSpec{
							{
								Port:     6443,
								Protocol: infrav1.ELBProtocolTCP,
							},
						},
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (ingressRules)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Description: "ingress rule",
								Protocol:    infrav1.SecurityGroupProtocolTCP,
								FromPort:    6443,
								ToPort:      6443,
							},
						},
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (disableHostsRewrite)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						DisableHostsRewrite: true,
						LoadBalancerType:    infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No options are allowed when LoadBalancer is disabled (preserveClientIP)",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						PreserveClientIP: true,
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		// The SSHKeyName tests were moved to sshkeyname_test.go
		{
			name: "Supported schemes are 'internet-facing, Internet-facing, internal, or nil', rest will be rejected",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{Scheme: &unsupportedIncorrectScheme},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid tags are rejected",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts bucket name with acceptable characters",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name:                           "abcdefghijklmnoprstuwxyz-0123456789",
						ControlPlaneIAMInstanceProfile: "control-plane.cluster-api-provider-aws.sigs.k8s.io",
						NodesIAMInstanceProfiles:       []string{"nodes.cluster-api-provider-aws.sigs.k8s.io"},
					},
				},
			},
		},
		{
			name: "rejects empty bucket name",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects bucket name shorter than 3 characters",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name: "fo",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects bucket name longer than 63 characters",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name: strings.Repeat("a", 64),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects bucket name starting with not letter or number",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name: "-foo",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects bucket name ending with not letter or number",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name: "foo-",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects bucket name formatted as IP address",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name: "8.8.8.8",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "requires bucket control plane IAM instance profile to be not empty",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name:                           "foo",
						ControlPlaneIAMInstanceProfile: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "requires at least one bucket node IAM instance profile",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name:                           "foo",
						ControlPlaneIAMInstanceProfile: "foo",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "requires all bucket node IAM instance profiles to be not empty",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name:                           "foo",
						ControlPlaneIAMInstanceProfile: "foo",
						NodesIAMInstanceProfiles:       []string{""},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "does not return error when all IAM instance profiles are populated",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					S3Bucket: &infrav1.S3Bucket{
						Name:                           "foo",
						ControlPlaneIAMInstanceProfile: "foo",
						NodesIAMInstanceProfiles:       []string{"bar"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts vpc cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects invalid vpc cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts vpc secondary cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
							SecondaryCidrBlocks: []infrav1.VpcCidrBlock{
								{
									IPv4CidrBlock: "10.0.1.0/24",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects invalid vpc secondary cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
							SecondaryCidrBlocks: []infrav1.VpcCidrBlock{
								{
									IPv4CidrBlock: "10.0.1.0",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects vpc secondary cidr duplicate with vpc primary cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
							SecondaryCidrBlocks: []infrav1.VpcCidrBlock{
								{
									IPv4CidrBlock: "10.0.0.0/16",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts vpc ipv6 cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPv6: &infrav1.IPv6{
								CidrBlock: "2001:2345:5678::/64",
								PoolID:    "pool-id",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "reject invalid vpc ipv6 cidr",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPv6: &infrav1.IPv6{
								CidrBlock: "2001:2345:5678::",
								PoolID:    "pool-id",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts ipv6 enabled subnet",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: []infrav1.SubnetSpec{
							{
								ID:     "sub-1",
								IsIPv6: true,
							},
							{
								ID: "sub-2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts cidr block for subnets",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: []infrav1.SubnetSpec{
							{
								ID:        "sub-1",
								CidrBlock: "10.0.10.0/24",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects invalid cidr block for subnets",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: []infrav1.SubnetSpec{
							{
								ID:        "sub-1",
								CidrBlock: "10.0.10.0",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts ipv6 cidr block for subnets",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: []infrav1.SubnetSpec{
							{
								ID:            "sub-1",
								IPv6CidrBlock: "2022:1234:5678:9101::/64",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects invalid ipv6 cidr block for subnets",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						Subnets: []infrav1.SubnetSpec{
							{
								ID:            "sub-1",
								IPv6CidrBlock: "2022:1234:5678:9101::",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects ingress rules with cidr block and source security group id",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:               infrav1.SecurityGroupProtocolTCP,
								CidrBlocks:             []string{"test"},
								SourceSecurityGroupIDs: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects ingress rules with cidr block and source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								IPv6CidrBlocks:           []string{"test"},
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects ingress rules with cidr block, source security group id, role and nat gateway IP source",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								IPv6CidrBlocks:           []string{"test"},
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
								NatGatewaysIPsSource:     true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects ingress rules with source security role and nat gateway IP source",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
								NatGatewaysIPsSource:     true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects ingress rules with cidr block and nat gateway IP source",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:             infrav1.SecurityGroupProtocolTCP,
								IPv6CidrBlocks:       []string{"test"},
								NatGatewaysIPsSource: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts ingress rules with cidr block",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:   infrav1.SecurityGroupProtocolTCP,
								CidrBlocks: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts ingress rules with nat gateway IPs source",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:             infrav1.SecurityGroupProtocolTCP,
								NatGatewaysIPsSource: true,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts ingress rules with source security group role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts ingress rules with source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						IngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects ipamPool if id or name not set",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPAMPool: &infrav1.IPAMPool{},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects cidrBlock and ipamPool if set together",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							CidrBlock: "10.0.0.0/16",
							IPAMPool:  &infrav1.IPAMPool{},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts CP ingress rules with source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects CP ingress rules with cidr block and source security group id",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
							{
								Protocol:               infrav1.SecurityGroupProtocolTCP,
								CidrBlocks:             []string{"test"},
								SourceSecurityGroupIDs: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects CP ingress rules with cidr block and source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								IPv6CidrBlocks:           []string{"test"},
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts CP ingress rules with cidr block",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
							{
								Protocol:   infrav1.SecurityGroupProtocolTCP,
								CidrBlocks: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts CP ingress rules with source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalControlPlaneIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts node ingress rules with source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalNodeIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects node ingress rules with cidr block and source security group id",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalNodeIngressRules: []infrav1.IngressRule{
							{
								Protocol:               infrav1.SecurityGroupProtocolTCP,
								CidrBlocks:             []string{"test"},
								SourceSecurityGroupIDs: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects node ingress rules with cidr block and source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalNodeIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								IPv6CidrBlocks:           []string{"test"},
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts node ingress rules with cidr block",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalNodeIngressRules: []infrav1.IngressRule{
							{
								Protocol:   infrav1.SecurityGroupProtocolTCP,
								CidrBlocks: []string{"test"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts node ingress rules with source security group id and role",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						AdditionalNodeIngressRules: []infrav1.IngressRule{
							{
								Protocol:                 infrav1.SecurityGroupProtocolTCP,
								SourceSecurityGroupIDs:   []string{"test"},
								SourceSecurityGroupRoles: []infrav1.SecurityGroupRole{infrav1.SecurityGroupBastion},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "accepts cidrBlock for default node port ingress rule",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						NodePortIngressRuleCidrBlocks: []string{"10.0.0.0/16"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "reject invalid cidrBlock for default node port ingress rule",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						NodePortIngressRuleCidrBlocks: []string{"10.0.0.0"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects targetGroupIPType when LoadBalancer is disabled",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
						LoadBalancerType:  infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects targetGroupIPType with Classic Load Balancer",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeClassic,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts targetGroupIPType IPv4 with Network Load Balancer",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects targetGroupIPType IPv6 with VPC IPv6 disabled",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv6,
					},
					NetworkSpec: infrav1.NetworkSpec{},
				},
			},
			wantErr: true,
		},
		{
			name: "accepts targetGroupIPType IPv6 with NLB and VPC IPv6 enabled",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv6,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPv6: &infrav1.IPv6{
								CidrBlock: "2001:db8::/56",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "rejects additionalListener targetGroupIPType with Classic Load Balancer",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeClassic,
						AdditionalListeners: []infrav1.AdditionalListenerSpec{
							{
								Port:              22623,
								Protocol:          infrav1.ELBProtocolTCP,
								TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "rejects additionalListener targetGroupIPType IPv6 with VPC IPv6 disabled",
			cluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
						AdditionalListeners: []infrav1.AdditionalListenerSpec{
							{
								Port:              8443,
								Protocol:          infrav1.ELBProtocolTCP,
								TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv6,
							},
						},
					},
					NetworkSpec: infrav1.NetworkSpec{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utilfeature.SetFeatureGateDuringTest(t, feature.Gates, feature.BootstrapFormatIgnition, true)

			cluster := tt.cluster.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, cluster); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			c := &infrav1.AWSCluster{}
			key := client.ObjectKey{
				Name:      cluster.Name,
				Namespace: "default",
			}

			g := NewWithT(t)
			g.Eventually(func() bool {
				err := testEnv.Get(ctx, key, c)
				return err == nil
			}, 10*time.Second).Should(BeTrue(), fmt.Sprintf("Eventually failed getting the newly created cluster %q", cluster.Name))

			if tt.expect != nil {
				tt.expect(g, c.Spec.ControlPlaneLoadBalancer)
			}
		})
	}
}

func TestAWSClusterValidateUpdate(t *testing.T) {
	tests := []struct {
		name       string
		oldCluster *infrav1.AWSCluster
		newCluster *infrav1.AWSCluster
		wantErr    bool
	}{
		{
			name: "Control Plane LB type is immutable when switching from disabled to any",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeClassic,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Control Plane LB type is immutable when switching from any to disabled",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeClassic,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeDisabled,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "region is immutable",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Region: "us-east-1",
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Region: "us-east-2",
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer name is immutable",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: aws.String("old-apiserver"),
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: aws.String("new-apiserver"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer name is immutable, even if it is nil",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: nil,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: aws.String("example-apiserver"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer scheme is immutable",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternal,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternetFacing,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer scheme is immutable when left empty",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternal,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer scheme can be set to default when left empty",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Scheme: &infrav1.ELBSchemeInternetFacing,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "controlPlaneLoadBalancer crossZoneLoadBalancer is mutable",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						CrossZoneLoadBalancing: false,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						CrossZoneLoadBalancing: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "controlPlaneEndpoint is immutable",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1beta1.APIEndpoint{
						Host: "example.com",
						Port: int32(8000),
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1beta1.APIEndpoint{
						Host: "foo.example.com",
						Port: int32(9000),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneEndpoint can be updated if it is empty",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1beta1.APIEndpoint{},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1beta1.APIEndpoint{
						Host: "example.com",
						Port: int32(8000),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "removal of externally managed annotation is not allowed",
			oldCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{clusterv1beta1.ManagedByAnnotation: ""},
				},
			},
			newCluster: &infrav1.AWSCluster{},
			wantErr:    true,
		},
		{
			name:       "adding externally managed annotation is allowed",
			oldCluster: &infrav1.AWSCluster{},
			newCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{clusterv1beta1.ManagedByAnnotation: ""},
				},
			},
			wantErr: false,
		},
		{
			name: "VPC id is immutable cannot be emptied once set",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{ID: "managed-or-unmanaged-vpc"},
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			wantErr: true,
		},
		{
			name: "VPC id is immutable, cannot be set to a different value once set",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{ID: "managed-or-unmanaged-vpc"},
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{ID: "a-new-vpc"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid keys are not accepted during update",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					AdditionalTags: infrav1.Tags{
						"key-1": "value-1",
						"key-2": "value-2",
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					AdditionalTags: infrav1.Tags{
						"key-1":                    "value-1",
						"":                         "value-2",
						strings.Repeat("CAPI", 33): "value-3",
						"key-4":                    strings.Repeat("CAPI", 65),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should fail if controlPlaneLoadBalancer healthcheckprotocol is updated and not classic elb",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeNLB,
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeNLB,
						HealthCheckProtocol: &infrav1.ELBProtocolSSL,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should pass if controlPlaneLoadBalancer healthcheckprotocol is updated and a classic elb",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeClassic,
						HealthCheckProtocol: &infrav1.ELBProtocolSSL,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeClassic,
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should pass if old secondary lb is absent",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					SecondaryControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						Name: ptr.To("test-lb"),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should pass if controlPlaneLoadBalancer healthcheckprotocol is same after update",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Should fail if controlPlaneLoadBalancer healthcheckprotocol is changed to non-default if it was not set before update for non-classic elb",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeNLB,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeNLB,
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Should pass if controlPlaneLoadBalancer healthcheckprotocol is changed to non-default if it was not set before update for classic elb",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType: infrav1.LoadBalancerTypeClassic,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:    infrav1.LoadBalancerTypeNLB,
						HealthCheckProtocol: &infrav1.ELBProtocolTCP,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "correct GC tasks annotation",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						infrav1.ExternalResourceGCTasksAnnotation: "load-balancer,target-group,security-group",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty GC tasks annotation",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						infrav1.ExternalResourceGCTasksAnnotation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "incorrect GC tasks annotation",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			newCluster: &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						infrav1.ExternalResourceGCTasksAnnotation: "load-balancer,INVALID,security-group",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should failed if controlPlaneLoadBalancer targetGroupIPType is changed",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv6,
					},
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							IPv6: &infrav1.IPv6{
								CidrBlock: "2001:db8::/56",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should pass controlPlaneLoadBalancer targetGroupIPType is the same on update",
			oldCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
					},
				},
			},
			newCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
						LoadBalancerType:  infrav1.LoadBalancerTypeNLB,
						TargetGroupIPType: &infrav1.TargetGroupIPTypeIPv4,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cluster := tt.oldCluster.DeepCopy()
			cluster.ObjectMeta.GenerateName = "cluster-"
			cluster.ObjectMeta.Namespace = "default"

			if err := testEnv.Create(ctx, cluster); err != nil {
				t.Errorf("failed to create cluster: %v", err)
			}
			cluster.ObjectMeta.Annotations = tt.newCluster.Annotations
			cluster.Spec = tt.newCluster.Spec
			if err := testEnv.Update(ctx, cluster); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		},
		)
	}
}

func TestAWSClusterDefaultCNIIngressRules(t *testing.T) {
	AZUsageLimit := 3
	defaultVPCSpec := infrav1.VPCSpec{
		AvailabilityZoneUsageLimit: &AZUsageLimit,
		AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
		SubnetSchema:               &infrav1.SubnetSchemaPreferPrivate,
	}
	g := NewWithT(t)
	tests := []struct {
		name          string
		beforeCluster *infrav1.AWSCluster
		afterCluster  *infrav1.AWSCluster
	}{
		{
			name: "CNI ingressRules are updated cni spec undefined",
			beforeCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			afterCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &infrav1.CNISpec{
							CNIIngressRules: infrav1.CNIIngressRules{
								{
									Description: "bgp (calico)",
									Protocol:    infrav1.SecurityGroupProtocolTCP,
									FromPort:    179,
									ToPort:      179,
								},
								{
									Description: "IP-in-IP (calico)",
									Protocol:    infrav1.SecurityGroupProtocolIPinIP,
									FromPort:    -1,
									ToPort:      65535,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CNIIngressRules are not added for empty CNISpec",
			beforeCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &infrav1.CNISpec{},
					},
				},
			},
			afterCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &infrav1.CNISpec{},
					},
				},
			},
		},
		{
			name: "CNI ingressRules are unmodified when they exist",
			beforeCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &infrav1.CNISpec{
							CNIIngressRules: infrav1.CNIIngressRules{
								{
									Description: "Antrea 1",
									Protocol:    infrav1.SecurityGroupProtocolTCP,
									FromPort:    10349,
									ToPort:      10349,
								},
							},
						},
					},
				},
			},
			afterCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &infrav1.CNISpec{
							CNIIngressRules: infrav1.CNIIngressRules{
								{
									Description: "Antrea 1",
									Protocol:    infrav1.SecurityGroupProtocolTCP,
									FromPort:    10349,
									ToPort:      10349,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cluster := tt.beforeCluster.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			g.Expect(testEnv.Create(ctx, cluster)).To(Succeed())
			g.Expect(cluster.Spec.NetworkSpec).To(Equal(tt.afterCluster.Spec.NetworkSpec))
		})
	}
}

func TestAWSClusterValidateAllowedCIDRBlocks(t *testing.T) {
	tests := []struct {
		name    string
		awsc    *infrav1.AWSCluster
		wantErr bool
	}{
		{
			name: "allow valid CIDRs",
			awsc: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks: []string{
							"192.168.0.0/16",
							"192.168.0.1/32",
							"2001:1234:5678:9a40::/56",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disableIngressRules allowed with empty CIDR block",
			awsc: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks:   []string{},
						DisableIngressRules: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disableIngressRules not allowed with CIDR blocks",
			awsc: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks: []string{
							"192.168.0.0/16",
							"192.168.0.1/32",
							"2001:1234:5678:9a40::/56",
						},
						DisableIngressRules: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CIDR block with invalid network",
			awsc: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks: []string{
							"100.200.300.400/99",
							"2001:1234:5678:9a40::/129",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CIDR block with garbage string",
			awsc: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks: []string{
							"abcdefg",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cluster := tt.awsc.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			if err := testEnv.Create(ctx, cluster); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAllowedCIDRBlocks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSClusterDefaultAllowedCIDRBlocks(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name          string
		beforeCluster *infrav1.AWSCluster
		afterCluster  *infrav1.AWSCluster
	}{
		{
			name: "empty AllowedCIDRBlocks is defaulted to allow open ingress to bastion host",
			beforeCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{},
			},
			afterCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks: []string{
							"0.0.0.0/0",
							"::/0",
						},
					},
				},
			},
		},
		{
			name: "AllowedCIDRBlocks change not allowed if DisableIngressRules is true",
			beforeCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						AllowedCIDRBlocks:   []string{"0.0.0.0/0", "::/0"},
						DisableIngressRules: true,
						Enabled:             true,
					},
				},
			},
			afterCluster: &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					Bastion: infrav1.Bastion{
						DisableIngressRules: true,
						Enabled:             true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cluster := tt.beforeCluster.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			err := testEnv.Create(ctx, cluster)
			if err != nil {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(cluster.Spec.Bastion).To(Equal(tt.afterCluster.Spec.Bastion))
			}
		})
	}
}

// defaultValidateTest returns a new testing function to be used in tests to
// make sure defaulting webhooks also pass validation tests on create,
// update and delete.
// NOTE: This is a copy of the DefaultValidateTest function in the cluster-api
// package, but it has been modified to allow warnings to be returned.
func defaultValidateTest(ctx context.Context, object runtime.Object, webhook defaulting.DefaulterValidator, allowWarnings bool) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		createCopy := object.DeepCopyObject()
		updateCopy := object.DeepCopyObject()
		deleteCopy := object.DeepCopyObject()
		defaultingUpdateCopy := updateCopy.DeepCopyObject()

		t.Run("validate-on-create", func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(webhook.Default(ctx, createCopy)).To(Succeed())
			warnings, err := webhook.ValidateCreate(ctx, createCopy)
			g.Expect(err).ToNot(HaveOccurred())
			if !allowWarnings {
				g.Expect(warnings).To(BeEmpty())
			}
		})
		t.Run("validate-on-update", func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(webhook.Default(ctx, defaultingUpdateCopy)).To(Succeed())
			g.Expect(webhook.Default(ctx, updateCopy)).To(Succeed())
			warnings, err := webhook.ValidateUpdate(ctx, updateCopy, defaultingUpdateCopy)
			g.Expect(err).ToNot(HaveOccurred())
			if !allowWarnings {
				g.Expect(warnings).To(BeEmpty())
			}
		})
		t.Run("validate-on-delete", func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(webhook.Default(ctx, deleteCopy)).To(Succeed())
			warnings, err := webhook.ValidateDelete(ctx, deleteCopy)
			g.Expect(err).ToNot(HaveOccurred())
			if !allowWarnings {
				g.Expect(warnings).To(BeEmpty())
			}
		})
	}
}
