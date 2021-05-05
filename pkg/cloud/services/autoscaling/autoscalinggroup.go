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

package asg

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

// SDKToAutoScalingGroup converts an AWS EC2 SDK AutoScalingGroup to the CAPA AutoScalingGroup type.
func (s *Service) SDKToAutoScalingGroup(v *autoscaling.Group) (*expinfrav1.AutoScalingGroup, error) {
	i := &expinfrav1.AutoScalingGroup{
		ID:   aws.StringValue(v.AutoScalingGroupARN),
		Name: aws.StringValue(v.AutoScalingGroupName),
		// TODO(rudoi): this is just terrible
		DesiredCapacity:   aws.Int32(int32(aws.Int64Value(v.DesiredCapacity))),
		MaxSize:           int32(aws.Int64Value(v.MaxSize)),
		MinSize:           int32(aws.Int64Value(v.MinSize)),
		CapacityRebalance: aws.BoolValue(v.CapacityRebalance),
		//TODO: determine what additional values go here and what else should be in the struct
	}

	if v.MixedInstancesPolicy != nil {
		i.MixedInstancesPolicy = &expinfrav1.MixedInstancesPolicy{
			InstancesDistribution: &expinfrav1.InstancesDistribution{
				OnDemandBaseCapacity:                v.MixedInstancesPolicy.InstancesDistribution.OnDemandBaseCapacity,
				OnDemandPercentageAboveBaseCapacity: v.MixedInstancesPolicy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity,
			},
		}

		for _, override := range v.MixedInstancesPolicy.LaunchTemplate.Overrides {
			i.MixedInstancesPolicy.Overrides = append(i.MixedInstancesPolicy.Overrides, expinfrav1.Overrides{InstanceType: aws.StringValue(override.InstanceType)})
		}

		onDemandAllocationStrategy := aws.StringValue(v.MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy)
		if onDemandAllocationStrategy == string(expinfrav1.OnDemandAllocationStrategyPrioritized) {
			i.MixedInstancesPolicy.InstancesDistribution.OnDemandAllocationStrategy = expinfrav1.OnDemandAllocationStrategyPrioritized
		}

		spotAllocationStrategy := aws.StringValue(v.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy)
		if spotAllocationStrategy == string(expinfrav1.SpotAllocationStrategyLowestPrice) {
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyLowestPrice
		} else {
			i.MixedInstancesPolicy.InstancesDistribution.SpotAllocationStrategy = expinfrav1.SpotAllocationStrategyCapacityOptimized
		}
	}

	if v.Status != nil {
		i.Status = expinfrav1.ASGStatus(*v.Status)
	}

	if len(v.Tags) > 0 {
		i.Tags = converters.ASGTagsToMap(v.Tags)
	}

	if len(v.Instances) > 0 {
		for _, autoscalingInstance := range v.Instances {
			tmp := &infrav1.Instance{
				ID: aws.StringValue(autoscalingInstance.InstanceId),
			}
			i.Instances = append(i.Instances, *tmp)
		}
	}

	return i, nil
}

// ASGIfExists returns the existing autoscaling group or nothing if it doesn't exist.
func (s *Service) ASGIfExists(name *string) (*expinfrav1.AutoScalingGroup, error) {
	if name == nil {
		s.scope.Info("Autoscaling Group does not have a name")
		return nil, nil
	}

	s.scope.Info("Looking for asg by name", "name", *name)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{name},
	}

	out, err := s.ASGClient.DescribeAutoScalingGroups(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeAutoScalingGroups", "failed to describe ASG %q: %v", *name, err)
		return nil, errors.Wrapf(err, "failed to describe AutoScaling Group: %q", *name)
	}
	//TODO: double check if you're handling nil vals
	return s.SDKToAutoScalingGroup(out.AutoScalingGroups[0])

}

// GetASGByName returns the existing ASG or nothing if it doesn't exist.
func (s *Service) GetASGByName(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
	s.scope.V(2).Info("Looking for existing AutoScalingGroup by name")

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(scope.Name()),
		},
	}

	out, err := s.ASGClient.DescribeAutoScalingGroups(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeInstances", "Failed to describe instances by tags: %v", err)
		return nil, errors.Wrap(err, "failed to describe instances by tags")
	case len(out.AutoScalingGroups) == 0:
		record.Eventf(scope.AWSMachinePool, "FailedDescribeInstances", "No Auto Scaling Groups with %s found", scope.Name())
		return nil, nil
	}

	return s.SDKToAutoScalingGroup(out.AutoScalingGroups[0])
}

