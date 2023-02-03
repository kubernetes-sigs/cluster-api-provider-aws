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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	subnetsVPCID = "vpc-subnets"
)

func TestReconcileSubnets(t *testing.T) {
	testCases := []struct {
		name          string
		input         ScopeBuilder
		expect        func(m *mocks.MockEC2APIMockRecorder)
		errorExpected bool
	}{
		{
			name: "Unmanaged VPC, 2 existing subnets in vpc, 2 subnet in spec, subnets match, with routes, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID: "subnet-1",
					},
					{
						ID: "subnet-2",
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []*ec2.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-1"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-2"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)
			},
		},
		{
			name: "IPv6 enabled vpc with default subnets should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID:            "subnet-1",
						IsIPv6:        true,
						IPv6CidrBlock: "2001:db8:1234:1a03::/64",
					},
					{
						ID:            "subnet-2",
						IsIPv6:        true,
						IPv6CidrBlock: "2001:db8:1234:1a02::/64",
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.10.0/24"),
								Ipv6CidrBlockAssociationSet: []*ec2.SubnetIpv6CidrBlockAssociation{
									{
										AssociationId: aws.String("amazon"),
										Ipv6CidrBlock: aws.String("2001:db8:1234:1a01::/64"),
										Ipv6CidrBlockState: &ec2.SubnetCidrBlockState{
											State: aws.String(ec2.SubnetCidrBlockStateCodeAssociated),
										},
									},
								},
								MapPublicIpOnLaunch:         aws.Bool(false),
								AssignIpv6AddressOnCreation: aws.Bool(true),
							},
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-2"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.20.0/24"),
								Ipv6CidrBlockAssociationSet: []*ec2.SubnetIpv6CidrBlockAssociation{
									{
										AssociationId: aws.String("amazon"),
										Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
										Ipv6CidrBlockState: &ec2.SubnetCidrBlockState{
											State: aws.String(ec2.SubnetCidrBlockStateCodeAssociated),
										},
									},
								},
								MapPublicIpOnLaunch:         aws.Bool(false),
								AssignIpv6AddressOnCreation: aws.Bool(true),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []*ec2.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-1"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-2"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)
			},
		},
		{
			name: "Unmanaged VPC, 2 existing subnets in vpc, 2 subnet in spec, subnets match, no routes, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID:   "subnet-1",
						Tags: map[string]string{"foo": "bar"}, // adding additional tag here which won't be added in unmanaged subnet hence not present in expect calls
					},
					{
						ID: "subnet-2",
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-1"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-2"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)
			},
			errorExpected: false,
		},
		{
			name: "Unmanaged VPC, 2 existing matching subnets, subnet tagging fails, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID: "subnet-1",
					},
					{
						ID: "subnet-2",
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []*ec2.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-1"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, fmt.Errorf("tagging failed"))
			},
		},
		{
			name: "Unmanaged VPC, 2 existing subnets in vpc, 0 subnet in spec, should fail",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)
			},
			errorExpected: true,
		},
		{
			name: "Unmanaged VPC, 0 existing subnets in vpc, 2 subnets in spec, should fail",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
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
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)
			},
			errorExpected: true,
		},
		{
			name: "Unmanaged VPC, 2 subnets exist, 2 private subnet in spec, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.10.0/24",
						IsPublic:         false,
					},
					{
						AvailabilityZone: "us-east-1b",
						CidrBlock:        "10.0.20.0/24",
						IsPublic:         false,
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-1"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)

				m.CreateTags(gomock.Eq(&ec2.CreateTagsInput{
					Resources: aws.StringSlice([]string{"subnet-2"}),
					Tags: []*ec2.Tag{
						{
							Key:   aws.String("kubernetes.io/cluster/test-cluster"),
							Value: aws.String("shared"),
						},
						{
							Key:   aws.String("kubernetes.io/role/internal-elb"),
							Value: aws.String("1"),
						},
					},
				})).
					Return(&ec2.CreateTagsOutput{}, nil)
			},
			errorExpected: false,
		},
		{
			name: "Managed VPC, no subnets exist, 1 private and 1 public subnet in spec, create both",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: []infrav1.SubnetSpec{
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
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				firstSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.1.0.0/16"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1a"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
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

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.2.0.0/16"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
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

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(secondSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(secondSubnet)
			},
		},
		{
			name: "Managed VPC, no subnets exist, 1 private subnet in spec (no public subnet), should fail",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: []infrav1.SubnetSpec{
					{
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.1.0.0/16",
						IsPublic:         false,
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)
			},
			errorExpected: true,
		},
		{
			name: "Managed VPC, no existing subnets exist, one az, expect one private and one public from default",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock: defaultVPCCidr,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
							},
						},
					}, nil)

				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				firstSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(firstSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(secondSubnet)
			},
		},
		{
			name: "Managed IPv6 VPC, no existing subnets exist, one az, expect one private and one public from default",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock: defaultVPCCidr,
					IPv6: &infrav1.IPv6{
						CidrBlock: "2001:db8:1234:1a01::/56",
						PoolID:    "amazon",
					},
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
							},
						},
					}, nil)

				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				firstSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a03::/64"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-1"),
							CidrBlock:                   aws.String("10.0.0.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []*ec2.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/64"),
									Ipv6CidrBlockState: &ec2.SubnetCidrBlockState{
										State: aws.String(ec2.SubnetCidrBlockStateCodeAssociated),
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(firstSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a02::/64"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-2"),
							CidrBlock:                   aws.String("10.0.128.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []*ec2.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
									Ipv6CidrBlockState: &ec2.SubnetCidrBlockState{
										State: aws.String(ec2.SubnetCidrBlockStateCodeAssociated),
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(secondSubnet)
			},
		},
		{
			name: "Managed VPC, no existing subnets exist, two az's, expect two private and two public from default",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock: defaultVPCCidr,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
							},
						},
					}, nil)

				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				zone1PublicSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/19"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/19"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.64.0/18"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.64.0/18"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PrivateSubnet)

				// zone 2

				zone2PublicSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.32.0/19"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.32.0/19"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PrivateSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				zone2PrivateSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/18"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/18"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone2PublicSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone2PrivateSubnet)
			},
		},
		{
			name: "Managed VPC, no existing subnets exist, two az's, max num azs is 1, expect one private and one public from default",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock:                  defaultVPCCidr,
					AvailabilityZoneUsageLimit: aws.Int(1),
					AvailabilityZoneSelection:  &infrav1.AZSelectionSchemeOrdered,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
							},
						},
					}, nil)

				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				zone1PublicSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/17"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/17"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PrivateSubnet)
			},
		},
		{
			name: "Managed VPC, existing public subnet, 2 subnets in spec, should create 1 subnet",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.0.0/17",
						IsPublic:         true,
					},
					{
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.128.0/17",
						IsPublic:         false,
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.0.0/17"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("public"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-subnet-public"),
									},
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("shared"),
									},
								},
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1a"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:            aws.String(subnetsVPCID),
							SubnetId:         aws.String("subnet-2"),
							CidrBlock:        aws.String("10.0.128.0/17"),
							AvailabilityZone: aws.String("us-east-1a"),
						},
					}, nil)

				m.WaitUntilSubnetAvailable(gomock.Any())

				// Public subnet
				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "Managed VPC, existing public subnet, 2 subnets in spec, should create 1 subnet, custom Name tag",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.0.0/17",
						IsPublic:         true,
					},
					{
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.128.0/17",
						IsPublic:         false,
						Tags:             map[string]string{"Name": "custom-sub"},
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.0.0/17"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
										Value: aws.String("owned"),
									},
									{
										Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
										Value: aws.String("public"),
									},
									{
										Key:   aws.String("Name"),
										Value: aws.String("test-cluster-subnet-public"),
									},
									{
										Key:   aws.String("kubernetes.io/cluster/test-cluster"),
										Value: aws.String("shared"),
									},
								},
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("custom-sub"), // must use the provided `Name` tag, not generate a name
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:            aws.String(subnetsVPCID),
							SubnetId:         aws.String("subnet-2"),
							CidrBlock:        aws.String("10.0.128.0/17"),
							AvailabilityZone: aws.String("us-east-1a"),
						},
					}, nil)

				m.WaitUntilSubnetAvailable(gomock.Any())

				// Public subnet
				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name: "With ManagedControlPlaneScope, Managed VPC, no existing subnets exist, two az's, expect two private and two public from default, created with tag including eksClusterName not a name of Cluster resource",
			input: NewManagedControlPlaneScope().
				WithEKSClusterName("test-eks-cluster").
				WithNetwork(&infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: subnetsVPCID,
						Tags: infrav1.Tags{
							infrav1.ClusterTagKey("test-cluster"): "owned",
						},
						CidrBlock: defaultVPCCidr,
					},
					Subnets: []infrav1.SubnetSpec{},
				}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
							},
						},
					}, nil)

				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				zone1PublicSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/19"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/19"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.64.0/18"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.64.0/18"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone1PrivateSubnet)

				// zone 2

				zone2PublicSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.32.0/19"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("public"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.32.0/19"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PrivateSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				zone2PrivateSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/18"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("subnet"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("shared"),
								},
								{
									Key:   aws.String("kubernetes.io/role/internal-elb"),
									Value: aws.String("1"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("private"),
								},
							},
						},
					},
				})).
					Return(&ec2.CreateSubnetOutput{
						Subnet: &ec2.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/18"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone2PublicSubnet)

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(zone2PrivateSubnet)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mocks.NewMockEC2API(mockCtrl)

			scope, err := tc.input.Build()
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock
			err = s.reconcileSubnets()

			if tc.errorExpected && err == nil {
				t.Fatal("expected error reconciling but not no error")
			}
			if !tc.errorExpected && err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}

func TestDiscoverSubnets(t *testing.T) {
	testCases := []struct {
		name   string
		input  *infrav1.NetworkSpec
		mocks  func(m *mocks.MockEC2APIMockRecorder)
		expect []infrav1.SubnetSpec
	}{
		{
			name: "provided VPC finds internet routes",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.10.0/24",
						IsPublic:         true,
					},
					{
						ID:               "subnet-2",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.11.0/24",
						IsPublic:         false,
					},
				},
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.10.0/24"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("Name"),
										Value: aws.String("provided-subnet-public"),
									},
								},
							},
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-2"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.11.0/24"),
								Tags: []*ec2.Tag{
									{
										Key:   aws.String("Name"),
										Value: aws.String("provided-subnet-private"),
									},
								},
							},
						},
					}, nil)

				m.DescribeRouteTables(gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []*ec2.RouteTable{
							{
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-1"),
									},
								},
								Routes: []*ec2.Route{
									{
										DestinationCidrBlock: aws.String("10.0.10.0/24"),
										GatewayId:            aws.String("local"),
									},
									{
										DestinationCidrBlock: aws.String("0.0.0.0/0"),
										GatewayId:            aws.String("igw-0"),
									},
								},
								RouteTableId: aws.String("rtb-1"),
							},
							{
								Associations: []*ec2.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-2"),
									},
								},
								Routes: []*ec2.Route{
									{
										DestinationCidrBlock: aws.String("10.0.11.0/24"),
										GatewayId:            aws.String("local"),
									},
								},
								RouteTableId: aws.String("rtb-2"),
							},
						},
					}, nil)

				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String(subnetsVPCID)},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(&ec2.CreateTagsOutput{}, nil).AnyTimes()
			},
			expect: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
					RouteTableID:     aws.String("rtb-1"),
					Tags: infrav1.Tags{
						"Name": "provided-subnet-public",
					},
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.11.0/24",
					IsPublic:         false,
					RouteTableID:     aws.String("rtb-2"),
					Tags: infrav1.Tags{
						"Name": "provided-subnet-private",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
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

			tc.mocks(ec2Mock.EXPECT())

			s := NewService(scope)
			s.EC2Client = ec2Mock

			if err := s.reconcileSubnets(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			subnets := s.scope.Subnets()
			out := make(map[string]infrav1.SubnetSpec)
			for _, sn := range subnets {
				out[sn.ID] = sn
			}
			for _, exp := range tc.expect {
				sn, ok := out[exp.ID]
				if !ok {
					t.Errorf("Expected to find subnet %s in %+v", exp.ID, subnets)
					continue
				}

				if !cmp.Equal(sn, exp) {
					expected, err := json.MarshalIndent(exp, "", "\t")
					if err != nil {
						t.Fatalf("got an unexpected error: %v", err)
					}
					actual, err := json.MarshalIndent(sn, "", "\t")
					if err != nil {
						t.Fatalf("got an unexpected error: %v", err)
					}
					t.Errorf("Expected %s, got %s", string(expected), string(actual))
				}
				delete(out, exp.ID)
			}
			if len(out) > 0 {
				t.Errorf("Got unexpected subnets: %+v", out)
			}
		})
	}
}

