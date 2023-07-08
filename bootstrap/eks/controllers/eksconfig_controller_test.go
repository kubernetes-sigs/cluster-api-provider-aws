/*
Copyright 2020 The Kubernetes Authors.

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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
)

func TestEKSConfigReconcilerReturnEarlyIfClusterInfraNotReady(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newEKSConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		InfrastructureReady: false,
	}

	reconciler := EKSConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		result, err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(result).To(Equal(reconcile.Result{}))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
}

func TestEKSConfigReconcilerReturnEarlyIfClusterControlPlaneNotInitialized(t *testing.T) {
	g := NewWithT(t)

	cluster := newCluster("cluster")
	machine := newMachine(cluster, "machine")
	config := newEKSConfig(machine)

	cluster.Status = clusterv1.ClusterStatus{
		InfrastructureReady: true,
	}

	reconciler := EKSConfigReconciler{
		Client: testEnv.Client,
	}

	g.Eventually(func(gomega Gomega) {
		result, err := reconciler.joinWorker(context.Background(), cluster, config, configOwner("Machine"))
		gomega.Expect(result).To(Equal(reconcile.Result{}))
		gomega.Expect(err).NotTo(HaveOccurred())
	}).Should(Succeed())
}

func configOwner(kind string) *bsutil.ConfigOwner {
	unstructuredOwner := unstructured.Unstructured{
		Object: map[string]interface{}{"kind": kind},
	}
	configOwner := bsutil.ConfigOwner{Unstructured: &unstructuredOwner}
	return &configOwner
}
