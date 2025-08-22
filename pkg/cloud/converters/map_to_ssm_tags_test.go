/*
Copyright 2021 The Kubernetes Authors.

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
package converters

import (
	"reflect"
	"testing"

	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func TestMapToSSMTags(t *testing.T) {
	tests := []struct {
		name     string
		input    infrav1.Tags
		expected []ssmtypes.Tag
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: []ssmtypes.Tag{},
		},
		{
			name:     "empty map",
			input:    infrav1.Tags{},
			expected: []ssmtypes.Tag{},
		},
		{
			name:  "single key-value",
			input: infrav1.Tags{"k1": "v1"},
			expected: []ssmtypes.Tag{
				{Key: strPtr("k1"), Value: strPtr("v1")},
			},
		},
		{
			name:  "multiple keys",
			input: infrav1.Tags{"k1": "v1", "k2": "v2"},
			expected: []ssmtypes.Tag{
				{Key: strPtr("k1"), Value: strPtr("v1")},
				{Key: strPtr("k2"), Value: strPtr("v2")},
			},
		},
		{
			name:  "empty value string",
			input: infrav1.Tags{"k1": ""},
			expected: []ssmtypes.Tag{
				{Key: strPtr("k1"), Value: strPtr("")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapToSSMTags(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("got %d tags, want %d", len(got), len(tt.expected))
				return
			}
			gotMap := make(map[string]string)
			for _, tag := range got {
				gotMap[*tag.Key] = *tag.Value
			}
			expectedMap := make(map[string]string)
			for _, tag := range tt.expected {
				expectedMap[*tag.Key] = *tag.Value
			}
			if !reflect.DeepEqual(gotMap, expectedMap) {
				t.Errorf("MapToSSMTags() = %v, want %v", gotMap, expectedMap)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
