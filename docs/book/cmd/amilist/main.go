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

// Package main provides a Lambda function to list AMIs and upload them to an S3 bucket.
package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/ami"
)

var (
	ctx = context.TODO()
	svc *s3manager.Uploader
)

const (
	bucket = "cluster-api-aws-amis.sigs.k8s.io"
)

func init() {
	sess, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	s3Client := s3.NewFromConfig(sess)

	svc = s3manager.NewUploader(s3Client)
}

func main() {
	lambda.Start(LambdaHandler)
}

// LambdaHandler defines a Lambda function handler.
func LambdaHandler() error {
	amis, err := ami.List(
		ami.ListInput{},
	)

	if err != nil {
		ctrl.Log.Error(err, "error fetching AMIs")
		return err
	}

	data, err := json.MarshalIndent(amis, "", "  ")
	if err != nil {
		ctrl.Log.Error(err, "error marshalling")
		return err
	}

	_, err = svc.Upload(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucket),
		Key:    aws.String("amis.json"),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		ctrl.Log.Error(err, "error uploading data")
	}

	return err
}
