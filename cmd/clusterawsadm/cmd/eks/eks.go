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

package eks

import (
	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/eks/addons"
)

// RootCmd is an EKS root CLI command.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "eks",
		Short: "Commands related to EKS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	newCmd.AddCommand(addons.RootCmd())

	return newCmd
}
