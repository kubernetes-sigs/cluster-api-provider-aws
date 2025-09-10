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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

func TestNodeadmUserdata(t *testing.T) {
	format.TruncatedDiff = false
	g := NewWithT(t)

	type args struct {
		input *NodeadmInput
	}

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
			name: "with kubelet config",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					KubeletConfig: &runtime.RawExtension{
						Raw: []byte(`
evictionHard:
  memory.available: "2000Mi"
`),
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "evictionHard:") &&
					strings.Contains(output, "memory.available: 2000Mi") &&
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
					PreNodeadmCommands: []string{
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
					AMIImageID:        "ami-123456",
					ServiceCIDR:       "192.168.0.0/16",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "cidr: 192.168.0.0/16") &&
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
					ClusterName:        "test-cluster",
					APIServerEndpoint:  "https://example.com",
					CACert:             "test-ca-cert",
					Boundary:           "CUSTOMBOUNDARY123",
					PreNodeadmCommands: []string{"echo 'pre-bootstrap'"},
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
					strings.Count(output, fmt.Sprintf("--%s", boundary)) == 5 // 3 parts * 2 boundaries each except cloud-config
			},
		},
		{
			name: "boundary verification - only node config part with default boundary",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
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
					ClusterName:        "test-cluster",
					APIServerEndpoint:  "https://example.com",
					CACert:             "test-ca-cert",
					PreNodeadmCommands: []string{"echo 'test'"},
					NTP: &eksbootstrapv1.NTP{
						Enabled: ptr.To(true),
						Servers: []string{"time.google.com"},
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				boundary := "//" // default boundary
				return strings.Contains(output, fmt.Sprintf(`boundary=%q`, boundary)) &&
					strings.Contains(output, "Content-Type: application/node.eks.aws") &&
					strings.Contains(output, "Content-Type: text/x-shellscript") &&
					strings.Contains(output, "Content-Type: text/cloud-config") &&
					strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1") &&
					strings.Count(output, fmt.Sprintf("--%s", boundary)) == 5 // 3 parts * 2 boundaries each except cloud-config
			},
		},
		{
			name: "verify other kubelet flags are preserved with node-labels",
			args: args{
				input: &NodeadmInput{
					ClusterName:       "test-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "test-ca-cert",
					KubeletFlags: []string{
						"--node-labels=tier=workers",
						"--register-with-taints=dedicated=gpu:NoSchedule",
						"--max-pods=58",
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "--node-labels") &&
					strings.Contains(output, "tier=workers") &&
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
					ClusterName: "test-cluster",
					CACert:      "test-ca-cert",
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
					// Missing CACert
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
