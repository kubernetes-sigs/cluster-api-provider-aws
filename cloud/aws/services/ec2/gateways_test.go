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

package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

func TestReconcileInternetGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.Network
		expect func(m *mock_ec2iface.MockEC2API)
	}{
		{
			name: "has igw",
			input: &v1alpha1.Network{
				VPC: v1alpha1.VPC{
					ID: "vpc-gateways",
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeInternetGateways(gomock.AssignableToTypeOf(&ec2.DescribeInternetGatewaysInput{})).
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
			},
		},
		{
			name: "no igw attached, creates one",
			input: &v1alpha1.Network{
				VPC: v1alpha1.VPC{
					ID: "vpc-gateways",
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeInternetGateways(gomock.AssignableToTypeOf(&ec2.DescribeInternetGatewaysInput{})).
					Return(&ec2.DescribeInternetGatewaysOutput{}, nil)

				m.EXPECT().
					CreateInternetGateway(gomock.AssignableToTypeOf(&ec2.CreateInternetGatewayInput{})).
					Return(&ec2.CreateInternetGatewayOutput{
						InternetGateway: &ec2.InternetGateway{InternetGatewayId: aws.String("igw-1")},
					}, nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.EXPECT().
					AttachInternetGateway(gomock.Eq(&ec2.AttachInternetGatewayInput{
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
			tc.expect(ec2Mock)

			s := NewService(ec2Mock)
			if err := s.reconcileInternetGateways("test-cluster", tc.input); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
