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
	"strings"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/converters"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
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
	s.scope.V(2).Info("Reconciling subnets")

	subnets := s.scope.Subnets()
	defer func() {
		s.scope.ClusterConfig.NetworkSpec.Subnets = subnets
	}()

	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets()
	if err != nil {
		return err
	}

	// If the subnets are empty, populate the slice with the default configuration.
	// Adds a single private and public subnet in the first available zone.
	if len(existing) < 2 && len(subnets) < 2 {
		zones, err := s.getAvailableZones()
		if err != nil {
			return err
		}

		if len(subnets.FilterPrivate()) == 0 {
			if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
				return errors.New("expected at least one private subnet available for use, got 0")
			}

			subnets = append(subnets, &v1alpha1.SubnetSpec{
				CidrBlock:        defaultPrivateSubnetCidr,
				AvailabilityZone: zones[0],
				IsPublic:         false,
			})
		}

		if len(subnets.FilterPublic()) == 0 {
			if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
				return errors.New("expected at least one public subnet available for use, got 0")
			}

			subnets = append(subnets, &v1alpha1.SubnetSpec{
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
			if (sn.ID != "" && exsn.ID == sn.ID) || (sn.CidrBlock == exsn.CidrBlock) {
				if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
					// TODO(vincepri): Validate provided subnet passes some basic checks.
					exsn.DeepCopyInto(sn)
					continue LoopExisting
				}

				// Make sure tags are up to date.
				err = tags.Ensure(exsn.Tags, &tags.ApplyParams{
					EC2Client:   s.scope.EC2,
					BuildParams: s.getSubnetTagParams(exsn.ID, exsn.IsPublic),
				})

				if err != nil {
					return errors.Wrapf(err, "failed to ensure tags on subnet %q", exsn.ID)
				}

				// TODO(vincepri): check if subnet needs to be updated.
				exsn.DeepCopyInto(sn)
				continue LoopExisting
			}
		}

		// TODO(vincepri): delete extra subnets that exist and are managed by us.
		subnets = append(subnets, exsn)
	}

	// Proceed to create the rest of the subnets that don't have an ID.
	if !s.scope.VPC().IsUnmanaged(s.scope.Name()) {
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
	}

	s.scope.V(2).Info("Subnets available", "subnets", subnets)
	return nil
}

func (s *Service) deleteSubnets() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping subnets deletion in unmanaged mode")
		return nil
	}

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
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			filter.EC2.SubnetStates(ec2.SubnetStatePending, ec2.SubnetStateAvailable),
		},
	}

	if s.scope.VPC().ID == "" {
		input.Filters = append(input.Filters, filter.EC2.Cluster(s.scope.Name()))
	} else {
		input.Filters = append(input.Filters, filter.EC2.VPC(s.scope.VPC().ID))
	}

	out, err := s.scope.EC2.DescribeSubnets(input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe subnets in vpc %q", s.scope.VPC().ID)
	}

	routeTables, err := s.describeVpcRouteTablesBySubnet()
	if err != nil {
		return nil, err
	}

	natGateways, err := s.describeNatGatewaysBySubnet()
	if err != nil {
		return nil, err
	}

	subnets := make([]*v1alpha1.SubnetSpec, 0, len(out.Subnets))
	// Besides what the AWS API tells us directly about the subnets, we also want to discover whether the subnet is "public" (i.e. directly connected to the internet) and if there are any associated NAT gateways.
	// We also look for a tag indicating that a particular subnet should be public, to try and determine whether a managed VPC's subnet should have such a route, but does not.
	for _, ec2sn := range out.Subnets {
		spec := &v1alpha1.SubnetSpec{
			ID:               *ec2sn.SubnetId,
			CidrBlock:        *ec2sn.CidrBlock,
			AvailabilityZone: *ec2sn.AvailabilityZone,
			Tags:             converters.TagsToMap(ec2sn.Tags),
		}

		// A subnet is public if it's tagged as such...
		if spec.Tags.GetRole() == v1alpha1.PublicRoleTagValue {
			spec.IsPublic = true
		}

		// ... or if it has an internet route
		rt := routeTables[*ec2sn.SubnetId]
		if rt != nil {
			spec.RouteTableID = rt.RouteTableId
			for _, route := range rt.Routes {
				if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw") {
					spec.IsPublic = true
				}
			}
		}

		ngw := natGateways[*ec2sn.SubnetId]
		if ngw != nil {
			spec.NatGatewayID = ngw.NatGatewayId
		}
		subnets = append(subnets, spec)
	}

	return subnets, nil
}

func (s *Service) createSubnet(sn *v1alpha1.SubnetSpec) (*v1alpha1.SubnetSpec, error) {
	out, err := s.scope.EC2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(s.scope.VPC().ID),
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

	applyTagsParams := &tags.ApplyParams{
		EC2Client:   s.scope.EC2,
		BuildParams: s.getSubnetTagParams(*out.Subnet.SubnetId, sn.IsPublic),
	}

	if err := tags.Apply(applyTagsParams); err != nil {
		return nil, errors.Wrapf(err, "failed to tag subnet %q", *out.Subnet.SubnetId)
	}

	if sn.IsPublic {
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

	s.scope.V(2).Info("Created new subnet in VPC with cidr and availability zone ",
		"subnet-id", *out.Subnet.SubnetId,
		"vpc-id", *out.Subnet.VpcId,
		"cidr-block", *out.Subnet.CidrBlock,
		"availability-zone", *out.Subnet.AvailabilityZone)

	record.Eventf(s.scope.Cluster, "CreatedSubnet", "Created new managed Subnet %q", *out.Subnet.SubnetId)

	return &v1alpha1.SubnetSpec{
		ID:               *out.Subnet.SubnetId,
		AvailabilityZone: *out.Subnet.AvailabilityZone,
		CidrBlock:        *out.Subnet.CidrBlock,
		IsPublic:         sn.IsPublic,
	}, nil
}

func (s *Service) deleteSubnet(id string) error {
	_, err := s.scope.EC2.DeleteSubnet(&ec2.DeleteSubnetInput{
		SubnetId: aws.String(id),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %q", id)
	}

	s.scope.V(2).Info("Deleted subnet in vpc", "subnet-id", id, "vpc-id", s.scope.VPC().ID)
	record.Eventf(s.scope.Cluster, "DeletedSubnet", "Deleted managed Subnet %q", id)
	return nil
}

func (s *Service) getSubnetTagParams(id string, public bool) v1alpha1.BuildParams {
	var role string
	if public {
		role = v1alpha1.PublicRoleTagValue
	} else {
		role = v1alpha1.PrivateRoleTagValue
	}

	var name strings.Builder
	name.WriteString(s.scope.Name())
	name.WriteString("-subnet-")
	name.WriteString(role)

	return v1alpha1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   v1alpha1.ResourceLifecycleOwned,
		Name:        aws.String(name.String()),
		Role:        aws.String(role),
	}
}
