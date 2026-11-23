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

import "fmt"

type semver struct {
	major, minor, patch int
	raw                 string
}

func (v semver) less(o semver) bool {
	if v.major != o.major {
		return v.major < o.major
	}
	if v.minor != o.minor {
		return v.minor < o.minor
	}
	return v.patch < o.patch
}

// PreviousVersion returns the tag that changelog generation should start
// from for version, chosen by semantic version rather than git reachability
// (a tag from a diverged release branch, e.g. the latest patch of the prior
// minor series, may not be reachable from the commit being released):
//   - if version's patch is 0 (a main-targeted major/minor release): the
//     greatest tag with a strictly smaller (major, minor) — the end of the
//     prior release series.
//   - if version's patch is > 0 (a release-branch-targeted patch release):
//     the greatest tag with the same major.minor and a strictly smaller
//     patch (this naturally includes the series' own X.Y.0 tag).
//
// tags may contain non-semver or malformed entries, which are ignored.
func PreviousVersion(tags []string, version string) (string, error) {
	major, minor, patch, err := ParseVersion(version)
	if err != nil {
		return "", err
	}

	var best *semver
	for _, t := range tags {
		tMajor, tMinor, tPatch, err := ParseVersion(t)
		if err != nil {
			continue
		}
		candidate := semver{major: tMajor, minor: tMinor, patch: tPatch, raw: t}

		var isCandidate bool
		if patch == 0 {
			isCandidate = tMajor < major || (tMajor == major && tMinor < minor)
		} else {
			isCandidate = tMajor == major && tMinor == minor && tPatch < patch
		}
		if !isCandidate {
			continue
		}
		if best == nil || best.less(candidate) {
			best = &candidate
		}
	}

	if best == nil {
		return "", fmt.Errorf("no previous version found for %q among the given tags", version)
	}
	return best.raw, nil
}
