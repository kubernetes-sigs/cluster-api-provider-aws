/*
Copyright 2018 The Kubernetes Authors.

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

package deployer_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2" // nolint
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/deployer"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha2"
)

type scopeGetter struct {
	actuators.AWSClients
}

func (s *scopeGetter) ClusterScope(params actuators.ClusterScopeParams) (*actuators.ClusterScope, error) {
	params.AWSClients = s.AWSClients
	return actuators.NewClusterScope(params)
}

func TestGetIP(t *testing.T) {
	testcases := []struct {
		name       string
		cluster    *clusterv1.Cluster
		expectedIP string
		elbExpects func(*mock_elbiface.MockELBAPIMockRecorder)
	}{
		{
			name: "sunny day test",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test", Namespace: "default"},
				Spec: clusterv1.ClusterSpec{
					ProviderSpec: clusterv1.ProviderSpec{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderSpec{}),
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderStatus{}),
				},
			},
			expectedIP: "something",
			elbExpects: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []*string{aws.String("test-apiserver")},
				}).Return(&elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							Scheme:  aws.String("internet-facing"),
							DNSName: aws.String("something"),
							VPCId:   aws.String("test-vpc"),
						},
					},
				}, nil)

				m.DescribeLoadBalancerAttributes(gomock.Any()).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{},
					}, nil)
			},
		},
		{
			name: "lookup IP if the status is empty",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test", Namespace: "default"},
				Spec: clusterv1.ClusterSpec{
					ProviderSpec: clusterv1.ProviderSpec{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderSpec{}),
					},
				},
			},
			expectedIP: "dunno",
			elbExpects: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
					LoadBalancerNames: []*string{aws.String("test-apiserver")},
				}).Return(&elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							Scheme:  aws.String("internet-facing"),
							DNSName: aws.String("dunno"),
							VPCId:   aws.String("test-vpc"),
						},
					},
				}, nil)

				m.DescribeLoadBalancerAttributes(gomock.Any()).
					Return(&elb.DescribeLoadBalancerAttributesOutput{
						LoadBalancerAttributes: &elb.LoadBalancerAttributes{},
					}, nil)
			},
		},
		{
			name: "return the IP if it is stored in the status",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test", Namespace: "default"},
				Spec: clusterv1.ClusterSpec{
					ProviderSpec: clusterv1.ProviderSpec{
						Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderSpec{}),
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

			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			deployer := deployer.New(deployer.Params{ClusterScopeGetter: &scopeGetter{
				actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			}})

			if tc.elbExpects != nil {
				tc.elbExpects(elbMock.EXPECT())
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
