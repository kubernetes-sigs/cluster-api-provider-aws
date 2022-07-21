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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/annotations"
)

const (
	serviceNameTag    = "kubernetes.io/service-name"
	eksClusterNameTag = "aws:eks:cluster-name"
)

// Reconcile will perform any setup operations for garbage collection. Default behaviour is to mark
// a cluster as requiring garbage collection unless it explicitly opts out.
func (s *Service) Reconcile(ctx context.Context) error {
	s.scope.Info("reconciling garbage collection")

	val, found := annotations.Get(s.scope.InfraCluster(), annotations.ExternalResourceGCAnnotation)
	if !found {
		val = "true"
	}

	shouldGC, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("converting value %s of annotation %s to bool: %w", val, annotations.ExternalResourceGCAnnotation, err)
	}

	if shouldGC {
		s.scope.V(2).Info("Enabling garbage collection for cluster")
		controllerutil.AddFinalizer(s.scope.InfraCluster(), expinfrav1.ExternalResourceGCFinalizer)

		if patchErr := s.scope.PatchObject(); patchErr != nil {
			return fmt.Errorf("patching infra cluster after adding gc finalizer: %w", patchErr)
		}
	}

	return nil
}

// ReconcileDelete performs any operations that relate to the reconciliation of a cluster delete as it relates to
// the external resources created by the workload cluster. For example, it will delete the ELB/NLBs that where created
// as a result of Services of type load balancer.
func (s *Service) ReconcileDelete(ctx context.Context) error {
	s.scope.Info("reconciling deletion for garbage collection")

	if err := s.deleteResources(ctx); err != nil {
		return fmt.Errorf("deleting workload services of type load balancer: %w", err)
	}

	s.scope.V(2).Info("Removing garbage collection finalizer cluster")
	controllerutil.RemoveFinalizer(s.scope.InfraCluster(), expinfrav1.ExternalResourceGCFinalizer)

	if patchErr := s.scope.PatchObject(); patchErr != nil {
		return fmt.Errorf("patching infra cluster after removing gc finalizer: %w", patchErr)
	}

	return nil
}

func (s *Service) deleteResources(ctx context.Context) error {
	s.scope.Info("deleting aws resources created by tenant cluster")

	serviceTag := infrav1.ClusterAWSCloudProviderTagKey(s.scope.KubernetesClusterName())
	awsInput := rgapi.GetResourcesInput{
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
