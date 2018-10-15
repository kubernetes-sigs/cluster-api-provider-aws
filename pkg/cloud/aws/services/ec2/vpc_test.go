// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ec2

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
)

func TestReconcileVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		input  *v1alpha1.VPC
		output *v1alpha1.VPC
		expect func(m *mock_ec2iface.MockEC2API)
	}{
		{
			name:   "vpc exists",
			input:  &v1alpha1.VPC{ID: "vpc-exists"},
			output: &v1alpha1.VPC{ID: "vpc-exists", CidrBlock: "10.0.0.0/8"},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
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
							},
						},
					}, nil)
			},
		},
		{
			name:   "vpc does not exist",
			input:  &v1alpha1.VPC{ID: "vpc-new", CidrBlock: "10.1.0.0/16"},
			output: &v1alpha1.VPC{ID: "vpc-new", CidrBlock: "10.1.0.0/16"},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
						VpcIds: []*string{
							aws.String("vpc-new"),
						},
						Filters: []*ec2.Filter{
							{
								Name:   aws.String("state"),
								Values: aws.StringSlice([]string{ec2.VpcStatePending, ec2.VpcStateAvailable}),
							},
						},
					})).
					Return(&ec2.DescribeVpcsOutput{}, nil)

				m.EXPECT().
					CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).
					Return(&ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							State:     aws.String("available"),
							VpcId:     aws.String("vpc-new"),
							CidrBlock: aws.String("10.1.0.0/16"),
						},
					}, nil)

				m.EXPECT().
					WaitUntilVpcAvailable(gomock.Eq(&ec2.DescribeVpcsInput{
						VpcIds: []*string{aws.String("vpc-new")},
					})).
					Return(nil)

				m.EXPECT().
					CreateTags(gomock.AssignableToTypeOf(&ec2.CreateTagsInput{})).
					Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)

			s := NewService(ec2Mock)
			if err := s.reconcileVPC("test-cluster", tc.input); err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			reflect.DeepEqual(tc.input, tc.output)
		})
	}
}
