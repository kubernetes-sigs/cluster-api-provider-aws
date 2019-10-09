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
	"bytes"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func getMachineScope() *scope.MachineScope {
	ms, err := scope.NewMachineScope(
		scope.MachineScopeParams{
			Client:     fake.NewFakeClient(),
			Cluster:    &clusterv1.Cluster{},
			Machine:    &clusterv1.Machine{},
			AWSCluster: &infrav1.AWSCluster{},
			AWSMachine: &infrav1.AWSMachine{},
		},
	)
	Expect(err).To(BeNil())
	return ms
}

func getClusterScope() *scope.ClusterScope {
	cs, err := scope.NewClusterScope(
		scope.ClusterScopeParams{
			Cluster:    &clusterv1.Cluster{},
			AWSCluster: &infrav1.AWSCluster{},
		},
	)

	Expect(err).To(BeNil())
	return cs
}

var _ = Describe("AWSMachineReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile an AWSMachine", func() {
		It("should not error with minimal set up", func() {
			reconciler := &AWSMachineReconciler{
				Client: k8sClient,
				Log:    log.Log,
			}
			By("Calling reconcile")
			instance := &infrav1.AWSMachine{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
			result, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				},
			})
			Expect(err).To(BeNil())
			Expect(result.RequeueAfter).To(BeZero())
		})
	})

	Context("reconcileNormal", func() {

		It("should exit immediately on an error state", func() {
			reconciler := AWSMachineReconciler{}
			er := errors.CreateMachineError
			em := "Couldn't create machine"

			ms := getMachineScope()
			ms.AWSMachine.Status.ErrorReason = &er
			ms.AWSMachine.Status.ErrorMessage = &em

			buf := new(bytes.Buffer)
			klog.LogToOutput(buf)

			_, err := reconciler.reconcileNormal(context.Background(), ms, nil)
			Expect(err).To(BeNil())
			Expect(buf).To(ContainSubstring("Error state detected, skipping reconciliation"))
		})

		It("should add our finalizer to the machine", func() {
			reconciler := AWSMachineReconciler{}
			ms := getMachineScope()
			_, _ = reconciler.reconcileNormal(context.Background(), ms, getClusterScope())

			Expect(ms.AWSMachine.Finalizers).To(ContainElement(infrav1.MachineFinalizer))
		})

		It("should exit immediately if cluster infra isn't ready", func() {
			reconciler := AWSMachineReconciler{}

			ms := getMachineScope()
			ms.Cluster.Status.InfrastructureReady = false

			buf := new(bytes.Buffer)
			klog.LogToOutput(buf)

			_, err := reconciler.reconcileNormal(context.Background(), ms, nil)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("Cluster infrastructure is not ready yet"))
		})

		It("should exit immediately if cluster infra isn't ready", func() {
			reconciler := AWSMachineReconciler{}

			ms := getMachineScope()
			ms.Cluster.Status.InfrastructureReady = true
			ms.Machine.Spec.Bootstrap.Data = nil

			buf := new(bytes.Buffer)
			klog.LogToOutput(buf)

			_, err := reconciler.reconcileNormal(context.Background(), ms, nil)
			Expect(err).To(BeNil())
			Expect(buf.String()).To(ContainSubstring("Bootstrap data is not yet available"))
		})
	})
})
