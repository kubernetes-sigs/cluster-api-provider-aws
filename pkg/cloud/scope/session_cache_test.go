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

package scope

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
)

// sessionMetadataStub is a minimal cloud.SessionMetadata implementation that lets a
// test drive sessionForClusterWithRegion as if it were called from different
// controllers for the same cluster.
type sessionMetadataStub struct {
	namespace      string
	infraName      string
	infraCluster   cloud.ClusterObject
	identityRef    *infrav1.AWSIdentityReference
	controllerName string
}

func (s *sessionMetadataStub) Namespace() string                          { return s.namespace }
func (s *sessionMetadataStub) InfraClusterName() string                   { return s.infraName }
func (s *sessionMetadataStub) InfraCluster() cloud.ClusterObject          { return s.infraCluster }
func (s *sessionMetadataStub) IdentityRef() *infrav1.AWSIdentityReference { return s.identityRef }
func (s *sessionMetadataStub) ControllerName() string                     { return s.controllerName }

// TestSessionRebuiltPerControllerAfterCredentialRotation ensures that rotating the
// credentials behind an identity invalidates the cached session of every controller,
// not just the first one to observe the change.
//
// providerCache is global while sessionCache is keyed per controller. Gating the
// session rebuild on a providerCache miss meant that whichever controller reconciled
// first repopulated providerCache, and every other controller then took the
// "nothing changed" path and kept returning a session built from credentials that no
// longer exist.
func TestSessionRebuiltPerControllerAfterCredentialRotation(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	sessionCache.Clear()
	providerCache.Clear()
	t.Cleanup(func() {
		sessionCache.Clear()
		providerCache.Clear()
	})

	scheme, err := setupScheme()
	g.Expect(err).NotTo(HaveOccurred())

	awsCluster := &infrav1.AWSCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       "AWSCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rotating-cluster",
			Namespace: "default",
		},
	}
	staticIdentity := &infrav1.AWSClusterStaticIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       string(infrav1.ClusterStaticIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "rotating-identity",
		},
		Spec: infrav1.AWSClusterStaticIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
			SecretRef: "source-secret-ref",
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source-secret-ref",
			Namespace: system.GetManagerNamespace(),
		},
		Data: map[string][]byte{
			"AccessKeyID":     []byte("AKIAIOSFODNN7EXAMPLE"),
			"SecretAccessKey": []byte("original-secret-access-key"),
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRESTMapper(newTestRESTMapper()).
		WithObjects(awsCluster, staticIdentity, secret).
		Build()

	newScope := func(controllerName string) cloud.SessionMetadata {
		return &sessionMetadataStub{
			namespace:      awsCluster.Namespace,
			infraName:      awsCluster.Name,
			infraCluster:   awsCluster,
			controllerName: controllerName,
			identityRef: &infrav1.AWSIdentityReference{
				Name: staticIdentity.Name,
				Kind: infrav1.ClusterStaticIdentityKind,
			},
		}
	}

	clusterController := newScope("awscluster")
	machineController := newScope("awsmachine")

	accessKeyFor := func(scope cloud.SessionMetadata) string {
		t.Helper()
		cfg, _, err := sessionForClusterWithRegion(k8sClient, scope, "us-west-2", logger.NewLogger(klog.Background()))
		g.Expect(err).NotTo(HaveOccurred())
		creds, err := cfg.Credentials.Retrieve(ctx)
		g.Expect(err).NotTo(HaveOccurred())
		return creds.AccessKeyID
	}

	// Both controllers build and cache a session from the original credentials.
	g.Expect(accessKeyFor(clusterController)).To(Equal("AKIAIOSFODNN7EXAMPLE"))
	g.Expect(accessKeyFor(machineController)).To(Equal("AKIAIOSFODNN7EXAMPLE"))

	// Re-read before mutating: buildAWSClusterStaticIdentity patches the secret, so the
	// local copy is stale by now.
	rotated := &corev1.Secret{}
	g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(secret), rotated)).To(Succeed())
	rotated.Data["AccessKeyID"] = []byte("AKIAI44QH8DHBEXAMPLE")
	rotated.Data["SecretAccessKey"] = []byte("rotated-secret-access-key")
	g.Expect(k8sClient.Update(ctx, rotated)).To(Succeed())

	// The first controller to reconcile after the rotation repopulates providerCache.
	g.Expect(accessKeyFor(clusterController)).To(Equal("AKIAI44QH8DHBEXAMPLE"))

	// Any other controller must also rebuild instead of serving its pre-rotation
	// session. Before the fix this returned the stale AKIAIOSFODNN7EXAMPLE forever.
	g.Expect(accessKeyFor(machineController)).To(Equal("AKIAI44QH8DHBEXAMPLE"))
}
