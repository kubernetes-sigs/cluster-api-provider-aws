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

package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	amiv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/ami/v1alpha1"
)

// Utility functions for the bootstrap.aws.infrastructure.cluster.x-k8s.io API group.

// NewSchemeAndCodecs is a utility function that returns a Scheme and CodecFactory
// that understand the types in the bootstrap.aws.infrastructure.cluster.x-k8s.io API group.
func NewSchemeAndCodecs() (*runtime.Scheme, *serializer.CodecFactory, error) {
	scheme := runtime.NewScheme()
	if err := amiv1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	codecs := serializer.NewCodecFactory(scheme)
	return scheme, &codecs, nil
}
