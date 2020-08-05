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

package v1alpha3

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaultingWebhook(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		resourceNS   string
		expectName   bool
		expectHash   bool
		expect       string
	}{
		{
			name:         "less than 100 chars",
			resourceName: "cluster1",
			resourceNS:   "default",
			expectName:   true,
			expectHash:   false,
			expect:       "default_cluster1",
		},
		{
			name:         "less than 100 chars, dot in name",
			resourceName: "team1.cluster1",
			resourceNS:   "default",
			expectName:   true,
			expectHash:   false,
			expect:       "default_team1_cluster1",
		},
		{
			name:         "more than 100 chars",
			resourceName: "ABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDE",
			resourceNS:   "default",
			expectName:   false,
			expectHash:   true,
			expect:       "capi_",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				ObjectMeta: v1.ObjectMeta{
					Name:      tc.resourceName,
					Namespace: tc.resourceNS,
				},
			}
			mcp.Default()

			actual := mcp.Spec.EKSClusterName
			g.Expect(actual).ToNot(BeEmpty())

			if tc.expectName {
				g.Expect(actual).To(Equal(tc.expect))
			}
			if tc.expectHash {
				g.Expect(strings.HasPrefix(actual, "capa_")).To(BeTrue())
			}
		})
	}
}

func TestValidatingWebhookCreate(t *testing.T) {
	tests := []struct {
		name           string
		eksClusterName string
		expectError    bool
	}{
		{
			name:           "ekscluster specified",
			eksClusterName: "default_cluster1",
			expectError:    false,
		},
		{
			name:           "ekscluster NOT specified",
			eksClusterName: "",
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mcp := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: tc.eksClusterName,
				},
			}
			err := mcp.ValidateCreate()

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestValidatingWebhookUpdate(t *testing.T) {
	tests := []struct {
		name           string
		oldClusterName string
		newClusterName string
		expectError    bool
	}{
		{
			name:           "ekscluster specified, same cluster names",
			oldClusterName: "default_cluster1",
			newClusterName: "default_cluster1",
			expectError:    false,
		},
		{
			name:           "ekscluster specified, different cluster names",
			oldClusterName: "default_cluster1",
			newClusterName: "default_cluster2",
			expectError:    true,
		},
		{
			name:           "old ekscluster specified, no new cluster name",
			oldClusterName: "default_cluster1",
			newClusterName: "",
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			newMCP := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: tc.newClusterName,
				},
			}
			oldMCP := &AWSManagedControlPlane{
				Spec: AWSManagedControlPlaneSpec{
					EKSClusterName: tc.oldClusterName,
				},
			}
			err := newMCP.ValidateUpdate(oldMCP)

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
