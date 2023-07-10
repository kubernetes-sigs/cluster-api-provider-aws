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
	"io"
	"net/url"
	"sort"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"

	api "sigs.k8s.io/promo-tools/v3/api/files"
)

// FilestorePromoter manages the promotion of files.
type FilestorePromoter struct {
	Source *api.Filestore
	Dest   *api.Filestore

	Files []api.File

	// Confirm, if set, will trigger a PRODUCTION artifact promotion.
	Confirm bool

	// UseServiceAccount must be true, for service accounts to be used
	// This gives some protection against a hostile manifest.
	UseServiceAccount bool
}

//counterfeiter:generate . syncFilestore
type syncFilestore interface {
	// OpenReader opens an io.ReadCloser for the specified file
	OpenReader(ctx context.Context, name string) (io.ReadCloser, error)

	// UploadFile uploads a local file to the specified destination
	UploadFile(ctx context.Context, dest string, localFile string) error

	// ListFiles returns all the file artifacts in the filestore, recursively.
	ListFiles(ctx context.Context) (map[string]*SyncFileInfo, error)
}

var supportedProviders = []Provider{
	GoogleCloudStorage,
	S3Storage,
}

func openFilestore(
	ctx context.Context,
	filestore *api.Filestore,
	useServiceAccount, confirm bool,
) (syncFilestore, error) {
	for _, provider := range supportedProviders {
		scheme := provider.Scheme()
		if strings.HasPrefix(filestore.Base, scheme+"://") {
			return provider.OpenFilestore(ctx, filestore, useServiceAccount, confirm)
		}
	}

	expected := []string{}
	for _, provider := range supportedProviders {
		expected = append(expected, provider.Scheme())
	}
	sort.Strings(expected)
	return nil, fmt.Errorf(
		"unrecognized scheme %q (supported schemes: %s)",
		filestore.Base,
		strings.Join(expected, ", "),
	)
}

// openGCSFilestore opens a filestore backed by Google Cloud Storage (GCS)
func (p *gcsProvider) OpenFilestore(
	ctx context.Context,
	filestore *api.Filestore,
	useServiceAccount, confirm bool,
) (syncFilestore, error) {
	u, err := url.Parse(filestore.Base)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing filestore base %q: %v",
			filestore.Base,
			err,
		)
	}

	if u.Scheme != p.Scheme() {
		return nil, fmt.Errorf("unrecognized scheme %q, expected %s", filestore.Base, p.Scheme())
	}

	withAuth, err := useStorageClientAuth(filestore, useServiceAccount, confirm)
	if err != nil {
		return nil, err
	}

	var opts []option.ClientOption
	if withAuth {
		logrus.Infof(
			"requesting an authenticated storage client for %s",
			filestore.Base,
		)

		ts := &gcloudTokenSource{ServiceAccount: filestore.ServiceAccount}
		opts = append(opts, option.WithTokenSource(ts))
	} else {
		logrus.Warnf(
			"requesting an UNAUTHENTICATED storage client for %s",
			filestore.Base,
		)

		opts = append(opts, option.WithoutAuthentication())
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("error building GCS client: %v", err)
	}

	prefix := strings.TrimPrefix(u.Path, "/")
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	bucket := u.Host

	s := &gcsSyncFilestore{
		provider:  p,
		filestore: filestore,
		client:    client,
		bucket:    bucket,
		prefix:    prefix,
	}
	return s, nil
}

func useStorageClientAuth(
	filestore *api.Filestore,
	useServiceAccount, confirm bool,
) (bool, error) {
	withAuth := false

	// Source filestores should be world-readable, so authentication should
	// not be required.
	if filestore.Src {
		return withAuth, nil
	}

	if confirm {
		if filestore.ServiceAccount == "" {
			return withAuth, fmt.Errorf("cannot execute a production file promotion without a service account")
		}

		withAuth = true
	} else if useServiceAccount {
		if filestore.ServiceAccount == "" {
			return withAuth, fmt.Errorf("requested an authenticated file promotion, but a service account was not specified")
		}

		withAuth = true
	}

	return withAuth, nil
}

// computeNeededOperations determines the list of files that need to be copied
func (p *FilestorePromoter) computeNeededOperations(
	source, dest map[string]*SyncFileInfo,
	destFilestore syncFilestore,
) ([]SyncFileOp, error) {
	ops := make([]SyncFileOp, 0)

	for i := range p.Files {
		f := &p.Files[i]
		relativePath := f.Name
		sourceFile := source[relativePath]
		if sourceFile == nil {
			absolutePath := joinFilepath(p.Source, relativePath)
			logrus.Debugf(
				"file %q not found in source (%q)",
				relativePath,
				absolutePath,
			)

			continue
		}

		destFile := dest[relativePath]
		if destFile == nil {
			destFile = &SyncFileInfo{}
			destFile.RelativePath = sourceFile.RelativePath
			destFile.AbsolutePath = joinFilepath(
				p.Dest,
				sourceFile.RelativePath)
			destFile.filestore = destFilestore
			ops = append(ops, &copyFileOp{
				Source:       sourceFile,
				Dest:         destFile,
				ManifestFile: f,
			})
			continue
		}

		changed := false
		if destFile.MD5 != sourceFile.MD5 {
			logrus.Warnf("MD5 mismatch on source %q vs dest %q: %q vs %q",
				sourceFile.AbsolutePath,
				destFile.AbsolutePath,
				sourceFile.MD5,
				destFile.MD5)
			changed = true
		}

		if destFile.Size != sourceFile.Size {
			logrus.Warnf("Size mismatch on source %q vs dest %q: %d vs %d",
				sourceFile.AbsolutePath,
				destFile.AbsolutePath,
				sourceFile.Size,
				destFile.Size)
			changed = true
		}

		if !changed {
			logrus.Infof("metadata match for %q", destFile.AbsolutePath)
			continue
		}
		ops = append(ops, &copyFileOp{
			Source:       sourceFile,
			Dest:         destFile,
			ManifestFile: f,
		})
	}

	return ops, nil
}

func joinFilepath(filestore *api.Filestore, relativePath string) string {
	s := strings.TrimSuffix(filestore.Base, "/")
	s += "/"
	s += strings.TrimPrefix(relativePath, "/")
	return s
}

// BuildOperations builds the required operations to sync from the
// Source Filestore to the Dest Filestore.
func (p *FilestorePromoter) BuildOperations(
	ctx context.Context,
) ([]SyncFileOp, error) {
	sourceFilestore, err := openFilestore(
		ctx,
		p.Source,
		p.UseServiceAccount,
		p.Confirm,
	)
	if err != nil {
		return nil, err
	}
	if sourceFilestore == nil {
		return nil, fmt.Errorf("source filestore cannot be nil")
	}

	destFilestore, err := openFilestore(
		ctx,
		p.Dest,
		p.UseServiceAccount,
		p.Confirm,
	)
	if err != nil {
		return nil, err
	}
	if destFilestore == nil {
		return nil, fmt.Errorf("destination filestore cannot be nil")
	}

	sourceFiles, err := sourceFilestore.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	destFiles, err := destFilestore.ListFiles(ctx)
	if err != nil {
		return nil, err
	}

	return p.computeNeededOperations(sourceFiles, destFiles, destFilestore)
}
