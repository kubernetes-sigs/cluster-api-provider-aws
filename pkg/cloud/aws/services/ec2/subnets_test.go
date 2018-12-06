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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
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
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeAvailabilityZones(gomock.AssignableToTypeOf(&ec2.DescribeAvailabilityZonesInput{})).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{
								RegionName: aws.String("us-east-1"),
								ZoneName:   aws.String("us-east-1a"),
							},
						},
					}, nil)

				m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
						{
							Name:   aws.String("tag-key"),
							Values: []*string{aws.String("kubernetes.io/cluster/test-cluster")},
						},
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
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
						},
					}, nil)

				m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
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

				m.WaitUntilSubnetAvailable(gomock.Any())

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.ModifySubnetAttribute(&ec2.ModifySubnetAttributeInput{
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				describeCall := m.DescribeSubnets(gomock.Eq(&ec2.DescribeSubnetsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []*string{aws.String(subnetsVPCID)},
						},
						{
							Name:   aws.String("tag-key"),
							Values: []*string{aws.String("kubernetes.io/cluster/test-cluster")},
						},
						{
							Name:   aws.String("state"),
							Values: []*string{aws.String("pending"), aws.String("available")},
						},
					},
				})).
					Return(&ec2.DescribeSubnetsOutput{}, nil)

				firstSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
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

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(firstSubnet)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				secondSubnet := m.CreateSubnet(gomock.Eq(&ec2.CreateSubnetInput{
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

				m.WaitUntilSubnetAvailable(gomock.Any()).
					After(secondSubnet)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{}))

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

			scope.ClusterStatus = &v1alpha1.AWSClusterProviderStatus{
				Network: *tc.input,
			}

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			if err := s.reconcileSubnets(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
