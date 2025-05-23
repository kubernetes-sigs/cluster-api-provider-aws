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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func describeVpcAttributeTrue(_ context.Context, input *ec2.DescribeVpcAttributeInput, _ ...request.Option) (*ec2.DescribeVpcAttributeOutput, error) {
	result := &ec2.DescribeVpcAttributeOutput{
		VpcId: input.VpcId,
	}
	switch input.Attribute {
	case types.VpcAttributeNameEnableDnsHostnames:
		result.EnableDnsHostnames = &types.AttributeBooleanValue{Value: aws.Bool(true)}
	case types.VpcAttributeNameEnableDnsSupport:
		result.EnableDnsSupport = &types.AttributeBooleanValue{Value: aws.Bool(true)}
	}
	return result, nil
}

func describeVpcAttributeFalse(_ context.Context, input *ec2.DescribeVpcAttributeInput, _ ...request.Option) (*ec2.DescribeVpcAttributeOutput, error) {
	result := &ec2.DescribeVpcAttributeOutput{
		VpcId: input.VpcId,
	}
	switch input.Attribute {
	case types.VpcAttributeNameEnableDnsHostnames:
		result.EnableDnsHostnames = &types.AttributeBooleanValue{Value: aws.Bool(false)}
	case types.VpcAttributeNameEnableDnsSupport:
		result.EnableDnsSupport = &types.AttributeBooleanValue{Value: aws.Bool(false)}
	}
	return result, nil
}

