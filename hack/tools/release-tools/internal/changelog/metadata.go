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
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var releaseSeriesEntryRE = regexp.MustCompile(`(?m)^\s*-\s*major:\s*(\d+)\s*\n\s*minor:\s*(\d+)\s*\n\s*contract:\s*(\S+)\s*$`)

// releaseSeriesEntry is one clusterctl metadata.yaml releaseSeries item.
type releaseSeriesEntry struct {
	major, minor int
	contract     string
}

// parseReleaseSeries extracts the releaseSeries entries from a metadata.yaml's
// raw content, in file order.
func parseReleaseSeries(content []byte) ([]releaseSeriesEntry, error) {
	matches := releaseSeriesEntryRE.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no releaseSeries entries found")
	}
	entries := make([]releaseSeriesEntry, 0, len(matches))
	for _, m := range matches {
		major, _ := strconv.Atoi(string(m[1]))
		minor, _ := strconv.Atoi(string(m[2]))
		entries = append(entries, releaseSeriesEntry{major: major, minor: minor, contract: string(m[3])})
	}
	return entries, nil
}

// DefaultContract returns the contract value of the most recently added
// releaseSeries entry in the metadata.yaml at path, i.e. the value a new
// release should inherit unless explicitly overridden.
func DefaultContract(path string) (string, error) {
	content, err := os.ReadFile(path) //nolint:gosec // path is an operator-supplied CLI flag
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", path, err)
	}
	entries, err := parseReleaseSeries(content)
	if err != nil {
		return "", fmt.Errorf("parsing %s: %w", path, err)
	}
	return entries[len(entries)-1].contract, nil
}

// AppendReleaseSeries appends a new {major, minor, contract} entry to the
// metadata.yaml at path, unless an entry for that major.minor already exists
// (in which case it is a no-op, making this safe to re-run). The file's
// existing content (including its header comment) is preserved; the new
// entry is appended in the same style as existing ones.
func AppendReleaseSeries(path string, major, minor int, contract string) error {
	content, err := os.ReadFile(path) //nolint:gosec // path is an operator-supplied CLI flag
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	entries, err := parseReleaseSeries(content)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}
	for _, e := range entries {
		if e.major == major && e.minor == minor {
			return nil
		}
	}

	updated := append([]byte(nil), content...)
	if len(updated) > 0 && updated[len(updated)-1] != '\n' {
		updated = append(updated, '\n')
	}
	updated = append(updated, []byte(fmt.Sprintf("  - major: %d\n    minor: %d\n    contract: %s\n", major, minor, contract))...)

	if err := os.WriteFile(path, updated, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
