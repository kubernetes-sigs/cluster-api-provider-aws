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
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/hash"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.S3Scope
	S3Client  s3iface.S3API
	STSClient stsiface.STSAPI
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

	bucketName, err := s.bucketName()
	if err != nil {
		return errors.Wrap(err, "generating bucket name")
	}

	if err := s.createBucketIfNotExist(bucketName); err != nil {
		return errors.Wrap(err, "ensuring bucket exists")
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

	bucketName, err := s.bucketName()
	if err != nil {
		return errors.Wrap(err, "generating bucket name")
	}

	s.scope.Info("Deleting S3 Bucket", "name", bucketName)

	if _, err := s.S3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}); err != nil {
		return errors.Wrap(err, "deleting S3 bucket")
	}

	return nil
}

func (s *Service) Create(m *scope.MachineScope, data []byte) (string, error) {
	if !s.bucketManagementEnabled() {
		return "", errors.New("requested object creation but bucket management is not enabled")
	}

	if m == nil {
		return "", errors.New("machine scope can't be nil")
	}

	if len(data) == 0 {
		return "", errors.New("got empty data")
	}

	bucket, err := s.bucketName()
	if err != nil {
		return "", errors.Wrap(err, "generating bucket name")
	}
	key := s.bootstrapDataKey(m)

	if _, err := s.S3Client.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return "", errors.Wrap(err, "putting object")
	}

	objectURL := &url.URL{
		Scheme: "s3",
		Host:   bucket,
		Path:   key,
	}

	return objectURL.String(), nil
}

func (s *Service) Delete(m *scope.MachineScope) error {
	if !s.bucketManagementEnabled() {
		return errors.New("requested object creation but bucket management is not enabled")
	}

	if m == nil {
		return errors.New("machine scope can't be nil")
	}

	bucketName, err := s.bucketName()
	if err != nil {
		return errors.Wrap(err, "generating bucket name")
	}

	if _, err := s.S3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s.bootstrapDataKey(m)),
	}); err != nil {
		return errors.Wrap(err, "deleting object")
	}

	return nil
}

func (s *Service) createBucketIfNotExist(bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := s.S3Client.CreateBucket(input)
	if err == nil {
		return nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.Wrap(err, "creating S3 bucket")
	}

	switch aerr.Code() {
	// If bucket already exists, all good.
	// TODO: This will fail if bucket is shared with other cluster.
	case s3.ErrCodeBucketAlreadyOwnedByYou:
		return nil
	default:
		return errors.Wrap(aerr, "creating S3 bucket")
	}
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

	return nil
}

func (s *Service) bucketPolicy(bucketName string) (string, error) {
	accountID, err := s.STSClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "getting account ID")
	}

	bucket := s.scope.Bucket()

	statements := []map[string]interface{}{
		{
			"Sid":    "control-plane",
			"Effect": "Allow",
			"Principal": map[string]interface{}{
				"AWS": fmt.Sprintf("arn:aws:iam::%s:role/%s", *accountID.Account, bucket.ControlPlaneIAMInstanceProfile),
			},
			"Action": []string{
				"s3:GetObject",
			},
			"Resource": fmt.Sprintf("arn:aws:s3:::%s/control-plane/*", bucketName),
		},
	}

	for _, iamInstanceProfile := range bucket.NodesIAMInstanceProfiles {
		statements = append(statements, map[string]interface{}{
			"Sid":    iamInstanceProfile,
			"Effect": "Allow",
			"Principal": map[string]interface{}{
				"AWS": fmt.Sprintf("arn:aws:iam::%s:role/%s", *accountID.Account, iamInstanceProfile),
			},
			"Action": []string{
				"s3:GetObject",
			},
			"Resource": fmt.Sprintf("arn:aws:s3:::%s/node/*", bucketName),
		})
	}

	policy := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statements,
	}

	policyRaw, err := json.Marshal(policy)
	if err != nil {
		return "", errors.Wrap(err, "building bucket policy")
	}

	return string(policyRaw), nil
}

func (s *Service) bucketManagementEnabled() bool {
	return s.scope.Bucket().Create
}

const (
	s3MaxBucketNameLength = 63
)

func (s *Service) bucketName() (string, error) {
	if name := s.scope.Bucket().Name; name != "" {
		return name, nil
	}

	name := fmt.Sprintf("%s-%s", s.scope.Namespace(), s.scope.KubernetesClusterName())
	if len(name) < s3MaxBucketNameLength {
		return name, nil
	}

	suffix := "-k8s"

	shortName, err := hash.Base36TruncatedHash(name, s3MaxBucketNameLength-len(suffix))
	if err != nil {
		return "", errors.Wrap(err, "unable to create S3 bucket name")
	}

	return fmt.Sprintf("%s%s", shortName, suffix), nil
}

func (s *Service) bootstrapDataKey(m *scope.MachineScope) string {
	// Use machine name as object key.
	return path.Join(m.Role(), m.Name())
}
