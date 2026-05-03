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
	// K8sRepo is the OWNER/NAME slug of the upstream Kubernetes repository.
	K8sRepo = "kubernetes/kubernetes"
	// tagsPerPage is the page size used when listing tags via the GitHub API.
	tagsPerPage = 100
)

// StableTagRe matches only stable release tags like v1.35.1.
// Pre-release suffixes (alpha, beta, rc) are intentionally not matched.
var StableTagRe = regexp.MustCompile(`^v1\.(\d+)\.(\d+)$`)

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
	allTags, err := FetchAllTags(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("fetching tags: %w", err)
	}
	stableTags := FilterStableTags(allTags)
	if len(stableTags) == 0 {
		return nil, fmt.Errorf("no stable tags found")
	}
	patchesByMinor := GroupByMinor(stableTags)
	resolvedMinors, err := ResolveRequestedMinors(patchesByMinor, latestN, requestedMinors)
	if err != nil {
		return nil, err
	}
	return BuildSupportedVersions(patchesByMinor, resolvedMinors), nil
}

// ResolveRequestedMinors resolves requested minor inputs into final MAJOR.MINOR keys.
//
// Arguments:
// patchesByMinor: Map keyed by MAJOR.MINOR to stable patches.
// latestN: Number of latest minors to return when requestedMinors is empty.
// requestedMinors: Optional explicit MAJOR.MINOR inputs.
//
// Returns:
// Ordered minor keys to render, or an error for invalid/unknown/duplicate inputs.
func ResolveRequestedMinors(patchesByMinor map[string][]string, latestN int, requestedMinors []string) ([]string, error) {
	if len(requestedMinors) == 0 {
		return TopMinors(patchesByMinor, latestN), nil
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
			return nil, fmt.Errorf("duplicate Kubernetes version %q", minor)
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
		return "", fmt.Errorf("invalid version %q: expected format MAJOR.MINOR (e.g. 1.36)", raw)
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
// *SupportedVersions: a list of minor-to-patches version mapping.
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
		Versions: versions,
	}
}

// FetchAllTags retrieves all tags from the Kubernetes GitHub repository using
// the go-github client. Pagination is driven by the API's Link header rather
// than an "empty page" sentinel, and rate-limit / cancellation are handled by
// the supplied context.
//
// Arguments:
// ctx: Context controlling cancellation and deadlines for the API requests.
// token: Optional GitHub token used to increase API rate limits.
//
// Returns:
// []string with all tag names in API order, or an error if any request fails.
func FetchAllTags(ctx context.Context, token string) ([]string, error) {
	owner, repo, ok := strings.Cut(K8sRepo, "/")
	if !ok {
		return nil, fmt.Errorf("invalid K8sRepo %q: expected OWNER/NAME format", K8sRepo)
	}

	client := github.NewClient(newGitHubHTTPClient(ctx, token))

	var names []string
	opt := &github.ListOptions{PerPage: tagsPerPage}
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("listing %s tags: %w", K8sRepo, err)
		}
		for _, t := range tags {
			names = append(names, t.GetName())
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return names, nil
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

// parseSemver parses a version string with semver.ParseTolerant. A MAJOR.MINOR
// input (for example "1.35") is padded to MAJOR.MINOR.0 so that minor-only
// values can also be compared via the semver ordering.
func parseSemver(v string) (semver.Version, error) {
	if strings.Count(v, ".") == 1 {
		v += ".0"
	}
	return semver.ParseTolerant(v)
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
		a, errA := parseSemver(minors[i])
		b, errB := parseSemver(minors[j])
		if errA != nil || errB != nil {
			return minors[i] > minors[j]
		}
		return a.GT(b)
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
		a, errA := parseSemver(patches[i])
		b, errB := parseSemver(patches[j])
		if errA != nil || errB != nil {
			return patches[i] > patches[j]
		}
		return a.GT(b)
	})
}
