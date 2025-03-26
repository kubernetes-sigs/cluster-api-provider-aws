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

package mm

import (
	"context"

	"github.com/spf13/cobra"

	"sigs.k8s.io/promo-tools/v4/image/manifest"
)

// TODO(cip-mm): Remove in the next minor release.

var MMCmd = &cobra.Command{
	Short: "[DEPRECATED] mm → Manifest Modifier",
	Long: `[DEPRECATED] mm → Manifest Modifier

This tool **m**odifies promoter **m**anifests. For now it dumps some filtered
subset of a staging GCR and merges those contents back into a given promoter
manifest.`,
	Use:           "mm",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

type modifyOptions struct {
	baseDir      string
	stagingRepo  string
	filterImage  string
	filterDigest string
	filterTag    string
}

var modifyOpts = &modifyOptions{}

func init() {
	MMCmd.PersistentFlags().StringVar(
		&modifyOpts.baseDir,
		"base_dir",
		"",
		"the manifest directory to look at and modify",
	)
	MMCmd.PersistentFlags().StringVar(
		&modifyOpts.stagingRepo,
		"staging_repo",
		"",
		"the staging repo which we want to read from",
	)
	MMCmd.PersistentFlags().StringVar(
		&modifyOpts.filterImage,
		"filter_image",
		"",
		"filter staging repo by this image name",
	)
	MMCmd.PersistentFlags().StringVar(
		&modifyOpts.filterDigest,
		"filter_digest",
		"",
		"filter images by this digest",
	)
	MMCmd.PersistentFlags().StringVar(
		&modifyOpts.filterTag,
		"filter_tag",
		"",
		"filter images by this tag",
	)
}

func run() error {
	opt := manifest.GrowOptions{}
	if err := opt.Populate(
		modifyOpts.baseDir,
		modifyOpts.stagingRepo,
		[]string{modifyOpts.filterImage},
		[]string{modifyOpts.filterDigest},
		[]string{modifyOpts.filterTag},
	); err != nil {
		return err
	}

	if err := opt.Validate(); err != nil {
		return err
	}

	ctx := context.Background()
	return manifest.Grow(ctx, &opt)
}
