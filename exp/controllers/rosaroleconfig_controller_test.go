/*
Copyright The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestROSARoleConfigReconciler_Reconcile(t *testing.T) {
	g := NewWithT(t)

	ctx := context.TODO()

	rosaRoleConfig := &expinfrav1.ROSARoleConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-rosa-role",
			Namespace: "test-namespace"},
		Spec: expinfrav1.ROSARoleConfigSpec{},
	}

	// Setup the reconciler with these mocks
	reconciler := &ROSARoleConfigReconciler{
		Client: testEnv.Client,
	}

	// Call the Reconcile function
	req := ctrl.Request{}
	req.NamespacedName = types.NamespacedName{Name: rosaRoleConfig.Name, Namespace: rosaRoleConfig.Namespace}
	_, errReconcile := reconciler.Reconcile(ctx, req)

	// Assertions
	g.Expect(errReconcile).ToNot(HaveOccurred())
}
