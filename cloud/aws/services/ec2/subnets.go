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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

func (s *Service) reconcileSubnets(subnets []*v1alpha1.Subnet, vpc *v1alpha1.VPC) error {
	existing, err := s.describeVpcSubnets(vpc.ID)
	if err != nil {
		return err
	}

	for _, sn := range subnets {
		if sn.VpcID == "" {
			sn.VpcID = vpc.ID
		}

		if xsn, ok := existing[sn.ID]; ok {
			// TODO(vincepri): check tags, check attr match.
			xsn.DeepCopyInto(sn)
			delete(existing, sn.ID)
			continue
		}

		nsn, err := s.createSubnet(sn)
		if err != nil {
			return err
		}

		nsn.DeepCopyInto(sn)
		delete(existing, sn.ID)
	}

	// TODO(vincepri): at this point `existing` contains all the other subnets that aren't
	// in our state and that we might have to delete if they are managed by us (e.g. they contain a specific tag).

	return nil
}

func (s *Service) describeVpcSubnets(vpcID string) (map[string]*v1alpha1.Subnet, error) {
	out, err := s.ec2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe subnets in vpc %q", vpcID)
	}

	subnets := make(map[string]*v1alpha1.Subnet, len(out.Subnets))
	for _, ec2sn := range out.Subnets {
		subnets[*ec2sn.SubnetId] = &v1alpha1.Subnet{
			ID:               *ec2sn.SubnetId,
			VpcID:            *ec2sn.VpcId,
			CidrBlock:        *ec2sn.CidrBlock,
			AvailabilityZone: *ec2sn.AvailabilityZone,
			IsPublic:         *ec2sn.MapPublicIpOnLaunch,
		}
	}

	return subnets, nil
}

func (s *Service) createSubnet(sn *v1alpha1.Subnet) (*v1alpha1.Subnet, error) {
	out, err := s.ec2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(sn.VpcID),
		CidrBlock:        aws.String(sn.CidrBlock),
		AvailabilityZone: aws.String(sn.AvailabilityZone),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}

	wReq := &ec2.DescribeSubnetsInput{SubnetIds: []*string{out.Subnet.SubnetId}}
	if err := s.ec2.WaitUntilSubnetAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for subnet %q", *out.Subnet.SubnetId)
	}

	if sn.IsPublic {
		attReq := &ec2.ModifySubnetAttributeInput{
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
		}

		if _, err := s.ec2.ModifySubnetAttribute(attReq); err != nil {
			return nil, errors.Wrapf(err, "failed to set subnet %q attributes", *out.Subnet.SubnetId)
		}
	}

	return &v1alpha1.Subnet{
		ID:               *out.Subnet.SubnetId,
		VpcID:            *out.Subnet.VpcId,
		AvailabilityZone: *out.Subnet.AvailabilityZone,
		CidrBlock:        *out.Subnet.CidrBlock,
		IsPublic:         *out.Subnet.MapPublicIpOnLaunch,
	}, nil
}

func (s *Service) deleteSubnet(sn *v1alpha1.Subnet) error {
	_, err := s.ec2.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(sn.ID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %q", sn.ID)
	}

	return nil
}
