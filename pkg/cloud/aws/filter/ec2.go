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

package filter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

const (
	filterNameTagKey        = "tag-key"
	filterNameVpcID         = "vpc-id"
	filterNameState         = "state"
	filterNameVpcAttachment = "attachment.vpc-id"
)

var (
	EC2 = new(ec2Filters)
)

type ec2Filters struct{}

// Returns an ec2Filters filter using the Cluster API per-cluster tag
func (ec2Filters) Cluster(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameTagKey),
		Values: aws.StringSlice([]string{tags.ClusterKey(clusterName)}),
	}
}

// Returns an ec2Filters filter for name
func (ec2Filters) Name(name string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("tag:Name"),
		Values: aws.StringSlice([]string{name}),
	}
}

// Returns an ec2Filters filter using the Cluster API per-cluster tag where
// the resource is owned
func (ec2Filters) ClusterOwned(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", tags.ClusterKey(clusterName))),
		Values: aws.StringSlice([]string{string(tags.ResourceLifecycleOwned)}),
	}
}

// Returns an ec2Filters filter using the Cluster API per-cluster tag where
// the resource is shared
func (ec2Filters) ClusterShared(clusterName string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", tags.ClusterKey(clusterName))),
		Values: aws.StringSlice([]string{string(tags.ResourceLifecycleShared)}),
	}
}

// Returns an ec2Filters filter using cluster-api-provider-aws managed tag
func (ec2Filters) ProviderManaged() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameTagKey),
		Values: aws.StringSlice([]string{TagNameAWSProviderManaged}),
	}
}

// Returns an ec2Filters filter using cluster-api-provider-aws role tag
func (ec2Filters) ProviderRole(role string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", TagNameAWSClusterAPIRole)),
		Values: aws.StringSlice([]string{role}),
	}
}

// Returns an ec2Filters filter for the specified VPC ID
func (ec2Filters) Vpc(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcID),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// Returns an ec2Filters filter for the specified VPC ID
func (ec2Filters) VpcAttachment(vpcID string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameVpcAttachment),
		Values: aws.StringSlice([]string{vpcID}),
	}
}

// Returns an ec2Filters filter for the state to be available
func (ec2Filters) Available() *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String(filterNameState),
		Values: aws.StringSlice([]string{"available"}),
	}
}

func (ec2Filters) NATGatewayStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

func (ec2Filters) InstanceStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: aws.StringSlice(states),
	}
}

func (ec2Filters) VPCStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

func (ec2Filters) SubnetsStates(states ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   aws.String("state"),
		Values: aws.StringSlice(states),
	}
}

// // Add additional cluster tag filters, to match on our tags
// func (ec2Filters) addFilterTags(clusterName string, filters []*ec2.Filter) []*ec2.Filter {
// 	filters = append(filters, filters.EC2.Cluster(clusterName))
// 	return filters
// }
