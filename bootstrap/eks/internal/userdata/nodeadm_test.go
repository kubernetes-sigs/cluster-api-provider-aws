/*
Copyright 2025 The Kubernetes Authors.

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

package userdata

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"k8s.io/utils/ptr"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestSetKubeletFlags(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name           string
		in             *NodeadmInput
		wantNodeLabels []string
		wantOtherFlags []string
	}{
		{
			name:           "empty kubelet flags",
			in:             &NodeadmInput{},
			wantNodeLabels: []string{"eks.amazonaws.com/capacityType=ON_DEMAND"},
			wantOtherFlags: nil,
		},
		{
			name: "unrelated kubelet flag preserved",
			in: &NodeadmInput{
				KubeletFlags: []string{"--register-with-taints=dedicated=infra:NoSchedule"},
			},
			wantNodeLabels: []string{"eks.amazonaws.com/capacityType=ON_DEMAND"},
			wantOtherFlags: []string{"--register-with-taints=dedicated=infra:NoSchedule"},
		},
		{
			name: "existing node-labels augmented",
			in: &NodeadmInput{
				KubeletFlags:  []string{"--node-labels=app=foo"},
				AMIImageID:    "ami-12345",
				NodeGroupName: "ng-1",
			},
			wantNodeLabels: []string{
				"app=foo",
				"eks.amazonaws.com/nodegroup-image=ami-12345",
				"eks.amazonaws.com/nodegroup=ng-1",
				"eks.amazonaws.com/capacityType=ON_DEMAND",
			},
			wantOtherFlags: nil,
		},
		{
			name: "existing eks-specific labels present",
			in: &NodeadmInput{
				KubeletFlags:  []string{"--node-labels=app=foo,eks.amazonaws.com/nodegroup=ng-1,eks.amazonaws.com/nodegroup-image=ami-12345,eks.amazonaws.com/capacityType=SPOT"},
				AMIImageID:    "ami-12345",
				NodeGroupName: "ng-1",
			},
			wantNodeLabels: []string{
				"app=foo",
				"eks.amazonaws.com/nodegroup=ng-1",
				"eks.amazonaws.com/nodegroup-image=ami-12345",
				"eks.amazonaws.com/capacityType=SPOT",
			},
			wantOtherFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.setKubeletFlags()
			var gotNodeLabels []string
			var gotOtherFlags []string
			for _, flag := range tt.in.KubeletFlags {
				if strings.HasPrefix(flag, "--node-labels=") {
					labels := strings.TrimPrefix(flag, "--node-labels=")
					gotNodeLabels = append(gotNodeLabels, strings.Split(labels, ",")...)
				} else {
					gotOtherFlags = append(gotOtherFlags, flag)
				}
			}
			g.Expect(gotNodeLabels).To(ContainElements(tt.wantNodeLabels), "expected node-labels to contain %v, got %v", tt.wantNodeLabels, gotNodeLabels)
			g.Expect(gotOtherFlags).To(ContainElements(tt.wantOtherFlags), "expected kubelet flags to contain %v, got %v", tt.wantOtherFlags, gotOtherFlags)
		})
	}
}

func TestNodeadmUserdata(t *testing.T) {
	format.TruncatedDiff = false
	g := NewWithT(t)

	type args struct {
		input *NodeadmInput
	}

	onDemandCapacity := v1beta2.ManagedMachinePoolCapacityTypeOnDemand
	spotCapacity := v1beta2.ManagedMachinePoolCapacityTypeSpot

	tests := []struct {
		name         string
		args         args
		expectErr    bool
		verifyOutput func(output string) bool
	}{
		{
			name: "basic nodeadm userdata",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "MIME-Version: 1.0") &&
					strings.Contains(output, "name: test-cluster") &&
					strings.Contains(output, "apiServerEndpoint: https://example.com") &&
					strings.Contains(output, "certificateAuthority: test-ca-cert") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "with kubelet flags",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					KubeletFlags: []string{
						"--node-labels=node-role.undistro.io/infra=true",
						"--register-with-taints=dedicated=infra:NoSchedule",
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "node-role.undistro.io/infra=true") &&
					strings.Contains(output, "register-with-taints") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "with pre bootstrap commands",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					PreBootstrapCommands: []string{
						"echo 'pre-bootstrap'",
						"yum install -y htop",
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "echo 'pre-bootstrap'") &&
					strings.Contains(output, "yum install -y htop") &&
					strings.Contains(output, "#!/bin/bash") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "with custom AMI",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://test-endpoint.eks.amazonaws.com",
					CACert:            "test-cert",
					NodeGroupName:     "test-nodegroup",
					AMIImageID:        "ami-123456",
					ServiceCIDR:       "192.168.0.0/16",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "cidr: 192.168.0.0/16") &&
					strings.Contains(output, "nodegroup-image=ami-123456") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "cloud-config part when NTP is set",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					NTP: &eksbootstrapv1.NTP{
						Enabled: ptr.To(true),
						Servers: []string{"time.google.com"},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Contains(output, "#cloud-config") &&
					strings.Contains(output, "time.google.com") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "cloud-config part when Users is set",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					Users: []eksbootstrapv1.User{
						{
							Name:              "testuser",
							SSHAuthorizedKeys: []string{"ssh-rsa AAAAB3..."},
						},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Contains(output, "#cloud-config") &&
					strings.Contains(output, "testuser") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "cloud-config part when DiskSetup is set",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					DiskSetup: &eksbootstrapv1.DiskSetup{
						Filesystems: []eksbootstrapv1.Filesystem{
							{
								Device:     "/dev/disk/azure/scsi1/lun0",
								Filesystem: "ext4",
								Label:      "etcd_disk",
							},
						},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Contains(output, "#cloud-config") &&
					strings.Contains(output, "/dev/disk/azure/scsi1/lun0") &&
					strings.Contains(output, "ext4") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "cloud-config part when Mounts is set",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					Mounts: []eksbootstrapv1.MountPoints{
						{"/dev/disk/scsi1/lun0"},
						{"/mnt/etcd"},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Contains(output, "#cloud-config") &&
					strings.Contains(output, "/dev/disk/scsi1/lun0") &&
					strings.Contains(output, "/mnt/etcd") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "boundary verification - all three parts with custom boundary",
			args: args{
				input: &NodeadmInput{
					ClusterName:          "test-cluster",
					APIServerEndpoint:    "https://example.com",
					CACert:               "test-ca-cert",
					NodeGroupName:        "test-nodegroup",
					Boundary:             "CUSTOMBOUNDARY123",
					PreBootstrapCommands: []string{"echo 'pre-bootstrap'"},
					NTP: &eksbootstrapv1.NTP{
						Enabled: ptr.To(true),
						Servers: []string{"time.google.com"},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				boundary := "CUSTOMBOUNDARY123"
				return strings.Contains(output, fmt.Sprintf(`boundary=%q`, boundary)) &&
					strings.Contains(output, fmt.Sprintf("--%s", boundary)) &&
					strings.Contains(output, fmt.Sprintf("--%s--", boundary)) &&
					strings.Contains(output, "Content-Type: application/node.eks.aws") &&
					strings.Contains(output, "Content-Type: text/x-shellscript") &&
					strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Count(output, fmt.Sprintf("--%s", boundary)) == 6 // 3 parts * 2 boundaries each
			},
		},
		{
			name: "boundary verification - only node config part with default boundary",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				boundary := "//" // default boundary
				return strings.Contains(output, fmt.Sprintf(`boundary=%q`, boundary)) &&
					strings.Contains(output, fmt.Sprintf("--%s", boundary)) &&
					strings.Contains(output, fmt.Sprintf("--%s--", boundary)) &&
					strings.Contains(output, "Content-Type: application/node.eks.aws") &&
					!strings.Contains(output, "Content-Type: text/x-shellscript") &&
					!strings.Contains(output, "Content-Type: text/cloud-config")
			},
		},
		{
			name: "boundary verification - all 3 parts",
			args: args{
				input: &NodeadmInput{
					ClusterName:          "test-cluster",
					APIServerEndpoint:    "https://example.com",
					CACert:               "test-ca-cert",
					NodeGroupName:        "test-nodegroup",
					PreBootstrapCommands: []string{"echo 'test'"},
					NTP: &eksbootstrapv1.NTP{
						Enabled: ptr.To(true),
						Servers: []string{"time.google.com"},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				boundary := "//" // default boundary
				fmt.Println(output)
				return strings.Contains(output, fmt.Sprintf(`boundary=%q`, boundary)) &&
					strings.Contains(output, "Content-Type: application/node.eks.aws") &&
					strings.Contains(output, "Content-Type: text/x-shellscript") &&
					strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Count(output, fmt.Sprintf("--%s", boundary)) == 6 // 3 parts * 2 boundaries each
			},
		},
		{
			name: "node-labels without capacityType - should add ON_DEMAND",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					AMIImageID:        "ami-123456",
					KubeletFlags: []string{
						"--node-labels=app=my-app,environment=production",
					},
					CapacityType: nil, // Should default to ON_DEMAND
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "app=my-app") &&
					strings.Contains(output, "environment=production") &&
					strings.Contains(output, "eks.amazonaws.com/capacityType=ON_DEMAND") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup-image=ami-123456") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup=test-nodegroup") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "node-labels with capacityType set to SPOT",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					AMIImageID:        "ami-123456",
					KubeletFlags: []string{
						"--node-labels=workload=batch",
					},
					CapacityType: &spotCapacity,
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "workload=batch") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup-image=ami-123456") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup=test-nodegroup") &&
					strings.Contains(output, "eks.amazonaws.com/capacityType=SPOT") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "no existing node-labels - should only add generated labels",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					AMIImageID:        "ami-789012",
					KubeletFlags: []string{
						"--max-pods=100",
					},
					CapacityType: &spotCapacity,
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "--node-labels") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup-image=ami-789012") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup=test-nodegroup") &&
					strings.Contains(output, "eks.amazonaws.com/capacityType=SPOT") &&
					strings.Contains(output, `"--max-pods=100"`) &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "verify other kubelet flags are preserved with node-labels",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					NodeGroupName:     "test-nodegroup",
					KubeletFlags: []string{
						"--node-labels=tier=workers",
						"--register-with-taints=dedicated=gpu:NoSchedule",
						"--max-pods=58",
					},
					CapacityType: &onDemandCapacity,
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "--node-labels") &&
					strings.Contains(output, "tier=workers") &&
					strings.Contains(output, "eks.amazonaws.com/nodegroup=test-nodegroup") &&
					strings.Contains(output, "eks.amazonaws.com/capacityType=ON_DEMAND") &&
					strings.Contains(output, `"--register-with-taints=dedicated=gpu:NoSchedule"`) &&
					strings.Contains(output, `"--max-pods=58"`) &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1")
			},
		},
		{
			name: "missing required fields",
			args: args{
				input: &NodeadmInput{
					ClusterName: "test-cluster",
					// Missing APIServerEndpoint, CACert, NodeGroupName
				},
			},
			expectErr: true,
		},
		{
			name: "missing API server endpoint",
			args: args{
				input: &NodeadmInput{
					ClusterName:   "test-cluster",
					CACert:        "test-ca-cert",
					NodeGroupName: "test-nodegroup",
					// Missing APIServerEndpoint
				},
			},
			expectErr: true,
		},
		{
			name: "missing CA certificate",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					NodeGroupName:     "test-nodegroup",
					// Missing CACert
				},
			},
			expectErr: true,
		},
		{
			name: "missing node group name",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					// Missing NodeGroupName
				},
			},
			expectErr: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			bytes, err := NewNodeadmUserdata(testcase.args.input)
			if testcase.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			if testcase.verifyOutput != nil {
				g.Expect(testcase.verifyOutput(string(bytes))).To(BeTrue(), "Output verification failed for: %s", testcase.name)
			}
		})
	}
}
