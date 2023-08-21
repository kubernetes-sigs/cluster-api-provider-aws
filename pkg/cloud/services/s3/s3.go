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
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	"net/url"
	iam "sigs.k8s.io/cluster-api-provider-aws/iam/api/v1beta1"
	awslogs "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/logs"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/utils"
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

var ExternalBucketError = errors.New("external bucket")

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

	if !s.isBucketDeleteRequired() {
		return nil
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

	if aerr.Code() == "BucketRegionError" {
		// TODO: retry object delete with correct region
		s.scope.V(1).Info("retry delete with correct region BucketRegionError")
		s3Client, err := s.accessBucketFromDifferentRegion()
		if err != nil {
			return errors.Wrap(err, "deleting S3 object")
		}

		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err == nil {
			return nil
		}
	}

	aerr, ok = err.(awserr.Error)
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

func (s *Service) accessBucketFromDifferentRegion() (s3Client *s3.S3, err error) {
	s3region, err := s3manager.GetBucketRegion(context.Background(), s.scope.Session(), s.scope.Bucket().Name, s.scope.Region())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			s.scope.Info("unable to find bucket", "bucket", s.scope.Bucket().Name, "region not found")
		}
		return nil, err
	}
	s.scope.Info("accessBucketFromDifferentRegion ", "bucket: ", s.scope.Bucket().Name, "region: ", s3region)
	s3Client = s3.New(s.scope.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(s.scope)).WithLogger(awslogs.NewWrapLogr(s.scope)).WithEndpointResolver(utils.CustomEndpointResolverForAWSIRSA(s3region)))
	return s3Client, nil
}

func (s *Service) isBucketDeleteRequired() bool {
	bucketTagging, _ := s.S3Client.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: aws.String(s.scope.Bucket().Name),
	})

	if len(bucketTagging.TagSet) == 0 {
		s.scope.Info("Skipping deletion for external bucket")
		return false
	}

	for _, t := range bucketTagging.TagSet {
		k := "sigs.k8s.io/cluster-api-provider-aws/cluster/" + s.scope.InfraClusterName()
		if *t.Key == k && *t.Value == "owned" {
			s.scope.Info("found self owned bucket. Proceed for deletion")
			return true
		}
	}
	return false
}

func (s *Service) createBucketIfNotExist(bucketName string) error {

	// TODO: Add tag for CAPA created bucket and for external bucket
	input := &s3.CreateBucketInput{
		Bucket:          aws.String(bucketName),
		ObjectOwnership: aws.String(s3.ObjectOwnershipBucketOwnerPreferred),
	}

	_, err := s.S3Client.CreateBucket(input)
	if err == nil {
		s.scope.Info("Created bucket", "bucket_name", bucketName)

		s.scope.Info("Add bucket tagging to", "bucket_name", bucketName)
		k := "sigs.k8s.io/cluster-api-provider-aws/cluster/" + s.scope.InfraClusterName()
		tag := s3.Tag{
			Key:   aws.String(k),
			Value: aws.String("owned"),
		}

		tagSet := make([]*s3.Tag, 0)
		tagSet = append(tagSet, &tag)
		tInput := s3.PutBucketTaggingInput{
			Bucket:  aws.String(bucketName),
			Tagging: &s3.Tagging{TagSet: tagSet},
		}

		_, err = s.S3Client.PutBucketTagging(&tInput)
		if err != nil {
			return errors.Wrap(err, "updating S3 bucket with tag")
		}
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
		aerr, ok := err.(awserr.Error)
		if !ok {
			return errors.Wrap(err, "enabling bucket public access in different region")
		}
		if aerr.Code() == "BucketRegionError" {
			// TODO: Should we try to modify access and policy for bucket or just avoid doing modification without errors
			s.scope.Info("accessBucketFromDifferentRegion BucketRegionError")
			s3Client, err := s.accessBucketFromDifferentRegion()
			if err != nil {
				s.scope.Info("Error while accessing bucket from different region ", "bucket: ", s.scope.Bucket().Name, err)
				return errors.Wrap(err, "enabling bucket public access in different region")
			}
			_, err = s3Client.PutPublicAccessBlock(input)
			if err == nil {
				s.scope.V(1).Info("Accessing bucket from different region, PutPublicAccessBlock")
				return nil
			}
			return nil
		}
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
		aerr, ok := err.(awserr.Error)
		if !ok {
			return errors.Wrap(err, "enabling bucket public access in different region")
		}
		if aerr.Code() == "BucketRegionError" {

			// TODO: Should we try to modify access and policy for bucket or just avoid doing modification without errors
			s.scope.Info("Access bucket from different region")
			s3Client, err := s.accessBucketFromDifferentRegion()
			if err != nil {
				s.scope.Info("Error while accessing bucket from different region ", "bucket: ", s.scope.Bucket().Name, err)
				return errors.Wrap(err, "enabling bucket public access in different region")
			}
			_, err = s3Client.PutBucketPolicy(input)
			if err == nil {
				s.scope.V(1).Info("Accessing bucket from different region PutBucketPolicy")
				return nil
			}
			return nil
		}
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
