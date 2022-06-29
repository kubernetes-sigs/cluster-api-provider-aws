/*
Copyright 2022 The Kubernetes Authors.

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
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/secret"
)

func TestCreateScope(t *testing.T) {
	RegisterTestingT(t)

	scheme, err := setupScheme()
	Expect(ekscontrolplanev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(err).NotTo(HaveOccurred())

	testCases := []struct {
		name         string
		logger       *logr.Logger
		cluster      *clusterv1.Cluster
		infraCluster client.Object
		client       client.Client
		opts         []ExternalResourceGCScopeOption
		expectError  bool
	}{
		{
			name:        "no cluster, should error",
			expectError: true,
		},
		{
			name:        "with cluster, no infra cluster,  no client, no options should error",
			cluster:     newCluster("cluster1"),
			expectError: true,
		},
		{
			name:         "with cluster, with infra cluster,  no client, no options should error",
			cluster:      newCluster("cluster1"),
			infraCluster: &unstructured.Unstructured{},
			expectError:  true,
		},
		{
			name:         "with cluster, with infra cluster, with client, no options should not error",
			cluster:      newCluster("cluster1"),
			infraCluster: newManagedInfra("cluster1"),
			expectError:  false,
			client:       newFakeClient(scheme, newKubeconfigSecret("cluster1")),
		},
		{
			name:         "with cluster, with infra cluster, with client, with client option should not error",
			cluster:      newCluster("cluster1"),
			infraCluster: newManagedInfra("cluster1"),
			client:       newFakeClient(scheme, newKubeconfigSecret("cluster1")),
			expectError:  false,
			opts:         []ExternalResourceGCScopeOption{WithExternalResourceGCScopeTenantClient(newFakeClient(scheme))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			// httpmock.Activate()
			// defer httpmock.Deactivate()

			// resp, err := httpmock.NewJsonResponder(200, &metav1.APIGroup{})
			// g.Expect(err).NotTo(HaveOccurred())

			// httpmock.RegisterResponder(
			// 	"GET",
			// 	"http://test-cluster-api.nodomain.example.com:6443/apis?timeout=1m0s",
			// 	resp)
			// httpmock.RegisterResponder(
			// 	"GET",
			// 	"http://test-cluster-api.nodomain.example.com:6443/api?timeout=1m0s",
			// 	resp)

			params := ExternalResourceGCScopeParams{
				Client:       tc.client,
				Cluster:      tc.cluster,
				InfraCluster: tc.infraCluster,
				Logger:       tc.logger,
			}

			s, err := NewExternalResourceGCScope(params, tc.opts...)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(s).NotTo(BeNil())
		})
	}
}

func TestRemoteClient(t *testing.T) {
	RegisterTestingT(t)

	scheme, err := setupScheme()
	Expect(ekscontrolplanev1.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(err).NotTo(HaveOccurred())

	httpmock.Activate()
	defer httpmock.Deactivate()

	resp, err := httpmock.NewJsonResponder(200, &metav1.APIGroup{})
	Expect(err).NotTo(HaveOccurred())

	httpmock.RegisterResponder(
		"GET",
		"http://test-cluster-api.nodomain.example.com:6443/apis?timeout=1m0s",
		resp)
	httpmock.RegisterResponder(
		"GET",
		"http://test-cluster-api.nodomain.example.com:6443/api?timeout=1m0s",
		resp)

	params := ExternalResourceGCScopeParams{
		Client:       newFakeClient(scheme, newKubeconfigSecret("cluster1")),
		Cluster:      newCluster("cluster1"),
		InfraCluster: newManagedInfra("cluster1"),
	}

	s, err := NewExternalResourceGCScope(params)
	Expect(err).NotTo(HaveOccurred())
	Expect(s).NotTo(BeNil())

	remoteClient, err := s.RemoteClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(remoteClient).NotTo(BeNil())
}

func newFakeClient(scheme *runtime.Scheme, objs ...client.Object) client.Client {
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func newManagedInfra(clusterName string) *ekscontrolplanev1.AWSManagedControlPlane {
	return &ekscontrolplanev1.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: ekscontrolplanev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: "default",
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{},
	}
}

func newKubeconfigSecret(clusterName string) *corev1.Secret {
	validKubeConfig := `
clusters:
- cluster:
    server: http://test-cluster-api.nodomain.example.com:6443
  name: test-cluster-api
contexts:
- context:
    cluster: test-cluster-api
    user: kubernetes-admin
  name: kubernetes-admin@test-cluster-api
current-context: kubernetes-admin@test-cluster-api
kind: Config
preferences: {}
users:
- name: kubernetes-admin
`
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-kubeconfig", clusterName),
			Namespace: metav1.NamespaceDefault,
		},
		Data: map[string][]byte{
			secret.KubeconfigDataName: []byte(validKubeConfig),
		},
	}
}
