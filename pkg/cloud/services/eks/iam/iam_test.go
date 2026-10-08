/*
Copyright 2026 The Kubernetes Authors.

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

package iam

import "testing"

func TestIsEKSOIDCIssuerHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{name: "legacy commercial", host: "oidc.eks.us-west-2.amazonaws.com", want: true},
		{name: "legacy GovCloud", host: "oidc.eks.us-gov-west-1.amazonaws.com", want: true},
		{name: "legacy China", host: "oidc.eks.cn-north-1.amazonaws.com.cn", want: true},
		{name: "dual-stack commercial", host: "oidc-eks.us-west-2.api.aws", want: true},
		{name: "dual-stack China", host: "oidc-eks.cn-north-1.api.amazonwebservices.com.cn", want: true},
		{name: "case insensitive", host: "OIDC-EKS.US-WEST-2.API.AWS", want: true},
		{name: "trailing dot", host: "oidc.eks.us-west-2.amazonaws.com.", want: true},
		{name: "non-EKS host", host: "example.com", want: false},
		{name: "missing region", host: "oidc.eks.amazonaws.com", want: false},
		{name: "nested region", host: "oidc.eks.us-west-2.example.amazonaws.com", want: false},
		{name: "lookalike suffix", host: "oidc.eks.us-west-2.amazonaws.com.example.com", want: false},
		{name: "lookalike dual-stack suffix", host: "oidc-eks.us-west-2.api.aws.example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEKSOIDCIssuerHost(tt.host); got != tt.want {
				t.Fatalf("isEKSOIDCIssuerHost(%q) = %t, want %t", tt.host, got, tt.want)
			}
		})
	}
}
