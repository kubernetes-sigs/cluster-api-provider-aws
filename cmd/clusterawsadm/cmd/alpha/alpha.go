/*
Copyright 2018 The Kubernetes Authors.

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

package alpha

import (
	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/bootstrap"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/migrate"
)

// AlphaCmd is the top-level alpha set of commands
func AlphaCmd() *cobra.Command { // nolint
	newCmd := &cobra.Command{
		Use:   "alpha",
		Short: "alpha commands",
		Long:  `Alpha commands may not be supported in future releases`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	newCmd.AddCommand(bootstrap.RootCmd())
	newCmd.AddCommand(migrate.MigrateCmd())
	return newCmd
}
