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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
)

func (s *Service) deleteElasticLoadbalancingResources(ctx context.Context, resources []*rgapi.ResourceTagMapping) error {
	s.scope.V(2).Info("Deleting load balancers")
	if err := s.deleteLoadBalancers(ctx, resources); err != nil {
		return fmt.Errorf("deleting load balancers: %w", err)
	}
	if err := s.deleteTargetGroups(ctx, resources); err != nil {
		return fmt.Errorf("deleting target groups: %w", err)
	}

	s.scope.V(2).Info("Finished deleting elasticloadbalancing resources")

	return nil
}

func (s *Service) deleteLoadBalancers(ctx context.Context, resources []*rgapi.ResourceTagMapping) error {
	for i := range resources {
		res := resources[i]

		lbServiceName := getTagValue(serviceNameTag, res)
		if lbServiceName == "" {
			s.scope.V(2).Info("Resource wasn't created for a Service via CCM, skipping load balancer deletion")

			continue
		}

		parsedARN, err := arn.Parse(*res.ResourceARN)
		if err != nil {
			return fmt.Errorf("parsing arn %s: %w", *res.ResourceARN, err)
		}

		if strings.HasPrefix(parsedARN.Resource, "loadbalancer/app/") {
			s.scope.V(2).Info("Deleting ALB for Service", "service", lbServiceName, "arn", parsedARN.String())
			return s.deleteLoadBalancerV2(ctx, &parsedARN)
		}
		if strings.HasPrefix(parsedARN.Resource, "loadbalancer/net/") {
			s.scope.V(2).Info("Deleting NLB for Service", "service", lbServiceName, "arn", parsedARN.String())
			return s.deleteLoadBalancerV2(ctx, &parsedARN)
		}
		if strings.HasPrefix(parsedARN.Resource, "loadbalancer/") {
			s.scope.V(2).Info("Deleting classic ELB for Service", "service", lbServiceName, "arn", parsedARN.String())
			return s.deleteLoadBalancer(ctx, &parsedARN)
		}
	}

	s.scope.V(2).Info("Finished processing tagged resources for load balancers")

	return nil
}

func (s *Service) deleteTargetGroups(ctx context.Context, resources []*rgapi.ResourceTagMapping) error {
	for i := range resources {
		res := resources[i]

		lbServiceName := getTagValue(serviceNameTag, res)
		if lbServiceName == "" {
			s.scope.V(2).Info("Resource wasn't created for a Service via CCM, skipping load balancer deletion")

			continue
		}

		parsedARN, err := arn.Parse(*res.ResourceARN)
		if err != nil {
			return fmt.Errorf("parsing arn %s: %w", *res.ResourceARN, err)
		}

		if strings.HasPrefix(parsedARN.Resource, "targetgroup/") {
			s.scope.V(2).Info("Deleting target group for Service", "service", lbServiceName, "arn", parsedARN.String())
			return s.deleteTargetGroup(ctx, &parsedARN)
		}
	}

	s.scope.V(2).Info("Finished processing tagged resources for target groups")

	return nil
}

func (s *Service) deleteLoadBalancerV2(ctx context.Context, lbARN *arn.ARN) error {
	input := elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(lbARN.String()),
	}

	s.scope.V(2).Info("Deleting v2 load balancer", "arn", lbARN.String())
	_, err := s.elbv2Client.DeleteLoadBalancerWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("deleting v2 load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteLoadBalancer(ctx context.Context, lbARN *arn.ARN) error {
	name := strings.ReplaceAll(lbARN.Resource, "loadbalancer/", "")
	input := elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	s.scope.V(2).Info("Deleting classic load balancer", "name", name, "arn", lbARN.String())
	_, err := s.elbClient.DeleteLoadBalancerWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("deleting classic load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteTargetGroup(ctx context.Context, lbARN *arn.ARN) error {
	name := strings.ReplaceAll(lbARN.Resource, "targetgroup/", "")
	input := elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(lbARN.String()),
	}

	s.scope.V(2).Info("Deleting target group", "name", name, "arn", lbARN.String())
	_, err := s.elbv2Client.DeleteTargetGroupWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("deleting target group: %w", err)
	}

	return nil
}