// CreateASG runs an autoscaling group.
func (s *Service) CreateASG(scope *scope.MachinePoolScope) (*expinfrav1.AutoScalingGroup, error) {
	subnetIDs := make([]string, len(scope.AWSMachinePool.Spec.Subnets))
	for i, v := range scope.AWSMachinePool.Spec.Subnets {
		subnetIDs[i] = aws.StringValue(v.ID)
	}
	// subnetIDs := []string{}
	// for _, v := range scope.AWSMachinePool.Spec.Subnets {
	// 	subnetIDs = append(subnetIDs, aws.StringValue(v.ID))
	// }

	input := &expinfrav1.AutoScalingGroup{
		Name:                 scope.Name(),
		MaxSize:              scope.AWSMachinePool.Spec.MaxSize,
		MinSize:              scope.AWSMachinePool.Spec.MinSize,
		Subnets:              subnetIDs,
		DefaultCoolDown:      scope.AWSMachinePool.Spec.DefaultCoolDown,
		CapacityRebalance:    scope.AWSMachinePool.Spec.CapacityRebalance,
		MixedInstancesPolicy: scope.AWSMachinePool.Spec.MixedInstancesPolicy,
	}

	if scope.MachinePool.Spec.Replicas != nil && !scope.MachinePool.Spec.ExternallyManagedReplicaCount {
		// only set DesiredCapacity if CAPI is managing the Replica Count
		input.DesiredCapacity = scope.MachinePool.Spec.Replicas
	}

	if scope.AWSMachinePool.Status.LaunchTemplateID == "" {
		return nil, errors.New("AWSMachinePool has no LaunchTemplateID for some reason")
	}

	// Make sure to use the MachinePoolScope here to get the merger of AWSCluster and AWSMachinePool tags
	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	input.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	s.scope.Info("Running instance")
	if err := s.runPool(input, scope.AWSMachinePool.Status.LaunchTemplateID); err != nil {
		// Only record the failure event if the error is not related to failed dependencies.
		// This is to avoid spamming failure events since the machine will be requeued by the actuator.
		// if !awserrors.IsFailedDependency(errors.Cause(err)) {
		// 	record.Warnf(scope.AWSMachinePool, "FailedCreate", "Failed to create instance: %v", err)
		// }
		s.scope.Error(err, "unable to create AutoScalingGroup")
		return nil, err
	}
	record.Eventf(scope.AWSMachinePool, "SuccessfulCreate", "Created new ASG: %s", scope.Name)

	return nil, nil
}

func (s *Service) runPool(i *expinfrav1.AutoScalingGroup, launchTemplateID string) error {
	input := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(i.Name),
		MaxSize:              aws.Int64(int64(i.MaxSize)),
		MinSize:              aws.Int64(int64(i.MinSize)),
		VPCZoneIdentifier:    aws.String(strings.Join(i.Subnets, ", ")),
		DefaultCooldown:      aws.Int64(int64(i.DefaultCoolDown.Duration.Seconds())),
		CapacityRebalance:    aws.Bool(i.CapacityRebalance),
	}

	if i.DesiredCapacity != nil {
		input.DesiredCapacity = aws.Int64(int64(aws.Int32Value(i.DesiredCapacity)))
	}

	if i.MixedInstancesPolicy != nil {
		input.MixedInstancesPolicy = createSDKMixedInstancesPolicy(i.Name, i.MixedInstancesPolicy)
	} else {
		input.LaunchTemplate = &autoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(launchTemplateID),
			Version:          aws.String(expinfrav1.LaunchTemplateLatestVersion),
		}
	}

	if i.Tags != nil {
		input.Tags = BuildTagsFromMap(i.Name, i.Tags)
	}

	_, err := s.ASGClient.CreateAutoScalingGroup(input)
	if err != nil {
		return errors.Wrap(err, "failed to create autoscaling group")
	}

	return nil
}

