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
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type vpcs interface {
	DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
}

type VPC struct {
	ID string
}

func (s *Service) ReconcileVPC(id string) (*VPC, error) {
	// Does it exist and look in good working order? ok exit no error
	vpc, err := s.lookupVPCByID(id)
	if err != nil {
		return nil, err
	}
	// If it doesn't exist, create it
	return vpc, nil
}

func (s *Service) lookupVPCByID(id string) (*VPC, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			&id,
		},
	}
	out, err := s.VPCs.DescribeVpcs(input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe vpcs: %v", err)
	}
	if len(out.Vpcs) > 1 {
		return nil, errors.New(fmt.Sprintf("found more than one vpc. Please clean up extra VPCs: %s", out.GoString()))
	}
	return &VPC{
		ID: *out.Vpcs[0].VpcId,
	}, nil
}
