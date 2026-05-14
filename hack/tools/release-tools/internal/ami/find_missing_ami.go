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

// Package ami provides utilities for computing missing CAPA AMI combinations.
package ami

import (
	"strings"

	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/ami/k8sreleases"
)

// PublishedAMI represents a single already-published AMI entry.
type PublishedAMI struct {
	KubernetesVersion string
	OS                string
	Region            string
}

// MissingAMI represents an expected AMI combination that has not been published.
type MissingAMI struct {
	KubernetesVersion string `json:"kubernetesVersion"`
	OS                string `json:"os"`
	Region            string `json:"region"`
}

// MissingAMIReport is the result returned by FindMissingAMIs.
type MissingAMIReport struct {
	Items []MissingAMI `json:"items"`
}

// FindMissingAMIs returns the combinations of (k8s patch version × OS × region)
// that are expected based on k8sVersions but absent from published.
//
// Arguments:
// k8sVersions: Kubernetes minor/patch versions that should have AMIs published.
// published:   Already-published AMIs to subtract from the expected set.
// osList:      Operating systems each patch version must cover.
// regions:     AWS regions each patch version/OS pair must cover.
//
// Both published entries and k8sVersions patches are normalised to a "v" prefix
// before comparison, so "1.36.0" and "v1.36.0" are treated as equivalent.
func FindMissingAMIs(
	k8sVersions *k8sreleases.SupportedVersions,
	published []PublishedAMI,
	osList []string,
	regions []string,
) *MissingAMIReport {
	// Build a lookup set of published AMIs keyed by "v<version>/<os>/<region>".
	lookup := make(map[string]struct{}, len(published))
	for _, p := range published {
		lookup[normalizeVersion(p.KubernetesVersion)+"/"+p.OS+"/"+p.Region] = struct{}{}
	}

	var missing []MissingAMI
	for _, mv := range k8sVersions.Versions {
		for _, patch := range mv.Patches {
			ver := normalizeVersion(patch)
			for _, osName := range osList {
				for _, region := range regions {
					if _, ok := lookup[ver+"/"+osName+"/"+region]; !ok {
						missing = append(missing, MissingAMI{
							KubernetesVersion: ver,
							OS:                osName,
							Region:            region,
						})
					}
				}
			}
		}
	}

	return &MissingAMIReport{Items: missing}
}

// normalizeVersion ensures a version string has a "v" prefix, trimming
// whitespace and any existing prefix before adding one.
func normalizeVersion(v string) string {
	return "v" + strings.TrimPrefix(strings.TrimSpace(v), "v")
}
