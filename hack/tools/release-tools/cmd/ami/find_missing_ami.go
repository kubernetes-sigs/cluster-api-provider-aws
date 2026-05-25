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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	missingami "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/ami"
	detect "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/ami/k8sreleases"
	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/printer"
)

// findMissingAmiCmd returns the `find-missing-ami` cobra command.
func findMissingAmiCmd() *cobra.Command {
	var (
		token   string
		version string
		osFlag  string
		region  string
		output  string
	)

	cmd := &cobra.Command{
		Use:   "find-missing-ami",
		Short: "List expected AMIs that have not been published",
		Long: `Compute the AMIs that should be published (k8s versions × OS × regions)
but are not yet present in the published inventory.

Published AMIs are read from stdin in the JSON format produced by:
  clusterawsadm ami list --owner-id <ID> -o json

Kubernetes versions are validated against the CAPA AMI publication policy
(latest ` + fmt.Sprintf("%d", defaultLatestVersionCount) + ` minor versions). Requesting a version outside
this window returns an error.`,
		Example: `  
  # Check all policy versions for missing AMIs
    release-tool ami find-missing-ami \
      --token "$GITHUB_TOKEN" \
      --os ubuntu-24.04,ubuntu-22.04 \
      --region ap-southeast-2

  # Check a specific version (must be within the CAPA AMI publication policy)
    release-tool ami find-missing-ami \
      --token "$GITHUB_TOKEN" \
      --version 1.36 \
      --os ubuntu-24.04 \
      --region ap-southeast-2,us-east-1`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return assertNoDuplicateFlags(os.Args,
				[]string{"token"},
				[]string{"version"},
				[]string{"os"},
				[]string{"region"},
				[]string{"output", "o"},
			)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			policyVersions, err := detect.DetectK8sVersions(cmd.Context(), token, defaultLatestVersionCount, nil)
			if err != nil {
				return fmt.Errorf("fetching Kubernetes versions: %w", err)
			}

			k8sVersions := policyVersions
			if requested := splitVersions(version); len(requested) > 0 {
				k8sVersions, err = filterByPolicy(policyVersions, requested)
				if err != nil {
					return err
				}
			}

			published, err := readPublishedAMIs(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("reading published AMIs from stdin: %w", err)
			}

			report := missingami.FindMissingAMIs(
				k8sVersions,
				published,
				splitVersions(osFlag),
				splitVersions(region),
			)

			p, err := printer.New(output, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return p.Print(missingAMIsPrinterInput(output, report))
		},
	}

	cmd.Flags().StringVar(&token, "token", os.Getenv("GITHUB_TOKEN"),
		"GitHub personal access token; defaults to the GITHUB_TOKEN environment variable")
	cmd.Flags().StringVar(&version, "version", "",
		"Comma-separated MAJOR.MINOR versions to check (e.g. 1.36,1.35). Must be within the CAPA AMI publication policy.")
	cmd.Flags().StringVar(&osFlag, "os", "ubuntu-24.04,ubuntu-22.04,flatcar-stable",
		"Comma-separated list of operating systems to check")
	cmd.Flags().StringVar(&region, "region", "ap-southeast-2",
		"Comma-separated list of AWS regions to check")
	cmd.Flags().StringVarP(&output, "output", "o", string(printer.TypeTable),
		"Output format: table, json, or yaml")

	return cmd
}

// filterByPolicy validates that each requested minor is within the CAPA AMI
// publication policy window and returns a filtered SupportedVersions containing
// only those minors.
func filterByPolicy(policy *detect.SupportedVersions, requested []string) (*detect.SupportedVersions, error) {
	policySet := make(map[string]detect.MinorVersion, len(policy.Versions))
	policyMinors := make([]string, 0, len(policy.Versions))
	for _, mv := range policy.Versions {
		policySet[mv.Minor] = mv
		policyMinors = append(policyMinors, mv.Minor)
	}

	filtered := make([]detect.MinorVersion, 0, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, raw := range requested {
		minor, err := detect.ParseMinorInput(raw)
		if err != nil {
			return nil, err
		}
		if _, dup := seen[minor]; dup {
			return nil, fmt.Errorf("duplicate version %q", raw)
		}
		seen[minor] = struct{}{}

		mv, ok := policySet[minor]
		if !ok {
			return nil, fmt.Errorf("version %q is outside the CAPA AMI publication policy; supported: %s",
				raw, strings.Join(policyMinors, ", "))
		}
		filtered = append(filtered, mv)
	}
	return &detect.SupportedVersions{Versions: filtered}, nil
}

// readPublishedAMIs parses the clusterawsadm ami list -o json output from r.
// Returns an empty list when r is an interactive terminal or contains no data.
func readPublishedAMIs(r io.Reader) ([]missingami.PublishedAMI, error) {
	// If r is an interactive terminal, treat as empty to avoid blocking.
	if f, ok := r.(*os.File); ok {
		if stat, err := f.Stat(); err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
			return nil, nil
		}
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}

	var list struct {
		Items []struct {
			Spec struct {
				KubernetesVersion string `json:"kubernetesVersion"`
				OS                string `json:"os"`
				Region            string `json:"region"`
			} `json:"spec"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parsing AMI JSON: %w", err)
	}

	result := make([]missingami.PublishedAMI, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, missingami.PublishedAMI{
			KubernetesVersion: item.Spec.KubernetesVersion,
			OS:                item.Spec.OS,
			Region:            item.Spec.Region,
		})
	}
	return result, nil
}

// missingAMIsPrinterInput returns a *printer.Table for table output or the
// domain report directly for json/yaml.
func missingAMIsPrinterInput(format string, report *missingami.MissingAMIReport) interface{} {
	normalized := printer.Type(strings.ToLower(strings.TrimSpace(format)))
	if normalized == printer.TypeTable || normalized == "" {
		return missingAMIsToTable(report)
	}
	return report
}

// missingAMIsToTable converts a MissingAMIReport into a *printer.Table.
func missingAMIsToTable(report *missingami.MissingAMIReport) *printer.Table {
	rows := make([][]string, 0, len(report.Items))
	for _, item := range report.Items {
		rows = append(rows, []string{item.KubernetesVersion, item.OS, item.Region})
	}
	return &printer.Table{
		Columns: []string{"KUBERNETES VERSION", "OS", "REGION"},
		Rows:    rows,
	}
}
