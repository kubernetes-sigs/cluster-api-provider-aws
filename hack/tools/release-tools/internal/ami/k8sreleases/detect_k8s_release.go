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

// Package k8sreleases provides utilities for detecting stable Kubernetes
// releases from the kubernetes/kubernetes GitHub repository.
package k8sreleases

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/google/go-github/v85/github"
	"golang.org/x/oauth2"
)

const (
	// k8sRepo is the OWNER/NAME slug of the upstream Kubernetes repository.
	k8sRepo = "kubernetes/kubernetes"
	// tagsPerPage is the page size used when listing tags via the GitHub API.
	tagsPerPage = 100
)

// stableTagRe matches only stable release tags like v1.35.1.
// Pre-release suffixes (alpha, beta, rc) are intentionally not matched.
var stableTagRe = regexp.MustCompile(`^v1\.(\d+)\.(\d+)$`)

// MinorVersion groups all patch releases under a single Kubernetes minor version.
type MinorVersion struct {
	Minor   string   `json:"minor"`
	Patches []string `json:"patches"`
}

// SupportedVersions is the structured result of a version query.
type SupportedVersions struct {
	Versions []MinorVersion `json:"versions"`
}

// DetectK8sVersions returns stable patch releases for either explicit minor
// inputs or the top latestN minors when no explicit minors are supplied.
//
// Arguments:
// ctx: Context controlling cancellation and deadlines for the GitHub API requests.
// token: GitHub token used to query Kubernetes tags API.
// latestN: Number of latest minor versions to return when requestedMinors is empty.
// requestedMinors: Optional explicit MAJOR.MINOR values (for example, 1.36).
//
// Returns:
// Supported minor-to-patch mappings with generation timestamp, or an error.
func DetectK8sVersions(ctx context.Context, token string, latestN int, requestedMinors []string) (*SupportedVersions, error) {
	if len(requestedMinors) == 0 && latestN <= 0 {
		return nil, fmt.Errorf("either requestedMinors or a positive latestN must be provided")
	}
	stableTags, err := fetchTags(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("fetching tags: %w", err)
	}
	if len(stableTags) == 0 {
		return nil, fmt.Errorf("no stable tags found")
	}
	patchesByMinor := groupByMinor(stableTags)
	resolvedMinors, err := resolveRequestedMinors(patchesByMinor, latestN, requestedMinors)
	if err != nil {
		return nil, err
	}
	return buildSupportedVersions(patchesByMinor, resolvedMinors), nil
}

// resolveRequestedMinors resolves requested minor inputs into final MAJOR.MINOR keys.
//
// Arguments:
// patchesByMinor: Map keyed by MAJOR.MINOR to stable patches.
// latestN: Number of latest minors to return when requestedMinors is empty.
// requestedMinors: Optional explicit MAJOR.MINOR inputs.
//
// Returns:
// Ordered minor keys to render, or an error for invalid/unknown/duplicate inputs.
func resolveRequestedMinors(patchesByMinor map[string][]string, latestN int, requestedMinors []string) ([]string, error) {
	if len(requestedMinors) == 0 {
		return topMinors(patchesByMinor, latestN), nil
	}
	seen := make(map[string]struct{}, len(requestedMinors))
	resolved := make([]string, 0, len(requestedMinors))
	for _, raw := range requestedMinors {
		minor, err := parseMinorInput(raw)
		if err != nil {
			return nil, err
		}
		if _, exists := patchesByMinor[minor]; !exists {
			return nil, fmt.Errorf("unknown Kubernetes version %q", minor)
		}
		if _, duplicate := seen[minor]; duplicate {
			return nil, fmt.Errorf("duplicate Kubernetes version %q", minor)
		}
		seen[minor] = struct{}{}
		resolved = append(resolved, minor)
	}
	return resolved, nil
}

// parseMinorInput normalizes and validates one MAJOR.MINOR input.
//
// Arguments:
// raw: Raw user-provided minor version, optionally with "v" prefix.
//
// Returns:
// Normalized MAJOR.MINOR value (without "v"), or an error.
func parseMinorInput(raw string) (string, error) {
	normalized := strings.TrimPrefix(strings.TrimSpace(raw), "v")
	parts := strings.SplitN(normalized, ".", 3)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("invalid version %q: expected format MAJOR.MINOR (e.g. 1.36)", raw)
	}
	return normalized, nil
}

