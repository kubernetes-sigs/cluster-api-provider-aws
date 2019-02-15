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
	"path"
	"reflect"
)

// Map defines a map of tags.
type Map map[string]string

// Equals returns true if the maps are equal.
func (m Map) Equals(other Map) bool {
	return reflect.DeepEqual(m, other)
}

// HasOwned returns true if the tags contains a tag that marks the resource as owned by the cluster.
func (m Map) HasOwned(cluster string) bool {
	value, ok := m[path.Join(NameKubernetesClusterPrefix, cluster)]
	return ok && ResourceLifecycle(value) == ResourceLifecycleOwned
}

// HasManaged returns true if the map contains NameAWSProviderManaged key set to true.
func (m Map) HasManaged() bool {
	value, ok := m[NameAWSProviderManaged]
	return ok && value == "true"
}

// GetRole returns the Cluster API role for the tagged resource
func (m Map) GetRole() string {
	return m[NameAWSClusterAPIRole]
}

// Difference returns the difference between this map and the other map.
// Items are considered equals if key and value are equals.
func (m Map) Difference(other Map) Map {
	res := make(Map, len(m))

	for key, value := range m {
		if otherValue, ok := other[key]; ok && value == otherValue {
			continue
		}
		res[key] = value
	}

	return res
}

// ResourceLifecycle configures the lifecycle of a resource
type ResourceLifecycle string

const (
	// ResourceLifecycleOwned is the value we use when tagging resources to indicate
	// that the resource is considered owned and managed by the cluster,
	// and in particular that the lifecycle is tied to the lifecycle of the cluster.
	ResourceLifecycleOwned = ResourceLifecycle("owned")

	// ResourceLifecycleShared is the value we use when tagging resources to indicate
	// that the resource is shared between multiple clusters, and should not be destroyed
	// if the cluster is destroyed.
	ResourceLifecycleShared = ResourceLifecycle("shared")

	// NameKubernetesClusterPrefix is the tag name we use to differentiate multiple
	// logically independent clusters running in the same AZ.
	// The tag key = NameKubernetesClusterPrefix + clusterID
	// The tag value is an ownership value
	NameKubernetesClusterPrefix = "kubernetes.io/cluster/"

	// NameAWSProviderPrefix is the tag prefix we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses NameKubernetesClusterPrefix
	NameAWSProviderPrefix = "sigs.k8s.io/cluster-api-provider-aws/"

	// NameAWSProviderManaged is the tag name we use to differentiate
	// cluster-api-provider-aws owned components from other tooling that
	// uses NameKubernetesClusterPrefix
	NameAWSProviderManaged = NameAWSProviderPrefix + "managed"

	// NameAWSClusterAPIRole is the tag name we use to mark roles for resources
	// dedicated to this cluster api provider implementation.
	NameAWSClusterAPIRole = NameAWSProviderPrefix + "role"

	// ValueAPIServerRole describes the value for the apiserver role
	ValueAPIServerRole = "apiserver"

	// ValueBastionRole describes the value for the bastion role
	ValueBastionRole = "bastion"

	// ValueCommonRole describes the value for the common role
	ValueCommonRole = "common"

	// ValuePublicRole describes the value for the public role
	ValuePublicRole = "public"

	// ValuePrivateRole describes the value for the private role
	ValuePrivateRole = "private"
)
