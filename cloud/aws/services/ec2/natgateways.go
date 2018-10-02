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
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

func (s *Service) reconcileNatGateways(clusterName string, subnets v1alpha1.Subnets, vpc *v1alpha1.VPC) error {
	glog.V(2).Infof("Reconciling NAT gateways")

	if len(subnets.FilterPrivate()) == 0 {
		glog.V(2).Infof("No private subnets available, skipping NAT gateways")
		return nil
	} else if len(subnets.FilterPublic()) == 0 {
		glog.V(2).Infof("No public subnets available. Cannot create NAT gateways for private subnets, this might be a configuration error.")
		return nil
	}

	existing, err := s.describeNatGatewaysBySubnet(vpc.ID)
	if err != nil {
		return err
	}

	for _, sn := range subnets.FilterPublic() {
		if sn.ID == "" {
			continue
		}

		if _, ok := existing[sn.ID]; ok {
			continue
		}

		ng, err := s.createNatGateway(clusterName, sn.ID)
		if err != nil {
			return err
		}

		glog.V(2).Infof("Created new NAT gateway in subnet %q: %s", sn.ID, *ng.NatGatewayId)

		sn.NatGatewayID = ng.NatGatewayId
	}

	return nil
}

func (s *Service) describeNatGatewaysBySubnet(vpcID string) (map[string]*ec2.NatGateway, error) {
	describeNatGatewayInput := &ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	gateways := make(map[string]*ec2.NatGateway)

	err := s.EC2.DescribeNatGatewaysPages(describeNatGatewayInput,
		func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, r := range page.NatGateways {
				gateways[*r.SubnetId] = r
			}
			return !lastPage
		})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe NAT gateways with VPC ID %q", vpcID)
	}

	return gateways, nil
}

func (s *Service) createNatGateway(clusterName string, subnetID string) (*ec2.NatGateway, error) {
	ip, err := s.allocateAddress()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create IP address for NAT gateway for subnet ID %q", subnetID)
	}

	out, err := s.EC2.CreateNatGateway(&ec2.CreateNatGatewayInput{
		SubnetId:     aws.String(subnetID),
		AllocationId: aws.String(ip),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create NAT gateway for subnet ID %q", subnetID)
	}

	wReq := &ec2.DescribeNatGatewaysInput{NatGatewayIds: []*string{out.NatGateway.NatGatewayId}}
	if err := s.EC2.WaitUntilNatGatewayAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for nat gateway %q in subnet %q", *out.NatGateway.NatGatewayId, subnetID)
	}

	if err := s.createTags(clusterName, *out.NatGateway.NatGatewayId, ResourceLifecycleOwned, nil); err != nil {
		return nil, errors.Wrapf(err, "failed to tag nat gateway %q", *out.NatGateway.NatGatewayId)
	}

	return out.NatGateway, nil
}

func (s *Service) getNatGatewayForSubnet(subnets v1alpha1.Subnets, sn *v1alpha1.Subnet) (string, error) {
	if sn.IsPublic {
		return "", errors.Errorf("cannot get NAT gateway for a public subnet, got id %q", sn.ID)
	}

	azGateways := make(map[string][]string)
	for _, psn := range subnets.FilterPublic() {
		if psn.NatGatewayID == nil {
			continue
		}

		azGateways[psn.AvailabilityZone] = append(azGateways[psn.AvailabilityZone], *psn.NatGatewayID)
	}

	if gws, ok := azGateways[sn.AvailabilityZone]; ok && len(gws) > 0 {
		return gws[0], nil
	}

	return "", errors.Errorf("no nat gateways are available in availability zone %q for subnet %q", sn.AvailabilityZone, sn.ID)
}
