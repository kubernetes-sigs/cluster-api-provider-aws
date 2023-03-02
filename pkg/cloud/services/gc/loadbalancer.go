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
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
)

func (s *Service) deleteLoadBalancers(ctx context.Context, resources []*AWSResource) error {
	for _, resource := range resources {
		if !s.isELBResourceToDelete(resource, "loadbalancer") {
			s.scope.Debug("Resource not a load balancer for deletion", "arn", resource.ARN.String())
			continue
		}

		switch {
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/app/"):
			s.scope.Debug("Deleting ALB for Service", "arn", resource.ARN.String())
			if err := s.deleteLoadBalancerV2(ctx, resource.ARN.String()); err != nil {
				return fmt.Errorf("deleting ALB: %w", err)
			}
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/net/"):
			s.scope.Debug("Deleting NLB for Service", "arn", resource.ARN.String())
			if err := s.deleteLoadBalancerV2(ctx, resource.ARN.String()); err != nil {
				return fmt.Errorf("deleting NLB: %w", err)
			}
		case strings.HasPrefix(resource.ARN.Resource, "loadbalancer/"):
			name := strings.ReplaceAll(resource.ARN.Resource, "loadbalancer/", "")
			s.scope.Debug("Deleting classic ELB for Service", "arn", resource.ARN.String(), "name", name)
			if err := s.deleteLoadBalancer(ctx, name); err != nil {
				return fmt.Errorf("deleting classic ELB: %w", err)
			}
		default:
			s.scope.Trace("Unexpected elasticloadbalancing resource, ignoring", "arn", resource.ARN.String())
		}
	}

	s.scope.Debug("Finished processing tagged resources for load balancers")

	return nil
}

func (s *Service) deleteTargetGroups(ctx context.Context, resources []*AWSResource) error {
	for _, resource := range resources {
		if !s.isELBResourceToDelete(resource, "targetgroup") {
			s.scope.Trace("Resource not a target group for deletion", "arn", resource.ARN.String())
			continue
		}

		name := strings.ReplaceAll(resource.ARN.Resource, "targetgroup/", "")
		if err := s.deleteTargetGroup(ctx, resource.ARN.String()); err != nil {
			return fmt.Errorf("deleting target group %s: %w", name, err)
		}
	}
	s.scope.Debug("Finished processing resources for target group deletion")

	return nil
}

func (s *Service) isELBResourceToDelete(resource *AWSResource, resourceName string) bool {
	if !s.isMatchingResource(resource, elb.ServiceName, resourceName) {
		return false
	}

	if serviceName := resource.Tags[serviceNameTag]; serviceName == "" {
		s.scope.Debug("Resource wasn't created for a Service via CCM", "arn", resource.ARN.String(), "resource_name", resourceName)
		return false
	}

	return true
}

func (s *Service) deleteLoadBalancerV2(ctx context.Context, lbARN string) error {
	input := elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(lbARN),
	}

	s.scope.Debug("Deleting v2 load balancer", "arn", lbARN)
	if _, err := s.elbv2Client.DeleteLoadBalancerWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting v2 load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteLoadBalancer(ctx context.Context, name string) error {
	input := elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	}

	s.scope.Debug("Deleting classic load balancer", "name", name)
	if _, err := s.elbClient.DeleteLoadBalancerWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting classic load balancer: %w", err)
	}

	return nil
}

func (s *Service) deleteTargetGroup(ctx context.Context, targetGroupARN string) error {
	input := elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(targetGroupARN),
	}

	s.scope.Debug("Deleting target group", "arn", targetGroupARN)
	if _, err := s.elbv2Client.DeleteTargetGroupWithContext(ctx, &input); err != nil {
		return fmt.Errorf("deleting target group: %w", err)
	}

	return nil
}

func (s *Service) describeTargetgroups() ([]string, error) {
	groups, err := s.elbv2Client.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		return nil, err
	}

	targetGroups := []string{}
	for _, group := range groups.TargetGroups {
		if err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, *group.TargetGroupArn)
	}
	return targetGroups, nil
}

