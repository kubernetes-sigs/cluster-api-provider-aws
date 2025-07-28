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
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	subnetsVPCID = "vpc-subnets"
)

func TestReconcileSubnets(t *testing.T) {
	// SubnetSpecs for different zone types.
	stubSubnetsAvailabilityZone := []infrav1.SubnetSpec{
		{ID: "subnet-private-us-east-1a", AvailabilityZone: "us-east-1a", CidrBlock: "10.0.1.0/24", IsPublic: false},
		{ID: "subnet-public-us-east-1a", AvailabilityZone: "us-east-1a", CidrBlock: "10.0.2.0/24", IsPublic: true},
	}
	stubAdditionalSubnetsAvailabilityZone := []infrav1.SubnetSpec{
		{ID: "subnet-private-us-east-1b", AvailabilityZone: "us-east-1b", CidrBlock: "10.0.3.0/24", IsPublic: false},
		{ID: "subnet-public-us-east-1b", AvailabilityZone: "us-east-1b", CidrBlock: "10.0.4.0/24", IsPublic: true},
	}
	stubSubnetsLocalZone := []infrav1.SubnetSpec{
		{ID: "subnet-private-us-east-1-nyc-1a", AvailabilityZone: "us-east-1-nyc-1a", CidrBlock: "10.0.5.0/24", IsPublic: false},
		{ID: "subnet-public-us-east-1-nyc-1a", AvailabilityZone: "us-east-1-nyc-1a", CidrBlock: "10.0.6.0/24", IsPublic: true},
	}
	stubSubnetsWavelengthZone := []infrav1.SubnetSpec{
		{ID: "subnet-private-us-east-1-wl1-nyc-wlz-1", AvailabilityZone: "us-east-1-wl1-nyc-wlz-1", CidrBlock: "10.0.7.0/24", IsPublic: false},
		{ID: "subnet-public-us-east-1-wl1-nyc-wlz-1", AvailabilityZone: "us-east-1-wl1-nyc-wlz-1", CidrBlock: "10.0.8.0/24", IsPublic: true},
	}
	stubSubnetsAllZones := slices.Concat(stubSubnetsAvailabilityZone, stubSubnetsLocalZone, stubSubnetsWavelengthZone)

	// NetworkSpec with subnets in zone type availability-zone
	stubNetworkSpecWithSubnets := &infrav1.NetworkSpec{
		VPC: infrav1.VPCSpec{
			ID: subnetsVPCID,
			Tags: infrav1.Tags{
				infrav1.ClusterTagKey("test-cluster"): "owned",
			},
		},
		Subnets: stubSubnetsAvailabilityZone,
	}
	// NetworkSpec with subnets in zone types availability-zone, local-zone and wavelength-zone
	stubNetworkSpecWithSubnetsEdge := stubNetworkSpecWithSubnets.DeepCopy()
	stubNetworkSpecWithSubnetsEdge.Subnets = stubSubnetsAllZones

	testCases := []struct {
		name                         string
		input                        ScopeBuilder
		expect                       func(m *mocks.MockEC2APIMockRecorder)
		errorExpected                bool
		errorMessageExpected         string
		tagUnmanagedNetworkResources bool
		optionalExpectSubnets        infrav1.Subnets
	}{
		{
			name: "Unmanaged VPC, disable TagUnmanagedNetworkResources, 2 existing subnets in vpc, 2 subnet in spec, subnets match, with routes, should succeed",
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
			}).WithTagUnmanagedNetworkResources(false),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
			tagUnmanagedNetworkResources: false,
		},
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
			tagUnmanagedNetworkResources: true,
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.10.0/24"),
								Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
									{
										AssociationId: aws.String("amazon"),
										Ipv6CidrBlock: aws.String("2001:db8:1234:1a01::/64"),
										Ipv6CidrBlockState: &types.SubnetCidrBlockState{
											State: types.SubnetCidrBlockStateCodeAssociated,
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
								Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
									{
										AssociationId: aws.String("amazon"),
										Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
										Ipv6CidrBlockState: &types.SubnetCidrBlockState{
											State: types.SubnetCidrBlockStateCodeAssociated,
										},
									},
								},
								MapPublicIpOnLaunch:         aws.Bool(false),
								AssignIpv6AddressOnCreation: aws.Bool(true),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
			tagUnmanagedNetworkResources: true,
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
			errorExpected:                false,
			tagUnmanagedNetworkResources: true,
		},
		{
			name: "Unmanaged VPC, one existing matching subnets, subnet tagging fails, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID: "subnet-1",
					},
				},
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a"), ZoneType: aws.String("availability-zone")},
				}).AnyTimes()

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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
			},
			tagUnmanagedNetworkResources: true,
		},
		{
			name: "Unmanaged VPC, one existing matching subnets, subnet tagging fails with subnet update, should succeed",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{
					{
						ID: "subnet-1",
					},
				},
			}).WithTagUnmanagedNetworkResources(true),
			optionalExpectSubnets: infrav1.Subnets{
				{
					ID:               "subnet-1",
					ResourceID:       "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:               aws.String(subnetsVPCID),
								SubnetId:            aws.String("subnet-1"),
								AvailabilityZone:    aws.String("us-east-1a"),
								CidrBlock:           aws.String("10.0.10.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a")},
				}).AnyTimes()

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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
			tagUnmanagedNetworkResources: true,
		},
		{
			name: "Unmanaged VPC, 2 existing matching subnets, subnet tagging fails with subnet update, should succeed",
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
			}).WithTagUnmanagedNetworkResources(true),
			optionalExpectSubnets: infrav1.Subnets{
				{
					ID:               "subnet-1",
					ResourceID:       "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					ResourceID:       "subnet-2",
					AvailabilityZone: "us-east-1b",
					CidrBlock:        "10.0.11.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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
								AvailabilityZone:    aws.String("us-east-1b"),
								CidrBlock:           aws.String("10.0.11.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-2"),
										RouteTableId: aws.String("rt-00000"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a")}, {ZoneName: aws.String("us-east-1b")},
				}).AnyTimes()

				subnet1tag := m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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
					Return(&ec2.CreateTagsOutput{}, fmt.Errorf("tagging failed")).After(subnet1tag)
			},
			tagUnmanagedNetworkResources: true,
		},
		{
			name: "Unmanaged VPC, 2 existing matching subnets, subnet tagging fails second call, should succeed",
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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
								AvailabilityZone:    aws.String("us-east-1b"),
								CidrBlock:           aws.String("10.0.20.0/24"),
								MapPublicIpOnLaunch: aws.Bool(false),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-1"),
										RouteTableId: aws.String("rt-12345"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
							{
								VpcId: aws.String(subnetsVPCID),
								Associations: []types.RouteTableAssociation{
									{
										SubnetId:     aws.String("subnet-2"),
										RouteTableId: aws.String("rt-22222"),
									},
								},
								Routes: []types.Route{
									{
										GatewayId: aws.String("igw-12345"),
									},
								},
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a")}, {ZoneName: aws.String("us-east-1b")},
				}).AnyTimes()

				secondSubnetTag := m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a"), ZoneType: aws.String("availability-zone")},
					{ZoneName: aws.String("us-east-1b"), ZoneType: aws.String("availability-zone")},
				}).AnyTimes()

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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
					Return(&ec2.CreateTagsOutput{}, fmt.Errorf("tagging failed")).After(secondSubnetTag)
			},
			tagUnmanagedNetworkResources: true,
		},
		{
			name: "Unmanaged VPC, 2 existing subnets in vpc, 0 subnet in spec, should fail",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
				},
				Subnets: []infrav1.SubnetSpec{},
			}).WithTagUnmanagedNetworkResources(true),
			expect:                       func(m *mocks.MockEC2APIMockRecorder) {},
			errorExpected:                true,
			tagUnmanagedNetworkResources: true,
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)
			},
			errorExpected:                true,
			tagUnmanagedNetworkResources: true,
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
			}).WithTagUnmanagedNetworkResources(true),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags: []types.Tag{
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

				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags: []types.Tag{
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

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
			errorExpected:                false,
			tagUnmanagedNetworkResources: true,
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
				describeCall := m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				}), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)
				// Create the first subnet
				firstSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.1.0.0/16"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1a"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.1.0.0/16"),
							AvailabilityZone:    aws.String("us-east-1a"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				// Wait until first subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).After(firstSubnet)

				// Create the second subnet
				secondSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.2.0.0/16"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.2.0.0/16"),
							AvailabilityZone:    aws.String("us-east-1a"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				// Wait until second subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).After(secondSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(secondSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
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
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				firstSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				// Wait until first subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				// Wait until second subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).After(secondSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
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
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				firstSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a03::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-1"),
							CidrBlock:                   aws.String("10.0.0.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				// Wait until first subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a02::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-2"),
							CidrBlock:                   aws.String("10.0.128.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				// Wait until second subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(secondSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
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
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				// Zone1
				m.DescribeAvailabilityZones(context.TODO(), gomock.Eq(&ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1b"},
				})).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).MaxTimes(2)

				zone1PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/19"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/19"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				// Wait until first subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.64.0/18"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.64.0/18"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				// Wait until second subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PrivateSubnet)

				// zone 2
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				zone2PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.32.0/19"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.32.0/19"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PrivateSubnet)

				// Wait until first subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				zone2PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/18"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/18"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone2PublicSubnet)

				// Wait until second subnet is available
				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
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
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				zone1PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/17"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/17"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).After(zone1PrivateSubnet)
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.0.0/17"),
								Tags: []types.Tag{
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
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1a"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:            aws.String(subnetsVPCID),
							SubnetId:         aws.String("subnet-2"),
							CidrBlock:        aws.String("10.0.128.0/17"),
							AvailabilityZone: aws.String("us-east-1a"),
						},
					}, nil)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil)

				// Public subnet
				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.0.0/17"),
								Tags: []types.Tag{
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
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1a"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("custom-sub"), // must use the provided `Name` tag, not generate a name
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:            aws.String(subnetsVPCID),
							SubnetId:         aws.String("subnet-2"),
							CidrBlock:        aws.String("10.0.128.0/17"),
							AvailabilityZone: aws.String("us-east-1a"),
						},
					}, nil)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil)

				// Public subnet
				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()
			},
		},
		{
			name: "Managed VPC, existing public and private subnets, 2 subnets in spec, custom tags in spec should be created",
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
						IsPublic:         true,
						Tags:             map[string]string{"this-tag-is-in-the-spec": "but-its-not-on-aws"},
					},
					{
						ID:               "subnet-2",
						AvailabilityZone: "us-east-1a",
						IsPublic:         false,
						Tags:             map[string]string{"subnet-2-this-tag-is-in-the-spec": "subnet-2-but-its-not-on-aws"},
					},
				},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				tagsOnSubnet1 := []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-subnet-public"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("public"),
					},
				}
				tagsOnSubnet2 := []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-subnet-private"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
					},
					{
						Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
						Value: aws.String("private"),
					},
				}
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.0.0/17"),
								Tags:             tagsOnSubnet1,
							},
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-2"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.128.0/17"),
								Tags:             tagsOnSubnet2,
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				// Public subnet
				expectedAppliedAwsTagsForSubnet1 := []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-subnet-public-us-east-1a"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
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
					{
						Key:   aws.String("this-tag-is-in-the-spec"),
						Value: aws.String("but-its-not-on-aws"),
					}}
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-1"},
					Tags:      expectedAppliedAwsTagsForSubnet1,
				})).
					Return(nil, nil)

				// Private subnet
				expectedAppliedAwsTagsForSubnet2 := []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-cluster-subnet-private-us-east-1a"),
					},
					{
						Key:   aws.String("kubernetes.io/cluster/test-cluster"),
						Value: aws.String("owned"),
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
					{
						Key:   aws.String("subnet-2-this-tag-is-in-the-spec"),
						Value: aws.String("subnet-2-but-its-not-on-aws"),
					}}
				m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
					Resources: []string{"subnet-2"},
					Tags:      expectedAppliedAwsTagsForSubnet2,
				})).
					Return(nil, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()
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
				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				// Zone 1 subnet.
				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				zone1PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/19"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.0.0/19"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.64.0/18"),
					AvailabilityZone: aws.String("us-east-1b"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.64.0/18"),
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PrivateSubnet)

				// zone 2
				zone2PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.32.0/19"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.32.0/19"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PrivateSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Eq(&ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				})).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				zone2PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/18"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-eks-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.128.0/18"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone2PublicSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone2PrivateSubnet)
			},
		},
		{ // Edge Zones
			name: "Managed VPC, local zones, no existing subnets exist, two az's, one LZ, expect two private and two public from default, and one private and public from Local Zone",
			input: func() *ClusterScopeBuilder {
				stubNetworkSpecEdgeLocalZonesOnly := stubNetworkSpecWithSubnets.DeepCopy()
				stubNetworkSpecEdgeLocalZonesOnly.Subnets = stubSubnetsAvailabilityZone
				stubNetworkSpecEdgeLocalZonesOnly.Subnets = append(stubNetworkSpecEdgeLocalZonesOnly.Subnets, stubAdditionalSubnetsAvailabilityZone...)
				stubNetworkSpecEdgeLocalZonesOnly.Subnets = append(stubNetworkSpecEdgeLocalZonesOnly.Subnets, stubSubnetsLocalZone...)
				return NewClusterScope().WithNetwork(stubNetworkSpecEdgeLocalZonesOnly)
			}(),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := stubMockDescribeSubnetsWithContextManaged(m)
				stubMockDescribeRouteTables(m)
				stubMockDescribeNatGateways(m)
				stubMockDescribeAvailabilityZonesWithContextCustomZones(m, []types.AvailabilityZone{
					{ZoneName: aws.String("us-east-1a"), ZoneType: aws.String("availability-zone")},
					{ZoneName: aws.String("us-east-1b"), ZoneType: aws.String("availability-zone")},
					{ZoneName: aws.String("us-east-1-nyc-1a"), ZoneType: aws.String("local-zone"), ParentZoneName: aws.String("us-east-1a")},
					{ZoneName: aws.String("us-east-1-wl1-nyc-wlz-1"), ZoneType: aws.String("wavelength-zone"), ParentZoneName: aws.String("us-east-1a")},
				}).AnyTimes()

				// Zone 1a subnets
				az1aPrivate := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "private", "10.0.1.0/24", false).
					After(describeCall)

				az1aPublic := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "public", "10.0.2.0/24", false).
					After(az1aPrivate)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1a").
					After(az1aPublic)

				// Zone 1b subnets
				az1bPrivate := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1b", "private", "10.0.3.0/24", false).
					After(az1aPublic)

				az1bPublic := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1b", "public", "10.0.4.0/24", false).
					After(az1bPrivate)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1b").
					After(az1bPublic)

				// Local zone 1-nyc-1a.
				lz1Private := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-nyc-1a", "private", "10.0.5.0/24", true).
					After(az1bPublic)

				lz1Public := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-nyc-1a", "public", "10.0.6.0/24", true).After(lz1Private)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1-nyc-1a").
					After(lz1Public)
			},
		},
		{
			name:  "Managed VPC, edge zones, custom names, no existing subnets exist, one AZ, LZ and WL, expect one private and one public subnets from each of default zones, Local Zone, and Wavelength",
			input: NewClusterScope().WithNetwork(stubNetworkSpecWithSubnetsEdge),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := stubMockDescribeSubnetsWithContextManaged(m)
				stubMockDescribeRouteTables(m)
				stubMockDescribeNatGateways(m)
				stubMockDescribeAvailabilityZonesWithContextAllZones(m)

				// AZone 1a subnets
				az1Private := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "private", "10.0.1.0/24", false).
					After(describeCall)

				az1Public := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "public", "10.0.2.0/24", false).After(az1Private)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1a").After(az1Public)

				// Local zone 1-nyc-1a.
				lz1Private := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-nyc-1a", "private", "10.0.5.0/24", true).
					After(describeCall)

				lz1Public := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-nyc-1a", "public", "10.0.6.0/24", true).After(lz1Private)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1-nyc-1a").After(lz1Public)

				// Wavelength zone nyc-1.
				wz1Private := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-wl1-nyc-wlz-1", "private", "10.0.7.0/24", true).
					After(describeCall)

				stubGenMockCreateSubnet(m, "test-cluster", "us-east-1-wl1-nyc-wlz-1", "public", "10.0.8.0/24", true).After(wz1Private)
			},
		},
		{
			name:  "Managed VPC, edge zones, error when retrieving zone information for subnet's AvailabilityZone",
			input: NewClusterScope().WithNetwork(stubNetworkSpecWithSubnetsEdge),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				stubMockDescribeSubnetsWithContextManaged(m)
				stubMockDescribeRouteTables(m)
				stubMockDescribeNatGateways(m)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{},
					}, nil)
			},
			errorExpected:        true,
			errorMessageExpected: `expected the zone attributes to be populated to subnet: unable to update zone information for subnet 'subnet-private-us-east-1a' and zone 'us-east-1a'`,
		},
		{
			name: "Managed VPC, edge zones, error when IPv6 subnet",
			input: func() *ClusterScopeBuilder {
				net := stubNetworkSpecWithSubnetsEdge.DeepCopy()
				// Only AZ and LZ to simplify the goal
				net.Subnets = infrav1.Subnets{}
				for i := range stubSubnetsAvailabilityZone {
					net.Subnets = append(net.Subnets, *stubSubnetsAvailabilityZone[i].DeepCopy())
				}
				for i := range stubSubnetsLocalZone {
					lz := stubSubnetsLocalZone[i].DeepCopy()
					lz.IsIPv6 = true
					net.Subnets = append(net.Subnets, *lz)
				}
				return NewClusterScope().WithNetwork(net)
			}(),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describe := stubMockDescribeSubnetsWithContextManaged(m)
				stubMockDescribeRouteTables(m)
				stubMockDescribeNatGateways(m)
				stubMockDescribeAvailabilityZonesWithContextAllZones(m)

				az1Private := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "private", "10.0.1.0/24", false).After(describe)

				az1Public := stubGenMockCreateSubnet(m, "test-cluster", "us-east-1a", "public", "10.0.2.0/24", false).After(az1Private)
				stubMockModifySubnetAttributeWithContext(m, "subnet-public-us-east-1a").After(az1Public)
			},
			errorExpected:        true,
			errorMessageExpected: `failed to create subnet: IPv6 is not supported with zone type "local-zone"`,
		},
		{
			name: "Unmanaged VPC, edge zones, existing subnets, one AZ, LZ and WL, expect one private and one public subnets from each of default zones, Local Zone, and Wavelength",
			input: func() *ClusterScopeBuilder {
				net := stubNetworkSpecWithSubnetsEdge.DeepCopy()
				net.VPC = infrav1.VPCSpec{
					ID: subnetsVPCID,
				}
				net.Subnets = infrav1.Subnets{
					{ResourceID: "subnet-az-1a-private"},
					{ResourceID: "subnet-az-1a-public"},
					{ResourceID: "subnet-lz-1a-private"},
					{ResourceID: "subnet-lz-1a-public"},
					{ResourceID: "subnet-wl-1a-private"},
					{ResourceID: "subnet-wl-1a-public"},
				}
				return NewClusterScope().WithNetwork(net)
			}(),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				stubMockDescribeSubnetsWithContextUnmanaged(m)
				stubMockDescribeAvailabilityZonesWithContextAllZones(m)
				stubMockDescribeRouteTablesWithWavelength(m,
					[]string{"subnet-az-1a-private", "subnet-lz-1a-private", "subnet-wl-1a-private"},
					[]string{"subnet-az-1a-public", "subnet-lz-1a-public"},
					[]string{"subnet-wl-1a-public"})

				stubMockDescribeNatGateways(m)
				stubMockCreateTags(m, "test-cluster", "subnet-az-1a-private", "us-east-1a", "private", false).AnyTimes()
			},
		},
		{
			name: "Managed VPC, no existing subnets exist, one az, prefer public subnet schema, expect one private and one public from default",
			input: NewClusterScope().WithNetwork(&infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID: subnetsVPCID,
					Tags: infrav1.Tags{
						infrav1.ClusterTagKey("test-cluster"): "owned",
					},
					CidrBlock:    defaultVPCCidr,
					SubnetSchema: &infrav1.SubnetSchemaPreferPublic,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				firstSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-1"),
							CidrBlock:           aws.String("10.0.128.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:               aws.String(subnetsVPCID),
							SubnetId:            aws.String("subnet-2"),
							CidrBlock:           aws.String("10.0.0.0/17"),
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(secondSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
		},
		{
			name: "Managed IPv6 VPC, no existing subnets exist, one az, prefer public subnet schema, expect one private and one public from default",
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
					SubnetSchema: &infrav1.SubnetSchemaPreferPublic,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				firstSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a02::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-1"),
							CidrBlock:                   aws.String("10.0.128.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(firstSubnet)

				secondSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/17"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a03::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-2"),
							CidrBlock:                   aws.String("10.0.0.0/17"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(firstSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(secondSubnet)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)
			},
		},
		{
			name: "Managed IPv6 VPC, no existing subnets exist, two az's, prefer public subnet schema, expect two private and two public from default",
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
					SubnetSchema: &infrav1.SubnetSchemaPreferPublic,
				},
				Subnets: []infrav1.SubnetSpec{},
			}),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeCall := m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				// Zone1
				m.DescribeAvailabilityZones(context.TODO(), gomock.Eq(&ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1b"},
				})).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1b"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).MaxTimes(2)

				zone1PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.64.0/18"),
					AvailabilityZone: aws.String("us-east-1b"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a02::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-1"),
							CidrBlock:                   aws.String("10.0.64.0/18"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a02::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(describeCall)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone1PublicSubnet)

				zone1PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.0.0/19"),
					AvailabilityZone: aws.String("us-east-1b"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a04::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1b"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-2"),
							CidrBlock:                   aws.String("10.0.0.0/19"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a04::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1b"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PublicSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone1PrivateSubnet)

				// zone 2
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1c"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1c"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil).AnyTimes()

				zone2PublicSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.128.0/18"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a03::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-public-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-1"),
							CidrBlock:                   aws.String("10.0.128.0/18"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone1PrivateSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-1"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-1"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					AssignIpv6AddressOnCreation: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-2"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)
				m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &types.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: aws.String("subnet-1"),
				}).
					Return(&ec2.ModifySubnetAttributeOutput{}, nil).
					After(zone2PublicSubnet)

				zone2PrivateSubnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
					VpcId:            aws.String(subnetsVPCID),
					CidrBlock:        aws.String("10.0.32.0/19"),
					AvailabilityZone: aws.String("us-east-1c"),
					Ipv6CidrBlock:    aws.String("2001:db8:1234:1a05::/64"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeSubnet,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-subnet-private-us-east-1c"),
								},
								{
									Key:   aws.String("kubernetes.io/cluster/test-cluster"),
									Value: aws.String("owned"),
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
						Subnet: &types.Subnet{
							VpcId:                       aws.String(subnetsVPCID),
							SubnetId:                    aws.String("subnet-2"),
							CidrBlock:                   aws.String("10.0.32.0/19"),
							AssignIpv6AddressOnCreation: aws.Bool(true),
							Ipv6CidrBlockAssociationSet: []types.SubnetIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a05::/64"),
									Ipv6CidrBlockState: &types.SubnetCidrBlockState{
										State: types.SubnetCidrBlockStateCodeAssociated,
									},
								},
							},
							AvailabilityZone:    aws.String("us-east-1c"),
							MapPublicIpOnLaunch: aws.Bool(false),
						},
					}, nil).
					After(zone2PublicSubnet)

				m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{
					SubnetIds: []string{"subnet-2"},
				}), gomock.Any()).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							VpcId:    aws.String(subnetsVPCID),
							SubnetId: aws.String("subnet-2"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
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
			if tc.errorExpected && err != nil && len(tc.errorMessageExpected) > 0 {
				if err.Error() != tc.errorMessageExpected {
					t.Fatalf("got an unexpected error message:\nwant: %v\n got: %v\n", tc.errorMessageExpected, err.Error())
				}
			}
			if !tc.errorExpected && err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
			if tc.errorExpected && err != nil && len(tc.errorMessageExpected) > 0 {
				if err.Error() != tc.errorMessageExpected {
					t.Fatalf("got an unexpected error message: %v", err)
				}
			}
			if len(tc.optionalExpectSubnets) > 0 {
				if !cmp.Equal(s.scope.Subnets(), tc.optionalExpectSubnets) {
					t.Errorf("got unexpect Subnets():\n%v", cmp.Diff(s.scope.Subnets(), tc.optionalExpectSubnets))
				}
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
						ZoneType:         ptr.To[infrav1.ZoneType]("availability-zone"),
					},
					{
						ID:               "subnet-2",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.11.0/24",
						IsPublic:         false,
						ZoneType:         ptr.To[infrav1.ZoneType]("availability-zone"),
					},
				},
			},
			mocks: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								VpcId:            aws.String(subnetsVPCID),
								SubnetId:         aws.String("subnet-1"),
								AvailabilityZone: aws.String("us-east-1a"),
								CidrBlock:        aws.String("10.0.10.0/24"),
								Tags: []types.Tag{
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
								Tags: []types.Tag{
									{
										Key:   aws.String("Name"),
										Value: aws.String("provided-subnet-private"),
									},
								},
							},
						},
					}, nil)

				m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName: aws.String("us-east-1a"),
								ZoneType: aws.String("availability-zone"),
							},
						},
					}, nil)

				m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
					Return(&ec2.DescribeRouteTablesOutput{
						RouteTables: []types.RouteTable{
							{
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-1"),
									},
								},
								Routes: []types.Route{
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
								Associations: []types.RouteTableAssociation{
									{
										SubnetId: aws.String("subnet-2"),
									},
								},
								Routes: []types.Route{
									{
										DestinationCidrBlock: aws.String("10.0.11.0/24"),
										GatewayId:            aws.String("local"),
									},
								},
								RouteTableId: aws.String("rtb-2"),
							},
						},
					}, nil)

				m.DescribeNatGateways(context.TODO(),
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []types.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []string{subnetsVPCID},
							},
							{
								Name:   aws.String("state"),
								Values: []string{"pending", "available"},
							},
						},
					}),
					gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)

				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(&ec2.CreateTagsOutput{}, nil).AnyTimes()
			},
			expect: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					ResourceID:       "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
					RouteTableID:     aws.String("rtb-1"),
					ZoneType:         ptr.To[infrav1.ZoneType]("availability-zone"),
				},
				{
					ID:               "subnet-2",
					ResourceID:       "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.11.0/24",
					IsPublic:         false,
					RouteTableID:     aws.String("rtb-2"),
					ZoneType:         ptr.To[infrav1.ZoneType]("availability-zone"),
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
				m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
						{
							Name:   aws.String("vpc-id"),
							Values: []string{subnetsVPCID},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
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

				m.DeleteSubnet(context.TODO(), &ec2.DeleteSubnetInput{
					SubnetId: aws.String("subnet-1"),
				}).
					Return(nil, nil)

				m.DeleteSubnet(context.TODO(), &ec2.DeleteSubnetInput{
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

// Test helpers.

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

func (b *ClusterScopeBuilder) WithTagUnmanagedNetworkResources(value bool) *ClusterScopeBuilder {
	b.customizers = append(b.customizers, func(p *scope.ClusterScopeParams) {
		p.TagUnmanagedNetworkResources = value
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

func TestService_retrieveZoneInfo(t *testing.T) {
	type testCase struct {
		name           string
		inputZoneNames []string
		expect         func(m *mocks.MockEC2APIMockRecorder)
		want           []types.AvailabilityZone
		wantErrMessage string
	}

	testCases := []*testCase{
		{
			name:           "empty zones",
			inputZoneNames: []string{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{},
					}, nil)
			},
			want: []types.AvailabilityZone{},
		},
		{
			name:           "error describing zones",
			inputZoneNames: []string{},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{},
					}, nil).Return(nil, awserrors.NewNotFound("FailedDescribeAvailableZones"))
			},
			wantErrMessage: `failed to describe availability zones: FailedDescribeAvailableZones`,
		},
		{
			name:           "get type availability zones",
			inputZoneNames: []string{"us-east-1a", "us-east-1b"},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1a", "us-east-1b"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName:       aws.String("us-east-1a"),
								ZoneType:       aws.String("availability-zone"),
								ParentZoneName: nil,
							},
							{
								ZoneName:       aws.String("us-east-1b"),
								ZoneType:       aws.String("availability-zone"),
								ParentZoneName: nil,
							},
						},
					}, nil)
			},
			want: []types.AvailabilityZone{
				{
					ZoneName:       aws.String("us-east-1a"),
					ZoneType:       aws.String("availability-zone"),
					ParentZoneName: nil,
				},
				{
					ZoneName:       aws.String("us-east-1b"),
					ZoneType:       aws.String("availability-zone"),
					ParentZoneName: nil,
				},
			},
		},
		{
			name:           "get type local zones",
			inputZoneNames: []string{"us-east-1-nyc-1a", "us-east-1-bos-1a"},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1-nyc-1a", "us-east-1-bos-1a"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName:       aws.String("us-east-1-nyc-1a"),
								ZoneType:       aws.String("local-zone"),
								ParentZoneName: aws.String("us-east-1a"),
							},
							{
								ZoneName:       aws.String("us-east-1-bos-1a"),
								ZoneType:       aws.String("local-zone"),
								ParentZoneName: aws.String("us-east-1b"),
							},
						},
					}, nil)
			},
			want: []types.AvailabilityZone{
				{
					ZoneName:       aws.String("us-east-1-nyc-1a"),
					ZoneType:       aws.String("local-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
				{
					ZoneName:       aws.String("us-east-1-bos-1a"),
					ZoneType:       aws.String("local-zone"),
					ParentZoneName: aws.String("us-east-1b"),
				},
			},
		},
		{
			name:           "get type wavelength zones",
			inputZoneNames: []string{"us-east-1-wl1-nyc-wlz-1", "us-east-1-wl1-bos-wlz-1"},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1-wl1-nyc-wlz-1", "us-east-1-wl1-bos-wlz-1"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName:       aws.String("us-east-1-wl1-nyc-wlz-1"),
								ZoneType:       aws.String("wavelength-zone"),
								ParentZoneName: aws.String("us-east-1a"),
							},
							{
								ZoneName:       aws.String("us-east-1-wl1-bos-wlz-1"),
								ZoneType:       aws.String("wavelength-zone"),
								ParentZoneName: aws.String("us-east-1b"),
							},
						},
					}, nil)
			},
			want: []types.AvailabilityZone{
				{
					ZoneName:       aws.String("us-east-1-wl1-nyc-wlz-1"),
					ZoneType:       aws.String("wavelength-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
				{
					ZoneName:       aws.String("us-east-1-wl1-bos-wlz-1"),
					ZoneType:       aws.String("wavelength-zone"),
					ParentZoneName: aws.String("us-east-1b"),
				},
			},
		},
		{
			name:           "get all zone types",
			inputZoneNames: []string{"us-east-1a", "us-east-1-nyc-1a", "us-east-1-wl1-nyc-wlz-1"},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
					ZoneNames: []string{"us-east-1a", "us-east-1-nyc-1a", "us-east-1-wl1-nyc-wlz-1"},
				}).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []types.AvailabilityZone{
							{
								ZoneName:       aws.String("us-east-1a"),
								ZoneType:       aws.String("availability-zone"),
								ParentZoneName: nil,
							},
							{
								ZoneName:       aws.String("us-east-1-nyc-1a"),
								ZoneType:       aws.String("local-zone"),
								ParentZoneName: aws.String("us-east-1a"),
							},
							{
								ZoneName:       aws.String("us-east-1-wl1-nyc-wlz-1"),
								ZoneType:       aws.String("wavelength-zone"),
								ParentZoneName: aws.String("us-east-1a"),
							},
						},
					}, nil)
			},
			want: []types.AvailabilityZone{
				{
					ZoneName:       aws.String("us-east-1a"),
					ZoneType:       aws.String("availability-zone"),
					ParentZoneName: nil,
				},
				{
					ZoneName:       aws.String("us-east-1-nyc-1a"),
					ZoneType:       aws.String("local-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
				{
					ZoneName:       aws.String("us-east-1-wl1-nyc-wlz-1"),
					ZoneType:       aws.String("wavelength-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

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

			got, err := s.retrieveZoneInfo(tc.inputZoneNames)
			if err != nil {
				if tc.wantErrMessage != err.Error() {
					t.Errorf("Service.retrieveZoneInfo() error != wanted, got: '%v', want: '%v'", err, tc.wantErrMessage)
				}
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Service.retrieveZoneInfo() = %v, want %v", got, tc.want)
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

// Stub functions to generate AWS mock calls.

func stubGetTags(prefix, role, zone string, isEdge bool) []types.Tag {
	tags := []types.Tag{
		{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("%s-subnet-%s-%s", prefix, role, zone))},
		{Key: aws.String("kubernetes.io/cluster/test-cluster"), Value: aws.String("owned")},
	}
	// tags are returned ordered, inserting LB subnets to prevent diffs...
	if !isEdge {
		lbLabel := "internal-elb"
		if role == "public" {
			lbLabel = "elb"
		}
		tags = append(tags, types.Tag{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/role/%s", lbLabel)),
			Value: aws.String("1"),
		})
	}
	// ... then appending the rest of tags
	tags = append(tags, []types.Tag{
		{Key: aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"), Value: aws.String("owned")},
		{Key: aws.String("sigs.k8s.io/cluster-api-provider-aws/role"), Value: aws.String(role)},
	}...)

	return tags
}

func stubGenMockCreateSubnet(m *mocks.MockEC2APIMockRecorder, prefix, zone, role, cidr string, isEdge bool) *gomock.Call {
	subnet := m.CreateSubnet(context.TODO(), gomock.Eq(&ec2.CreateSubnetInput{
		VpcId:            aws.String(subnetsVPCID),
		CidrBlock:        aws.String(cidr),
		AvailabilityZone: aws.String(zone),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeSubnet,
				Tags:         stubGetTags(prefix, role, zone, isEdge),
			},
		},
	})).
		Return(&ec2.CreateSubnetOutput{
			Subnet: &types.Subnet{
				VpcId:               aws.String(subnetsVPCID),
				SubnetId:            aws.String(fmt.Sprintf("subnet-%s-%s", role, zone)),
				CidrBlock:           aws.String(cidr),
				AvailabilityZone:    aws.String(zone),
				MapPublicIpOnLaunch: aws.Bool(false),
			},
		}, nil)

	// Wait for the subnet to become available
	return m.DescribeSubnets(gomock.Any(), gomock.Eq(&ec2.DescribeSubnetsInput{SubnetIds: []string{fmt.Sprintf("subnet-%s-%s", role, zone)}}), gomock.Any()).
		Return(&ec2.DescribeSubnetsOutput{
			Subnets: []types.Subnet{
				{
					VpcId:            aws.String(subnetsVPCID),
					SubnetId:         aws.String(fmt.Sprintf("subnet-%s-%s", role, zone)),
					CidrBlock:        aws.String(cidr),
					AvailabilityZone: aws.String(zone),
					State:            types.SubnetStateAvailable,
				},
			},
		}, nil).After(subnet)
}

func stubMockCreateTags(m *mocks.MockEC2APIMockRecorder, prefix, name, zone, role string, isEdge bool) *gomock.Call {
	return m.CreateTags(context.TODO(), gomock.Eq(&ec2.CreateTagsInput{
		Resources: []string{name},
		Tags:      stubGetTags(prefix, role, zone, isEdge),
	})).
		Return(&ec2.CreateTagsOutput{}, nil)
}

func stubMockDescribeRouteTables(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
		Return(&ec2.DescribeRouteTablesOutput{}, nil)
}

func stubMockDescribeRouteTablesWithWavelength(m *mocks.MockEC2APIMockRecorder, privSubnets, pubSubnetsIGW, pubSubnetsCarrier []string) *gomock.Call {
	routes := []types.RouteTable{}

	// create public route table
	pubTable := types.RouteTable{
		Routes: []types.Route{
			{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				GatewayId:            aws.String("igw-0"),
			},
		},
		RouteTableId: aws.String("rtb-public"),
	}
	for _, sub := range pubSubnetsIGW {
		pubTable.Associations = append(pubTable.Associations, types.RouteTableAssociation{
			SubnetId: aws.String(sub),
		})
	}
	routes = append(routes, pubTable)

	// create public carrier route table
	pubCarrierTable := types.RouteTable{
		Routes: []types.Route{
			{
				DestinationCidrBlock: aws.String("0.0.0.0/0"),
				CarrierGatewayId:     aws.String("cagw-0"),
			},
		},
		RouteTableId: aws.String("rtb-carrier"),
	}
	for _, sub := range pubSubnetsCarrier {
		pubCarrierTable.Associations = append(pubCarrierTable.Associations, types.RouteTableAssociation{
			SubnetId: aws.String(sub),
		})
	}
	routes = append(routes, pubCarrierTable)

	// create private route table
	privTable := types.RouteTable{
		Routes: []types.Route{
			{
				DestinationCidrBlock: aws.String("10.0.11.0/24"),
				GatewayId:            aws.String("vpc-natgw-1a"),
			},
		},
		RouteTableId: aws.String("rtb-private"),
	}
	for _, sub := range privSubnets {
		privTable.Associations = append(privTable.Associations, types.RouteTableAssociation{
			SubnetId: aws.String(sub),
		})
	}
	routes = append(routes, privTable)

	return m.DescribeRouteTables(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeRouteTablesInput{})).
		Return(&ec2.DescribeRouteTablesOutput{
			RouteTables: routes,
		}, nil)
}

func stubMockDescribeSubnets(m *mocks.MockEC2APIMockRecorder, out *ec2.DescribeSubnetsOutput, filterKey, filterValue string) *gomock.Call {
	return m.DescribeSubnets(context.TODO(), gomock.Eq(&ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{"pending", "available"},
			},
			{
				Name:   aws.String(filterKey),
				Values: []string{filterValue},
			},
		},
	})).
		Return(out, nil)
}

