/*
Copyright 2020 The Kubernetes Authors.

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

package cip

import (
	"fmt"

	"github.com/spf13/cobra"

	"sigs.k8s.io/promo-tools/v3/internal/legacy/cli"
)

// auditCmd represents the base command when called without any subcommands
// TODO: Update command description.
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run the image auditor",
	Long: `cip audit - Image auditor

Start an audit server that responds to Pub/Sub push events.
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.RunAuditCmd(auditOpts); err != nil {
			return fmt.Errorf("run `cip audit`: %w", err)
		}
		return nil
	},
}

var auditOpts = &cli.AuditOptions{}

func init() {
	auditCmd.PersistentFlags().StringVar(
		&auditOpts.ProjectID,
		"project",
		"",
		"GCP project name (used for labeling error reporting logs in GCP)",
	)

	auditCmd.PersistentFlags().StringVar(
		&auditOpts.RepoURL,
		"url",
		"",
		"repository URL for promoter manifests",
	)

	auditCmd.PersistentFlags().StringVar(
		&auditOpts.RepoBranch,
		"branch",
		"",
		"git branch of the promoter manifest repo to checkout",
	)

	auditCmd.PersistentFlags().StringVar(
		&auditOpts.ManifestPath,
		"path",
		"",
		"manifest path (relative to the root of promoter manifest repo)",
	)

	auditCmd.PersistentFlags().BoolVar(
		&auditOpts.Verbose,
		"verbose",
		auditOpts.Verbose,
		"include extra logging information",
	)

	CipCmd.AddCommand(auditCmd)
}