// buildSupportedVersions builds a SupportedVersions struct from the given patches by minor.
//
// Arguments:
// patchesByMinor: Map keyed by MAJOR.MINOR to stable patches.
// minors: List of MAJOR.MINOR values to include in the result.
//
// Returns:
// *SupportedVersions: a list of minor-to-patches version mapping.
func buildSupportedVersions(patchesByMinor map[string][]string, minors []string) *SupportedVersions {
	versions := make([]MinorVersion, 0, len(minors))
	for _, minor := range minors {
		patches := append([]string(nil), patchesByMinor[minor]...)
		sortPatchesDesc(patches)
		versions = append(versions, MinorVersion{
			Minor:   minor,
			Patches: patches,
		})
	}
	return &SupportedVersions{
		Versions: versions,
	}
}

// fetchTags retrieves all stable release tags from the Kubernetes GitHub
// repository. It paginates through all results and filters out pre-release
// tags (alpha, beta, rc), returning only stable MAJOR.MINOR.PATCH tags.
//
// Arguments:
// ctx: Context controlling cancellation and deadlines for the API requests.
// token: Optional GitHub token used to increase API rate limits.
//
// Returns:
// []string of stable tag names, or an error if any request fails.
func fetchTags(ctx context.Context, token string) ([]string, error) {
	owner, repo, ok := strings.Cut(k8sRepo, "/")
	if !ok {
		return nil, fmt.Errorf("invalid k8sRepo %q: expected OWNER/NAME format", k8sRepo)
	}

	client := github.NewClient(newGitHubHTTPClient(ctx, token))

	var stable []string
	opt := &github.ListOptions{PerPage: tagsPerPage}
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("listing %s tags: %w", k8sRepo, err)
		}
		for _, t := range tags {
			if name := t.GetName(); stableTagRe.MatchString(name) {
				stable = append(stable, name)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return stable, nil
}

// newGitHubHTTPClient returns an HTTP client suitable for the GitHub API.
// When a token is supplied the client is wrapped in an oauth2 bearer-token
// transport; otherwise an unauthenticated client (lower rate limit) is used.
func newGitHubHTTPClient(ctx context.Context, token string) *http.Client {
	if token == "" {
		return nil
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	httpClient.Timeout = 30 * time.Second
	return httpClient
}

// groupByMinor groups stable patch tags by MAJOR.MINOR key.
//
// Arguments:
// tags: Stable tags expected with a leading "v" prefix (for example, "v1.35.4").
//
// Returns:
// map[string][]string keyed by MAJOR.MINOR without "v", with patch versions as values.
func groupByMinor(tags []string) map[string][]string {
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

// parseSemver parses a version string with semver.ParseTolerant. A MAJOR.MINOR
// input (for example "1.35") is padded to MAJOR.MINOR.0 so that minor-only
// values can also be compared via the semver ordering.
func parseSemver(v string) (semver.Version, error) {
	if strings.Count(v, ".") == 1 {
		v += ".0"
	}
	return semver.ParseTolerant(v)
}

// topMinors selects the latest/highest minor versions and returns them in descending order.
//
// Arguments:
// groups: Map keyed by MAJOR.MINOR containing grouped patch versions.
// n: Number of latest minor versions to return.
//
// Returns:
// []string of up to n latest MAJOR.MINOR versions sorted descending.
func topMinors(groups map[string][]string, n int) []string {
	minors := make([]string, 0, len(groups))
	for m := range groups {
		minors = append(minors, m)
	}
	sort.Slice(minors, func(i, j int) bool {
		a, _ := parseSemver(minors[i])
		b, _ := parseSemver(minors[j])
		return a.Compare(b) > 0
	})
	if n > len(minors) {
		n = len(minors)
	}
	return minors[:n]
}

// sortPatchesDesc sorts patch versions in-place from newest to oldest.
//
// Arguments:
// patches: Slice of MAJOR.MINOR.PATCH versions to sort.
//
// Returns:
// None. The input slice is modified in place.
func sortPatchesDesc(patches []string) {
	sort.Slice(patches, func(i, j int) bool {
		a, _ := parseSemver(patches[i])
		b, _ := parseSemver(patches[j])
		return a.Compare(b) > 0
	})
}
