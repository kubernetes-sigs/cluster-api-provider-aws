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

package filter

const (
	// TagNameKubernetesClusterPrefix is the tag name we use to differentiate multiple
	// logically independent clusters running in the same AZ.
	// The tag key = TagNameKubernetesClusterPrefix + clusterID
	// The tag value is an ownership value
	TagNameKubernetesClusterPrefix = "kubernetes.io/cluster/"

	// TagNameAWSProviderManaged is the tag name we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses TagNameKubernetesClusterPrefix
	TagNameAWSProviderManaged = "sigs.k8s.io/cluster-api-provider-aws/managed"

	// TagNameAWSClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	TagNameAWSClusterAPIRole = "sigs.k8s.io/cluster-api-provider-aws/role"

	// TagValueAPIServerRole describes the value for the apiserver role
	TagValueAPIServerRole = "apiserver"

	// TagValueBastionRole describes the value for the bastion role
	TagValueBastionRole = "bastion"

	// TagValueCommonRole describes the value for the common role
	TagValueCommonRole = "common"
)
