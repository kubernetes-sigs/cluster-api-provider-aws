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
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

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
	log := logger.NewLogger(klog.Background())

	roleProviderFor := func(accessKeyID string) *AWSRolePrincipalTypeProvider {
		sourceProvider := NewAWSStaticPrincipalTypeProvider(
			&infrav1.AWSClusterStaticIdentity{
				ObjectMeta: metav1.ObjectMeta{Name: "source-identity"},
				Spec:       infrav1.AWSClusterStaticIdentitySpec{SecretRef: "source-secret-ref"},
			},
			&corev1.Secret{
				Data: map[string][]byte{
					"AccessKeyID":     []byte(accessKeyID),
					"SecretAccessKey": []byte("secret-access-key-for-" + accessKeyID),
				},
			},
		)

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
			sourceProvider, "us-west-2", log,
		)
	}

	hashOf := func(p *AWSRolePrincipalTypeProvider) string {
		t.Helper()
		h, err := p.Hash()
		g.Expect(err).NotTo(HaveOccurred())
		return h
	}

	original := hashOf(roleProviderFor("AKIAIOSFODNN7EXAMPLE"))
	rotated := hashOf(roleProviderFor("AKIAI44QH8DHBEXAMPLE"))
	unchanged := hashOf(roleProviderFor("AKIAIOSFODNN7EXAMPLE"))

	// Rotating the source credentials must invalidate the cached provider.
	g.Expect(rotated).NotTo(Equal(original))

	// Identical input must still hash identically, otherwise the cache never hits.
	g.Expect(unchanged).To(Equal(original))
}
