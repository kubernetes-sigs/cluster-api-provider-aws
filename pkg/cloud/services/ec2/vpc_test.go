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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
)

func describeVpcAttributeTrue(input *ec2.DescribeVpcAttributeInput) (*ec2.DescribeVpcAttributeOutput, error) {
	result := &ec2.DescribeVpcAttributeOutput{
		VpcId: input.VpcId,
	}
	switch aws.StringValue(input.Attribute) {
	case "enableDnsHostnames":
		result.EnableDnsHostnames = &ec2.AttributeBooleanValue{Value: aws.Bool(true)}
	case "enableDnsSupport":
		result.EnableDnsSupport = &ec2.AttributeBooleanValue{Value: aws.Bool(true)}
	}
	return result, nil
}

func describeVpcAttributeFalse(input *ec2.DescribeVpcAttributeInput) (*ec2.DescribeVpcAttributeOutput, error) {
	result := &ec2.DescribeVpcAttributeOutput{
		VpcId: input.VpcId,
	}
	switch aws.StringValue(input.Attribute) {
	case "enableDnsHostnames":
		result.EnableDnsHostnames = &ec2.AttributeBooleanValue{Value: aws.Bool(false)}
	case "enableDnsSupport":
		result.EnableDnsSupport = &ec2.AttributeBooleanValue{Value: aws.Bool(false)}
	}
	return result, nil
}

func TestReconcileVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *infrav1.VPCSpec
		output *infrav1.VPCSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name:   "managed vpc exists",
			input:  &infrav1.VPCSpec{ID: "vpc-exists"},
			output: &infrav1.VPCSpec{ID: "vpc-exists", CidrBlock: "10.0.0.0/8"},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []*string{
						aws.String("vpc-exists"),
					},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
						},
					},
				})).
					Return(&ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							{
								State:     aws.String("available"),
								VpcId:     aws.String("vpc-exists"),
								CidrBlock: aws.String("10.0.0.0/8"),
								Tags: []*ec2.Tag{
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
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:   "managed vpc exists and ipv6 enabled",
			input:  &infrav1.VPCSpec{ID: "vpc-exists", EnableIPv6: true},
			output: &infrav1.VPCSpec{ID: "vpc-exists", CidrBlock: "10.0.0.0/8", EnableIPv6: true, Ipv6CidrBlock: aws.String("2001:10:10:10::/56")},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []*string{
						aws.String("vpc-exists"),
					},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
						},
						{
							Name:   aws.String("ipv6-cidr-block-association.state"),
							Values: aws.StringSlice([]string{ec2.VpcCidrBlockStateCodeAssociated}),
						},
					},
				})).
					Return(&ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							{
								State:     aws.String("available"),
								VpcId:     aws.String("vpc-exists"),
								CidrBlock: aws.String("10.0.0.0/8"),
								Ipv6CidrBlockAssociationSet: []*ec2.VpcIpv6CidrBlockAssociation{
									{
										Ipv6CidrBlock: aws.String("2001:10:10:10::/56"),
										Ipv6CidrBlockState: &ec2.VpcCidrBlockState{
											State: aws.String(ec2.VpcCidrBlockStateCodeAssociated),
										},
									},
								},
								Tags: []*ec2.Tag{
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
										Value: aws.String("owned"),
									},
								},
							},
						},
					}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:   "managed vpc does not exist",
			input:  &infrav1.VPCSpec{},
			output: &infrav1.VPCSpec{ID: "vpc-new", CidrBlock: "10.1.0.0/16"},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
						},
						{
							Name:   aws.String("tag-key"),
							Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
						},
					},
				})).
					Return(&ec2.DescribeVpcsOutput{}, nil)

				m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).
					Return(&ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							State:     aws.String("available"),
							VpcId:     aws.String("vpc-new"),
							CidrBlock: aws.String("10.1.0.0/16"),
						},
					}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).
					Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)

				m.WaitUntilVpcAvailable(gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []*string{aws.String("vpc-new")},
				})).
					Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
		{
			name:   "managed vpc does not exist and ipv6 enabled",
			input:  &infrav1.VPCSpec{EnableIPv6: true},
			output: &infrav1.VPCSpec{ID: "vpc-new", CidrBlock: "10.1.0.0/16", EnableIPv6: true, Ipv6CidrBlock: aws.String("2001:10:10:10::/56")},
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
						},
						{
							Name:   aws.String("tag-key"),
							Values: aws.StringSlice([]string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"}),
						},
						{
							Name:   aws.String("ipv6-cidr-block-association.state"),
							Values: aws.StringSlice([]string{ec2.VpcCidrBlockStateCodeAssociated}),
						},
					},
				})).
					Return(&ec2.DescribeVpcsOutput{}, nil)

				m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).
					Return(&ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							State:     aws.String("available"),
							VpcId:     aws.String("vpc-new"),
							CidrBlock: aws.String("10.1.0.0/16"),
							Ipv6CidrBlockAssociationSet: []*ec2.VpcIpv6CidrBlockAssociation{
								{
									Ipv6CidrBlock: aws.String("2001:10:10:10::/56"),
									Ipv6CidrBlockState: &ec2.VpcCidrBlockState{
										State: aws.String(ec2.VpcCidrBlockStateCodeAssociated),
									},
								},
							},
						},
					}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).
					Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)

				m.WaitUntilVpcAvailable(gomock.Eq(&ec2.DescribeVpcsInput{
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("ipv6-cidr-block-association.state"),
							Values: aws.StringSlice([]string{ec2.VpcCidrBlockStateCodeAssociated}),
						},
					},
					VpcIds: []*string{aws.String("vpc-new")},
				})).
					Return(nil)

				m.CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
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
						NetworkSpec: infrav1.NetworkSpec{
							VPC: *tc.input,
						},
					},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			if err := s.reconcileVPC(); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			reflect.DeepEqual(tc.input, tc.output)
		})
	}
}
