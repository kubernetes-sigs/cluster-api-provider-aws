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

package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.EC2Scope
	EC2Client ec2iface.EC2API

	// SSMClient is used to look up the official EKS AMI ID
	SSMClient ssmiface.SSMAPI
	ec2ServiceFactory func(scope.EC2Scope) services.EC2Interface
}

// NewService returns a new service given the ec2 api client.
func NewService(clusterScope scope.EC2Scope) *Service {
	return &Service{
		scope:     clusterScope,
		EC2Client: scope.NewEC2Client(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
		SSMClient: scope.NewSSMClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster()),
	}
}
