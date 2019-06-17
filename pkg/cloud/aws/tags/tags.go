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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

// ApplyParams are function parameters used to apply tags on an aws resource.
type ApplyParams struct {
	v1alpha1.BuildParams
	EC2Client ec2iface.EC2API
}

// Apply tags a resource with tags including the cluster tag.
func Apply(params *ApplyParams) error {
	tags := v1alpha1.Build(params.BuildParams)

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
func Ensure(current v1alpha1.Tags, params *ApplyParams) error {
	want := v1alpha1.Build(params.BuildParams)
	if !current.Equals(want) {
		return Apply(params)
	}
	return nil
}
