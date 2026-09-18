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
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runGit runs a git command in dir with a fixed, deterministic commit
// identity, failing the test on error.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.CommandContext(context.Background(), "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func skipIfNoGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available")
	}
}

// setupIntegrationRepo creates a repo with a v2.13.0 tag on main, then a
// v2.13.2 tag on a diverged release-2.13 branch that main's subsequent
// history does not include — reproducing the scenario where
// "git describe --abbrev=0" on main would miss the true previous version.
func setupIntegrationRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-q", "-b", "main")

	metadata := `# maps release series of major.minor to cluster-api contract version
apiVersion: clusterctl.cluster.x-k8s.io/v1alpha3
kind: Metadata
releaseSeries:
  - major: 2
    minor: 13
    contract: v1beta1
`
	writeAndCommit(t, dir, "metadata.yaml", metadata, "initial metadata")
	runGit(t, dir, "tag", "-a", "v2.13.0", "-m", "v2.13.0")

	runGit(t, dir, "checkout", "-q", "-b", "release-2.13")
	writeAndCommit(t, dir, "fix.txt", "a patch fix\n", "a patch fix")
	runGit(t, dir, "tag", "-a", "v2.13.2", "-m", "v2.13.2")

	runGit(t, dir, "checkout", "-q", "main")
	writeAndCommit(t, dir, "unrelated.txt", "more work on main\n", "unrelated main work")

	return dir
}

func writeAndCommit(t *testing.T, dir, relPath, content, message string) {
	t.Helper()
	path := filepath.Join(dir, relPath)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", relPath, err)
	}
	runGit(t, dir, "add", relPath)
	runGit(t, dir, "commit", "-q", "-m", message)
}

func gitTags(t *testing.T, dir string) []string {
	t.Helper()
	out := runGit(t, dir, "tag", "-l", "v*")
	var tags []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			tags = append(tags, line)
		}
	}
	return tags
}

func TestIntegration_MainMinorRelease(t *testing.T) {
	skipIfNoGit(t)
	dir := setupIntegrationRepo(t)

	// The bug this guards against: on main, the nearest *reachable* tag is
	// v2.13.0, but the true previous version is v2.13.2 on release-2.13.
	previous, err := PreviousVersion(gitTags(t, dir), "v2.14.0")
	if err != nil {
		t.Fatalf("PreviousVersion() unexpected error: %v", err)
	}
	if previous != "v2.13.2" {
		t.Fatalf("PreviousVersion() = %q, want %q (reachability-based selection would wrongly give v2.13.0)", previous, "v2.13.2")
	}

	// Simulate `make release-notes-pr VERSION=v2.14.0`: write the changelog
	// with front-matter, then update metadata.yaml, as notes-pr does.
	changelogPath := filepath.Join(dir, "CHANGELOG")
	if err := os.MkdirAll(changelogPath, 0o755); err != nil {
		t.Fatalf("mkdir CHANGELOG: %v", err)
	}
	notesPath := filepath.Join(changelogPath, "v2.14.0.md")
	if err := os.WriteFile(notesPath, WithFrontMatter("v1beta1", []byte("# Changelog\n\n- a new feature\n")), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	if err := AppendReleaseSeries(filepath.Join(dir, "metadata.yaml"), 2, 14, "v1beta1"); err != nil {
		t.Fatalf("AppendReleaseSeries() unexpected error: %v", err)
	}
	runGit(t, dir, "add", "CHANGELOG/v2.14.0.md", "metadata.yaml")
	runGit(t, dir, "commit", "-q", "-m", "Add CHANGELOG/v2.14.0.md")

	// Pre-merge check: metadata.yaml must already match the changelog.
	if err := VerifyMetadata(notesPath, filepath.Join(dir, "metadata.yaml"), 2, 14); err != nil {
		t.Fatalf("VerifyMetadata() unexpected error: %v", err)
	}

	// Post-merge (release-trigger.yaml): validate, then branch + tag.
	trigger, err := ValidateTrigger("main", "v2.14.0")
	if err != nil {
		t.Fatalf("ValidateTrigger() unexpected error: %v", err)
	}
	if !trigger.CreateReleaseBranch || trigger.ReleaseBranch != "release-2.14" {
		t.Fatalf("ValidateTrigger() = %+v, want CreateReleaseBranch=true, ReleaseBranch=release-2.14", trigger)
	}

	runGit(t, dir, "checkout", "-q", "-b", trigger.ReleaseBranch)
	runGit(t, dir, "tag", "-a", "v2.14.0", "-F", notesPath)

	// The release branch inherits the correct metadata.yaml from the merge
	// commit it was cut from — no separate metadata step is needed.
	branchMetadata, err := os.ReadFile(filepath.Join(dir, "metadata.yaml")) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading metadata.yaml on release branch: %v", err)
	}
	if !strings.Contains(string(branchMetadata), "minor: 14") {
		t.Errorf("release-2.14's metadata.yaml missing the new series entry: %s", branchMetadata)
	}

	tags := gitTags(t, dir)
	found := false
	for _, tg := range tags {
		if tg == "v2.14.0" {
			found = true
		}
	}
	if !found {
		t.Errorf("tag v2.14.0 was not created, tags: %v", tags)
	}
}

