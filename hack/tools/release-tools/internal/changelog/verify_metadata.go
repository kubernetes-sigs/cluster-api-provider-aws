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
)

// VerifyMetadata checks that the metadata.yaml at metadataPath already has a
// releaseSeries entry for {major, minor} whose contract matches the
// "contract" front-matter value in the changelog at changelogPath. This is
// used as a pre-merge check for main-targeted (major/minor) CHANGELOG PRs,
// which are expected to include their own metadata.yaml update (see
// AppendReleaseSeries / the notes-pr command).
func VerifyMetadata(changelogPath, metadataPath string, major, minor int) error {
	changelogContent, err := os.ReadFile(changelogPath) //nolint:gosec // changelogPath is an operator-supplied CLI flag
	if err != nil {
		return fmt.Errorf("reading %s: %w", changelogPath, err)
	}
	wantContract, _ := ParseFrontMatter(changelogContent)
	if wantContract == "" {
		return fmt.Errorf("%s has no \"contract\" front-matter value", changelogPath)
	}

	metadataContent, err := os.ReadFile(metadataPath) //nolint:gosec // metadataPath is an operator-supplied CLI flag
	if err != nil {
		return fmt.Errorf("reading %s: %w", metadataPath, err)
	}
	entries, err := parseReleaseSeries(metadataContent)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", metadataPath, err)
	}

	for _, e := range entries {
		if e.major == major && e.minor == minor {
			if e.contract != wantContract {
				return fmt.Errorf("%s has a releaseSeries entry for %d.%d with contract %q, "+
					"but %s declares contract %q; they must match",
					metadataPath, major, minor, e.contract, changelogPath, wantContract)
			}
			return nil
		}
	}

	return fmt.Errorf("%s has no releaseSeries entry for %d.%d; "+
		"main-targeted CHANGELOG PRs must include this metadata.yaml update (see `make release-notes-pr`)",
		metadataPath, major, minor)
}
