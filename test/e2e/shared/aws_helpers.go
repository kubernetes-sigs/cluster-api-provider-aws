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
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/elb"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

type WaitForLoadBalancerToExistForServiceInput struct {
	AWSSession       client.ConfigProvider
	ServiceName      string
	ServiceNamespace string
	ClusterName      string
	Type             infrav1.LoadBalancerType
}

func WaitForLoadBalancerToExistForService(input WaitForLoadBalancerToExistForServiceInput, intervals ...interface{}) {
	By(fmt.Sprintf("Waiting for AWS load balancer of type %s to exist for service %s/%s", input.Type, input.ServiceNamespace, input.ServiceName))

	Eventually(func() bool {
		arns, err := GetLoadBalancerARNs(GetLoadBalancerARNsInput{ //nolint: gosimple
			AWSSession:       input.AWSSession,
			ServiceName:      input.ServiceName,
			ServiceNamespace: input.ServiceNamespace,
			ClusterName:      input.ClusterName,
			Type:             input.Type,
		})
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "Error getting loadbalancer arns: %v\n", err)

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
	Type             infrav1.LoadBalancerType
}

func GetLoadBalancerARNs(input GetLoadBalancerARNsInput) ([]string, error) {
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
		fmt.Fprintf(GinkgoWriter, "Error querying resources by tags: %v\n", err)
		return nil, fmt.Errorf("describing resource tags: %w", err)
	}

	matchingARNs := []string{}
	for _, resARN := range descOutput.ARNs {
		parsedArn, err := arn.Parse(resARN)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "Error parsing arn %q: %v\n", resARN, err)
			return nil, fmt.Errorf("parsing resource arn %q: %w", resARN, err)
		}

		if parsedArn.Service != "elasticloadbalancing" {
			continue
		}

		switch input.Type {
		case infrav1.LoadBalancerTypeALB:
			if strings.HasPrefix(parsedArn.Resource, "loadbalancer/app/") {
				matchingARNs = append(matchingARNs, resARN)
			}
		case infrav1.LoadBalancerTypeNLB:
			if strings.HasPrefix(parsedArn.Resource, "loadbalancer/net/") {
				matchingARNs = append(matchingARNs, resARN)
			}
		case infrav1.LoadBalancerTypeELB:
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

type CheckClassicElbHealthCheckInput struct {
	AWSSession       client.ConfigProvider
	LoadBalancerName string
	ExpectedTarget   string
}

func CheckClassicElbHealthCheck(input CheckClassicElbHealthCheckInput, intervals ...interface{}) {
	Byf("Checking the health check for the classic load balancer %s", input.LoadBalancerName)

	elbSvc := elb.New(input.AWSSession)

	Eventually(func() error {
		out, err := elbSvc.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{
			LoadBalancerNames: []*string{
				aws.String(input.LoadBalancerName),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get list of load balancers: %w", err)
		}
		if len(out.LoadBalancerDescriptions) == 0 {
			return fmt.Errorf("no load balancers found")
		}

		lb := out.LoadBalancerDescriptions[0]
		if lb.HealthCheck.Target == nil {
			return fmt.Errorf("health check target is nil")
		}

		if *lb.HealthCheck.Target != input.ExpectedTarget {
			return fmt.Errorf("health check target does not match expected target")
		}

		return nil
	}, intervals...).Should(Succeed(), "eventually failed trying to check the health check for the classic load balancer")
}
