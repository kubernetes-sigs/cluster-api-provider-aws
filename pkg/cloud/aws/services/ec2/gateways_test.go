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

package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface" //nolint
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestReconcileInternetGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name: "has igw",
			input: &v1alpha1.NetworkSpec{
				VPC: v1alpha1.VPCSpec{
					ID: "vpc-gateways",
					Tags: v1alpha1.Tags{
						v1alpha1.ClusterTagKey("test-cluster"): "owned",
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
			input: &v1alpha1.NetworkSpec{
				VPC: v1alpha1.VPCSpec{
					ID: "vpc-gateways",
					Tags: v1alpha1.Tags{
						v1alpha1.ClusterTagKey("test-cluster"): "owned",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeInternetGateways(gomock.AssignableToTypeOf(&ec2.DescribeInternetGatewaysInput{})).
					Return(&ec2.DescribeInternetGatewaysOutput{}, nil)

				m.CreateInternetGateway(gomock.AssignableToTypeOf(&ec2.CreateInternetGatewayInput{})).
					Return(&ec2.CreateInternetGatewayOutput{
						InternetGateway: &ec2.InternetGateway{InternetGatewayId: aws.String("igw-1")},
					}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

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
			elbMock := mock_elbiface.NewMockELBAPI(mockCtrl)

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			scope.ClusterConfig = &v1alpha1.AWSClusterProviderSpec{
				NetworkSpec: *tc.input,
			}

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			if err := s.reconcileInternetGateways(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
