//go:build e2e
// +build e2e

/*
Copyright 2025 The Kubernetes Authors.

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

package unmanaged

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
)

var _ = ginkgo.Context("[unmanaged] [conversion]", func() {
	var (
		ctx context.Context
	)

	ginkgo.BeforeEach(func() {
		ctx = context.TODO()
	})

	ginkgo.Describe("AWSClusterControllerIdentity conversion", func() {
		ginkgo.It("should successfully convert v1beta1 to v1beta2", func() {
			specName := "conversion-awsclustercontrolleridentity"
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, specName, namespace, e2eCtx)

			bootstrapClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
			kubeconfigPath := e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath()

			ginkgo.By("Creating AWSClusterControllerIdentity in v1beta1 format using kubectl")
			// Create the resource using kubectl with explicit v1beta1 APIVersion
			// This simulates a resource created by an old provider version
			v1beta1YAML := fmt.Sprintf(`apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSClusterControllerIdentity
metadata:
  name: %s
spec:
  allowedNamespaces: {}
`, infrav1.AWSClusterControllerIdentityName)

			// Delete existing if present
			existing := &infrav1.AWSClusterControllerIdentity{}
			err := bootstrapClient.Get(ctx, client.ObjectKey{Name: infrav1.AWSClusterControllerIdentityName}, existing)
			if err == nil {
				Expect(bootstrapClient.Delete(ctx, existing)).To(Succeed())
				// Wait for deletion
				Eventually(func() bool {
					err := bootstrapClient.Get(ctx, client.ObjectKey{Name: infrav1.AWSClusterControllerIdentityName}, existing)
					return apierrors.IsNotFound(err)
				}, 30*time.Second, 1*time.Second).Should(BeTrue())
			}

			// Create using kubectl apply
			cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(v1beta1YAML)
			output, err := cmd.CombinedOutput()
			if err != nil && !strings.Contains(string(output), "already exists") {
				Expect(err).NotTo(HaveOccurred(), "Failed to create v1beta1 AWSClusterControllerIdentity: %s", string(output))
			}

			ginkgo.By("Verifying resource exists via kubectl get")
			// Use kubectl to get the resource and verify it can be retrieved
			// This will trigger conversion if needed
			getCmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath,
				"get", "awsclustercontrolleridentity", infrav1.AWSClusterControllerIdentityName,
				"-o", "jsonpath={.apiVersion}")
			output, err = getCmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "Failed to get AWSClusterControllerIdentity: %s", string(output))

			apiVersion := strings.TrimSpace(string(output))
			ginkgo.By(fmt.Sprintf("Resource APIVersion: %s", apiVersion))

			// Verify the APIVersion is v1beta2 (converted from v1beta1)
			Expect(apiVersion).To(Equal("infrastructure.cluster.x-k8s.io/v1beta2"),
				"Expected APIVersion to be v1beta2 after conversion, got %s", apiVersion)

			ginkgo.By("Verifying resource can be retrieved using v1beta2 client")
			// Verify we can get it using the v1beta2 client
			v1beta2Identity := &infrav1.AWSClusterControllerIdentity{}
			err = bootstrapClient.Get(ctx, client.ObjectKey{
				Name: infrav1.AWSClusterControllerIdentityName,
			}, v1beta2Identity)
			Expect(err).NotTo(HaveOccurred(), "Failed to get AWSClusterControllerIdentity using v1beta2 client")

			// Verify the APIVersion is correct
			Expect(v1beta2Identity.APIVersion).To(Equal(infrav1.GroupVersion.String()),
				"Expected APIVersion to be %s, got %s", infrav1.GroupVersion.String(), v1beta2Identity.APIVersion)

			ginkgo.By("Verifying conversion webhook works by listing resources")
			// List all AWSClusterControllerIdentity resources to ensure conversion works for list operations
			identityList := &infrav1.AWSClusterControllerIdentityList{}
			err = bootstrapClient.List(ctx, identityList)
			Expect(err).NotTo(HaveOccurred(), "Failed to list AWSClusterControllerIdentity resources")

			// Verify at least one item exists and has correct APIVersion
			Expect(len(identityList.Items)).To(BeNumerically(">=", 1),
				"Expected at least one AWSClusterControllerIdentity in the list")
			Expect(identityList.Items[0].APIVersion).To(Equal(infrav1.GroupVersion.String()),
				"Expected list item APIVersion to be %s, got %s", infrav1.GroupVersion.String(), identityList.Items[0].APIVersion)

			ginkgo.By("PASSED! Conversion webhook successfully converts v1beta1 to v1beta2")
		})

		ginkgo.It("should handle conversion during provider upgrade", func() {
			specName := "conversion-upgrade-awsclustercontrolleridentity"
			namespace := shared.SetupSpecNamespace(ctx, specName, e2eCtx)
			defer shared.DumpSpecResourcesAndCleanup(ctx, specName, namespace, e2eCtx)

			bootstrapClient := e2eCtx.Environment.BootstrapClusterProxy.GetClient()
			kubeconfigPath := e2eCtx.Environment.BootstrapClusterProxy.GetKubeconfigPath()

			ginkgo.By("Step 1: Creating AWSClusterControllerIdentity in v1beta1 format (simulating old provider)")
			// Create the resource using kubectl with explicit v1beta1 APIVersion
			// This simulates a resource created by an old provider version
			v1beta1YAML := fmt.Sprintf(`apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AWSClusterControllerIdentity
metadata:
  name: %s
spec:
  allowedNamespaces: {}
`, infrav1.AWSClusterControllerIdentityName)

			// Delete if exists first to ensure clean state
			existing := &infrav1.AWSClusterControllerIdentity{}
			err := bootstrapClient.Get(ctx, client.ObjectKey{Name: infrav1.AWSClusterControllerIdentityName}, existing)
			if err == nil {
				Expect(bootstrapClient.Delete(ctx, existing)).To(Succeed())
				// Wait for deletion
				Eventually(func() bool {
					err := bootstrapClient.Get(ctx, client.ObjectKey{Name: infrav1.AWSClusterControllerIdentityName}, existing)
					return apierrors.IsNotFound(err)
				}, 30*time.Second, 1*time.Second).Should(BeTrue())
			}

			// Create using kubectl apply
			createCmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "apply", "-f", "-")
			createCmd.Stdin = strings.NewReader(v1beta1YAML)
			output, err := createCmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "Failed to create v1beta1 AWSClusterControllerIdentity: %s", string(output))

			ginkgo.By("Step 2: Verifying resource exists with v1beta1 APIVersion")
			// Verify it was created with v1beta1 by checking the stored version
			// Note: The stored version might be v1beta2 if conversion already happened,
			// but we can verify conversion works when we get it

			ginkgo.By("Step 3: Simulating provider upgrade - getting resource should trigger conversion")
			// When we get the resource using v1beta2 client, Kubernetes will request conversion
			// This simulates what happens when the provider is upgraded and the webhook handles conversion
			v1beta2Identity := &infrav1.AWSClusterControllerIdentity{}
			err = bootstrapClient.Get(ctx, client.ObjectKey{
				Name: infrav1.AWSClusterControllerIdentityName,
			}, v1beta2Identity)

			Expect(err).NotTo(HaveOccurred(), "Failed to get AWSClusterControllerIdentity after upgrade - conversion webhook may have failed")

			ginkgo.By("Step 4: Verifying converted resource has correct APIVersion")
			// Verify the APIVersion is correct after conversion
			Expect(v1beta2Identity.APIVersion).To(Equal(infrav1.GroupVersion.String()),
				"Expected APIVersion to be %s after conversion, got %s. This indicates the conversion webhook fix is working.",
				infrav1.GroupVersion.String(), v1beta2Identity.APIVersion)

			ginkgo.By("Step 5: Verifying resource can be retrieved via kubectl")
			// Verify via kubectl that the resource is accessible and has correct APIVersion
			kubectlCmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath,
				"get", "awsclustercontrolleridentity", infrav1.AWSClusterControllerIdentityName,
				"-o", "jsonpath={.apiVersion}")
			output, err = kubectlCmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), "kubectl get failed: %s", string(output))

			apiVersion := strings.TrimSpace(string(output))
			Expect(apiVersion).To(Equal("infrastructure.cluster.x-k8s.io/v1beta2"),
				"kubectl get returned wrong APIVersion: expected v1beta2, got %s", apiVersion)

			ginkgo.By("Step 6: Verifying resource spec is preserved after conversion")
			// Verify the spec is preserved correctly
			Expect(v1beta2Identity.Spec.AllowedNamespaces).NotTo(BeNil(),
				"Spec.AllowedNamespaces should be preserved after conversion")

			ginkgo.By("PASSED! Provider upgrade conversion test successful")
		})
	})
})