// describeLoadBalancers gets all elastic LBs.
func (s *Service) describeLoadBalancers() ([]string, error) {
	var names []string
	err := s.elbClient.DescribeLoadBalancersPages(&elb.DescribeLoadBalancersInput{}, func(r *elb.DescribeLoadBalancersOutput, last bool) bool {
		for _, lb := range r.LoadBalancerDescriptions {
			names = append(names, *lb.LoadBalancerName)
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return names, nil
}

// describeLoadBalancersV2 gets all network and application LBs.
func (s *Service) describeLoadBalancersV2() ([]string, error) {
	var arns []string
	err := s.elbv2Client.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{}, func(r *elbv2.DescribeLoadBalancersOutput, last bool) bool {
		for _, lb := range r.LoadBalancers {
			arns = append(arns, *lb.LoadBalancerArn)
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return arns, nil
}

// getProviderOwnedLoadBalancers gets cloud provider created LB(ELB) for this cluster, filtering by tag: kubernetes.io/cluster/<cluster-name>:owned.
func (s *Service) getProviderOwnedLoadBalancers() ([]*AWSResource, error) {
	names, err := s.describeLoadBalancers()
	if err != nil {
		return nil, err
	}

	return s.filterProviderOwnedLB(names)
}

// getProviderOwnedLoadBalancersV2 gets cloud provider created LBv2(NLB and ALB) for this cluster, filtering by tag: kubernetes.io/cluster/<cluster-name>:owned.
func (s *Service) getProviderOwnedLoadBalancersV2() ([]*AWSResource, error) {
	arns, err := s.describeLoadBalancersV2()
	if err != nil {
		return nil, err
	}

	return s.filterProviderOwnedLBV2(arns)
}

// getProviderOwnedTargetgroups gets cloud provider created target groups of v2 LBs(NLB and ALB) for this cluster, filtering by tag: kubernetes.io/cluster/<cluster-name>:owned.
func (s *Service) getProviderOwnedTargetgroups() ([]*AWSResource, error) {
	targetGroups, err := s.describeTargetgroups()
	if err != nil {
		return nil, err
	}

	return s.filterProviderOwnedLBV2(targetGroups)
}

// filterProviderOwnedLB filters LB resource tags by tag: kubernetes.io/cluster/<cluster-name>:owned.
func (s *Service) filterProviderOwnedLB(names []string) ([]*AWSResource, error) {
	resources := []*AWSResource{}
	lbChunks := chunkResources(names)
	for _, chunk := range lbChunks {
		output, err := s.elbClient.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: aws.StringSlice(chunk)})
		if err != nil {
			return nil, err
		}

		for _, tagDesc := range output.TagDescriptions {
			for _, tag := range tagDesc.Tags {
				serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())
				if *tag.Key == serviceTag && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
					arn := composeArn(elbService, elbResourcePrefix+*tagDesc.LoadBalancerName)
					resource, err := composeAWSResource(aws.String(arn), converters.ELBTagsToMap(tagDesc.Tags))
					if err != nil {
						return nil, err
					}
					resources = append(resources, resource)
					break
				}
			}
		}
	}

	return resources, nil
}

// filterProviderOwnedLBV2 filters LBv2 resource tags by tag: kubernetes.io/cluster/<cluster-name>:owned.
func (s *Service) filterProviderOwnedLBV2(arns []string) ([]*AWSResource, error) {
	resources := []*AWSResource{}
	lbChunks := chunkResources(arns)
	for _, chunk := range lbChunks {
		output, err := s.elbv2Client.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: aws.StringSlice(chunk)})
		if err != nil {
			return nil, err
		}

		for _, tagDesc := range output.TagDescriptions {
			for _, tag := range tagDesc.Tags {
				serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())
				if *tag.Key == serviceTag && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
					resource, err := composeAWSResource(tagDesc.ResourceArn, converters.V2TagsToMap(tagDesc.Tags))
					if err != nil {
						return nil, err
					}
					resources = append(resources, resource)
					break
				}
			}
		}
	}

	return resources, nil
}
