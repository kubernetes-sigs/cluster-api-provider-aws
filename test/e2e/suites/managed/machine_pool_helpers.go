// +build e2e

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

package managed

import (
	"context"

	. "github.com/onsi/gomega"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	clusterv1exp "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

type deleteMachinePoolInput struct {
	MachinePool *clusterv1exp.MachinePool
	Deleter     framework.Deleter
}

func deleteMachinePool(ctx context.Context, input deleteMachinePoolInput) {
	shared.Byf("Deleting machine pool %s", input.MachinePool.Name)
	Expect(input.Deleter.Delete(ctx, input.MachinePool)).To(Succeed())
}

type waitForMachinePoolDeletedInput struct {
	MachinePool *clusterv1exp.MachinePool
	Getter      framework.Getter
}

func waitForMachinePoolDeleted(ctx context.Context, input waitForMachinePoolDeletedInput, intervals ...interface{}) {
	shared.Byf("Waiting for machine pool %s to be deleted", input.MachinePool.GetName())
	Eventually(func() bool {
		mp := &clusterv1exp.MachinePool{}
		key := client.ObjectKey{
			Namespace: input.MachinePool.GetNamespace(),
			Name:      input.MachinePool.GetName(),
		}
		err := input.Getter.Get(ctx, key, mp)
		notFound := apierrors.IsNotFound(err)
		return notFound
	}, intervals...).Should(BeTrue())
}
