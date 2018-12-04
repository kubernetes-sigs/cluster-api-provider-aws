/*
Copyright 2018 The Kubernetes Authors.

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

package versioninfo

import (
	"testing"
)

func TestIsRepoAtRelease(t *testing.T) {
	testCases := []struct {
		name                  string
		inputGitTreeState     string
		inputGitReleaseCommit string
		inputGitCommit        string
		expected              bool
	}{
		{
			name:                  "GitTreeState=clean and GitReleaseCommit==GitCommit",
			inputGitTreeState:     "clean",
			inputGitReleaseCommit: "1d5a033234f40b",
			inputGitCommit:        "1d5a033234f40b",
			expected:              true,
		},
		{
			name:                  "GitTreeState!=clean and GitReleaseCommit==GitCommit",
			inputGitTreeState:     "dirty",
			inputGitReleaseCommit: "1d5a033234f40b",
			inputGitCommit:        "1d5a033234f40b",
			expected:              false,
		},
		{
			name:                  "GitTreeState=clean and GitReleaseCommit!=GitCommit",
			inputGitTreeState:     "clean",
			inputGitReleaseCommit: "1d5a033234f40b",
			inputGitCommit:        "7adf033234f40b",
			expected:              false,
		},
		{
			name:                  "GitTreeState!=clean and GitReleaseCommit!=GitCommit",
			inputGitTreeState:     "dirty",
			inputGitReleaseCommit: "1d5a033234f40b",
			inputGitCommit:        "7adf033234f40b",
			expected:              false,
		},
	}

	for _, tc := range testCases {
		GitTreeState = tc.inputGitTreeState
		GitReleaseCommit = tc.inputGitReleaseCommit
		GitCommit = tc.inputGitCommit
		actual := isRepoAtRelease()

		if tc.expected != actual {
			t.Fatalf("isRepoAtRelease failed: [%s] want %v, got %v", tc.name, tc.expected, actual)
		}
	}
}
