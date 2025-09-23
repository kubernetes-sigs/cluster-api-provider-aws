/*
Copyright 2020 The Kubernetes Authors.

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

package endpoints

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFlags(t *testing.T) {
	testCases := []struct {
		name                        string
		flagToParse                 string
		expectedError               error
		expectedServiceEndpointsMap map[string]serviceEndpoint
	}{
		{
			name:                        "no configuration",
			flagToParse:                 "",
			expectedServiceEndpointsMap: make(map[string]serviceEndpoint),
		},
		{
			name:                        "single region, single service",
			flagToParse:                 "us-iso:ec2=https://localhost:8080",
			expectedServiceEndpointsMap: map[string]serviceEndpoint{"EC2": {ServiceID: "EC2", URL: "https://localhost:8080", SigningRegion: "us-iso"}},
		},
		{
			name:        "single region, multiple services with v1 service IDs",
			flagToParse: "us-iso:s3=https://s3.com,elasticloadbalancing=https://elb.com,ec2=https://ec2.com,tagging=https://tagging.com,sqs=https://sqs.com,events=https://events.com,eks=https://eks.com,ssm=https://ssm.com,sts=https://sts.com,secretsmanager=https://secretmanager.com",
			expectedServiceEndpointsMap: map[string]serviceEndpoint{
				"EC2":                         {ServiceID: "EC2", URL: "https://ec2.com", SigningRegion: "us-iso"},
				"EKS":                         {ServiceID: "EKS", URL: "https://eks.com", SigningRegion: "us-iso"},
				"Elastic Load Balancing":      {ServiceID: "Elastic Load Balancing", URL: "https://elb.com", SigningRegion: "us-iso"},
				"Elastic Load Balancing v2":   {ServiceID: "Elastic Load Balancing v2", URL: "https://elb.com", SigningRegion: "us-iso"},
				"EventBridge":                 {ServiceID: "EventBridge", URL: "https://events.com", SigningRegion: "us-iso"},
				"Resource Groups Tagging API": {ServiceID: "Resource Groups Tagging API", URL: "https://tagging.com", SigningRegion: "us-iso"},
				"S3":                          {ServiceID: "S3", URL: "https://s3.com", SigningRegion: "us-iso"},
				"SQS":                         {ServiceID: "SQS", URL: "https://sqs.com", SigningRegion: "us-iso"},
				"SSM":                         {ServiceID: "SSM", URL: "https://ssm.com", SigningRegion: "us-iso"},
				"STS":                         {ServiceID: "STS", URL: "https://sts.com", SigningRegion: "us-iso"},
				"Secrets Manager":             {ServiceID: "Secrets Manager", URL: "https://secretmanager.com", SigningRegion: "us-iso"},
			},
		},
		{
			name:        "single region, multiple services with v2 service IDs",
			flagToParse: "us-iso:S3=https://s3.com,Elastic Load Balancing=https://elb.com,Elastic Load Balancing v2=https://elbv2.com,EC2=https://ec2.com,Resource Groups Tagging API=https://tagging.com,SQS=https://sqs.com,EventBridge=https://events.com,EKS=https://eks.com,SSM=https://ssm.com,STS=https://sts.com,Secrets Manager=https://secretmanager.com",
			expectedServiceEndpointsMap: map[string]serviceEndpoint{
				"EC2":                         {ServiceID: "EC2", URL: "https://ec2.com", SigningRegion: "us-iso"},
				"EKS":                         {ServiceID: "EKS", URL: "https://eks.com", SigningRegion: "us-iso"},
				"Elastic Load Balancing":      {ServiceID: "Elastic Load Balancing", URL: "https://elb.com", SigningRegion: "us-iso"},
				"Elastic Load Balancing v2":   {ServiceID: "Elastic Load Balancing v2", URL: "https://elbv2.com", SigningRegion: "us-iso"},
				"EventBridge":                 {ServiceID: "EventBridge", URL: "https://events.com", SigningRegion: "us-iso"},
				"Resource Groups Tagging API": {ServiceID: "Resource Groups Tagging API", URL: "https://tagging.com", SigningRegion: "us-iso"},
				"S3":                          {ServiceID: "S3", URL: "https://s3.com", SigningRegion: "us-iso"},
				"SQS":                         {ServiceID: "SQS", URL: "https://sqs.com", SigningRegion: "us-iso"},
				"SSM":                         {ServiceID: "SSM", URL: "https://ssm.com", SigningRegion: "us-iso"},
				"STS":                         {ServiceID: "STS", URL: "https://sts.com", SigningRegion: "us-iso"},
				"Secrets Manager":             {ServiceID: "Secrets Manager", URL: "https://secretmanager.com", SigningRegion: "us-iso"},
			},
		},
		{
			name:          "single region, duplicate service",
			flagToParse:   "us-iso:ec2=https://localhost:8080,ec2=https://elbhost:8080",
			expectedError: errServiceEndpointDuplicateServiceID,
		},
		{
			name:          "single region, non-valid URI",
			flagToParse:   "us-iso:ec2=fdsfs",
			expectedError: errServiceEndpointURL,
		},
		{
			name:          "multiples regions",
			flagToParse:   "us-iso:ec2=https://localhost:8080,sts=https://elbhost:8080;gb-iso:ec2=https://localhost:8080,sts=https://elbhost:8080",
			expectedError: nil,
			expectedServiceEndpointsMap: map[string]serviceEndpoint{
				"EC2": {ServiceID: "EC2", URL: "https://localhost:8080", SigningRegion: "gb-iso"},
				"STS": {ServiceID: "STS", URL: "https://elbhost:8080", SigningRegion: "gb-iso"},
			},
		},
		{
			name:          "invalid config",
			flagToParse:   "us-isoec2=localhost",
			expectedError: errServiceEndpointSigningRegion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer t.Cleanup(func() {
				serviceEndpointsMap = make(map[string]serviceEndpoint)
			})

			err := ParseFlag(tc.flagToParse)

			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("did not expect correct error: got %v, expected %v", err, tc.expectedError)
			}

			if err == nil {
				if !cmp.Equal(serviceEndpointsMap, tc.expectedServiceEndpointsMap) {
					t.Fatalf("expected serviceEndpointsMap: %#v, but got: %#v", tc.expectedServiceEndpointsMap, serviceEndpointsMap)
				}
			}
		})
	}
}

func TestGetPartitionFromRegion(t *testing.T) {
	testCases := []struct {
		name          string
		region        string
		expectedValue string
	}{
		{
			name:          "should return gov partition",
			region:        "us-gov-east-1",
			expectedValue: "aws-us-gov",
		},
		{
			name:          "should return cn partition",
			region:        "cn-north-1",
			expectedValue: "aws-cn",
		},
		{
			name:          "should return iso partition",
			region:        "us-iso-east-1",
			expectedValue: "aws-iso",
		},
		{
			name:          "should return iso-b partition",
			region:        "us-isob-east-1",
			expectedValue: "aws-iso-b",
		},
		{
			name:          "should return default partition for valid region",
			region:        "us-west-2",
			expectedValue: "aws",
		},
		{
			name:          "should return default partition for invalid region",
			region:        "us-west-3",
			expectedValue: "aws",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := GetPartitionFromRegion(tc.region)

			if value != tc.expectedValue {
				t.Fatalf("did not get expected value: got %v, expected %v", value, tc.expectedValue)
			}
		})
	}
}
