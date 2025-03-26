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

	reg "sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry"
	"sigs.k8s.io/promo-tools/v4/internal/legacy/dockerregistry/schema"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
)

// This file has all the promoter implementation functions
// related to image promotion.

// ParseManifests reads the manifest file or manifest directory
// and parses them to return a slice of Manifest objects.
func (di *DefaultPromoterImplementation) ParseManifests(opts *options.Options) (mfests []schema.Manifest, err error) {
	// If the options have a manifest file defined, we use that one
	if opts.Manifest != "" {
		mfest, err := schema.ParseManifestFromFile(opts.Manifest)
		if err != nil {
			return mfests, fmt.Errorf("parsing the manifest file: %w", err)
		}

		mfests = []schema.Manifest{mfest}
		// The thin manifests
	} else if opts.ThinManifestDir != "" {
		mfests, err = schema.ParseThinManifestsFromDir(opts.ThinManifestDir, opts.UseProwManifestDiff)
		if err != nil {
			return nil, fmt.Errorf("parsing thin manifest directory: %w", err)
		}
	}
	return mfests, nil
}

// MakeSyncContext takes a slice of manifests and creates a sync context
// object based on them and the promoter options
func (di DefaultPromoterImplementation) MakeSyncContext(
	opts *options.Options, mfests []schema.Manifest,
) (*reg.SyncContext, error) {
	sc, err := reg.MakeSyncContext(
		mfests, opts.Threads, opts.Confirm, opts.UseServiceAcct,
	)
	if err != nil {
		return nil, fmt.Errorf("creating sync context: %w", err)
	}
	return sc, err
}

// GetPromotionEdges checks the manifests and determines from
// them the promotion edges, ie the images that need to be
// promoted.
func (di *DefaultPromoterImplementation) GetPromotionEdges(
	sc *reg.SyncContext, mfests []schema.Manifest,
) (promotionEdges map[reg.PromotionEdge]interface{}, err error) {
	// First, get the "edges" from the manifests
	promotionEdges, err = reg.ToPromotionEdges(mfests)
	if err != nil {
		return nil, fmt.Errorf("converting list of manifests to edges for promotion: %w", err)
	}

	// Run the promotion edge filtering
	promotionEdges, ok, err := sc.FilterPromotionEdges(promotionEdges, true)
	if err != nil {
		return nil, fmt.Errorf("filtering promotion edges: %w", err)
	}
	if !ok {
		// If any funny business was detected during a comparison of the manifests
		// with the state of the registries, then exit immediately.
		return nil, errors.New("encountered errors during edge filtering")
	}
	return promotionEdges, nil
}

// PromoteImages starts an image promotion of a set of edges
func (di *DefaultPromoterImplementation) PromoteImages(
	sc *reg.SyncContext,
	promotionEdges map[reg.PromotionEdge]interface{},
) error {
	if err := sc.Promote(promotionEdges, nil); err != nil {
		return fmt.Errorf("running image promotion: %w", err)
	}
	di.PrintSection("END (PROMOTION)", true)
	return nil
}
