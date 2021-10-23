package filter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

// ASG exposes the autoscaling sdk related filters.
var ASG = new(asgFilters)

type asgFilters struct{}

// Name returns a filter based on the resource name.
func (asgFilters) Name(name string) *autoscaling.Filter {
	return &autoscaling.Filter{
		Name:   aws.String("tag:Name"),
		Values: aws.StringSlice([]string{name}),
	}
}

// ClusterOwned returns a filter using the Cluster API per-cluster tag where
// the resource is owned.
func (asgFilters) ClusterOwned(clusterName string) *autoscaling.Filter {
	return &autoscaling.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.ClusterTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
	}
}

// ProviderRole returns a filter using cluster-api-provider-aws role tag.
func (asgFilters) ProviderRole(role string) *autoscaling.Filter {
	return &autoscaling.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.NameAWSClusterAPIRole)),
		Values: aws.StringSlice([]string{role}),
	}
}

// ProviderOwned returns a filter using the cloud provider tag where the resource is owned.
func (asgFilters) ProviderOwned(clusterName string) *autoscaling.Filter {
	return &autoscaling.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", infrav1.ClusterAWSCloudProviderTagKey(clusterName))),
		Values: aws.StringSlice([]string{string(infrav1.ResourceLifecycleOwned)}),
	}
}
