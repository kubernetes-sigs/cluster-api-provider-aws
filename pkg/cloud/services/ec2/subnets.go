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
	"math/rand"
	"strings"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/internal/cidr"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	internalLoadBalancerTag = "kubernetes.io/role/internal-elb"
	externalLoadBalancerTag = "kubernetes.io/role/elb"
)

func (s *Service) reconcileSubnets() error {
	s.scope.V(2).Info("Reconciling subnets")

	subnets := s.scope.Subnets()
	defer func() {
		s.scope.AWSCluster.Spec.NetworkSpec.Subnets = subnets
	}()

	// Describe subnets in the vpc.
	existing, err := s.describeVpcSubnets()
	if err != nil {
		return err
	}

	// If we have an unmanaged VPC then we expect either of the following to be true:
	// 1) No subnets specified and no existing subnets in the VPC
	// 2) Subnets specified and they must already exist
	// Otherwise this is an error
	subnetErrs := []string{}
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		if len(subnets) == 0 {
			if len(existing) > 0 {
				return errors.New("no subnets specified but subnets already exist in the VPC")
			}
		} else {
			for _, subnet := range subnets {
				existingSubnet := existing.FindByID(subnet.ID)
				if existingSubnet == nil {
					subnetErrs = append(subnetErrs, fmt.Sprintf("subnet %s specified but it doesn't exist in vpc %s", subnet.ID, s.scope.VPC().ID))
				}
			}
		}
	}

	// If there are subnets specified we need at least 1 private and 1 public subnet
	if len(subnets) > 0 {
		privateSubnets := subnets.FilterPrivate()
		if len(privateSubnets) < 1 {
			subnetErrs = append(subnetErrs, "expected at least 1 private subnet but got 0")
		}
		publicSubnets := subnets.FilterPublic()
		if len(publicSubnets) < 1 {
			subnetErrs = append(subnetErrs, "expected at least 1 public subnet but got 0")
		}
	}

	if len(subnetErrs) > 0 {
		return errors.Errorf("subnet errors encountered: %s", subnetErrs)
	}

	if len(subnets) == 0 {
		// If we have no subnets then create subnets. There will be 1 public and 1 private subnet
		// for each az in a region up to a maximum of 3 azs
		subnets, err = s.getDefaultSubnets()
		if err != nil {
			return errors.Wrap(err, "failed getting default subnets")
		}
	}

LoopExisting:
	for i := range existing {
		exsn := existing[i]
		// Check if the subnet already exists in the state, in that case reconcile it.
		for j := range subnets {
			sn := subnets[j]
			// Two subnets are defined equal to each other if their id is equal
			// or if they are in the same vpc and the cidr block is the same.
			if (sn.ID != "" && exsn.ID == sn.ID) || (sn.CidrBlock == exsn.CidrBlock) {
				if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
					// TODO(vincepri): Validate provided subnet passes some basic checks.
					exsn.DeepCopyInto(sn)
					continue LoopExisting
				}

				// Make sure tags are up to date.
				if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
					if err := tags.Ensure(exsn.Tags, &tags.ApplyParams{
						EC2Client:   s.scope.EC2,
						BuildParams: s.getSubnetTagParams(exsn.ID, exsn.IsPublic, exsn.AvailabilityZone, sn.Tags),
					}); err != nil {
						return false, err
					}
					return true, nil
				}, awserrors.SubnetNotFound); err != nil {
					record.Warnf(s.scope.AWSCluster, "FailedTagSubnet", "Failed tagging managed Subnet %q: %v", exsn.ID, err)
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
			subnet.ID = nsn.ID
		}
	}

	s.scope.V(2).Info("Subnets available", "subnets", subnets)
	return nil
}

