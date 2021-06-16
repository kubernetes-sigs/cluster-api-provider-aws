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

package network

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

func TestReconcileInternetGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name: "has igw",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-gateways",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInternetGateways(gomock.AssignableToTypeOf(&ec2.DescribeInternetGatewaysInput{})).
					Return(&ec2.DescribeInternetGatewaysOutput{
						InternetGateways: []*ec2.InternetGateway{
							{
								InternetGatewayId: aws.String("igw-0"),
								Attachments: []*ec2.InternetGatewayAttachment{
									{
										State: aws.String(ec2.AttachmentStatusAttached),
										VpcId: aws.String("vpc-gateways"),
									},
								},
							},
						},
					}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "no igw attached, creates one",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: "vpc-gateways",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInternetGateways(gomock.AssignableToTypeOf(&ec2.DescribeInternetGatewaysInput{})).
					Return(&ec2.DescribeInternetGatewaysOutput{}, nil)

				m.CreateInternetGateway(gomock.AssignableToTypeOf(&ec2.CreateInternetGatewayInput{})).
					Return(&ec2.CreateInternetGatewayOutput{
						InternetGateway: &ec2.InternetGateway{
							InternetGatewayId: aws.String("igw-1"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String(infrav1.ClusterTagKey("test-cluster")),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-igw"),
								},
							},
						},
					}, nil)

				m.AttachInternetGateway(gomock.Eq(&ec2.AttachInternetGatewayInput{
					InternetGatewayId: aws.String("igw-1"),
					VpcId:             aws.String("vpc-gateways"),
				})).
					Return(&ec2.AttachInternetGatewayOutput{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec: infrav1.AWSClusterSpec{
						NetworkSpec: *tc.input,
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			if err := s.reconcileInternetGateways(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
