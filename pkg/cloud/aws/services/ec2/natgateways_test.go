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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
)

const (
	ElasticIPAllocationID = "elastic-ip-allocation-id"
)

func TestReconcileNatGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  []*v1alpha1.Subnet
		expect func(m *mock_ec2iface.MockEC2API)
	}{
		{
			name: "single private subnet exists, should create no NAT gateway",
			input: []*v1alpha1.Subnet{
				{
					ID:               "subnet-1",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().CreateNatGateway(gomock.Any()).Times(0)
			},
		},
		{
			name: "no private subnet exists, should create no NAT gateway",
			input: []*v1alpha1.Subnet{
				{
					ID:               "subnet-1",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {

				m.EXPECT().
					DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).Times(0)

				m.EXPECT().CreateNatGateway(gomock.Any()).Times(0)
			},
		},
		{
			name: "public & private subnet exists, should create 1 NAT gateway",
			input: []*v1alpha1.Subnet{
				{
					ID:               "subnet-1",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {

				m.EXPECT().
					DescribeNatGatewaysPages(
						gomock.Eq(&ec2.DescribeNatGatewaysInput{
							Filter: []*ec2.Filter{
								{
									Name:   aws.String("vpc-id"),
									Values: []*string{aws.String(subnetsVPCID)},
								},
								{
									Name:   aws.String("state"),
									Values: []*string{aws.String("available")},
								},
							},
						}),
						gomock.Any()).Return(nil)

				m.EXPECT().
					DescribeAddresses(gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.EXPECT().
					AllocateAddress(&ec2.AllocateAddressInput{Domain: aws.String("vpc")}).
					Return(&ec2.AllocateAddressOutput{
						AllocationId: aws.String(ElasticIPAllocationID),
					}, nil)

				m.EXPECT().
					CreateNatGateway(&ec2.CreateNatGatewayInput{
						AllocationId: aws.String(ElasticIPAllocationID),
						SubnetId:     aws.String("subnet-1"),
					}).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
					},
				}, nil)

				m.EXPECT().
					WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
						NatGatewayIds: []*string{aws.String("natgateway")},
					}).Return(nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "two public & 1 private subnet, and one NAT gateway exists",
			input: []*v1alpha1.Subnet{
				{
					ID:               "subnet-1",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
				{
					ID:               "subnet-3",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1b",
					CidrBlock:        "10.0.13.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {

				m.EXPECT().
					DescribeNatGatewaysPages(
						gomock.Eq(&ec2.DescribeNatGatewaysInput{
							Filter: []*ec2.Filter{
								{
									Name:   aws.String("vpc-id"),
									Values: []*string{aws.String(subnetsVPCID)},
								},
								{
									Name:   aws.String("state"),
									Values: []*string{aws.String("available")},
								},
							},
						}),
						gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool)
					funct(&ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{{
						NatGatewayId: aws.String("gateway"),
						SubnetId:     aws.String("subnet-1"),
					}}}, true)
				}).Return(nil)

				m.EXPECT().
					DescribeAddresses(gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.EXPECT().
					AllocateAddress(&ec2.AllocateAddressInput{Domain: aws.String("vpc")}).
					Return(&ec2.AllocateAddressOutput{
						AllocationId: aws.String(ElasticIPAllocationID),
					}, nil)

				m.EXPECT().
					CreateNatGateway(&ec2.CreateNatGatewayInput{
						AllocationId: aws.String(ElasticIPAllocationID),
						SubnetId:     aws.String("subnet-3"),
					}).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
					},
				}, nil)

				m.EXPECT().
					WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
						NatGatewayIds: []*string{aws.String("natgateway")},
					}).Return(nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "public & private subnet, and one NAT gateway exists",
			input: []*v1alpha1.Subnet{
				{
					ID:               "subnet-1",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {

				m.EXPECT().
					DescribeNatGatewaysPages(
						gomock.Eq(&ec2.DescribeNatGatewaysInput{
							Filter: []*ec2.Filter{
								{
									Name:   aws.String("vpc-id"),
									Values: []*string{aws.String(subnetsVPCID)},
								},
								{
									Name:   aws.String("state"),
									Values: []*string{aws.String("available")},
								},
							},
						}),
						gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool)
					funct(&ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{{
						NatGatewayId: aws.String("gateway"),
						SubnetId:     aws.String("subnet-1"),
					}}}, true)
				}).Return(nil)

				m.EXPECT().DescribeAddresses(gomock.Any()).Times(0)

				m.EXPECT().AllocateAddress(gomock.Any()).Times(0)

				m.EXPECT().CreateNatGateway(gomock.Any()).Times(0)
			},
		},
		{
			name: "public & private subnet declared, but doesn't exist yet",
			input: []*v1alpha1.Subnet{
				{
					ID:               "",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "",
					VpcID:            subnetsVPCID,
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)

			s := NewService(ec2Mock)
			if err := s.reconcileNatGateways("test-cluster", tc.input, &v1alpha1.VPC{ID: subnetsVPCID}); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
