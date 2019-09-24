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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	"sigs.k8s.io/cluster-api/util"
)

// GetRunningInstanceByTags returns the existing instance or nothing if it doesn't exist.
func (s *Service) GetRunningInstanceByTags(scope *scope.MachineScope) (*infrav1.Instance, error) {
	s.scope.V(2).Info("Looking for existing machine instance by tags")

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			filter.EC2.VPC(s.scope.VPC().ID),
			filter.EC2.ClusterOwned(s.scope.Name()),
			filter.EC2.Name(scope.Name()),
			filter.EC2.InstanceStates(ec2.InstanceStateNamePending, ec2.InstanceStateNameRunning),
		},
	}

	out, err := s.scope.EC2.DescribeInstances(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "failed to describe instances by tags")
	}

	// TODO: currently just returns the first matched instance, need to
	// better rationalize how to find the right instance to return if multiple
	// match
	for _, res := range out.Reservations {
		for _, inst := range res.Instances {
			return s.SDKToInstance(inst)
		}
	}

	return nil, nil
}

// InstanceIfExists returns the existing instance or nothing if it doesn't exist.
func (s *Service) InstanceIfExists(id *string) (*infrav1.Instance, error) {
	if id == nil {
		s.scope.Info("Instance does not have an instance id")
		return nil, nil
	}

	s.scope.V(2).Info("Looking for instance by id", "instance-id", *id)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{id},
	}

	out, err := s.scope.EC2.DescribeInstances(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, errors.Wrapf(err, "failed to describe instance: %q", *id)
	}

	if len(out.Reservations) > 0 && len(out.Reservations[0].Instances) > 0 {
		return s.SDKToInstance(out.Reservations[0].Instances[0])
	}

	return nil, nil
}

// CreateInstance runs an ec2 instance.
func (s *Service) CreateInstance(scope *scope.MachineScope) (*infrav1.Instance, error) {
	s.scope.V(2).Info("Creating an instance for a machine")

	input := &infrav1.Instance{
		Type:              scope.AWSMachine.Spec.InstanceType,
		IAMProfile:        scope.AWSMachine.Spec.IAMInstanceProfile,
		RootDeviceSize:    scope.AWSMachine.Spec.RootDeviceSize,
		NetworkInterfaces: scope.AWSMachine.Spec.NetworkInterfaces,
	}

	// Make sure to use the MachineScope here to get the merger of AWSCluster and AWSMachine tags
	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	input.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String(scope.Role()),
		Additional:  additionalTags,
	})

	var err error
	// Pick image from the machine configuration, or use a default one.
	if scope.AWSMachine.Spec.AMI.ID != nil {
		input.ImageID = *scope.AWSMachine.Spec.AMI.ID
	} else {
		input.ImageID, err = s.defaultAMILookup(scope.AWSMachine.Spec.ImageLookupOrg, "ubuntu", "18.04", *scope.Machine.Spec.Version)
		if err != nil {
			return nil, err
		}
	}

	// Pick subnet from the machine configuration, or based on the availability zone specified,
	// or default to the first private subnet available.
	// TODO(vincepri): Move subnet picking logic to its own function/method.
	if scope.AWSMachine.Spec.Subnet != nil && scope.AWSMachine.Spec.Subnet.ID != nil {
		input.SubnetID = *scope.AWSMachine.Spec.Subnet.ID
	} else if scope.AWSMachine.Spec.AvailabilityZone != nil {
		sns := s.scope.Subnets().FilterPrivate().FilterByZone(*scope.AWSMachine.Spec.AvailabilityZone)
		if len(sns) == 0 {
			return nil, awserrors.NewFailedDependency(
				errors.Errorf("failed to run machine %q, no subnets available in availaibility zone %q",
					scope.Name(),
					*scope.AWSMachine.Spec.AvailabilityZone,
				),
			)
		}
		input.SubnetID = sns[0].ID
	} else if input.SubnetID == "" {
		sns := s.scope.Subnets().FilterPrivate()
		if len(sns) == 0 {
			return nil, awserrors.NewFailedDependency(
				errors.Errorf("failed to run machine %q, no subnets available", scope.Name()),
			)
		}
		input.SubnetID = sns[0].ID
	}

	if s.scope.Network().APIServerELB.DNSName == "" {
		return nil, awserrors.NewFailedDependency(
			errors.New("failed to run controlplane, APIServer ELB not available"),
		)
	}

	// Set userdata.
	input.UserData = aws.String(*scope.Machine.Spec.Bootstrap.Data)

	// Set security groups.
	ids, err := s.GetCoreSecurityGroups(scope)
	if err != nil {
		return nil, err
	}
	input.SecurityGroupIDs = append(input.SecurityGroupIDs, ids...)

	// Pick SSH key, if any.
	if scope.AWSMachine.Spec.SSHKeyName != "" {
		input.SSHKeyName = aws.String(scope.AWSMachine.Spec.SSHKeyName)
	} else {
		input.SSHKeyName = aws.String(defaultSSHKeyName)
	}

	s.scope.V(2).Info("Running instance", "machine-role", scope.Role())
	out, err := s.runInstance(scope.Role(), input)
	if err != nil {
		// Only record the failure event if the error is not related to failed dependencies.
		// This is to avoid spamming failure events since the machine will be requeued by the actuator.
		if !awserrors.IsFailedDependency(errors.Cause(err)) {
			record.Warnf(scope.Machine, "FailedCreate", "Failed to create instance: %v", err)
		}
		return nil, err
	}

	if len(input.NetworkInterfaces) > 0 {
		for _, id := range input.NetworkInterfaces {
			s.scope.V(2).Info("Attaching security groups to provided network interface", "groups", input.SecurityGroupIDs, "interface", id)
			if err := s.attachSecurityGroupsToNetworkInterface(input.SecurityGroupIDs, id); err != nil {
				return nil, err
			}
		}
	}

	record.Eventf(scope.Machine, "SuccessfulCreate", "Created new %s instance with id %q", scope.Role(), out.ID)
	return out, nil
}

