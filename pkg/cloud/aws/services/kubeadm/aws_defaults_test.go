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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmv1beta1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/kubeadm"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest"
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
		{
			name: "test taint values are passed through",
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
				Machine: &clusterv1.Machine{
					Spec: clusterv1.MachineSpec{
						Taints: []corev1.Taint{
							{
								Key:    "testKey",
								Value:  "someValue",
								Effect: "someEffect",
							},
							{
								Key:   "key",
								Value: "Value",
							},
						},
					},
				},
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

			if len(out.NodeRegistration.Taints) != len(tc.joinMachine.Spec.Taints) {
				t.Fatalf("The config and the machine should have the same number of taints: actual %v and expected %v", out.NodeRegistration.Taints, tc.joinMachine.Spec.Taints)
			}

			for i, v := range out.NodeRegistration.Taints {
				if tc.joinMachine.Spec.Taints[i] != v {
					t.Fatalf("Expected a taint of %v but got %v instead", tc.joinMachine.Spec.Taints[i], v)
				}
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
					Logger: &cloudtest.Log{},
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
					Logger:        &cloudtest.Log{},
				},
				Machine: &clusterv1.Machine{},
			},
			joinConfiguration: &kubeadmv1beta1.JoinConfiguration{},
		},
		{
			name: "should append nodeRole label to provided labels",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
					Logger:        &cloudtest.Log{},
				},
				Machine: &clusterv1.Machine{},
			},
			joinConfiguration: &kubeadmv1beta1.JoinConfiguration{
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					KubeletExtraArgs: map[string]string{
						"node-labels": "foo=bar,pizza=pepperoni",
					},
				},
			},
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

			if _, ok := tc.joinConfiguration.NodeRegistration.KubeletExtraArgs["node-labels"]; ok && out.NodeRegistration.KubeletExtraArgs["node-labels"] != "foo=bar,pizza=pepperoni,node-role.kubernetes.io/node=" {
				t.Fatal("did not properly append nodeRole label")
			}
		})
	}
}

func TestSetControlPlaneJoinConfigurationOverrides(t *testing.T) {
	testcases := []struct {
		name       string
		joinconfig *kubeadmv1beta1.JoinConfiguration
	}{
		{
			name: "simple override test",
			joinconfig: &kubeadmv1beta1.JoinConfiguration{
				ControlPlane: &kubeadmv1beta1.JoinControlPlane{
					LocalAPIEndpoint: kubeadmv1beta1.APIEndpoint{
						AdvertiseAddress: "fake.endpoint",
						BindPort:         420,
					},
				},
			},
		},
		{
			name: "nil join config",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := kubeadm.SetControlPlaneJoinConfigurationOverrides(tc.joinconfig)

			// ignore assertions if nil was passed in
			if tc.joinconfig == nil {
				return
			}

			// Assertion assumes the test case does not start with the override value
			if tc.joinconfig.ControlPlane.LocalAPIEndpoint.AdvertiseAddress == out.ControlPlane.LocalAPIEndpoint.AdvertiseAddress {
				t.Fatal("did not properly override the local api endpoint advertise address")
			}
			if tc.joinconfig.ControlPlane.LocalAPIEndpoint.BindPort == out.ControlPlane.LocalAPIEndpoint.BindPort {
				t.Fatal("did not properly override the local api endpoint bind port")
			}
		})
	}
}

func TestSetDefaultClusterConfiguration(t *testing.T) {
	testcases := []struct {
		name                 string
		joinMachine          *jm
		clusterConfiguration *kubeadmv1beta1.ClusterConfiguration
	}{
		{
			name: "simple",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
				},
				Machine: &clusterv1.Machine{},
			},
			clusterConfiguration: &kubeadmv1beta1.ClusterConfiguration{},
		},
		{
			name: "handle nil",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
				},
				Machine: &clusterv1.Machine{},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := kubeadm.SetDefaultClusterConfiguration(tc.joinMachine, tc.clusterConfiguration)

			// Ignore the assertions of the base is nil
			if tc.clusterConfiguration == nil {
				return
			}

			// Assertion: set a control plane only if the user hasn't specified one
			if tc.clusterConfiguration.ControlPlaneEndpoint == "" {
				if out.ControlPlaneEndpoint == "" {
					t.Fatal("Failed to set a control plane endpoint when one was not provided")
				}
			}

			if tc.clusterConfiguration.ControlPlaneEndpoint != "" && out.ControlPlaneEndpoint != tc.clusterConfiguration.ControlPlaneEndpoint {
				t.Fatal("Overrode ControlPlaneEndpoint. This user supplied value should not be overridden")
			}

			// Assertion: always add local ipv4 lookup and the ELB's DNS name to the CertSANs
			if len(out.APIServer.CertSANs)-len(tc.clusterConfiguration.APIServer.CertSANs) != 2 {
				t.Fatal("expected one additional CertSAN in the output")
			}
		})
	}
}

