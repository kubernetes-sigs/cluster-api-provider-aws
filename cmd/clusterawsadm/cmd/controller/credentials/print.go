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

// Package credentials provides a CLI utilities for AWS credentials.
package credentials

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/controller"
)

// PrintCredentialsCmd is a CLI command that will print credentials the controller is using.
func PrintCredentialsCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "print-credentials",
		Short: "print credentials the controller is using",
		Long: templates.LongDesc(`
			print credentials the controller is using
		`),
		Example: templates.Examples(`
		# print credentials
		clusterawsadm controller print-credentials --kubeconfig=kubeconfig --namespace=capa-system
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := controller.GetClient(kubeconfigPath, kubeconfigContext)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get client-go client for the cluster: %s\n", err.Error())
				return err
			}

			secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), controller.BootstrapCredsSecret, metav1.GetOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get bootstrap credentials secret: %s\n", err.Error())
				return err
			}

			controller.PrintBootstrapCredentials(secret)
			return nil
		},
	}
	addKubeconfigFlag(newCmd)
	addKubeconfigContextFlag(newCmd)
	addNamespaceFlag(newCmd)
	return newCmd
}
