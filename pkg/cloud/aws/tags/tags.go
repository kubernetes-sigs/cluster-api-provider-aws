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
)

// ApplyParams are function parameters used to apply tags on an aws resource.
type ApplyParams struct {
	BuildParams
	EC2Client ec2iface.EC2API
}

// Apply tags a resource with tags including the cluster tag.
func Apply(params *ApplyParams) error {
	tags := Build(params.BuildParams)

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
func Ensure(current Map, params *ApplyParams) error {
	want := Build(params.BuildParams)
	if !current.Equals(want) {
		return Apply(params)
	}
	return nil
}

// BuildParams is used to build tags around an aws resource.
type BuildParams struct {
	// Lifecycle determines the resource lifecycle.
	Lifecycle ResourceLifecycle

	// ClusterName is the cluster associated with the resource.
	ClusterName string

	// ResourceID is the unique identifier of the resource to be tagged.
	ResourceID string

	// Name is the name of the resource, it's applied as the tag "Name" on AWS.
	// +optional
	Name *string

	// Role is the role associated to the resource.
	// +optional
	Role *string

	// Any additional tags to be added to the resource.
	// +optional
	Additional Map
}

// Build builds tags including the cluster tag and returns them in map form.
func Build(params BuildParams) Map {
	tags := make(Map)
	for k, v := range params.Additional {
		tags[k] = v
	}

	tags[ClusterKey(params.ClusterName)] = string(params.Lifecycle)
	if params.Role != nil {
		tags[NameAWSClusterAPIRole] = *params.Role
	}

	if params.Name != nil {
		tags["Name"] = *params.Name
	}

	return tags
}
