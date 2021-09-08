/*
Copyright 2019 The Kubernetes Authors.

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

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/internal/userdata"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func TestEKSConfigReconciler(t *testing.T) {
	t.Run("Should reconcile an EKSConfig and create data Secret", func(t *testing.T) {
		g := NewWithT(t)
		amcp := newAMCP("test-cluster")
		cluster := newCluster(amcp.Name)
		machine := newMachine(cluster, "test-machine")
		config := newEKSConfig(machine)
		t.Log(dump("amcp", amcp))
		t.Log(dump("config", config))
		t.Log(dump("machine", machine))
		t.Log(dump("cluster", cluster))
		expectedUserData, err := newUserData(cluster.Name, map[string]string{"test-arg": "test-value"})
		g.Expect(err).To(BeNil())
		g.Expect(testEnv.Client.Create(ctx, amcp)).To(Succeed())

		reconciler := EKSConfigReconciler{
			Client: testEnv.Client,
		}
		t.Log(fmt.Sprintf("Calling reconcile on cluster '%s' and config '%s' should requeue", cluster.Name, config.Name))
		g.Eventually(func(gomega Gomega) {
			result, err := reconciler.joinWorker(ctx, cluster, config)
			gomega.Expect(err).NotTo(HaveOccurred())
			gomega.Expect(result.Requeue).To(BeFalse())
		}).Should(Succeed())

		t.Log(fmt.Sprintf("Secret '%s' should exist and be correct", config.Name))
		secretList := &corev1.SecretList{}
		testEnv.Client.List(ctx, secretList)
		t.Log(dump("secrets", secretList))
		secret := &corev1.Secret{}
		g.Eventually(func(gomega Gomega) {
			gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{
				Name:      config.Name,
				Namespace: "default",
			}, secret)).To(Succeed())
		}).Should(Succeed())

		g.Expect(string(secret.Data["value"])).To(Equal(string(expectedUserData)))
	})

	t.Run("Should reconcile an EKSConfig and update data Secret", func(t *testing.T) {
		g := NewWithT(t)
		amcp := newAMCP("test-cluster")
		cluster := newCluster(amcp.Name)
		machine := newMachine(cluster, "test-machine")
		config := newEKSConfig(machine)
		t.Log(dump("amcp", amcp))
		t.Log(dump("config", config))
		t.Log(dump("machine", machine))
		t.Log(dump("cluster", cluster))
		oldUserData, err := newUserData(cluster.Name, map[string]string{"test-arg": "test-value"})
		g.Expect(err).To(BeNil())
		expectedUserData, err := newUserData(cluster.Name, map[string]string{"test-arg": "updated-test-value"})
		g.Expect(err).To(BeNil())
		g.Expect(testEnv.Client.Create(ctx, amcp)).To(Succeed())

		amcpList := &ekscontrolplanev1.AWSManagedControlPlaneList{}
		testEnv.Client.List(ctx, amcpList)
		t.Log(dump("stored-amcps", amcpList))

		reconciler := EKSConfigReconciler{
			Client: testEnv.Client,
		}
		t.Log(fmt.Sprintf("Calling reconcile on cluster '%s' and config '%s' should requeue", cluster.Name, config.Name))
		g.Eventually(func(gomega Gomega) {
			result, err := reconciler.joinWorker(ctx, cluster, config)
			gomega.Expect(err).NotTo(HaveOccurred())
			gomega.Expect(result.Requeue).To(BeFalse())
		}).Should(Succeed())

		t.Log(fmt.Sprintf("Secret '%s' should exist and be correct", config.Name))
		secretList := &corev1.SecretList{}
		testEnv.Client.List(ctx, secretList)
		t.Log(dump("secrets", secretList))

		secret := &corev1.Secret{}
		g.Eventually(func(gomega Gomega) {
			gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{
				Name:      config.Name,
				Namespace: "default",
			}, secret)).To(Succeed())
			gomega.Expect(string(secret.Data["value"])).To(Equal(string(oldUserData)))
		}).Should(Succeed())

		// Secret already exists in testEnv so we update it
		config.Spec.KubeletExtraArgs = map[string]string{
			"test-arg": "updated-test-value",
		}
		t.Log(dump("config", config))
		g.Eventually(func(gomega Gomega) {
			result, err := reconciler.joinWorker(ctx, cluster, config)
			gomega.Expect(err).NotTo(HaveOccurred())
			gomega.Expect(result.Requeue).To(BeFalse())
		}).Should(Succeed())

		t.Log(fmt.Sprintf("Secret '%s' should exist and be up to date", config.Name))

		testEnv.Client.List(ctx, secretList)
		t.Log(dump("secrets", secretList))
		g.Eventually(func(gomega Gomega) {
			gomega.Expect(testEnv.Client.Get(ctx, client.ObjectKey{
				Name:      config.Name,
				Namespace: "default",
			}, secret)).To(Succeed())
			gomega.Expect(string(secret.Data["value"])).To(Equal(string(expectedUserData)))
		}).Should(Succeed())
	})
}

// newCluster return a CAPI cluster object.
func newCluster(name string) *clusterv1.Cluster {
	cluster := &clusterv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clusterv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Spec: clusterv1.ClusterSpec{
			ControlPlaneRef: &corev1.ObjectReference{
				Name:      name,
				Kind:      "AWSManagedControlPlane",
				Namespace: "default",
			},
		},
		Status: clusterv1.ClusterStatus{
			InfrastructureReady: true,
		},
	}
	conditions.MarkTrue(cluster, clusterv1.ControlPlaneInitializedCondition)
	return cluster
}

func dump(desc string, o interface{}) string {
	dat, _ := yaml.Marshal(o)
	return fmt.Sprintf("%s:\n%s", desc, string(dat))
}

// newMachine return a CAPI machine object; if cluster is not nil, the machine is linked to the cluster as well.
func newMachine(cluster *clusterv1.Cluster, name string) *clusterv1.Machine {
	generatedName := fmt.Sprintf("%s-%s", name, util.RandomString(5))
	machine := &clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: clusterv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      generatedName,
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				ConfigRef: &corev1.ObjectReference{
					Kind:       "EKSConfig",
					APIVersion: bootstrapv1.GroupVersion.String(),
				},
			},
		},
	}
	if cluster != nil {
		machine.Spec.ClusterName = cluster.Name
		machine.ObjectMeta.Labels = map[string]string{
			clusterv1.ClusterLabelName: cluster.Name,
		}
	}
	return machine
}

// newEKSConfig return an EKSConfig object; if machine is not nil, the EKSConfig is linked to the machine as well.
func newEKSConfig(machine *clusterv1.Machine) *bootstrapv1.EKSConfig {
	config := &bootstrapv1.EKSConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EKSConfig",
			APIVersion: bootstrapv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
		},
		Spec: bootstrapv1.EKSConfigSpec{
			KubeletExtraArgs: map[string]string{
				"test-arg": "test-value",
			},
		},
	}
	if machine != nil {
		config.ObjectMeta.Name = machine.Name
		config.ObjectMeta.UID = types.UID(fmt.Sprintf("%s uid", machine.Name))
		config.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			{
				Kind:       "Machine",
				APIVersion: clusterv1.GroupVersion.String(),
				Name:       machine.Name,
				UID:        types.UID(fmt.Sprintf("%s uid", machine.Name)),
			},
		}
		machine.Spec.Bootstrap.ConfigRef.Name = config.Name
		machine.Spec.Bootstrap.ConfigRef.Namespace = config.Namespace
	}
	return config
}

// newUserData generates user data ready to be passed to r.storeBootstrapData.
func newUserData(clusterName string, kubeletExtraArgs map[string]string) ([]byte, error) {
	return userdata.NewNode(&userdata.NodeInput{
		ClusterName:      clusterName,
		KubeletExtraArgs: kubeletExtraArgs,
	})
}

// newAMCP returns an EKS AWSManagedControlPlane object.
func newAMCP(name string) *ekscontrolplanev1.AWSManagedControlPlane {
	generatedName := fmt.Sprintf("%s-%s", name, util.RandomString(5))
	return &ekscontrolplanev1.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: ekscontrolplanev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      generatedName,
			Namespace: "default",
		},
		Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
			EKSClusterName: generatedName,
		},
	}
}
