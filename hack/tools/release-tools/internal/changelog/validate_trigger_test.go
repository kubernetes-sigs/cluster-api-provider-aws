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

func TestValidateTrigger(t *testing.T) {
	tests := []struct {
		name         string
		targetBranch string
		version      string
		wantErr      bool
		want         Trigger
	}{
		{
			name:         "main minor release is valid",
			targetBranch: "main",
			version:      "v2.10.0",
			want: Trigger{
				ReleaseBranch:       "release-2.10",
				CreateReleaseBranch: true,
				Major:               2,
				Minor:               10,
				Patch:               0,
			},
		},
		{
			name:         "main major release is valid",
			targetBranch: "main",
			version:      "v3.0.0",
			want: Trigger{
				ReleaseBranch:       "release-3.0",
				CreateReleaseBranch: true,
				Major:               3,
				Minor:               0,
				Patch:               0,
			},
		},
		{
			name:         "main rejects a non-zero patch",
			targetBranch: "main",
			version:      "v2.10.1",
			wantErr:      true,
		},
		{
			name:         "release branch patch release is valid",
			targetBranch: "release-2.9",
			version:      "v2.9.3",
			want: Trigger{
				ReleaseBranch:       "release-2.9",
				CreateReleaseBranch: false,
				Major:               2,
				Minor:               9,
				Patch:               3,
			},
		},
		{
			name:         "release branch rejects patch 0",
			targetBranch: "release-2.9",
			version:      "v2.9.0",
			wantErr:      true,
		},
		{
			name:         "release branch rejects a minor bump",
			targetBranch: "release-2.9",
			version:      "v2.10.0",
			wantErr:      true,
		},
		{
			name:         "release branch rejects a major bump",
			targetBranch: "release-2.9",
			version:      "v3.9.0",
			wantErr:      true,
		},
		{
			name:         "unsupported target branch",
			targetBranch: "some-feature-branch",
			version:      "v2.9.3",
			wantErr:      true,
		},
		{
			name:         "invalid semver",
			targetBranch: "main",
			version:      "2.10.0",
			wantErr:      true,
		},
		{
			name:         "pre-release semver rejected",
			targetBranch: "main",
			version:      "v2.10.0-beta.0",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateTrigger(tt.targetBranch, tt.version)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateTrigger() expected error, got nil (result: %+v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateTrigger() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ValidateTrigger() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
