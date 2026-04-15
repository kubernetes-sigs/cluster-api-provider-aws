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

package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestIsSymbolicLaunchTemplateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{name: "$Latest is symbolic", version: "$Latest", want: true},
		{name: "$Default is symbolic", version: "$Default", want: true},
		{name: "concrete version 1", version: "1", want: false},
		{name: "concrete version 42", version: "42", want: false},
		{name: "empty string", version: "", want: false},
		{name: "lowercase $latest is not symbolic", version: "$latest", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSymbolicLaunchTemplateVersion(tt.version); got != tt.want {
				t.Errorf("isSymbolicLaunchTemplateVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestLaunchTemplateNeedsUpdate(t *testing.T) {
	tests := []struct {
		name          string
		statusID      *string
		statusVersion *string
		ngID          *string
		ngVersion     *string
		want          bool
	}{
		{
			name: "no launch template in status or nodegroup — no update",
			want: false,
		},
		{
			name:     "status has ID but nodegroup has no LT — update needed (BYO LT being applied for first time)",
			statusID: aws.String("lt-12345"),
			want:     true,
		},
		{
			name:     "same ID, no version — no update",
			statusID: aws.String("lt-12345"),
			ngID:     aws.String("lt-12345"),
			want:     false,
		},
		{
			name:     "ID changed — update needed",
			statusID: aws.String("lt-99999"),
			ngID:     aws.String("lt-12345"),
			want:     true,
		},
		{
			name:          "same ID, concrete version unchanged — no update",
			statusID:      aws.String("lt-12345"),
			statusVersion: aws.String("3"),
			ngID:          aws.String("lt-12345"),
			ngVersion:     aws.String("3"),
			want:          false,
		},
		{
			name:          "same ID, concrete version changed — update needed",
			statusID:      aws.String("lt-12345"),
			statusVersion: aws.String("4"),
			ngID:          aws.String("lt-12345"),
			ngVersion:     aws.String("3"),
			want:          true,
		},
		{
			name:          "same ID, status version is $Latest — no update (symbolic version skipped)",
			statusID:      aws.String("lt-12345"),
			statusVersion: aws.String("$Latest"),
			ngID:          aws.String("lt-12345"),
			ngVersion:     aws.String("7"),
			want:          false,
		},
		{
			name:          "same ID, status version is $Default — no update (symbolic version skipped)",
			statusID:      aws.String("lt-12345"),
			statusVersion: aws.String("$Default"),
			ngID:          aws.String("lt-12345"),
			ngVersion:     aws.String("2"),
			want:          false,
		},
		{
			name:          "status has version but no ID, nodegroup has neither — no update",
			statusVersion: aws.String("5"),
			want:          false,
		},
		{
			name:      "status has no version, nodegroup has version — no update",
			statusID:  aws.String("lt-12345"),
			ngID:      aws.String("lt-12345"),
			ngVersion: aws.String("2"),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := launchTemplateNeedsUpdate(tt.statusID, tt.statusVersion, tt.ngID, tt.ngVersion); got != tt.want {
				t.Errorf("launchTemplateNeedsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
