// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
)

const (
	// TagNameKubernetesClusterPrefix is the tag name we use to differentiate multiple
	// logically independent clusters running in the same AZ.
	// The tag key = TagNameKubernetesClusterPrefix + clusterID
	// The tag value is an ownership value
	TagNameKubernetesClusterPrefix = "kubernetes.io/cluster/"

	// TagNameAWSProviderManaged is the tag name we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses TagNameKubernetesClusterPrefix
	TagNameAWSProviderManaged = "sigs.k8s.io/cluster-api-provider-aws/managed"

	// TagNameAWSClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	TagNameAWSClusterAPIRole = "sigs.k8s.io/cluster-api-provider-aws/role"

	TagValueCommonRole    = "common"
	TagValueAPIServerRole = "apiserver"
)

// ResourceLifecycle configures the lifecycle of a resource
type ResourceLifecycle string

const (
	// ResourceLifecycleOwned is the value we use when tagging resources to indicate
	// that the resource is considered owned and managed by the cluster,
	// and in particular that the lifecycle is tied to the lifecycle of the cluster.
	ResourceLifecycleOwned = ResourceLifecycle("owned")

	// ResourceLifecycleShared is the value we use when tagging resources to indicate
	// that the resource is shared between multiple clusters, and should not be destroyed
	// if the cluster is destroyed.
	ResourceLifecycleShared = ResourceLifecycle("shared")
)

func (s *Service) clusterTagKey(clusterName string) string {
	return fmt.Sprintf("%s%s", TagNameKubernetesClusterPrefix, clusterName)
}

// buildTags builds tags including the cluster tag
func (s *Service) buildTags(clusterName string, lifecycle ResourceLifecycle, name string, role string, additionalTags map[string]string) map[string]string {
	tags := make(map[string]string)
	for k, v := range additionalTags {
		tags[k] = v
	}

	tags[s.clusterTagKey(clusterName)] = string(lifecycle)
	if lifecycle == ResourceLifecycleOwned {
		tags[TagNameAWSProviderManaged] = "true"
	}

	if role != "" {
		tags[TagNameAWSClusterAPIRole] = role
	}

	if name != "" {
		tags["Name"] = name
	}

	return tags
}

// tagsToMap converts a []*elb.Tag into a map[string]string.
func tagsToMap(src []*elb.Tag) map[string]string {
	// Create an array of exactly the length we require to hopefully avoid some
	// allocations while looping.
	tags := make(map[string]string)

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// mapToTags converts a map[string]string to a []*elb.Tag
func mapToTags(src map[string]string) []*elb.Tag {
	// Create an array of exactly the length we require to hopefully avoid some
	// allocations while looping.
	tags := make([]*elb.Tag, 0, len(src))

	for k, v := range src {
		tag := &elb.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}
