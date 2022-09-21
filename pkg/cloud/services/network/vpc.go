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
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const (
	defaultVPCCidr = "10.0.0.0/16"
)

func (s *Service) reconcileVPC() error {
	s.scope.V(2).Info("Reconciling VPC")

	// If the ID is not nil, VPC is either managed or unmanaged but should exist in the AWS.
	if s.scope.VPC().ID != "" {
		vpc, err := s.describeVPCByID()
		if err != nil {
			return errors.Wrap(err, ".spec.vpc.id is set but VPC resource is missing in AWS; failed to describe VPC resources. (might be in creation process)")
		}

		s.scope.VPC().CidrBlock = vpc.CidrBlock
		s.scope.VPC().Tags = vpc.Tags
		s.scope.VPC().EnableIPv6 = vpc.EnableIPv6
		s.scope.VPC().IPv6CidrBlock = vpc.IPv6CidrBlock
		s.scope.VPC().IPv6Pool = vpc.IPv6Pool

		// If VPC is unmanaged, return early.
		if vpc.IsUnmanaged(s.scope.Name()) {
			s.scope.V(2).Info("Working on unmanaged VPC", "vpc-id", vpc.ID)
			if err := s.scope.PatchObject(); err != nil {
				return errors.Wrap(err, "failed to patch unmanaged VPC fields")
			}
			record.Eventf(s.scope.InfraCluster(), "SuccessfulSetVPCAttributes", "Set managed VPC attributes for %q", vpc.ID)
			return nil
		}

		// if the VPC is managed, make managed sure attributes are configured.
		if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			if err := s.ensureManagedVPCAttributes(vpc); err != nil {
				return false, err
			}
			return true, nil
		}, awserrors.VPCNotFound); err != nil {
			return errors.Wrapf(err, "failed to to set vpc attributes for %q", vpc.ID)
		}

		return nil
	}

	// .spec.vpc.id is nil, Create a new managed vpc.
	if !conditions.Has(s.scope.InfraCluster(), infrav1.VpcReadyCondition) {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcCreationStartedReason, clusterv1.ConditionSeverityInfo, "")
		if err := s.scope.PatchObject(); err != nil {
			return errors.Wrap(err, "failed to patch conditions")
		}
	}
	vpc, err := s.createVPC()
	if err != nil {
		return errors.Wrap(err, "failed to create new vpc")
	}
	s.scope.Info("Created VPC", "vpc-id", vpc.ID)

	s.scope.VPC().CidrBlock = vpc.CidrBlock
	s.scope.VPC().IPv6CidrBlock = vpc.IPv6CidrBlock
	s.scope.VPC().IPv6Pool = vpc.IPv6Pool
	s.scope.VPC().EnableIPv6 = vpc.EnableIPv6
	s.scope.VPC().Tags = vpc.Tags
	s.scope.VPC().ID = vpc.ID

	// Make sure attributes are configured
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := s.ensureManagedVPCAttributes(vpc); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.VPCNotFound); err != nil {
		return errors.Wrapf(err, "failed to to set vpc attributes for %q", vpc.ID)
	}

	return nil
}

