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

func (s *Service) reconcileGateways(in *v1alpha1.Network) error {
	// Nat gateway.
	ngs, err := s.describeVpcNatGateways(in.VPC)
	if IsNotFound(err) {
		ps := in.Subnets.FilterPublic()
		if len(ps) == 0 {
			return errors.New("cannot create nat gateway without a public subnet")
		}

		ng, err := s.createNatGateway(in.VPC, ps[0])
		if err != nil {
			return err
		}

		ngs = []*ec2.NatGateway{ng}
	} else if err != nil {
		return err
	}

	// Internet gateway.
	igs, err := s.describeVpcInternetGateways(in.VPC)
	if IsNotFound(err) {
		ig, err := s.createInternetGateway(in.VPC)
		if err != nil {
			return nil
		}
		igs = []*ec2.InternetGateway{ig}
	} else if err != nil {
		return err
	}

	in.NatGatewayID = ngs[0].NatGatewayId
	in.InternetGatewayID = igs[0].InternetGatewayId
	return nil
}

func (s *Service) createNatGateway(vpc *v1alpha1.VPC, publicSubnet *v1alpha1.Subnet) (*ec2.NatGateway, error) {
	eip, err := s.ec2.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String(ec2.DomainTypeVpc),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to allocate elastic ip for nat gateway")
	}

	out, err := s.ec2.CreateNatGateway(&ec2.CreateNatGatewayInput{
		AllocationId: eip.AllocationId,
		SubnetId:     aws.String(publicSubnet.ID),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create nat gateway")
	}

	return out.NatGateway, nil
}

func (s *Service) createInternetGateway(vpc *v1alpha1.VPC) (*ec2.InternetGateway, error) {
	ig, err := s.ec2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create internet gateway")
	}

	_, err = s.ec2.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: ig.InternetGateway.InternetGatewayId,
		VpcId:             aws.String(vpc.ID),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to attach internet gateway %q to vpc %q", *ig.InternetGateway.InternetGatewayId, vpc.ID)
	}

	return ig.InternetGateway, nil
}

func (s *Service) describeVpcNatGateways(vpc *v1alpha1.VPC) ([]*ec2.NatGateway, error) {
	out, err := s.ec2.DescribeNatGateways(&ec2.DescribeNatGatewaysInput{
		Filter: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc.ID)},
			},
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe nat gateways in vpc %q", vpc.ID)
	}

	if len(out.NatGateways) == 0 {
		return nil, NewNotFound(errors.Errorf("no nat gateways found in vpc %q", vpc.ID))
	} else if len(out.NatGateways) > 1 {
		return nil, NewNotFound(errors.Errorf("multiple nat gateways found with supplied filters: %s", out.GoString()))
	}

	return out.NatGateways, nil
}

func (s *Service) describeVpcInternetGateways(vpc *v1alpha1.VPC) ([]*ec2.InternetGateway, error) {
	out, err := s.ec2.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []*string{aws.String(vpc.ID)},
			},
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe nat gateways in vpc %q", vpc.ID)
	}

	if len(out.InternetGateways) == 0 {
		return nil, NewNotFound(errors.Errorf("no nat gateways found in vpc %q", vpc.ID))
	}

	return out.InternetGateways, nil
}