func stubMockDescribeSubnetsWithContextUnmanaged(m *mocks.MockEC2APIMockRecorder) *gomock.Call {
	return stubMockDescribeSubnets(m, &ec2.DescribeSubnetsOutput{
		Subnets: []types.Subnet{
			{SubnetId: aws.String("subnet-az-1a-private"), AvailabilityZone: aws.String("us-east-1a")},
			{SubnetId: aws.String("subnet-az-1a-public"), AvailabilityZone: aws.String("us-east-1a")},
			{SubnetId: aws.String("subnet-lz-1a-private"), AvailabilityZone: aws.String("us-east-1-nyc-1a")},
			{SubnetId: aws.String("subnet-lz-1a-public"), AvailabilityZone: aws.String("us-east-1-nyc-1a")},
			{SubnetId: aws.String("subnet-wl-1a-private"), AvailabilityZone: aws.String("us-east-1-wl1-nyc-wlz-1")},
			{SubnetId: aws.String("subnet-wl-1a-public"), AvailabilityZone: aws.String("us-east-1-wl1-nyc-wlz-1")},
		},
	}, "vpc-id", subnetsVPCID)
}

func stubMockDescribeSubnetsWithContextManaged(m *mocks.MockEC2APIMockRecorder) *gomock.Call {
	return stubMockDescribeSubnets(m, &ec2.DescribeSubnetsOutput{}, "vpc-id", subnetsVPCID)
}

