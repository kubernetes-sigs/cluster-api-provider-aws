/*
Copyright 2025 The Kubernetes Authors.

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
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

// AllocateDedicatedHost allocates a single dedicated host based on the specification.
// This function always allocates exactly one dedicated host per call.
// The dedicated host will inherit additional tags defined in the AWSMachineTemplate.
func (s *Service) AllocateDedicatedHost(ctx context.Context, spec *infrav1.DynamicHostAllocationSpec, instanceType, availabilityZone string, scope *scope.MachineScope) (string, error) {
	s.scope.Debug("Allocating single dedicated host", "instanceType", instanceType, "availabilityZone", availabilityZone)
	input := &ec2.AllocateHostsInput{
		InstanceType:     aws.String(instanceType),
		AvailabilityZone: aws.String(availabilityZone),
		Quantity:         aws.Int32(1),
	}

	// Build tags for the dedicated host
	// Only include additionalTags from the machine and dedicated host specific tags
	additionalTags := scope.AdditionalTags()

	// Start with additional tags from the machine (AWSMachineTemplate additionalTags)
	dedicatedHostTags := make(map[string]string)
	for key, value := range additionalTags {
		dedicatedHostTags[key] = value
	}

	// Merge in dedicated host specific tags from the spec
	// Dedicated host specific tags take precedence over additional tags
	for key, value := range spec.Tags {
		dedicatedHostTags[key] = value
	}

	// Add tags to the allocation request
	if len(dedicatedHostTags) > 0 {
		var tagSpecs []types.TagSpecification
		var tags []types.Tag
		for key, value := range dedicatedHostTags {
			tags = append(tags, types.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpecs = append(tagSpecs, types.TagSpecification{
			ResourceType: types.ResourceTypeDedicatedHost,
			Tags:         tags,
		})
		input.TagSpecifications = tagSpecs
	}

	s.scope.Info("Allocating dedicated host", "input", input, "machine", scope.Name())
	output, err := s.EC2Client.AllocateHosts(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to allocate dedicated host: %+v", input))
	}

	// Ensure we got exactly one host as expected
	if len(output.HostIds) != 1 {
		return "", errors.Errorf("expected one dedicated host, but got %d hosts", len(output.HostIds))
	}

	hostID := output.HostIds[0]
	s.scope.Info("Successfully allocated single dedicated host",
		"hostID", hostID,
		"availabilityZone", availabilityZone,
		"machine", scope.Name(),
		"instanceType", instanceType)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulAllocateDedicatedHost", "Allocated dedicated host %s in %s for machine %s", hostID, availabilityZone, scope.Name())

	return hostID, nil
}

// ReleaseDedicatedHost releases a dedicated host with retry logic.
// This function implements custom retry logic for dedicated host release operations
// since they are expensive and should be retried on transient failures.
func (s *Service) ReleaseDedicatedHost(ctx context.Context, hostID string) error {
	s.scope.Debug("Releasing dedicated host", "hostID", hostID)

	input := &ec2.ReleaseHostsInput{
		HostIds: []string{hostID},
	}

	// Retry configuration optimized for dedicated host release operations
	// These are expensive resources, so we want to be more aggressive with retries
	maxAttempts := 5
	maxBackoff := 30 * time.Second

	// Retryable error codes for dedicated host operations
	retryableErrors := []string{
		"RequestLimitExceeded",
		"Throttling",
		"ServiceUnavailable",
		"InternalError",
		"InvalidHostState",
		"HostInUse",
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		output, err := s.EC2Client.ReleaseHosts(ctx, input)
		if err == nil {
			s.scope.Info("Successfully released dedicated host", "attempt", attempt, "output", "result", s.getReleaseHostsOutput(output))
			record.Eventf(s.scope.InfraCluster(), "SuccessfulReleaseDedicatedHost", "Released dedicated host %s", hostID)
			return nil
		}

		lastErr = err

		// Check if this is a retryable error
		errorCode := s.getErrorCode(err)
		isRetryable := false
		for _, retryableCode := range retryableErrors {
			if errorCode == retryableCode {
				isRetryable = true
				break
			}
		}

		if !isRetryable {
			s.scope.Error(err, "Non-retryable error releasing dedicated host", "errorCode", errorCode, "result", s.getReleaseHostsOutput(output))
			record.Warnf(s.scope.InfraCluster(), "FailedReleaseDedicatedHost", "Failed to release dedicated host %s: %v", hostID, err)
			return errors.Wrap(err, "failed to release dedicated host")
		}

		// If this is the last attempt, don't wait
		if attempt >= maxAttempts {
			break
		}

		// Calculate exponential backoff delay with jitter
		baseDelay := time.Duration(1<<(attempt-1)) * time.Second
		if baseDelay > maxBackoff {
			baseDelay = maxBackoff
		}
		// Add jitter (±25% of the delay)
		jitter := time.Duration(float64(baseDelay) * 0.25)
		delay := baseDelay + time.Duration(float64(jitter)*(2*rand.Float64()-1))

		s.scope.Info("Retrying dedicated host release after error",
			"hostID", hostID,
			"attempt", attempt,
			"maxAttempts", maxAttempts,
			"errorCode", errorCode,
			"delay", delay,
			"result", s.getReleaseHostsOutput(output))

		// Wait with context cancellation support
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	// All retries exhausted
	s.scope.Error(lastErr, "Failed to release dedicated host after all retries",
		"hostID", hostID,
		"maxAttempts", maxAttempts)
	record.Warnf(s.scope.InfraCluster(), "FailedReleaseDedicatedHost", "Failed to release dedicated host %s after %d attempts: %v", hostID, maxAttempts, lastErr)
	return errors.Wrap(lastErr, "failed to release dedicated host after all retries")
}

// getErrorCode extracts the error code from an AWS error.
func (s *Service) getErrorCode(err error) string {
	if smithyErr := awserrors.ParseSmithyError(err); smithyErr != nil {
		return smithyErr.ErrorCode()
	}
	if code, ok := awserrors.Code(err); ok {
		return code
	}
	return "Unknown"
}

// DescribeDedicatedHost describes a specific dedicated host.
func (s *Service) DescribeDedicatedHost(ctx context.Context, hostID string) (*infrav1.DedicatedHostInfo, error) {
	input := &ec2.DescribeHostsInput{
		HostIds: []string{hostID},
	}

	output, err := s.EC2Client.DescribeHosts(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe dedicated host")
	}

	if len(output.Hosts) == 0 {
		return nil, errors.Errorf("dedicated host %s not found", hostID)
	}

	host := output.Hosts[0]
	hostInfo := s.convertToHostInfo(host)

	return hostInfo, nil
}

// convertToHostInfo converts an AWS Host to the DedicatedHostInfo struct.
func (s *Service) convertToHostInfo(host types.Host) *infrav1.DedicatedHostInfo {
	hostInfo := &infrav1.DedicatedHostInfo{
		HostID:           aws.ToString(host.HostId),
		AvailabilityZone: aws.ToString(host.AvailabilityZone),
		State:            string(host.State),
		Tags:             make(map[string]string),
	}

	// Parse properties from HostProperties
	if host.HostProperties != nil {
		if host.HostProperties.InstanceFamily != nil {
			hostInfo.InstanceFamily = *host.HostProperties.InstanceFamily
		}
		if host.HostProperties.InstanceType != nil {
			hostInfo.InstanceType = *host.HostProperties.InstanceType
		}
		if host.HostProperties.TotalVCpus != nil {
			hostInfo.TotalCapacity = *host.HostProperties.TotalVCpus
		}
	}

	// Calculate available capacity from instances
	instanceCount := len(host.Instances)
	if instanceCount > math.MaxInt32 {
		instanceCount = math.MaxInt32
	}
	// bounds check ensures instanceCount <= math.MaxInt32, preventing integer overflow
	usedCapacity := int32(instanceCount) //nolint:gosec
	hostInfo.AvailableCapacity = hostInfo.TotalCapacity - usedCapacity

	// Convert tags
	for _, tag := range host.Tags {
		if tag.Key != nil && tag.Value != nil {
			hostInfo.Tags[*tag.Key] = *tag.Value
		}
	}

	return hostInfo
}

func (s *Service) getReleaseHostsOutput(output *ec2.ReleaseHostsOutput) string {
	var errs []string

	if output.Successful != nil {
		return strings.Join(output.Successful, ", ")
	} else if output.Unsuccessful != nil {
		for _, err := range output.Unsuccessful {
			errs = append(errs, fmt.Sprintf("Resource ID: %s, Error Code: %s, Error Message: %s", aws.ToString(err.ResourceId), aws.ToString(err.Error.Code), aws.ToString(err.Error.Message)))
		}
		return strings.Join(errs, ", ")
	}

	return ""
}
