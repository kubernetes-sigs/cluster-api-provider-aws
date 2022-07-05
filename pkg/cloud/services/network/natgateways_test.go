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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/test/mocks"
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
				m.CreateNatGateway(gomock.Any()).Times(0)
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
				m.DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).Times(0)
				m.CreateNatGateway(gomock.Any()).Times(0)
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

				m.AllocateAddress(&ec2.AllocateAddressInput{
					Domain: aws.String("vpc"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("elastic-ip"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-eip-apiserver"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("apiserver"),
								},
							},
						},
					},
				}).Return(&ec2.AllocateAddressOutput{
					AllocationId: aws.String(ElasticIPAllocationID),
				}, nil)

				m.CreateNatGateway(&ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-1"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("natgateway"),
							Tags: []*ec2.Tag{
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
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
						SubnetId:     aws.String("subnet-1"),
					},
				}, nil)

				m.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				}).Return(nil)
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

				m.AllocateAddress(&ec2.AllocateAddressInput{
					Domain: aws.String("vpc"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("elastic-ip"),
							Tags: []*ec2.Tag{
								{
									Key:   aws.String("Name"),
									Value: aws.String("test-cluster-eip-apiserver"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
									Value: aws.String("owned"),
								},
								{
									Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
									Value: aws.String("apiserver"),
								},
							},
						},
					},
				}).Return(&ec2.AllocateAddressOutput{
					AllocationId: aws.String(ElasticIPAllocationID),
				}, nil)

				m.CreateNatGateway(&ec2.CreateNatGatewayInput{
					AllocationId: aws.String(ElasticIPAllocationID),
					SubnetId:     aws.String("subnet-3"),
					TagSpecifications: []*ec2.TagSpecification{
						{
							ResourceType: aws.String("natgateway"),
							Tags: []*ec2.Tag{
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
					NatGateway: &ec2.NatGateway{
						NatGatewayId: aws.String("natgateway"),
						SubnetId:     aws.String("subnet-3"),
					},
				}, nil)

				m.WaitUntilNatGatewayAvailable(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				}).Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
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
						Tags: []*ec2.Tag{
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
					}}}, true)
				}).Return(nil)

				m.DescribeAddresses(gomock.Any()).Times(0)
				m.AllocateAddress(gomock.Any()).Times(0)
				m.CreateNatGateway(gomock.Any()).Times(0)
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
				m.DescribeNatGatewaysPages(gomock.Any(), gomock.Any()).
					Return(nil).
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
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			ctx := context.TODO()
			client.Create(ctx, awsCluster)
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
				m.DescribeNatGatewaysPages(
					gomock.Eq(&ec2.DescribeNatGatewaysInput{
						Filter: []*ec2.Filter{
							{
								Name:   aws.String("vpc-id"),
								Values: []*string{aws.String("managed-vpc")},
							},
							{
								Name:   aws.String("state"),
								Values: []*string{aws.String("pending"), aws.String("available")},
							},
						},
					}),
					gomock.Any()).Return(nil)
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}),
					gomock.Any()).Do(mockDescribeNatGatewaysOutput).Return(nil)

				m.DeleteNatGateway(gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []*ec2.NatGateway{
						{
							State: aws.String("available"),
						},
					},
				}, nil)
				m.DescribeNatGateways(gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []*ec2.NatGateway{
						{
							State: aws.String("deleted"),
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).Do(mockDescribeNatGatewaysOutput).Return(nil)

				m.DeleteNatGateway(gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
				})).Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []*ec2.NatGateway{
						{
							State: aws.String("unknown"),
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).Do(mockDescribeNatGatewaysOutput).Return(nil)

				m.DeleteNatGateway(gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).Return(awserrors.NewFailedDependency("failed dependency"))
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).Do(mockDescribeNatGatewaysOutput).Return(nil)

				m.DeleteNatGateway(gomock.Eq(&ec2.DeleteNatGatewayInput{
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
				m.DescribeNatGatewaysPages(
					gomock.AssignableToTypeOf(&ec2.DescribeNatGatewaysInput{}), gomock.Any()).Do(mockDescribeNatGatewaysOutput).Return(nil)

				m.DeleteNatGateway(gomock.Eq(&ec2.DeleteNatGatewayInput{
					NatGatewayId: aws.String("natgateway"),
				})).Return(&ec2.DeleteNatGatewayOutput{}, nil)

				m.DescribeNatGateways(gomock.Eq(&ec2.DescribeNatGatewaysInput{
					NatGatewayIds: []*string{aws.String("natgateway")},
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

var mockDescribeNatGatewaysOutput = func(_, y interface{}) {
	funct := y.(func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool)
	funct(&ec2.DescribeNatGatewaysOutput{NatGateways: []*ec2.NatGateway{{
		NatGatewayId: aws.String("natgateway"),
		SubnetId:     aws.String("subnet-1"),
	}}}, true)
}
