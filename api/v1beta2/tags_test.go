/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestTagsMerge(t *testing.T) {
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
			if e, a := tc.expected, tags; !cmp.Equal(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}

func TestTagsDifference(t *testing.T) {
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
			if e, a := tc.expected, out; !cmp.Equal(e, a) {
				t.Errorf("expected %#v, got %#v", e, a)
			}
		})
	}
}

func TestTagsValidate(t *testing.T) {
	tests := []struct {
		name     string
		self     Tags
		expected []*field.Error
	}{
		{
			name: "no errors",
			self: Tags{
				"validKey": "validValue",
			},
			expected: nil,
		},
		{
			name: "no errors - spaces allowed",
			self: Tags{
				"validKey": "valid Value",
			},
			expected: nil,
		},
		{
			name: "key cannot be empty",
			self: Tags{
				"": "validValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot be empty",
					Field:    "spec.additionalTags",
					BadValue: "",
				},
			},
		},
		{
			name: "key cannot be empty - second element",
			self: Tags{
				"validKey": "validValue",
				"":         "secondValidValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot be empty",
					Field:    "spec.additionalTags",
					BadValue: "",
				},
			},
		},
		{
			name: "key with 128 characters is accepted",
			self: Tags{
				strings.Repeat("CAPI", 32): "validValue",
			},
			expected: nil,
		},
		{
			name: "key too long",
			self: Tags{
				strings.Repeat("CAPI", 33): "validValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot be longer than 128 characters",
					Field:    "spec.additionalTags",
					BadValue: strings.Repeat("CAPI", 33),
				},
			},
		},
		{
			name: "value too long",
			self: Tags{
				"validKey": strings.Repeat("CAPI", 65),
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "value cannot be longer than 256 characters",
					Field:    "spec.additionalTags",
					BadValue: strings.Repeat("CAPI", 65),
				},
			},
		},
		{
			name: "multiple errors are appended",
			self: Tags{
				"validKey":                 strings.Repeat("CAPI", 65),
				strings.Repeat("CAPI", 33): "validValue",
				"":                         "thirdValidValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "value cannot be longer than 256 characters",
					Field:    "spec.additionalTags",
					BadValue: strings.Repeat("CAPI", 65),
				},
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot be longer than 128 characters",
					Field:    "spec.additionalTags",
					BadValue: strings.Repeat("CAPI", 33),
				},
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot be empty",
					Field:    "spec.additionalTags",
					BadValue: "",
				},
			},
		},
		{
			name: "key has aws: prefix",
			self: Tags{
				"aws:key": "validValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "user created tag's key cannot have prefix aws:",
					Field:    "spec.additionalTags",
					BadValue: "aws:key",
				},
			},
		},
		{
			name: "key has wrong characters",
			self: Tags{
				"wrong*key": "validValue",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ .",
					Field:    "spec.additionalTags",
					BadValue: "wrong*key",
				},
			},
		},
		{
			name: "value has wrong characters",
			self: Tags{
				"validKey": "wrong*value",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "value cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ .",
					Field:    "spec.additionalTags",
					BadValue: "wrong*value",
				},
			},
		},
		{
			name: "value and key both has wrong characters",
			self: Tags{
				"wrong*key": "wrong*value",
			},
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "key cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ .",
					Field:    "spec.additionalTags",
					BadValue: "wrong*key",
				},
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "value cannot have characters other than alphabets, numbers, spaces and _ . : / = + - @ .",
					Field:    "spec.additionalTags",
					BadValue: "wrong*value",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.self.Validate()
			sort.Slice(out, getSortFieldErrorsFunc(out))
			sort.Slice(tc.expected, getSortFieldErrorsFunc(tc.expected))

			if !cmp.Equal(out, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, out)
			}
		})
	}
}

func getSortFieldErrorsFunc(errs []*field.Error) func(i, j int) bool {
	return func(i, j int) bool {
		if errs[i].Detail != errs[j].Detail {
			return errs[i].Detail < errs[j].Detail
		}
		iBV, ok := errs[i].BadValue.(string)
		if !ok {
			panic("unexpected error converting BadValue to string")
		}
		jBV, ok := errs[j].BadValue.(string)
		if !ok {
			panic("unexpected error converting BadValue to string")
		}
		return iBV < jBV
	}
}
