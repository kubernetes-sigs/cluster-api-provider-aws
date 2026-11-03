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

const sampleMetadata = `# maps release series of major.minor to cluster-api contract version
apiVersion: clusterctl.cluster.x-k8s.io/v1alpha3
kind: Metadata
releaseSeries:
  - major: 2
    minor: 8
    contract: v1beta1
  - major: 2
    minor: 9
    contract: v1beta1
`

func writeSampleMetadata(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "metadata.yaml")
	if err := os.WriteFile(path, []byte(sampleMetadata), 0o600); err != nil {
		t.Fatalf("writing sample metadata.yaml: %v", err)
	}
	return path
}

func TestDefaultContract(t *testing.T) {
	path := writeSampleMetadata(t)

	got, err := DefaultContract(path)
	if err != nil {
		t.Fatalf("DefaultContract() unexpected error: %v", err)
	}
	if got != "v1beta1" {
		t.Errorf("DefaultContract() = %q, want %q", got, "v1beta1")
	}
}

func TestAppendReleaseSeriesAddsNewEntry(t *testing.T) {
	path := writeSampleMetadata(t)

	if err := AppendReleaseSeries(path, 2, 10, "v1beta1"); err != nil {
		t.Fatalf("AppendReleaseSeries() unexpected error: %v", err)
	}

	entries, err := parseReleaseSeries(mustRead(t, path))
	if err != nil {
		t.Fatalf("parseReleaseSeries() unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}
	last := entries[2]
	if last.major != 2 || last.minor != 10 || last.contract != "v1beta1" {
		t.Errorf("last entry = %+v, want {2 10 v1beta1}", last)
	}
}

func TestAppendReleaseSeriesIsIdempotent(t *testing.T) {
	path := writeSampleMetadata(t)

	if err := AppendReleaseSeries(path, 2, 9, "v1beta1"); err != nil {
		t.Fatalf("AppendReleaseSeries() unexpected error: %v", err)
	}

	got := string(mustRead(t, path))
	if got != sampleMetadata {
		t.Errorf("metadata.yaml was modified for an already-existing entry:\n%s", got)
	}
}

func TestAppendReleaseSeriesPreservesHeaderComment(t *testing.T) {
	path := writeSampleMetadata(t)

	if err := AppendReleaseSeries(path, 2, 10, "v1beta1"); err != nil {
		t.Fatalf("AppendReleaseSeries() unexpected error: %v", err)
	}

	got := string(mustRead(t, path))
	if !strings.HasPrefix(got, "# maps release series") {
		t.Errorf("header comment was not preserved:\n%s", got)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path) //nolint:gosec // test-controlled path
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return content
}
