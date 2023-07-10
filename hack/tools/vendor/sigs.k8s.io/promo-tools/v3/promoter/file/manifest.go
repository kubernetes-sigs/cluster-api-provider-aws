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

package file

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "sigs.k8s.io/promo-tools/v3/api/files"
)

// ManifestPromoter promotes files as described in Manifest.
type ManifestPromoter struct {
	Manifest *api.Manifest

	// Confirm, if set, will trigger a PRODUCTION artifact promotion.
	Confirm bool

	// UseServiceAccount must be true, for service accounts to be used
	// This gives some protection against a hostile manifest.
	UseServiceAccount bool
}

// BuildOperations builds the required operations to sync from the
// Source Filestore to all Dest Filestores in the manifest.
func (p *ManifestPromoter) BuildOperations(
	ctx context.Context,
) ([]SyncFileOp, error) {
	source, err := getSourceFilestore(p.Manifest)
	if err != nil {
		return nil, err
	}

	var operations []SyncFileOp

	for i := range p.Manifest.Filestores {
		filestore := &p.Manifest.Filestores[i]
		if filestore.Src {
			continue
		}
		logrus.Infof("processing destination %q", filestore.Base)
		fp := &FilestorePromoter{
			Source:            source,
			Dest:              filestore,
			Files:             p.Manifest.Files,
			Confirm:           p.Confirm,
			UseServiceAccount: p.UseServiceAccount,
		}
		ops, err := fp.BuildOperations(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"error building promotion operations for %q: %v",
				filestore.Base, err)
		}
		operations = append(operations, ops...)
	}

	return operations, nil
}

// getSourceFilestore returns the Filestore with the source attribute
// It returns an error if the source filestore cannot be found, or if
// multiple filestores are marked as the source.
func getSourceFilestore(manifest *api.Manifest) (*api.Filestore, error) {
	var source *api.Filestore
	for i := range manifest.Filestores {
		filestore := &manifest.Filestores[i]
		if filestore.Src {
			if source != nil {
				return nil, fmt.Errorf("found multiple source filestores")
			}
			source = filestore
		}
	}
	if source == nil {
		return nil, fmt.Errorf("source filestore not found")
	}
	return source, nil
}
