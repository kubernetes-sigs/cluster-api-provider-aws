/*
Copyright 2018 The Kubernetes Authors.

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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

// TagsToMap converts a []*ec2.Tag into a tags.Map.
func TagsToMap(src []*ec2.Tag) tags.Map {
	tags := make(tags.Map, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapToTags converts a tags.Map to a []*ec2.Tag
func MapToTags(src tags.Map) []*ec2.Tag {
	tags := make([]*ec2.Tag, 0, len(src))

	for k, v := range src {
		tag := &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	return tags
}

// ELBTagsToMap converts a []*elb.Tag into a tags.Map.
func ELBTagsToMap(src []*elb.Tag) tags.Map {
	tags := make(tags.Map, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapToELBTags converts a tags.Map to a []*elb.Tag
func MapToELBTags(src tags.Map) []*elb.Tag {
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
