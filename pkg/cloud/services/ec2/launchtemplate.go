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
	"encoding/base64"
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// GetLaunchTemplate returns the existing LaunchTemplate or nothing if it doesn't exist.
// For now by name until we need the input to be something different
func (s *Service) GetLaunchTemplate(id string) (*expinfrav1.AWSLaunchTemplate, error) {
	if id == "" {
		return nil, nil
	}

	s.scope.V(2).Info("Looking for existing LaunchTemplates")

	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		Versions:         aws.StringSlice([]string{expinfrav1.LaunchTemplateLatestVersion}),
	}

	out, err := s.EC2Client.DescribeLaunchTemplateVersions(input)
	switch {
	case awserrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		s.scope.Info("", "aerr", err.Error())
	}

	if len(out.LaunchTemplateVersions) == 0 {
		return nil, nil
	}

	return s.SDKToLaunchTemplate(out.LaunchTemplateVersions[0])
}

// CreateLaunchTemplate generates a launch template to be used with the autoscaling group
func (s *Service) CreateLaunchTemplate(scope *scope.MachinePoolScope, imageID *string, userData []byte) (string, error) {
	s.scope.Info("Create a new launch template")

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userData)
	if err != nil {
		return "", errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateName: aws.String(scope.Name()),
	}

	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	if len(tags) > 0 {
		spec := &ec2.TagSpecification{ResourceType: aws.String(ec2.ResourceTypeLaunchTemplate)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		input.TagSpecifications = append(input.TagSpecifications, spec)
	}

	result, err := s.EC2Client.CreateLaunchTemplate(input)
	if err != nil {
		return "", err
	}
	return aws.StringValue(result.LaunchTemplate.LaunchTemplateId), nil
}

func (s *Service) CreateLaunchTemplateVersion(scope *scope.MachinePoolScope, imageID *string, userData []byte) error {
	s.scope.V(2).Info("creating new launch template version", "machine-pool", scope.Name())

	launchTemplateData, err := s.createLaunchTemplateData(scope, imageID, userData)
	if err != nil {
		return errors.Wrapf(err, "unable to form launch template data")
	}

	input := &ec2.CreateLaunchTemplateVersionInput{
		LaunchTemplateData: launchTemplateData,
		LaunchTemplateId:   aws.String(scope.AWSMachinePool.Status.LaunchTemplateID),
	}

	_, err = s.EC2Client.CreateLaunchTemplateVersion(input)
	if err != nil {
		return errors.Wrapf(err, "unable to create launch template version")
	}

	return nil
}

func (s *Service) createLaunchTemplateData(scope *scope.MachinePoolScope, imageID *string, userData []byte) (*ec2.RequestLaunchTemplateData, error) {
	lt := scope.AWSMachinePool.Spec.AWSLaunchTemplate

	// An explicit empty string for SSHKeyName means do not specify a key in the ASG launch
	var sshKeyNamePtr *string
	if lt.SSHKeyName != nil && *lt.SSHKeyName != "" {
		sshKeyNamePtr = lt.SSHKeyName
	}

	data := &ec2.RequestLaunchTemplateData{
		InstanceType: aws.String(lt.InstanceType),
		IamInstanceProfile: &ec2.LaunchTemplateIamInstanceProfileSpecificationRequest{
			Name: aws.String(lt.IamInstanceProfile),
		},
		KeyName:  sshKeyNamePtr,
		UserData: pointer.StringPtr(base64.StdEncoding.EncodeToString(userData)),
	}

	ids, err := s.GetCoreNodeSecurityGroups(scope)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		data.SecurityGroupIds = append(data.SecurityGroupIds, aws.String(id))
	}

	// add additional security groups as well
	for _, additionalGroup := range scope.AWSMachinePool.Spec.AWSLaunchTemplate.AdditionalSecurityGroups {
		data.SecurityGroupIds = append(data.SecurityGroupIds, additionalGroup.ID)
	}

	// set the AMI ID
	data.ImageId = imageID

	// Set up root volume
	if lt.RootVolume != nil {
		rootDeviceName, err := s.checkRootVolume(lt.RootVolume, *data.ImageId)
		if err != nil {
			return nil, err
		}

		ebsRootDevice := &ec2.LaunchTemplateEbsBlockDeviceRequest{
			DeleteOnTermination: aws.Bool(true),
			VolumeSize:          aws.Int64(lt.RootVolume.Size),
			Encrypted:           aws.Bool(lt.RootVolume.Encrypted),
		}

		if lt.RootVolume.IOPS != 0 {
			ebsRootDevice.Iops = aws.Int64(lt.RootVolume.IOPS)
		}

		if lt.RootVolume.EncryptionKey != "" {
			ebsRootDevice.Encrypted = aws.Bool(true)
			ebsRootDevice.KmsKeyId = aws.String(lt.RootVolume.EncryptionKey)
		}

		if lt.RootVolume.Type != "" {
			ebsRootDevice.VolumeType = aws.String(lt.RootVolume.Type)
		}

		data.BlockDeviceMappings = []*ec2.LaunchTemplateBlockDeviceMappingRequest{
			{
				DeviceName: rootDeviceName,
				Ebs:        ebsRootDevice,
			},
		}
	}

	data.TagSpecifications = s.buildLaunchTemplateTagSpecificationRequest(scope)

	return data, nil
}

// DeleteLaunchTemplate delete a launch template
func (s *Service) DeleteLaunchTemplate(id string) error {
	s.scope.V(2).Info("Deleting launch template", "id", id)

	input := &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateId: aws.String(id),
	}

	if _, err := s.EC2Client.DeleteLaunchTemplate(input); err != nil {
		return errors.Wrapf(err, "failed to delete launch template %q", id)
	}

	s.scope.V(2).Info("Deleted launch template", "id", id)
	return nil
}

