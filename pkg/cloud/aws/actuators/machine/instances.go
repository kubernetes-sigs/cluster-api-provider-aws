package machine

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	awsclient "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/client"
)

// Scan machine tags, and return a deduped tags list
func removeDuplicatedTags(tags []*ec2.Tag) []*ec2.Tag {
	m := make(map[string]bool)
	result := []*ec2.Tag{}

	// look for duplicates
	for _, entry := range tags {
		if _, value := m[*entry.Key]; !value {
			m[*entry.Key] = true
			result = append(result, entry)
		}
	}
	return result
}

// removeStoppedMachine removes all instances of a specific machine that are in a stopped state.
func removeStoppedMachine(machine *clusterv1.Machine, client awsclient.Client) error {
	instances, err := GetStoppedInstances(machine, client)
	if err != nil {
		glog.Errorf("error getting stopped instances: %v", err)
		return fmt.Errorf("error getting stopped instances: %v", err)
	}

	if len(instances) == 0 {
		glog.Infof("no stopped instances found for machine %v", machine.Name)
		return nil
	}

	return TerminateInstances(client, instances)
}

func buildEC2Filters(inputFilters []providerconfigv1.Filter) []*ec2.Filter {
	filters := make([]*ec2.Filter, len(inputFilters))
	for i, f := range inputFilters {
		values := make([]*string, len(f.Values))
		for j, v := range f.Values {
			values[j] = aws.String(v)
		}
		filters[i] = &ec2.Filter{
			Name:   aws.String(f.Name),
			Values: values,
		}
	}
	return filters
}

func getSecurityGroupsIDs(securityGroups []providerconfigv1.AWSResourceReference, client awsclient.Client) ([]*string, error) {
	var securityGroupIDs []*string
	for _, g := range securityGroups {
		// ID has priority
		if g.ID != nil {
			securityGroupIDs = append(securityGroupIDs, g.ID)
		} else if g.Filters != nil {
			glog.Info("Describing security groups based on filters")
			// Get groups based on filters
			describeSecurityGroupsRequest := ec2.DescribeSecurityGroupsInput{
				Filters: buildEC2Filters(g.Filters),
			}
			describeSecurityGroupsResult, err := client.DescribeSecurityGroups(&describeSecurityGroupsRequest)
			if err != nil {
				glog.Errorf("error describing security groups: %v", err)
				return nil, fmt.Errorf("error describing security groups: %v", err)
			}
			for _, g := range describeSecurityGroupsResult.SecurityGroups {
				groupID := *g.GroupId
				securityGroupIDs = append(securityGroupIDs, &groupID)
			}
		}
	}

	if len(securityGroups) == 0 {
		glog.Info("No security group found")
	}

	return securityGroupIDs, nil
}

func getSubnetIDs(subnet providerconfigv1.AWSResourceReference, availabilityZone string, client awsclient.Client) ([]*string, error) {
	var subnetIDs []*string
	// ID has priority
	if subnet.ID != nil {
		subnetIDs = append(subnetIDs, subnet.ID)
	} else {
		var filters []providerconfigv1.Filter
		if availabilityZone != "" {
			filters = append(filters, providerconfigv1.Filter{Name: "availabilityZone", Values: []string{availabilityZone}})
		}
		filters = append(filters, subnet.Filters...)
		glog.Info("Describing subnets based on filters")
		describeSubnetRequest := ec2.DescribeSubnetsInput{
			Filters: buildEC2Filters(filters),
		}
		describeSubnetResult, err := client.DescribeSubnets(&describeSubnetRequest)
		if err != nil {
			glog.Errorf("error describing subnetes: %v", err)
			return nil, fmt.Errorf("error describing subnets: %v", err)
		}
		for _, n := range describeSubnetResult.Subnets {
			subnetID := *n.SubnetId
			subnetIDs = append(subnetIDs, &subnetID)
		}
	}
	if len(subnetIDs) == 0 {
		return nil, fmt.Errorf("no subnet IDs were found")
	}
	return subnetIDs, nil
}

func getAMI(AMI providerconfigv1.AWSResourceReference, client awsclient.Client) (*string, error) {
	if AMI.ID != nil {
		amiID := AMI.ID
		glog.Infof("Using AMI %s", *amiID)
		return amiID, nil
	}
	if len(AMI.Filters) > 0 {
		glog.Info("Describing AMI based on filters")
		describeImagesRequest := ec2.DescribeImagesInput{
			Filters: buildEC2Filters(AMI.Filters),
		}
		describeAMIResult, err := client.DescribeImages(&describeImagesRequest)
		if err != nil {
			glog.Errorf("error describing AMI: %v", err)
			return nil, fmt.Errorf("error describing AMI: %v", err)
		}
		if len(describeAMIResult.Images) < 1 {
			glog.Errorf("no image for given filters not found")
			return nil, fmt.Errorf("no image for given filters not found")
		}
		latestImage := describeAMIResult.Images[0]
		latestTime, err := time.Parse(time.RFC3339, *latestImage.CreationDate)
		if err != nil {
			glog.Errorf("unable to parse time for %q AMI: %v", *latestImage.ImageId, err)
			return nil, fmt.Errorf("unable to parse time for %q AMI: %v", *latestImage.ImageId, err)
		}
		for _, image := range describeAMIResult.Images[1:] {
			imageTime, err := time.Parse(time.RFC3339, *image.CreationDate)
			if err != nil {
				glog.Errorf("unable to parse time for %q AMI: %v", *image.ImageId, err)
				return nil, fmt.Errorf("unable to parse time for %q AMI: %v", *image.ImageId, err)
			}
			if latestTime.Before(imageTime) {
				latestImage = image
				latestTime = imageTime
			}
		}
		return latestImage.ImageId, nil
	}
	return nil, fmt.Errorf("AMI ID or AMI filters need to be specified")
}
