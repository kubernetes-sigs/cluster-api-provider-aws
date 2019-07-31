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

package awscluster

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrastructurev1alpha2 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/infrastructure/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	RegisterTestingT(t)
	ctx := context.Background()
	instance := &infrastructurev1alpha2.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	Expect(err).NotTo(HaveOccurred())
	c := mgr.GetClient()
	Expect(add(mgr, newReconciler(mgr))).To(Succeed())

	stopMgr, mgrStopped := StartTestManager(mgr)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the AWSCluster object and expect the Reconcile and Deployment to be created
	err = c.Create(ctx, instance)
	Expect(err).NotTo(HaveOccurred())
	defer c.Delete(context.TODO(), instance)

	Eventually(func() bool {
		key := client.ObjectKey{Name: instance.Name, Namespace: instance.Namespace}
		if err := c.Get(ctx, key, instance); err != nil {
			return false
		}
		return true
	}, timeout).Should(BeTrue())

}
