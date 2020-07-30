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

package network

import (
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

func (s *Service) applyTags(params *infrav1.BuildParams) error {
	tags := infrav1.Build(*params)

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

	_, err := s.EC2Client.CreateTags(createTagsInput)
	return errors.Wrapf(err, "failed to tag resource %q in cluster %q", params.ResourceID, params.ClusterName)
}
