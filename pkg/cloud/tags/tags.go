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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

// ApplyParams are function parameters used to apply tags on an aws resource.
type ApplyParams struct {
	infrav1.BuildParams
	EC2Client ec2iface.EC2API
}

// Apply tags a resource with tags including the cluster tag.
func Apply(params *ApplyParams) error {
	tags := infrav1.Build(params.BuildParams)

	awsTags := make([]*ec2.Tag, 0, len(tags))
	for k, v := range tags {
		tag := &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		awsTags = append(awsTags, tag)
	}

	createTagsInput := &ec2.CreateTagsInput{
		Resources: aws.StringSlice([]string{params.ResourceID}),
		Tags:      awsTags,
	}

	_, err := params.EC2Client.CreateTags(createTagsInput)
	return errors.Wrapf(err, "failed to tag resource %q in cluster %q", params.ResourceID, params.ClusterName)
}

// Ensure applies the tags if the current tags differ from the params.
func Ensure(current infrav1.Tags, params *ApplyParams) error {
	diff := computeDiff(current, params.BuildParams)
	if len(diff) > 0 {
		return Apply(params)
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
