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

package ami

import (
	"fmt"
	"strconv"
)

const (
	// TagDistribution is the EC2 tag key holding the AMI's base OS distribution.
	TagDistribution = "distribution"
	// TagDistributionVersion is the EC2 tag key holding the distribution version.
	TagDistributionVersion = "distribution_version"
	// TagKubernetesVersion is the EC2 tag key holding the Kubernetes version.
	TagKubernetesVersion = "kubernetes_version"
	// TagBuildTimestamp is the EC2 tag key holding the build's timestamp,
	// used to pick the newest AMI within a duplicate group.
	TagBuildTimestamp = "build_timestamp"
)

// AMI is a minimal, SDK-independent view of an EC2 AMI used for duplicate detection.
type AMI struct {
	Region string            `json:"region"`
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Tags   map[string]string `json:"tags,omitempty"`
}

// DuplicateGroup is a set of AMIs in one region that share the same
// distribution/distribution_version/kubernetes_version tags.
type DuplicateGroup struct {
	Region     string `json:"region"`
	GroupKey   string `json:"groupKey"`
	Keep       AMI    `json:"keep"`
	Duplicates []AMI  `json:"duplicates"`
}

// UngroupableAMI is an AMI that can't be safely deduplicated: it is missing
// one of the grouping tags, or its build_timestamp tag is missing/invalid.
type UngroupableAMI struct {
	AMI    AMI    `json:"ami"`
	Reason string `json:"reason"`
}

// DuplicatesReport is the result returned by FindDuplicateAMIs.
type DuplicatesReport struct {
	Groups      []DuplicateGroup `json:"groups"`
	Ungroupable []UngroupableAMI `json:"ungroupable,omitempty"`
}

// EntryStatus classifies a single flattened duplicate-detection result row.
type EntryStatus string

const (
	// StatusKeep marks the AMI kept from a duplicate group (highest build_timestamp).
	StatusKeep EntryStatus = "KEEP"
	// StatusDuplicate marks an AMI superseded by a StatusKeep AMI in the same group.
	StatusDuplicate EntryStatus = "DUPLICATE"
	// StatusUngroupable marks an AMI that couldn't be evaluated for duplication.
	StatusUngroupable EntryStatus = "UNGROUPABLE"
)

// Entry is one flattened, filterable row of a DuplicatesReport.
type Entry struct {
	Region         string      `json:"region"`
	GroupKey       string      `json:"groupKey,omitempty"`
	Status         EntryStatus `json:"status"`
	AMI            AMI         `json:"ami"`
	BuildTimestamp string      `json:"buildTimestamp,omitempty"`
	Reason         string      `json:"reason,omitempty"`
}

// Entries flattens the report into KEEP/DUPLICATE/UNGROUPABLE rows, one per AMI.
func (r *DuplicatesReport) Entries() []Entry {
	entries := make([]Entry, 0, len(r.Ungroupable))
	for _, group := range r.Groups {
		entries = append(entries, Entry{
			Region:         group.Region,
			GroupKey:       group.GroupKey,
			Status:         StatusKeep,
			AMI:            group.Keep,
			BuildTimestamp: group.Keep.Tags[TagBuildTimestamp],
		})
		for _, dup := range group.Duplicates {
			entries = append(entries, Entry{
				Region:         group.Region,
				GroupKey:       group.GroupKey,
				Status:         StatusDuplicate,
				AMI:            dup,
				BuildTimestamp: dup.Tags[TagBuildTimestamp],
			})
		}
	}
	for _, u := range r.Ungroupable {
		entries = append(entries, Entry{
			Region: u.AMI.Region,
			Status: StatusUngroupable,
			AMI:    u.AMI,
			Reason: u.Reason,
		})
	}
	return entries
}

// FindDuplicateAMIs groups amis by region and by their distribution,
// distribution_version and kubernetes_version tags, keeping only the AMI
// with the highest build_timestamp tag in each group. Any AMI that can't be
// grouped (missing tags, or a missing/invalid build_timestamp) is returned
// separately rather than silently dropped.
func FindDuplicateAMIs(amis []AMI) *DuplicatesReport {
	buckets := make(map[string][]AMI)
	var ungroupable []UngroupableAMI

	for _, a := range amis {
		key, ok, reason := groupKey(a)
		if !ok {
			ungroupable = append(ungroupable, UngroupableAMI{AMI: a, Reason: reason})
			continue
		}

		if _, ok := buildTimestamp(a); !ok {
			ungroupable = append(ungroupable, UngroupableAMI{
				AMI:    a,
				Reason: fmt.Sprintf("missing or invalid %q tag", TagBuildTimestamp),
			})
			continue
		}

		buckets[a.Region+"|"+key] = append(buckets[a.Region+"|"+key], a)
	}

	var groups []DuplicateGroup
	for _, bucketImages := range buckets {
		if len(bucketImages) < 2 {
			continue
		}

		keep := bucketImages[0]
		keepTs, _ := buildTimestamp(keep)
		duplicates := make([]AMI, 0, len(bucketImages)-1)
		for _, a := range bucketImages[1:] {
			ts, _ := buildTimestamp(a)
			if ts > keepTs {
				duplicates = append(duplicates, keep)
				keep, keepTs = a, ts
			} else {
				duplicates = append(duplicates, a)
			}
		}

		key, _, _ := groupKey(keep)
		groups = append(groups, DuplicateGroup{Region: keep.Region, GroupKey: key, Keep: keep, Duplicates: duplicates})
	}

	return &DuplicatesReport{Groups: groups, Ungroupable: ungroupable}
}

// groupKey returns the dedup group key for an AMI, built from its
// distribution, distribution_version and kubernetes_version tags. ok is
// false if any of those tags is missing, in which case reason explains why.
func groupKey(a AMI) (key string, ok bool, reason string) {
	distribution, found := a.Tags[TagDistribution]
	if !found {
		return "", false, fmt.Sprintf("missing %q tag", TagDistribution)
	}
	distributionVersion, found := a.Tags[TagDistributionVersion]
	if !found {
		return "", false, fmt.Sprintf("missing %q tag", TagDistributionVersion)
	}
	kubernetesVersion, found := a.Tags[TagKubernetesVersion]
	if !found {
		return "", false, fmt.Sprintf("missing %q tag", TagKubernetesVersion)
	}

	return distribution + "|" + distributionVersion + "|" + kubernetesVersion, true, ""
}

// buildTimestamp parses the AMI's build_timestamp tag as an integer.
func buildTimestamp(a AMI) (int64, bool) {
	value, found := a.Tags[TagBuildTimestamp]
	if !found {
		return 0, false
	}
	ts, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return ts, true
}
