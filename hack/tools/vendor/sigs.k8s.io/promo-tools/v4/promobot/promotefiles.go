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

package promobot

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/sirupsen/logrus"

	api "sigs.k8s.io/promo-tools/v4/api/files"
	"sigs.k8s.io/promo-tools/v4/promoter/file"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// PromoteFilesOptions holds the flag-values for a file promotion
type PromoteFilesOptions struct {
	// FilestoresPath is the path to the manifest file containing the filestores section
	FilestoresPath string

	// FilesPath specifies a path to manifest files containing the files section.
	FilesPath string

	// ManifestsPath specifies a path containing filestores and files for
	// multiple projects.
	//
	// Example layout:
	//
	// ├── filestores
	// │   ├── project1
	// │   │   └── filepromoter-manifest.yaml
	// │   └── project2
	// │       └── filepromoter-manifest.yaml
	// └── manifests
	//     ├── project1
	//     │   ├── blue.yaml
	//     │   ├── green.yaml
	//     │   └── red.yaml
	//     └── project2
	//         ├── blue.yaml
	//         ├── green.yaml
	//         └── red.yaml
	ManifestsPath string

	// Confirm, if set, will trigger a PRODUCTION artifact promotion.
	Confirm bool

	// UseServiceAccount must be true, for service accounts to be used
	// This gives some protection against a hostile manifest.
	UseServiceAccount bool

	// Out is the destination for "normal" output (such as dry-run)
	Out io.Writer
}

// PopulateDefaults sets the default values for PromoteFilesOptions
func (o *PromoteFilesOptions) PopulateDefaults() {
	o.Confirm = false
	o.UseServiceAccount = false
	o.Out = os.Stdout
}

// RunPromoteFiles executes a file promotion command
func RunPromoteFiles(ctx context.Context, options PromoteFilesOptions) error {
	manifests, err := ReadManifests(options)
	if err != nil {
		return err
	}

	if options.Confirm {
		fmt.Fprintf(
			options.Out,
			"********** START **********\n",
		)
	} else {
		fmt.Fprintf(
			options.Out,
			"********** START (DRY RUN) **********\n",
		)
	}

	var ops []file.SyncFileOp
	for _, manifest := range manifests {
		promoter := &file.ManifestPromoter{
			Manifest:          manifest,
			Confirm:           options.Confirm,
			UseServiceAccount: options.UseServiceAccount,
		}

		o, err := promoter.BuildOperations(ctx)
		if err != nil {
			return fmt.Errorf("error building operations: %v", err)
		}

		ops = append(ops, o...)
	}

	// So that we can support future parallel execution, an error
	// in one operation does not prevent us attempting the
	// remaining operations
	var errors []error
	for _, op := range ops {
		if _, err := fmt.Fprintf(options.Out, "%v\n", op); err != nil {
			errors = append(
				errors,
				fmt.Errorf("error writing to output: %v", err),
			)
		}

		if options.Confirm {
			if err := op.Run(ctx); err != nil {
				logrus.Warnf("error copying file: %v", err)
				errors = append(errors, err)
			}
		}
	}

	if len(errors) != 0 {
		fmt.Fprintf(
			options.Out,
			"********** FINISHED WITH ERRORS **********\n")
		for _, err := range errors {
			fmt.Fprintf(options.Out, "%v\n", err)
		}

		return errors[0]
	}

	if options.Confirm {
		fmt.Fprintf(
			options.Out,
			"********** FINISHED **********\n",
		)
	} else {
		fmt.Fprintf(
			options.Out,
			"********** FINISHED (DRY RUN) **********\n",
		)
	}

	return nil
}

// ReadManifests reads a set of manifests.
func ReadManifests(options PromoteFilesOptions) ([]*api.Manifest, error) {
	manifests := make([]*api.Manifest, 0)
	mPath := options.ManifestsPath
	if mPath == "" {
		m, err := ReadManifest(options)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, m)
		return manifests, nil
	}

	filestoresDir := filepath.Join(mPath, "filestores")
	manifestsDir := filepath.Join(mPath, "manifests")

	var projects []string

	// TODO: Consider using filepath.WalkDir() instead
	if err := filepath.Walk(
		filestoresDir,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() || info.Name() == "filestores" {
				return nil
			}

			projects = append(projects, info.Name())
			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("error listing projects: %w", err)
	}

	for _, prj := range projects {
		filestores := filepath.Join(
			filestoresDir,
			prj,
			"filepromoter-manifest.yaml",
		)

		files := filepath.Join(
			manifestsDir,
			prj,
		)

		prjOpts := &PromoteFilesOptions{
			FilestoresPath:    filestores,
			FilesPath:         files,
			Confirm:           options.Confirm,
			UseServiceAccount: options.UseServiceAccount,
			Out:               options.Out,
		}

		m, err := ReadManifest(*prjOpts)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, m)
	}

	return manifests, nil
}

// ReadManifest reads a manifest.
func ReadManifest(options PromoteFilesOptions) (*api.Manifest, error) {
	merged := &api.Manifest{}

	filestores, err := readFilestores(options.FilestoresPath)
	if err != nil {
		return nil, err
	}
	merged.Filestores = filestores

	files, err := readFiles(options.FilesPath)
	if err != nil {
		return nil, err
	}
	merged.Files = files

	// Validate the merged manifest
	if err := merged.Validate(); err != nil {
		return nil, fmt.Errorf("error validating merged manifest: %v", err)
	}

	return merged, nil
}

// readFilestores reads a filestores manifest
func readFilestores(p string) ([]api.Filestore, error) {
	if p == "" {
		return nil, fmt.Errorf("FilestoresPath is required")
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest %q: %v", p, err)
	}

	manifest, err := api.ParseManifest(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing manifest %q: %v", p, err)
	}

	if len(manifest.Files) != 0 {
		return nil, fmt.Errorf(
			"files should not be present in filestore manifest %q",
			p)
	}

	return manifest.Filestores, nil
}

// readFiles reads and merges the file manifests from the file or directory filesPath
func readFiles(filesPath string) ([]api.File, error) {
	// We first list and sort the paths, for a consistent ordering
	var paths []string
	err := filepath.Walk(filesPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		paths = append(paths, p)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error listing file manifests: %w", err)
	}

	sort.Strings(paths)

	var files []api.File
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("error reading file %q: %w", p, err)
		}

		manifest, err := api.ParseManifest(b)
		if err != nil {
			return nil, fmt.Errorf("error parsing manifest %q: %v", p, err)
		}

		if len(manifest.Filestores) != 0 {
			return nil, fmt.Errorf("filestores should not be present in manifest %q", p)
		}

		files = append(files, manifest.Files...)
	}

	return files, nil
}
