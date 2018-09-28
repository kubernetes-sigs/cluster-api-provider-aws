// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
)

func TestCreateInstanceReal(t *testing.T) {
	testcases := []struct {
		name    string
		cluster clusterv1.Cluster
		machine clusterv1.Machine
	}{
		{
			name: "controlplane",
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: clusterv1.ClusterNetworkingConfig{
						ServiceDomain: "cluster.local",
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{
								"10.96.0.0/12",
							},
						},
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{
								"192.168.0.0/16",
							},
						},
					},
					ProviderConfig: clusterv1.ProviderConfig{
						Value: &runtime.RawExtension{
							Raw: []byte(`apiVersion: "awsproviderconfig/v1alpha1"
kind: AWSClusterProviderConfig
`),
						},
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: &runtime.RawExtension{
						Raw: []byte(`apiVersion: "awsproviderconfig/v1alpha1"
kind: AWSClusterProviderStatus
network:
  subnets:
  - id: subnet-97344ab8
    public: false
`),
					},
				},
			},
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-controlplane-1",
				},
				Spec: clusterv1.MachineSpec{
					Versions: clusterv1.MachineVersionInfo{
						Kubelet:      "v1.11.3",
						ControlPlane: "v1.11.3",
					},
					ProviderConfig: clusterv1.ProviderConfig{
						Value: &runtime.RawExtension{
							Raw: []byte(`apiVersion: "awsproviderconfig/v1alpha1"
kind: AWSMachineProviderConfig
nodeRole: "controlplane"
`),
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))
			ec2client := ec2.New(sess)
			s := ec2svc.NewService(ec2client)

			codec, err := providerconfigv1.NewCodec()
			if err != nil {
				t.Fatalf("Could not create codec: %v", err)
			}

			clusterProviderCfg := &providerconfigv1.AWSClusterProviderConfig{}
			err = codec.DecodeFromProviderConfig(tc.cluster.Spec.ProviderConfig, clusterProviderCfg)
			if err != nil {
				t.Fatalf("Could not decode cluster providerConfig: %v", err)
			}

			clusterStatus := &providerconfigv1.AWSClusterProviderStatus{}
			err = codec.DecodeProviderStatus(tc.cluster.Status.ProviderStatus, clusterStatus)
			if err != nil {
				t.Fatalf("Could not decode cluster providerStatus: %v", err)
			}

			machineProviderCfg := &providerconfigv1.AWSMachineProviderConfig{}
			err = codec.DecodeFromProviderConfig(tc.machine.Spec.ProviderConfig, machineProviderCfg)
			if err != nil {
				t.Fatalf("Could not decode machine providerConfig: %v", err)
			}

			instance, err := s.CreateInstance(&tc.machine, machineProviderCfg, &tc.cluster, clusterProviderCfg, clusterStatus)
			if err != nil {
				t.Fatalf("Failed to create instance: %v", err)
			}
			t.Logf("Instance: %v", instance)
		})
	}
}
