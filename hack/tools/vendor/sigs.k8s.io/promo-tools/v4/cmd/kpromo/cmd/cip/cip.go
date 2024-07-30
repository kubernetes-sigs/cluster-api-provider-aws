/*
Copyright 2019 The Kubernetes Authors.

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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/promo-tools/v4/internal/legacy/cli"
	promoter "sigs.k8s.io/promo-tools/v4/promoter/image"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
)

// CipCmd represents the base command when called without any subcommands
// TODO: Update command description.
var CipCmd = &cobra.Command{
	Use:   "cip",
	Short: "Promote images from a staging registry to production",
	Long: `cip - Kubernetes container image promoter

Promote images from a staging registry to production
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.RunPromoteCmd(runOpts); err != nil {
			return fmt.Errorf("run `cip run`: %w", err)
		}
		return nil
	},
}

var runOpts = &options.Options{}

// TODO: Function 'init' is too long (171 > 60) (funlen)
//
//nolint:funlen
func init() {
	CipCmd.PersistentFlags().BoolVar(
		&runOpts.Confirm,
		"confirm",
		runOpts.Confirm,
		"initiate a PRODUCTION image promotion",
	)

	// TODO: Move this into a default options function in pkg/promobot
	CipCmd.PersistentFlags().StringVar(
		&runOpts.Manifest,
		cli.PromoterManifestFlag,
		runOpts.Manifest,
		"the manifest file to load",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.ThinManifestDir,
		cli.PromoterThinManifestDirFlag,
		runOpts.ThinManifestDir,
		`recursively read in all manifests within a folder, but all manifests
MUST be 'thin' manifests named 'promoter-manifest.yaml', which are like regular
manifests but instead of defining the 'images: ...' field directly, the
'imagesPath' field must be defined that points to another YAML file containing
the 'images: ...' contents`,
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.UseProwManifestDiff,
		"use-prow-manifest-diff",
		runOpts.UseProwManifestDiff,
		"use only the latest diff for the manifest dir. Works only on prow.",
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.JSONLogSummary,
		"json-log-summary",
		runOpts.JSONLogSummary,
		"only log a JSON summary of important errors",
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.ParseOnly,
		"parse-only",
		runOpts.ParseOnly,
		"only check that the given manifest file is parsable as a Manifest",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.KeyFiles,
		"key-files",
		runOpts.KeyFiles,
		`CSV of service account key files that must be activated for the
promotion (<json-key-file-path>,...)`,
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.Snapshot,
		cli.PromoterSnapshotFlag,
		runOpts.Snapshot,
		"read all images in a repository and print to stdout",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SnapshotTag,
		"snapshot-tag",
		runOpts.SnapshotTag,
		"only snapshot images with the given tag",
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.MinimalSnapshot,
		"minimal-snapshot",
		runOpts.MinimalSnapshot,
		fmt.Sprintf(`(only works with '--%s' or '--%s') discard tagless images
from snapshot output if they are referenced by a manifest list`,
			cli.PromoterSnapshotFlag,
			cli.PromoterManifestBasedSnapshotOfFlag,
		),
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.OutputFormat,
		cli.PromoterOutputFlag,
		options.DefaultOptions.OutputFormat,
		fmt.Sprintf(`(only works with '--%s' or '--%s') choose output
format of the snapshot (allowed values: %q)`,
			cli.PromoterSnapshotFlag,
			cli.PromoterManifestBasedSnapshotOfFlag,
			promoter.AllowedOutputFormats,
		),
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SnapshotSvcAcct,
		"snapshot-service-account",
		runOpts.SnapshotSvcAcct,
		fmt.Sprintf(
			"service account to use for '--%s'",
			cli.PromoterSnapshotFlag,
		),
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.ManifestBasedSnapshotOf,
		cli.PromoterManifestBasedSnapshotOfFlag,
		runOpts.ManifestBasedSnapshotOf,
		fmt.Sprintf(`read all images in either '--%s' or '--%s' and print all
images that should be promoted to the given registry (assuming the given,
registry is empty); this is like '--%s', but instead of reading over the
network from a registry, it reads from the local manifests only`,
			cli.PromoterManifestFlag,
			cli.PromoterThinManifestDirFlag,
			cli.PromoterSnapshotFlag,
		),
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.UseServiceAcct,
		"use-service-account",
		runOpts.UseServiceAcct,
		"pass '--account=...' to all gcloud calls",
	)

	// This flag does nothing, but we don't want to remove it in case it breaks someone.
	// Instead, we just keep it hidden.
	unused := 0
	CipCmd.PersistentFlags().IntVar(
		&unused,
		"max-image-size",
		0,
		"the maximum image size (in MiB) allowed for promotion",
	)
	if err := CipCmd.PersistentFlags().MarkHidden("max-image-size"); err != nil {
		logrus.Infof("Failed to mark max-image-size flag as hidden: %v", err)
	}

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SignerAccount,
		"signer-account",
		options.DefaultOptions.SignerAccount,
		"service account to use as signing identity",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SignCheckIdentity,
		"certificate-identity",
		options.DefaultOptions.SignCheckIdentity,
		"identity to look for when verifying signatures",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SignCheckIssuer,
		"certificate-oidc-issuer",
		options.DefaultOptions.SignCheckIssuer,
		"OIDC Issuer that will be used for the signing identity, used for verify the images",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SignCheckIdentityRegexp,
		"certificate-identity-regexp",
		"",
		"A regular expression alternative to --certificate-identity. Accepts the Go regular expression syntax described at https://golang.org/s/re2syntax. Either --certificate-identity or --certificate-identity-regexp must be set for keyless flows.",
	)

	CipCmd.PersistentFlags().StringVar(
		&runOpts.SignCheckIssuerRegexp,
		"certificate-oidc-issuer-regexp",
		"",
		"A regular expression alternative to --certificate-oidc-issuer. Accepts the Go regular expression syntax described at https://golang.org/s/re2syntax.",
	)

	CipCmd.PersistentFlags().BoolVar(
		&runOpts.SignImages,
		"sign",
		options.DefaultOptions.SignImages,
		"when true, sign promoted images",
	)

	CipCmd.PersistentFlags().IntVar(
		&runOpts.MaxSignatureCopies,
		"max-signature-copies",
		options.DefaultOptions.MaxSignatureCopies,
		"maximum number of concurrent signature copies",
	)

	CipCmd.PersistentFlags().IntVar(
		&runOpts.MaxSignatureOps,
		"max-signature-ops",
		options.DefaultOptions.MaxSignatureOps,
		"maximum number of concurrent signature operations",
	)

	CipCmd.PersistentFlags().IntVar(
		&runOpts.SeverityThreshold,
		"vuln-severity-threshold",
		options.DefaultOptions.SeverityThreshold,
		`Using this flag will cause the promoter to only run the vulnerability
check. Found vulnerabilities at or above this threshold will result in the
vulnerability check failing [severity levels between 0 and 5; 0 - UNSPECIFIED,
1 - MINIMAL, 2 - LOW, 3 - MEDIUM, 4 - HIGH, 5 - CRITICAL]`,
	)
}
