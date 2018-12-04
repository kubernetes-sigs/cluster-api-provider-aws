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
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitBranch is the branch from which this binary was built
	GitBranch string
	// GitReleaseTag is the git tag from which this binary is released
	GitReleaseTag string
	// GitReleaseCommit is the commit corresponding to the GitReleaseTag
	GitReleaseCommit string
	// BuildTime is the time at which this binary was built
	BuildTime string
	// GitTreeState indicates if the git tree, from which this binary was built, was clean or dirty
	GitTreeState string
	// GitCommit is the git commit at which this binary binary was built
	GitCommit string
	// GitMajor is the major version of the release
	GitMajor string
	// GitMinor is the minor version of the release
	GitMinor string

	printLongHand bool
)

func isRepoAtRelease() bool {
	return GitTreeState == "clean" && GitReleaseCommit == GitCommit
}

func printShortDirtyVersionInfo() {
	fmt.Printf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q, GitReleaseCommit:%q, GitTreeState:%q\n",
		GitReleaseTag, GitMajor, GitMinor, GitReleaseCommit, GitTreeState)
}

func printShortCleanVersionInfo() {
	fmt.Printf("Version Info: GitReleaseTag: %q, MajorVersion: %q, MinorVersion:%q\n",
		GitReleaseTag, GitMajor, GitMinor)
}

func printVerboseVersionInfo() {
	fmt.Println("Version Info:")
	fmt.Printf("GitReleaseTag: %q, Major: %q, Minor: %q, GitRelaseCommit: %q\n", GitReleaseTag, GitMajor, GitMinor, GitReleaseCommit)
	fmt.Printf("Git Branch: %q\n", GitBranch)
	fmt.Printf("Git commit: %q\n", GitCommit)
	fmt.Printf("Git tree state: %q\n", GitTreeState)
}

// VersionCmd is the version command for the binary
func VersionCmd() *cobra.Command { // nolint
	vc := &cobra.Command{
		Use:   "version",
		Short: "Print version of this binary",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if printLongHand {
				printVerboseVersionInfo()
			} else if isRepoAtRelease() {
				printShortCleanVersionInfo()
			} else {
				printShortDirtyVersionInfo()
			}
		},
	}
	vc.Flags().BoolVarP(&printLongHand, "long", "l", false, "Print longhand version info")

	return vc
}
