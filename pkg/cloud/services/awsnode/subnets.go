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

package awsnode

import "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"

func (s *Service) secondarySubnets() []*v1alpha3.SubnetSpec {
	subnets := []*v1alpha3.SubnetSpec{}
	for _, sub := range s.scope.Subnets() {
		if val, ok := sub.Tags[v1alpha3.NameAWSSubnetAssociation]; ok && val == v1alpha3.SecondarySubnetTagValue {
			subnets = append(subnets, sub)
		}
	}
	return subnets
}