func (s *Service) DeleteASGAndWait(name string) error {
	if err := s.DeleteASG(name); err != nil {
		return err
	}

	s.scope.V(2).Info("Waiting for ASG to be deleted", "name", name)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: aws.StringSlice([]string{name}),
	}

	if err := s.ASGClient.WaitUntilGroupNotExists(input); err != nil {
		return errors.Wrapf(err, "failed to wait for ASG %q deletion", name)
	}

	return nil
}

func (s *Service) DeleteASG(name string) error {
	s.scope.V(2).Info("Attempting to delete ASG", "name", name)

	input := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		ForceDelete:          aws.Bool(true),
	}

	if _, err := s.ASGClient.DeleteAutoScalingGroup(input); err != nil {
		return errors.Wrapf(err, "failed to delete ASG %q", name)
	}

	s.scope.V(2).Info("Deleted ASG", "name", name)
	return nil
}

func (s *Service) UpdateASG(scope *scope.MachinePoolScope) error {
	subnetIDs := make([]string, len(scope.AWSMachinePool.Spec.Subnets))
	for i, v := range scope.AWSMachinePool.Spec.Subnets {
		subnetIDs[i] = aws.StringValue(v.ID)
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(scope.Name()), //TODO: define dynamically - borrow logic from ec2
		MaxSize:              aws.Int64(int64(scope.AWSMachinePool.Spec.MaxSize)),
		MinSize:              aws.Int64(int64(scope.AWSMachinePool.Spec.MinSize)),
		VPCZoneIdentifier:    aws.String(strings.Join(subnetIDs, ", ")),
		CapacityRebalance:    aws.Bool(scope.AWSMachinePool.Spec.CapacityRebalance),
	}

	if scope.MachinePool.Spec.Replicas != nil && !scope.MachinePool.Spec.ExternallyManagedReplicaCount {
		input.DesiredCapacity = aws.Int64(int64(*scope.MachinePool.Spec.Replicas))
	}

	if scope.AWSMachinePool.Spec.MixedInstancesPolicy != nil {
		input.MixedInstancesPolicy = createSDKMixedInstancesPolicy(scope.Name(), scope.AWSMachinePool.Spec.MixedInstancesPolicy)
	} else {
		input.LaunchTemplate = &autoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(scope.AWSMachinePool.Status.LaunchTemplateID),
			Version:          aws.String(expinfrav1.LaunchTemplateLatestVersion),
		}
	}

	if _, err := s.ASGClient.UpdateAutoScalingGroup(input); err != nil {
		return errors.Wrapf(err, "failed to update ASG %q", scope.Name())
	}

	return nil
}

func createSDKMixedInstancesPolicy(name string, i *expinfrav1.MixedInstancesPolicy) *autoscaling.MixedInstancesPolicy {
	mixedInstancesPolicy := &autoscaling.MixedInstancesPolicy{
		LaunchTemplate: &autoscaling.LaunchTemplate{
			LaunchTemplateSpecification: &autoscaling.LaunchTemplateSpecification{
				LaunchTemplateName: aws.String(name),
				Version:            aws.String(expinfrav1.LaunchTemplateLatestVersion),
			},
		},
	}

	if i.InstancesDistribution != nil {
		mixedInstancesPolicy.InstancesDistribution = &autoscaling.InstancesDistribution{
			OnDemandAllocationStrategy:          aws.String(string(i.InstancesDistribution.OnDemandAllocationStrategy)),
			OnDemandBaseCapacity:                i.InstancesDistribution.OnDemandBaseCapacity,
			OnDemandPercentageAboveBaseCapacity: i.InstancesDistribution.OnDemandPercentageAboveBaseCapacity,
			SpotAllocationStrategy:              aws.String(string(i.InstancesDistribution.SpotAllocationStrategy)),
		}
	}

	for _, override := range i.Overrides {
		mixedInstancesPolicy.LaunchTemplate.Overrides = append(mixedInstancesPolicy.LaunchTemplate.Overrides, &autoscaling.LaunchTemplateOverrides{
			InstanceType: aws.String(override.InstanceType),
		})
	}

	return mixedInstancesPolicy
}

