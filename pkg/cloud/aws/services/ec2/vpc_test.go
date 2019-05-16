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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb/mock_elbiface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestReconcileVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.VPCSpec
		output *v1alpha1.VPCSpec
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name:   "managed vpc exists",
			input:  &v1alpha1.VPCSpec{ID: "vpc-exists"},
			output: &v1alpha1.VPCSpec{ID: "vpc-exists", CidrBlock: "10.0.0.0/8"},
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
			},
		},
		{
			name: "managed vpc does not exist",
			input: &v1alpha1.VPCSpec{
				CidrBlock: "10.1.0.0/16",
			},
			output: &v1alpha1.VPCSpec{ID: "vpc-new", CidrBlock: "10.1.0.0/16"},
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

				m.ModifyVpcAttribute(gomock.AssignableToTypeOf(&ec2.ModifyVpcAttributeInput{})).Return(&ec2.ModifyVpcAttributeOutput{}, nil)

				m.WaitUntilVpcAvailable(gomock.Eq(&ec2.DescribeVpcsInput{
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

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
				},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
					ELB: elbMock,
				},
			})

			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			scope.ClusterConfig = &v1alpha1.AWSClusterProviderSpec{
				NetworkSpec: v1alpha1.NetworkSpec{
					VPC: *tc.input,
				},
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
