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
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/providerconfig/v1alpha1"
)

const (
	defaultPrivateSubnetCidr = "10.0.0.0/24"
	defaultPublicSubnetCidr  = "10.0.1.0/24"
)

func (s *Service) reconcileSubnets(clusterName string, network *v1alpha1.Network) error {
	glog.V(2).Infof("Reconciling subnets")

	// Make sure all subnets have a vpc id.
	for _, sn := range network.Subnets {
		if sn.VpcID == "" {
			sn.VpcID = network.VPC.ID
		}
	}

	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets(clusterName, network.VPC.ID)
	if err != nil {
		return err
	}

	// If the subnets are empty, populate the slice with the default configuration.
	// Adds a single private and public subnet in the first available zone.
	if len(network.Subnets) < 2 {
		zones, err := s.getAvailableZones()
		if err != nil {
			return err
		}

		if len(network.Subnets.FilterPrivate()) == 0 {
			network.Subnets = append(network.Subnets, &v1alpha1.Subnet{
				VpcID:            network.VPC.ID,
				CidrBlock:        defaultPrivateSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         false,
			})
		}

		if len(network.Subnets.FilterPublic()) == 0 {
			network.Subnets = append(network.Subnets, &v1alpha1.Subnet{
				VpcID:            network.VPC.ID,
				CidrBlock:        defaultPublicSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         true,
			})
		}

	}

LoopExisting:
	for _, exsn := range existing {
		// Check if the subnet already exists in the state, in that case reconcile it.
		for _, sn := range network.Subnets {
			// Two subnets are defined equal to each other if their id is equal
			// or if they are in the same vpc and the cidr block is the same.
			if (sn.ID != "" && exsn.ID == sn.ID) || (sn.VpcID == exsn.VpcID && sn.CidrBlock == exsn.CidrBlock) {
				// TODO(vincepri): check if subnet needs to be updated.
				exsn.DeepCopyInto(sn)
				continue LoopExisting
			}
		}

		// TODO(vincepri): delete extra subnets that exist and are managed by us.
		network.Subnets = append(network.Subnets, exsn)
	}

	// Proceed to create the rest of the subnets that don't have an ID.
	for _, subnet := range network.Subnets {
		if subnet.ID != "" {
			continue
		}

		nsn, err := s.createSubnet(clusterName, subnet)
		if err != nil {
			return err
		}

		nsn.DeepCopyInto(subnet)
	}

	glog.V(2).Infof("Subnets available: %v", network.Subnets)
	return nil
}

func (s *Service) deleteSubnets(clusterName string, network *v1alpha1.Network) error {
	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets(clusterName, network.VPC.ID)
	if err != nil {
		return err
	}

	for _, sn := range existing {
		input := &ec2.DeleteSubnetInput{
			SubnetId: aws.String(sn.ID),
		}

		if _, err := s.EC2.DeleteSubnet(input); err != nil {
			return errors.Wrapf(err, "failed to delete subnet %q", sn.ID)
		}

		glog.Infof("deleted subnet %q in VPC %q", sn.ID, network.VPC.ID)
	}
	return nil
}

func (s *Service) describeVpcSubnets(clusterName string, vpcID string) (v1alpha1.Subnets, error) {
	out, err := s.EC2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			s.filterVpc(vpcID),
			s.filterCluster(clusterName),
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

func (s *Service) createSubnet(clusterName string, sn *v1alpha1.Subnet) (*v1alpha1.Subnet, error) {
	mapPublicIP := sn.IsPublic

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

	suffix := "private"
	role := TagValueCommonRole
	if mapPublicIP {
		suffix = "public"
		role = TagValueBastionRole
	}
	name := fmt.Sprintf("%s-subnet-%s", clusterName, suffix)
	if err := s.createTags(clusterName, *out.Subnet.SubnetId, ResourceLifecycleOwned, name, role, nil); err != nil {
		return nil, errors.Wrapf(err, "failed to tag subnet %q", *out.Subnet.SubnetId)
	}

	if mapPublicIP {
		attReq := &ec2.ModifySubnetAttributeInput{
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
			SubnetId: out.Subnet.SubnetId,
		}

		if _, err := s.EC2.ModifySubnetAttribute(attReq); err != nil {
			return nil, errors.Wrapf(err, "failed to set subnet %q attributes", *out.Subnet.SubnetId)
		}
	}

	glog.V(2).Infof("Created new subnet %q in VPC %q with cidr %q and availability zone %q",
		*out.Subnet.SubnetId, *out.Subnet.VpcId, *out.Subnet.CidrBlock, *out.Subnet.AvailabilityZone)

	return &v1alpha1.Subnet{
		ID:               *out.Subnet.SubnetId,
		VpcID:            *out.Subnet.VpcId,
		AvailabilityZone: *out.Subnet.AvailabilityZone,
		CidrBlock:        *out.Subnet.CidrBlock,
		IsPublic:         mapPublicIP,
	}, nil
}

func (s *Service) deleteSubnet(sn *v1alpha1.Subnet) error {
	_, err := s.EC2.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(sn.ID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %q", sn.ID)
	}

	glog.V(2).Infof("Deleted subnet %q", sn.ID)
	return nil
}
