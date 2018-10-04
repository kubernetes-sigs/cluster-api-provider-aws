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

const (
	subnetsVPCID = "vpc-subnets"
)

func TestReconcileSubnets(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.Network
		expect func(m *mock_ec2iface.MockEC2API)
	}{
		{
			name: "single private subnet exists, should create public with defaults",
			input: &v1alpha1.Network{
				VPC: v1alpha1.VPC{ID: subnetsVPCID},
				Subnets: []*v1alpha1.Subnet{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.10.0/24",
						IsPublic:         false,
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeAvailabilityZones(gomock.AssignableToTypeOf(&ec2.DescribeAvailabilityZonesInput{})).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								RegionName: aws.String("us-east-1"),
								ZoneName:   aws.String("us-east-1a"),
							},
						},
					}, nil)

				m.EXPECT().
					DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("tag-key"),
								Values: []*string{aws.String("kubernetes.io/cluster/test-cluster")},
							},
						},
					})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							&ec2.Subnet{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.EXPECT().
					CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
						VpcId:            aws.String(subnetsVPCID),
						CidrBlock:        aws.String(defaultPublicSubnetCidr),
						AvailabilityZone: aws.String("us-east-1a"),
					})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.1.0.0/16"),
							AvailabilityZone:    aws.String("us-east-1a"),
							MapPublicIpOnLaunch: aws.Bool(true),
						},
					}, nil)

				m.EXPECT().
					WaitUntilSubnetAvailable(gomock.Any())

				m.EXPECT().
					CreateTags(gomock.Eq(&ec2.CreateTagsInput{
						Resources: aws.StringSlice([]string{"subnet-2"}),
						Tags: []*ec2.Tag{&ec2.Tag{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						}},
					})).
					Return(nil, nil)

				m.EXPECT().
					ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
						MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
							Value: aws.Bool(true),
						},
						SubnetId: aws.String("subnet-2"),
					}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil)

			},
		},
		{
			name: "no subnet exist, create private and public",
			input: &v1alpha1.Network{
				VPC: v1alpha1.VPC{ID: subnetsVPCID},
				Subnets: []*v1alpha1.Subnet{
					{
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.1.0.0/16",
						IsPublic:         false,
					},
					{
						AvailabilityZone: "us-east-1b",
						CidrBlock:        "10.2.0.0/16",
						IsPublic:         true,
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				describeCall := m.EXPECT().
					DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("tag-key"),
								Values: []*string{aws.String("kubernetes.io/cluster/test-cluster")},
							},
						},
					})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				firstSubnet := m.EXPECT().
					CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
						VpcId:            aws.String(subnetsVPCID),
						CidrBlock:        aws.String("10.1.0.0/16"),
						AvailabilityZone: aws.String("us-east-1a"),
					})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.1.0.0/16"),
							AvailabilityZone:    aws.String("us-east-1a"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.EXPECT().
					WaitUntilSubnetAvailable(gomock.Any()).
					After(firstSubnet)

				m.EXPECT().
					CreateTags(gomock.Eq(&ec2.CreateTagsInput{
						Resources: aws.StringSlice([]string{"subnet-1"}),
						Tags: []*ec2.Tag{&ec2.Tag{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						}},
					})).
					Return(nil, nil)

				secondSubnet := m.EXPECT().
					CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
						VpcId:            aws.String(subnetsVPCID),
						CidrBlock:        aws.String("10.2.0.0/16"),
						AvailabilityZone: aws.String("us-east-1b"),
					})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.2.0.0/16"),
							AvailabilityZone:    aws.String("us-east-1a"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				m.EXPECT().
					WaitUntilSubnetAvailable(gomock.Any()).
					After(secondSubnet)

				m.EXPECT().
					CreateTags(gomock.Eq(&ec2.CreateTagsInput{
						Resources: aws.StringSlice([]string{"subnet-2"}),
						Tags: []*ec2.Tag{&ec2.Tag{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("owned"),
						}},
					})).
					Return(nil, nil)

				m.EXPECT().
					ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
						MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
							Value: aws.Bool(true),
						},
						SubnetId: aws.String("subnet-2"),
					}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(secondSubnet)

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)

			s := NewService(ec2Mock)
			if err := s.reconcileSubnets("test-cluster", tc.input); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
