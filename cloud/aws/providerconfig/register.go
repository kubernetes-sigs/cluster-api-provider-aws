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

package providerconfig

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeBuilder with scheme builder
	SchemeBuilder runtime.SchemeBuilder
	// AddToScheme is method for adding objects to the scheme
	AddToScheme        = SchemeBuilder.AddToScheme
	localSchemeBuilder = &SchemeBuilder
)

func init() {
	localSchemeBuilder.Register(addKnownTypes)
}

// GroupName is group name of the cluster kinds
const GroupName = "aws.cluster.k8s.io"

// SchemeGroupVersion is scheme group version of the cluster kinds
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Kind returns group kind for a given kind/object
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource returns group resource for a given resource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AWSMachineProviderConfig{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AWSClusterProviderConfig{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AWSMachineProviderStatus{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&AWSClusterProviderStatus{},
	)
	return nil
}
