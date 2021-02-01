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

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
)

func TestParseFlags(t *testing.T) {

	testCases := []struct {
		name           string
		flagToParse    string
		expectedOutput []scope.ServiceEndpoint
		expectedError  error
	}{
		{
			name:           "no configuration",
			flagToParse:    "",
			expectedOutput: nil,
			expectedError:  nil,
		},
		{
			name:        "single region, single service",
			flagToParse: "us-iso:ec2=https://localhost:8080",
			expectedOutput: []scope.ServiceEndpoint{
				{
					ServiceID:     "ec2",
					URL:           "https://localhost:8080",
					SigningRegion: "us-iso",
				},
			},
			expectedError: nil,
		},
		{
			name:        "single region, multiple services",
			flagToParse: "us-iso:ec2=https://localhost:8080,sts=https://elbhost:8080",
			expectedOutput: []scope.ServiceEndpoint{
				{
					ServiceID:     "ec2",
					URL:           "https://localhost:8080",
					SigningRegion: "us-iso",
				},
				{
					ServiceID:     "sts",
					URL:           "https://elbhost:8080",
					SigningRegion: "us-iso",
				},
			},
			expectedError: nil,
		},
		{
			name:           "single region, duplicate service",
			flagToParse:    "us-iso:ec2=https://localhost:8080,ec2=https://elbhost:8080",
			expectedOutput: nil,
			expectedError:  errServiceEndpointDuplicateServiceID,
		},
		{
			name:           "single region, non-valid URI",
			flagToParse:    "us-iso:ec2=fdsfs",
			expectedOutput: nil,
			expectedError:  errServiceEndpointURL,
		},
		{
			name:        "multiples regions",
			flagToParse: "us-iso:ec2=https://localhost:8080,sts=https://elbhost:8080;gb-iso:ec2=https://localhost:8080,sts=https://elbhost:8080",
			expectedOutput: []scope.ServiceEndpoint{
				{
					ServiceID:     "ec2",
					URL:           "https://localhost:8080",
					SigningRegion: "us-iso",
				},
				{
					ServiceID:     "sts",
					URL:           "https://elbhost:8080",
					SigningRegion: "us-iso",
				},
				{
					ServiceID:     "ec2",
					URL:           "https://localhost:8080",
					SigningRegion: "gb-iso",
				},
				{
					ServiceID:     "sts",
					URL:           "https://elbhost:8080",
					SigningRegion: "gb-iso",
				},
			},
			expectedError: nil,
		},
		{
			name:           "invalid config",
			flagToParse:    "us-isoec2=localhost",
			expectedOutput: nil,
			expectedError:  errServiceEndpointSigningRegion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ParseFlag(tc.flagToParse)

			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("did not expect correct error: got %v, expected %v", err, tc.expectedError)
			}

			if !endpointsEqual(out, tc.expectedOutput) {
				t.Fatalf("did not expect correct output: got %v, expected %v", out, tc.expectedOutput)
			}
		})
	}
}

func endpointsEqual(a, b []scope.ServiceEndpoint) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
