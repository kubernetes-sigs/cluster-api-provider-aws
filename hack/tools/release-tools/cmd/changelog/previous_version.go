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
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/changelog"
)

// previousVersionCmd returns the `previous-version` cobra command.
func previousVersionCmd() *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "previous-version",
		Short: "Pick the tag changelog generation should start from, by semantic version",
		Long: `Reads a newline-separated list of tags from stdin (e.g. the output of
"git tag -l 'v*'") and prints the tag that changelog generation for
--version should start from, chosen by semantic version rather than git
reachability.

This matters because "git describe --abbrev=0" only considers reachable
tags: on main, the latest patch of the prior minor series may live only on
a diverged release branch and would be skipped, causing a changelog to
repeat part of the previous release's window.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var tags []string
			scanner := bufio.NewScanner(cmd.InOrStdin())
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					tags = append(tags, line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("reading tags from stdin: %w", err)
			}

			previous, err := changelog.PreviousVersion(tags, version)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), previous)
			return err
		},
	}

	cmd.Flags().StringVar(&version, "version", "", "the version being released, e.g. v2.10.0 (required)")
	_ = cmd.MarkFlagRequired("version")

	return cmd
}
