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

package v1beta1

import (
	"reflect"
	"testing"
)

func TestTags_Merge(t *testing.T) {
	tests := []struct {
		name     string
		other    Tags
		expected Tags
	}{
		{
			name:  "nil other",
			other: nil,
			expected: Tags{
				"a": "b",
				"c": "d",
			},
		},
		{
			name:  "empty other",
			other: Tags{},
			expected: Tags{
				"a": "b",
				"c": "d",
			},
		},
		{
			name: "disjoint",
			other: Tags{
				"1": "2",
				"3": "4",
			},
			expected: Tags{
				"a": "b",
				"c": "d",
				"1": "2",
				"3": "4",
			},
		},
		{
			name: "overlapping, other wins",
			other: Tags{
				"1": "2",
				"3": "4",
				"a": "hello",
			},
			expected: Tags{
				"a": "hello",
				"c": "d",
				"1": "2",
				"3": "4",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tags := Tags{
				"a": "b",
				"c": "d",
			}

			tags.Merge(tc.other)
			if e, a := tc.expected, tags; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}

func TestTags_Difference(t *testing.T) {
	tests := []struct {
		name     string
		self     Tags
		input    Tags
		expected Tags
	}{
		{
			name:     "self and input are nil",
			self:     nil,
			input:    nil,
			expected: Tags{},
		},
		{
			name: "input is nil",
			self: Tags{
				"a": "b",
				"c": "d",
			},
			input: nil,
			expected: Tags{
				"a": "b",
				"c": "d",
			},
		},
		{
			name: "similar input",
			self: Tags{
				"a": "b",
				"c": "d",
			},
			input: Tags{
				"a": "b",
				"c": "d",
			},
			expected: Tags{},
		},
		{
			name: "input with extra tags",
			self: Tags{
				"a": "b",
				"c": "d",
			},
			input: Tags{
				"a": "b",
				"c": "d",
				"e": "f",
			},
			expected: Tags{},
		},
		{
			name: "same keys, different values",
			self: Tags{
				"a": "b",
				"c": "d",
			},
			input: Tags{
				"a": "b1",
				"c": "d",
				"e": "f",
			},
			expected: Tags{
				"a": "b",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.self.Difference(tc.input)
			if e, a := tc.expected, out; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}
