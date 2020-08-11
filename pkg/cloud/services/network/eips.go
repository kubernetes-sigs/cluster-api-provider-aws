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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

func (s *Service) getOrAllocateAddresses(num int, role string) (eips []string, err error) {
	out, err := s.describeAddresses(role)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeAddresses", "Failed to query addresses for role %q: %v", role, err)
		return nil, errors.Wrap(err, "failed to query addresses")
	}

	for _, address := range out.Addresses {
		if address.AssociationId == nil {
			eips = append(eips, aws.StringValue(address.AllocationId))
		}
	}

	for len(eips) < num {
		ip, err := s.allocateAddress(role)
		if err != nil {
			return nil, err
		}
		eips = append(eips, ip)
	}

	return eips, nil
}

func (s *Service) allocateAddress(role string) (string, error) {
	out, err := s.EC2Client.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAllocateEIP", "Failed to allocate Elastic IP for %q: %v", role, err)
		return "", errors.Wrap(err, "failed to allocate Elastic IP")
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		buildParams := s.getEIPTagParams(*out.AllocationId, role)
		tagsBuilder := tags.New(&buildParams, tags.WithEC2(s.EC2Client))
		if err := tagsBuilder.Apply(); err != nil {
			return false, err
		}
		return true, nil
	}, awserrors.EIPNotFound); err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedAllocateAddress", "Failed to tag elastic IP %q: %v", aws.StringValue(out.AllocationId), err)
		return "", errors.Wrapf(err, "failed to tag Elastic IP %q", aws.StringValue(out.AllocationId))
	}

	return aws.StringValue(out.AllocationId), nil
}

func (s *Service) describeAddresses(role string) (*ec2.DescribeAddressesOutput, error) {
	x := []*ec2.Filter{filter.EC2.Cluster(s.scope.Name())}
	if role != "" {
		x = append(x, filter.EC2.ProviderRole(role))
	}

	return s.EC2Client.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: x,
	})
}

func (s *Service) disassociateAddress(ip *ec2.Address) error {
	err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		_, err := s.EC2Client.DisassociateAddress(&ec2.DisassociateAddressInput{
			AssociationId: ip.AssociationId,
		})
		if err != nil {
			cause, _ := awserrors.Code(errors.Cause(err))
			if cause != awserrors.AssociationIDNotFound {
				return false, err
			}
		}
		return true, nil
	}, awserrors.AuthFailure)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedDisassociateEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
		return errors.Wrapf(err, "failed to disassociate Elastic IP %q", *ip.AllocationId)
	}
	return nil
}

func (s *Service) releaseAddresses() error {
	out, err := s.EC2Client.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{filter.EC2.Cluster(s.scope.Name())},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to describe elastic IPs %q", err)
	}

	for i := range out.Addresses {
		ip := out.Addresses[i]
		if ip.AssociationId != nil {
			_, err := s.EC2Client.DisassociateAddress(&ec2.DisassociateAddressInput{
				AssociationId: ip.AssociationId,
			})
			if err != nil {
				record.Warnf(s.scope.InfraCluster(), "FailedDisassociateEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
				return errors.Errorf("failed to disassociate Elastic IP %q with allocation ID %q: Still associated with association ID %q", *ip.PublicIp, *ip.AllocationId, *ip.AssociationId)
			}
		}

		err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
			_, err := s.EC2Client.ReleaseAddress(&ec2.ReleaseAddressInput{AllocationId: ip.AllocationId})
			if err != nil {
				if ip.AssociationId != nil {
					if s.disassociateAddress(ip) != nil {
						return false, err
					}
				}
				return false, err
			}

			return true, nil
		}, awserrors.AuthFailure, awserrors.InUseIPAddress)
		if err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedReleaseEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
			return errors.Wrapf(err, "failed to release ElasticIP %q", *ip.AllocationId)
		}

		s.scope.Info("released ElasticIP", "eip", *ip.PublicIp, "allocation-id", *ip.AllocationId)
	}
	return nil
}

func (s *Service) getEIPTagParams(allocationID, role string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-eip-%s", s.scope.Name(), role)

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  allocationID,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(role),
		Additional:  s.scope.AdditionalTags(),
	}
}
