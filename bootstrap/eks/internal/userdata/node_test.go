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
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"k8s.io/utils/pointer"
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
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster
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
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --kubelet-extra-args '--node-labels=node-role.undistro.io/infra=true --register-with-taints=dedicated=infra:NoSchedule'
`),
		},
		{
			name: "with container runtime",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					ContainerRuntime: pointer.String("containerd"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --container-runtime containerd
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
					ContainerRuntime: pointer.String("containerd"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --kubelet-extra-args '--node-labels=node-role.undistro.io/infra=true --register-with-taints=dedicated=infra:NoSchedule' --container-runtime containerd
`),
		},
		{
			name: "with ipv6",
			args: args{
				input: &NodeInput{
					ClusterName:     "test-cluster",
					ServiceIPV6Cidr: pointer.String("fe80:0000:0000:0000:0204:61ff:fe9d:f156/24"),
					IPFamily:        pointer.String("ipv6"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --ip-family ipv6 --service-ipv6-cidr fe80:0000:0000:0000:0204:61ff:fe9d:f156/24
`),
		},
		{
			name: "without max pods",
			args: args{
				input: &NodeInput{
					ClusterName: "test-cluster",
					UseMaxPods:  pointer.Bool(false),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --use-max-pods false
`),
		},
		{
			name: "with api retry attempts",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					APIRetryAttempts: pointer.Int(5),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --aws-api-retry-attempts 5
`),
		},
		{
			name: "with pause container",
			args: args{
				input: &NodeInput{
					ClusterName:           "test-cluster",
					PauseContainerAccount: pointer.String("12345678"),
					PauseConatinerVersion: pointer.String("v1"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --pause-container-account 12345678 --pause-container-version v1
`),
		},
		{
			name: "with dns cluster ip",
			args: args{
				input: &NodeInput{
					ClusterName:  "test-cluster",
					DNSClusterIP: pointer.String("192.168.0.1"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --dns-cluster-ip 192.168.0.1
`),
		},
		{
			name: "with docker json",
			args: args{
				input: &NodeInput{
					ClusterName:      "test-cluster",
					DockerConfigJson: pointer.String("{'debug':true}"),
				},
			},
			expectedBytes: []byte(`#!/bin/bash
/etc/eks/bootstrap.sh test-cluster --docker-config-json {'debug':true}
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

const sampleDockerjson = `
{
	"allow-nondistributable-artifacts": [],
	"api-cors-header": "",
	"authorization-plugins": [],
	"bip": "",
	"bridge": "",
	"cgroup-parent": "",
	"cluster-advertise": "",
	"cluster-store": "",
	"cluster-store-opts": {},
	"containerd": "/run/containerd/containerd.sock",
	"containerd-namespace": "docker",
	"containerd-plugin-namespace": "docker-plugins",
	"data-root": "",
	"debug": true,
	"default-address-pools": [
	  {
		"base": "172.30.0.0/16",
		"size": 24
	  },
	  {
		"base": "172.31.0.0/16",
		"size": 24
	  }
	],
	"default-cgroupns-mode": "private",
	"default-gateway": "",
	"default-gateway-v6": "",
	"default-runtime": "runc",
	"default-shm-size": "64M",
	"default-ulimits": {
	  "nofile": {
		"Hard": 64000,
		"Name": "nofile",
		"Soft": 64000
	  }
	},
	"dns": [],
	"dns-opts": [],
	"dns-search": [],
	"exec-opts": [],
	"exec-root": "",
	"experimental": false,
	"features": {},
	"fixed-cidr": "",
	"fixed-cidr-v6": "",
	"group": "",
	"hosts": [],
	"icc": false,
	"init": false,
	"init-path": "/usr/libexec/docker-init",
	"insecure-registries": [],
	"ip": "0.0.0.0",
	"ip-forward": false,
	"ip-masq": false,
	"iptables": false,
	"ip6tables": false,
	"ipv6": false,
	"labels": [],
	"live-restore": true,
	"log-driver": "json-file",
	"log-level": "",
	"log-opts": {
	  "cache-disabled": "false",
	  "cache-max-file": "5",
	  "cache-max-size": "20m",
	  "cache-compress": "true",
	  "env": "os,customer",
	  "labels": "somelabel",
	  "max-file": "5",
	  "max-size": "10m"
	},
	"max-concurrent-downloads": 3,
	"max-concurrent-uploads": 5,
	"max-download-attempts": 5,
	"mtu": 0,
	"no-new-privileges": false,
	"node-generic-resources": [
	  "NVIDIA-GPU=UUID1",
	  "NVIDIA-GPU=UUID2"
	],
	"oom-score-adjust": -500,
	"pidfile": "",
	"raw-logs": false,
	"registry-mirrors": [],
	"runtimes": {
	  "cc-runtime": {
		"path": "/usr/bin/cc-runtime"
	  },
	  "custom": {
		"path": "/usr/local/bin/my-runc-replacement",
		"runtimeArgs": [
		  "--debug"
		]
	  }
	},
	"seccomp-profile": "",
	"selinux-enabled": false,
	"shutdown-timeout": 15,
	"storage-driver": "",
	"storage-opts": [],
	"swarm-default-advertise-addr": "",
	"tls": true,
	"tlscacert": "",
	"tlscert": "",
	"tlskey": "",
	"tlsverify": true,
	"userland-proxy": false,
	"userland-proxy-path": "/usr/libexec/docker-proxy",
	"userns-remap": ""
  }
`
