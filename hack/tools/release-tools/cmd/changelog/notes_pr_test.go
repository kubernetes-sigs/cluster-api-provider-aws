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

const notesPRSampleMetadata = `# maps release series of major.minor to cluster-api contract version
apiVersion: clusterctl.cluster.x-k8s.io/v1alpha3
kind: Metadata
releaseSeries:
  - major: 2
    minor: 9
    contract: v1beta1
`

func TestNotesPRCmdMinorReleaseUpdatesMetadata(t *testing.T) {
	dir := t.TempDir()
	changelogPath := filepath.Join(dir, "v2.10.0.md")
	metadataPath := filepath.Join(dir, "metadata.yaml")
	if err := os.WriteFile(changelogPath, []byte("# Changelog\n\n- did a thing\n"), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	if err := os.WriteFile(metadataPath, []byte(notesPRSampleMetadata), 0o600); err != nil {
		t.Fatalf("writing metadata: %v", err)
	}

	cmd := notesPRCmd()
	cmd.SetArgs([]string{
		"--input", changelogPath,
		"--version", "v2.10.0",
		"--metadata", metadataPath,
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	changelogContent, err := os.ReadFile(changelogPath) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading changelog: %v", err)
	}
	if !strings.Contains(string(changelogContent), "contract: v1beta1") {
		t.Errorf("changelog missing contract front-matter: %s", changelogContent)
	}

	metadataContent, err := os.ReadFile(metadataPath) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading metadata: %v", err)
	}
	if !strings.Contains(string(metadataContent), "major: 2\n    minor: 10\n    contract: v1beta1") {
		t.Errorf("metadata.yaml missing new release series entry: %s", metadataContent)
	}
}

func TestNotesPRCmdPatchReleaseLeavesMetadataUntouched(t *testing.T) {
	dir := t.TempDir()
	changelogPath := filepath.Join(dir, "v2.9.1.md")
	metadataPath := filepath.Join(dir, "metadata.yaml")
	if err := os.WriteFile(changelogPath, []byte("# Changelog\n"), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	if err := os.WriteFile(metadataPath, []byte(notesPRSampleMetadata), 0o600); err != nil {
		t.Fatalf("writing metadata: %v", err)
	}

	cmd := notesPRCmd()
	cmd.SetArgs([]string{
		"--input", changelogPath,
		"--version", "v2.9.1",
		"--metadata", metadataPath,
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	metadataContent, err := os.ReadFile(metadataPath) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading metadata: %v", err)
	}
	if string(metadataContent) != notesPRSampleMetadata {
		t.Errorf("metadata.yaml was modified for a patch release:\n%s", metadataContent)
	}
}

func TestNotesPRCmdRequiresFlags(t *testing.T) {
	cmd := notesPRCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing required flags, got nil")
	}
}
