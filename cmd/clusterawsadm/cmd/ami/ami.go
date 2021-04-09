/*
Copyright 2020 The Kubernetes Authors.

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

package ami

import (
	"github.com/spf13/cobra"
	cp "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/ami/copy"
	ls "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd/ami/list"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/cmd"
)

// RootCmd is the root of the `ami command`
func RootCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "ami [command]",
		Short: "AMI commands",
		Args:  cobra.NoArgs,
		Long: cmd.LongDesc(`
			All AMI related actions such as:
			# Copy AMIs based on Kubernetes version, OS etc from an AWS account where AMIs are stored
            to the current AWS account (use case: air-gapped deployments)
			# (to be implemented) List available AMIs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	newCmd.AddCommand(cp.CopyAMICmd())
	newCmd.AddCommand(cp.EncryptedCopyAMICmd())
	newCmd.AddCommand(ls.ListAMICmd())
	return newCmd
}
