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

package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestGenerateELBName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "test",
			expected: "test-apiserver",
		},
		{
			name:     "0123456789012345678901",
			expected: "0123456789012345678901-apiserver",
		},
		{
			name:     "01234567890123456789012",
			expected: "26o3cjil5at5qn27vukn5x09b3ql-k8s",
		},
		{
			name:     "anotherverylongtoolongname",
			expected: "t8gnrbbifaaf5d0k4xmwui3xwvip-k8s",
		},
		{
			name:     "anotherverylongtoolongnameanotherverylongtoolongname",
			expected: "tph1huzox1f10z9ow1inrootjws8-k8s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elbName, err := GenerateELBName(tt.name)
			if err != nil {
				t.Error(err)
			}

			if elbName != tt.expected {
				t.Errorf("expected ELB name: %v, got name: %v", tt.expected, elbName)
			}

			if len(elbName) > 32 {
				t.Errorf("ELB name too long: %v vs. %s", len(elbName), "32")
			}
		})
	}
}

func TestGetAPIServerClassicELBSpec_ControlPlaneLoadBalancer(t *testing.T) {
	tests := []struct {
		name   string
		lb     *infrav1.AWSLoadBalancerSpec
		mocks  func(m *mock_ec2iface.MockEC2APIMockRecorder)
		expect func(t *testing.T, res *infrav1.ClassicELB)
	}{
		{
			name:  "nil load balancer config",
			lb:    nil,
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if res.Attributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer not to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with cross zone enabled",
			lb: &infrav1.AWSLoadBalancerSpec{
				CrossZoneLoadBalancing: true,
			},
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if !res.Attributes.CrossZoneLoadBalancing {
					t.Error("Expected load balancer to have cross-zone load balancing enabled")
				}
			},
		},
		{
			name: "load balancer config with subnets specified",
			lb: &infrav1.AWSLoadBalancerSpec{
				Subnets: []string{"subnet-1", "subnet-2"},
			},
			mocks: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []*string{
						aws.String("subnet-1"),
						aws.String("subnet-2"),
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
							},
							{
								SubnetId:         aws.String("subnet-2"),
								AvailabilityZone: aws.String("us-east-1b"),
							},
						},
					}, nil)
			},
			expect: func(t *testing.T, res *infrav1.ClassicELB) {
				if len(res.SubnetIDs) != 2 {
					t.Errorf("Expected load balancer to be configured for 2 subnets, got %v", len(res.SubnetIDs))
				}
				if len(res.AvailabilityZones) != 2 {
					t.Errorf("Expected load balancer to be configured for 2 availability zones, got %v", len(res.AvailabilityZones))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "foo",
						Name:      "bar",
					},
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						ControlPlaneLoadBalancer: tc.lb,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			tc.mocks(ec2Mock.EXPECT())

			s := &Service{
				scope:     clusterScope,
				EC2Client: ec2Mock,
			}

			spec, err := s.getAPIServerClassicELBSpec()
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, spec)
		})
	}
}
