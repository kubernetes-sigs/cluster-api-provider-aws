// +build e2e

/*
Copyright 2021 The Kubernetes Authors.

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

package shared

import (
	"context"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
)

const (
	credsSecretName = "e2e-account-creds"
	capaNamespace   = "capa-system"
	eksNamespace    = "capa-eks-control-plane-system"
	idName          = "e2e-account"
)

func SetupStaticCredentials(ctx context.Context, namespace *corev1.Namespace, e2eCtx *E2EContext) {
	Expect(ctx).NotTo(BeNil(), "ctx is required for SetupStaticCredentials")
	Expect(namespace).NotTo(BeNil(), "namespace is required for SetupStaticCredentials")
	Expect(e2eCtx).NotTo(BeNil(), "e2eCtx is required for SetupStaticCredentials")
	Expect(e2eCtx.Environment.BootstrapAccessKey).NotTo(BeNil(), "e2eCtx.Environment.BootstrapAccessKey is required for SetupStaticCredentials")

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      credsSecretName,
			Namespace: capaNamespace,
		},
		StringData: map[string]string{
			"AccessKeyID":     *e2eCtx.Environment.BootstrapAccessKey.AccessKeyId,
			"SecretAccessKey": *e2eCtx.Environment.BootstrapAccessKey.SecretAccessKey,
		},
	}

	client := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	Byf("Creating credentials secret %s in namespace %s", secret.Name, secret.Namespace)
	Eventually(func() error {
		return client.Create(ctx, secret)
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-create-identity")...).Should(Succeed())

	id := &infrav1.AWSClusterStaticIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      idName,
			Namespace: namespace.Name,
		},
		Spec: infrav1.AWSClusterStaticIdentitySpec{
			SecretRef: credsSecretName,
			AWSClusterIdentitySpec: infrav1.AWSClusterIdentitySpec{
				AllowedNamespaces: &infrav1.AllowedNamespaces{
					NamespaceList: []string{namespace.Name},
				},
			},
		},
	}

	Byf("Creating AWSClusterStaticIdentity %s in namespace %s", id.Name, namespace.Name)
	Eventually(func() error {
		return client.Create(ctx, id)
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-create-identity")...).Should(Succeed())
}

func CleanupStaticCredentials(ctx context.Context, namespace *corev1.Namespace, e2eCtx *E2EContext) {
	Expect(ctx).NotTo(BeNil(), "ctx is required for SetupStaticCredentials")
	Expect(namespace).NotTo(BeNil(), "namespace is required for SetupStaticCredentials")
	Expect(e2eCtx).NotTo(BeNil(), "e2eCtx is required for SetupStaticCredentials")

	id := &infrav1.AWSClusterStaticIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      idName,
			Namespace: namespace.Name,
		},
	}

	Byf("Deleting AWSClusterStaticIdentity %s in namespace %s", idName, namespace.Name)
	client := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
	Eventually(func() error {
		return client.Delete(ctx, id)
	}, e2eCtx.E2EConfig.GetIntervals("", "wait-create-identity")...).Should(Succeed())

	//NOTE: secrets should be cleared up when the namespaces are deleted
}
