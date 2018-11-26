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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

func (s *Service) reconcileInternetGateways(clusterName string, in *v1alpha1.Network) error {
	klog.V(2).Infof("Reconciling internet gateways")

	igs, err := s.describeVpcInternetGateways(clusterName, &in.VPC)
	if awserrors.IsNotFound(err) {
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

func (s *Service) deleteInternetGateways(clusterName string, in *v1alpha1.Network) error {
	igs, err := s.describeVpcInternetGateways(clusterName, &in.VPC)
	if awserrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, ig := range igs {
		detachReq := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
			VpcId:             aws.String(in.VPC.ID),
		}

		if _, err := s.scope.EC2.DetachInternetGateway(detachReq); err != nil {
			return errors.Wrapf(err, "failed to detach internet gateway %q", *ig.InternetGatewayId)
		}

		klog.Infof("detached internet gateway %q from VPC %q", *ig.InternetGatewayId, in.VPC.ID)

		deleteReq := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
		}

		if _, err = s.scope.EC2.DeleteInternetGateway(deleteReq); err != nil {
			return errors.Wrapf(err, "failed to delete internet gateway %q", *ig.InternetGatewayId)
		}

		klog.Infof("Deleted internet gateway %q", in.VPC.ID)
	}
	return nil
}

func (s *Service) createInternetGateway(clusterName string, vpc *v1alpha1.VPC) (*ec2.InternetGateway, error) {
	ig, err := s.scope.EC2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create internet gateway")
	}

	name := fmt.Sprintf("%s-igw", clusterName)

	applyTagsParams := &tags.ApplyParams{
		EC2Client: s.scope.EC2,
		BuildParams: tags.BuildParams{
			ClusterName: clusterName,
			ResourceID:  *ig.InternetGateway.InternetGatewayId,
			Lifecycle:   tags.ResourceLifecycleOwned,
			Name:        aws.String(name),
			Role:        aws.String(tags.ValueCommonRole),
		},
	}

	if err := tags.Apply(applyTagsParams); err != nil {
		return nil, errors.Wrapf(err, "failed to tag internet gateway %q", *ig.InternetGateway.InternetGatewayId)
	}

	klog.Infof("created internet gateway %q", vpc.ID)

	_, err = s.scope.EC2.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
		InternetGatewayId: ig.InternetGateway.InternetGatewayId,
		VpcId:             aws.String(vpc.ID),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to attach internet gateway %q to vpc %q", *ig.InternetGateway.InternetGatewayId, vpc.ID)
	}

	klog.Infof("attached internet gateway %q to VPC %q", *ig.InternetGateway.InternetGatewayId, vpc.ID)
	return ig.InternetGateway, nil
}

func (s *Service) describeVpcInternetGateways(clusterName string, vpc *v1alpha1.VPC) ([]*ec2.InternetGateway, error) {
	out, err := s.scope.EC2.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCAttachment(vpc.ID),
			filter.EC2.Cluster(clusterName),
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe nat gateways in vpc %q", vpc.ID)
	}

	if len(out.InternetGateways) == 0 {
		return nil, awserrors.NewNotFound(errors.Errorf("no nat gateways found in vpc %q", vpc.ID))
	}

	return out.InternetGateways, nil
}
