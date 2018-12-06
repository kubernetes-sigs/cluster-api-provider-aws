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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

func SDKToInstance(v *ec2.Instance) *v1alpha1.Instance {
	i := &v1alpha1.Instance{
		ID:           aws.StringValue(v.InstanceId),
		State:        v1alpha1.InstanceState(*v.State.Name),
		Type:         aws.StringValue(v.InstanceType),
		SubnetID:     aws.StringValue(v.SubnetId),
		ImageID:      aws.StringValue(v.ImageId),
		KeyName:      v.KeyName,
		PrivateIP:    v.PrivateIpAddress,
		PublicIP:     v.PublicIpAddress,
		ENASupport:   v.EnaSupport,
		EBSOptimized: v.EbsOptimized,
	}

	for _, sg := range v.SecurityGroups {
		i.SecurityGroupIDs = append(i.SecurityGroupIDs, *sg.GroupId)
	}

	// TODO: Handle returned IAM instance profile, since we are currently
	// using a string representing the name, but the InstanceProfile returned
	// from the sdk only returns ARN and ID.

	if len(v.Tags) > 0 {
		i.Tags = TagsToMap(v.Tags)
	}

	return i
}
