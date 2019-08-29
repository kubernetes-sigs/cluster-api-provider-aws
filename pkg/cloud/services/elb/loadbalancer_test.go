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

package elb

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
)

func TestDeleteLoadBalancers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		expect func(m *mock_elbiface.MockELBAPIMockRecorder)
	}{
		{
			name: "elb exists, it deletes it",
			expect: func(m *mock_elbiface.MockELBAPIMockRecorder) {
				m.DescribeLoadBalancers(gomock.Any()).Return(&elb.DescribeLoadBalancersOutput{
					LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
						{
							LoadBalancerName: aws.String("test-cluster-apiserver"),
							VPCId:            aws.String("test-vpc"),
							Scheme:           aws.String(string(infrav1.ClassicELBSchemeInternetFacing)),
						},
					},
				}, nil)

				m.DescribeLoadBalancerAttributes(gomock.Any()).Return(&elb.DescribeLoadBalancerAttributesOutput{
					LoadBalancerAttributes: &elb.LoadBalancerAttributes{},
				}, nil)

				m.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
					LoadBalancerName: aws.String("test-cluster-apiserver"),
				})

				m.DescribeLoadBalancers(gomock.Any()).Return(&elb.DescribeLoadBalancersOutput{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSClients: scope.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
				AWSCluster: &infrav1.AWSCluster{
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: infrav1.NetworkSpec{
							VPC: infrav1.VPCSpec{
								ID: "test-vpc",
							},
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(elbMock.EXPECT())
			s := NewService(scope)
			if err := s.DeleteLoadbalancers(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
