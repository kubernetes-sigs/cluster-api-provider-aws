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

package k8srelease

import (
	"fmt"
	"strings"
	"time"
)

const CapaModeToken = "capa"

// DetectK8sVersions returns stable patch releases for CAPA mode or explicit minor inputs.
//
// Arguments:
// token: GitHub token used to query Kubernetes tags API.
// requestedMinors: Either "capa" or one/more MAJOR.MINOR values (for example, 1.36).
//
// Returns:
// Supported minor-to-patch mappings with generation timestamp, or an error.
func DetectK8sVersions(token string, requestedMinors ...string) (*SupportedVersions, error) {
	if len(requestedMinors) == 0 {
		return nil, fmt.Errorf("at least one requested minor is required")
	}
	allTags, err := FetchAllTags(token)
	if err != nil {
		return nil, fmt.Errorf("fetching tags: %w", err)
	}
	stableTags := FilterStableTags(allTags)
	if len(stableTags) == 0 {
		return nil, fmt.Errorf("no stable tags found")
	}
	patchesByMinor := GroupByMinor(stableTags)
	resolvedMinors, err := ResolveRequestedMinors(patchesByMinor, requestedMinors)
	if err != nil {
		return nil, err
	}
	return BuildSupportedVersions(patchesByMinor, resolvedMinors), nil
}

// ResolveRequestedMinors resolves requested minor inputs into final MAJOR.MINOR keys.
//
// Arguments:
// patchesByMinor: Map keyed by MAJOR.MINOR to stable patches.
// requestedMinors: User-requested values where first arg may be "capa".
//
// Returns:
// Ordered minor keys to render, or an error for invalid/unknown/duplicate inputs.
func ResolveRequestedMinors(patchesByMinor map[string][]string, requestedMinors []string) ([]string, error) {
	if len(requestedMinors) == 0 {
		return nil, fmt.Errorf("at least one requested minor is required")
	}
	if requestedMinors[0] == CapaModeToken {
		return TopMinors(patchesByMinor, LatestVersionCount), nil
	}
	seen := make(map[string]struct{}, len(requestedMinors))
	resolved := make([]string, 0, len(requestedMinors))
	for _, raw := range requestedMinors {
		minor, err := ParseMinorInput(raw)
		if err != nil {
			return nil, err
		}
		if _, exists := patchesByMinor[minor]; !exists {
			return nil, fmt.Errorf("no stable releases found for Kubernetes %s", minor)
		}
		if _, duplicate := seen[minor]; duplicate {
			return nil, fmt.Errorf("duplicate minor version %q is not allowed", minor)
		}
		seen[minor] = struct{}{}
		resolved = append(resolved, minor)
	}
	return resolved, nil
}

// ParseMinorInput normalizes and validates one MAJOR.MINOR input.
//
// Arguments:
// raw: Raw user-provided minor version, optionally with "v" prefix.
//
// Returns:
// Normalized MAJOR.MINOR value (without "v"), or an error.
func ParseMinorInput(raw string) (string, error) {
	normalized := strings.TrimPrefix(strings.TrimSpace(raw), "v")
	parts := strings.SplitN(normalized, ".", 3)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid version %q: expected format MAJOR.MINOR (e.g. 1.6) or %q", raw, CapaModeToken)
	}
	return normalized, nil
}

// BuildSupportedVersions builds a SupportedVersions struct from the given patches by minor.
//
// Arguments:
// patchesByMinor: Map keyed by MAJOR.MINOR to stable patches.
// minors: List of MAJOR.MINOR values to include in the result.
//
// Returns:
// *SupportedVersions: a list of minor-to-patches version mapping, or an error.
func BuildSupportedVersions(patchesByMinor map[string][]string, minors []string) *SupportedVersions {
	versions := make([]MinorVersion, 0, len(minors))
	for _, minor := range minors {
		patches := append([]string(nil), patchesByMinor[minor]...)
		SortPatchesDesc(patches)
		versions = append(versions, MinorVersion{
			Minor:   minor,
			Patches: patches,
		})
	}
	return &SupportedVersions{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Versions:    versions,
	}
}
