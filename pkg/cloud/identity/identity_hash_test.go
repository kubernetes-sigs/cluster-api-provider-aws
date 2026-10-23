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

package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

const (
	originalAccessKeyID = "AKIAIOSFODNN7EXAMPLE" //nolint:gosec // G101: test fixture, not real credentials
	rotatedAccessKeyID  = "AKIAI44QH8DHBEXAMPLE" //nolint:gosec // G101: test fixture, not real credentials
	testRegion          = "us-west-2"
)

func staticProviderFor(accessKeyID string) *AWSStaticPrincipalTypeProvider {
	return NewAWSStaticPrincipalTypeProvider(
		&infrav1.AWSClusterStaticIdentity{
			ObjectMeta: metav1.ObjectMeta{Name: "source-identity"},
			Spec:       infrav1.AWSClusterStaticIdentitySpec{SecretRef: "source-secret-ref"},
		},
		&corev1.Secret{
			Data: map[string][]byte{
				"AccessKeyID":     []byte(accessKeyID),
				"SecretAccessKey": []byte(secretAccessKeyFor(accessKeyID)),
			},
		},
	)
}

func secretAccessKeyFor(accessKeyID string) string {
	return "secret-access-key-for-" + accessKeyID
}

func roleProviderFor(accessKeyID, region string) *AWSRolePrincipalTypeProvider {
	return NewAWSRolePrincipalTypeProvider(
		&infrav1.AWSClusterRoleIdentity{
			ObjectMeta: metav1.ObjectMeta{Name: "role-identity"},
			Spec: infrav1.AWSClusterRoleIdentitySpec{
				AWSRoleSpec: infrav1.AWSRoleSpec{
					RoleArn: "arn:aws:iam::123456789012:role/bootstrapper",
				},
				ExternalID: "external-id",
			},
		},
		staticProviderFor(accessKeyID), region, logger.NewLogger(klog.Background()),
	)
}

func hashOf(t *testing.T, p AWSPrincipalTypeProvider) string {
	t.Helper()
	h, err := p.Hash()
	NewWithT(t).Expect(err).NotTo(HaveOccurred())
	return h
}

// TestRolePrincipalHashCoversSourceCredentials verifies that a role provider's hash
// changes when the credentials it uses to assume the role change.
//
// Hash gob-encodes the provider, and gob only encodes exported fields. sourceProvider
// is unexported, so before this was mixed in explicitly the hash of a role provider was
// determined solely by the AWSClusterRoleIdentity object. Rotating the source secret
// while leaving that object untouched produced an identical hash, providerCache returned
// the provider built from the previous credentials, and every caller kept using them.
func TestRolePrincipalHashCoversSourceCredentials(t *testing.T) {
	g := NewWithT(t)

	original := hashOf(t, roleProviderFor(originalAccessKeyID, testRegion))
	rotated := hashOf(t, roleProviderFor(rotatedAccessKeyID, testRegion))
	unchanged := hashOf(t, roleProviderFor(originalAccessKeyID, testRegion))

	// Rotating the source credentials must invalidate the cached provider.
	g.Expect(rotated).NotTo(Equal(original))

	// Identical input must still hash identically, otherwise the cache never hits.
	g.Expect(unchanged).To(Equal(original))
}

// TestRolePrincipalHashCoversRegion verifies that two role providers built from the same
// AWSClusterRoleIdentity but for different regions do not collide in providerCache.
//
// region is unexported and therefore not gob-encoded. Two clusters in different regions
// sharing one AWSClusterRoleIdentity hashed identically, so the second cluster reused the
// cached provider of the first and assumed the role against the first cluster's regional
// STS endpoint. That breaks VPC STS endpoints and aws:RequestedRegion policy conditions,
// and misreports the calling region in CloudTrail.
func TestRolePrincipalHashCoversRegion(t *testing.T) {
	g := NewWithT(t)

	usEast := hashOf(t, roleProviderFor(originalAccessKeyID, "us-east-1"))
	euWest := hashOf(t, roleProviderFor(originalAccessKeyID, "eu-west-1"))

	g.Expect(euWest).NotTo(Equal(usEast))
}

// TestHashIsDigestOfProvider verifies that Hash returns a hex-encoded SHA-256 digest.
//
// It previously returned sha256.New().Sum(gobBytes), which appends the digest of the
// empty string to its argument instead of hashing it. The result was the full gob payload
// with 32 constant bytes attached, so the plaintext SecretAccessKey of a static identity
// ended up in every providerCache key and, once sessions record the hashes they were built
// from, in every session cache entry as well.
func TestHashIsDigestOfProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider AWSPrincipalTypeProvider
	}{
		{name: "static", provider: staticProviderFor(originalAccessKeyID)},
		{name: "role", provider: roleProviderFor(originalAccessKeyID, testRegion)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			hash := hashOf(t, tt.provider)

			decoded, err := hex.DecodeString(hash)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(decoded).To(HaveLen(sha256.Size))

			g.Expect(hash).NotTo(ContainSubstring(secretAccessKeyFor(originalAccessKeyID)),
				"hash must not carry the plaintext secret access key")
		})
	}
}
