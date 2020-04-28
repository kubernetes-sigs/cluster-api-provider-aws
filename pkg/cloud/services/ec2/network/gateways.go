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
	"fmt"

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
)

func (s *Service) reconcileInternetGateways() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping internet gateways reconcile in unmanaged mode")
		return nil
	}

	s.scope.V(2).Info("Reconciling internet gateways")

	igs, err := s.describeVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
			return errors.Errorf("failed to validate network: no internet gateways found in VPC %q", s.scope.VPC().ID)
		}

		ig, err := s.createInternetGateway()
		if err != nil {
			return nil
		}
		igs = []*ec2.InternetGateway{ig}
	} else if err != nil {
		return err
	}

	gateway := igs[0]
	s.scope.VPC().InternetGatewayID = gateway.InternetGatewayId

	// Make sure tags are up to date.
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := tags.Ensure(converters.TagsToMap(gateway.Tags), &tags.ApplyParams{
			EC2Client:   s.scope.EC2,
			BuildParams: s.getGatewayTagParams(*gateway.InternetGatewayId),
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.Target, "FailedTagInternetGateway", "Failed to tag managed Internet Gateway %q: %v", gateway.InternetGatewayId, err)
		return errors.Wrapf(err, "failed to tag internet gateway %q", *gateway.InternetGatewayId)
	}

	return nil
}

func (s *Service) deleteInternetGateways() error {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping internet gateway deletion in unmanaged mode")
		return nil
	}

	igs, err := s.describeVpcInternetGateways()
	if awserrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, ig := range igs {
		detachReq := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
			VpcId:             aws.String(s.scope.VPC().ID),
		}

		if _, err := s.scope.EC2.DetachInternetGateway(detachReq); err != nil {
			record.Warnf(s.scope.Target, "FailedDetachInternetGateway", "Failed to detach Internet Gateway %q from VPC %q: %v", *ig.InternetGatewayId, s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to detach internet gateway %q", *ig.InternetGatewayId)
		}

		record.Eventf(s.scope.Target, "SuccessfulDetachInternetGateway", "Detached Internet Gateway %q from VPC %q", *ig.InternetGatewayId, s.scope.VPC().ID)
		s.scope.Info("Detached internet gateway from VPC", "internet-gateway-id", *ig.InternetGatewayId, "vpc-id", s.scope.VPC().ID)

		deleteReq := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: ig.InternetGatewayId,
		}

		if _, err = s.scope.EC2.DeleteInternetGateway(deleteReq); err != nil {
			record.Warnf(s.scope.Target, "FailedDeleteInternetGateway", "Failed to delete Internet Gateway %q previously attached to VPC %q: %v", *ig.InternetGatewayId, s.scope.VPC().ID, err)
			return errors.Wrapf(err, "failed to delete internet gateway %q", *ig.InternetGatewayId)
		}

		record.Eventf(s.scope.Target, "SuccessfulDeleteInternetGateway", "Deleted Internet Gateway %q previously attached to VPC %q", *ig.InternetGatewayId, s.scope.VPC().ID)
		s.scope.Info("Deleted internet gateway in VPC", "internet-gateway-id", *ig.InternetGatewayId, "vpc-id", s.scope.VPC().ID)
	}

	return nil
}

func (s *Service) createInternetGateway() (*ec2.InternetGateway, error) {
	ig, err := s.scope.EC2.CreateInternetGateway(&ec2.CreateInternetGatewayInput{})
	if err != nil {
		record.Warnf(s.scope.Target, "FailedCreateInternetGateway", "Failed to create new managed Internet Gateway: %v", err)
		return nil, errors.Wrap(err, "failed to create internet gateway")
	}
	record.Eventf(s.scope.Target, "SuccessfulCreateInternetGateway", "Created new managed Internet Gateway %q", *ig.InternetGateway.InternetGatewayId)
	s.scope.Info("Created internet gateway for VPC", "vpc-id", s.scope.VPC().ID)

	tagParams := s.getGatewayTagParams(*ig.InternetGateway.InternetGatewayId)
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := tags.Apply(&tags.ApplyParams{
			EC2Client:   s.scope.EC2,
			BuildParams: tagParams,
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.Target, "FailedTagInternetGateway", "Failed to tag managed Internet Gateway %q: %v", *ig.InternetGateway.InternetGatewayId, err)
		return nil, errors.Wrapf(err, "failed to tag internet gateway %q", *ig.InternetGateway.InternetGatewayId)
	}

	// Update the tags, so that when ig.InternetGateway is returned it has the
	// latest tag data rather than returning empty tags.
	ig.InternetGateway.Tags = converters.MapToTags(infrav1.Build(tagParams))
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if _, err := s.scope.EC2.AttachInternetGateway(&ec2.AttachInternetGatewayInput{
			InternetGatewayId: ig.InternetGateway.InternetGatewayId,
			VpcId:             aws.String(s.scope.VPC().ID),
		}); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.InternetGatewayNotFound); err != nil {
		record.Warnf(s.scope.Target, "FailedAttachInternetGateway", "Failed to attach managed Internet Gateway %q to vpc %q: %v", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID, err)
		return nil, errors.Wrapf(err, "failed to attach internet gateway %q to vpc %q", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID)
	}
	record.Eventf(s.scope.Target, "SuccessfulAttachInternetGateway", "Internet Gateway %q attached to VPC %q", *ig.InternetGateway.InternetGatewayId, s.scope.VPC().ID)
	s.scope.Info("attached internet gateway to VPC", "internet-gateway-id", *ig.InternetGateway.InternetGatewayId, "vpc-id", s.scope.VPC().ID)

	return ig.InternetGateway, nil
}

func (s *Service) describeVpcInternetGateways() ([]*ec2.InternetGateway, error) {
	out, err := s.scope.EC2.DescribeInternetGateways(&ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCAttachment(s.scope.VPC().ID),
		},
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe internet gateways in vpc %q", s.scope.VPC().ID)
	}

	if len(out.InternetGateways) == 0 {
		return nil, awserrors.NewNotFound(errors.Errorf("no internet gateways found in vpc %q", s.scope.VPC().ID))
	}

	return out.InternetGateways, nil
}

func (s *Service) getGatewayTagParams(id string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-igw", s.scope.Name())

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
