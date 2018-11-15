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

package deployer_test

import (
	"testing"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"

	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1" // nolint
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/mocks"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestGetIP(t *testing.T) {
	testcases := []struct {
		name       string
		cluster    *clusterv1.Cluster
		expectedIP string
		elbExpects func(*mocks.MockELBInterface)
	}{
		{
			name: "sunny day test",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderConfig{}),
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderStatus{}),
				},
			},
			expectedIP: "something",
			elbExpects: func(m *mocks.MockELBInterface) {
				m.EXPECT().
					GetAPIServerDNSName("test").
					Return("something", nil)
			},
		},
		{
			name: "lookup IP if the status is empty",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderConfig{}),
					},
				},
			},
			expectedIP: "dunno",
			elbExpects: func(m *mocks.MockELBInterface) {
				m.EXPECT().
					GetAPIServerDNSName("test").
					Return("dunno", nil)
			},
		},
		{
			name: "return the IP if it is stored in the status",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderConfig{}),
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderStatus{
						Network: providerv1.Network{
							APIServerELB: providerv1.ClassicELB{
								DNSName: "banana",
							},
						},
					}),
				},
			},
			expectedIP: "banana",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			services := mocks.NewSDKGetter(mockCtrl)
			deployer := deployer.New(services)

			if tc.elbExpects != nil {
				tc.elbExpects(services.ELBMock)
			}

			ip, err := deployer.GetIP(tc.cluster, nil)
			if err != nil {
				t.Fatalf("failed to get API server address: %v", err)
			}

			if ip != tc.expectedIP {
				t.Fatalf("got the wrong IP. Found %v, wanted %v", ip, tc.expectedIP)
			}
		})
	}
}