func (s *Service) ensureManagedVPCAttributes(vpc *infrav1.VPCSpec) error {
	var (
		errs    []error
		updated bool
	)

	// Cannot get or set both attributes at the same time.
	descAttrInput := &ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String(vpc.ID),
		Attribute: aws.String("enableDnsHostnames"),
	}
	vpcAttr, err := s.EC2Client.DescribeVpcAttribute(descAttrInput)
	if err != nil {
		// If the returned error is a 'NotFound' error it should trigger retry
		if code, ok := awserrors.Code(errors.Cause(err)); ok && code == awserrors.VPCNotFound {
			return err
		}
		errs = append(errs, errors.Wrap(err, "failed to describe enableDnsHostnames vpc attribute"))
	} else if !aws.BoolValue(vpcAttr.EnableDnsHostnames.Value) {
		attrInput := &ec2.ModifyVpcAttributeInput{
			VpcId:              aws.String(vpc.ID),
			EnableDnsHostnames: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		if _, err := s.EC2Client.ModifyVpcAttribute(attrInput); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to set enableDnsHostnames vpc attribute"))
		} else {
			updated = true
		}
	}

	descAttrInput = &ec2.DescribeVpcAttributeInput{
		VpcId:     aws.String(vpc.ID),
		Attribute: aws.String("enableDnsSupport"),
	}
	vpcAttr, err = s.EC2Client.DescribeVpcAttribute(descAttrInput)
	if err != nil {
		// If the returned error is a 'NotFound' error it should trigger retry
		if code, ok := awserrors.Code(errors.Cause(err)); ok && code == awserrors.VPCNotFound {
			return err
		}
		errs = append(errs, errors.Wrap(err, "failed to describe enableDnsSupport vpc attribute"))
	} else if !aws.BoolValue(vpcAttr.EnableDnsSupport.Value) {
		attrInput := &ec2.ModifyVpcAttributeInput{
			VpcId:            aws.String(vpc.ID),
			EnableDnsSupport: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}
		if _, err := s.EC2Client.ModifyVpcAttribute(attrInput); err != nil {
			errs = append(errs, errors.Wrap(err, "failed to set enableDnsSupport vpc attribute"))
		} else {
			updated = true
		}
	}

	if len(errs) > 0 {
		record.Warnf(s.scope.InfraCluster(), "FailedSetVPCAttributes", "Failed to set managed VPC attributes for %q: %v", vpc.ID, err)
		return kerrors.NewAggregate(errs)
	}

	if updated {
		record.Eventf(s.scope.InfraCluster(), "SuccessfulSetVPCAttributes", "Set managed VPC attributes for %q", vpc.ID)
	}

	return nil
}

func (s *Service) createVPC() (*infrav1.VPCSpec, error) {
	input := &ec2.CreateVpcInput{
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeVpc, s.getVPCTagParams(services.TemporaryResourceID)),
		},
	}

	// setup BYOIP
	if s.scope.VPC().IPv6CidrBlock != "" {
		input.Ipv6CidrBlock = aws.String(s.scope.VPC().IPv6CidrBlock)
		input.Ipv6Pool = aws.String(s.scope.VPC().IPv6Pool)
		input.AmazonProvidedIpv6CidrBlock = aws.Bool(false)
	} else {
		input.AmazonProvidedIpv6CidrBlock = aws.Bool(s.scope.VPC().EnableIPv6)
	}

	if s.scope.VPC().CidrBlock == "" {
		s.scope.VPC().CidrBlock = defaultVPCCidr
	}
	input.CidrBlock = &s.scope.VPC().CidrBlock

	out, err := s.EC2Client.CreateVpc(input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "Failed to create new managed VPC: %v", err)
		return nil, errors.Wrap(err, "failed to create vpc")
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateVPC", "Created new managed VPC %q", *out.Vpc.VpcId)
	s.scope.V(2).Info("Created new VPC with cidr", "vpc-id", *out.Vpc.VpcId, "cidr-block", *out.Vpc.CidrBlock)

	if !s.scope.VPC().EnableIPv6 {
		return &infrav1.VPCSpec{
			ID:        *out.Vpc.VpcId,
			CidrBlock: *out.Vpc.CidrBlock,
			Tags:      converters.TagsToMap(out.Vpc.Tags),
		}, nil
	}

	// BYOIP was defined, no need to look up the VPC.
	if s.scope.VPC().EnableIPv6 && s.scope.VPC().IPv6CidrBlock != "" {
		return &infrav1.VPCSpec{
			ID:            *out.Vpc.VpcId,
			CidrBlock:     *out.Vpc.CidrBlock,
			EnableIPv6:    true,
			IPv6CidrBlock: s.scope.VPC().IPv6CidrBlock,
			IPv6Pool:      s.scope.VPC().IPv6Pool,
			Tags:          converters.TagsToMap(out.Vpc.Tags),
		}, nil
	}

	// We have to describe the VPC again because the `create` output will **NOT** contain the associated IPv6 address.
	vpc, err := s.EC2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: aws.StringSlice([]string{aws.StringValue(out.Vpc.VpcId)}),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "DescribeVpcs", "Failed to describe the new ipv6 vpc: %v", err)
		return nil, errors.Wrap(err, "failed to describe new ipv6 vpc")
	}
	if len(vpc.Vpcs) == 0 {
		record.Warnf(s.scope.InfraCluster(), "DescribeVpcs", "Failed to find the new ipv6 vpc, returned list was empty.")
		return nil, errors.New("failed to find new ipv6 vpc; returned list was empty")
	}
	for _, set := range vpc.Vpcs[0].Ipv6CidrBlockAssociationSet {
		if *set.Ipv6CidrBlockState.State == ec2.SubnetCidrBlockStateCodeAssociated {
			return &infrav1.VPCSpec{
				EnableIPv6:    true,
				ID:            *vpc.Vpcs[0].VpcId,
				CidrBlock:     *out.Vpc.CidrBlock,
				IPv6CidrBlock: aws.StringValue(set.Ipv6CidrBlock),
				IPv6Pool:      aws.StringValue(set.Ipv6Pool),
				Tags:          converters.TagsToMap(vpc.Vpcs[0].Tags),
			}, nil
		}
	}

	return nil, fmt.Errorf("no IPv6 associated CIDR block sets found for IPv6 enabled cluster with vpc id %s", *out.Vpc.VpcId)
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

	if _, err := s.EC2Client.DeleteVpc(input); err != nil {
		// Ignore if it's already deleted
		if code, ok := awserrors.Code(err); ok && code == awserrors.VPCNotFound {
			s.scope.V(4).Info("Skipping VPC deletion, VPC not found")
			return nil
		}

		// Ignore if VPC ID is not present,
		if code, ok := awserrors.Code(err); ok && code == awserrors.VPCMissingParameter {
			s.scope.V(4).Info("Skipping VPC deletion, VPC ID not present")
			return nil
		}

		record.Warnf(s.scope.InfraCluster(), "FailedDeleteVPC", "Failed to delete managed VPC %q: %v", vpc.ID, err)
		return errors.Wrapf(err, "failed to delete vpc %q", vpc.ID)
	}

	s.scope.Info("Deleted VPC", "vpc-id", vpc.ID)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteVPC", "Deleted managed VPC %q", vpc.ID)
	return nil
}