// GetCoreSecurityGroups looks up the security group IDs managed by this actuator
// They are considered "core" to its proper functioning
func (s *Service) GetCoreSecurityGroups(scope *scope.MachineScope) ([]string, error) {
	// These are common across both controlplane and node machines
	sgRoles := []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupNode,
		infrav1.SecurityGroupLB,
	}
	switch scope.Role() {
	case "node":
		// Just the common security groups above
	case "control-plane":
		sgRoles = append(sgRoles, infrav1.SecurityGroupControlPlane)
	default:
		return nil, errors.Errorf("Unknown node role %q", scope.Role())
	}
	ids := make([]string, 0, len(sgRoles))
	for _, sg := range sgRoles {
		if _, ok := s.scope.SecurityGroups()[sg]; !ok {
			return nil, awserrors.NewFailedDependency(
				errors.Errorf("%s security group not available", sg),
			)
		}
		ids = append(ids, s.scope.SecurityGroups()[sg].ID)
	}
	return ids, nil
}

// TerminateInstance terminates an EC2 instance.
// Returns nil on success, error in all other cases.
func (s *Service) TerminateInstance(instanceID string) error {
	s.scope.V(2).Info("Attempting to terminate instance", "instance-id", instanceID)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	if _, err := s.scope.EC2.TerminateInstances(input); err != nil {
		return errors.Wrapf(err, "failed to terminate instance with id %q", instanceID)
	}

	s.scope.V(2).Info("Terminated instance", "instance-id", instanceID)
	return nil
}

// TerminateInstanceAndWait terminates and waits
// for an EC2 instance to terminate.
func (s *Service) TerminateInstanceAndWait(instanceID string) error {
	if err := s.TerminateInstance(instanceID); err != nil {
		return err
	}

	s.scope.V(2).Info("Waiting for EC2 instance to terminate", "instance-id", instanceID)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	}

	if err := s.scope.EC2.WaitUntilInstanceTerminated(input); err != nil {
		return errors.Wrapf(err, "failed to wait for instance %q termination", instanceID)
	}

	return nil
}

