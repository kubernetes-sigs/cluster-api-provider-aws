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
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/test/mocks"
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
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				privateRouteTable := m.CreateRouteTable(matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-1")}}, nil)

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

				publicRouteTable := m.CreateRouteTable(matchRouteTableInput(&ec2.CreateRouteTableInput{VpcId: aws.String("vpc-routetables")})).
					Return(&ec2.CreateRouteTableOutput{RouteTable: &ec2.RouteTable{RouteTableId: aws.String("rt-2")}}, nil)

				m.CreateRoute(gomock.Eq(&ec2.CreateRouteInput{
					GatewayId:            aws.String("igw-01"),
					DestinationCidrBlock: aws.String("0.0.0.0/0"),
					RouteTableId:         aws.String("rt-2"),
				})).
					After(publicRouteTable)

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
										NatGatewayId:         aws.String("nat-01"),
									},
									// Extra (managed outside of CAPA) route with Managed Prefix List destination.
									{
										DestinationPrefixListId: aws.String("pl-foobar"),
									},
								},
								Tags: []*ec2.Tag{
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

func TestDeleteRouteTables(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	describeRouteTableOutput := &ec2.DescribeRouteTablesOutput{
		RouteTables: []*ec2.RouteTable{
			{
				RouteTableId: aws.String("route-table-private"),
				Associations: []*ec2.RouteTableAssociation{
					{
						SubnetId: nil,
					},
				},
				Routes: []*ec2.Route{
					{
						DestinationCidrBlock: aws.String("0.0.0.0/0"),
						NatGatewayId:         aws.String("outdated-nat-01"),
					},
				},
			},
			{
				RouteTableId: aws.String("route-table-public"),
				Associations: []*ec2.RouteTableAssociation{
					{
						SubnetId:                aws.String("subnet-routetables-public"),
						RouteTableAssociationId: aws.String("route-table-public"),
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
						Value: aws.String("test-cluster-rt-public-us-east-1a"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
				},
			},
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
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)

				m.DisassociateRouteTable(gomock.Eq(&ec2.DisassociateRouteTableInput{
					AssociationId: aws.String("route-table-public"),
				})).Return(&ec2.DisassociateRouteTableOutput{}, nil)

				m.DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-public"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)
			},
		},
		{
			name:  "Should return error if describe route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
			wantErr: true,
		},
		{
			name:  "Should return error if delete route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(nil, awserrors.NewNotFound("not found"))
			},
			wantErr: true,
		},
		{
			name:  "Should return error if disassociate route table fails",
			input: &infrav1.NetworkSpec{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(describeRouteTableOutput, nil)

				m.DeleteRouteTable(gomock.Eq(&ec2.DeleteRouteTableInput{
					RouteTableId: aws.String("route-table-private"),
				})).Return(&ec2.DeleteRouteTableOutput{}, nil)

				m.DisassociateRouteTable(gomock.Eq(&ec2.DisassociateRouteTableInput{
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
