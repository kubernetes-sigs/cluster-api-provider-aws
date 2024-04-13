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
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	api "sigs.k8s.io/promo-tools/v4/api/files"
)

// GoogleCloudStorage is the provider for Google Cloud Storage (GCS).
var GoogleCloudStorage = &gcsProvider{}

type gcsProvider struct{}

var _ Provider = &gcsProvider{}

func (p *gcsProvider) Scheme() string {
	return api.GCSScheme
}

type gcsSyncFilestore struct {
	provider  *gcsProvider
	filestore *api.Filestore
	client    *storage.Client
	bucket    string
	prefix    string
}

// OpenReader opens an io.ReadCloser for the specified file.
func (s *gcsSyncFilestore) OpenReader(
	ctx context.Context,
	name string,
) (io.ReadCloser, error) {
	absolutePath := s.prefix + name
	return s.client.Bucket(s.bucket).Object(absolutePath).NewReader(ctx)
}

// UploadFile uploads a local file to the specified destination.
func (s *gcsSyncFilestore) UploadFile(ctx context.Context, dest, localFile string) error {
	absolutePath := s.prefix + dest

	gcsURL := s.provider.Scheme() + api.Backslash + s.bucket + "/" + absolutePath

	in, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("error opening %q: %v", localFile, err)
	}
	defer func() {
		if err := in.Close(); err != nil {
			logrus.Warnf("error closing %q: %v", localFile, err)
		}
	}()

	// Compute crc32 checksum for upload integrity
	var fileCRC32C uint32
	{
		hasher := crc32.New(crc32.MakeTable(crc32.Castagnoli))
		if _, err := io.Copy(hasher, in); err != nil {
			return fmt.Errorf("error computing crc32 checksum: %v", err)
		}
		fileCRC32C = hasher.Sum32()

		if _, err := in.Seek(0, 0); err != nil {
			return fmt.Errorf("error rewinding in file: %v", err)
		}
	}

	logrus.Infof("uploading to %s", gcsURL)

	w := s.client.Bucket(s.bucket).Object(absolutePath).NewWriter(ctx)

	w.CRC32C = fileCRC32C
	w.SendCRC32C = true

	// Much bigger chunk size for faster uploading
	w.ChunkSize = 128 * 1024 * 1024

	if _, err := io.Copy(w, in); err != nil {
		if err2 := w.Close(); err2 != nil {
			logrus.Warnf("error closing upload stream: %v", err)
			// TODO: Try to delete the possibly partially written file?
		}
		return fmt.Errorf("error uploading to %q: %v", gcsURL, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("error uploading to %q: %v", gcsURL, err)
	}

	return nil
}

// ListFiles returns all the file artifacts in the filestore, recursively.
func (s *gcsSyncFilestore) ListFiles(
	ctx context.Context,
) (map[string]*SyncFileInfo, error) {
	files := make(map[string]*SyncFileInfo)

	q := &storage.Query{Prefix: s.prefix}
	logrus.Infof("listing files in bucket %s with prefix %q", s.bucket, s.prefix)
	it := s.client.Bucket(s.bucket).Objects(ctx, q)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf(
				"error listing objects in %q: %v",
				s.filestore.Base, err)
		}
		name := obj.Name
		if !strings.HasPrefix(name, s.prefix) {
			return nil, fmt.Errorf(
				"found object %q without prefix %q",
				name, s.prefix)
		}

		file := &SyncFileInfo{}
		file.AbsolutePath = s.provider.Scheme() + api.Backslash + s.bucket + "/" + obj.Name
		file.RelativePath = strings.TrimPrefix(name, s.prefix)
		if obj.MD5 == nil {
			return nil, fmt.Errorf("MD5 not set on file %q", file.AbsolutePath)
		}

		file.MD5 = hex.EncodeToString(obj.MD5)
		file.Size = obj.Size
		file.filestore = s

		files[file.RelativePath] = file
	}

	return files, nil
}
