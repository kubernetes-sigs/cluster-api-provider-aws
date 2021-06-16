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

	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const (
	defaultVPCCidr = "10.0.0.0/16"
)

func (s *Service) reconcileVPC() error {
	s.scope.V(2).Info("Reconciling VPC")

	vpc, err := s.describeVPC()
	if awserrors.IsNotFound(err) { // nolint:nestif
		// Create a new managed vpc.
		if !conditions.Has(s.scope.InfraCluster(), infrav1.VpcReadyCondition) {
			conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcCreationStartedReason, clusterv1.ConditionSeverityInfo, "")
			if err := s.scope.PatchObject(); err != nil {
				return errors.Wrap(err, "failed to patch conditions")
			}
		}
		vpc, err = s.createVPC()
		if err != nil {
			return errors.Wrap(err, "failed to create new vpc")
		}
	} else if err != nil {
		return errors.Wrap(err, "failed to describe VPCs")
	}

	// This function creates a new infrav1.VPCSpec, populates it with data from AWS, and then deep copies into the
	// AWSCluster's VPC spec (see the DeepCopyInto lines below). This is potentially problematic, as it completely
	// overwrites the data for the VPC spec as retrieved from the apiserver. This is a temporary band-aid to restore
	// recently-added fields that descripe user intent and do not come from AWS resource descriptions.
	//
	// FIXME(ncdc): rather than copying these values from the scope to vpc, find a better way to merge AWS information
	// with data in the scope retrieved from the apiserver. Could use something like mergo.
	//
	// NOTE: it may look like we are losing InternetGatewayID because it's not populated by describeVPC/createVPC or
	// restored here, but that's ok. It is restored by reconcileInternetGateways, which is invoked after this.
	vpc.AvailabilityZoneSelection = s.scope.VPC().AvailabilityZoneSelection
	vpc.AvailabilityZoneUsageLimit = s.scope.VPC().AvailabilityZoneUsageLimit

	if vpc.IsUnmanaged(s.scope.Name()) {
		vpc.DeepCopyInto(s.scope.VPC())
		s.scope.V(2).Info("Working on unmanaged VPC", "vpc-id", vpc.ID)
		return nil
	}

	// Make sure attributes are configured
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getVPCTagParams(vpc.ID)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Ensure(vpc.Tags); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.VPCNotFound); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedTagVPC", "Failed to tag managed VPC %q: %v", vpc.ID, err)
		return errors.Wrapf(err, "failed to tag vpc %q", vpc.ID)
	}

	// Make sure attributes are configured
	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		if err := s.ensureManagedVPCAttributes(vpc); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.VPCNotFound); err != nil {
		return errors.Wrapf(err, "failed to to set vpc attributes for %q", vpc.ID)
	}

	vpc.DeepCopyInto(s.scope.VPC())
	s.scope.V(2).Info("Working on managed VPC", "vpc-id", vpc.ID)
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
	if s.scope.VPC().IsUnmanaged(s.scope.Name()) {
		return nil, errors.Errorf("cannot create a managed vpc in unmanaged mode")
	}

	if s.scope.VPC().CidrBlock == "" {
		s.scope.VPC().CidrBlock = defaultVPCCidr
	}

	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(s.scope.VPC().CidrBlock),
		TagSpecifications: []*ec2.TagSpecification{
			tags.BuildParamsToTagSpecification(ec2.ResourceTypeVpc, s.getVPCTagParams(services.TemporaryResourceID)),
		},
	}

	out, err := s.EC2Client.CreateVpc(input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedCreateVPC", "Failed to create new managed VPC: %v", err)
		return nil, errors.Wrap(err, "failed to create vpc")
	}

	record.Eventf(s.scope.InfraCluster(), "SuccessfulCreateVPC", "Created new managed VPC %q", *out.Vpc.VpcId)
	s.scope.V(2).Info("Created new VPC with cidr", "vpc-id", *out.Vpc.VpcId, "cidr-block", *out.Vpc.CidrBlock)

	// TODO: we should attempt to record the VPC ID as soon as possible by setting s.scope.VPC().ID
	// however, the logic used for determining managed vs unmanaged VPCs relies on the tags and will
	// need to be updated to accommodate for the recording of the VPC ID prior to the tagging.

	wReq := &ec2.DescribeVpcsInput{VpcIds: []*string{out.Vpc.VpcId}}
	if err := s.EC2Client.WaitUntilVpcAvailable(wReq); err != nil {
		return nil, errors.Wrapf(err, "failed to wait for vpc %q", *out.Vpc.VpcId)
	}

	return &infrav1.VPCSpec{
		ID:        *out.Vpc.VpcId,
		CidrBlock: *out.Vpc.CidrBlock,
		Tags:      converters.TagsToMap(out.Vpc.Tags),
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

	if _, err := s.EC2Client.DeleteVpc(input); err != nil {
		// Ignore if it's already deleted
		if code, ok := awserrors.Code(err); ok && code == awserrors.VPCNotFound {
			s.scope.V(4).Info("Skipping VPC deletion, VPC not found")
			return nil
		}
		record.Warnf(s.scope.InfraCluster(), "FailedDeleteVPC", "Failed to delete managed VPC %q: %v", vpc.ID, err)
		return errors.Wrapf(err, "failed to delete vpc %q", vpc.ID)
	}

	s.scope.V(2).Info("Deleted VPC", "vpc-id", vpc.ID)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulDeleteVPC", "Deleted managed VPC %q", vpc.ID)
	return nil
}

func (s *Service) describeVPC() (*infrav1.VPCSpec, error) {
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

	return &infrav1.VPCSpec{
		ID:        *out.Vpcs[0].VpcId,
		CidrBlock: *out.Vpcs[0].CidrBlock,
		Tags:      converters.TagsToMap(out.Vpcs[0].Tags),
	}, nil
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