func stubMockDescribeNatGateways(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeNatGateways(context.TODO(),
		gomock.Eq(&ec2.DescribeNatGatewaysInput{
			Filter: []types.Filter{
				{Name: aws.String("vpc-id"), Values: []string{subnetsVPCID}},
				{Name: aws.String("state"), Values: []string{"pending", "available"}},
			},
		}),
		gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)
}

func stubMockModifySubnetAttributeWithContext(m *mocks.MockEC2APIMockRecorder, name string) *gomock.Call {
	return m.ModifySubnetAttribute(context.TODO(), &ec2.ModifySubnetAttributeInput{
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{Value: aws.Bool(true)},
		SubnetId:            aws.String(name),
	}).
		Return(&ec2.ModifySubnetAttributeOutput{}, nil)
}

func stubMockDescribeAvailabilityZonesWithContextAllZones(m *mocks.MockEC2APIMockRecorder) {
	m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
		Return(&ec2.DescribeAvailabilityZonesOutput{
			AvailabilityZones: []types.AvailabilityZone{
				{
					ZoneName:       aws.String("us-east-1a"),
					ZoneType:       aws.String("availability-zone"),
					ParentZoneName: nil,
				},
				{
					ZoneName:       aws.String("us-east-1-nyc-1a"),
					ZoneType:       aws.String("local-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
				{
					ZoneName:       aws.String("us-east-1-wl1-nyc-wlz-1"),
					ZoneType:       aws.String("wavelength-zone"),
					ParentZoneName: aws.String("us-east-1a"),
				},
			},
		}, nil).AnyTimes()
}

func stubMockDescribeAvailabilityZonesWithContextCustomZones(m *mocks.MockEC2APIMockRecorder, zones []types.AvailabilityZone) *gomock.Call {
	return m.DescribeAvailabilityZones(context.TODO(), gomock.Any()).
		Return(&ec2.DescribeAvailabilityZonesOutput{
			AvailabilityZones: zones,
		}, nil).AnyTimes()
}
