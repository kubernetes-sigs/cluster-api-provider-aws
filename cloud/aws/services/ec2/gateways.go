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

func (s *Service) reconcileInternetGateways(clusterName string, in *v1alpha1.Network) error {
	glog.V(2).Infof("Reconciling internet gateways")

	igs, err := s.describeVpcInternetGateways(&in.VPC)
	if IsNotFound(err) {
		ig, err := s.createInternetGateway(clusterName, &in.VPC)
		if err != nil {
			return nil
		}
		igs = []*ec2.InternetGateway{ig}
	} else if err != nil {
		return err
	}

	in.InternetGatewayID = igs[0].InternetGatewayId
	return nil
}

func (s *Service) createInternetGateway(clusterName string, vpc *v1alpha1.VPC) (*ec2.InternetGateway, error) {
	ig, err := s.EC2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create internet gateway")
	}

	if err := s.createTags(clusterName, *ig.InternetGateway.InternetGatewayId, ResourceLifecycleOwned, nil); err != nil {
		return nil, errors.Wrapf(err, "failed to tag internet gateway %q", *ig.InternetGateway.InternetGatewayId)
	}

	_, err = s.EC2.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: ig.InternetGateway.InternetGatewayId,
		VpcId:             aws.String(vpc.ID),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to attach internet gateway %q to vpc %q", *ig.InternetGateway.InternetGatewayId, vpc.ID)
	}

	return ig.InternetGateway, nil
}

func (s *Service) describeVpcInternetGateways(vpc *v1alpha1.VPC) ([]*ec2.InternetGateway, error) {
	out, err := s.EC2.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
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
