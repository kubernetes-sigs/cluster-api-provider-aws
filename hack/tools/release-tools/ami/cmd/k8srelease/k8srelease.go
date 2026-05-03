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

// Package k8srelease wires the `release-tool ami detect-k8s-release` cobra
// command. The business logic lives in
// release-tools/ami/k8srelease (imported here as detect).
package k8srelease

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	detect "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/ami/k8srelease"
	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/printer"
)

// defaultLatestVersionCount mirrors the CAPA AMI publication policy of tracking
// the latest three minor Kubernetes versions:
// https://cluster-api-aws.sigs.k8s.io/topics/images/built-amis#ami-publication-policy
const defaultLatestVersionCount = 3

// Cmd returns the `detect-k8s-release` cobra command.
func Cmd() *cobra.Command {
	var (
		token              string
		latestVersionCount int
		version            string
		output             string
	)

	cmd := &cobra.Command{
		Use:   "detect-k8s-release",
		Short: "Detect supported Kubernetes release versions",
		Long: `Query the kubernetes/kubernetes GitHub repository for stable release tags.

By default returns the latest 3 minor Kubernetes versions (the CAPA AMI build
policy). Use --latest-version-count to change the count, or --version to
request specific minors. The --version flag may be specified at most once;
to request multiple minors, pass them as a comma-separated list.`,
		Example: `  
  # Set GITHUB_TOKEN environment variable
  export GITHUB_TOKEN=ghp_xxx

  # Default: latest 3 minor versions, table output
  release-tool ami detect-k8s-release

  # Latest 4 minor versions, JSON output
  release-tool ami detect-k8s-release --latest-version-count 4 -o json

  # Specific minors (comma-separated, single --version flag)
  release-tool ami detect-k8s-release --version 1.34,1.30`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := assertNoDuplicateFlags(os.Args,
				[]string{"version"},
				[]string{"latest-version-count"},
				[]string{"token"},
				[]string{"output", "o"},
			); err != nil {
				return err
			}
			minors := splitVersions(version)
			if len(minors) > 0 && cmd.Flags().Changed("latest-version-count") {
				return fmt.Errorf("--version and --latest-version-count are mutually exclusive")
			}

			// Fall back to $GITHUB_TOKEN when --token isn't supplied. The flag's
			// default is intentionally empty so the token does not leak into
			// `--help` output or other places cobra prints flag defaults.
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}

			result, err := detect.DetectK8sVersions(cmd.Context(), token, latestVersionCount, minors)
			if err != nil {
				return err
			}

			p, err := printer.New(output, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			input := selectPrinterInput(output, result)
			return p.Print(input)
		},
	}

	cmd.Flags().StringVar(&token, "token", "",
		"GitHub personal access token; defaults to the GITHUB_TOKEN environment variable")
	cmd.Flags().IntVar(&latestVersionCount, "latest-version-count", defaultLatestVersionCount,
		"Number of latest minor Kubernetes versions to return (ignored when --version is set)")
	cmd.Flags().StringVar(&version, "version", "",
		"Comma-separated MAJOR.MINOR versions to fetch (e.g. 1.34,1.30). May be specified at most once. Overrides --latest-version-count.")
	cmd.Flags().StringVarP(&output, "output", "o", string(printer.TypeTable),
		"Output format: table, json, or yaml")
	return cmd
}

// selectPrinterInput returns the right input shape for the requested output
// format: table mode needs a *printer.Table, json/yaml accept the domain
// object directly. An empty format value defaults to table for parity with
// the cobra flag default.
func selectPrinterInput(format string, result *detect.SupportedVersions) interface{} {
	normalized := printer.Type(strings.ToLower(strings.TrimSpace(format)))
	if normalized == printer.TypeTable || normalized == "" {
		return supportedVersionsToTable(result)
	}
	return result
}

// supportedVersionsToTable converts a SupportedVersions into the generic
// *printer.Table consumed by the table printer.
func supportedVersionsToTable(result *detect.SupportedVersions) *printer.Table {
	rows := make([][]string, 0, len(result.Versions))
	for _, v := range result.Versions {
		rows = append(rows, []string{v.Minor, strings.Join(v.Patches, ", ")})
	}
	return &printer.Table{
		Columns: []string{"MINOR VERSION", "PATCH VERSIONS"},
		Rows:    rows,
	}
}

// splitVersions parses a comma-separated --version value into a clean slice,
// trimming whitespace and dropping empty entries.
func splitVersions(in string) []string {
	out := make([]string, 0)
	for _, v := range strings.Split(in, ",") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// assertNoDuplicateFlags scans the command-line args for repeated occurrences
// of any of the supplied flag alias groups and returns an error on the first
// duplicate found. Each group represents the long and (optional) short forms
// of a single flag (e.g. {"output", "o"}); using ANY of the aliases more than
// once across the group counts as a duplicate, so `--output json -o yaml` is
// rejected. We inspect args directly because pflag's scalar setters
// (StringVar, IntVar, ...) silently overwrite on repeat, so by the time RunE
// executes the duplication is no longer visible from the parsed flag value
// alone.
func assertNoDuplicateFlags(args []string, aliasGroups ...[]string) error {
	for _, aliases := range aliasGroups {
		count := 0
		for _, a := range args {
			for _, n := range aliases {
				prefix := "--"
				if len(n) == 1 {
					prefix = "-"
				}
				full := prefix + n
				if a == full || strings.HasPrefix(a, full+"=") {
					count++
					break
				}
			}
		}
		if count <= 1 {
			continue
		}
		primary := aliases[0]
		if primary == "version" {
			return fmt.Errorf("--version may only be specified once; use a comma-separated value for multiple inputs (e.g. 1.34,1.30)")
		}
		return fmt.Errorf("--%s may only be specified once", primary)
	}
	return nil
}
