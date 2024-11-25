/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
)

func (s *Service) reconcileSecurityGroups(cluster *eks.Cluster) error {
	s.scope.Info("Reconciling EKS security groups", "cluster-name", ptr.Deref(cluster.Name, ""))

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[infrav1.SecurityGroupRole]infrav1.SecurityGroup)
	}

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:aws:eks:cluster-name"),
				Values: []*string{cluster.Name},
			},
		},
	}

	output, err := s.EC2Client.DescribeSecurityGroupsWithContext(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("describing security groups: %w", err)
	}

	if len(output.SecurityGroups) == 0 {
		return ErrNoSecurityGroup
	}

	sg := infrav1.SecurityGroup{
		ID:   *output.SecurityGroups[0].GroupId,
		Name: *output.SecurityGroups[0].GroupName,
		Tags: converters.TagsToMap(output.SecurityGroups[0].Tags),
	}
	s.scope.ControlPlane.Status.Network.SecurityGroups[infrav1.SecurityGroupNode] = sg

	input = &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			cluster.ResourcesVpcConfig.ClusterSecurityGroupId,
		},
	}

	output, err = s.EC2Client.DescribeSecurityGroupsWithContext(context.TODO(), input)
	if err != nil || len(output.SecurityGroups) == 0 {
		return fmt.Errorf("describing EKS cluster security group: %w", err)
	}

	clusterSecurityGroup := infrav1.SecurityGroup{
		ID:   aws.StringValue(cluster.ResourcesVpcConfig.ClusterSecurityGroupId),
		Name: *output.SecurityGroups[0].GroupName,
		Tags: converters.TagsToMap(output.SecurityGroups[0].Tags),
	}
	s.scope.ControlPlane.Status.Network.SecurityGroups[ekscontrolplanev1.SecurityGroupCluster] = clusterSecurityGroup

	additionalTags := s.scope.ControlPlane.Spec.AdditionalTags
	if !reflect.DeepEqual(sg.Tags, desiredTags(sg.Tags, additionalTags)) {
		if err = s.updateTagsForEKSManagedSecurityGroup(&sg.ID, sg.Tags, desiredTags(sg.Tags, additionalTags)); err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(clusterSecurityGroup.Tags, desiredTags(clusterSecurityGroup.Tags, additionalTags)) {
		if err = s.updateTagsForEKSManagedSecurityGroup(&clusterSecurityGroup.ID, clusterSecurityGroup.Tags, desiredTags(clusterSecurityGroup.Tags, additionalTags)); err != nil {
			return err
		}
	}

	return nil
}

// desiredTags will return the default tags of EKS cluster with the spec.additionalTags from AWSManagedControlPlane
func desiredTags(existingTags, additionalTags infrav1.Tags) infrav1.Tags {
	merged := make(infrav1.Tags)

	for key, value := range existingTags {
		// since the cluster is created/managed by CAPA, existing tags will contain the default EKS security group tags
		if key == "aws:eks:cluster-name" || key == "Name" || strings.Contains(key, infrav1.NameKubernetesAWSCloudProviderPrefix) {
			merged[key] = value
		}
	}

	// additional tags from spec.additionalTags of AWSManagedControlPlane will be added/updated to desired tags considering them as source of truth
	for key, value := range additionalTags {
		merged[key] = value
	}

	return merged
}

// updateTagsForEKSManagedSecurityGroup will update the tags in the EKS security group with the desired tags via create/update/delete operations
func (s *Service) updateTagsForEKSManagedSecurityGroup(securityGroupID *string, existingTags, desiredTags infrav1.Tags) error {
	tagsToDelete, newTags := getTagUpdates(existingTags, desiredTags)

	// Create tags for updating or adding
	if len(newTags) > 0 {
		desiredTags := converters.MapToTags(newTags)
		_, err := s.EC2Client.CreateTags(&ec2.CreateTagsInput{
			Resources: []*string{securityGroupID},
			Tags:      desiredTags,
		})
		if err != nil {
			return fmt.Errorf("failed to create/update tags: %v", err)
		}
	}

	// Delete tags
	if len(tagsToDelete) > 0 {
		var ec2TagKeys []*ec2.Tag
		for _, key := range tagsToDelete {
			// the default tags added to EKS cluster via AWS will not be deleted
			if key != "aws:eks:cluster-name" && key != "Name" && !strings.Contains(key, infrav1.NameKubernetesAWSCloudProviderPrefix) {
				ec2TagKeys = append(ec2TagKeys, &ec2.Tag{
					Key: aws.String(key),
				})
			}
		}

		_, err := s.EC2Client.DeleteTags(&ec2.DeleteTagsInput{
			Resources: []*string{securityGroupID},
			Tags:      ec2TagKeys,
		})
		if err != nil {
			return fmt.Errorf("failed to delete tags: %v", err)
		}
	}

	return nil
}
