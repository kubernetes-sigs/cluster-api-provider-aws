package versioninfo

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitBranch is the branch from which this binary was built
	GitBranch string
	// ReleaseTag is the git tag from which this binary is released
	ReleaseTag string
	// BuildTime is the time at which this binary was built
	BuildTime string
	// GitTreeState indicates if the git tree, from which this binary was built, was clean or dirty
	GitTreeState string
	// GitCommit is the git commit at which this binary binary was built
	GitCommit string
	// GitMajor is the major version of the release
	GitMajor string
	// GitMinor is the minor verison of the release
	GitMinor string
)

// VersionCmd is the version command for the binary
func VersionCmd() *cobra.Command { // nolint
	return &cobra.Command{
		Use:   "version",
		Short: "Print version of this binary",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Git Branch: %s\n", GitBranch)
			fmt.Printf("Git commit: %s\n", GitCommit)
			fmt.Printf("Release tag: %s\n", ReleaseTag)
			fmt.Printf("Git tree state: %s\n", GitTreeState)
			fmt.Printf("Git version: Major:%q, Minor:%q", GitMajor, GitMinor)
		},
	}
}
