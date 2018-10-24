// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bootstrap

import (
	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/bootstrap/cloudformation"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/alpha/bootstrap/credentials"
)

// Cmd is the root of the `alpha bootstrap command`
func Cmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap",
		Long:  `Bootstrapping commands`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	newCmd.AddCommand(cloudformation.Cmd())
	newCmd.AddCommand(credentials.Cmd())
	return newCmd
}
