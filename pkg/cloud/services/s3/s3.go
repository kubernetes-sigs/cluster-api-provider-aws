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

package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"

	iam "sigs.k8s.io/cluster-api-provider-aws/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.S3Scope
	S3Client  s3iface.S3API
	STSClient stsiface.STSAPI
}

var DisabledError = errors.New("s3 management disabled")

func IsDisabledError(err error) bool {
	return err == DisabledError
}

var EmptyBucketError = errors.New("empty bucket name")

func IsEmptyBucketError(err error) bool {
	return err == EmptyBucketError
}

var EmptyKeyError = errors.New("empty key")

func IsEmptyKeyError(err error) bool {
	return err == EmptyKeyError
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

func (s *Service) ReconcileBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()

	if err := s.createBucketIfNotExist(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket exists")
	}

	if err := s.ensureBucketAccess(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket ACL ")
	}

	if err := s.ensureBucketPolicy(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket policy")
	}

	return nil
}

func (s *Service) DeleteBucket() error {
	if !s.bucketManagementEnabled() {
		return nil
	}

	bucketName := s.bucketName()
	if bucketName == "" {
		return EmptyBucketError
	}

	log := s.scope.WithValues("name", bucketName)

	log.Info("Deleting S3 Bucket")

	_, err := s.S3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
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

// Create will add a file to the s3 bucket which is private and server side encrypted.
func (s *Service) Create(key string, data []byte) (string, error) {
	if !s.bucketManagementEnabled() {
		return "", DisabledError
	}

	// server side encryption, acl defaults to private
	return s.create(&s3.PutObjectInput{
		Body:                 aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:               aws.String(s.scope.Bucket().Name),
		Key:                  aws.String(key),
		ServerSideEncryption: aws.String("aws:kms"),
	})
}

// CreatePublic will add file to the s3 bucket which is public and open to the world to access.
func (s *Service) CreatePublic(key string, data []byte) (string, error) {
	// acl public-read
	if !s.bucketManagementEnabled() {
		return "", DisabledError
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
		return "", EmptyBucketError
	}

	if aws.StringValue(putInput.Key) == "" {
		return "", EmptyKeyError
	}

	s.scope.Info("Creating public object", "bucket_name", aws.StringValue(putInput.Bucket), "key", aws.StringValue(putInput.Key))

	if _, err := s.S3Client.PutObject(putInput); err != nil {
		return "", errors.Wrap(err, "putting object")
	}

	objectURL := &url.URL{
		Scheme: "s3",
		Host:   aws.StringValue(putInput.Bucket),
		Path:   aws.StringValue(putInput.Key),
	}

	return objectURL.String(), nil
}

func (s *Service) Delete(key string) error {
	if !s.bucketManagementEnabled() {
		return DisabledError
	}

	if key == "" {
		return EmptyKeyError
	}

	bucketName := s.bucketName()
	if bucketName == "" {
		return EmptyBucketError
	}

	s.scope.Info("Deleting object", "bucket_name", bucketName, "key", key)

	_, err := s.S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.Wrap(err, "deleting S3 object")
	}

	switch aerr.Code() {
	case s3.ErrCodeNoSuchBucket:
	default:
		return errors.Wrap(aerr, "deleting S3 object")
	}

	return nil
}

func (s *Service) createBucketIfNotExist(bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket:          aws.String(bucketName),
		ObjectOwnership: aws.String(s3.ObjectOwnershipBucketOwnerPreferred),
	}

	_, err := s.S3Client.CreateBucket(input)
	if err == nil {
		s.scope.Info("Created bucket", "bucket_name", bucketName)

		return nil
	}

	aerr, ok := err.(awserr.Error)
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
	f := false
	input := &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls: aws.Bool(f),
		},
	}

	if _, err := s.S3Client.PutPublicAccessBlock(input); err != nil {
		return errors.Wrap(err, "enabling bucket public access")
	}

	s.scope.V(4).Info("Updated bucket ACL to allow public access", "bucket_name", bucketName)

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

	s.scope.V(4).Info("Updated bucket policy", "bucket_name", bucketName)

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

	statements := []iam.StatementEntry{
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
		{
			Sid:    "control-plane",
			Effect: iam.EffectAllow,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", *accountID.Account, bucket.ControlPlaneIAMInstanceProfile)},
			},
			Action:   []string{"s3:GetObject"},
			Resource: []string{fmt.Sprintf("arn:aws:s3:::%s/control-plane/*", bucketName)},
		},
	}

	for _, iamInstanceProfile := range bucket.NodesIAMInstanceProfiles {
		statements = append(statements, iam.StatementEntry{
			Sid:    iamInstanceProfile,
			Effect: iam.EffectAllow,
			Principal: map[iam.PrincipalType]iam.PrincipalID{
				iam.PrincipalAWS: []string{fmt.Sprintf("arn:aws:iam::%s:role/%s", *accountID.Account, iamInstanceProfile)},
			},
			Action:   []string{"s3:GetObject"},
			Resource: []string{fmt.Sprintf("arn:aws:s3:::%s/node/*", bucketName)},
		})
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