// SDKToLaunchTemplate converts an AWS EC2 SDK instance to the CAPA instance type.
func (s *Service) SDKToLaunchTemplate(d *ec2.LaunchTemplateVersion) (*expinfrav1.AWSLaunchTemplate, error) {
	v := d.LaunchTemplateData
	i := &expinfrav1.AWSLaunchTemplate{
		Name: aws.StringValue(d.LaunchTemplateName),
		AMI: infrav1.AWSResourceReference{
			ID: v.ImageId,
		},
		IamInstanceProfile: aws.StringValue(v.IamInstanceProfile.Name),
		InstanceType:       aws.StringValue(v.InstanceType),
		SSHKeyName:         v.KeyName,
		VersionNumber:      d.VersionNumber,
	}

	// Extract IAM Instance Profile name from ARN
	if v.IamInstanceProfile != nil && v.IamInstanceProfile.Arn != nil {
		split := strings.Split(aws.StringValue(v.IamInstanceProfile.Arn), "instance-profile/")
		if len(split) > 1 && split[1] != "" {
			i.IamInstanceProfile = split[1]
		}
	}

	for _, id := range v.SecurityGroupIds {
		// This will include the core security groups as well, making the "Additional" a bit
		// dishonest. However, including the core groups drastically simplifies comparison with
		// the incoming security groups.
		i.AdditionalSecurityGroups = append(i.AdditionalSecurityGroups, infrav1.AWSResourceReference{ID: id})
	}

	return i, nil
}

// LaunchTemplateNeedsUpdate checks if a new launch template version is needed
func (s *Service) LaunchTemplateNeedsUpdate(scope *scope.MachinePoolScope, incoming *expinfrav1.AWSLaunchTemplate, existing *expinfrav1.AWSLaunchTemplate) (bool, error) {
	if incoming.IamInstanceProfile != existing.IamInstanceProfile {
		return true, nil
	}

	if incoming.InstanceType != existing.InstanceType {
		return true, nil
	}

	incomingIDs := make([]string, len(incoming.AdditionalSecurityGroups))
	for i, ref := range incoming.AdditionalSecurityGroups {
		incomingIDs[i] = aws.StringValue(ref.ID)
	}

	coreIDs, err := s.GetCoreNodeSecurityGroups(scope)
	if err != nil {
		return false, err
	}

	incomingIDs = append(incomingIDs, coreIDs...)

	existingIDs := make([]string, len(existing.AdditionalSecurityGroups))
	for i, ref := range existing.AdditionalSecurityGroups {
		existingIDs[i] = aws.StringValue(ref.ID)
	}

	sort.Strings(incomingIDs)
	sort.Strings(existingIDs)

	if !reflect.DeepEqual(incomingIDs, existingIDs) {
		return true, nil
	}

	return false, nil
}

func (s *Service) DiscoverLaunchTemplateAMI(scope *scope.MachinePoolScope) (*string, error) {
	lt := scope.AWSMachinePool.Spec.AWSLaunchTemplate

	if lt.AMI.ID != nil {
		return lt.AMI.ID, nil
	}

	if scope.MachinePool.Spec.Template.Spec.Version == nil {
		err := errors.New("Either AWSMachinePool's spec.awslaunchtemplate.ami.id or MachinePool's spec.template.spec.version must be defined")
		s.scope.Error(err, "")
		return nil, err
	}

	var lookupAMI string
	var err error

	imageLookupFormat := lt.ImageLookupFormat
	if imageLookupFormat == "" {
		imageLookupFormat = scope.InfraCluster.ImageLookupFormat()
	}

	imageLookupOrg := lt.ImageLookupOrg
	if imageLookupOrg == "" {
		imageLookupOrg = scope.InfraCluster.ImageLookupOrg()
	}

	imageLookupBaseOS := lt.ImageLookupBaseOS
	if imageLookupBaseOS == "" {
		imageLookupBaseOS = scope.InfraCluster.ImageLookupBaseOS()
	}

	if scope.IsEKSManaged() && imageLookupFormat == "" && imageLookupOrg == "" && imageLookupBaseOS == "" {
		lookupAMI, err = s.eksAMILookup(*scope.MachinePool.Spec.Template.Spec.Version)
		if err != nil {
			return nil, err
		}
	} else {
		lookupAMI, err = s.defaultAMIIDLookup(imageLookupFormat, imageLookupOrg, imageLookupBaseOS, *scope.MachinePool.Spec.Template.Spec.Version)
		if err != nil {
			return nil, err
		}
	}

	return aws.String(lookupAMI), nil
}

func (s *Service) buildLaunchTemplateTagSpecificationRequest(scope *scope.MachinePoolScope) []*ec2.LaunchTemplateTagSpecificationRequest {
	tagSpecifications := make([]*ec2.LaunchTemplateTagSpecificationRequest, 0)
	additionalTags := scope.AdditionalTags()
	// Set the cloud provider tag
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(scope.Name()),
		Role:        aws.String("node"),
		Additional:  additionalTags,
	})

	if len(tags) > 0 {
		// tag instances
		spec := &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeInstance)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpecifications = append(tagSpecifications, spec)

		// tag EBS volumes
		spec = &ec2.LaunchTemplateTagSpecificationRequest{ResourceType: aws.String(ec2.ResourceTypeVolume)}
		for key, value := range tags {
			spec.Tags = append(spec.Tags, &ec2.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
		tagSpecifications = append(tagSpecifications, spec)

	}
	return tagSpecifications
}