func TestIntegration_ReleaseBranchPatch(t *testing.T) {
	skipIfNoGit(t)
	dir := setupIntegrationRepo(t)
	runGit(t, dir, "checkout", "-q", "release-2.13")

	beforeMetadata, err := os.ReadFile(filepath.Join(dir, "metadata.yaml")) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading metadata.yaml: %v", err)
	}

	changelogDir := filepath.Join(dir, "CHANGELOG")
	if err := os.MkdirAll(changelogDir, 0o755); err != nil {
		t.Fatalf("mkdir CHANGELOG: %v", err)
	}
	notesPath := filepath.Join(changelogDir, "v2.13.3.md")
	if err := os.WriteFile(notesPath, WithFrontMatter("v1beta1", []byte("# Changelog\n\n- a bugfix\n")), 0o600); err != nil {
		t.Fatalf("writing changelog: %v", err)
	}
	runGit(t, dir, "add", "CHANGELOG/v2.13.3.md")
	runGit(t, dir, "commit", "-q", "-m", "Add CHANGELOG/v2.13.3.md")

	trigger, err := ValidateTrigger("release-2.13", "v2.13.3")
	if err != nil {
		t.Fatalf("ValidateTrigger() unexpected error: %v", err)
	}
	if trigger.CreateReleaseBranch {
		t.Fatalf("ValidateTrigger() = %+v, want CreateReleaseBranch=false for a patch release", trigger)
	}

	runGit(t, dir, "tag", "-a", "v2.13.3", "-F", notesPath)

	afterMetadata, err := os.ReadFile(filepath.Join(dir, "metadata.yaml")) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading metadata.yaml: %v", err)
	}
	if !bytes.Equal(beforeMetadata, afterMetadata) {
		t.Errorf("metadata.yaml was modified for a patch release")
	}
}

func TestIntegration_RejectionCasesLeaveNoState(t *testing.T) {
	skipIfNoGit(t)
	dir := setupIntegrationRepo(t)
	tagsBefore := gitTags(t, dir)

	tests := []struct {
		name   string
		branch string
		ver    string
	}{
		{"main rejects non-zero patch", "main", "v2.14.1"},
		{"release branch rejects mismatched minor", "release-2.13", "v2.14.0"},
		{"release branch rejects patch 0", "release-2.13", "v2.13.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ValidateTrigger(tt.branch, tt.ver); err == nil {
				t.Fatalf("ValidateTrigger(%q, %q) expected error, got nil", tt.branch, tt.ver)
			}
		})
	}

	// A rejected ValidateTrigger call must be side-effect free: it doesn't
	// touch git state on its own, so the tag set is exactly as it was.
	tagsAfter := gitTags(t, dir)
	if strings.Join(tagsBefore, ",") != strings.Join(tagsAfter, ",") {
		t.Errorf("tag set changed after rejected triggers: before=%v after=%v", tagsBefore, tagsAfter)
	}
}
