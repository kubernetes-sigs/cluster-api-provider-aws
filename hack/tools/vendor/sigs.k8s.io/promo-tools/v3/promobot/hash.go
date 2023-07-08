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

package promobot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"

	api "sigs.k8s.io/promo-tools/v3/api/files"
	"sigs.k8s.io/release-utils/hash"
)

// GenerateManifestOptions holds the parameters for a hash-files operation.
type GenerateManifestOptions struct {
	// BaseDir is the directory containing the files to hash
	BaseDir string

	// Prefix exports only files matching the specified prefix.
	//
	// If we were instead to change BaseDir, we would also
	// restrict the files, but the relative paths would also
	// change.
	Prefix string
}

// PopulateDefaults sets the default values for GenerateManifestOptions.
func (o *GenerateManifestOptions) PopulateDefaults() {
	// There are no fields with non-empty default values
	// (but we still want to follow the PopulateDefaults pattern)
}

// GenerateManifest generates a manifest containing the files in options.BaseDir
func GenerateManifest(_ context.Context, options GenerateManifestOptions) (*api.Manifest, error) {
	manifest := &api.Manifest{}

	if options.BaseDir == "" {
		return nil, xerrors.New("must specify BaseDir")
	}

	basedir := options.BaseDir
	if !strings.HasSuffix(basedir, "/") {
		basedir += "/"
	}

	if err := filepath.Walk(basedir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasPrefix(p, basedir) {
			return fmt.Errorf("expected path %q to have prefix %q", p, basedir)
		}

		if !strings.HasPrefix(p, filepath.Join(basedir, options.Prefix)) {
			return nil
		}

		if !info.IsDir() {
			relativePath := strings.TrimPrefix(p, basedir)
			sha256, err := hash.SHA256ForFile(p)
			if err != nil {
				return fmt.Errorf("error hashing file %q: %w", p, err)
			}
			manifest.Files = append(manifest.Files, api.File{
				Name:   relativePath,
				SHA256: sha256,
			})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error walking path %q: %w", options.BaseDir, err)
	}

	return manifest, nil
}