func (s *Service) runInstance(role string, i *infrav1.Instance) (*infrav1.Instance, error) {
	input := &ec2.RunInstancesInput{
		InstanceType: aws.String(i.Type),
		ImageId:      aws.String(i.ImageID),
		KeyName:      i.SSHKeyName,
		EbsOptimized: i.EBSOptimized,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
	}

	if i.UserData != nil {
		var buf bytes.Buffer

		decoded, err := base64.StdEncoding.DecodeString(*i.UserData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode bootstrapData")
		}

		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(decoded); err != nil {
			return nil, errors.Wrap(err, "failed to gzip userdata")
		}

		if err := gz.Close(); err != nil {
			return nil, errors.Wrap(err, "failed to gzip userdata")
		}

		s.scope.V(2).Info("userData size", "bytes", buf.Len(), "role", role)

		input.UserData = aws.String(base64.StdEncoding.EncodeToString(buf.Bytes()))
	}

	if len(i.NetworkInterfaces) > 0 {
		instances := make([]*ec2.InstanceNetworkInterfaceSpecification, 0, len(i.NetworkInterfaces))

		for index, id := range i.NetworkInterfaces {
			instances = append(instances, &ec2.InstanceNetworkInterfaceSpecification{
				NetworkInterfaceId: aws.String(id),
				DeviceIndex:        aws.Int64(int64(index)),
			})
		}

		input.NetworkInterfaces = instances
	} else {
		input.SubnetId = aws.String(i.SubnetID)

		if len(i.SecurityGroupIDs) > 0 {
			input.SecurityGroupIds = aws.StringSlice(i.SecurityGroupIDs)
		}
	}

	if i.IAMProfile != "" {
		input.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
			Name: aws.String(i.IAMProfile),
		}
	}

	if i.RootDeviceSize != 0 {
		rootDeviceName, err := s.getImageRootDevice(i.ImageID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get root volume from image %q", i.ImageID)
		}

		input.BlockDeviceMappings = []*ec2.BlockDeviceMapping{
			{
				DeviceName: rootDeviceName,
				Ebs: &ec2.EbsBlockDevice{
					DeleteOnTermination: aws.Bool(true),
					VolumeSize:          aws.Int64(i.RootDeviceSize),
				},
			},
		}
	}

	if len(i.Tags) > 0 {
		spec := &ec2.TagSpecification{ResourceType: aws.String(ec2.ResourceTypeInstance)}
		for key, value := range i.Tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}

		input.TagSpecifications = append(input.TagSpecifications, spec)
	}

	out, err := s.scope.EC2.RunInstances(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run instance")
	}

	if len(out.Instances) == 0 {
		return nil, errors.Errorf("no instance returned for reservation %v", out.GoString())
	}

	waitTimeout := 1 * time.Minute
	s.scope.V(2).Info("Waiting for instance to be in running state", "instance-id", *out.Instances[0].InstanceId, "timeout", waitTimeout.String())
	ctx, cancel := context.WithTimeout(aws.BackgroundContext(), waitTimeout)
	defer cancel()

	if err := s.scope.EC2.WaitUntilInstanceRunningWithContext(
		ctx,
		&ec2.DescribeInstancesInput{InstanceIds: []*string{out.Instances[0].InstanceId}},
		request.WithWaiterLogger(&awslog{s.scope.Logger}),
	); err != nil {
		s.scope.V(2).Info("Could not determine if Machine is running. Machine state might be unavailable until next renconciliation.")
	}

	return converters.SDKToInstance(out.Instances[0]), nil
}

// An internal type to satisfy aws' log interface.
type awslog struct {
	logr.Logger
}

func (a *awslog) Log(args ...interface{}) {
	a.WithName("aws-logger").Info("AWS context", args...)
}

// GetInstanceSecurityGroups returns a map from ENI id to the security groups applied to that ENI
// While some security group operations take place at the "instance" level, these are in fact an API convenience for manipulating the first ("primary") ENI's properties.
func (s *Service) GetInstanceSecurityGroups(instanceID string) (map[string][]string, error) {
	enis, err := s.getInstanceENIs(instanceID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ENIs for instance %q", instanceID)
	}

	out := make(map[string][]string)
	for _, eni := range enis {
		var groups []string
		for _, group := range eni.Groups {
			groups = append(groups, aws.StringValue(group.GroupId))
		}
		out[aws.StringValue(eni.NetworkInterfaceId)] = groups
	}
	return out, nil
}