func TestDeleteSubnets(t *testing.T) {
	testCases := []struct {
		name          string
		input         *infrav1.NetworkSpec
		expect        func(m *mocks.MockEC2APIMockRecorder)
		errorExpected bool
	}{
		{
			name: "managed vpc - success",
			input: &infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID: "subnet-1",
					},
					{
						ID: "subnet-2",
					},
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []*ec2.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(true),
							},
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-2"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DeleteSubnet(&ec2.DeleteSubnetInput{
					SubnetId: aws.String("subnet-1"),
				}).
					Return(nil, nil)

				m.DeleteSubnet(&ec2.DeleteSubnetInput{
					SubnetId: aws.String("subnet-2"),
				}).
					Return(nil, nil)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

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

			err = s.deleteSubnets()
			if tc.errorExpected && err == nil {
				t.Fatal("expected error but not no error")
			}
			if !tc.errorExpected && err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}

// Test helpers

type ScopeBuilder interface {
	Build() (scope.NetworkScope, error)
}

func NewClusterScope() *ClusterScopeBuilder {
	return &ClusterScopeBuilder{
		customizers: []func(p *scope.ClusterScopeParams){},
	}
}

type ClusterScopeBuilder struct {
	customizers []func(p *scope.ClusterScopeParams)
}

