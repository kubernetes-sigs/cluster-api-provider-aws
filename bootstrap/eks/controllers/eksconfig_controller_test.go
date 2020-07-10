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
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ = Describe("EKSConfigReconciler", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Reconcile a EKSConfig", func() {
		It("should wait until infrastructure is ready", func() {
			cluster := newCluster("cluster1")
			Expect(testEnv.Create(context.Background(), cluster)).To(Succeed())

			machine := newMachine(cluster, "my-machine")
			Expect(testEnv.Create(context.Background(), machine)).To(Succeed())

			config := newEKSConfig(machine, "my-machine-config")
			Expect(testEnv.Create(context.Background(), config)).To(Succeed())

			reconciler := EKSConfigReconciler{
				Log:    log.Log,
				Client: testEnv,
			}
			By("Calling reconcile should requeue")
			result, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: "default",
					Name:      "my-machine-config",
				},
			})
			Expect(err).To(Succeed())
			Expect(result.Requeue).To(BeFalse())
		})
	})
})

// getEKSConfig returns a EKSConfig object from the cluster
func getEKSConfig(c client.Client, name string) (*bootstrapv1.EKSConfig, error) {
	ctx := context.Background()
	controlplaneConfigKey := client.ObjectKey{
		Namespace: "default",
		Name:      name,
	}
	config := &bootstrapv1.EKSConfig{}
	err := c.Get(ctx, controlplaneConfigKey, config)
	return config, err
}

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
	}
}

// newEKSConfig return a CABPK EKSConfig object; if machine is not nil, the EKSConfig is linked to the machine as well
func newEKSConfig(machine *clusterv1.Machine, name string) *bootstrapv1.EKSConfig {
	config := &bootstrapv1.EKSConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EKSConfig",
			APIVersion: bootstrapv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
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
