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

// Package controller provides the controller command.
package controller

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/controller/credentials"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/controller/rollout"
)

// RootCmd is the root of the `controller command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "controller [command]",
		Short: "controller commands",
		Args:  cobra.NoArgs,
		Long: templates.LongDesc(`
			All controller related actions such as:
			# Zero controller credentials and rollout controllers
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	newCmd.AddCommand(credentials.ZeroCredentialsCmd())
	newCmd.AddCommand(credentials.UpdateCredentialsCmd())
	newCmd.AddCommand(credentials.PrintCredentialsCmd())
	newCmd.AddCommand(rollout.RolloutControllersCmd())

	return newCmd
}
