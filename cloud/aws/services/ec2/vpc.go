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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

type vpcs interface {
	DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error)
	DeleteVpc(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error)
}

type VPC struct {
	ID        string
	CidrBlock string
}

func (s *Service) ReconcileVPC(v v1alpha1.VPC) (*VPC, error) {
	// Does it exist and look in good working order? ok exit no error
	vpc, err := s.lookupVPCByID(v.ID)
	if err != nil {
		if IsNotFound(err) {
			return s.createVPC(v)
		}

		return nil, err
	}

	// TODO(vincepri): tag vpc with https://docs.aws.amazon.com/sdk-for-go/api/service/resourcegroupstaggingapi/#ResourceGroupsTaggingAPI.TagResources

	return vpc, nil
}

func (s *Service) createVPC(v v1alpha1.VPC) (*VPC, error) {
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(v.CidrBlock),
	}

	out, err := s.VPCs.CreateVpc(input)
	if err != nil {
		return nil, err
	}

	return &VPC{
		ID:        *out.Vpc.VpcId,
		CidrBlock: *out.Vpc.CidrBlock,
	}, nil
}

func (s *Service) deleteVPC(v v1alpha1.VPC) error {
	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(v.ID),
	}

	_, err := s.VPCs.DeleteVpc(input)
	return err
}

func (s *Service) lookupVPCByID(id string) (*VPC, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			&id,
		},
	}

	out, err := s.VPCs.DescribeVpcs(input)
	if err != nil {
		return nil, err
	}

	if len(out.Vpcs) == 0 {
		return nil, NewNotFound(fmt.Errorf("could not find vpc: %q", id))
	} else if len(out.Vpcs) > 1 {
		return nil, NewConflict(fmt.Errorf("found more than one vpc with supplied filters. Please clean up extra VPCs: %s", out.GoString()))
	}

	return &VPC{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
	}, nil
}
