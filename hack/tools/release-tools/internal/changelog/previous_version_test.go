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

package changelog

import "testing"

func TestPreviousVersion(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "minor release skips a diverged patch series not reachable from main",
			tags:    []string{"v2.11.0", "v2.12.0", "v2.13.0", "v2.13.1", "v2.13.2"},
			version: "v2.14.0",
			want:    "v2.13.2",
		},
		{
			name:    "major release picks the greatest tag from the prior major",
			tags:    []string{"v1.9.0", "v1.10.0", "v1.10.1", "v2.0.0"},
			version: "v3.0.0",
			want:    "v2.0.0",
		},
		{
			name:    "patch release picks the latest existing patch in the same series",
			tags:    []string{"v2.9.0", "v2.9.1", "v2.9.2", "v2.10.0"},
			version: "v2.9.3",
			want:    "v2.9.2",
		},
		{
			name:    "first patch after a minor falls back to the series' own X.Y.0",
			tags:    []string{"v2.9.0", "v2.10.0"},
			version: "v2.9.1",
			want:    "v2.9.0",
		},
		{
			name:    "ignores malformed tags",
			tags:    []string{"not-a-version", "v2.9.0", "latest"},
			version: "v2.9.1",
			want:    "v2.9.0",
		},
		{
			name:    "no candidate found",
			tags:    []string{"v2.9.0"},
			version: "v1.0.0",
			wantErr: true,
		},
		{
			name:    "invalid target version",
			tags:    []string{"v2.9.0"},
			version: "2.9.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PreviousVersion(tt.tags, tt.version)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("PreviousVersion() expected error, got nil (result: %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("PreviousVersion() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("PreviousVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
