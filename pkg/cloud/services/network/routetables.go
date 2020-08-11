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

package network

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	"sigs.k8s.io/cluster-api/util/conditions"
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

	for i := range s.scope.Subnets() {
		// We need to compile the minimum routes for this subnet first, so we can compare it or create them.
		var routes []*ec2.Route
		sn := s.scope.Subnets()[i]
		if sn.IsPublic {
			if s.scope.VPC().InternetGatewayID == nil {
				return errors.Errorf("failed to create routing tables: internet gateway for %q is nil", s.scope.VPC().ID)
			}
			routes = append(routes, s.getGatewayPublicRoute())
		} else {
			natGatewayID, err := s.getNatGatewayForSubnet(sn)
			if err != nil {
				return err
			}
			routes = append(routes, s.getNatGatewayPrivateRoute(natGatewayID))
		}

		if rt, ok := subnetRouteMap[sn.ID]; ok {
			s.scope.V(2).Info("Subnet is already associated with route table", "subnet-id", sn.ID, "route-table-id", *rt.RouteTableId)
			// TODO(vincepri): check that everything is in order, e.g. routes match the subnet type.

			// For managed environments we need to reconcile the routes of our tables if there is a mistmatch.
			// For example, a gateway can be deleted and our controller will re-create it, then we replace the route
			// for the subnet to allow traffic to flow.
			for _, currentRoute := range rt.Routes {
				for i := range routes {
					// Routes destination cidr blocks must be unique within a routing table.
					// If there is a mistmatch, we replace the routing association.
					specRoute := routes[i]
					if *currentRoute.DestinationCidrBlock == *specRoute.DestinationCidrBlock &&
						((currentRoute.GatewayId != nil && *currentRoute.GatewayId != *specRoute.GatewayId) ||
							(currentRoute.NatGatewayId != nil && *currentRoute.NatGatewayId != *specRoute.NatGatewayId)) {
						if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
							if _, err := s.EC2Client.ReplaceRoute(&ec2.ReplaceRouteInput{
								RouteTableId:         rt.RouteTableId,
								DestinationCidrBlock: specRoute.DestinationCidrBlock,
								GatewayId:            specRoute.GatewayId,
								NatGatewayId:         specRoute.NatGatewayId,
							}); err != nil {
								return false, err
							}
							return true, nil
						}); err != nil {
							record.Warnf(s.scope.InfraCluster(), "FailedReplaceRoute", "Failed to replace outdated route on managed RouteTable %q: %v", *rt.RouteTableId, err)
							return errors.Wrapf(err, "failed to replace outdated route on route table %q", *rt.RouteTableId)
						}
					}
				}
			}

			// Make sure tags are up to date.
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				buildParams := s.getRouteTableTagParams(*rt.RouteTableId, sn.IsPublic, sn.AvailabilityZone)
				tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
				if err := tagsBuilder.Ensure(converters.TagsToMap(rt.Tags)); err != nil {
					return false, err
				}
				return true, nil
			}, awserrors.RouteTableNotFound); err != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedTagRouteTable", "Failed to tag managed RouteTable %q: %v", *rt.RouteTableId, err)
				return errors.Wrapf(err, "failed to ensure tags on route table %q", *rt.RouteTableId)
			}

			// Not recording "SuccessfulTagRouteTable" here as we don't know if this was a no-op or an actual change
			continue
		}

		// For each subnet that doesn't have a routing table associated with it,
		// create a new table with the appropriate default routes and associate it to the subnet.
		rt, err := s.createRouteTableWithRoutes(routes, sn.IsPublic, sn.AvailabilityZone)
		if err != nil {
			return err
		}

		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if err := s.associateRouteTable(rt, sn.ID); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.RouteTableNotFound, awserrors.SubnetNotFound); err != nil {
			return err
		}

		s.scope.V(2).Info("Subnet has been associated with route table", "subnet-id", sn.ID, "route-table-id", rt.ID)
		sn.RouteTableID = aws.String(rt.ID)
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition)
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

			if _, err := s.EC2Client.DisassociateRouteTable(&ec2.DisassociateRouteTableInput{AssociationId: as.RouteTableAssociationId}); err != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedDisassociateRouteTable", "Failed to disassociate managed RouteTable %q from Subnet %q: %v", *rt.RouteTableId, *as.SubnetId, err)
				return errors.Wrapf(err, "failed to disassociate route table %q from subnet %q", *rt.RouteTableId, *as.SubnetId)
			}

			record.Eventf(s.scope.InfraCluster(), "SuccessfulDisassociateRouteTable", "Disassociated managed RouteTable %q from subnet %q", *rt.RouteTableId, *as.SubnetId)
			s.scope.Info("Deleted association between route table and subnet", "route-table-id", *rt.RouteTableId, "subnet-id", *as.SubnetId)
		}

		if _, err := s.EC2Client.DeleteRouteTable(&ec2.DeleteRouteTableInput{RouteTableId: rt.RouteTableId}); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDeleteRouteTable", "Failed to delete managed RouteTable %q: %v", *rt.RouteTableId, err)
			return errors.Wrapf(err, "failed to delete route table %q", *rt.RouteTableId)
		}

		record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteRouteTable", "Deleted managed RouteTable %q", *rt.RouteTableId)
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

	out, err := s.EC2Client.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: filters,
	})
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeVPCRouteTable", "Failed to describe route tables in vpc %q: %v", s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to describe route tables in vpc %q", s.scope.VPC().ID)
	}

	return out.RouteTables, nil
}

