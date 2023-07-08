/*
Copyright 2022 The Kubernetes Authors.

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

package imagepromoter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	reg "sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/registry"
	"sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry/schema"
	options "sigs.k8s.io/promo-tools/v3/promoter/image/options"
	"sigs.k8s.io/promo-tools/v3/types/image"
)

// Run a snapshot
func (di *DefaultPromoterImplementation) Snapshot(opts *options.Options, rii registry.RegInvImage) error {
	// Run the snapshot
	var snapshot string
	switch strings.ToLower(opts.OutputFormat) {
	case "csv":
		snapshot = rii.ToCSV()
	case "yaml":
		snapshot = rii.ToYAML(registry.YamlMarshalingOpts{})
	default:
		// In the previous cli/run it took any malformed format string. Now we err.
		return fmt.Errorf("invalid snapshot output format: %s", opts.OutputFormat)
	}

	// TODO: Maybe store the snapshot somewhere?
	di.PrintSection("END (SNAPSHOT)", opts.Confirm)
	fmt.Println(snapshot)
	return nil
}

func (di *DefaultPromoterImplementation) GetSnapshotSourceRegistry(
	opts *options.Options,
) (*registry.Context, error) {
	// Build the source registry:
	srcRegistry := &registry.Context{
		ServiceAccount: opts.SnapshotSvcAcct,
		Src:            true,
	}

	// The only difference when running from Snapshot or
	// ManifestBasedSnapshotOf will be the Name property
	// of the source registry
	if opts.Snapshot != "" {
		srcRegistry.Name = image.Registry(opts.Snapshot)
	} else if opts.ManifestBasedSnapshotOf == "" {
		srcRegistry.Name = image.Registry(opts.ManifestBasedSnapshotOf)
	} else {
		return nil, errors.New(
			"when snapshotting, Snapshot or ManifestBasedSnapshotOf have to be set",
		)
	}

	return srcRegistry, nil
}

// GetSnapshotManifest creates the manifest list from the
// specified snapshot source
func (di *DefaultPromoterImplementation) GetSnapshotManifests(
	opts *options.Options,
) ([]schema.Manifest, error) {
	// Build the source registry:
	srcRegistry, err := di.GetSnapshotSourceRegistry(opts)
	if err != nil {
		return nil, fmt.Errorf("building source registry for snapshot")
	}

	// Add it to a new manifest and return it:
	return []schema.Manifest{
		{
			Registries: []registry.Context{
				*srcRegistry,
			},
			Images: []registry.Image{},
		},
	}, nil
}

// AppendManifestToSnapshot checks if a manifest was specified in the
// options passed to the promoter. If one is found, we parse it and
// append it to the list of manifests generated for the snapshot
// during GetSnapshotManifests()
func (di *DefaultPromoterImplementation) AppendManifestToSnapshot(
	opts *options.Options, mfests []schema.Manifest,
) ([]schema.Manifest, error) {
	// If no manifest was passed in the options, we return the
	// same list of manifests unchanged
	if opts.Manifest == "" {
		logrus.Info("No manifest defined, not appending to snapshot")
		return mfests, nil
	}

	// Parse the specified manifest and append it to the list
	mfest, err := schema.ParseManifestFromFile(opts.Manifest)
	if err != nil {
		return nil, fmt.Errorf("parsing specified manifest: %w", err)
	}

	return append(mfests, mfest), nil
}

func (di *DefaultPromoterImplementation) GetRegistryImageInventory(
	opts *options.Options, mfests []schema.Manifest,
) (registry.RegInvImage, error) {
	// I'm pretty sure the registry context here can be the same for
	// both snapshot sources and when running in the original cli/run,
	// In the 2nd case (Snapshot), it was recreated like we do here.
	sc, err := di.MakeSyncContext(opts, mfests)
	if err != nil {
		return nil, fmt.Errorf("making sync context for registry inventory: %w", err)
	}

	srcRegistry, err := di.GetSnapshotSourceRegistry(opts)
	if err != nil {
		return nil, fmt.Errorf("creting source registry for image inventory: %w", err)
	}

	if len(opts.ManifestBasedSnapshotOf) > 0 {
		promotionEdges, err := reg.ToPromotionEdges(mfests)
		if err != nil {
			return nil, fmt.Errorf("converting list of manifests to edges for promotion: %w", err)
		}

		// Create the registry inventory
		rii := reg.EdgesToRegInvImage(
			promotionEdges,
			opts.ManifestBasedSnapshotOf,
		)

		if opts.MinimalSnapshot {
			if err := sc.ReadRegistriesGGCR(
				[]registry.Context{*srcRegistry},
				true,
			); err != nil {
				return nil, fmt.Errorf("reading registry for minimal snapshot: %w", err)
			}

			sc.ReadGCRManifestLists(reg.MkReadManifestListCmdReal)
			rii = sc.RemoveChildDigestEntries(rii)
		}

		return rii, nil
	}

	if err := sc.ReadRegistriesGGCR(
		[]registry.Context{*srcRegistry}, true,
	); err != nil {
		return nil, fmt.Errorf("reading registries: %w", err)
	}

	rii := sc.Inv[mfests[0].Registries[0].Name]
	if opts.SnapshotTag != "" {
		rii = reg.FilterByTag(rii, opts.SnapshotTag)
	}

	if opts.MinimalSnapshot {
		logrus.Info("removing tagless child digests of manifest lists")
		sc.ReadGCRManifestLists(reg.MkReadManifestListCmdReal)
		rii = sc.RemoveChildDigestEntries(rii)
	}
	return rii, nil
}
