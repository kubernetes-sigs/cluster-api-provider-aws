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

import "testing"

func TestGroupKey(t *testing.T) {
	tests := []struct {
		name    string
		tags    map[string]string
		wantKey string
		wantOK  bool
	}{
		{
			name:    "all tags present",
			tags:    map[string]string{"distribution": "ubuntu", "distribution_version": "22.04", "kubernetes_version": "v1.36.3"},
			wantKey: "ubuntu|22.04|v1.36.3",
			wantOK:  true,
		},
		{
			name:   "missing distribution",
			tags:   map[string]string{"distribution_version": "22.04", "kubernetes_version": "v1.36.3"},
			wantOK: false,
		},
		{
			name:   "no tags",
			tags:   nil,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, ok, reason := groupKey(AMI{ID: "ami-1", Tags: tt.tags})
			if ok != tt.wantOK {
				t.Fatalf("groupKey() ok = %v, want %v (reason: %q)", ok, tt.wantOK, reason)
			}
			if ok && key != tt.wantKey {
				t.Errorf("groupKey() = %q, want %q", key, tt.wantKey)
			}
			if !ok && reason == "" {
				t.Errorf("groupKey() reason is empty when ok = false")
			}
		})
	}
}

func TestBuildTimestamp(t *testing.T) {
	tests := []struct {
		name   string
		tags   map[string]string
		wantTs int64
		wantOK bool
	}{
		{name: "valid", tags: map[string]string{"build_timestamp": "1785150077"}, wantTs: 1785150077, wantOK: true},
		{name: "missing", tags: nil, wantOK: false},
		{name: "not numeric", tags: map[string]string{"build_timestamp": "not-a-number"}, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, ok := buildTimestamp(AMI{ID: "ami-1", Tags: tt.tags})
			if ok != tt.wantOK {
				t.Fatalf("buildTimestamp() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && ts != tt.wantTs {
				t.Errorf("buildTimestamp() = %d, want %d", ts, tt.wantTs)
			}
		})
	}
}

func TestFindDuplicateAMIs(t *testing.T) {
	baseTags := map[string]string{"distribution": "ubuntu", "distribution_version": "22.04", "kubernetes_version": "v1.36.3"}
	withTags := func(overrides map[string]string) map[string]string {
		tags := make(map[string]string, len(baseTags)+len(overrides))
		for k, v := range baseTags {
			tags[k] = v
		}
		for k, v := range overrides {
			tags[k] = v
		}
		return tags
	}

	amis := []AMI{
		{Region: "us-east-1", ID: "ami-oldest", Tags: withTags(map[string]string{"build_timestamp": "50"})},
		{Region: "us-east-1", ID: "ami-newer", Tags: withTags(map[string]string{"build_timestamp": "200"})},
		{Region: "us-east-1", ID: "ami-older", Tags: withTags(map[string]string{"build_timestamp": "100"})},
		{Region: "us-east-1", ID: "ami-singleton", Tags: map[string]string{"distribution": "flatcar", "distribution_version": "stable", "kubernetes_version": "v1.36.3", "build_timestamp": "1"}},
		{Region: "us-east-1", ID: "ami-missing-tag", Tags: map[string]string{"distribution_version": "22.04", "kubernetes_version": "v1.36.3", "build_timestamp": "1"}},
		{Region: "us-east-1", ID: "ami-missing-ts", Tags: baseTags},
		// Same group key/timestamp scheme, but a different region: must not
		// be grouped with the us-east-1 AMIs.
		{Region: "eu-west-1", ID: "ami-other-region", Tags: withTags(map[string]string{"build_timestamp": "999"})},
	}

	report := FindDuplicateAMIs(amis)

	if len(report.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(report.Groups))
	}
	got := report.Groups[0]
	if got.Region != "us-east-1" {
		t.Errorf("group region = %q, want %q", got.Region, "us-east-1")
	}
	if got.Keep.ID != "ami-newer" {
		t.Errorf("Keep = %q, want %q", got.Keep.ID, "ami-newer")
	}
	if len(got.Duplicates) != 2 {
		t.Fatalf("Duplicates has %d entries, want 2", len(got.Duplicates))
	}

	wantUngroupableIDs := map[string]bool{"ami-missing-tag": true, "ami-missing-ts": true}
	if len(report.Ungroupable) != len(wantUngroupableIDs) {
		t.Fatalf("ungroupable has %d entries, want %d", len(report.Ungroupable), len(wantUngroupableIDs))
	}
	for _, u := range report.Ungroupable {
		if !wantUngroupableIDs[u.AMI.ID] {
			t.Errorf("unexpected ungroupable AMI %q", u.AMI.ID)
		}
		if u.Reason == "" {
			t.Errorf("ungroupable AMI %q has empty reason", u.AMI.ID)
		}
	}
}

func TestEntries(t *testing.T) {
	report := &DuplicatesReport{
		Groups: []DuplicateGroup{
			{
				Region:   "us-east-1",
				GroupKey: "ubuntu|22.04|v1.36.3",
				Keep:     AMI{ID: "ami-newer", Tags: map[string]string{"build_timestamp": "200"}},
				Duplicates: []AMI{
					{ID: "ami-older", Tags: map[string]string{"build_timestamp": "100"}},
					{ID: "ami-oldest", Tags: map[string]string{"build_timestamp": "50"}},
				},
			},
		},
		Ungroupable: []UngroupableAMI{
			{AMI: AMI{ID: "ami-missing-tag"}, Reason: `missing "distribution" tag`},
		},
	}

	entries := report.Entries()
	if len(entries) != 4 {
		t.Fatalf("got %d entries, want 4", len(entries))
	}

	byID := make(map[string]Entry, len(entries))
	for _, e := range entries {
		byID[e.AMI.ID] = e
	}

	if got := byID["ami-newer"]; got.Status != StatusKeep || got.BuildTimestamp != "200" || got.GroupKey != "ubuntu|22.04|v1.36.3" {
		t.Errorf("ami-newer entry = %+v, want Status=KEEP BuildTimestamp=200 GroupKey=ubuntu|22.04|v1.36.3", got)
	}
	if got := byID["ami-older"]; got.Status != StatusDuplicate || got.BuildTimestamp != "100" {
		t.Errorf("ami-older entry = %+v, want Status=DUPLICATE BuildTimestamp=100", got)
	}
	if got := byID["ami-oldest"]; got.Status != StatusDuplicate {
		t.Errorf("ami-oldest entry = %+v, want Status=DUPLICATE", got)
	}
	if got := byID["ami-missing-tag"]; got.Status != StatusUngroupable || got.Reason == "" {
		t.Errorf("ami-missing-tag entry = %+v, want Status=UNGROUPABLE with a reason", got)
	}
}
