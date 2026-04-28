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

// Package k8srelease provides utilities for detecting stable Kubernetes
// releases from the kubernetes/kubernetes GitHub repository.
package k8srelease

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// GitHubTagsURL is the endpoint used to fetch Kubernetes release tags.
	GitHubTagsURL = "https://api.github.com/repos/kubernetes/kubernetes/tags"
	// LatestVersionCount is the number of minor versions tracked per the CAPA AMI
	// publication policy: https://cluster-api-aws.sigs.k8s.io/topics/images/built-amis#ami-publication-policy
	LatestVersionCount = 3
	// CapaModeToken is the CLI token selecting CAPA policy mode.
	CapaModeToken = "capa"
)

// StableTagRe matches only stable release tags like v1.35.1.
// Pre-release suffixes (alpha, beta, rc) are intentionally not matched.
var StableTagRe = regexp.MustCompile(`^v1\.(\d+)\.(\d+)$`)

// GitHubTag is one entry from the GitHub tags API response.
type GitHubTag struct {
	Name string `json:"name"`
}

// MinorVersion groups all patch releases under a single Kubernetes minor version.
type MinorVersion struct {
	Minor   string   `json:"minor"`
	Patches []string `json:"patches"`
}

// SupportedVersions is the structured result of a CAPA-policy version query.
type SupportedVersions struct {
	GeneratedAt string         `json:"generated_at"`
	Versions    []MinorVersion `json:"versions"`
}

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
	if requestedMinors[0] == CapaModeToken {
		if len(requestedMinors) > 1 {
			return nil, fmt.Errorf("invalid inputs to detect k8s releases version(s)")
		}
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
			return nil, fmt.Errorf("unknown Kubernetes version %q", minor)
		}
		if _, duplicate := seen[minor]; duplicate {
			return nil, fmt.Errorf("invalid inputs to detect k8s releases version(s)")
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

// ToTable converts SupportedVersions to a table representation for CLI output.
//
// Arguments:
// s: SupportedVersions receiver containing minor and patch data.
//
// Returns:
// *metav1.Table where each row contains one minor version and comma-separated patches.
func (s *SupportedVersions) ToTable() *metav1.Table {
	table := &metav1.Table{
		TypeMeta: metav1.TypeMeta{
			APIVersion: metav1.SchemeGroupVersion.String(),
			Kind:       "Table",
		},
		ColumnDefinitions: []metav1.TableColumnDefinition{
			{Name: "Minor Version", Type: "string"},
			{Name: "Patch Versions", Type: "string"},
		},
	}
	for _, v := range s.Versions {
		table.Rows = append(table.Rows, metav1.TableRow{
			Cells: []interface{}{v.Minor, strings.Join(v.Patches, ", ")},
		})
	}
	return table
}

// FetchAllTags retrieves all tags from the Kubernetes GitHub repository using paginated API requests.
//
// Arguments:
// token: Optional GitHub token used to increase API rate limits.
//
// Returns:
// []string with all tag names in API order, or an error if any request/parse step fails.
func FetchAllTags(token string) ([]string, error) {
	var names []string
	client := &http.Client{Timeout: 30 * time.Second}

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s?per_page=100&page=%d", GitHubTagsURL, page)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching tags page %d: %w", page, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading response body page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status %d on page %d: %s", resp.StatusCode, page, body)
		}

		var tags []GitHubTag
		if err := json.Unmarshal(body, &tags); err != nil {
			return nil, fmt.Errorf("decoding tags page %d: %w", page, err)
		}

		if len(tags) == 0 {
			break
		}

		for _, t := range tags {
			names = append(names, t.Name)
		}
	}

	return names, nil
}

// FilterStableTags filters tags down to stable Kubernetes patch releases only.
//
// Arguments:
// tags: Raw Kubernetes tag names returned by the GitHub API.
//
// Returns:
// []string containing only stable tags that match MAJOR.MINOR.PATCH format.
func FilterStableTags(tags []string) []string {
	var stable []string
	for _, tag := range tags {
		if StableTagRe.MatchString(tag) {
			stable = append(stable, tag)
		}
	}
	return stable
}

// GroupByMinor groups stable patch tags by MAJOR.MINOR key.
//
// Arguments:
// tags: Stable tags expected with a leading "v" prefix (for example, "v1.35.4").
//
// Returns:
// map[string][]string keyed by MAJOR.MINOR without "v", with patch versions as values.
func GroupByMinor(tags []string) map[string][]string {
	groups := make(map[string][]string)
	for _, tag := range tags {
		ver := strings.TrimPrefix(tag, "v")
		parts := strings.SplitN(ver, ".", 3)
		if len(parts) != 3 {
			continue
		}
		minor := parts[0] + "." + parts[1]
		groups[minor] = append(groups[minor], ver)
	}
	return groups
}

func versionGreater(a, b string, parts int) bool {
	aParts := strings.SplitN(a, ".", parts)
	bParts := strings.SplitN(b, ".", parts)

	if len(aParts) != parts || len(bParts) != parts {
		// fallback to string comparison if format is unexpected
		return a > b
	}

	for i := range parts {
		av, _ := strconv.Atoi(aParts[i])
		bv, _ := strconv.Atoi(bParts[i])

		if av != bv {
			return av > bv
		}
	}

	return false
}

// MinorGreater compares two MAJOR.MINOR versions and reports whether a is newer than b.
//
// Arguments:
// a: Left minor version in MAJOR.MINOR format.
// b: Right minor version in MAJOR.MINOR format.
//
// Returns:
// true when version a is greater than version b; otherwise false.
func MinorGreater(a, b string) bool {
	return versionGreater(a, b, 2)
}

// PatchGreater compares two MAJOR.MINOR.PATCH versions and reports whether a is newer than b.
//
// Arguments:
// a: Left patch version in MAJOR.MINOR.PATCH format.
// b: Right patch version in MAJOR.MINOR.PATCH format.
//
// Returns:
// true when version a is greater than version b; otherwise false.
func PatchGreater(a, b string) bool {
	return versionGreater(a, b, 3)
}

// TopMinors selects the latest/highest minor versions and returns them in descending order.
//
// Arguments:
// groups: Map keyed by MAJOR.MINOR containing grouped patch versions.
// n: Number of latest minor versions to return.
//
// Returns:
// []string of up to n latest MAJOR.MINOR versions sorted descending.
func TopMinors(groups map[string][]string, n int) []string {
	minors := make([]string, 0, len(groups))
	for m := range groups {
		minors = append(minors, m)
	}
	sort.Slice(minors, func(i, j int) bool {
		return MinorGreater(minors[i], minors[j])
	})
	if n > len(minors) {
		n = len(minors)
	}
	return minors[:n]
}

// SortPatchesDesc sorts patch versions in-place from newest to oldest.
//
// Arguments:
// patches: Slice of MAJOR.MINOR.PATCH versions to sort.
//
// Returns:
// None. The input slice is modified in place.
func SortPatchesDesc(patches []string) {
	sort.Slice(patches, func(i, j int) bool {
		return PatchGreater(patches[i], patches[j])
	})
}
