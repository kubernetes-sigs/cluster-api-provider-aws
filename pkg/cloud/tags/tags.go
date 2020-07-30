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

package tags

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

var (
	ErrBuildParamsRequired = errors.New("no build params supplied")
	ErrApplyFuncRequired   = errors.New("no tags apply function supplied")
)

// TagsApplyFunc is used to define a function that will apply tags
type TagsApplyFunc func(params *infrav1.BuildParams) error

// Apply tags a resource with tags including the cluster tag.
func Apply(params *infrav1.BuildParams, fn TagsApplyFunc) error {
	if params == nil {
		return ErrBuildParamsRequired
	}
	if fn == nil {
		return ErrApplyFuncRequired
	}

	if err := fn(params); err != nil {
		return fmt.Errorf("failed applying tags: %w", err)
	}
	return nil
}

// BuildParamsToTagSpecification builds a TagSpecification for the specified resource type
func BuildParamsToTagSpecification(ec2ResourceType string, params infrav1.BuildParams) *ec2.TagSpecification {
	tags := infrav1.Build(params)

	tagSpec := &ec2.TagSpecification{ResourceType: aws.String(ec2ResourceType)}

	// For testing, we need sorted keys
	sortedKeys := make([]string, 0, len(tags))
	for k := range tags {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		tagSpec.Tags = append(tagSpec.Tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(tags[key]),
		})
	}

	return tagSpec
}

// Ensure applies the tags if the current tags differ from the params.
func Ensure(current infrav1.Tags, params *infrav1.BuildParams, fn TagsApplyFunc) error {
	diff := computeDiff(current, *params)
	if len(diff) > 0 {
		return Apply(params, fn)
	}
	return nil
}

func computeDiff(current infrav1.Tags, buildParams infrav1.BuildParams) infrav1.Tags {
	want := infrav1.Build(buildParams)

	// Some tags could be external set by some external entities
	// and that means even if there is no change in cluster
	// managed tags, tags would be updated as "current" and
	// "want" would be different due to external tags.
	// This fix makes sure that tags are updated only if
	// there is a change in cluster managed tags.
	return want.Difference(current)
}