func TestReconcileVPC(t *testing.T) {
	usageLimit := 3
	selection := infrav1.AZSelectionSchemeOrdered
	managedVPCTags := []types.Tag{
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
			Value: aws.String("common"),
		},
		{
			Key:   aws.String("Name"),
			Value: aws.String("test-cluster-vpc"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
			Value: aws.String("owned"), // = managed by CAPA
		},
	}

	testCases := []struct {
		name              string
		input             *infrav1.VPCSpec
		want              *infrav1.VPCSpec
		additionalTags    map[string]string
		expect            func(m *mocks.MockEC2APIMockRecorder)
		wantErrContaining *string // nil to assert success
	}{
		{
			name:  "Should update tags with aws VPC resource tags, if managed vpc exists",
			input: &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			want: &infrav1.VPCSpec{
				ID:        "vpc-exists",
				CidrBlock: "10.0.0.0/8",
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			wantErrContaining: nil,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []string{
						"vpc-exists",
					},
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{string(types.VpcStatePending), string(types.VpcStateAvailable)},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      managedVPCTags,
						},
					},
				}, nil)

				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			// I need additional tags in scope and make sure they are applied
			name:  "Should ensure tags after creation remain the same",
			input: &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			additionalTags: map[string]string{
				"additional": "tags",
			},
			want: &infrav1.VPCSpec{
				ID:        "vpc-exists",
				CidrBlock: "10.0.0.0/8",
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			wantErrContaining: nil,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []string{
						"vpc-exists",
					},
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{string(types.VpcStatePending), string(types.VpcStateAvailable)},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      managedVPCTags,
						},
					},
				}, nil)
				m.CreateTags(context.TODO(), &ec2.CreateTagsInput{
					Resources: []string{"vpc-exists"},
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-cluster-vpc"),
						},
						{
							Key:   aws.String("additional"),
							Value: aws.String("tags"),
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
				})
				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:              "Should create a new VPC if managed vpc does not exist",
			input:             &infrav1.VPCSpec{AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: nil,
			want: &infrav1.VPCSpec{
				ID:        "vpc-new",
				CidrBlock: "10.1.0.0/16",
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeVPCByNameCall := m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("tag:Name"),
							Values: []string{"test-cluster-vpc"},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{}}, nil)
				m.CreateVpc(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).After(describeVPCByNameCall).Return(&ec2.CreateVpcOutput{
					Vpc: &types.Vpc{
						State:     types.VpcStateAvailable,
						VpcId:     aws.String("vpc-new"),
						CidrBlock: aws.String("10.1.0.0/16"),
						Tags:      managedVPCTags,
					},
				}, nil)

				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)
			},
		},
		{
			name: "Should amend attributes of existing VPC",
			input: &infrav1.VPCSpec{
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
				IPv6:                       &infrav1.IPv6{},
			},
			wantErrContaining: nil,
			want: &infrav1.VPCSpec{
				ID:        "vpc-new",
				CidrBlock: "10.1.0.0/16",
				IPv6: &infrav1.IPv6{
					CidrBlock: "2001:db8:1234:1a03::/56",
					PoolID:    "amazon",
				},
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{
					VpcIds: []string{"vpc-new"},
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							CidrBlock: aws.String("10.1.0.0/16"),
							Ipv6CidrBlockAssociationSet: []types.VpcIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/56"),
									Ipv6CidrBlockState: &types.VpcCidrBlockState{
										State: types.VpcCidrBlockStateCodeAssociated,
									},
									Ipv6Pool: aws.String("amazon"),
								},
							},
							State: types.VpcStateAvailable,
							Tags:  managedVPCTags,
							VpcId: aws.String("vpc-new"),
						},
					},
				}, nil)

				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)
			},
		},
		{
			name: "Should create a new IPv6 VPC with BYOIP set up if managed IPv6 vpc does not exist",
			input: &infrav1.VPCSpec{
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
				IPv6: &infrav1.IPv6{
					CidrBlock: "2001:db8:1234:1a03::/56",
					PoolID:    "my-pool",
				},
			},
			wantErrContaining: nil,
			want: &infrav1.VPCSpec{
				ID:        "vpc-new",
				CidrBlock: "10.1.0.0/16",
				IPv6: &infrav1.IPv6{
					CidrBlock: "2001:db8:1234:1a03::/56",
					PoolID:    "my-pool",
				},
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeVPCByNameCall := m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("tag:Name"),
							Values: []string{"test-cluster-vpc"},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{}}, nil)
				m.CreateVpc(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateVpcInput{
					AmazonProvidedIpv6CidrBlock: aws.Bool(false),
					Ipv6Pool:                    aws.String("my-pool"),
					Ipv6CidrBlock:               aws.String("2001:db8:1234:1a03::/56"),
				})).After(describeVPCByNameCall).Return(&ec2.CreateVpcOutput{
					Vpc: &types.Vpc{
						State:     types.VpcStateAvailable,
						VpcId:     aws.String("vpc-new"),
						CidrBlock: aws.String("10.1.0.0/16"),
						Tags:      managedVPCTags,
					},
				}, nil)

				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)
			},
		},
		{
			name: "Describing the VPC fails with IPv6 VPC should return an error",
			input: &infrav1.VPCSpec{
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
				IPv6:                       &infrav1.IPv6{},
			},
			wantErrContaining: aws.String("nope"),
			want:              nil,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{
					VpcIds: []string{"vpc-new"},
				})).Return(nil, errors.New("nope"))
			},
		},
		{
			name: "should set up IPv6 associations if found VPC is IPv6 enabled",
			input: &infrav1.VPCSpec{
				ID:                         "unmanaged-vpc-exists",
				IPv6:                       &infrav1.IPv6{},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			want: &infrav1.VPCSpec{
				ID:        "unmanaged-vpc-exists",
				CidrBlock: "10.0.0.0/8",
				Tags:      nil,
				IPv6: &infrav1.IPv6{
					PoolID:    "my-pool",
					CidrBlock: "2001:db8:1234:1a03::/56",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Ipv6CidrBlockAssociationSet: []types.VpcIpv6CidrBlockAssociation{
								{
									AssociationId: aws.String("amazon"),
									Ipv6CidrBlock: aws.String("2001:db8:1234:1a03::/56"),
									Ipv6CidrBlockState: &types.VpcCidrBlockState{
										State: types.VpcCidrBlockStateCodeAssociated,
									},
									Ipv6Pool: aws.String("my-pool"),
								},
							},
						},
					},
				}, nil)
			},
		},
		{
			name:              "managed vpc id exists, but vpc resource is missing",
			input:             &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("VPC resource is missing in AWS"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []string{
						"vpc-exists",
					},
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{string(types.VpcStatePending), string(types.VpcStateAvailable)},
						},
					},
				})).Return(nil, awserr.New("404", "http not found err", errors.New("err")))
			},
		},
		{
			name:  "Should patch vpc spec successfully, if unmanaged vpc exists",
			input: &infrav1.VPCSpec{ID: "unmanaged-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			want: &infrav1.VPCSpec{
				ID:                         "unmanaged-vpc-exists",
				CidrBlock:                  "10.0.0.0/8",
				Tags:                       nil,
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
						},
					},
				}, nil)
			},
		},
		{
			name:  "Should retry if vpc not found error occurs during attributes configuration for managed vpc",
			input: &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			want: &infrav1.VPCSpec{
				ID:        "managed-vpc-exists",
				CidrBlock: "10.0.0.0/8",
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
					"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      managedVPCTags,
						},
					},
				}, nil)
				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).Return(nil, awserr.New("InvalidVpcID.NotFound", "not found", nil))
				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:              "Should return error if failed to set vpc attributes for managed vpc",
			input:             &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("failed to set vpc attributes"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []string{
						"managed-vpc-exists",
					},
					Filters: []types.Filter{
						{
							Name:   aws.String("state"),
							Values: []string{string(types.VpcStatePending), string(types.VpcStateAvailable)},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							State:     types.VpcStateAvailable,
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      managedVPCTags,
						},
					},
				}, nil)
				m.DescribeVpcAttribute(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).AnyTimes().Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
		},
		{
			name:              "Should return error if failed to create vpc",
			input:             &infrav1.VPCSpec{AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("failed to create new managed VPC"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				describeVPCByNameCall := m.DescribeVpcs(context.TODO(), gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []types.Filter{
						{
							Name:   aws.String("tag:Name"),
							Values: []string{"test-cluster-vpc"},
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{}}, nil)
				m.CreateVpc(context.TODO(), gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).After(describeVPCByNameCall).Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
		},
		{
			name:              "Should return error if describe vpc returns empty list",
			input:             &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("VPC resource is missing in AWS"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{},
				}, nil)
			},
		},
		{
			name:              "Should return error if describe vpc returns more than 1 vpcs",
			input:             &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("Only one VPC per cluster name is supported"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							VpcId: aws.String("vpc_1"),
						},
						{
							VpcId: aws.String("vpc_2"),
						},
					},
				}, nil)
			},
		},
		{
			name:              "Should return error if vpc state is not available/pending",
			input:             &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErrContaining: aws.String("could not find available or pending vpc"),
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(context.TODO(), gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []types.Vpc{
						{
							VpcId: aws.String("vpc"),
							State: types.VpcState("deleting"),
						},
					},
				}, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			g := NewWithT(t)
			clusterScope, err := getClusterScope(tc.input, tc.additionalTags)
			g.Expect(err).NotTo(HaveOccurred())
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())
			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			err = s.reconcileVPC()
			if tc.wantErrContaining != nil {
				g.Expect(err).ToNot(BeNil())
				g.Expect(err.Error()).To(ContainSubstring(*tc.wantErrContaining))
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(tc.want).To(Equal(&clusterScope.AWSCluster.Spec.NetworkSpec.VPC))
		})
	}
}

func TestDeleteVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tags := map[string]string{
		"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
	}

	testCases := []struct {
		name           string
		input          *infrav1.VPCSpec
		additionalTags map[string]string
		wantErr        bool
		expect         func(m *mocks.MockEC2APIMockRecorder)
	}{
		{
			name:  "Should not delete vpc if vpc is unmanaged",
			input: &infrav1.VPCSpec{ID: "unmanaged-vpc"},
		},
		{
			name: "Should return error if delete vpc failed",
			input: &infrav1.VPCSpec{
				ID:   "managed-vpc",
				Tags: tags,
			},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteVpc(context.TODO(), gomock.Eq(&ec2.DeleteVpcInput{
					VpcId: aws.String("managed-vpc"),
				})).Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
		},
		{
			name: "Should return without error if delete vpc succeeded",
			input: &infrav1.VPCSpec{
				ID:   "managed-vpc",
				Tags: tags,
			},
			wantErr: false,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteVpc(context.TODO(), gomock.Eq(&ec2.DeleteVpcInput{
					VpcId: aws.String("managed-vpc"),
				})).Return(&ec2.DeleteVpcOutput{}, nil)
			},
		},
		{
			name: "Should not delete vpc if vpc not found",
			input: &infrav1.VPCSpec{
				ID:   "managed-vpc",
				Tags: tags,
			},
			wantErr: false,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteVpc(context.TODO(), gomock.Eq(&ec2.DeleteVpcInput{
					VpcId: aws.String("managed-vpc"),
				})).Return(nil, awserr.New("InvalidVpcID.NotFound", "not found", nil))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			clusterScope, err := getClusterScope(tc.input, tc.additionalTags)
			g.Expect(err).NotTo(HaveOccurred())
			if tc.expect != nil {
				tc.expect(ec2Mock.EXPECT())
			}
			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			err = s.deleteVPC()
			if tc.wantErr {
				g.Expect(err).ToNot(BeNil())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func getClusterScope(vpcSpec *infrav1.VPCSpec, additionalTags map[string]string) (*scope.ClusterScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)

	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: infrav1.AWSClusterSpec{
			NetworkSpec: infrav1.NetworkSpec{
				VPC: *vpcSpec,
			},
			AdditionalTags: additionalTags,
		},
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(awsCluster).WithStatusSubresource(awsCluster).Build()

	return scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		AWSCluster: awsCluster,
		Client:     client,
	})
}
