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

package ami

import (
	"strings"
	"testing"
)

func TestReadPublishedAMIs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantCount int
	}{
		{
			name:      "empty input returns empty list",
			input:     "",
			wantCount: 0,
		},
		{
			name:      "whitespace-only input returns empty list",
			input:     "  \n\t",
			wantCount: 0,
		},
		{
			name:      "null items returns empty list",
			input:     `{"kind":"AWSAMIList","items":null}`,
			wantCount: 0,
		},
		{
			name:      "empty items array returns empty list",
			input:     `{"kind":"AWSAMIList","items":[]}`,
			wantCount: 0,
		},
		{
			name: "items are parsed",
			input: `{"items":[
				{"spec":{"kubernetesVersion":"v1.32.0","os":"ubuntu-24.04","region":"us-east-1"}},
				{"spec":{"kubernetesVersion":"v1.31.4","os":"flatcar-stable","region":"ap-southeast-2"}}
			]}`,
			wantCount: 2,
		},
		{
			name:    "plain text input returns parse error",
			input:   "No AMIs found\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPublishedAMIs(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("readPublishedAMIs() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("readPublishedAMIs() unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d published AMIs, want %d", len(got), tt.wantCount)
			}
		})
	}
}
