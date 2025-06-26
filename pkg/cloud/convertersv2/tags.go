/*
Copyright 2025 The Kubernetes Authors.

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

// Package convertersv2 provides conversion functions for AWS SDK V2 types to CAPA types.
package convertersv2

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	v2ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv1types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// TagsToMap converts a []v2ec2types.Tag into a infrav1.Tags.
func TagsToMap(src []v2ec2types.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapPtrToMap converts a [string]*string into a infrav1.Tags.
func MapPtrToMap(src map[string]*string) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for k, v := range src {
		tags[k] = *v
	}

	return tags
}

// MapToTags converts a infrav1.Tags to a []v2ec2types.Tag.
func MapToTags(src infrav1.Tags) []v2ec2types.Tag {
	tags := make([]v2ec2types.Tag, 0, len(src))

	for k, v := range src {
		tag := v2ec2types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// ELBTagsToMap converts a []elbv1types.Tag into a infrav1.Tags.
func ELBTagsToMap(src []elbv1types.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// V2TagsToMap converts a []elbv2types.Tag into a infrav1.Tags.
func V2TagsToMap(src []elbv2types.Tag) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}

// MapToELBTags converts a infrav1.Tags to a []elbv1types.Tag.
func MapToELBTags(src infrav1.Tags) []elbv1types.Tag {
	tags := make([]elbv1types.Tag, 0, len(src))

	for k, v := range src {
		tag := elbv1types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// MapToV2Tags converts a infrav1.Tags to a []elbv2types.Tag.
func MapToV2Tags(src infrav1.Tags) []elbv2types.Tag {
	tags := make([]elbv2types.Tag, 0, len(src))

	for k, v := range src {
		tag := elbv2types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// MapToSecretsManagerTags converts a infrav1.Tags to a []secretsmanagertypes.Tag.
func MapToSecretsManagerTags(src infrav1.Tags) []secretsmanagertypes.Tag {
	tags := make([]secretsmanagertypes.Tag, 0, len(src))

	for k, v := range src {
		tag := secretsmanagertypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// MapToSSMTags converts a infrav1.Tags to a []ssm.Tag.
func MapToSSMTags(src infrav1.Tags) []ssmtypes.Tag {
	tags := make([]ssmtypes.Tag, 0, len(src))

	for k, v := range src {
		tag := ssmtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// MapToIAMTags converts a infrav1.Tags to a []iamtypes.Tag.
func MapToIAMTags(src infrav1.Tags) []iamtypes.Tag {
	tags := make([]iamtypes.Tag, 0, len(src))

	for k, v := range src {
		tag := iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		tags = append(tags, tag)
	}

	// Sort so that unit tests can expect a stable order
	sort.Slice(tags, func(i, j int) bool { return *tags[i].Key < *tags[j].Key })

	return tags
}

// ASGTagsToMap converts a []autoscalingtypes.TagDescription into a infrav1.Tags.
func ASGTagsToMap(src []autoscalingtypes.TagDescription) infrav1.Tags {
	tags := make(infrav1.Tags, len(src))

	for _, t := range src {
		tags[*t.Key] = *t.Value
	}

	return tags
}
