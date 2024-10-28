package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

// runInstancesForInputAllowed will set DryRun=true on the input, and test if the RunInstances call will succeed,
// returning a bool and any additional error. Note that the aws-sdk behavior is to return an error whether the DryRun works
// or not.
func (s *Service) runInstancesForInputAllowed(ctx context.Context, input *ec2.RunInstancesInput) (bool, error) {
	input.DryRun = ptr.To(true)
	_, err := s.EC2Client.RunInstancesWithContext(ctx, input)
	input.DryRun = nil

	// This is the success path, the API returns an error with the code 'DryRunOperation'
	if awserrors.IsDryRunOperationError(err) {
		s.scope.Debug("DryRun validation passed for RunInstances")
		return true, nil
	} else if awserrors.IsPermissionsError(err) {
		// This is the failure path, signifying we lack permissions to perform RunInstances with the configured input
		s.scope.Debug("DryRun validation failed for RunInstances")
		return false, nil
	}

	// Any other error scenario means failure
	return false, err
}

func dropNetworkInterfaceTags(input *ec2.RunInstancesInput) []*ec2.TagSpecification {
	var tagSpecifications []*ec2.TagSpecification
	for _, spec := range input.TagSpecifications {
		if *spec.ResourceType != ec2.ResourceTypeNetworkInterface {
			tagSpecifications = append(tagSpecifications, spec)
		}
	}

	return tagSpecifications
}
