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
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/filter"
)

const (
	defaultRegion = "us-east-1"
)

func (s *Service) getRegion() string {
	switch x := s.scope.EC2.(type) {
	case *ec2.EC2:
		return *x.Config.Region
	default:
		return defaultRegion
	}
}

func (s *Service) getAvailableZones() ([]string, error) {
	out, err := s.scope.EC2.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{filter.EC2.Available()},
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to describe availability zones")
	}

	zones := make([]string, 0, len(out.AvailabilityZones))
	for _, zone := range out.AvailabilityZones {
		zones = append(zones, *zone.ZoneName)
	}

	sort.Strings(zones)
	return zones, nil
}
