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
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

// AllocateDedicatedHost allocates a new dedicated host based on the specification.
func (s *Service) AllocateDedicatedHost(ctx context.Context, spec *infrav1.DynamicHostAllocationSpec, availabilityZone string) (string, error) {
	s.scope.Debug("Allocating dedicated host", "instanceFamily", spec.InstanceFamily, "availabilityZone", availabilityZone)

	// Determine the instance type for the dedicated host
	instanceType := s.determineInstanceType(spec)

	// Set default quantity if not specified
	quantity := int32(1)
	if spec.Quantity != nil {
		quantity = *spec.Quantity
	}

	// Use the specified AZ or fall back to the provided one
	targetAZ := availabilityZone
	if spec.AvailabilityZone != nil {
		targetAZ = *spec.AvailabilityZone
	}

	input := &ec2.AllocateHostsInput{
		InstanceFamily:   aws.String(spec.InstanceFamily),
		AvailabilityZone: aws.String(targetAZ),
		Quantity:         aws.Int32(quantity),
	}

	// Set instance type if specified
	if instanceType != "" {
		input.InstanceType = aws.String(instanceType)
	}

	// Add tags if specified
	if len(spec.Tags) > 0 {
		var tagSpecs []types.TagSpecification
		var tags []types.Tag
		for key, value := range spec.Tags {
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

	s.scope.Info("Allocating dedicated host", "input", input)
	output, err := s.EC2Client.AllocateHosts(ctx, input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAllocateDedicatedHost", "Failed to allocate dedicated host: %v", err)
		return "", errors.Wrap(err, fmt.Sprintf("failed to allocate dedicated host: %+v", input))
	}

	if len(output.HostIds) == 0 {
		return "", errors.New("no dedicated host ID returned from allocation")
	}

	hostID := output.HostIds[0]
	s.scope.Info("Successfully allocated dedicated host", "hostID", hostID, "availabilityZone", targetAZ)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulAllocateDedicatedHost", "Allocated dedicated host %s in %s", hostID, targetAZ)

	return hostID, nil
}

// ReleaseDedicatedHost releases a dedicated host.
func (s *Service) ReleaseDedicatedHost(ctx context.Context, hostID string) error {
	s.scope.Debug("Releasing dedicated host", "hostID", hostID)

	input := &ec2.ReleaseHostsInput{
		HostIds: []string{hostID},
	}

	_, err := s.EC2Client.ReleaseHosts(ctx, input)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedReleaseDedicatedHost", "Failed to release dedicated host %s: %v", hostID, err)
		return errors.Wrap(err, "failed to release dedicated host")
	}

	s.scope.Info("Successfully released dedicated host", "hostID", hostID)
	record.Eventf(s.scope.InfraCluster(), "SuccessfulReleaseDedicatedHost", "Released dedicated host %s", hostID)

	return nil
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


// determineInstanceType determines the instance type for the dedicated host.
func (s *Service) determineInstanceType(spec *infrav1.DynamicHostAllocationSpec) string {
	if spec.InstanceType != nil {
		return *spec.InstanceType
	}

	// If no specific instance type is provided, let AWS determine based on instance family
	return ""
}

// convertToHostInfo converts an AWS Host to our DedicatedHostInfo struct.
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

// ValidateHostCompatibility validates that an instance type is compatible with a host.
func (s *Service) ValidateHostCompatibility(ctx context.Context, hostID string, instanceType string) error {
	hostInfo, err := s.DescribeDedicatedHost(ctx, hostID)
	if err != nil {
		return errors.Wrap(err, "failed to describe host for compatibility check")
	}

	// Extract instance family from instance type (e.g., "m5.large" -> "m5")
	instanceFamily := extractInstanceFamily(instanceType)

	if hostInfo.InstanceFamily != instanceFamily {
		return errors.Errorf("instance type %s (family %s) is not compatible with host %s (family %s)",
			instanceType, instanceFamily, hostID, hostInfo.InstanceFamily)
	}

	// Check if host has available capacity
	if hostInfo.AvailableCapacity <= 0 {
		return errors.Errorf("host %s has no available capacity", hostID)
	}

	return nil
}

// extractInstanceFamily extracts the instance family from an instance type.
func extractInstanceFamily(instanceType string) string {
	// Instance types follow the pattern: family + generation + size (e.g., "m5.large")
	parts := strings.Split(instanceType, ".")
	if len(parts) < 2 {
		return instanceType
	}

	// Extract just the family part (e.g., "m5" from "m5.large")
	family := parts[0]

	// Remove generation number to get the base family (e.g., "m5" -> "m")
	// Actually, for dedicated hosts, we typically want to keep the generation
	return family
}
