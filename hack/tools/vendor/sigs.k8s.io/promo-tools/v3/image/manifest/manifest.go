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

package manifest

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	reg "sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/schema"
	"sigs.k8s.io/promo-tools/v3/types/image"
)

const (
	// This is a banned tag. It is not allowed to be manipulated with this tool.
	latestTag = "latest"
)

// GrowOptions holds the  parameters for modifying manifests.
type GrowOptions struct {
	// BaseDir is the directory containing the thin promoter manifests.
	BaseDir string
	// StagingRepo is the staging subproject repo to read from. If no filters
	// are provided, all images are attempted to be promoted as-is without any
	// modifications.
	StagingRepo image.Registry
	// FilterImage is the image (name) to filter by. Optional.
	FilterImages []image.Name
	// FilterDigest is the image digest to filter by. Optional.
	FilterDigests []image.Digest
	// FilterTags is the image tag to filter by. Optional.
	FilterTags []image.Tag
}

// Populate sets the values for GrowOptions.
func (o *GrowOptions) Populate(
	baseDir,
	stagingRepo string,
	filterImages,
	filterDigests,
	filterTags []string,
) error {
	baseDirAbsPath, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf(
			"cannot resolve %q to absolute path: %w", baseDir, err)
	}

	o.BaseDir = baseDirAbsPath
	o.StagingRepo = image.Registry(stagingRepo)
	o.FilterImages = toImageNames(filterImages)
	o.FilterDigests = toDigests(filterDigests)
	o.FilterTags = toTags(filterTags)

	return nil
}

func toImageNames(imageStrings []string) []image.Name {
	imgNames := []image.Name{}
	for _, imgString := range imageStrings {
		imgName := image.Name(imgString)
		imgNames = append(imgNames, imgName)
	}

	return imgNames
}

func toTags(tagStrings []string) []image.Tag {
	tags := []image.Tag{}
	for _, tagString := range tagStrings {
		tag := image.Tag(tagString)
		tags = append(tags, tag)
	}

	return tags
}

func toDigests(digestStrings []string) []image.Digest {
	digests := []image.Digest{}
	for _, digestString := range digestStrings {
		digest := image.Digest(digestString)
		digests = append(digests, digest)
	}

	return digests
}

// Validate validates the options.
func (o *GrowOptions) Validate() error {
	if o.BaseDir == "" {
		return xerrors.New("must specify --base_dir")
	}

	if o.StagingRepo == "" {
		return xerrors.New("must specify --staging_repo")
	}

	if containsTag(o.FilterTags, latestTag) {
		return fmt.Errorf(
			"--filter_tag cannot be %q (anti-pattern)", latestTag)
	}

	return nil
}

func containsTag(tags []image.Tag, check string) bool {
	for _, tag := range tags {
		if tag == image.Tag(check) {
			return true
		}
	}

	return false
}

// Grow modifies a manifest by adding images into it.
func Grow(
	_ context.Context,
	o *GrowOptions,
) error {
	var err error
	var riiCombined registry.RegInvImage

	// (1) Scan the BaseDir and find the promoter manifest to modify.
	mfest, err := Find(o)
	if err != nil {
		return err
	}

	// (2) Scan the StagingRepo, and whittle the read results down with some
	// filters (Filter* fields in GrowOptions).
	riiUnfiltered, err := ReadStagingRepo(o)
	if err != nil {
		return err
	}

	// (3) Apply some filters.
	riiFiltered, err := ApplyFilters(o, riiUnfiltered)
	if err != nil {
		return err
	}

	// (4) Inject (2)'s output into (1)'s manifest's images to create a larger
	// RegInvImage.
	riiCombined = Union(mfest.ToRegInvImage(), riiFiltered)

	// (5) Write back RegInvImage as Manifest ([]Image field}) back onto disk.
	err = Write(mfest, riiCombined)

	return err
}

// Write writes images as YAML out to the expected path of the given
// (thin) schema.
func Write(m schema.Manifest, rii registry.RegInvImage) error {
	// Chop off trailing "promoter-manifest.yaml".
	p := path.Dir(m.Filepath)
	// Get staging repo directory name as it is laid out in the thin manifest
	// dir.
	stagingRepoName := path.Base(p)
	// Construct path to the images.yaml.
	imagesPath := path.Join(p, "..", "..",
		"images", stagingRepoName, "images.yaml")
	logrus.Infoln("RENDER", imagesPath)

	// Write the file.
	err := os.WriteFile(
		imagesPath, []byte(rii.ToYAML(registry.YamlMarshalingOpts{})), 0o644)
	return err
}

// Find finds the manifest to modify.
func Find(o *GrowOptions) (schema.Manifest, error) {
	var err error
	var manifests []schema.Manifest
	manifests, err = schema.ParseThinManifestsFromDir(o.BaseDir, false)
	if err != nil {
		return schema.Manifest{}, err
	}

	logrus.Infof("%d manifests parsed", len(manifests))
	for _, manifest := range manifests {
		if manifest.SrcRegistry.Name == o.StagingRepo {
			return manifest, nil
		}
	}
	return schema.Manifest{},
		fmt.Errorf("could not find Manifest for %q", o.StagingRepo)
}

