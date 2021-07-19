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
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestAWSCluster_ValidateCreate(t *testing.T) {
	tests := []struct {
		name    string
		cluster *AWSCluster
		wantErr bool
	}{
		// The SSHKeyName tests were moved to sshkeyname_test.go
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster := tt.cluster.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			ctx := context.TODO()
			if err := testEnv.Create(ctx, cluster); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAWSCluster_ValidateUpdate(t *testing.T) {
	tests := []struct {
		name       string
		oldCluster *AWSCluster
		newCluster *AWSCluster
		wantErr    bool
	}{
		{
			name: "region is immutable",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					Region: "us-east-1",
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					Region: "us-east-2",
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer scheme is immutable",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneLoadBalancer: &AWSLoadBalancerSpec{
						Scheme: &ClassicELBSchemeInternal,
					},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneLoadBalancer: &AWSLoadBalancerSpec{
						Scheme: &ClassicELBSchemeInternetFacing,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneLoadBalancer crossZoneLoadBalancer is mutable",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneLoadBalancer: &AWSLoadBalancerSpec{
						CrossZoneLoadBalancing: false,
					},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneLoadBalancer: &AWSLoadBalancerSpec{
						CrossZoneLoadBalancing: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "controlPlaneEndpoint is immutable",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1.APIEndpoint{
						Host: "example.com",
						Port: int32(8000),
					},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1.APIEndpoint{
						Host: "foo.example.com",
						Port: int32(9000),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "controlPlaneEndpoint can be updated if it is empty",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1.APIEndpoint{},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					ControlPlaneEndpoint: clusterv1.APIEndpoint{
						Host: "example.com",
						Port: int32(8000),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "VPC id is immutable cannot be emptied once set",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: VPCSpec{ID: "managed-or-unmanaged-vpc"},
					},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{},
			},
			wantErr: true,
		},
		{
			name: "VPC id is immutable, cannot be set to a different value once set",
			oldCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: VPCSpec{ID: "managed-or-unmanaged-vpc"},
					},
				},
			},
			newCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: VPCSpec{ID: "a-new-vpc"},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cluster := tt.oldCluster.DeepCopy()
			cluster.ObjectMeta = metav1.ObjectMeta{
				GenerateName: "cluster-",
				Namespace:    "default",
			}
			if err := testEnv.Create(ctx, cluster); err != nil {
				t.Errorf("failed to create cluster: %v", err)
			}
			cluster.Spec = tt.newCluster.Spec
			if err := testEnv.Update(ctx, cluster); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		},
		)
	}
}

func TestAWSCluster_DefaultCNIIngressRules(t *testing.T) {
	AZUsageLimit := 3
	defaultVPCSpec := VPCSpec{
		AvailabilityZoneUsageLimit: &AZUsageLimit,
		AvailabilityZoneSelection:  &AZSelectionSchemeOrdered,
	}
	g := NewWithT(t)
	tests := []struct {
		name          string
		beforeCluster *AWSCluster
		afterCluster  *AWSCluster
	}{
		{
			name: "CNI ingressRules are updated cni spec undefined",
			beforeCluster: &AWSCluster{
				Spec: AWSClusterSpec{},
			},
			afterCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &CNISpec{
							CNIIngressRules: CNIIngressRules{
								{
									Description: "bgp (calico)",
									Protocol:    SecurityGroupProtocolTCP,
									FromPort:    179,
									ToPort:      179,
								},
								{
									Description: "IP-in-IP (calico)",
									Protocol:    SecurityGroupProtocolIPinIP,
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
			beforeCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &CNISpec{},
					},
				},
			},
			afterCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &CNISpec{},
					},
				},
			},
		},
		{
			name: "CNI ingressRules are unmodified when they exist",
			beforeCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &CNISpec{
							CNIIngressRules: CNIIngressRules{
								{
									Description: "Antrea 1",
									Protocol:    SecurityGroupProtocolTCP,
									FromPort:    10349,
									ToPort:      10349,
								},
							},
						},
					},
				},
			},
			afterCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					NetworkSpec: NetworkSpec{
						VPC: defaultVPCSpec,
						CNI: &CNISpec{
							CNIIngressRules: CNIIngressRules{
								{
									Description: "Antrea 1",
									Protocol:    SecurityGroupProtocolTCP,
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

func TestAWSCluster_ValidateAllowedCIDRBlocks(t *testing.T) {
	tests := []struct {
		name    string
		awsc    *AWSCluster
		wantErr bool
	}{
		{
			name: "allow valid CIDRs",
			awsc: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks: []string{
							"192.168.0.0/16",
							"192.168.0.1/32",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disableIngressRules allowed with empty CIDR block",
			awsc: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks:   []string{},
						DisableIngressRules: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "disableIngressRules not allowed with CIDR blocks",
			awsc: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks: []string{
							"192.168.0.0/16",
							"192.168.0.1/32",
						},
						DisableIngressRules: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CIDR block with invalid network",
			awsc: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks: []string{
							"100.200.300.400/99",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CIDR block with garbage string",
			awsc: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
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

func TestAWSCluster_DefaultAllowedCIDRBlocks(t *testing.T) {
	g := NewWithT(t)
	tests := []struct {
		name          string
		beforeCluster *AWSCluster
		afterCluster  *AWSCluster
	}{
		{
			name: "empty AllowedCIDRBlocks is defaulted to allow open ingress to bastion host",
			beforeCluster: &AWSCluster{
				Spec: AWSClusterSpec{},
			},
			afterCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks: []string{
							"0.0.0.0/0",
						},
					},
				},
			},
		},
		{
			name: "empty AllowedCIDRBlocks is kept if DisableIngressRules is true",
			beforeCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks:   []string{},
						DisableIngressRules: true,
						Enabled:             true,
					},
				},
			},
			afterCluster: &AWSCluster{
				Spec: AWSClusterSpec{
					Bastion: Bastion{
						AllowedCIDRBlocks:   []string{},
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
			g.Expect(testEnv.Create(ctx, cluster)).To(Succeed())
			g.Expect(cluster.Spec.Bastion).To(Equal(tt.afterCluster.Spec.Bastion))
		})
	}
}
