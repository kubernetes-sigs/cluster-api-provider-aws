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

package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

func TestReconcileVPC(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		vpc    *v1alpha1.VPC
		expect func(m *mock_ec2iface.MockEC2API)
	}{
		{
			name: "vpc exists",
			vpc:  &v1alpha1.VPC{ID: "vpc-exists"},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
						VpcIds: []*string{
							aws.String("vpc-exists"),
						},
					})).
					Return(&ec2.DescribeVpcsOutput{
						Vpcs: []*ec2.Vpc{
							&ec2.Vpc{
								VpcId:     aws.String("vpc-exists"),
								CidrBlock: aws.String("10.0.0.0/16"),
							},
						},
					}, nil)
			},
		},
		{
			name: "vpc does not exist",
			vpc:  &v1alpha1.VPC{ID: "vpc-new"},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeVpcs(gomock.Eq(&ec2.DescribeVpcsInput{
						VpcIds: []*string{
							aws.String("vpc-new"),
						},
					})).
					Return(&ec2.DescribeVpcsOutput{}, nil)

				m.EXPECT().
					CreateVpc(gomock.AssignableToTypeOf(&ec2.CreateVpcInput{})).
					Return(&ec2.CreateVpcOutput{
						Vpc: &ec2.Vpc{
							VpcId:     aws.String("vpc-new"),
							CidrBlock: aws.String("10.0.0.0/16"),
						},
					}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)

			s := ec2svc.Service{
				VPCs: ec2Mock,
			}

			vpc, err := s.ReconcileVPC(tc.vpc)
			if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			if vpc.ID != tc.vpc.ID {
				t.Fatalf("Expected an id of %v but found %v", tc.vpc.ID, vpc.ID)
			}
		})
	}
}
