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
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	usageLimit := 3
	selection := infrav1.AZSelectionSchemeOrdered
	tags := []*ec2.Tag{
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
	}

	testCases := []struct {
		name    string
		input   *infrav1.VPCSpec
		want    *infrav1.VPCSpec
		expect  func(m *mocks.MockEC2APIMockRecorder)
		wantErr bool
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
			wantErr: false,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
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
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							State:     aws.String("available"),
							VpcId:     aws.String("vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      tags,
						},
					},
				}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:    "Should create a new VPC if managed vpc does not exist",
			input:   &infrav1.VPCSpec{AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: false,
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
				m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).Return(&ec2.CreateVpcOutput{
					Vpc: &ec2.Vpc{
						State:     aws.String("available"),
						VpcId:     aws.String("vpc-new"),
						CidrBlock: aws.String("10.1.0.0/16"),
						Tags:      tags,
					},
				}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)
			},
		},
		{
			name:    "managed vpc id exists, but vpc resource is missing",
			input:   &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
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
				})).Return(nil, awserr.New("404", "http not found err", errors.New("err")))
			},
		},
		{
			name:  "Should patch vpc spec successfully, if unmanaged vpc exists",
			input: &infrav1.VPCSpec{ID: "unmanaged-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			want: &infrav1.VPCSpec{
				ID:        "unmanaged-vpc-exists",
				CidrBlock: "10.0.0.0/8",
				Tags: map[string]string{
					"sigs.k8s.io/cluster-api-provider-aws/role": "common",
					"Name": "test-cluster-vpc",
				},
				AvailabilityZoneUsageLimit: &usageLimit,
				AvailabilityZoneSelection:  &selection,
			},
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							State:     aws.String("available"),
							VpcId:     aws.String("unmanaged-vpc-exists"),
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
							},
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
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							State:     aws.String("available"),
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      tags,
						},
					},
				}, nil)
				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).Return(nil, awserr.New("InvalidVpcID.NotFound", "not found", nil))
				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeTrue).AnyTimes()
			},
		},
		{
			name:    "Should return error if failed to set vpc attributes for managed vpc",
			input:   &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
					VpcIds: []*string{
						aws.String("managed-vpc-exists"),
					},
					Filters: []*ec2.Filter{
						{
							Name:   aws.String("state"),
							Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
						},
					},
				})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							State:     aws.String("available"),
							VpcId:     aws.String("unmanaged-vpc-exists"),
							CidrBlock: aws.String("10.0.0.0/8"),
							Tags:      tags,
						},
					},
				}, nil)
				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).AnyTimes().Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
		},
		{
			name:    "Should return error if failed to create vpc",
			input:   &infrav1.VPCSpec{AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).Return(nil, awserrors.NewFailedDependency("failed dependency"))
			},
		},
		{
			name:    "Should return error if describe vpc returns empty list",
			input:   &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{},
				}, nil)
			},
		},
		{
			name:    "Should return error if describe vpc returns more than 1 vpcs",
			input:   &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
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
			name:    "Should return error if vpc state is not available/pending",
			input:   &infrav1.VPCSpec{ID: "managed-vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			wantErr: true,
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DescribeVpcs(gomock.AssignableToTypeOf(&ec2.DescribeVpcsInput{})).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{
							VpcId: aws.String("vpc"),
							State: aws.String("deleting"),
						},
					},
				}, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			clusterScope, err := getClusterScope(tc.input)
			g.Expect(err).NotTo(HaveOccurred())
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock.EXPECT())
			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			err = s.reconcileVPC()
			if tc.wantErr {
				g.Expect(err).ToNot(BeNil())
				return
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(tc.want).To(Equal(&clusterScope.AWSCluster.Spec.NetworkSpec.VPC))
		})
	}
}

func Test_DeleteVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tags := map[string]string{
		"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster": "owned",
	}

	testCases := []struct {
		name    string
		input   *infrav1.VPCSpec
		wantErr bool
		expect  func(m *mocks.MockEC2APIMockRecorder)
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
				m.DeleteVpc(gomock.Eq(&ec2.DeleteVpcInput{
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
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteVpc(gomock.Eq(&ec2.DeleteVpcInput{
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
			expect: func(m *mocks.MockEC2APIMockRecorder) {
				m.DeleteVpc(gomock.Eq(&ec2.DeleteVpcInput{
					VpcId: aws.String("managed-vpc"),
				})).Return(nil, awserr.New("InvalidVpcID.NotFound", "not found", nil))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			clusterScope, err := getClusterScope(tc.input)
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

func getClusterScope(vpcSpec *infrav1.VPCSpec) (*scope.ClusterScope, error) {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	awsCluster := &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec: infrav1.AWSClusterSpec{
			NetworkSpec: infrav1.NetworkSpec{
				VPC: *vpcSpec,
			},
		},
	}
	client.Create(context.TODO(), awsCluster)
	return scope.NewClusterScope(scope.ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		},
		AWSCluster: awsCluster,
		Client:     client,
	})
}
