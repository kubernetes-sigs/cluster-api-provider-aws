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

package gh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	"sigs.k8s.io/promo-tools/v4/gh2gcs"
	"sigs.k8s.io/release-sdk/gcli"
	"sigs.k8s.io/release-sdk/github"
)

// GHCmd represents the base command when called without any subcommands
var GHCmd = &cobra.Command{
	Use:           "gh --org kubernetes --repo release --bucket <bucket> --release-dir <release-dir> [--tags v0.0.0] [--include-prereleases] [--output-dir <temp-dir>] [--download-only] [--config <config-file>]",
	Short:         "Uploads GitHub releases to Google Cloud Storage",
	Example:       "gh --org kubernetes --repo release --bucket k8s-staging-release-test --release-dir release --tags v0.0.0,v0.0.1",
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkRequiredFlags(cmd.Flags())
	},
	RunE: func(*cobra.Command, []string) error {
		return run(opts)
	},
}

type options struct {
	downloadOnly       bool
	includePrereleases bool
	org                string
	repo               string
	bucket             string
	releaseDir         string
	outputDir          string
	tags               []string
	configFilePath     string
}

var opts = &options{}

var (
	orgFlag                = "org"
	repoFlag               = "repo"
	tagsFlag               = "tags"
	configFlag             = "config"
	includePrereleasesFlag = "include-prereleases"
	bucketFlag             = "bucket"
	releaseDirFlag         = "release-dir"
	outputDirFlag          = "output-dir"
	downloadOnlyFlag       = "download-only"

	// requiredFlags only if the config flag is not set
	requiredFlags = []string{
		orgFlag,
		repoFlag,
		bucketFlag,
	}
)

func init() {
	GHCmd.PersistentFlags().StringVar(
		&opts.org,
		orgFlag,
		"",
		"GitHub org/user",
	)

	GHCmd.PersistentFlags().StringVar(
		&opts.repo,
		repoFlag,
		"",
		"GitHub repo",
	)

	GHCmd.PersistentFlags().StringSliceVar(
		&opts.tags,
		tagsFlag,
		[]string{},
		"release tags to upload to GCS",
	)

	GHCmd.PersistentFlags().BoolVar(
		&opts.includePrereleases,
		includePrereleasesFlag,
		false,
		"specifies whether prerelease assets should be uploaded to GCS",
	)

	GHCmd.PersistentFlags().StringVar(
		&opts.bucket,
		bucketFlag,
		"",
		"GCS bucket to upload to",
	)

	GHCmd.PersistentFlags().StringVar(
		&opts.releaseDir,
		releaseDirFlag,
		"",
		fmt.Sprintf("directory to upload to within the specified GCS bucket, defaults to %q", filepath.Join(orgFlag, repoFlag)),
	)

	GHCmd.PersistentFlags().StringVar(
		&opts.outputDir,
		outputDirFlag,
		"",
		"local directory for releases to be downloaded to",
	)

	GHCmd.PersistentFlags().BoolVar(
		&opts.downloadOnly,
		downloadOnlyFlag,
		false,
		"only download the releases, do not push them to GCS. Requires the output-dir flag to also be set",
	)

	GHCmd.PersistentFlags().StringVar(
		&opts.configFilePath,
		configFlag,
		"",
		"config file to set all the branch/repositories the user wants to",
	)
}

func checkRequiredFlags(flags *pflag.FlagSet) error {
	if flags.Lookup(configFlag).Changed {
		return nil
	}

	checkRequiredFlags := []string{}
	flags.VisitAll(func(flag *pflag.Flag) {
		for _, requiredflag := range requiredFlags {
			if requiredflag == flag.Name && !flag.Changed {
				checkRequiredFlags = append(checkRequiredFlags, requiredflag)
			}
		}
	})

	if len(checkRequiredFlags) != 0 {
		return errors.New("Required flag(s) `" + strings.Join(checkRequiredFlags, ", ") + "` not set")
	}

	return nil
}

func run(opts *options) error {
	if err := opts.SetAndValidate(); err != nil {
		return fmt.Errorf("validating gh2gcs options: %w", err)
	}

	if err := gcli.PreCheck(); err != nil {
		return fmt.Errorf("pre-checking for GCP package usage: %w", err)
	}

	releaseCfgs := &gh2gcs.Config{}
	if opts.configFilePath != "" {
		logrus.Infof("Reading the config file %s...", opts.configFilePath)
		data, err := os.ReadFile(opts.configFilePath)
		if err != nil {
			return fmt.Errorf("failed to read the file: %w", err)
		}

		logrus.Info("Parsing the config...")
		err = yaml.UnmarshalStrict(data, &releaseCfgs)
		if err != nil {
			return fmt.Errorf("failed to decode the file: %w", err)
		}
	} else {
		// TODO: Expose certain GCSCopyOptions for user configuration
		newReleaseCfg := &gh2gcs.ReleaseConfig{
			Org:                opts.org,
			Repo:               opts.repo,
			Tags:               opts.tags,
			IncludePrereleases: opts.includePrereleases,
			GCSBucket:          opts.bucket,
			ReleaseDir:         opts.releaseDir,
		}

		releaseCfgs.ReleaseConfigs = append(releaseCfgs.ReleaseConfigs, *newReleaseCfg)
	}

	// Create a real GitHub API client
	gh := github.New()

	for _, releaseCfg := range releaseCfgs.ReleaseConfigs {
		if len(releaseCfg.Tags) == 0 {
			releaseTags, err := gh.GetReleaseTags(releaseCfg.Org, releaseCfg.Repo, releaseCfg.IncludePrereleases)
			if err != nil {
				return fmt.Errorf("getting release tags: %w", err)
			}

			releaseCfg.Tags = releaseTags
		}

		if err := gh2gcs.DownloadReleases(&releaseCfg, gh, opts.outputDir); err != nil {
			return fmt.Errorf("downloading release assets: %w", err)
		}
		logrus.Infof("Files downloaded to %s directory", opts.outputDir)

		if !opts.downloadOnly {
			if err := gh2gcs.Upload(&releaseCfg, gh, opts.outputDir); err != nil {
				return fmt.Errorf("uploading release assets to GCS: %w", err)
			}
		}
	}

	return nil
}

// SetAndValidate sets some default options and verifies if options are valid
func (o *options) SetAndValidate() error {
	logrus.Info("Validating gh2gcs options...")

	if o.releaseDir == "" {
		o.releaseDir = filepath.Join(o.org, o.repo)
		logrus.Infof("Using default release dir: %s", o.releaseDir)
	}

	if o.outputDir == "" {
		if opts.downloadOnly {
			return fmt.Errorf("when %s flag is set you need to specify the %s", downloadOnlyFlag, outputDirFlag)
		}

		tmpDir, err := os.MkdirTemp("", "gh")
		if err != nil {
			return fmt.Errorf("creating temp directory: %w", err)
		}

		o.outputDir = tmpDir
	}

	return nil
}