func (s *Service) createRouteTableWithRoutes(routes []*ec2.Route, isPublic bool, zone string) (*infrav1.RouteTable, error) {
	out, err := s.EC2Client.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(s.scope.VPC().ID),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateRouteTable", "Failed to create managed RouteTable: %v", err)
		return nil, errors.Wrapf(err, "failed to create route table in vpc %q", s.scope.VPC().ID)
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateRouteTable", "Created managed RouteTable %q", *out.RouteTable.RouteTableId)

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getRouteTableTagParams(*out.RouteTable.RouteTableId, isPublic, zone)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Apply(); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.RouteTableNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedTagRouteTable", "Failed to tag managed RouteTable %q: %v", *out.RouteTable.RouteTableId, err)
		return nil, errors.Wrapf(err, "failed to tag route table %q", *out.RouteTable.RouteTableId)
	}
	record.Eventf(s.scope.InfraCluster(), "SuccessfulTagRouteTable", "Tagged managed RouteTable %q", *out.RouteTable.RouteTableId)

	for i := range routes {
		route := routes[i]
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if _, err := s.EC2Client.CreateRoute(&ec2.CreateRouteInput{
				RouteTableId:                out.RouteTable.RouteTableId,
				DestinationCidrBlock:        route.DestinationCidrBlock,
				DestinationIpv6CidrBlock:    route.DestinationIpv6CidrBlock,
				EgressOnlyInternetGatewayId: route.EgressOnlyInternetGatewayId,
				GatewayId:                   route.GatewayId,
				InstanceId:                  route.InstanceId,
				NatGatewayId:                route.NatGatewayId,
				NetworkInterfaceId:          route.NetworkInterfaceId,
				VpcPeeringConnectionId:      route.VpcPeeringConnectionId,
			}); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.RouteTableNotFound, awserrors.NATGatewayNotFound, awserrors.GatewayNotFound); err != nil {
			// TODO(vincepri): cleanup the route table if this fails.
			record.Warnf(s.scope.InfraCluster(), "FailedCreateRoute", "Failed to create route %s for RouteTable %q: %v", route.GoString(), *out.RouteTable.RouteTableId, err)
			return nil, errors.Wrapf(err, "failed to create route in route table %q: %s", *out.RouteTable.RouteTableId, route.GoString())
		}
		record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateRoute", "Created route %s for RouteTable %q", route.GoString(), *out.RouteTable.RouteTableId)
	}

	return &infrav1.RouteTable{
		ID: *out.RouteTable.RouteTableId,
	}, nil
}

func (s *Service) associateRouteTable(rt *infrav1.RouteTable, subnetID string) error {
	_, err := s.EC2Client.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rt.ID),
		SubnetId:     aws.String(subnetID),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAssociateRouteTable", "Failed to associate managed RouteTable %q with Subnet %q: %v", rt.ID, subnetID, err)
		return errors.Wrapf(err, "failed to associate route table %q to subnet %q", rt.ID, subnetID)
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulAssociateRouteTable", "Associated managed RouteTable %q with subnet %q", rt.ID, subnetID)

	return nil
}

func (s *Service) getNatGatewayPrivateRoute(natGatewayID string) *ec2.Route {
	return &ec2.Route{
		DestinationCidrBlock: aws.String(anyIPv4CidrBlock),
		NatGatewayId:         aws.String(natGatewayID),
	}
}

func (s *Service) getGatewayPublicRoute() *ec2.Route {
	return &ec2.Route{
		DestinationCidrBlock: aws.String(anyIPv4CidrBlock),
		GatewayId:            aws.String(*s.scope.VPC().InternetGatewayID),
	}
}

func (s *Service) getRouteTableTagParams(id string, public bool, zone string) infrav1.BuildParams {
	var name strings.Builder

	name.WriteString(s.scope.Name())
	name.WriteString("-rt-")
	if public {
		name.WriteString("public")
	} else {
		name.WriteString("private")
	}
	name.WriteString("-")
	name.WriteString(zone)

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name.String()),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
