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
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	rgapi "github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/annotations"
)

const (
	serviceNameTag    = "kubernetes.io/service-name"
	eksClusterNameTag = "aws:eks:cluster-name"
)

// ReconcileDelete is responsible for determining if the infra cluster needs to be garbage collected. If
// does then it will perform garbage collection. For example, it will delete the ELB/NLBs that where created
// as a result of Services of type load balancer.
func (s *Service) ReconcileDelete(ctx context.Context) error {
	s.scope.Info("reconciling deletion for garbage collection")

	val, found := annotations.Get(s.scope.InfraCluster(), expinfrav1.ExternalResourceGCAnnotation)
	if !found {
		val = "true"
	}

	shouldGC, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("converting value %s of annotation %s to bool: %w", val, expinfrav1.ExternalResourceGCAnnotation, err)
	}

	if shouldGC {
		if err := s.deleteResources(ctx); err != nil {
			return fmt.Errorf("deleting workload services of type load balancer: %w", err)
		}
	}

	return nil
}

func (s *Service) deleteResources(ctx context.Context) error {
	s.scope.Info("deleting aws resources created by tenant cluster")

	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())
	awsInput := rgapi.GetResourcesInput{
		ResourceTypeFilters: nil,
		TagFilters: []*rgapi.TagFilter{
			{
				Key:    aws.String(serviceTag),
				Values: []*string{aws.String(string(infrav1.ResourceLifecycleOwned))},
			},
		},
	}

	awsOutput, err := s.resourceTaggingClient.GetResourcesWithContext(ctx, &awsInput)
	if err != nil {
		return fmt.Errorf("getting tagged resources: %w", err)
	}

	resources := map[string][]*rgapi.ResourceTagMapping{}

	for i := range awsOutput.ResourceTagMappingList {
		mapping := awsOutput.ResourceTagMappingList[i]
		parsedArn, err := arn.Parse(*mapping.ResourceARN)
		if err != nil {
			return fmt.Errorf("parsing resource arn %s: %w", *mapping.ResourceARN, err)
		}

		_, found := s.cleanupFuncs[parsedArn.Service]
		if !found {
			s.scope.V(2).Info("skipping clean-up of tagged resource for service", "service", parsedArn.Service, "arn", mapping.ResourceARN)

			continue
		}

		resources[parsedArn.Service] = append(resources[parsedArn.Service], mapping)
	}

	for svcName, svcResources := range resources {
		cleanupFunc := s.cleanupFuncs[svcName]

		s.scope.V(2).Info("Calling clean-up function for service", "service_name", svcName)
		if deleteErr := cleanupFunc(ctx, svcResources); deleteErr != nil {
			return fmt.Errorf("deleting resources for service %s: %w", svcName, deleteErr)
		}
	}

	return nil
}

func getTagValue(tagName string, mapping *rgapi.ResourceTagMapping) string {
	for _, tag := range mapping.Tags {
		if *tag.Key == tagName {
			return *tag.Value
		}
	}

	return ""
}
