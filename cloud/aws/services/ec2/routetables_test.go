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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

func TestReconcileRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.Network
		expect func(m *mock_ec2iface.MockEC2API)
		err    error
	}{
		{
			name: "no routes existing, single private and single public, same AZ",
			input: &v1alpha1.Network{
				InternetGatewayID: aws.String("igw-01"),
				VPC: v1alpha1.VPC{
					ID: "vpc-routetables",
				},
				Subnets: v1alpha1.Subnets{
					&v1alpha1.Subnet{
						VpcID:            "vpc-routetables",
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&v1alpha1.Subnet{
						VpcID:            "vpc-routetables",
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.EXPECT().
					CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.EXPECT().
					CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
						NatGatewayId:         aws.String("nat-01"),
						DestinationCidrBlock: aws.String("0.0.0.0/0"),
						RouteTableId:         aws.String("rt-1"),
					})).
					After(privateRouteTable)

				m.EXPECT().
					AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
						RouteTableId: aws.String("rt-1"),
						SubnetId:     aws.String("subnet-routetables-private"),
					})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.EXPECT().
					CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.EXPECT().
					CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
						GatewayId:            aws.String("igw-01"),
						DestinationCidrBlock: aws.String("0.0.0.0/0"),
						RouteTableId:         aws.String("rt-2"),
					})).
					After(publicRouteTable)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.EXPECT().
					AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
						RouteTableId: aws.String("rt-2"),
						SubnetId:     aws.String("subnet-routetables-public"),
					})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(publicRouteTable)
			},
		},
		{
			name: "subnets in different availability zones, returns error",
			input: &v1alpha1.Network{
				InternetGatewayID: aws.String("igw-01"),
				VPC: v1alpha1.VPC{
					ID: "vpc-routetables",
				},
				Subnets: v1alpha1.Subnets{
					&v1alpha1.Subnet{
						VpcID:            "vpc-routetables",
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&v1alpha1.Subnet{
						VpcID:            "vpc-routetables",
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1b",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)
			},
			err: errors.New(`no nat gateways are available in availability zone "us-east-1a"`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)

			s := NewService(ec2Mock)
			if err := s.reconcileRouteTables("test-cluster", tc.input); err != nil && tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("was expecting error to look like '%v', but got '%v'", tc.err, err)
				}
			} else if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