// BuildTags takes the tag configuration from the resources and returns a slice of autoscaling Tags
// usable in autoscaling API calls
func BuildTags(name string, params infrav1.BuildParams) []*autoscaling.Tag {
	tags := make([]*autoscaling.Tag, 0)
	resourceName := aws.String(name)
	propagateAtLaunch := aws.Bool(false)
	resourceType := aws.String("auto-scaling-group")
	if params.Additional != nil {
		for k, v := range params.Additional {
			tags = append(tags, &autoscaling.Tag{
				Key:   aws.String(k),
				Value: aws.String(v),
				// We set the instance tags in the LaunchTemplate, disabling propagation to prevent the two
				// resources from clobbering each other's tags
				PropagateAtLaunch: propagateAtLaunch,
				ResourceId:        resourceName,
				ResourceType:      resourceType,
			})
		}
	}

	tags = append(tags, &autoscaling.Tag{
		Key:               aws.String(infrav1.ClusterTagKey(params.ClusterName)),
		Value:             aws.String(string(params.Lifecycle)),
		PropagateAtLaunch: propagateAtLaunch,
		ResourceId:        resourceName,
		ResourceType:      resourceType,
	})

	if params.Role != nil {
		tags = append(tags, &autoscaling.Tag{
			Key:               aws.String(infrav1.NameAWSClusterAPIRole),
			Value:             params.Role,
			PropagateAtLaunch: propagateAtLaunch,
			ResourceId:        resourceName,
			ResourceType:      resourceType,
		})
	}

	if params.Name != nil {
		tags = append(tags, &autoscaling.Tag{
			Key:               aws.String("Name"),
			Value:             params.Name,
			PropagateAtLaunch: propagateAtLaunch,
			ResourceId:        resourceName,
			ResourceType:      resourceType,
		})
	}

	return tags
}

// BuildTagsFromMap takes a map of keys and values and returns them as autoscaling group tags
func BuildTagsFromMap(asgName string, inTags map[string]string) []*autoscaling.Tag {
	if inTags == nil {
		return nil
	}
	tags := make([]*autoscaling.Tag, 0)
	for k, v := range inTags {
		tags = append(tags, &autoscaling.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
			// We set the instance tags in the LaunchTemplate, disabling propagation to prevent the two
			// resources from clobbering the tags set in the LaunchTemplate
			PropagateAtLaunch: aws.Bool(false),
			ResourceId:        aws.String(asgName),
			ResourceType:      aws.String("auto-scaling-group"),
		})
	}

	return tags
}

// UpdateResourceTags updates the tags for an autoscaling group.
// This will be called if there is anything to create (update) or delete.
// We may not always have to perform each action, so we check what we're
// receiving to avoid calling AWS if we don't need to.
func (s *Service) UpdateResourceTags(resourceID *string, create, remove map[string]string) error {
	s.scope.V(2).Info("Attempting to update tags on resource", "resource-id", *resourceID)
	s.scope.Info("updating tags on resource", "resource-id", *resourceID, "create", create, "remove", remove)

	// If we have anything to create or update
	if len(create) > 0 {
		s.scope.V(2).Info("Attempting to create tags on resource", "resource-id", *resourceID)

		createOrUpdateTagsInput := &autoscaling.CreateOrUpdateTagsInput{}

		createOrUpdateTagsInput.Tags = mapToTags(create, resourceID)

		if _, err := s.ASGClient.CreateOrUpdateTags(createOrUpdateTagsInput); err != nil {
			return errors.Wrapf(err, "failed to update tags on AutoScalingGroup %q", *resourceID)
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		s.scope.V(2).Info("Attempting to delete tags on resource", "resource-id", *resourceID)

		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := mapToTags(remove, resourceID)

		// Create the DeleteTags input
		input := &autoscaling.DeleteTagsInput{
			Tags: removeTagsInput,
		}

		// Delete tags in AWS.
		if _, err := s.ASGClient.DeleteTags(input); err != nil {
			return errors.Wrapf(err, "failed to delete tags on AutoScalingGroup %q: %v", *resourceID, remove)
		}
	}

	return nil
}

func mapToTags(input map[string]string, resourceID *string) []*autoscaling.Tag {
	tags := make([]*autoscaling.Tag, 0)
	for k, v := range input {
		tags = append(tags, &autoscaling.Tag{
			Key:               aws.String(k),
			PropagateAtLaunch: aws.Bool(false),
			ResourceId:        resourceID,
			ResourceType:      aws.String("auto-scaling-group"),
			Value:             aws.String(v),
		})
	}
	return tags
}
