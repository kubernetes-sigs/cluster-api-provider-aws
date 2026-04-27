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
	GitHubTagsURL = "https://api.github.com/repos/kubernetes/kubernetes/tags"
	// LatestVersionCount is the number of minor versions tracked per the CAPA AMI
	// publication policy: https://cluster-api-aws.sigs.k8s.io/topics/images/built-amis#ami-publication-policy
	LatestVersionCount = 3
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

// DetectVersionsForMinor detects all stable patch releases for one Kubernetes minor version.
//
// Arguments:
// minor: Target Kubernetes minor version in MAJOR.MINOR format (for example, "1.36").
// token: GitHub token used to query Kubernetes tags API.
//
// Returns:
// *MinorVersion with the requested minor and sorted patch versions, or an error.
func DetectVersionsForMinor(minor, token string) (*MinorVersion, error) {
	result, err := DetectK8sVersions(token, minor)
	if err != nil {
		return nil, err
	}
	if len(result.Versions) == 0 {
		return nil, fmt.Errorf("no stable releases found for Kubernetes %s", minor)
	}
	return &result.Versions[0], nil
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


// MinorGreater compares two MAJOR.MINOR versions and reports whether a is newer than b.
//
// Arguments:
// a: Left minor version in MAJOR.MINOR format.
// b: Right minor version in MAJOR.MINOR format.
//
// Returns:
// true when version a is greater than version b; otherwise false.
func MinorGreater(a, b string) bool {
	aParts := strings.SplitN(a, ".", 2)
	bParts := strings.SplitN(b, ".", 2)
	if len(aParts) != 2 || len(bParts) != 2 {
		return a > b
	}
	aMaj, _ := strconv.Atoi(aParts[0])
	bMaj, _ := strconv.Atoi(bParts[0])
	if aMaj != bMaj {
		return aMaj > bMaj
	}
	aMin, _ := strconv.Atoi(aParts[1])
	bMin, _ := strconv.Atoi(bParts[1])
	return aMin > bMin
}

// TopMinors selects the highest minor versions and returns them in descending order.
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

// PatchGreater compares two MAJOR.MINOR.PATCH versions and reports whether a is newer than b.
//
// Arguments:
// a: Left patch version in MAJOR.MINOR.PATCH format.
// b: Right patch version in MAJOR.MINOR.PATCH format.
//
// Returns:
// true when version a is greater than version b; otherwise false.
func PatchGreater(a, b string) bool {
	aParts := strings.SplitN(a, ".", 3)
	bParts := strings.SplitN(b, ".", 3)
	if len(aParts) != 3 || len(bParts) != 3 {
		return a > b
	}
	for i := range 3 {
		av, _ := strconv.Atoi(aParts[i])
		bv, _ := strconv.Atoi(bParts[i])
		if av != bv {
			return av > bv
		}
	}
	return false
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
