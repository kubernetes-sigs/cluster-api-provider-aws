/*
Copyright 2026 The Kubernetes Authors.

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

package template

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestToYAML(t *testing.T) {
	tests := []struct {
		name       string
		raw        *runtime.RawExtension
		want       string
		wantErrSub string
	}{
		{
			name: "nil raw extension",
			raw:  nil,
			want: "",
		},
		{
			name: "empty raw extension",
			raw:  &runtime.RawExtension{},
			want: "",
		},
		{
			name: "json raw extension",
			raw:  &runtime.RawExtension{Raw: []byte(`{"maxPods":110}`)},
			want: "maxPods: 110\n",
		},
		{
			name: "yaml raw extension",
			raw:  &runtime.RawExtension{Raw: []byte("maxPods: 110\n")},
			want: "maxPods: 110\n",
		},
		{
			name:       "invalid raw extension",
			raw:        &runtime.RawExtension{Raw: []byte(`{"maxPods":`)},
			wantErrSub: "runtime object raw is neither json nor yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToYAML(tc.raw)
			if tc.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErrSub, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("unexpected YAML: got %q, want %q", got, tc.want)
			}
		})
	}
}
