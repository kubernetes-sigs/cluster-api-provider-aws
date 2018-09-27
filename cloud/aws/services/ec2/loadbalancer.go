// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
)

func (s *Service) reconcileLoadbalancer(lb v1alpha1.LoadBalancer) error {
	existing, err := s.describeLoadbalancers(lb.Name)

	// compare existing and "current" to do reconcile
	return nil
}

func (s *Service) describeLoadbalancers(apiServerLBName string) (v1alpha1.LoadBalancer, error) {
	req := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(apiServerLBName),
		},
	}

	out, err := s.ELB.DescribeLoadBalancers(req)

	lb := v1alpha1.LoadBalancer{}
	// Parse out fields from response and populate lb
	// lb.Name := out.LoadBalancerDescription[0].
	return lb, nil
}
