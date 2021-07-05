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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/diff"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

	testCases := []struct {
		name        string
		input       *infrav1.VPCSpec
		expected    *infrav1.VPCSpec
		expect      func(m *mock_ec2iface.MockEC2APIMockRecorder)
		expectError bool
	}{
		{
			name:  "if unmanaged vpc exists, updates tags with aws VPC resource tags",
			input: &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			expected: &infrav1.VPCSpec{
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
			expectError: false,
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
			name:        "if managed vpc does not exist, creates a new VPC",
			input:       &infrav1.VPCSpec{AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			expectError: false,
			expected: &infrav1.VPCSpec{
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
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).
					Return(&ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							State:     aws.String("available"),
							VpcId:     aws.String("vpc-new"),
							CidrBlock: aws.String("10.1.0.0/16"),
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
					}, nil)

				m.DescribeVpcAttribute(gomock.AssignableToTypeOf(&ec2.DescribeVpcAttributeInput{})).
					DoAndReturn(describeVpcAttributeFalse).MinTimes(1)

				m.ModifyVpcAttribute(gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).
					Return(&ec2.ModifyVpcAttributeOutput{}, nil).Times(2)
			},
		},
		{
			name:        "managed vpc id exists, but vpc resource is missing",
			input:       &infrav1.VPCSpec{ID: "vpc-exists", AvailabilityZoneUsageLimit: &usageLimit, AvailabilityZoneSelection: &selection},
			expectError: true,
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
					Return(nil, awserr.New("404", "http not found err", errors.New("err")))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			awsCluster := &infrav1.AWSCluster{
				Spec: infrav1.AWSClusterSpec{
					NetworkSpec: infrav1.NetworkSpec{
						VPC: *tc.input,
					},
				},
			}
			client := fake.NewFakeClientWithScheme(scheme)
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
			g := NewWithT(t)

			err = s.reconcileVPC()
			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			} else {
				g.Expect(err).To(BeNil())
			}

			if !reflect.DeepEqual(tc.expected, &clusterScope.AWSCluster.Spec.NetworkSpec.VPC) {
				t.Errorf("Actual/expected mismatch: %s", diff.ObjectDiff(tc.expected, clusterScope.AWSCluster.Spec.NetworkSpec.VPC))
			}
		})
	}
}
