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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-aws/test/e2e/shared"
)

type deleteMachineDeploymentInput struct {
	MachineDeployment *clusterv1.MachineDeployment
	Deleter           framework.Deleter
}

func deleteMachineDeployment(ctx context.Context, input deleteMachineDeploymentInput) {
	shared.Byf("Deleting machine deployment %s", input.MachineDeployment.Name)
	Expect(input.Deleter.Delete(ctx, input.MachineDeployment)).To(Succeed())
}

type waitForMachineDeploymentDeletedInput struct {
	MachineDeployment *clusterv1.MachineDeployment
	Getter            framework.Getter
}

func waitForMachineDeploymentDeleted(ctx context.Context, input waitForMachineDeploymentDeletedInput, intervals ...interface{}) {
	shared.Byf("Waiting for machine deployment %s to be deleted", input.MachineDeployment.GetName())
	Eventually(func() bool {
		mp := &clusterv1.MachineDeployment{}
		key := client.ObjectKey{
			Namespace: input.MachineDeployment.GetNamespace(),
			Name:      input.MachineDeployment.GetName(),
		}
		err := input.Getter.Get(ctx, key, mp)
		notFound := apierrors.IsNotFound(err)
		return notFound
	}, intervals...).Should(BeTrue())
}

type deleteMachineInput struct {
	Machine *clusterv1.Machine
	Deleter framework.Deleter
}

func deleteMachine(ctx context.Context, input deleteMachineInput) {
	shared.Byf("Deleting machine %s", input.Machine.Name)
	Expect(input.Deleter.Delete(ctx, input.Machine)).To(Succeed())
}

type waitForMachineDeletedInput struct {
	Machine *clusterv1.Machine
	Getter  framework.Getter
}

func waitForMachineDeleted(ctx context.Context, input waitForMachineDeletedInput, intervals ...interface{}) {
	shared.Byf("Waiting for machine %s to be deleted", input.Machine.GetName())
	Eventually(func() bool {
		mp := &clusterv1.Machine{}
		key := client.ObjectKey{
			Namespace: input.Machine.GetNamespace(),
			Name:      input.Machine.GetName(),
		}
		err := input.Getter.Get(ctx, key, mp)
		notFound := apierrors.IsNotFound(err)
		return notFound
	}, intervals...).Should(BeTrue())
}
