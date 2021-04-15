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

package scope

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2/klogr"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/identity"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestIsClusterPermittedToUsePrincipal(t *testing.T) {
	testCases := []struct {
		name             string
		clusterNamespace string
		allowedNs        *infrav1.AllowedNamespaces
		setup            func(client.Client, *testing.T)
		expectedResult   bool
		expectErr        bool
	}{
		{
			name:             "All clusters are permitted to use identity if allowedNamespaces is empty",
			clusterNamespace: "default",
			allowedNs:        &infrav1.AllowedNamespaces{},
			expectedResult:   true,
			expectErr:        false,
		},
		{
			name:             "No clusters are permitted to use identity if allowedNamespaces is nil",
			clusterNamespace: "default",
			allowedNs:        nil,
			expectedResult:   false,
			expectErr:        false,
		},
		{
			name:             "A namespace is permitted if allowedNamespaces list has it",
			clusterNamespace: "match",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: []string{"match"},
				Selector:      metav1.LabelSelector{},
			},
			setup: func(c client.Client, t *testing.T) {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "match",
					},
				}
				ns.SetGroupVersionKind(infrav1.GroupVersion.WithKind("Namespace"))
				err := c.Create(context.Background(), ns)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResult: true,
			expectErr:      false,
		},
		{
			name:             "A namespace is not permitted if allowedNamespaces list does not have it",
			clusterNamespace: "default",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: []string{"nomatch"},
				Selector:      metav1.LabelSelector{},
			},
			setup: func(c client.Client, t *testing.T) {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				}
				ns.SetGroupVersionKind(infrav1.GroupVersion.WithKind("Namespace"))
				err := c.Create(context.Background(), ns)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResult: false,
			expectErr:      false,
		},
		{
			name:             "A namespace is not permitted if allowedNamespaces list and selector do not have it",
			clusterNamespace: "default",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: []string{"nomatch"},
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{"ns": "nomatchlabel"},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "match",
					},
				}
				ns.SetGroupVersionKind(infrav1.GroupVersion.WithKind("Namespace"))
				err := c.Create(context.Background(), ns)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResult: false,
			expectErr:      false,
		},
		{
			name:             "A namespace is not permitted if allowedNamespaces list and selector do not have it",
			clusterNamespace: "default",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: nil,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{"ns": "nomatchlabel"},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				}
				ns.SetGroupVersionKind(infrav1.GroupVersion.WithKind("Namespace"))
				err := c.Create(context.Background(), ns)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResult: false,
			expectErr:      false,
		},
		{
			name:             "A namespace is permitted if allowedNamespaces list does not have it but selector matches its label",
			clusterNamespace: "default",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: []string{"noMatch"},
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{"ns": "matchlabel"},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "default",
						Labels: map[string]string{"ns": "matchlabel"},
					},
				}
				ns.SetGroupVersionKind(infrav1.GroupVersion.WithKind("Namespace"))
				err := c.Create(context.Background(), ns)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResult: true,
			expectErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			k8sClient := fake.NewFakeClientWithScheme(scheme)
			if tc.setup != nil {
				tc.setup(k8sClient, t)
			}
			result, err := isClusterPermittedToUsePrincipal(k8sClient, tc.allowedNs, tc.clusterNamespace)
			if tc.expectErr {
				g.Expect(err).ToNot(BeNil())
			} else {
				g.Expect(err).To(BeNil())
			}

			if tc.expectedResult != result {
				t.Fatal("Did not get expected result")
			}

		})
	}
}

