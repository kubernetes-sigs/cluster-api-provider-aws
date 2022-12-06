//go:build e2e
// +build e2e

/*
Copyright 2022 The Kubernetes Authors.

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

package shared

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/test/framework"
)

// WaitForDeploymentsAvailableInput is the input for WaitForDeploymentsAvailable.
type WaitForDeploymentsAvailableInput struct {
	Getter    framework.Getter
	Name      string
	Namespace string
}

// WaitForDeploymentsAvailable will wait for the specified intervel for a Deployment to have status.Available = True.
// NOTE: this is based on a version from the Cluster API test framework.
func WaitForDeploymentsAvailable(ctx context.Context, input WaitForDeploymentsAvailableInput, intervals ...interface{}) {
	By(fmt.Sprintf("Waiting for deployment %s/%s to be available", input.Namespace, input.Name))
	deployment := &appsv1.Deployment{}
	Eventually(func() bool {
		key := client.ObjectKey{
			Namespace: input.Namespace,
			Name:      input.Name,
		}
		if err := input.Getter.Get(ctx, key, deployment); err != nil {
			return false
		}
		for _, c := range deployment.Status.Conditions {
			if c.Type == appsv1.DeploymentAvailable && c.Status == corev1.ConditionTrue {
				return true
			}
		}
		return false
	}, intervals...).Should(BeTrue(), func() string { return DescribeFailedDeployment(input, deployment) })
}

// DescribeFailedDeployment returns detailed output to help debug a deployment failure in e2e.
func DescribeFailedDeployment(input WaitForDeploymentsAvailableInput, deployment *appsv1.Deployment) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Deployment %s/%s failed to get status.Available = True condition",
		input.Namespace, input.Name))
	if deployment == nil {
		b.WriteString("\nDeployment: nil\n")
	} else {
		b.WriteString(fmt.Sprintf("\nDeployment:\n%s\n", framework.PrettyPrint(deployment)))
	}
	return b.String()
}
