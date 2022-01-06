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
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
	controlplanev1alpha3 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/yaml"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("EKSConfigReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile an EKSConfig", func() {
		It("should wait until infrastructure is ready", func() {
			cluster := newCluster("test-cluster")
			machine := newMachine(cluster, "test-machine")
			config := newEKSConfig(machine, "test-config")

			bytes, err := yaml.Marshal(machine)
			Expect(err).To(BeNil())

			owner := &unstructured.Unstructured{}
			err = yaml.Unmarshal(bytes, owner)
			Expect(err).To(BeNil())

			reconciler := EKSConfigReconciler{
				Log:    log.Log,
				Client: k8sClient,
			}

			By("Calling reconcile should requeue")
			result, err := reconciler.joinWorker(context.Background(), log.Log, cluster, config)
			Expect(err).To(Succeed())
			Expect(result.Requeue).To(BeFalse())
		})
		It("should update an existing secret", func() {
			cluster := newCluster("cluster")
			machine := newMachine(cluster, "machine")
			config := newEKSConfig(machine, "cfg")
			config.Status.DataSecretName = pointer.StringPtr("cfg")

			cluster.Status = clusterv1.ClusterStatus{
				InfrastructureReady:     true,
				ControlPlaneInitialized: true,
			}

			secret := newSecret("cfg")
			err := k8sClient.Create(context.Background(), secret)
			Expect(err).NotTo(HaveOccurred())

			awsmcp := newAWSManagedCotrolPlane("cluster")
			err = k8sClient.Create(context.Background(), awsmcp)
			Expect(err).NotTo(HaveOccurred())

			reconciler := EKSConfigReconciler{
				Log:    log.Log,
				Client: k8sClient,
			}

			result, err := reconciler.joinWorker(context.Background(), log.Log, cluster, config)

			// try to get the new EKSConfig
			createdEKSConfig := &bootstrapv1.EKSConfig{}
			err = k8sClient.Get(context.Background(), types.NamespacedName{Namespace: config.Namespace, Name: config.Name}, createdEKSConfig)
			fmt.Println(createdEKSConfig.Status)

			updateSecret := &corev1.Secret{}
			err = k8sClient.Get(context.Background(), types.NamespacedName{Namespace: secret.Namespace, Name: secret.Name}, updateSecret)
			Expect(string(updateSecret.Data["value"])).NotTo(Equal("fake-data"))
			Expect(result).To(Equal(reconcile.Result{}))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

// newCluster return a CAPI cluster object
func newCluster(name string) *clusterv1.Cluster {
	return &clusterv1.Cluster{
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
	}
}

// newCluster return a CAPI cluster object
func newAWSManagedCotrolPlane(name string) *controlplanev1alpha3.AWSManagedControlPlane {
	return &controlplanev1alpha3.AWSManagedControlPlane{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSManagedControlPlane",
			APIVersion: controlplanev1alpha3.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
	}
}

// newMachine return a CAPI machine object; if cluster is not nil, the machine is linked to the cluster as well
func newMachine(cluster *clusterv1.Cluster, name string) *clusterv1.Machine {
	machine := &clusterv1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: clusterv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
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

// newEKSConfig return an EKSConfig object; if machine is not nil, the EKSConfig is linked to the machine as well
func newEKSConfig(machine *clusterv1.Machine, name string) *bootstrapv1.EKSConfig {
	config := &bootstrapv1.EKSConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EKSConfig",
			APIVersion: bootstrapv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
			UID:       types.UID(fmt.Sprintf("%s uid", name)),
		},
	}
	if machine != nil {
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

// newSecret return an EKSConfig object; if machine is not nil, the EKSConfig is linked to the machine as well
func newSecret(name string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Data: map[string][]byte{
			"value": []byte("fake-data"),
		},
	}
	return secret
}
