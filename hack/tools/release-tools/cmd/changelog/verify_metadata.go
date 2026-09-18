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
	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/changelog"
)

// verifyMetadataCmd returns the `verify-metadata` cobra command.
func verifyMetadataCmd() *cobra.Command {
	var (
		changelogPath string
		metadataPath  string
		major, minor  int
	)

	cmd := &cobra.Command{
		Use:   "verify-metadata",
		Short: "Verify a main-targeted CHANGELOG PR includes the matching metadata.yaml update",
		Long: `Fails unless metadata.yaml already has a releaseSeries entry for
{major, minor} whose contract matches the "contract" front-matter value in
the given CHANGELOG file. Used as a pre-merge CI check for main-targeted
(major/minor) CHANGELOG PRs, which are expected to include their own
metadata.yaml update (see "notes-pr").`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return changelog.VerifyMetadata(changelogPath, metadataPath, major, minor)
		},
	}

	cmd.Flags().StringVar(&changelogPath, "changelog", "", "path to the CHANGELOG/vX.Y.Z.md file being verified (required)")
	cmd.Flags().StringVar(&metadataPath, "metadata", "metadata.yaml", "path to metadata.yaml")
	cmd.Flags().IntVar(&major, "major", 0, "major version of the release series (required)")
	cmd.Flags().IntVar(&minor, "minor", 0, "minor version of the release series (required)")
	_ = cmd.MarkFlagRequired("changelog")
	_ = cmd.MarkFlagRequired("major")
	_ = cmd.MarkFlagRequired("minor")

	return cmd
}
