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

package userdata

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"k8s.io/utils/ptr"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
)

func TestNewNode(t *testing.T) {
	format.TruncatedDiff = false
	g := NewWithT(t)

	type args struct {
		input *NodeInput
	}

	tests := []struct {
		name          string
		args          args
		expectedBytes []byte
		expectErr     bool
	}{
		{
			name: "only cluster name",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
`),
			expectErr: false,
		},
		{
			name: "sample-with-values",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					KubeletExtraArgs: map[string]string{
						"node-labels":          "node-role.undistro.io/infra=true",
						"register-with-taints": "dedicated=infra:NoSchedule",
					},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --kubelet-extra-args '--node-labels=node-role.undistro.io/infra=true --register-with-taints=dedicated=infra:NoSchedule'
`),
		},
		{
			name: "with container runtime",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					ContainerRuntime: ptr.To[string]("containerd"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --container-runtime containerd
`),
		},
		{
			name: "with kubelet extra args and container runtime",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					KubeletExtraArgs: map[string]string{
						"node-labels":          "node-role.undistro.io/infra=true",
						"register-with-taints": "dedicated=infra:NoSchedule",
					},
					ContainerRuntime: ptr.To[string]("containerd"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --kubelet-extra-args '--node-labels=node-role.undistro.io/infra=true --register-with-taints=dedicated=infra:NoSchedule' --container-runtime containerd
`),
		},
		{
			name: "with ipv6",
			args: args{
				input: &NodeInput{
					ClusterName:     "test-cluster",
					ServiceIPV6Cidr: ptr.To[string]("fe80:0000:0000:0000:0204:61ff:fe9d:f156/24"),
					IPFamily:        ptr.To[string]("ipv6"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --ip-family ipv6 --service-ipv6-cidr fe80:0000:0000:0000:0204:61ff:fe9d:f156/24
`),
		},
		{
			name: "without max pods",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					UseMaxPods:  ptr.To[bool](false),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --use-max-pods false
`),
		},
		{
			name: "with api retry attempts",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					APIRetryAttempts: ptr.To[int](5),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --aws-api-retry-attempts 5
`),
		},
		{
			name: "with pause container",
			args: args{
				input: &NodeInput{
					ClusterName:           "test-cluster",
					PauseContainerAccount: ptr.To[string]("12345678"),
					PauseContainerVersion: ptr.To[string]("v1"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --pause-container-account 12345678 --pause-container-version v1
`),
		},
		{
			name: "with dns cluster ip",
			args: args{
				input: &NodeInput{
					ClusterName:  "test-cluster",
					DNSClusterIP: ptr.To[string]("192.168.0.1"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --dns-cluster-ip 192.168.0.1
`),
		},
		{
			name: "with docker json",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					DockerConfigJSON: ptr.To[string]("{\"debug\":true}"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster --docker-config-json '{"debug":true}'
`),
		},
		{
			name: "with pre-bootstrap command",
			args: args{
				input: &NodeInput{
					ClusterName:          "test-cluster",
					PreBootstrapCommands: []string{"date", "echo \"testing\""},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - "date"
  - "echo \"testing\""
  - /etc/eks/bootstrap.sh test-cluster
`),
		},
		{
			name: "with post-bootstrap command",
			args: args{
				input: &NodeInput{
					ClusterName:           "test-cluster",
					PostBootstrapCommands: []string{"date", "echo \"testing\""},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
  - "date"
  - "echo \"testing\""
`),
		},
		{
			name: "with pre & post-bootstrap command",
			args: args{
				input: &NodeInput{
					ClusterName:           "test-cluster",
					PreBootstrapCommands:  []string{"echo \"testing pre\""},
					PostBootstrapCommands: []string{"echo \"testing post\""},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - "echo \"testing pre\""
  - /etc/eks/bootstrap.sh test-cluster
  - "echo \"testing post\""
`),
		},
		{
			name: "with bootstrap override command",
			args: args{
				input: &NodeInput{
					ClusterName:              "test-cluster",
					BootstrapCommandOverride: ptr.To[string]("/custom/mybootstrap.sh"),
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /custom/mybootstrap.sh test-cluster
`),
		},
		{
			name: "with disk setup and mount points",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					DiskSetup: &eksbootstrapv1.DiskSetup{
						Filesystems: []eksbootstrapv1.Filesystem{
							{
								Device:     "/dev/sdb",
								Filesystem: "ext4",
								Label:      "vol2",
							},
						},
						Partitions: []eksbootstrapv1.Partition{
							{
								Device: "/dev/sdb",
								Layout: true,
							},
						},
					},
					Mounts: []eksbootstrapv1.MountPoints{
						[]string{"LABEL=vol2", "/mnt/vol2", "ext4", "defaults"},
						[]string{"LABEL=vol2", "/opt/data", "ext4", "defaults"},
					},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
disk_setup:
  /dev/sdb:
    layout: true
fs_setup:
  - label: vol2
    filesystem: ext4
    device: /dev/sdb
mounts:
  -
    - LABEL=vol2
    - /mnt/vol2
    - ext4
    - defaults
  -
    - LABEL=vol2
    - /opt/data
    - ext4
    - defaults
`),
		},
		{
			name: "with files",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					Files: []eksbootstrapv1.File{
						{
							Path:    "/etc/sysctl.d/91-fs-inotify.conf",
							Content: "fs.inotify.max_user_instances=256",
						},
					},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
  - path: /etc/sysctl.d/91-fs-inotify.conf
    content: |
      fs.inotify.max_user_instances=256
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
`),
		},
		{
			name: "with ntp",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					NTP: &eksbootstrapv1.NTP{
						Enabled: aws.Bool(true),
						Servers: []string{"time1.google.com", "time2.google.com", "time3.google.com", "time4.google.com"},
					},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
ntp:
  enabled: true
  servers:
    - time1.google.com
    - time2.google.com
    - time3.google.com
    - time4.google.com
`),
		},
		{
			name: "with users",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					Users: []eksbootstrapv1.User{
						{
							Name:  "testuser",
							Shell: aws.String("/bin/bash"),
						},
					},
				},
			},
			expectedBytes: []byte(`#cloud-config
write_files:
runcmd:
  - /etc/eks/bootstrap.sh test-cluster
users:
  - name: testuser
    shell: /bin/bash
`),
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			bytes, err := NewNode(testcase.args.input)
			if testcase.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(bytes)).To(Equal(string(testcase.expectedBytes)))
		})
	}
}

func TestNewNodeAL2023(t *testing.T) {
	g := NewWithT(t)

	type args struct {
		input *NodeInput
	}

	tests := []struct {
		name         string
		args         args
		expectErr    bool
		verifyOutput func(output string) bool
	}{
		{
			name: "AL2023 with shell script and node config",
			args: args{
				input: &NodeInput{
					AMIFamilyType:     AMIFamilyAL2023,
					ClusterName:       "my-cluster",
					APIServerEndpoint: "https://example.com",
					CACert:            "Y2VydGlmaWNhdGVBdXRob3JpdHk=",
					NodeGroupName:     "test-nodegroup",
					DNSClusterIP:      ptr.To[string]("10.100.0.10"),
					Boundary:          "BOUNDARY",
					KubeletExtraArgs: map[string]string{
						"node-labels": "app=my-app,environment=production",
					},
					PreBootstrapCommands: []string{
						"# Install additional packages",
						"yum install -y htop jq iptables-services",
						"",
						"# Pre-cache commonly used container images",
						"nohup docker pull public.ecr.aws/eks-distro/kubernetes/pause:3.2 &",
						"",
						"# Configure HTTP proxy if needed",
						`cat > /etc/profile.d/http-proxy.sh << 'EOF'
export HTTP_PROXY="http://proxy.example.com:3128"
export HTTPS_PROXY="http://proxy.example.com:3128"
export NO_PROXY="localhost,127.0.0.1,169.254.169.254,.internal"
EOF`,
					},
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				// Verify MIME structure
				if !strings.Contains(output, "MIME-Version: 1.0") ||
					!strings.Contains(output, `Content-Type: multipart/mixed; boundary="BOUNDARY"`) {
					return false
				}

				// Verify shell script content
				if !strings.Contains(output, "#!/bin/bash") ||
					!strings.Contains(output, "yum install -y htop jq iptables-services") ||
					!strings.Contains(output, "docker pull public.ecr.aws/eks-distro/kubernetes/pause:3.2") {
					return false
				}

				// Verify node config content
				if !strings.Contains(output, "apiVersion: node.eks.aws/v1alpha1") ||
					!strings.Contains(output, "name: my-cluster") ||
					!strings.Contains(output, "apiServerEndpoint: https://example.com") ||
					!strings.Contains(output, `"--node-labels=app=my-app,environment=production"`) ||
					!strings.Contains(output, "cidr: 172.20.0.0/16") {
					return false
				}

				return true
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			bytes, err := NewNode(testcase.args.input)
			if testcase.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			if testcase.verifyOutput != nil {
				g.Expect(testcase.verifyOutput(string(bytes))).To(BeTrue(), "Output verification failed")
			}
		})
	}
}

func TestGenerateAL2023UserData(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name         string
		input        *NodeInput
		expectErr    bool
		verifyOutput func(output string) bool
	}{
		{
			name: "valid AL2023 input",
			input: &NodeInput{
				AMIFamilyType:     AMIFamilyAL2023,
				ClusterName:       "test-cluster",
				APIServerEndpoint: "https://test-endpoint.eks.amazonaws.com",
				CACert:            "test-cert",
				NodeGroupName:     "test-nodegroup",
				UseMaxPods:        ptr.To[bool](false),
				DNSClusterIP:      ptr.To[string]("172.20.0.10"),
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "name: test-cluster") &&
					strings.Contains(output, "maxPods: 58") &&
					strings.Contains(output, "nodegroup=test-nodegroup") &&
					strings.Contains(output, "cidr: 172.20.0.0/16") &&
					strings.Contains(output, "clusterDNS:\n      - 172.20.0.10")
			},
		},
		{
			name: "AL2023 with custom DNS and AMI",
			input: &NodeInput{
				AMIFamilyType:     AMIFamilyAL2023,
				ClusterName:       "test-cluster",
				APIServerEndpoint: "https://test-endpoint.eks.amazonaws.com",
				CACert:            "test-cert",
				NodeGroupName:     "test-nodegroup",
				UseMaxPods:        ptr.To[bool](true),
				DNSClusterIP:      ptr.To[string]("10.100.0.10"),
				AMIImageID:        "ami-123456",
				ServiceCIDR:       "192.168.0.0/16",
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "cidr: 192.168.0.0/16") &&
					strings.Contains(output, "maxPods: 110") &&
					strings.Contains(output, "nodegroup-image=ami-123456")
			},
		},
		{
			name: "AL2023 with custom labels and commands",
			input: &NodeInput{
				AMIFamilyType:     AMIFamilyAL2023,
				ClusterName:       "test-cluster",
				APIServerEndpoint: "https://test-endpoint.eks.amazonaws.com",
				CACert:            "test-cert",
				NodeGroupName:     "test-nodegroup",
				KubeletExtraArgs: map[string]string{
					"node-labels": "app=my-app,environment=production",
				},
				PreBootstrapCommands: []string{
					"echo 'pre-bootstrap'",
				},
				PostBootstrapCommands: []string{
					"echo 'post-bootstrap'",
				},
			},
			expectErr: false,
			verifyOutput: func(output string) bool {
				return strings.Contains(output, "echo 'pre-bootstrap'") &&
					strings.Contains(output, "echo 'post-bootstrap'") &&
					strings.Contains(output, `"--node-labels=app=my-app,environment=production"`) &&
					strings.Contains(output, "cidr: 172.20.0.0/16")
			},
		},
		{
			name: "AL2023 missing required fields",
			input: &NodeInput{
				AMIFamilyType: AMIFamilyAL2023,
				ClusterName:   "test-cluster",
				// Missing APIServerEndpoint, CACert, NodeGroupName
			},
			expectErr: true,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			bytes, err := generateAL2023UserData(testcase.input)
			if testcase.expectErr {
				g.Expect(err).To(HaveOccurred())
				return
			}

			g.Expect(err).NotTo(HaveOccurred())
			if testcase.verifyOutput != nil {
				g.Expect(testcase.verifyOutput(string(bytes))).To(BeTrue(), "Output verification failed")
			}
		})
	}
}
