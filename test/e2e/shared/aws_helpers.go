//go:build e2e
// +build e2e

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

package shared

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/client"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

type LoadBalancerType string

var (
	LoadBalancerTypeELB = LoadBalancerType("elb")
	LoadBalancerTypeALB = LoadBalancerType("alb")
	LoadBalancerTypeNLB = LoadBalancerType("nlb")
)

type WaitForLoadBalancerToExistForServiceInput struct {
	AWSSession       client.ConfigProvider
	ServiceName      string
	ServiceNamespace string
	ClusterName      string
	Type             LoadBalancerType
}

func WaitForLoadBalancerToExistForService(ctx context.Context, input WaitForLoadBalancerToExistForServiceInput, intervals ...interface{}) {
	By(fmt.Sprintf("Waiting for AWS load balancer of type %s to exist for service %s/%s", input.Type, input.ServiceNamespace, input.ServiceName))

	Eventually(func() bool {
		arns, err := GetLoadBalancerARNs(ctx, GetLoadBalancerARNsInput{ //nolint: gosimple
			AWSSession:       input.AWSSession,
			ServiceName:      input.ServiceName,
			ServiceNamespace: input.ServiceNamespace,
			ClusterName:      input.ClusterName,
			Type:             input.Type,
		})
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "error getting loadbalancer arns: %v\n", err)

			return false
		}
		if len(arns) == 0 {
			return false
		}

		return true
	}, intervals...).Should(BeTrue(), "failed to wait for loadbalancer")
}

type GetLoadBalancerARNsInput struct {
	AWSSession       client.ConfigProvider
	ServiceName      string
	ServiceNamespace string
	ClusterName      string
	Type             LoadBalancerType
}

func GetLoadBalancerARNs(ctx context.Context, input GetLoadBalancerARNsInput) ([]string, error) {
	By(fmt.Sprintf("Getting AWS load balancer ARNs of type %s for service %s/%s", input.Type, input.ServiceNamespace, input.ServiceName))

	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(input.ClusterName)
	tags := map[string][]string{
		"kubernetes.io/service-name": {fmt.Sprintf("%s/%s", input.ServiceNamespace, input.ServiceName)},
		serviceTag:                   {string(infrav1.ResourceLifecycleOwned)},
	}
	descInput := &DescribeResourcesByTagsInput{
		AWSSession: input.AWSSession,
		Tags:       tags,
	}

	descOutput, err := DescribeResourcesByTags(*descInput)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "error querying resources by tags: %v\n", err)

		return nil, fmt.Errorf("describing resource tags: %w", err)
	}

	matchingARNs := []string{}
	for i := range descOutput.ARNs {
		resARN := descOutput.ARNs[i]
		parsedArn, err := arn.Parse(resARN)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "error parsing arn %s: %v\n", resARN, err)

			return nil, fmt.Errorf("parsing resource arn %s: %w", resARN, err)
		}

		if parsedArn.Service != "elasticloadbalancing" {
			continue
		}

		switch input.Type {
		case LoadBalancerTypeALB:
			if strings.HasPrefix(parsedArn.Resource, "loadbalancer/app/") {
				matchingARNs = append(matchingARNs, resARN)
			}
		case LoadBalancerTypeNLB:
			if strings.HasPrefix(parsedArn.Resource, "loadbalancer/net/") {
				matchingARNs = append(matchingARNs, resARN)
			}
		case LoadBalancerTypeELB:
			if strings.HasPrefix(parsedArn.Resource, "loadbalancer/") {
				matchingARNs = append(matchingARNs, resARN)
			}
		}
	}

	return matchingARNs, nil
}

type DescribeResourcesByTagsInput struct {
	AWSSession client.ConfigProvider
	Tags       map[string][]string
}

type DescribeResourcesByTagsOutput struct {
	ARNs []string
}

func DescribeResourcesByTags(input DescribeResourcesByTagsInput) (*DescribeResourcesByTagsOutput, error) {
	if len(input.Tags) == 0 {
		return nil, errors.New("you must supply tags")
	}

	awsInput := rgapi.GetResourcesInput{
		TagFilters: []*rgapi.TagFilter{},
	}

	for k, v := range input.Tags {
		awsInput.TagFilters = append(awsInput.TagFilters, &rgapi.TagFilter{
			Key:    aws.String(k),
			Values: aws.StringSlice(v),
		})
	}

	rgSvc := rgapi.New(input.AWSSession)
	awsOutput, err := rgSvc.GetResources(&awsInput)
	if err != nil {
		return nil, fmt.Errorf("getting resources by tags: %w", err)
	}

	output := &DescribeResourcesByTagsOutput{
		ARNs: []string{},
	}
	for _, res := range awsOutput.ResourceTagMappingList {
		output.ARNs = append(output.ARNs, *res.ResourceARN)
	}

	return output, nil
}
