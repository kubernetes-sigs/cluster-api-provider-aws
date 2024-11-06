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

package credentials

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/util"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/controller/credentials"
)

// UpdateCredentialsCmd is a CLI command that will update credentials the controller is using.
func UpdateCredentialsCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "update-credentials",
		Short: "update credentials the controller is using (i.e., update controller bootstrap secret)",
		Long: templates.LongDesc(`
			Update credentials the controller is started with
		`),
		Example: templates.Examples(`
		# update credentials: AWS_B64ENCODED_CREDENTIALS environment variable must be set and be used to update the bootstrap secret
		# Kubeconfig file will be searched in default locations
		clusterawsadm controller update-credentials --namespace=capa-system
		# Provided kubeconfig file will be used
		clusterawsadm controller update-credentials --kubeconfig=kubeconfig  --namespace=capa-system
		# Kubeconfig in the default location will be retrieved and the provided context will be used
		clusterawsadm controller update-credentials --kubeconfig-context=mgmt-cluster  --namespace=capa-system
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			encodedCreds, err := util.GetEnv("AWS_B64ENCODED_CREDENTIALS")
			if err != nil {
				return err
			}

			return credentials.UpdateCredentials(credentials.UpdateCredentialsInput{
				KubeconfigPath:    kubeconfigPath,
				KubeconfigContext: kubeconfigContext,
				Credentials:       encodedCreds,
				Namespace:         namespace,
			})
		},
	}
	addKubeconfigFlag(newCmd)
	addKubeconfigContextFlag(newCmd)
	addNamespaceFlag(newCmd)
	return newCmd
}
