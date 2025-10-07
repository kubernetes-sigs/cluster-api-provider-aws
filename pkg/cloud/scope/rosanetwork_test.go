/*
 Copyright The Kubernetes Authors.

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
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
)

func TestNewROSANetworkScope(t *testing.T) {
	g := NewGomegaWithT(t)

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	infrav1.AddToScheme(scheme)
	expinfrav1.AddToScheme(scheme)

	clusterControllerIdentity := &infrav1.AWSClusterControllerIdentity{
		TypeMeta: metav1.TypeMeta{
			Kind: string(infrav1.ControllerIdentityKind),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: infrav1.AWSClusterControllerIdentitySpec{
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}

	staticSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "static-secret",
			Namespace: system.GetManagerNamespace(),
		},
		Data: map[string][]byte{
			"AccessKeyID":     []byte("access-key-id"),
			"SecretAccessKey": []byte("secret-access-key"),
		},
	}

	clusterStaticIdentity := &infrav1.AWSClusterStaticIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster-static-identity",
		},
		Spec: infrav1.AWSClusterStaticIdentitySpec{
			SecretRef: "static-secret",
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{},
			},
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(clusterControllerIdentity, staticSecret, clusterStaticIdentity).Build()

	rosaNetwork := expinfrav1.ROSANetwork{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ROSANetwork",
			APIVersion: "v1beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-rosa-net",
			Namespace: "test-namespace",
		},
		Spec: expinfrav1.ROSANetworkSpec{
			IdentityRef: &infrav1.AWSIdentityReference{
				Name: "default",
				Kind: "AWSClusterControllerIdentity",
			},
		},
		Status: expinfrav1.ROSANetworkStatus{},
	}

	rosaNetScopeParams := ROSANetworkScopeParams{
		Client:         fakeClient,
		ControllerName: "test-rosanet-controller",
		Logger:         logger.NewLogger(klog.Background()),
		ROSANetwork:    &rosaNetwork,
	}

	rosaNetScope, err := NewROSANetworkScope(rosaNetScopeParams)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rosaNetScope.ControllerName()).To(Equal("test-rosanet-controller"))
	g.Expect(rosaNetScope.InfraCluster()).To(Equal(&rosaNetwork))
	g.Expect(rosaNetScope.InfraClusterName()).To(Equal("test-rosa-net"))
	g.Expect(rosaNetScope.Namespace()).To(Equal("test-namespace"))
	g.Expect(rosaNetScope.IdentityRef()).To(Equal(rosaNetwork.Spec.IdentityRef))
	g.Expect(rosaNetScope.Session()).ToNot(BeNil())

	// AWSClusterStaticIdentity
	rosaNetwork.Spec.IdentityRef.Name = "cluster-static-identity"
	rosaNetwork.Spec.IdentityRef.Kind = "AWSClusterStaticIdentity"
	rosaNetScope, err = NewROSANetworkScope(rosaNetScopeParams)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rosaNetScope.Session()).ToNot(BeNil())
}