func TestPrincipalParsing(t *testing.T) {
	// Create the scope.
	clusterScope, _ := NewClusterScope(ClusterScopeParams{
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
		},
		AWSCluster: &infrav1.AWSCluster{},
	},
	)

	testCases := []struct {
		name        string
		awsCluster  infrav1.AWSCluster
		identityRef *corev1.ObjectReference
		identity    runtime.Object
		setup       func(client.Client, *testing.T)
		expect      func([]identity.AWSPrincipalTypeProvider)
		expectError bool
	}{
		{
			name: "Default case - no Principal specified",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster1",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				Spec: infrav1.AWSClusterSpec{},
			},
			setup: func(c client.Client, t *testing.T) {
			},
			expect: func(providers []identity.AWSPrincipalTypeProvider) {
				if len(providers) != 0 {
					t.Fatalf("Expected 0 providers, got %v", len(providers))
				}
			},
		},
		{
			name: "Can get a session for a static Principal",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster2",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				Spec: infrav1.AWSClusterSpec{
					IdentityRef: &infrav1.AWSIdentityReference{
						Name: "static-identity",
						Kind: infrav1.ClusterStaticIdentityKind,
					},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				identity := &infrav1.AWSClusterStaticIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "static-identity",
					},
					Spec: infrav1.AWSClusterStaticIdentitySpec{
						SecretRef: corev1.SecretReference{
							Name:      "static-credentials-secret",
							Namespace: "default",
						},
						AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
							AllowedNamespaces: &infrav1.AllowedNamespaces{},
						},
					},
				}
				identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterStaticIdentity"))
				err := c.Create(context.Background(), identity)
				if err != nil {
					t.Fatal(err)
				}

				credentialsSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "static-credentials-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"AccessKeyID":     []byte("1234567890"),
						"SecretAccessKey": []byte("abcdefghijklmnop"),
						"SessionToken":    []byte("asdfasdfasdf"),
					},
				}
				credentialsSecret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Kind: "Secret", Version: "v1"})
				err = c.Create(context.Background(), credentialsSecret)
				if err != nil {
					t.Fatal(err)
				}
			},
			expect: func(providers []identity.AWSPrincipalTypeProvider) {
				if len(providers) != 1 {
					t.Fatalf("Expected 1 provider, got %v", len(providers))
				}
				provider := providers[0]
				p, ok := provider.(*identity.AWSStaticPrincipalTypeProvider)
				if !ok {
					t.Fatal("Expected providers to be of type AWSStaticPrincipalTypeProvider")
				}
				if p.AccessKeyID != "1234567890" {
					t.Fatalf("Expected AccessKeyID to be '%s', got '%s'", "1234567890", p.AccessKeyID)
				}
				if p.SecretAccessKey != "abcdefghijklmnop" {
					t.Fatalf("Expected SecretAccessKey to be '%s', got '%s'", "abcdefghijklmnop", p.SecretAccessKey)
				}
				if p.SessionToken != "asdfasdfasdf" {
					t.Fatalf("Expected SessionToken to be '%s', got '%s'", "asdfasdfasdf", p.SessionToken)
				}
			},
		},
		{
			name: "Can build a chain identity",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				Spec: infrav1.AWSClusterSpec{
					IdentityRef: &infrav1.AWSIdentityReference{
						Name: "role-identity",
						Kind: infrav1.ClusterRoleIdentityKind,
					},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				staticPrincipal := &infrav1.AWSClusterStaticIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "static-identity",
					},
					Spec: infrav1.AWSClusterStaticIdentitySpec{
						SecretRef: corev1.SecretReference{
							Name:      "static-credentials-secret",
							Namespace: "default",
						},
						AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
							AllowedNamespaces: &infrav1.AllowedNamespaces{},
						},
					},
				}
				staticPrincipal.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterStaticIdentity"))
				err := c.Create(context.Background(), staticPrincipal)
				if err != nil {
					t.Fatal(err)
				}

				credentialsSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "static-credentials-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"AccessKeyID":     []byte("1234567890"),
						"SecretAccessKey": []byte("abcdefghijklmnop"),
						"SessionToken":    []byte("asdfasdfasdf"),
					},
				}
				credentialsSecret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Kind: "Secret", Version: "v1"})
				err = c.Create(context.Background(), credentialsSecret)
				if err != nil {
					t.Fatal(err)
				}

				roleIdentity := &infrav1.AWSClusterRoleIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "role-identity",
					},
					Spec: infrav1.AWSClusterRoleIdentitySpec{
						AWSRoleSpec: infrav1.AWSRoleSpec{
							RoleArn:     "role-arn",
							SessionName: "test-session",
						},
						SourceIdentityRef: &infrav1.AWSIdentityReference{
							Name: "static-identity",
							Kind: infrav1.ClusterStaticIdentityKind,
						},
						AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
							AllowedNamespaces: &infrav1.AllowedNamespaces{},
						},
					},
				}
				roleIdentity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
				err = c.Create(context.Background(), roleIdentity)
				if err != nil {
					t.Fatal(err)
				}
			},
			expect: func(providers []identity.AWSPrincipalTypeProvider) {
				if len(providers) != 1 {
					t.Fatalf("Expected 1 providers, got %v", len(providers))
				}
			},
		},
		{
			name: "Can get a session for a role Principal",
			awsCluster: infrav1.AWSCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "AWSCluster",
				},
				Spec: infrav1.AWSClusterSpec{
					IdentityRef: &infrav1.AWSIdentityReference{
						Name: "role-identity",
						Kind: infrav1.ClusterRoleIdentityKind,
					},
				},
			},
			setup: func(c client.Client, t *testing.T) {
				identity := &infrav1.AWSClusterRoleIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "role-identity",
					},
					Spec: infrav1.AWSClusterRoleIdentitySpec{
						AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
							AllowedNamespaces: &infrav1.AllowedNamespaces{},
						},
						AWSRoleSpec: infrav1.AWSRoleSpec{
							RoleArn: "role-arn",
						},
					},
				}
				identity.SetGroupVersionKind(infrav1.GroupVersion.WithKind("AWSClusterRoleIdentity"))
				err := c.Create(context.Background(), identity)
				if err != nil {
					t.Fatal(err)
				}
			},
			expect: func(providers []identity.AWSPrincipalTypeProvider) {
				if len(providers) != 1 {
					t.Fatalf("Expected 1 providers, got %v", len(providers))
				}
				provider := providers[0]
				p, ok := provider.(*identity.AWSRolePrincipalTypeProvider)
				if !ok {
					t.Fatal("Expected providers to be of type AWSRolePrincipalTypeProvider")
				}
				if p.Principal.Spec.RoleArn != "role-arn" {
					t.Fatal(errors.Errorf("Expected Role Provider ARN to be 'role-arn', got '%s'", p.Principal.Spec.RoleArn))
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme, err := setupScheme()
			if err != nil {
				t.Fatal(err)
			}
			k8sClient := fake.NewFakeClientWithScheme(scheme)
			tc.setup(k8sClient, t)
			clusterScope.AWSCluster = &tc.awsCluster
			providers, err := getProvidersForCluster(context.Background(), k8sClient, clusterScope, klogr.New())
			if tc.expectError {
				if err == nil {
					t.Fatal("Expected an error but didn't get one")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				tc.expect(providers)
			}
		})
	}
}
