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

package resource

import (
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/resource/list"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

// RootCmd is the root of the `resource command`.
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "resource [command]",
		Short: "Commands related to AWS resources",
		Args:  cobra.NoArgs,
		Long: cmd.LongDesc(`
			All AWS resources related actions such as:
			# List of AWS resources created by CAPA
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	newCmd.AddCommand(list.ListAWSResourceCmd())

	return newCmd
}
