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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/tags"
)

func (s *Service) reconcileTags(cluster *eks.Cluster) error {
	clusterTags := converters.MapPtrToMap(cluster.Tags)
	buildParams := s.getEKSTagParams(*cluster.Arn)
	tagsBuilder := tags.New(buildParams, tags.WithEKS(s.EKSClient))
	if err := tagsBuilder.Ensure(clusterTags); err != nil {
		return fmt.Errorf("failed ensuring tags on cluster: %w", err)
	}

	return nil
}

func (s *Service) getEKSTagParams(id string) *infrav1.BuildParams {
	name := s.scope.KubernetesClusterName()

	return &infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		ResourceID:  id,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(infrav1.CommonRoleTagValue),
		Additional:  s.scope.AdditionalTags(),
	}
}

func getTagUpdates(currentTags map[string]string, tags map[string]string) (untagKeys []string, newTags map[string]string) {
	untagKeys = []string{}
	newTags = make(map[string]string)
	for key := range currentTags {
		if _, ok := tags[key]; !ok {
			untagKeys = append(untagKeys, key)
		}
	}
	for key, value := range tags {
		if currentV, ok := currentTags[key]; !ok || value != currentV {
			newTags[key] = value
		}
	}
	return untagKeys, newTags
}

func (s *NodegroupService) reconcileTags(ng *eks.Nodegroup) error {
	tags := ngTags(s.scope.ClusterName(), s.scope.AdditionalTags())
	return updateTags(s.EKSClient, ng.NodegroupArn, aws.StringValueMap(ng.Tags), tags)
}

func (s *FargateService) reconcileTags(fp *eks.FargateProfile) error {
	tags := ngTags(s.scope.ClusterName(), s.scope.AdditionalTags())
	return updateTags(s.EKSClient, fp.FargateProfileArn, aws.StringValueMap(fp.Tags), tags)
}

func updateTags(client eksiface.EKSAPI, arn *string, existingTags, desiredTags map[string]string) error {
	untagKeys, newTags := getTagUpdates(existingTags, desiredTags)

	if len(newTags) > 0 {
		tagInput := &eks.TagResourceInput{
			ResourceArn: arn,
			Tags:        aws.StringMap(newTags),
		}
		_, err := client.TagResource(tagInput)
		if err != nil {
			return err
		}
	}

	if len(untagKeys) > 0 {
		untagInput := &eks.UntagResourceInput{
			ResourceArn: arn,
			TagKeys:     aws.StringSlice(untagKeys),
		}
		_, err := client.UntagResource(untagInput)
		if err != nil {
			return err
		}
	}

	return nil
}
