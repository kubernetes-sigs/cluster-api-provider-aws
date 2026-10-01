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
	"testing"
)

func TestVerifyMetadataCmd(t *testing.T) {
	dir := t.TempDir()
	changelogPath := filepath.Join(dir, "v2.10.0.md")
	metadataPath := filepath.Join(dir, "metadata.yaml")
	if err := os.WriteFile(changelogPath, []byte("---\ncontract: v1beta1\n---\n# Changelog\n"), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	if err := os.WriteFile(metadataPath, []byte(notesPRSampleMetadata+"  - major: 2\n    minor: 10\n    contract: v1beta1\n"), 0o600); err != nil {
		t.Fatalf("writing metadata: %v", err)
	}

	cmd := verifyMetadataCmd()
	cmd.SetArgs([]string{
		"--changelog", changelogPath,
		"--metadata", metadataPath,
		"--major", "2",
		"--minor", "10",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}
}

func TestVerifyMetadataCmdMissingEntry(t *testing.T) {
	dir := t.TempDir()
	changelogPath := filepath.Join(dir, "v2.10.0.md")
	metadataPath := filepath.Join(dir, "metadata.yaml")
	if err := os.WriteFile(changelogPath, []byte("---\ncontract: v1beta1\n---\n# Changelog\n"), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	if err := os.WriteFile(metadataPath, []byte(notesPRSampleMetadata), 0o600); err != nil {
		t.Fatalf("writing metadata: %v", err)
	}

	cmd := verifyMetadataCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{
		"--changelog", changelogPath,
		"--metadata", metadataPath,
		"--major", "2",
		"--minor", "10",
	})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing metadata entry, got nil")
	}
}

func TestVerifyMetadataCmdRequiresFlags(t *testing.T) {
	cmd := verifyMetadataCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing required flags, got nil")
	}
}
