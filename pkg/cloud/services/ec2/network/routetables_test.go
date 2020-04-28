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
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestReconcileRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
		err    error
	}{
		{
			name: "no routes existing, single private and single public, same AZ",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-routetables",
					InternetGatewayID: aws.String("igw-01"),
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
					NatGatewayId:         aws.String("nat-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-1"),
					SubnetId:     aws.String("subnet-routetables-private"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.CreateRouteTable(gomock.Eq(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.AssociateRouteTable(gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-2"),
					SubnetId:     aws.String("subnet-routetables-public"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(publicRouteTable)
			},
		},
		{
			name: "subnets in different availability zones, returns error",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					InternetGatewayID: aws.String("igw-01"),
					ID:                "vpc-routetables",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1b",
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)
			},
			err: errors.New(`no nat gateways available in "us-east-1a"`),
		},
		{
			name: "routes exist, but the nat gateway ID is incorrect, replaces it",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					InternetGatewayID: aws.String("igw-01"),
					ID:                "vpc-routetables",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					&infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
						RouteTableID:     aws.String("route-table-1"),
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{
							{
								RouteTableId: aws.String("route-table-private"),
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-private"),
									},
								},
								Routes: []*ec2.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										NatGatewayId:         aws.String("outdated-nat-01"),
									},
								},
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-private"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
							{
								RouteTableId: aws.String("route-table-public"),
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-public"),
									},
								},
								Routes: []*ec2.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										GatewayId:            aws.String("igw-01"),
									},
								},
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-public"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.ReplaceRoute(gomock.Eq(
					&ec2.ReplaceRouteInput{
						DestinationCidrBlock: aws.String("0.0.0.0/0"),
						RouteTableId:         aws.String("route-table-private"),
						NatGatewayId:         aws.String("nat-01"),
					},
				)).
					Return(nil, nil)
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
						NetworkSpec: *tc.input,
					},
				},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope.NetworkScope)
			if err := s.reconcileRouteTables(); err != nil && tc.err != nil {
				if !strings.Contains(err.Error(), tc.err.Error()) {
					t.Fatalf("was expecting error to look like '%v', but got '%v'", tc.err, err)
				}
			} else if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
