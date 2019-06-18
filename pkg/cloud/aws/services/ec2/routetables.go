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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

const (
	anyIPv4CidrBlock       = "0.0.0.0/0"
	mainRouteTableInVPCKey = "main"
)

func (s *Service) reconcileRouteTables() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping routing tables reconcile in unmanaged mode")
		return nil
	}

	s.scope.V(2).Info("Reconciling routing tables")

	subnetRouteMap, err := s.describeVpcRouteTablesBySubnet()
	if err != nil {
		return err
	}

	for _, sn := range s.scope.Subnets() {
		if igw, ok := subnetRouteMap[sn.ID]; ok {
			s.scope.V(2).Info("Subnet is already associated with route table", "subnet-id", sn.ID, "route-table-id", *igw.RouteTableId)
			// TODO(vincepri): if the route table ids are both non-empty and they don't match, replace the association.
			// TODO(vincepri): check that everything is in order, e.g. routes match the subnet type.

			// Make sure tags are up to date.
			err := tags.Ensure(converters.TagsToMap(igw.Tags), &tags.ApplyParams{
				EC2Client:   s.scope.EC2,
				BuildParams: s.getRouteTableTagParams(*igw.RouteTableId, sn.IsPublic),
			})

			if err != nil {
				return errors.Wrapf(err, "failed to ensure tags on route table %q", *igw.RouteTableId)
			}

			continue
		}

		// For each subnet that doesn't have a routing table associated with it,
		// create a new table with the appropriate default routes and associate it to the subnet.
		var routes []*ec2.Route
		if sn.IsPublic {
			if s.scope.VPC().InternetGatewayID == nil {
				return errors.Errorf("failed to create routing tables: internet gateway for %q is nil", s.scope.VPC().ID)
			}

			routes = s.getDefaultPublicRoutes()
		} else {
			natGatewayID, err := s.getNatGatewayForSubnet(sn)
			if err != nil {
				return err
			}

			routes = s.getDefaultPrivateRoutes(natGatewayID)
		}

		rt, err := s.createRouteTableWithRoutes(routes, sn.IsPublic)
		if err != nil {
			return err
		}

		if err := s.associateRouteTable(rt, sn.ID); err != nil {
			return err
		}

		s.scope.V(2).Info("Subnet has been associated with route table", "subnet-id", sn.ID, "route-table-id", rt.ID)
		sn.RouteTableID = aws.String(rt.ID)
	}

	return nil
}

func (s *Service) describeVpcRouteTablesBySubnet() (map[string]*ec2.RouteTable, error) {
	rts, err := s.describeVpcRouteTables()
	if err != nil {
		return nil, err
	}

	// Amazon allows a subnet to be associated only with a single routing table
	// https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html.
	res := make(map[string]*ec2.RouteTable)
	for _, rt := range rts {
		for _, as := range rt.Associations {
			if as.Main != nil && *as.Main {
				res[mainRouteTableInVPCKey] = rt
			}
			if as.SubnetId == nil {
				continue
			}

			res[*as.SubnetId] = rt
		}
	}

	return res, nil
}

func (s *Service) deleteRouteTables() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping routing tables deletion in unmanaged mode")
		return nil
	}

	rts, err := s.describeVpcRouteTables()
	if err != nil {
		return errors.Wrapf(err, "failed to describe route tables in vpc %q", s.scope.VPC().ID)
	}

	for _, rt := range rts {
		for _, as := range rt.Associations {
			if as.SubnetId == nil {
				continue
			}

			if _, err := s.scope.EC2.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{AssociationId: as.RouteTableAssociationId}); err != nil {
				return errors.Wrapf(err, "failed to disassociate route table %q from subnet %q", *rt.RouteTableId, *as.SubnetId)
			}

			s.scope.Info("Deleted association between route table and subnet", "route-table-id", *rt.RouteTableId, "subnet-id", *as.SubnetId)
		}

		if _, err := s.scope.EC2.DeleteRouteTable(&ec2.DeleteRouteTableInput{RouteTableId: rt.RouteTableId}); err != nil {
			return errors.Wrapf(err, "failed to delete route table %q", *rt.RouteTableId)
		}

		s.scope.Info("Deleted route table", "route-table-id", *rt.RouteTableId)
	}
	return nil
}

func (s *Service) describeVpcRouteTables() ([]*ec2.RouteTable, error) {
	filters := []*ec2.Filter{
		filter.EC2.VPC(s.scope.VPC().ID),
	}

	if !s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		filters = append(filters, filter.EC2.Cluster(s.scope.Name()))
	}

	out, err := s.scope.EC2.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: filters,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe route tables in vpc %q", s.scope.VPC().ID)
	}

	return out.RouteTables, nil
}

func (s *Service) createRouteTableWithRoutes(routes []*ec2.Route, isPublic bool) (*v1alpha1.RouteTable, error) {
	out, err := s.scope.EC2.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(s.scope.VPC().ID),
	})

	applyTagsParams := &tags.ApplyParams{
		EC2Client:   s.scope.EC2,
		BuildParams: s.getRouteTableTagParams(*out.RouteTable.RouteTableId, isPublic),
	}

	if err := tags.Apply(applyTagsParams); err != nil {
		return nil, errors.Wrapf(err, "failed to tag route table %q", *out.RouteTable.RouteTableId)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create route table in vpc %q", s.scope.VPC().ID)
	}

	for _, route := range routes {
		_, err := s.scope.EC2.CreateRoute(&ec2.CreateRouteInput{
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
	_, err := s.scope.EC2.AssociateRouteTable(&ec2.AssociateRouteTableInput{
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

func (s *Service) getDefaultPublicRoutes() []*ec2.Route {
	return []*ec2.Route{
		{
			DestinationCidrBlock: aws.String(anyIPv4CidrBlock),
			GatewayId:            aws.String(*s.scope.VPC().InternetGatewayID),
		},
	}
}

func (s *Service) getRouteTableTagParams(id string, public bool) v1alpha1.BuildParams {
	var name strings.Builder

	name.WriteString(s.scope.Name())
	name.WriteString("-rt-")
	if public {
		name.WriteString("public")
	} else {
		name.WriteString("private")
	}

	return v1alpha1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   v1alpha1.ResourceLifecycleOwned,
		Name:        aws.String(name.String()),
		Role:        aws.String(v1alpha1.CommonRoleTagValue),
	}
}
