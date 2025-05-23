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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestReconcileRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		expect func(m *mocks.MockEC2APIMockRecorder)
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
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					NatGatewayId:         aws.String("nat-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-1"),
					SubnetId:     aws.String("subnet-routetables-private"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-2"),
					SubnetId:     aws.String("subnet-routetables-public"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(publicRouteTable)
			},
		},
		{
			name: "no routes existing, single private and single public IPv6 enabled subnets, same AZ",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-routetables",
					InternetGatewayID: aws.String("igw-01"),
					IPv6: &infrav1.IPv6{
						EgressOnlyInternetGatewayID: aws.String("eigw-01"),
						CidrBlock:                   "2001:db8:1234::/56",
						PoolID:                      "my-pool",
					},
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:1234:1::/64",
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:1234:2::/64",
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					NatGatewayId:         aws.String("nat-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					DestinationIpv6CidrBlock:    aws.String("::/0"),
					EgressOnlyInternetGatewayId: aws.String("eigw-01"),
					RouteTableId:                aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-1"),
					SubnetId:     aws.String("subnet-routetables-private"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					DestinationIpv6CidrBlock: aws.String("::/0"),
					GatewayId:                aws.String("igw-01"),
					RouteTableId:             aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-2"),
					SubnetId:     aws.String("subnet-routetables-public"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(publicRouteTable)
			},
		},
		{
			name: "no routes existing, single private and single public IPv6 enabled subnets with existing Egress only IWG, same AZ",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:                "vpc-routetables",
					InternetGatewayID: aws.String("igw-01"),
					IPv6: &infrav1.IPv6{
						CidrBlock:                   "2001:db8:1234::/56",
						PoolID:                      "my-pool",
						EgressOnlyInternetGatewayID: aws.String("eigw-01"),
					},
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:1234:1::/64",
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						IsIPv6:           true,
						IPv6CidrBlock:    "2001:db8:1234:2::/64",
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					NatGatewayId:         aws.String("nat-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					DestinationIpv6CidrBlock:    aws.String("::/0"),
					EgressOnlyInternetGatewayId: aws.String("eigw-01"),
					RouteTableId:                aws.String("rt-1"),
				})).
					After(privateRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
					RouteTableId: aws.String("rt-1"),
					SubnetId:     aws.String("subnet-routetables-private"),
				})).
					Return(&ec2.AssociateRouteTableOutput{}, nil).
					After(privateRouteTable)

				publicRouteTable := m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					DestinationIpv6CidrBlock: aws.String("::/0"),
					GatewayId:                aws.String("igw-01"),
					RouteTableId:             aws.String("rt-2"),
				})).
					After(publicRouteTable)

				m.AssociateRouteTable(context.TODO(), gomock.Eq(&ec2.AssociateRouteTableInput{
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
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1b",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
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
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
						RouteTableID:     aws.String("route-table-1"),
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								RouteTableId: aws.String("route-table-private"),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-private"),
									},
								},
								Routes: []types.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										NatGatewayId:         aws.String("outdated-nat-01"),
									},
								},
								Tags: []types.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-private-us-east-1a"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
							{
								RouteTableId: aws.String("route-table-public"),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-public"),
									},
								},
								Routes: []types.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										GatewayId:            aws.String("igw-01"),
									},
								},
								Tags: []types.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-public-us-east-1a"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.ReplaceRoute(context.TODO(), gomock.Eq(
					&ec2.ReplaceRouteInput{
						DestinationCidrBlock: aws.String("0.0.0.0/0"),
						RouteTableId:         aws.String("route-table-private"),
						NatGatewayId:         aws.String("nat-01"),
					},
				)).
					Return(nil, nil)
			},
		},
		{
			name: "extra routes exist, do nothing",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					InternetGatewayID: aws.String("igw-01"),
					ID:                "vpc-routetables",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-private",
						IsPublic:         false,
						AvailabilityZone: "us-east-1a",
					},
					infrav1.SubnetSpec{
						ID:               "subnet-routetables-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								RouteTableId: aws.String("route-table-private"),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-private"),
									},
								},
								Routes: []types.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										NatGatewayId:         aws.String("nat-01"),
									},
									// Extra (managed outside of CAPA) route with Managed Prefix List destination.
									{
										DestinationPrefixListId: aws.String("pl-foobar"),
									},
								},
								Tags: []types.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-private-us-east-1a"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
							{
								RouteTableId: aws.String("route-table-public"),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-routetables-public"),
									},
								},
								Routes: []types.Route{
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										GatewayId:            aws.String("igw-01"),
									},
								},
								Tags: []types.Tag{
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("common"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-rt-public-us-east-1a"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)
			},
		},
		{
			name: "failed to create route, delete route table and fail",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					InternetGatewayID: aws.String("igw-01"),
					ID:                "vpc-rtbs",
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: infrav1.Subnets{
					infrav1.SubnetSpec{
						ID:               "subnet-rtbs-public",
						IsPublic:         true,
						NatGatewayID:     aws.String("nat-01"),
						AvailabilityZone: "us-east-1a",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.CreateRouteTable(context.TODO(), matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-rtbs")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &types.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

				m.CreateRoute(context.TODO(), gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-1"),
				})).
					Return(nil, awserrors.NewNotFound("MissingParameter"))

				m.DeleteRouteTable(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteRouteTableInput{})).
					Return(&ec2.DeleteRouteTableOutput{}, nil)
			},
			err: errors.New(`failed to create route in route table "rt-1"`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

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

// Delete Route Table(s).
var (
	stubEc2RouteTablePrivate = types.RouteTable{
		RouteTableId: aws.String("route-table-private"),
		Associations: []types.RouteTableAssociation{
			{
				SubnetId: nil,
			},
		},
		Routes: []types.Route{
			{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				NatGatewayId:         aws.String("outdated-nat-01"),
			},
		},
	}
	stubEc2RouteTablePublicWithAssociations = types.RouteTable{
		RouteTableId: aws.String("route-table-public"),
		Associations: []types.RouteTableAssociation{
			{
				SubnetId:                aws.String("subnet-routetables-public"),
				RouteTableAssociationId: aws.String("route-table-public"),
			},
		},
		Routes: []types.Route{
			{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				GatewayId:            aws.String("igw-01"),
			},
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("kubernetes.io/cluster/test-cluster"),
				Value: aws.String("owned"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
				Value: aws.String("common"),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-cluster-rt-public-us-east-1a"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Value: aws.String("owned"),
			},
		},
	}
)

func TestDeleteRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	describeRouteTableOutput := &ec2.DescribeRouteTablesOutput{
		RouteTables: []types.RouteTable{
			stubEc2RouteTablePrivate,
			stubEc2RouteTablePublicWithAssociations,
		},
	}

	testCases := []struct {
		name    string
		input   *infrav1.NetworkSpec
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name: "Should skip deletion if vpc is unmanaged",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:   "vpc-routetables",
					Tags: infrav1.Tags{},
				},
			},
		},
		{
			name:  "Should delete route table successfully",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(context.TODO(), gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)

				m.DisassociateRouteTable(context.TODO(), gomock.Eq(&ec2.DisassociateRouteTableInput{
					AssociationId: aws.String("route-table-public"),
				})).Return(&ec2.DisassociateRouteTableOutput{}, nil)

				m.DeleteRouteTable(context.TODO(), gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-public"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)
			},
		},
		{
			name:  "Should return error if describe route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
			wantErr: true,
		},
		{
			name:  "Should return error if delete route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(context.TODO(), gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(nil, awserrors.NewNotFound("not found"))
			},
			wantErr: true,
		},
		{
			name:  "Should return error if disassociate route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(context.TODO(), gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)

				m.DisassociateRouteTable(context.TODO(), gomock.Eq(&ec2.DisassociateRouteTableInput{
					AssociationId: aws.String("route-table-public"),
				})).Return(nil, awserrors.NewNotFound("not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

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
			g.Expect(err).NotTo(HaveOccurred())
			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			s := NewService(scope)
			s.EC2Client = ec2Mock

			err = s.deleteRouteTables()
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestDeleteRouteTable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name    string
		input   types.RouteTable
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
	}{
		{
			name:  "Should delete route table successfully",
			input: stubEc2RouteTablePrivate,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteRouteTable(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteRouteTableInput{})).
					Return(&ec2.DeleteRouteTableOutput{}, nil)
			},
		},
		{
			name:  "Should return error if delete route table fails",
			input: stubEc2RouteTablePrivate,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteRouteTable(context.TODO(), gomock.AssignableToTypeOf(&ec2.DeleteRouteTableInput{})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
			wantErr: true,
		},
		{
			name:  "Should return error if disassociate route table fails",
			input: stubEc2RouteTablePublicWithAssociations,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DisassociateRouteTable(context.TODO(), gomock.Eq(&ec2.DisassociateRouteTableInput{
					AssociationId: aws.String("route-table-public"),
				})).Return(nil, awserrors.NewNotFound("not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

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
					Spec:       infrav1.AWSClusterSpec{},
				},
			})
			g.Expect(err).NotTo(HaveOccurred())
			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			s := NewService(scope)
			s.EC2Client = ec2Mock

			err = s.deleteRouteTable(tc.input)
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

type routeTableInputMatcher struct {
	routeTableInput *ec2.CreateRouteTableInput
}

func (r routeTableInputMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*ec2.CreateRouteTableInput)
	if !ok {
		fmt.Println("heeeeyy")
		return false
	}
	if *actual.VpcId != *r.routeTableInput.VpcId {
		return false
	}

	return true
}

func (r routeTableInputMatcher) String() string {
	return fmt.Sprintf("partially matches %v", r.routeTableInput)
}

func matchRouteTableInput(input *ec2.CreateRouteTableInput) gomock.Matcher {
	return routeTableInputMatcher{routeTableInput: input}
}

func TestService_getRoutesForSubnet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	defaultSubnets := infrav1.Subnets{
		{
			ResourceID:       "subnet-az-2z-private",
			AvailabilityZone: "us-east-2z",
			IsPublic:         false,
		},
		{
			ResourceID:       "subnet-az-2z-public",
			AvailabilityZone: "us-east-2z",
			IsPublic:         true,
			NatGatewayID:     ptr.To("nat-gw-fromZone-us-east-2z"),
		},
		{
			ResourceID:       "subnet-az-1a-private",
			AvailabilityZone: "us-east-1a",
			IsPublic:         false,
		},
		{
			ResourceID:       "subnet-az-1a-public",
			AvailabilityZone: "us-east-1a",
			IsPublic:         true,
			NatGatewayID:     ptr.To("nat-gw-fromZone-us-east-1a"),
		},
		{
			ResourceID:       "subnet-lz-invalid2z-private",
			AvailabilityZone: "us-east-2-inv-1z",
			IsPublic:         false,
			ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
			ParentZoneName:   ptr.To("us-east-2a"),
		},
		{
			ResourceID:       "subnet-lz-invalid1a-public",
			AvailabilityZone: "us-east-2-nyc-1z",
			IsPublic:         true,
			ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
			ParentZoneName:   ptr.To("us-east-2z"),
		},
		{
			ResourceID:       "subnet-lz-1a-private",
			AvailabilityZone: "us-east-1-nyc-1a",
			IsPublic:         false,
			ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
			ParentZoneName:   ptr.To("us-east-1a"),
		},
		{
			ResourceID:       "subnet-lz-1a-public",
			AvailabilityZone: "us-east-1-nyc-1a",
			IsPublic:         true,
			ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
			ParentZoneName:   ptr.To("us-east-1a"),
		},
		{
			ResourceID:       "subnet-wl-invalid2z-private",
			AvailabilityZone: "us-east-2-wl1-inv-wlz-1",
			IsPublic:         false,
			ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
			ParentZoneName:   ptr.To("us-east-2z"),
		},
		{
			ResourceID:       "subnet-wl-invalid2z-public",
			AvailabilityZone: "us-east-2-wl1-inv-wlz-1",
			IsPublic:         true,
			ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
			ParentZoneName:   ptr.To("us-east-2z"),
		},
		{
			ResourceID:       "subnet-wl-1a-private",
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
			IsPublic:         false,
			ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
			ParentZoneName:   ptr.To("us-east-1a"),
		},
		{
			ResourceID:       "subnet-wl-1a-public",
			AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
			IsPublic:         true,
			ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
			ParentZoneName:   ptr.To("us-east-1a"),
		},
	}

	vpcName := "vpc-test-for-routes"
	defaultNetwork := infrav1.NetworkSpec{
		VPC: infrav1.VPCSpec{
			ID:                vpcName,
			InternetGatewayID: aws.String("vpc-igw"),
			CarrierGatewayID:  aws.String("vpc-cagw"),
			IPv6: &infrav1.IPv6{
				CidrBlock:                   "2001:db8:1234:1::/64",
				EgressOnlyInternetGatewayID: aws.String("vpc-eigw"),
			},
		},
		Subnets: defaultSubnets,
	}

	tests := []struct {
		name                string
		specOverrideNet     *infrav1.NetworkSpec
		specOverrideSubnets *infrav1.Subnets
		inputSubnet         *infrav1.SubnetSpec
		want                []*ec2.CreateRouteInput
		wantErr             bool
		wantErrMessage      string
	}{
		{
			name:                "empty subnet should have empty routes",
			specOverrideSubnets: &infrav1.Subnets{},
			inputSubnet: &infrav1.SubnetSpec{
				ID: "subnet-1-private",
			},
			want:           []*ec2.CreateRouteInput{},
			wantErrMessage: `no nat gateways available in "" for private subnet "subnet-1-private"`,
		},
		{
			name:           "empty subnet should have empty routes",
			inputSubnet:    &infrav1.SubnetSpec{},
			want:           []*ec2.CreateRouteInput{},
			wantErrMessage: `no nat gateways available in "" for private subnet ""`,
		},
		// public subnets ipv4
		{
			name: "public ipv4 subnet, availability zone, must have ipv4 default route to igw",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-public",
				AvailabilityZone: "us-east-1a",
				IsIPv6:           false,
				IsPublic:         true,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					GatewayId:            aws.String("vpc-igw"),
				},
			},
		},
		{
			name: "public ipv6 subnet, availability zone, must have ipv6 default route to igw",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-public",
				AvailabilityZone: "us-east-1a",
				IsPublic:         true,
				IsIPv6:           true,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					GatewayId:            aws.String("vpc-igw"),
				},
				{
					DestinationIpv6CidrBlock: aws.String("::/0"),
					GatewayId:                aws.String("vpc-igw"),
				},
			},
		},
		{
			name: "public ipv4 subnet, local zone, must have ipv4 default route to igw",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-public",
				AvailabilityZone: "us-east-1-nyc-1a",
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				IsPublic:         true,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					GatewayId:            aws.String("vpc-igw"),
				},
			},
		},
		{
			name: "public ipv4 subnet, wavelength zone, must have ipv4 default route to carrier gateway",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-public",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				IsPublic:         true,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					CarrierGatewayId:     aws.String("vpc-cagw"),
				},
			},
		},
		// public subnet ipv4, GW not found.
		{
			name: "public ipv4 subnet, availability zone, must return error when no internet gateway available",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				net.VPC.InternetGatewayID = nil
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-public",
				AvailabilityZone: "us-east-1a",
				IsPublic:         true,
			},
			wantErrMessage: `failed to create routing tables: internet gateway for VPC "vpc-test-for-routes" is not present`,
		},
		{
			name: "public ipv4 subnet, local zone, must return error when no internet gateway available",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				net.VPC.InternetGatewayID = nil
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-public",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         true,
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `failed to create routing tables: internet gateway for VPC "vpc-test-for-routes" is not present`,
		},
		{
			name: "public ipv4 subnet, wavelength zone, must return error when no Carrier Gateway found",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				net.VPC.CarrierGatewayID = nil
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-public",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				IsPublic:         true,
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `failed to create carrier routing table: carrier gateway for VPC "vpc-test-for-routes" is not present`,
		},
		// public subnet ipv6, unsupported
		{
			name: "public ipv6 subnet, local zone, must return error for unsupported ip version",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-public",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         true,
				IsIPv6:           true,
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "local-zone"`,
		},
		{
			name: "public ipv6 subnet, wavelength zone, must return error for unsupported ip version",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-public",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				IsPublic:         true,
				IsIPv6:           true,
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErr:        true,
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "wavelength-zone"`,
		},
		// private subnets
		{
			name: "private ipv4 subnet, availability zone, must have ipv4 default route to nat gateway",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-private",
				AvailabilityZone: "us-east-1a",
				IsPublic:         false,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					NatGatewayId:         aws.String("nat-gw-fromZone-us-east-1a"),
				},
			},
		},
		{
			name: "private ipv4 subnet, local zone, must have ipv4 default route to nat gateway",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-private",
				AvailabilityZone: "us-east-1-nyc-1a",
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
				IsPublic:         false,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					NatGatewayId:         aws.String("nat-gw-fromZone-us-east-1a"),
				},
			},
		},
		{
			name: "private ipv4 subnet, wavelength zone, must have ipv4 default route to nat gateway",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-private",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
				IsPublic:         false,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					NatGatewayId:         aws.String("nat-gw-fromZone-us-east-1a"),
				},
			},
		},
		// egress-only subnet ipv6
		{
			name: "egress-only ipv6 subnet, availability zone, must have ipv6 default route to egress-only gateway",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-private",
				AvailabilityZone: "us-east-1a",
				IsIPv6:           true,
				IsPublic:         false,
			},
			want: []*ec2.CreateRouteInput{
				{
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					NatGatewayId:         aws.String("nat-gw-fromZone-us-east-1a"),
				},
				{
					DestinationIpv6CidrBlock:    aws.String("::/0"),
					EgressOnlyInternetGatewayId: aws.String("vpc-eigw"),
				},
			},
		},
		{
			name: "private ipv6 subnet, availability zone, non-ipv6 block, must return error",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				net.VPC.IPv6 = nil
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-private",
				AvailabilityZone: "us-east-1a",
				IsIPv6:           true,
				IsPublic:         false,
			},
			wantErrMessage: `ipv6 block missing for ipv6 enabled subnet, can't create route for egress only internet gateway`,
		},
		// private subnet ipv6, unsupported
		{
			name: "private ipv6 subnet, local zone, must return unsupported",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-private",
				AvailabilityZone: "us-east-1-nyc-a",
				IsIPv6:           true,
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "local-zone"`,
		},
		{
			name: "private ipv6 subnet, wavelength zone, must return unsupported",
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-private",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
				IsIPv6:           true,
				IsPublic:         false,
			},
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "wavelength-zone"`,
		},
		// private subnet, gateway not found
		{
			name: "private ipv4 subnet, availability zone, must return error when invalid gateway",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				for i := range net.Subnets {
					if net.Subnets[i].AvailabilityZone == "us-east-1a" && net.Subnets[i].IsPublic {
						net.Subnets[i].NatGatewayID = nil
					}
				}
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-az-1a-private",
				AvailabilityZone: "us-east-1a",
				IsPublic:         false,
			},
			wantErrMessage: `no nat gateways available in "us-east-1a" for private subnet "subnet-az-1a-private"`,
		},
		{
			name: "private ipv4 subnet, local zone, must return error when invalid gateway",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := defaultNetwork.DeepCopy()
				for i := range net.Subnets {
					if net.Subnets[i].AvailabilityZone == "us-east-1a" && net.Subnets[i].IsPublic {
						net.Subnets[i].NatGatewayID = nil
					}
				}
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-lz-1a-private",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsIPv6:           true,
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneType("local-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "local-zone"`,
		},
		{
			name: "private ipv4 subnet, wavelength zone, must return error when invalid gateway",
			specOverrideNet: func() *infrav1.NetworkSpec {
				net := new(infrav1.NetworkSpec)
				*net = defaultNetwork
				net.VPC.CarrierGatewayID = nil
				return net
			}(),
			inputSubnet: &infrav1.SubnetSpec{
				ResourceID:       "subnet-wl-1a-private",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				IsIPv6:           true,
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneType("wavelength-zone")),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			wantErrMessage: `can't determine routes for unsupported ipv6 subnet in zone type "wavelength-zone"`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			cluster := scope.ClusterScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster-routes"},
				},
				AWSCluster: &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test"},
					Spec:       infrav1.AWSClusterSpec{},
				},
			}
			cluster.AWSCluster.Spec.NetworkSpec = defaultNetwork
			if tc.specOverrideNet != nil {
				cluster.AWSCluster.Spec.NetworkSpec = *tc.specOverrideNet
			}
			if tc.specOverrideSubnets != nil {
				cluster.AWSCluster.Spec.NetworkSpec.Subnets = *tc.specOverrideSubnets
			}

			scope, err := scope.NewClusterScope(cluster)
			if err != nil {
				t.Errorf("Service.getRoutesForSubnet() error setting up the test case: %v", err)
			}

			s := NewService(scope)
			got, err := s.getRoutesForSubnet(tc.inputSubnet)

			wantErr := tc.wantErr
			if len(tc.wantErrMessage) > 0 {
				wantErr = true
			}
			if wantErr && err == nil {
				t.Fatal("expected error but got no error")
			}
			if err != nil {
				if !wantErr {
					t.Fatalf("got an unexpected error: %v", err)
				}
				if wantErr && len(tc.wantErrMessage) > 0 && err.Error() != tc.wantErrMessage {
					t.Fatalf("got an unexpected error message:\nwant: %v\n got: %v\n", tc.wantErrMessage, err)
				}
			}
			if len(tc.want) > 0 {
				if !cmp.Equal(got, tc.want, cmp.AllowUnexported(ec2.CreateRouteInput{})) {
					t.Errorf("got unexpect routes:\n%v", cmp.Diff(got, tc.want))
				}
			}
		})
	}
}
