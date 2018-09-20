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

const (
	defaultPrivateSubnetCidr = "10.0.0.0/24"
	defaultPublicSubnetCidr  = "10.0.1.0/24"
)

func (s *Service) reconcileSubnets(subnets v1alpha1.Subnets, vpc *v1alpha1.VPC) error {
	// Make sure all subnets have a vpc id.
	for _, sn := range subnets {
		if sn.VpcID == "" {
			sn.VpcID = vpc.ID
		}
	}

	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets(vpc.ID)
	if err != nil {
		return err
	}

	// If the subnets are empty, populate the slice with the default configuration.
	// Adds a single private and public subnet in the first available zone.
	if len(subnets) < 2 {
		zones, err := s.getAvailableZones()
		if err != nil {
			return err
		}

		if len(subnets.FilterPrivate()) == 0 {
			subnets = append(subnets, &v1alpha1.Subnet{
				VpcID:            vpc.ID,
				CidrBlock:        defaultPrivateSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         false,
			})
		}

		if len(subnets.FilterPublic()) == 0 {
			subnets = append(subnets, &v1alpha1.Subnet{
				VpcID:            vpc.ID,
				CidrBlock:        defaultPublicSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         true,
			})
		}

	}

	for _, exsn := range existing {
		// Check if the subnet already exists in the state, in that case reconcile it.
		for _, sn := range subnets {
			// Two subnets are defined equal to each other if their id is equal
			// or if they are in the same vpc and the cidr block is the same.
			if (sn.ID != "" && exsn.ID == sn.ID) || (sn.VpcID == exsn.VpcID && sn.CidrBlock == exsn.CidrBlock) {
				// TODO(vincepri): check if subnet needs to be updated.
				exsn.DeepCopyInto(sn)
				break
			}
		}

		// TODO(vincepri): delete extra subnets that exist and are managed by us.
		subnets = append(subnets, exsn)
	}

	// Proceed to create the rest of the subnets that don't have an ID.
	for _, subnet := range subnets {
		if subnet.ID != "" {
			continue
		}

		nsn, err := s.createSubnet(subnet)
		if err != nil {
			return err
		}

		nsn.DeepCopyInto(subnet)
	}

	return nil
}

func (s *Service) describeVpcSubnets(vpcID string) (v1alpha1.Subnets, error) {
	out, err := s.EC2.DescribeSubnets(&ec2.DescribeSubnetsInput{
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

	subnets := make([]*v1alpha1.Subnet, 0, len(out.Subnets))
	for _, ec2sn := range out.Subnets {
		subnets = append(subnets, &v1alpha1.Subnet{
			ID:               *ec2sn.SubnetId,
			VpcID:            *ec2sn.VpcId,
			CidrBlock:        *ec2sn.CidrBlock,
			AvailabilityZone: *ec2sn.AvailabilityZone,
			IsPublic:         *ec2sn.MapPublicIpOnLaunch,
		})
	}

	return subnets, nil
}

func (s *Service) createSubnet(sn *v1alpha1.Subnet) (*v1alpha1.Subnet, error) {
	out, err := s.EC2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(sn.VpcID),
		CidrBlock:        aws.String(sn.CidrBlock),
		AvailabilityZone: aws.String(sn.AvailabilityZone),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}

	wReq := &ec2.DescribeSubnetsInput{SubnetIds: []*string{out.Subnet.SubnetId}}
	if err := s.EC2.WaitUntilSubnetAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for subnet %q", *out.Subnet.SubnetId)
	}

	if sn.IsPublic {
		attReq := &ec2.ModifySubnetAttributeInput{
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
		}

		if _, err := s.EC2.ModifySubnetAttribute(attReq); err != nil {
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
	_, err := s.EC2.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(sn.ID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %q", sn.ID)
	}

	return nil
}