func (s *Service) describeVPCByID() (*infrav1.VPCSpec, error) {
	if s.scope.VPC().ID == "" {
		return nil, errors.New("VPC ID is not set, failed to describe VPCs by ID")
	}

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPCStates(ec2.VpcStatePending, ec2.VpcStateAvailable),
		},
	}

	input.VpcIds = []*string{aws.String(s.scope.VPC().ID)}

	out, err := s.EC2Client.DescribeVpcs(input)
	if err != nil {
		if awserrors.IsNotFound(err) {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed to query ec2 for VPCs")
	}

	if len(out.Vpcs) == 0 {
		return nil, awserrors.NewNotFound(fmt.Sprintf("could not find vpc %q", s.scope.VPC().ID))
	} else if len(out.Vpcs) > 1 {
		return nil, awserrors.NewConflict(fmt.Sprintf("found %v VPCs with matching tags for %v. Only one VPC per cluster name is supported. Ensure duplicate VPCs are deleted for this AWS account and there are no conflicting instances of Cluster API Provider AWS. filtered VPCs: %v", len(out.Vpcs), s.scope.Name(), out.GoString()))
	}

	switch *out.Vpcs[0].State {
	case ec2.VpcStateAvailable, ec2.VpcStatePending:
	default:
		return nil, awserrors.NewNotFound("could not find available or pending vpc")
	}

	vpc := &infrav1.VPCSpec{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
		Tags:      converters.TagsToMap(out.Vpcs[0].Tags),
	}
	for _, set := range out.Vpcs[0].Ipv6CidrBlockAssociationSet {
		if *set.Ipv6CidrBlockState.State == ec2.SubnetCidrBlockStateCodeAssociated {
			vpc.IPv6CidrBlock = aws.StringValue(set.Ipv6CidrBlock)
			vpc.IPv6Pool = aws.StringValue(set.Ipv6Pool)
			vpc.EnableIPv6 = true
			break
		}
	}
	return vpc, nil
}

func (s *Service) getVPCTagParams(id string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-vpc", s.scope.Name())

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}
