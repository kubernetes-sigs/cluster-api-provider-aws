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
	ElasticIPAllocationID = "elastic-ip-allocation-id"
)

func TestReconcileNatGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  []*v1alpha1.Subnet
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.CreateNatGateway(gomock.Any()).Times(0)
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).Times(0)
				m.CreateNatGateway(gomock.Any()).Times(0)
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {

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

				m.DescribeAddresses(gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.AllocateAddress(&ec2.AllocateAddressInput{Domain: aws.String("vpc")}).
					Return(&ec2.AllocateAddressOutput{
						AllocationId: aws.String(ElasticIPAllocationID),
					}, nil)

				m.CreateNatGateway(&ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-1"),
				}).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
					},
				}, nil)

				m.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				}).Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
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
					gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool)
					funct(&ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{{
						NatGatewayId: aws.String("gateway"),
						SubnetId:     aws.String("subnet-1"),
					}}}, true)
				}).Return(nil)

				m.DescribeAddresses(gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.AllocateAddress(&ec2.AllocateAddressInput{Domain: aws.String("vpc")}).
					Return(&ec2.AllocateAddressOutput{
						AllocationId: aws.String(ElasticIPAllocationID),
					}, nil)

				m.CreateNatGateway(&ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-3"),
				}).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
					},
				}, nil)

				m.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				}).Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
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
					gomock.Any()).Do(func(_, y interface{}) {
					funct := y.(func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool)
					funct(&ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{{
						NatGatewayId: aws.String("gateway"),
						SubnetId:     aws.String("subnet-1"),
					}}}, true)
				}).Return(nil)

				m.DescribeAddresses(gomock.Any()).Times(0)
				m.AllocateAddress(gomock.Any()).Times(0)
				m.CreateNatGateway(gomock.Any()).Times(0)
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
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
				Network: v1alpha1.Network{
					VPC: v1alpha1.VPC{
						ID: subnetsVPCID,
					},
					Subnets: tc.input,
				},
			}

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			if err := s.reconcileNatGateways(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}
