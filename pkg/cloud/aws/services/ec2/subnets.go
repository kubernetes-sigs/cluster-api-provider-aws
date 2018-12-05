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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	defaultPrivateSubnetCidr = "10.0.0.0/24"
	defaultPublicSubnetCidr  = "10.0.1.0/24"
)

func (s *Service) reconcileSubnets() error {
	klog.V(2).Infof("Reconciling subnets")

	subnets := s.scope.Subnets()
	defer func() {
		s.scope.Network().Subnets = subnets
	}()

	// Make sure all subnets have a vpc id.
	for _, sn := range subnets {
		if sn.VpcID == "" {
			sn.VpcID = s.scope.VPC().ID
		}
	}

	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets()
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
				VpcID:            s.scope.VPC().ID,
				CidrBlock:        defaultPrivateSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         false,
			})
		}

		if len(subnets.FilterPublic()) == 0 {
			subnets = append(subnets, &v1alpha1.Subnet{
				VpcID:            s.scope.VPC().ID,
				CidrBlock:        defaultPublicSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         true,
			})
		}
	}

LoopExisting:
	for _, exsn := range existing {
		// Check if the subnet already exists in the state, in that case reconcile it.
		for _, sn := range subnets {
			// Two subnets are defined equal to each other if their id is equal
			// or if they are in the same vpc and the cidr block is the same.
			if (sn.ID != "" && exsn.ID == sn.ID) || (sn.VpcID == exsn.VpcID && sn.CidrBlock == exsn.CidrBlock) {
				// TODO(vincepri): check if subnet needs to be updated.
				exsn.DeepCopyInto(sn)
				continue LoopExisting
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

	klog.V(2).Infof("Subnets available: %v", subnets)
	return nil
}

func (s *Service) deleteSubnets() error {
	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets()
	if err != nil {
		return err
	}

	for _, sn := range existing {
		if err := s.deleteSubnet(sn.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) describeVpcSubnets() (v1alpha1.Subnets, error) {
	out, err := s.scope.EC2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
			filter.EC2.Cluster(s.scope.Name()),
			filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe subnets in vpc %q", s.scope.VPC().ID)
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
	mapPublicIP := sn.IsPublic

	out, err := s.scope.EC2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(sn.VpcID),
		CidrBlock:        aws.String(sn.CidrBlock),
		AvailabilityZone: aws.String(sn.AvailabilityZone),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnet")
	}

	wReq := &ec2.DescribeSubnetsInput{SubnetIds: []*string{out.Subnet.SubnetId}}
	if err := s.scope.EC2.WaitUntilSubnetAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for subnet %q", *out.Subnet.SubnetId)
	}

	suffix := "private"
	role := tags.ValueCommonRole
	if mapPublicIP {
		suffix = "public"
		role = tags.ValueBastionRole
	}
	name := fmt.Sprintf("%s-subnet-%s", s.scope.Name(), suffix)

	applyTagsParams := &tags.ApplyParams{
		EC2Client: s.scope.EC2,
		BuildParams: tags.BuildParams{
			ClusterName: s.scope.Name(),
			ResourceID:  *out.Subnet.SubnetId,
			Lifecycle:   tags.ResourceLifecycleOwned,
			Name:        aws.String(name),
			Role:        aws.String(role),
		},
	}

	if err := tags.Apply(applyTagsParams); err != nil {
		return nil, errors.Wrapf(err, "failed to tag subnet %q", *out.Subnet.SubnetId)
	}

	if mapPublicIP {
		attReq := &ec2.ModifySubnetAttributeInput{
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
			SubnetId: out.Subnet.SubnetId,
		}

		if _, err := s.scope.EC2.ModifySubnetAttribute(attReq); err != nil {
			return nil, errors.Wrapf(err, "failed to set subnet %q attributes", *out.Subnet.SubnetId)
		}
	}

	klog.V(2).Infof("Created new subnet %q in VPC %q with cidr %q and availability zone %q",
		*out.Subnet.SubnetId, *out.Subnet.VpcId, *out.Subnet.CidrBlock, *out.Subnet.AvailabilityZone)

	record.Eventf(s.scope.Cluster, "CreatedSubnet", "Created new managed Subnet %q", *out.Subnet.SubnetId)

	return &v1alpha1.Subnet{
		ID:               *out.Subnet.SubnetId,
		VpcID:            *out.Subnet.VpcId,
		AvailabilityZone: *out.Subnet.AvailabilityZone,
		CidrBlock:        *out.Subnet.CidrBlock,
		IsPublic:         mapPublicIP,
	}, nil
}

func (s *Service) deleteSubnet(id string) error {
	_, err := s.scope.EC2.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %q", id)
	}

	klog.V(2).Infof("Deleted subnet %q in vpc %q", id, s.scope.VPC().ID)
	record.Eventf(s.scope.Cluster, "DeletedSubnet", "Deleted managed Subnet %q", id)
	return nil
}
