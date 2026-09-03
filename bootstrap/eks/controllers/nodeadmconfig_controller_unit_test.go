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

package controllers

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
)

func setupScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = clusterv1beta1.AddToScheme(scheme)
	_ = eksbootstrapv1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	return scheme
}

func TestGenerateCustomHybridUserdataResolvesFiles(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	version := "1.29"
	config := &eksbootstrapv1.NodeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "default",
		},
		Spec: eksbootstrapv1.NodeadmConfigSpec{
			Files: []eksbootstrapv1.File{
				{
					Path: "/etc/resolved-file",
					ContentFrom: &eksbootstrapv1.FileSource{
						Secret: eksbootstrapv1.SecretFileSource{Name: "file-secret", Key: "content"},
					},
				},
			},
			Hybrid: &eksbootstrapv1.HybridOptions{
				CustomUserData: &eksbootstrapv1.CustomUserDataOptions{
					Template: `{{ range .Files }}{{.Path}}={{.Content}}{{ end }}`,
				},
			},
		},
	}
	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			EKSClusterName: "test-cluster",
			Region:         "us-west-2",
			Version:        &version,
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "file-secret", Namespace: "default"},
		Data:       map[string][]byte{"content": []byte("resolved-content")},
	}
	r := &NodeadmConfigReconciler{
		Client: fake.NewClientBuilder().WithScheme(setupScheme()).WithObjects(secret).Build(),
	}

	data, err := r.generateCustomHybridUserdata(ctx, config, controlPlane, "activation-id", "activation-code")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(data)).To(ContainSubstring("/etc/resolved-file=resolved-content"))
}

func TestGetActivationFromSecret(t *testing.T) {
	tests := []struct {
		name           string
		secretName     string
		secretData     map[string][]byte
		secretExists   bool
		expectError    bool
		expectedID     string
		expectedCode   string
		expectCondFail bool
	}{
		{
			name:       "valid secret with activation credentials",
			secretName: "test-activation",
			secretData: map[string][]byte{
				"activationId":   []byte("act-12345"),
				"activationCode": []byte("code-secret"),
			},
			secretExists:   true,
			expectError:    false,
			expectedID:     "act-12345",
			expectedCode:   "code-secret",
			expectCondFail: false,
		},
		{
			name:           "secret not found",
			secretName:     "missing-secret",
			secretExists:   false,
			expectError:    true,
			expectCondFail: true,
		},
		{
			name:       "secret missing activationId",
			secretName: "incomplete-secret",
			secretData: map[string][]byte{
				"activationCode": []byte("code-secret"),
			},
			secretExists: true,
			expectError:  true,
		},
		{
			name:       "secret missing activationCode",
			secretName: "incomplete-secret2",
			secretData: map[string][]byte{
				"activationId": []byte("act-12345"),
			},
			secretExists: true,
			expectError:  true,
		},
		{
			name:       "secret with empty activationId",
			secretName: "empty-id-secret",
			secretData: map[string][]byte{
				"activationId":   []byte(""),
				"activationCode": []byte("code-secret"),
			},
			secretExists: true,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme := setupScheme()

			// Build fake client with optional secret
			objs := []runtime.Object{}
			if tc.secretExists {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      tc.secretName,
						Namespace: "default",
					},
					Data: tc.secretData,
				}
				objs = append(objs, secret)
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(objs...).
				Build()

			reconciler := &NodeadmConfigReconciler{
				Client: fakeClient,
			}

			// Create a test config for condition tracking
			config := &eksbootstrapv1.NodeadmConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "default",
				},
			}

			activationID, activationCode, err := reconciler.getActivationFromSecret(
				context.Background(),
				"default",
				tc.secretName,
				config,
			)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				if tc.expectCondFail {
					cond := v1beta1conditions.Get(config, eksbootstrapv1.SSMActivationReadyCondition)
					g.Expect(cond).NotTo(BeNil())
					g.Expect(cond.Status).To(Equal(corev1.ConditionFalse))
				}
			} else {
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(activationID).To(Equal(tc.expectedID))
				g.Expect(activationCode).To(Equal(tc.expectedCode))
				cond := v1beta1conditions.Get(config, eksbootstrapv1.SSMActivationReadyCondition)
				g.Expect(cond).NotTo(BeNil())
				g.Expect(cond.Status).To(Equal(corev1.ConditionTrue))
			}
		})
	}
}

func TestStoreActivationSecret(t *testing.T) {
	tests := []struct {
		name               string
		secretName         string
		activationID       string
		activationCode     string
		existingSecret     bool
		existingSecretData map[string][]byte
		expectError        bool
	}{
		{
			name:           "create new secret",
			secretName:     "new-activation-secret",
			activationID:   "act-12345",
			activationCode: "code-secret",
			existingSecret: false,
			expectError:    false,
		},
		{
			name:           "update existing secret",
			secretName:     "existing-activation-secret",
			activationID:   "act-67890",
			activationCode: "code-new-secret",
			existingSecret: true,
			existingSecretData: map[string][]byte{
				"activationId":   []byte("old-id"),
				"activationCode": []byte("old-code"),
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			scheme := setupScheme()

			// Build fake client
			objs := []runtime.Object{}
			if tc.existingSecret {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      tc.secretName,
						Namespace: "default",
					},
					Data: tc.existingSecretData,
				}
				objs = append(objs, secret)
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(objs...).
				Build()

			reconciler := &NodeadmConfigReconciler{
				Client: fakeClient,
			}

			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
			}

			config := &eksbootstrapv1.NodeadmConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config",
					Namespace: "default",
					UID:       types.UID("test-uid"),
				},
			}

			err := reconciler.storeActivationSecret(
				context.Background(),
				cluster,
				config,
				tc.secretName,
				tc.activationID,
				tc.activationCode,
			)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())

				// Verify the secret was created/updated correctly
				secret := &corev1.Secret{}
				err = fakeClient.Get(context.Background(), types.NamespacedName{
					Name:      tc.secretName,
					Namespace: "default",
				}, secret)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(string(secret.Data[ssmActivationIDKey])).To(Equal(tc.activationID))
				g.Expect(string(secret.Data[ssmActivationCodeKey])).To(Equal(tc.activationCode))
				g.Expect(secret.Labels[clusterv1.ClusterNameLabel]).To(Equal(cluster.Name))
				g.Expect(secret.OwnerReferences).To(HaveLen(1))
				g.Expect(secret.OwnerReferences[0].Name).To(Equal(config.Name))
			}
		})
	}
}

func TestReconcileDeleteWithoutFinalizer(t *testing.T) {
	g := NewWithT(t)
	scheme := setupScheme()

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()

	reconciler := &NodeadmConfigReconciler{
		Client: fakeClient,
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cluster",
			Namespace: "default",
		},
	}

	// Config without finalizer should return early
	config := &eksbootstrapv1.NodeadmConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-config",
			Namespace:  "default",
			Finalizers: []string{}, // No finalizer
		},
	}

	result, err := reconciler.reconcileDelete(context.Background(), cluster, config)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result.Requeue).To(BeFalse())
	g.Expect(result.RequeueAfter).To(BeZero())
}