func TestSetInitConfigurationOverrides(t *testing.T) {
	testcasts := []struct {
		name              string
		joinMachine       *jm
		initConfiguration *kubeadmv1beta1.InitConfiguration
	}{
		// Assumption: Do not use any default values as input values for testcases.
		{
			name: "assertions pass with an empty configuration",
			joinMachine: &jm{
				Scope: &actuators.Scope{},
			},
			initConfiguration: &kubeadmv1beta1.InitConfiguration{},
		},
		{
			name:        "nil configuration should not break",
			joinMachine: &jm{},
		},
		{
			name: "overrides actually work",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					Logger: &cloudtest.Log{},
				},
			},
			initConfiguration: &kubeadmv1beta1.InitConfiguration{
				NodeRegistration: kubeadmv1beta1.NodeRegistrationOptions{
					Name:      "moonshine",
					CRISocket: "/dev/null",
					KubeletExtraArgs: map[string]string{
						"cloud-provider": "gcp",
					},
				},
			},
		},
	}

	for _, tc := range testcasts {
		t.Run(tc.name, func(t *testing.T) {
			out := kubeadm.SetInitConfigurationOverrides(tc.joinMachine, tc.initConfiguration)

			// Ignore assertions if the initConfiguration is nil
			if tc.initConfiguration == nil {
				return
			}

			// Assertion: Override the user provided value, if any, to the dynamic lookup value
			if tc.initConfiguration.NodeRegistration.Name != "" && tc.initConfiguration.NodeRegistration.Name == out.NodeRegistration.Name {
				t.Fatal("Provided node name was not properly overridden")
			}

			// Assertion: Override the user provided CRI socket with containerd's CRI socket
			if tc.initConfiguration.NodeRegistration.CRISocket != "" && tc.initConfiguration.NodeRegistration.CRISocket == out.NodeRegistration.CRISocket {
				t.Fatal("Provided CRISocket was not properly overridden")
			}

			// Assertion: Kubelet extra args are populated with the aws cloud provier
			if out.NodeRegistration.KubeletExtraArgs["cloud-provider"] != kubeadm.CloudProvider {
				t.Fatal("Cloud provider was not successfully set to aws")
			}

		})
	}
}

func TestSetClusterConfigurationOverrides(t *testing.T) {
	testcases := []struct {
		name                 string
		joinMachine          *jm
		clusterConfiguration *kubeadmv1beta1.ClusterConfiguration
	}{
		{
			name: "simple test",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "hello",
						},
						Spec: clusterv1.ClusterSpec{
							ClusterNetwork: clusterv1.ClusterNetworkingConfig{
								ServiceDomain: "hogwarts",
								Pods: clusterv1.NetworkRanges{
									CIDRBlocks: []string{"123/127"},
								},
								Services: clusterv1.NetworkRanges{
									CIDRBlocks: []string{"456/128"},
								},
							},
						},
					},
					Logger: &cloudtest.Log{},
				},
				Machine: &clusterv1.Machine{},
			},
			clusterConfiguration: &kubeadmv1beta1.ClusterConfiguration{
				ClusterName: "some test",
				APIServer: kubeadmv1beta1.APIServer{
					ControlPlaneComponent: kubeadmv1beta1.ControlPlaneComponent{
						ExtraArgs: map[string]string{"cloud-provider": "not-aws"},
					},
				},
				ControllerManager: kubeadmv1beta1.ControlPlaneComponent{
					ExtraArgs: map[string]string{"cloud-provider": "azure"},
				},
			},
		},
		{
			name: "nil cluster configuration test",
			joinMachine: &jm{
				Scope: &actuators.Scope{
					ClusterStatus: &v1alpha1.AWSClusterProviderStatus{},
					Cluster: &clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "hello",
						},
						Spec: clusterv1.ClusterSpec{
							ClusterNetwork: clusterv1.ClusterNetworkingConfig{
								ServiceDomain: "hogwarts",
								Pods: clusterv1.NetworkRanges{
									CIDRBlocks: []string{"123/127"},
								},
								Services: clusterv1.NetworkRanges{
									CIDRBlocks: []string{"456/128"},
								},
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

			out := kubeadm.SetClusterConfigurationOverrides(tc.joinMachine, tc.clusterConfiguration)
			// Assertion: Forces apiserver cloud-provider to aws
			if out.APIServer.ControlPlaneComponent.ExtraArgs["cloud-provider"] != kubeadm.CloudProvider {
				t.Fatal("cloud-provider argument was not set properly on the apiserver")
			}

			// Assertion: Forces controller-manager cloud-provider to aws
			if out.ControllerManager.ExtraArgs["cloud-provider"] != kubeadm.CloudProvider {
				t.Fatal("cloud-provider argument was not set properly on the controller-manager")
			}

			// Assertion: Sets the kubeadm cluster name to match the cluster-api cluster name
			if out.ClusterName != tc.joinMachine.GetScope().Name() {
				t.Fatal("The cluster name was not set properly")
			}

			// Assertion: Sets necessary networking attributes from the ClusterSpec
			if out.Networking.DNSDomain != tc.joinMachine.GetScope().Cluster.Spec.ClusterNetwork.ServiceDomain {
				t.Fatal("Failed to set the DNSDomain properly")
			}
			if out.Networking.PodSubnet != tc.joinMachine.GetScope().Cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0] {
				t.Fatal("Failed to set the pod subnet")
			}
			if out.Networking.ServiceSubnet != tc.joinMachine.GetScope().Cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0] {
				t.Fatal("Failed to set the service subnet")
			}

			// Assertion: Sets the Kubernetes version from the ClusterSpec
			if out.KubernetesVersion != tc.joinMachine.GetMachine().Spec.Versions.ControlPlane {
				t.Fatal("Failed to set the kubernetes version properly")
			}
		})
	}
}

/** Helpers **/

func contains(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}
