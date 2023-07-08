/*
Copyright 2023 The Kubernetes Authors.

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
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"

	api "sigs.k8s.io/promo-tools/v3/api/files"
)

// S3Storage is the provider for Amazon S3.
var S3Storage = &s3Provider{}

type s3Provider struct{}

func (p *s3Provider) Scheme() string {
	return api.S3Scheme
}

type s3SyncFilestore struct {
	provider  *s3Provider
	filestore *api.Filestore
	client    *s3.S3
	bucket    string
	prefix    string
}

// openS3Filestore opens a filestore backed by Amazon S3 (S3)

func (p *s3Provider) OpenFilestore(ctx context.Context, filestore *api.Filestore, useServiceAccount, config bool) (syncFilestore, error) { //nolint: revive
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

	bucket := u.Host
	// We send requests direct to the bucket region;
	// it's more efficient and it's required for regional buckets.
	bucketRegion, err := p.findRegionForBucket(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("finding region for bucket: %w", err)
	}

	awsConfig := aws.NewConfig()
	awsConfig = awsConfig.WithRegion(bucketRegion)
	if !useServiceAccount {
		awsConfig = awsConfig.WithCredentials(credentials.AnonymousCredentials)
	}
	awsConfig = awsConfig.WithCredentialsChainVerboseErrors(true)

	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("error building S3 session: %w", err)
	}

	client := s3.New(awsSession)

	prefix := strings.TrimPrefix(u.Path, "/")
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	s := &s3SyncFilestore{
		provider:  p,
		filestore: filestore,
		client:    client,
		bucket:    bucket,
		prefix:    prefix,
	}
	return s, nil
}

// findRegionForBucket returns the region in which the bucket is located.
func (p *s3Provider) findRegionForBucket(ctx context.Context, bucket string) (string, error) {
	// Pick a region to query, defaulting to the "normal" AWS_REGION env var if set.
	lookupRegion := os.Getenv("AWS_REGION")
	if lookupRegion == "" {
		// We have to query some region, us-east-2 is pretty reliable.
		lookupRegion = "us-east-2"
	}

	// This is an unauthenticated request (it just does a HEAD on the bucket),
	// but we still force anonymous credentials to be safe.
	awsConfig := aws.NewConfig()
	awsConfig = awsConfig.WithRegion(lookupRegion)
	awsConfig = awsConfig.WithCredentials(credentials.AnonymousCredentials)
	awsConfig = awsConfig.WithCredentialsChainVerboseErrors(true)

	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return "", fmt.Errorf("error creating AWS session: %w", err)
	}

	bucketRegion, err := s3manager.GetBucketRegion(ctx, awsSession, bucket, lookupRegion)
	if err != nil {
		return "", fmt.Errorf("error finding s3 region for bucket %q: %w", bucket, err)
	}

	return bucketRegion, nil
}

// OpenReader opens an io.ReadCloser for the specified file.
func (s *s3SyncFilestore) OpenReader(
	ctx context.Context,
	name string,
) (io.ReadCloser, error) {
	key := s.prefix + name
	req := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	}
	obj, err := s.client.GetObjectWithContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error reading object %q: %w", key, err)
	}
	return obj.Body, nil
}

// UploadFile uploads a local file to the specified destination.
func (s *s3SyncFilestore) UploadFile(ctx context.Context, dest, localFile string) error {
	key := s.prefix + dest

	s3URL := s.provider.Scheme() + "://" + s.bucket + "/" + key

	stat, err := os.Stat(localFile)
	if err != nil {
		return fmt.Errorf("error getting stat of %q: %w", localFile, err)
	}

	f, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("error opening %q: %w", localFile, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Warnf("error closing %q: %v", localFile, err)
		}
	}()

	contentLength := stat.Size()

	req := &s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &key,
		ContentLength: &contentLength,
		Metadata:      make(map[string]*string),
	}

	// Compute hashes for upload integrity and for metadata
	hashes, err := computeHashes(f)
	if err != nil {
		return err
	}

	req.ChecksumSHA256 = aws.String(base64.StdEncoding.EncodeToString(hashes.SHA256))

	// TODO: Any more hashes?  Very cheap to compute now...
	req.Metadata["content-hash-md5"] = aws.String(hex.EncodeToString(hashes.MD5))
	req.Metadata["content-hash-sha256"] = aws.String(hex.EncodeToString(hashes.SHA256))
	req.Metadata["content-hash-sha512"] = aws.String(hex.EncodeToString(hashes.SHA512))

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("error rewinding file: %w", err)
	}
	req.Body = f

	logrus.Infof("uploading to %s", s3URL)

	response, err := s.client.PutObjectWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("error uploading %q: %w", s3URL, err)
	}

	logrus.Debugf("uploaded to %s %v", s3URL, response)

	if got, want := aws.StringValue(response.ChecksumSHA256), base64.StdEncoding.EncodeToString(hashes.SHA256); got != want {
		// AWS should check this for us, but we double check it here.
		return fmt.Errorf("checksum mismatch on upload of %q: got %q, want %q", s3URL, got, want)
	}

	expectedETag := "\"" + hex.EncodeToString(hashes.MD5) + "\""
	if got, want := aws.StringValue(response.ETag), expectedETag; got != want {
		// We do a simple upload so that the etag is the md5, but we double check that it worked here
		return fmt.Errorf("unexpected etag after upload of %q: got %q, want %q", s3URL, got, want)
	}

	return nil
}

// ListFiles returns all the file artifacts in the filestore, recursively.
func (s *s3SyncFilestore) ListFiles(
	ctx context.Context,
) (map[string]*SyncFileInfo, error) {
	prefix := s.provider.Scheme() + "://" + s.bucket + "/" + s.prefix

	logrus.Infof("listing files under %s", prefix)

	files := make(map[string]*SyncFileInfo)

	req := &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &s.prefix,
	}

	var errors []error
	objectCallback := func(obj *s3.Object) error {
		name := aws.StringValue(obj.Key)
		if !strings.HasPrefix(name, s.prefix) {
			return fmt.Errorf(
				"found object %q without prefix %q",
				name, s.prefix)
		}

		file := &SyncFileInfo{}
		file.AbsolutePath = s.provider.Scheme() + "://" + s.bucket + "/" + name
		file.RelativePath = strings.TrimPrefix(name, s.prefix)

		md5 := aws.StringValue(obj.ETag)
		md5 = strings.Trim(md5, "\"")
		if md5 == "" {
			return fmt.Errorf("MD5 not set on file %q", file.AbsolutePath)
		}

		// Check that this at least looks like an md5
		if len(md5) != 32 {
			return fmt.Errorf("unexpected md5 (%q) on file %q", md5, file.AbsolutePath)
		}

		file.MD5 = md5
		file.Size = aws.Int64Value(obj.Size)
		file.filestore = s

		files[file.RelativePath] = file
		return nil
	}
	pageCallback := func(page *s3.ListObjectsV2Output, hasNextPage bool) bool {
		for _, obj := range page.Contents {
			err := objectCallback(obj)
			if err != nil {
				errors = append(errors, err)
				// stop iteration immediately on error
				return false
			}
		}
		return true
	}
	if err := s.client.ListObjectsV2PagesWithContext(ctx, req, pageCallback); err != nil {
		return nil, fmt.Errorf("error listing objects: %w", err)
	}

	if len(errors) != 0 {
		return nil, errors[0]
	}

	return files, nil
}

type Hashes struct {
	SHA256 []byte
	SHA512 []byte
	MD5    []byte
	Length int64
}

func computeHashes(in io.ReadSeeker) (*Hashes, error) {
	hasherSHA256 := sha256.New()
	hasherSHA512 := sha512.New()
	hasherMD5 := md5.New()

	hasher := io.MultiWriter(hasherMD5, hasherSHA256, hasherSHA512)

	n, err := io.Copy(hasher, in)
	if err != nil {
		return nil, fmt.Errorf("error hashing: %w", err)
	}

	if _, err := in.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("error rewinding file: %w", err)
	}
	return &Hashes{
		SHA256: hasherSHA256.Sum(nil),
		SHA512: hasherSHA512.Sum(nil),
		MD5:    hasherMD5.Sum(nil),
		Length: n,
	}, nil
}
