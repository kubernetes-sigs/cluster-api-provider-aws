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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
)

type vpcs struct{}

func (v *vpcs) CreateVpc(i *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	return &ec2.CreateVpcOutput{
		Vpc: &ec2.Vpc{
			VpcId:     aws.String("newVpc"),
			CidrBlock: aws.String("10.0.0.0/16"),
		},
	}, nil
}

func (v *vpcs) DeleteVpc(i *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	return &ec2.DeleteVpcOutput{}, nil
}

func (v *vpcs) DescribeVpcs(i *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	if len(i.VpcIds) == 0 {
		return nil, errors.New("invalid request")
	}

	if *i.VpcIds[0] == "oldVpc" {
		return &ec2.DescribeVpcsOutput{
			Vpcs: []*ec2.Vpc{
				&ec2.Vpc{
					VpcId:     aws.String("oldVpc"),
					CidrBlock: aws.String("10.0.0.0/16"),
				},
			},
		}, nil
	}

	return &ec2.DescribeVpcsOutput{}, nil
}

func TestReconcileVPC(t *testing.T) {
	testCases := []struct {
		name string
		id   string
	}{
		{
			name: "vpc exists",
			id:   "oldVpc",
		},
		{
			name: "vpc does not exist",
			id:   "newVpc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := ec2svc.Service{
				Instances: nil,
				VPCs:      &vpcs{},
			}
			vpc, err := s.ReconcileVPC(v1alpha1.VPC{ID: tc.id})
			if err != nil {
				t.Fatalf("got an unexpected error: %v", err)
			}

			if vpc.ID != tc.id {
				t.Fatalf("Expected an id of %v but found %v", tc.id, vpc.ID)
			}
		})
	}
}