func (b *ClusterScopeBuilder) WithNetwork(n *infrav1.NetworkSpec) *ClusterScopeBuilder {
	b.customizers = append(b.customizers, func(p *scope.ClusterScopeParams) {
		p.AWSCluster.Spec.NetworkSpec = *n
	})

	return b
}

func (b *ClusterScopeBuilder) Build() (scope.NetworkScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	param := &scope.ClusterScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		AWSCluster: &infrav1.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			Spec:       infrav1.AWSClusterSpec{},
		},
	}

	for _, customizer := range b.customizers {
		customizer(param)
	}

	return scope.NewClusterScope(*param)
}

func NewManagedControlPlaneScope() *ManagedControlPlaneScopeBuilder {
	return &ManagedControlPlaneScopeBuilder{
		customizers: []func(p *scope.ManagedControlPlaneScopeParams){},
	}
}

type ManagedControlPlaneScopeBuilder struct {
	customizers []func(p *scope.ManagedControlPlaneScopeParams)
}

func (b *ManagedControlPlaneScopeBuilder) WithNetwork(n *infrav1.NetworkSpec) *ManagedControlPlaneScopeBuilder {
	b.customizers = append(b.customizers, func(p *scope.ManagedControlPlaneScopeParams) {
		p.ControlPlane.Spec.NetworkSpec = *n
	})

	return b
}

func (b *ManagedControlPlaneScopeBuilder) WithEKSClusterName(name string) *ManagedControlPlaneScopeBuilder {
	b.customizers = append(b.customizers, func(p *scope.ManagedControlPlaneScopeParams) {
		p.ControlPlane.Spec.EKSClusterName = name
	})

	return b
}

func (b *ManagedControlPlaneScopeBuilder) Build() (scope.NetworkScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	param := &scope.ManagedControlPlaneScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			Spec:       ekscontrolplanev1.AWSManagedControlPlaneSpec{},
		},
	}

	for _, customizer := range b.customizers {
		customizer(param)
	}

	return scope.NewManagedControlPlaneScope(*param)
}
