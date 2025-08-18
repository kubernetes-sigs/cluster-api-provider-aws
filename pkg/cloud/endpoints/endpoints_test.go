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
)

func TestParseFlags(t *testing.T) {
	testCases := []struct {
		name          string
		flagToParse   string
		expectedError error
	}{
		{
			name:          "no configuration",
			flagToParse:   "",
			expectedError: nil,
		},
		{
			name:          "single region, single service",
			flagToParse:   "us-iso:ec2=https://localhost:8080",
			expectedError: nil,
		},
		{
			name:          "single region, multiple services",
			flagToParse:   "us-iso:ec2=https://localhost:8080,sts=https://elbhost:8080",
			expectedError: nil,
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
		},
		{
			name:          "invalid config",
			flagToParse:   "us-isoec2=localhost",
			expectedError: errServiceEndpointSigningRegion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ParseFlag(tc.flagToParse)

			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("did not expect correct error: got %v, expected %v", err, tc.expectedError)
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
