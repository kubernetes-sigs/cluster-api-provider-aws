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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/secret"
)

func TestCreateScope(t *testing.T) {
	RegisterTestingT(t)

	scheme, err := setupScheme()
	Expect(err).NotTo(HaveOccurred())

	testCases := []struct {
		name        string
		logger      *logr.Logger
		cluster     *clusterv1.Cluster
		client      client.Client
		opts        []RemoteScopeOption
		expectError bool
	}{
		{
			name:        "no cluster, should error",
			expectError: true,
		},
		{
			name:        "with cluster, no client, no options should error",
			cluster:     newCluster("cluster1"),
			expectError: true,
		},
		{
			name:        "with cluster, with client, no options should not error",
			cluster:     newCluster("cluster1"),
			expectError: false,
			client:      newFakeClient(scheme, newKubeconfigSecret("cluster1")),
		},
		{
			name:        "with cluster, no client, with client option should not error",
			cluster:     newCluster("cluster1"),
			expectError: false,
			opts:        []RemoteScopeOption{WithRemoteScopeTenantClient(newFakeClient(scheme))},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			httpmock.Activate()
			defer httpmock.Deactivate()

			resp, err := httpmock.NewJsonResponder(200, &metav1.APIGroup{})
			g.Expect(err).NotTo(HaveOccurred())

			httpmock.RegisterResponder(
				"GET",
				"http://test-cluster-api.nodomain.example.com:6443/apis?timeout=1m0s",
				resp)
			httpmock.RegisterResponder(
				"GET",
				"http://test-cluster-api.nodomain.example.com:6443/api?timeout=1m0s",
				resp)

			params := RemoteClusterScopeParams{
				Client:  tc.client,
				Cluster: tc.cluster,
				Logger:  tc.logger,
			}

			s, err := NewRemoteClusterScope(params, tc.opts...)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(s).NotTo(BeNil())
		})
	}
}

func newFakeClient(scheme *runtime.Scheme, objs ...client.Object) client.Client {
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
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
