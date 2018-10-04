package ec2

import (
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

// Add additional cluster tag filters, to match on our tags
func (s *Service) addFilterTags(clusterName string, filters []*ec2.Filter) []*ec2.Filter {
	filters = append(filters, s.filterCluster(clusterName))
	return filters
}
