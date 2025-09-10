/*
Copyright 2025 The Kubernetes Authors.

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
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

func TestNodeadmConfigReconciler_CreateSecret(t *testing.T) {
	g := NewWithT(t)

	amcp := newAMCP("test-cluster")
	endpoint := clusterv1.APIEndpoint{Host: "https://9.9.9.9", Port: 6443}
	amcp.Spec.ControlPlaneEndpoint = endpoint
	cluster := newCluster(amcp.Name)
	newStatus := cluster.Status
	amcpStatus := amcp.Status
	g.Expect(testEnv.Client.Create(ctx, amcp)).To(Succeed())
	g.Expect(testEnv.Client.Create(ctx, cluster)).To(Succeed())
	cluster.Status = newStatus
	g.Expect(testEnv.Client.Status().Update(ctx, cluster)).To(Succeed())
	amcp.Status = amcpStatus
	g.Expect(testEnv.Client.Status().Update(ctx, amcp)).To(Succeed())
	kubeconfigSecret := newKubeconfigSecret("https://9.9.9.9:6443", cluster)
	g.Expect(testEnv.Client.Create(ctx, kubeconfigSecret)).To(Succeed())

	machine := newMachine(cluster, "test-machine")
	cfg := newNodeadmConfig(machine)
	g.Expect(testEnv.Client.Create(ctx, cfg)).To(Succeed())

	reconciler := NodeadmConfigReconciler{Client: testEnv.Client}

	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(ctx, cluster, cfg, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}, time.Second*15, time.Second*5).Should(Succeed())

	secret := &corev1.Secret{}
	g.Eventually(func(gomega Gomega) {
		gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{Name: cfg.Name, Namespace: "default"}, secret)).To(Succeed())
	}, time.Minute, time.Second*5).Should(Succeed())

	g.Expect(string(secret.Data["value"])).To(ContainSubstring("apiVersion: node.eks.aws/v1alpha1"))
	g.Expect(string(secret.Data["value"])).To(ContainSubstring("apiServerEndpoint: https://9.9.9.9"))
}

func TestNodeadmConfigReconciler_UpdateSecret_ForMachinePool(t *testing.T) {
	g := NewWithT(t)

	amcp := newAMCP("test-cluster")
	endpoint := clusterv1.APIEndpoint{Host: "https://9.9.9.9", Port: 6443}
	amcp.Spec.ControlPlaneEndpoint = endpoint
	cluster := newCluster(amcp.Name)
	newStatus := cluster.Status
	amcpStatus := amcp.Status
	g.Expect(testEnv.Client.Create(ctx, amcp)).To(Succeed())
	g.Expect(testEnv.Client.Create(ctx, cluster)).To(Succeed())
	cluster.Status = newStatus
	g.Expect(testEnv.Client.Status().Update(ctx, cluster)).To(Succeed())
	amcp.Status = amcpStatus
	g.Expect(testEnv.Client.Status().Update(ctx, amcp)).To(Succeed())
	kubeconfigSecret := newKubeconfigSecret("https://9.9.9.9:6443", cluster)
	g.Expect(testEnv.Client.Create(ctx, kubeconfigSecret)).To(Succeed())

	mp := newMachinePool(cluster, "test-mp")
	cfg := newNodeadmConfig(nil)
	cfg.ObjectMeta.Name = mp.Name
	cfg.ObjectMeta.UID = types.UID(fmt.Sprintf("%s uid", mp.Name))
	cfg.ObjectMeta.OwnerReferences = []metav1.OwnerReference{{
		Kind:       "MachinePool",
		APIVersion: expclusterv1.GroupVersion.String(),
		Name:       mp.Name,
		UID:        types.UID(fmt.Sprintf("%s uid", mp.Name)),
	}}
	cfg.Status.DataSecretName = &mp.Name
	// initial kubelet flags
	cfg.Spec.Kubelet = &eksbootstrapv1.KubeletOptions{Flags: []string{"--register-with-taints=dedicated=infra:NoSchedule"}}

	reconciler := NodeadmConfigReconciler{Client: testEnv.Client}

	// first reconcile creates secret
	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(ctx, cluster, cfg, configOwner("MachinePool"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}, time.Minute, time.Second*5).Should(Succeed())

	secret := &corev1.Secret{}
	g.Eventually(func(gomega Gomega) {
		gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{Name: cfg.Name, Namespace: "default"}, secret)).To(Succeed())
	}, time.Minute, time.Second*5).Should(Succeed())
	oldData := append([]byte(nil), secret.Data["value"]...)

	// change flags to force different userdata
	cfg.Spec.Kubelet.Flags = []string{"--register-with-taints=dedicated=db:NoSchedule"}

	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(ctx, cluster, cfg, configOwner("MachinePool"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}, time.Minute, time.Second*5).Should(Succeed())

	g.Eventually(func(gomega Gomega) {
		gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{Name: cfg.Name, Namespace: "default"}, secret)).To(Succeed())
		gomega.Expect(secret.Data["value"]).NotTo(Equal(oldData))
	}, time.Minute, time.Second*5).Should(Succeed())
}

func TestNodeadmConfigReconciler_ResolvesSecretFileReference(t *testing.T) {
	g := NewWithT(t)

	amcp := newAMCP("test-cluster")
	endpoint := clusterv1.APIEndpoint{Host: "https://9.9.9.9", Port: 6443}
	amcp.Spec.ControlPlaneEndpoint = endpoint
	cluster := newCluster(amcp.Name)
	newStatus := cluster.Status
	amcpStatus := amcp.Status
	g.Expect(testEnv.Client.Create(ctx, amcp)).To(Succeed())
	g.Expect(testEnv.Client.Create(ctx, cluster)).To(Succeed())
	cluster.Status = newStatus
	g.Expect(testEnv.Client.Status().Update(ctx, cluster)).To(Succeed())
	amcp.Status = amcpStatus
	g.Expect(testEnv.Client.Status().Update(ctx, amcp)).To(Succeed())
	kubeconfigSecret := newKubeconfigSecret("https://9.9.9.9:6443", cluster)
	g.Expect(testEnv.Client.Create(ctx, kubeconfigSecret)).To(Succeed())

	//nolint:gosec // test constant
	secretPath := "/etc/secret.txt"
	secretContent := "secretValue"
	machine := newMachine(cluster, "test-machine")
	cfg := newNodeadmConfig(machine)
	cfg.Spec.Files = append(cfg.Spec.Files, eksbootstrapv1.File{
		ContentFrom: &eksbootstrapv1.FileSource{Secret: eksbootstrapv1.SecretFileSource{Name: "my-secret2", Key: "secretKey"}},
		Path:        secretPath,
	})
	// ensure cloud-config part is rendered
	cfg.Spec.NTP = &eksbootstrapv1.NTP{Enabled: func() *bool { b := true; return &b }()}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "my-secret2"},
		Data:       map[string][]byte{"secretKey": []byte(secretContent)},
	}

	g.Expect(testEnv.Client.Create(ctx, secret)).To(Succeed())

	expectedContains := []string{
		"#cloud-config",
		secretContent,
	}

	reconciler := NodeadmConfigReconciler{Client: testEnv.Client}
	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(ctx, cluster, cfg, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}, time.Minute, time.Second*5).Should(Succeed())

	got := &corev1.Secret{}
	g.Eventually(func(gomega Gomega) {
		gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{Name: cfg.Name, Namespace: "default"}, got)).To(Succeed())
	}, time.Minute, time.Second*5).Should(Succeed())

	for _, s := range expectedContains {
		g.Expect(string(got.Data["value"])).To(ContainSubstring(s), "userdata should contain %q", s)
	}
}
