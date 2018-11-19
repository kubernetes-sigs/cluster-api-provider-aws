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
)

const (
	filterNameTagKey        = "tag-key"
	filterNameVpcID         = "vpc-id"
	filterNameState         = "state"
	filterNameVpcAttachment = "attachment.vpc-id"
)

// Returns an EC2 filter using the Cluster API per-cluster tag
func (s *Service) filterCluster(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameTagKey),
		Values: aws.StringSlice([]string{s.clusterTagKey(clusterName)}),
	}
}

// Returns an EC2 filter for name
func (s *Service) filterName(name string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("tag:Name"),
		Values: aws.StringSlice([]string{name}),
	}
}

// Returns an EC2 filter using the Cluster API per-cluster tag where
// the resource is owned
func (s *Service) filterClusterOwned(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", s.clusterTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(ResourceLifecycleOwned)}),
	}
}

// Returns an EC2 filter using the Cluster API per-cluster tag where
// the resource is shared
func (s *Service) filterClusterShared(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", s.clusterTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(ResourceLifecycleShared)}),
	}
}

// Returns an EC2 filter using cluster-api-provider-aws managed tag
func (s *Service) filterAWSProviderManaged() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameTagKey),
		Values: aws.StringSlice([]string{TagNameAWSProviderManaged}),
	}
}

// Returns an EC2 filter using cluster-api-provider-aws role tag
func (s *Service) filterAWSProviderRole(role string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", TagNameAWSClusterAPIRole)),
		Values: aws.StringSlice([]string{role}),
	}
}

// Returns an EC2 filter for the specified VPC ID
func (s *Service) filterVpc(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcID),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// Returns an EC2 filter for the specified VPC ID
func (s *Service) filterVpcAttachment(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcAttachment),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// Returns an EC2 filter for the state to be available
func (s *Service) filterAvailable() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameState),
		Values: aws.StringSlice([]string{"available"}),
	}
}

func (s *Service) filterNATGatewayStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

func (s *Service) filterInstanceStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: aws.StringSlice(states),
	}
}

func (s *Service) filterVPCStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

func (s *Service) filterSubnetsStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

// Add additional cluster tag filters, to match on our tags
func (s *Service) addFilterTags(clusterName string, filters []*ec2.Filter) []*ec2.Filter {
	filters = append(filters, s.filterCluster(clusterName))
	return filters
}
