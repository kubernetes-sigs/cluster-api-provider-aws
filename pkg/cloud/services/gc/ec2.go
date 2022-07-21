/*
Copyright 2022 The Kubernetes Authors.

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

package gc

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
)

func (s *Service) deleteEC2Resources(ctx context.Context, resources []*rgapi.ResourceTagMapping) error {
	for i := range resources {
		res := resources[i]

		parsedARN, err := arn.Parse(*res.ResourceARN)
		if err != nil {
			return fmt.Errorf("parsing arn %s: %w", *res.ResourceARN, err)
		}

		if strings.HasPrefix(parsedARN.Resource, "security-group/") {
			s.scope.V(2).Info("Deleting Security group", "arn", parsedARN.String())
			return s.deleteSecurityGroup(ctx, &parsedARN, res)
		}
	}

	s.scope.V(2).Info("Finished deleting ec2 resources")

	return nil
}

func (s *Service) deleteSecurityGroup(ctx context.Context, lbARN *arn.ARN, mapping *rgapi.ResourceTagMapping) error {
	eksClusterName := getTagValue(eksClusterNameTag, mapping)
	if eksClusterName != "" {
		s.scope.V(2).Info("Security group created by EKS directly, skipping deletion", "cluster_name", eksClusterName)

		return nil
	}

	//TODO: should we check for the security group name start with k8s-elb-

	groupID := strings.ReplaceAll(lbARN.Resource, "security-group/", "")
	input := ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(groupID),
	}

	s.scope.V(2).Info("Deleting security group", "group_id", groupID, "arn", lbARN.String())
	_, err := s.ec2Client.DeleteSecurityGroupWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("deleting security group: %w", err)
	}

	return nil
}
