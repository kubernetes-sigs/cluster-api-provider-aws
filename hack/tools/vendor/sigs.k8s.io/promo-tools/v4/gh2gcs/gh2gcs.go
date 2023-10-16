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

package gh2gcs

import (
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"sigs.k8s.io/release-sdk/github"
	"sigs.k8s.io/release-sdk/object"
)

// Config contains a slice of `ReleaseConfig` to be used when unmarshalling a
// yaml config containing multiple repository configs.
type Config struct {
	ReleaseConfigs []ReleaseConfig `yaml:"releaseConfigs"`
}

// ReleaseConfig contains source (GitHub) and destination (GCS) information
// to perform a copy/upload operation.
type ReleaseConfig struct {
	// GitHub options
	Org                string   `yaml:"org"`
	Repo               string   `yaml:"repo"`
	Tags               []string `yaml:"tags"`
	IncludePrereleases bool     `yaml:"includePrereleases"`

	// GCS options
	GCSBucket  string `yaml:"gcsBucket"`
	ReleaseDir string `yaml:"releaseDir"`
}

// DownloadReleases downloads release assets to a local directory
// Assets to download are derived from the tags specified in `ReleaseConfig`.
func DownloadReleases(releaseCfg *ReleaseConfig, ghClient *github.GitHub, outputDir string) error {
	tags := releaseCfg.Tags
	tagsString := strings.Join(tags, ", ")

	logrus.Infof(
		"Downloading assets for the following %s/%s release tags: %s",
		releaseCfg.Org,
		releaseCfg.Repo,
		tagsString,
	)

	return ghClient.DownloadReleaseAssets(releaseCfg.Org, releaseCfg.Repo, tags, outputDir)
}

// Upload copies a set of release assets from local directory to GCS
// Assets to upload are derived from the tags specified in `ReleaseConfig`.
func Upload(cfg *ReleaseConfig, _ *github.GitHub, outputDir string) error {
	uploadBase := filepath.Join(outputDir, cfg.Org, cfg.Repo)

	tags := cfg.Tags
	for _, tag := range tags {
		srcDir := filepath.Join(uploadBase, tag)
		gcsPath := filepath.Join(cfg.GCSBucket, cfg.ReleaseDir, tag)

		gcs := object.NewGCS()
		if err := gcs.CopyToRemote(srcDir, gcsPath); err != nil {
			return err
		}
	}

	return nil
}
