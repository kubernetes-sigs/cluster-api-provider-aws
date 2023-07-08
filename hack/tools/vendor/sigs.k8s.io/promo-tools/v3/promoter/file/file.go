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
	"os"

	"github.com/sirupsen/logrus"

	api "sigs.k8s.io/promo-tools/v3/api/files"
	"sigs.k8s.io/release-utils/hash"
)

// SyncFileInfo tracks a file during the synchronization operation.
type SyncFileInfo struct {
	RelativePath string
	AbsolutePath string

	// Some backends (GCS and S3) expose the MD5 of the content in metadata
	// This can allow skipping unnecessary copies.
	// Note: with multipart uploads or compression, the value is unobvious.
	MD5 string

	Size int64

	filestore syncFilestore
}

// copyFileOp manages copying a single file.
type copyFileOp struct {
	Source *SyncFileInfo
	Dest   *SyncFileInfo

	ManifestFile *api.File
}

// Run implements SyncFileOp.Run
func (o *copyFileOp) Run(ctx context.Context) error {
	// Download to our temp file
	f, err := os.CreateTemp("", "promoter")
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	tempFilename := f.Name()

	defer func() {
		if f != nil {
			if err := f.Close(); err != nil {
				logrus.Warnf(
					"error closing temp file %q: %v",
					tempFilename, err)
			}
		}

		if err := os.Remove(tempFilename); err != nil {
			logrus.Warnf(
				"unable to remove temp file %q: %v",
				tempFilename, err)
		}
	}()

	in, err := o.Source.filestore.OpenReader(ctx, o.Source.RelativePath)
	if err != nil {
		return fmt.Errorf("error reading %q: %v", o.Source.AbsolutePath, err)
	}
	defer in.Close()

	if _, err := io.Copy(f, in); err != nil {
		return fmt.Errorf(
			"error downloading %s: %v",
			o.Source.AbsolutePath, err)
	}
	// We close the file to be sure it is fully written
	if err := f.Close(); err != nil {
		return fmt.Errorf("error writing temp file %q: %v", tempFilename, err)
	}
	f = nil

	// Verify the source hash
	sha256, err := hash.SHA256ForFile(tempFilename)
	if err != nil {
		return err
	}
	if sha256 != o.ManifestFile.SHA256 {
		return fmt.Errorf(
			"sha256 did not match for file %q: actual=%q expected=%q",
			o.Source.AbsolutePath, sha256, o.ManifestFile.SHA256)
	}

	// Upload to the destination
	return o.Dest.filestore.UploadFile(ctx, o.Dest.RelativePath, tempFilename)
}

// String is the pretty-printer for an operation, as used by dry-run.
func (o *copyFileOp) String() string {
	return fmt.Sprintf(
		"COPY %q to %q",
		o.Source.AbsolutePath, o.Dest.AbsolutePath)
}