// UpdateInstanceSecurityGroups modifies the security groups of the given
// EC2 instance.
func (s *Service) UpdateInstanceSecurityGroups(instanceID string, ids []string) error {
	s.scope.V(2).Info("Attempting to update security groups on instance", "instance-id", instanceID)

	enis, err := s.getInstanceENIs(instanceID)
	if err != nil {
		return errors.Wrapf(err, "failed to get ENIs for instance %q", instanceID)
	}

	s.scope.V(3).Info("Found ENIs on instance", "number-of-enis", len(enis), "instance-id", instanceID)

	for _, eni := range enis {
		if err := s.attachSecurityGroupsToNetworkInterface(ids, aws.StringValue(eni.NetworkInterfaceId)); err != nil {
			return errors.Wrapf(err, "failed to modify network interfaces on instance %q", instanceID)
		}
	}

	return nil
}

// UpdateResourceTags updates the tags for an instance.
// This will be called if there is anything to create (update) or delete.
// We may not always have to perform each action, so we check what we're
// receiving to avoid calling AWS if we don't need to.
func (s *Service) UpdateResourceTags(resourceID *string, create map[string]string, remove map[string]string) error {
	s.scope.V(2).Info("Attempting to update tags on resource", "resource-id", *resourceID)

	// If we have anything to create or update
	if len(create) > 0 {
		s.scope.V(2).Info("Attempting to create tags on resource", "resource-id", *resourceID)

		// Convert our create map into an array of *ec2.Tag
		createTagsInput := converters.MapToTags(create)

		// Create the CreateTags input.
		input := &ec2.CreateTagsInput{
			Resources: []*string{resourceID},
			Tags:      createTagsInput,
		}

		// Create/Update tags in AWS.
		if _, err := s.scope.EC2.CreateTags(input); err != nil {
			return errors.Wrapf(err, "failed to create tags for resource %q: %+v", *resourceID, create)
		}
	}

	// If we have anything to remove
	if len(remove) > 0 {
		s.scope.V(2).Info("Attempting to delete tags on resource", "resource-id", *resourceID)

		// Convert our remove map into an array of *ec2.Tag
		removeTagsInput := converters.MapToTags(remove)

		// Create the DeleteTags input
		input := &ec2.DeleteTagsInput{
			Resources: []*string{resourceID},
			Tags:      removeTagsInput,
		}

		// Delete tags in AWS.
		if _, err := s.scope.EC2.DeleteTags(input); err != nil {
			return errors.Wrapf(err, "failed to delete tags for resource %q: %v", *resourceID, remove)
		}
	}

	return nil
}

func (s *Service) getInstanceENIs(instanceID string) ([]*ec2.NetworkInterface, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("attachment.instance-id"),
				Values: []*string{aws.String(instanceID)},
			},
		},
	}

	output, err := s.scope.EC2.DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, err
	}

	return output.NetworkInterfaces, nil
}

func (s *Service) getImageRootDevice(imageID string) (*string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []*string{aws.String(imageID)},
	}

	output, err := s.scope.EC2.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	if len(output.Images) == 0 {
		return nil, errors.Errorf("no images returned when looking up ID %q", imageID)
	}

	return output.Images[0].RootDeviceName, nil
}

func (s *Service) getInstanceRootDeviceSize(instance *ec2.Instance) (*int64, error) {

	for _, bdm := range instance.BlockDeviceMappings {
		if aws.StringValue(bdm.DeviceName) == aws.StringValue(instance.RootDeviceName) {
			input := &ec2.DescribeVolumesInput{
				VolumeIds: []*string{bdm.Ebs.VolumeId},
			}

			out, err := s.scope.EC2.DescribeVolumes(input)
			if err != nil {
				return nil, err
			}

			if len(out.Volumes) == 0 {
				return nil, errors.Errorf("no volumes found for id %q", aws.StringValue(bdm.Ebs.VolumeId))
			}

			return out.Volumes[0].Size, nil
		}
	}
	return nil, nil
}

