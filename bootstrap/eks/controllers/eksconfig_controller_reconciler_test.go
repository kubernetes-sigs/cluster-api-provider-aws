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
	"testing"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/yaml"
)

func TestEKSConfigReconciler(t *testing.T) {
	t.Run("Should reconcile an EKSConfig and wait until infrastructure is ready", func(t *testing.T) {
		g := NewWithT(t)
		cluster := newCluster("test-cluster")
		machine := newMachine(cluster, "test-machine")
		config := newEKSConfig(machine, "test-config")

		bytes, err := yaml.Marshal(machine)
		g.Expect(err).To(BeNil())

		owner := &unstructured.Unstructured{}
		err = yaml.Unmarshal(bytes, owner)
		g.Expect(err).To(BeNil())

		reconciler := EKSConfigReconciler{
			Client: testEnv.Client,
		}

		t.Log("Calling reconcile should requeue")
		result, err := reconciler.joinWorker(context.Background(), cluster, config)
		g.Expect(err).To(Succeed())
		g.Expect(result.Requeue).To(BeFalse())
	})
}

// newCluster return a CAPI cluster object.
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

// newMachine return a CAPI machine object; if cluster is not nil, the machine is linked to the cluster as well.
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

// newEKSConfig return an EKSConfig object; if machine is not nil, the EKSConfig is linked to the machine as well.
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
