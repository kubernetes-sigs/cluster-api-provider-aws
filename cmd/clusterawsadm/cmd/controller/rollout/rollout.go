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

package rollout

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/controller/rollout"
)

var (
	namespace         string
	kubeconfigPath    string
	kubeconfigContext string
)

// RolloutControllersCmd is a CLI command that initiates rollout and restart on capa-controller-manager deployment.
func RolloutControllersCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "rollout-controller",
		Short: "initiates rollout and restart on capa-controller-manager deployment",
		Long: templates.LongDesc(`
			initiates rollout and restart on capa-controller-manager deployment
		`),
		Example: templates.Examples(`
		# rollout controller deployment
		clusterawsadm controller rollout-controller --kubeconfig=kubeconfig --namespace=capa-system
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rollout.RolloutControllers(rollout.RolloutControllersInput{
				KubeconfigPath:    kubeconfigPath,
				KubeconfigContext: kubeconfigContext,
				Namespace:         namespace,
			})
		},
	}
	addKubeconfigFlag(newCmd)
	addKubeconfigContextFlag(newCmd)
	addNamespaceFlag(newCmd)
	return newCmd
}
