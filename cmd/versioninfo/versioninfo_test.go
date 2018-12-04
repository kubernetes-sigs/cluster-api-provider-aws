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
