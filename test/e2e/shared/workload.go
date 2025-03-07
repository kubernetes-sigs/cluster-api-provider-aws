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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api/test/framework"
)

// WaitForDeploymentsAvailableInput is the input for WaitForDeploymentsAvailable.
type WaitForDeploymentsAvailableInput struct {
	Getter    framework.Getter
	Name      string
	Namespace string
}

// GetDeploymentInput is the input for GetDeployment.
type GetDeploymentInput struct {
	Getter    framework.Getter
	Name      string
	Namespace string
}

// ReconfigureDeploymentInput is the input for ReconfigureDeployment.
type ReconfigureDeploymentInput struct {
	Getter       framework.Getter
	ClientSet    *kubernetes.Clientset
	Name         string
	Namespace    string
	WaitInterval []interface{}
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

// GetDeployment gets the deployment object.
func GetDeployment(ctx context.Context, input GetDeploymentInput) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	key := client.ObjectKey{
		Namespace: input.Namespace,
		Name:      input.Name,
	}
	getErr := input.Getter.Get(ctx, key, deployment)
	if getErr != nil {
		return nil, getErr
	}
	return deployment, nil
}

// ReconfigureDeployment updates the deployment object.
func ReconfigureDeployment(ctx context.Context, input ReconfigureDeploymentInput, updateBy func(dep *appsv1.Deployment) (*appsv1.Deployment, error), validateBy func(dep *appsv1.Deployment) error) {
	By(fmt.Sprintf("Get %s deployment object", input.Name))
	deployment, getErr := GetDeployment(ctx, GetDeploymentInput{
		Getter:    input.Getter,
		Name:      input.Name,
		Namespace: input.Namespace,
	})
	Expect(getErr).To(BeNil())

	By(fmt.Sprintf("Update %s deployment object spec", input.Name))
	_, specErr := updateBy(deployment)
	Expect(specErr).To(BeNil())

	By(fmt.Sprintf("Update %s deployment object", input.Name))
	_, updateErr := input.ClientSet.AppsV1().Deployments(input.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	Expect(updateErr).To(BeNil())

	By(fmt.Sprintf("Wait for %s deployment to be available after reconfiguring", input.Name))
	WaitForDeploymentsAvailable(ctx, WaitForDeploymentsAvailableInput{
		Getter:    input.Getter,
		Name:      input.Name,
		Namespace: input.Namespace,
	}, input.WaitInterval...)

	if validateBy != nil {
		By(fmt.Sprintf("Validate %s deployment updated as expected", input.Name))
		updatedDeployment, err := GetDeployment(ctx, GetDeploymentInput{
			Getter:    input.Getter,
			Name:      input.Name,
			Namespace: input.Namespace,
		})
		Expect(err).To(BeNil())

		vaErr := validateBy(updatedDeployment)
		Expect(vaErr).To(BeNil())
	}
}

// EnableAlternativeGCStrategy enables AlternativeGCStrategy in CAPA controller manager args field.
func EnableAlternativeGCStrategy(dep *appsv1.Deployment) (*appsv1.Deployment, error) {
	for i, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AlternativeGCStrategy") {
			dep.Spec.Template.Spec.Containers[0].Args[i] = strings.Replace(arg, "AlternativeGCStrategy=false", "AlternativeGCStrategy=true", 1)
			return dep, nil
		}
	}
	return nil, fmt.Errorf("fail to find AlternativeGCStrategy to enable")
}

// DisableAlternativeGCStrategy disables AlternativeGCStrategy in CAPA controller manager args field.
func DisableAlternativeGCStrategy(dep *appsv1.Deployment) (*appsv1.Deployment, error) {
	for i, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AlternativeGCStrategy") {
			dep.Spec.Template.Spec.Containers[0].Args[i] = strings.Replace(arg, "AlternativeGCStrategy=true", "AlternativeGCStrategy=false", 1)
			return dep, nil
		}
	}
	return nil, fmt.Errorf("fail to find AlternativeGCStrategy to disable")
}

// ValidateAlternativeGCStrategyEnabled validates AlternativeGCStrategy in CAPA controller manager args field is set to true.
func ValidateAlternativeGCStrategyEnabled(dep *appsv1.Deployment) error {
	for _, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AlternativeGCStrategy=true") {
			return nil
		}
	}
	return fmt.Errorf("fail to validate AlternativeGCStrategy set to true")
}

// ValidateAlternativeGCStrategyDisabled validates AlternativeGCStrategy in CAPA controller manager args field is set to false.
func ValidateAlternativeGCStrategyDisabled(dep *appsv1.Deployment) error {
	for _, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AlternativeGCStrategy=false") {
			return nil
		}
	}
	return fmt.Errorf("fail to validate AlternativeGCStrategy set to false")
}

// ValidateAWSClusterEnabled validates AWSCluster in CAPA controller manager args field is set to true.
func ValidateAWSClusterEnabled(dep *appsv1.Deployment) error {
	for _, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AWSCluster=true") {
			return nil
		}
	}
	return fmt.Errorf("fail to validate AWSCluster set to true")
}

// ValidateAWSMachineEnabled validates AWSMachine in CAPA controller manager args field is set to true.
func ValidateAWSMachineEnabled(dep *appsv1.Deployment) error {
	for _, arg := range dep.Spec.Template.Spec.Containers[0].Args {
		if strings.Contains(arg, "feature-gates") && strings.Contains(arg, "AWSMachine=true") {
			return nil
		}
	}
	return fmt.Errorf("fail to validate AWSMachine set to true")
}
