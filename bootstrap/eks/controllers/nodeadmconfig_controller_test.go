package controllers

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestNodeadmConfigReconcilerReturnEarlyIfClusterInfraNotReady(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newNodeadmConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		InfrastructureReady: false,
	}

	reconciler := NodeadmConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
}

func TestNodeadmConfigReconcilerReturnEarlyIfClusterControlPlaneNotInitialized(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newNodeadmConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		InfrastructureReady: true,
	}

	reconciler := NodeadmConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		_, err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
}

func newNodeadmConfig(machine *clusterv1.Machine) *eksbootstrapv1.NodeadmConfig {
	config := &eksbootstrapv1.NodeadmConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeadmConfig",
			APIVersion: eksbootstrapv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
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
		config.Status.DataSecretName = &machine.Name
		machine.Spec.Bootstrap.ConfigRef.Name = config.Name
		machine.Spec.Bootstrap.ConfigRef.Namespace = config.Namespace
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
		config.Status.DataSecretName = &machine.Name
		machine.Spec.Bootstrap.ConfigRef.Name = config.Name
		machine.Spec.Bootstrap.ConfigRef.Namespace = config.Namespace
	}
	return config
}
