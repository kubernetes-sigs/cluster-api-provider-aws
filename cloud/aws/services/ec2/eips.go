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
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/wait"
)

func (s *Service) getOrAllocateAddress(clusterName string, role string) (string, error) {
	out, err := s.describeAddresses(clusterName, role)
	if err != nil {
		return "", errors.Wrap(err, "failed to query addresses")
	}

	// TODO: better handle multiple addresses returned
	for _, address := range out.Addresses {
		if address.AssociationId == nil {
			return aws.StringValue(address.AllocationId), nil
		}
	}
	return s.allocateAddress(clusterName, role)
}

func (s *Service) allocateAddress(clusterName string, role string) (string, error) {
	out, err := s.EC2.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})

	if err != nil {
		return "", errors.Wrap(err, "failed to create Elastic IP address")
	}

	name := fmt.Sprintf("%s-eip-%s", clusterName, role)
	if err := s.createTags(clusterName, *out.AllocationId, ResourceLifecycleOwned, name, role, nil); err != nil {
		return "", errors.Wrapf(err, "failed to tag elastic IP %q", aws.StringValue(out.AllocationId))
	}

	return aws.StringValue(out.AllocationId), nil
}

func (s *Service) describeAddresses(clusterName string, role string) (*ec2.DescribeAddressesOutput, error) {
	filters := []*ec2.Filter{s.filterCluster(clusterName)}
	if role != "" {
		filters = append(filters, s.filterAWSProviderRole(role))
	}

	return s.EC2.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: filters,
	})
}

func (s *Service) releaseAddresses(clusterName string) error {
	out, err := s.EC2.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{s.filterCluster(clusterName)},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to describe elastic IPs %q", err)
	}

	for _, ip := range out.Addresses {
		if ip.AssociationId != nil {
			return errors.Errorf("failed to release elastic IP %q with allocation ID %q: Still associated with association ID %q", *ip.PublicIp, *ip.AllocationId, *ip.AssociationId)
		}

		releaseAddressInput := &ec2.ReleaseAddressInput{
			AllocationId: ip.AllocationId,
		}

		delete := func() (bool, error) {
			_, err := s.EC2.ReleaseAddress(releaseAddressInput)
			if err != nil {
				return false, err
			}

			return true, nil
		}

		retryableErrors := []string{
			errorAuthFailure,
			errorInUseIPAddress,
		}

		err := wait.WaitForWithRetryable(wait.NewBackoff(), delete, retryableErrors)
		if err != nil {
			return errors.Wrapf(err, "failed to release ElasticIP %q", *ip.AllocationId)
		}

		glog.Infof("released ElasticIP %q with allocation ID %q", *ip.PublicIp, *ip.AllocationId)
	}
	return nil
}