// SDKToInstance converts an AWS EC2 SDK instance to the CAPA instance type.
// converters.SDKToInstance populates all instance fields except for rootVolumeSize,
// because EC2.DescribeInstances does not return the size of storage devices. An
// additional call to EC2 is required to get this value.
func (s *Service) SDKToInstance(v *ec2.Instance) (*infrav1.Instance, error) {
	i := &infrav1.Instance{
		ID:           aws.StringValue(v.InstanceId),
		State:        infrav1.InstanceState(*v.State.Name),
		Type:         aws.StringValue(v.InstanceType),
		SubnetID:     aws.StringValue(v.SubnetId),
		ImageID:      aws.StringValue(v.ImageId),
		SSHKeyName:   v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	// Extract IAM Instance Profile name from ARN
	// TODO: Handle this comparison more safely, perhaps by querying IAM for the
	// instance profile ARN and comparing to the ARN returned by EC2
	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Arn != nil {
		split := strings.Split(aws.StringValue(v.IamInstanceProfile.Arn), "instance-profile/")
		if len(split) > 1 && split[1] != "" {
			i.IAMProfile = split[1]
		}
	}

	for _, sg := range v.SecurityGroups {
		i.SecurityGroupIDs = append(i.SecurityGroupIDs, *sg.GroupId)
	}

	if len(v.Tags) > 0 {
		i.Tags = converters.TagsToMap(v.Tags)
	}

	rootSize, err := s.getInstanceRootDeviceSize(v)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get root volume size for instance: %q", aws.StringValue(v.InstanceId))
	}

	i.RootDeviceSize = aws.Int64Value(rootSize)
	return i, nil
}

func (s *Service) getNetworkInterfaceSecurityGroups(interfaceID string) ([]string, error) {
	input := &ec2.DescribeNetworkInterfaceAttributeInput{
		Attribute:          aws.String("groupSet"),
		NetworkInterfaceId: aws.String(interfaceID),
	}

	output, err := s.scope.EC2.DescribeNetworkInterfaceAttribute(input)
	if err != nil {
		return nil, err
	}

	groups := make([]string, len(output.Groups))
	for i := range output.Groups {
		groups[i] = aws.StringValue(output.Groups[i].GroupId)
	}

	return groups, nil
}

func (s *Service) attachSecurityGroupsToNetworkInterface(groups []string, interfaceID string) error {
	existingGroups, err := s.getNetworkInterfaceSecurityGroups(interfaceID)
	if err != nil {
		return errors.Wrapf(err, "failed to look up network interface security groups: %+v", err)
	}

	var totalGroups []string
	copy(existingGroups, totalGroups)

	for _, group := range groups {
		if !util.Contains(existingGroups, group) {
			totalGroups = append(totalGroups, group)
		}
	}

	// no new groups to attach
	if len(existingGroups) == len(totalGroups) {
		return nil
	}

	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(interfaceID),
		Groups:             aws.StringSlice(totalGroups),
	}

	if _, err := s.scope.EC2.ModifyNetworkInterfaceAttribute(input); err != nil {
		return errors.Wrapf(err, "failed to modify interface %q to have security groups %v", interfaceID, totalGroups)
	}
	return nil
}

// DetachSecurityGroupsFromNetworkInterface looks up an ENI by interfaceID and
// detaches a list of Security Groups from that ENI.
func (s *Service) DetachSecurityGroupsFromNetworkInterface(groups []string, interfaceID string) error {
	existingGroups, err := s.getNetworkInterfaceSecurityGroups(interfaceID)
	if err != nil {
		return errors.Wrapf(err, "failed to look up network interface security groups")
	}

	remainingGroups := existingGroups
	for _, group := range groups {
		remainingGroups = util.Filter(remainingGroups, group)
	}

	input := &ec2.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(interfaceID),
		Groups:             aws.StringSlice(remainingGroups),
	}

	if _, err := s.scope.EC2.ModifyNetworkInterfaceAttribute(input); err != nil {
		return errors.Wrapf(err, "failed to modify interface %q", interfaceID)
	}
	return nil
}
