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
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

const (
	anyIPv4CidrBlock = "0.0.0.0/0"
)

func (s *Service) reconcileRouteTables(clusterName string, in *v1alpha1.Network) error {
	klog.V(2).Infof("Reconciling routing tables")

	subnetRouteMap, err := s.describeVpcRouteTablesBySubnet(clusterName, in.VPC.ID)
	if err != nil {
		return err
	}

	for _, sn := range in.Subnets {
		if igw, ok := subnetRouteMap[sn.ID]; ok {
			klog.V(2).Infof("Subnet %q is already associated with route table %q", sn.ID, *igw.RouteTableId)
			// TODO(vincepri): if the route table ids are both non-empty and they don't match, replace the association.
			// TODO(vincepri): check that everything is in order, e.g. routes match the subnet type.
			continue
		}

		// For each subnet that doesn't have a routing table associated with it,
		// create a new table with the appropriate default routes and associate it to the subnet.
		var routes []*ec2.Route
		if sn.IsPublic {
			if in.InternetGatewayID == nil {
				return errors.Errorf("failed to create routing tables: internet gateway for %q is nil", in.VPC.ID)
			}

			routes = s.getDefaultPublicRoutes(*in.InternetGatewayID)
		} else {
			natGatewayID, err := s.getNatGatewayForSubnet(in.Subnets, sn)
			if err != nil {
				return err
			}

			routes = s.getDefaultPrivateRoutes(natGatewayID)
		}

		rt, err := s.createRouteTableWithRoutes(clusterName, &in.VPC, routes, sn.IsPublic)
		if err != nil {
			return err
		}

		if err := s.associateRouteTable(rt, sn.ID); err != nil {
			return err
		}

		klog.V(2).Infof("Subnet %q has been associated with route table %q", sn.ID, rt.ID)
		sn.RouteTableID = aws.String(rt.ID)
	}

	return nil
}

func (s *Service) describeVpcRouteTablesBySubnet(clusterName string, vpcID string) (map[string]*ec2.RouteTable, error) {
	rts, err := s.describeVpcRouteTables(clusterName, vpcID)
	if err != nil {
		return nil, err
	}

	// Amazon allows a subnet to be associated only with a single routing table
	// https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html.
	res := make(map[string]*ec2.RouteTable)
	for _, rt := range rts {
		for _, as := range rt.Associations {
			if as.SubnetId == nil {
				continue
			}

			res[*as.SubnetId] = rt
		}
	}

	return res, nil
}

func (s *Service) deleteRouteTables(clusterName string, in *v1alpha1.Network) error {
	rts, err := s.describeVpcRouteTables(clusterName, in.VPC.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to describe route tables in vpc %q", in.VPC.ID)
	}

	for _, rt := range rts {
		for _, as := range rt.Associations {
			if as.SubnetId == nil {
				continue
			}

			if _, err := s.EC2.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{AssociationId: as.RouteTableAssociationId}); err != nil {
				return errors.Wrapf(err, "failed to disassociate route table %q from subnet %q", *rt.RouteTableId, *as.SubnetId)
			}

			klog.Infof("Deleted association between route table %q and subnet %q", *rt.RouteTableId, *as.SubnetId)
		}

		if _, err := s.EC2.DeleteRouteTable(&ec2.DeleteRouteTableInput{RouteTableId: rt.RouteTableId}); err != nil {
			return errors.Wrapf(err, "failed to delete route table %q", *rt.RouteTableId)
		}

		klog.Infof("Deleted route table %q", *rt.RouteTableId)
	}
	return nil
}

func (s *Service) describeVpcRouteTables(clusterName string, vpcID string) ([]*ec2.RouteTable, error) {
	out, err := s.EC2.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			s.filterVpc(vpcID),
			s.filterCluster(clusterName),
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe route tables in vpc %q", vpcID)
	}

	return out.RouteTables, nil
}

// TODO: dedup some of the public/private logic shared with createSubnet
func (s *Service) createRouteTableWithRoutes(clusterName string, vpc *v1alpha1.VPC, routes []*ec2.Route, isPublic bool) (*v1alpha1.RouteTable, error) {
	out, err := s.EC2.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpc.ID),
	})

	suffix := "private"
	role := TagValueCommonRole
	if isPublic {
		suffix = "public"
		role = TagValueBastionRole
	}
	name := fmt.Sprintf("%s-rt-%s", clusterName, suffix)
	if err := s.createTags(clusterName, *out.RouteTable.RouteTableId, ResourceLifecycleOwned, name, role, nil); err != nil {
		return nil, errors.Wrapf(err, "failed to tag route table %q", *out.RouteTable.RouteTableId)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create route table in vpc %q", vpc.ID)
	}

	for _, route := range routes {
		_, err := s.EC2.CreateRoute(&ec2.CreateRouteInput{
			RouteTableId:                out.RouteTable.RouteTableId,
			DestinationCidrBlock:        route.DestinationCidrBlock,
			DestinationIpv6CidrBlock:    route.DestinationIpv6CidrBlock,
			EgressOnlyInternetGatewayId: route.EgressOnlyInternetGatewayId,
			GatewayId:                   route.GatewayId,
			InstanceId:                  route.InstanceId,
			NatGatewayId:                route.NatGatewayId,
			NetworkInterfaceId:          route.NetworkInterfaceId,
			VpcPeeringConnectionId:      route.VpcPeeringConnectionId,
		})

		if err != nil {
			// TODO(vincepri): cleanup the route table if this fails.
			return nil, errors.Wrapf(err, "failed to create route in route table %q: %s", *out.RouteTable.RouteTableId, route.GoString())
		}
	}

	return &v1alpha1.RouteTable{
		ID: *out.RouteTable.RouteTableId,
	}, nil
}

func (s *Service) associateRouteTable(rt *v1alpha1.RouteTable, subnetID string) error {
	_, err := s.EC2.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rt.ID),
		SubnetId:     aws.String(subnetID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to associate route table %q to subnet %q", rt.ID, subnetID)
	}

	return nil
}

func (s *Service) getDefaultPrivateRoutes(natGatewayID string) []*ec2.Route {
	return []*ec2.Route{
		{
			DestinationCidrBlock: aws.String(anyIPv4CidrBlock),
			NatGatewayId:         aws.String(natGatewayID),
		},
	}
}

func (s *Service) getDefaultPublicRoutes(internetGatewayID string) []*ec2.Route {
	return []*ec2.Route{
		{
			DestinationCidrBlock: aws.String(anyIPv4CidrBlock),
			GatewayId:            aws.String(internetGatewayID),
		},
	}
}
