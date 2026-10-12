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

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestVerifyMetadata(t *testing.T) {
	tests := []struct {
		name       string
		changelog  string
		metadata   string
		major      int
		minor      int
		wantErr    bool
		wantErrSub string
	}{
		{
			name:      "matching entry passes",
			changelog: "---\ncontract: v1beta1\n---\n# Changelog\n",
			metadata:  sampleMetadata + "  - major: 2\n    minor: 10\n    contract: v1beta1\n",
			major:     2,
			minor:     10,
		},
		{
			name:       "missing entry fails",
			changelog:  "---\ncontract: v1beta1\n---\n# Changelog\n",
			metadata:   sampleMetadata,
			major:      2,
			minor:      10,
			wantErr:    true,
			wantErrSub: "no releaseSeries entry",
		},
		{
			name:       "contract mismatch fails",
			changelog:  "---\ncontract: v1beta2\n---\n# Changelog\n",
			metadata:   sampleMetadata + "  - major: 2\n    minor: 10\n    contract: v1beta1\n",
			major:      2,
			minor:      10,
			wantErr:    true,
			wantErrSub: "must match",
		},
		{
			name:       "changelog missing contract front-matter fails",
			changelog:  "# Changelog\n",
			metadata:   sampleMetadata + "  - major: 2\n    minor: 10\n    contract: v1beta1\n",
			major:      2,
			minor:      10,
			wantErr:    true,
			wantErrSub: "front-matter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			changelogPath := writeFile(t, dir, "changelog.md", tt.changelog)
			metadataPath := writeFile(t, dir, "metadata.yaml", tt.metadata)

			err := VerifyMetadata(changelogPath, metadataPath, tt.major, tt.minor)
			if tt.wantErr {
				if err == nil {
					t.Fatal("VerifyMetadata() expected error, got nil")
				}
				if tt.wantErrSub != "" && !strings.Contains(err.Error(), tt.wantErrSub) {
					t.Errorf("VerifyMetadata() error = %q, want substring %q", err.Error(), tt.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("VerifyMetadata() unexpected error: %v", err)
			}
		})
	}
}
