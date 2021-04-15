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

package controlleridentitycreator

import (
	"context"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("AWSInstanceStateController", func() {
	It("should maintain list of cluster queue URLs and reconcile failing machines", func() {
		ctx := context.Background()

		instance := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
		instance.Default()

		// Create the AWSCluster object and expect the Reconcile and Deployment to be created
		Expect(testEnv.Create(ctx, instance)).To(Succeed())

		By("Ensuring AWSClusterControllerIdentity instance is created")
		Eventually(func() bool {
			cp := &infrav1.AWSClusterControllerIdentity{}
			key := client.ObjectKey{
				Name: infrav1.AWSClusterControllerIdentityName,
			}
			err := testEnv.Get(ctx, key, cp)
			if err != nil {
				return false
			}
			if reflect.DeepEqual(*cp.Spec.AllowedNamespaces, infrav1.AllowedNamespaces{}) {
				return true
			}
			return false
		}, 10*time.Second).Should(Equal(true))

	})
})
