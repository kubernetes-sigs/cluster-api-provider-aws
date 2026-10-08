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

import (
	"fmt"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/changelog"
)

// validateTriggerCmd returns the `validate-trigger` cobra command.
func validateTriggerCmd() *cobra.Command {
	var (
		branch  string
		version string
	)

	cmd := &cobra.Command{
		Use:   "validate-trigger",
		Short: "Validate a CHANGELOG PR's version against the branch it targeted",
		Long: `Validate that a version added via a CHANGELOG/*.md file is allowed for the
branch that PR targeted:
  - "main" only accepts major/minor releases (patch must be 0).
  - a "release-X.Y" branch only accepts patch releases for its own X.Y.

On success, prints "key=value" lines (suitable for appending to
$GITHUB_OUTPUT) describing the release branch to tag on and whether it
still needs to be created.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			trigger, err := changelog.ValidateTrigger(branch, version)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "release_branch=%s\n", trigger.ReleaseBranch)
			fmt.Fprintf(out, "create_release_branch=%t\n", trigger.CreateReleaseBranch)
			fmt.Fprintf(out, "major=%d\n", trigger.Major)
			fmt.Fprintf(out, "minor=%d\n", trigger.Minor)
			fmt.Fprintf(out, "patch=%d\n", trigger.Patch)
			return nil
		},
	}

	cmd.Flags().StringVar(&branch, "branch", "", "branch the CHANGELOG PR was merged into (required)")
	cmd.Flags().StringVar(&version, "version", "", "version parsed from the CHANGELOG filename, e.g. v2.10.0 (required)")
	_ = cmd.MarkFlagRequired("branch")
	_ = cmd.MarkFlagRequired("version")

	return cmd
}
