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

package v1alpha1_test

import (
	"reflect"
	"testing"

	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestClusterConfigFromProviderConfigNoError(t *testing.T) {
	testcases := []struct {
		name     string
		input    clusterv1.ProviderConfig
		expected providerv1.AWSClusterProviderConfig
	}{
		{
			name: "ProviderConfig.Value is nil should return an empty AWSClusterProviderConfig",
			input: clusterv1.ProviderConfig{
				Value: nil,
			},
			expected: providerv1.AWSClusterProviderConfig{},
		},
		{
			name: "ProviderConfig actually has data",
			input: clusterv1.ProviderConfig{
				Value: cloudtest.RuntimeRawExtension(t, &providerv1.AWSClusterProviderConfig{
					Region: "us-west-1",
				}),
			},
			expected: providerv1.AWSClusterProviderConfig{
				Region: "us-west-1",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := providerv1.ClusterConfigFromProviderConfig(tc.input)
			if err != nil {
				t.Fatalf("did not expect an error but got: %v", err)
			}
			if !reflect.DeepEqual(actual, &tc.expected) {
				t.Fatalf("expected equality but got\n%#v\nand expected\n%#v", actual, tc.expected)
			}
		})
	}
}
