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
)

func addKubeconfigContextFlag(c *cobra.Command) {
	c.Flags().StringVar(&kubeconfigContext, "kubeconfig-context", "", "Context to be used within the kubeconfig file. If empty, current context will be used.")
}

func addKubeconfigFlag(c *cobra.Command) {
	c.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file to use for the management cluster. If empty, default discovery rules apply.")
}

func addNamespaceFlag(c *cobra.Command) {
	c.Flags().StringVar(&namespace, "namespace", "capa-system", "Namespace the controllers are in. If empty, default value (capa-system) is used")
}