// ReadStagingRepo reads the StagingRepo, and applies whatever filters are
// available to the resulting RegInvImage. This RegInvImage is what we want to
// inject into the "images.yaml" of a thin schema.
func ReadStagingRepo(
	o *GrowOptions,
) (registry.RegInvImage, error) {
	stagingRepoRC := registry.Context{
		Name: o.StagingRepo,
	}

	manifests := []schema.Manifest{
		{
			Registries: []registry.Context{
				stagingRepoRC,
			},
			Images: []registry.Image{},
		},
	}

	sc, err := reg.MakeSyncContext(
		manifests,
		10,
		true,
		false)
	if err != nil {
		return registry.RegInvImage{}, err
	}
	sc.ReadRegistries(
		[]registry.Context{stagingRepoRC},
		// Read all registries recursively, because we want to produce a
		// complete snapshot.
		true,
		reg.MkReadRepositoryCmdReal)

	return sc.Inv[manifests[0].Registries[0].Name], nil
}

// ApplyFilters applies the filters in the options to whittle down the given
// rii.
func ApplyFilters(o *GrowOptions, rii registry.RegInvImage) (registry.RegInvImage, error) {
	// If nothing to filter, short-circuit.
	if len(rii) == 0 {
		return rii, nil
	}

	// Now perform some filtering, if any.
	if len(o.FilterImages) > 0 {
		rii = FilterByImages(rii, o.FilterImages)
	}

	if len(o.FilterTags) > 0 {
		// TODO(manifest): Should func be pulled into this package?
		rii = FilterByTags(rii, o.FilterTags)
	}

	if len(o.FilterDigests) > 0 {
		rii = FilterByDigests(rii, o.FilterDigests)
	}

	// Remove any other tags that should still be filtered.
	excludeTags := map[image.Tag]bool{latestTag: true}
	rii = ExcludeTags(rii, excludeTags)

	if len(rii) == 0 {
		return registry.RegInvImage{}, xerrors.New(
			"no images survived filtering; double-check your --filter_* flag(s) for typos",
		)
	}

	return rii, nil
}

// FilterByImages removes all images in RegInvImage that do not match the
// filterImage.
func FilterByImages(rii registry.RegInvImage, filterImages []image.Name) registry.RegInvImage {
	filtered := make(registry.RegInvImage)
	for imageName, digestTags := range rii {
		for _, filterImage := range filterImages {
			if imageName == filterImage {
				filtered[imageName] = digestTags
			}
		}
	}
	return filtered
}

// FilterByTags removes all images in RegInvImage that do not match the
// filterTag.
// TODO(manifest): Dedupe with `FilterByTag` in legacy/dockerregistry/inventory.go
func FilterByTags(rii registry.RegInvImage, filterTags []image.Tag) registry.RegInvImage {
	filtered := make(registry.RegInvImage)

	for imageName, digestTags := range rii {
		for digest, tags := range digestTags {
			for _, tag := range tags {
				for _, filterTag := range filterTags {
					if tag == filterTag {
						if filtered[imageName] == nil {
							filtered[imageName] = make(registry.DigestTags)
						}

						filtered[imageName][digest] = append(
							filtered[imageName][digest],
							tag,
						)
					}
				}
			}
		}
	}

	return filtered
}

// FilterByDigests removes all images in RegInvImage that do not match the
// filterDigest.
func FilterByDigests(rii registry.RegInvImage, filterDigests []image.Digest) registry.RegInvImage {
	filtered := make(registry.RegInvImage)
	for imageName, digestTags := range rii {
		for digest, tags := range digestTags {
			for _, filterDigest := range filterDigests {
				if digest == filterDigest {
					if filtered[imageName] == nil {
						filtered[imageName] = make(registry.DigestTags)
					}
					filtered[imageName][digest] = tags
				}
			}
		}
	}

	return filtered
}

// ExcludeTags removes tags in rii that match excludedTags.
func ExcludeTags(rii registry.RegInvImage, excludedTags map[image.Tag]bool) registry.RegInvImage {
	filtered := make(registry.RegInvImage)
	for imageName, digestTags := range rii {
		for digest, tags := range digestTags {
			for _, tag := range tags {
				if _, excludeMe := excludedTags[tag]; excludeMe {
					continue
				}
				if filtered[imageName] == nil {
					filtered[imageName] = make(registry.DigestTags)
				}
				filtered[imageName][digest] = append(
					filtered[imageName][digest],
					tag)
			}
		}
	}
	return filtered
}

// Union inject b's contents into a. However, it does so in a special way.
func Union(a, b registry.RegInvImage) registry.RegInvImage {
	for imageName, digestTags := range b {
		// If a does not have this image at all, then it's a simple
		// injection.
		if a[imageName] == nil {
			a[imageName] = digestTags
			continue
		}
		for digest, tags := range digestTags {
			// If a has the image but not this digest, inject just this digest
			// and all associated tags.
			if a[imageName][digest] == nil {
				a[imageName][digest] = tags
				continue
			}
			// If c has the digest already, try to inject those tags in b that
			// are not already in a.
			tagSlice := registry.TagSlice{}
			for tag := range tags.Union(a[imageName][digest]) {
				if tag == "latest" {
					continue
				}
				tagSlice = append(tagSlice, tag)
			}
			a[imageName][digest] = tagSlice
		}
	}

	return a
}