func (s *Service) getDefaultSubnets() (infrav1.Subnets, error) {
	zones, err := s.getAvailableZones()
	if err != nil {
		return nil, err
	}

	//TODO: should this be configurable or are we ok with 3 AZs max for the default?
	if len(zones) > 3 {
		s.scope.V(2).Info("region has more than 3 availability zones, picking 3 at random", "region", s.scope.Region())
		rand.Shuffle(len(zones), func(i, j int) {
			zones[i], zones[j] = zones[j], zones[i]
		})
		zones = append(zones, zones[0], zones[1], zones[2])
	}

	// 1 private subnet for each AZ plus 1 other subnet that will be further sub-divided for the public subnets
	numSubnets := len(zones) + 1
	subnetCIDRs, err := cidr.SplitIntoSubnetsIPv4(s.scope.VPC().CidrBlock, numSubnets)
	if err != nil {
		return nil, errors.Wrapf(err, "failed splitting VPC CIDR %s into subnets", s.scope.VPC().CidrBlock)
	}

	publicSubnetCIDRs, err := cidr.SplitIntoSubnetsIPv4(subnetCIDRs[0].String(), len(zones))
	if err != nil {
		return nil, errors.Wrapf(err, "failed splitting CIDR %s into public subnets", subnetCIDRs[0].String())
	}
	privateSubnetCIDRs := append(subnetCIDRs[:0], subnetCIDRs[1:]...)

	subnets := infrav1.Subnets{}
	for i, zone := range zones {
		subnets = append(subnets, &infrav1.SubnetSpec{
			CidrBlock:        publicSubnetCIDRs[i].String(),
			AvailabilityZone: zone,
			IsPublic:         true,
		})
		subnets = append(subnets, &infrav1.SubnetSpec{
			CidrBlock:        privateSubnetCIDRs[i].String(),
			AvailabilityZone: zone,
			IsPublic:         false,
		})
	}

	return subnets, nil
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

func (s *Service) describeVpcSubnets() (infrav1.Subnets, error) {
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

	subnets := make([]*infrav1.SubnetSpec, 0, len(out.Subnets))
	// Besides what the AWS API tells us directly about the subnets, we also want to discover whether the subnet is "public" (i.e. directly connected to the internet) and if there are any associated NAT gateways.
	// We also look for a tag indicating that a particular subnet should be public, to try and determine whether a managed VPC's subnet should have such a route, but does not.
	for _, ec2sn := range out.Subnets {
		spec := &infrav1.SubnetSpec{
			ID:               *ec2sn.SubnetId,
			CidrBlock:        *ec2sn.CidrBlock,
			AvailabilityZone: *ec2sn.AvailabilityZone,
			Tags:             converters.TagsToMap(ec2sn.Tags),
		}

		// A subnet is public if it's tagged as such...
		if spec.Tags.GetRole() == infrav1.PublicRoleTagValue {
			spec.IsPublic = true
		}

		// ... or if it has an internet route
		rt := routeTables[*ec2sn.SubnetId]
		if rt == nil {
			// If there is no explicit association, subnet defaults to main route table as implicit association
			rt = routeTables[mainRouteTableInVPCKey]
		}
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

func (s *Service) createSubnet(sn *infrav1.SubnetSpec) (*infrav1.SubnetSpec, error) {
	out, err := s.scope.EC2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(s.scope.VPC().ID),
		CidrBlock:        aws.String(sn.CidrBlock),
		AvailabilityZone: aws.String(sn.AvailabilityZone),
	})

	if err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedCreateSubnet", "Failed creating new managed Subnet %v", err)
		return nil, errors.Wrap(err, "failed to create subnet")
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulCreateSubnet", "Created new managed Subnet %q", *out.Subnet.SubnetId)

	wReq := &ec2.DescribeSubnetsInput{SubnetIds: []*string{out.Subnet.SubnetId}}
	if err := s.scope.EC2.WaitUntilSubnetAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for subnet %q", *out.Subnet.SubnetId)
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := tags.Apply(&tags.ApplyParams{
			EC2Client:   s.scope.EC2,
			BuildParams: s.getSubnetTagParams(*out.Subnet.SubnetId, sn.IsPublic, sn.AvailabilityZone, sn.Tags),
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.SubnetNotFound); err != nil {
		record.Warnf(s.scope.AWSCluster, "FailedTagSubnet", "Failed tagging managed Subnet %q: %v", *out.Subnet.SubnetId, err)
		return nil, errors.Wrapf(err, "failed to tag subnet %q", *out.Subnet.SubnetId)
	}

	record.Eventf(s.scope.AWSCluster, "SuccessfulTagSubnet", "Tagged managed Subnet %q", *out.Subnet.SubnetId)

	if sn.IsPublic {
		attReq := &ec2.ModifySubnetAttributeInput{
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
			SubnetId: out.Subnet.SubnetId,
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.scope.EC2.ModifySubnetAttribute(attReq); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.SubnetNotFound); err != nil {
			record.Warnf(s.scope.AWSCluster, "FailedModifySubnetAttributes", "Failed modifying managed Subnet %q attributes: %v", *out.Subnet.SubnetId, err)
			return nil, errors.Wrapf(err, "failed to set subnet %q attributes", *out.Subnet.SubnetId)
		}
		record.Eventf(s.scope.AWSCluster, "SuccessfulModifySubnetAttributes", "Modified managed Subnet %q attributes", *out.Subnet.SubnetId)
	}

	s.scope.V(2).Info("Created new subnet in VPC with cidr and availability zone ",
		"subnet-id", *out.Subnet.SubnetId,
		"vpc-id", *out.Subnet.VpcId,
		"cidr-block", *out.Subnet.CidrBlock,
		"availability-zone", *out.Subnet.AvailabilityZone)

	return &infrav1.SubnetSpec{
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
		record.Warnf(s.scope.AWSCluster, "FailedDeleteSubnet", "Failed to delete managed Subnet %q: %v", id, err)
		return errors.Wrapf(err, "failed to delete subnet %q", id)
	}

	s.scope.V(2).Info("Deleted subnet in vpc", "subnet-id", id, "vpc-id", s.scope.VPC().ID)
	record.Eventf(s.scope.AWSCluster, "SuccessfulDeleteSubnet", "Deleted managed Subnet %q", id)
	return nil
}

func (s *Service) getSubnetTagParams(id string, public bool, zone string, manualTags infrav1.Tags) infrav1.BuildParams {
	var role string
	additionalTags := s.scope.AdditionalTags()

	if public {
		role = infrav1.PublicRoleTagValue
		additionalTags[externalLoadBalancerTag] = "1"
	} else {
		role = infrav1.PrivateRoleTagValue
		additionalTags[internalLoadBalancerTag] = "1"
	}

	// Add tag needed for Service type=LoadBalancer
	additionalTags[infrav1.NameKubernetesAWSCloudProviderPrefix+s.scope.Name()] = string(infrav1.ResourceLifecycleShared)

	for k, v := range manualTags {
		additionalTags[k] = v
	}

	var name strings.Builder
	name.WriteString(s.scope.Name())
	name.WriteString("-subnet-")
	name.WriteString(role)
	name.WriteString("-")
	name.WriteString(strings.Replace(zone, "-", "", -1))

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name.String()),
		Role:        aws.String(role),
		Additional:  additionalTags,
	}
}
