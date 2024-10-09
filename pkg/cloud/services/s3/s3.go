/*
Copyright 2021 The Kubernetes Authors.

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

// Package s3 provides a way to interact with AWS S3.
package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iam "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
)

// AwsDefaultRegion is the default AWS region.
const (
	AwsDefaultRegion   string = "us-east-1"
	forbiddenErrorCode string = "Forbidden"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.S3Scope
	S3Client  s3iface.S3API
	STSClient stsiface.STSAPI
}

var errS3ManagementDisabled = errors.New("s3 management disabled")

// IsS3ManagementDisabledError returns true if the given error is of type errS3ManagementDisabled.
func IsS3ManagementDisabledError(err error) bool {
	return errors.Is(err, errS3ManagementDisabled)
}

var errBucketNameUndefined = errors.New("bucket name not defined")

// IsBucketNameUndefinedError returns true if the given error is of type errBucketNameUndefined.
func IsBucketNameUndefinedError(err error) bool {
	return errors.Is(err, errBucketNameUndefined)
}

var errObjectKeyIsEmpty = errors.New("empty key")

// IsObjectKeyIsEmptyError returns true if the given error is of type errObjectKeyIsEmpty.
func IsObjectKeyIsEmptyError(err error) bool {
	return errors.Is(err, errObjectKeyIsEmpty)
}

// NewService returns a new service given the api clients.
func NewService(s3Scope scope.S3Scope) *Service {
	s3Client := scope.NewS3Client(s3Scope, s3Scope, s3Scope, s3Scope.InfraCluster())
	STSClient := scope.NewSTSClient(s3Scope, s3Scope, s3Scope, s3Scope.InfraCluster())

	return &Service{
		scope:     s3Scope,
		S3Client:  s3Client,
		STSClient: STSClient,
	}
}

// ReconcileBucket reconciles the S3 bucket.
func (s *Service) ReconcileBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()

	if err := s.createBucketIfNotExist(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket exists")
	}

	if err := s.tagBucket(bucketName); err != nil {
		return errors.Wrap(err, "tagging bucket")
	}

	if err := s.ensureBucketAccess(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket ACL ")
	}

	if err := s.ensureBucketPolicy(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket policy")
	}

	return nil
}

// DeleteBucket deletes the S3 bucket.
func (s *Service) DeleteBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()
	if bucketName == "" {
		return errBucketNameUndefined
	}

	log := s.scope.WithValues("name", bucketName)

	log.Info("Deleting S3 Bucket")

	_, err := s.S3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return nil
	}

	var aerr awserr.Error
	ok := errors.As(err, &aerr)
	if !ok {
		return errors.Wrap(err, "deleting S3 bucket")
	}

	switch aerr.Code() {
	case s3.ErrCodeNoSuchBucket:
		log.Info("Bucket already removed")
	case "BucketNotEmpty":
		log.Info("Bucket not empty, skipping removal")
	default:
		return errors.Wrap(aerr, "deleting S3 bucket")
	}

	return nil
}

// Create will add a machine bootstrap file to the S3 bucket which is private and server side encrypted.
func (s *Service) Create(m *scope.MachineScope, data []byte) (string, error) {
	if m == nil {
		return "", errors.New("machine scope can't be nil")
	}
	return s.CreatePrivateKey(s.bootstrapDataKey(m), data)
}

// CreatePrivateKey will add file to the S3 bucket which is private and server side encrypted.
func (s *Service) CreatePrivateKey(key string, data []byte) (string, error) {
	if !s.bucketManagementEnabled() {
		return "", errS3ManagementDisabled
	}

	if len(data) == 0 {
		return "", errors.New("got empty data")
	}

	// server side encryption, acl defaults to private
	return s.create(&s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:               aws.String(s.scope.Bucket().Name),
		Key:                  aws.String(key),
		ServerSideEncryption: aws.String("aws:kms"),
	})
}

// CreatePublicKey will add file to the s3 bucket which is public and open to the world to access.
func (s *Service) CreatePublicKey(key string, data []byte) (string, error) {
	// acl public-read
	if !s.bucketManagementEnabled() {
		return "", errS3ManagementDisabled
	}

	if len(data) == 0 {
		return "", errors.New("got empty data")
	}

	return s.create(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket: aws.String(s.scope.Bucket().Name),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
	})
}

func (s *Service) create(putInput *s3.PutObjectInput) (string, error) {
	if aws.StringValue(putInput.Bucket) == "" {
		return "", errBucketNameUndefined
	}

	if aws.StringValue(putInput.Key) == "" {
		return "", errObjectKeyIsEmpty
	}

	s.scope.Info("Creating public object", "bucket_name", aws.StringValue(putInput.Bucket), "key", aws.StringValue(putInput.Key))

	if _, err := s.S3Client.PutObject(putInput); err != nil {
		return "", errors.Wrap(err, "putting object")
	}

	if exp := s.scope.Bucket().PresignedURLDuration; exp != nil {
		s.scope.Info("Generating presigned URL", "bucket_name", aws.StringValue(putInput.Bucket), "key", aws.StringValue(putInput.Key))
		req, _ := s.S3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: putInput.Bucket,
			Key:    putInput.Key,
		})
		return req.Presign(exp.Duration)
	}

	objectURL := &url.URL{
		Scheme: "s3",
		Host:   aws.StringValue(putInput.Bucket),
		Path:   aws.StringValue(putInput.Key),
	}

	return objectURL.String(), nil
}

// Delete deletes the object from the S3 bucket.
func (s *Service) Delete(m *scope.MachineScope) error {
	if m == nil {
		return errors.New("machine scope can't be nil")
	}
	return s.DeleteKey(s.bootstrapDataKey(m))
}

// DeleteKey takes a key which is a s3 path to an object e.g. /path/file.ext.
func (s *Service) DeleteKey(key string) error {
	if !s.bucketManagementEnabled() {
		return errS3ManagementDisabled
	}

	if key == "" {
		return errObjectKeyIsEmpty
	}

	bucket := s.bucketName()
	if bucket == "" {
		return errBucketNameUndefined
	}

	_, err := s.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case forbiddenErrorCode:
				// In the case that the IAM policy does not have sufficient
				// permissions to get the object, we will attempt to delete it
				// anyway for backwards compatibility reasons.
				s.scope.Debug("Received 403 forbidden from S3 HeadObject call. If GetObject permission has been granted to the controller but not ListBucket, object is already deleted. Attempting deletion anyway in case GetObject permission hasn't been granted to the controller but DeleteObject has.", "bucket", bucket, "key", key)

				if err := s.deleteObject(bucket, key); err != nil {
					return err
				}

				s.scope.Debug("Delete object call succeeded despite missing GetObject permission", "bucket", bucket, "key", key)

				return nil
			case "NotFound":
				s.scope.Debug("Either bucket or object does not exist", "bucket", bucket, "key", key)
				return nil
			case s3.ErrCodeNoSuchKey:
				s.scope.Debug("Object already deleted", "bucket", bucket, "key", key)
				return nil
			case s3.ErrCodeNoSuchBucket:
				s.scope.Debug("Bucket does not exist", "bucket", bucket)
				return nil
			}
		}
		return err
	}

	s.scope.Info("Deleting object", "bucket_name", bucket, "key", key)

	return s.deleteObject(bucket, key)
}

func (s *Service) deleteObject(bucket, key string) error {
	if _, err := s.S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		if ptr.Deref(s.scope.Bucket().BestEffortDeleteObjects, false) {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case forbiddenErrorCode, "AccessDenied":
					s.scope.Debug("Ignoring deletion error", "bucket", bucket, "key", key, "error", aerr.Message())
					return nil
				}
			}
		}
		return errors.Wrap(err, "deleting S3 object")
	}

	return nil
}

func (s *Service) createBucketIfNotExist(bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket:          aws.String(bucketName),
		ObjectOwnership: aws.String(s3.ObjectOwnershipBucketOwnerPreferred),
	}

	// See https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html#AmazonS3-CreateBucket-request-LocationConstraint.
	if s.scope.Region() != AwsDefaultRegion {
		input.CreateBucketConfiguration = &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(s.scope.Region()),
		}
	}

	_, err := s.S3Client.CreateBucket(input)
	if err == nil {
		s.scope.Info("Created bucket", "bucket_name", bucketName)

		return nil
	}

	var aerr awserr.Error
	ok := errors.As(err, &aerr)
	if !ok {
		return errors.Wrap(err, "creating S3 bucket")
	}

	switch aerr.Code() {
	// If bucket already exists, all good.
	//
	// TODO: This will fail if bucket is shared with other cluster.
	case s3.ErrCodeBucketAlreadyOwnedByYou:
		return nil
	default:
		return errors.Wrap(aerr, "creating S3 bucket")
	}
}

func (s *Service) ensureBucketAccess(bucketName string) error {
	input := &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls: aws.Bool(false),
		},
	}

	if _, err := s.S3Client.PutPublicAccessBlock(input); err != nil {
		return errors.Wrap(err, "enabling bucket public access")
	}

	s.scope.GetLogger().Info("Updated bucket ACL to allow public access", "bucket_name", bucketName)

	return nil
}

func (s *Service) ensureBucketPolicy(bucketName string) error {
	bucketPolicy, err := s.bucketPolicy(bucketName)
	if err != nil {
		return errors.Wrap(err, "generating Bucket policy")
	}

	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(bucketPolicy),
	}

	if _, err := s.S3Client.PutBucketPolicy(input); err != nil {
		return errors.Wrap(err, "creating S3 bucket policy")
	}

	s.scope.Trace("Updated bucket policy", "bucket_name", bucketName)

	return nil
}

func (s *Service) tagBucket(bucketName string) error {
	taggingInput := &s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketName),
		Tagging: &s3.Tagging{
			TagSet: nil,
		},
	}

	tags := infrav1.Build(infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        nil,
		Role:        aws.String("node"),
		Additional:  s.scope.AdditionalTags(),
	})

	for key, value := range tags {
		taggingInput.Tagging.TagSet = append(taggingInput.Tagging.TagSet, &s3.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	sort.Slice(taggingInput.Tagging.TagSet, func(i, j int) bool {
		return *taggingInput.Tagging.TagSet[i].Key < *taggingInput.Tagging.TagSet[j].Key
	})

	_, err := s.S3Client.PutBucketTagging(taggingInput)
	if err != nil {
		return err
	}

	s.scope.Trace("Tagged bucket", "bucket_name", bucketName)

	return nil
}

// bucketPolicy grants access to get/put objects the cluster needs including a per cluster subdir in case two clusters share the same bucket.
// /<clustername> contains cluster wide object e.g. oidc configs for irsa.
// /control-plane contains ignite configs for control-plane nodes stored per node id.
// /node contains ignite configs for worker nodes stored per node id.
func (s *Service) bucketPolicy(bucketName string) (string, error) {
	accountID, err := s.STSClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "getting account ID")
	}

	bucket := s.scope.Bucket()
	partition := system.GetPartitionFromRegion(s.scope.Region())

	statements := []iam.StatementEntry{
		{
			Sid:    "ForceSSLOnlyAccess",
			Effect: iam.EffectDeny,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{"*"},
			},
			Action:   []string{"s3:*"},
			Resource: []string{fmt.Sprintf("arn:%s:s3:::%s/*", partition, bucketName)},
			Condition: iam.Conditions{
				"Bool": map[string]interface{}{
					"aws:SecureTransport": false,
				},
			},
		},
		{
			// grant access to the /<clustername> folder to the control plane nodes
			Sid:    s.scope.Name(),
			Effect: iam.EffectAllow,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", *accountID.Account, bucket.ControlPlaneIAMInstanceProfile)},
			},
			Action:   []string{"s3:GetObject", "s3:PutObject"},
			Resource: []string{fmt.Sprintf("arn:aws:s3:::%s/%s/*", bucketName, s.scope.Name())},
		},
	}

	if bucket.PresignedURLDuration == nil {
		if bucket.ControlPlaneIAMInstanceProfile != "" {
			statements = append(statements, iam.StatementEntry{
				Sid:    "control-plane",
				Effect: iam.EffectAllow,
				Principal: map[iam.PrincipalType]iam.PrincipalID{
					iam.PrincipalAWS: []string{fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, *accountID.Account, bucket.ControlPlaneIAMInstanceProfile)},
				},
				Action:   []string{"s3:GetObject"},
				Resource: []string{fmt.Sprintf("arn:%s:s3:::%s/control-plane/*", partition, bucketName)},
			})
		}

		for _, iamInstanceProfile := range bucket.NodesIAMInstanceProfiles {
			statements = append(statements, iam.StatementEntry{
				Sid:    iamInstanceProfile,
				Effect: iam.EffectAllow,
				Principal: map[iam.PrincipalType]iam.PrincipalID{
					iam.PrincipalAWS: []string{fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, *accountID.Account, iamInstanceProfile)},
				},
				Action:   []string{"s3:GetObject"},
				Resource: []string{fmt.Sprintf("arn:%s:s3:::%s/node/*", partition, bucketName)},
			})
		}
	}

	policy := iam.PolicyDocument{
		Version:   "2012-10-17",
		Statement: statements,
	}

	policyRaw, err := json.Marshal(policy)
	if err != nil {
		return "", errors.Wrap(err, "building bucket policy")
	}

	return string(policyRaw), nil
}

func (s *Service) bucketManagementEnabled() bool {
	return s.scope.Bucket() != nil
}

func (s *Service) bucketName() string {
	return s.scope.Bucket().Name
}

func (s *Service) bootstrapDataKey(m *scope.MachineScope) string {
	// Use machine name as object key.
	return path.Join(m.Role(), m.Name())
}
