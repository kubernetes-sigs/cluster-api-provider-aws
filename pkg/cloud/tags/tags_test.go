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

package tags

import (
	"reflect"
	"testing"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

func TestTags_ComputeDiff(t *testing.T) {
	pName := "test"
	pRole := "testrole"
	bp := infrav1.BuildParams{
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		ClusterName: "testcluster",
		Name:        &pName,
		Role:        &pRole,
		Additional:  map[string]string{"k1": "v1"},
	}

	tests := []struct {
		name     string
		input    infrav1.Tags
		expected infrav1.Tags
	}{
		{
			name:  "input is nil",
			input: nil,
			expected: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
			},
		},
		{
			name: "same input",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
			},
			expected: infrav1.Tags{},
		},
		{
			name: "input with external tags",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v1",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         pRole,
				"k2":                                  "v2",
			},
			expected: infrav1.Tags{},
		},
		{
			name: "input with modified values",
			input: infrav1.Tags{
				"Name":                                pName,
				"k1":                                  "v2",
				infrav1.ClusterTagKey(bp.ClusterName): string(infrav1.ResourceLifecycleOwned),
				infrav1.NameAWSClusterAPIRole:         "testrole2",
				"k2":                                  "v2",
			},
			expected: infrav1.Tags{
				"k1":                          "v1",
				infrav1.NameAWSClusterAPIRole: pRole,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := computeDiff(tc.input, bp)
			if e, a := tc.expected, out; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}
