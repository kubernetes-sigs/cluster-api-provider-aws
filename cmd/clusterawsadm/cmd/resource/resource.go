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

// Package resource provides commands related to AWS resources.
package resource

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/cmd/resource/list"
)

// RootCmd is the root of the `resource command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "resource [command]",
		Short: "Commands related to AWS resources",
		Args:  cobra.NoArgs,
		Long: templates.LongDesc(`
			All AWS resources related actions such as:
			# List of AWS resources created by CAPA
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	newCmd.AddCommand(list.ListAWSResourceCmd())

	return newCmd
}
