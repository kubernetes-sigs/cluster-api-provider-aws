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

package instancestate

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func createAWSCluster(name string) *infrav1.AWSCluster {
	return &infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
	}
}

func persistObject(g *WithT, o client.Object) {
	ctx := context.TODO()
	g.Expect(k8sClient.Create(ctx, o)).Should(Succeed())
	metaObj, err := meta.Accessor(o)
	g.Expect(err).NotTo(HaveOccurred())
	lookupKey := types.NamespacedName{Name: metaObj.GetName(), Namespace: metaObj.GetNamespace()}

	g.Eventually(func() bool {
		err := k8sClient.Get(ctx, lookupKey, o)
		return err == nil
	}, time.Second*10).Should(BeTrue())
}

func deleteAWSCluster(g *WithT, name string) {
	ctx := context.TODO()
	awsLookupKey := types.NamespacedName{Name: name, Namespace: "default"}
	awsCluster := &infrav1.AWSCluster{}
	err := k8sClient.Get(ctx, awsLookupKey, awsCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// already deleted
			return
		}
		Fail("Unexpected error when fetching cluster")
	}
	g.Expect(k8sClient.Delete(ctx, awsCluster)).To(Succeed())
}
