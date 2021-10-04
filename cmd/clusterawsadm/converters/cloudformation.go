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

package converters

import (
	"github.com/awslabs/goformation/v4/cloudformation/tags"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

// MapToCloudFormationTags converts a infrav1.Tags to []tags.Tag.
func MapToCloudFormationTags(src infrav1.Tags) []tags.Tag {
	cfnTags := make([]tags.Tag, 0, len(src))

	for k, v := range src {
		tag := tags.Tag{
			Key:   k,
			Value: v,
		}

		cfnTags = append(cfnTags, tag)
	}

	return cfnTags
}
