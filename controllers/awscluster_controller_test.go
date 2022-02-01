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
	"testing"

	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAWSClusterReconciler(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	reconciler := &AWSClusterReconciler{
		Client: testEnv.Client,
	}

	instance := &infrav1.AWSCluster{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}
	instance.Default()

	// Create the AWSCluster object and expect the Reconcile and Deployment to be created
	g.Expect(testEnv.Create(ctx, instance)).To(Succeed())

	// Calling reconcile should not error and not requeue the request with insufficient set up
	result, err := reconciler.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		},
	})
	g.Expect(err).To(BeNil())
	g.Expect(result.RequeueAfter).To(BeZero())
}
