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

	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

const (
	defaultVPCCidr = "10.0.0.0/16"
)

func (s *Service) reconcileVPC() error {
	s.scope.V(2).Info("Reconciling VPC")

	vpc, err := s.describeVPC()
	if awserrors.IsNotFound(err) {
		// Create a new managed vpc.
		vpc, err = s.createVPC()
		if err != nil {
			return errors.Wrap(err, "failed to create new vpc")
		}

	} else if err != nil {
		return errors.Wrap(err, "failed to describe VPCs")
	}

	if vpc.IsUnmanaged(s.scope.Name()) {
		vpc.DeepCopyInto(s.scope.VPC())
		s.scope.V(2).Info("Working on unmanaged VPC", "vpc-id", vpc.ID)
		return nil
	}

	// Make sure tags are up to date.
	err = tags.Ensure(vpc.Tags, &tags.ApplyParams{
		EC2Client:   s.scope.EC2,
		BuildParams: s.getVPCTagParams(vpc.ID),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to tag vpc %q", vpc.ID)
	}

	vpc.DeepCopyInto(s.scope.VPC())
	s.scope.V(2).Info("Working on managed VPC", "vpc-id", vpc.ID)
	return nil
}

func (s *Service) createVPC() (*v1alpha1.VPCSpec, error) {
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		return nil, errors.Errorf("cannot create a managed vpc in unmanaged mode")
	}

	if s.scope.VPC().CidrBlock == "" {
		s.scope.VPC().CidrBlock = defaultVPCCidr
	}

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(s.scope.VPC().CidrBlock),
	}

	out, err := s.scope.EC2.CreateVpc(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vpc")
	}

	// Cannot set both attributes at the same time.
	attrInput := &ec2.ModifyVpcAttributeInput{
		VpcId:              out.Vpc.VpcId,
		EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	}

	if _, err = s.scope.EC2.ModifyVpcAttribute(attrInput); err != nil {
		return nil, errors.Wrap(err, "failed to set vpc attributes")
	}

	attrInput = &ec2.ModifyVpcAttributeInput{
		VpcId:            out.Vpc.VpcId,
		EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
	}

	if _, err = s.scope.EC2.ModifyVpcAttribute(attrInput); err != nil {
		return nil, errors.Wrap(err, "failed to set vpc attributes")
	}

	wReq := &ec2.DescribeVpcsInput{VpcIds: []*string{out.Vpc.VpcId}}
	if err := s.scope.EC2.WaitUntilVpcAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for vpc %q", *out.Vpc.VpcId)
	}

	s.scope.V(2).Info("Created new VPC with cidr", "vpc-id", *out.Vpc.VpcId, "cidr-block", *out.Vpc.CidrBlock)
	record.Eventf(s.scope.Cluster, "CreatedVPC", "Created new managed VPC %q", *out.Vpc.VpcId)

	tagParams := s.getVPCTagParams(*out.Vpc.VpcId)
	tagApply := &tags.ApplyParams{
		EC2Client:   s.scope.EC2,
		BuildParams: tagParams,
	}

	if err := tags.Apply(tagApply); err != nil {
		return nil, err
	}

	return &v1alpha1.VPCSpec{
		ID:        *out.Vpc.VpcId,
		CidrBlock: *out.Vpc.CidrBlock,
		Tags:      v1alpha1.Build(tagParams),
	}, nil
}

func (s *Service) deleteVPC() error {
	vpc := s.scope.VPC()

	if vpc.IsUnmanaged(s.scope.Name()) {
		s.scope.V(4).Info("Skipping VPC deletion in unmanaged mode")
		return nil
	}

	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpc.ID),
	}

	_, err := s.scope.EC2.DeleteVpc(input)
	if err != nil {
		// Ignore if it's already deleted
		if code, ok := awserrors.Code(err); code == "InvalidVpcID.NotFound" && ok {
			s.scope.V(4).Info("Skipping VPC deletion, VPC not found")
			return nil
		}
		return errors.Wrapf(err, "failed to delete vpc %q", vpc.ID)
	}

	s.scope.V(2).Info("Deleted VPC", "vpc-id", vpc.ID)
	record.Eventf(s.scope.Cluster, "DeletedVPC", "Deleted managed VPC %q", vpc.ID)
	return nil
}

func (s *Service) describeVPC() (*v1alpha1.VPCSpec, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCStates(ec2.VpcStatePending, ec2.VpcStateAvailable),
		},
	}

	if s.scope.VPC().ID == "" {
		// Try to find a previously created and tagged VPC
		input.Filters = append(input.Filters, filter.EC2.Cluster(s.scope.Name()))
	} else {
		input.VpcIds = []*string{aws.String(s.scope.VPC().ID)}
	}

	out, err := s.scope.EC2.DescribeVpcs(input)
	if err != nil {
		if awserrors.IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed to query ec2 for VPCs")
	}

	if len(out.Vpcs) == 0 {
		return nil, awserrors.NewNotFound(errors.Errorf("could not find vpc %q", s.scope.VPC().ID))
	} else if len(out.Vpcs) > 1 {
		return nil, awserrors.NewConflict(errors.Errorf("found more than one vpc with supplied filters. Please clean up extra VPCs: %s", out.GoString()))
	}

	switch *out.Vpcs[0].State {
	case ec2.VpcStateAvailable, ec2.VpcStatePending:
	default:
		return nil, awserrors.NewNotFound(errors.Errorf("could not find available or pending vpc"))
	}

	return &v1alpha1.VPCSpec{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
		Tags:      converters.TagsToMap(out.Vpcs[0].Tags),
	}, nil
}

func (s *Service) getVPCTagParams(id string) v1alpha1.BuildParams {
	name := fmt.Sprintf("%s-vpc", s.scope.Name())

	return v1alpha1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   v1alpha1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(v1alpha1.CommonRoleTagValue),
	}
}
