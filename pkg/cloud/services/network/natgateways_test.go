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
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
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

const (
	ElasticIPAllocationID = "elastic-ip-allocation-id"
)

func TestReconcileNatGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  []infrav1.SubnetSpec
		expect func(m *mocks.MockEC2APIMockRecorder)
	}{
		{
			name: "single private subnet exists, should create no NAT gateway",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.CreateNatGateway(context.TODO(), gomock.Any()).Times(0)
			},
		},
		{
			name: "no private subnet exists, should create no NAT gateway",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.Any(), gomock.Any()).Times(0)
				m.CreateNatGateway(context.TODO(), gomock.Any()).Times(0)
			},
		},
		{
			name: "public & private subnet exists, should create 1 NAT gateway",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
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

				m.DescribeAddresses(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.AllocateAddress(context.TODO(), &ec2.AllocateAddressInput{
					Domain: types.DomainTypeVpc,
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeElasticIp,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-eip-common"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
							},
						},
					},
				}).Return(&ec2.AllocateAddressOutput{
					AllocationId: aws.String(ElasticIPAllocationID),
				}, nil)

				m.CreateNatGateway(context.TODO(), &ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-1"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeNatgateway,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-nat"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
							},
						},
					},
				},
				).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &types.NatGateway{
						NatGatewayId: aws.String("natgateway"),
						SubnetId:     aws.String("subnet-1"),
					},
				}, nil)

				m.DescribeNatGateways(gomock.Any(), &ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
				}, gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							State:        types.NatGatewayStateAvailable,
							NatGatewayId: aws.String("natgateway"),
							SubnetId:     aws.String("subnet-1"),
						},
					},
				}, nil)
			},
		},
		{
			name: "two public & 1 private subnet, and one NAT gateway exists",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
				{
					ID:               "subnet-3",
					AvailabilityZone: "us-east-1b",
					CidrBlock:        "10.0.13.0/24",
					IsPublic:         true,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
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
					gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("gateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DescribeAddresses(context.TODO(), gomock.Any()).
					Return(&ec2.DescribeAddressesOutput{}, nil)

				m.AllocateAddress(context.TODO(), &ec2.AllocateAddressInput{
					Domain: types.DomainTypeVpc,
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeElasticIp,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-eip-common"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
							},
						},
					},
				}).Return(&ec2.AllocateAddressOutput{
					AllocationId: aws.String(ElasticIPAllocationID),
				}, nil)

				m.CreateNatGateway(context.TODO(), &ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-3"),
					TagSpecifications: []types.TagSpecification{
						{
							ResourceType: types.ResourceTypeNatgateway,
							Tags: []types.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-nat"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
							},
						},
					},
				}).Return(&ec2.CreateNatGatewayOutput{
					NatGateway: &types.NatGateway{
						NatGatewayId: aws.String("natgateway"),
						SubnetId:     aws.String("subnet-3"),
					},
				}, nil)

				m.DescribeNatGateways(gomock.Any(), &ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
				}, gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							State:        types.NatGatewayStateAvailable,
							NatGatewayId: aws.String("natgateway"),
							SubnetId:     aws.String("subnet-3"),
						},
					},
				}, nil)

				m.CreateTags(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil).Times(1)
			},
		},
		{
			name: "public & private subnet, and one NAT gateway exists",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
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
					gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("gateway"),
							SubnetId:     aws.String("subnet-1"),
							Tags: []types.Tag{
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("common"),
								},
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-nat"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
							},
						},
					}}, nil)

				m.DescribeAddresses(context.TODO(), gomock.Any()).Times(0)
				m.AllocateAddress(context.TODO(), gomock.Any()).Times(0)
				m.CreateNatGateway(context.TODO(), gomock.Any()).Times(0)
			},
		},
		{
			name: "public & private subnet declared, but don't exist yet",
			input: []infrav1.SubnetSpec{
				{
					ID:               "",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{}, nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: subnetsVPCID,
							Tags: infrav1.Tags{
								infrav1.ClusterTagKey("test-cluster"): "owned",
							},
						},
						Subnets: tc.input,
					},
				},
			}
			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: awsCluster,
				Client:     client,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			if err := s.reconcileNatGateways(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteNatGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name           string
		input          []infrav1.SubnetSpec
		isUnmanagedVPC bool
		expect         func(m *mocks.MockEC2APIMockRecorder)
		wantErr        bool
	}{
		{
			name:           "Should skip deletion if vpc is unmanaged",
			isUnmanagedVPC: true,
		},
		{
			name: "Should skip deletion if no private subnet is present",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
			},
		},
		{
			name: "Should skip deletion if no public subnet is present",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         false,
				},
			},
		},
		{
			name: "Should skip deletion if no existing natgateway present",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					Filter: []types.Filter{
						{
							Name:   aws.String("vpc-id"),
							Values: []string{"managed-vpc"},
						},
						{
							Name:   aws.String("state"),
							Values: []string{"pending", "available"},
						},
					},
				}), gomock.Any()).Return(&ec2.DescribeNatGatewaysOutput{}, nil)
			},
		},
		{
			name: "Should successfully delete natgateways",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("natgateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DeleteNatGateway(context.TODO(), gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
				})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							State: types.NatGatewayStateAvailable,
						},
					},
				}, nil)
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							State: types.NatGatewayStateDeleted,
						},
					},
				}, nil)
			},
		},
		{
			name: "Should return error if natgateway has unknown state",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.12.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-3",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.14.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("natgateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DeleteNatGateway(context.TODO(), gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
				})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							State: types.NatGatewayStatePending,
						},
					},
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "Should return error if describe natgateway output and error, both are nil",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.14.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("natgateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DeleteNatGateway(context.TODO(), gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
				})).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "Should return error if describe natgateway pages fails",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.14.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
			wantErr: true,
		},
		{
			name: "Should return error if delete natgateway fails",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.14.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("natgateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DeleteNatGateway(context.TODO(), gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
			wantErr: true,
		},
		{
			name: "Should return error if describe natgateway fails",
			input: []infrav1.SubnetSpec{
				{
					ID:               "subnet-1",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.10.0/24",
					IsPublic:         true,
				},
				{
					ID:               "subnet-2",
					AvailabilityZone: "us-east-1a",
					CidrBlock:        "10.0.14.0/24",
					IsPublic:         false,
				},
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeNatGateways(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("natgateway"),
								SubnetId:     aws.String("subnet-1"),
							},
						},
					}, nil)

				m.DeleteNatGateway(context.TODO(), gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(context.TODO(), gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []string{"natgateway"},
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
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: "managed-vpc",
							Tags: infrav1.Tags{
								infrav1.ClusterTagKey("test-cluster"): "owned",
							},
						},
						Subnets: tc.input,
					},
				},
			}
			if tc.isUnmanagedVPC {
				awsCluster.Spec.NetworkSpec.VPC.Tags = infrav1.Tags{}
			}
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: awsCluster,
				Client:     client,
			})
			g.Expect(err).NotTo(HaveOccurred())
			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}

			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			err = s.deleteNatGateways()
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestGetdNatGatewayForEdgeSubnet(t *testing.T) {
	subnetsSpec := infrav1.Subnets{
		{
			ID:               "subnet-az-1x-private",
			AvailabilityZone: "us-east-1x",
			IsPublic:         false,
		},
		{
			ID:               "subnet-az-1x-public",
			AvailabilityZone: "us-east-1x",
			IsPublic:         true,
			NatGatewayID:     aws.String("natgw-az-1b-last"),
		},
		{
			ID:               "subnet-az-1a-private",
			AvailabilityZone: "us-east-1a",
			IsPublic:         false,
		},
		{
			ID:               "subnet-az-1a-public",
			AvailabilityZone: "us-east-1a",
			IsPublic:         true,
			NatGatewayID:     aws.String("natgw-az-1b-first"),
		},
		{
			ID:               "subnet-az-1b-private",
			AvailabilityZone: "us-east-1b",
			IsPublic:         false,
		},
		{
			ID:               "subnet-az-1b-public",
			AvailabilityZone: "us-east-1b",
			IsPublic:         true,
			NatGatewayID:     aws.String("natgw-az-1b-second"),
		},
		{
			ID:               "subnet-az-1p-private",
			AvailabilityZone: "us-east-1p",
			IsPublic:         false,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name             string
		spec             infrav1.Subnets
		input            infrav1.SubnetSpec
		expect           string
		expectErr        bool
		expectErrMessage string
	}{
		{
			name: "zone availability-zone, valid nat gateway",
			input: infrav1.SubnetSpec{
				ID:               "subnet-az-1b-private",
				AvailabilityZone: "us-east-1b",
				IsPublic:         false,
			},
			expect: "natgw-az-1b-second",
		},
		{
			name: "zone availability-zone, valid nat gateway",
			input: infrav1.SubnetSpec{
				ID:               "subnet-az-1a-private",
				AvailabilityZone: "us-east-1a",
				IsPublic:         false,
			},
			expect: "natgw-az-1b-first",
		},
		{
			name: "zone availability-zone, valid nat gateway",
			input: infrav1.SubnetSpec{
				ID:               "subnet-az-1x-private",
				AvailabilityZone: "us-east-1x",
				IsPublic:         false,
			},
			expect: "natgw-az-1b-last",
		},
		{
			name: "zone local-zone, valid nat gateway from parent",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-nyc1a-private",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
				ParentZoneName:   aws.String("us-east-1a"),
			},
			expect: "natgw-az-1b-first",
		},
		{
			name: "zone local-zone, valid nat gateway from parent",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-nyc1a-private",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
				ParentZoneName:   aws.String("us-east-1x"),
			},
			expect: "natgw-az-1b-last",
		},
		{
			name: "zone local-zone, valid nat gateway from fallback",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-nyc1a-private",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
				ParentZoneName:   aws.String("us-east-1-notAvailable"),
			},
			expect: "natgw-az-1b-first",
		},
		{
			name: "edge zones without NAT GW support, no public subnet and NAT Gateway for the parent zone, return first nat gateway available",
			input: infrav1.SubnetSpec{
				ID:               "subnet-7",
				AvailabilityZone: "us-east-1-nyc-1a",
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
			},
			expect: "natgw-az-1b-first",
		},
		{
			name: "edge zones without NAT GW support, no public subnet and NAT Gateway for the parent zone, return first nat gateway available",
			input: infrav1.SubnetSpec{
				ID:               "subnet-7",
				CidrBlock:        "10.0.10.0/24",
				AvailabilityZone: "us-east-1-nyc-1a",
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
				ParentZoneName:   aws.String("us-east-1-notFound"),
			},
			expect: "natgw-az-1b-first",
		},
		{
			name: "edge zones without NAT GW support, valid public subnet and NAT Gateway for the parent zone, return parent's zone nat gateway",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-7",
				AvailabilityZone: "us-east-1-nyc-1a",
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
				ParentZoneName:   aws.String("us-east-1b"),
			},
			expect: "natgw-az-1b-second",
		},
		{
			name: "wavelength zones without Nat GW support, public subnet and Nat Gateway for the parent zone, return parent's zone nat gateway",
			input: infrav1.SubnetSpec{
				ID:               "subnet-7",
				CidrBlock:        "10.0.10.0/24",
				AvailabilityZone: "us-east-1-wl1-nyc-wlz-1",
				ZoneType:         ptr.To(infrav1.ZoneTypeWavelengthZone),
				ParentZoneName:   aws.String("us-east-1x"),
			},
			expect: "natgw-az-1b-last",
		},
		// errors
		{
			name: "error if the subnet is public",
			input: infrav1.SubnetSpec{
				ID:               "subnet-az-1-public",
				AvailabilityZone: "us-east-1a",
				IsPublic:         true,
			},
			expectErr:        true,
			expectErrMessage: `cannot get NAT gateway for a public subnet, got id "subnet-az-1-public"`,
		},
		{
			name: "error if the subnet is public",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-1-public",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         true,
			},
			expectErr:        true,
			expectErrMessage: `cannot get NAT gateway for a public subnet, got id "subnet-lz-1-public"`,
		},
		{
			name: "error if there are no nat gateways available in the subnets",
			spec: infrav1.Subnets{},
			input: infrav1.SubnetSpec{
				ID:               "subnet-az-1-private",
				AvailabilityZone: "us-east-1p",
				IsPublic:         false,
			},
			expectErr:        true,
			expectErrMessage: `no nat gateways available in "us-east-1p" for private subnet "subnet-az-1-private"`,
		},
		{
			name: "error if there are no nat gateways available in the subnets",
			spec: infrav1.Subnets{},
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-1",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         false,
				ZoneType:         ptr.To(infrav1.ZoneTypeLocalZone),
			},
			expectErr:        true,
			expectErrMessage: `no nat gateways available in "us-east-1-nyc-1a" for private edge subnet "subnet-lz-1", current state: map[]`,
		},
		{
			name: "error if the subnet is public",
			input: infrav1.SubnetSpec{
				ID:               "subnet-lz-1",
				AvailabilityZone: "us-east-1-nyc-1a",
				IsPublic:         true,
			},
			expectErr:        true,
			expectErrMessage: `cannot get NAT gateway for a public subnet, got id "subnet-lz-1"`,
		},
	}

	for idx, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			subnets := subnetsSpec
			if tc.spec != nil {
				subnets = tc.spec
			}
			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			awsCluster := &infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: infrav1.VPCSpec{
							ID: subnetsVPCID,
							Tags: infrav1.Tags{
								infrav1.ClusterTagKey("test-cluster"): "owned",
							},
						},
						Subnets: subnets,
					},
				},
			}

			client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

			clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSCluster: awsCluster,
				Client:     client,
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
				return
			}

			s := NewService(clusterScope)

			id, err := s.getNatGatewayForSubnet(&testCases[idx].input)

			if tc.expectErr && err == nil {
				t.Fatal("expected error but got no error")
			}
			if err != nil && len(tc.expectErrMessage) > 0 {
				if err.Error() != tc.expectErrMessage {
					t.Fatalf("got an unexpected error message:\nwant: %v\n got: %v\n", tc.expectErrMessage, err.Error())
				}
			}
			if !tc.expectErr && err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}
			if len(tc.expect) > 0 {
				g.Expect(id).To(Equal(tc.expect))
			}
		})
	}
}
