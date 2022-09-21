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

package asg

import (
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the asg client.
type Service struct {
	scope     cloud.ClusterScoper
	ASGClient autoscalingiface.AutoScalingAPI
	EC2Client ec2iface.EC2API
}

// NewService returns a new service given the asg api client.
func NewService(clusterScope cloud.ClusterScoper) *Service {
	return &Service{
		scope:     clusterScope,
		ASGClient: scope.NewASGClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		EC2Client: scope.NewEC2Client(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
	}
}
