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

func (s *Service) reconcileRouteTables(in *v1alpha1.Network) error {
	// TODO(vincepri): add reconcile code when NAT gateways are ready.
	return nil
}

func (s *Service) describeVpcRouteTables(vpcID string) ([]*ec2.RouteTable, error) {
	out, err := s.ec2.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe route tables in vpc %q", vpcID)
	}

	return out.RouteTables, nil
}

func (s *Service) createRouteTable(rt *v1alpha1.RouteTable, vpc *v1alpha1.VPC) (*v1alpha1.RouteTable, error) {
	out, err := s.ec2.CreateRouteTable(&ec2.CreateRouteTableInput{
		VpcId: aws.String(vpc.ID),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create route table")
	}

	return &v1alpha1.RouteTable{
		ID: *out.RouteTable.RouteTableId,
	}, nil
}

func (s *Service) associateRouteTable(rt *v1alpha1.RouteTable, subnetID string) error {
	_, err := s.ec2.AssociateRouteTable(&ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rt.ID),
		SubnetId:     aws.String(subnetID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to associate route table %q to subnet %q", rt.ID, subnetID)
	}

	return nil
}
