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

// Package changelog implements the version/branch validation and metadata.yaml
// automation used by the CHANGELOG-PR release trigger.
package changelog

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	semverRE        = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)
	releaseBranchRE = regexp.MustCompile(`^release-(\d+)\.(\d+)$`)
)

// Trigger describes the outcome of validating a version against the branch a
// CHANGELOG PR was merged into.
type Trigger struct {
	// ReleaseBranch is the branch the tag must be created on: the triggering
	// branch itself for a release-branch-targeted patch, or the new
	// "release-X.Y" branch for a main-targeted major/minor release.
	ReleaseBranch string
	// CreateReleaseBranch is true when ReleaseBranch does not exist yet and
	// must be created (main-targeted major/minor releases only).
	CreateReleaseBranch bool
	Major               int
	Minor               int
	Patch               int
}

// ParseVersion parses a "vMAJOR.MINOR.PATCH" string into its components.
func ParseVersion(version string) (major, minor, patch int, err error) {
	m := semverRE.FindStringSubmatch(version)
	if m == nil {
		return 0, 0, 0, fmt.Errorf("invalid semver version %q, expected vMAJOR.MINOR.PATCH", version)
	}
	major, _ = strconv.Atoi(m[1])
	minor, _ = strconv.Atoi(m[2])
	patch, _ = strconv.Atoi(m[3])
	return major, minor, patch, nil
}

// ValidateTrigger checks that version is valid semver and consistent with the
// rules for the branch a CHANGELOG PR targeted:
//   - targeting "main": version must be a major/minor release (patch == 0).
//   - targeting "release-X.Y": version's major.minor must equal X.Y exactly,
//     and patch must be > 0; this is a patch release.
//
// Any other target branch, or a version that violates these rules, is an
// error.
func ValidateTrigger(targetBranch, version string) (Trigger, error) {
	major, minor, patch, err := ParseVersion(version)
	if err != nil {
		return Trigger{}, err
	}

	if targetBranch == "main" {
		if patch != 0 {
			return Trigger{}, fmt.Errorf("version %q targets main but has a non-zero patch (%d); "+
				"main only accepts major/minor releases (patch must be 0)", version, patch)
		}
		return Trigger{
			ReleaseBranch:       fmt.Sprintf("release-%d.%d", major, minor),
			CreateReleaseBranch: true,
			Major:               major,
			Minor:               minor,
			Patch:               patch,
		}, nil
	}

	if bm := releaseBranchRE.FindStringSubmatch(targetBranch); bm != nil {
		branchMajor, _ := strconv.Atoi(bm[1])
		branchMinor, _ := strconv.Atoi(bm[2])
		if major != branchMajor || minor != branchMinor {
			return Trigger{}, fmt.Errorf("version %q does not match branch %q; "+
				"a release branch can only cut patch releases for its own major.minor", version, targetBranch)
		}
		if patch == 0 {
			return Trigger{}, fmt.Errorf("version %q has patch 0 but targets release branch %q; "+
				"release branches can only cut patch releases (patch must be > 0)", version, targetBranch)
		}
		return Trigger{
			ReleaseBranch:       targetBranch,
			CreateReleaseBranch: false,
			Major:               major,
			Minor:               minor,
			Patch:               patch,
		}, nil
	}

	return Trigger{}, fmt.Errorf("unsupported target branch %q; CHANGELOG PRs may only target "+
		"\"main\" or a \"release-X.Y\" branch", targetBranch)
}
