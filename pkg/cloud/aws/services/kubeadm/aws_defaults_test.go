/*
Copyright 2019 The Kubernetes Authors.

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

package kubeadm_test

import (
	"fmt"
	"testing"

	kubeadmv1beta1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/kubeadm"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type jm struct {
	*actuators.Scope
	*clusterv1.Machine
}

func (j *jm) GetScope() *actuators.Scope {
	return j.Scope
}
func (j *jm) GetMachine() *clusterv1.Machine {
	return j.Machine
}

func TestSetJoinNodeConfigurationBootstrapSettings(t *testing.T) {
	testcases := []struct {
		name              string
		caCertHash        string
		bootstrapToken    string
		joinMachine       *jm
		joinConfiguration *kubeadmv1beta1.JoinConfiguration
	}{
		{
			name: "test default values",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{
						Network: v1alpha1.Network{
							APIServerELB: v1alpha1.ClassicELB{
								DNSName: "test.test",
							},
						},
					},
				},
				Machine: &clusterv1.Machine{},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := kubeadm.SetJoinNodeConfigurationOverrides(tc.caCertHash, tc.bootstrapToken, tc.joinMachine, tc.joinConfiguration)

			// The apiserver endpoint should always match the incoming machine's network
			expected := fmt.Sprintf("%v:%v", tc.joinMachine.GetScope().Network().APIServerELB.DNSName, kubeadm.APIServerBindPort)
			if out.Discovery.BootstrapToken.APIServerEndpoint != expected {
				t.Fatalf("join configuration apiserver endpoint: %q but expected %q", out.Discovery.BootstrapToken.APIServerEndpoint, expected)
			}

			// The bootstrap token on the new join node configuration should be
			// the same that is passed in
			if out.Discovery.BootstrapToken.Token != tc.bootstrapToken {
				t.Fatalf("bootstrap tokens did  not match: got %q but expected %q", out.Discovery.BootstrapToken.Token, tc.bootstrapToken)
			}

			// The passed in caCertHash should be appended to the existing hashes
			if !contains(out.Discovery.BootstrapToken.CACertHashes, tc.caCertHash) {
				t.Fatalf("did not find %q in %v", tc.caCertHash, out.Discovery.BootstrapToken.CACertHashes)
			}
		})
	}
}

func TestSetJoinNodeConfigurationOverrides(t *testing.T) {
	testcases := []struct {
		name              string
		caCertHash        string
		bootstrapToken    string
		joinMachine       *jm
		joinConfiguration *kubeadmv1beta1.JoinConfiguration
	}{
		{
			name: "test node registration override values",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{
						Network: v1alpha1.Network{
							APIServerELB: v1alpha1.ClassicELB{
								DNSName: "test.test",
							},
						},
					},
				},
				Machine: &clusterv1.Machine{},
			},
			joinConfiguration: &kubeadmv1beta1.JoinConfiguration{
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:      "cat",
					CRISocket: "unix:///var/run/docker/docker.sock",
					KubeletExtraArgs: map[string]string{
						"cloud-provider": "gcp",
					},
				},
			},
		},
		{
			name: "should set the cloud-provider extra args even if args are nil",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
				},
				Machine: &clusterv1.Machine{},
			},
			joinConfiguration: &kubeadmv1beta1.JoinConfiguration{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := kubeadm.SetJoinNodeConfigurationOverrides(tc.caCertHash, tc.bootstrapToken, tc.joinMachine, tc.joinConfiguration)

			if tc.joinConfiguration.NodeRegistration.Name != kubeadm.HostnameLookup &&
				tc.joinConfiguration.NodeRegistration.Name == out.NodeRegistration.Name {
				t.Fatal("did not properly override the NodeRegistration.Name")
			}

			if tc.joinConfiguration.NodeRegistration.CRISocket != kubeadm.ContainerdSocket &&
				tc.joinConfiguration.NodeRegistration.CRISocket == out.NodeRegistration.CRISocket {
				t.Fatal("did not properly override the CRISocket")
			}

			if out.NodeRegistration.KubeletExtraArgs["cloud-provider"] != kubeadm.CloudProvider {
				t.Fatal("did not properly set the cloud-provider on the kubelet extra args")
			}
		})
	}
}

func contains(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}
