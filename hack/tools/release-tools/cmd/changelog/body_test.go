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
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestBodyCmdStripsFrontMatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "v2.10.0.md")
	if err := os.WriteFile(path, []byte("---\ncontract: v1beta1\n---\n# Changelog\n\n- did a thing\n"), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}

	cmd := bodyCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--input", path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	want := "# Changelog\n\n- did a thing\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

func TestBodyCmdRequiresInput(t *testing.T) {
	cmd := bodyCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing --input, got nil")
	}
}

func TestBodyCmdMissingFile(t *testing.T) {
	cmd := bodyCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetArgs([]string{"--input", filepath.Join(t.TempDir(), "does-not-exist.md")})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() expected error for missing file, got nil")
	}
}
