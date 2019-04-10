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

package machine

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	clusterFilterKeyPrefix = "kubernetes.io/cluster/"
	clusterFilterValue     = "owned"
)

func awsTagFilter(name string) *string {
	return aws.String(fmt.Sprint("tag:", name))
}

func clusterFilterKey(name string) string {
	return fmt.Sprint(clusterFilterKeyPrefix, name)
}

func clusterFilter(name string) *ec2.Filter {
	return &ec2.Filter{
		Name:   awsTagFilter(clusterFilterKey(name)),
		Values: aws.StringSlice([]string{clusterFilterValue}),
	}
}
