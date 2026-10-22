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
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/identity"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

func TestIsClusterPermittedToUsePrincipal(t *testing.T) {
	testCases := []struct {
		name             string
		clusterNamespace string
		allowedNs        *infrav1.AllowedNamespaces
		setup            func(*testing.T, client.Client)
		obj              client.Object
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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
		{
			name:             "Cluster-scoped resources bypass allowedNamespaces check even with restricted allowedNamespaces",
			clusterNamespace: "default",
			allowedNs: &infrav1.AllowedNamespaces{
				NamespaceList: []string{"team-a", "team-b"},
			},
			setup: func(t *testing.T, c client.Client) {
				t.Helper()
			},
			obj: func() client.Object {
				obj := &infrav1.AWSClusterRoleIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster-scoped-identity",
					},
					Spec: infrav1.AWSClusterRoleIdentitySpec{
						AWSRoleSpec: infrav1.AWSRoleSpec{
							RoleArn: "arn:aws:iam::123456789012:role/test-role",
						},
					},
				}
				obj.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   infrav1.GroupVersion.Group,
					Version: infrav1.GroupVersion.Version,
					Kind:    "AWSClusterRoleIdentity",
				})
				return obj
			}(),
			expectedResult: true,
			expectErr:      false,
		},
		{
			name:             "Cluster-scoped resources bypass allowedNamespaces check even with nil allowedNamespaces",
			clusterNamespace: "default",
			allowedNs:        nil,
			setup: func(t *testing.T, c client.Client) {
				t.Helper()
			},
			obj: func() client.Object {
				obj := &infrav1.AWSClusterRoleIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster-scoped-identity-2",
					},
					Spec: infrav1.AWSClusterRoleIdentitySpec{
						AWSRoleSpec: infrav1.AWSRoleSpec{
							RoleArn: "arn:aws:iam::123456789012:role/test-role",
						},
					},
				}
				obj.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   infrav1.GroupVersion.Group,
					Version: infrav1.GroupVersion.Version,
					Kind:    "AWSClusterRoleIdentity",
				})
				return obj
			}(),
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
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRESTMapper(newTestRESTMapper()).
				Build()
			if tc.setup != nil {
				tc.setup(t, k8sClient)
			}

			// Create a test object (namespace-scoped AWSCluster) if not provided
			testObj := tc.obj
			if testObj == nil {
				testObj = &infrav1.AWSCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-cluster",
						Namespace: tc.clusterNamespace,
					},
				}
			}

			result, err := isClusterPermittedToUsePrincipal(k8sClient, tc.allowedNs, tc.clusterNamespace, testObj)
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
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRESTMapper(newTestRESTMapper()).
		Build()
	clusterScope, _ := NewClusterScope(ClusterScopeParams{
		Client: cl,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
		},
		AWSCluster: &infrav1.AWSCluster{Spec: infrav1.AWSClusterSpec{Region: "us-west-2"}},
	},
	)

	testCases := []struct {
		name        string
		awsCluster  infrav1.AWSCluster
		identityRef *corev1.ObjectReference
		identity    runtime.Object
		setup       func(*testing.T, client.Client)
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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()
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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

				identity := &infrav1.AWSClusterStaticIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "static-identity",
					},
					Spec: infrav1.AWSClusterStaticIdentitySpec{ //nolint:gosec // G101: test fixture, not real credentials
						SecretRef: "static-credentials-secret",
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
						Namespace: system.GetManagerNamespace(),
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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

				staticPrincipal := &infrav1.AWSClusterStaticIdentity{
					ObjectMeta: metav1.ObjectMeta{
						Name: "static-identity",
					},
					Spec: infrav1.AWSClusterStaticIdentitySpec{ //nolint:gosec // G101: test fixture, not real credentials
						SecretRef: "static-credentials-secret",
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
						Namespace: system.GetManagerNamespace(),
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
			setup: func(t *testing.T, c client.Client) {
				t.Helper()

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
			k8sClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRESTMapper(newTestRESTMapper()).
				Build()
			tc.setup(t, k8sClient)
			clusterScope.AWSCluster = &tc.awsCluster
			providers, err := getProvidersForCluster(context.Background(), k8sClient, clusterScope, clusterScope.Region(), logger.NewLogger(klog.Background()))
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

// Session conditions are saved by controller cleanup, including failed initialization.
func TestSessionConditionPersistence(t *testing.T) {
	type sessionScope interface {
		InitSession() error
		Close() error
		PatchSessionConditions(context.Context) error
	}
	for _, resource := range []struct {
		name        string
		object      func(*infrav1.AWSIdentityReference) cloud.ClusterObject
		createScope func(client.Client, cloud.ClusterObject) (sessionScope, error)
	}{
		{name: "AWSCluster", object: func(ref *infrav1.AWSIdentityReference) cloud.ClusterObject {
			return &infrav1.AWSCluster{Spec: infrav1.AWSClusterSpec{Region: "eu-west-1", IdentityRef: ref}}
		}, createScope: func(c client.Client, obj cloud.ClusterObject) (sessionScope, error) {
			return NewClusterScope(ClusterScopeParams{Client: c, Cluster: &clusterv1.Cluster{}, AWSCluster: obj.(*infrav1.AWSCluster), ControllerName: "session-persistence-test"})
		}},
		{name: "AWSManagedControlPlane", object: func(ref *infrav1.AWSIdentityReference) cloud.ClusterObject {
			return &ekscontrolplanev1.AWSManagedControlPlane{Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{Region: "eu-west-1", IdentityRef: ref}}
		}, createScope: func(c client.Client, obj cloud.ClusterObject) (sessionScope, error) {
			return NewManagedControlPlaneScope(ManagedControlPlaneScopeParams{Client: c, Cluster: &clusterv1.Cluster{}, ControlPlane: obj.(*ekscontrolplanev1.AWSManagedControlPlane), ControllerName: "session-persistence-test"})
		}},
		{name: "ROSAControlPlane", object: func(ref *infrav1.AWSIdentityReference) cloud.ClusterObject {
			return &rosacontrolplanev1.ROSAControlPlane{Spec: rosacontrolplanev1.RosaControlPlaneSpec{Region: "eu-west-1", IdentityRef: ref}}
		}, createScope: func(c client.Client, obj cloud.ClusterObject) (sessionScope, error) {
			return NewROSAControlPlaneScope(ROSAControlPlaneScopeParams{Client: c, Cluster: &clusterv1.Cluster{}, ControlPlane: obj.(*rosacontrolplanev1.ROSAControlPlane), ControllerName: "session-persistence-test"})
		}},
	} {
		for _, state := range []string{"success", "missing identity", "denied identity"} {
			t.Run(resource.name+"/"+state, func(t *testing.T) {
				g := NewWithT(t)
				ctx := context.Background()
				scheme := runtime.NewScheme()
				g.Expect(infrav1.AddToScheme(scheme)).To(Succeed())
				g.Expect(ekscontrolplanev1.AddToScheme(scheme)).To(Succeed())
				g.Expect(rosacontrolplanev1.AddToScheme(scheme)).To(Succeed())
				var ref *infrav1.AWSIdentityReference
				if state == "missing identity" {
					ref = &infrav1.AWSIdentityReference{Kind: infrav1.ClusterStaticIdentityKind, Name: "missing"}
				}
				if state == "denied identity" {
					ref = &infrav1.AWSIdentityReference{Kind: infrav1.ClusterRoleIdentityKind, Name: "denied"}
				}
				obj := resource.object(ref)
				obj.SetName(strings.ToLower(strings.ReplaceAll("session-test-"+resource.name+"-"+state, " ", "-")))
				obj.SetNamespace("default")
				v1beta1conditions.MarkTrue(obj, "Unrelated")
				denied := &infrav1.AWSClusterRoleIdentity{ObjectMeta: metav1.ObjectMeta{Name: "denied"}}
				mapper := newTestRESTMapper().(*meta.DefaultRESTMapper)
				mapper.Add(ekscontrolplanev1.GroupVersion.WithKind("AWSManagedControlPlane"), meta.RESTScopeNamespace)
				mapper.Add(rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"), meta.RESTScopeNamespace)
				c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithStatusSubresource(obj).WithObjects(obj, denied).Build()
				scoper, err := resource.createScope(c, obj)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(v1beta1conditions.Get(obj, infrav1.PrincipalCredentialRetrievedCondition)).To(BeNil())
				err = scoper.InitSession()
				if state == "success" {
					g.Expect(err).ToNot(HaveOccurred())
				} else {
					g.Expect(err).To(HaveOccurred())
				}
				stored := obj.DeepCopyObject().(cloud.ClusterObject)
				g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
				// Session initialization does not write: even failures wait for controller cleanup.
				g.Expect(v1beta1conditions.Get(stored, infrav1.PrincipalCredentialRetrievedCondition)).To(BeNil())
				g.Expect(scoper.Close()).To(Succeed())
				g.Expect(c.Get(ctx, client.ObjectKeyFromObject(obj), stored)).To(Succeed())
				reported := v1beta1conditions.Get(stored, infrav1.PrincipalCredentialRetrievedCondition)
				g.Expect(reported).ToNot(BeNil())
				if state == "success" {
					g.Expect(reported.Status).To(Equal(corev1.ConditionTrue))
				} else {
					g.Expect(reported.Status).To(Equal(corev1.ConditionFalse))
					g.Expect(reported.Reason).To(Equal(infrav1.PrincipalCredentialRetrievalFailedReason))
				}
				if state == "denied identity" {
					g.Expect(v1beta1conditions.IsFalse(stored, infrav1.PrincipalUsageAllowedCondition)).To(BeTrue())
				}
				g.Expect(v1beta1conditions.IsTrue(stored, "Unrelated")).To(BeTrue())
				if state == "success" {
					// Backfill when a successful session is cached but the condition is absent on the object.
					v1beta1conditions.Delete(stored, infrav1.PrincipalCredentialRetrievedCondition)
					g.Expect(c.Status().Update(ctx, stored)).To(Succeed())
					scoper, err = resource.createScope(c, stored)
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(scoper.InitSession()).To(Succeed())
					v1beta1conditions.MarkFalse(stored, "Unrelated", "ChangedInMemory", "Info", "")
					stored.SetAnnotations(map[string]string{"unrelated-change": "unsaved"})
					g.Expect(scoper.PatchSessionConditions(ctx)).To(Succeed())
					g.Expect(c.Get(ctx, client.ObjectKeyFromObject(stored), stored)).To(Succeed())
					g.Expect(v1beta1conditions.IsTrue(stored, infrav1.PrincipalCredentialRetrievedCondition)).To(BeTrue())
					g.Expect(v1beta1conditions.IsTrue(stored, "Unrelated")).To(BeTrue())
					g.Expect(stored.GetAnnotations()).To(BeEmpty())
					transition := v1beta1conditions.Get(stored, infrav1.PrincipalCredentialRetrievedCondition).LastTransitionTime
					scoper, err = resource.createScope(c, stored)
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(scoper.InitSession()).To(Succeed())
					g.Expect(scoper.PatchSessionConditions(ctx)).To(Succeed())
					g.Expect(c.Get(ctx, client.ObjectKeyFromObject(stored), stored)).To(Succeed())
					g.Expect(v1beta1conditions.Get(stored, infrav1.PrincipalCredentialRetrievedCondition).LastTransitionTime).To(Equal(transition))
				}
			})
		}
	}
}
